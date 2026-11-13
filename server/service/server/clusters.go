package server

import (
	"errors"
	"fmt"

	"ark-server-commander/database"
	"ark-server-commander/models"
	"ark-server-commander/service/docker_manager"
	"ark-server-commander/utils"
)

// UnusedCluster 没有服务器使用的集群数据
type UnusedCluster struct {
	ID   string `json:"id"`
	Size int64  `json:"size"` // 数据目录 clusters/<id> 的大小，来自磁盘统计缓存
}

// unusedClusters 存在卷或目录、但没有任何服务器（包括其他用户的）设置该集群 ID 的集群
func unusedClusters(sizes map[string]int64) ([]UnusedCluster, error) {
	dm, err := docker_manager.GetDockerManager()
	if err != nil {
		return nil, err
	}
	ids, err := dm.ClusterIDs()
	if err != nil {
		return nil, err
	}
	var used []string
	if err := database.DB.Model(&models.Server{}).Where("cluster_id <> ''").Distinct().Pluck("cluster_id", &used).Error; err != nil {
		return nil, fmt.Errorf("读取服务器集群配置失败: %w", err)
	}
	inUse := make(map[string]bool, len(used))
	for _, id := range used {
		inUse[id] = true
	}
	result := []UnusedCluster{}
	for _, id := range ids {
		if !inUse[id] {
			result = append(result, UnusedCluster{ID: id, Size: sizes[utils.GetClusterVolumeName(id)]})
		}
	}
	return result, nil
}

// DeleteCluster 删除不再使用的集群卷和数据目录
func (s *ServerService) DeleteCluster(clusterID string) error {
	if !utils.ClusterIDPattern.MatchString(clusterID) {
		return ErrInvalidClusterID
	}
	var count int64
	if err := database.DB.Model(&models.Server{}).Where("cluster_id = ?", clusterID).Count(&count).Error; err != nil {
		return fmt.Errorf("读取服务器集群配置失败: %w", err)
	}
	if count > 0 {
		return ErrClusterInUse
	}
	dm, err := docker_manager.GetDockerManager()
	if err != nil {
		return err
	}
	if err := dm.RemoveCluster(clusterID); err != nil {
		if errors.Is(err, docker_manager.ErrVolumeInUse) {
			return fmt.Errorf("%w：有容器仍挂载着它，先启动或重建改过集群 ID 的服务器", ErrClusterInUse)
		}
		return err
	}
	invalidateVolumeSizes()
	return nil
}
