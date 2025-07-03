package utils

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/docker/docker/api/types/container"
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
	volumeName := fmt.Sprintf("ark-server-%d", serverID)

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

// CreateContainer 创建并启动ARK服务器容器
// serverID: 服务器ID
// serverName: 服务器名称
// port: 游戏端口
// queryPort: 查询端口
// rconPort: RCON端口
// adminPassword: 管理员密码
// mapName: 地图名称
// 返回: 容器ID和错误信息
func (dm *DockerManager) CreateContainer(serverID uint, serverName string, port, queryPort, rconPort int, adminPassword, mapName string) (string, error) {
	containerName := fmt.Sprintf("ark-server-%d", serverID)
	volumeName := fmt.Sprintf("ark-server-%d", serverID)

	// 检查容器是否已存在
	if exists, err := dm.ContainerExists(containerName); err != nil {
		return "", fmt.Errorf("检查容器是否存在失败: %v", err)
	} else if exists {
		// 如果容器已存在，先删除它
		if err := dm.RemoveContainer(containerName); err != nil {
			return "", fmt.Errorf("删除已存在的容器失败: %v", err)
		}
	}

	// 构建容器配置
	containerConfig := &container.Config{
		Image: "tbro98/ase-server:latest",
		Env: []string{
			"TZ=Asia/Shanghai",
		},
		ExposedPorts: nat.PortSet{
			"7777/udp": struct{}{},
			"7777/tcp": struct{}{},
			"7778/udp": struct{}{},
			"7778/tcp": struct{}{},
			nat.Port(fmt.Sprintf("%d/udp", queryPort)): struct{}{},
			nat.Port(fmt.Sprintf("%d/tcp", queryPort)): struct{}{},
			nat.Port(fmt.Sprintf("%d/udp", rconPort)):  struct{}{},
			nat.Port(fmt.Sprintf("%d/tcp", rconPort)):  struct{}{},
		},
	}

	// 构建主机配置
	hostConfig := &container.HostConfig{
		RestartPolicy: container.RestartPolicy{
			Name: "unless-stopped",
		},
		PortBindings: nat.PortMap{
			"7777/udp": {
				{HostPort: fmt.Sprintf("%d", port)},
			},
			"7777/tcp": {
				{HostPort: fmt.Sprintf("%d", port)},
			},
			"7778/udp": {
				{HostPort: fmt.Sprintf("%d", port+1)},
			},
			"7778/tcp": {
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

	// 启动容器
	fmt.Printf("正在启动Docker容器: %s\n", containerName)
	err = dm.client.ContainerStart(dm.ctx, resp.ID, container.StartOptions{})
	if err != nil {
		return "", fmt.Errorf("启动Docker容器失败: %v", err)
	}

	fmt.Printf("Docker容器创建并启动成功: %s (ID: %s)\n", containerName, resp.ID)
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
	volumeName := fmt.Sprintf("ark-server-%d", serverID)

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
		Image: "alpine:latest",
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
	volumeName := fmt.Sprintf("ark-server-%d", serverID)

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
		Image: "alpine:latest",
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
		Image: "alpine:latest",
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
	return fmt.Sprintf("ark-server-%d", serverID)
}

// GetServerVolumeName 获取服务器卷名称
// serverID: 服务器ID
// 返回: 卷名称
func GetServerVolumeName(serverID uint) string {
	return fmt.Sprintf("ark-server-%d", serverID)
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
