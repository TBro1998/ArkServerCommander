package server

import (
	"errors"
	"testing"

	"ark-server-commander/models"
	"ark-server-commander/service/docker_manager"
)

func TestApplyPhase(t *testing.T) {
	cases := []struct {
		name  string
		lines []string
		want  models.RuntimeInfo
	}{
		{"empty", nil, models.RuntimeInfo{}},
		{"first download", []string{
			" [*] Server files not found in /home/steam/arkserver - downloading, this takes a while...",
			" [*] Updating Steam client",
		}, models.RuntimeInfo{Phase: models.PhaseDownloading}},
		{"download progress", []string{
			" [*] Updating ARK server",
			" Update state (0x61) downloading, progress: 12.34 (1234 / 10000)",
			" Update state (0x61) downloading, progress: 45.67 (4567 / 10000)",
		}, models.RuntimeInfo{Phase: models.PhaseDownloading, Progress: 45.67}},
		{"mods", []string{
			" [i] ARK server update completed",
			" [i] Mod 111 already installed (1/3)",
			" [*] Downloading mod: 222 (2/3)",
		}, models.RuntimeInfo{Phase: models.PhaseInstallingMods, ModCurrent: 2, ModTotal: 3}},
		{"installing mod", []string{" [*] Installing mod: 333 (3/3)"},
			models.RuntimeInfo{Phase: models.PhaseInstallingMods, ModCurrent: 3, ModTotal: 3}},
		{"launching", []string{
			" [*] Downloading mod: 222 (2/3)",
			" [*] Starting ARK server...",
			"some proton output",
		}, models.RuntimeInfo{Phase: models.PhaseLaunching}},
		{"no mods", []string{" [i] Server files already installed - skipping download", " [i] Skipping mod update"},
			models.RuntimeInfo{Phase: models.PhaseLaunching}},
		{"game output only", []string{"wine: something", "Proton: Upgrading prefix"}, models.RuntimeInfo{}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var got models.RuntimeInfo
			applyPhase(&got, tc.lines)
			if got != tc.want {
				t.Fatalf("got %+v, want %+v", got, tc.want)
			}
		})
	}
}

func TestParsePlayers(t *testing.T) {
	if got := parsePlayers("No Players Connected"); len(got) != 0 {
		t.Fatalf("expected no players, got %+v", got)
	}
	out := "0. Alice, 76561198000000001\n1. Bob, The Builder, 76561198000000002\n\n"
	got := parsePlayers(out)
	if len(got) != 2 {
		t.Fatalf("expected 2 players, got %+v", got)
	}
	if got[0].Name != "Alice" || got[0].SteamID != "76561198000000001" || got[0].Index != 0 {
		t.Fatalf("unexpected first player %+v", got[0])
	}
	// 名称里带逗号时，最后一个逗号之后才是 SteamID
	if got[1].Name != "Bob, The Builder" || got[1].SteamID != "76561198000000002" {
		t.Fatalf("unexpected second player %+v", got[1])
	}
}

func TestNormalizeServerFields(t *testing.T) {
	cluster, mods, err := normalizeServerFields(" my-cluster_1 ", " 123, 456 ,，789 ")
	if err != nil || cluster != "my-cluster_1" || mods != "123,456,789" {
		t.Fatalf("got %q %q %v", cluster, mods, err)
	}
	if _, _, err := normalizeServerFields("bad id", ""); !errors.Is(err, ErrInvalidClusterID) {
		t.Fatalf("expected ErrInvalidClusterID, got %v", err)
	}
	if _, _, err := normalizeServerFields("", "123,abc"); !errors.Is(err, ErrInvalidModIDs) {
		t.Fatalf("expected ErrInvalidModIDs, got %v", err)
	}
	if c, m, err := normalizeServerFields("", ""); err != nil || c != "" || m != "" {
		t.Fatalf("empty values should be allowed, got %q %q %v", c, m, err)
	}
}

func TestValidateSchedule(t *testing.T) {
	ok := []models.ScheduleRequest{
		{Action: models.ScheduleRestart, Cron: "0 4 * * *", WarnMinutes: 10},
		{Action: models.ScheduleBackup, Cron: "@every 6h", Keep: 5},
		{Action: models.ScheduleCommand, Cron: "*/30 * * * *", Command: "SaveWorld"},
		{Action: models.ScheduleStop, Cron: "0 2 * * 1"},
	}
	for _, req := range ok {
		if _, err := validateSchedule(req); err != nil {
			t.Fatalf("%+v should be valid: %v", req, err)
		}
	}
	bad := []models.ScheduleRequest{
		{Action: models.ScheduleRestart, Cron: "61 * * * *"},
		{Action: "explode", Cron: "0 4 * * *"},
		{Action: models.ScheduleCommand, Cron: "0 4 * * *", Command: ""},
		{Action: models.ScheduleCommand, Cron: "0 4 * * *", Command: "a\nb"},
		{Action: models.ScheduleRestart, Cron: "0 4 * * *", WarnMinutes: 90},
	}
	for _, req := range bad {
		if _, err := validateSchedule(req); err == nil {
			t.Fatalf("%+v should be invalid", req)
		}
	}
	// 备份任务不需要倒计时
	got, _ := validateSchedule(models.ScheduleRequest{Action: models.ScheduleBackup, Cron: "@daily", WarnMinutes: 5})
	if got.WarnMinutes != 0 {
		t.Fatalf("backup warn minutes should be reset, got %d", got.WarnMinutes)
	}
}

func TestBackupNamePattern(t *testing.T) {
	valid := []string{"20261111-153000_manual.tar.gz", "20261111-153000_pre-restore.tar.gz", "20261111-153000_scheduled-2.tar.gz"}
	for _, n := range valid {
		if !backupNamePattern.MatchString(n) {
			t.Fatalf("%q should match", n)
		}
	}
	invalid := []string{"../20261111-153000_manual.tar.gz", "20261111-153000_manual.tar.gz.tmp", "foo.tar.gz", "20261111-153000_other.tar.gz"}
	for _, n := range invalid {
		if backupNamePattern.MatchString(n) {
			t.Fatalf("%q should not match", n)
		}
	}
}

func TestPickSeedSource(t *testing.T) {
	installed := docker_manager.ServerFilesState{Installed: true, Executable: true}
	candidates := []models.Server{
		{ID: 2, Status: "running"},
		{ID: 3, Status: "stopped"},
		{ID: 4, Status: "stopped"},
		{ID: 5, Status: "stopped"},
	}
	states := map[uint]docker_manager.ServerFilesState{
		2: installed,
		3: {Installed: false, Executable: true}, // 正在更新（标记已清除），不能作为来源
		4: installed,
		5: installed,
	}

	// 已停止的优先，其次按 ID
	if got, ok := pickSeedSource(false, candidates, states); !ok || got.ID != 4 {
		t.Fatalf("expected server 4, got %+v %v", got, ok)
	}
	// 只有运行中的服务器已安装时也可以作为来源
	if got, ok := pickSeedSource(false, candidates[:2], states); !ok || got.ID != 2 {
		t.Fatalf("expected server 2, got %+v %v", got, ok)
	}
	// 目标已有服务端文件（例如更新游戏刚清除了标记）时不复制
	if _, ok := pickSeedSource(true, candidates, states); ok {
		t.Fatal("should not copy when the target already has server files")
	}
	// 没有已安装完成的服务器
	if _, ok := pickSeedSource(false, candidates[1:2], states); ok {
		t.Fatal("should not copy from an incomplete install")
	}
}
