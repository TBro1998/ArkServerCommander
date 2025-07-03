package docker_manager

import (
	"ark-server-manager/utils"
	"fmt"

	"github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/client"
)

// CreateVolume 创建Docker卷
// serverID: 服务器ID
// 返回: 卷名称和错误信息
func (dm *DockerManager) CreateVolume(serverID uint) (string, error) {
	volumeName := utils.GetServerVolumeName(serverID)

	// 检查卷是否已存在
	exists, err := dm.VolumeExists(volumeName)
	if err != nil {
		return "", fmt.Errorf("检查卷是否存在失败: %v", err)
	}
	if exists {
		fmt.Printf("Docker卷 %s 已存在，跳过创建\n", volumeName)
		return volumeName, nil
	}

	// 创建卷
	fmt.Printf("正在创建Docker卷: %s\n", volumeName)
	volumeCreateBody := volume.CreateOptions{
		Name: volumeName,
	}

	_, err = dm.client.VolumeCreate(dm.ctx, volumeCreateBody)
	if err != nil {
		return "", fmt.Errorf("创建Docker卷失败: %v", err)
	}

	fmt.Printf("Docker卷创建成功: %s\n", volumeName)
	return volumeName, nil
}

// RemoveVolume 删除Docker卷
// volumeName: 卷名称
// 返回: 错误信息
func (dm *DockerManager) RemoveVolume(volumeName string) error {
	// 检查卷是否存在
	exists, err := dm.VolumeExists(volumeName)
	if err != nil {
		return fmt.Errorf("检查卷是否存在失败: %v", err)
	}
	if !exists {
		fmt.Printf("Docker卷 %s 不存在，跳过删除\n", volumeName)
		return nil
	}

	fmt.Printf("正在删除Docker卷: %s\n", volumeName)
	err = dm.client.VolumeRemove(dm.ctx, volumeName, false)
	if err != nil {
		return fmt.Errorf("删除Docker卷失败: %v", err)
	}

	fmt.Printf("Docker卷删除成功: %s\n", volumeName)
	return nil
}

// VolumeExists 检查Docker卷是否存在
// volumeName: 卷名称
// 返回: 是否存在和错误信息
func (dm *DockerManager) VolumeExists(volumeName string) (bool, error) {
	// 尝试获取卷信息
	_, err := dm.client.VolumeInspect(dm.ctx, volumeName)
	if err != nil {
		if client.IsErrNotFound(err) {
			return false, nil // 卷不存在
		}
		return false, fmt.Errorf("检查Docker卷失败: %v", err)
	}

	return true, nil
}
