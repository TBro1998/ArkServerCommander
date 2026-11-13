package docker_manager

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"ark-server-commander/config"
	"ark-server-commander/utils"

	"github.com/containerd/errdefs"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/volume"
)

// ErrVolumeInUse 卷仍被容器挂载（例如改了集群 ID 但还没重建容器的服务器）
var ErrVolumeInUse = errors.New("卷仍被容器使用")

const clusterVolumePrefix = "ase-cluster-"

// ClusterIDs 列出存在集群卷或集群目录（数据目录 clusters/ 下）的集群 ID，按名称排序
func (dm *DockerManager) ClusterIDs() ([]string, error) {
	ids := map[string]bool{}

	// name 过滤是子串匹配，再按前缀确认
	resp, err := dm.client.VolumeList(dm.ctx, volume.ListOptions{Filters: filters.NewArgs(filters.Arg("name", clusterVolumePrefix))})
	if err != nil {
		return nil, fmt.Errorf("列出集群卷失败: %w", err)
	}
	for _, v := range resp.Volumes {
		if id, ok := strings.CutPrefix(v.Name, clusterVolumePrefix); ok && utils.ClusterIDPattern.MatchString(id) {
			ids[id] = true
		}
	}

	entries, err := os.ReadDir(filepath.Join(config.DataDir, "clusters"))
	if err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("读取集群数据目录失败: %w", err)
	}
	for _, e := range entries {
		if e.IsDir() && utils.ClusterIDPattern.MatchString(e.Name()) {
			ids[e.Name()] = true
		}
	}

	result := make([]string, 0, len(ids))
	for id := range ids {
		result = append(result, id)
	}
	sort.Strings(result)
	return result, nil
}

// RemoveCluster 删除集群卷和数据目录 clusters/<id> 里的全部数据（跨服传送的角色和物品）。
// 卷仍被容器挂载时返回 ErrVolumeInUse，目录不动
func (dm *DockerManager) RemoveCluster(clusterID string) error {
	if !utils.ClusterIDPattern.MatchString(clusterID) {
		return fmt.Errorf("无效的集群ID: %s", clusterID)
	}
	name := utils.GetClusterVolumeName(clusterID)
	if err := dm.client.VolumeRemove(dm.ctx, name, false); err != nil && !errdefs.IsNotFound(err) {
		if errdefs.IsConflict(err) {
			return fmt.Errorf("%w: %v", ErrVolumeInUse, err)
		}
		return fmt.Errorf("删除集群卷失败: %w", err)
	}
	if err := os.RemoveAll(clusterVolumeSpec(clusterID).localDir()); err != nil {
		return fmt.Errorf("删除集群数据目录失败: %w", err)
	}
	return nil
}
