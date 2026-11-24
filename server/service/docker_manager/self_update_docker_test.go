package docker_manager

import (
	"os"
	"slices"
	"strings"
	"testing"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/network"
)

// 用真实 Docker 验证替换容器的流程（改名、按旧配置创建、停旧启新、等端口响应、删旧容器，以及失败时回滚）。
// 默认跳过，ASC_SELF_UPDATE_DOCKER_TEST=1 时运行，需要本地有 busybox:1.36、busybox:1.37、alpine:3.20。
// 用 busybox 的 httpd 模拟管理端的 8080 / 8081，所有容器、网络、镜像标签都以 asc-su-test 开头，不碰别的容器。

const (
	suTestName    = "asc-su-test"
	suTestNetwork = "asc-su-test-net"
	suTestRef     = "asc-su-test-image:latest"
)

func setupSelfUpdateDockerTest(t *testing.T) *DockerManager {
	t.Helper()
	if os.Getenv("ASC_SELF_UPDATE_DOCKER_TEST") != "1" {
		t.Skip("设置 ASC_SELF_UPDATE_DOCKER_TEST=1 才运行")
	}
	dm, err := GetDockerManager()
	if err != nil {
		t.Fatal(err)
	}
	cleanup := func() {
		for _, name := range []string{suTestName, suTestName + "-old"} {
			if info, err := dm.client.ContainerInspect(dm.ctx, name); err == nil {
				_ = dm.client.ContainerRemove(dm.ctx, info.ID, container.RemoveOptions{Force: true, RemoveVolumes: true})
			}
		}
		_ = dm.client.NetworkRemove(dm.ctx, suTestNetwork)
		_, _ = dm.client.ImageRemove(dm.ctx, suTestRef, image.RemoveOptions{})
	}
	cleanup()
	t.Cleanup(cleanup)
	if _, err := dm.client.NetworkCreate(dm.ctx, suTestNetwork, network.CreateOptions{}); err != nil {
		t.Fatal(err)
	}
	return dm
}

// startFakeCommander 用 image 创建并启动模拟的管理端：主网络 bridge，另连一个自定义网络，镜像声明式的匿名卷 /server-files
func startFakeCommander(t *testing.T, dm *DockerManager, image string) container.InspectResponse {
	t.Helper()
	if err := dm.client.ImageTag(dm.ctx, image, suTestRef); err != nil {
		t.Fatal(err)
	}
	resp, err := dm.client.ContainerCreate(dm.ctx,
		&container.Config{
			Image:   suTestRef,
			Env:     []string{"TZ=Asia/Shanghai"},
			Cmd:     []string{"sh", "-c", "mkdir -p /www && httpd -p 8081 -h /www && exec httpd -f -p 8080 -h /www"},
			Labels:  map[string]string{"com.docker.compose.project": "asc-su-test", "com.docker.compose.image": "sha256:old"},
			Volumes: map[string]struct{}{"/server-files": {}},
		},
		&container.HostConfig{RestartPolicy: container.RestartPolicy{Name: container.RestartPolicyUnlessStopped}},
		nil, nil, suTestName)
	if err != nil {
		t.Fatal(err)
	}
	if err := dm.client.NetworkConnect(dm.ctx, suTestNetwork, resp.ID, &network.EndpointSettings{Aliases: []string{"commander"}}); err != nil {
		t.Fatal(err)
	}
	if err := dm.client.ContainerStart(dm.ctx, resp.ID, container.StartOptions{}); err != nil {
		t.Fatal(err)
	}
	if err := dm.waitCommanderReady(resp.ID, selfUpdateHealthTimeout); err != nil {
		t.Fatalf("模拟的管理端没有起来: %v", err)
	}
	info, err := dm.client.ContainerInspect(dm.ctx, resp.ID)
	if err != nil {
		t.Fatal(err)
	}
	return info
}

func anonymousVolume(t *testing.T, info container.InspectResponse) string {
	t.Helper()
	m, ok := findMount(info.Mounts, "/server-files")
	if !ok || m.Name == "" {
		t.Fatalf("没有 /server-files 卷: %+v", info.Mounts)
	}
	return m.Name
}

func TestSelfUpdateDockerReplace(t *testing.T) {
	dm := setupSelfUpdateDockerTest(t)
	old := startFakeCommander(t, dm, "busybox:1.36")
	vol := anonymousVolume(t, old)

	if err := dm.client.ImageTag(dm.ctx, "busybox:1.37", suTestRef); err != nil {
		t.Fatal(err)
	}
	newID, err := dm.GetImageID(suTestRef)
	if err != nil || newID == old.Image {
		t.Fatalf("新镜像应和旧镜像不同: %v", err)
	}

	state, err := dm.replaceContainer(old.ID, newID)
	if err != nil || state != SelfUpdateSucceeded {
		t.Fatalf("替换失败: %s %v", state, err)
	}

	cur, err := dm.client.ContainerInspect(dm.ctx, suTestName)
	if err != nil {
		t.Fatal(err)
	}
	if cur.ID == old.ID || cur.Image != newID || !cur.State.Running {
		t.Errorf("新容器不对: id=%s image=%s running=%v", cur.ID, cur.Image, cur.State.Running)
	}
	if got := anonymousVolume(t, cur); got != vol {
		t.Errorf("/server-files 应接回原来的卷 %s，得到 %s", vol, got)
	}
	if cur.Config.Hostname == old.Config.Hostname || !strings.HasPrefix(cur.ID, cur.Config.Hostname) {
		t.Errorf("hostname 应是新容器自己的 ID，得到 %s", cur.Config.Hostname)
	}
	if cur.Config.Labels["com.docker.compose.image"] != newID || cur.Config.Labels["com.docker.compose.project"] != "asc-su-test" {
		t.Errorf("标签不对: %v", cur.Config.Labels)
	}
	if cur.HostConfig.RestartPolicy.Name != container.RestartPolicyUnlessStopped {
		t.Errorf("重启策略应保留: %v", cur.HostConfig.RestartPolicy)
	}
	ep := cur.NetworkSettings.Networks[suTestNetwork]
	if ep == nil || !slices.Contains(ep.Aliases, "commander") {
		t.Errorf("应连上自定义网络并保留别名: %+v", cur.NetworkSettings.Networks)
	}
	if _, err := dm.client.ContainerInspect(dm.ctx, old.ID); err == nil {
		t.Error("旧容器应已删除")
	}
}

func TestSelfUpdateDockerRollback(t *testing.T) {
	dm := setupSelfUpdateDockerTest(t)
	old := startFakeCommander(t, dm, "busybox:1.36")

	// alpine 的 busybox 没有 httpd，新容器起不来
	if err := dm.client.ImageTag(dm.ctx, "alpine:3.20", suTestRef); err != nil {
		t.Fatal(err)
	}
	newID, err := dm.GetImageID(suTestRef)
	if err != nil {
		t.Fatal(err)
	}

	state, err := dm.replaceContainer(old.ID, newID)
	if err == nil || state != SelfUpdateRolledBack {
		t.Fatalf("应回滚，得到 %s %v", state, err)
	}
	t.Logf("回滚原因: %v", err)

	cur, err := dm.client.ContainerInspect(dm.ctx, suTestName)
	if err != nil {
		t.Fatal(err)
	}
	if cur.ID != old.ID || !cur.State.Running {
		t.Errorf("旧容器应改回原名并重新运行: id=%s running=%v", cur.ID, cur.State.Running)
	}
	if _, err := dm.client.ContainerInspect(dm.ctx, suTestName+"-old"); err == nil {
		t.Error("不应留下 -old 容器")
	}
	if err := dm.waitCommanderReady(cur.ID, selfUpdateHealthTimeout); err != nil {
		t.Errorf("回滚后的旧容器应正常响应: %v", err)
	}
}
