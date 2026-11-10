package docker_manager

import (
	"ark-server-commander/utils"
	"errors"
	"fmt"

	"github.com/containerd/errdefs"
	"github.com/docker/docker/api/types/volume"
	"go.uber.org/zap"
)

// serverVolumeNames 返回服务器的全部卷：服务端文件卷、游戏数据卷、插件卷
func serverVolumeNames(serverID uint) []string {
	return []string{
		utils.GetServerFilesVolumeName(serverID),
		utils.GetServerVolumeName(serverID),
		utils.GetServerPluginsVolumeName(serverID),
	}
}

// CreateVolume 创建Docker卷（包括服务端文件卷、游戏数据卷和插件卷）
// serverID: 服务器ID
// 返回: 错误信息；任一卷创建失败时会清理已创建的卷
func (dm *DockerManager) CreateVolume(serverID uint) error {
	for _, name := range serverVolumeNames(serverID) {
		if err := dm.createSingleVolume(name); err != nil {
			if rmErr := dm.RemoveVolume(serverID); rmErr != nil {
				utils.Warn("清理已创建的卷失败", zap.Uint("server_id", serverID), zap.Error(rmErr))
			}
			return fmt.Errorf("创建卷 %s 失败: %v", name, err)
		}
	}

	utils.Info("Docker卷创建成功", zap.Strings("volumes", serverVolumeNames(serverID)))
	return nil
}

// createSingleVolume 创建单个Docker卷
// volumeName: 卷名称
// 返回: 错误信息
func (dm *DockerManager) createSingleVolume(volumeName string) error {
	// 检查卷是否已存在
	exists, err := dm.VolumeExists(volumeName)
	if err != nil {
		return fmt.Errorf("检查卷是否存在失败: %v", err)
	}
	if exists {
		utils.Debug("Docker卷已存在，跳过创建", zap.String("volume", volumeName))
		return nil
	}

	// 创建卷
	utils.Infof("正在创建Docker卷: %s", volumeName)
	volumeCreateBody := volume.CreateOptions{
		Name: volumeName,
	}

	_, err = dm.client.VolumeCreate(dm.ctx, volumeCreateBody)
	if err != nil {
		return fmt.Errorf("创建Docker卷失败: %v", err)
	}

	utils.Info("Docker卷创建成功", zap.String("volume", volumeName))
	return nil
}

// RemoveVolume 删除Docker卷（包括服务端文件卷、游戏数据卷和插件卷）
// serverID: 服务器ID
// 返回: 错误信息；每个卷都会尝试删除，任一失败都会返回错误
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
// volumeName: 卷名称
// 返回: 错误信息
func (dm *DockerManager) removeSingleVolume(volumeName string) error {
	// 检查卷是否存在
	exists, err := dm.VolumeExists(volumeName)
	if err != nil {
		return fmt.Errorf("检查卷是否存在失败: %v", err)
	}
	if !exists {
		utils.Debug("Docker卷不存在，跳过删除", zap.String("volume", volumeName))
		return nil
	}

	utils.Infof("正在删除Docker卷: %s", volumeName)
	err = dm.client.VolumeRemove(dm.ctx, volumeName, false)
	if err != nil {
		return fmt.Errorf("删除Docker卷失败: %v", err)
	}

	utils.Info("Docker卷删除成功", zap.String("volume", volumeName))
	return nil
}

// VolumeExists 检查Docker卷是否存在
// volumeName: 卷名称
// 返回: 是否存在和错误信息
func (dm *DockerManager) VolumeExists(volumeName string) (bool, error) {
	// 尝试获取卷信息
	_, err := dm.client.VolumeInspect(dm.ctx, volumeName)
	if err != nil {
		if errdefs.IsNotFound(err) {
			return false, nil // 卷不存在
		}
		return false, fmt.Errorf("检查Docker卷失败: %v", err)
	}

	return true, nil
}
