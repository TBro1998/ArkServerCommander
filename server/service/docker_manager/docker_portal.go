package docker_manager

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"ark-server-commander/config"
	"ark-server-commander/utils"

	"github.com/containerd/errdefs"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	"go.uber.org/zap"
)

// 玩家站容器（ark-player-portal，每个管理端一个）：只读挂载数据目录 portal/，提供网页和只读接口。
// 它对公网开放，所以和管理端隔离：不提权、去掉所有 capability、根文件系统只读、不挂 Docker socket、
// 只看得到快照目录。用 bridge 网络和端口映射，联机房节点照样能经 127.0.0.1:<端口> 访问到它

// portalDirRel 玩家站快照在数据目录下的相对路径
const portalDirRel = "portal"

// portalContainerPort 玩家站在容器里监听的端口（镜像默认的 PORTAL_LISTEN）
const portalContainerPort = "8080/tcp"

// PortalDataDir 玩家站快照在后端本地的目录
func PortalDataDir() string {
	return filepath.Join(config.DataDir, portalDirRel)
}

// PortalContainerSpec 玩家站容器的配置
type PortalContainerSpec struct {
	Image    string
	HostPort int
	Env      []string
	Labels   map[string]string
}

// CreatePortalContainer 按配置（重新）创建并启动玩家站容器：已存在的同名容器先删除。
// 快照目录先在本地建好，否则 Docker 会以 root 自动创建
func (dm *DockerManager) CreatePortalContainer(spec PortalContainerSpec) error {
	name := utils.PortalContainerName
	if err := os.MkdirAll(PortalDataDir(), 0o755); err != nil {
		return fmt.Errorf("创建玩家站数据目录失败: %v", err)
	}
	if err := dm.RemovePortalContainer(); err != nil {
		return err
	}

	port := nat.Port(portalContainerPort)
	cfg := &container.Config{
		Image:        spec.Image,
		Env:          append([]string{"TZ=" + config.ContainerTZ}, spec.Env...),
		Labels:       spec.Labels,
		ExposedPorts: nat.PortSet{port: struct{}{}},
	}
	hostCfg := &container.HostConfig{
		RestartPolicy:  container.RestartPolicy{Name: container.RestartPolicyUnlessStopped},
		Binds:          []string{hostPath(PortalDataDir()) + ":/data:ro"},
		PortBindings:   nat.PortMap{port: []nat.PortBinding{{HostIP: "0.0.0.0", HostPort: strconv.Itoa(spec.HostPort)}}},
		ReadonlyRootfs: true,
		CapDrop:        []string{"ALL"},
		SecurityOpt:    []string{"no-new-privileges"},
		Resources: container.Resources{
			Memory:    256 << 20,
			PidsLimit: ptr(int64(256)),
		},
	}
	resp, err := dm.client.ContainerCreate(dm.ctx, cfg, hostCfg, nil, nil, name)
	if err != nil {
		return fmt.Errorf("创建玩家站容器失败: %v", err)
	}
	if err := dm.client.ContainerStart(dm.ctx, resp.ID, container.StartOptions{}); err != nil {
		_ = dm.client.ContainerRemove(dm.ctx, resp.ID, container.RemoveOptions{Force: true})
		return fmt.Errorf("启动玩家站容器失败: %v", err)
	}
	utils.Info("玩家站容器已启动", zap.String("container_name", name), zap.Int("port", spec.HostPort))
	return nil
}

// RemovePortalContainer 强制删除玩家站容器，不存在时不算错误（玩家站不保存任何状态）
func (dm *DockerManager) RemovePortalContainer() error {
	err := dm.client.ContainerRemove(dm.ctx, utils.PortalContainerName, container.RemoveOptions{Force: true})
	if err != nil && !errdefs.IsNotFound(err) {
		return fmt.Errorf("删除玩家站容器失败: %v", err)
	}
	if err == nil {
		utils.Info("玩家站容器已删除", zap.String("container_name", utils.PortalContainerName))
	}
	return nil
}

func ptr[T any](v T) *T { return &v }
