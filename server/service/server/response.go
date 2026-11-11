package server

import (
	"strings"

	"ark-server-commander/models"
	"ark-server-commander/utils"
)

// toResponse 把数据库记录转换为接口响应（不含配置文件内容和启动参数），
// 运行中的服务器附带后台监控采集的实时信息
func toResponse(server models.Server) models.ServerResponse {
	response := models.ServerResponse{
		ID:            server.ID,
		Identifier:    server.Identifier,
		SessionName:   server.SessionName,
		ClusterID:     server.ClusterID,
		Port:          server.Port,
		QueryPort:     server.QueryPort,
		RCONPort:      server.RCONPort,
		Map:           server.Map,
		MaxPlayers:    server.MaxPlayers,
		GameModIds:    server.GameModIds,
		Status:        server.Status,
		LastError:     server.LastError,
		AutoRestart:   server.AutoRestart,
		AutoUpdate:    server.AutoUpdate,
		ConfigPending: server.PendingGameUserSettings != "" || server.PendingGameIni != "",
		UserID:        server.UserID,
		CreatedAt:     server.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:     server.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
	if server.Status == "running" {
		response.Runtime = runtimeMonitor.Get(server.ID)
	}
	return response
}

// normalizeServerFields 校验并规范化集群 ID 和 Mod ID 列表
func normalizeServerFields(clusterID, modIDs string) (string, string, error) {
	clusterID = strings.TrimSpace(clusterID)
	if clusterID != "" && !utils.ClusterIDPattern.MatchString(clusterID) {
		return "", "", ErrInvalidClusterID
	}
	modIDs = utils.NormalizeModIDs(modIDs)
	if modIDs != "" && !utils.ModIDsPattern.MatchString(modIDs) {
		return "", "", ErrInvalidModIDs
	}
	return clusterID, modIDs, nil
}
