package docker_manager

import (
	"ark-server-commander/utils"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"
	"sync"

	"github.com/containerd/errdefs"
	"github.com/docker/docker/api/types/image"
	"go.uber.org/zap"
)

// ImagePullProgress 镜像拉取进度信息
type ImagePullProgress struct {
	Status         string `json:"status"`
	Progress       string `json:"progress"`
	ProgressDetail *struct {
		Current int64 `json:"current"`
		Total   int64 `json:"total"`
	} `json:"progressDetail"`
	ID    string `json:"id"`
	Error string `json:"error"`
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
	HasUpdate    bool                    `json:"has_update"`    // 最近一次检查更新的结果
	CurrentLayer string                  `json:"current_layer"` // 当前下载的层级
	Layers       map[string]*LayerStatus `json:"layers"`        // 所有层级的状态
}

// 全局镜像状态管理
var (
	imagePullStatus       = make(map[string]bool)                    // 记录镜像是否正在拉取
	imagePullCurrentLayer = make(map[string]string)                  // 记录当前下载的层级
	imagePullLayers       = make(map[string]map[string]*LayerStatus) // 记录所有层级的状态
	imageUpdateAvailable  = make(map[string]bool)                    // 记录最近一次检查更新的结果，拉取成功后清除
	imagePullError        = make(map[string]string)                  // 最近一次拉取失败的原因，下次拉取开始时清除
	imagePullMutex        sync.RWMutex
)

// ClaimPull 原子地占用镜像的拉取权，返回 false 表示该镜像已经在拉取中。
// 检查和占用在同一把锁里完成，并发的请求只有一个能成功。
// 占用成功后，调用方必须调用 PullClaimedImage，它会在结束时释放占用
func ClaimPull(imageName string) bool {
	imagePullMutex.Lock()
	defer imagePullMutex.Unlock()
	if imagePullStatus[imageName] {
		return false
	}
	imagePullStatus[imageName] = true
	imagePullCurrentLayer[imageName] = ""
	imagePullLayers[imageName] = make(map[string]*LayerStatus)
	delete(imagePullError, imageName)
	return true
}

// finishPull 结束拉取：释放占用，并记录失败原因供镜像状态接口返回
func finishPull(imageName string, err error) {
	imagePullMutex.Lock()
	defer imagePullMutex.Unlock()
	imagePullStatus[imageName] = false
	delete(imagePullCurrentLayer, imageName)
	delete(imagePullLayers, imageName)
	if err != nil {
		imagePullError[imageName] = err.Error()
	} else {
		delete(imagePullError, imageName)
	}
}

// PullClaimedImage 拉取已经通过 ClaimPull 占用的镜像并记录进度，结束时释放占用
func (dm *DockerManager) PullClaimedImage(imageName string) error {
	err := dm.pullImage(imageName)
	finishPull(imageName, err)
	return err
}

// pullImage 拉取镜像并更新进度信息，不负责占用和释放
func (dm *DockerManager) pullImage(imageName string) error {
	utils.Info("开始拉取Docker镜像", zap.String("image", imageName))

	// 拉取镜像
	reader, err := dm.client.ImagePull(dm.ctx, imageName, image.PullOptions{})
	if err != nil {
		return fmt.Errorf("拉取Docker镜像失败: %v", err)
	}
	defer reader.Close()

	// 逐条解析拉取进度（每条是一个 JSON 对象）
	decoder := json.NewDecoder(reader)
	for {
		var progressInfo ImagePullProgress
		if err := decoder.Decode(&progressInfo); err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("读取镜像拉取进度失败: %v", err)
		}

		// 拉取过程中的错误（如网络中断、镜像不存在）出现在进度流里，而不是 ImagePull 的返回值
		if progressInfo.Error != "" {
			return fmt.Errorf("拉取Docker镜像失败: %s", progressInfo.Error)
		}

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
	}

	imagePullMutex.Lock()
	imageUpdateAvailable[imageName] = false
	imagePullMutex.Unlock()

	utils.Info("Docker镜像拉取成功", zap.String("image", imageName))

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
	status.Ready = exists

	// 检查是否正在拉取中（镜像已存在时可能是在更新）
	imagePullMutex.RLock()
	status.HasUpdate = imageUpdateAvailable[imageName]
	status.Pulling = imagePullStatus[imageName]
	if msg := imagePullError[imageName]; msg != "" && !status.Pulling {
		status.Error = "拉取镜像失败: " + msg
	}
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
// 比较镜像仓库中该标签当前的 manifest digest 和本地镜像的 RepoDigests，结果会记录下来供 GetImageStatus 返回
// imageName: 镜像名称
// 返回: 是否有更新和错误信息
func (dm *DockerManager) CheckImageUpdate(imageName string) (bool, error) {
	hasUpdate, err := dm.checkImageUpdate(imageName)
	if err != nil {
		return false, err
	}

	imagePullMutex.Lock()
	imageUpdateAvailable[imageName] = hasUpdate
	imagePullMutex.Unlock()

	return hasUpdate, nil
}

func (dm *DockerManager) checkImageUpdate(imageName string) (bool, error) {
	localImage, err := dm.client.ImageInspect(dm.ctx, imageName)
	if err != nil {
		// 本地镜像不存在，需要下载
		if errdefs.IsNotFound(err) {
			return true, nil
		}
		return false, fmt.Errorf("获取本地镜像信息失败: %v", err)
	}

	// 只查询镜像仓库中的 manifest，不下载镜像层
	remote, err := dm.client.DistributionInspect(dm.ctx, imageName, "")
	if err != nil {
		return false, fmt.Errorf("获取远程镜像信息失败: %v", err)
	}

	// RepoDigests 形如 tbro98/ase-server@sha256:...，本地构建的镜像没有 RepoDigests，会被视为有更新
	remoteDigest := "@" + remote.Descriptor.Digest.String()
	for _, repoDigest := range localImage.RepoDigests {
		if strings.HasSuffix(repoDigest, remoteDigest) {
			return false, nil
		}
	}

	return true, nil
}

// IsContainerImageOutdated 判断容器创建时使用的镜像是否已不是 imageName 当前指向的镜像
// containerName: 容器名称
// imageName: 镜像名称
// 返回: 是否过期和错误信息
func (dm *DockerManager) IsContainerImageOutdated(containerName, imageName string) (bool, error) {
	containerInfo, err := dm.client.ContainerInspect(dm.ctx, containerName)
	if err != nil {
		return false, fmt.Errorf("获取Docker容器信息失败: %v", err)
	}

	imageInfo, err := dm.client.ImageInspect(dm.ctx, imageName)
	if err != nil {
		return false, fmt.Errorf("获取Docker镜像信息失败: %v", err)
	}

	return containerInfo.Image != imageInfo.ID, nil
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
