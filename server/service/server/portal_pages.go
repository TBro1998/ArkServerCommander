package server

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"unicode/utf8"

	"ark-server-commander/database"
	"ark-server-commander/models"
	"ark-server-commander/utils"

	"gorm.io/gorm"
)

// 玩家站的通用页面：服主写的规则、新手指南等，加上每台服务器自动生成的「服务器设置」页。
// 页面由几种通用的块组成（heading 标题、paragraph 段落、list 列表、kv 名称和值），玩家站网页和联机助手只认这几种块，
// 以后加新的内容只要换成块，客户端不用改。文本可以是字符串，也可以是 {"zh": "...", "en": "..."}（客户端按界面语言选）。
// 格式见 ArkPlayerPortal/docs/snapshot.md 的「页面」

const (
	maxPortalPages     = 50
	maxPageTitleRunes  = 64
	maxPageBodyRunes   = 20000
	pageKindCustom     = "custom"
	pageKindSettings   = "settings"
	settingsPagePrefix = "settings-"
)

var (
	ErrPageNotFound = errors.New("页面不存在")
	ErrInvalidPage  = errors.New("页面无效：标题不能为空、最多 64 字、不能换行，内容最多 20000 字")
	ErrTooManyPages = errors.New("页面最多 50 个，请先删除不用的")
)

// i18n 两种语言的文本
type i18n struct {
	Zh string `json:"zh"`
	En string `json:"en"`
}

type portalBlock struct {
	Type  string `json:"type"`            // heading | paragraph | list | kv
	Text  any    `json:"text,omitempty"`  // heading、paragraph：字符串或 i18n
	Items []any  `json:"items,omitempty"` // list：字符串或 i18n；kv：portalKV
}

type portalKV struct {
	Label any `json:"label"` // 字符串或 i18n
	Value any `json:"value"` // 字符串、数字或布尔值（客户端把布尔值显示成开 / 关）
}

type portalPage struct {
	ID        string        `json:"id"`
	Kind      string        `json:"kind"` // custom：服主写的页面；settings：服务器设置（server_key 是哪台）
	Title     any           `json:"title"`
	ServerKey string        `json:"server_key,omitempty"`
	Blocks    []portalBlock `json:"blocks"`
}

// parsePageBody 把服主写的纯文本变成块：以「# 」开头的行是标题，以「- 」或「* 」开头的连续行是列表，
// 其余的连续行是一个段落（保留换行），空行分隔。不解析 HTML 和其他 Markdown
func parsePageBody(body string) []portalBlock {
	blocks := []portalBlock{}
	var para []string
	var list []any
	flush := func() {
		if len(para) > 0 {
			blocks = append(blocks, portalBlock{Type: "paragraph", Text: strings.Join(para, "\n")})
			para = nil
		}
		if len(list) > 0 {
			blocks = append(blocks, portalBlock{Type: "list", Items: list})
			list = nil
		}
	}
	for _, line := range strings.Split(strings.ReplaceAll(body, "\r\n", "\n"), "\n") {
		t := strings.TrimSpace(line)
		switch {
		case t == "":
			flush()
		case strings.HasPrefix(t, "# ") || strings.HasPrefix(t, "## "):
			flush()
			blocks = append(blocks, portalBlock{Type: "heading", Text: strings.TrimSpace(strings.TrimLeft(t, "#"))})
		case strings.HasPrefix(t, "- ") || strings.HasPrefix(t, "* "):
			if len(para) > 0 {
				blocks = append(blocks, portalBlock{Type: "paragraph", Text: strings.Join(para, "\n")})
				para = nil
			}
			list = append(list, strings.TrimSpace(t[2:]))
		default:
			if len(list) > 0 {
				blocks = append(blocks, portalBlock{Type: "list", Items: list})
				list = nil
			}
			para = append(para, t)
		}
	}
	flush()
	return blocks
}

// ---- 服务器设置一览 ----

type settingKind int

const (
	settingRate  settingKind = iota // 倍率，没有设置时是 1
	settingFloat                    // 数字，只在设置了时显示
	settingInt
	settingBool
)

type portalSettingDef struct {
	file    cfgFile
	key     string
	label   i18n
	kind    settingKind
	always  bool   // 没有设置时也显示（倍率显示 1，布尔值显示 def）
	def     string // always 的布尔值、数字的默认值
	inverse bool   // 布尔值取反显示（ServerPVE → PvP）
}

type portalSettingGroup struct {
	title i18n
	items []portalSettingDef
}

// portalSettingGroups 玩家关心的服务器设置。键名和默认值依据 ark.wiki.gg 的 Server configuration（两个游戏相同的部分）
var portalSettingGroups = []portalSettingGroup{
	{i18n{"倍率", "Rates"}, []portalSettingDef{
		{file: cfgGUS, key: "XPMultiplier", label: i18n{"经验", "XP"}, kind: settingRate, always: true},
		{file: cfgGUS, key: "TamingSpeedMultiplier", label: i18n{"驯服速度", "Taming speed"}, kind: settingRate, always: true},
		{file: cfgGUS, key: "HarvestAmountMultiplier", label: i18n{"采集量", "Harvest amount"}, kind: settingRate, always: true},
		{file: cfgGUS, key: "ResourcesRespawnPeriodMultiplier", label: i18n{"资源刷新间隔", "Resource respawn period"}, kind: settingFloat},
		{file: cfgGUS, key: "DayCycleSpeedScale", label: i18n{"昼夜速度", "Day cycle speed"}, kind: settingFloat},
		{file: cfgGUS, key: "PlayerCharacterFoodDrainMultiplier", label: i18n{"玩家饥饿速度", "Player food drain"}, kind: settingFloat},
		{file: cfgGUS, key: "PlayerCharacterWaterDrainMultiplier", label: i18n{"玩家口渴速度", "Player water drain"}, kind: settingFloat},
		{file: cfgGame, key: "CropGrowthSpeedMultiplier", label: i18n{"作物生长速度", "Crop growth speed"}, kind: settingFloat},
	}},
	{i18n{"繁殖", "Breeding"}, []portalSettingDef{
		{file: cfgGame, key: "BabyMatureSpeedMultiplier", label: i18n{"幼崽成长速度", "Baby mature speed"}, kind: settingRate, always: true},
		{file: cfgGame, key: "EggHatchSpeedMultiplier", label: i18n{"孵蛋速度", "Egg hatch speed"}, kind: settingRate, always: true},
		{file: cfgGame, key: "MatingIntervalMultiplier", label: i18n{"交配间隔", "Mating interval"}, kind: settingFloat},
		{file: cfgGame, key: "BabyCuddleIntervalMultiplier", label: i18n{"印记间隔", "Imprint interval"}, kind: settingFloat},
		{file: cfgGame, key: "BabyImprintAmountMultiplier", label: i18n{"每次印记量", "Imprint amount"}, kind: settingFloat},
	}},
	{i18n{"游戏规则", "Gameplay"}, []portalSettingDef{
		{file: cfgGUS, key: "ServerPVE", label: i18n{"PvP", "PvP"}, kind: settingBool, always: true, def: "false", inverse: true},
		{file: cfgGUS, key: "MaxTamedDinos", label: i18n{"驯服恐龙上限", "Tamed dino limit"}, kind: settingInt},
		{file: cfgGame, key: "MaxNumberOfPlayersInTribe", label: i18n{"部落人数上限", "Max tribe size"}, kind: settingInt},
		{file: cfgGUS, key: "AllowThirdPersonPlayer", label: i18n{"第三人称", "Third person"}, kind: settingBool},
		{file: cfgGUS, key: "ServerCrosshair", label: i18n{"准星", "Crosshair"}, kind: settingBool},
		{file: cfgGUS, key: "ShowMapPlayerLocation", label: i18n{"地图显示自己的位置", "Show own position on map"}, kind: settingBool},
		{file: cfgGUS, key: "AllowFlyerCarryPvE", label: i18n{"PvE 飞行坐骑抓取", "Flyer carry in PvE"}, kind: settingBool},
		{file: cfgGUS, key: "AlwaysAllowStructurePickup", label: i18n{"随时拾取建筑", "Always allow structure pickup"}, kind: settingBool},
		{file: cfgGUS, key: "DisableStructureDecayPvE", label: i18n{"关闭 PvE 建筑腐烂", "No structure decay in PvE"}, kind: settingBool},
		{file: cfgGame, key: "bAllowUnlimitedRespecs", label: i18n{"无限洗点", "Unlimited respecs"}, kind: settingBool},
		{file: cfgGame, key: "bDisableStructurePlacementCollision", label: i18n{"建筑无碰撞放置", "No placement collision"}, kind: settingBool},
	}},
}

// buildSettingsPage 一台服务器的「服务器设置」页：每组一个标题和一个 kv 块，只列设置了的项（倍率等总是列出）
func buildSettingsPage(key string, lookup func(cfgFile, string) string) portalPage {
	page := portalPage{
		ID:        settingsPagePrefix + key,
		Kind:      pageKindSettings,
		Title:     i18n{"服务器设置", "Server settings"},
		ServerKey: key,
		Blocks:    []portalBlock{},
	}
	for gi, g := range portalSettingGroups {
		var items []any
		if gi == len(portalSettingGroups)-1 {
			// 野生恐龙最高等级：OverrideOfficialDifficulty × 30（没有设置时取决于地图，不显示）
			if d, err := strconv.ParseFloat(strings.TrimSpace(lookup(cfgGUS, "OverrideOfficialDifficulty")), 64); err == nil && d > 0 && !math.IsInf(d, 0) {
				items = append(items, portalKV{Label: i18n{"野生恐龙最高等级", "Max wild dino level"}, Value: int(math.Round(d * 30))})
			}
		}
		for _, def := range g.items {
			raw := strings.TrimSpace(lookup(def.file, def.key))
			if v, ok := settingValue(def, raw); ok {
				items = append(items, portalKV{Label: def.label, Value: v})
			}
		}
		if len(items) > 0 {
			page.Blocks = append(page.Blocks, portalBlock{Type: "heading", Text: g.title}, portalBlock{Type: "kv", Items: items})
		}
	}
	return page
}

func settingValue(def portalSettingDef, raw string) (any, bool) {
	switch def.kind {
	case settingRate:
		return parseRate(raw), true
	case settingFloat:
		f, err := strconv.ParseFloat(raw, 64)
		if err != nil || math.IsNaN(f) || math.IsInf(f, 0) {
			return nil, false
		}
		return f, true
	case settingInt:
		f, err := strconv.ParseFloat(raw, 64)
		if err != nil || math.IsNaN(f) || math.IsInf(f, 0) {
			return nil, false
		}
		return int64(f), true
	case settingBool:
		if raw == "" {
			if !def.always {
				return nil, false
			}
			raw = def.def
		}
		b := parseIniBool(raw)
		if def.inverse {
			b = !b
		}
		return b, true
	}
	return nil, false
}

// portalModList 服务器的 Mod，ASE 的带上创意工坊名称（查不到时为空，客户端显示 ID）
type portalMod struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

func portalModList(game string, ids []string) []portalMod {
	out := make([]portalMod, 0, len(ids))
	titles := map[string]string{}
	if game == utils.GameASE && len(ids) > 0 {
		list := ids
		if len(list) > 100 {
			list = list[:100]
		}
		for _, info := range lookupMods(list, false) {
			titles[info.ID] = info.Title
		}
	}
	for _, id := range ids {
		out = append(out, portalMod{ID: id, Title: titles[id]})
	}
	return out
}

// ---- 服主的页面（接口） ----

func normalizePage(req models.PortalPageRequest) (models.PortalPageRequest, error) {
	req.Title = strings.TrimSpace(req.Title)
	req.Body = strings.TrimSpace(strings.ReplaceAll(req.Body, "\r\n", "\n"))
	if req.Title == "" || utf8.RuneCountInString(req.Title) > maxPageTitleRunes || strings.ContainsAny(req.Title, "\r\n") ||
		utf8.RuneCountInString(req.Body) > maxPageBodyRunes {
		return req, ErrInvalidPage
	}
	return req, nil
}

// ListPortalPages 服主写的页面，按排序
func (s *ServerService) ListPortalPages() ([]models.PortalPage, error) {
	var list []models.PortalPage
	if err := database.DB.Order("sort, id").Find(&list).Error; err != nil {
		return nil, fmt.Errorf("读取页面失败: %w", err)
	}
	return list, nil
}

// CreatePortalPage 新建页面，排在最后
func (s *ServerService) CreatePortalPage(req models.PortalPageRequest) (*models.PortalPage, error) {
	req, err := normalizePage(req)
	if err != nil {
		return nil, err
	}
	var n int64
	if err := database.DB.Model(&models.PortalPage{}).Count(&n).Error; err != nil {
		return nil, fmt.Errorf("保存页面失败: %w", err)
	}
	if n >= maxPortalPages {
		return nil, ErrTooManyPages
	}
	var last models.PortalPage
	_ = database.DB.Order("sort desc").Limit(1).Find(&last).Error
	p := models.PortalPage{Title: req.Title, Body: req.Body, Sort: last.Sort + 1}
	if err := database.DB.Create(&p).Error; err != nil {
		return nil, fmt.Errorf("保存页面失败: %w", err)
	}
	portalManager.wake()
	return &p, nil
}

func findPortalPage(id string) (models.PortalPage, error) {
	var p models.PortalPage
	if err := database.DB.First(&p, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return p, ErrPageNotFound
		}
		return p, fmt.Errorf("读取页面失败: %w", err)
	}
	return p, nil
}

// UpdatePortalPage 修改页面的标题和内容
func (s *ServerService) UpdatePortalPage(id string, req models.PortalPageRequest) (*models.PortalPage, error) {
	req, err := normalizePage(req)
	if err != nil {
		return nil, err
	}
	p, err := findPortalPage(id)
	if err != nil {
		return nil, err
	}
	if err := database.DB.Model(&p).Updates(map[string]any{"title": req.Title, "body": req.Body}).Error; err != nil {
		return nil, fmt.Errorf("保存页面失败: %w", err)
	}
	portalManager.wake()
	return &p, nil
}

// MovePortalPage 把页面向前（-1）或向后（1）移动一位
func (s *ServerService) MovePortalPage(id string, delta int) error {
	list, err := s.ListPortalPages()
	if err != nil {
		return err
	}
	idx := -1
	for i, p := range list {
		if strconv.FormatUint(uint64(p.ID), 10) == id {
			idx = i
		}
	}
	if idx < 0 {
		return ErrPageNotFound
	}
	j := idx + delta
	if delta == 0 || j < 0 || j >= len(list) {
		return nil
	}
	list[idx], list[j] = list[j], list[idx]
	err = database.DB.Transaction(func(tx *gorm.DB) error {
		for i, p := range list {
			if err := tx.Model(&models.PortalPage{}).Where("id = ?", p.ID).Update("sort", i+1).Error; err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("保存页面顺序失败: %w", err)
	}
	portalManager.wake()
	return nil
}

// DeletePortalPage 删除页面
func (s *ServerService) DeletePortalPage(id string) error {
	res := database.DB.Delete(&models.PortalPage{}, "id = ?", id)
	if res.Error != nil {
		return fmt.Errorf("删除页面失败: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return ErrPageNotFound
	}
	portalManager.wake()
	return nil
}

// customPortalPages 快照里服主写的页面
func customPortalPages() []portalPage {
	var list []models.PortalPage
	if err := database.DB.Order("sort, id").Find(&list).Error; err != nil {
		return nil
	}
	out := make([]portalPage, 0, len(list))
	for _, p := range list {
		out = append(out, portalPage{
			ID:     strconv.FormatUint(uint64(p.ID), 10),
			Kind:   pageKindCustom,
			Title:  p.Title,
			Blocks: parsePageBody(p.Body),
		})
	}
	return out
}
