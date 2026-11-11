package docker_manager

import (
	"ark-server-commander/config"
	"ark-server-commander/utils"
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
)

// 数据目录布局（后端容器内为 config.DataDir，镜像内是 /data，宿主机上对应 ./data）：
//   servers/<id>/files     服务端文件卷
//   servers/<id>/saved     游戏数据卷（存档、配置）
//   servers/<id>/plugins   插件卷
//   clusters/<cluster_id>  集群共享卷
// 这些目录以“绑定到目录的命名卷”挂给 ARK 容器：卷名不变，宿主机上可以直接管理文件。

// hostDataDir 数据目录在宿主机上的路径。ARK 容器由宿主机的 Docker 创建，绑定路径必须是宿主机路径
var hostDataDir string

// InitStorage 确定数据目录在宿主机上的路径，必须在创建卷之前调用。
// 优先使用 HOST_DATA_DIR；否则后端运行在容器里时，检查自身容器的挂载，取目标为数据目录的那一项的 Source；
// 后端直接运行在宿主机上时，数据目录本身就是宿主机路径
func InitStorage() error {
	if config.HostDataDir != "" {
		hostDataDir = path.Clean(config.HostDataDir)
		return nil
	}
	if _, err := os.Stat("/.dockerenv"); err != nil {
		hostDataDir = filepath.ToSlash(config.DataDir)
		return nil
	}
	hostname, err := os.Hostname()
	if err != nil {
		return fmt.Errorf("读取容器主机名失败，请设置 HOST_DATA_DIR: %w", err)
	}
	dm, err := GetDockerManager()
	if err != nil {
		return err
	}
	info, err := dm.client.ContainerInspect(dm.ctx, hostname)
	if err != nil {
		return fmt.Errorf("无法探测数据目录在宿主机上的路径，请设置 HOST_DATA_DIR: %w", err)
	}
	want := path.Clean(filepath.ToSlash(config.DataDir))
	for _, m := range info.Mounts {
		if path.Clean(m.Destination) == want && m.Source != "" {
			hostDataDir = m.Source
			return nil
		}
	}
	return fmt.Errorf("容器没有把 %s 挂载到宿主机，请挂载数据目录或设置 HOST_DATA_DIR", want)
}

// volumeSpec 一个卷及其在数据目录下的相对路径
type volumeSpec struct {
	Name string // 卷名
	Rel  string // 相对数据目录的路径，用 / 分隔
}

// localDir 目录在后端容器内的路径（后端用它创建和删除目录）
func (s volumeSpec) localDir() string {
	return filepath.Join(config.DataDir, filepath.FromSlash(s.Rel))
}

// hostDir 目录在宿主机上的路径（Docker 卷的 device 参数用它）
func (s volumeSpec) hostDir() string {
	return path.Join(hostDataDir, s.Rel)
}

// serverVolumeSpecs 服务器的三个卷：服务端文件卷、游戏数据卷、插件卷
func serverVolumeSpecs(serverID uint) []volumeSpec {
	base := fmt.Sprintf("servers/%d", serverID)
	return []volumeSpec{
		{Name: utils.GetServerFilesVolumeName(serverID), Rel: base + "/files"},
		{Name: utils.GetServerVolumeName(serverID), Rel: base + "/saved"},
		{Name: utils.GetServerPluginsVolumeName(serverID), Rel: base + "/plugins"},
	}
}

// clusterVolumeSpec 集群共享卷
func clusterVolumeSpec(clusterID string) volumeSpec {
	return volumeSpec{Name: utils.GetClusterVolumeName(clusterID), Rel: "clusters/" + clusterID}
}

// RemoveServerDirs 删除服务器在数据目录下的全部文件。
// 删除命名卷不会删除绑定目录里的内容，所以要在 RemoveVolume 之后单独调用
func RemoveServerDirs(serverID uint) error {
	if serverID == 0 {
		return errors.New("无效的服务器 ID")
	}
	return os.RemoveAll(filepath.Join(config.DataDir, "servers", fmt.Sprint(serverID)))
}
