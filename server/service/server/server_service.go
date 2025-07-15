package server

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"ark-server-commander/database"
	"ark-server-commander/models"
	"ark-server-commander/service/docker_manager"
	"ark-server-commander/utils"
)

// ServerService 服务器管理业务逻辑服务
type ServerService struct{}

// NewServerService 创建服务器服务实例
func NewServerService() *ServerService {
	return &ServerService{}
}

// GetServers 获取用户的所有服务器
func (s *ServerService) GetServers(userID uint) ([]models.ServerResponse, error) {
	var servers []models.Server
	if err := database.DB.Where("user_id = ?", userID).Find(&servers).Error; err != nil {
		return nil, fmt.Errorf("获取服务器列表失败: %w", err)
	}

	dockerManager, err := docker_manager.GetDockerManager()
	if err != nil {
		return nil, fmt.Errorf("获取Docker管理器失败: %w", err)
	}

	var serverResponses []models.ServerResponse
	for _, server := range servers {
		// 获取Docker容器实时状态
		containerName := utils.GetServerContainerName(server.ID)
		realTimeStatus := server.Status

		// 检查容器是否存在
		containerExists, err := dockerManager.ContainerExists(containerName)
		if err == nil && containerExists {
			if dockerStatus, err := dockerManager.GetContainerStatus(containerName); err == nil {
				realTimeStatus = dockerStatus

				// 如果实时状态与数据库状态不同，更新数据库（异步）
				if realTimeStatus != server.Status {
					go func(s models.Server, status string) {
						database.DB.Model(&s).Update("status", status)
					}(server, realTimeStatus)
				}
			}
		} else if err == nil && !containerExists && server.Status == "running" {
			// 如果容器不存在但数据库状态是运行中，更新为停止状态
			realTimeStatus = "stopped"
			go func(s models.Server) {
				database.DB.Model(&s).Update("status", "stopped")
			}(server)
		}

		serverResponses = append(serverResponses, models.ServerResponse{
			ID:            server.ID,
			Identifier:    server.Identifier,
			SessionName:   server.SessionName,
			ClusterID:     server.ClusterID,
			Port:          server.Port,
			QueryPort:     server.QueryPort,
			RCONPort:      server.RCONPort,
			AdminPassword: server.AdminPassword,
			Map:           server.Map,
			MaxPlayers:    server.MaxPlayers,
			GameModIds:    server.GameModIds,
			Status:        realTimeStatus,
			AutoRestart:   server.AutoRestart,
			UserID:        server.UserID,
			CreatedAt:     server.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:     server.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return serverResponses, nil
}

// CreateServer 创建新服务器
func (s *ServerService) CreateServer(userID uint, req models.ServerRequest) (*models.ServerResponse, error) {
	// 检查服务器标识是否已存在
	var existingServer models.Server
	if err := database.DB.Where("identifier = ? AND user_id = ?", req.Identifier, userID).First(&existingServer).Error; err == nil {
		return nil, fmt.Errorf("服务器标识已存在")
	}

	// 设置默认值
	if req.Map == "" {
		req.Map = "TheIsland"
	}
	if req.MaxPlayers == 0 {
		req.MaxPlayers = 70
	}
	if req.AutoRestart == nil {
		defaultVal := true
		req.AutoRestart = &defaultVal
	}

	// 开始数据库事务
	tx := database.DB.Begin()
	if tx.Error != nil {
		return nil, fmt.Errorf("数据库事务启动失败: %w", tx.Error)
	}

	// 创建服务器
	server := models.Server{
		Identifier:    req.Identifier,
		SessionName:   req.SessionName,
		ClusterID:     req.ClusterID,
		Port:          req.Port,
		QueryPort:     req.QueryPort,
		RCONPort:      req.RCONPort,
		AdminPassword: req.AdminPassword,
		Map:           req.Map,
		MaxPlayers:    req.MaxPlayers,
		GameModIds:    req.GameModIds,
		Status:        "stopped",
		AutoRestart:   *req.AutoRestart,
		UserID:        userID,
	}

	if req.ServerArgs != nil {
		argsJson, err := json.Marshal(req.ServerArgs)
		if err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("启动参数格式错误: %w", err)
		}
		server.ServerArgsJSON = string(argsJson)
	} else {
		server.ServerArgsJSON = "{}"
	}

	if err := tx.Create(&server).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("服务器创建失败: %w", err)
	}

	// 检查Docker环境
	if err := docker_manager.CheckDockerStatus(); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("Docker环境检查失败: %w", err)
	}

	// 创建Docker卷
	dockerManager, err := docker_manager.GetDockerManager()
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("获取Docker管理器失败: %w", err)
	}

	volumeName, err := dockerManager.CreateVolume(server.ID)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("创建Docker卷失败: %w", err)
	}

	// 处理配置文件
	var gameUserSettings string
	var gameIni string

	if req.GameUserSettings != "" {
		if err := utils.ValidateINIContent(req.GameUserSettings); err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("GameUserSettings.ini格式错误: %w", err)
		}
		gameUserSettings = req.GameUserSettings
	} else {
		gameUserSettings = utils.GetDefaultGameUserSettings(server.Identifier, server.Map, 70)
	}

	if req.GameIni != "" {
		if err := utils.ValidateINIContent(req.GameIni); err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("Game.ini格式错误: %w", err)
		}
		gameIni = req.GameIni
	} else {
		gameIni = utils.GetDefaultGameIni()
	}

	// 写入配置文件
	if err := dockerManager.WriteConfigFile(server.ID, utils.GameUserSettingsFileName, gameUserSettings); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("写入GameUserSettings.ini失败: %w", err)
	}

	if err := dockerManager.WriteConfigFile(server.ID, utils.GameIniFileName, gameIni); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("写入Game.ini失败: %w", err)
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		dockerManager.RemoveVolume(volumeName)
		return nil, fmt.Errorf("数据库提交失败: %w", err)
	}

	// 构建响应
	response := models.ServerResponse{
		ID:            server.ID,
		Identifier:    server.Identifier,
		SessionName:   server.SessionName,
		ClusterID:     server.ClusterID,
		Port:          server.Port,
		QueryPort:     server.QueryPort,
		RCONPort:      server.RCONPort,
		AdminPassword: server.AdminPassword,
		Map:           server.Map,
		GameModIds:    server.GameModIds,
		Status:        server.Status,
		AutoRestart:   server.AutoRestart,
		UserID:        server.UserID,
		CreatedAt:     server.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:     server.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	// 读取配置文件内容
	if gameUserSettings, err := dockerManager.ReadConfigFile(uint(server.ID), utils.GameUserSettingsFileName); err == nil {
		response.GameUserSettings = gameUserSettings
	}
	if gameIni, err := dockerManager.ReadConfigFile(uint(server.ID), utils.GameIniFileName); err == nil {
		response.GameIni = gameIni
	}

	return &response, nil
}

// GetServer 获取单个服务器信息
func (s *ServerService) GetServer(userID uint, serverID string) (*models.ServerResponse, error) {
	id, err := strconv.ParseUint(serverID, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("无效的服务器ID")
	}

	var server models.Server
	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&server).Error; err != nil {
		return nil, fmt.Errorf("服务器不存在")
	}

	// 解析启动参数
	var serverArgs *models.ServerArgs
	if server.ServerArgsJSON != "" && server.ServerArgsJSON != "{}" {
		serverArgs = models.NewServerArgs()
		if err := json.Unmarshal([]byte(server.ServerArgsJSON), serverArgs); err != nil {
			serverArgs = models.FromServer(server)
		}
	} else {
		serverArgs = models.FromServer(server)
	}

	response := models.ServerResponse{
		ID:            server.ID,
		Identifier:    server.Identifier,
		SessionName:   server.SessionName,
		ClusterID:     server.ClusterID,
		Port:          server.Port,
		QueryPort:     server.QueryPort,
		RCONPort:      server.RCONPort,
		AdminPassword: server.AdminPassword,
		Map:           server.Map,
		MaxPlayers:    server.MaxPlayers,
		GameModIds:    server.GameModIds,
		Status:        server.Status,
		AutoRestart:   server.AutoRestart,
		UserID:        server.UserID,
		CreatedAt:     server.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:     server.UpdatedAt.Format("2006-01-02 15:04:05"),
		ServerArgs:    serverArgs,
		GeneratedArgs: serverArgs.GenerateArgsString(server),
	}

	// 读取配置文件内容
	dockerManager, err := docker_manager.GetDockerManager()
	if err != nil {
		return nil, fmt.Errorf("获取Docker管理器失败: %w", err)
	}

	if gameUserSettings, err := dockerManager.ReadConfigFile(uint(id), utils.GameUserSettingsFileName); err == nil {
		response.GameUserSettings = gameUserSettings
	}
	if gameIni, err := dockerManager.ReadConfigFile(uint(id), utils.GameIniFileName); err == nil {
		response.GameIni = gameIni
	}

	return &response, nil
}

// GetServerRCON 获取服务器RCON连接信息
func (s *ServerService) GetServerRCON(userID uint, serverID string) (map[string]interface{}, error) {
	id, err := strconv.ParseUint(serverID, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("无效的服务器ID")
	}

	var server models.Server
	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&server).Error; err != nil {
		return nil, fmt.Errorf("服务器不存在")
	}

	return map[string]interface{}{
		"server_id":         server.ID,
		"server_identifier": server.Identifier,
		"rcon_port":         server.RCONPort,
		"admin_password":    server.AdminPassword,
	}, nil
}

// ExecuteRCONCommand 执行RCON命令
func (s *ServerService) ExecuteRCONCommand(userID uint, serverID string, command string) (map[string]interface{}, error) {
	id, err := strconv.ParseUint(serverID, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("无效的服务器ID")
	}

	var server models.Server
	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&server).Error; err != nil {
		return nil, fmt.Errorf("服务器不存在")
	}

	if server.Status != "running" {
		return nil, fmt.Errorf("服务器未运行，无法执行RCON命令")
	}

	dockerManager, err := docker_manager.GetDockerManager()
	if err != nil {
		return nil, fmt.Errorf("获取Docker管理器失败: %w", err)
	}

	containerName := utils.GetServerContainerName(server.ID)
	containerExists, err := dockerManager.ContainerExists(containerName)
	if err != nil {
		return nil, fmt.Errorf("检查容器状态失败: %w", err)
	}
	if !containerExists {
		return nil, fmt.Errorf("服务器容器不存在，无法执行RCON命令")
	}

	rconCommand := fmt.Sprintf("echo '%s' | /ark/rcon -H localhost -P %d -p '%s'", command, server.RCONPort, server.AdminPassword)
	output, err := dockerManager.ExecuteCommand(containerName, rconCommand)
	if err != nil {
		return nil, fmt.Errorf("执行RCON命令失败: %w", err)
	}

	return map[string]interface{}{
		"command": command,
		"output":  output,
	}, nil
}

// UpdateServer 更新服务器配置
func (s *ServerService) UpdateServer(userID uint, serverID string, req models.ServerUpdateRequest) (*models.ServerResponse, bool, error) {
	id, err := strconv.ParseUint(serverID, 10, 32)
	if err != nil {
		return nil, false, fmt.Errorf("无效的服务器ID")
	}

	var server models.Server
	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&server).Error; err != nil {
		return nil, false, fmt.Errorf("服务器不存在")
	}

	// 检查标识是否冲突
	if req.Identifier != "" && req.Identifier != server.Identifier {
		var existingServer models.Server
		if err := database.DB.Where("identifier = ? AND user_id = ? AND id != ?", req.Identifier, userID, id).First(&existingServer).Error; err == nil {
			return nil, false, fmt.Errorf("服务器标识已存在")
		}
		server.Identifier = req.Identifier
	}

	// 更新字段
	if req.SessionName != "" {
		server.SessionName = req.SessionName
	}
	if req.ClusterID != "" {
		server.ClusterID = req.ClusterID
	}
	if req.Port > 0 {
		server.Port = req.Port
	}
	if req.QueryPort > 0 {
		server.QueryPort = req.QueryPort
	}
	if req.RCONPort > 0 {
		server.RCONPort = req.RCONPort
	}
	if req.AdminPassword != "" {
		server.AdminPassword = req.AdminPassword
	}
	if req.Map != "" {
		server.Map = req.Map
	}
	if req.MaxPlayers > 0 {
		server.MaxPlayers = req.MaxPlayers
	}
	if req.GameModIds != "" {
		server.GameModIds = req.GameModIds
	}

	// 检查启动参数是否发生变化
	argsChanged := false
	if req.ServerArgs != nil {
		argsJson, err := json.Marshal(req.ServerArgs)
		if err != nil {
			return nil, false, fmt.Errorf("启动参数格式错误: %w", err)
		}
		newArgsJSON := string(argsJson)
		if server.ServerArgsJSON != newArgsJSON {
			argsChanged = true
			server.ServerArgsJSON = newArgsJSON
		}
	}

	if err := database.DB.Save(&server).Error; err != nil {
		return nil, false, fmt.Errorf("服务器更新失败: %w", err)
	}

	// 处理配置文件更新
	if req.GameUserSettings != "" || req.GameIni != "" {
		dockerManager, err := docker_manager.GetDockerManager()
		if err != nil {
			return nil, false, fmt.Errorf("获取Docker管理器失败: %w", err)
		}

		if req.GameUserSettings != "" {
			if err := utils.ValidateINIContent(req.GameUserSettings); err != nil {
				return nil, false, fmt.Errorf("GameUserSettings.ini格式错误: %w", err)
			}
			if err := dockerManager.WriteConfigFile(uint(id), utils.GameUserSettingsFileName, req.GameUserSettings); err != nil {
				return nil, false, fmt.Errorf("写入GameUserSettings.ini失败: %w", err)
			}
		}

		if req.GameIni != "" {
			if err := utils.ValidateINIContent(req.GameIni); err != nil {
				return nil, false, fmt.Errorf("Game.ini格式错误: %w", err)
			}
			if err := dockerManager.WriteConfigFile(uint(id), utils.GameIniFileName, req.GameIni); err != nil {
				return nil, false, fmt.Errorf("写入Game.ini失败: %w", err)
			}
		}
	}

	// 构建响应
	response := models.ServerResponse{
		ID:            server.ID,
		Identifier:    server.Identifier,
		SessionName:   server.SessionName,
		ClusterID:     server.ClusterID,
		Port:          server.Port,
		QueryPort:     server.QueryPort,
		RCONPort:      server.RCONPort,
		AdminPassword: server.AdminPassword,
		Map:           server.Map,
		MaxPlayers:    server.MaxPlayers,
		GameModIds:    server.GameModIds,
		Status:        server.Status,
		AutoRestart:   server.AutoRestart,
		UserID:        server.UserID,
		CreatedAt:     server.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:     server.UpdatedAt.Format("2006-01-02 15:04:05"),
		ServerArgs:    models.FromServer(server),
	}

	// 读取配置文件内容
	dockerManager, err := docker_manager.GetDockerManager()
	if err != nil {
		return nil, false, fmt.Errorf("获取Docker管理器失败: %w", err)
	}

	if gameUserSettings, err := dockerManager.ReadConfigFile(uint(id), utils.GameUserSettingsFileName); err == nil {
		response.GameUserSettings = gameUserSettings
	}
	if gameIni, err := dockerManager.ReadConfigFile(uint(id), utils.GameIniFileName); err == nil {
		response.GameIni = gameIni
	}

	return &response, argsChanged, nil
}

// DeleteServer 删除服务器
func (s *ServerService) DeleteServer(userID uint, serverID string) error {
	id, err := strconv.ParseUint(serverID, 10, 32)
	if err != nil {
		return fmt.Errorf("无效的服务器ID")
	}

	var server models.Server
	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&server).Error; err != nil {
		return fmt.Errorf("服务器不存在")
	}

	if server.Status == "running" {
		return fmt.Errorf("无法删除正在运行的服务器，请先停止服务器")
	}

	// 软删除服务器
	if err := database.DB.Delete(&server).Error; err != nil {
		return fmt.Errorf("服务器删除失败: %w", err)
	}

	// 删除Docker容器
	dockerManager, err := docker_manager.GetDockerManager()
	if err != nil {
		return fmt.Errorf("获取Docker管理器失败: %w", err)
	}

	containerName := utils.GetServerContainerName(server.ID)
	containerExists, err := dockerManager.ContainerExists(containerName)
	if err != nil {
		fmt.Printf("Warning: Failed to check container existence: %v\n", err)
	} else if containerExists {
		if err := dockerManager.RemoveContainer(containerName); err != nil {
			fmt.Printf("Warning: Failed to remove Docker container: %v\n", err)
		}
	}

	return nil
}

// StartServer 启动服务器
func (s *ServerService) StartServer(userID uint, serverID string) error {
	id, err := strconv.ParseUint(serverID, 10, 32)
	if err != nil {
		return fmt.Errorf("无效的服务器ID")
	}

	var server models.Server
	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&server).Error; err != nil {
		return fmt.Errorf("服务器不存在")
	}

	if server.Status == "running" {
		return fmt.Errorf("服务器已在运行中")
	}

	if server.Status == "starting" {
		return fmt.Errorf("服务器正在启动中")
	}

	// 更新服务器状态为启动中
	server.Status = "starting"
	if err := database.DB.Save(&server).Error; err != nil {
		return fmt.Errorf("更新服务器状态失败: %w", err)
	}

	// 启动Docker容器
	dockerManager, err := docker_manager.GetDockerManager()
	if err != nil {
		return fmt.Errorf("获取Docker管理器失败: %w", err)
	}

	containerName := utils.GetServerContainerName(server.ID)

	go func() {
		if err := s.startServerAsync(server, dockerManager, containerName); err != nil {
			fmt.Printf("Failed to start server: %v\n", err)
			database.DB.Model(&server).Update("status", "stopped")
		}
	}()

	return nil
}

// startServerAsync 异步启动服务器
func (s *ServerService) startServerAsync(server models.Server, dockerManager *docker_manager.DockerManager, containerName string) error {
	// 检查镜像是否存在
	imageName := "tbro98/ase-server:latest"
	exists, err := dockerManager.ImageExists(imageName)
	if err != nil {
		return fmt.Errorf("检查镜像是否存在失败: %w", err)
	}
	if !exists {
		return fmt.Errorf("镜像 %s 不存在", imageName)
	}

	// 检查容器是否存在
	containerExists, err := dockerManager.ContainerExists(containerName)
	if err != nil {
		return fmt.Errorf("检查容器是否存在失败: %w", err)
	}

	needRecreateContainer := false

	if containerExists {
		// 检查是否需要重建容器
		envVars, err := dockerManager.GetContainerEnvVars(containerName)
		if err != nil {
			needRecreateContainer = true
		} else {
			// 获取当前服务器的启动参数
			var serverArgs *models.ServerArgs
			if server.ServerArgsJSON != "" && server.ServerArgsJSON != "{}" {
				serverArgs = models.NewServerArgs()
				if err := json.Unmarshal([]byte(server.ServerArgsJSON), serverArgs); err != nil {
					serverArgs = models.FromServer(server)
				}
			} else {
				serverArgs = models.FromServer(server)
			}
			currentArgsString := serverArgs.GenerateArgsString(server)

			// 比较环境变量
			if containerArgsString, exists := envVars["SERVER_ARGS"]; exists {
				if containerArgsString != currentArgsString {
					needRecreateContainer = true
				}
			} else {
				needRecreateContainer = true
			}

			// 检查其他参数
			if !needRecreateContainer {
				if server.GameModIds != envVars["GameModIds"] {
					needRecreateContainer = true
				}
			}
		}

		if needRecreateContainer {
			if err := dockerManager.RemoveContainer(containerName); err != nil {
				return fmt.Errorf("删除现有容器失败: %w", err)
			}
		}
	}

	// 创建容器
	if !containerExists || needRecreateContainer {
		_, err = dockerManager.CreateContainer(server.ID, server.Identifier, server.Port, server.QueryPort, server.RCONPort, server.AdminPassword, server.Map, server.GameModIds, server.AutoRestart)
		if err != nil {
			return fmt.Errorf("创建容器失败: %w", err)
		}
	}

	// 启动容器
	if err := dockerManager.StartContainer(containerName); err != nil {
		return fmt.Errorf("启动容器失败: %w", err)
	}

	// 等待容器启动
	for i := 0; i < 30; i++ {
		time.Sleep(1 * time.Second)
		status, err := dockerManager.GetContainerStatus(containerName)
		if err != nil {
			continue
		}

		if status == "running" {
			if err := database.DB.Model(&server).Update("status", "running").Error; err != nil {
				fmt.Printf("Failed to update server status to running: %v\n", err)
			}
			return nil
		}
	}

	return fmt.Errorf("容器启动超时")
}

// StopServer 停止服务器
func (s *ServerService) StopServer(userID uint, serverID string) error {
	id, err := strconv.ParseUint(serverID, 10, 32)
	if err != nil {
		return fmt.Errorf("无效的服务器ID")
	}

	var server models.Server
	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&server).Error; err != nil {
		return fmt.Errorf("服务器不存在")
	}

	if server.Status == "stopped" {
		return fmt.Errorf("服务器已经停止")
	}

	if server.Status == "stopping" {
		return fmt.Errorf("服务器正在停止中")
	}

	// 更新服务器状态为停止中
	server.Status = "stopping"
	if err := database.DB.Save(&server).Error; err != nil {
		return fmt.Errorf("更新服务器状态失败: %w", err)
	}

	// 停止Docker容器
	dockerManager, err := docker_manager.GetDockerManager()
	if err != nil {
		return fmt.Errorf("获取Docker管理器失败: %w", err)
	}

	containerName := utils.GetServerContainerName(server.ID)

	go func() {
		s.stopServerAsync(server, dockerManager, containerName)
	}()

	return nil
}

// stopServerAsync 异步停止服务器
func (s *ServerService) stopServerAsync(server models.Server, dockerManager *docker_manager.DockerManager, containerName string) {
	// 检查容器是否存在
	containerExists, err := dockerManager.ContainerExists(containerName)
	if err != nil {
		fmt.Printf("Failed to check container existence: %v\n", err)
		database.DB.Model(&server).Update("status", "stopped")
		return
	}

	if !containerExists {
		database.DB.Model(&server).Update("status", "stopped")
		return
	}

	// 停止容器
	if err := dockerManager.StopContainer(containerName); err != nil {
		fmt.Printf("Failed to stop Docker container: %v\n", err)
	}

	// 验证容器状态
	for i := 0; i < 30; i++ {
		time.Sleep(1 * time.Second)
		status, err := dockerManager.GetContainerStatus(containerName)
		if err != nil {
			continue
		}

		if status == "stopped" || status == "not_found" {
			break
		}
	}

	// 更新状态为已停止
	if err := database.DB.Model(&server).Update("status", "stopped").Error; err != nil {
		fmt.Printf("Failed to update server status to stopped: %v\n", err)
	}
}

// GetServerFolderInfo 获取服务器文件夹信息
func (s *ServerService) GetServerFolderInfo(userID uint, serverID string) (map[string]interface{}, error) {
	id, err := strconv.ParseUint(serverID, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("无效的服务器ID")
	}

	var server models.Server
	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&server).Error; err != nil {
		return nil, fmt.Errorf("服务器不存在")
	}

	// 获取Docker卷信息
	dockerManager, err := docker_manager.GetDockerManager()
	if err != nil {
		return nil, fmt.Errorf("获取Docker管理器失败: %w", err)
	}

	volumeName := utils.GetServerVolumeName(server.ID)
	pluginsVolumeName := utils.GetServerPluginsVolumeName(server.ID)
	containerName := utils.GetServerContainerName(server.ID)

	// 检查卷是否存在
	volumeExists, _ := dockerManager.VolumeExists(volumeName)
	pluginsVolumeExists, _ := dockerManager.VolumeExists(pluginsVolumeName)

	// 检查容器是否存在
	containerExists, _ := dockerManager.ContainerExists(containerName)

	// 获取容器状态
	containerStatus := "not_found"
	if containerExists {
		if status, err := dockerManager.GetContainerStatus(containerName); err == nil {
			containerStatus = status
		}
	}

	return map[string]interface{}{
		"server_id":             server.ID,
		"server_name":           server.Identifier,
		"volume_name":           volumeName,
		"volume_exists":         volumeExists,
		"plugins_volume_name":   pluginsVolumeName,
		"plugins_volume_exists": pluginsVolumeExists,
		"container_name":        containerName,
		"container_exists":      containerExists,
		"container_status":      containerStatus,
	}, nil
}

// GetImageStatus 获取镜像状态
func (s *ServerService) GetImageStatus() (map[string]interface{}, error) {
	dockerManager, err := docker_manager.GetDockerManager()
	if err != nil {
		return nil, fmt.Errorf("获取Docker管理器失败: %w", err)
	}

	requiredImages := []string{
		"tbro98/ase-server:latest",
		"alpine:latest",
	}

	imageStatuses := make(map[string]*docker_manager.ImageStatus)
	allReady := true
	anyPulling := false
	pullingCount := 0

	for _, imageName := range requiredImages {
		status := dockerManager.GetImageStatus(imageName)
		imageStatuses[imageName] = status

		if !status.Ready {
			allReady = false
		}

		if status.Pulling {
			anyPulling = true
			pullingCount++
		}
	}

	// 生成总体状态描述
	var overallStatus string
	if allReady {
		overallStatus = "所有镜像已就绪"
	} else if anyPulling {
		overallStatus = "正在下载镜像"
	} else {
		overallStatus = "镜像未就绪，等待下载"
	}

	return map[string]interface{}{
		"images":            imageStatuses,
		"any_pulling":       anyPulling,
		"any_not_ready":     !allReady,
		"can_create_server": true,
		"can_start_server":  allReady,
		"overall_status":    overallStatus,
		"pulling_count":     pullingCount,
		"total_images":      len(requiredImages),
	}, nil
}
