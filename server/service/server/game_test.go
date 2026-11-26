package server

import (
	"errors"
	"fmt"
	"slices"
	"testing"

	"ark-server-commander/database"
	"ark-server-commander/models"
	"ark-server-commander/utils"
)

func TestImageForGame(t *testing.T) {
	cases := []struct {
		game   string
		arkApi bool
		want   string
	}{
		{"", false, utils.ArkServerLinuxImage},
		{utils.GameASE, true, utils.ArkServerProtonImage},
		{utils.GameASA, false, utils.AsaServerWineImage},
		{utils.GameASA, true, utils.AsaServerProtonImage},
	}
	for _, c := range cases {
		image := utils.ArkServerImageFor(c.game, c.arkApi)
		if image != c.want {
			t.Errorf("ArkServerImageFor(%q, %v) = %s, want %s", c.game, c.arkApi, image, c.want)
		}
		if utils.ImageGame(image) != utils.NormalizeGame(c.game) || utils.ImageUsesArkApi(image) != c.arkApi {
			t.Errorf("image %s: game %q arkapi %v", image, utils.ImageGame(image), utils.ImageUsesArkApi(image))
		}
	}
	// 服务器用到的每个镜像都要能在首页镜像管理里拉取和更新
	for _, image := range utils.ArkServerImages {
		if !slices.Contains(managedImages, image) {
			t.Errorf("managedImages 缺少 %s", image)
		}
	}
	if utils.ArkAppUpdateFor(utils.GameASA, "") != utils.AsaServerAppID {
		t.Errorf("ASA APP_UPDATE = %q", utils.ArkAppUpdateFor(utils.GameASA, ""))
	}
	if utils.ArkAppUpdateFor("", utils.GameBranchPreAquatica) != "376030 -beta preaquatica" {
		t.Errorf("ASE preaquatica APP_UPDATE = %q", utils.ArkAppUpdateFor("", utils.GameBranchPreAquatica))
	}
	// ASA 只有正式版
	if utils.ValidGameBranch(utils.GameASA, utils.GameBranchPreAquatica) || !utils.ValidGameBranch(utils.GameASA, "") {
		t.Error("ASA must only accept the public branch")
	}
	if utils.ValidGame("asx") || !utils.ValidGame("") {
		t.Error("ValidGame")
	}
}

// ASA 没有查询端口：不占用端口，也不和别的服务器的查询端口冲突
func TestCheckPortsWithoutQueryPort(t *testing.T) {
	setupTestDB(t)
	createTestServer(t, 7777, 27015, 32330)
	asa := createTestServer(t, 7790, 0, 32340)
	if err := checkPorts(0, 7792, 0, 32341); err != nil {
		t.Fatalf("a second ASA server without query port should not conflict: %v", err)
	}
	if err := checkPorts(0, 7794, 27015, 32342); !errors.Is(err, ErrPortConflict) {
		t.Fatalf("expected query port conflict with the ASE server, got %v", err)
	}
	if err := checkPorts(asa.ID, 7790, 0, 32340); err != nil {
		t.Fatalf("editing the ASA server itself: %v", err)
	}
}

func TestServerDefaultsPerGame(t *testing.T) {
	s := NewServerService()
	asa, err := s.ServerDefaults(utils.GameASA)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := asa.ServerArgs.CommandLineArgs["newsaveformat"]; ok {
		t.Error("ASA defaults must not carry ASE-only args")
	}
	ase, err := s.ServerDefaults("")
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := ase.ServerArgs.CommandLineArgs["newsaveformat"]; !ok {
		t.Error("ASE defaults lost their args")
	}
	if _, err := s.ServerDefaults("x"); !errors.Is(err, ErrInvalidGame) {
		t.Fatalf("expected ErrInvalidGame, got %v", err)
	}
	if DefaultMapFor(utils.GameASA) != "TheIsland_WP" || DefaultMapFor("") != "TheIsland" {
		t.Error("DefaultMapFor")
	}
}

// 联机节点只为 ASA 服务器放行游戏端口（ASE 是游戏端口、+1 和查询端口）
func TestLinkNodePortsASA(t *testing.T) {
	ports := linkNodePorts([]models.Server{
		{ID: 1, Port: 7777, QueryPort: 27015},
		{ID: 2, Game: utils.GameASA, Port: 7790},
	})
	want := []string{"7777", "7778", "27015", "7790"}
	if !slices.Equal(ports, want) {
		t.Fatalf("linkNodePorts = %v, want %v", ports, want)
	}
	if p := buildLinkServerProfile(models.Server{ID: 2, Game: utils.GameASA, Port: 7790}); p.Game != utils.GameASA || p.QueryPort != 0 {
		t.Fatalf("ASA link profile = %+v", p)
	}
	if p := buildPublishProfile(models.Server{ID: 1, Port: 7777, QueryPort: 27015}); p.Game != utils.GameASE {
		t.Fatalf("old servers must publish as ase, got %q", p.Game)
	}
}

func TestAsaServerModsHaveNoUpdateInfo(t *testing.T) {
	mods := asaServerMods(" 123, 456 ,")
	if len(mods) != 2 || mods[0].ID != "123" || !mods[1].Installed || mods[1].HasUpdate {
		t.Fatalf("asaServerMods = %+v", mods)
	}
}

// 跨游戏导入配置时地图、Mod 和启动参数不导入
func TestValidateServerConfigGame(t *testing.T) {
	cfg := &models.ServerConfig{Format: models.ServerConfigFormat, Version: models.ServerConfigVersion, Game: "xbox"}
	if err := validateServerConfig(cfg); !errors.Is(err, ErrInvalidServerConfig) {
		t.Fatalf("expected invalid config for unknown game, got %v", err)
	}
	cfg.Game = ""
	if err := validateServerConfig(cfg); err != nil {
		t.Fatalf("old exports without game must stay valid: %v", err)
	}
}

func TestImportConfigAcrossGames(t *testing.T) {
	setupTestDB(t)
	useTempDataDir(t)
	s := NewServerService()
	dst := createTestServer(t, 7790, 0, 32340)
	database.DB.Model(&dst).Updates(map[string]any{"game": utils.GameASA, "map": "TheIsland_WP", "game_mod_ids": "900"})

	// ASE 导出的配置导入到 ASA 服务器：地图和 Mod 不导入，集群和选项照常
	cfg := &models.ServerConfig{Format: models.ServerConfigFormat, Version: models.ServerConfigVersion,
		Map: "Ragnarok", GameModIds: "111", ClusterID: "c1", MaxPlayers: 20}
	if _, err := s.ImportConfig(1, fmt.Sprint(dst.ID), cfg, []string{models.ConfigPartMap, models.ConfigPartMods, models.ConfigPartCluster, models.ConfigPartOptions}); err != nil {
		t.Fatal(err)
	}
	var updated models.Server
	database.DB.First(&updated, dst.ID)
	if updated.Map != "TheIsland_WP" || updated.GameModIds != "900" || updated.ClusterID != "c1" || updated.MaxPlayers != 20 {
		t.Fatalf("after cross-game import: %+v", updated)
	}

	// 同一个游戏照常导入
	cfg.Game, cfg.Map, cfg.GameModIds = utils.GameASA, "Ragnarok_WP", "901"
	if _, err := s.ImportConfig(1, fmt.Sprint(dst.ID), cfg, []string{models.ConfigPartMap, models.ConfigPartMods}); err != nil {
		t.Fatal(err)
	}
	updated = models.Server{}
	database.DB.First(&updated, dst.ID)
	if updated.Map != "Ragnarok_WP" || updated.GameModIds != "901" {
		t.Fatalf("after same-game import: %+v", updated)
	}
}

// 创建时的校验（在创建卷之前）：ASE 必须有查询端口，ASA 只有正式版。
// 创建本身要用真实的 Docker 卷，ASA 的查询端口写回 0 在本机管理端上验证过
func TestCreateServerGameValidation(t *testing.T) {
	setupTestDB(t)
	s := NewServerService()
	if _, err := s.CreateServer(1, models.ServerRequest{Port: 7792, RCONPort: 27091, AdminPassword: "pw", MaxPlayers: 10}); !errors.Is(err, ErrInvalidPort) {
		t.Fatalf("expected ErrInvalidPort for ASE without query port, got %v", err)
	}
	if _, err := s.CreateServer(1, models.ServerRequest{Game: utils.GameASA, GameBranch: utils.GameBranchPreAquatica, Port: 7794, RCONPort: 27092, AdminPassword: "pw", MaxPlayers: 10}); !errors.Is(err, ErrInvalidGameBranch) {
		t.Fatalf("expected ErrInvalidGameBranch, got %v", err)
	}
	if _, err := s.CreateServer(1, models.ServerRequest{Game: "xbox", Port: 7796, QueryPort: 27093, RCONPort: 27094, AdminPassword: "pw", MaxPlayers: 10}); !errors.Is(err, ErrInvalidGame) {
		t.Fatalf("expected ErrInvalidGame, got %v", err)
	}
}
