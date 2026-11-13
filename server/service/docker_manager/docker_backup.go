package docker_manager

import (
	"archive/tar"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"path"
	"strings"
	"time"

	"ark-server-commander/utils"

	"github.com/docker/docker/api/types/container"
)

const (
	// shooterGameDir 数据卷挂载点 Saved 的上一级目录，导入存档时以它为解压目标
	shooterGameDir = "/home/steam/arkserver/ShooterGame"
	// backupRoot 备份归档里的顶层目录
	backupRoot = "Saved/"
)

// ErrInvalidBackup 备份文件不是本系统导出的格式
var ErrInvalidBackup = errors.New("备份文件格式无效")

// backupSkipped 备份时跳过的目录：日志没有恢复价值，集群目录是单独的共享卷
func backupSkipped(name string) bool {
	return strings.HasPrefix(name, backupRoot+"Logs/") || strings.HasPrefix(name, backupRoot+"clusters/")
}

// ExportSaved 把服务器数据卷（Saved 目录，含存档和配置文件）打包成 tar.gz 写入 w
func (dm *DockerManager) ExportSaved(serverID uint, w io.Writer) error {
	containerID, cleanup, err := dm.newConfigHelper(serverID)
	if err != nil {
		return err
	}
	defer cleanup()

	reader, _, err := dm.client.CopyFromContainer(dm.ctx, containerID, configSavedDir)
	if err != nil {
		return fmt.Errorf("读取数据卷失败: %v", err)
	}
	defer reader.Close()

	gz := gzip.NewWriter(w)
	tw := tar.NewWriter(gz)
	tr := tar.NewReader(reader)
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("读取数据卷失败: %v", err)
		}
		if backupSkipped(header.Name) {
			continue
		}
		// 只保留普通文件和目录（与 ValidateBackup 一致）：Linux 版服务端的 Config/LinuxServer 是入口脚本建的符号链接，
		// 打进归档会让这个备份无法恢复，启动时脚本会重新创建它
		if header.Typeflag != tar.TypeReg && header.Typeflag != tar.TypeDir {
			continue
		}
		if err := tw.WriteHeader(header); err != nil {
			return fmt.Errorf("写入备份失败: %v", err)
		}
		if header.Typeflag == tar.TypeReg {
			if _, err := io.Copy(tw, tr); err != nil {
				return fmt.Errorf("写入备份失败: %v", err)
			}
		}
	}
	if err := tw.Close(); err != nil {
		return fmt.Errorf("写入备份失败: %v", err)
	}
	if err := gz.Close(); err != nil {
		return fmt.Errorf("写入备份失败: %v", err)
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

// ImportSaved 用备份覆盖服务器的存档和配置：先清空 SavedArks 和 Config，再解压归档。
// 调用方需保证容器已停止，r 是 ValidateBackup 校验过的 tar.gz
func (dm *DockerManager) ImportSaved(serverID uint, r io.Reader) error {
	volumeBind := fmt.Sprintf("%s:/saved", utils.GetServerVolumeName(serverID))
	if _, err := dm.RunHelper([]string{volumeBind}, []string{"rm", "-rf", "/saved/SavedArks", "/saved/Config"}, 2*time.Minute); err != nil {
		return fmt.Errorf("清理旧存档失败: %w", err)
	}

	containerID, cleanup, err := dm.newConfigHelper(serverID)
	if err != nil {
		return err
	}
	defer cleanup()

	// Docker 接受 gzip 压缩的 tar，归档里的 Saved/ 正好落在数据卷的挂载点上
	if err := dm.client.CopyToContainer(dm.ctx, containerID, shooterGameDir, r, container.CopyToContainerOptions{}); err != nil {
		return fmt.Errorf("写入存档失败: %v", err)
	}
	return nil
}

// RemoveInstallMarkers 删除服务端文件卷里的安装标记，下次启动时入口脚本会重新执行
// steamcmd app_update validate（game）或重新下载全部 Mod（mods）
func (dm *DockerManager) RemoveInstallMarkers(serverID uint, game, mods bool) error {
	cmd := []string{"rm", "-f"}
	if game {
		cmd = append(cmd, "/files/.install_complete")
	}
	if mods {
		cmd = append(cmd, "/files/.installed_mods")
	}
	if len(cmd) == 2 {
		return nil
	}
	bind := fmt.Sprintf("%s:/files", utils.GetServerFilesVolumeName(serverID))
	if _, err := dm.RunHelper([]string{bind}, cmd, time.Minute); err != nil {
		return fmt.Errorf("清除安装标记失败: %w", err)
	}
	return nil
}
