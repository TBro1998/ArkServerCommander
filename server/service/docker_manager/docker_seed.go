package docker_manager

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

// 服务端主程序相对安装目录的路径，存在才说明游戏本体已下载。Windows 版（Proton 镜像）和 Linux 版的主程序不同，
// 同一个文件卷换过镜像后两个都可能存在
const (
	windowsServerExecutable = "ShooterGame/Binaries/Win64/ShooterGameServer.exe"
	linuxServerExecutable   = "ShooterGame/Binaries/Linux/ShooterGameServer"
)

// installMarkerFile 入口脚本下载完成后写入的安装标记，位于服务端文件卷根目录
const installMarkerFile = ".install_complete"

// ServerFilesState 服务端文件卷的安装状态
type ServerFilesState struct {
	Installed bool // 有 .install_complete 标记（入口脚本下载完成后写入）
	Windows   bool // Windows 版服务端主程序存在
	Linux     bool // Linux 版服务端主程序存在
}

// HasExecutable 任一平台的服务端主程序存在，说明游戏本体已经下载过
func (s ServerFilesState) HasExecutable() bool {
	return s.Windows || s.Linux
}

// HasExecutableFor 指定平台的服务端主程序存在（arkApi 为 true 表示 Windows 版）
func (s ServerFilesState) HasExecutableFor(arkApi bool) bool {
	if arkApi {
		return s.Windows
	}
	return s.Linux
}

// isRegularFile 路径存在且是普通文件（跟随符号链接）
func isRegularFile(p string) bool {
	info, err := os.Stat(p)
	return err == nil && info.Mode().IsRegular()
}

// CheckServerFiles 检查多台服务器的服务端文件目录，返回各自的安装状态。
// 目录不存在的服务器不会出现在结果里
func (dm *DockerManager) CheckServerFiles(serverIDs []uint) (map[uint]ServerFilesState, error) {
	result := map[uint]ServerFilesState{}
	for _, id := range serverIDs {
		dir := ServerFilesDir(id)
		info, err := os.Stat(dir)
		if errors.Is(err, fs.ErrNotExist) {
			continue
		}
		if err != nil {
			return nil, fmt.Errorf("检查服务端文件失败: %w", err)
		}
		if !info.IsDir() {
			continue
		}
		result[id] = ServerFilesState{
			Installed: isRegularFile(filepath.Join(dir, installMarkerFile)),
			Windows:   isRegularFile(filepath.Join(dir, filepath.FromSlash(windowsServerExecutable))),
			Linux:     isRegularFile(filepath.Join(dir, filepath.FromSlash(linuxServerExecutable))),
		}
	}
	return result, nil
}

// CopyServerFiles 把 srcID 的服务端文件（含已安装的 Mod 和 .installed_mods）复制到 dstID 的文件卷，
// 并删除 .install_complete，让入口脚本启动时执行 steamcmd validate：只校验本地文件、下载有差异的部分
// （来源是另一个平台时也一样，validate 会补上这个平台的主程序）。ArkApi 文件在文件卷里，会一起复制。
// 存档和插件在单独的卷里，不会被复制
func (dm *DockerManager) CopyServerFiles(srcID, dstID uint) error {
	if srcID == dstID {
		return errors.New("来源和目标是同一台服务器")
	}
	if err := copyTree(ServerFilesDir(srcID), ServerFilesDir(dstID)); err != nil {
		return fmt.Errorf("复制服务端文件失败: %w", err)
	}
	if err := os.Remove(filepath.Join(ServerFilesDir(dstID), installMarkerFile)); err != nil && !errors.Is(err, fs.ErrNotExist) {
		return fmt.Errorf("复制服务端文件失败: %w", err)
	}
	return nil
}
