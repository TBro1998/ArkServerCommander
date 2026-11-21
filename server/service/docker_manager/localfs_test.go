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
	"strings"
	"testing"
	"time"
	"unicode/utf16"

	"ark-server-commander/config"
	"ark-server-commander/utils"
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

	// ARK 在内容含中文时把配置写成带 BOM 的 UTF-16，读取时要转成 UTF-8
	text := "[ServerSettings]\r\nServerAdminPassword=pw\r\n[SessionSettings]\r\nSessionName=测试-孤岛\r\n"
	units := utf16.Encode([]rune(text))
	le := []byte{0xFF, 0xFE}
	be := []byte{0xFE, 0xFF}
	for _, u := range units {
		le = append(le, byte(u), byte(u>>8))
		be = append(be, byte(u>>8), byte(u))
	}
	for name, data := range map[string][]byte{"utf-16le": le, "utf-16be": be} {
		if err := os.WriteFile(filepath.Join(configDir, "GameUserSettings.ini"), data, 0o644); err != nil {
			t.Fatal(err)
		}
		got, err := dm.ReadConfigFile(1, "GameUserSettings.ini")
		if err != nil || got != strings.TrimSpace(text) {
			t.Fatalf("%s: read = %q, %v", name, got, err)
		}
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

func TestPrepareGameBranch(t *testing.T) {
	useTempDataDir(t)
	var dm *DockerManager
	dir := ServerFilesDir(1)
	marker := filepath.Join(dir, installMarkerFile)
	manifest := filepath.Join(dir, appManifestFile)
	writeFile(t, marker, "linux")
	writeFile(t, filepath.Join(dir, filepath.FromSlash(linuxServerExecutable)), "bin")
	writeFile(t, manifest, "\"AppState\"\n{\n\t\"appid\"\t\t\"376030\"\n\t\"UserConfig\"\n\t{\n\t\t\"BetaKey\"\t\t\"preaquatica\"\n\t}\n}\n")

	// 旧文件卷没有记录，视为正式版：只补写记录，不清除安装标记
	if err := dm.PrepareGameBranch(1, utils.GameBranchPublic); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(marker); err != nil {
		t.Fatalf("marker removed for the same branch: %v", err)
	}
	if got := readGameBranch(dir); got != utils.GameBranchPublic {
		t.Fatalf("branch = %q", got)
	}

	// 切到 preaquatica：清除安装标记，appmanifest 不动（APP_UPDATE 带 -beta）
	if err := dm.PrepareGameBranch(1, utils.GameBranchPreAquatica); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(marker); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("marker should be removed: %v", err)
	}
	if got := readGameBranch(dir); got != utils.GameBranchPreAquatica {
		t.Fatalf("branch = %q", got)
	}
	states, _ := dm.CheckServerFiles([]uint{1})
	if states[1].Branch != utils.GameBranchPreAquatica {
		t.Fatalf("state branch = %q", states[1].Branch)
	}

	// 同一版本再次启动：不清除标记
	writeFile(t, marker, "linux")
	if err := dm.PrepareGameBranch(1, utils.GameBranchPreAquatica); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(marker); err != nil {
		t.Fatalf("marker removed for the same branch: %v", err)
	}

	// 切回正式版：清除标记，去掉 appmanifest 里的 BetaKey，其他内容保留
	if err := dm.PrepareGameBranch(1, utils.GameBranchPublic); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(marker); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("marker should be removed: %v", err)
	}
	data, err := os.ReadFile(manifest)
	if err != nil {
		t.Fatal(err)
	}
	if want := "\"AppState\"\n{\n\t\"appid\"\t\t\"376030\"\n\t\"UserConfig\"\n\t{\n\t}\n}\n"; string(data) != want {
		t.Fatalf("manifest = %q", data)
	}

	// 文件卷目录不存在时不处理
	if err := dm.PrepareGameBranch(9, utils.GameBranchPreAquatica); err != nil {
		t.Fatal(err)
	}
}
