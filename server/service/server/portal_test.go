package server

import (
	"database/sql"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	"ark-server-commander/database"
	"ark-server-commander/models"
	"ark-server-commander/service/docker_manager"
	"ark-server-commander/utils"
)

const testArkShopConfig = "\xEF\xBB\xBF" + `{
  "Mysql": {"UseMysql": true, "MysqlHost": "db", "MysqlUser": "root", "MysqlPass": "s3cret", "MysqlDB": "shop", "MysqlPort": 3306},
  "General": {"Discord": {"URL": "https://discord.com/api/webhooks/secret"}, "DbPathOverride": "C:/secret"},
  "Kits": {
    "starter": {"DefaultAmount": 2, "Price": 90, "Description": "Starter", "OnlyFromSpawn": true,
      "Items": [{"Amount": 1, "Quality": 0, "ForceBlueprint": false, "Blueprint": "Blueprint'/Game/Saddle.Saddle'"}],
      "Dinos": [{"Level": 10, "Blueprint": "Blueprint'/Game/Carno.Carno'"}]},
    "vip": {"DefaultAmount": 1, "Description": "VIP", "Permissions": "Admins, Premiums",
      "Commands": [{"Command": "GiveEngrams"}, {"Command": "SetCheatPlayer true", "DisplayAs": "All engrams"}]},
    "cheap": {"Price": -1, "DefaultAmount": 0}
  },
  "ShopItems": {
    "zeta": {"Type": "item", "Description": "Ingots\u0007", "Price": 15, "Items": [{"Amount": 100, "Blueprint": "Blueprint'/Game/Ingot.Ingot'"}]},
    "alpha": {"Type": "dino", "Description": "Carno", "Level": 15, "Price": 50, "MinLevel": 10, "Neutered": true, "Gender": "Male",
      "SaddleBlueprint": "Blueprint'/Game/Saddle.Saddle'", "Blueprint": "Blueprint'/Game/Carno.Carno'"},
    "crate": {"Type": "beacon", "Description": "Crate", "Price": 100, "ClassName": "SupplyCrate_Level25_Double_C"},
    "exp": {"Type": "experience", "Description": "XP", "Price": 55, "Amount": 1000.5, "GiveToDino": true},
    "tek": {"Type": "unlockengram", "Price": 20, "Items": [{"Blueprint": "Blueprint'/Game/TekBoots.TekBoots'"}]},
    "admin": {"Type": "command", "Price": 1000, "Items": [{"Command": "Destroyall", "ExecuteAsAdmin": true}, {"Command": "GiveEngrams", "DisplayAs": "Engrams"}]},
    "weird": {"Type": "Mystery", "Price": 1}
  },
  "SellItems": {"stone": {"Type": "item", "Description": "Stone", "Price": 10, "Amount": 100, "Blueprint": "Blueprint'/Game/Stone.Stone'"}},
  "Messages": {"BuyCmd": "/购买", "KitCmd": "/kit with space", "PointsCmd": ""}
}`

func TestBuildPortalShop(t *testing.T) {
	shop, err := buildPortalShop([]byte(testArkShopConfig))
	if err != nil {
		t.Fatal(err)
	}
	raw, _ := json.Marshal(shop)
	for _, secret := range []string{"s3cret", "Mysql", "webhooks", "C:/secret", "Destroyall", "SetCheatPlayer", "ExecuteAsAdmin"} {
		if strings.Contains(string(raw), secret) {
			t.Errorf("商店里不应当出现 %q：%s", secret, raw)
		}
	}

	want := portalShopCommands{Buy: "/购买", Kit: "/kit", BuyKit: "/buykit", Points: "/points", Shop: "/shop", Sell: "/sell"}
	if shop.Commands != want {
		t.Errorf("命令名：%+v", shop.Commands)
	}

	var ids []string
	for _, it := range shop.Items {
		ids = append(ids, it.ID)
	}
	if !slices.Equal(ids, []string{"zeta", "alpha", "crate", "exp", "tek", "admin", "weird"}) {
		t.Fatalf("商品应当按原文顺序：%v", ids)
	}
	byID := map[string]portalShopItem{}
	for _, it := range shop.Items {
		byID[it.ID] = it
	}
	if it := byID["zeta"]; it.Description != "Ingots" || len(it.Items) != 1 || it.Items[0].Amount != 100 || it.Dinos != nil {
		t.Errorf("item：%+v", it)
	}
	if it := byID["alpha"]; len(it.Dinos) != 1 || it.Dinos[0].Gender != "male" || !it.Dinos[0].Neutered || it.Dinos[0].Level != 15 || it.MinLevel != 10 {
		t.Errorf("dino：%+v", it)
	}
	if byID["crate"].Beacon != "SupplyCrate_Level25_Double_C" {
		t.Errorf("beacon：%+v", byID["crate"])
	}
	if e := byID["exp"].Experience; e == nil || e.Amount != 1000.5 || !e.ToDino {
		t.Errorf("experience：%+v", e)
	}
	if !slices.Equal(byID["tek"].Engrams, []string{"Blueprint'/Game/TekBoots.TekBoots'"}) {
		t.Errorf("unlockengram：%+v", byID["tek"])
	}
	if !slices.Equal(byID["admin"].Commands, []string{"Engrams"}) {
		t.Errorf("command 只应当有 DisplayAs：%+v", byID["admin"].Commands)
	}
	if it := byID["weird"]; it.Type != "mystery" || it.Items != nil || it.Commands != nil {
		t.Errorf("不认识的类型只写 type：%+v", it)
	}

	if len(shop.Kits) != 3 {
		t.Fatalf("礼包：%+v", shop.Kits)
	}
	starter, vip, cheap := shop.Kits[0], shop.Kits[1], shop.Kits[2]
	if starter.Price == nil || *starter.Price != 90 || !starter.OnlyFromSpawn || len(starter.Dinos) != 1 || len(starter.Items) != 1 {
		t.Errorf("starter：%+v", starter)
	}
	if vip.Price != nil || !slices.Equal(vip.Permissions, []string{"Admins", "Premiums"}) || !slices.Equal(vip.Commands, []string{"All engrams"}) {
		t.Errorf("vip：%+v", vip)
	}
	if cheap.Price != nil {
		t.Errorf("Price 为 -1 时不能购买：%+v", cheap)
	}
	if len(shop.Sell) != 1 || shop.Sell[0].Amount != 100 || shop.Sell[0].Price != 10 {
		t.Errorf("回收：%+v", shop.Sell)
	}

	if _, err := buildPortalShop([]byte("{bad")); err == nil {
		t.Error("无效的 JSON 应当报错")
	}
	empty, err := buildPortalShop([]byte(`{}`))
	if err != nil || empty.Items == nil || empty.Kits == nil || empty.Sell == nil {
		t.Errorf("空配置应当返回空数组：%+v %v", empty, err)
	}
}

func TestReadArkShopConfig(t *testing.T) {
	setupTestDB(t)
	useTempDataDir(t)
	srv := createTestServer(t, 7777, 27015, 32330)
	if _, err := readArkShopConfig(srv.ID); !errors.Is(err, errNoArkShop) {
		t.Fatalf("没有插件目录时应当是 errNoArkShop：%v", err)
	}
	dir := filepath.Join(docker_manager.ServerPluginsDir(srv.ID), "arkshop") // 目录名不区分大小写
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if _, err := readArkShopConfig(srv.ID); !errors.Is(err, errNoArkShop) {
		t.Fatalf("没有 config.json 时应当是 errNoArkShop：%v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "config.json"), []byte(`{"ShopItems":{}}`), 0o644); err != nil {
		t.Fatal(err)
	}
	if data, err := readArkShopConfig(srv.ID); err != nil || string(data) != `{"ShopItems":{}}` {
		t.Fatalf("读取失败：%q %v", data, err)
	}
}

// fakePortalDocker 记录玩家站的 Docker 操作
type fakePortalDocker struct {
	ready   bool
	imageID string
	info    docker_manager.LinkContainerInfo
	created []docker_manager.PortalContainerSpec
	removed int
}

func (f *fakePortalDocker) ImageReady(string) (bool, error) { return f.ready, nil }
func (f *fakePortalDocker) ImageID(string) (string, error)  { return f.imageID, nil }
func (f *fakePortalDocker) Inspect(string) (docker_manager.LinkContainerInfo, error) {
	return f.info, nil
}
func (f *fakePortalDocker) Create(spec docker_manager.PortalContainerSpec) error {
	f.created = append(f.created, spec)
	f.info = docker_manager.LinkContainerInfo{Exists: true, Running: true, Labels: spec.Labels}
	return nil
}
func (f *fakePortalDocker) Remove() error {
	f.removed++
	f.info = docker_manager.LinkContainerInfo{}
	return nil
}

func newTestPortalManager(t *testing.T) (*PortalManager, *fakePortalDocker, string) {
	t.Helper()
	dir := filepath.Join(t.TempDir(), "portal")
	fd := &fakePortalDocker{ready: true, imageID: "sha256:aaa"}
	m := &PortalManager{
		status:   models.PortalStatus{State: models.PortalDisabled},
		dir:      func() string { return dir },
		docker:   fd,
		now:      func() time.Time { return time.Date(2026, 11, 27, 12, 0, 0, 0, time.UTC) },
		wakeCh:   make(chan struct{}, 1),
		written:  map[string]string{},
		versions: map[uint]string{},
		players:  map[uint]playersCache{},
	}
	return m, fd, dir
}

func enablePortal(t *testing.T, req models.PortalSettingsRequest) {
	t.Helper()
	req.Enabled = true
	if _, err := (&ServerService{}).UpdatePortal(req); err != nil {
		t.Fatal(err)
	}
}

func readSnapshot(t *testing.T, dir string) portalSnapshot {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(dir, "portal.json"))
	if err != nil {
		t.Fatal(err)
	}
	var snap portalSnapshot
	if err := json.Unmarshal(data, &snap); err != nil {
		t.Fatal(err)
	}
	return snap
}

func TestPortalSnapshotAndContainer(t *testing.T) {
	setupTestDB(t)
	useTempDataDir(t)
	config := models.PortalSettingsRequest{Port: 8090, PublicURL: "https://ark.example.com", Name: "示例方舟", QQGroup: "123456",
		Links: []models.PortalLink{{Title: "规则", URL: "https://example.com/rules"}}}
	a := createTestServer(t, 7777, 27015, 32330)
	b := createTestServer(t, 7779, 27017, 32332)
	c := createTestServer(t, 7781, 27019, 32334)
	database.DB.Model(&a).Updates(map[string]any{"arkapi": true, "cluster_id": "c1", "published": true})
	database.DB.Model(&b).Updates(map[string]any{"cluster_id": "c1"})
	config.HiddenServers = []uint{c.ID}

	shopDir := filepath.Join(docker_manager.ServerPluginsDir(a.ID), "ArkShop")
	if err := os.MkdirAll(shopDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(shopDir, "config.json"), []byte(testArkShopConfig), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := (&ServerService{}).CreateAnnouncement(models.PortalAnnouncementRequest{Title: "维护", Body: "周三\r\n停服", Pinned: true}); err != nil {
		t.Fatal(err)
	}

	m, fd, dir := newTestPortalManager(t)

	// 没有开启：不写快照、删除容器
	m.tick()
	if _, err := os.Stat(dir); !os.IsNotExist(err) || fd.removed != 1 || m.Status().State != models.PortalDisabled {
		t.Fatalf("没有开启时不应当有快照和容器：%v %+v", err, m.Status())
	}

	enablePortal(t, config)
	m.tick()
	if len(fd.created) != 1 || fd.created[0].HostPort != 8090 || fd.created[0].Image != utils.PortalImage {
		t.Fatalf("应当创建玩家站容器：%+v", fd.created)
	}
	if st := m.Status(); st.State != models.PortalRunning || st.URL != "https://ark.example.com" {
		t.Fatalf("状态：%+v", st)
	}

	snap := readSnapshot(t, dir)
	if snap.Site.Name != "示例方舟" || snap.Site.QQGroup != "123456" || len(snap.Site.Links) != 1 {
		t.Errorf("站点资料：%+v", snap.Site)
	}
	if len(snap.Announcements) != 1 || snap.Announcements[0].Body != "周三\n停服" || !snap.Announcements[0].Pinned {
		t.Errorf("公告：%+v", snap.Announcements)
	}
	if len(snap.Servers) != 2 {
		t.Fatalf("隐藏的服务器不应当出现：%+v", snap.Servers)
	}
	sa, sb := snap.Servers[0], snap.Servers[1]
	if !sa.Shop || sb.Shop || sa.ClusterGroup != 1 || sb.ClusterGroup != 1 || sa.Status != "offline" || sa.Mods == nil {
		t.Errorf("服务器：%+v %+v", sa, sb)
	}
	if sb.Address != "" || sb.RoomID != 0 {
		t.Errorf("没有发布的服务器不应当有地址：%+v", sb)
	}
	raw, _ := os.ReadFile(filepath.Join(dir, "portal.json"))
	for _, secret := range []string{"password\":\"", "admin_password", "c1", "s3cret"} {
		if strings.Contains(string(raw), secret) {
			t.Errorf("portal.json 不应当包含 %q", secret)
		}
	}
	shopRaw, err := os.ReadFile(filepath.Join(dir, "shops", sa.Key+".json"))
	if err != nil || strings.Contains(string(shopRaw), "s3cret") || !strings.Contains(string(shopRaw), "/购买") {
		t.Fatalf("商店文件：%s %v", shopRaw, err)
	}

	// 没有变化：不重写，不重建容器
	info, _ := os.Stat(filepath.Join(dir, "portal.json"))
	m.now = func() time.Time { return time.Date(2026, 11, 27, 13, 0, 0, 0, time.UTC) }
	m.tick()
	if again, _ := os.Stat(filepath.Join(dir, "portal.json")); !again.ModTime().Equal(info.ModTime()) || readSnapshot(t, dir).UpdatedAt != snap.UpdatedAt {
		t.Error("内容没变时不应当重写")
	}
	if len(fd.created) != 1 {
		t.Errorf("参数没变时不应当重建容器：%d", len(fd.created))
	}

	// 停用 ArkShop：删除商店文件；改端口：重建容器
	if err := os.RemoveAll(shopDir); err != nil {
		t.Fatal(err)
	}
	config.Port = 8091
	enablePortal(t, config)
	m.tick()
	if _, err := os.Stat(filepath.Join(dir, "shops", sa.Key+".json")); !os.IsNotExist(err) {
		t.Error("ArkShop 不在了应当删除商店文件")
	}
	if readSnapshot(t, dir).Servers[0].Shop {
		t.Error("shop 应当变成 false")
	}
	if len(fd.created) != 2 || fd.created[1].HostPort != 8091 {
		t.Errorf("端口变了应当重建：%+v", fd.created)
	}

	// 镜像还没下载好
	fd.ready = false
	fd.info = docker_manager.LinkContainerInfo{}
	m.tick()
	if m.Status().State != models.PortalStarting {
		t.Errorf("镜像没下载好时应当是 starting：%+v", m.Status())
	}

	// 关闭：删除容器和快照
	config.Port = 8091
	if _, err := (&ServerService{}).UpdatePortal(models.PortalSettingsRequest{Port: 8091}); err != nil {
		t.Fatal(err)
	}
	m.tick()
	if _, err := os.Stat(dir); !os.IsNotExist(err) || m.Status().State != models.PortalDisabled {
		t.Errorf("关闭后应当删除快照：%v %+v", err, m.Status())
	}
}

func TestPortalSettingsValidation(t *testing.T) {
	setupTestDB(t)
	createTestServer(t, 7777, 27015, 32330)
	s := &ServerService{}
	bad := []models.PortalSettingsRequest{
		{Port: 80},
		{Port: 8090, QQGroup: "abc"},
		{Port: 8090, Discord: "javascript:alert(1)"},
		{Port: 8090, PublicURL: "ftp://x"},
		{Port: 8090, Links: []models.PortalLink{{Title: "x", URL: "javascript:alert(1)"}}},
		{Port: 8090, Name: strings.Repeat("名", 65)},
	}
	for _, req := range bad {
		if _, err := s.UpdatePortal(req); !errors.Is(err, ErrInvalidPortalSettings) {
			t.Errorf("%+v 应当是 ErrInvalidPortalSettings：%v", req, err)
		}
	}
	// 开启时端口不能和服务器冲突
	if _, err := s.UpdatePortal(models.PortalSettingsRequest{Enabled: true, Port: 7778}); !errors.Is(err, ErrPortConflict) {
		t.Errorf("和游戏端口+1 冲突应当报错：%v", err)
	}
	info, err := s.UpdatePortal(models.PortalSettingsRequest{Enabled: true, Port: 8090, PublicURL: "https://ark.example.com/",
		Links: []models.PortalLink{{}, {Title: "规则", URL: "https://example.com"}}, HiddenServers: []uint{1}})
	if err != nil {
		t.Fatal(err)
	}
	if info.Settings.PublicURL != "https://ark.example.com" || len(info.Links) != 1 || !info.Servers[0].Hidden {
		t.Errorf("保存结果：%+v", info)
	}
	// 开着玩家站时，新建服务器不能占用它的端口
	if err := checkPorts(0, 8090, 27100, 32400); !errors.Is(err, ErrPortConflict) {
		t.Errorf("新服务器占用玩家站端口应当报错：%v", err)
	}
	if got := portalLinkTCPPort(); got != "8090" {
		t.Errorf("联机房节点应当放行玩家站端口：%s", got)
	}
	// 登记联机房带端口；发布到收录站带地址，隐藏的服务器不带
	if p := buildLinkRoomProfile(models.LinkRoom{}, nil, nil); p.PortalPort != 8090 {
		t.Errorf("登记联机房应当带玩家站端口：%d", p.PortalPort)
	}
	var srv models.Server
	database.DB.First(&srv, 1)
	if got := buildPublishProfile(srv).PortalURL; got != "" {
		t.Errorf("隐藏的服务器不应当带玩家站地址：%q", got)
	}
	if _, err := s.UpdatePortal(models.PortalSettingsRequest{Enabled: true, Port: 8090, PublicURL: "https://ark.example.com"}); err != nil {
		t.Fatal(err)
	}
	if got := buildPublishProfile(srv).PortalURL; got != "https://ark.example.com" {
		t.Errorf("发布应当带玩家站地址：%q", got)
	}
	if _, err := s.UpdatePortal(models.PortalSettingsRequest{Port: 8090}); err != nil {
		t.Fatal(err)
	}
	if buildPublishProfile(srv).PortalURL != "" || portalPortForLink() != 0 || portalLinkTCPPort() != "0" {
		t.Error("关闭玩家站后不应当带地址和端口")
	}
}

func TestAnnouncements(t *testing.T) {
	setupTestDB(t)
	s := &ServerService{}
	if _, err := s.CreateAnnouncement(models.PortalAnnouncementRequest{Title: "  "}); !errors.Is(err, ErrInvalidAnnouncement) {
		t.Errorf("空标题应当报错：%v", err)
	}
	if _, err := s.CreateAnnouncement(models.PortalAnnouncementRequest{Title: "a\nb"}); !errors.Is(err, ErrInvalidAnnouncement) {
		t.Errorf("标题有换行应当报错：%v", err)
	}
	a, err := s.CreateAnnouncement(models.PortalAnnouncementRequest{Title: "开服", Body: "欢迎"})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := s.UpdateAnnouncement("999", models.PortalAnnouncementRequest{Title: "x"}); !errors.Is(err, ErrAnnouncementNotFound) {
		t.Errorf("不存在的公告：%v", err)
	}
	id := jsonID(a.ID)
	if _, err := s.UpdateAnnouncement(id, models.PortalAnnouncementRequest{Title: "开服啦", Pinned: true}); err != nil {
		t.Fatal(err)
	}
	list, _ := s.ListAnnouncements()
	if len(list) != 1 || list[0].Title != "开服啦" || !list[0].Pinned || list[0].Body != "" {
		t.Errorf("修改结果：%+v", list)
	}
	if err := s.DeleteAnnouncement(id); err != nil {
		t.Fatal(err)
	}
	if err := s.DeleteAnnouncement(id); !errors.Is(err, ErrAnnouncementNotFound) {
		t.Errorf("重复删除：%v", err)
	}
}

func jsonID(id uint) string {
	b, _ := json.Marshal(id)
	return string(b)
}

func TestParsePageBody(t *testing.T) {
	blocks := parsePageBody("# 服规\r\n不许偷家。\n第二行\n\n- 禁止外挂\n- 禁止辱骂\n结尾<script>")
	raw, _ := json.Marshal(blocks)
	want := `[{"type":"heading","text":"服规"},{"type":"paragraph","text":"不许偷家。\n第二行"},{"type":"list","items":["禁止外挂","禁止辱骂"]},{"type":"paragraph","text":"结尾\u003cscript\u003e"}]`
	if string(raw) != want {
		t.Fatalf("blocks:\n%s\nwant\n%s", raw, want)
	}
	if b := parsePageBody(""); len(b) != 0 {
		t.Errorf("空内容应当没有块：%+v", b)
	}
}

func TestBuildSettingsPage(t *testing.T) {
	values := map[string]string{
		"XPMultiplier": "3", "OverrideOfficialDifficulty": "5.0", "ServerPVE": "True", "MaxNumberOfPlayersInTribe": "6",
		"AllowThirdPersonPlayer": "False", "DayCycleSpeedScale": "abc",
	}
	page := buildSettingsPage("3", func(_ cfgFile, key string) string { return values[key] })
	if page.ID != "settings-3" || page.Kind != pageKindSettings || page.ServerKey != "3" {
		t.Fatalf("page: %+v", page)
	}
	got := map[string]any{}
	for _, b := range page.Blocks {
		if b.Type != "kv" {
			continue
		}
		for _, it := range b.Items {
			kv := it.(portalKV)
			got[kv.Label.(i18n).En] = kv.Value
		}
	}
	for label, want := range map[string]any{
		"XP": 3.0, "Taming speed": 1.0, "Baby mature speed": 1.0, "Max wild dino level": 150, "PvP": false,
		"Max tribe size": int64(6), "Third person": false,
	} {
		if got[label] != want {
			t.Errorf("%s = %v (%T), want %v", label, got[label], got[label], want)
		}
	}
	if _, ok := got["Day cycle speed"]; ok {
		t.Error("无效的数字不应当显示")
	}
	if _, ok := got["Crosshair"]; ok {
		t.Error("没有设置的开关不应当显示")
	}
}

// writeArkShopDB 用 ArkShop 的表结构建一个 SQLite 库
func writeArkShopDB(t *testing.T, path string) {
	t.Helper()
	db, err := sql.Open("sqlite", path)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	for _, q := range []string{
		"CREATE TABLE Players (Id INTEGER PRIMARY KEY AUTOINCREMENT, SteamId INTEGER DEFAULT 0, Kits TEXT DEFAULT '{}', Points INTEGER DEFAULT 0, TotalSpent INTEGER DEFAULT 0)",
		`INSERT INTO Players (SteamId, Kits, Points, TotalSpent) VALUES (76561198000000001, '{"starter":{"Amount":2}}', 120, 30)`,
		"INSERT INTO Players (SteamId, Kits, Points, TotalSpent) VALUES (12, '{}', 5, 0)",
	} {
		if _, err := db.Exec(q); err != nil {
			t.Fatal(err)
		}
	}
}

func TestPortalPlayersAndIdentity(t *testing.T) {
	setupTestDB(t)
	useTempDataDir(t)
	srv := createTestServer(t, 7777, 27015, 32330)
	database.DB.Model(&srv).Update("arkapi", true)
	shopDir := filepath.Join(docker_manager.ServerPluginsDir(srv.ID), "ArkShop")
	if err := os.MkdirAll(shopDir, 0o755); err != nil {
		t.Fatal(err)
	}
	config := `{"Mysql":{"UseMysql":false},"General":{"DbPathOverride":""},"ShopItems":{}}`
	if err := os.WriteFile(filepath.Join(shopDir, "config.json"), []byte(config), 0o644); err != nil {
		t.Fatal(err)
	}
	writeArkShopDB(t, filepath.Join(shopDir, "ArkShop.db"))
	if _, err := (&ServerService{}).CreatePortalPage(models.PortalPageRequest{Title: "服规", Body: "- 禁止外挂"}); err != nil {
		t.Fatal(err)
	}
	enablePortal(t, models.PortalSettingsRequest{Port: 8090})

	m, fd, dir := newTestPortalManager(t)
	m.players = map[uint]playersCache{}
	m.tick()

	data, err := os.ReadFile(filepath.Join(dir, "players", "1.json"))
	if err != nil {
		t.Fatal(err)
	}
	var players portalPlayers
	if err := json.Unmarshal(data, &players); err != nil {
		t.Fatal(err)
	}
	p, ok := players.Players["76561198000000001"]
	if !ok || p.Points != 120 || p.TotalSpent != 30 || p.Kits["starter"] != 2 || len(players.Players) != 1 {
		t.Fatalf("players: %+v", players)
	}
	snap := readSnapshot(t, dir)
	if !snap.Servers[0].Points || snap.Servers[0].SettingsPage != "settings-1" {
		t.Errorf("server: %+v", snap.Servers[0])
	}
	if len(snap.Pages) != 2 || snap.Pages[0].Kind != pageKindCustom || snap.Pages[1].ID != "settings-1" {
		t.Errorf("pages: %+v", snap.Pages)
	}

	// 身份密钥生成后不再变
	key1, err := os.ReadFile(filepath.Join(dir, identityKeyFile))
	if err != nil || len(key1) == 0 {
		t.Fatalf("identity key: %v", err)
	}
	m.tick()
	if key2, _ := os.ReadFile(filepath.Join(dir, identityKeyFile)); string(key2) != string(key1) {
		t.Error("身份密钥不应当变化")
	}
	if env := strings.Join(fd.created[0].Env, " "); !strings.Contains(env, "PORTAL_ARKFGF_KEYS_URL=") || !strings.Contains(env, "PORTAL_PUBLIC_URL=") {
		t.Errorf("container env: %v", fd.created[0].Env)
	}

	// 改成 MySQL：不再读点数，删除点数文件
	config = `{"Mysql":{"UseMysql":true},"ShopItems":{}}`
	if err := os.WriteFile(filepath.Join(shopDir, "config.json"), []byte(config), 0o644); err != nil {
		t.Fatal(err)
	}
	m.tick()
	if _, err := os.Stat(filepath.Join(dir, "players", "1.json")); !os.IsNotExist(err) {
		t.Error("用 MySQL 时应当删除点数文件")
	}
}

func TestPortalPagesCRUD(t *testing.T) {
	setupTestDB(t)
	s := &ServerService{}
	if _, err := s.CreatePortalPage(models.PortalPageRequest{Title: ""}); !errors.Is(err, ErrInvalidPage) {
		t.Errorf("空标题：%v", err)
	}
	a, _ := s.CreatePortalPage(models.PortalPageRequest{Title: "A"})
	b, _ := s.CreatePortalPage(models.PortalPageRequest{Title: "B"})
	if err := s.MovePortalPage(jsonID(b.ID), -1); err != nil {
		t.Fatal(err)
	}
	list, _ := s.ListPortalPages()
	if len(list) != 2 || list[0].Title != "B" || list[1].Title != "A" {
		t.Fatalf("顺序：%+v", list)
	}
	if _, err := s.UpdatePortalPage(jsonID(a.ID), models.PortalPageRequest{Title: "A2", Body: "x"}); err != nil {
		t.Fatal(err)
	}
	if err := s.DeletePortalPage(jsonID(a.ID)); err != nil {
		t.Fatal(err)
	}
	if err := s.DeletePortalPage(jsonID(a.ID)); !errors.Is(err, ErrPageNotFound) {
		t.Errorf("重复删除：%v", err)
	}
}
