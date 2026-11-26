package server

import (
	"strings"
	"sync"
	"time"

	"ark-server-commander/database"
	"ark-server-commander/models"
	"ark-server-commander/service/docker_manager"
	"ark-server-commander/utils"

	"go.uber.org/zap"
)

// modCheckInterval 后台检查 Mod 更新的间隔；Steam 结果本身缓存 6 小时，详情页可以强制刷新
const modCheckInterval = 30 * time.Minute

// ServerMod 服务器配置的 Mod 及其安装状态
type ServerMod struct {
	ModInfo
	Installed   bool  `json:"installed"`              // 服务端文件卷里已安装
	InstalledAt int64 `json:"installed_at,omitempty"` // 安装时间（Unix 秒），未知时没有
	HasUpdate   bool  `json:"has_update"`             // 创意工坊上的更新时间晚于安装时间
}

// outdatedMods 每台服务器有更新的 Mod ID，由后台检查和详情页查询写入，随服务器列表返回
var outdatedMods sync.Map // uint -> []string

// GetServerMods 返回服务器配置的 Mod（按加载顺序）及安装、更新状态。refresh 为 true 时忽略缓存重新查询 Steam。
// ASA 的 Mod 来自 CurseForge，由服务端启动时自己下载和更新，管理端查不到安装状态（CurseForge 接口要 API Key），
// 只按顺序返回 ID
func (s *ServerService) GetServerMods(userID uint, serverID string, refresh bool) ([]ServerMod, error) {
	server, err := getUserServer(userID, serverID)
	if err != nil {
		return nil, err
	}
	if utils.NormalizeGame(server.Game) == utils.GameASA {
		return asaServerMods(server.GameModIds), nil
	}
	return checkServerMods(server.ID, server.GameModIds, refresh)
}

// asaServerMods ASA 服务器的 Mod 列表：只有 ID，视为已安装（服务端每次启动都会下载缺少的和有更新的）
func asaServerMods(modIDs string) []ServerMod {
	result := []ServerMod{}
	ids := utils.NormalizeModIDs(modIDs)
	if ids == "" {
		return result
	}
	for _, id := range strings.Split(ids, ",") {
		result = append(result, ServerMod{ModInfo: ModInfo{ID: id}, Installed: true})
	}
	return result
}

// checkServerMods 对比创意工坊更新时间和安装时间，结果同时记入 outdatedMods
func checkServerMods(serverID uint, modIDs string, refresh bool) ([]ServerMod, error) {
	ids := utils.NormalizeModIDs(modIDs)
	if ids == "" {
		outdatedMods.Delete(serverID)
		return []ServerMod{}, nil
	}
	installed, err := docker_manager.InstalledMods(serverID)
	if err != nil {
		return nil, err
	}

	infos := lookupMods(strings.Split(ids, ","), refresh)
	result := make([]ServerMod, 0, len(infos))
	var outdated []string
	for _, info := range infos {
		m := ServerMod{ModInfo: info}
		if ts, ok := installed[info.ID]; ok {
			m.Installed = true
			m.InstalledAt = ts
			m.HasUpdate = info.Checked && info.Found && ts > 0 && info.TimeUpdated > ts
		}
		if m.HasUpdate {
			outdated = append(outdated, info.ID)
		}
		result = append(result, m)
	}
	if len(outdated) > 0 {
		outdatedMods.Store(serverID, outdated)
	} else {
		outdatedMods.Delete(serverID)
	}
	return result, nil
}

// modsOutdatedFor 返回服务器有更新的 Mod，只保留仍在当前配置里的
func modsOutdatedFor(serverID uint, modIDs string) []string {
	v, ok := outdatedMods.Load(serverID)
	if !ok {
		return nil
	}
	configured := make(map[string]bool)
	for _, id := range strings.Split(modIDs, ",") {
		configured[id] = true
	}
	var result []string
	for _, id := range v.([]string) {
		if configured[id] {
			result = append(result, id)
		}
	}
	return result
}

// forgetOutdatedMods 清除指定 Mod 的更新标记（已安排重新下载）；ids 为空时清除全部
func forgetOutdatedMods(serverID uint, ids []string) {
	if len(ids) == 0 {
		outdatedMods.Delete(serverID)
		return
	}
	v, ok := outdatedMods.Load(serverID)
	if !ok {
		return
	}
	remove := make(map[string]bool, len(ids))
	for _, id := range ids {
		remove[id] = true
	}
	var kept []string
	for _, id := range v.([]string) {
		if !remove[id] {
			kept = append(kept, id)
		}
	}
	if len(kept) > 0 {
		outdatedMods.Store(serverID, kept)
	} else {
		outdatedMods.Delete(serverID)
	}
}

// StartModUpdateChecker 后台定期检查所有服务器的 Mod 是否有更新
func StartModUpdateChecker() {
	go func() {
		// 等启动时的其他工作完成后再查，避免和容器状态校正抢 Steam / 磁盘
		time.Sleep(time.Minute)
		checkAllMods()
		ticker := time.NewTicker(modCheckInterval)
		defer ticker.Stop()
		for range ticker.C {
			checkAllMods()
		}
	}()
}

func checkAllMods() {
	var servers []models.Server
	// ASA 的 Mod 不在创意工坊上，由服务端自己更新
	if err := database.DB.Select("id", "game_mod_ids").Where("game_mod_ids <> '' AND game <> ?", utils.GameASA).Find(&servers).Error; err != nil {
		utils.Warn("检查 Mod 更新失败：读取服务器列表出错", zap.Error(err))
		return
	}
	for _, server := range servers {
		if _, err := checkServerMods(server.ID, server.GameModIds, false); err != nil {
			utils.Warn("检查 Mod 更新失败", zap.Uint("server_id", server.ID), zap.Error(err))
		}
	}
}
