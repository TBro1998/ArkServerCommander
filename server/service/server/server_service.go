package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"

	"ark-server-commander/database"
	"ark-server-commander/models"
	"ark-server-commander/service/docker_manager"
	"ark-server-commander/utils"

	"go.uber.org/zap"
)

// managedImages 由本系统管理的镜像，只允许对这些镜像执行拉取和更新
var managedImages = []string{
	utils.ArkServerLinuxImage,
	utils.ArkServerProtonImage,
}

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

	// 一次列出所有服务器容器，避免对每台服务器分别 inspect
	containers, listErr := dockerManager.ListServerContainers()
	if listErr != nil {
		utils.Warn("列出服务器容器失败，状态沿用数据库记录", zap.Error(listErr))
	}
	// 两个 ARK 镜像当前的 ID；镜像还没下载时没有记录，对应服务器不显示"旧镜像"
	currentImageIDs := map[string]string{}
	for _, name := range utils.ArkServerImages {
		if id, err := dockerManager.GetImageID(name); err == nil {
			currentImageIDs[name] = id
		}
	}

	var serverResponses []models.ServerResponse
	for _, server := range servers {
		realTimeStatus := server.Status
		imageOutdated := false

		if listErr == nil {
			dockerStatus := "not_found"
			if c, ok := containers[utils.GetServerContainerName(server.ID)]; ok {
				dockerStatus = c.Status
				// 容器用的不是服务器当前应使用的镜像（镜像更新过，或切换了 ArkApi）时，下次启动会重建
				if want, ok := currentImageIDs[utils.ArkServerImageFor(server.ArkApi)]; ok {
					imageOutdated = c.ImageID != want
				}
			}
			realTimeStatus = syncStatusWithDocker(server, dockerStatus)
		}

		response := toResponse(server)
		response.Status = realTimeStatus
		response.ImageOutdated = imageOutdated
		response.Runtime = runtimeFor(server.ID, realTimeStatus)
		serverResponses = append(serverResponses, response)
	}

	return serverResponses, nil
}

// CreateServer 创建新服务器
func (s *ServerService) CreateServer(userID uint, req models.ServerRequest) (*models.ServerResponse, error) {
	if err := checkPorts(0, req.Port, req.QueryPort, req.RCONPort); err != nil {
		return nil, err
	}

	clusterID, modIDs, err := normalizeServerFields(req.ClusterID, req.GameModIds)
	if err != nil {
		return nil, err
	}
	req.ClusterID, req.GameModIds = clusterID, modIDs

	// 设置默认值
	if req.SessionName == "" {
		req.SessionName = "ARK Server"
	}
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
		AutoUpdate:    req.AutoUpdate != nil && *req.AutoUpdate,
		ArkApi:        req.ArkApi != nil && *req.ArkApi,
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
	// AutoRestart 的默认值是 true，GORM 创建时会把 false 换成默认值，需要显式写回
	if !*req.AutoRestart {
		if err := tx.Model(&server).Update("auto_restart", false).Error; err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("服务器创建失败: %w", err)
		}
		server.AutoRestart = false
	}

	// 创建Docker卷
	dockerManager, err := docker_manager.GetDockerManager()
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("获取Docker管理器失败: %w", err)
	}

	if err := dockerManager.CreateVolume(server.ID); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("创建Docker卷失败: %w", err)
	}

	// 卷创建之后，只要事务没有提交成功，就清理已创建的卷，避免泄漏
	committed := false
	defer func() {
		if committed {
			return
		}
		if err := dockerManager.RemoveVolume(server.ID); err != nil {
			utils.Warn("创建服务器失败后清理Docker卷失败", zap.Uint("server_id", server.ID), zap.Error(err))
		}
		// 绑定目录不会随卷删除；事务回滚后 ID 可能被复用，残留文件会串到下一台服务器
		if err := docker_manager.RemoveServerDirs(server.ID); err != nil {
			utils.Warn("创建服务器失败后清理数据目录失败", zap.Uint("server_id", server.ID), zap.Error(err))
		}
	}()

	// 处理配置文件
	var gameUserSettings string
	var gameIni string

	if req.GameUserSettings != "" {
		if err = utils.ValidateINIContent(req.GameUserSettings); err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("GameUserSettings.ini格式错误: %w", err)
		}
		gameUserSettings = req.GameUserSettings
	} else {
		gameUserSettings = utils.GetDefaultGameUserSettings(server.SessionName, server.Map, 70)
	}

	if req.GameIni != "" {
		if err = utils.ValidateINIContent(req.GameIni); err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("game.ini格式错误: %w", err)
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
		return nil, fmt.Errorf("数据库提交失败: %w", err)
	}
	committed = true

	// 构建响应
	response := toResponse(server)
	response.ServerArgs = serverArgsFromJSON(server)

	// 读取配置文件内容
	fillConfigFiles(&response, server, dockerManager)

	return &response, nil
}

// portUse 记录一个宿主机端口及其用途，用于冲突提示
type portUse struct {
	port int
	name string
}

// serverPorts 返回服务器会占用的所有宿主机端口。
// ARK 除了游戏端口，还会使用游戏端口+1（原始 UDP 套接字），创建容器时这两个都会绑定
func serverPorts(port, queryPort, rconPort int) []portUse {
	return []portUse{
		{port, "游戏端口"},
		{port + 1, "游戏端口+1"},
		{queryPort, "查询端口"},
		{rconPort, "RCON端口"},
	}
}

// checkPorts 检查端口是否合法：本服务器内部不能重复，也不能和其他服务器占用的端口冲突。
// excludeID 是正在编辑的服务器（新建时传 0）。所有用户的服务器共用宿主机端口，所以不按用户过滤；
// 已删除的服务器容器已经移除，不占用端口（软删除记录会被 GORM 自动排除）。
// 注意：这里只是提前拦截，宿主机上其他程序占用的端口仍然要等到启动容器时才能发现
func checkPorts(excludeID uint, port, queryPort, rconPort int) error {
	mine := serverPorts(port, queryPort, rconPort)
	for i, a := range mine {
		if a.port < 1 || a.port > 65535 {
			return fmt.Errorf("%w: %s %d 超出有效范围（1-65535）", ErrInvalidPort, a.name, a.port)
		}
		for _, b := range mine[i+1:] {
			if a.port == b.port {
				return fmt.Errorf("%w: %s 与 %s 都是 %d", ErrPortConflict, a.name, b.name, a.port)
			}
		}
	}

	var others []models.Server
	query := database.DB.Select("id", "session_name", "port", "query_port", "rcon_port")
	if excludeID != 0 {
		query = query.Where("id <> ?", excludeID)
	}
	if err := query.Find(&others).Error; err != nil {
		return fmt.Errorf("检查端口冲突失败: %w", err)
	}

	for _, other := range others {
		for _, theirs := range serverPorts(other.Port, other.QueryPort, other.RCONPort) {
			for _, a := range mine {
				if a.port == theirs.port {
					return fmt.Errorf("%w: %s %d 已被服务器「%s」的%s占用", ErrPortConflict, a.name, a.port, other.SessionName, theirs.name)
				}
			}
		}
	}
	return nil
}

// serverArgsFromJSON 把数据库里的启动参数解析出来；没有保存过时返回默认值
func serverArgsFromJSON(server models.Server) *models.ServerArgs {
	if server.ServerArgsJSON == "" || server.ServerArgsJSON == "{}" {
		return models.FromServer(server)
	}
	args := models.NewServerArgs()
	if err := json.Unmarshal([]byte(server.ServerArgsJSON), args); err != nil {
		return models.FromServer(server)
	}
	return args
}

// fillConfigFiles 把配置文件内容填进响应；有待应用的修改时返回待应用的内容，
// 这样编辑页看到的是用户最后一次保存的内容，而不是卷里即将被替换的旧文件
func fillConfigFiles(response *models.ServerResponse, server models.Server, dockerManager *docker_manager.DockerManager) {
	response.ConfigPending = server.PendingGameUserSettings != "" || server.PendingGameIni != ""

	if server.PendingGameUserSettings != "" {
		response.GameUserSettings = server.PendingGameUserSettings
	} else if content, err := dockerManager.ReadConfigFile(server.ID, utils.GameUserSettingsFileName); err == nil {
		response.GameUserSettings = content
	}

	if server.PendingGameIni != "" {
		response.GameIni = server.PendingGameIni
	} else if content, err := dockerManager.ReadConfigFile(server.ID, utils.GameIniFileName); err == nil {
		response.GameIni = content
	}
}

// syncAdminPassword 把管理员密码写入 GameUserSettings.ini 的 [ServerSettings]。
// 密码不再放进启动参数和环境变量（docker inspect 能直接看到），ARK 会从这个文件读取 ServerAdminPassword，RCON 同样使用它。
// 每次启动前都以数据库里的值为准覆盖；此时容器已停止，ARK 停服时的写回不会再改动它
func syncAdminPassword(server models.Server, dockerManager *docker_manager.DockerManager) error {
	if strings.ContainsAny(server.AdminPassword, "\r\n") {
		return fmt.Errorf("管理员密码不能包含换行符")
	}

	current, err := dockerManager.ReadConfigFile(server.ID, utils.GameUserSettingsFileName)
	if err != nil && !errors.Is(err, docker_manager.ErrConfigNotFound) {
		// 读取失败时不能当作文件为空，否则会用只含密码的内容覆盖掉整个配置文件
		return fmt.Errorf("读取GameUserSettings.ini失败: %w", err)
	}

	updated := utils.SetIniValue(current, "ServerSettings", "ServerAdminPassword", server.AdminPassword)
	if strings.TrimSpace(updated) == strings.TrimSpace(current) {
		return nil
	}
	if err := dockerManager.WriteConfigFile(server.ID, utils.GameUserSettingsFileName, updated); err != nil {
		return fmt.Errorf("写入管理员密码失败: %w", err)
	}
	return nil
}

// applyPendingConfig 在启动容器前把暂存的配置文件写入数据卷。
// 此时容器已停止，ARK 停服时的写回已经完成，写入的内容不会再被覆盖
func applyPendingConfig(serverID uint, dockerManager *docker_manager.DockerManager) error {
	// 重新读取，拿到启动请求之后才保存的修改
	var server models.Server
	if err := database.DB.Select("id", "pending_game_user_settings", "pending_game_ini").First(&server, serverID).Error; err != nil {
		return fmt.Errorf("读取待应用的配置失败: %w", err)
	}

	files := []struct {
		column, fileName, content string
	}{
		{"pending_game_user_settings", utils.GameUserSettingsFileName, server.PendingGameUserSettings},
		{"pending_game_ini", utils.GameIniFileName, server.PendingGameIni},
	}
	for _, f := range files {
		if f.content == "" {
			continue
		}
		if err := dockerManager.WriteConfigFile(serverID, f.fileName, f.content); err != nil {
			return fmt.Errorf("写入%s失败: %w", f.fileName, err)
		}
		// 只清除刚写入的内容；写入期间用户又保存了新内容时保留，留到下次启动
		database.DB.Model(&models.Server{}).
			Where("id = ? AND "+f.column+" = ?", serverID, f.content).
			Update(f.column, "")
		utils.Info("已应用待生效的配置文件", zap.Uint("server_id", serverID), zap.String("file", f.fileName))
	}
	return nil
}

// GetServer 获取单个服务器信息
func (s *ServerService) GetServer(userID uint, serverID string) (*models.ServerResponse, error) {
	id, err := strconv.ParseUint(serverID, 10, 32)
	if err != nil {
		return nil, ErrInvalidServerID
	}

	var server models.Server
	if err = database.DB.Where("id = ? AND user_id = ?", id, userID).First(&server).Error; err != nil {
		return nil, ErrServerNotFound
	}

	// 解析启动参数
	serverArgs := serverArgsFromJSON(server)

	response := toResponse(server)
	response.ServerArgs = serverArgs
	response.GeneratedArgs = serverArgs.GenerateArgsString(server)

	// 读取配置文件内容
	dockerManager, err := docker_manager.GetDockerManager()
	if err != nil {
		return nil, fmt.Errorf("获取Docker管理器失败: %w", err)
	}

	fillConfigFiles(&response, server, dockerManager)

	return &response, nil
}

// GetServerRCON 获取服务器RCON连接信息
func (s *ServerService) GetServerRCON(userID uint, serverID string) (map[string]interface{}, error) {
	id, err := strconv.ParseUint(serverID, 10, 32)
	if err != nil {
		return nil, ErrInvalidServerID
	}

	var server models.Server
	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&server).Error; err != nil {
		return nil, ErrServerNotFound
	}

	return map[string]interface{}{
		"server_id":      server.ID,
		"rcon_port":      server.RCONPort,
		"admin_password": server.AdminPassword,
	}, nil
}

// UpdateServer 更新服务器配置
func (s *ServerService) UpdateServer(userID uint, serverID string, req models.ServerUpdateRequest) (*models.ServerResponse, bool, error) {
	id, err := strconv.ParseUint(serverID, 10, 32)
	if err != nil {
		return nil, false, ErrInvalidServerID
	}

	var server models.Server
	if err = database.DB.Where("id = ? AND user_id = ?", id, userID).First(&server).Error; err != nil {
		return nil, false, ErrServerNotFound
	}

	// 更新字段：字符串字段为空表示不修改；集群ID和模组ID用指针，传空字符串表示清空
	oldPorts := [3]int{server.Port, server.QueryPort, server.RCONPort}
	if req.SessionName != "" {
		server.SessionName = req.SessionName
	}
	if req.ClusterID != nil {
		server.ClusterID = *req.ClusterID
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
	if req.GameModIds != nil {
		server.GameModIds = *req.GameModIds
	}
	if req.AutoRestart != nil {
		server.AutoRestart = *req.AutoRestart
	}
	if req.AutoUpdate != nil {
		server.AutoUpdate = *req.AutoUpdate
	}
	// 切换 ArkApi 即切换镜像，下次启动时发现容器的镜像不对会重建
	if req.ArkApi != nil {
		server.ArkApi = *req.ArkApi
	}
	clusterID, modIDs, err := normalizeServerFields(server.ClusterID, server.GameModIds)
	if err != nil {
		return nil, false, err
	}
	server.ClusterID, server.GameModIds = clusterID, modIDs

	// 端口有变化时才检查，避免历史数据里已有的冲突让无关字段也保存不了
	if [3]int{server.Port, server.QueryPort, server.RCONPort} != oldPorts {
		if err := checkPorts(server.ID, server.Port, server.QueryPort, server.RCONPort); err != nil {
			return nil, false, err
		}
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

	// 先校验配置文件，避免数据库已经更新、写文件时才报错，造成部分更新
	writingConfig := req.GameUserSettings != "" || req.GameIni != ""
	if req.GameUserSettings != "" {
		if err := utils.ValidateINIContent(req.GameUserSettings); err != nil {
			return nil, false, fmt.Errorf("GameUserSettings.ini格式错误: %w", err)
		}
	}
	if req.GameIni != "" {
		if err := utils.ValidateINIContent(req.GameIni); err != nil {
			return nil, false, fmt.Errorf("Game.ini格式错误: %w", err)
		}
	}

	// 服务器未停止时，ARK 会在停服时把内存里的配置写回 GameUserSettings.ini，直接写文件会被覆盖。
	// 这时先把内容暂存到数据库，下次启动前（容器已停止、写回已完成）再写入数据卷
	deferConfig := writingConfig && server.Status != "stopped"
	if req.GameUserSettings != "" {
		// 直接写入时清掉旧的暂存内容，否则下次启动会用旧内容覆盖这次的修改
		server.PendingGameUserSettings = ""
		if deferConfig {
			server.PendingGameUserSettings = req.GameUserSettings
		}
	}
	if req.GameIni != "" {
		server.PendingGameIni = ""
		if deferConfig {
			server.PendingGameIni = req.GameIni
		}
	}

	// 保存时不写 status / last_error：读取之后启停协程可能已经改了它们，整行写回会把状态覆盖成旧值
	if err := database.DB.Omit("Status", "LastError").Save(&server).Error; err != nil {
		return nil, false, fmt.Errorf("服务器更新失败: %w", err)
	}
	// 重新读取，响应里的状态要和数据库一致
	if err := database.DB.First(&server, server.ID).Error; err != nil {
		return nil, false, fmt.Errorf("读取服务器失败: %w", err)
	}

	if writingConfig && !deferConfig {
		dockerManager, err := docker_manager.GetDockerManager()
		if err != nil {
			return nil, false, fmt.Errorf("获取Docker管理器失败: %w", err)
		}

		if req.GameUserSettings != "" {
			if err := dockerManager.WriteConfigFile(uint(id), utils.GameUserSettingsFileName, req.GameUserSettings); err != nil {
				return nil, false, fmt.Errorf("写入GameUserSettings.ini失败: %w", err)
			}
		}
		if req.GameIni != "" {
			if err := dockerManager.WriteConfigFile(uint(id), utils.GameIniFileName, req.GameIni); err != nil {
				return nil, false, fmt.Errorf("写入Game.ini失败: %w", err)
			}
		}
	}

	// 构建响应
	response := toResponse(server)
	response.ServerArgs = serverArgsFromJSON(server)

	// 读取配置文件内容
	dockerManager, err := docker_manager.GetDockerManager()
	if err != nil {
		return nil, false, fmt.Errorf("获取Docker管理器失败: %w", err)
	}

	fillConfigFiles(&response, server, dockerManager)

	return &response, argsChanged, nil
}

// DeleteServer 删除服务器：只允许删除已停止的服务器。
// 先删容器再软删数据库记录：容器删除失败时记录仍在，用户可以重试，不会留下没人管理的孤儿容器。
// deleteData 为 false 时数据卷和备份保留；为 true 时删除服务器的三个卷和全部备份（集群卷是共享的，不删除）
func (s *ServerService) DeleteServer(userID uint, serverID string, deleteData bool) error {
	id, err := strconv.ParseUint(serverID, 10, 32)
	if err != nil {
		return ErrInvalidServerID
	}

	var server models.Server
	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&server).Error; err != nil {
		return ErrServerNotFound
	}

	switch server.Status {
	case "stopped":
	case "running":
		return ErrServerRunning
	default:
		return ErrServerBusy
	}

	dockerManager, err := docker_manager.GetDockerManager()
	if err != nil {
		return fmt.Errorf("获取Docker管理器失败: %w", err)
	}

	containerName := utils.GetServerContainerName(server.ID)
	exists, err := dockerManager.ContainerExists(containerName)
	if err != nil {
		return fmt.Errorf("检查Docker容器失败: %w", err)
	}
	if exists {
		if err := dockerManager.RemoveContainer(containerName); err != nil {
			return fmt.Errorf("删除Docker容器失败: %w", err)
		}
	}

	// 只在状态仍是 stopped 时删除，避免删除期间有人启动了服务器
	res := database.DB.Where("status = ?", "stopped").Delete(&models.Server{}, server.ID)
	if res.Error != nil {
		return fmt.Errorf("服务器删除失败: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return ErrServerBusy
	}

	// 删除服务器关联的定时任务
	database.DB.Where("server_id = ?", server.ID).Delete(&models.Schedule{})
	runtimeMonitor.Forget(server.ID)

	if deleteData {
		var errs []error
		if err := dockerManager.RemoveVolume(server.ID); err != nil {
			errs = append(errs, err)
		}
		// 卷绑定在数据目录 servers/<id>/ 下，删除卷不会删除文件，需要单独清理
		if err := docker_manager.RemoveServerDirs(server.ID); err != nil {
			errs = append(errs, fmt.Errorf("删除数据目录失败: %w", err))
		}
		if err := os.RemoveAll(serverBackupDir(server.ID)); err != nil {
			errs = append(errs, fmt.Errorf("删除备份失败: %w", err))
		}
		if len(errs) > 0 {
			// 记录已经删除，这里只能提示用户手动清理
			return fmt.Errorf("服务器已删除，但清理数据失败，请手动删除数据目录 servers/%d 和对应的 Docker 卷: %w",
				server.ID,
				errors.Join(errs...))
		}
		utils.Info("已删除服务器数据", zap.Uint("server_id", server.ID))
	}
	return nil
}

// transitionStatus 原子地把服务器状态切换为 to，仅当当前状态不在 blocked 中时才会成功。
// 切换成功时清空上一次的错误信息。
// 返回 false 表示当前状态不允许切换，或已被并发请求抢先切换。
func transitionStatus(serverID uint, blocked []string, to string) (bool, error) {
	res := database.DB.Model(&models.Server{}).
		Where("id = ? AND status NOT IN ?", serverID, blocked).
		Updates(map[string]any{"status": to, "last_error": ""})
	return res.RowsAffected > 0, res.Error
}

// finishOperation 启停操作结束时写入最终状态和错误信息（errMsg 为空表示成功）
func finishOperation(serverID uint, status, errMsg string) {
	err := database.DB.Model(&models.Server{}).Where("id = ?", serverID).
		Updates(map[string]any{"status": status, "last_error": errMsg}).Error
	if err != nil {
		utils.Error("更新服务器状态失败", zap.Uint("server_id", serverID), zap.String("status", status), zap.Error(err))
	}
}

// syncStatusWithDocker 用容器的实际状态校正数据库中的状态，返回应展示给用户的状态。
// starting / stopping 表示后台操作还在进行，此时容器状态处于过渡中，不做校正，
// 否则会把"启动中"改回"已停止"，让并发保护失效。
// 只同步 running / stopped 两种确定状态；用 status 作为条件更新，避免覆盖并发写入的新状态。
func syncStatusWithDocker(server models.Server, dockerStatus string) string {
	if server.Status == "starting" || server.Status == "stopping" {
		return server.Status
	}

	actual := dockerStatus
	if actual == "not_found" {
		actual = "stopped"
	}
	if actual != "running" && actual != "stopped" {
		// restarting、paused 等状态只展示，不写入数据库
		return dockerStatus
	}

	if actual != server.Status {
		err := database.DB.Model(&models.Server{}).
			Where("id = ? AND status = ?", server.ID, server.Status).
			Update("status", actual).Error
		if err != nil {
			utils.Warn("同步服务器状态失败", zap.Uint("server_id", server.ID), zap.Error(err))
		}
	}
	return actual
}

// ReconcileStatuses 后端启动时调用：上次退出时还处于 starting / stopping 的服务器，
// 对应的后台协程已经不存在了，按容器实际状态改为 running 或 stopped
func (s *ServerService) ReconcileStatuses() error {
	var servers []models.Server
	if err := database.DB.Where("status IN ?", []string{"starting", "stopping"}).Find(&servers).Error; err != nil {
		return fmt.Errorf("获取服务器列表失败: %w", err)
	}
	if len(servers) == 0 {
		return nil
	}

	dockerManager, err := docker_manager.GetDockerManager()
	if err != nil {
		return fmt.Errorf("获取Docker管理器失败: %w", err)
	}

	for _, server := range servers {
		status := "stopped"
		if dockerStatus, err := dockerManager.GetContainerStatus(utils.GetServerContainerName(server.ID)); err == nil && dockerStatus == "running" {
			status = "running"
		}
		finishOperation(server.ID, status, "后端重启时操作被中断")
		utils.Info("已校正中断操作的服务器状态", zap.Uint("server_id", server.ID), zap.String("status", status))
	}
	return nil
}

// StartServer 启动服务器
func (s *ServerService) StartServer(userID uint, serverID string) error {
	id, err := strconv.ParseUint(serverID, 10, 32)
	if err != nil {
		return ErrInvalidServerID
	}

	var server models.Server
	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&server).Error; err != nil {
		return ErrServerNotFound
	}

	if server.Status == "running" {
		return ErrAlreadyRunning
	}

	if server.Status == "starting" {
		return ErrAlreadyStarting
	}

	// 原子地切换为启动中，避免并发请求重复启动
	ok, err := transitionStatus(server.ID, []string{"running", "starting"}, "starting")
	if err != nil {
		return fmt.Errorf("更新服务器状态失败: %w", err)
	}
	if !ok {
		// 被并发请求抢先，重新读取状态以返回准确的错误
		var latest models.Server
		if err := database.DB.Select("status").First(&latest, server.ID).Error; err == nil && latest.Status == "running" {
			return ErrAlreadyRunning
		}
		return ErrAlreadyStarting
	}
	server.Status = "starting"

	// 启动Docker容器
	dockerManager, err := docker_manager.GetDockerManager()
	if err != nil {
		return fmt.Errorf("获取Docker管理器失败: %w", err)
	}

	containerName := utils.GetServerContainerName(server.ID)

	go s.runStart(server, dockerManager, containerName)

	return nil
}

// startServerAsync 异步启动服务器
func (s *ServerService) startServerAsync(server models.Server, dockerManager *docker_manager.DockerManager, containerName string) error {
	// ARK 镜像按是否启用 ArkApi 选择，没有时自动下载
	imageName := utils.ArkServerImageFor(server.ArkApi)
	if err := ensureServerImage(server.ID, imageName, dockerManager); err != nil {
		return err
	}

	if err := applyPendingConfig(server.ID, dockerManager); err != nil {
		return err
	}

	if err := syncAdminPassword(server, dockerManager); err != nil {
		return err
	}

	// 还没下载过服务端时，从另一台已安装的服务器复制，省去首次下载
	seedServerFiles(server, dockerManager)

	// 启用了 ArkApi 但文件卷里还没有（新服务器，或刚打开开关）：安装版本库里最新的版本。
	// 放在复制服务端文件之后，复制来源带有 ArkApi 时沿用它的版本
	if server.ArkApi {
		if err := ensureArkApiInstalled(server.ID); err != nil {
			return err
		}
	}

	// 开启了自动更新：清除安装标记，入口脚本会重新校验游戏本体并重新下载 Mod
	if server.AutoUpdate {
		if err := dockerManager.RemoveInstallMarkers(server.ID, true, server.GameModIds != ""); err != nil {
			return err
		}
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
			currentArgsString := serverArgsFromJSON(server).GenerateArgsString(server)

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

			// 集群卷挂载与当前 cluster_id 不一致（包括升级前创建、没有挂载集群卷的旧容器）
			if !needRecreateContainer && envVars["CLUSTER_VOLUME"] != docker_manager.ExpectedClusterVolume(server.ClusterID) {
				needRecreateContainer = true
			}

			// 自动重启开关改过，重启策略需要重建容器才能生效
			if !needRecreateContainer {
				wantPolicy := "unless-stopped"
				if !server.AutoRestart {
					wantPolicy = "no"
				}
				if policy, err := dockerManager.GetContainerRestartPolicy(containerName); err != nil || policy != wantPolicy {
					needRecreateContainer = true
				}
			}

			// 镜像更新后，旧容器仍指向旧镜像，需要重建才能使用新镜像；切换了 ArkApi（换了镜像）时同样重建
			if !needRecreateContainer {
				outdated, err := dockerManager.IsContainerImageOutdated(containerName, imageName)
				if err != nil || outdated {
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
		_, err = dockerManager.CreateContainer(server.ID, server.Port, server.QueryPort, server.RCONPort, server.AdminPassword, server.Map, server.GameModIds, server.AutoRestart)
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
			return nil
		}
	}

	// 没能进入运行状态（例如反复崩溃重启），停掉容器，避免重启策略让它在"已停止"的记录下继续运行
	if err := dockerManager.StopContainer(containerName); err != nil {
		utils.Warn("停止启动超时的容器失败", zap.String("container", containerName), zap.Error(err))
	}
	return fmt.Errorf("容器启动超时，请查看容器日志")
}

// runStart 执行启动并记录最终状态。失败时如果容器仍在运行，记录为 running，保证状态与容器一致
func (s *ServerService) runStart(server models.Server, dockerManager *docker_manager.DockerManager, containerName string) error {
	err := s.startServerAsync(server, dockerManager, containerName)
	if err == nil {
		finishOperation(server.ID, "running", "")
		return nil
	}

	utils.Error("启动服务器失败", zap.Uint("server_id", server.ID), zap.Error(err))
	status := "stopped"
	if dockerStatus, statusErr := dockerManager.GetContainerStatus(containerName); statusErr == nil && dockerStatus == "running" {
		status = "running"
	}
	finishOperation(server.ID, status, err.Error())
	return err
}

// StopServer 停止服务器
func (s *ServerService) StopServer(userID uint, serverID string) error {
	id, err := strconv.ParseUint(serverID, 10, 32)
	if err != nil {
		return ErrInvalidServerID
	}

	var server models.Server
	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&server).Error; err != nil {
		return ErrServerNotFound
	}

	if server.Status == "stopped" {
		return ErrAlreadyStopped
	}

	if server.Status == "stopping" {
		return ErrAlreadyStopping
	}

	// 原子地切换为停止中，避免并发请求重复停止
	ok, err := transitionStatus(server.ID, []string{"stopped", "stopping"}, "stopping")
	if err != nil {
		return fmt.Errorf("更新服务器状态失败: %w", err)
	}
	if !ok {
		// 被并发请求抢先，重新读取状态以返回准确的错误
		var latest models.Server
		if err := database.DB.Select("status").First(&latest, server.ID).Error; err == nil && latest.Status == "stopped" {
			return ErrAlreadyStopped
		}
		return ErrAlreadyStopping
	}
	server.Status = "stopping"

	// 停止Docker容器
	dockerManager, err := docker_manager.GetDockerManager()
	if err != nil {
		return fmt.Errorf("获取Docker管理器失败: %w", err)
	}

	containerName := utils.GetServerContainerName(server.ID)

	go s.stopServerAsync(server, dockerManager, containerName)

	return nil
}

// stopServerAsync 异步停止服务器，结束时按容器实际状态写入 stopped 或 running（停止失败时）
// 返回 true 表示容器已停止
func (s *ServerService) stopServerAsync(server models.Server, dockerManager *docker_manager.DockerManager, containerName string) bool {
	if status, err := dockerManager.GetContainerStatus(containerName); err == nil && (status == "stopped" || status == "not_found") {
		finishOperation(server.ID, "stopped", "")
		return true
	}

	// 停服前先通过 RCON 保存世界：容器的主进程收不到 SIGTERM，docker stop 超时后会被直接 SIGKILL，
	// 不先保存的话最近一次自动存档之后的进度会丢。保存失败不阻止停止，只记录日志
	if status, err := dockerManager.GetContainerStatus(containerName); err == nil && status == "running" {
		if out, err := dockerManager.RconCommand(containerName, server.RCONPort, server.AdminPassword, "SaveWorld", time.Minute); err != nil {
			utils.Warn("停服前保存世界失败，继续停止容器", zap.Uint("server_id", server.ID), zap.Error(err))
		} else {
			utils.Info("停服前已保存世界", zap.Uint("server_id", server.ID), zap.String("response", out))
		}
	}

	stopErr := dockerManager.StopContainer(containerName)
	if stopErr != nil {
		utils.Error("停止Docker容器失败", zap.Uint("server_id", server.ID), zap.Error(stopErr))
	}

	// 确认容器已经停下；StopContainer 返回时通常已停止，这里兼容 Docker 返回较慢的情况
	var status string
	var statusErr error
	for i := 0; i < 30; i++ {
		status, statusErr = dockerManager.GetContainerStatus(containerName)
		if statusErr == nil && (status == "stopped" || status == "not_found") {
			finishOperation(server.ID, "stopped", "")
			return true
		}
		time.Sleep(1 * time.Second)
	}

	errMsg := "停止超时，容器仍在运行"
	switch {
	case stopErr != nil:
		errMsg = stopErr.Error()
	case statusErr != nil:
		errMsg = statusErr.Error()
	}

	// 无法确认已停止时，记录为容器当前的确定状态，而不是假装已停止
	final := "running"
	if statusErr != nil {
		final = "stopped"
	}
	finishOperation(server.ID, final, errMsg)
	return false
}

// CheckImageUpdates 检查所有管理的镜像更新
func (s *ServerService) CheckImageUpdates() (map[string]bool, error) {
	dockerManager, err := docker_manager.GetDockerManager()
	if err != nil {
		return nil, fmt.Errorf("获取Docker管理器失败: %w", err)
	}

	updateStatus := make(map[string]bool)
	var lastErr error
	for _, imageName := range managedImages {
		// 还没下载的镜像（例如没有服务器启用 ArkApi 时的 Proton 镜像）不算"有更新"，界面上显示的是下载按钮
		if exists, err := dockerManager.ImageExists(imageName); err == nil && !exists {
			continue
		}
		hasUpdate, err := dockerManager.CheckImageUpdate(imageName)
		if err != nil {
			utils.Warn("检查镜像更新失败", zap.String("image", imageName), zap.Error(err))
			lastErr = err
			continue
		}
		updateStatus[imageName] = hasUpdate
	}

	// 全部检查失败通常是无法访问镜像仓库，需要告知用户
	if len(updateStatus) == 0 && lastErr != nil {
		return nil, fmt.Errorf("检查镜像更新失败: %w", lastErr)
	}

	return updateStatus, nil
}

// PullImage 手动拉取指定镜像
func (s *ServerService) PullImage(imageName string) error {
	if !slices.Contains(managedImages, imageName) {
		return fmt.Errorf("%w: %s", ErrImageNotAllowed, imageName)
	}

	dockerManager, err := docker_manager.GetDockerManager()
	if err != nil {
		return fmt.Errorf("获取Docker管理器失败: %w", err)
	}

	// 同步占用拉取权：并发的请求只有一个能开始拉取，其余得到“正在下载”的提示
	if !docker_manager.ClaimPull(imageName) {
		return ErrImagePulling
	}

	go func() {
		if err := dockerManager.PullClaimedImage(imageName); err != nil {
			utils.Error("拉取镜像失败", zap.String("image", imageName), zap.Error(err))
		} else {
			utils.Info("镜像拉取完成", zap.String("image", imageName))
		}
	}()

	return nil
}

// UpdateImage 更新指定镜像
// 拉取完成后，已停止的服务器会在下次启动时自动用新镜像重建容器；
// restartRunning 为 true 时，正在运行且使用旧镜像的服务器会立即重启以应用新镜像
func (s *ServerService) UpdateImage(imageName string, userID uint, restartRunning bool) ([]models.ServerResponse, error) {
	if !slices.Contains(managedImages, imageName) {
		return nil, fmt.Errorf("%w: %s", ErrImageNotAllowed, imageName)
	}

	// 获取受影响的服务器
	affectedServers, err := s.GetAffectedServers(imageName, userID)
	if err != nil {
		return nil, fmt.Errorf("获取受影响服务器失败: %w", err)
	}

	dockerManager, err := docker_manager.GetDockerManager()
	if err != nil {
		return nil, fmt.Errorf("获取Docker管理器失败: %w", err)
	}

	// 同步占用拉取权，避免检查和拉取之间被并发请求抢先
	if !docker_manager.ClaimPull(imageName) {
		return nil, ErrImagePulling
	}

	// 异步更新镜像
	go func() {
		utils.Info("开始更新镜像", zap.String("image", imageName))
		if err := dockerManager.PullClaimedImage(imageName); err != nil {
			utils.Error("更新镜像失败", zap.String("image", imageName), zap.Error(err))
			return
		}
		utils.Info("镜像更新完成", zap.String("image", imageName))

		if !restartRunning || !slices.Contains(utils.ArkServerImages, imageName) {
			return
		}

		// 只重启使用这个镜像的服务器
		var running []models.Server
		if err := database.DB.Where("user_id = ? AND status = ? AND arkapi = ?", userID, "running", imageName == utils.ArkServerProtonImage).Find(&running).Error; err != nil {
			utils.Error("获取运行中的服务器失败", zap.Error(err))
			return
		}

		for _, server := range running {
			containerName := utils.GetServerContainerName(server.ID)
			outdated, err := dockerManager.IsContainerImageOutdated(containerName, imageName)
			if err != nil {
				utils.Warn("检查容器镜像失败", zap.String("container", containerName), zap.Error(err))
				continue
			}
			if outdated {
				go s.restartServer(server, dockerManager)
			}
		}
	}()

	return affectedServers, nil
}

// restartServer 停止并重新启动服务器，启动时会按当前镜像和配置重建容器
func (s *ServerService) restartServer(server models.Server, dockerManager *docker_manager.DockerManager) {
	if err := s.restartWith(server, dockerManager, nil); err == nil {
		utils.Info("服务器已重启", zap.Uint("server_id", server.ID))
	}
}

// restartWith 停止（如果在运行）并重新启动服务器；beforeStart 在容器停止之后、启动之前执行（可为 nil），
// 失败时不再启动，错误写入 last_error。整个过程同步执行，调用方需要时自行放进协程
func (s *ServerService) restartWith(server models.Server, dockerManager *docker_manager.DockerManager, beforeStart func() error) error {
	containerName := utils.GetServerContainerName(server.ID)

	var latest models.Server
	if err := database.DB.Select("status").First(&latest, server.ID).Error; err != nil {
		return err
	}

	if latest.Status != "stopped" {
		ok, err := transitionStatus(server.ID, []string{"stopped", "stopping", "starting"}, "stopping")
		if err != nil || !ok {
			utils.Warn("服务器状态已变化，跳过重启", zap.Uint("server_id", server.ID), zap.Error(err))
			return ErrServerBusy
		}
		server.Status = "stopping"
		if !s.stopServerAsync(server, dockerManager, containerName) {
			utils.Warn("停止服务器失败，跳过启动", zap.Uint("server_id", server.ID))
			return fmt.Errorf("停止服务器失败")
		}
	}

	ok, err := transitionStatus(server.ID, []string{"running", "starting", "stopping"}, "starting")
	if err != nil || !ok {
		utils.Warn("服务器状态已变化，跳过启动", zap.Uint("server_id", server.ID), zap.Error(err))
		return ErrServerBusy
	}
	server.Status = "starting"

	if beforeStart != nil {
		if err := beforeStart(); err != nil {
			finishOperation(server.ID, "stopped", err.Error())
			return err
		}
	}

	// 重新读取，使用停止期间可能已经保存的最新配置
	if err := database.DB.First(&server, server.ID).Error; err != nil {
		finishOperation(server.ID, "stopped", "读取服务器配置失败: "+err.Error())
		return err
	}
	return s.runStart(server, dockerManager, containerName)
}

// GetAffectedServers 获取使用指定镜像的服务器列表
func (s *ServerService) GetAffectedServers(imageName string, userID uint) ([]models.ServerResponse, error) {
	affected := []models.ServerResponse{}
	if !slices.Contains(utils.ArkServerImages, imageName) {
		return affected, nil
	}
	servers, err := s.GetServers(userID)
	if err != nil {
		return nil, err
	}
	for _, server := range servers {
		if server.Image == imageName {
			affected = append(affected, server)
		}
	}
	return affected, nil
}

// RecreateContainer 按当前配置重建服务器容器。
// 运行中的服务器会先停止，停止失败则放弃重建；重建在后台完成，期间状态为 stopping，
// 结束后服务器处于已停止状态，需要手动启动。失败原因记录在 last_error
func (s *ServerService) RecreateContainer(userID uint, serverID string) error {
	id, err := strconv.ParseUint(serverID, 10, 32)
	if err != nil {
		return ErrInvalidServerID
	}

	var server models.Server
	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&server).Error; err != nil {
		return ErrServerNotFound
	}

	if server.Status == "starting" || server.Status == "stopping" {
		return ErrServerBusy
	}

	dockerManager, err := docker_manager.GetDockerManager()
	if err != nil {
		return fmt.Errorf("获取Docker管理器失败: %w", err)
	}

	wasRunning := server.Status == "running"

	// 原子占位，防止重建期间被启动、停止或再次重建
	ok, err := transitionStatus(server.ID, []string{"starting", "stopping"}, "stopping")
	if err != nil {
		return fmt.Errorf("更新服务器状态失败: %w", err)
	}
	if !ok {
		return ErrServerBusy
	}
	server.Status = "stopping"

	go s.runRecreate(server, dockerManager, wasRunning)
	return nil
}

// runRecreate 在后台完成停止（如果在运行）和重建，结束时状态为 stopped，失败原因写入 last_error
func (s *ServerService) runRecreate(server models.Server, dockerManager *docker_manager.DockerManager, wasRunning bool) {
	containerName := utils.GetServerContainerName(server.ID)

	if wasRunning {
		if !s.stopServerAsync(server, dockerManager, containerName) {
			// stopServerAsync 已经记录了失败原因和容器的实际状态
			return
		}
		// 停止完成后状态已经是 stopped，重新占位，防止重建期间被启动
		ok, err := transitionStatus(server.ID, []string{"starting", "stopping"}, "stopping")
		if err != nil || !ok {
			utils.Warn("服务器状态已变化，放弃重建容器", zap.Uint("server_id", server.ID), zap.Error(err))
			return
		}
	}

	// 重新读取，使用停止期间可能已经保存的最新配置
	if err := database.DB.First(&server, server.ID).Error; err != nil {
		finishOperation(server.ID, "stopped", "读取服务器配置失败: "+err.Error())
		return
	}

	// 容器已停止，把暂存的配置文件写入数据卷
	if err := applyPendingConfig(server.ID, dockerManager); err != nil {
		utils.Error("重建容器前应用配置失败", zap.Uint("server_id", server.ID), zap.Error(err))
		finishOperation(server.ID, "stopped", err.Error())
		return
	}

	// CreateContainer 发现同名容器时会先删除
	_, err := dockerManager.CreateContainer(server.ID, server.Port, server.QueryPort, server.RCONPort, server.AdminPassword, server.Map, server.GameModIds, server.AutoRestart)
	if err != nil {
		utils.Error("重建容器失败", zap.Uint("server_id", server.ID), zap.Error(err))
		finishOperation(server.ID, "stopped", "重建容器失败: "+err.Error())
		return
	}

	finishOperation(server.ID, "stopped", "")
	utils.Info("服务器容器重建完成", zap.Uint("server_id", server.ID))
}

// GetImageStatus 获取镜像状态
func (s *ServerService) GetImageStatus() (map[string]interface{}, error) {
	dockerManager, err := docker_manager.GetDockerManager()
	if err != nil {
		return nil, fmt.Errorf("获取Docker管理器失败: %w", err)
	}

	required, err := requiredImages()
	if err != nil {
		return nil, err
	}

	imageStatuses := make(map[string]*docker_manager.ImageStatus)
	allReady := true
	anyPulling := false
	pullingCount := 0

	for _, imageName := range managedImages {
		status := dockerManager.GetImageStatus(imageName)
		status.Required = required[imageName]
		imageStatuses[imageName] = status

		// 没有服务器用到的 ARK 镜像不算未就绪，启动用到它的服务器时会自动下载
		if !status.Ready && status.Required {
			allReady = false
		}

		if status.Pulling {
			anyPulling = true
			pullingCount++
		}
	}

	// 总体状态用状态码（ready / pulling / not_ready），文案由前端翻译
	var overallStatus string
	if allReady {
		overallStatus = "ready"
	} else if anyPulling {
		overallStatus = "pulling"
	} else {
		overallStatus = "not_ready"
	}

	// 配置文件和存档由后端直接读写，ARK 镜像在启动时自动下载，创建和启动都不依赖镜像；
	// can_create_server / can_start_server 保留给前端，恒为 true
	return map[string]interface{}{
		"images":            imageStatuses,
		"any_pulling":       anyPulling,
		"any_not_ready":     !allReady,
		"can_create_server": true,
		"can_start_server":  true,
		"overall_status":    overallStatus,
		"pulling_count":     pullingCount,
		"total_images":      len(managedImages),
	}, nil
}
