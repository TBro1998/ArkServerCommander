package docker_manager

import (
	"ark-server-commander/database"
	"ark-server-commander/models"
	"ark-server-commander/utils"
	"encoding/json"
	"fmt"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	"go.uber.org/zap"
)

// CreateContainerWithRollback 创建容器（带回滚机制）
// 如果创建过程中任何步骤失败，会自动清理已创建的资源
func (dm *DockerManager) CreateContainerWithRollback(serverID uint, serverName string, port, queryPort, rconPort int, adminPassword, mapName, gameModIds string, autoRestart bool) (containerID string, err error) {
	// 创建回滚管理器
	rollback := NewRollbackManager()

	// 使用 defer 确保发生错误时执行回滚
	defer func() {
		if err != nil && rollback.Count() > 0 {
			utils.Warn("容器创建失败，开始执行回滚", zap.Error(err))
			if rollbackErr := rollback.Rollback(); rollbackErr != nil {
				utils.Error("回滚过程中发生错误", zap.Error(rollbackErr))
			}
		}
	}()

	containerName := utils.GetServerContainerName(serverID)
	volumeName := utils.GetServerVolumeName(serverID)
	imageName := "tbro98/ase-server:latest"

	utils.Info("开始创建容器（带回滚保护）",
		zap.String("container", containerName),
		zap.Uint("server_id", serverID))

	// 步骤1: 检查并删除已存在的容器
	if exists, checkErr := dm.ContainerExists(containerName); checkErr != nil {
		err = fmt.Errorf("检查容器是否存在失败: %w", checkErr)
		return "", err
	} else if exists {
		utils.Info("发现已存在的容器，准备删除", zap.String("container", containerName))
		if removeErr := dm.RemoveContainer(containerName); removeErr != nil {
			err = fmt.Errorf("删除已存在的容器失败: %w", removeErr)
			return "", err
		}
	}

	// 步骤2: 检查镜像是否存在
	imageExists, checkErr := dm.ImageExists(imageName)
	if checkErr != nil {
		err = fmt.Errorf("检查镜像是否存在失败: %w", checkErr)
		return "", err
	}
	if !imageExists {
		err = fmt.Errorf("镜像 %s 不存在，请等待镜像下载完成", imageName)
		return "", err
	}

	// 步骤3: 获取服务器信息
	var server models.Server
	if dbErr := database.DB.Where("id = ?", serverID).First(&server).Error; dbErr != nil {
		err = fmt.Errorf("获取服务器信息失败: %w", dbErr)
		return "", err
	}

	// 步骤4: 构建启动参数
	serverArgs := models.NewServerArgs()
	if server.ServerArgsJSON != "" && server.ServerArgsJSON != "{}" {
		_ = json.Unmarshal([]byte(server.ServerArgsJSON), serverArgs)
	} else {
		serverArgs = models.FromServer(server)
	}
	argsString := serverArgs.GenerateArgsString(server)

	// 步骤5: 构建环境变量
	envVars := []string{
		"TZ=Asia/Shanghai",
		fmt.Sprintf("SERVER_ARGS=%s", argsString),
	}
	if server.GameModIds != "" {
		envVars = append(envVars, fmt.Sprintf("GameModIds=%s", server.GameModIds))
	}

	// 步骤6: 构建容器配置
	containerConfig := &container.Config{
		Image: imageName,
		Env:   envVars,
		ExposedPorts: nat.PortSet{
			nat.Port(fmt.Sprintf("%d/udp", port)):      struct{}{},
			nat.Port(fmt.Sprintf("%d/tcp", port)):      struct{}{},
			nat.Port(fmt.Sprintf("%d/udp", port+1)):    struct{}{},
			nat.Port(fmt.Sprintf("%d/tcp", port+1)):    struct{}{},
			nat.Port(fmt.Sprintf("%d/udp", queryPort)): struct{}{},
			nat.Port(fmt.Sprintf("%d/tcp", queryPort)): struct{}{},
			nat.Port(fmt.Sprintf("%d/udp", rconPort)):  struct{}{},
			nat.Port(fmt.Sprintf("%d/tcp", rconPort)):  struct{},
		},
	}

	// 步骤7: 设置重启策略
	restartPolicyName := container.RestartPolicyMode("unless-stopped")
	if !autoRestart {
		restartPolicyName = container.RestartPolicyMode("no")
	}

	// 步骤8: 构建主机配置
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
			fmt.Sprintf("%s:/home/steam/arkserver/ShooterGame/Binaries/Win64/ArkApi/Plugins", utils.GetServerPluginsVolumeName(serverID)),
		},
	}

	// 步骤9: 创建容器
	utils.Info("正在创建Docker容器", zap.String("container", containerName))
	resp, createErr := dm.client.ContainerCreate(dm.ctx, containerConfig, hostConfig, nil, nil, containerName)
	if createErr != nil {
		err = fmt.Errorf("创建Docker容器失败: %w", createErr)
		return "", err
	}

	containerID = resp.ID

	// 添加回滚操作：删除刚创建的容器
	rollback.AddAction("container", containerName, "删除容器", func() error {
		return dm.RemoveContainer(containerName)
	})

	utils.Info("Docker容器创建成功",
		zap.String("container_name", containerName),
		zap.String("container_id", containerID))

	// 成功完成，清空回滚操作
	rollback.Clear()

	return containerID, nil
}
