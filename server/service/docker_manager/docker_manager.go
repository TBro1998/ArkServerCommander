package docker_manager

import (
	"ark-server-commander/config"
	"ark-server-commander/database"
	"ark-server-commander/models"
	"ark-server-commander/utils"
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/containerd/errdefs"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"go.uber.org/zap"
)

// DockerManager Docker管理器结构体
type DockerManager struct {
	client *client.Client
	ctx    context.Context
}

// 单例相关变量
var (
	instance *DockerManager
	once     sync.Once
)

// GetDockerManager 获取Docker管理器单例实例
func GetDockerManager() (*DockerManager, error) {
	var err error
	once.Do(func() {
		// 创建新实例
		cli, clientErr := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
		if clientErr != nil {
			err = fmt.Errorf("创建Docker客户端失败: %v", clientErr)
			return
		}

		instance = &DockerManager{
			client: cli,
			ctx:    context.Background(),
		}
	})

	if err != nil {
		return nil, err
	}

	return instance, nil
}

// CloseDockerManager 关闭Docker管理器（通常在程序退出时调用）
func CloseDockerManager() error {
	if instance != nil && instance.client != nil {
		err := instance.client.Close()
		instance = nil
		return err
	}
	return nil
}

// CheckDockerStatus 检查Docker环境状态
// 返回: Docker是否可用和错误信息
func CheckDockerStatus() error {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return fmt.Errorf("docker客户端创建失败: %v", err)
	}
	defer cli.Close()

	// 检查Docker服务是否运行
	_, err = cli.Ping(context.Background())
	if err != nil {
		return fmt.Errorf("docker服务未运行或无法连接: %v", err)
	}

	return nil
}

// CreateContainer 创建ARK服务器容器（不自动启动）
// serverID: 服务器ID
// serverName: 服务器名称
// port: 游戏端口
// queryPort: 查询端口
// rconPort: RCON端口
// adminPassword: 管理员密码
// mapName: 地图名称
// gameModIds: 游戏模组ID列表，用逗号分隔
// autoRestart: 是否自动重启
// 返回: 容器ID和错误信息
func (dm *DockerManager) CreateContainer(serverID uint, serverName string, port, queryPort, rconPort int, adminPassword, mapName, gameModIds string, autoRestart bool) (string, error) {
	containerName := utils.GetServerContainerName(serverID)
	volumeName := utils.GetServerVolumeName(serverID)

	// 查数据库获取Server对象：镜像由是否启用 ArkApi 决定
	var server models.Server
	if err := database.DB.Where("id = ?", serverID).First(&server).Error; err != nil {
		return "", fmt.Errorf("获取服务器信息失败: %v", err)
	}
	imageName := utils.ArkServerImageFor(server.ArkApi)

	// 检查容器是否已存在
	if exists, err := dm.ContainerExists(containerName); err != nil {
		return "", fmt.Errorf("检查容器是否存在失败: %v", err)
	} else if exists {
		// 如果容器已存在，先删除它
		if err := dm.RemoveContainer(containerName); err != nil {
			return "", fmt.Errorf("删除已存在的容器失败: %v", err)
		}
	}

	// 检查镜像是否存在
	exists, err := dm.ImageExists(imageName)
	if err != nil {
		return "", fmt.Errorf("检查镜像是否存在失败: %v", err)
	}
	if !exists {
		return "", fmt.Errorf("镜像 %s 不存在，请等待镜像下载完成", imageName)
	}

	// 反序列化ServerArgsJSON
	serverArgs := models.NewServerArgs()
	if server.ServerArgsJSON != "" && server.ServerArgsJSON != "{}" {
		_ = json.Unmarshal([]byte(server.ServerArgsJSON), serverArgs)
	} else {
		serverArgs = models.FromServer(server)
	}
	argsString := serverArgs.GenerateArgsString(server)

	// 构建环境变量
	envVars := []string{
		"TZ=" + config.ContainerTZ,
		fmt.Sprintf("SERVER_ARGS=%s", argsString),
	}
	if server.GameModIds != "" {
		envVars = append(envVars, fmt.Sprintf("GameModIds=%s", server.GameModIds))
	}
	// 记录挂载的集群卷，启动时据此判断是否需要重建容器（旧容器没有这个变量）
	if cv := ExpectedClusterVolume(server.ClusterID); cv != "" {
		envVars = append(envVars, "CLUSTER_VOLUME="+cv)
	}

	// 构建容器配置
	containerConfig := &container.Config{
		Image: imageName,
		Env:   envVars,
		// 每个参数作为独立的一项传给入口脚本（"$@"），名称、密码里的空格不会被拆开；
		// SERVER_ARGS 保留用于重建判断，也兼容还不认识 Cmd 的旧镜像
		Cmd: serverArgs.GenerateArgs(server),
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
		// 服务端文件卷挂到安装目录，数据卷和插件卷嵌套挂载在其中。
		// 插件卷只对 Proton 版有用，Linux 版也照样挂载（只是多出 Win64/ArkApi/Plugins 目录），切换时挂载不变
		Binds: []string{
			fmt.Sprintf("%s:/home/steam/arkserver", utils.GetServerFilesVolumeName(serverID)),
			fmt.Sprintf("%s:/home/steam/arkserver/ShooterGame/Saved", volumeName),
			fmt.Sprintf("%s:/home/steam/arkserver/ShooterGame/Binaries/Win64/ArkApi/Plugins", utils.GetServerPluginsVolumeName(serverID)),
		},
	}

	// 集群：同一 cluster_id 的服务器共享 ARK 默认的集群目录 Saved/clusters，跨服传送的角色和物品存放在这里。
	// 先把卷建成绑定到数据目录 clusters/<cluster_id> 的命名卷，否则 Docker 会自动创建一个普通卷
	if cv := ExpectedClusterVolume(server.ClusterID); cv != "" {
		if err := dm.EnsureClusterVolume(server.ClusterID); err != nil {
			return "", fmt.Errorf("创建集群卷失败: %v", err)
		}
		hostConfig.Binds = append(hostConfig.Binds, cv+":/home/steam/arkserver/ShooterGame/Saved/clusters")
	}

	// 创建容器
	utils.Infof("正在创建Docker容器: %s", containerName)
	resp, err := dm.client.ContainerCreate(dm.ctx, containerConfig, hostConfig, nil, nil, containerName)
	if err != nil {
		return "", fmt.Errorf("创建Docker容器失败: %v", err)
	}

	utils.Info("Docker容器创建成功，处于停止状态",
		zap.String("container_name", containerName),
		zap.String("container_id", resp.ID),
	)
	return resp.ID, nil
}

// StartContainer 启动容器
// containerName: 容器名称
// 返回: 错误信息
func (dm *DockerManager) StartContainer(containerName string) error {
	utils.Infof("正在启动容器: %s", containerName)
	err := dm.client.ContainerStart(dm.ctx, containerName, container.StartOptions{})
	if err != nil {
		return fmt.Errorf("启动Docker容器失败: %v", err)
	}

	utils.Info("容器启动成功", zap.String("container_name", containerName))
	return nil
}

// StopContainer 停止容器
// containerName: 容器名称
// 返回: 错误信息
func (dm *DockerManager) StopContainer(containerName string) error {
	utils.Infof("正在停止容器: %s", containerName)

	// 设置30秒超时时间
	timeout := 30
	err := dm.client.ContainerStop(dm.ctx, containerName, container.StopOptions{
		Timeout: &timeout,
	})
	if err != nil {
		return fmt.Errorf("停止Docker容器失败: %v", err)
	}

	utils.Info("容器停止成功", zap.String("container_name", containerName))
	return nil
}

// RemoveContainer 删除容器
// containerName: 容器名称
// 返回: 错误信息
func (dm *DockerManager) RemoveContainer(containerName string) error {
	// 先尝试停止容器
	dm.StopContainer(containerName)

	// 删除容器
	utils.Infof("正在删除容器: %s", containerName)
	err := dm.client.ContainerRemove(dm.ctx, containerName, container.RemoveOptions{
		Force: true,
	})
	if err != nil {
		return fmt.Errorf("删除Docker容器失败: %v", err)
	}

	utils.Info("容器删除成功", zap.String("container_name", containerName))
	return nil
}

// ContainerExists 检查容器是否存在
// containerName: 容器名称
// 返回: 是否存在和错误信息
func (dm *DockerManager) ContainerExists(containerName string) (bool, error) {
	// 使用inspect命令检查容器是否存在
	_, err := dm.client.ContainerInspect(dm.ctx, containerName)
	if err != nil {
		if errdefs.IsNotFound(err) {
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
		if errdefs.IsNotFound(err) {
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

// GetContainerEnvVars 获取容器的环境变量
// containerName: 容器名称
// 返回: 环境变量映射和错误信息
func (dm *DockerManager) GetContainerEnvVars(containerName string) (map[string]string, error) {
	containerInfo, err := dm.client.ContainerInspect(dm.ctx, containerName)
	if err != nil {
		if errdefs.IsNotFound(err) {
			return nil, fmt.Errorf("容器不存在: %s", containerName)
		}
		return nil, fmt.Errorf("获取Docker容器信息失败: %v", err)
	}

	// 解析环境变量
	envVars := make(map[string]string)
	for _, env := range containerInfo.Config.Env {
		// 环境变量格式: KEY=VALUE
		for i, char := range env {
			if char == '=' {
				key := env[:i]
				value := env[i+1:]
				envVars[key] = value
				break
			}
		}
	}

	return envVars, nil
}

// ExpectedClusterVolume 服务器应挂载的集群卷名，没有设置（或格式不合法的旧数据）时返回空字符串
func ExpectedClusterVolume(clusterID string) string {
	if clusterID == "" || !utils.ClusterIDPattern.MatchString(clusterID) {
		return ""
	}
	return utils.GetClusterVolumeName(clusterID)
}

// GetContainerRestartPolicy 返回容器的重启策略名称（如 unless-stopped、no）
func (dm *DockerManager) GetContainerRestartPolicy(containerName string) (string, error) {
	info, err := dm.client.ContainerInspect(dm.ctx, containerName)
	if err != nil {
		return "", fmt.Errorf("获取Docker容器信息失败: %v", err)
	}
	if info.HostConfig == nil {
		return "", nil
	}
	return string(info.HostConfig.RestartPolicy.Name), nil
}
