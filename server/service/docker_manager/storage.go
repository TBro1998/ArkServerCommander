package docker_manager

import (
	"ark-server-commander/config"
	"ark-server-commander/utils"
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"go.uber.org/zap"
)

// 数据目录布局（后端容器内为 config.DataDir，镜像内是 /data，宿主机上对应 ./data）：
//   servers/<id>/files     服务端文件卷（启用了单独的服务端文件卷时，新服务器改为 /server-files/<id>，见 files_volume.go）
//   servers/<id>/saved     游戏数据卷（存档、配置）
//   servers/<id>/plugins   插件卷
//   clusters/<cluster_id>  集群共享卷
// 这些目录以“绑定到目录的命名卷”挂给 ARK 容器：卷名不变，宿主机上可以直接管理文件。

// hostMount 后端本地的一个目录和它在宿主机上的路径（都用 / 分隔）
type hostMount struct {
	local, host string
}

// hostMounts 后端本地目录到宿主机路径的对应关系。ARK 容器由宿主机的 Docker 创建，绑定路径必须是宿主机路径。
// 后端运行在容器里时来自自身容器的挂载（可以有嵌套挂载）；为空表示后端直接运行在宿主机上
var hostMounts []hostMount

// storageReady InitStorage 是否已完成
var storageReady bool

// filesRoot 单独存放服务端文件的目录（镜像声明的卷 /server-files），为空表示服务端文件都在数据目录里
var filesRoot string

// filesDirCache 服务器 ID -> 服务端文件目录。启用单独的服务端文件卷后，数据目录里已有服务端文件的旧服务器继续用原位置
var filesDirCache sync.Map

// InitStorage 确定数据目录在宿主机上的路径，并决定服务端文件放在哪里，必须在创建卷之前调用。
// 后端直接运行在宿主机上时，本地路径就是宿主机路径；在容器里时检查自身容器的挂载，HOST_DATA_DIR 优先于探测到的数据目录挂载
func InitStorage() error {
	storageReady = false
	hostMounts = nil
	filesRoot = ""
	filesDirCache.Clear()
	dataDir := path.Clean(filepath.ToSlash(config.DataDir))
	_, statErr := os.Stat("/.dockerenv")

	if statErr != nil {
		if config.HostDataDir != "" {
			hostMounts = []hostMount{{local: dataDir, host: path.Clean(config.HostDataDir)}}
		}
		storageReady = true
		return nil
	}

	self, err := selfContainer()
	if err != nil && config.HostDataDir == "" {
		return fmt.Errorf("无法探测数据目录在宿主机上的路径，请设置 HOST_DATA_DIR: %w", err)
	}
	if err == nil {
		hostMounts = containerHostMounts(self.Mounts)
	}
	if config.HostDataDir != "" {
		filtered := hostMounts[:0]
		for _, m := range hostMounts {
			if m.local != dataDir {
				filtered = append(filtered, m)
			}
		}
		hostMounts = append(filtered, hostMount{local: dataDir, host: path.Clean(config.HostDataDir)})
	}
	if _, ok := mapHostPath(hostMounts, dataDir); !ok {
		return fmt.Errorf("容器没有把 %s 挂载到宿主机，请挂载数据目录或设置 HOST_DATA_DIR", dataDir)
	}
	storageReady = true

	if err == nil {
		initFilesVolume(self)
	}
	return nil
}

// selfContainer 后端自身容器的信息
func selfContainer() (container.InspectResponse, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return container.InspectResponse{}, fmt.Errorf("读取容器主机名失败: %w", err)
	}
	dm, err := GetDockerManager()
	if err != nil {
		return container.InspectResponse{}, err
	}
	return dm.client.ContainerInspect(dm.ctx, hostname)
}

// containerHostMounts 容器挂载的目标路径和宿主机路径
func containerHostMounts(points []container.MountPoint) []hostMount {
	var mounts []hostMount
	for _, m := range points {
		if m.Source == "" || m.Destination == "" {
			continue
		}
		mounts = append(mounts, hostMount{local: path.Clean(m.Destination), host: m.Source})
	}
	return mounts
}

// findMount 目标路径为 dest 的挂载
func findMount(points []container.MountPoint, dest string) (container.MountPoint, bool) {
	for _, m := range points {
		if path.Clean(m.Destination) == dest {
			return m, true
		}
	}
	return container.MountPoint{}, false
}

// isBindMount 挂载是否为绑定的宿主机目录（Docker Desktop 上经文件共享访问，读写慢）
func isBindMount(m container.MountPoint) bool {
	return m.Type == mount.TypeBind
}

// hostPath 后端本地路径在宿主机上的路径
func hostPath(local string) string {
	local = path.Clean(filepath.ToSlash(local))
	if len(hostMounts) == 0 {
		return local
	}
	p, _ := mapHostPath(hostMounts, local)
	return p
}

// mapHostPath 按目标路径最长的那个挂载换算；没有挂载覆盖这个路径时返回 false
func mapHostPath(mounts []hostMount, local string) (string, bool) {
	best := -1
	for i, m := range mounts {
		if local == m.local || strings.HasPrefix(local, strings.TrimSuffix(m.local, "/")+"/") {
			if best < 0 || len(m.local) > len(mounts[best].local) {
				best = i
			}
		}
	}
	if best < 0 {
		return "", false
	}
	m := mounts[best]
	return path.Join(m.host, strings.TrimPrefix(local, m.local)), true
}

// volumeSpec 一个卷及其在后端本地的目录
type volumeSpec struct {
	Name  string // 卷名
	Local string // 后端本地的目录
}

// localDir 目录在后端容器内的路径（后端用它创建和删除目录）
func (s volumeSpec) localDir() string {
	return s.Local
}

// hostDir 目录在宿主机上的路径（Docker 卷的 device 参数用它）
func (s volumeSpec) hostDir() string {
	return hostPath(s.Local)
}

// serverVolumeSpecs 服务器的三个卷：服务端文件卷、游戏数据卷、插件卷
func serverVolumeSpecs(serverID uint) []volumeSpec {
	return []volumeSpec{
		{Name: utils.GetServerFilesVolumeName(serverID), Local: ServerFilesDir(serverID)},
		{Name: utils.GetServerVolumeName(serverID), Local: ServerSavedDir(serverID)},
		{Name: utils.GetServerPluginsVolumeName(serverID), Local: ServerPluginsDir(serverID)},
	}
}

// clusterVolumeSpec 集群共享卷
func clusterVolumeSpec(clusterID string) volumeSpec {
	return volumeSpec{Name: utils.GetClusterVolumeName(clusterID), Local: filepath.Join(config.DataDir, "clusters", clusterID)}
}

// legacyServerFilesDir 数据目录里的服务端文件目录
func legacyServerFilesDir(serverID uint) string {
	return filepath.Join(config.DataDir, "servers", fmt.Sprint(serverID), "files")
}

// ServerFilesDir 服务端文件卷在后端本地的目录（对应容器内 /home/steam/arkserver）。
// 默认在数据目录的 servers/<id>/files；启用了单独的服务端文件卷时为 /server-files/<id>，
// 但数据目录里已经有服务端文件（启用之前下载的）的服务器继续用原位置，不用重新下载
func ServerFilesDir(serverID uint) string {
	if filesRoot == "" {
		return legacyServerFilesDir(serverID)
	}
	if v, ok := filesDirCache.Load(serverID); ok {
		return v.(string)
	}
	dir := filepath.Join(filesRoot, fmt.Sprint(serverID))
	if legacy := legacyServerFilesDir(serverID); dirNotEmpty(legacy) {
		dir = legacy
	}
	filesDirCache.Store(serverID, dir)
	return dir
}

// dirNotEmpty 目录存在并且至少有一项内容
func dirNotEmpty(dir string) bool {
	f, err := os.Open(dir)
	if err != nil {
		return false
	}
	defer f.Close()
	names, _ := f.Readdirnames(1)
	return len(names) > 0
}

// ServerSavedDir 游戏数据卷在后端本地的目录（对应容器内 ShooterGame/Saved）
func ServerSavedDir(serverID uint) string {
	return filepath.Join(config.DataDir, "servers", fmt.Sprint(serverID), "saved")
}

// ServerPluginsDir 插件卷在后端本地的目录（对应容器内 ArkApi/Plugins）
func ServerPluginsDir(serverID uint) string {
	return filepath.Join(config.DataDir, "servers", fmt.Sprint(serverID), "plugins")
}

// ServerDisabledPluginsDir 停用的插件存放目录，不挂载进容器，ArkApi 看不到
func ServerDisabledPluginsDir(serverID uint) string {
	return filepath.Join(config.DataDir, "servers", fmt.Sprint(serverID), "plugins-disabled")
}

// ServerDir 服务器在数据目录下的根目录
func ServerDir(serverID uint) string {
	return filepath.Join(config.DataDir, "servers", fmt.Sprint(serverID))
}

// RemoveServerDirs 删除服务器在数据目录下的全部文件（包括单独存放的服务端文件）。
// 删除命名卷不会删除绑定目录里的内容，所以要在 RemoveVolume 之后单独调用
func RemoveServerDirs(serverID uint) error {
	if serverID == 0 {
		return errors.New("无效的服务器 ID")
	}
	if filesRoot != "" {
		if err := os.RemoveAll(filepath.Join(filesRoot, fmt.Sprint(serverID))); err != nil {
			return err
		}
	}
	filesDirCache.Delete(serverID)
	if err := os.RemoveAll(filepath.Join(config.DataDir, "servers", fmt.Sprint(serverID))); err != nil {
		return err
	}
	utils.Debug("已删除服务器数据目录", zap.Uint("server_id", serverID))
	return nil
}
