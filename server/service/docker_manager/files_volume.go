package docker_manager

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"time"

	"ark-server-commander/config"
	"ark-server-commander/utils"

	"github.com/containerd/errdefs"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/volume"
	"go.uber.org/zap"
)

// 单独的服务端文件卷：Windows / macOS 的 Docker Desktop 经文件共享访问宿主机磁盘，steamcmd 下载和预分配 20 GB 的
// 服务端文件非常慢。镜像声明了卷 /server-files（compose 不用改，Docker 自动创建匿名卷，放在 Docker 自己的文件系统里），
// 在 Docker Desktop 上并且数据目录是绑定的宿主机目录时，新服务器的服务端文件放到 /server-files/<id>，其余数据仍在数据目录。
// SERVER_FILES_VOLUME=on / off 可以强制启用或关闭。
//
// 匿名卷没有名字：docker compose up 重建容器时会接回原来的卷，但 down 之后再 up（或手动删掉容器）会换一个新的空卷，
// ARK 服务器的服务端文件卷还绑定在旧卷的目录里。所以把用过的卷名记在数据目录的 .server-files-volume.json，
// 发现换了卷时用临时容器把旧卷的内容复制到新卷；各服务器的服务端文件卷在下次创建容器时按新路径重建（见 ensureBindVolume），
// 没有服务器再用旧卷后删除旧卷

// filesVolumeRecordFile 记录服务端文件卷名的文件，在数据目录里
const filesVolumeRecordFile = ".server-files-volume.json"

// filesVolumeRecord 用过的服务端文件卷
type filesVolumeRecord struct {
	Volume  string   `json:"volume"`            // 当前使用的卷
	Retired []string `json:"retired,omitempty"` // 已迁移走、等各服务器不再使用后删除的旧卷
}

// filesVolumeMu 保护记录文件的读写和旧卷清理
var filesVolumeMu sync.Mutex

// filesVolumeMigrateTimeout 复制旧卷内容的超时时间
const filesVolumeMigrateTimeout = 2 * time.Hour

// wantFilesVolume 是否启用单独的服务端文件卷
func wantFilesVolume(mode string, dockerDesktop, dataIsBind bool) bool {
	switch strings.ToLower(strings.TrimSpace(mode)) {
	case "on", "true", "1":
		return true
	case "off", "false", "0":
		return false
	}
	return dockerDesktop && dataIsBind
}

// initFilesVolume 决定是否启用单独的服务端文件卷，启用时检查卷是否换过并迁移
func initFilesVolume(self container.InspectResponse) {
	m, ok := findMount(self.Mounts, config.ServerFilesMount)
	if !ok {
		return
	}
	dataMount, _ := findMount(self.Mounts, path.Clean(filepath.ToSlash(config.DataDir)))
	if !wantFilesVolume(config.ServerFilesVolume, isDockerDesktop(), isBindMount(dataMount)) {
		return
	}
	if _, ok := mapHostPath(hostMounts, config.ServerFilesMount); !ok {
		utils.Warn("服务端文件卷没有宿主机路径，服务端文件仍放在数据目录")
		return
	}
	filesRoot = config.ServerFilesMount
	utils.Info("服务端文件放在单独的卷里（Docker 自己的文件系统）",
		zap.String("dir", filesRoot), zap.String("volume", m.Name), zap.String("host", hostPath(filesRoot)))

	if m.Name != "" {
		migrateFilesVolume(m.Name, self.Image)
	}
	CleanupRetiredFilesVolumes()
}

// isDockerDesktop Docker 是否为 Docker Desktop（Windows / macOS）
func isDockerDesktop() bool {
	dm, err := GetDockerManager()
	if err != nil {
		return false
	}
	info, err := dm.client.Info(dm.ctx)
	if err != nil {
		utils.Warn("读取 Docker 信息失败", zap.Error(err))
		return false
	}
	return strings.Contains(info.OperatingSystem, "Docker Desktop")
}

func filesVolumeRecordPath() string {
	return filepath.Join(config.DataDir, filesVolumeRecordFile)
}

func readFilesVolumeRecord() filesVolumeRecord {
	var rec filesVolumeRecord
	data, err := os.ReadFile(filesVolumeRecordPath())
	if err != nil {
		return rec
	}
	if err := json.Unmarshal(data, &rec); err != nil {
		utils.Warn("服务端文件卷记录无法解析，重新记录", zap.Error(err))
		return filesVolumeRecord{}
	}
	return rec
}

func writeFilesVolumeRecord(rec filesVolumeRecord) error {
	data, err := json.MarshalIndent(rec, "", "  ")
	if err != nil {
		return err
	}
	tmp := filesVolumeRecordPath() + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, filesVolumeRecordPath())
}

// migrateFilesVolume 当前卷和记录的不同（容器被删除后重新创建，拿到了新的匿名卷）时，把旧卷的内容复制过来
func migrateFilesVolume(current, image string) {
	filesVolumeMu.Lock()
	defer filesVolumeMu.Unlock()

	rec := readFilesVolumeRecord()
	if rec.Volume != "" && rec.Volume != current {
		old := rec.Volume
		exists, err := volumeExists(old)
		switch {
		case err != nil:
			utils.Warn("检查旧的服务端文件卷失败，下次启动再迁移", zap.String("volume", old), zap.Error(err))
			return
		case exists:
			utils.Info("服务端文件卷换了（管理端容器被删除后重新创建），正在从旧卷复制服务端文件，可能需要几分钟",
				zap.String("old", old), zap.String("new", current))
			start := time.Now()
			if err := copyVolume(image, old, current); err != nil {
				utils.Error("复制旧的服务端文件卷失败，下次启动再试；在这之前启动的服务器会重新下载服务端文件",
					zap.String("old", old), zap.Error(err))
				return
			}
			utils.Info("服务端文件已复制到新卷", zap.Duration("elapsed", time.Since(start)))
			if !slices.Contains(rec.Retired, old) {
				rec.Retired = append(rec.Retired, old)
			}
		default:
			utils.Warn("旧的服务端文件卷已不存在，服务器启动时会重新下载服务端文件", zap.String("volume", old))
		}
	}
	rec.Volume = current
	rec.Retired = slices.DeleteFunc(rec.Retired, func(v string) bool { return v == current })
	if err := writeFilesVolumeRecord(rec); err != nil {
		utils.Warn("记录服务端文件卷失败", zap.Error(err))
	}
}

func volumeExists(name string) (bool, error) {
	dm, err := GetDockerManager()
	if err != nil {
		return false, err
	}
	return dm.VolumeExists(name)
}

// copyVolume 用临时容器（管理端自己的镜像）把 src 卷的内容复制到 dst 卷，保留属主、权限和时间
func copyVolume(image, src, dst string) error {
	dm, err := GetDockerManager()
	if err != nil {
		return err
	}
	name := fmt.Sprintf("ase-files-migrate-%d", time.Now().Unix())
	resp, err := dm.client.ContainerCreate(dm.ctx,
		&container.Config{
			Image:      image,
			Entrypoint: []string{"sh", "-c"},
			Cmd:        []string{"cp -a /src/. " + config.ServerFilesMount + "/"},
		},
		// 目标卷挂在镜像声明的卷的位置，否则 Docker 会给这个临时容器再建一个匿名卷，删除容器后留下来
		&container.HostConfig{Binds: []string{src + ":/src:ro", dst + ":" + config.ServerFilesMount}},
		nil, nil, name)
	if err != nil {
		return fmt.Errorf("创建复制用的容器失败: %w", err)
	}
	defer func() {
		if err := dm.client.ContainerRemove(dm.ctx, resp.ID, container.RemoveOptions{Force: true}); err != nil {
			utils.Warn("删除复制用的容器失败", zap.String("container", name), zap.Error(err))
		}
	}()
	if err := dm.client.ContainerStart(dm.ctx, resp.ID, container.StartOptions{}); err != nil {
		return fmt.Errorf("启动复制用的容器失败: %w", err)
	}
	waitCh, errCh := dm.client.ContainerWait(dm.ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case w := <-waitCh:
		if w.StatusCode != 0 {
			lines, _ := dm.ContainerLogTail(name, 20)
			return fmt.Errorf("复制失败（退出码 %d）: %s", w.StatusCode, strings.Join(lines, "; "))
		}
		return nil
	case err := <-errCh:
		return fmt.Errorf("等待复制完成失败: %w", err)
	case <-time.After(filesVolumeMigrateTimeout):
		return errors.New("复制超时")
	}
}

// CleanupRetiredFilesVolumes 删除已迁移走、并且没有服务器的服务端文件卷还绑定在里面的旧卷。
// 运行中的服务器还在用旧卷时保留，等它下次创建容器（按新路径重建服务端文件卷）后再删
func CleanupRetiredFilesVolumes() {
	filesVolumeMu.Lock()
	defer filesVolumeMu.Unlock()
	if filesRoot == "" {
		return
	}
	rec := readFilesVolumeRecord()
	if len(rec.Retired) == 0 {
		return
	}
	dm, err := GetDockerManager()
	if err != nil {
		return
	}
	list, err := dm.client.VolumeList(dm.ctx, volume.ListOptions{Filters: filters.NewArgs(filters.Arg("name", "ase-server-files-"))})
	if err != nil {
		utils.Warn("列出服务端文件卷失败", zap.Error(err))
		return
	}
	var devices []string
	for _, v := range list.Volumes {
		if v != nil {
			devices = append(devices, v.Options["device"])
		}
	}

	var keep []string
	for _, old := range rec.Retired {
		v, err := dm.client.VolumeInspect(dm.ctx, old)
		if errdefs.IsNotFound(err) {
			continue
		}
		if err != nil {
			keep = append(keep, old)
			continue
		}
		if volumeInUse(v.Mountpoint, devices) {
			keep = append(keep, old)
			continue
		}
		if err := dm.client.VolumeRemove(dm.ctx, old, false); err != nil {
			utils.Warn("删除旧的服务端文件卷失败", zap.String("volume", old), zap.Error(err))
			keep = append(keep, old)
			continue
		}
		utils.Info("已删除旧的服务端文件卷", zap.String("volume", old))
	}
	if len(keep) != len(rec.Retired) {
		rec.Retired = keep
		if err := writeFilesVolumeRecord(rec); err != nil {
			utils.Warn("记录服务端文件卷失败", zap.Error(err))
		}
	}
}

// volumeInUse 有服务端文件卷绑定在 mountpoint 下面
func volumeInUse(mountpoint string, devices []string) bool {
	if mountpoint == "" {
		return true
	}
	mp := strings.TrimSuffix(path.Clean(mountpoint), "/")
	for _, d := range devices {
		d = path.Clean(d)
		if d == mp || strings.HasPrefix(d, mp+"/") {
			return true
		}
	}
	return false
}
