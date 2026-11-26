package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"ark-server-commander/config"
	"ark-server-commander/database"
	"ark-server-commander/models"
	"ark-server-commander/service/docker_manager"
	"ark-server-commander/utils"
)

// fakeCatalogSite 同时扮演 arkfgf.com（目录接口）和 ark-server-api.com（下载地址）
type fakeCatalogSite struct {
	mu       sync.Mutex
	requests int
	fail     bool
	version  string
	// oldSite 模拟还不认识 game 参数的网站：总是返回 ASE 的目录，也没有 game 字段
	oldSite bool
	ts      *httptest.Server
}

func newFakeCatalogSite(t *testing.T) *fakeCatalogSite {
	t.Helper()
	f := &fakeCatalogSite{version: "3.05"}
	shop := zipBytes(t, map[string]string{"ArkShop/ArkShop.dll": "dll", "ArkShop/config.json": `{"default":true}`})
	api := zipBytes(t, map[string]string{"ArkApi/version.dll": "x", "ArkApi/config.json": "{}"})
	asaShop := zipBytes(t, map[string]string{"ArkShop/ArkShop.dll": "asa dll", "ArkShop/config.json": `{"asa":true}`})
	asaAPI := zipBytes(t, map[string]string{"AsaApiLoader.exe": "x", "ArkApi/AsaApi.dll": "x", "config.json": "{}"})
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/plugins/catalog", func(w http.ResponseWriter, r *http.Request) {
		f.mu.Lock()
		defer f.mu.Unlock()
		f.requests++
		if f.fail {
			w.WriteHeader(http.StatusServiceUnavailable)
			fmt.Fprint(w, `{"error":"x","code":"catalog_unavailable"}`)
			return
		}
		base := f.ts.URL
		if r.URL.Query().Get("game") == "asa" && !f.oldSite {
			json.NewEncoder(w).Encode(PluginCatalog{
				Game:        "asa",
				Source:      "ark-server-api.com",
				GeneratedAt: time.Now().UTC(),
				Plugins: []CatalogPlugin{
					{ID: 34, Title: "ASA Arkshop", Category: "official", Version: "1.81", PageURL: base + "/resources/asa-arkshop.34/",
						DownloadURL: base + "/resources/asa-arkshop.34/download", Plugins: []CatalogPluginInfo{{Name: "ArkShop"}}},
				},
				ArkApi: &CatalogArkApi{Version: "2.99", DownloadURL: base + "/resources/asa-server-api.31/download", PageURL: base + "/resources/asa-server-api.31/"},
			})
			return
		}
		game := "ase"
		if f.oldSite {
			game = ""
		}
		json.NewEncoder(w).Encode(PluginCatalog{
			Game:        game,
			Source:      "ark-server-api.com",
			GeneratedAt: time.Now().UTC(),
			Plugins: []CatalogPlugin{
				{ID: 36, Title: "ASE ArkShop", Category: "official", Version: f.version, PageURL: base + "/resources/ase-arkshop.36/",
					DownloadURL: base + "/resources/ase-arkshop.36/download", Plugins: []CatalogPluginInfo{{Name: "ArkShop", Dependencies: []string{"Permissions"}}}},
				{ID: 191, Title: "Paid", Category: "paid", PageURL: base + "/resources/paid.191/"},
				{ID: 77, Title: "Elsewhere", Category: "free", PageURL: "https://example.com/", DownloadURL: "https://example.com/x.zip"},
			},
			ArkApi: &CatalogArkApi{Version: "9.99", DownloadURL: base + "/resources/ase-server-api.32/download", PageURL: base + "/resources/ase-server-api.32/"},
		})
	})
	mux.HandleFunc("/resources/ase-arkshop.36/download", func(w http.ResponseWriter, r *http.Request) { w.Write(shop) })
	mux.HandleFunc("/resources/ase-server-api.32/download", func(w http.ResponseWriter, r *http.Request) { w.Write(api) })
	mux.HandleFunc("/resources/asa-arkshop.34/download", func(w http.ResponseWriter, r *http.Request) { w.Write(asaShop) })
	mux.HandleFunc("/resources/asa-server-api.31/download", func(w http.ResponseWriter, r *http.Request) { w.Write(asaAPI) })
	f.ts = httptest.NewServer(mux)
	t.Cleanup(f.ts.Close)

	oldAllow := allowCatalogDownload
	host := strings.TrimPrefix(f.ts.URL, "http://")
	allowCatalogDownload = func(u *url.URL) bool { return u.Host == host }
	t.Cleanup(func() { allowCatalogDownload = oldAllow })
	for _, c := range []*pluginCatalogCache{asePluginCatalog, asaPluginCatalog} {
		oldAPI := c.api
		c.api = func() string { return f.ts.URL + "/api/v1" }
		c.forget()
		t.Cleanup(func() {
			c.api = oldAPI
			c.forget()
		})
	}
	return f
}

func TestPluginCatalogCache(t *testing.T) {
	f := newFakeCatalogSite(t)
	s := NewServerService()

	c, err := s.GetPluginCatalog("", false)
	if err != nil {
		t.Fatal(err)
	}
	if len(c.Plugins) != 3 || c.Source != "ark-server-api.com" || c.Game != "ase" {
		t.Fatalf("unexpected catalog %+v", c)
	}
	// 缓存内不再请求，refresh 时重新获取
	s.GetPluginCatalog("", false)
	if f.requests != 1 {
		t.Fatalf("expected 1 request, got %d", f.requests)
	}
	s.GetPluginCatalog("", true)
	if f.requests != 2 {
		t.Fatalf("expected refresh to fetch again, got %d", f.requests)
	}
	// 网站出错时返回旧数据
	f.mu.Lock()
	f.fail = true
	f.mu.Unlock()
	if c, err := s.GetPluginCatalog("", true); err != nil || len(c.Plugins) != 3 {
		t.Fatalf("expected stale catalog, got %v %v", c, err)
	}
	// 没有旧数据时报错
	asePluginCatalog.forget()
	if _, err := s.GetPluginCatalog("", false); !errors.Is(err, ErrPluginCatalogUnavailable) {
		t.Fatalf("expected ErrPluginCatalogUnavailable, got %v", err)
	}
	if _, err := s.GetPluginCatalog("atlas", false); !errors.Is(err, ErrInvalidGame) {
		t.Fatalf("expected ErrInvalidGame, got %v", err)
	}
	// ARKFGF_API 关闭
	asePluginCatalog.api = func() string { return "" }
	if _, err := s.GetPluginCatalog("", false); !errors.Is(err, ErrArkfgfDisabled) {
		t.Fatalf("expected ErrArkfgfDisabled, got %v", err)
	}
}

func TestInstallCatalogPlugin(t *testing.T) {
	setupTestDB(t)
	old := config.DataDir
	config.DataDir = t.TempDir()
	t.Cleanup(func() { config.DataDir = old })
	f := newFakeCatalogSite(t)

	srv := createTestServer(t, 7777, 27015, 32330)
	id := fmt.Sprint(srv.ID)
	s := NewServerService()

	if _, err := s.InstallCatalogPlugin(1, id, 404); !errors.Is(err, ErrCatalogPluginNotFound) {
		t.Fatalf("expected ErrCatalogPluginNotFound, got %v", err)
	}
	if _, err := s.InstallCatalogPlugin(1, id, 191); !errors.Is(err, ErrCatalogPluginNotInstallable) {
		t.Fatalf("paid: expected ErrCatalogPluginNotInstallable, got %v", err)
	}
	// 下载地址不在允许的站点上
	if _, err := s.InstallCatalogPlugin(1, id, 77); !errors.Is(err, ErrCatalogPluginNotInstallable) {
		t.Fatalf("other host: expected ErrCatalogPluginNotInstallable, got %v", err)
	}
	database.DB.Model(&models.Server{}).Where("id = ?", srv.ID).Update("status", "running")
	if _, err := s.InstallCatalogPlugin(1, id, 36); !errors.Is(err, ErrServerNotStopped) {
		t.Fatalf("expected ErrServerNotStopped, got %v", err)
	}
	database.DB.Model(&models.Server{}).Where("id = ?", srv.ID).Update("status", "stopped")

	names, err := s.InstallCatalogPlugin(1, id, 36)
	if err != nil {
		t.Fatal(err)
	}
	if fmt.Sprint(names) != "[ArkShop]" {
		t.Fatalf("unexpected names %v", names)
	}
	overview, _ := s.ListPlugins(1, id)
	if len(overview.Plugins) != 1 {
		t.Fatalf("unexpected plugins %+v", overview.Plugins)
	}
	src := overview.Plugins[0].Source
	if src == nil || src.ResourceID != 36 || src.Version != "3.05" || src.Source != "ark-server-api.com" {
		t.Fatalf("unexpected source %+v", src)
	}

	// 停用后更新：保留配置和停用状态，来源版本更新
	if err := s.SavePluginConfig(1, id, "ArkShop", `{"mine":1}`); err != nil {
		t.Fatal(err)
	}
	if err := s.SetPluginEnabled(1, id, "ArkShop", false); err != nil {
		t.Fatal(err)
	}
	f.mu.Lock()
	f.version = "3.06"
	f.mu.Unlock()
	if _, err := s.GetPluginCatalog("", true); err != nil {
		t.Fatal(err)
	}
	if _, err := s.InstallCatalogPlugin(1, id, 36); err != nil {
		t.Fatal(err)
	}
	overview, _ = s.ListPlugins(1, id)
	p := overview.Plugins[0]
	if p.Enabled || p.Source == nil || p.Source.Version != "3.06" {
		t.Fatalf("unexpected plugin after update %+v (source %+v)", p, p.Source)
	}
	if content, _ := s.GetPluginConfig(1, id, "ArkShop"); content != `{"mine":1}` {
		t.Fatalf("config not kept: %q", content)
	}

	// 上传覆盖安装后不再算目录安装的
	z := makeZip(t, map[string]string{"ArkShop/ArkShop.dll": "own build"})
	if _, err := s.InstallPlugins(1, id, z, z.Size(), ""); err != nil {
		t.Fatal(err)
	}
	overview, _ = s.ListPlugins(1, id)
	if overview.Plugins[0].Source != nil {
		t.Fatal("uploaded plugin should have no catalog source")
	}
	if _, err := os.Stat(filepath.Join(docker_manager.ServerDisabledPluginsDir(srv.ID), "ArkShop", pluginSourceFile)); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("source file should be gone, got %v", err)
	}
}

func TestArkApiFromCatalog(t *testing.T) {
	setupTestDB(t)
	old := config.DataDir
	config.DataDir = t.TempDir()
	t.Cleanup(func() { config.DataDir = old })
	newFakeCatalogSite(t)
	// 不查询 GitHub
	aseApi.releases.Lock()
	oldList, oldErr, oldChecked := aseApi.releases.list, aseApi.releases.err, aseApi.releases.checked
	aseApi.releases.list, aseApi.releases.err, aseApi.releases.checked = nil, "", time.Now()
	aseApi.releases.Unlock()
	t.Cleanup(func() {
		aseApi.releases.Lock()
		aseApi.releases.list, aseApi.releases.err, aseApi.releases.checked = oldList, oldErr, oldChecked
		aseApi.releases.Unlock()
	})

	s := NewServerService()
	versions, _ := s.ListArkApiVersions("", false)
	if versions.Latest != "9.99" {
		t.Fatalf("expected catalog version to be latest, got %+v", versions)
	}
	var found *ArkApiVersion
	for i := range versions.Versions {
		if versions.Versions[i].Version == "9.99" {
			found = &versions.Versions[i]
		}
	}
	if found == nil || !found.OnCatalog || found.Local {
		t.Fatalf("unexpected catalog version %+v", found)
	}

	srv := createTestServer(t, 7777, 27015, 32330)
	got, err := s.InstallArkApi(1, fmt.Sprint(srv.ID), "9.99")
	if err != nil {
		t.Fatal(err)
	}
	if v, installed := installedArkApiVersion(srv); got != "9.99" || v != "9.99" || !installed {
		t.Fatalf("unexpected installed version %q %q %v", got, v, installed)
	}
	// 下载后进了版本库
	if _, err := os.Stat(aseApi.libraryPath("9.99")); err != nil {
		t.Fatalf("catalog version not saved to library: %v", err)
	}
}

// ASA 的目录单独请求（?game=asa）、单独缓存；不认识 game 参数的旧网站返回的 ASE 目录不能用
func TestAsaPluginCatalog(t *testing.T) {
	f := newFakeCatalogSite(t)
	s := NewServerService()

	c, err := s.GetPluginCatalog("asa", false)
	if err != nil {
		t.Fatal(err)
	}
	if c.Game != "asa" || len(c.Plugins) != 1 || c.Plugins[0].ID != 34 {
		t.Fatalf("unexpected asa catalog %+v", c)
	}
	ase, err := s.GetPluginCatalog("ase", false)
	if err != nil || ase.Game != "ase" || len(ase.Plugins) != 3 {
		t.Fatalf("unexpected ase catalog %+v %v", ase, err)
	}
	if f.requests != 2 {
		t.Fatalf("expected one request per game, got %d", f.requests)
	}

	asaPluginCatalog.forget()
	f.mu.Lock()
	f.oldSite = true
	f.mu.Unlock()
	if _, err := s.GetPluginCatalog("asa", false); !errors.Is(err, ErrPluginCatalogUnavailable) {
		t.Fatalf("old site: expected ErrPluginCatalogUnavailable, got %v", err)
	}
	// ASE 不受影响（旧网站没有 game 字段）
	if c, err := s.GetPluginCatalog("", true); err != nil || len(c.Plugins) != 3 || c.Game != "ase" {
		t.Fatalf("old site ase catalog: %+v %v", c, err)
	}
}

// ASA 服务器从 ASA 目录安装；ASE 目录里的资源 ID 在 ASA 服务器上找不到
func TestInstallAsaCatalogPlugin(t *testing.T) {
	setupTestDB(t)
	old := config.DataDir
	config.DataDir = t.TempDir()
	t.Cleanup(func() { config.DataDir = old })
	newFakeCatalogSite(t)

	srv := createTestServer(t, 7777, 27015, 32330)
	database.DB.Model(&models.Server{}).Where("id = ?", srv.ID).Updates(map[string]any{"game": utils.GameASA, "query_port": 0})
	id := fmt.Sprint(srv.ID)
	s := NewServerService()

	if _, err := s.InstallCatalogPlugin(1, id, 36); !errors.Is(err, ErrCatalogPluginNotFound) {
		t.Fatalf("ase resource on asa server: expected ErrCatalogPluginNotFound, got %v", err)
	}
	names, err := s.InstallCatalogPlugin(1, id, 34)
	if err != nil {
		t.Fatal(err)
	}
	if fmt.Sprint(names) != "[ArkShop]" {
		t.Fatalf("unexpected names %v", names)
	}
	overview, _ := s.ListPlugins(1, id)
	if len(overview.Plugins) != 1 || overview.Plugins[0].Source == nil || overview.Plugins[0].Source.ResourceID != 34 ||
		overview.Plugins[0].Source.Version != "1.81" {
		t.Fatalf("unexpected plugins %+v", overview.Plugins)
	}
}

// ASA 目录里的 AsaApi 本体出现在 AsaApi 版本库里，安装时下载；ArkApi 版本库里没有它
func TestAsaApiFromCatalog(t *testing.T) {
	setupTestDB(t)
	old := config.DataDir
	config.DataDir = t.TempDir()
	t.Cleanup(func() { config.DataDir = old })
	newFakeCatalogSite(t)
	// 不查询 GitHub
	for _, a := range []*apiFlavor{aseApi, asaApi} {
		a.releases.Lock()
		oldList, oldErr, oldChecked := a.releases.list, a.releases.err, a.releases.checked
		a.releases.list, a.releases.err, a.releases.checked = nil, "", time.Now()
		a.releases.Unlock()
		t.Cleanup(func() {
			a.releases.Lock()
			a.releases.list, a.releases.err, a.releases.checked = oldList, oldErr, oldChecked
			a.releases.Unlock()
		})
	}

	s := NewServerService()
	versions, _ := s.ListArkApiVersions(utils.GameASA, false)
	if versions.Latest != "2.99" || len(versions.Versions) != 1 || !versions.Versions[0].OnCatalog {
		t.Fatalf("expected asa catalog version, got %+v", versions)
	}
	aseVersions, _ := s.ListArkApiVersions(utils.GameASE, false)
	for _, v := range aseVersions.Versions {
		if v.Version == "2.99" {
			t.Fatal("AsaApi listed in ArkApi library")
		}
	}

	srv := createTestServer(t, 7777, 27015, 32330)
	database.DB.Model(&models.Server{}).Where("id = ?", srv.ID).Updates(map[string]any{"game": utils.GameASA, "query_port": 0})
	got, err := s.InstallArkApi(1, fmt.Sprint(srv.ID), "2.99")
	if err != nil {
		t.Fatal(err)
	}
	if got != "2.99" {
		t.Fatalf("unexpected version %q", got)
	}
	if _, err := os.Stat(asaApi.libraryPath("2.99")); err != nil {
		t.Fatalf("catalog version not saved to library: %v", err)
	}
}
