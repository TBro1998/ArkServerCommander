package utils

import "fmt"

// GetServerContainerName 获取服务器容器名称
// serverID: 服务器ID
// 返回: 容器名称
func GetServerContainerName(serverID uint) string {
	return fmt.Sprintf("ase-server-%d", serverID)
}

// GetServerVolumeName 获取服务器卷名称
// serverID: 服务器ID
// 返回: 卷名称
func GetServerVolumeName(serverID uint) string {
	return fmt.Sprintf("ase-server-%d", serverID)
}

// GetServerPluginsVolumeName 获取服务器插件卷名称
// serverID: 服务器ID
// 返回: 插件卷名称
func GetServerPluginsVolumeName(serverID uint) string {
	return fmt.Sprintf("ase-server-plugins-%d", serverID)
}

// GetServerFilesVolumeName 获取服务端文件卷名称
// 镜像不再包含服务端文件，容器首次启动时由 steamcmd 下载到该卷，重建容器不会重新下载
// serverID: 服务器ID
// 返回: 服务端文件卷名称
func GetServerFilesVolumeName(serverID uint) string {
	return fmt.Sprintf("ase-server-files-%d", serverID)
}

// ArkServerImage ARK服务器镜像（源码见子模块 ASE-Server-Docker）
const ArkServerImage = "tbro98/ase-server:latest"
