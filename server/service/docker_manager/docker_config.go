package docker_manager

import (
	"ark-server-manager/utils"
	"fmt"
	"io"
	"strings"

	"github.com/docker/docker/api/types/container"
)

// ReadConfigFile 从容器卷中读取配置文件
// serverID: 服务器ID
// fileName: 文件名
// 返回: 文件内容和错误信息
func (dm *DockerManager) ReadConfigFile(serverID uint, fileName string) (string, error) {
	volumeName := utils.GetServerVolumeName(serverID)
	alpineImage := "alpine:latest"

	// 检查Alpine镜像是否存在
	exists, err := dm.ImageExists(alpineImage)
	if err != nil {
		return "", fmt.Errorf("检查Alpine镜像失败: %v", err)
	}

	if !exists {
		return "", fmt.Errorf("Alpine镜像不存在，请确保后端启动时已成功拉取镜像")
	}

	// 配置文件路径（所有配置文件都在 Config/WindowsServer 目录）
	configPath := fmt.Sprintf("/home/steam/arkserver/ShooterGame/Saved/Config/WindowsServer/%s", fileName)

	// 创建临时容器读取文件
	containerConfig := &container.Config{
		Image: alpineImage,
		Cmd:   []string{"cat", configPath},
	}

	hostConfig := &container.HostConfig{
		Binds: []string{
			fmt.Sprintf("%s:/home/steam/arkserver/ShooterGame/Saved", volumeName),
		},
	}

	// 创建临时容器
	resp, err := dm.client.ContainerCreate(dm.ctx, containerConfig, hostConfig, nil, nil, "")
	if err != nil {
		return "", fmt.Errorf("创建临时容器失败: %v", err)
	}

	// 启动容器
	err = dm.client.ContainerStart(dm.ctx, resp.ID, container.StartOptions{})
	if err != nil {
		return "", fmt.Errorf("启动临时容器失败: %v", err)
	}

	// 等待容器完成
	waitCh, errCh := dm.client.ContainerWait(dm.ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			return "", fmt.Errorf("等待容器完成失败: %v", err)
		}
	case <-waitCh:
		// 容器已完成
	}

	// 获取容器日志（输出）
	logs, err := dm.client.ContainerLogs(dm.ctx, resp.ID, container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
	})
	if err != nil {
		return "", fmt.Errorf("获取容器日志失败: %v", err)
	}
	defer logs.Close()

	output, err := io.ReadAll(logs)
	if err != nil {
		return "", fmt.Errorf("读取容器输出失败: %v", err)
	}

	// 删除临时容器
	dm.client.ContainerRemove(dm.ctx, resp.ID, container.RemoveOptions{
		Force: true,
	})

	// 移除Docker日志前缀（前8个字节）
	if len(output) > 8 {
		output = output[8:]
	}

	return string(output), nil
}

// WriteConfigFile 向容器卷中写入配置文件
// serverID: 服务器ID
// fileName: 文件名
// content: 文件内容
// 返回: 错误信息
func (dm *DockerManager) WriteConfigFile(serverID uint, fileName, content string) error {
	volumeName := utils.GetServerVolumeName(serverID)
	alpineImage := "alpine:latest"

	// 检查Alpine镜像是否存在
	exists, err := dm.ImageExists(alpineImage)
	if err != nil {
		return fmt.Errorf("检查Alpine镜像失败: %v", err)
	}

	if !exists {
		return fmt.Errorf("Alpine镜像不存在，请确保后端启动时已成功拉取镜像")
	}

	// 配置文件路径（所有配置文件都在 Config/WindowsServer 目录）
	configPath := fmt.Sprintf("/home/steam/arkserver/ShooterGame/Saved/Config/WindowsServer/%s", fileName)
	configDir := "/home/steam/arkserver/ShooterGame/Saved/Config/WindowsServer"

	// 首先确保目录存在
	containerConfig := &container.Config{
		Image: alpineImage,
		Cmd:   []string{"mkdir", "-p", configDir},
	}

	hostConfig := &container.HostConfig{
		Binds: []string{
			fmt.Sprintf("%s:/home/steam/arkserver/ShooterGame/Saved", volumeName),
		},
	}

	// 创建临时容器创建目录
	resp, err := dm.client.ContainerCreate(dm.ctx, containerConfig, hostConfig, nil, nil, "")
	if err != nil {
		return fmt.Errorf("创建临时容器失败: %v", err)
	}

	// 启动容器
	err = dm.client.ContainerStart(dm.ctx, resp.ID, container.StartOptions{})
	if err != nil {
		return fmt.Errorf("启动临时容器失败: %v", err)
	}

	// 等待容器完成
	waitCh, errCh := dm.client.ContainerWait(dm.ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			return fmt.Errorf("等待容器完成失败: %v", err)
		}
	case <-waitCh:
		// 容器已完成
	}

	// 删除临时容器
	dm.client.ContainerRemove(dm.ctx, resp.ID, container.RemoveOptions{
		Force: true,
	})

	// 写入文件内容
	containerConfig = &container.Config{
		Image: alpineImage,
		Cmd:   []string{"sh", "-c", fmt.Sprintf("echo '%s' > %s", strings.ReplaceAll(content, "'", "'\"'\"'"), configPath)},
	}

	// 创建临时容器写入文件
	resp, err = dm.client.ContainerCreate(dm.ctx, containerConfig, hostConfig, nil, nil, "")
	if err != nil {
		return fmt.Errorf("创建临时容器失败: %v", err)
	}

	// 启动容器
	err = dm.client.ContainerStart(dm.ctx, resp.ID, container.StartOptions{})
	if err != nil {
		return fmt.Errorf("启动临时容器失败: %v", err)
	}

	// 等待容器完成
	waitCh, errCh = dm.client.ContainerWait(dm.ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			return fmt.Errorf("等待容器完成失败: %v", err)
		}
	case <-waitCh:
		// 容器已完成
	}

	// 删除临时容器
	dm.client.ContainerRemove(dm.ctx, resp.ID, container.RemoveOptions{
		Force: true,
	})

	fmt.Printf("配置文件写入成功: %s\n", fileName)
	return nil
}
