package docker_manager

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"ark-server-commander/config"
	"ark-server-commander/utils"

	"github.com/containerd/errdefs"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"go.uber.org/zap"
)

// P2P 联机房的 EasyTier 节点容器（ase-link-<id>）。用宿主机网络，入站到虚拟 IP 的游戏流量被投递到宿主机的
// 127.0.0.1:<端口>，再经 Docker 的端口映射到达 ARK 容器。凭证文件放在数据目录 servers/<id>/link/，挂到容器 /data

// LinkContainerSpec 联机房节点容器的配置
type LinkContainerSpec struct {
	ServerID uint
	Image    string
	Args     []string          // easytier-core 的参数（镜像的入口已经是 easytier-core）
	Labels   map[string]string // 记录参数摘要，参数变化时重建
}

// LinkContainerInfo 联机房节点容器的状态
type LinkContainerInfo struct {
	Exists  bool
	Running bool
	ImageID string
	Labels  map[string]string
}

// linkDirRel 联机房数据在数据目录下的相对路径
func linkDirRel(serverID uint) string {
	return fmt.Sprintf("servers/%d/link", serverID)
}

// LinkDataDir 联机房数据（凭证文件）在后端本地的目录
func LinkDataDir(serverID uint) string {
	return filepath.Join(config.DataDir, filepath.FromSlash(linkDirRel(serverID)))
}

// RemoveLinkDir 删除联机房数据目录（凭证文件）
func RemoveLinkDir(serverID uint) error {
	if serverID == 0 {
		return errors.New("无效的服务器 ID")
	}
	return os.RemoveAll(LinkDataDir(serverID))
}

// InspectLinkContainer 读取联机房节点容器的状态，不存在时 Exists 为 false
func (dm *DockerManager) InspectLinkContainer(name string) (LinkContainerInfo, error) {
	info, err := dm.client.ContainerInspect(dm.ctx, name)
	if err != nil {
		if errdefs.IsNotFound(err) {
			return LinkContainerInfo{}, nil
		}
		return LinkContainerInfo{}, fmt.Errorf("获取Docker容器信息失败: %v", err)
	}
	out := LinkContainerInfo{Exists: true, ImageID: info.Image}
	if info.State != nil {
		out.Running = info.State.Running
	}
	if info.Config != nil {
		out.Labels = info.Config.Labels
	}
	return out, nil
}

// CreateLinkContainer 按配置（重新）创建并启动联机房节点容器：已存在的同名容器先删除。
// 数据目录先在本地建好，否则 Docker 会以 root 自动创建
func (dm *DockerManager) CreateLinkContainer(spec LinkContainerSpec) error {
	name := utils.GetLinkContainerName(spec.ServerID)
	if err := os.MkdirAll(LinkDataDir(spec.ServerID), 0o755); err != nil {
		return fmt.Errorf("创建联机房数据目录失败: %v", err)
	}
	if _, err := dm.client.ContainerInspect(dm.ctx, name); err == nil {
		if err := dm.client.ContainerRemove(dm.ctx, name, container.RemoveOptions{Force: true}); err != nil {
			return fmt.Errorf("删除旧的联机房容器失败: %v", err)
		}
	} else if !errdefs.IsNotFound(err) {
		return fmt.Errorf("获取Docker容器信息失败: %v", err)
	}

	cfg := &container.Config{
		Image:  spec.Image,
		Cmd:    spec.Args,
		Env:    []string{"TZ=" + config.ContainerTZ},
		Labels: spec.Labels,
	}
	hostCfg := &container.HostConfig{
		NetworkMode:   "host",
		RestartPolicy: container.RestartPolicy{Name: container.RestartPolicyUnlessStopped},
		Binds:         []string{path.Join(hostDataDir, linkDirRel(spec.ServerID)) + ":/data"},
	}
	resp, err := dm.client.ContainerCreate(dm.ctx, cfg, hostCfg, nil, nil, name)
	if err != nil {
		return fmt.Errorf("创建联机房容器失败: %v", err)
	}
	if err := dm.client.ContainerStart(dm.ctx, resp.ID, container.StartOptions{}); err != nil {
		_ = dm.client.ContainerRemove(dm.ctx, resp.ID, container.RemoveOptions{Force: true})
		return fmt.Errorf("启动联机房容器失败: %v", err)
	}
	utils.Info("联机房容器已启动", zap.String("container_name", name))
	return nil
}

// RemoveLinkContainer 强制删除联机房节点容器，不存在时不算错误。
// 节点不保存需要落盘的状态（凭证在 /data），不用先优雅停止
func (dm *DockerManager) RemoveLinkContainer(name string) error {
	err := dm.client.ContainerRemove(dm.ctx, name, container.RemoveOptions{Force: true})
	if err != nil && !errdefs.IsNotFound(err) {
		return fmt.Errorf("删除联机房容器失败: %v", err)
	}
	if err == nil {
		utils.Info("联机房容器已删除", zap.String("container_name", name))
	}
	return nil
}

// ListLinkContainers 列出所有联机房节点容器（ase-link-<id>，含已停止的）的名称
func (dm *DockerManager) ListLinkContainers() ([]string, error) {
	list, err := dm.client.ContainerList(dm.ctx, container.ListOptions{
		All:     true,
		Filters: filters.NewArgs(filters.Arg("name", "^/ase-link-[0-9]+$")),
	})
	if err != nil {
		return nil, fmt.Errorf("获取容器列表失败: %v", err)
	}
	var names []string
	for _, c := range list {
		for _, n := range c.Names {
			names = append(names, strings.TrimPrefix(n, "/"))
		}
	}
	return names, nil
}
