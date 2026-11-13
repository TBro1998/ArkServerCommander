package docker_manager

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"ark-server-commander/config"
)

func TestParseInstalledMods(t *testing.T) {
	got := parseInstalledMods(strings.NewReader("111\n222 1700000000\n\n  333   1700000100  \n222 1700000200\n"))
	want := map[string]int64{"111": 0, "222": 1700000200, "333": 1700000100}
	if len(got) != len(want) {
		t.Fatalf("got %v, want %v", got, want)
	}
	for id, ts := range want {
		if got[id] != ts {
			t.Errorf("%s: got %d, want %d", id, got[id], ts)
		}
	}
}

func TestInstalledModsAndRemove(t *testing.T) {
	old := config.DataDir
	config.DataDir = t.TempDir()
	t.Cleanup(func() { config.DataDir = old })

	// 服务端文件还没下载
	mods, err := InstalledMods(1)
	if err != nil || len(mods) != 0 {
		t.Fatalf("expected empty result, got %v, %v", mods, err)
	}

	dir := ServerFilesDir(1)
	modsDir := filepath.Join(dir, "ShooterGame", "Content", "Mods")
	if err := os.MkdirAll(modsDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, installedModsFile), []byte("111\n222 1700000000\n333 1700000100\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	// 旧格式的 Mod 用 .mod 文件的修改时间
	modTime := time.Unix(1600000000, 0)
	modFile := filepath.Join(modsDir, "111.mod")
	if err := os.WriteFile(modFile, nil, 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Chtimes(modFile, modTime, modTime); err != nil {
		t.Fatal(err)
	}

	mods, err = InstalledMods(1)
	if err != nil {
		t.Fatal(err)
	}
	if mods["111"] != 1600000000 || mods["222"] != 1700000000 || mods["333"] != 1700000100 {
		t.Fatalf("unexpected result %v", mods)
	}

	if err := RemoveInstalledMods(1, []string{"222", "999"}); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(filepath.Join(dir, installedModsFile))
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "111\n333 1700000100\n" {
		t.Fatalf("unexpected file content %q", data)
	}
}
