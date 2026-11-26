package models

// ServerConfigFormat 导出的服务器配置文件的格式标识
const (
	ServerConfigFormat  = "ark-server-commander/server-config"
	ServerConfigVersion = 1
)

// 导入、复制配置时可以选择的部分
const (
	ConfigPartConfigFiles = "config_files" // GameUserSettings.ini 和 Game.ini（保留目标服务器的名称、端口、人数和管理员密码）
	ConfigPartServerArgs  = "server_args"  // 启动参数
	ConfigPartMods        = "mods"         // Mod 列表
	ConfigPartMap         = "map"          // 地图
	ConfigPartCluster     = "cluster"      // 集群 ID
	ConfigPartOptions     = "options"      // 最大玩家数、ArkApi、自动重启、自动更新
)

// ConfigParts 所有可选的部分
var ConfigParts = []string{ConfigPartConfigFiles, ConfigPartServerArgs, ConfigPartMods, ConfigPartMap, ConfigPartCluster, ConfigPartOptions}

// ServerConfig 导出的服务器配置：不含端口和管理员密码（配置文件里的 ServerAdminPassword 也已去掉），
// 可以导入到新建或已有的服务器，或复制到其他服务器
type ServerConfig struct {
	Format     string `json:"format"`
	Version    int    `json:"version"`
	ExportedAt string `json:"exported_at,omitempty"`
	// 来源服务器的游戏（ase / asa），旧版导出的没有，视为 ase。导入到另一个游戏的服务器时，
	// 地图、Mod 和启动参数不通用，不会导入
	Game        string `json:"game,omitempty"`
	SessionName string `json:"session_name,omitempty"` // 来源服务器的名称，只用于显示，不会导入
	Map         string `json:"map,omitempty"`
	MaxPlayers  int    `json:"max_players,omitempty"`
	ClusterID   string `json:"cluster_id"`
	GameModIds  string `json:"game_mod_ids"`
	ArkApi      bool   `json:"arkapi"`
	AutoRestart bool   `json:"auto_restart"`
	AutoUpdate  bool   `json:"auto_update"`
	// 配置文件内容
	GameUserSettings string `json:"game_user_settings"`
	GameIni          string `json:"game_ini"`
	// 启动参数
	ServerArgs *ServerArgs `json:"server_args,omitempty"`
}

// ApplyConfigRequest 把配置导入到一台已有的服务器
type ApplyConfigRequest struct {
	Config ServerConfig `json:"config" binding:"required"`
	Parts  []string     `json:"parts"` // 为空表示全部
}

// CopyConfigRequest 把一台服务器的配置复制到其他服务器
type CopyConfigRequest struct {
	TargetIDs []uint   `json:"target_ids" binding:"required,min=1"`
	Parts     []string `json:"parts"` // 为空表示全部
}

// CopyConfigResult 复制到一台服务器的结果
type CopyConfigResult struct {
	ServerID    uint   `json:"server_id"`
	SessionName string `json:"session_name,omitempty"`
	Success     bool   `json:"success"`
	Error       string `json:"error,omitempty"`
	Code        string `json:"code,omitempty"`
	// 目标服务器未停止，配置文件已暂存，下次启动时生效
	ConfigPending bool `json:"config_pending,omitempty"`
}
