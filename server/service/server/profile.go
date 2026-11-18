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

	var args models.ServerArgs
	_ = json.Unmarshal([]byte(server.ServerArgsJSON), &args)
	var gus, game string
	if dm, err := docker_manager.GetDockerManager(); err == nil {
		gus, _ = dm.ReadConfigFile(server.ID, "GameUserSettings.ini")
		game, _ = dm.ReadConfigFile(server.ID, "Game.ini")
	}
	lookup := func(key, file, section string) string {
		for k, v := range args.QueryParams {
			if strings.EqualFold(k, key) {
				return v
			}
		}
		value := ""
		for _, e := range utils.FindIniEntries(file, key) {
			if strings.EqualFold(e.Section, section) {
				value = e.Value // 同名键以最后一个为准
			}
		}
		return value
	}
	const gusSection, gameSection = "ServerSettings", "/script/shootergame.shootergamemode"
	profile.PvP = !parseIniBool(lookup("serverPVE", gus, gusSection))
	profile.Password = strings.TrimSpace(lookup("ServerPassword", gus, gusSection)) != ""
	profile.RateXP = parseRate(lookup("XPMultiplier", gus, gusSection))
	profile.RateTame = parseRate(lookup("TamingSpeedMultiplier", gus, gusSection))
	profile.RateHarvest = parseRate(lookup("HarvestAmountMultiplier", gus, gusSection))
	profile.RateBreeding = parseRate(lookup("BabyMatureSpeedMultiplier", game, gameSection))
	return profile
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
