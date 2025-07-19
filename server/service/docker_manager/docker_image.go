package docker_manager

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/containerd/errdefs"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
)

// ImagePullProgress 镜像拉取进度信息
type ImagePullProgress struct {
	Status         string `json:"status"`
	Progress       string `json:"progress"`
	ProgressDetail *struct {
		Current int64 `json:"current"`
		Total   int64 `json:"total"`
	} `json:"progressDetail"`
	ID string `json:"id"`
}

// LayerStatus 层级状态信息
type LayerStatus struct {
	ID       string `json:"id"`       // 层级ID
	Size     int64  `json:"size"`     // 层级大小
	Progress int64  `json:"progress"` // 已下载大小
	Status   string `json:"status"`   // 层级状态：downloading, extracting, verifying, complete
}

// ImageStatus 镜像状态信息
type ImageStatus struct {
	Exists       bool                    `json:"exists"`        // 镜像是否存在
	Pulling      bool                    `json:"pulling"`       // 是否正在拉取
	Ready        bool                    `json:"ready"`         // 是否准备就绪
	Error        string                  `json:"error"`         // 错误信息
	CurrentLayer string                  `json:"current_layer"` // 当前下载的层级
	Layers       map[string]*LayerStatus `json:"layers"`        // 所有层级的状态
}

// 全局镜像状态管理
var (
	imagePullStatus       = make(map[string]bool)                    // 记录镜像是否正在拉取
	imagePullCurrentLayer = make(map[string]string)                  // 记录当前下载的层级
	imagePullLayers       = make(map[string]map[string]*LayerStatus) // 记录所有层级的状态
	imagePullMutex        sync.RWMutex
)

// PullImageWithProgress 拉取Docker镜像并显示进度
// imageName: 镜像名称
// 返回: 错误信息
func (dm *DockerManager) PullImageWithProgress(imageName string) error {
	fmt.Printf("正在拉取Docker镜像: %s\n", imageName)

	// 设置拉取状态
	imagePullMutex.Lock()
	imagePullStatus[imageName] = true
	imagePullCurrentLayer[imageName] = ""
	imagePullLayers[imageName] = make(map[string]*LayerStatus)
	imagePullMutex.Unlock()

	// 确保在函数结束时清理状态
	defer func() {
		imagePullMutex.Lock()
		imagePullStatus[imageName] = false
		delete(imagePullCurrentLayer, imageName)
		delete(imagePullLayers, imageName)
		imagePullMutex.Unlock()
	}()

	// 拉取镜像
	reader, err := dm.client.ImagePull(dm.ctx, imageName, image.PullOptions{})
	if err != nil {
		return fmt.Errorf("拉取Docker镜像失败: %v", err)
	}
	defer reader.Close()

	// 读取并显示拉取进度
	buffer := make([]byte, 1024)
	for {
		n, err := reader.Read(buffer)
		if n > 0 {
			// 解析JSON进度信息
			progress := string(buffer[:n])
			lines := strings.Split(progress, "\n")

			for _, line := range lines {
				if line == "" {
					continue
				}

				var progressInfo ImagePullProgress
				if err := json.Unmarshal([]byte(line), &progressInfo); err == nil {
					// 更新进度信息
					imagePullMutex.Lock()

					// 处理层级信息
					layerID := progressInfo.ID
					if layerID != "" {
						imagePullCurrentLayer[imageName] = layerID

						// 确保层级存在
						if imagePullLayers[imageName][layerID] == nil {
							imagePullLayers[imageName][layerID] = &LayerStatus{
								ID:       layerID,
								Size:     0,
								Progress: 0,
								Status:   "pending",
							}
						}

						// 更新层级大小信息
						if progressInfo.ProgressDetail != nil {
							if progressInfo.ProgressDetail.Total > 0 {
								imagePullLayers[imageName][layerID].Size = progressInfo.ProgressDetail.Total
							}
							if progressInfo.ProgressDetail.Current > 0 {
								imagePullLayers[imageName][layerID].Progress = progressInfo.ProgressDetail.Current
							}
						}

						// 如果ProgressDetail中没有大小信息，尝试从Progress字符串中解析
						if imagePullLayers[imageName][layerID].Size == 0 && progressInfo.Progress != "" {
							if size := parseSizeFromProgress(progressInfo.Progress); size > 0 {
								imagePullLayers[imageName][layerID].Size = size
							}
						}

						// 更新层级状态
						if strings.Contains(progressInfo.Status, "Downloading") {
							imagePullLayers[imageName][layerID].Status = "downloading"
						} else if strings.Contains(progressInfo.Status, "Extracting") {
							imagePullLayers[imageName][layerID].Status = "extracting"
						} else if strings.Contains(progressInfo.Status, "Verifying") {
							imagePullLayers[imageName][layerID].Status = "verifying"
						} else if strings.Contains(progressInfo.Status, "Pull complete") || strings.Contains(progressInfo.Status, "complete") {
							imagePullLayers[imageName][layerID].Status = "complete"
							// 完成时，如果没有大小信息，设置一个默认值
							if imagePullLayers[imageName][layerID].Size == 0 {
								imagePullLayers[imageName][layerID].Size = imagePullLayers[imageName][layerID].Progress
								if imagePullLayers[imageName][layerID].Size == 0 {
									imagePullLayers[imageName][layerID].Size = 1024 * 1024 // 默认1MB
								}
							}
							imagePullLayers[imageName][layerID].Progress = imagePullLayers[imageName][layerID].Size
						}
					}

					imagePullMutex.Unlock()

					// 打印进度到控制台
					if strings.Contains(progressInfo.Status, "Downloading") || strings.Contains(progressInfo.Status, "Extracting") {
						fmt.Printf("\r[%s] %s", imageName, progressInfo.Status)
					}
				}
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("读取镜像拉取进度失败: %v", err)
		}
	}

	fmt.Printf("\nDocker镜像拉取成功: %s\n", imageName)

	return nil
}

// ImageExists 检查Docker镜像是否存在
// imageName: 镜像名称
// 返回: 是否存在和错误信息
func (dm *DockerManager) ImageExists(imageName string) (bool, error) {
	// 尝试获取镜像信息
	_, err := dm.client.ImageInspect(dm.ctx, imageName)
	if err != nil {
		if errdefs.IsNotFound(err) {
			return false, nil // 镜像不存在
		}
		return false, fmt.Errorf("检查Docker镜像失败: %v", err)
	}

	return true, nil
}

// GetImageStatus 获取镜像状态和进度信息
// imageName: 镜像名称
// 返回: 镜像状态信息
func (dm *DockerManager) GetImageStatus(imageName string) *ImageStatus {
	status := &ImageStatus{
		Exists:       false,
		Pulling:      false,
		Ready:        false,
		Error:        "",
		CurrentLayer: "",
		Layers:       make(map[string]*LayerStatus),
	}

	// 检查镜像是否存在
	exists, err := dm.ImageExists(imageName)
	if err != nil {
		status.Error = fmt.Sprintf("检查镜像失败: %v", err)
		return status
	}

	status.Exists = exists
	if exists {
		status.Ready = true
		return status
	}

	// 检查是否正在拉取中
	imagePullMutex.RLock()
	status.Pulling = imagePullStatus[imageName]
	if status.Pulling {
		status.CurrentLayer = imagePullCurrentLayer[imageName]

		// 复制层级状态
		if layers, exists := imagePullLayers[imageName]; exists {
			for layerID, layerStatus := range layers {
				status.Layers[layerID] = &LayerStatus{
					ID:       layerStatus.ID,
					Size:     layerStatus.Size,
					Progress: layerStatus.Progress,
					Status:   layerStatus.Status,
				}
			}
		}
	}
	imagePullMutex.RUnlock()

	return status
}

// IsImagePulling 检查镜像是否正在拉取中
func IsImagePulling(imageName string) bool {
	imagePullMutex.RLock()
	defer imagePullMutex.RUnlock()
	return imagePullStatus[imageName]
}

// WaitForImage 等待镜像拉取完成（已废弃，请使用 GetImageStatus）
// imageName: 镜像名称
// timeout: 超时时间（秒）
// 返回: 是否成功和错误信息
func (dm *DockerManager) WaitForImage(imageName string, timeout int) (bool, error) {
	log.Printf("等待镜像 %s 拉取完成...", imageName)

	// 每秒检查一次镜像是否存在
	for i := 0; i < timeout; i++ {
		exists, err := dm.ImageExists(imageName)
		if err != nil {
			return false, fmt.Errorf("检查镜像失败: %v", err)
		}

		if exists {
			log.Printf("镜像 %s 已准备就绪", imageName)
			return true, nil
		}

		// 检查是否正在拉取中
		if IsImagePulling(imageName) {
			log.Printf("镜像 %s 正在拉取中，继续等待...", imageName)
		}

		// 等待1秒后再次检查
		time.Sleep(1 * time.Second)
	}

	return false, fmt.Errorf("等待镜像 %s 超时（%d秒）", imageName, timeout)
}

// parseSizeFromProgress 从进度字符串中解析大小信息
// progress: 进度字符串，如 "1.5MB/2.0MB"
// 返回: 解析出的大小（字节）
func parseSizeFromProgress(progress string) int64 {
	// 移除空格
	progress = strings.TrimSpace(progress)

	// 查找 "/" 分隔符
	parts := strings.Split(progress, "/")
	if len(parts) != 2 {
		return 0
	}

	// 解析总大小部分
	totalStr := strings.TrimSpace(parts[1])
	return parseSizeString(totalStr)
}

// CheckImageUpdate 检查镜像是否有更新
// imageName: 镜像名称
// 返回: 是否有更新和错误信息
func (dm *DockerManager) CheckImageUpdate(imageName string) (bool, error) {
	// 获取本地镜像信息
	localImageInfo, err := dm.GetImageInfo(imageName)
	if err != nil {
		// 如果本地镜像不存在，认为有更新（需要下载）
		if errdefs.IsNotFound(err) {
			return true, nil
		}
		return false, fmt.Errorf("获取本地镜像信息失败: %v", err)
	}

	// 尝试拉取最新镜像信息（不实际下载）
	// 这里简化实现，实际应该比较镜像的digest
	// 由于Docker API限制，我们暂时返回false（没有更新）
	// 在实际生产环境中，可以通过比较镜像的创建时间或digest来判断
	_ = localImageInfo

	// TODO: 实现真正的镜像更新检查逻辑
	// 可以通过Docker Registry API来获取远程镜像信息并比较
	return false, nil
}

// GetImageInfo 获取本地镜像详细信息
// imageName: 镜像名称
// 返回: 镜像信息和错误信息
func (dm *DockerManager) GetImageInfo(imageName string) (*ImageInfo, error) {
	imageInspect, err := dm.client.ImageInspect(dm.ctx, imageName)
	if err != nil {
		return nil, err
	}

	return &ImageInfo{
		ID:      imageInspect.ID,
		Tags:    imageInspect.RepoTags,
		Size:    imageInspect.Size,
		Created: imageInspect.Created,
	}, nil
}

// ImageInfo 镜像信息结构体
type ImageInfo struct {
	ID      string   `json:"id"`      // 镜像ID
	Tags    []string `json:"tags"`    // 镜像标签
	Size    int64    `json:"size"`    // 镜像大小
	Created string   `json:"created"` // 创建时间
}

// RemoveOldImage 删除旧版本镜像
// imageName: 镜像名称
// keepLatest: 是否保留最新版本
// 返回: 错误信息
func (dm *DockerManager) RemoveOldImage(imageName string, keepLatest bool) error {
	if !keepLatest {
		// 删除指定镜像
		_, err := dm.client.ImageRemove(dm.ctx, imageName, image.RemoveOptions{
			Force:         true,
			PruneChildren: true,
		})
		if err != nil {
			return fmt.Errorf("删除镜像失败: %v", err)
		}
		fmt.Printf("镜像 %s 已删除\n", imageName)
	} else {
		fmt.Println("保留最新镜像，跳过删除")
	}

	return nil
}

// GetContainersByImage 获取使用指定镜像的所有容器
// imageName: 镜像名称
// 返回: 容器信息列表和错误信息
func (dm *DockerManager) GetContainersByImage(imageName string) ([]ContainerInfo, error) {
	containers, err := dm.client.ContainerList(dm.ctx, container.ListOptions{
		All: true,
	})
	if err != nil {
		return nil, fmt.Errorf("获取容器列表失败: %v", err)
	}

	var result []ContainerInfo
	for _, c := range containers {
		if c.Image == imageName {
			result = append(result, ContainerInfo{
				ID:     c.ID,
				Name:   c.Names[0],
				Image:  c.Image,
				Status: c.Status,
				State:  c.State,
			})
		}
	}

	return result, nil
}

// ContainerInfo 容器信息结构体
type ContainerInfo struct {
	ID     string `json:"id"`     // 容器ID
	Name   string `json:"name"`   // 容器名称
	Image  string `json:"image"`  // 镜像名称
	Status string `json:"status"` // 容器状态
	State  string `json:"state"`  // 容器状态
}

// GetImageHistory 获取镜像历史信息
// imageName: 镜像名称
// 返回: 历史信息和错误信息
func (dm *DockerManager) GetImageHistory(imageName string) ([]ImageHistoryEntry, error) {
	history, err := dm.client.ImageHistory(dm.ctx, imageName)
	if err != nil {
		return nil, fmt.Errorf("获取镜像历史失败: %v", err)
	}

	var result []ImageHistoryEntry
	for _, h := range history {
		result = append(result, ImageHistoryEntry{
			ID:        h.ID,
			Created:   h.Created,
			CreatedBy: h.CreatedBy,
			Size:      h.Size,
			Comment:   h.Comment,
		})
	}

	return result, nil
}

// ImageHistoryEntry 镜像历史条目
type ImageHistoryEntry struct {
	ID        string `json:"id"`         // 层ID
	Created   int64  `json:"created"`    // 创建时间
	CreatedBy string `json:"created_by"` // 创建命令
	Size      int64  `json:"size"`       // 层大小
	Comment   string `json:"comment"`    // 注释
}

// parseSizeString 解析大小字符串
// sizeStr: 大小字符串，如 "2.0MB", "1.5GB"
// 返回: 字节数
func parseSizeString(sizeStr string) int64 {
	sizeStr = strings.ToLower(strings.TrimSpace(sizeStr))

	// 移除单位
	var multiplier int64 = 1
	if strings.HasSuffix(sizeStr, "kb") {
		multiplier = 1024
		sizeStr = strings.TrimSuffix(sizeStr, "kb")
	} else if strings.HasSuffix(sizeStr, "mb") {
		multiplier = 1024 * 1024
		sizeStr = strings.TrimSuffix(sizeStr, "mb")
	} else if strings.HasSuffix(sizeStr, "gb") {
		multiplier = 1024 * 1024 * 1024
		sizeStr = strings.TrimSuffix(sizeStr, "gb")
	} else if strings.HasSuffix(sizeStr, "b") {
		multiplier = 1
		sizeStr = strings.TrimSuffix(sizeStr, "b")
	}

	// 解析数字
	if size, err := strconv.ParseFloat(sizeStr, 64); err == nil {
		return int64(size * float64(multiplier))
	}

	return 0
}
