package models

import (
	"encoding/json"
	"fmt"
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

// ServerArgsResponse 启动参数响应结构
type ServerArgsResponse struct {
	QueryParams     map[string]string      `json:"query_params"`
	CommandLineArgs map[string]interface{} `json:"command_line_args"`
	CustomArgs      []string               `json:"custom_args"`
	GeneratedArgs   string                 `json:"generated_args"` // 生成的完整启动参数字符串
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

	return args
}

// GenerateArgsString 生成完整的启动参数字符串
// 从服务器基础参数中获取：游戏端口、查询端口、RCON端口、管理员密码、地图、模组ID
// 从启动参数中获取：其他自定义参数
func (sa *ServerArgs) GenerateArgsString(server Server) string {
	var queryParams []string
	var commandLineParams []string

	// 添加地图名称（从服务器基础参数获取）
	result := server.Map

	// 添加基础查询参数（从服务器基础参数获取，不在启动参数中重复设置）
	queryParams = append(queryParams, "?listen")
	queryParams = append(queryParams, fmt.Sprintf("?Port=%d", server.Port))
	queryParams = append(queryParams, fmt.Sprintf("?QueryPort=%d", server.QueryPort))
	queryParams = append(queryParams, fmt.Sprintf("?MaxPlayers=%d", server.MaxPlayers))
	queryParams = append(queryParams, "?RCONEnabled=True")
	queryParams = append(queryParams, fmt.Sprintf("?RCONPort=%d", server.RCONPort))
	queryParams = append(queryParams, fmt.Sprintf("?ServerAdminPassword=%s", server.AdminPassword))

	// 添加服务器名称（SessionName）
	if server.SessionName != "" {
		queryParams = append(queryParams, fmt.Sprintf("?SessionName=%s", server.SessionName))
	}

	if server.GameModIds != "" {
		queryParams = append(queryParams, fmt.Sprintf("?GameModIds=%s", server.GameModIds))
	}

	// 添加自定义查询参数
	for key, value := range sa.QueryParams {
		// 跳过基础参数，避免重复
		if key == "listen" || key == "Port" || key == "QueryPort" || key == "MaxPlayers" ||
			key == "RCONEnabled" || key == "RCONPort" || key == "ServerAdminPassword" || key == "GameModIds" {
			continue
		}

		// 如果查询参数的值为空或为"False"，则不添加该参数
		if value == "" || strings.ToLower(value) == "false" {
			continue
		}

		queryParams = append(queryParams, fmt.Sprintf("?%s=%s", key, value))
	}

	// 添加命令行参数
	for key, value := range sa.CommandLineArgs {
		switch v := value.(type) {
		case bool:
			if v {
				commandLineParams = append(commandLineParams, fmt.Sprintf("-%s", key))
			}
		case string:
			if v != "" {
				commandLineParams = append(commandLineParams, fmt.Sprintf("-%s=%s", key, v))
			}
			// 如果字符串为空，则不追加该参数
		case int, int32, int64, float32, float64:
			// 检查数值是否为0，如果是0则不追加该参数
			if v != 0 {
				commandLineParams = append(commandLineParams, fmt.Sprintf("-%s=%v", key, v))
			}
		default:
			commandLineParams = append(commandLineParams, fmt.Sprintf("-%s=%v", key, v))
		}
	}

	// 添加集群ID（ClusterID）
	if server.ClusterID != "" {
		commandLineParams = append(commandLineParams, fmt.Sprintf("-clusterid=%s", server.ClusterID))
	}

	// 添加自定义参数
	commandLineParams = append(commandLineParams, sa.CustomArgs...)

	// 组合参数：地图 + 查询参数（无空格连接）+ 空格 + 命令行参数
	if len(queryParams) > 0 {
		result += strings.Join(queryParams, "")
	}

	if len(commandLineParams) > 0 {
		if len(queryParams) > 0 {
			result += " " // 查询参数和命令行参数之间加空格
		}
		result += strings.Join(commandLineParams, " ")
	}

	return result
}

// MarshalJSON 自定义JSON序列化
func (sa *ServerArgs) MarshalJSON() ([]byte, error) {
	type Alias ServerArgs
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(sa),
	})
}

// UnmarshalJSON 自定义JSON反序列化
func (sa *ServerArgs) UnmarshalJSON(data []byte) error {
	type Alias ServerArgs
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(sa),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	return nil
}
