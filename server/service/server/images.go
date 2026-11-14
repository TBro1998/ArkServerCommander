package server

import (
	"fmt"
	"sync"
	"time"

	"ark-server-commander/database"
	"ark-server-commander/models"
	"ark-server-commander/service/docker_manager"
	"ark-server-commander/utils"

	"go.uber.org/zap"
)

// pullingImages 启动前正在下载镜像的服务器 ID -> 镜像名，用于在界面上显示"下载镜像"阶段
var pullingImages sync.Map

// pullingImageFor 返回服务器启动时正在下载的镜像
func pullingImageFor(serverID uint) (string, bool) {
	v, ok := pullingImages.Load(serverID)
	if !ok {
		return "", false
	}
	return v.(string), true
}

// imageWaitTimeout 等待别处发起的拉取结束的最长时间
const imageWaitTimeout = 2 * time.Hour

// ensureServerImage 服务器使用的 ARK 镜像不存在时下载它（启用 ArkApi 的服务器第一次启动时通常要下载 Proton 镜像）。
// 已经有人在下载同一个镜像（首页镜像管理、另一台服务器的启动）时等它结束
func ensureServerImage(serverID uint, imageName string, dm *docker_manager.DockerManager) error {
	exists, err := dm.ImageExists(imageName)
	if err != nil {
		return fmt.Errorf("检查镜像失败: %w", err)
	}
	if exists {
		return nil
	}

	pullingImages.Store(serverID, imageName)
	defer pullingImages.Delete(serverID)

	if docker_manager.ClaimPull(imageName) {
		utils.Info("服务器使用的镜像不存在，开始下载", zap.Uint("server_id", serverID), zap.String("image", imageName))
		if err := dm.PullClaimedImage(imageName); err != nil {
			return fmt.Errorf("下载镜像 %s 失败: %w", imageName, err)
		}
		return nil
	}

	deadline := time.Now().Add(imageWaitTimeout)
	for docker_manager.IsImagePulling(imageName) {
		if time.Now().After(deadline) {
			return fmt.Errorf("等待镜像 %s 下载超时", imageName)
		}
		time.Sleep(2 * time.Second)
	}
	if exists, err := dm.ImageExists(imageName); err != nil {
		return fmt.Errorf("检查镜像失败: %w", err)
	} else if !exists {
		return fmt.Errorf("镜像 %s 下载失败，请在首页镜像管理查看原因", imageName)
	}
	return nil
}

// requiredImages 当前需要的镜像：ARK 镜像只算有服务器在用的，
// 没有服务器时算 Linux 版（新建服务器默认不启用 ArkApi）
func requiredImages() (map[string]bool, error) {
	var servers []models.Server
	if err := database.DB.Select("id", "arkapi").Find(&servers).Error; err != nil {
		return nil, fmt.Errorf("获取服务器列表失败: %w", err)
	}
	required := map[string]bool{}
	if len(servers) == 0 {
		required[utils.ArkServerLinuxImage] = true
	}
	for _, s := range servers {
		required[utils.ArkServerImageFor(s.ArkApi)] = true
	}
	return required, nil
}
