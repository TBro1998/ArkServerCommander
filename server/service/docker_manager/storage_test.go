package docker_manager

import (
	"os"
	"path/filepath"
	"testing"

	"ark-server-commander/config"
)

func TestMapHostPath(t *testing.T) {
	mounts := []hostMount{
		{local: "/data", host: "E:/ArkServer/data"},
		{local: "/server-files", host: "/var/lib/docker/volumes/ark-server-files/_data"},
		{local: "/data/clusters", host: "/var/lib/docker/volumes/ark-clusters/_data"},
		{local: "/var/run/docker.sock", host: "/var/run/docker.sock"},
	}
	cases := map[string]string{
		"/data":                 "E:/ArkServer/data",
		"/data/servers/1/saved": "E:/ArkServer/data/servers/1/saved",
		"/data/link":            "E:/ArkServer/data/link",
		"/server-files/1":       "/var/lib/docker/volumes/ark-server-files/_data/1",
		"/data/clusters/abc":    "/var/lib/docker/volumes/ark-clusters/_data/abc",
		"/data/clustersx/abc":   "E:/ArkServer/data/clustersx/abc",
	}
	for local, want := range cases {
		got, ok := mapHostPath(mounts, local)
		if !ok || got != want {
			t.Errorf("mapHostPath(%q) = %q, %v, want %q", local, got, ok, want)
		}
	}
	if _, ok := mapHostPath(mounts, "/elsewhere/1"); ok {
		t.Error("path outside every mount should not map")
	}
	if _, ok := mapHostPath(mounts, "/server-filesx/1"); ok {
		t.Error("prefix without a path separator should not match")
	}
}

func TestServerFilesDirWithFilesVolume(t *testing.T) {
	oldData, oldRoot := config.DataDir, filesRoot
	t.Cleanup(func() {
		config.DataDir, filesRoot = oldData, oldRoot
		filesDirCache.Clear()
	})
	config.DataDir = t.TempDir()
	filesRoot = ""
	filesDirCache.Clear()
	legacy := filepath.Join(config.DataDir, "servers", "3", "files")
	if got := ServerFilesDir(3); got != legacy {
		t.Errorf("without files volume: got %q, want %q", got, legacy)
	}

	filesRoot = t.TempDir()
	filesDirCache.Clear()
	// 新服务器（数据目录里没有服务端文件，空目录也算没有）放在单独的卷里
	if err := os.MkdirAll(legacy, 0o755); err != nil {
		t.Fatal(err)
	}
	if got, want := ServerFilesDir(3), filepath.Join(filesRoot, "3"); got != want {
		t.Errorf("new server: got %q, want %q", got, want)
	}
	specs := serverVolumeSpecs(3)
	if specs[0].localDir() != filepath.Join(filesRoot, "3") || specs[1].localDir() != ServerSavedDir(3) || specs[2].localDir() != ServerPluginsDir(3) {
		t.Errorf("unexpected volume specs: %+v", specs)
	}

	// 启用之前已经下载到数据目录的服务器继续用原位置
	old := filepath.Join(config.DataDir, "servers", "4", "files")
	if err := os.MkdirAll(filepath.Join(old, "ShooterGame"), 0o755); err != nil {
		t.Fatal(err)
	}
	if got := ServerFilesDir(4); got != old {
		t.Errorf("legacy server: got %q, want %q", got, old)
	}

	// 删除服务器数据时两处都删
	newDir := filepath.Join(filesRoot, "3")
	if err := os.MkdirAll(newDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := RemoveServerDirs(3); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(newDir); !os.IsNotExist(err) {
		t.Errorf("files volume dir should be removed, stat err = %v", err)
	}
	if _, err := os.Stat(filepath.Join(config.DataDir, "servers", "3")); !os.IsNotExist(err) {
		t.Errorf("data dir should be removed, stat err = %v", err)
	}
}

func TestWantFilesVolume(t *testing.T) {
	cases := []struct {
		mode          string
		desktop, bind bool
		want          bool
	}{
		{"auto", true, true, true},
		{"auto", true, false, false},
		{"auto", false, true, false},
		{"", true, true, true},
		{"on", false, false, true},
		{"OFF", true, true, false},
	}
	for _, c := range cases {
		if got := wantFilesVolume(c.mode, c.desktop, c.bind); got != c.want {
			t.Errorf("wantFilesVolume(%q, %v, %v) = %v, want %v", c.mode, c.desktop, c.bind, got, c.want)
		}
	}
}

func TestVolumeInUse(t *testing.T) {
	mp := "/var/lib/docker/volumes/abc/_data"
	devices := []string{"E:/data/servers/1/saved", "/var/lib/docker/volumes/abcd/_data/2"}
	if volumeInUse(mp, devices) {
		t.Error("no device under the old volume")
	}
	if !volumeInUse(mp, append(devices, mp+"/1")) {
		t.Error("device under the old volume should keep it")
	}
	if !volumeInUse("", devices) {
		t.Error("unknown mountpoint should be treated as in use")
	}
}
