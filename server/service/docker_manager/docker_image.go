package docker_manager

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/containerd/errdefs"
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

// ImageStatus 镜像状态信息
type ImageStatus struct {
	Exists         bool   `json:"exists"`          // 镜像是否存在
	Pulling        bool   `json:"pulling"`         // 是否正在拉取
	Progress       string `json:"progress"`        // 下载进度信息
	Ready          bool   `json:"ready"`           // 是否准备就绪
	Error          string `json:"error"`           // 错误信息
	Percent        int    `json:"percent"`         // 拉取进度百分比
	Size           int64  `json:"size"`            // 镜像总大小（字节）
	DownloadedSize int64  `json:"downloaded_size"` // 已下载大小（字节）
}

// 全局镜像状态管理
var (
	imagePullStatus     = make(map[string]bool)   // 记录镜像是否正在拉取
	imagePullProgress   = make(map[string]string) // 记录镜像拉取进度
	imagePullPercent    = make(map[string]int)    // 记录镜像拉取百分比
	imagePullSize       = make(map[string]int64)  // 记录镜像总大小
	imagePullDownloaded = make(map[string]int64)  // 记录已下载大小
	imagePullMutex      sync.RWMutex
)

// PullImageWithProgress 拉取Docker镜像并显示进度
// imageName: 镜像名称
// 返回: 错误信息
func (dm *DockerManager) PullImageWithProgress(imageName string) error {
	fmt.Printf("正在拉取Docker镜像: %s\n", imageName)

	// 设置拉取状态
	imagePullMutex.Lock()
	imagePullStatus[imageName] = true
	imagePullProgress[imageName] = "正在准备下载..."
	imagePullPercent[imageName] = 5
	imagePullSize[imageName] = 0
	imagePullDownloaded[imageName] = 0
	imagePullMutex.Unlock()

	// 确保在函数结束时清理状态
	defer func() {
		imagePullMutex.Lock()
		imagePullStatus[imageName] = false
		delete(imagePullProgress, imageName)
		delete(imagePullPercent, imageName)
		delete(imagePullSize, imageName)
		delete(imagePullDownloaded, imageName)
		imagePullMutex.Unlock()
	}()

	// 拉取镜像
	reader, err := dm.client.ImagePull(dm.ctx, imageName, image.PullOptions{})
	if err != nil {
		imagePullMutex.Lock()
		imagePullProgress[imageName] = fmt.Sprintf("拉取失败: %v", err)
		imagePullMutex.Unlock()
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
					if progressInfo.Progress != "" {
						imagePullProgress[imageName] = progressInfo.Progress
					} else if progressInfo.Status != "" {
						// 翻译状态信息为中文
						statusText := progressInfo.Status
						if strings.Contains(statusText, "Downloading") {
							statusText = "正在下载"
						} else if strings.Contains(statusText, "Extracting") {
							statusText = "正在解压"
						} else if strings.Contains(statusText, "Verifying") {
							statusText = "正在验证"
						} else if strings.Contains(statusText, "Pull complete") {
							statusText = "拉取完成"
						} else if strings.Contains(statusText, "Already up to date") {
							statusText = "镜像已是最新"
						}
						imagePullProgress[imageName] = statusText
					}

					// 解析进度详情
					if progressInfo.ProgressDetail != nil {
						current := progressInfo.ProgressDetail.Current
						total := progressInfo.ProgressDetail.Total

						if total > 0 {
							// 更新大小信息
							imagePullSize[imageName] = total
							imagePullDownloaded[imageName] = current

							// 根据状态计算百分比
							var percent int
							if strings.Contains(progressInfo.Status, "Downloading") {
								// 下载阶段：0-70%
								downloadPercent := int((float64(current) / float64(total)) * 70)
								percent = downloadPercent
							} else if strings.Contains(progressInfo.Status, "Extracting") {
								// 解压阶段：70-90%
								extractPercent := int((float64(current) / float64(total)) * 20)
								percent = 70 + extractPercent
							} else if strings.Contains(progressInfo.Status, "Verifying") {
								// 验证阶段：90-95%
								percent = 90 + int((float64(current)/float64(total))*5)
							} else if strings.Contains(progressInfo.Status, "Pull complete") || strings.Contains(progressInfo.Status, "complete") {
								// 完成：100%
								percent = 100
							} else {
								// 其他状态：基于当前进度
								percent = int((float64(current) / float64(total)) * 100)
							}

							imagePullPercent[imageName] = percent
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
			imagePullMutex.Lock()
			imagePullProgress[imageName] = fmt.Sprintf("读取进度失败: %v", err)
			imagePullMutex.Unlock()
			return fmt.Errorf("读取镜像拉取进度失败: %v", err)
		}
	}

	fmt.Printf("\nDocker镜像拉取成功: %s\n", imageName)

	// 设置完成状态
	imagePullMutex.Lock()
	imagePullProgress[imageName] = "镜像拉取完成"
	imagePullPercent[imageName] = 100
	imagePullMutex.Unlock()

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
		Exists:         false,
		Pulling:        false,
		Progress:       "",
		Ready:          false,
		Error:          "",
		Percent:        0,
		Size:           0,
		DownloadedSize: 0,
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
		status.Progress = "镜像已存在"
		status.Percent = 100

		// 获取镜像大小信息
		if imageInfo, err := dm.client.ImageInspect(dm.ctx, imageName); err == nil {
			log.Printf("镜像 %s 的Size: %d bytes", imageName, imageInfo.Size)

			// 使用Size字段（这是镜像的实际大小）
			if imageInfo.Size > 0 {
				status.Size = imageInfo.Size
				log.Printf("使用Size字段作为镜像大小: %d bytes", status.Size)
			}

			// 如果通过ImageInspect获取的大小仍然很小，尝试使用ImageList获取更准确的大小
			if status.Size < 1024*1024 {
				log.Printf("镜像大小过小，尝试使用ImageList获取更准确的大小")
				if accurateSize, err := dm.getImageSizeFromList(imageName); err == nil && accurateSize > 0 {
					status.Size = accurateSize
					log.Printf("通过ImageList获取到更准确的大小: %d bytes", status.Size)
				}
			}
		}

		return status
	}

	// 检查是否正在拉取中
	imagePullMutex.RLock()
	status.Pulling = imagePullStatus[imageName]
	if status.Pulling {
		status.Progress = imagePullProgress[imageName]
		status.Percent = imagePullPercent[imageName]
		status.Size = imagePullSize[imageName]
		status.DownloadedSize = imagePullDownloaded[imageName]
	}
	imagePullMutex.RUnlock()

	if !status.Pulling {
		status.Progress = "镜像不存在，需要下载"
		status.Percent = 0
		status.Size = 0
		status.DownloadedSize = 0
	}

	return status
}

// IsImagePulling 检查镜像是否正在拉取中
func IsImagePulling(imageName string) bool {
	imagePullMutex.RLock()
	defer imagePullMutex.RUnlock()
	return imagePullStatus[imageName]
}

// GetImagePullProgress 获取镜像拉取进度
func GetImagePullProgress(imageName string) string {
	imagePullMutex.RLock()
	defer imagePullMutex.RUnlock()
	return imagePullProgress[imageName]
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

// getImageSizeFromList 通过ImageList获取镜像大小
// imageName: 镜像名称
// 返回: 镜像大小（字节）和错误信息
func (dm *DockerManager) getImageSizeFromList(imageName string) (int64, error) {
	// 获取镜像列表
	images, err := dm.client.ImageList(dm.ctx, image.ListOptions{})
	if err != nil {
		return 0, fmt.Errorf("获取镜像列表失败: %v", err)
	}

	// 查找指定镜像
	for _, img := range images {
		// 检查镜像标签是否匹配
		for _, tag := range img.RepoTags {
			if tag == imageName {
				log.Printf("在ImageList中找到镜像 %s，大小: %d bytes", imageName, img.Size)
				return img.Size, nil
			}
		}

		// 如果没有标签，检查镜像ID
		if img.ID == imageName {
			log.Printf("在ImageList中找到镜像 %s（通过ID），大小: %d bytes", imageName, img.Size)
			return img.Size, nil
		}
	}

	return 0, fmt.Errorf("在镜像列表中未找到镜像: %s", imageName)
}
