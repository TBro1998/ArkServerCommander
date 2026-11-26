package server

import (
	"archive/zip"
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"ark-server-commander/assets"
	"ark-server-commander/config"
	"ark-server-commander/models"
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
		_, root, err := aseApi.readZip(bytes.NewReader(data), int64(len(data)))
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
		if got := aseApi.versionFromZip(in[0], in[1]); got != want {
			t.Fatalf("aseApi.versionFromZip(%q, %q) = %q, want %q", in[0], in[1], got, want)
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
	got := aseApi.parseReleases(body)
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
	if err := aseApi.install(id, v1, "1.0"); err != nil {
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
	if v, installed := installedArkApiVersion(models.Server{ID: id}); v != "1.0" || !installed {
		t.Fatalf("installed = %q %v", v, installed)
	}

	// 更新：替换 dll，保留用户改过的 config.json
	os.WriteFile(filepath.Join(win64, "config.json"), []byte(`{"mine":1}`), 0o644)
	v2 := zipBytes(t, map[string]string{"version.dll": "v2", "config.json": `{"default":true}`})
	if err := aseApi.install(id, v2, "2.0"); err != nil {
		t.Fatal(err)
	}
	if b, _ := os.ReadFile(filepath.Join(win64, "version.dll")); string(b) != "v2" {
		t.Fatalf("version.dll = %q", b)
	}
	if b, _ := os.ReadFile(filepath.Join(win64, "config.json")); string(b) != `{"mine":1}` {
		t.Fatalf("config.json should be kept, got %q", b)
	}
	if v, _ := installedArkApiVersion(models.Server{ID: id}); v != "2.0" {
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
	if err := ensureArkApiInstalled(models.Server{ID: id}); err != nil {
		t.Fatal(err)
	}
	if v, installed := installedArkApiVersion(models.Server{ID: id}); v != assets.BuiltinArkApiVersion || !installed {
		t.Fatalf("installed = %q %v", v, installed)
	}

	// 已安装时不动
	os.WriteFile(filepath.Join(docker_manager.ServerFilesDir(id), arkApiVersionFile), []byte("old\n"), 0o644)
	if err := ensureArkApiInstalled(models.Server{ID: id}); err != nil {
		t.Fatal(err)
	}
	if v, _ := installedArkApiVersion(models.Server{ID: id}); v != "old" {
		t.Fatalf("should not reinstall, got %q", v)
	}

	// 版本库里有更新的版本时用它
	if err := aseApi.saveZip("99.1", zipBytes(t, map[string]string{"version.dll": "new"})); err != nil {
		t.Fatal(err)
	}
	const other = 4
	if err := ensureArkApiInstalled(models.Server{ID: other}); err != nil {
		t.Fatal(err)
	}
	if v, _ := installedArkApiVersion(models.Server{ID: other}); v != "99.1" {
		t.Fatalf("expected 99.1, got %q", v)
	}
}

func TestArkApiLibrary(t *testing.T) {
	useTempDataDir(t)
	s := NewServerService()

	data := zipBytes(t, map[string]string{"version.dll": "x"})
	version, err := s.UploadArkApi("", bytes.NewReader(data), int64(len(data)), "ArkApi_3.60.zip", "")
	if err != nil || version != "3.60" {
		t.Fatalf("upload: %q %v", version, err)
	}
	// 识别不出版本号时要求填写
	if _, err := s.UploadArkApi("", bytes.NewReader(data), int64(len(data)), "upload.zip", ""); !errors.Is(err, ErrInvalidArkApi) {
		t.Fatalf("expected ErrInvalidArkApi, got %v", err)
	}
	if _, err := s.UploadArkApi("", bytes.NewReader(data), int64(len(data)), "upload.zip", "../x"); !errors.Is(err, ErrInvalidArkApi) {
		t.Fatalf("expected ErrInvalidArkApi for bad version, got %v", err)
	}

	local := aseApi.localVersions()
	if !local["3.60"].Local || !local["3.60"].Deletable || !local[assets.BuiltinArkApiVersion].Builtin {
		t.Fatalf("local versions: %+v", local)
	}
	if aseApi.latestLocal() != "3.60" {
		t.Fatalf("latest local = %q", aseApi.latestLocal())
	}

	if err := s.DeleteArkApiVersion("", assets.BuiltinArkApiVersion); !errors.Is(err, ErrArkApiBuiltin) {
		t.Fatalf("expected ErrArkApiBuiltin, got %v", err)
	}
	if err := s.DeleteArkApiVersion("", "3.60"); err != nil {
		t.Fatal(err)
	}
	if err := s.DeleteArkApiVersion("", "3.60"); !errors.Is(err, ErrArkApiNotFound) {
		t.Fatalf("expected ErrArkApiNotFound, got %v", err)
	}
}

func TestBuiltinArkApiZip(t *testing.T) {
	zr, err := zip.NewReader(bytes.NewReader(assets.BuiltinArkApiZip), int64(len(assets.BuiltinArkApiZip)))
	if err != nil {
		t.Fatal(err)
	}
	_, root, err := aseApi.readZip(bytes.NewReader(assets.BuiltinArkApiZip), int64(len(assets.BuiltinArkApiZip)))
	if err != nil || root != "" {
		t.Fatalf("root %q err %v", root, err)
	}
	if len(zr.File) < 2 {
		t.Fatalf("unexpected builtin zip: %d files", len(zr.File))
	}
}

func TestAsaApiInstall(t *testing.T) {
	useTempDataDir(t)
	// 不查询 GitHub：发布列表为空
	asaApi.releases.Lock()
	asaApi.releases.list, asaApi.releases.err, asaApi.releases.checked = nil, "", time.Now()
	asaApi.releases.Unlock()
	t.Cleanup(func() {
		asaApi.releases.Lock()
		asaApi.releases.checked = time.Time{}
		asaApi.releases.Unlock()
	})
	asa := models.Server{ID: 11, Game: utils.GameASA}

	// AsaApi 没有内置版本，本地和 GitHub 上都没有时启动前安装失败
	if err := ensureArkApiInstalled(asa); !errors.Is(err, ErrArkApiNotFound) {
		t.Fatalf("expected ErrArkApiNotFound, got %v", err)
	}
	if v := asaApi.localVersions(); len(v) != 0 {
		t.Fatalf("AsaApi should have no builtin version, got %+v", v)
	}

	// ArkApi 的压缩包不能当 AsaApi 上传，反过来也不行
	s := NewServerService()
	ark := zipBytes(t, map[string]string{"version.dll": "x"})
	if _, err := s.UploadArkApi(utils.GameASA, bytes.NewReader(ark), int64(len(ark)), "AsaApi_2.11.zip", ""); !errors.Is(err, ErrInvalidArkApi) {
		t.Fatalf("expected ErrInvalidArkApi for an ArkApi zip, got %v", err)
	}
	data := zipBytes(t, map[string]string{
		"AsaApi_2.11/AsaApiLoader.exe":             "loader",
		"AsaApi_2.11/version.dll":                  "inject",
		"AsaApi_2.11/config.json":                  `{"asa":true}`,
		"AsaApi_2.11/ArkApi/Plugins/Perm/Perm.dll": "plugin",
	})
	if _, err := s.UploadArkApi("", bytes.NewReader(data), int64(len(data)), "AsaApi_2.11.zip", "9.9"); !errors.Is(err, ErrInvalidArkApi) {
		t.Fatalf("expected ErrInvalidArkApi for an AsaApi zip uploaded as ArkApi, got %v", err)
	}
	version, err := s.UploadArkApi(utils.GameASA, bytes.NewReader(data), int64(len(data)), "AsaApi_2.11.zip", "")
	if err != nil || version != "2.11" {
		t.Fatalf("upload = %q %v", version, err)
	}
	// 两个版本库分开
	if _, ok := aseApi.localVersions()["2.11"]; ok {
		t.Fatal("AsaApi upload must not appear in the ArkApi library")
	}

	// 启动前安装本地最新版本，标志文件是 AsaApiLoader.exe
	if err := ensureArkApiInstalled(asa); err != nil {
		t.Fatal(err)
	}
	win64 := arkApiWin64Dir(asa.ID)
	if b, _ := os.ReadFile(filepath.Join(win64, "AsaApiLoader.exe")); string(b) != "loader" {
		t.Fatalf("AsaApiLoader.exe = %q", b)
	}
	if _, err := os.Stat(filepath.Join(win64, "ArkApi", "Plugins", "Perm")); !os.IsNotExist(err) {
		t.Fatal("plugins inside the AsaApi zip should be skipped")
	}
	if v, installed := installedArkApiVersion(asa); v != "2.11" || !installed {
		t.Fatalf("installed = %q %v", v, installed)
	}
	// ASE 服务器只看 version.dll：同样的文件卷里有 version.dll 也算装了 ArkApi，没有 AsaApiLoader.exe 的 ASA 服务器不算
	os.Remove(filepath.Join(win64, "AsaApiLoader.exe"))
	if _, installed := installedArkApiVersion(asa); installed {
		t.Fatal("ASA server without AsaApiLoader.exe should not count as installed")
	}

	list, err := s.ListArkApiVersions(utils.GameASA, false)
	if err != nil || list.Name != "AsaApi" || list.Latest != "2.11" {
		t.Fatalf("list = %+v %v", list, err)
	}
	if _, err := s.ListArkApiVersions("xbox", false); !errors.Is(err, ErrInvalidGame) {
		t.Fatalf("expected ErrInvalidGame, got %v", err)
	}
	if err := s.DeleteArkApiVersion(utils.GameASA, "2.11"); err != nil {
		t.Fatal(err)
	}
}
