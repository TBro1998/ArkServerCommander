package server

import (
	"encoding/json"
	"math"
	"strconv"
	"strings"

	"ark-server-commander/models"
	"ark-server-commander/service/docker_manager"
	"ark-server-commander/utils"
)

// serverProfile 服务器资料：发布到收录站（publish.go）和登记联机房（link.go）共用
type serverProfile struct {
	Name         string
	Map          string
	MaxPlayers   int
	Password     bool // 设置了加入密码（ServerPassword）
	PvP          bool
	ArkApi       bool
	Cluster      string
	Mods         []string
	RateXP       float64
	RateTame     float64
	RateHarvest  float64
	RateBreeding float64
	Listing      models.ListingProfile // 收录资料（简介、标签、QQ 群、Discord）
}

// buildServerProfile 从服务器记录、启动参数和配置文件得到服务器资料。
// 倍率、PvE 和加入密码先看启动参数里的 ?Key=Value（优先级高于配置文件），再看 GameUserSettings.ini / Game.ini
func buildServerProfile(server models.Server) serverProfile {
	profile := serverProfile{
		Name:       server.SessionName,
		Map:        server.Map,
		MaxPlayers: server.MaxPlayers,
		ArkApi:     server.ArkApi,
		Cluster:    server.ClusterID,
		Mods:       []string{},
		Listing:    listingOf(server),
	}
	for _, id := range strings.Split(server.GameModIds, ",") {
		if id = strings.TrimSpace(id); id != "" {
			profile.Mods = append(profile.Mods, id)
		}
	}

	lookup := serverConfigLookup(server)
	profile.PvP = !parseIniBool(lookup(cfgGUS, "serverPVE"))
	profile.Password = strings.TrimSpace(lookup(cfgGUS, "ServerPassword")) != ""
	profile.RateXP = parseRate(lookup(cfgGUS, "XPMultiplier"))
	profile.RateTame = parseRate(lookup(cfgGUS, "TamingSpeedMultiplier"))
	profile.RateHarvest = parseRate(lookup(cfgGUS, "HarvestAmountMultiplier"))
	profile.RateBreeding = parseRate(lookup(cfgGame, "BabyMatureSpeedMultiplier"))
	return profile
}

// 配置项所在的文件：GameUserSettings.ini 的 [ServerSettings]、Game.ini 的 [/script/shootergame.shootergamemode]
type cfgFile int

const (
	cfgGUS cfgFile = iota
	cfgGame
)

// serverConfigLookup 读一次服务器的启动参数和两个配置文件，返回查值函数：先看启动参数里的 ?Key=Value（优先级高于配置文件），
// 再看配置文件对应 section 里的键（同名键以最后一个为准），没有设置时返回空字符串
func serverConfigLookup(server models.Server) func(file cfgFile, key string) string {
	var args models.ServerArgs
	_ = json.Unmarshal([]byte(server.ServerArgsJSON), &args)
	var gus, game string
	if dm, err := docker_manager.GetDockerManager(); err == nil {
		gus, _ = dm.ReadConfigFile(server.ID, "GameUserSettings.ini")
		game, _ = dm.ReadConfigFile(server.ID, "Game.ini")
	}
	return func(file cfgFile, key string) string {
		for k, v := range args.QueryParams {
			if strings.EqualFold(k, key) {
				return v
			}
		}
		content, section := gus, "ServerSettings"
		if file == cfgGame {
			content, section = game, "/script/shootergame.shootergamemode"
		}
		value := ""
		for _, e := range utils.FindIniEntries(content, key) {
			if strings.EqualFold(e.Section, section) {
				value = e.Value
			}
		}
		return value
	}
}

func parseIniBool(v string) bool {
	return strings.EqualFold(strings.TrimSpace(v), "true") || strings.TrimSpace(v) == "1"
}

// parseRate 没有设置时游戏按 1 倍处理
func parseRate(v string) float64 {
	f, err := strconv.ParseFloat(strings.TrimSpace(v), 64)
	if err != nil || math.IsNaN(f) || math.IsInf(f, 0) || f <= 0 {
		return 1
	}
	return f
}
