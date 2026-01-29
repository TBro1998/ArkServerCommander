package docker_manager

import (
	"ark-server-commander/utils"
	"fmt"

	"github.com/containerd/errdefs"
	"github.com/docker/docker/api/types/volume"
	"go.uber.org/zap"
)

// CreateVolume 创建Docker卷（包括游戏数据卷和插件卷）
// serverID: 服务器ID
// 返回: 游戏数据卷名称和错误信息
func (dm *DockerManager) CreateVolume(serverID uint) (string, error) {
	volumeName := utils.GetServerVolumeName(serverID)
	pluginsVolumeName := utils.GetServerPluginsVolumeName(serverID)

	// 创建游戏数据卷
	if err := dm.createSingleVolume(volumeName); err != nil {
		return "", fmt.Errorf("创建游戏数据卷失败: %v", err)
	}

	// 创建插件卷
	if err := dm.createSingleVolume(pluginsVolumeName); err != nil {
		// 如果插件卷创建失败，清理已创建的游戏数据卷
		dm.RemoveVolume(volumeName)
		return "", fmt.Errorf("创建插件卷失败: %v", err)
	}

	utils.Info("Docker卷创建成功", zap.String("data_volume", volumeName), zap.String("plugins_volume", pluginsVolumeName))
	return volumeName, nil
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

// RemoveVolume 删除Docker卷（包括游戏数据卷和插件卷）
// volumeName: 游戏数据卷名称
// 返回: 错误信息
func (dm *DockerManager) RemoveVolume(volumeName string) error {
	// 从游戏数据卷名称推导出插件卷名称
	// 这里假设volumeName是游戏数据卷名称，需要推导出插件卷名称
	// 由于当前的设计，我们需要通过serverID来推导，但这里只有volumeName
	// 所以我们先删除游戏数据卷，然后尝试删除可能的插件卷

	// 删除游戏数据卷
	if err := dm.removeSingleVolume(volumeName); err != nil {
		return err
	}

	// 尝试删除插件卷（通过替换卷名称）
	pluginsVolumeName := volumeName + "-plugins"
	if err := dm.removeSingleVolume(pluginsVolumeName); err != nil {
		// 插件卷删除失败不影响主流程，只记录警告
		utils.Warn("删除插件卷失败", zap.String("volume", pluginsVolumeName), zap.Error(err))
	}

	return nil
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
