package server

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"ark-server-commander/config"
	"ark-server-commander/models"
	"ark-server-commander/service/docker_manager"
	"ark-server-commander/utils"
)

// TestLinkNodeDockerSmoke 用真实的 Docker 和 easytier/easytier 镜像启动一个联机房节点（参数由 linkContainerSpec 生成，
// peers 指向不存在的地址），在节点里签发和吊销凭证，最后删除容器。需要本机 Docker，默认跳过：
//
//	ASC_LINK_DOCKER_TEST=1 go test ./service/server/ -run TestLinkNodeDockerSmoke -v
func TestLinkNodeDockerSmoke(t *testing.T) {
	if os.Getenv("ASC_LINK_DOCKER_TEST") != "1" {
		t.Skip("设置 ASC_LINK_DOCKER_TEST=1 才运行（需要 Docker）")
	}
	utils.InitLogger()
	oldDir, oldHost := config.DataDir, config.HostDataDir
	config.DataDir, config.HostDataDir = t.TempDir(), ""
	t.Cleanup(func() { config.DataDir, config.HostDataDir = oldDir, oldHost })
	if err := docker_manager.InitStorage(); err != nil {
		t.Fatal(err)
	}

	ops := dockerLinkOps{}
	deadline := time.Now().Add(5 * time.Minute)
	for {
		ready, err := ops.ImageReady(utils.EasyTierImage)
		if err != nil {
			t.Fatal(err)
		}
		if ready {
			break
		}
		if time.Now().After(deadline) {
			t.Fatal("下载镜像超时")
		}
		time.Sleep(2 * time.Second)
	}

	server := models.Server{
		ID: 59001, SessionName: "冒烟 测试", Port: 7777, QueryPort: 27015,
		LinkNetworkName: "ark-smoke0000000000", LinkNetworkSecret: strings.Repeat("ab", 32), LinkRPCPort: 39001,
	}
	name := utils.GetLinkContainerName(server.ID)
	t.Cleanup(func() { _ = ops.Remove(name) })
	spec := linkContainerSpec(server, []string{"tcp://127.0.0.1:1"})
	if err := ops.Create(spec); err != nil {
		t.Fatal(err)
	}

	m := &LinkManager{docker: ops, now: time.Now, newID: newUUID}
	var cred linkCredential
	var err error
	for range 20 { // 节点启动后 RPC 端口要等一会儿才可用
		time.Sleep(time.Second)
		if cred, err = m.generateCredential(name, server.LinkRPCPort); err == nil {
			break
		}
	}
	if err != nil {
		t.Fatal(err)
	}
	if cred.Secret == "" || len(cred.ID) != 36 {
		t.Fatalf("credential = %+v", cred)
	}
	info, err := ops.Inspect(name)
	if err != nil || !info.Running || info.Labels[linkCommandHashLabel] != spec.Labels[linkCommandHashLabel] {
		t.Fatalf("info = %+v err = %v", info, err)
	}
	if err := m.revokeCredential(name, server.LinkRPCPort, cred.ID); err != nil {
		t.Fatal(err)
	}
	// 不存在的凭证也算吊销完成
	if err := m.revokeCredential(name, server.LinkRPCPort, newUUID()); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(docker_manager.LinkDataDir(server.ID), "cred.json")); err != nil {
		t.Logf("凭证文件还没写入数据目录: %v", err)
	}
	if err := ops.Remove(name); err != nil {
		t.Fatal(err)
	}
	if info, _ := ops.Inspect(name); info.Exists {
		t.Fatal("container should be removed")
	}
}
