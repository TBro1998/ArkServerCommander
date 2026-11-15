package server

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"

	"ark-server-commander/database"
	"ark-server-commander/models"
	"ark-server-commander/service/docker_manager"
	"ark-server-commander/utils"
)

func TestServerConfigCopyAndImport(t *testing.T) {
	setupTestDB(t)
	useTempDataDir(t)

	dm, err := docker_manager.GetDockerManager()
	if err != nil {
		t.Fatal(err)
	}
	s := NewServerService()

	src := createTestServer(t, 7777, 27015, 32330)
	database.DB.Model(&src).Updates(map[string]any{"map": "Ragnarok", "game_mod_ids": "111,222", "cluster_id": "c1", "arkapi": true})
	dst := createTestServer(t, 7787, 27025, 32340)
	running := createTestServer(t, 7797, 27035, 32350)
	database.DB.Model(&running).Update("status", "running")

	write := func(id uint, file, content string) {
		t.Helper()
		if err := os.MkdirAll(docker_manager.ServerSavedDir(id), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := dm.WriteConfigFile(id, file, content); err != nil {
			t.Fatal(err)
		}
	}
	read := func(id uint, file string) string {
		t.Helper()
		content, err := dm.ReadConfigFile(id, file)
		if err != nil {
			t.Fatal(err)
		}
		return content
	}
	write(src.ID, utils.GameUserSettingsFileName, "[ServerSettings]\nSessionName=Source\nServerAdminPassword=srcpw\nMaxPlayers=70\nXPMultiplier=3\n\n[SessionSettings]\nSessionName=Source\nPort=7777\n")
	write(src.ID, utils.GameIniFileName, "[/script/shootergame.shootergamemode]\nbAllowUnlimitedRespecs=True\n")
	dstGUS := "[ServerSettings]\nSessionName=Target\nServerAdminPassword=dstpw\nMaxPlayers=20\nXPMultiplier=1\n\n[SessionSettings]\nSessionName=Target\nPort=7787\n"
	write(dst.ID, utils.GameUserSettingsFileName, dstGUS)
	write(running.ID, utils.GameUserSettingsFileName, dstGUS)

	// 导出不含管理员密码
	cfg, err := s.ExportConfig(1, fmt.Sprint(src.ID))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(cfg.GameUserSettings, "srcpw") || cfg.Map != "Ragnarok" || cfg.GameModIds != "111,222" || cfg.Format != models.ServerConfigFormat {
		t.Fatalf("导出内容不对: %+v", cfg)
	}

	// 只复制配置文件和 Mod：目标保留自己的名称、端口、人数和密码，地图、集群不变
	results, err := s.CopyConfig(1, fmt.Sprint(src.ID), models.CopyConfigRequest{
		TargetIDs: []uint{src.ID, dst.ID, running.ID, 999},
		Parts:     []string{models.ConfigPartConfigFiles, models.ConfigPartMods},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 3 || !results[0].Success || !results[1].Success || results[2].Success || results[2].Code != "server_not_found" {
		t.Fatalf("结果不对: %+v", results)
	}
	if results[0].ConfigPending || !results[1].ConfigPending {
		t.Fatalf("运行中的服务器应暂存配置: %+v", results)
	}
	// 导出时去掉的密码写回到 section 开头，其余键留在原位置
	want := "[ServerSettings]\nServerAdminPassword=dstpw\nSessionName=Target\nMaxPlayers=20\nXPMultiplier=3\n\n[SessionSettings]\nSessionName=Target\nPort=7787"
	if got := strings.TrimSpace(read(dst.ID, utils.GameUserSettingsFileName)); got != want {
		t.Fatalf("目标配置文件:\n%q\n期望:\n%q", got, want)
	}
	if got := read(dst.ID, utils.GameIniFileName); !strings.Contains(got, "bAllowUnlimitedRespecs=True") {
		t.Fatalf("Game.ini 没有复制: %q", got)
	}
	var updated models.Server
	updated = models.Server{}
	database.DB.First(&updated, dst.ID)
	if updated.GameModIds != "111,222" || updated.Map == "Ragnarok" || updated.ClusterID != "" || updated.ArkApi {
		t.Fatalf("只应复制 Mod: %+v", updated)
	}
	updated = models.Server{}
	database.DB.First(&updated, running.ID)
	if strings.TrimSpace(updated.PendingGameUserSettings) != want || strings.TrimSpace(read(running.ID, utils.GameUserSettingsFileName)) != strings.TrimSpace(dstGUS) {
		t.Fatalf("运行中的服务器应只写入暂存配置: %q", updated.PendingGameUserSettings)
	}

	// 导入全部：选项里的人数同步到配置文件
	cfg.MaxPlayers = 50
	if _, err := s.ImportConfig(1, fmt.Sprint(dst.ID), cfg, nil); err != nil {
		t.Fatal(err)
	}
	updated = models.Server{}
	database.DB.First(&updated, dst.ID)
	if updated.Map != "Ragnarok" || updated.ClusterID != "c1" || !updated.ArkApi || updated.MaxPlayers != 50 || updated.SessionName != "test" || updated.Port != 7787 {
		t.Fatalf("导入后: %+v", updated)
	}
	if got := read(dst.ID, utils.GameUserSettingsFileName); !strings.Contains(got, "MaxPlayers=50") || strings.Contains(got, "MaxPlayers=20") {
		t.Fatalf("人数没有同步: %q", got)
	}

	// 无效的配置和部分
	bad := *cfg
	bad.Format = "other"
	if _, err := s.ImportConfig(1, fmt.Sprint(dst.ID), &bad, nil); !errors.Is(err, ErrInvalidServerConfig) {
		t.Fatalf("期望 ErrInvalidServerConfig，实际 %v", err)
	}
	if _, err := s.CopyConfig(1, fmt.Sprint(src.ID), models.CopyConfigRequest{TargetIDs: []uint{dst.ID}, Parts: []string{"ports"}}); !errors.Is(err, ErrInvalidServerConfig) {
		t.Fatalf("期望 ErrInvalidServerConfig，实际 %v", err)
	}
	if _, err := s.CopyConfig(1, fmt.Sprint(src.ID), models.CopyConfigRequest{TargetIDs: []uint{src.ID}}); !errors.Is(err, ErrInvalidServerConfig) {
		t.Fatalf("只有来源服务器时应报错，实际 %v", err)
	}
}
