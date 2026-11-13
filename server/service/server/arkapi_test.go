package server

import (
	"archive/zip"
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"ark-server-commander/assets"
	"ark-server-commander/config"
	"ark-server-commander/service/docker_manager"
	"ark-server-commander/utils"

	"go.uber.org/zap"
)

func zipBytes(t *testing.T, files map[string]string) []byte {
	t.Helper()
	r := makeZip(t, files)
	data := make([]byte, r.Size())
	if _, err := r.ReadAt(data, 0); err != nil {
		t.Fatal(err)
	}
	return data
}

func useTempDataDir(t *testing.T) {
	t.Helper()
	if utils.Logger == nil {
		utils.Logger = zap.NewNop()
	}
	old := config.DataDir
	config.DataDir = t.TempDir()
	t.Cleanup(func() { config.DataDir = old })
}

func TestCompareVersions(t *testing.T) {
	cases := []struct {
		a, b string
		want int
	}{
		{"3.56", "3.54", 1},
		{"3.6", "3.56", -1},
		{"3.54", "3.54", 0},
		{"3.54.1", "3.54", 1},
		{"10.0", "9.9", 1},
		{"3.56-fix", "3.56", 1},
	}
	for _, c := range cases {
		if got := compareVersions(c.a, c.b); got != c.want {
			t.Fatalf("compareVersions(%q, %q) = %d, want %d", c.a, c.b, got, c.want)
		}
	}
}

func TestReadArkApiZipRoot(t *testing.T) {
	cases := []struct {
		name     string
		files    map[string]string
		wantRoot string
		wantErr  bool
	}{
		{"根目录", map[string]string{"version.dll": "x", "msdia140.dll": "x", "config.json": "{}"}, "", false},
		{"带前缀", map[string]string{"ArkApi_3.54/version.dll": "x", "ArkApi_3.54/ArkApi/Plugins/A/A.dll": "x"}, "ArkApi_3.54/", false},
		{"多个时取最浅的", map[string]string{"a/b/version.dll": "x", "a/version.dll": "x"}, "a/", false},
		{"没有 version.dll", map[string]string{"ArkShop/ArkShop.dll": "x"}, "", true},
		{"越界路径", map[string]string{"../version.dll": "x"}, "", true},
	}
	for _, c := range cases {
		data := zipBytes(t, c.files)
		_, root, err := readArkApiZip(bytes.NewReader(data), int64(len(data)))
		if c.wantErr {
			if !errors.Is(err, ErrInvalidArkApi) {
				t.Fatalf("%s: expected ErrInvalidArkApi, got %v", c.name, err)
			}
			continue
		}
		if err != nil || root != c.wantRoot {
			t.Fatalf("%s: got root %q err %v, want %q", c.name, root, err, c.wantRoot)
		}
	}
}

func TestArkApiVersionFromZip(t *testing.T) {
	cases := map[[2]string]string{
		{"ArkApi_3.56.zip", ""}:          "3.56",
		{"arkapi-v3.54.zip", ""}:         "3.54",
		{"download.zip", "ArkApi 3.53/"}: "3.53",
		{"download.zip", ""}:             "",
	}
	for in, want := range cases {
		if got := arkApiVersionFromZip(in[0], in[1]); got != want {
			t.Fatalf("arkApiVersionFromZip(%q, %q) = %q, want %q", in[0], in[1], got, want)
		}
	}
}

func TestParseArkApiReleases(t *testing.T) {
	var body []githubRelease
	add := func(tag string, draft, pre bool, assets ...string) {
		r := githubRelease{TagName: tag, Draft: draft, Prerelease: pre}
		for _, a := range assets {
			r.Assets = append(r.Assets, struct {
				Name string `json:"name"`
				URL  string `json:"browser_download_url"`
			}{a, "https://example.com/" + a})
		}
		body = append(body, r)
	}
	add("v3.54", false, false, "Source.zip", "ArkApi_3.54.zip")
	add("v3.55", true, false, "ArkApi_3.55.zip")
	add("v3.60-beta", false, true, "ArkApi.zip")
	add("v3.50", false, false) // 没有附件
	got := parseArkApiReleases(body)
	if len(got) != 1 || got[0].version != "3.54" || got[0].download != "https://example.com/ArkApi_3.54.zip" {
		t.Fatalf("got %+v", got)
	}
}

func TestInstallArkApi(t *testing.T) {
	useTempDataDir(t)
	const id = 7
	win64 := arkApiWin64Dir(id)

	v1 := zipBytes(t, map[string]string{
		"ArkApi_1.0/version.dll":            "v1",
		"ArkApi_1.0/config.json":            `{"default":true}`,
		"ArkApi_1.0/ArkApi/Plugins/A/A.dll": "plugin",
		"ArkApi_1.0/ArkApi/Logs/readme.txt": "logs",
		"ArkApi_1.0/msdia140.dll":           "dia",
	})
	if err := installArkApi(id, v1, "1.0"); err != nil {
		t.Fatal(err)
	}
	if b, _ := os.ReadFile(filepath.Join(win64, "version.dll")); string(b) != "v1" {
		t.Fatalf("version.dll = %q", b)
	}
	// 没有配置时用压缩包里的
	if b, _ := os.ReadFile(filepath.Join(win64, "config.json")); string(b) != `{"default":true}` {
		t.Fatalf("config.json = %q", b)
	}
	// 插件目录的内容不解压，其他子目录照常
	if _, err := os.Stat(filepath.Join(win64, "ArkApi", "Plugins", "A")); !os.IsNotExist(err) {
		t.Fatal("plugins inside the ArkApi zip should be skipped")
	}
	if _, err := os.Stat(filepath.Join(win64, "ArkApi", "Logs", "readme.txt")); err != nil {
		t.Fatal(err)
	}
	if v, installed := installedArkApiVersion(id); v != "1.0" || !installed {
		t.Fatalf("installed = %q %v", v, installed)
	}

	// 更新：替换 dll，保留用户改过的 config.json
	os.WriteFile(filepath.Join(win64, "config.json"), []byte(`{"mine":1}`), 0o644)
	v2 := zipBytes(t, map[string]string{"version.dll": "v2", "config.json": `{"default":true}`})
	if err := installArkApi(id, v2, "2.0"); err != nil {
		t.Fatal(err)
	}
	if b, _ := os.ReadFile(filepath.Join(win64, "version.dll")); string(b) != "v2" {
		t.Fatalf("version.dll = %q", b)
	}
	if b, _ := os.ReadFile(filepath.Join(win64, "config.json")); string(b) != `{"mine":1}` {
		t.Fatalf("config.json should be kept, got %q", b)
	}
	if v, _ := installedArkApiVersion(id); v != "2.0" {
		t.Fatalf("version = %q", v)
	}
	entries, _ := os.ReadDir(win64)
	for _, e := range entries {
		if strings.HasSuffix(e.Name(), ".arkapi-tmp") {
			t.Fatalf("temporary file left behind: %s", e.Name())
		}
	}
}

func TestEnsureArkApiInstalledUsesLatestLocal(t *testing.T) {
	useTempDataDir(t)
	const id = 3

	// 版本库里没有更新的版本时装内置版本
	if err := ensureArkApiInstalled(id); err != nil {
		t.Fatal(err)
	}
	if v, installed := installedArkApiVersion(id); v != assets.BuiltinArkApiVersion || !installed {
		t.Fatalf("installed = %q %v", v, installed)
	}

	// 已安装时不动
	os.WriteFile(filepath.Join(docker_manager.ServerFilesDir(id), arkApiVersionFile), []byte("old\n"), 0o644)
	if err := ensureArkApiInstalled(id); err != nil {
		t.Fatal(err)
	}
	if v, _ := installedArkApiVersion(id); v != "old" {
		t.Fatalf("should not reinstall, got %q", v)
	}

	// 版本库里有更新的版本时用它
	if err := saveArkApiZip("99.1", zipBytes(t, map[string]string{"version.dll": "new"})); err != nil {
		t.Fatal(err)
	}
	const other = 4
	if err := ensureArkApiInstalled(other); err != nil {
		t.Fatal(err)
	}
	if v, _ := installedArkApiVersion(other); v != "99.1" {
		t.Fatalf("expected 99.1, got %q", v)
	}
}

func TestArkApiLibrary(t *testing.T) {
	useTempDataDir(t)
	s := NewServerService()

	data := zipBytes(t, map[string]string{"version.dll": "x"})
	version, err := s.UploadArkApi(bytes.NewReader(data), int64(len(data)), "ArkApi_3.60.zip", "")
	if err != nil || version != "3.60" {
		t.Fatalf("upload: %q %v", version, err)
	}
	// 识别不出版本号时要求填写
	if _, err := s.UploadArkApi(bytes.NewReader(data), int64(len(data)), "upload.zip", ""); !errors.Is(err, ErrInvalidArkApi) {
		t.Fatalf("expected ErrInvalidArkApi, got %v", err)
	}
	if _, err := s.UploadArkApi(bytes.NewReader(data), int64(len(data)), "upload.zip", "../x"); !errors.Is(err, ErrInvalidArkApi) {
		t.Fatalf("expected ErrInvalidArkApi for bad version, got %v", err)
	}

	local := localArkApiVersions()
	if !local["3.60"].Local || !local["3.60"].Deletable || !local[assets.BuiltinArkApiVersion].Builtin {
		t.Fatalf("local versions: %+v", local)
	}
	if latestLocalArkApi() != "3.60" {
		t.Fatalf("latest local = %q", latestLocalArkApi())
	}

	if err := s.DeleteArkApiVersion(assets.BuiltinArkApiVersion); !errors.Is(err, ErrArkApiBuiltin) {
		t.Fatalf("expected ErrArkApiBuiltin, got %v", err)
	}
	if err := s.DeleteArkApiVersion("3.60"); err != nil {
		t.Fatal(err)
	}
	if err := s.DeleteArkApiVersion("3.60"); !errors.Is(err, ErrArkApiNotFound) {
		t.Fatalf("expected ErrArkApiNotFound, got %v", err)
	}
}

func TestBuiltinArkApiZip(t *testing.T) {
	zr, err := zip.NewReader(bytes.NewReader(assets.BuiltinArkApiZip), int64(len(assets.BuiltinArkApiZip)))
	if err != nil {
		t.Fatal(err)
	}
	_, root, err := readArkApiZip(bytes.NewReader(assets.BuiltinArkApiZip), int64(len(assets.BuiltinArkApiZip)))
	if err != nil || root != "" {
		t.Fatalf("root %q err %v", root, err)
	}
	if len(zr.File) < 2 {
		t.Fatalf("unexpected builtin zip: %d files", len(zr.File))
	}
}
