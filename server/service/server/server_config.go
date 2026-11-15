package server

import (
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"time"

	"ark-server-commander/database"
	"ark-server-commander/models"
	"ark-server-commander/service/docker_manager"
	"ark-server-commander/utils"
)

// identityIniKeys 每台服务器自己的值：导出时去掉管理员密码，导入或复制配置文件时保留目标服务器原来的这些键。
// 游戏以启动参数里的名称、端口和人数为准，配置文件里的这些键也保持和服务器一致
var identityIniKeys = []string{"SessionName", "ServerAdminPassword", "Port", "QueryPort", "RCONPort", "MaxPlayers"}

// readCurrentConfig 读取服务器当前的配置文件：有暂存内容时取暂存内容，没有文件时为空。
// 读取失败时报错，不能当作空文件，否则复制时会丢掉目标服务器原有的名称和密码
func readCurrentConfig(server models.Server, dockerManager *docker_manager.DockerManager, file string) (string, error) {
	pending := server.PendingGameUserSettings
	if file == utils.GameIniFileName {
		pending = server.PendingGameIni
	}
	if pending != "" {
		return pending, nil
	}
	content, err := dockerManager.ReadConfigFile(server.ID, file)
	if err != nil {
		if errors.Is(err, docker_manager.ErrConfigNotFound) {
			return "", nil
		}
		return "", fmt.Errorf("读取%s失败: %w", file, err)
	}
	return content, nil
}

// findServer 按 ID 读取当前用户的服务器
func findServer(userID uint, serverID string) (models.Server, error) {
	var server models.Server
	id, err := strconv.ParseUint(serverID, 10, 32)
	if err != nil {
		return server, ErrInvalidServerID
	}
	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&server).Error; err != nil {
		return server, ErrServerNotFound
	}
	return server, nil
}

// buildServerConfig 生成服务器的配置，配置文件里去掉管理员密码
func buildServerConfig(server models.Server) (*models.ServerConfig, error) {
	dockerManager, err := docker_manager.GetDockerManager()
	if err != nil {
		return nil, fmt.Errorf("获取Docker管理器失败: %w", err)
	}
	gus, err := readCurrentConfig(server, dockerManager, utils.GameUserSettingsFileName)
	if err != nil {
		return nil, err
	}
	gameIni, err := readCurrentConfig(server, dockerManager, utils.GameIniFileName)
	if err != nil {
		return nil, err
	}
	return &models.ServerConfig{
		Format:           models.ServerConfigFormat,
		Version:          models.ServerConfigVersion,
		ExportedAt:       time.Now().In(appLocation()).Format(time.RFC3339),
		SessionName:      server.SessionName,
		Map:              server.Map,
		MaxPlayers:       server.MaxPlayers,
		ClusterID:        server.ClusterID,
		GameModIds:       server.GameModIds,
		ArkApi:           server.ArkApi,
		AutoRestart:      server.AutoRestart,
		AutoUpdate:       server.AutoUpdate,
		GameUserSettings: utils.RemoveIniKeys(gus, "ServerAdminPassword"),
		GameIni:          gameIni,
		ServerArgs:       serverArgsFromJSON(server),
	}, nil
}

// ExportConfig 导出服务器配置（不含端口和管理员密码），可以导入到其他服务器
func (s *ServerService) ExportConfig(userID uint, serverID string) (*models.ServerConfig, error) {
	server, err := findServer(userID, serverID)
	if err != nil {
		return nil, err
	}
	return buildServerConfig(server)
}

// normalizeConfigParts 校验要应用的部分，为空表示全部
func normalizeConfigParts(parts []string) (map[string]bool, error) {
	if len(parts) == 0 {
		parts = models.ConfigParts
	}
	set := make(map[string]bool, len(parts))
	for _, p := range parts {
		if !slices.Contains(models.ConfigParts, p) {
			return nil, fmt.Errorf("%w：未知的部分 %q", ErrInvalidServerConfig, p)
		}
		set[p] = true
	}
	return set, nil
}

// validateServerConfig 检查导入的配置文件格式
func validateServerConfig(cfg *models.ServerConfig) error {
	if cfg.Format != models.ServerConfigFormat {
		return fmt.Errorf("%w：不是本程序导出的服务器配置", ErrInvalidServerConfig)
	}
	if cfg.Version < 1 || cfg.Version > models.ServerConfigVersion {
		return fmt.Errorf("%w：不支持的版本 %d", ErrInvalidServerConfig, cfg.Version)
	}
	if cfg.MaxPlayers != 0 && (cfg.MaxPlayers < 1 || cfg.MaxPlayers > 200) {
		return fmt.Errorf("%w：最大玩家数需在 1-200 之间", ErrInvalidServerConfig)
	}
	if strings.ContainsAny(cfg.Map, " \t\r\n?") {
		return fmt.Errorf("%w：地图名称无效", ErrInvalidServerConfig)
	}
	return nil
}

// applyConfig 把配置里选中的部分写入目标服务器，规则和 UpdateServer 相同（未停止时配置文件先暂存）
func (s *ServerService) applyConfig(userID uint, target models.Server, cfg *models.ServerConfig, parts map[string]bool) (*models.ServerResponse, error) {
	req := models.ServerUpdateRequest{}

	if parts[models.ConfigPartOptions] {
		req.MaxPlayers = cfg.MaxPlayers
		req.ArkApi = &cfg.ArkApi
		req.AutoRestart = &cfg.AutoRestart
		req.AutoUpdate = &cfg.AutoUpdate
	}
	if parts[models.ConfigPartMap] {
		req.Map = cfg.Map
	}
	if parts[models.ConfigPartMods] {
		req.GameModIds = &cfg.GameModIds
	}
	if parts[models.ConfigPartCluster] {
		req.ClusterID = &cfg.ClusterID
	}
	if parts[models.ConfigPartServerArgs] && cfg.ServerArgs != nil {
		args := &models.ServerArgsRequest{
			QueryParams:     cfg.ServerArgs.QueryParams,
			CommandLineArgs: cfg.ServerArgs.CommandLineArgs,
			CustomArgs:      cfg.ServerArgs.CustomArgs,
		}
		if args.QueryParams == nil {
			args.QueryParams = map[string]string{}
		}
		if args.CommandLineArgs == nil {
			args.CommandLineArgs = map[string]any{}
		}
		if args.CustomArgs == nil {
			args.CustomArgs = []string{}
		}
		req.ServerArgs = args
	}
	if parts[models.ConfigPartConfigFiles] {
		if cfg.GameUserSettings != "" {
			dockerManager, err := docker_manager.GetDockerManager()
			if err != nil {
				return nil, fmt.Errorf("获取Docker管理器失败: %w", err)
			}
			current, err := readCurrentConfig(target, dockerManager, utils.GameUserSettingsFileName)
			if err != nil {
				return nil, err
			}
			gus := utils.KeepIniKeys(cfg.GameUserSettings, current, identityIniKeys...)
			// 同时改了人数时，配置文件里的人数也跟着改
			if req.MaxPlayers > 0 {
				for _, e := range utils.FindIniEntries(gus, "MaxPlayers") {
					gus = utils.SetIniValue(gus, e.Section, e.Key, strconv.Itoa(req.MaxPlayers))
				}
			}
			req.GameUserSettings = gus
		}
		req.GameIni = cfg.GameIni
	}

	response, _, err := s.UpdateServer(userID, strconv.FormatUint(uint64(target.ID), 10), req)
	return response, err
}

// ImportConfig 把导出的配置导入到一台已有的服务器，parts 为空表示全部。
// 端口、名称和管理员密码保持不变；服务器未停止时配置文件先暂存，下次启动时生效
func (s *ServerService) ImportConfig(userID uint, serverID string, cfg *models.ServerConfig, parts []string) (*models.ServerResponse, error) {
	if err := validateServerConfig(cfg); err != nil {
		return nil, err
	}
	partSet, err := normalizeConfigParts(parts)
	if err != nil {
		return nil, err
	}
	target, err := findServer(userID, serverID)
	if err != nil {
		return nil, err
	}
	return s.applyConfig(userID, target, cfg, partSet)
}

// CopyConfig 把一台服务器的配置复制到其他服务器，逐台应用，某台失败不影响其它的
func (s *ServerService) CopyConfig(userID uint, sourceID string, req models.CopyConfigRequest) ([]models.CopyConfigResult, error) {
	partSet, err := normalizeConfigParts(req.Parts)
	if err != nil {
		return nil, err
	}
	source, err := findServer(userID, sourceID)
	if err != nil {
		return nil, err
	}
	cfg, err := buildServerConfig(source)
	if err != nil {
		return nil, err
	}

	results := []models.CopyConfigResult{}
	seen := map[uint]bool{source.ID: true}
	for _, id := range req.TargetIDs {
		if seen[id] {
			continue
		}
		seen[id] = true
		result := models.CopyConfigResult{ServerID: id}
		target, err := findServer(userID, strconv.FormatUint(uint64(id), 10))
		if err == nil {
			result.SessionName = target.SessionName
			var updated *models.ServerResponse
			if updated, err = s.applyConfig(userID, target, cfg, partSet); err == nil {
				result.Success = true
				result.ConfigPending = updated.ConfigPending
			}
		}
		if err != nil {
			result.Error = err.Error()
			result.Code = ErrorCode(err)
		}
		results = append(results, result)
	}
	if len(results) == 0 {
		return nil, fmt.Errorf("%w：没有可以复制到的服务器", ErrInvalidServerConfig)
	}
	return results, nil
}
