package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"ark-server-commander/config"
	"ark-server-commander/models"
	"ark-server-commander/service/docker_manager"
	"ark-server-commander/utils"

	"go.uber.org/zap"
)

var httpClient = &http.Client{Timeout: 10 * time.Second}

// ---------- 公网地址 ----------

var publicHost struct {
	sync.Mutex
	value   string
	fetched time.Time
}

// PublicHost 返回玩家连接用的地址：优先使用 PUBLIC_HOST，否则探测公网 IP（缓存 1 小时，失败时 5 分钟后重试）
func PublicHost() string {
	if config.PublicHost != "" {
		return config.PublicHost
	}
	publicHost.Lock()
	defer publicHost.Unlock()
	ttl := time.Hour
	if publicHost.value == "" {
		ttl = 5 * time.Minute
	}
	if !publicHost.fetched.IsZero() && time.Since(publicHost.fetched) < ttl {
		return publicHost.value
	}
	publicHost.fetched = time.Now()
	for _, endpoint := range []string{"https://api.ipify.org", "https://ifconfig.me/ip"} {
		resp, err := httpClient.Get(endpoint)
		if err != nil {
			continue
		}
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 64))
		resp.Body.Close()
		ip := strings.TrimSpace(string(body))
		if resp.StatusCode == http.StatusOK && net.ParseIP(ip) != nil {
			publicHost.value = ip
			return ip
		}
	}
	utils.Warn("探测公网 IP 失败，可通过 PUBLIC_HOST 手动设置")
	return publicHost.value
}

// ---------- 磁盘占用 ----------

var volumeSizes struct {
	sync.Mutex
	sizes   map[string]int64
	fetched time.Time
	err     error
}

// cachedVolumeSizes 卷大小计算很慢，缓存 5 分钟；refresh 为 true 时强制刷新
func cachedVolumeSizes(refresh bool) (map[string]int64, time.Time, error) {
	volumeSizes.Lock()
	defer volumeSizes.Unlock()
	if !refresh && volumeSizes.sizes != nil && time.Since(volumeSizes.fetched) < 5*time.Minute {
		return volumeSizes.sizes, volumeSizes.fetched, nil
	}
	// 卷绑定在数据目录下，直接统计目录大小（Docker 的 DiskUsage 查不到绑定卷）
	sizes, err := docker_manager.LocalVolumeSizes()
	if err != nil {
		return volumeSizes.sizes, volumeSizes.fetched, err
	}
	volumeSizes.sizes = sizes
	volumeSizes.fetched = time.Now()
	return sizes, volumeSizes.fetched, nil
}

// ServerDiskUsage 一台服务器的磁盘占用（字节）
type ServerDiskUsage struct {
	ServerID    uint   `json:"server_id"`
	SessionName string `json:"session_name"`
	Files       int64  `json:"files"`   // 服务端文件卷
	Data        int64  `json:"data"`    // 存档和配置（数据卷）
	Plugins     int64  `json:"plugins"` // 插件卷
	Backups     int64  `json:"backups"` // 备份文件
	Total       int64  `json:"total"`
}

// Overview 首页概览
type Overview struct {
	PublicHost     string            `json:"public_host"`
	TotalServers   int               `json:"total_servers"`
	RunningServers int               `json:"running_servers"`
	ReadyServers   int               `json:"ready_servers"`
	OnlinePlayers  int               `json:"online_players"`
	CPUPercent     float64           `json:"cpu_percent"`
	MemoryBytes    uint64            `json:"memory_bytes"`
	Disk           []ServerDiskUsage `json:"disk"`
	ClusterVolumes map[string]int64  `json:"cluster_volumes"`
	DiskTotal      int64             `json:"disk_total"`
	DiskUpdatedAt  string            `json:"disk_updated_at,omitempty"`
	DiskError      string            `json:"disk_error,omitempty"`
	RecentErrors   []ServerError     `json:"recent_errors"`
}

// ServerError 最近一次启停失败的服务器
type ServerError struct {
	ServerID    uint   `json:"server_id"`
	SessionName string `json:"session_name"`
	Error       string `json:"error"`
}

// GetOverview 汇总服务器状态、在线人数、资源和磁盘占用
func (s *ServerService) GetOverview(userID uint, refreshDisk bool) (*Overview, error) {
	servers, err := s.GetServers(userID)
	if err != nil {
		return nil, err
	}

	ov := &Overview{
		PublicHost:     PublicHost(),
		TotalServers:   len(servers),
		Disk:           []ServerDiskUsage{},
		ClusterVolumes: map[string]int64{},
		RecentErrors:   []ServerError{},
	}
	for _, sv := range servers {
		if sv.Status == "running" {
			ov.RunningServers++
		}
		if rt := sv.Runtime; rt != nil {
			if rt.Phase == models.PhaseReady {
				ov.ReadyServers++
			}
			ov.OnlinePlayers += rt.Players
			ov.CPUPercent += rt.CPUPercent
			ov.MemoryBytes += rt.MemoryBytes
		}
		if sv.LastError != "" {
			ov.RecentErrors = append(ov.RecentErrors, ServerError{ServerID: sv.ID, SessionName: sv.SessionName, Error: sv.LastError})
		}
	}

	sizes, fetched, err := cachedVolumeSizes(refreshDisk)
	if err != nil {
		ov.DiskError = err.Error()
	}
	if !fetched.IsZero() {
		ov.DiskUpdatedAt = fetched.Format(time.RFC3339)
	}
	clusters := map[string]bool{}
	for _, sv := range servers {
		usage := ServerDiskUsage{
			ServerID:    sv.ID,
			SessionName: sv.SessionName,
			Files:       sizes[utils.GetServerFilesVolumeName(sv.ID)],
			Data:        sizes[utils.GetServerVolumeName(sv.ID)],
			Plugins:     sizes[utils.GetServerPluginsVolumeName(sv.ID)],
			Backups:     backupDirSize(sv.ID),
		}
		usage.Total = usage.Files + usage.Data + usage.Plugins + usage.Backups
		ov.DiskTotal += usage.Total
		ov.Disk = append(ov.Disk, usage)
		if cv := docker_manager.ExpectedClusterVolume(sv.ClusterID); cv != "" {
			clusters[cv] = true
		}
	}
	for cv := range clusters {
		ov.ClusterVolumes[cv] = sizes[cv]
		ov.DiskTotal += sizes[cv]
	}
	return ov, nil
}

// ---------- 端口建议 ----------

// SuggestedPorts 推荐的端口组合
type SuggestedPorts struct {
	Port      int `json:"port"`
	QueryPort int `json:"query_port"`
	RCONPort  int `json:"rcon_port"`
}

// SuggestPorts 从 ARK 默认端口开始，找第一组与现有服务器都不冲突的端口。
// 游戏端口每次加 2（还要占用 +1），查询端口和 RCON 端口每次加 1
func (s *ServerService) SuggestPorts() SuggestedPorts {
	for i := 0; i < 500; i++ {
		p := SuggestedPorts{Port: 7777 + 2*i, QueryPort: 27015 + i, RCONPort: 32330 + i}
		if checkPorts(0, p.Port, p.QueryPort, p.RCONPort) == nil {
			return p
		}
	}
	return SuggestedPorts{Port: 7777, QueryPort: 27015, RCONPort: 32330}
}

// ---------- 创意工坊 Mod 信息 ----------

// ModInfo 创意工坊 Mod 的基本信息
type ModInfo struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Checked     bool   `json:"checked"` // 是否从 Steam 拿到了结果；为 false 时表示查询失败（例如无法访问 Steam），Found 无意义
	Found       bool   `json:"found"`
	FileSize    int64  `json:"file_size,omitempty"`
	TimeUpdated int64  `json:"time_updated,omitempty"`
	PreviewURL  string `json:"preview_url,omitempty"`
}

var modCache sync.Map // id -> cachedMod

type cachedMod struct {
	info    ModInfo
	fetched time.Time
}

// GetModInfo 查询 Mod 名称等信息（Steam 公开接口，不需要 API Key），结果缓存 6 小时
func (s *ServerService) GetModInfo(ids string) ([]ModInfo, error) {
	ids = utils.NormalizeModIDs(ids)
	if ids == "" {
		return []ModInfo{}, nil
	}
	if !utils.ModIDsPattern.MatchString(ids) {
		return nil, ErrInvalidModIDs
	}
	list := strings.Split(ids, ",")
	if len(list) > 100 {
		list = list[:100]
	}

	var missing []string
	for _, id := range list {
		if c, ok := modCache.Load(id); !ok || expired(c.(cachedMod)) {
			missing = append(missing, id)
		}
	}
	if len(missing) > 0 {
		if err := fetchMods(missing); err != nil {
			utils.Warn("查询创意工坊 Mod 信息失败", zap.Error(err))
			// 失败结果也缓存一段时间，无法访问 Steam 时不必每次都等到超时；已有成功结果的保留
			now := time.Now()
			for _, id := range missing {
				if c, ok := modCache.Load(id); ok && c.(cachedMod).info.Checked {
					continue
				}
				modCache.Store(id, cachedMod{info: ModInfo{ID: id}, fetched: now})
			}
		}
	}

	result := make([]ModInfo, 0, len(list))
	for _, id := range list {
		if c, ok := modCache.Load(id); ok {
			result = append(result, c.(cachedMod).info)
		} else {
			result = append(result, ModInfo{ID: id})
		}
	}
	return result, nil
}

// expired 成功结果缓存 6 小时，失败结果 5 分钟后重试
func expired(c cachedMod) bool {
	if c.info.Checked {
		return time.Since(c.fetched) > 6*time.Hour
	}
	return time.Since(c.fetched) > 5*time.Minute
}

func fetchMods(ids []string) error {
	form := url.Values{}
	form.Set("itemcount", strconv.Itoa(len(ids)))
	for i, id := range ids {
		form.Set(fmt.Sprintf("publishedfileids[%d]", i), id)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		"https://api.steampowered.com/ISteamRemoteStorage/GetPublishedFileDetails/v1/",
		strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Steam 接口返回 %d", resp.StatusCode)
	}

	var body struct {
		Response struct {
			PublishedFileDetails []struct {
				PublishedFileID string `json:"publishedfileid"`
				Result          int    `json:"result"`
				Title           string `json:"title"`
				FileSize        any    `json:"file_size"`
				TimeUpdated     int64  `json:"time_updated"`
				PreviewURL      string `json:"preview_url"`
			} `json:"publishedfiledetails"`
		} `json:"response"`
	}
	if err := json.NewDecoder(io.LimitReader(resp.Body, 4<<20)).Decode(&body); err != nil {
		return err
	}
	now := time.Now()
	for _, d := range body.Response.PublishedFileDetails {
		info := ModInfo{ID: d.PublishedFileID, Checked: true, Found: d.Result == 1, Title: d.Title, TimeUpdated: d.TimeUpdated, PreviewURL: d.PreviewURL}
		switch v := d.FileSize.(type) {
		case float64:
			info.FileSize = int64(v)
		case string:
			info.FileSize, _ = strconv.ParseInt(v, 10, 64)
		}
		modCache.Store(d.PublishedFileID, cachedMod{info: info, fetched: now})
	}
	return nil
}
