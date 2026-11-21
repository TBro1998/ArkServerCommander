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

// LinkContainerName P2P 联机房的 EasyTier 节点容器名称（每个管理端一个房间、一个节点）。
// 不匹配 ase-server-<id>，不会被当成 ARK 服务器容器（ListServerContainers 按名称过滤）。
// 旧版每台服务器一个节点（ase-link-<id>），升级后由联机房后台删除
const LinkContainerName = "ase-link-room"

// EasyTierImage P2P 联机房节点使用的 EasyTier 镜像（入口是 easytier-core）
const EasyTierImage = "easytier/easytier:v2.6.4"

// ASE 服务端在 Steam 上的 AppID
const ArkServerAppID = "376030"

// 游戏版本（Steam 分支）：主流的两个版本是正式版（public）和海洋 DLC 更新前的 preaquatica 测试分支。
// 两个版本的服务端文件不同，复制服务端文件时只能在同一版本的服务器之间复制
const (
	GameBranchPublic      = ""            // 正式版（Steam 默认分支）
	GameBranchPreAquatica = "preaquatica" // 海洋 DLC 更新前的版本（-beta preaquatica）
)

// GameBranches 可选的游戏版本
var GameBranches = []string{GameBranchPublic, GameBranchPreAquatica}

// ValidGameBranch 是否是可选的游戏版本
func ValidGameBranch(branch string) bool {
	for _, b := range GameBranches {
		if b == branch {
			return true
		}
	}
	return false
}

// ArkAppUpdateFor 容器环境变量 APP_UPDATE 的值，入口脚本按空格拆开传给 steamcmd 的 app_update。
// 正式版和镜像里的默认值相同（只有 AppID）
func ArkAppUpdateFor(branch string) string {
	if branch == GameBranchPublic {
		return ArkServerAppID
	}
	return ArkServerAppID + " -beta " + branch
}
