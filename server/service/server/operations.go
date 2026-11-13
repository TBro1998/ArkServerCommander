package server

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"ark-server-commander/database"
	"ark-server-commander/models"
	"ark-server-commander/service/docker_manager"
	"ark-server-commander/utils"

	"go.uber.org/zap"
)

// getUserServer 读取属于该用户的服务器
func getUserServer(userID uint, serverID string) (models.Server, error) {
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

// RestartServer 重启服务器：运行中的先停止（停服前保存世界），再按当前配置启动。
// 已停止的服务器直接启动。在后台执行，状态依次为 stopping → starting → running
func (s *ServerService) RestartServer(userID uint, serverID string) error {
	server, err := getUserServer(userID, serverID)
	if err != nil {
		return err
	}
	if server.Status == "starting" || server.Status == "stopping" {
		return ErrServerBusy
	}
	dockerManager, err := docker_manager.GetDockerManager()
	if err != nil {
		return fmt.Errorf("获取Docker管理器失败: %w", err)
	}
	if err := s.claimForRestart(server); err != nil {
		return err
	}
	go s.restartWith(server, dockerManager, nil)
	return nil
}

// claimForRestart 在返回响应前检查能否开始重启，让用户立即看到"正在启停"的错误。
// 这只是提前检查，真正防止并发的是 restartWith 里基于 transitionStatus 的原子切换
func (s *ServerService) claimForRestart(server models.Server) error {
	if server.Status == "stopped" {
		return nil
	}
	var latest models.Server
	if err := database.DB.Select("status").First(&latest, server.ID).Error; err != nil {
		return err
	}
	if latest.Status == "starting" || latest.Status == "stopping" {
		return ErrServerBusy
	}
	return nil
}

// UpdateGame 更新游戏本体（steamcmd 校验并更新）和/或 Mod：mods 为 true 时重新下载全部 Mod，
// 否则只重新下载 modIDs 里的 Mod（必须是服务器配置的 Mod）。
// 运行中的服务器会先停止；清除安装标记后启动，入口脚本在启动时完成下载。在后台执行
func (s *ServerService) UpdateGame(userID uint, serverID string, game, mods bool, modIDs []string) error {
	if mods {
		modIDs = nil
	}
	if !game && !mods && len(modIDs) == 0 {
		return fmt.Errorf("%w: 至少选择更新游戏本体或 Mod", ErrInvalidCommand)
	}
	server, err := getUserServer(userID, serverID)
	if err != nil {
		return err
	}
	if len(modIDs) > 0 {
		configured := make(map[string]bool)
		for _, id := range strings.Split(server.GameModIds, ",") {
			configured[id] = true
		}
		for _, id := range modIDs {
			if !configured[id] {
				return fmt.Errorf("%w: Mod %s 不在服务器的 Mod 列表中", ErrInvalidModIDs, id)
			}
		}
	}
	if server.Status == "starting" || server.Status == "stopping" {
		return ErrServerBusy
	}
	dockerManager, err := docker_manager.GetDockerManager()
	if err != nil {
		return fmt.Errorf("获取Docker管理器失败: %w", err)
	}
	if err := s.claimForRestart(server); err != nil {
		return err
	}
	go func() {
		if err := s.restartWith(server, dockerManager, func() error {
			if err := dockerManager.RemoveInstallMarkers(server.ID, game, mods && server.GameModIds != ""); err != nil {
				return err
			}
			if err := docker_manager.RemoveInstalledMods(server.ID, modIDs); err != nil {
				return err
			}
			if mods || len(modIDs) > 0 {
				forgetOutdatedMods(server.ID, modIDs)
			}
			return nil
		}); err != nil {
			utils.Warn("更新游戏失败", zap.Uint("server_id", server.ID), zap.Error(err))
		}
	}()
	return nil
}

// GetLogs 返回服务器容器最近 tail 行日志；容器还没创建时返回空列表
func (s *ServerService) GetLogs(userID uint, serverID string, tail int) ([]string, error) {
	server, err := getUserServer(userID, serverID)
	if err != nil {
		return nil, err
	}
	dockerManager, err := docker_manager.GetDockerManager()
	if err != nil {
		return nil, fmt.Errorf("获取Docker管理器失败: %w", err)
	}
	containerName := utils.GetServerContainerName(server.ID)
	if exists, err := dockerManager.ContainerExists(containerName); err != nil {
		return nil, err
	} else if !exists {
		return []string{}, nil
	}
	return dockerManager.ContainerLogTail(containerName, tail)
}

// FollowLogs 持续输出服务器容器日志，直到 ctx 取消。容器不存在时立即返回
func (s *ServerService) FollowLogs(ctx context.Context, userID uint, serverID string, tail int, onLine func(string) error) error {
	server, err := getUserServer(userID, serverID)
	if err != nil {
		return err
	}
	dockerManager, err := docker_manager.GetDockerManager()
	if err != nil {
		return fmt.Errorf("获取Docker管理器失败: %w", err)
	}
	containerName := utils.GetServerContainerName(server.ID)
	if exists, err := dockerManager.ContainerExists(containerName); err != nil {
		return err
	} else if !exists {
		return nil
	}
	return dockerManager.FollowContainerLogs(ctx, containerName, tail, onLine)
}

// rconForServer 校验服务器在运行，返回执行 RCON 命令的函数
func rconForServer(server models.Server) (func(command string, timeout time.Duration) (string, error), error) {
	if server.Status != "running" {
		return nil, ErrServerNotRunning
	}
	dockerManager, err := docker_manager.GetDockerManager()
	if err != nil {
		return nil, fmt.Errorf("获取Docker管理器失败: %w", err)
	}
	containerName := utils.GetServerContainerName(server.ID)
	return func(command string, timeout time.Duration) (string, error) {
		out, err := dockerManager.RconCommand(containerName, server.RCONPort, server.AdminPassword, command, timeout)
		if err != nil {
			// 把常见的底层错误换成用户能理解的提示
			msg := err.Error()
			switch {
			case strings.Contains(msg, "Connection refused"), strings.Contains(msg, "timed out"):
				return "", ErrRconUnavailable
			case strings.Contains(msg, "RCON auth failed"):
				return "", ErrRconAuthFailed
			}
		}
		return out, err
	}, nil
}

// validateCommand RCON 命令不能为空、不能换行（一次只执行一条）
func validateCommand(command string) (string, error) {
	command = strings.TrimSpace(command)
	if command == "" || len(command) > 1000 || strings.ContainsAny(command, "\r\n\x00") {
		return "", ErrInvalidCommand
	}
	return command, nil
}

// RconExec 在运行中的服务器上执行一条 RCON 命令
func (s *ServerService) RconExec(userID uint, serverID, command string) (string, error) {
	command, err := validateCommand(command)
	if err != nil {
		return "", err
	}
	server, err := getUserServer(userID, serverID)
	if err != nil {
		return "", err
	}
	exec, err := rconForServer(server)
	if err != nil {
		return "", err
	}
	utils.Info("执行RCON命令", zap.Uint("server_id", server.ID), zap.String("command", command))
	return exec(command, time.Minute)
}

// Player 在线玩家
type Player struct {
	Index   int    `json:"index"`
	Name    string `json:"name"`
	SteamID string `json:"steam_id"`
}

// listPlayersPattern ListPlayers 的输出格式："0. 玩家名, 76561198000000000"
var listPlayersPattern = regexp.MustCompile(`^\s*(\d+)\.\s*(.*),\s*(\d+)\s*$`)

// parsePlayers 解析 ListPlayers 的输出，没有玩家时输出 "No Players Connected"
func parsePlayers(output string) []Player {
	players := []Player{}
	for _, line := range strings.Split(output, "\n") {
		m := listPlayersPattern.FindStringSubmatch(strings.TrimSpace(line))
		if m == nil {
			continue
		}
		index, _ := strconv.Atoi(m[1])
		players = append(players, Player{Index: index, Name: strings.TrimSpace(m[2]), SteamID: m[3]})
	}
	return players
}

// ListPlayers 通过 RCON 获取在线玩家
func (s *ServerService) ListPlayers(userID uint, serverID string) ([]Player, error) {
	server, err := getUserServer(userID, serverID)
	if err != nil {
		return nil, err
	}
	exec, err := rconForServer(server)
	if err != nil {
		return nil, err
	}
	out, err := exec("ListPlayers", 20*time.Second)
	if err != nil {
		return nil, err
	}
	return parsePlayers(out), nil
}

// broadcast 向服务器广播一条消息，失败只记录日志
func broadcast(server models.Server, message string) {
	exec, err := rconForServer(server)
	if err != nil {
		return
	}
	message = strings.NewReplacer("\r", " ", "\n", " ").Replace(message)
	if _, err := exec("Broadcast "+message, 20*time.Second); err != nil {
		utils.Warn("广播消息失败", zap.Uint("server_id", server.ID), zap.Error(err))
	}
}
