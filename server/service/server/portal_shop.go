package server

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"strings"
	"unicode/utf8"

	"ark-server-commander/models"
	"ark-server-commander/service/docker_manager"
)

// 玩家站的商店：从服务器的 ArkShop 插件配置（ArkApi/Plugins/ArkShop/config.json）整理出给玩家看的部分，
// 格式见 ArkPlayerPortal/docs/snapshot.md。Mysql 整段、General 里的设置一律不读；
// 命令类商品和礼包里的命令只取 DisplayAs，命令本身可能是以管理员身份执行的，不能公开

const (
	arkShopPluginName = "ArkShop"
	// 商品、礼包、回收项各自最多取这么多，描述等文本截断到这么长
	maxShopEntries     = 2000
	maxShopContents    = 100
	maxShopTextRunes   = 500
	maxShopShortRunes  = 200
	defaultBuyCmd      = "/buy"
	defaultKitCmd      = "/kit"
	defaultBuyKitCmd   = "/buykit"
	defaultPointsCmd   = "/points"
	defaultShopCmd     = "/shop"
	defaultSellCmd     = "/sell"
	shopCommandMaxSize = 32
)

type portalShop struct {
	UpdatedAt string             `json:"updated_at"`
	Commands  portalShopCommands `json:"commands"`
	Items     []portalShopItem   `json:"items"`
	Kits      []portalShopKit    `json:"kits"`
	Sell      []portalShopSell   `json:"sell"`
}

type portalShopCommands struct {
	Buy    string `json:"buy"`
	Kit    string `json:"kit"`
	BuyKit string `json:"buy_kit"`
	Points string `json:"points"`
	Shop   string `json:"shop"`
	Sell   string `json:"sell"`
}

type portalShopEntry struct {
	Blueprint     string `json:"blueprint"`
	Amount        int    `json:"amount"`
	Quality       int    `json:"quality"`
	BlueprintOnly bool   `json:"blueprint_only"`
}

type portalShopDino struct {
	Blueprint       string `json:"blueprint"`
	Level           int    `json:"level"`
	Neutered        bool   `json:"neutered"`
	Gender          string `json:"gender"`
	SaddleBlueprint string `json:"saddle_blueprint"`
}

type portalShopExperience struct {
	Amount float64 `json:"amount"`
	ToDino bool    `json:"to_dino"`
}

type portalShopItem struct {
	ID          string                `json:"id"`
	Type        string                `json:"type"`
	Description string                `json:"description"`
	Price       int                   `json:"price"`
	MinLevel    int                   `json:"min_level"`
	MaxLevel    int                   `json:"max_level"`
	Permissions []string              `json:"permissions"`
	Items       []portalShopEntry     `json:"items,omitempty"`
	Dinos       []portalShopDino      `json:"dinos,omitempty"`
	Beacon      string                `json:"beacon,omitempty"`
	Experience  *portalShopExperience `json:"experience,omitempty"`
	Engrams     []string              `json:"engrams,omitempty"`
	Commands    []string              `json:"commands,omitempty"`
}

type portalShopKit struct {
	Name          string            `json:"name"`
	Description   string            `json:"description"`
	Price         *int              `json:"price"`
	DefaultAmount int               `json:"default_amount"`
	MinLevel      int               `json:"min_level"`
	MaxLevel      int               `json:"max_level"`
	Permissions   []string          `json:"permissions"`
	OnlyFromSpawn bool              `json:"only_from_spawn"`
	Items         []portalShopEntry `json:"items"`
	Dinos         []portalShopDino  `json:"dinos"`
	Commands      []string          `json:"commands"`
}

type portalShopSell struct {
	ID          string `json:"id"`
	Description string `json:"description"`
	Price       int    `json:"price"`
	Amount      int    `json:"amount"`
	Blueprint   string `json:"blueprint"`
}

// arkShop* 是 ArkShop config.json 里读取的字段（只声明要用的，其余字段一律忽略）。
// 数字用 float64 接收：ArkShop 用 nlohmann::json 读取，整数、小数都能写
type arkShopItemEntry struct {
	Amount         float64 `json:"Amount"`
	Quality        float64 `json:"Quality"`
	ForceBlueprint bool    `json:"ForceBlueprint"`
	Blueprint      string  `json:"Blueprint"`
	DisplayAs      string  `json:"DisplayAs"` // command 类商品、礼包的 Commands
}

type arkShopDino struct {
	Level           float64 `json:"Level"`
	Blueprint       string  `json:"Blueprint"`
	Neutered        bool    `json:"Neutered"`
	Gender          string  `json:"Gender"`
	SaddleBlueprint string  `json:"SaddleBlueprint"`
}

type arkShopItem struct {
	Type        string             `json:"Type"`
	Description string             `json:"Description"`
	Price       float64            `json:"Price"`
	MinLevel    float64            `json:"MinLevel"`
	MaxLevel    float64            `json:"MaxLevel"`
	Permissions string             `json:"Permissions"`
	Items       []arkShopItemEntry `json:"Items"`
	ClassName   string             `json:"ClassName"`
	Amount      float64            `json:"Amount"`
	GiveToDino  bool               `json:"GiveToDino"`
	arkShopDino
}

type arkShopKit struct {
	Description   string             `json:"Description"`
	Price         *float64           `json:"Price"`
	DefaultAmount float64            `json:"DefaultAmount"`
	MinLevel      float64            `json:"MinLevel"`
	MaxLevel      float64            `json:"MaxLevel"`
	Permissions   string             `json:"Permissions"`
	OnlyFromSpawn bool               `json:"OnlyFromSpawn"`
	Items         []arkShopItemEntry `json:"Items"`
	Dinos         []arkShopDino      `json:"Dinos"`
	Commands      []arkShopItemEntry `json:"Commands"`
}

type arkShopSell struct {
	Description string  `json:"Description"`
	Price       float64 `json:"Price"`
	Amount      float64 `json:"Amount"`
	Blueprint   string  `json:"Blueprint"`
}

var errNoArkShop = errors.New("没有安装或已停用 ArkShop")

// readArkShopConfig 读取服务器的 ArkShop config.json（只看启用的插件目录，插件目录名不区分大小写）。
// 目录内容由 ARK 容器写入，用 os.Root 打开，符号链接不能指到插件目录外
func readArkShopConfig(serverID uint) ([]byte, error) {
	root, err := os.OpenRoot(docker_manager.ServerPluginsDir(serverID))
	if err != nil {
		return nil, errNoArkShop
	}
	defer root.Close()
	entries, err := fs.ReadDir(root.FS(), ".")
	if err != nil {
		return nil, errNoArkShop
	}
	for _, e := range entries {
		if !e.IsDir() || !strings.EqualFold(e.Name(), arkShopPluginName) {
			continue
		}
		f, err := root.Open(e.Name() + "/config.json")
		if err != nil {
			return nil, errNoArkShop
		}
		defer f.Close()
		data, err := io.ReadAll(io.LimitReader(f, maxPluginConfigSize+1))
		if err != nil {
			return nil, fmt.Errorf("读取 ArkShop 的 config.json 失败：%w", err)
		}
		if len(data) > maxPluginConfigSize {
			return nil, errors.New("ArkShop 的 config.json 太大")
		}
		return data, nil
	}
	return nil, errNoArkShop
}

// buildPortalShop 把 ArkShop 的 config.json 整理成玩家站的商店
func buildPortalShop(data []byte) (*portalShop, error) {
	data = bytes.TrimPrefix(data, []byte("\xEF\xBB\xBF"))
	var cfg map[string]json.RawMessage
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("ArkShop 的 config.json 不是有效的 JSON：%w", err)
	}
	shop := &portalShop{Items: []portalShopItem{}, Kits: []portalShopKit{}, Sell: []portalShopSell{}}

	var messages map[string]any
	_ = json.Unmarshal(cfg["Messages"], &messages)
	cmd := func(key, def string) string {
		if v, ok := messages[key].(string); ok {
			if v = strings.TrimSpace(v); v != "" && len(v) <= shopCommandMaxSize && !strings.ContainsAny(v, " \t\r\n") {
				return v
			}
		}
		return def
	}
	shop.Commands = portalShopCommands{
		Buy:    cmd("BuyCmd", defaultBuyCmd),
		Kit:    cmd("KitCmd", defaultKitCmd),
		BuyKit: cmd("BuyKitCmd", defaultBuyKitCmd),
		Points: cmd("PointsCmd", defaultPointsCmd),
		Shop:   cmd("ShopCmd", defaultShopCmd),
		Sell:   cmd("SellCmd", defaultSellCmd),
	}

	for _, kv := range orderedObject(cfg["ShopItems"]) {
		var it arkShopItem
		if json.Unmarshal(kv.value, &it) != nil || len(shop.Items) >= maxShopEntries {
			continue
		}
		shop.Items = append(shop.Items, convertShopItem(kv.key, it))
	}
	for _, kv := range orderedObject(cfg["Kits"]) {
		var k arkShopKit
		if json.Unmarshal(kv.value, &k) != nil || len(shop.Kits) >= maxShopEntries {
			continue
		}
		shop.Kits = append(shop.Kits, convertShopKit(kv.key, k))
	}
	for _, kv := range orderedObject(cfg["SellItems"]) {
		var s arkShopSell
		if json.Unmarshal(kv.value, &s) != nil || len(shop.Sell) >= maxShopEntries {
			continue
		}
		shop.Sell = append(shop.Sell, portalShopSell{
			ID:          shopText(kv.key, maxShopShortRunes),
			Description: shopText(s.Description, maxShopTextRunes),
			Price:       int(s.Price),
			Amount:      int(s.Amount),
			Blueprint:   shopText(s.Blueprint, maxShopTextRunes),
		})
	}
	return shop, nil
}

func convertShopItem(id string, it arkShopItem) portalShopItem {
	out := portalShopItem{
		ID:          shopText(id, maxShopShortRunes),
		Type:        shopText(strings.ToLower(strings.TrimSpace(it.Type)), 32),
		Description: shopText(it.Description, maxShopTextRunes),
		Price:       int(it.Price),
		MinLevel:    int(it.MinLevel),
		MaxLevel:    int(it.MaxLevel),
		Permissions: splitPermissions(it.Permissions),
	}
	switch out.Type {
	case "item":
		out.Items = shopEntries(it.Items)
	case "dino":
		out.Dinos = []portalShopDino{shopDino(it.arkShopDino)}
	case "beacon":
		out.Beacon = shopText(it.ClassName, maxShopShortRunes)
	case "experience":
		out.Experience = &portalShopExperience{Amount: it.Amount, ToDino: it.GiveToDino}
	case "unlockengram":
		for _, e := range it.Items {
			if len(out.Engrams) < maxShopContents && strings.TrimSpace(e.Blueprint) != "" {
				out.Engrams = append(out.Engrams, shopText(e.Blueprint, maxShopTextRunes))
			}
		}
	case "command":
		out.Commands = shopCommands(it.Items)
	}
	return out
}

func convertShopKit(name string, k arkShopKit) portalShopKit {
	out := portalShopKit{
		Name:          shopText(name, maxShopShortRunes),
		Description:   shopText(k.Description, maxShopTextRunes),
		DefaultAmount: int(k.DefaultAmount),
		MinLevel:      int(k.MinLevel),
		MaxLevel:      int(k.MaxLevel),
		Permissions:   splitPermissions(k.Permissions),
		OnlyFromSpawn: k.OnlyFromSpawn,
		Items:         shopEntries(k.Items),
		Dinos:         []portalShopDino{},
		Commands:      shopCommands(k.Commands),
	}
	// ArkShop：Price 没有设置或为 -1 时礼包不能用点数购买，只能用免费次数领取
	if k.Price != nil && *k.Price != -1 {
		p := int(*k.Price)
		out.Price = &p
	}
	for _, d := range k.Dinos {
		if len(out.Dinos) < maxShopContents {
			out.Dinos = append(out.Dinos, shopDino(d))
		}
	}
	return out
}

func shopEntries(in []arkShopItemEntry) []portalShopEntry {
	out := []portalShopEntry{}
	for _, e := range in {
		if len(out) >= maxShopContents || strings.TrimSpace(e.Blueprint) == "" {
			continue
		}
		out = append(out, portalShopEntry{
			Blueprint:     shopText(e.Blueprint, maxShopTextRunes),
			Amount:        int(e.Amount),
			Quality:       int(e.Quality),
			BlueprintOnly: e.ForceBlueprint,
		})
	}
	return out
}

// shopCommands 命令只取给玩家看的 DisplayAs，没有 DisplayAs 的不显示
func shopCommands(in []arkShopItemEntry) []string {
	out := []string{}
	for _, e := range in {
		if d := strings.TrimSpace(e.DisplayAs); d != "" && len(out) < maxShopContents {
			out = append(out, shopText(d, maxShopTextRunes))
		}
	}
	return out
}

func shopDino(d arkShopDino) portalShopDino {
	gender := strings.ToLower(strings.TrimSpace(d.Gender))
	if gender != "male" && gender != "female" && gender != "random" {
		gender = ""
	}
	return portalShopDino{
		Blueprint:       shopText(d.Blueprint, maxShopTextRunes),
		Level:           int(d.Level),
		Neutered:        d.Neutered,
		Gender:          gender,
		SaddleBlueprint: shopText(d.SaddleBlueprint, maxShopTextRunes),
	}
}

// splitPermissions ArkShop 的 Permissions 是逗号分隔的权限组
func splitPermissions(v string) []string {
	out := []string{}
	for p := range strings.SplitSeq(v, ",") {
		if p = strings.TrimSpace(p); p != "" && len(out) < 32 {
			out = append(out, shopText(p, 64))
		}
	}
	return out
}

// shopText 去掉控制字符（保留换行）和无效的 UTF-8，截断到 n 个字符
func shopText(s string, n int) string {
	s = strings.ToValidUTF8(s, "")
	s = strings.Map(func(r rune) rune {
		if r == '\n' || r >= 0x20 && r != 0x7f {
			return r
		}
		return -1
	}, strings.TrimSpace(s))
	if utf8.RuneCountInString(s) > n {
		s = string([]rune(s)[:n])
	}
	return s
}

type jsonKV struct {
	key   string
	value json.RawMessage
}

// orderedObject 按原文顺序取出 JSON 对象的键值（ShopItems、Kits、SellItems 的顺序就是服主排列的顺序）。
// 不是对象或格式不对时返回已经读到的部分
func orderedObject(raw json.RawMessage) []jsonKV {
	if len(raw) == 0 {
		return nil
	}
	dec := json.NewDecoder(bytes.NewReader(raw))
	if t, err := dec.Token(); err != nil || t != json.Delim('{') {
		return nil
	}
	var out []jsonKV
	for dec.More() {
		t, err := dec.Token()
		if err != nil {
			return out
		}
		key, _ := t.(string)
		var v json.RawMessage
		if err := dec.Decode(&v); err != nil {
			return out
		}
		out = append(out, jsonKV{key: key, value: v})
	}
	return out
}

// portalShopStatus 设置页显示的商店状态：有没有商店、读取失败的原因
func portalShopStatus(server models.Server) (bool, string) {
	if !server.ArkApi {
		return false, ""
	}
	data, err := readArkShopConfig(server.ID)
	if errors.Is(err, errNoArkShop) {
		return false, ""
	}
	if err == nil {
		_, err = buildPortalShop(data)
	}
	if err != nil {
		return false, err.Error()
	}
	return true, ""
}
