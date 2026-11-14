package docker_manager

import (
	"archive/tar"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
)

const (
	// backupRoot 备份归档里的顶层目录
	backupRoot = "Saved/"
)

// ErrInvalidBackup 备份文件不是本系统导出的格式
var ErrInvalidBackup = errors.New("备份文件格式无效")

// backupSkipped 备份时跳过的目录：日志没有恢复价值，集群目录是单独的共享卷
func backupSkipped(name string) bool {
	return strings.HasPrefix(name, backupRoot+"Logs/") || strings.HasPrefix(name, backupRoot+"clusters/")
}

// ExportSaved 把服务器数据卷（Saved 目录，含存档和配置文件）打包成 tar.gz 写入 w，直接读后端本地的数据目录
func (dm *DockerManager) ExportSaved(serverID uint, w io.Writer) error {
	root, err := openSavedRoot(serverID)
	if err != nil {
		return fmt.Errorf("读取数据卷失败: %v", err)
	}
	defer root.Close()

	gz := gzip.NewWriter(w)
	tw := tar.NewWriter(gz)
	err = fs.WalkDir(root.FS(), ".", func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		name := backupRoot
		if p != "." {
			name += p
		}
		if d.IsDir() {
			if p != "." {
				name += "/"
			}
			if backupSkipped(name) {
				return fs.SkipDir
			}
		} else if backupSkipped(name) {
			return nil
		}
		// 只保留普通文件和目录（与 ValidateBackup 一致）：Linux 版服务端的 Config/LinuxServer 是入口脚本建的符号链接，
		// 打进归档会让这个备份无法恢复，启动时脚本会重新创建它
		if !d.IsDir() && !d.Type().IsRegular() {
			return nil
		}
		if d.IsDir() {
			info, err := d.Info()
			if err != nil {
				return err
			}
			return tw.WriteHeader(&tar.Header{Typeflag: tar.TypeDir, Name: name, Mode: int64(info.Mode().Perm()), ModTime: info.ModTime()})
		}
		return addFileToTar(tw, root, filepath.FromSlash(p), name)
	})
	if err != nil {
		return fmt.Errorf("写入备份失败: %v", err)
	}
	if err := tw.Close(); err != nil {
		return fmt.Errorf("写入备份失败: %v", err)
	}
	if err := gz.Close(); err != nil {
		return fmt.Errorf("写入备份失败: %v", err)
	}
	return nil
}

// addFileToTar 把一个普通文件写入归档。大小以打开后的 Stat 为准，运行中备份时文件变大只取这么多字节，
// 变小则报错（写出的条目会不完整）
func addFileToTar(tw *tar.Writer, root *os.Root, rel, name string) error {
	f, err := root.Open(rel)
	if err != nil {
		return err
	}
	defer f.Close()
	info, err := f.Stat()
	if err != nil {
		return err
	}
	if !info.Mode().IsRegular() {
		return nil
	}
	if err := tw.WriteHeader(&tar.Header{
		Typeflag: tar.TypeReg,
		Name:     name,
		Mode:     int64(info.Mode().Perm()),
		Size:     info.Size(),
		ModTime:  info.ModTime(),
	}); err != nil {
		return err
	}
	if _, err := io.CopyN(tw, f, info.Size()); err != nil {
		return fmt.Errorf("读取 %s 失败（文件在备份过程中被修改）: %v", name, err)
	}
	return nil
}

// ValidateBackup 检查归档是否为本系统导出的格式：tar.gz，所有条目都在 Saved/ 下，且没有越界路径。
// 返回归档中是否包含存档目录 SavedArks
func ValidateBackup(r io.Reader) (hasSavedArks bool, err error) {
	gz, err := gzip.NewReader(r)
	if err != nil {
		return false, fmt.Errorf("%w: 不是 gzip 文件", ErrInvalidBackup)
	}
	defer gz.Close()

	tr := tar.NewReader(gz)
	count := 0
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return false, fmt.Errorf("%w: %v", ErrInvalidBackup, err)
		}
		name := header.Name
		clean := path.Clean(name)
		if path.IsAbs(name) || clean == ".." || strings.HasPrefix(clean, "../") ||
			(clean != "Saved" && !strings.HasPrefix(clean, backupRoot)) {
			return false, fmt.Errorf("%w: 包含不允许的路径 %q", ErrInvalidBackup, name)
		}
		switch header.Typeflag {
		case tar.TypeReg, tar.TypeDir:
		default:
			return false, fmt.Errorf("%w: 包含不支持的文件类型 %q", ErrInvalidBackup, name)
		}
		if strings.HasPrefix(clean, backupRoot+"SavedArks") {
			hasSavedArks = true
		}
		count++
	}
	if count == 0 {
		return false, fmt.Errorf("%w: 归档为空", ErrInvalidBackup)
	}
	return hasSavedArks, nil
}

// ImportSaved 用备份覆盖服务器的存档和配置：先清空 SavedArks 和 Config，再解压归档到数据目录。
// 调用方需保证容器已停止，r 是 ValidateBackup 校验过的 tar.gz（这里仍会逐项检查路径和类型）
func (dm *DockerManager) ImportSaved(serverID uint, r io.Reader) error {
	root, err := openSavedRoot(serverID)
	if err != nil {
		return fmt.Errorf("打开数据目录失败: %v", err)
	}
	defer root.Close()

	for _, dir := range []string{"SavedArks", "Config"} {
		if err := root.RemoveAll(dir); err != nil {
			return fmt.Errorf("清理旧存档失败: %w", err)
		}
	}

	gz, err := gzip.NewReader(r)
	if err != nil {
		return fmt.Errorf("%w: 不是 gzip 文件", ErrInvalidBackup)
	}
	defer gz.Close()

	tr := tar.NewReader(gz)
	for {
		header, err := tr.Next()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return fmt.Errorf("读取备份失败: %v", err)
		}
		clean := path.Clean(header.Name)
		if clean == "Saved" {
			continue
		}
		rel, ok := strings.CutPrefix(clean, backupRoot)
		if !ok || path.IsAbs(clean) {
			return fmt.Errorf("%w: 包含不允许的路径 %q", ErrInvalidBackup, header.Name)
		}
		local := filepath.FromSlash(rel)
		switch header.Typeflag {
		case tar.TypeDir:
			if err := root.MkdirAll(local, 0o755); err != nil {
				return fmt.Errorf("写入存档失败: %v", err)
			}
		case tar.TypeReg:
			if err := extractFile(root, local, header, tr); err != nil {
				return fmt.Errorf("写入存档失败: %v", err)
			}
		default:
			return fmt.Errorf("%w: 包含不支持的文件类型 %q", ErrInvalidBackup, header.Name)
		}
	}
}

// extractFile 把归档里的一个普通文件写到 root 下的 rel，保留权限和修改时间
func extractFile(root *os.Root, rel string, header *tar.Header, r io.Reader) error {
	if dir := filepath.Dir(rel); dir != "." {
		if err := root.MkdirAll(dir, 0o755); err != nil {
			return err
		}
	}
	mode := fs.FileMode(header.Mode).Perm()
	if mode == 0 {
		mode = 0o644
	}
	f, err := root.OpenFile(rel, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, mode)
	if err != nil {
		return err
	}
	if _, err := io.Copy(f, r); err != nil {
		f.Close()
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	if !header.ModTime.IsZero() {
		_ = root.Chtimes(rel, header.ModTime, header.ModTime)
	}
	return nil
}

// RemoveInstallMarkers 删除服务端文件卷里的安装标记，下次启动时入口脚本会重新执行
// steamcmd app_update validate（game）或重新下载全部 Mod（mods）
func (dm *DockerManager) RemoveInstallMarkers(serverID uint, game, mods bool) error {
	var names []string
	if game {
		names = append(names, installMarkerFile)
	}
	if mods {
		names = append(names, installedModsFile)
	}
	dir := ServerFilesDir(serverID)
	for _, name := range names {
		if err := os.Remove(filepath.Join(dir, name)); err != nil && !errors.Is(err, fs.ErrNotExist) {
			return fmt.Errorf("清除安装标记失败: %w", err)
		}
	}
	return nil
}
