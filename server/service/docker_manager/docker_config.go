package docker_manager

import (
	"ark-server-commander/utils"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"go.uber.org/zap"
)

const (
	// 配置文件相对 Saved 目录的位置（所有配置文件都在 Config/WindowsServer 目录；
	// Linux 版服务端读的 Config/LinuxServer 由入口脚本建成指向它的符号链接）
	configSubDir = "Config/WindowsServer"
	// 配置文件大小上限，防止异常文件占满内存
	maxConfigFileSize = 10 << 20
	// UTF-8 字节序标记，Windows 编辑器保存的配置文件可能带有
	utf8BOM = "\uFEFF"
)

// ErrConfigNotFound 配置文件在数据卷里不存在（新建的服务器还没有写过配置）
var ErrConfigNotFound = errors.New("文件不存在")

// openSavedRoot 打开服务器数据卷（Saved 目录）在后端本地的目录。
// 目录里的内容由 ARK 容器以 root 写入，用 os.Root 访问，符号链接不能指到目录外
func openSavedRoot(serverID uint) (*os.Root, error) {
	return os.OpenRoot(ServerSavedDir(serverID))
}

// checkConfigFileName 配置文件名只能是单个文件名，不能带路径
func checkConfigFileName(fileName string) error {
	if fileName == "" || fileName != filepath.Base(fileName) || fileName == "." || fileName == ".." ||
		strings.ContainsAny(fileName, `/\`) {
		return fmt.Errorf("无效的配置文件名: %q", fileName)
	}
	return nil
}

// configFilePath 配置文件相对 Saved 目录的路径
func configFilePath(fileName string) string {
	return filepath.Join(filepath.FromSlash(configSubDir), fileName)
}

// ReadConfigFile 从数据卷中读取配置文件（直接读后端本地的数据目录）
// serverID: 服务器ID
// fileName: 文件名
// 返回: 文件内容和错误信息
func (dm *DockerManager) ReadConfigFile(serverID uint, fileName string) (string, error) {
	if err := checkConfigFileName(fileName); err != nil {
		return "", err
	}

	root, err := openSavedRoot(serverID)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return "", fmt.Errorf("%w: %s", ErrConfigNotFound, fileName)
		}
		return "", fmt.Errorf("打开数据目录失败: %v", err)
	}
	defer root.Close()

	f, err := root.Open(configFilePath(fileName))
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return "", fmt.Errorf("%w: %s", ErrConfigNotFound, fileName)
		}
		return "", fmt.Errorf("读取配置文件失败: %v", err)
	}
	defer f.Close()

	info, err := f.Stat()
	if err != nil {
		return "", fmt.Errorf("读取配置文件失败: %v", err)
	}
	if !info.Mode().IsRegular() {
		return "", fmt.Errorf("路径不是文件: %s", fileName)
	}
	if info.Size() > maxConfigFileSize {
		return "", fmt.Errorf("配置文件过大: %s (%d 字节)", fileName, info.Size())
	}

	data, err := io.ReadAll(io.LimitReader(f, maxConfigFileSize))
	if err != nil {
		return "", fmt.Errorf("读取文件内容失败: %v", err)
	}

	// 去掉 UTF-8 BOM 和首尾空白，写入时会补上结尾换行
	return strings.TrimSpace(strings.TrimPrefix(string(data), utf8BOM)), nil
}

// WriteConfigFile 向数据卷中写入配置文件（覆盖已有内容）。
// 先写到同目录的临时文件再改名，写到一半失败不会留下残缺的配置
// serverID: 服务器ID
// fileName: 文件名
// content: 文件内容
// 返回: 错误信息
func (dm *DockerManager) WriteConfigFile(serverID uint, fileName, content string) error {
	if err := checkConfigFileName(fileName); err != nil {
		return err
	}
	if len(content) > maxConfigFileSize {
		return fmt.Errorf("配置文件过大: %s", fileName)
	}
	if !strings.HasSuffix(content, "\n") {
		content += "\n"
	}

	root, err := openSavedRoot(serverID)
	if err != nil {
		return fmt.Errorf("打开数据目录失败: %v", err)
	}
	defer root.Close()

	if err := root.MkdirAll(filepath.FromSlash(configSubDir), 0o755); err != nil {
		return fmt.Errorf("创建配置目录失败: %v", err)
	}
	target := configFilePath(fileName)
	tmp := target + ".asc-tmp"
	if err := root.WriteFile(tmp, []byte(content), 0o644); err != nil {
		_ = root.Remove(tmp)
		return fmt.Errorf("写入配置文件失败: %v", err)
	}
	if err := root.Rename(tmp, target); err != nil {
		_ = root.Remove(tmp)
		return fmt.Errorf("写入配置文件失败: %v", err)
	}

	utils.Info("配置文件写入成功", zap.Uint("server_id", serverID), zap.String("file", fileName))
	return nil
}
