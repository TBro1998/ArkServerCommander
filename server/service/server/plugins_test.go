package server

import (
	"archive/zip"
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"ark-server-commander/config"
	"ark-server-commander/database"
	"ark-server-commander/models"
	"ark-server-commander/service/docker_manager"
)

func makeZip(t *testing.T, files map[string]string) *bytes.Reader {
	t.Helper()
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	for name, body := range files {
		w, err := zw.Create(name)
		if err != nil {
			t.Fatal(err)
		}
		w.Write([]byte(body))
	}
	if err := zw.Close(); err != nil {
		t.Fatal(err)
	}
	return bytes.NewReader(buf.Bytes())
}

func zipFiles(t *testing.T, files map[string]string) []*zip.File {
	t.Helper()
	r := makeZip(t, files)
	zr, err := zip.NewReader(r, r.Size())
	if err != nil {
		t.Fatal(err)
	}
	return zr.File
}

func TestPluginRoots(t *testing.T) {
	cases := []struct {
		name    string
		files   map[string]string
		input   string
		want    map[string]string // 插件名 -> 前缀
		wantErr bool
	}{
		{
			name:  "目录与 dll 同名",
			files: map[string]string{"ArkShop/ArkShop.dll": "x", "ArkShop/config.json": "{}", "ArkShop/lib/mysql.dll": "x"},
			want:  map[string]string{"ArkShop": "ArkShop/"},
		},
		{
			name:  "带前缀的多个插件",
			files: map[string]string{"ArkApi/Plugins/ArkShop/ArkShop.dll": "x", "ArkApi/Plugins/Permissions/Permissions.dll": "x"},
			want:  map[string]string{"ArkShop": "ArkApi/Plugins/ArkShop/", "Permissions": "ArkApi/Plugins/Permissions/"},
		},
		{
			name:  "文件直接放在根目录，只有一个 dll",
			files: map[string]string{"AAC.dll": "x", "config.json": "{}"},
			want:  map[string]string{"AAC": ""},
		},
		{
			name:    "根目录多个 dll 无法识别",
			files:   map[string]string{"A.dll": "x", "B.dll": "x"},
			wantErr: true,
		},
		{
			name:  "根目录多个 dll，指定名称",
			files: map[string]string{"A.dll": "x", "B.dll": "x"},
			input: "B",
			want:  map[string]string{"B": ""},
		},
		{
			name:    "路径穿越",
			files:   map[string]string{"../evil/evil.dll": "x"},
			wantErr: true,
		},
		{
			name:  "Windows 打包的 zip 用反斜杠分隔",
			files: map[string]string{"MissionDisabler\\MissionDisabler.dll": "x", "MissionDisabler\\config.json": "{}"},
			want:  map[string]string{"MissionDisabler": "MissionDisabler/"},
		},
		{
			name:    "反斜杠的路径穿越",
			files:   map[string]string{"..\\evil\\evil.dll": "x"},
			wantErr: true,
		},
		{
			name:    "没有 dll",
			files:   map[string]string{"readme.txt": "x"},
			wantErr: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			roots, err := pluginRoots(zipFiles(t, c.files), c.input)
			if c.wantErr {
				if !errors.Is(err, ErrInvalidPlugin) {
					t.Fatalf("expected ErrInvalidPlugin, got %v", err)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			got := map[string]string{}
			for _, r := range roots {
				got[r.name] = r.prefix
			}
			if fmt.Sprint(got) != fmt.Sprint(c.want) {
				t.Fatalf("got %v, want %v", got, c.want)
			}
		})
	}
}

func TestPluginLifecycle(t *testing.T) {
	setupTestDB(t)
	old := config.DataDir
	config.DataDir = t.TempDir()
	t.Cleanup(func() { config.DataDir = old })

	srv := createTestServer(t, 7777, 27015, 32330)
	id := fmt.Sprint(srv.ID)
	s := NewServerService()

	z := makeZip(t, map[string]string{
		"ArkShop/ArkShop.dll":         "v1",
		"ArkShop/config.json":         `{"default":true}`,
		"ArkShop/PluginInfo.json":     "\xEF\xBB\xBF" + `{"FullName":"Ark Shop","Version":1.5,"MinApiVersion":3.5,"Dependencies":["Permissions"]}`,
		"Permissions/Permissions.dll": "x",
	})
	names, err := s.InstallPlugins(1, id, z, z.Size(), "")
	if err != nil {
		t.Fatal(err)
	}
	if len(names) != 2 {
		t.Fatalf("expected 2 plugins, got %v", names)
	}

	overview, err := s.ListPlugins(1, id)
	if err != nil {
		t.Fatal(err)
	}
	if len(overview.Plugins) != 2 || overview.Plugins[0].Name != "ArkShop" {
		t.Fatalf("unexpected plugins %+v", overview.Plugins)
	}
	shop := overview.Plugins[0]
	if shop.FullName != "Ark Shop" || shop.Version != "1.5" || shop.MinAPIVersion != "3.5" || !shop.HasDLL || !shop.HasConfig || !shop.Enabled {
		t.Fatalf("unexpected plugin info %+v", shop)
	}

	// 修改配置后覆盖安装，配置保留
	if err := s.SavePluginConfig(1, id, "ArkShop", `{"mine":1}`); err != nil {
		t.Fatal(err)
	}
	if err := s.SavePluginConfig(1, id, "ArkShop", `{bad`); !errors.Is(err, ErrInvalidPluginConfig) {
		t.Fatalf("expected ErrInvalidPluginConfig, got %v", err)
	}
	z = makeZip(t, map[string]string{"ArkShop/ArkShop.dll": "v2", "ArkShop/config.json": `{"default":true}`})
	if _, err := s.InstallPlugins(1, id, z, z.Size(), ""); err != nil {
		t.Fatal(err)
	}
	dir := filepath.Join(docker_manager.ServerPluginsDir(srv.ID), "ArkShop")
	if data, _ := os.ReadFile(filepath.Join(dir, "ArkShop.dll")); string(data) != "v2" {
		t.Fatalf("dll not upgraded: %q", data)
	}
	if content, _ := s.GetPluginConfig(1, id, "ArkShop"); content != `{"mine":1}` {
		t.Fatalf("config not kept: %q", content)
	}

	// 停用后移出插件卷，再启用移回来
	if err := s.SetPluginEnabled(1, id, "ArkShop", false); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(dir); !errors.Is(err, os.ErrNotExist) {
		t.Fatal("disabled plugin still in plugins volume")
	}
	overview, _ = s.ListPlugins(1, id)
	if overview.Plugins[0].Enabled {
		t.Fatal("expected plugin to be disabled")
	}
	if err := s.SetPluginEnabled(1, id, "ArkShop", true); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(dir); err != nil {
		t.Fatal("enabled plugin not restored")
	}

	// 运行中不能删除
	database.DB.Model(&models.Server{}).Where("id = ?", srv.ID).Update("status", "running")
	if err := s.DeletePlugin(1, id, "ArkShop"); !errors.Is(err, ErrServerNotStopped) {
		t.Fatalf("expected ErrServerNotStopped, got %v", err)
	}
	database.DB.Model(&models.Server{}).Where("id = ?", srv.ID).Update("status", "stopped")
	if err := s.DeletePlugin(1, id, "ArkShop"); err != nil {
		t.Fatal(err)
	}
	if err := s.DeletePlugin(1, id, "../saved"); !errors.Is(err, ErrPluginNotFound) {
		t.Fatalf("expected ErrPluginNotFound, got %v", err)
	}
	overview, _ = s.ListPlugins(1, id)
	if len(overview.Plugins) != 1 || overview.Plugins[0].Name != "Permissions" {
		t.Fatalf("unexpected plugins after delete %+v", overview.Plugins)
	}

	// ArkApi 还没有下载时不能编辑它的配置
	if _, err := s.GetPluginConfig(1, id, ""); !errors.Is(err, ErrArkApiNotInstalled) {
		t.Fatalf("expected ErrArkApiNotInstalled, got %v", err)
	}
}

func TestInstallPluginsFromURL(t *testing.T) {
	setupTestDB(t)
	old := config.DataDir
	config.DataDir = t.TempDir()
	t.Cleanup(func() { config.DataDir = old })

	srv := createTestServer(t, 7777, 27015, 32330)
	id := fmt.Sprint(srv.ID)
	s := NewServerService()

	z := makeZip(t, map[string]string{"ArkShop/ArkShop.dll": "v1", "ArkShop/config.json": `{"a":1}`})
	zipData, err := io.ReadAll(z)
	if err != nil {
		t.Fatal(err)
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/plugin.zip", func(w http.ResponseWriter, r *http.Request) {
		w.Write(zipData)
	})
	mux.HandleFunc("/missing.zip", func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	})
	mux.HandleFunc("/notzip.zip", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello"))
	})
	ts := httptest.NewServer(mux)
	defer ts.Close()

	// 只接受 http / https
	if _, err := s.InstallPluginsFromURL(1, id, "file:///etc/passwd", ""); !errors.Is(err, ErrInvalidPlugin) {
		t.Fatalf("expected ErrInvalidPlugin, got %v", err)
	}
	if _, err := s.InstallPluginsFromURL(1, id, "  ", ""); !errors.Is(err, ErrInvalidPlugin) {
		t.Fatalf("expected ErrInvalidPlugin, got %v", err)
	}

	// 404 和内容不是 zip 时报错，且不会创建插件
	if _, err := s.InstallPluginsFromURL(1, id, ts.URL+"/missing.zip", ""); !errors.Is(err, ErrPluginDownload) {
		t.Fatalf("expected ErrPluginDownload, got %v", err)
	}
	if _, err := s.InstallPluginsFromURL(1, id, ts.URL+"/notzip.zip", ""); !errors.Is(err, ErrInvalidPlugin) {
		t.Fatalf("expected ErrInvalidPlugin, got %v", err)
	}

	// 运行中不能安装，且不下载
	database.DB.Model(&models.Server{}).Where("id = ?", srv.ID).Update("status", "running")
	if _, err := s.InstallPluginsFromURL(1, id, ts.URL+"/plugin.zip", ""); !errors.Is(err, ErrServerNotStopped) {
		t.Fatalf("expected ErrServerNotStopped, got %v", err)
	}
	database.DB.Model(&models.Server{}).Where("id = ?", srv.ID).Update("status", "stopped")

	names, err := s.InstallPluginsFromURL(1, id, ts.URL+"/plugin.zip", "")
	if err != nil {
		t.Fatal(err)
	}
	if fmt.Sprint(names) != "[ArkShop]" {
		t.Fatalf("unexpected installed plugins %v", names)
	}
	overview, err := s.ListPlugins(1, id)
	if err != nil {
		t.Fatal(err)
	}
	if len(overview.Plugins) != 1 || overview.Plugins[0].Name != "ArkShop" || !overview.Plugins[0].HasDLL {
		t.Fatalf("unexpected plugins %+v", overview.Plugins)
	}

	// 覆盖安装保留配置
	if err := s.SavePluginConfig(1, id, "ArkShop", `{"mine":2}`); err != nil {
		t.Fatal(err)
	}
	if _, err := s.InstallPluginsFromURL(1, id, ts.URL+"/plugin.zip", ""); err != nil {
		t.Fatal(err)
	}
	if content, _ := s.GetPluginConfig(1, id, "ArkShop"); content != `{"mine":2}` {
		t.Fatalf("config not kept: %q", content)
	}
}

// Windows 上打包、用反斜杠分隔的 zip 解压到正确的目录结构
func TestInstallPluginWithBackslashPaths(t *testing.T) {
	setupTestDB(t)
	old := config.DataDir
	config.DataDir = t.TempDir()
	t.Cleanup(func() { config.DataDir = old })
	srv := createTestServer(t, 7777, 27015, 32330)
	s := NewServerService()

	z := makeZip(t, map[string]string{
		"AdvancedMessages\\AdvancedMessages.dll": "dll",
		"AdvancedMessages\\config.json":          "{}",
		"AdvancedMessages\\lib\\helper.dll":   "lib",
	})
	names, err := s.InstallPlugins(1, fmt.Sprint(srv.ID), z, z.Size(), "")
	if err != nil {
		t.Fatal(err)
	}
	if fmt.Sprint(names) != "[AdvancedMessages]" {
		t.Fatalf("unexpected names %v", names)
	}
	dir := filepath.Join(docker_manager.ServerPluginsDir(srv.ID), "AdvancedMessages")
	for _, f := range []string{"AdvancedMessages.dll", "config.json", filepath.Join("lib", "helper.dll")} {
		if _, err := os.Stat(filepath.Join(dir, f)); err != nil {
			t.Fatalf("missing %s: %v", f, err)
		}
	}
}
