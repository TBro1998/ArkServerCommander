package server

import (
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	"ark-server-commander/config"
	"ark-server-commander/database"
	"ark-server-commander/models"
	"ark-server-commander/service/docker_manager"
	"ark-server-commander/utils"

	"go.uber.org/zap"
)

// 管理端自我更新的服务层：检查条件、拉取新镜像、备份数据库、启动辅助容器。替换容器的过程见 docker_manager/self_update.go

// SelfUpdateStatus 自我更新的状态，随 GET /system/version 返回
type SelfUpdateStatus struct {
	Supported bool   `json:"supported"`
	Reason    string `json:"reason,omitempty"` // 不支持的原因：dev / not_container / custom_image / no_socket / unavailable
	// State idle / pulling（下载新镜像）/ restarting（辅助容器正在替换容器）/ error（下载阶段失败，替换阶段的结果看 last）
	State     string                           `json:"state"`
	Progress  int                              `json:"progress,omitempty"` // pulling 时的下载进度（0-100，估算）
	ErrorCode string                           `json:"error_code,omitempty"`
	Error     string                           `json:"error,omitempty"`
	Last      *docker_manager.SelfUpdateRecord `json:"last,omitempty"` // 最近一次更新的记录，state 为 interrupted 表示辅助容器没有写入结果就退出了
}

var selfUpdate struct {
	sync.Mutex
	state     string // "" / pulling / restarting / error
	imageRef  string
	errorCode string
	err       string
}

// selfUpdateBackupSuffix 更新前数据库备份的后缀，放在数据库旁边，每次更新覆盖
const selfUpdateBackupSuffix = ".before-update"

// GetSelfUpdateStatus 自我更新的状态
func (s *ServerService) GetSelfUpdateStatus() SelfUpdateStatus {
	st := SelfUpdateStatus{State: "idle"}
	if utils.Version == "dev" {
		st.Reason = "dev"
	} else if _, reason := docker_manager.FindSelfUpdateTarget(); reason != "" {
		st.Reason = reason
	} else {
		st.Supported = true
	}

	rec, recErr := docker_manager.ReadSelfUpdateRecord(config.DataDir)
	helperRunning := false
	if recErr == nil && rec.State == docker_manager.SelfUpdateStarted {
		helperRunning = docker_manager.SelfUpdateHelperRunning()
	}

	selfUpdate.Lock()
	state, ref, code, msg := selfUpdate.state, selfUpdate.imageRef, selfUpdate.errorCode, selfUpdate.err
	// 辅助容器在旧容器停止前失败（例如创建新容器失败）时，旧管理端还活着，记录里已经有结果
	if state == "restarting" && recErr == nil && rec.State != docker_manager.SelfUpdateStarted {
		selfUpdate.state = ""
		state = ""
	}
	selfUpdate.Unlock()

	switch {
	case state == "pulling":
		st.State = "pulling"
		if dm, err := docker_manager.GetDockerManager(); err == nil {
			st.Progress = pullPercent(dm.GetImageStatus(ref))
		}
	case state == "restarting" || helperRunning:
		st.State = "restarting"
	case state == "error":
		st.State = "error"
		st.ErrorCode, st.Error = code, msg
	}

	if recErr == nil {
		last := rec
		if last.State == docker_manager.SelfUpdateStarted && st.State != "restarting" {
			last.State = "interrupted"
		}
		st.Last = &last
	}
	return st
}

// pullPercent 按各层已下载的字节估算进度。层的大小要等开始下载才知道，所以只是估算
func pullPercent(status *docker_manager.ImageStatus) int {
	if status == nil || len(status.Layers) == 0 {
		return 0
	}
	var done, total int64
	for _, l := range status.Layers {
		switch l.Status {
		case "exists", "complete":
			continue
		}
		total += l.Size
		if l.Status == "downloading" {
			done += l.Progress
		} else if l.Status != "pending" {
			done += l.Size
		}
	}
	if total <= 0 {
		return 0
	}
	p := int(done * 100 / total)
	return min(max(p, 0), 99)
}

// StartSelfUpdate 开始更新管理端：检查条件后在后台拉取新镜像、启动辅助容器，立即返回。
// toVersion 是界面上看到的最新版本，只记录在更新记录里用于显示
func (s *ServerService) StartSelfUpdate(toVersion string) error {
	if utils.Version == "dev" {
		return fmt.Errorf("%w：开发版本不能自动更新", ErrSelfUpdateUnsupported)
	}
	target, reason := docker_manager.FindSelfUpdateTarget()
	if reason != "" {
		return fmt.Errorf("%w：%s", ErrSelfUpdateUnsupported, selfUpdateReasonText(reason))
	}
	if docker_manager.SelfUpdateHelperRunning() {
		return ErrSelfUpdateRunning
	}
	if err := checkSelfUpdateIdle(); err != nil {
		return err
	}

	selfUpdate.Lock()
	if selfUpdate.state == "pulling" || selfUpdate.state == "restarting" {
		selfUpdate.Unlock()
		return ErrSelfUpdateRunning
	}
	if !docker_manager.ClaimPull(target.ImageRef) {
		selfUpdate.Unlock()
		return ErrImagePulling
	}
	selfUpdate.state, selfUpdate.imageRef, selfUpdate.errorCode, selfUpdate.err = "pulling", target.ImageRef, "", ""
	selfUpdate.Unlock()

	go runSelfUpdate(target, toVersion)
	return nil
}

func runSelfUpdate(target docker_manager.SelfUpdateTarget, toVersion string) {
	fail := func(code string, err error) {
		utils.Error("更新管理端失败", zap.String("code", code), zap.Error(err))
		selfUpdate.Lock()
		selfUpdate.state, selfUpdate.errorCode, selfUpdate.err = "error", code, err.Error()
		selfUpdate.Unlock()
	}

	dm, err := docker_manager.GetDockerManager()
	if err != nil {
		docker_manager.ReleasePull(target.ImageRef, err)
		fail("self_update_failed", err)
		return
	}
	utils.Info("开始更新管理端，下载新镜像", zap.String("image", target.ImageRef))
	if err := dm.PullClaimedImage(target.ImageRef); err != nil {
		fail("self_update_pull_failed", err)
		return
	}
	newID, err := dm.GetImageID(target.ImageRef)
	if err != nil {
		fail("self_update_failed", err)
		return
	}
	if newID == target.ImageID {
		fail("self_update_latest", errors.New("镜像仓库里还没有新版本的镜像（刚发布时镜像可能还在构建，或者镜像加速器还没同步），请稍后再试"))
		return
	}
	// 拉取期间可能有服务器开始启停
	if err := checkSelfUpdateIdle(); err != nil {
		fail("self_update_busy", err)
		return
	}
	if err := backupDatabaseForUpdate(); err != nil {
		fail("self_update_failed", fmt.Errorf("更新前备份数据库失败: %w", err))
		return
	}

	rec := docker_manager.SelfUpdateRecord{
		State:       docker_manager.SelfUpdateStarted,
		FromVersion: utils.Version,
		ToVersion:   toVersion,
		OldImageID:  target.ImageID,
		NewImageID:  newID,
		StartedAt:   time.Now(),
	}
	if err := docker_manager.WriteSelfUpdateRecord(config.DataDir, rec); err != nil {
		fail("self_update_failed", fmt.Errorf("写入更新记录失败: %w", err))
		return
	}
	selfUpdate.Lock()
	selfUpdate.state = "restarting"
	selfUpdate.Unlock()
	if err := docker_manager.LaunchSelfUpdateHelper(target, newID); err != nil {
		now := time.Now()
		rec.State, rec.Error, rec.FinishedAt = docker_manager.SelfUpdateRolledBack, err.Error(), &now
		_ = docker_manager.WriteSelfUpdateRecord(config.DataDir, rec)
		fail("self_update_failed", err)
		return
	}
	utils.Info("已启动更新容器，管理端即将重启", zap.String("new_image", newID))
}

// checkSelfUpdateIdle 有服务器正在启停、正在备份或恢复、正在下载镜像时不更新（管理端重启会打断这些操作）
func checkSelfUpdateIdle() error {
	var n int64
	if err := database.DB.Model(&models.Server{}).Where("status IN ?", []string{"starting", "stopping"}).Count(&n).Error; err != nil {
		return err
	}
	if n > 0 {
		return fmt.Errorf("%w：有服务器正在启动或停止", ErrSelfUpdateBusy)
	}
	if anyBackupRunning() {
		return fmt.Errorf("%w：有服务器正在备份或恢复存档", ErrSelfUpdateBusy)
	}
	for _, img := range managedImages {
		if docker_manager.IsImagePulling(img) {
			return fmt.Errorf("%w：正在下载镜像 %s", ErrSelfUpdateBusy, img)
		}
	}
	return nil
}

// anyBackupRunning 是否有服务器正在备份或恢复
func anyBackupRunning() bool {
	busy := false
	backupLocks.Range(func(_, v any) bool {
		m := v.(*sync.Mutex)
		if !m.TryLock() {
			busy = true
			return false
		}
		m.Unlock()
		return true
	})
	return busy
}

// backupDatabaseForUpdate 更新前把数据库复制一份（新版本可能迁移表结构，回滚后需要时可以手动换回）
func backupDatabaseForUpdate() error {
	dst := config.DBPath + selfUpdateBackupSuffix
	if err := os.Remove(dst); err != nil && !os.IsNotExist(err) {
		return err
	}
	return database.DB.Exec("VACUUM INTO ?", dst).Error
}

func selfUpdateReasonText(reason string) string {
	switch reason {
	case docker_manager.SelfUpdateNotContainer:
		return "管理端没有运行在 Docker 容器里"
	case docker_manager.SelfUpdateCustomImage:
		return "管理端不是用 " + utils.CommanderImageRepo + ":latest 镜像运行的（固定了版本或自己构建的镜像），请手动更新"
	case docker_manager.SelfUpdateNoSocket:
		return "管理端容器没有挂载 /var/run/docker.sock"
	default:
		return "读取管理端容器信息失败"
	}
}

// StartSelfUpdateCleanup 启动时处理上一次自我更新：更新成功后删除旧镜像。
// 辅助容器还在用旧镜像（它要等新管理端启动完成才退出），所以隔一段时间重试几次
func StartSelfUpdateCleanup() {
	rec, err := docker_manager.ReadSelfUpdateRecord(config.DataDir)
	if err != nil || rec.ImageCleaned || rec.OldImageID == "" {
		return
	}
	if rec.State != docker_manager.SelfUpdateSucceeded && rec.State != docker_manager.SelfUpdateStarted {
		return
	}
	go func() {
		for i := 0; i < 20; i++ {
			time.Sleep(30 * time.Second)
			rec, err := docker_manager.ReadSelfUpdateRecord(config.DataDir)
			if err != nil || rec.ImageCleaned {
				return
			}
			if rec.State == docker_manager.SelfUpdateStarted {
				continue // 辅助容器还没写入结果
			}
			if rec.State != docker_manager.SelfUpdateSucceeded {
				return
			}
			if self := docker_manager.SelfImageID(); self == "" || self == rec.OldImageID {
				return
			}
			dm, err := docker_manager.GetDockerManager()
			if err != nil {
				return
			}
			if err := dm.RemoveImage(rec.OldImageID); err != nil {
				continue
			}
			utils.Info("已删除更新前的旧镜像", zap.String("image", rec.OldImageID))
			rec.ImageCleaned = true
			if err := docker_manager.WriteSelfUpdateRecord(config.DataDir, rec); err != nil {
				utils.Warn("写入更新记录失败", zap.Error(err))
			}
			return
		}
		utils.Info("更新前的旧镜像仍被使用，没有删除")
	}()
}
