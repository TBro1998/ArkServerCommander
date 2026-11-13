package utils

import (
	"fmt"
	"regexp"
	"strings"
)

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

// ARK 服务器镜像（源码见子模块 ASE-Server-Docker，同一个 dockerfile 的两个构建目标）。
// 不再使用 latest：它停在最后一个自带 ArkApi 的版本，留给旧版管理端
const (
	// ArkServerLinuxImage 原生 Linux 版服务端，不使用 ArkApi 的服务器用它
	ArkServerLinuxImage = "tbro98/ase-server:linux"
	// ArkServerProtonImage 用 GE-Proton 运行 Windows 版服务端，启用 ArkApi 的服务器用它
	ArkServerProtonImage = "tbro98/ase-server:proton"
)

// ArkServerImages 所有 ARK 服务器镜像
var ArkServerImages = []string{ArkServerLinuxImage, ArkServerProtonImage}

// ArkServerImageFor 服务器应使用的镜像：启用 ArkApi 时用 Proton 版（ArkApi 只能注入 Windows 版服务端）
func ArkServerImageFor(arkApi bool) string {
	if arkApi {
		return ArkServerProtonImage
	}
	return ArkServerLinuxImage
}

// AlpineImage 读写卷内文件的临时容器使用的镜像
const AlpineImage = "alpine:latest"

// GetClusterVolumeName 集群共享卷名称：cluster_id 相同的服务器挂载同一个卷，跨服传送的数据存放在这里
func GetClusterVolumeName(clusterID string) string {
	return "ase-cluster-" + clusterID
}

// ClusterIDPattern cluster_id 的允许格式，同时满足 ARK 参数和 Docker 卷名的要求
var ClusterIDPattern = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9_.-]{0,63}$`)

// ModIDsPattern Mod ID 列表格式：逗号分隔的数字
var ModIDsPattern = regexp.MustCompile(`^[0-9]+(,[0-9]+)*$`)

// NormalizeModIDs 去掉 Mod ID 列表里的空白和多余逗号，例如 " 123, 456 ," -> "123,456"
func NormalizeModIDs(ids string) string {
	parts := strings.FieldsFunc(ids, func(r rune) bool { return r == ',' || r == ' ' || r == '\t' || r == '\n' || r == '\r' || r == '，' })
	return strings.Join(parts, ",")
}
