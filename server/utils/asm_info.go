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

// ARK 服务器镜像。ASE 的源码见子模块 ASE-Server-Docker（同一个 dockerfile 的两个构建目标），
// ASA 的见子模块 ASA-Server-Docker。ASE 不再使用 latest：它停在最后一个自带 ArkApi 的版本，留给旧版管理端
const (
	// ArkServerLinuxImage ASE 原生 Linux 版服务端，不使用 ArkApi 的 ASE 服务器用它
	ArkServerLinuxImage = "tbro98/ase-server:linux"
	// ArkServerProtonImage ASE 用 GE-Proton 运行 Windows 版服务端，启用 ArkApi 的 ASE 服务器用它
	ArkServerProtonImage = "tbro98/ase-server:proton"
	// AsaServerProtonImage ASA 只有 Windows 版服务端，所有 ASA 服务器都用 GE-Proton 运行；
	// 是否加载 AsaApi 由容器环境变量 ASAAPI 决定（AsaApiEnv）
	AsaServerProtonImage = "tbro98/asa-server:proton"
)

// ArkServerImages 所有 ARK 服务器镜像
var ArkServerImages = []string{ArkServerLinuxImage, ArkServerProtonImage, AsaServerProtonImage}

// RetiredArkServerImages 不再使用的旧镜像：没有容器在用时删掉，免得一直占着磁盘。
// tbro98/asa-server:wine 是以前不启用 AsaApi 的 ASA 服务器用的 Wine 版，和 Proton 版只差几十 MB，已经去掉
var RetiredArkServerImages = []string{"tbro98/asa-server:wine"}

// ArkServerImageFor 服务器应使用的镜像：ASE 启用 ArkApi 时用 Proton 版（插件只能注入 Windows 版服务端），
// 否则用原生 Linux 版；ASA 只有一个镜像
func ArkServerImageFor(game string, arkApi bool) string {
	if NormalizeGame(game) == GameASA {
		return AsaServerProtonImage
	}
	if arkApi {
		return ArkServerProtonImage
	}
	return ArkServerLinuxImage
}

// AsaApiEnv ASA 容器的环境变量 ASAAPI：启用 AsaApi 时为 "1"，否则为 "0"（入口脚本不加载文件卷里留着的 AsaApi）。
// ASE 不用这个变量，返回空。值不同时重建容器
func AsaApiEnv(game string, arkApi bool) string {
	if NormalizeGame(game) != GameASA {
		return ""
	}
	if arkApi {
		return "1"
	}
	return "0"
}

// ImageGame 镜像属于哪个游戏，不是 ARK 服务器镜像时为空
func ImageGame(image string) string {
	switch image {
	case ArkServerLinuxImage, ArkServerProtonImage:
		return GameASE
	case AsaServerProtonImage:
		return GameASA
	}
	return ""
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

// 服务端在 Steam 上的 AppID
const (
	ArkServerAppID = "376030"  // ASE
	AsaServerAppID = "2430930" // ASA
)

// GameAppID 游戏的服务端 AppID
func GameAppID(game string) string {
	if NormalizeGame(game) == GameASA {
		return AsaServerAppID
	}
	return ArkServerAppID
}

// 游戏版本（Steam 分支）：ASE 主流的两个版本是蜗牛海洋正式版（public，发行商 Snail 的海洋 DLC 更新后的版本）
// 和通配符原版（preaquatica 测试分支，开发商 Wildcard 在海洋 DLC 更新前的版本）。ASA 只有正式版。
// 两个版本的服务端文件不同，复制服务端文件时只能在同一版本的服务器之间复制
const (
	GameBranchPublic      = ""            // 蜗牛海洋正式版（Steam 默认分支）
	GameBranchPreAquatica = "preaquatica" // 通配符原版（-beta preaquatica）
)

// GameBranches 可选的游戏版本（ASE）
var GameBranches = []string{GameBranchPublic, GameBranchPreAquatica}

// ValidGameBranch 是否是这个游戏可选的游戏版本
func ValidGameBranch(game, branch string) bool {
	if NormalizeGame(game) == GameASA {
		return branch == GameBranchPublic
	}
	for _, b := range GameBranches {
		if b == branch {
			return true
		}
	}
	return false
}

// ArkAppUpdateFor 容器环境变量 APP_UPDATE 的值，入口脚本按空格拆开传给 steamcmd 的 app_update。
// 正式版和镜像里的默认值相同（只有 AppID）
func ArkAppUpdateFor(game, branch string) string {
	appID := GameAppID(game)
	if branch == GameBranchPublic {
		return appID
	}
	return appID + " -beta " + branch
}

// PortalContainerName 玩家站容器名称（每个管理端一个）
const PortalContainerName = "ark-player-portal"

// PortalImage 玩家站镜像（子模块 ArkPlayerPortal）。固定版本号：快照格式和玩家站一起改，不做新旧版本兼容，
// 玩家站发了新版就在这里改版本号，随 Commander 一起发布。本地开发时可以把自己构建的镜像打上这个标签
const PortalImage = "tbro98/ark-player-portal:0.2.0"
