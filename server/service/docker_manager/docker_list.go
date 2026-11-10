package docker_manager

import (
	"fmt"
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
)

// ContainerSummary 服务器容器的概要信息
type ContainerSummary struct {
	Status  string // running / stopped / starting / unknown，与 GetContainerStatus 的取值一致
	ImageID string // 创建容器时使用的镜像 ID
}

// ListServerContainers 一次列出所有服务器容器（ase-server-<id>），key 是容器名。
// 比对每台服务器分别 inspect 少很多 Docker 请求，服务器多时列表接口明显更快
func (dm *DockerManager) ListServerContainers() (map[string]ContainerSummary, error) {
	list, err := dm.client.ContainerList(dm.ctx, container.ListOptions{
		All:     true,
		Filters: filters.NewArgs(filters.Arg("name", "^/ase-server-[0-9]+$")),
	})
	if err != nil {
		return nil, fmt.Errorf("获取容器列表失败: %v", err)
	}

	result := make(map[string]ContainerSummary, len(list))
	for _, c := range list {
		summary := ContainerSummary{Status: summaryStatus(c.State), ImageID: c.ImageID}
		for _, name := range c.Names {
			result[strings.TrimPrefix(name, "/")] = summary
		}
	}
	return result, nil
}

// summaryStatus 把 Docker 的容器状态转换成本系统的状态，与 GetContainerStatus 保持一致
func summaryStatus(state string) string {
	switch state {
	case "running", "paused":
		return "running"
	case "exited", "created":
		return "stopped"
	case "restarting":
		return "starting"
	default:
		return "unknown"
	}
}

// GetImageID 返回本地镜像的 ID
func (dm *DockerManager) GetImageID(imageName string) (string, error) {
	info, err := dm.client.ImageInspect(dm.ctx, imageName)
	if err != nil {
		return "", fmt.Errorf("获取Docker镜像信息失败: %v", err)
	}
	return info.ID, nil
}
