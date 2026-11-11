package docker_manager

import (
	"ark-server-commander/utils"
	"errors"
	"fmt"
	"os"

	"github.com/containerd/errdefs"
	"github.com/docker/docker/api/types/volume"
	"go.uber.org/zap"
)

// serverVolumeNames 返回服务器的全部卷：服务端文件卷、游戏数据卷、插件卷
func serverVolumeNames(serverID uint) []string {
	specs := serverVolumeSpecs(serverID)
	names := make([]string, 0, len(specs))
	for _, s := range specs {
		names = append(names, s.Name)
	}
	return names
}

// CreateVolume 创建服务器的 Docker 卷（服务端文件卷、游戏数据卷和插件卷）。
// 每个卷都绑定到数据目录下的 servers/<id>/ 子目录，宿主机上可以直接访问。
// 任一卷创建失败时会清理已创建的卷
func (dm *DockerManager) CreateVolume(serverID uint) error {
	for _, spec := range serverVolumeSpecs(serverID) {
		if err := dm.ensureBindVolume(spec); err != nil {
			if rmErr := dm.RemoveVolume(serverID); rmErr != nil {
				utils.Warn("清理已创建的卷失败", zap.Uint("server_id", serverID), zap.Error(rmErr))
			}
			return fmt.Errorf("创建卷 %s 失败: %v", spec.Name, err)
		}
	}

	utils.Info("Docker卷创建成功", zap.Strings("volumes", serverVolumeNames(serverID)))
	return nil
}

// EnsureClusterVolume 确保集群共享卷存在并绑定到数据目录下的 clusters/<cluster_id>。
// 必须在创建容器前调用，否则 Docker 会自动创建一个普通卷
func (dm *DockerManager) EnsureClusterVolume(clusterID string) error {
	if ExpectedClusterVolume(clusterID) == "" {
		return nil
	}
	return dm.ensureBindVolume(clusterVolumeSpec(clusterID))
}

// ensureBindVolume 创建数据目录下的目录，并创建指向该目录的命名卷；卷已存在时只确保目录存在
func (dm *DockerManager) ensureBindVolume(spec volumeSpec) error {
	if hostDataDir == "" {
		return errors.New("数据目录尚未初始化")
	}
	if err := os.MkdirAll(spec.localDir(), 0o755); err != nil {
		return fmt.Errorf("创建目录失败: %w", err)
	}

	exists, err := dm.VolumeExists(spec.Name)
	if err != nil {
		return fmt.Errorf("检查卷是否存在失败: %v", err)
	}
	if exists {
		utils.Debug("Docker卷已存在，跳过创建", zap.String("volume", spec.Name))
		return nil
	}

	utils.Infof("正在创建Docker卷: %s -> %s", spec.Name, spec.hostDir())
	_, err = dm.client.VolumeCreate(dm.ctx, volume.CreateOptions{
		Name:   spec.Name,
		Driver: "local",
		DriverOpts: map[string]string{
			"type":   "none",
			"o":      "bind",
			"device": spec.hostDir(),
		},
	})
	if err != nil {
		return fmt.Errorf("创建Docker卷失败: %v", err)
	}

	utils.Info("Docker卷创建成功", zap.String("volume", spec.Name))
	return nil
}

// RemoveVolume 删除服务器的 Docker 卷（服务端文件卷、游戏数据卷和插件卷）。
// 只删除卷本身，绑定目录里的文件保留，需要时再调用 RemoveServerDirs。
// 每个卷都会尝试删除，任一失败都会返回错误
func (dm *DockerManager) RemoveVolume(serverID uint) error {
	var errs []error
	for _, name := range serverVolumeNames(serverID) {
		if err := dm.removeSingleVolume(name); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

// removeSingleVolume 删除单个Docker卷
func (dm *DockerManager) removeSingleVolume(volumeName string) error {
	exists, err := dm.VolumeExists(volumeName)
	if err != nil {
		return fmt.Errorf("检查卷是否存在失败: %v", err)
	}
	if !exists {
		utils.Debug("Docker卷不存在，跳过删除", zap.String("volume", volumeName))
		return nil
	}

	utils.Infof("正在删除Docker卷: %s", volumeName)
	if err := dm.client.VolumeRemove(dm.ctx, volumeName, false); err != nil {
		return fmt.Errorf("删除Docker卷失败: %v", err)
	}

	utils.Info("Docker卷删除成功", zap.String("volume", volumeName))
	return nil
}

// VolumeExists 检查Docker卷是否存在
func (dm *DockerManager) VolumeExists(volumeName string) (bool, error) {
	_, err := dm.client.VolumeInspect(dm.ctx, volumeName)
	if err != nil {
		if errdefs.IsNotFound(err) {
			return false, nil
		}
		return false, fmt.Errorf("检查Docker卷失败: %v", err)
	}

	return true, nil
}
