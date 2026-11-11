package server

import (
	"errors"

	"ark-server-commander/service/docker_manager"
)

// 服务层哨兵错误，controller 通过 errors.Is 判断并映射 HTTP 状态码
var (
	ErrInvalidServerID  = errors.New("无效的服务器ID")
	ErrServerNotFound   = errors.New("服务器不存在")
	ErrIdentifierExists = errors.New("服务器标识已存在")
	ErrServerRunning    = errors.New("无法删除正在运行的服务器，请先停止服务器")
	ErrAlreadyRunning   = errors.New("服务器已在运行中")
	ErrAlreadyStarting  = errors.New("服务器正在启动中")
	ErrAlreadyStopped   = errors.New("服务器已经停止")
	ErrAlreadyStopping  = errors.New("服务器正在停止中")
	ErrImagePulling     = errors.New("镜像正在下载中，请等待完成")
	ErrPortConflict     = errors.New("端口配置冲突")
	ErrInvalidPort      = errors.New("端口无效")
	ErrServerBusy       = errors.New("服务器正在启动或停止中，请稍后再试")
	ErrImageNotAllowed  = errors.New("不允许操作该镜像")
	ErrInvalidClusterID = errors.New("集群ID只能包含字母、数字、下划线、点和连字符，且以字母或数字开头，最长64个字符")
	ErrInvalidModIDs    = errors.New("Mod ID 只能是逗号分隔的数字")
	ErrServerNotRunning = errors.New("服务器未运行")
	ErrServerNotStopped = errors.New("请先停止服务器")
	ErrInvalidCommand   = errors.New("命令不能为空、不能包含换行，且不超过1000个字符")
	ErrBackupNotFound   = errors.New("备份不存在")
	ErrInvalidBackup    = docker_manager.ErrInvalidBackup
	ErrBackupBusy       = errors.New("该服务器正在备份或恢复，请稍后再试")
	ErrScheduleNotFound = errors.New("定时任务不存在")
	ErrInvalidSchedule  = errors.New("定时任务配置无效")
	ErrRconUnavailable  = errors.New("RCON 端口没有响应，服务器可能还在下载、安装 Mod 或加载地图")
	ErrRconAuthFailed   = errors.New("RCON 认证失败，管理员密码与游戏内的不一致，重启服务器后会重新写入")
)

// errorCodes 错误码，前端据此显示对应语言的提示；中文界面直接显示 error 字段（包含具体细节）
var errorCodes = []struct {
	err  error
	code string
}{
	{ErrInvalidServerID, "invalid_server_id"},
	{ErrServerNotFound, "server_not_found"},
	{ErrIdentifierExists, "identifier_exists"},
	{ErrServerRunning, "server_running"},
	{ErrAlreadyRunning, "already_running"},
	{ErrAlreadyStarting, "already_starting"},
	{ErrAlreadyStopped, "already_stopped"},
	{ErrAlreadyStopping, "already_stopping"},
	{ErrImagePulling, "image_pulling"},
	{ErrPortConflict, "port_conflict"},
	{ErrInvalidPort, "invalid_port"},
	{ErrServerBusy, "server_busy"},
	{ErrImageNotAllowed, "image_not_allowed"},
	{ErrInvalidClusterID, "invalid_cluster_id"},
	{ErrInvalidModIDs, "invalid_mod_ids"},
	{ErrServerNotRunning, "server_not_running"},
	{ErrServerNotStopped, "server_not_stopped"},
	{ErrInvalidCommand, "invalid_command"},
	{ErrBackupNotFound, "backup_not_found"},
	{ErrInvalidBackup, "invalid_backup"},
	{ErrBackupBusy, "backup_busy"},
	{ErrScheduleNotFound, "schedule_not_found"},
	{ErrInvalidSchedule, "invalid_schedule"},
	{ErrRconUnavailable, "rcon_unavailable"},
	{ErrRconAuthFailed, "rcon_auth_failed"},
}

// ErrorCode 返回错误对应的错误码，未知错误返回空字符串
func ErrorCode(err error) string {
	for _, e := range errorCodes {
		if errors.Is(err, e.err) {
			return e.code
		}
	}
	return ""
}
