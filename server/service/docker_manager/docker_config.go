package docker_manager

import (
	"archive/tar"
	"ark-server-commander/utils"
	"bytes"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"time"

	"github.com/containerd/errdefs"
	"github.com/docker/docker/api/types/container"
	"go.uber.org/zap"
)

const (
	// 数据卷挂载到临时容器里的路径，与 ARK 容器内的路径一致
	configSavedDir = "/home/steam/arkserver/ShooterGame/Saved"
	// 配置文件相对 Saved 目录的位置（所有配置文件都在 Config/WindowsServer 目录）
	configSubDir = "Config/WindowsServer"
	// 用于访问数据卷的辅助镜像
	configHelperImage = "alpine:latest"
	// 配置文件大小上限，防止异常文件占满内存
	maxConfigFileSize = 10 << 20
	// UTF-8 字节序标记，Windows 编辑器保存的配置文件可能带有
	utf8BOM = "\uFEFF"
)

// newConfigHelper 创建一个挂载了服务器数据卷的临时容器。
// 容器只创建、不启动：CopyFromContainer / CopyToContainer 对已创建的容器同样可用，
// 省去启动和等待退出的开销，也不会留下运行中的临时容器。
// 返回容器 ID 和清理函数，调用方必须调用清理函数。
func (dm *DockerManager) newConfigHelper(serverID uint) (string, func(), error) {
	exists, err := dm.ImageExists(configHelperImage)
	if err != nil {
		return "", nil, fmt.Errorf("检查Alpine镜像失败: %v", err)
	}
	if !exists {
		return "", nil, fmt.Errorf("Alpine镜像不存在，请先在首页下载镜像")
	}

	resp, err := dm.client.ContainerCreate(dm.ctx,
		&container.Config{Image: configHelperImage, Cmd: []string{"true"}},
		&container.HostConfig{Binds: []string{
			fmt.Sprintf("%s:%s", utils.GetServerVolumeName(serverID), configSavedDir),
		}},
		nil, nil, "")
	if err != nil {
		return "", nil, fmt.Errorf("创建临时容器失败: %v", err)
	}

	cleanup := func() {
		if err := dm.client.ContainerRemove(dm.ctx, resp.ID, container.RemoveOptions{Force: true}); err != nil {
			utils.Warn("清理临时容器失败", zap.String("container_id", resp.ID), zap.Error(err))
		}
	}
	return resp.ID, cleanup, nil
}

// checkConfigFileName 配置文件名只能是单个文件名，不能带路径
func checkConfigFileName(fileName string) error {
	if fileName == "" || fileName != filepath.Base(fileName) || fileName == "." || fileName == ".." {
		return fmt.Errorf("无效的配置文件名: %q", fileName)
	}
	return nil
}

// ReadConfigFile 从数据卷中读取配置文件
// serverID: 服务器ID
// fileName: 文件名
// 返回: 文件内容和错误信息
func (dm *DockerManager) ReadConfigFile(serverID uint, fileName string) (string, error) {
	if err := checkConfigFileName(fileName); err != nil {
		return "", err
	}

	containerID, cleanup, err := dm.newConfigHelper(serverID)
	if err != nil {
		return "", err
	}
	defer cleanup()

	configPath := fmt.Sprintf("%s/%s/%s", configSavedDir, configSubDir, fileName)
	reader, _, err := dm.client.CopyFromContainer(dm.ctx, containerID, configPath)
	if err != nil {
		if errdefs.IsNotFound(err) {
			return "", fmt.Errorf("文件不存在: %s", fileName)
		}
		return "", fmt.Errorf("从数据卷复制文件失败: %v", err)
	}
	defer reader.Close()

	// 返回的是 tar 流，里面只有这一个文件
	tarReader := tar.NewReader(reader)
	header, err := tarReader.Next()
	if err != nil {
		if err == io.EOF {
			return "", fmt.Errorf("文件不存在: %s", fileName)
		}
		return "", fmt.Errorf("读取 tar 文件头失败: %v", err)
	}
	if header.Typeflag != tar.TypeReg {
		return "", fmt.Errorf("路径不是文件: %s", fileName)
	}
	if header.Size > maxConfigFileSize {
		return "", fmt.Errorf("配置文件过大: %s (%d 字节)", fileName, header.Size)
	}

	var buffer bytes.Buffer
	if _, err := io.Copy(&buffer, io.LimitReader(tarReader, maxConfigFileSize)); err != nil {
		return "", fmt.Errorf("读取文件内容失败: %v", err)
	}

	// 去掉 UTF-8 BOM 和首尾空白，写入时会补上结尾换行
	return strings.TrimSpace(strings.TrimPrefix(buffer.String(), utf8BOM)), nil
}

// WriteConfigFile 向数据卷中写入配置文件（覆盖已有内容）
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

	containerID, cleanup, err := dm.newConfigHelper(serverID)
	if err != nil {
		return err
	}
	defer cleanup()

	if !strings.HasSuffix(content, "\n") {
		content += "\n"
	}

	// 用 tar 把目录和文件一起写入：Config/WindowsServer 不存在时会被创建，已存在时只合并内容。
	// 内容直接放在 tar 里，不经过 shell，不需要转义，也没有命令行长度限制
	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)
	now := time.Now()
	for _, dir := range []string{"Config/", configSubDir + "/"} {
		if err := tw.WriteHeader(&tar.Header{Typeflag: tar.TypeDir, Name: dir, Mode: 0o755, ModTime: now}); err != nil {
			return fmt.Errorf("生成配置文件归档失败: %v", err)
		}
	}
	if err := tw.WriteHeader(&tar.Header{
		Typeflag: tar.TypeReg,
		Name:     configSubDir + "/" + fileName,
		Mode:     0o644,
		Size:     int64(len(content)),
		ModTime:  now,
	}); err != nil {
		return fmt.Errorf("生成配置文件归档失败: %v", err)
	}
	if _, err := tw.Write([]byte(content)); err != nil {
		return fmt.Errorf("生成配置文件归档失败: %v", err)
	}
	if err := tw.Close(); err != nil {
		return fmt.Errorf("生成配置文件归档失败: %v", err)
	}

	if err := dm.client.CopyToContainer(dm.ctx, containerID, configSavedDir, &buf, container.CopyToContainerOptions{}); err != nil {
		return fmt.Errorf("写入配置文件失败: %v", err)
	}

	utils.Info("配置文件写入成功", zap.Uint("server_id", serverID), zap.String("file", fileName))
	return nil
}
