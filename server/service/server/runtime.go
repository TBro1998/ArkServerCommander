package server

import (
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"ark-server-commander/database"
	"ark-server-commander/models"
	"ark-server-commander/service/docker_manager"
	"ark-server-commander/utils"

	"go.uber.org/zap"
)

// monitorInterval 后台采集运行信息的间隔
const monitorInterval = 10 * time.Second

// phaseLogTail 判断运行阶段时读取的日志行数
const phaseLogTail = 400

// RuntimeMonitor 定期采集运行中服务器的阶段、在线人数和资源占用，结果缓存在内存里
type RuntimeMonitor struct {
	mu        sync.RWMutex
	info      map[uint]*models.RuntimeInfo
	prevStats map[uint]*docker_manager.StatsSample
	inFlight  map[uint]bool
}

var runtimeMonitor = &RuntimeMonitor{
	info:      map[uint]*models.RuntimeInfo{},
	prevStats: map[uint]*docker_manager.StatsSample{},
	inFlight:  map[uint]bool{},
}

// Get 返回服务器最近一次采集的信息副本，没有时返回 nil
func (m *RuntimeMonitor) Get(serverID uint) *models.RuntimeInfo {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if info, ok := m.info[serverID]; ok {
		copied := *info
		return &copied
	}
	return nil
}

// Forget 删除服务器的缓存（停止或删除后调用）
func (m *RuntimeMonitor) Forget(serverID uint) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.info, serverID)
	delete(m.prevStats, serverID)
}

// StartRuntimeMonitor 启动后台采集协程
func StartRuntimeMonitor() {
	go func() {
		// 启动后立即采集一次，前端打开页面时就有数据
		runtimeMonitor.tick()
		ticker := time.NewTicker(monitorInterval)
		defer ticker.Stop()
		for range ticker.C {
			runtimeMonitor.tick()
		}
	}()
}

func (m *RuntimeMonitor) tick() {
	var servers []models.Server
	if err := database.DB.Select("id", "status", "query_port").Where("status = ?", "running").Find(&servers).Error; err != nil {
		utils.Warn("采集运行信息失败：读取服务器列表出错", zap.Error(err))
		return
	}

	dm, err := docker_manager.GetDockerManager()
	if err != nil {
		return
	}

	running := make(map[uint]bool, len(servers))
	for _, server := range servers {
		running[server.ID] = true
		m.mu.Lock()
		busy := m.inFlight[server.ID]
		if !busy {
			m.inFlight[server.ID] = true
		}
		m.mu.Unlock()
		if busy {
			// 上一次采集还没结束（例如 Docker 响应慢），跳过这一轮
			continue
		}
		go func(server models.Server) {
			defer func() {
				m.mu.Lock()
				delete(m.inFlight, server.ID)
				m.mu.Unlock()
			}()
			m.collect(dm, server)
		}(server)
	}

	// 清理已经不在运行的服务器
	m.mu.Lock()
	for id := range m.info {
		if !running[id] {
			delete(m.info, id)
			delete(m.prevStats, id)
		}
	}
	m.mu.Unlock()
}

func (m *RuntimeMonitor) collect(dm *docker_manager.DockerManager, server models.Server) {
	containerName := utils.GetServerContainerName(server.ID)
	info := &models.RuntimeInfo{}

	if lines, err := dm.ContainerLogTail(containerName, phaseLogTail); err == nil {
		applyPhase(info, lines)
	}

	// 游戏进程启动后才去探测查询端口：下载阶段探测没有意义，还会多一次 docker exec
	if info.Phase == models.PhaseLaunching || info.Phase == "" {
		if a2s, err := dm.QueryServerInfo(containerName, server.QueryPort); err == nil {
			info.Phase = models.PhaseReady
			info.Players = a2s.Players
			info.MaxPlayers = a2s.MaxPlayers
			info.ServerName = a2s.Name
		}
	}

	if sample, err := dm.ContainerStats(containerName); err == nil {
		m.mu.Lock()
		prev := m.prevStats[server.ID]
		m.prevStats[server.ID] = sample
		m.mu.Unlock()
		info.CPUPercent = docker_manager.CPUPercent(prev, sample)
		info.MemoryBytes = sample.MemoryUsage
		info.MemoryLimit = sample.MemoryLimit
	}

	info.UpdatedAt = time.Now().Format(time.RFC3339)
	m.mu.Lock()
	m.info[server.ID] = info
	m.mu.Unlock()
}

var (
	steamProgressPattern = regexp.MustCompile(`Update state \(0x[0-9a-fA-F]+\) [a-z ]+, progress: ([0-9.]+)`)
	modProgressPattern   = regexp.MustCompile(`(?:Downloading mod|Installing mod|Mod [0-9]+ already installed|Failed to install mod [0-9]+)\D*?(?:[0-9]+ )?\(([0-9]+)/([0-9]+)\)`)
)

// applyPhase 从入口脚本的日志判断当前阶段：从最后一行往前找，最近出现的标志决定阶段。
// 标志行见 ASE-Server-Docker/scripts/start_server.sh
func applyPhase(info *models.RuntimeInfo, lines []string) {
	for i := len(lines) - 1; i >= 0; i-- {
		line := lines[i]
		switch {
		case strings.Contains(line, "Starting ARK server"):
			info.Phase = models.PhaseLaunching
			return
		case strings.Contains(line, "Skipping mod update"),
			strings.Contains(line, "Server files already installed"):
			// 这两行之后紧跟着启动游戏，日志还没刷到 Starting 时也按启动中处理
			info.Phase = models.PhaseLaunching
			return
		}
		if m := modProgressPattern.FindStringSubmatch(line); m != nil {
			info.Phase = models.PhaseInstallingMods
			info.ModCurrent, _ = strconv.Atoi(m[1])
			info.ModTotal, _ = strconv.Atoi(m[2])
			return
		}
		if m := steamProgressPattern.FindStringSubmatch(line); m != nil {
			info.Phase = models.PhaseDownloading
			info.Progress, _ = strconv.ParseFloat(m[1], 64)
			return
		}
		if strings.Contains(line, "Server files not found") ||
			strings.Contains(line, "Updating ARK server") ||
			strings.Contains(line, "Updating Steam client") {
			info.Phase = models.PhaseDownloading
			return
		}
	}
}
