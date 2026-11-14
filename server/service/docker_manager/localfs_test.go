package docker_manager

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"errors"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"ark-server-commander/config"
)

func useTempDataDir(t *testing.T) string {
	t.Helper()
	old := config.DataDir
	config.DataDir = t.TempDir()
	t.Cleanup(func() { config.DataDir = old })
	return config.DataDir
}

func writeFile(t *testing.T, p, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestConfigFileReadWrite(t *testing.T) {
	useTempDataDir(t)
	var dm *DockerManager
	if err := os.MkdirAll(ServerSavedDir(1), 0o755); err != nil {
		t.Fatal(err)
	}

	if _, err := dm.ReadConfigFile(1, "Game.ini"); !errors.Is(err, ErrConfigNotFound) {
		t.Fatalf("expected ErrConfigNotFound, got %v", err)
	}
	if _, err := dm.ReadConfigFile(2, "Game.ini"); !errors.Is(err, ErrConfigNotFound) {
		t.Fatalf("missing saved dir: expected ErrConfigNotFound, got %v", err)
	}
	for _, bad := range []string{"", "..", "../x.ini", "a/b.ini", `a\b.ini`} {
		if err := dm.WriteConfigFile(1, bad, "x"); err == nil {
			t.Fatalf("expected error for file name %q", bad)
		}
	}

	if err := dm.WriteConfigFile(1, "Game.ini", "[a]\nk=v"); err != nil {
		t.Fatal(err)
	}
	configDir := filepath.Join(ServerSavedDir(1), "Config", "WindowsServer")
	raw, err := os.ReadFile(filepath.Join(configDir, "Game.ini"))
	if err != nil || string(raw) != "[a]\nk=v\n" {
		t.Fatalf("file on disk = %q, %v", raw, err)
	}
	// 覆盖已有文件，读取时去掉 BOM 和首尾空白
	if err := dm.WriteConfigFile(1, "Game.ini", utf8BOM+"\n[a]\nk=w\n\n"); err != nil {
		t.Fatal(err)
	}
	got, err := dm.ReadConfigFile(1, "Game.ini")
	if err != nil || got != "[a]\nk=w" {
		t.Fatalf("read = %q, %v", got, err)
	}
	entries, _ := os.ReadDir(configDir)
	if len(entries) != 1 {
		t.Fatalf("temporary file left behind: %v", entries)
	}
}

func TestExportImportSaved(t *testing.T) {
	useTempDataDir(t)
	var dm *DockerManager
	saved := ServerSavedDir(1)
	writeFile(t, filepath.Join(saved, "SavedArks", "TheIsland.ark"), "world")
	writeFile(t, filepath.Join(saved, "Config", "WindowsServer", "Game.ini"), "[a]\n")
	writeFile(t, filepath.Join(saved, "Logs", "ShooterGame.log"), "log")
	writeFile(t, filepath.Join(saved, "clusters", "x"), "cluster")
	if runtime.GOOS != "windows" {
		if err := os.Symlink("WindowsServer", filepath.Join(saved, "Config", "LinuxServer")); err != nil {
			t.Fatal(err)
		}
	}

	var buf bytes.Buffer
	if err := dm.ExportSaved(1, &buf); err != nil {
		t.Fatal(err)
	}
	names := map[string]bool{}
	gz, err := gzip.NewReader(bytes.NewReader(buf.Bytes()))
	if err != nil {
		t.Fatal(err)
	}
	tr := tar.NewReader(gz)
	for {
		h, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatal(err)
		}
		names[h.Name] = true
	}
	for _, want := range []string{"Saved/", "Saved/SavedArks/TheIsland.ark", "Saved/Config/WindowsServer/Game.ini"} {
		if !names[want] {
			t.Fatalf("missing %s in %v", want, names)
		}
	}
	for _, skip := range []string{"Saved/Logs/", "Saved/Logs/ShooterGame.log", "Saved/clusters/x", "Saved/Config/LinuxServer"} {
		if names[skip] {
			t.Fatalf("%s should be skipped", skip)
		}
	}
	if has, err := ValidateBackup(bytes.NewReader(buf.Bytes())); err != nil || !has {
		t.Fatalf("exported backup invalid: %v %v", has, err)
	}

	// 恢复到另一台服务器：旧的存档和配置被清空，其他目录保留
	dst := ServerSavedDir(2)
	writeFile(t, filepath.Join(dst, "SavedArks", "Old.ark"), "old")
	writeFile(t, filepath.Join(dst, "Config", "WindowsServer", "Extra.ini"), "old")
	writeFile(t, filepath.Join(dst, "Logs", "keep.log"), "keep")
	if err := dm.ImportSaved(2, bytes.NewReader(buf.Bytes())); err != nil {
		t.Fatal(err)
	}
	if data, err := os.ReadFile(filepath.Join(dst, "SavedArks", "TheIsland.ark")); err != nil || string(data) != "world" {
		t.Fatalf("restored save = %q, %v", data, err)
	}
	for _, gone := range []string{"SavedArks/Old.ark", "Config/WindowsServer/Extra.ini"} {
		if _, err := os.Stat(filepath.Join(dst, filepath.FromSlash(gone))); !errors.Is(err, os.ErrNotExist) {
			t.Fatalf("%s should be removed, err=%v", gone, err)
		}
	}
	if _, err := os.Stat(filepath.Join(dst, "Logs", "keep.log")); err != nil {
		t.Fatalf("Logs should be kept: %v", err)
	}
}

func TestImportSavedRejectsEscape(t *testing.T) {
	useTempDataDir(t)
	var dm *DockerManager
	if err := os.MkdirAll(ServerSavedDir(1), 0o755); err != nil {
		t.Fatal(err)
	}
	archive := makeArchive(t, []entry{{"Saved/../../evil", tar.TypeReg, "x"}})
	if err := dm.ImportSaved(1, archive); err == nil {
		t.Fatal("expected error for path outside Saved")
	}
	if _, err := os.Stat(filepath.Join(config.DataDir, "servers", "evil")); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("file written outside saved dir: %v", err)
	}
}

func TestCheckAndCopyServerFiles(t *testing.T) {
	useTempDataDir(t)
	var dm *DockerManager
	src := ServerFilesDir(1)
	modFile := filepath.Join(src, "ShooterGame", "Content", "Mods", "111.mod")
	writeFile(t, filepath.Join(src, installMarkerFile), "linux")
	writeFile(t, filepath.Join(src, installedModsFile), "111 1700000000\n")
	writeFile(t, filepath.Join(src, filepath.FromSlash(linuxServerExecutable)), "bin")
	writeFile(t, modFile, "mod")
	modTime := time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC)
	if err := os.Chtimes(modFile, modTime, modTime); err != nil {
		t.Fatal(err)
	}
	if runtime.GOOS != "windows" {
		if err := os.Chmod(filepath.Join(src, filepath.FromSlash(linuxServerExecutable)), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.Symlink("../Linux/ShooterGameServer", filepath.Join(src, "ShooterGame", "Content", "link")); err != nil {
			t.Fatal(err)
		}
	}
	// 目标已有的同名文件被覆盖，其他文件保留
	dst := ServerFilesDir(2)
	writeFile(t, filepath.Join(dst, installedModsFile), "stale")
	writeFile(t, filepath.Join(dst, "other"), "keep")

	states, err := dm.CheckServerFiles([]uint{1, 2, 3})
	if err != nil {
		t.Fatal(err)
	}
	if st := states[1]; !st.Installed || !st.Linux || st.Windows {
		t.Fatalf("server 1 state = %+v", st)
	}
	if st, ok := states[2]; !ok || st.Installed || st.HasExecutable() {
		t.Fatalf("server 2 state = %+v, %v", st, ok)
	}
	if _, ok := states[3]; ok {
		t.Fatal("server 3 has no files dir and should be absent")
	}

	if err := dm.CopyServerFiles(1, 2); err != nil {
		t.Fatal(err)
	}
	states, _ = dm.CheckServerFiles([]uint{2})
	if st := states[2]; st.Installed || !st.Linux {
		t.Fatalf("after copy state = %+v (marker must be removed)", st)
	}
	if data, _ := os.ReadFile(filepath.Join(dst, installedModsFile)); string(data) != "111 1700000000\n" {
		t.Fatalf(".installed_mods = %q", data)
	}
	if data, _ := os.ReadFile(filepath.Join(dst, "other")); string(data) != "keep" {
		t.Fatalf("existing file lost: %q", data)
	}
	info, err := os.Stat(filepath.Join(dst, "ShooterGame", "Content", "Mods", "111.mod"))
	if err != nil || !info.ModTime().Equal(modTime) {
		t.Fatalf("mod time not preserved: %v %v", info, err)
	}
	if _, err := os.Stat(filepath.Join(src, installMarkerFile)); err != nil {
		t.Fatalf("source marker must stay: %v", err)
	}
	if runtime.GOOS != "windows" {
		info, err := os.Stat(filepath.Join(dst, filepath.FromSlash(linuxServerExecutable)))
		if err != nil || info.Mode().Perm() != 0o755 {
			t.Fatalf("executable mode not preserved: %v %v", info, err)
		}
		link, err := os.Readlink(filepath.Join(dst, "ShooterGame", "Content", "link"))
		if err != nil || link != "../Linux/ShooterGameServer" {
			t.Fatalf("symlink not preserved: %q %v", link, err)
		}
	}

	if err := dm.RemoveInstallMarkers(1, true, true); err != nil {
		t.Fatal(err)
	}
	if err := dm.RemoveInstallMarkers(1, true, true); err != nil {
		t.Fatalf("removing missing markers should succeed: %v", err)
	}
	states, _ = dm.CheckServerFiles([]uint{1})
	if states[1].Installed {
		t.Fatal("marker should be removed")
	}
}
