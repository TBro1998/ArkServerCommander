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
