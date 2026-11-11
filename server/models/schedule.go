package models

import "time"

// 定时任务类型
const (
	ScheduleRestart = "restart" // 重启（停服前按 WarnMinutes 广播倒计时）
	ScheduleBackup  = "backup"  // 备份存档，只保留最近 Keep 份定时备份
	ScheduleUpdate  = "update"  // 更新游戏本体和 Mod 后重启
	ScheduleCommand = "command" // 执行一条 RCON 命令
	ScheduleStart   = "start"   // 启动
	ScheduleStop    = "stop"    // 停止（停服前广播倒计时）
)

// Schedule 服务器定时任务，按 Cron 表达式（5 段，时区为 CONTAINER_TZ）触发
type Schedule struct {
	ID          uint       `json:"id" gorm:"primarykey"`
	ServerID    uint       `json:"server_id" gorm:"index;not null"`
	UserID      uint       `json:"-" gorm:"not null"`
	Action      string     `json:"action" gorm:"not null"`
	Cron        string     `json:"cron" gorm:"not null"`
	Enabled     bool       `json:"enabled" gorm:"default:true"`
	WarnMinutes int        `json:"warn_minutes" gorm:"default:0"`
	Message     string     `json:"message" gorm:"default:''"` // 倒计时广播内容，{minutes} 会被替换为剩余分钟数
	Command     string     `json:"command" gorm:"default:''"` // command 类型要执行的 RCON 命令
	Keep        int        `json:"keep" gorm:"default:0"`     // backup 类型保留的份数，0 表示不清理
	LastRunAt   *time.Time `json:"last_run_at"`
	LastResult  string     `json:"last_result" gorm:"default:''"` // 最近一次执行结果，成功为空
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	NextRunAt   *time.Time `json:"next_run_at,omitempty" gorm:"-"`
}

// ScheduleRequest 创建或修改定时任务
type ScheduleRequest struct {
	Action      string `json:"action" binding:"required"`
	Cron        string `json:"cron" binding:"required"`
	Enabled     *bool  `json:"enabled"`
	WarnMinutes int    `json:"warn_minutes" binding:"min=0,max=60"`
	Message     string `json:"message"`
	Command     string `json:"command"`
	Keep        int    `json:"keep" binding:"min=0,max=1000"`
}
