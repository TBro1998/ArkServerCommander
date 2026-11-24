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
	"time"

	"ark-server-commander/config"
	"ark-server-commander/utils"

	"github.com/containerd/errdefs"
	"github.com/distribution/reference"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/go-connections/nat"
	dockerspec "github.com/moby/docker-image-spec/specs-go/v1"
	"go.uber.org/zap"
)

// 管理端自我更新：管理端停掉自己的容器后，进程也就没了，没法接着创建新容器，所以由一个单独的辅助容器完成替换。
//
// 管理端（服务层 self_update.go）拉取新镜像后，用自己当前的镜像（一定带 self-update 子命令，不依赖新镜像）启动辅助容器
// SelfUpdateHelperName，挂上 Docker socket 和数据目录，运行 `/app/main self-update <旧容器 ID> <新镜像 ID> /data`。
// 辅助容器（RunSelfUpdateHelper）：
//  1. 把旧容器改名为 <名称>-old，用原来的名称、配置和新镜像创建新容器（这一步失败就改回原名，服务不中断）
//  2. 停止旧容器，启动新容器，等新容器里后端和前端的端口都有响应
//  3. 成功就删除旧容器（不删卷）；失败就删除新容器，把旧容器改回原名再启动
//
// 结果写在数据目录的 .self-update.json，新的（或回滚后的旧的）管理端读取后显示在设置页。
// ARK 服务器和联机房节点是单独的容器，管理端重启不影响它们。

// SelfUpdateHelperName 辅助容器的名称
const SelfUpdateHelperName = "ark-commander-updater"

// selfUpdateRecordFile 更新记录，在数据目录里
const selfUpdateRecordFile = ".self-update.json"

// 更新记录的状态
const (
	SelfUpdateStarted    = "started"     // 已启动辅助容器，还没有结果
	SelfUpdateSucceeded  = "success"     // 新容器已运行，旧容器已删除
	SelfUpdateRolledBack = "rolled_back" // 新容器没有正常运行，已恢复旧容器
	SelfUpdateFailed     = "failed"      // 失败，并且没能恢复旧容器（需要手动处理）
)

const (
	selfUpdateStopTimeout   = 60 // 停止旧容器的超时（秒）
	selfUpdateHealthTimeout = 3 * time.Minute
)

// SelfUpdateRecord 最近一次自我更新的记录
type SelfUpdateRecord struct {
	State       string     `json:"state"`
	FromVersion string     `json:"from_version"`
	ToVersion   string     `json:"to_version,omitempty"` // GitHub 上的最新版本，只用于显示；实际版本以新管理端的版本号为准
	OldImageID  string     `json:"old_image_id"`
	NewImageID  string     `json:"new_image_id"`
	Error       string     `json:"error,omitempty"`
	StartedAt   time.Time  `json:"started_at"`
	FinishedAt  *time.Time `json:"finished_at,omitempty"`
	// ImageCleaned 旧镜像已删除（或不需要删除）
	ImageCleaned bool `json:"image_cleaned,omitempty"`
}

// SelfUpdateTarget 当前管理端容器，自我更新要替换的对象
type SelfUpdateTarget struct {
	ContainerID string
	ImageRef    string // 容器创建时用的镜像名（例如 tbro98/arkservercommander:latest），拉取和创建新容器都用它
	ImageID     string // 容器当前的镜像 ID
	socket      string // Docker socket 在宿主机上的路径
	dataHost    string // 数据目录在宿主机上的路径
}

// 不能自我更新的原因
const (
	SelfUpdateNotContainer = "not_container" // 后端没有运行在容器里
	SelfUpdateCustomImage  = "custom_image"  // 不是发布的镜像的 latest 标签（固定了版本、自己构建的镜像）
	SelfUpdateNoSocket     = "no_socket"     // 没有挂载 /var/run/docker.sock（例如用 DOCKER_HOST 连接）
	SelfUpdateUnavailable  = "unavailable"   // 读取自身容器失败
)

// dockerSocketPath 容器里 Docker socket 的路径
const dockerSocketPath = "/var/run/docker.sock"

// FindSelfUpdateTarget 检查能否自我更新；不能时返回原因（上面的常量）
func FindSelfUpdateTarget() (SelfUpdateTarget, string) {
	if _, err := os.Stat("/.dockerenv"); err != nil {
		return SelfUpdateTarget{}, SelfUpdateNotContainer
	}
	self, err := selfContainer()
	if err != nil || self.Config == nil {
		utils.Warn("读取管理端自身容器失败", zap.Error(err))
		return SelfUpdateTarget{}, SelfUpdateUnavailable
	}
	if !isCommanderLatest(self.Config.Image) {
		return SelfUpdateTarget{}, SelfUpdateCustomImage
	}
	sock, ok := findMount(self.Mounts, dockerSocketPath)
	if !ok || sock.Source == "" {
		return SelfUpdateTarget{}, SelfUpdateNoSocket
	}
	dataHost, ok := mapHostPath(hostMounts, path.Clean(filepath.ToSlash(config.DataDir)))
	if !ok {
		return SelfUpdateTarget{}, SelfUpdateUnavailable
	}
	return SelfUpdateTarget{
		ContainerID: self.ID,
		ImageRef:    self.Config.Image,
		ImageID:     self.Image,
		socket:      sock.Source,
		dataHost:    dataHost,
	}, ""
}

// isCommanderLatest 镜像名是否是发布的管理端镜像的 latest 标签（可以带镜像仓库地址，例如加速器）
func isCommanderLatest(image string) bool {
	named, err := reference.ParseNormalizedNamed(image)
	if err != nil {
		return false
	}
	if _, ok := named.(reference.Digested); ok {
		return false
	}
	if reference.Path(named) != utils.CommanderImageRepo {
		return false
	}
	tagged, ok := reference.TagNameOnly(named).(reference.Tagged)
	return ok && tagged.Tag() == "latest"
}

// SelfUpdateHelperRunning 辅助容器是否正在运行
func SelfUpdateHelperRunning() bool {
	dm, err := GetDockerManager()
	if err != nil {
		return false
	}
	info, err := dm.client.ContainerInspect(dm.ctx, SelfUpdateHelperName)
	return err == nil && info.State != nil && info.State.Running
}

// LaunchSelfUpdateHelper 用管理端当前的镜像启动辅助容器，替换 target 为 newImageID
func LaunchSelfUpdateHelper(target SelfUpdateTarget, newImageID string) error {
	dm, err := GetDockerManager()
	if err != nil {
		return err
	}
	// 上一次没删掉的（AutoRemove 失败或 Docker 重启过）
	if info, err := dm.client.ContainerInspect(dm.ctx, SelfUpdateHelperName); err == nil {
		if info.State != nil && info.State.Running {
			return errors.New("上一次更新还在进行")
		}
		if err := dm.client.ContainerRemove(dm.ctx, info.ID, container.RemoveOptions{Force: true}); err != nil {
			return fmt.Errorf("删除旧的更新容器失败: %w", err)
		}
	} else if !errdefs.IsNotFound(err) {
		return fmt.Errorf("检查更新容器失败: %w", err)
	}

	resp, err := dm.client.ContainerCreate(dm.ctx,
		&container.Config{
			Image:      target.ImageID,
			Entrypoint: []string{"/app/main", "self-update"},
			Cmd:        []string{target.ContainerID, newImageID, "/data"},
			Env:        []string{"LOG_FORMAT=console"},
		},
		&container.HostConfig{
			AutoRemove:  true,
			NetworkMode: "none",
			Binds:       []string{target.socket + ":" + dockerSocketPath, target.dataHost + ":/data"},
			// 镜像声明了卷 /server-files，不挂点东西 Docker 会给辅助容器建一个匿名卷
			Mounts: []mount.Mount{{Type: mount.TypeTmpfs, Target: config.ServerFilesMount}},
		},
		nil, nil, SelfUpdateHelperName)
	if err != nil {
		return fmt.Errorf("创建更新容器失败: %w", err)
	}
	if err := dm.client.ContainerStart(dm.ctx, resp.ID, container.StartOptions{}); err != nil {
		_ = dm.client.ContainerRemove(dm.ctx, resp.ID, container.RemoveOptions{Force: true})
		return fmt.Errorf("启动更新容器失败: %w", err)
	}
	return nil
}

// RunSelfUpdateHelper 辅助容器的入口：把容器 targetID 换成用 newImageID 创建的新容器，结果写入 dataDir 里的更新记录
func RunSelfUpdateHelper(targetID, newImageID, dataDir string) error {
	dm, err := GetDockerManager()
	if err != nil {
		return err
	}
	// 让管理端先把"已开始更新"的响应发出去
	time.Sleep(3 * time.Second)

	state, runErr := dm.replaceContainer(targetID, newImageID)
	if runErr != nil {
		utils.Error("更新管理端失败", zap.String("state", state), zap.Error(runErr))
	} else {
		utils.Info("管理端已更新", zap.String("image", newImageID))
	}

	rec, _ := ReadSelfUpdateRecord(dataDir)
	now := time.Now()
	rec.State = state
	rec.FinishedAt = &now
	rec.Error = ""
	if runErr != nil {
		rec.Error = runErr.Error()
	}
	if rec.NewImageID == "" {
		rec.NewImageID = newImageID
	}
	if err := WriteSelfUpdateRecord(dataDir, rec); err != nil {
		utils.Warn("写入更新记录失败", zap.Error(err))
	}
	return runErr
}

// replaceContainer 返回更新记录的状态，失败时附带原因
func (dm *DockerManager) replaceContainer(targetID, newImageID string) (string, error) {
	old, err := dm.client.ContainerInspect(dm.ctx, targetID)
	if err != nil {
		return SelfUpdateFailed, fmt.Errorf("读取管理端容器失败: %w", err)
	}
	if old.Config == nil || old.HostConfig == nil {
		return SelfUpdateFailed, errors.New("管理端容器的配置不完整")
	}
	name := strings.TrimPrefix(old.Name, "/")

	// 拉取之后到这里，标签可能又被别的操作指向了别的镜像
	refImage, err := dm.client.ImageInspect(dm.ctx, old.Config.Image)
	if err != nil {
		return SelfUpdateRolledBack, fmt.Errorf("读取新镜像失败: %w", err)
	}
	if refImage.ID != newImageID {
		return SelfUpdateRolledBack, fmt.Errorf("镜像 %s 已不是刚下载的版本，请重新更新", old.Config.Image)
	}
	var oldImageConfig *dockerspec.DockerOCIImageConfig
	if oldImage, err := dm.client.ImageInspect(dm.ctx, old.Image); err == nil {
		oldImageConfig = oldImage.Config
	} else {
		utils.Warn("读取旧镜像失败，按容器的完整配置创建新容器", zap.Error(err))
	}

	cfg, hostCfg, primary, extra := buildReplacement(old, oldImageConfig, newImageID)

	backupName := name + "-old"
	if info, err := dm.client.ContainerInspect(dm.ctx, backupName); err == nil {
		if info.State != nil && info.State.Running {
			return SelfUpdateRolledBack, fmt.Errorf("已有一个运行中的容器 %s，请先处理", backupName)
		}
		if err := dm.client.ContainerRemove(dm.ctx, info.ID, container.RemoveOptions{}); err != nil {
			return SelfUpdateRolledBack, fmt.Errorf("删除上次留下的 %s 失败: %w", backupName, err)
		}
	}

	if err := dm.client.ContainerRename(dm.ctx, old.ID, backupName); err != nil {
		return SelfUpdateRolledBack, fmt.Errorf("重命名旧容器失败: %w", err)
	}
	created, err := dm.client.ContainerCreate(dm.ctx, cfg, hostCfg, primary, nil, name)
	if err != nil {
		return dm.rollback(old.ID, "", name, false, fmt.Errorf("创建新容器失败: %w", err))
	}
	for netName, ep := range extra {
		if err := dm.client.NetworkConnect(dm.ctx, netName, created.ID, ep); err != nil {
			return dm.rollback(old.ID, created.ID, name, false, fmt.Errorf("连接网络 %s 失败: %w", netName, err))
		}
	}

	timeout := selfUpdateStopTimeout
	if err := dm.client.ContainerStop(dm.ctx, old.ID, container.StopOptions{Timeout: &timeout}); err != nil {
		return dm.rollback(old.ID, created.ID, name, true, fmt.Errorf("停止旧容器失败: %w", err))
	}
	if err := dm.client.ContainerStart(dm.ctx, created.ID, container.StartOptions{}); err != nil {
		return dm.rollback(old.ID, created.ID, name, true, fmt.Errorf("启动新容器失败: %w", err))
	}
	if err := dm.waitCommanderReady(created.ID, selfUpdateHealthTimeout); err != nil {
		return dm.rollback(old.ID, created.ID, name, true, err)
	}

	if err := dm.client.ContainerRemove(dm.ctx, old.ID, container.RemoveOptions{}); err != nil {
		// 新容器已经在运行，旧容器留着不影响使用
		utils.Warn("删除旧容器失败，可以手动删除", zap.String("container", backupName), zap.Error(err))
	}
	return SelfUpdateSucceeded, nil
}

// rollback 删除新容器，把旧容器改回原名；stopped 为 true 时旧容器已经停了，要重新启动
func (dm *DockerManager) rollback(oldID, newID, name string, stopped bool, cause error) (string, error) {
	var errs []string
	if newID != "" {
		if err := dm.client.ContainerRemove(dm.ctx, newID, container.RemoveOptions{Force: true}); err != nil {
			errs = append(errs, "删除新容器失败: "+err.Error())
		}
	}
	if err := dm.client.ContainerRename(dm.ctx, oldID, name); err != nil {
		errs = append(errs, "旧容器改回原名失败: "+err.Error())
	}
	if stopped {
		if err := dm.client.ContainerStart(dm.ctx, oldID, container.StartOptions{}); err != nil {
			errs = append(errs, "启动旧容器失败: "+err.Error())
		}
	}
	if len(errs) > 0 {
		return SelfUpdateFailed, fmt.Errorf("%w；恢复旧容器时出错：%s", cause, strings.Join(errs, "；"))
	}
	return SelfUpdateRolledBack, cause
}

// waitCommanderReady 等新容器里的后端（8081）和前端（8080）都有响应。只看端口是否有 HTTP 响应，不要求状态码，
// 免得新版本改了页面路径被误判为失败
func (dm *DockerManager) waitCommanderReady(id string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	urls := []string{"http://127.0.0.1:8081/health", "http://127.0.0.1:8080/"}
	for {
		info, err := dm.client.ContainerInspect(dm.ctx, id)
		if err != nil {
			return fmt.Errorf("读取新容器状态失败: %w", err)
		}
		if info.State == nil || !info.State.Running || info.RestartCount > 0 {
			return fmt.Errorf("新容器退出了: %s", dm.logSummary(id))
		}
		ready := true
		for _, u := range urls {
			// wget 退出码 0 是成功，8 是服务器返回了错误状态码（也说明端口在响应）
			res, err := dm.ExecCommand(id, []string{"wget", "-q", "-T", "5", "-O", "/dev/null", u}, nil, 15*time.Second)
			if err != nil || (res.ExitCode != 0 && res.ExitCode != 8) {
				ready = false
				break
			}
		}
		if ready {
			return nil
		}
		if time.Now().After(deadline) {
			return fmt.Errorf("新容器 %s 内没有正常响应: %s", timeout, dm.logSummary(id))
		}
		time.Sleep(3 * time.Second)
	}
}

// logSummary 容器日志的最后几行，放进错误信息
func (dm *DockerManager) logSummary(id string) string {
	lines, err := dm.ContainerLogTail(id, 8)
	if err != nil || len(lines) == 0 {
		return "没有日志"
	}
	return strings.Join(lines, " | ")
}

// buildReplacement 按旧容器的配置生成新容器的配置。
// 容器的配置里混着创建时从旧镜像继承来的值（Env、Cmd、标签、WORKDIR……），原样照抄会盖住新镜像改过的默认值，
// 所以和旧镜像相同的部分去掉，让新镜像的值生效（watchtower 也是这么做的）。oldImage 为 nil 时不做这一步。
// 返回主网络（随容器创建）和其余要在启动前连接的网络
func buildReplacement(old container.InspectResponse, oldImage *dockerspec.DockerOCIImageConfig, newImageID string) (
	*container.Config, *container.HostConfig, *network.NetworkingConfig, map[string]*network.EndpointSettings) {

	cfg := *old.Config
	cfg.Env = slices.Clone(old.Config.Env)
	cfg.Labels = make(map[string]string, len(old.Config.Labels))
	for k, v := range old.Config.Labels {
		cfg.Labels[k] = v
	}
	cfg.ExposedPorts = nil
	if old.Config.ExposedPorts != nil {
		cfg.ExposedPorts = make(nat.PortSet, len(old.Config.ExposedPorts))
	}
	for p := range old.Config.ExposedPorts {
		cfg.ExposedPorts[p] = struct{}{}
	}
	cfg.Volumes = nil
	if old.Config.Volumes != nil {
		cfg.Volumes = make(map[string]struct{}, len(old.Config.Volumes))
	}
	for v := range old.Config.Volumes {
		cfg.Volumes[v] = struct{}{}
	}

	// 没有指定 hostname 时 Docker 用容器 ID 的前 12 位；照抄的话新容器的 hostname 是旧容器的 ID，
	// 管理端靠 hostname 找自己的容器（selfContainer），会找不到
	if len(old.ID) >= 12 && cfg.Hostname == old.ID[:12] {
		cfg.Hostname = ""
	}

	if oldImage != nil {
		img := oldImage.ImageConfig
		cfg.Env = slices.DeleteFunc(cfg.Env, func(e string) bool { return slices.Contains(img.Env, e) })
		if slices.Equal(old.Config.Cmd, img.Cmd) {
			cfg.Cmd = nil
		}
		if slices.Equal(old.Config.Entrypoint, img.Entrypoint) {
			cfg.Entrypoint = nil
		}
		if cfg.WorkingDir == img.WorkingDir {
			cfg.WorkingDir = ""
		}
		if cfg.User == img.User {
			cfg.User = ""
		}
		if cfg.StopSignal == img.StopSignal {
			cfg.StopSignal = ""
		}
		for k, v := range img.Labels {
			if cfg.Labels[k] == v {
				delete(cfg.Labels, k)
			}
		}
		for p := range img.ExposedPorts {
			delete(cfg.ExposedPorts, nat.Port(p))
		}
		for v := range img.Volumes {
			delete(cfg.Volumes, v)
		}
		if hc := oldImage.Healthcheck; hc != nil && cfg.Healthcheck != nil && slices.Equal(cfg.Healthcheck.Test, hc.Test) {
			cfg.Healthcheck = nil
		}
	}
	// docker compose 用这个标签记录创建容器时的镜像 ID，不一致时下次 up 会重建容器
	if _, ok := cfg.Labels["com.docker.compose.image"]; ok {
		cfg.Labels["com.docker.compose.image"] = newImageID
	}

	hostCfg := *old.HostConfig
	hostCfg.Binds = slices.Clone(old.HostConfig.Binds)
	hostCfg.Mounts = slices.Clone(old.HostConfig.Mounts)
	// 镜像声明的卷（/server-files）没有写在 HostConfig 里，Docker 自动建的是匿名卷；
	// 新容器要按卷名接回原来的卷，否则会拿到一个新的空卷
	for _, m := range old.Mounts {
		if m.Type != mount.TypeVolume || m.Name == "" || mountCovered(&hostCfg, m.Destination) {
			continue
		}
		hostCfg.Mounts = append(hostCfg.Mounts, mount.Mount{Type: mount.TypeVolume, Source: m.Name, Target: m.Destination, ReadOnly: !m.RW})
	}

	var primary *network.NetworkingConfig
	extra := map[string]*network.EndpointSettings{}
	mode := string(old.HostConfig.NetworkMode)
	if old.NetworkSettings != nil {
		for netName, ep := range old.NetworkSettings.Networks {
			if ep == nil {
				continue
			}
			settings := &network.EndpointSettings{
				IPAMConfig: ep.IPAMConfig,
				Links:      slices.Clone(ep.Links),
				DriverOpts: ep.DriverOpts,
			}
			// 旧版 Docker 会把容器 ID 加进别名，不能带到新容器
			for _, a := range ep.Aliases {
				if a != old.ID && !(len(old.ID) >= 12 && a == old.ID[:12]) {
					settings.Aliases = append(settings.Aliases, a)
				}
			}
			if netName == mode {
				primary = &network.NetworkingConfig{EndpointsConfig: map[string]*network.EndpointSettings{netName: settings}}
				continue
			}
			// host / none / container: 模式没有别的网络；默认的 bridge 由 NetworkMode 决定，不用再连接
			if container.NetworkMode(mode).IsHost() || container.NetworkMode(mode).IsNone() ||
				container.NetworkMode(mode).IsContainer() || (netName == "bridge" && (mode == "default" || mode == "")) {
				continue
			}
			extra[netName] = settings
		}
	}
	return &cfg, &hostCfg, primary, extra
}

// mountCovered HostConfig 里是否已经有挂载到 dest 的项
func mountCovered(h *container.HostConfig, dest string) bool {
	dest = path.Clean(dest)
	for _, b := range h.Binds {
		parts := strings.Split(b, ":")
		if len(parts) >= 2 && path.Clean(parts[1]) == dest {
			return true
		}
	}
	for _, m := range h.Mounts {
		if path.Clean(m.Target) == dest {
			return true
		}
	}
	for t := range h.Tmpfs {
		if path.Clean(t) == dest {
			return true
		}
	}
	return false
}

func selfUpdateRecordPath(dataDir string) string {
	return filepath.Join(dataDir, selfUpdateRecordFile)
}

// ReadSelfUpdateRecord 读取更新记录，没有时返回 os.ErrNotExist
func ReadSelfUpdateRecord(dataDir string) (SelfUpdateRecord, error) {
	var rec SelfUpdateRecord
	data, err := os.ReadFile(selfUpdateRecordPath(dataDir))
	if err != nil {
		return rec, err
	}
	if err := json.Unmarshal(data, &rec); err != nil {
		return SelfUpdateRecord{}, err
	}
	return rec, nil
}

// WriteSelfUpdateRecord 写入更新记录（先写临时文件再改名）
func WriteSelfUpdateRecord(dataDir string, rec SelfUpdateRecord) error {
	data, err := json.MarshalIndent(rec, "", "  ")
	if err != nil {
		return err
	}
	p := selfUpdateRecordPath(dataDir)
	if err := os.WriteFile(p+".tmp", data, 0o644); err != nil {
		return err
	}
	return os.Rename(p+".tmp", p)
}

// RemoveImage 删除镜像（不强删）。镜像不存在算成功
func (dm *DockerManager) RemoveImage(id string) error {
	_, err := dm.client.ImageRemove(dm.ctx, id, image.RemoveOptions{})
	if err != nil && errdefs.IsNotFound(err) {
		return nil
	}
	return err
}

// SelfImageID 管理端自身容器当前的镜像 ID，不在容器里或读取失败时为空
func SelfImageID() string {
	if _, err := os.Stat("/.dockerenv"); err != nil {
		return ""
	}
	self, err := selfContainer()
	if err != nil {
		return ""
	}
	return self.Image
}
