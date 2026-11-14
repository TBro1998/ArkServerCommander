package server

import (
	"errors"
	"fmt"
	"maps"
	"slices"
	"strconv"
	"strings"

	"ark-server-commander/database"
	"ark-server-commander/models"
	"ark-server-commander/service/docker_manager"
	"ark-server-commander/utils"
)

// SetConfigValues 在 GameUserSettings.ini 或 Game.ini 的一个 section 里设置若干键，其它内容原样保留。
// 以暂存内容（如果有）或数据卷里的文件为基础修改，再按 UpdateServer 的规则保存：服务器未停止时存为暂存配置。
// 读取文件失败时直接报错，不能当作空文件，否则会用只含这几个键的内容覆盖整个配置文件。
// 同名键会合并为一行，Game.ini 里需要重复出现的键（如 OverrideNamedEngramEntries）不适合用这个方法
func (s *ServerService) SetConfigValues(userID uint, serverID, file, section string, values map[string]string) (*models.ServerResponse, error) {
	if file != utils.GameUserSettingsFileName && file != utils.GameIniFileName {
		return nil, ErrInvalidConfigValue
	}
	if strings.TrimSpace(section) == "" || strings.ContainsAny(section, "\r\n[]") || len(values) == 0 {
		return nil, ErrInvalidConfigValue
	}
	for key, value := range values {
		if strings.TrimSpace(key) == "" || strings.ContainsAny(key, "\r\n=[]") || strings.ContainsAny(value, "\r\n") {
			return nil, ErrInvalidConfigValue
		}
	}

	id, err := strconv.ParseUint(serverID, 10, 32)
	if err != nil {
		return nil, ErrInvalidServerID
	}
	var server models.Server
	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&server).Error; err != nil {
		return nil, ErrServerNotFound
	}

	current := server.PendingGameUserSettings
	if file == utils.GameIniFileName {
		current = server.PendingGameIni
	}
	if current == "" {
		dockerManager, err := docker_manager.GetDockerManager()
		if err != nil {
			return nil, fmt.Errorf("获取Docker管理器失败: %w", err)
		}
		current, err = dockerManager.ReadConfigFile(server.ID, file)
		if err != nil && !errors.Is(err, docker_manager.ErrConfigNotFound) {
			return nil, fmt.Errorf("读取%s失败: %w", file, err)
		}
	}

	updated := current
	for _, key := range slices.Sorted(maps.Keys(values)) {
		updated = utils.SetIniValue(updated, strings.TrimSpace(section), strings.TrimSpace(key), values[key])
	}

	req := models.ServerUpdateRequest{}
	if file == utils.GameUserSettingsFileName {
		req.GameUserSettings = updated
	} else {
		req.GameIni = updated
	}
	response, _, err := s.UpdateServer(userID, serverID, req)
	return response, err
}
