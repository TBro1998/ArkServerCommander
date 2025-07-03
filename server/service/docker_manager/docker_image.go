package docker_manager

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
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
	Exists   bool   `json:"exists"`   // 镜像是否存在
	Pulling  bool   `json:"pulling"`  // 是否正在拉取
	Progress string `json:"progress"` // 下载进度信息
	Ready    bool   `json:"ready"`    // 是否准备就绪
	Error    string `json:"error"`    // 错误信息
	Percent  int    `json:"percent"`  // 拉取进度百分比
}

// 全局镜像状态管理
var (
	imagePullStatus   = make(map[string]bool)   // 记录镜像是否正在拉取
	imagePullProgress = make(map[string]string) // 记录镜像拉取进度
	imagePullMutex    sync.RWMutex
)

// PullImage 拉取Docker镜像
// imageName: 镜像名称
// 返回: 错误信息
func (dm *DockerManager) PullImage(imageName string) error {
	fmt.Printf("正在拉取Docker镜像: %s\n", imageName)

	// 拉取镜像
	reader, err := dm.client.ImagePull(dm.ctx, imageName, image.PullOptions{})
	if err != nil {
		return fmt.Errorf("拉取Docker镜像失败: %v", err)
	}
	defer reader.Close()

	// 读取拉取进度（可选，用于显示进度）
	_, err = io.Copy(io.Discard, reader)
	if err != nil {
		return fmt.Errorf("读取镜像拉取进度失败: %v", err)
	}

	fmt.Printf("Docker镜像拉取成功: %s\n", imageName)
	return nil
}

// PullImageWithProgress 拉取Docker镜像并显示进度
// imageName: 镜像名称
// 返回: 错误信息
func (dm *DockerManager) PullImageWithProgress(imageName string) error {
	fmt.Printf("正在拉取Docker镜像: %s\n", imageName)

	// 设置拉取状态
	imagePullMutex.Lock()
	imagePullStatus[imageName] = true
	imagePullProgress[imageName] = "开始拉取镜像..."
	imagePullMutex.Unlock()

	// 确保在函数结束时清理状态
	defer func() {
		imagePullMutex.Lock()
		imagePullStatus[imageName] = false
		delete(imagePullProgress, imageName)
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
						imagePullProgress[imageName] = progressInfo.Status
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
	imagePullMutex.Unlock()

	return nil
}

// ImageExists 检查Docker镜像是否存在
// imageName: 镜像名称
// 返回: 是否存在和错误信息
func (dm *DockerManager) ImageExists(imageName string) (bool, error) {
	// 尝试获取镜像信息
	_, _, err := dm.client.ImageInspectWithRaw(dm.ctx, imageName)
	if err != nil {
		if client.IsErrNotFound(err) {
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
		Exists:   false,
		Pulling:  false,
		Progress: "",
		Ready:    false,
		Error:    "",
		Percent:  0,
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
		return status
	}

	// 检查是否正在拉取中
	imagePullMutex.RLock()
	status.Pulling = imagePullStatus[imageName]
	if status.Pulling {
		status.Progress = imagePullProgress[imageName]
		// 根据进度信息估算百分比
		if strings.Contains(status.Progress, "Downloading") {
			status.Percent = 50
		} else if strings.Contains(status.Progress, "Extracting") {
			status.Percent = 75
		} else if strings.Contains(status.Progress, "完成") {
			status.Percent = 100
		} else {
			status.Percent = 25
		}
	}
	imagePullMutex.RUnlock()

	if !status.Pulling {
		status.Progress = "镜像不存在，需要下载"
		status.Percent = 0
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
