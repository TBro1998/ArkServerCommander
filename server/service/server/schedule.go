package server

import (
	"fmt"
	"strings"
	"sync"
	"time"
	_ "time/tzdata" // 运行镜像里可能没有时区数据库

	"ark-server-commander/config"
	"ark-server-commander/database"
	"ark-server-commander/models"
	"ark-server-commander/service/docker_manager"
	"ark-server-commander/utils"

	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

// cronParser 标准 5 段 Cron 表达式（分 时 日 月 周），也支持 @daily、@every 1h 等写法
var cronParser = cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)

// 倒计时广播的默认内容
const (
	defaultRestartMessage = "Server will restart in {minutes} minute(s). / 服务器将在 {minutes} 分钟后重启"
	defaultStopMessage    = "Server will shut down in {minutes} minute(s). / 服务器将在 {minutes} 分钟后关闭"
)

var (
	locOnce sync.Once
	appLoc  *time.Location
)

// appLocation 定时任务和备份文件名使用的时区，与 ARK 容器一致（CONTAINER_TZ）
func appLocation() *time.Location {
	locOnce.Do(func() {
		loc, err := time.LoadLocation(config.ContainerTZ)
		if err != nil {
			utils.Warn("加载时区失败，使用 UTC", zap.String("tz", config.ContainerTZ), zap.Error(err))
			loc = time.UTC
		}
		appLoc = loc
	})
	return appLoc
}

func nowLocal() time.Time { return time.Now().In(appLocation()) }

// validateSchedule 校验并规范化定时任务
func validateSchedule(req models.ScheduleRequest) (models.ScheduleRequest, error) {
	req.Cron = strings.TrimSpace(req.Cron)
	if _, err := cronParser.Parse(req.Cron); err != nil {
		return req, fmt.Errorf("%w: Cron 表达式无效：%v", ErrInvalidSchedule, err)
	}
	switch req.Action {
	case models.ScheduleRestart, models.ScheduleUpdate, models.ScheduleStop:
	case models.ScheduleBackup, models.ScheduleStart:
		req.WarnMinutes = 0
	case models.ScheduleCommand:
		cmd, err := validateCommand(req.Command)
		if err != nil {
			return req, err
		}
		req.Command = cmd
		req.WarnMinutes = 0
	default:
		return req, fmt.Errorf("%w: 未知的任务类型 %q", ErrInvalidSchedule, req.Action)
	}
	if req.WarnMinutes < 0 || req.WarnMinutes > 60 {
		return req, fmt.Errorf("%w: 提前广播时间需要在 0 到 60 分钟之间", ErrInvalidSchedule)
	}
	req.Message = strings.TrimSpace(strings.NewReplacer("\r", " ", "\n", " ").Replace(req.Message))
	return req, nil
}

// withNextRun 计算下次执行时间（只用于展示）
func withNextRun(sch models.Schedule) models.Schedule {
	if !sch.Enabled {
		return sch
	}
	if spec, err := cronParser.Parse(sch.Cron); err == nil {
		next := spec.Next(nowLocal())
		sch.NextRunAt = &next
	}
	return sch
}

// ListSchedules 列出服务器的定时任务
func (s *ServerService) ListSchedules(userID uint, serverID string) ([]models.Schedule, error) {
	server, err := getUserServer(userID, serverID)
	if err != nil {
		return nil, err
	}
	var list []models.Schedule
	if err := database.DB.Where("server_id = ?", server.ID).Order("id").Find(&list).Error; err != nil {
		return nil, fmt.Errorf("读取定时任务失败: %w", err)
	}
	for i := range list {
		list[i] = withNextRun(list[i])
	}
	return list, nil
}

// CreateSchedule 新建定时任务
func (s *ServerService) CreateSchedule(userID uint, serverID string, req models.ScheduleRequest) (*models.Schedule, error) {
	server, err := getUserServer(userID, serverID)
	if err != nil {
		return nil, err
	}
	req, err = validateSchedule(req)
	if err != nil {
		return nil, err
	}
	sch := models.Schedule{
		ServerID:    server.ID,
		UserID:      userID,
		Action:      req.Action,
		Cron:        req.Cron,
		Enabled:     req.Enabled == nil || *req.Enabled,
		WarnMinutes: req.WarnMinutes,
		Message:     req.Message,
		Command:     req.Command,
		Keep:        req.Keep,
	}
	// Enabled 的默认值是 true，GORM 创建时会把零值 false 换成默认值（并回填到 sch），需要创建后显式写入
	wantEnabled := sch.Enabled
	if err := database.DB.Create(&sch).Error; err != nil {
		return nil, fmt.Errorf("创建定时任务失败: %w", err)
	}
	if !wantEnabled {
		if err := database.DB.Model(&sch).Update("enabled", false).Error; err != nil {
			return nil, fmt.Errorf("创建定时任务失败: %w", err)
		}
		sch.Enabled = false
	}
	result := withNextRun(sch)
	return &result, nil
}

func getSchedule(userID uint, serverID, scheduleID string) (models.Schedule, error) {
	var sch models.Schedule
	server, err := getUserServer(userID, serverID)
	if err != nil {
		return sch, err
	}
	if err := database.DB.Where("id = ? AND server_id = ?", scheduleID, server.ID).First(&sch).Error; err != nil {
		return sch, ErrScheduleNotFound
	}
	return sch, nil
}

// UpdateSchedule 修改定时任务
func (s *ServerService) UpdateSchedule(userID uint, serverID, scheduleID string, req models.ScheduleRequest) (*models.Schedule, error) {
	sch, err := getSchedule(userID, serverID, scheduleID)
	if err != nil {
		return nil, err
	}
	req, err = validateSchedule(req)
	if err != nil {
		return nil, err
	}
	enabled := sch.Enabled
	if req.Enabled != nil {
		enabled = *req.Enabled
	}
	updates := map[string]any{
		"action":       req.Action,
		"cron":         req.Cron,
		"enabled":      enabled,
		"warn_minutes": req.WarnMinutes,
		"message":      req.Message,
		"command":      req.Command,
		"keep":         req.Keep,
	}
	if err := database.DB.Model(&sch).Updates(updates).Error; err != nil {
		return nil, fmt.Errorf("修改定时任务失败: %w", err)
	}
	database.DB.First(&sch, sch.ID)
	result := withNextRun(sch)
	return &result, nil
}

// DeleteSchedule 删除定时任务
func (s *ServerService) DeleteSchedule(userID uint, serverID, scheduleID string) error {
	sch, err := getSchedule(userID, serverID, scheduleID)
	if err != nil {
		return err
	}
	return database.DB.Delete(&sch).Error
}

// RunScheduleNow 立即执行一次定时任务（后台执行，结果写入 last_result）
func (s *ServerService) RunScheduleNow(userID uint, serverID, scheduleID string) error {
	sch, err := getSchedule(userID, serverID, scheduleID)
	if err != nil {
		return err
	}
	go scheduler.run(s, sch)
	return nil
}

// Scheduler 每分钟检查一次到期的定时任务
type Scheduler struct {
	running sync.Map // schedule ID -> struct{}，同一个任务上一次还没执行完时不重复触发
}

var scheduler = &Scheduler{}

// StartScheduler 启动定时任务调度协程
func StartScheduler(s *ServerService) {
	go func() {
		last := nowLocal().Truncate(time.Minute)
		for {
			// 睡到下一分钟的开头
			next := last.Add(time.Minute)
			time.Sleep(time.Until(next))
			now := nowLocal().Truncate(time.Minute)
			// 处理睡眠期间跨过的每一分钟（例如系统休眠后醒来，最多补 5 分钟）
			for t := last.Add(time.Minute); !t.After(now); t = t.Add(time.Minute) {
				if now.Sub(t) <= 5*time.Minute {
					scheduler.tick(s, t)
				}
			}
			last = now
		}
	}()
}

func (sc *Scheduler) tick(s *ServerService, t time.Time) {
	var list []models.Schedule
	if err := database.DB.Where("enabled = ?", true).Find(&list).Error; err != nil {
		utils.Warn("读取定时任务失败", zap.Error(err))
		return
	}
	for _, sch := range list {
		spec, err := cronParser.Parse(sch.Cron)
		if err != nil {
			continue
		}
		// t 是整分钟；上一秒之后的下一次执行时间正好是 t，说明这一分钟该执行
		if spec.Next(t.Add(-time.Second)).Equal(t) {
			go sc.run(s, sch)
		}
	}
}

// run 执行一次定时任务，结果写入 last_run_at / last_result
func (sc *Scheduler) run(s *ServerService, sch models.Schedule) {
	if _, busy := sc.running.LoadOrStore(sch.ID, struct{}{}); busy {
		return
	}
	defer sc.running.Delete(sch.ID)

	err := sc.execute(s, sch)
	result := ""
	if err != nil {
		result = err.Error()
		utils.Warn("定时任务执行失败", zap.Uint("schedule_id", sch.ID), zap.String("action", sch.Action), zap.Error(err))
	} else {
		utils.Info("定时任务执行完成", zap.Uint("schedule_id", sch.ID), zap.String("action", sch.Action))
	}
	now := time.Now()
	database.DB.Model(&models.Schedule{}).Where("id = ?", sch.ID).
		Updates(map[string]any{"last_run_at": &now, "last_result": result})
}

func loadServer(id uint) (models.Server, error) {
	var server models.Server
	if err := database.DB.First(&server, id).Error; err != nil {
		return server, ErrServerNotFound
	}
	return server, nil
}

// countdown 停服前按分钟广播倒计时：剩余分钟数为 warn、5、1 时各广播一次（只广播不超过 warn 的）
func countdown(server models.Server, warn int, message, fallback string) {
	if warn <= 0 || server.Status != "running" {
		return
	}
	if message == "" {
		message = fallback
	}
	points := []int{warn}
	for _, p := range []int{10, 5, 1} {
		if p < warn {
			points = append(points, p)
		}
	}
	for i, minutes := range points {
		broadcast(server, strings.ReplaceAll(message, "{minutes}", fmt.Sprint(minutes)))
		wait := minutes
		if i+1 < len(points) {
			wait = minutes - points[i+1]
		}
		time.Sleep(time.Duration(wait) * time.Minute)
	}
}

func (sc *Scheduler) execute(s *ServerService, sch models.Schedule) error {
	server, err := loadServer(sch.ServerID)
	if err != nil {
		return err
	}
	dockerManager, err := docker_manager.GetDockerManager()
	if err != nil {
		return err
	}

	switch sch.Action {
	case models.ScheduleBackup:
		if server.Status == "starting" || server.Status == "stopping" {
			return ErrServerBusy
		}
		if _, err := createBackup(server, BackupScheduled); err != nil {
			return err
		}
		pruneBackups(server.ID, BackupScheduled, sch.Keep)
		return nil

	case models.ScheduleCommand:
		exec, err := rconForServer(server)
		if err != nil {
			return err
		}
		_, err = exec(sch.Command, time.Minute)
		return err

	case models.ScheduleStart:
		if server.Status != "stopped" {
			return nil
		}
		return s.StartServer(server.UserID, fmt.Sprint(server.ID))

	case models.ScheduleStop:
		if server.Status != "running" {
			return nil
		}
		countdown(server, sch.WarnMinutes, sch.Message, defaultStopMessage)
		return s.StopServer(server.UserID, fmt.Sprint(server.ID))

	case models.ScheduleRestart:
		// 已停止的服务器不在定时重启时启动，避免违背用户手动停服的意图
		if server.Status != "running" {
			return nil
		}
		countdown(server, sch.WarnMinutes, sch.Message, defaultRestartMessage)
		return s.restartWith(server, dockerManager, nil)

	case models.ScheduleUpdate:
		wasRunning := server.Status == "running"
		if server.Status == "starting" || server.Status == "stopping" {
			return ErrServerBusy
		}
		countdown(server, sch.WarnMinutes, sch.Message, defaultRestartMessage)
		markers := func() error {
			return dockerManager.RemoveInstallMarkers(server.ID, true, server.GameModIds != "")
		}
		if !wasRunning {
			// 已停止的服务器只清除标记，下次启动时更新
			return markers()
		}
		return s.restartWith(server, dockerManager, markers)
	}
	return fmt.Errorf("%w: 未知的任务类型 %q", ErrInvalidSchedule, sch.Action)
}
