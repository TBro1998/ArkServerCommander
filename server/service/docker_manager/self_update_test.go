package docker_manager

import (
	"slices"
	"testing"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/go-connections/nat"
	dockerspec "github.com/moby/docker-image-spec/specs-go/v1"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
)

func TestIsCommanderLatest(t *testing.T) {
	cases := map[string]bool{
		"tbro98/arkservercommander":                                 true,
		"tbro98/arkservercommander:latest":                          true,
		"docker.io/tbro98/arkservercommander:latest":                true,
		"docker.m.daocloud.io/tbro98/arkservercommander:latest":     true, // 镜像加速器
		"tbro98/arkservercommander:0.0.5":                           false,
		"tbro98/arkservercommander@sha256:" + sixtyFourHex():        false,
		"arkservercommander:dev":                                    false,
		"someone/arkservercommander:latest":                         false,
		"sha256:" + sixtyFourHex():                                  false,
		"":                                                          false,
		"tbro98/arkservercommander:latest@sha256:" + sixtyFourHex(): false,
	}
	for image, want := range cases {
		if got := isCommanderLatest(image); got != want {
			t.Errorf("isCommanderLatest(%q) = %v, want %v", image, got, want)
		}
	}
}

func sixtyFourHex() string {
	return "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"
}

const testContainerID = "abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789"

// composeContainer 模拟 docker compose 创建的管理端容器的 inspect 结果
func composeContainer() container.InspectResponse {
	return container.InspectResponse{
		ContainerJSONBase: &container.ContainerJSONBase{
			ID:    testContainerID,
			Name:  "/ark-commander",
			Image: "sha256:old",
			HostConfig: &container.HostConfig{
				NetworkMode:   "ark_default",
				Privileged:    true,
				RestartPolicy: container.RestartPolicy{Name: container.RestartPolicyUnlessStopped},
				PortBindings:  nat.PortMap{"8080/tcp": {{HostPort: "8080"}}},
				Binds:         []string{"/opt/ark/data:/data:rw", "/var/run/docker.sock:/var/run/docker.sock:rw"},
			},
		},
		Config: &container.Config{
			Hostname:     testContainerID[:12],
			Image:        "tbro98/arkservercommander:latest",
			Env:          []string{"TZ=Asia/Shanghai", "PATH=/usr/bin", "NODE_VERSION=26.0.0", "NODE_ENV=production"},
			Cmd:          []string{"./docker-entrypoint.sh"},
			Entrypoint:   []string{"docker-entrypoint.sh"},
			WorkingDir:   "/app",
			ExposedPorts: nat.PortSet{"8080/tcp": {}},
			Volumes:      map[string]struct{}{"/server-files": {}},
			Labels: map[string]string{
				"com.docker.compose.project": "ark",
				"com.docker.compose.service": "ark-commander",
				"com.docker.compose.image":   "sha256:old",
				"org.opencontainers.version": "1",
			},
		},
		Mounts: []container.MountPoint{
			{Type: mount.TypeBind, Source: "/opt/ark/data", Destination: "/data", RW: true},
			{Type: mount.TypeBind, Source: "/var/run/docker.sock", Destination: "/var/run/docker.sock", RW: true},
			{Type: mount.TypeVolume, Name: "f00dcafe", Source: "/var/lib/docker/volumes/f00dcafe/_data", Destination: "/server-files", RW: true},
		},
		NetworkSettings: &container.NetworkSettings{
			Networks: map[string]*network.EndpointSettings{
				"ark_default": {
					Aliases:   []string{"ark-commander", testContainerID[:12]},
					IPAddress: "172.18.0.2",
					NetworkID: "net1",
				},
				"monitoring": {Aliases: []string{"commander"}, IPAddress: "172.19.0.3"},
			},
		},
	}
}

func oldImageConfig() *dockerspec.DockerOCIImageConfig {
	return &dockerspec.DockerOCIImageConfig{ImageConfig: ocispec.ImageConfig{
		Env:          []string{"PATH=/usr/bin", "NODE_VERSION=26.0.0", "NODE_ENV=production"},
		Cmd:          []string{"./docker-entrypoint.sh"},
		Entrypoint:   []string{"docker-entrypoint.sh"},
		WorkingDir:   "/app",
		ExposedPorts: map[string]struct{}{"8080/tcp": {}},
		Volumes:      map[string]struct{}{"/server-files": {}},
		Labels:       map[string]string{"org.opencontainers.version": "1"},
	}}
}

func TestBuildReplacementDropsImageDefaults(t *testing.T) {
	old := composeContainer()
	cfg, hostCfg, primary, extra := buildReplacement(old, oldImageConfig(), "sha256:new")

	if cfg.Image != "tbro98/arkservercommander:latest" {
		t.Errorf("Image = %q，应保留镜像名", cfg.Image)
	}
	if cfg.Hostname != "" {
		t.Errorf("Hostname = %q，自动生成的 hostname 不能带到新容器", cfg.Hostname)
	}
	if !slices.Equal(cfg.Env, []string{"TZ=Asia/Shanghai"}) {
		t.Errorf("Env = %v，只应保留用户设置的变量", cfg.Env)
	}
	if cfg.Cmd != nil || cfg.Entrypoint != nil || cfg.WorkingDir != "" {
		t.Errorf("Cmd/Entrypoint/WorkingDir 和旧镜像相同，应交给新镜像: %v %v %q", cfg.Cmd, cfg.Entrypoint, cfg.WorkingDir)
	}
	if len(cfg.ExposedPorts) != 0 || len(cfg.Volumes) != 0 {
		t.Errorf("ExposedPorts/Volumes 来自镜像，应去掉: %v %v", cfg.ExposedPorts, cfg.Volumes)
	}
	if _, ok := cfg.Labels["org.opencontainers.version"]; ok {
		t.Error("镜像的标签应去掉")
	}
	if cfg.Labels["com.docker.compose.project"] != "ark" || cfg.Labels["com.docker.compose.image"] != "sha256:new" {
		t.Errorf("compose 标签不对: %v", cfg.Labels)
	}
	// 不能改到旧容器的配置
	if old.Config.Labels["com.docker.compose.image"] != "sha256:old" || len(old.Config.Env) != 4 || old.Config.Hostname == "" {
		t.Error("buildReplacement 改动了旧容器的配置")
	}

	if !hostCfg.Privileged || hostCfg.RestartPolicy.Name != container.RestartPolicyUnlessStopped || len(hostCfg.PortBindings) != 1 {
		t.Errorf("HostConfig 应原样保留: %+v", hostCfg)
	}
	if len(hostCfg.Binds) != 2 {
		t.Errorf("Binds = %v", hostCfg.Binds)
	}
	if len(hostCfg.Mounts) != 1 || hostCfg.Mounts[0].Source != "f00dcafe" || hostCfg.Mounts[0].Target != "/server-files" ||
		hostCfg.Mounts[0].Type != mount.TypeVolume || hostCfg.Mounts[0].ReadOnly {
		t.Errorf("匿名卷应按卷名接回: %+v", hostCfg.Mounts)
	}
	if len(old.HostConfig.Mounts) != 0 {
		t.Error("buildReplacement 改动了旧容器的 HostConfig")
	}

	if primary == nil || primary.EndpointsConfig["ark_default"] == nil {
		t.Fatalf("主网络应随容器创建: %+v", primary)
	}
	ep := primary.EndpointsConfig["ark_default"]
	if !slices.Equal(ep.Aliases, []string{"ark-commander"}) || ep.IPAddress != "" || ep.NetworkID != "" {
		t.Errorf("主网络的设置不对（别名去掉容器 ID，不带运行时字段）: %+v", ep)
	}
	if len(extra) != 1 || extra["monitoring"] == nil || !slices.Equal(extra["monitoring"].Aliases, []string{"commander"}) {
		t.Errorf("其他网络应在启动前连接: %+v", extra)
	}
}

func TestBuildReplacementKeepsUserOverrides(t *testing.T) {
	old := composeContainer()
	old.Config.Hostname = "commander"
	old.Config.Cmd = []string{"sh", "-c", "custom"}
	old.Config.Env = append(old.Config.Env, "NODE_ENV=development")
	cfg, _, _, _ := buildReplacement(old, oldImageConfig(), "sha256:new")
	if cfg.Hostname != "commander" {
		t.Errorf("自己设置的 hostname 应保留，得到 %q", cfg.Hostname)
	}
	if !slices.Equal(cfg.Cmd, []string{"sh", "-c", "custom"}) {
		t.Errorf("改过的 Cmd 应保留，得到 %v", cfg.Cmd)
	}
	if !slices.Contains(cfg.Env, "NODE_ENV=development") {
		t.Errorf("改过的环境变量应保留，得到 %v", cfg.Env)
	}

	// 读不到旧镜像时照抄容器的配置
	cfg, _, _, _ = buildReplacement(old, nil, "sha256:new")
	if len(cfg.Env) != len(old.Config.Env) || cfg.WorkingDir != "/app" {
		t.Errorf("没有旧镜像时应保留完整配置: %+v", cfg)
	}
}

func TestBuildReplacementNetworkModes(t *testing.T) {
	// docker run 默认的 bridge 网络
	old := composeContainer()
	old.HostConfig.NetworkMode = "default"
	old.NetworkSettings.Networks = map[string]*network.EndpointSettings{"bridge": {IPAddress: "172.17.0.2"}}
	_, _, primary, extra := buildReplacement(old, nil, "sha256:new")
	if primary != nil || len(extra) != 0 {
		t.Errorf("默认网络不用指定: %+v %+v", primary, extra)
	}

	old.HostConfig.NetworkMode = "host"
	old.NetworkSettings.Networks = map[string]*network.EndpointSettings{"host": {}}
	_, _, primary, extra = buildReplacement(old, nil, "sha256:new")
	if primary != nil || len(extra) != 0 {
		t.Errorf("host 网络由 NetworkMode 决定，不用指定: %+v %+v", primary, extra)
	}
}

func TestBuildReplacementCoveredMounts(t *testing.T) {
	old := composeContainer()
	// compose 新版本把匿名卷写进 HostConfig.Mounts，这时不能重复加
	old.HostConfig.Mounts = []mount.Mount{{Type: mount.TypeVolume, Source: "f00dcafe", Target: "/server-files"}}
	_, hostCfg, _, _ := buildReplacement(old, nil, "sha256:new")
	if len(hostCfg.Mounts) != 1 {
		t.Errorf("已有的挂载不应重复: %+v", hostCfg.Mounts)
	}

	old = composeContainer()
	old.HostConfig.Binds = append(old.HostConfig.Binds, "ark-files:/server-files")
	_, hostCfg, _, _ = buildReplacement(old, nil, "sha256:new")
	if len(hostCfg.Mounts) != 0 {
		t.Errorf("Binds 里已有的卷不应再加: %+v", hostCfg.Mounts)
	}
}

func TestSelfUpdateRecordRoundTrip(t *testing.T) {
	dir := t.TempDir()
	if _, err := ReadSelfUpdateRecord(dir); err == nil {
		t.Fatal("没有记录时应返回错误")
	}
	now := time.Now().UTC().Truncate(time.Second)
	rec := SelfUpdateRecord{State: SelfUpdateSucceeded, FromVersion: "0.0.5", OldImageID: "sha256:old", NewImageID: "sha256:new", StartedAt: now, FinishedAt: &now}
	if err := WriteSelfUpdateRecord(dir, rec); err != nil {
		t.Fatal(err)
	}
	got, err := ReadSelfUpdateRecord(dir)
	if err != nil {
		t.Fatal(err)
	}
	if got.State != rec.State || got.NewImageID != rec.NewImageID || got.FinishedAt == nil || !got.FinishedAt.Equal(now) {
		t.Errorf("读回的记录不一致: %+v", got)
	}
}
