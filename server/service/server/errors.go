package server

import "errors"

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
)
