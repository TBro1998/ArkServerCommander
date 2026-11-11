package server

import (
	"sort"
	"sync"
	"time"

	"ark-server-commander/database"
	"ark-server-commander/models"
	"ark-server-commander/service/docker_manager"
	"ark-server-commander/utils"

	"go.uber.org/zap"
)

// maxSeedCandidates 检查来源时最多挂载的服务器文件卷数量
const maxSeedCandidates = 20

// seeding 正在复制服务端文件的服务器 ID -> 来源服务器名称，用于在界面上显示"复制服务端文件"阶段
var seeding sync.Map

// seedSourceName 返回正在复制时的来源服务器名称
func seedSourceName(serverID uint) (string, bool) {
	v, ok := seeding.Load(serverID)
	if !ok {
		return "", false
	}
	return v.(string), true
}

// pickSeedSource 选出复制来源：必须已安装完成（有标记和主程序），优先已停止的服务器（文件不会在复制时变动），
// 其次按 ID 顺序。targetHasFiles 为 true（目标已有服务端主程序，例如"更新游戏"刚清除了标记）时不复制
func pickSeedSource(targetHasFiles bool, candidates []models.Server, states map[uint]docker_manager.ServerFilesState) (models.Server, bool) {
	if targetHasFiles {
		return models.Server{}, false
	}
	sorted := append([]models.Server(nil), candidates...)
	sort.SliceStable(sorted, func(i, j int) bool {
		si, sj := sorted[i].Status == "stopped", sorted[j].Status == "stopped"
		if si != sj {
			return si
		}
		return sorted[i].ID < sorted[j].ID
	})
	for _, c := range sorted {
		if st := states[c.ID]; st.Installed && st.Executable {
			return c, true
		}
	}
	return models.Server{}, false
}

// seedServerFiles 服务器还没有下载过服务端时，从另一台已安装的服务器复制服务端文件和 Mod，
// 首次启动就不用从 Steam 下载约 20 GB；复制后入口脚本会 validate，只下载与最新版本的差异。
// 任何一步失败都只记录日志，照常启动（入口脚本会自行从 Steam 下载）
func seedServerFiles(server models.Server, dm *docker_manager.DockerManager) {
	var others []models.Server
	if err := database.DB.Select("id", "session_name", "status").Where("id <> ?", server.ID).Find(&others).Error; err != nil {
		utils.Warn("读取服务器列表失败，跳过复制服务端文件", zap.Error(err))
		return
	}
	if len(others) == 0 {
		return
	}
	// 优先检查已停止的服务器，数量多时只检查前若干台
	sort.SliceStable(others, func(i, j int) bool { return others[i].Status == "stopped" && others[j].Status != "stopped" })
	if len(others) > maxSeedCandidates {
		others = others[:maxSeedCandidates]
	}

	ids := []uint{server.ID}
	for _, o := range others {
		ids = append(ids, o.ID)
	}
	states, err := dm.CheckServerFiles(ids)
	if err != nil {
		utils.Warn("检查服务端文件失败，跳过复制", zap.Uint("server_id", server.ID), zap.Error(err))
		return
	}

	source, ok := pickSeedSource(states[server.ID].Executable, others, states)
	if !ok {
		return
	}

	seeding.Store(server.ID, source.SessionName)
	defer seeding.Delete(server.ID)

	utils.Info("开始从已安装的服务器复制服务端文件",
		zap.Uint("server_id", server.ID), zap.Uint("source_id", source.ID), zap.String("source", source.SessionName))
	start := time.Now()
	if err := dm.CopyServerFiles(source.ID, server.ID); err != nil {
		utils.Warn("复制服务端文件失败，将从 Steam 下载", zap.Uint("server_id", server.ID), zap.Error(err))
		return
	}
	utils.Info("服务端文件复制完成", zap.Uint("server_id", server.ID), zap.Duration("elapsed", time.Since(start)))
}
