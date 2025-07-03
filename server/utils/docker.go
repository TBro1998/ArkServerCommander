package utils

import (
	"context"
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

// DockerManager Docker管理器结构体
type DockerManager struct {
	client *client.Client
	ctx    context.Context
}

// NewDockerManager 创建新的Docker管理器
func NewDockerManager() (*DockerManager, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, fmt.Errorf("创建Docker客户端失败: %v", err)
	}

	return &DockerManager{
		client: cli,
		ctx:    context.Background(),
	}, nil
}

// Close 关闭Docker客户端连接
func (dm *DockerManager) Close() error {
	return dm.client.Close()
}

// CreateVolume 创建Docker卷
// serverID: 服务器ID
// 返回: 卷名称和错误信息
func (dm *DockerManager) CreateVolume(serverID uint) (string, error) {
	volumeName := GetServerVolumeName(serverID)

	// 检查卷是否已存在
	exists, err := dm.VolumeExists(volumeName)
	if err != nil {
		return "", fmt.Errorf("检查卷是否存在失败: %v", err)
	}
	if exists {
		fmt.Printf("Docker卷 %s 已存在，跳过创建\n", volumeName)
		return volumeName, nil
	}

	// 创建卷
	fmt.Printf("正在创建Docker卷: %s\n", volumeName)
	volumeCreateBody := volume.CreateOptions{
		Name: volumeName,
	}

	_, err = dm.client.VolumeCreate(dm.ctx, volumeCreateBody)
	if err != nil {
		return "", fmt.Errorf("创建Docker卷失败: %v", err)
	}

	fmt.Printf("Docker卷创建成功: %s\n", volumeName)
	return volumeName, nil
}

// RemoveVolume 删除Docker卷
// volumeName: 卷名称
// 返回: 错误信息
func (dm *DockerManager) RemoveVolume(volumeName string) error {
	// 检查卷是否存在
	exists, err := dm.VolumeExists(volumeName)
	if err != nil {
		return fmt.Errorf("检查卷是否存在失败: %v", err)
	}
	if !exists {
		fmt.Printf("Docker卷 %s 不存在，跳过删除\n", volumeName)
		return nil
	}

	fmt.Printf("正在删除Docker卷: %s\n", volumeName)
	err = dm.client.VolumeRemove(dm.ctx, volumeName, false)
	if err != nil {
		return fmt.Errorf("删除Docker卷失败: %v", err)
	}

	fmt.Printf("Docker卷删除成功: %s\n", volumeName)
	return nil
}

// VolumeExists 检查Docker卷是否存在
// volumeName: 卷名称
// 返回: 是否存在和错误信息
func (dm *DockerManager) VolumeExists(volumeName string) (bool, error) {
	// 尝试获取卷信息
	_, err := dm.client.VolumeInspect(dm.ctx, volumeName)
	if err != nil {
		if client.IsErrNotFound(err) {
			return false, nil // 卷不存在
		}
		return false, fmt.Errorf("检查Docker卷失败: %v", err)
	}

	return true, nil
}

// PullImage 拉取Docker镜像
// imageName: 镜像名称
// 返回: 错误信息
func (dm *DockerManager) PullImage(imageName string) error {
	fmt.Printf("正在拉取Docker镜像: %s\n", imageName)

	// 拉取镜像
	reader, err := dm.client.ImagePull(dm.ctx, imageName, image.PullOptions{})
	if err != nil {
		return fmt.Errorf("拉取Docker镜像失败: %v", err)
	}
	defer reader.Close()

	// 读取拉取进度（可选，用于显示进度）
	_, err = io.Copy(io.Discard, reader)
	if err != nil {
		return fmt.Errorf("读取镜像拉取进度失败: %v", err)
	}

	fmt.Printf("Docker镜像拉取成功: %s\n", imageName)
	return nil
}

// PullImageWithProgress 拉取Docker镜像并显示进度
// imageName: 镜像名称
// 返回: 错误信息
func (dm *DockerManager) PullImageWithProgress(imageName string) error {
	fmt.Printf("正在拉取Docker镜像: %s\n", imageName)

	// 拉取镜像
	reader, err := dm.client.ImagePull(dm.ctx, imageName, image.PullOptions{})
	if err != nil {
		return fmt.Errorf("拉取Docker镜像失败: %v", err)
	}
	defer reader.Close()

	// 读取并显示拉取进度
	buffer := make([]byte, 1024)
	for {
		n, err := reader.Read(buffer)
		if n > 0 {
			// 解析JSON进度信息并显示
			progress := string(buffer[:n])
			if strings.Contains(progress, "Downloading") || strings.Contains(progress, "Extracting") {
				// 提取进度信息
				if strings.Contains(progress, "progress") {
					fmt.Print("\r正在下载... ")
				}
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("读取镜像拉取进度失败: %v", err)
		}
	}

	fmt.Printf("\nDocker镜像拉取成功: %s\n", imageName)
	return nil
}

// ImageExists 检查Docker镜像是否存在
// imageName: 镜像名称
// 返回: 是否存在和错误信息
func (dm *DockerManager) ImageExists(imageName string) (bool, error) {
	// 尝试获取镜像信息
	_, _, err := dm.client.ImageInspectWithRaw(dm.ctx, imageName)
	if err != nil {
		if client.IsErrNotFound(err) {
			return false, nil // 镜像不存在
		}
		return false, fmt.Errorf("检查Docker镜像失败: %v", err)
	}

	return true, nil
}

// WaitForImage 等待镜像拉取完成
// imageName: 镜像名称
// timeout: 超时时间（秒）
// 返回: 是否成功和错误信息
func (dm *DockerManager) WaitForImage(imageName string, timeout int) (bool, error) {
	log.Printf("等待镜像 %s 拉取完成...", imageName)

	// 每秒检查一次镜像是否存在
	for i := 0; i < timeout; i++ {
		exists, err := dm.ImageExists(imageName)
		if err != nil {
			return false, fmt.Errorf("检查镜像失败: %v", err)
		}

		if exists {
			log.Printf("镜像 %s 已准备就绪", imageName)
			return true, nil
		}

		// 检查是否正在拉取中
		if IsImagePulling(imageName) {
			log.Printf("镜像 %s 正在拉取中，继续等待...", imageName)
		}

		// 等待1秒后再次检查
		time.Sleep(1 * time.Second)
	}

	return false, fmt.Errorf("等待镜像 %s 超时（%d秒）", imageName, timeout)
}

// CreateContainer 创建ARK服务器容器（不自动启动）
// serverID: 服务器ID
// serverName: 服务器名称
// port: 游戏端口
// queryPort: 查询端口
// rconPort: RCON端口
// adminPassword: 管理员密码
// mapName: 地图名称
// autoRestart: 是否自动重启
// 返回: 容器ID和错误信息
func (dm *DockerManager) CreateContainer(serverID uint, serverName string, port, queryPort, rconPort int, adminPassword, mapName string, autoRestart bool) (string, error) {
	containerName := GetServerContainerName(serverID)
	volumeName := GetServerVolumeName(serverID)
	imageName := "tbro98/ase-server:latest"

	// 检查容器是否已存在
	if exists, err := dm.ContainerExists(containerName); err != nil {
		return "", fmt.Errorf("检查容器是否存在失败: %v", err)
	} else if exists {
		// 如果容器已存在，先删除它
		if err := dm.RemoveContainer(containerName); err != nil {
			return "", fmt.Errorf("删除已存在的容器失败: %v", err)
		}
	}

	// 检查镜像是否存在，如果不存在则等待拉取完成
	if exists, err := dm.ImageExists(imageName); err != nil {
		return "", fmt.Errorf("检查镜像是否存在失败: %v", err)
	} else if !exists {
		// 等待镜像拉取完成（最多等待60秒）
		if ready, err := dm.WaitForImage(imageName, 60); err != nil || !ready {
			return "", fmt.Errorf("镜像 %s 拉取超时或失败，请稍后重试", imageName)
		}
	}

	// 构建容器配置
	containerConfig := &container.Config{
		Image: imageName,
		Env: []string{
			"TZ=Asia/Shanghai",
		},
		ExposedPorts: nat.PortSet{
			nat.Port(fmt.Sprintf("%d/udp", port)):      struct{}{},
			nat.Port(fmt.Sprintf("%d/tcp", port)):      struct{}{},
			nat.Port(fmt.Sprintf("%d/udp", port+1)):    struct{}{},
			nat.Port(fmt.Sprintf("%d/tcp", port+1)):    struct{}{},
			nat.Port(fmt.Sprintf("%d/udp", queryPort)): struct{}{},
			nat.Port(fmt.Sprintf("%d/tcp", queryPort)): struct{}{},
			nat.Port(fmt.Sprintf("%d/udp", rconPort)):  struct{}{},
			nat.Port(fmt.Sprintf("%d/tcp", rconPort)):  struct{}{},
		},
	}

	// 根据autoRestart设置重启策略
	restartPolicyName := container.RestartPolicyMode("unless-stopped")
	if !autoRestart {
		restartPolicyName = container.RestartPolicyMode("no")
	}

	// 构建主机配置
	hostConfig := &container.HostConfig{
		RestartPolicy: container.RestartPolicy{
			Name: restartPolicyName,
		},
		PortBindings: nat.PortMap{
			nat.Port(fmt.Sprintf("%d/udp", port)): {
				{HostPort: fmt.Sprintf("%d", port)},
			},
			nat.Port(fmt.Sprintf("%d/tcp", port)): {
				{HostPort: fmt.Sprintf("%d", port)},
			},
			nat.Port(fmt.Sprintf("%d/udp", port+1)): {
				{HostPort: fmt.Sprintf("%d", port+1)},
			},
			nat.Port(fmt.Sprintf("%d/tcp", port+1)): {
				{HostPort: fmt.Sprintf("%d", port+1)},
			},
			nat.Port(fmt.Sprintf("%d/udp", queryPort)): {
				{HostPort: fmt.Sprintf("%d", queryPort)},
			},
			nat.Port(fmt.Sprintf("%d/tcp", queryPort)): {
				{HostPort: fmt.Sprintf("%d", queryPort)},
			},
			nat.Port(fmt.Sprintf("%d/udp", rconPort)): {
				{HostPort: fmt.Sprintf("%d", rconPort)},
			},
			nat.Port(fmt.Sprintf("%d/tcp", rconPort)): {
				{HostPort: fmt.Sprintf("%d", rconPort)},
			},
		},
		Binds: []string{
			fmt.Sprintf("%s:/home/steam/arkserver/ShooterGame/Saved", volumeName),
		},
	}

	// 创建容器
	fmt.Printf("正在创建Docker容器: %s\n", containerName)
	resp, err := dm.client.ContainerCreate(dm.ctx, containerConfig, hostConfig, nil, nil, containerName)
	if err != nil {
		return "", fmt.Errorf("创建Docker容器失败: %v", err)
	}

	fmt.Printf("Docker容器创建成功: %s (ID: %s)，容器处于停止状态，需要手动启动\n", containerName, resp.ID)
	return resp.ID, nil
}

// StartContainer 启动容器
// containerName: 容器名称
// 返回: 错误信息
func (dm *DockerManager) StartContainer(containerName string) error {
	fmt.Printf("正在启动容器: %s\n", containerName)
	err := dm.client.ContainerStart(dm.ctx, containerName, container.StartOptions{})
	if err != nil {
		return fmt.Errorf("启动Docker容器失败: %v", err)
	}

	fmt.Printf("容器启动成功: %s\n", containerName)
	return nil
}

// StopContainer 停止容器
// containerName: 容器名称
// 返回: 错误信息
func (dm *DockerManager) StopContainer(containerName string) error {
	fmt.Printf("正在停止容器: %s\n", containerName)

	// 设置30秒超时时间
	timeout := 30
	err := dm.client.ContainerStop(dm.ctx, containerName, container.StopOptions{
		Timeout: &timeout,
	})
	if err != nil {
		return fmt.Errorf("停止Docker容器失败: %v", err)
	}

	fmt.Printf("容器停止成功: %s\n", containerName)
	return nil
}

// RemoveContainer 删除容器
// containerName: 容器名称
// 返回: 错误信息
func (dm *DockerManager) RemoveContainer(containerName string) error {
	// 先尝试停止容器
	dm.StopContainer(containerName)

	// 删除容器
	fmt.Printf("正在删除容器: %s\n", containerName)
	err := dm.client.ContainerRemove(dm.ctx, containerName, container.RemoveOptions{
		Force: true,
	})
	if err != nil {
		return fmt.Errorf("删除Docker容器失败: %v", err)
	}

	fmt.Printf("容器删除成功: %s\n", containerName)
	return nil
}

// ContainerExists 检查容器是否存在
// containerName: 容器名称
// 返回: 是否存在和错误信息
func (dm *DockerManager) ContainerExists(containerName string) (bool, error) {
	// 使用inspect命令检查容器是否存在
	_, err := dm.client.ContainerInspect(dm.ctx, containerName)
	if err != nil {
		if client.IsErrNotFound(err) {
			return false, nil // 容器不存在
		}
		return false, fmt.Errorf("检查Docker容器失败: %v", err)
	}

	return true, nil
}

// GetContainerStatus 获取容器状态
// containerName: 容器名称
// 返回: 状态字符串和错误信息
func (dm *DockerManager) GetContainerStatus(containerName string) (string, error) {
	containerInfo, err := dm.client.ContainerInspect(dm.ctx, containerName)
	if err != nil {
		if client.IsErrNotFound(err) {
			return "not_found", nil
		}
		return "", fmt.Errorf("获取Docker容器状态失败: %v", err)
	}

	// 解析Docker状态为我们的状态格式
	state := containerInfo.State
	if state.Running {
		return "running", nil
	} else if state.Status == "exited" {
		return "stopped", nil
	} else if state.Status == "created" {
		return "stopped", nil
	} else if state.Status == "restarting" {
		return "starting", nil
	} else {
		return "unknown", nil
	}
}

// ExecuteCommand 在容器中执行命令
// containerName: 容器名称
// command: 要执行的命令
// 返回: 命令输出和错误信息
func (dm *DockerManager) ExecuteCommand(containerName string, command string) (string, error) {
	// 创建执行配置
	execConfig := container.ExecOptions{
		Cmd:          []string{"sh", "-c", command},
		AttachStdout: true,
		AttachStderr: true,
	}

	// 创建执行实例
	execResp, err := dm.client.ContainerExecCreate(dm.ctx, containerName, execConfig)
	if err != nil {
		return "", fmt.Errorf("创建执行实例失败: %v", err)
	}

	// 执行命令
	resp, err := dm.client.ContainerExecAttach(dm.ctx, execResp.ID, container.ExecAttachOptions{})
	if err != nil {
		return "", fmt.Errorf("执行命令失败: %v", err)
	}
	defer resp.Close()

	// 读取输出
	output, err := io.ReadAll(resp.Reader)
	if err != nil {
		return "", fmt.Errorf("读取命令输出失败: %v", err)
	}

	// 检查执行结果
	inspectResp, err := dm.client.ContainerExecInspect(dm.ctx, execResp.ID)
	if err != nil {
		return "", fmt.Errorf("检查执行结果失败: %v", err)
	}

	if inspectResp.ExitCode != 0 {
		return string(output), fmt.Errorf("命令执行失败，退出码: %d", inspectResp.ExitCode)
	}

	return string(output), nil
}

// ReadConfigFile 从容器卷中读取配置文件
// serverID: 服务器ID
// fileName: 文件名
// 返回: 文件内容和错误信息
func (dm *DockerManager) ReadConfigFile(serverID uint, fileName string) (string, error) {
	volumeName := GetServerVolumeName(serverID)
	alpineImage := "alpine:latest"

	// 检查Alpine镜像是否存在，如果不存在则等待拉取完成
	if exists, err := dm.ImageExists(alpineImage); err != nil {
		return "", fmt.Errorf("检查Alpine镜像是否存在失败: %v", err)
	} else if !exists {
		// 等待镜像拉取完成（最多等待30秒）
		if ready, err := dm.WaitForImage(alpineImage, 30); err != nil || !ready {
			return "", fmt.Errorf("Alpine镜像拉取超时或失败，请稍后重试")
		}
	}

	// 根据文件名确定路径
	var configPath string
	if fileName == ServerConfigFileName {
		// server.cfg 文件在根目录
		configPath = "/home/steam/arkserver/ShooterGame/Saved/server.cfg"
	} else {
		// 其他配置文件在 Config/WindowsServer 目录
		configPath = fmt.Sprintf("/home/steam/arkserver/ShooterGame/Saved/Config/WindowsServer/%s", fileName)
	}

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
	volumeName := GetServerVolumeName(serverID)
	alpineImage := "alpine:latest"

	// 检查Alpine镜像是否存在，如果不存在则等待拉取完成
	if exists, err := dm.ImageExists(alpineImage); err != nil {
		return fmt.Errorf("检查Alpine镜像是否存在失败: %v", err)
	} else if !exists {
		// 等待镜像拉取完成（最多等待30秒）
		if ready, err := dm.WaitForImage(alpineImage, 30); err != nil || !ready {
			return fmt.Errorf("Alpine镜像拉取超时或失败，请稍后重试")
		}
	}

	// 根据文件名确定路径
	var configPath string
	var configDir string

	if fileName == ServerConfigFileName {
		// server.cfg 文件在根目录
		configPath = "/home/steam/arkserver/ShooterGame/Saved/server.cfg"
		configDir = "/home/steam/arkserver/ShooterGame/Saved"
	} else {
		// 其他配置文件在 Config/WindowsServer 目录
		configPath = fmt.Sprintf("/home/steam/arkserver/ShooterGame/Saved/Config/WindowsServer/%s", fileName)
		configDir = "/home/steam/arkserver/ShooterGame/Saved/Config/WindowsServer"
	}

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

// GetServerContainerName 获取服务器容器名称
// serverID: 服务器ID
// 返回: 容器名称
func GetServerContainerName(serverID uint) string {
	return fmt.Sprintf("ase-server-%d", serverID)
}

// GetServerVolumeName 获取服务器卷名称
// serverID: 服务器ID
// 返回: 卷名称
func GetServerVolumeName(serverID uint) string {
	return fmt.Sprintf("ase-server-%d", serverID)
}

// CheckDockerStatus 检查Docker环境状态
// 返回: Docker是否可用和错误信息
func CheckDockerStatus() error {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return fmt.Errorf("Docker客户端创建失败: %v", err)
	}
	defer cli.Close()

	// 检查Docker服务是否运行
	_, err = cli.Ping(context.Background())
	if err != nil {
		return fmt.Errorf("Docker服务未运行或无法连接: %v", err)
	}

	return nil
}
