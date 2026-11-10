package models

import (
	"fmt"
	"maps"
	"slices"
	"strconv"
	"strings"
)

// ServerArgs ARK服务器启动参数结构
type ServerArgs struct {
	// 查询参数（以?开头的参数，不包含基础参数）
	QueryParams map[string]string `json:"query_params" gorm:"-"`

	// 命令行参数（以-开头的参数）
	CommandLineArgs map[string]interface{} `json:"command_line_args" gorm:"-"`

	// 自定义参数（用户添加的其他参数）
	CustomArgs []string `json:"custom_args" gorm:"-"`
}

// ServerArgsRequest 启动参数请求结构
type ServerArgsRequest struct {
	QueryParams     map[string]string      `json:"query_params"`
	CommandLineArgs map[string]interface{} `json:"command_line_args"`
	CustomArgs      []string               `json:"custom_args"`
}

// NewServerArgs 创建新的启动参数实例
func NewServerArgs() *ServerArgs {
	return &ServerArgs{
		QueryParams:     make(map[string]string),
		CommandLineArgs: make(map[string]interface{}),
		CustomArgs:      []string{},
	}
}

// FromServer 从Server模型创建ServerArgs
func FromServer(server Server) *ServerArgs {
	args := NewServerArgs()

	// 设置默认的查询参数（不包含基础参数，因为基础参数从Server模型获取）
	args.QueryParams["listen"] = ""

	// 设置默认的命令行参数
	args.CommandLineArgs["NoBattlEye"] = true
	args.CommandLineArgs["servergamelog"] = true
	args.CommandLineArgs["structurememopts"] = true
	args.CommandLineArgs["UseStructureStasisGrid"] = true
	args.CommandLineArgs["SecureSendArKPayload"] = true
	args.CommandLineArgs["UseItemDupeCheck"] = true
	args.CommandLineArgs["UseSecureSpawnRules"] = true
	args.CommandLineArgs["nosteamclient"] = true
	args.CommandLineArgs["game"] = true
	args.CommandLineArgs["server"] = true
	args.CommandLineArgs["log"] = true
	args.CommandLineArgs["MinimumTimeBetweenInventoryRetrieval"] = 3600
	args.CommandLineArgs["newsaveformat"] = true
	args.CommandLineArgs["usestore"] = true
	args.CommandLineArgs["BackupTransferPlayerDatas"] = true
	args.CommandLineArgs["converttostore"] = true

	return args
}

// GenerateArgs 生成启动参数列表，第一个元素是地图 URL（地图 + ?查询参数），其余是命令行参数。
// 从服务器基础参数中获取：游戏端口、查询端口、RCON端口、管理员密码、地图、模组ID
// 从启动参数中获取：其他自定义参数
// map 的遍历顺序是随机的，这里按 key 排序，保证相同配置每次生成的结果一致，
// 否则启动时会误判参数已变化而重建容器
func (sa *ServerArgs) GenerateArgs(server Server) []string {
	var url strings.Builder

	// 添加地图名称（从服务器基础参数获取）
	url.WriteString(server.Map)

	// 添加基础查询参数（从服务器基础参数获取，不在启动参数中重复设置）
	url.WriteString("?listen")
	fmt.Fprintf(&url, "?Port=%d", server.Port)
	fmt.Fprintf(&url, "?QueryPort=%d", server.QueryPort)
	fmt.Fprintf(&url, "?MaxPlayers=%d", server.MaxPlayers)
	url.WriteString("?RCONEnabled=True")
	fmt.Fprintf(&url, "?RCONPort=%d", server.RCONPort)

	// 添加服务器名称（SessionName）
	if server.SessionName != "" {
		fmt.Fprintf(&url, "?SessionName=%s", server.SessionName)
	}

	if server.GameModIds != "" {
		fmt.Fprintf(&url, "?GameModIds=%s", server.GameModIds)
	}

	// 添加自定义查询参数
	for _, key := range slices.Sorted(maps.Keys(sa.QueryParams)) {
		value := sa.QueryParams[key]
		// 跳过基础参数，避免重复
		if key == "listen" || key == "Port" || key == "QueryPort" || key == "MaxPlayers" ||
			key == "RCONEnabled" || key == "RCONPort" || key == "ServerAdminPassword" ||
			key == "SessionName" || key == "GameModIds" {
			continue
		}

		// 如果查询参数的值为空或为"False"，则不添加该参数
		if value == "" || strings.ToLower(value) == "false" {
			continue
		}

		fmt.Fprintf(&url, "?%s=%s", key, value)
	}

	args := []string{url.String()}

	// 添加命令行参数
	for _, key := range slices.Sorted(maps.Keys(sa.CommandLineArgs)) {
		switch v := sa.CommandLineArgs[key].(type) {
		case bool:
			if v {
				args = append(args, "-"+key)
			}
		case string:
			// 如果字符串为空，则不追加该参数
			if v != "" {
				args = append(args, fmt.Sprintf("-%s=%s", key, v))
			}
		case float64:
			// JSON 反序列化得到的数字都是 float64；为 0 时不追加该参数，避免大数被格式化成科学计数法
			if v != 0 {
				args = append(args, fmt.Sprintf("-%s=%s", key, strconv.FormatFloat(v, 'f', -1, 64)))
			}
		case int:
			if v != 0 {
				args = append(args, fmt.Sprintf("-%s=%d", key, v))
			}
		case nil:
		default:
			args = append(args, fmt.Sprintf("-%s=%v", key, v))
		}
	}

	// 添加集群ID（ClusterID）
	if server.ClusterID != "" {
		args = append(args, "-clusterid="+server.ClusterID)
	}

	// 添加自定义参数（一项里可能写了多个参数，按空白拆开）
	for _, custom := range sa.CustomArgs {
		args = append(args, strings.Fields(custom)...)
	}

	return args
}

// GenerateArgsString 生成完整的启动参数字符串（地图 URL 和命令行参数之间用空格分隔）
func (sa *ServerArgs) GenerateArgsString(server Server) string {
	return strings.Join(sa.GenerateArgs(server), " ")
}
