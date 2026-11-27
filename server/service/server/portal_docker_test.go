package server

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"ark-server-commander/config"
	"ark-server-commander/database"
	"ark-server-commander/models"
	"ark-server-commander/service/docker_manager"
	"ark-server-commander/utils"
)

// TestPortalDockerSmoke 用真实的 Docker 启动玩家站容器（镜像 utils.PortalImage 要在本机已有，可以把自己构建的镜像打上这个标签），
// 写快照后经宿主机端口请求网页和接口，最后删除容器。需要本机 Docker，默认跳过：
//
//	ASC_PORTAL_DOCKER_TEST=1 go test ./service/server/ -run TestPortalDockerSmoke -v
func TestPortalDockerSmoke(t *testing.T) {
	if os.Getenv("ASC_PORTAL_DOCKER_TEST") != "1" {
		t.Skip("设置 ASC_PORTAL_DOCKER_TEST=1 才运行（需要 Docker）")
	}
	utils.InitLogger()
	setupTestDB(t)
	oldDir, oldHost := config.DataDir, config.HostDataDir
	config.DataDir, config.HostDataDir = t.TempDir(), ""
	if os.Getenv("ASC_PORTAL_KEEP") != "" {
		config.DataDir, _ = os.MkdirTemp("", "portal-smoke")
		t.Logf("数据目录 %s", config.DataDir)
	}
	t.Cleanup(func() { config.DataDir, config.HostDataDir = oldDir, oldHost })
	if err := docker_manager.InitStorage(); err != nil {
		t.Fatal(err)
	}
	ops := dockerPortalOps{}
	if ready, err := ops.ImageReady(utils.PortalImage); err != nil || !ready {
		t.Fatalf("本机没有镜像 %s（%v），先构建并打上这个标签", utils.PortalImage, err)
	}

	srv := createTestServer(t, 7777, 27015, 32330)
	database.DB.Model(&srv).Update("arkapi", true)
	shopDir := filepath.Join(docker_manager.ServerPluginsDir(srv.ID), "ArkShop")
	if err := os.MkdirAll(shopDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(shopDir, "config.json"), []byte(testArkShopConfig), 0o644); err != nil {
		t.Fatal(err)
	}
	const port = 38090
	if _, err := (&ServerService{}).UpdatePortal(models.PortalSettingsRequest{Enabled: true, Port: port, Name: "冒烟测试"}); err != nil {
		t.Fatal(err)
	}

	m := &PortalManager{
		status: models.PortalStatus{State: models.PortalDisabled}, dir: docker_manager.PortalDataDir, docker: ops,
		now: time.Now, wakeCh: make(chan struct{}, 1), written: map[string]string{}, versions: map[uint]string{},
	}
	if os.Getenv("ASC_PORTAL_KEEP") == "" {
		t.Cleanup(func() { _ = ops.Remove() })
	}
	m.tick()
	if st := m.Status(); st.State != models.PortalRunning {
		t.Fatalf("玩家站没有运行：%+v", st)
	}

	get := func(path string) (int, string) {
		deadline := time.Now().Add(30 * time.Second)
		for {
			resp, err := http.Get(fmt.Sprintf("http://127.0.0.1:%d%s", port, path))
			if err == nil {
				body, _ := io.ReadAll(resp.Body)
				resp.Body.Close()
				// 玩家站每秒检查一次快照，Docker Desktop 的文件共享看到新文件可能晚一点，503、404 都再等等
				if (resp.StatusCode != http.StatusServiceUnavailable && resp.StatusCode != http.StatusNotFound) || time.Now().After(deadline) {
					return resp.StatusCode, string(body)
				}
			} else if time.Now().After(deadline) {
				t.Fatalf("请求 %s 失败：%v", path, err)
			}
			time.Sleep(500 * time.Millisecond)
		}
	}
	if code, body := get("/api/v1/portal"); code != 200 || !strings.Contains(body, "冒烟测试") {
		t.Fatalf("portal：%d %s", code, body)
	}
	if code, body := get("/api/v1/shops/1"); code != 200 || !strings.Contains(body, "/购买") || strings.Contains(body, "s3cret") {
		t.Fatalf("shop：%d %s", code, body)
	}
	if code, body := get("/"); code != 200 || !strings.Contains(body, "<div id=\"root\">") {
		t.Fatalf("网页：%d %.200s", code, body)
	}
}
