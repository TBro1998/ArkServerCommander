package utils

// ServerArgDefinition 启动参数定义
type ServerArgDefinition struct {
	Name        string `json:"name"`        // 参数名称
	Type        string `json:"type"`        // 参数类型：query, command_line
	Category    string `json:"category"`    // 参数分类
	Description string `json:"description"` // 参数描述
	Default     string `json:"default"`     // 默认值
	Required    bool   `json:"required"`    // 是否必需
	Platform    string `json:"platform"`    // 支持的平台：All, Windows, Linux
}

// GetQueryParamDefinitions 获取查询参数定义（以?开头的参数）
func GetQueryParamDefinitions() map[string]ServerArgDefinition {
	return map[string]ServerArgDefinition{
		"listen": {
			Name:        "listen",
			Type:        "query",
			Category:    "basic",
			Description: "启用服务器监听模式",
			Default:     "",
			Required:    true,
			Platform:    "All",
		},
		"Port": {
			Name:        "Port",
			Type:        "query",
			Category:    "basic",
			Description: "游戏服务器端口",
			Default:     "7777",
			Required:    true,
			Platform:    "All",
		},
		"QueryPort": {
			Name:        "QueryPort",
			Type:        "query",
			Category:    "basic",
			Description: "服务器查询端口",
			Default:     "27015",
			Required:    true,
			Platform:    "All",
		},
		"MaxPlayers": {
			Name:        "MaxPlayers",
			Type:        "query",
			Category:    "basic",
			Description: "最大玩家数量",
			Default:     "70",
			Required:    false,
			Platform:    "All",
		},
		"RCONEnabled": {
			Name:        "RCONEnabled",
			Type:        "query",
			Category:    "basic",
			Description: "启用RCON远程控制台",
			Default:     "True",
			Required:    false,
			Platform:    "All",
		},
		"RCONPort": {
			Name:        "RCONPort",
			Type:        "query",
			Category:    "basic",
			Description: "RCON端口",
			Default:     "32330",
			Required:    false,
			Platform:    "All",
		},
		"ServerAdminPassword": {
			Name:        "ServerAdminPassword",
			Type:        "query",
			Category:    "basic",
			Description: "服务器管理员密码",
			Default:     "password",
			Required:    true,
			Platform:    "All",
		},
		"GameModIds": {
			Name:        "GameModIds",
			Type:        "query",
			Category:    "mods",
			Description: "游戏模组ID列表，用逗号分隔",
			Default:     "",
			Required:    false,
			Platform:    "All",
		},
		"AllowAnyoneBabyImprintCuddle": {
			Name:        "AllowAnyoneBabyImprintCuddle",
			Type:        "query",
			Category:    "gameplay",
			Description: "允许任何人照顾幼龙（抚摸等），不仅仅是印记者",
			Default:     "false",
			Required:    false,
			Platform:    "All",
		},
		"AllowCrateSpawnsOnTopOfStructures": {
			Name:        "AllowCrateSpawnsOnTopOfStructures",
			Type:        "query",
			Category:    "gameplay",
			Description: "允许空投箱出现在建筑顶部",
			Default:     "false",
			Required:    false,
			Platform:    "All",
		},
		"AllowFlyerCarryPvE": {
			Name:        "AllowFlyerCarryPvE",
			Type:        "query",
			Category:    "gameplay",
			Description: "PvE模式下飞行生物可以携带野生恐龙",
			Default:     "false",
			Required:    false,
			Platform:    "All",
		},
		"AllowFlyingStaminaRecovery": {
			Name:        "AllowFlyingStaminaRecovery",
			Type:        "query",
			Category:    "gameplay",
			Description: "飞行生物在玩家站立时恢复耐力",
			Default:     "false",
			Required:    false,
			Platform:    "All",
		},
		"AllowMultipleAttachedC4": {
			Name:        "AllowMultipleAttachedC4",
			Type:        "query",
			Category:    "gameplay",
			Description: "允许在恐龙身上附加多个C4",
			Default:     "false",
			Required:    false,
			Platform:    "All",
		},
		"AllowThirdPersonPlayer": {
			Name:        "AllowThirdPersonPlayer",
			Type:        "query",
			Category:    "gameplay",
			Description: "启用第三人称视角",
			Default:     "True",
			Required:    false,
			Platform:    "All",
		},
		"AutoDestroyDecayedDinos": {
			Name:        "AutoDestroyDecayedDinos",
			Type:        "query",
			Category:    "gameplay",
			Description: "自动销毁可认领的腐烂恐龙",
			Default:     "false",
			Required:    false,
			Platform:    "All",
		},
		"ClampItemSpoilingTimes": {
			Name:        "ClampItemSpoilingTimes",
			Type:        "query",
			Category:    "gameplay",
			Description: "限制所有物品的腐烂时间到最大腐烂时间",
			Default:     "false",
			Required:    false,
			Platform:    "All",
		},
		"ClampItemStats": {
			Name:        "ClampItemStats",
			Type:        "query",
			Category:    "gameplay",
			Description: "启用/禁用物品属性限制",
			Default:     "false",
			Required:    false,
			Platform:    "All",
		},
		"CustomLiveTuningUrl": {
			Name:        "CustomLiveTuningUrl",
			Type:        "query",
			Category:    "advanced",
			Description: "LiveTuning文件的直接链接",
			Default:     "",
			Required:    false,
			Platform:    "All",
		},
		"customdynamicconfigurl": {
			Name:        "customdynamicconfigurl",
			Type:        "query",
			Category:    "advanced",
			Description: "动态配置文件的直接链接",
			Default:     "",
			Required:    false,
			Platform:    "All",
		},
		"DestroyUnconnectedWaterPipes": {
			Name:        "DestroyUnconnectedWaterPipes",
			Type:        "query",
			Category:    "gameplay",
			Description: "两天后自动销毁未连接的水管",
			Default:     "false",
			Required:    false,
			Platform:    "All",
		},
		"DisableImprintDinoBuff": {
			Name:        "DisableImprintDinoBuff",
			Type:        "query",
			Category:    "gameplay",
			Description: "禁用恐龙印记玩家属性加成",
			Default:     "false",
			Required:    false,
			Platform:    "All",
		},
		"DisableWeatherFog": {
			Name:        "DisableWeatherFog",
			Type:        "query",
			Category:    "gameplay",
			Description: "禁用雾天天气",
			Default:     "",
			Required:    false,
			Platform:    "All",
		},
		"EnableExtraStructurePreventionVolumes": {
			Name:        "EnableExtraStructurePreventionVolumes",
			Type:        "query",
			Category:    "gameplay",
			Description: "在特定资源丰富区域完全禁用建筑",
			Default:     "false",
			Required:    false,
			Platform:    "All",
		},
		"EnableFullDump": {
			Name:        "EnableFullDump",
			Type:        "query",
			Category:    "logging",
			Description: "启用完整的服务器日志转储",
			Default:     "",
			Required:    false,
			Platform:    "All",
		},
		"ExtinctionEventTimeInterval": {
			Name:        "ExtinctionEventTimeInterval",
			Type:        "query",
			Category:    "events",
			Description: "灭绝事件时间间隔（秒）",
			Default:     "2592000",
			Required:    false,
			Platform:    "All",
		},
		"FastDecayUnsnappedCoreStructures": {
			Name:        "FastDecayUnsnappedCoreStructures",
			Type:        "query",
			Category:    "gameplay",
			Description: "未对齐的基础结构以5倍速度腐烂",
			Default:     "false",
			Required:    false,
			Platform:    "All",
		},
		"ForceAllStructureLocking": {
			Name:        "ForceAllStructureLocking",
			Type:        "query",
			Category:    "gameplay",
			Description: "默认锁定所有结构",
			Default:     "false",
			Required:    false,
			Platform:    "All",
		},
		"ForceFlyerExplosives": {
			Name:        "ForceFlyerExplosives",
			Type:        "query",
			Category:    "gameplay",
			Description: "飞行生物可以携带C4飞行",
			Default:     "false",
			Required:    false,
			Platform:    "All",
		},
		"MapPlayerLocation": {
			Name:        "MapPlayerLocation",
			Type:        "query",
			Category:    "gameplay",
			Description: "在地图上显示玩家精确位置",
			Default:     "True",
			Required:    false,
			Platform:    "All",
		},
		"MaxPersonalTamedDinos": {
			Name:        "MaxPersonalTamedDinos",
			Type:        "query",
			Category:    "gameplay",
			Description: "每个部落的恐龙驯服限制",
			Default:     "500",
			Required:    false,
			Platform:    "All",
		},
		"MaxTamedDinos": {
			Name:        "MaxTamedDinos",
			Type:        "query",
			Category:    "gameplay",
			Description: "服务器上驯服恐龙的最大数量",
			Default:     "4000",
			Required:    false,
			Platform:    "All",
		},
		"MinimumDinoReuploadInterval": {
			Name:        "MinimumDinoReuploadInterval",
			Type:        "query",
			Category:    "gameplay",
			Description: "恐龙重新上传的最小间隔（秒）",
			Default:     "0",
			Required:    false,
			Platform:    "All",
		},
		"NonPermanentDiseases": {
			Name:        "NonPermanentDiseases",
			Type:        "query",
			Category:    "gameplay",
			Description: "疾病不是永久的（重生后会失去）",
			Default:     "false",
			Required:    false,
			Platform:    "All",
		},
		"OnlyAutoDestroyCoreStructures": {
			Name:        "OnlyAutoDestroyCoreStructures",
			Type:        "query",
			Category:    "gameplay",
			Description: "只自动销毁核心/基础结构",
			Default:     "false",
			Required:    false,
			Platform:    "All",
		},
		"OnlyDecayUnsnappedCoreStructures": {
			Name:        "OnlyDecayUnsnappedCoreStructures",
			Type:        "query",
			Category:    "gameplay",
			Description: "只有未对齐的核心结构会腐烂",
			Default:     "false",
			Required:    false,
			Platform:    "All",
		},
		"OverrideOfficialDifficulty": {
			Name:        "OverrideOfficialDifficulty",
			Type:        "query",
			Category:    "gameplay",
			Description: "覆盖默认服务器难度等级",
			Default:     "5.0",
			Required:    false,
			Platform:    "All",
		},
		"OverrideStructurePlatformPrevention": {
			Name:        "OverrideStructurePlatformPrevention",
			Type:        "query",
			Category:    "gameplay",
			Description: "允许在平台鞍座上建造炮塔或尖刺结构",
			Default:     "false",
			Required:    false,
			Platform:    "All",
		},
		"OxygenSwimSpeedStatMultiplier": {
			Name:        "OxygenSwimSpeedStatMultiplier",
			Type:        "query",
			Category:    "gameplay",
			Description: "氧气等级对游泳速度的倍数",
			Default:     "1.0",
			Required:    false,
			Platform:    "All",
		},
		"PersonalTamedDinosSaddleStructureCost": {
			Name:        "PersonalTamedDinosSaddleStructureCost",
			Type:        "query",
			Category:    "gameplay",
			Description: "平台鞍座占用的恐龙槽位数量",
			Default:     "19",
			Required:    false,
			Platform:    "All",
		},
		"PreventDownloadSurvivors": {
			Name:        "PreventDownloadSurvivors",
			Type:        "query",
			Category:    "transfer",
			Description: "阻止跨ARK数据转移幸存者",
			Default:     "False",
			Required:    false,
			Platform:    "All",
		},
		"PreventDownloadItems": {
			Name:        "PreventDownloadItems",
			Type:        "query",
			Category:    "transfer",
			Description: "阻止跨ARK数据转移物品",
			Default:     "False",
			Required:    false,
			Platform:    "All",
		},
		"PreventDownloadDinos": {
			Name:        "PreventDownloadDinos",
			Type:        "query",
			Category:    "transfer",
			Description: "阻止跨ARK数据转移恐龙",
			Default:     "False",
			Required:    false,
			Platform:    "All",
		},
		"PreventDiseases": {
			Name:        "PreventDiseases",
			Type:        "query",
			Category:    "gameplay",
			Description: "完全禁用疾病",
			Default:     "false",
			Required:    false,
			Platform:    "All",
		},
		"PreventMateBoost": {
			Name:        "PreventMateBoost",
			Type:        "query",
			Category:    "gameplay",
			Description: "禁用恐龙交配加成",
			Default:     "false",
			Required:    false,
			Platform:    "All",
		},
		"PreventOfflinePvP": {
			Name:        "PreventOfflinePvP",
			Type:        "query",
			Category:    "pvp",
			Description: "启用离线PvP防护",
			Default:     "false",
			Required:    false,
			Platform:    "All",
		},
		"PreventOfflinePvPInterval": {
			Name:        "PreventOfflinePvPInterval",
			Type:        "query",
			Category:    "pvp",
			Description: "离线PvP防护等待时间（秒）",
			Default:     "900",
			Required:    false,
			Platform:    "All",
		},
		"PreventSpawnAnimations": {
			Name:        "PreventSpawnAnimations",
			Type:        "query",
			Category:    "gameplay",
			Description: "禁用玩家角色重生动画",
			Default:     "false",
			Required:    false,
			Platform:    "All",
		},
		"PreventTribeAlliances": {
			Name:        "PreventTribeAlliances",
			Type:        "query",
			Category:    "gameplay",
			Description: "阻止部落创建联盟",
			Default:     "false",
			Required:    false,
			Platform:    "All",
		},
		"PreventUploadSurvivors": {
			Name:        "PreventUploadSurvivors",
			Type:        "query",
			Category:    "transfer",
			Description: "阻止上传幸存者",
			Default:     "False",
			Required:    false,
			Platform:    "All",
		},
		"PreventUploadItems": {
			Name:        "PreventUploadItems",
			Type:        "query",
			Category:    "transfer",
			Description: "阻止上传物品",
			Default:     "False",
			Required:    false,
			Platform:    "All",
		},
		"PreventUploadDinos": {
			Name:        "PreventUploadDinos",
			Type:        "query",
			Category:    "transfer",
			Description: "阻止上传恐龙",
			Default:     "False",
			Required:    false,
			Platform:    "All",
		},
		"RandomSupplyCratePoints": {
			Name:        "RandomSupplyCratePoints",
			Type:        "query",
			Category:    "gameplay",
			Description: "空投箱在随机位置出现",
			Default:     "",
			Required:    false,
			Platform:    "All",
		},
		"ServerCrosshair": {
			Name:        "ServerCrosshair",
			Type:        "query",
			Category:    "gameplay",
			Description: "禁用服务器准星",
			Default:     "False",
			Required:    false,
			Platform:    "All",
		},
		"ShowFloatingDamageText": {
			Name:        "ShowFloatingDamageText",
			Type:        "query",
			Category:    "gameplay",
			Description: "启用RPG风格的弹出伤害文本",
			Default:     "false",
			Required:    false,
			Platform:    "All",
		},
		"TheMaxStructuresInRange": {
			Name:        "TheMaxStructuresInRange",
			Type:        "query",
			Category:    "gameplay",
			Description: "服务器允许的最大结构数量",
			Default:     "10500",
			Required:    false,
			Platform:    "All",
		},
		"TribeLogDestroyedEnemyStructures": {
			Name:        "TribeLogDestroyedEnemyStructures",
			Type:        "query",
			Category:    "gameplay",
			Description: "在部落日志中显示敌方结构破坏",
			Default:     "false",
			Required:    false,
			Platform:    "All",
		},
		"UseOptimizedHarvestingHealth": {
			Name:        "UseOptimizedHarvestingHealth",
			Type:        "query",
			Category:    "gameplay",
			Description: "使用优化的收获健康系统",
			Default:     "false",
			Required:    false,
			Platform:    "All",
		},
	}
}

// GetCommandLineArgDefinitions 获取命令行参数定义（以-开头的参数）
func GetCommandLineArgDefinitions() map[string]ServerArgDefinition {
	return map[string]ServerArgDefinition{
		"allcores": {
			Name:        "allcores",
			Type:        "command_line",
			Category:    "performance",
			Description: "使用所有CPU核心",
			Default:     "",
			Required:    false,
			Platform:    "All",
		},
		"allowansel": {
			Name:        "allowansel",
			Type:        "command_line",
			Category:    "graphics",
			Description: "激活NVIDIA Ansel支持",
			Default:     "",
			Required:    false,
			Platform:    "Windows",
		},
		"automanagedmods": {
			Name:        "automanagedmods",
			Type:        "command_line",
			Category:    "mods",
			Description: "自动MOD下载/安装/更新",
			Default:     "",
			Required:    false,
			Platform:    "All",
		},
		"ClearOldItems": {
			Name:        "ClearOldItems",
			Type:        "command_line",
			Category:    "maintenance",
			Description: "一次性清除所有旧的不装备物品",
			Default:     "",
			Required:    false,
			Platform:    "All",
		},
		"culture": {
			Name:        "culture",
			Type:        "command_line",
			Category:    "localization",
			Description: "直接覆盖语言设置",
			Default:     "en",
			Required:    false,
			Platform:    "All",
		},
		"d3d10": {
			Name:        "d3d10",
			Type:        "command_line",
			Category:    "graphics",
			Description: "强制使用DX10而不是DX11",
			Default:     "",
			Required:    false,
			Platform:    "Windows",
		},
		"EnableIdlePlayerKick": {
			Name:        "EnableIdlePlayerKick",
			Type:        "command_line",
			Category:    "gameplay",
			Description: "踢出空闲玩家",
			Default:     "",
			Required:    false,
			Platform:    "All",
		},
		"ForceAllowCaveFlyers": {
			Name:        "ForceAllowCaveFlyers",
			Type:        "command_line",
			Category:    "gameplay",
			Description: "强制允许飞行生物进入洞穴",
			Default:     "",
			Required:    false,
			Platform:    "All",
		},
		"ForceRespawnDinos": {
			Name:        "ForceRespawnDinos",
			Type:        "command_line",
			Category:    "gameplay",
			Description: "启动时销毁所有野生生物",
			Default:     "",
			Required:    false,
			Platform:    "All",
		},
		"gameplaylogging": {
			Name:        "gameplaylogging",
			Type:        "command_line",
			Category:    "logging",
			Description: "输出游戏日志到文件",
			Default:     "",
			Required:    false,
			Platform:    "All",
		},
		"insecure": {
			Name:        "insecure",
			Type:        "command_line",
			Category:    "security",
			Description: "禁用Valve反作弊系统",
			Default:     "",
			Required:    false,
			Platform:    "Windows",
		},
		"lowmemory": {
			Name:        "lowmemory",
			Type:        "command_line",
			Category:    "performance",
			Description: "低内存模式，减少图形和音频效果",
			Default:     "",
			Required:    false,
			Platform:    "All",
		},
		"MapModID": {
			Name:        "MapModID",
			Type:        "command_line",
			Category:    "mods",
			Description: "通过ModID直接加载自定义地图",
			Default:     "",
			Required:    false,
			Platform:    "All",
		},
		"noantispeedhack": {
			Name:        "noantispeedhack",
			Type:        "command_line",
			Category:    "security",
			Description: "禁用反速度黑客检测",
			Default:     "",
			Required:    false,
			Platform:    "All",
		},
		"NoBattlEye": {
			Name:        "NoBattlEye",
			Type:        "command_line",
			Category:    "security",
			Description: "不使用BattleEye",
			Default:     "",
			Required:    false,
			Platform:    "All",
		},
		"NoBiomeWalls": {
			Name:        "NoBiomeWalls",
			Type:        "command_line",
			Category:    "graphics",
			Description: "消除生物群系变化区域墙效果",
			Default:     "",
			Required:    false,
			Platform:    "All",
		},
		"nocombineclientmoves": {
			Name:        "nocombineclientmoves",
			Type:        "command_line",
			Category:    "performance",
			Description: "禁用服务器玩家移动物理优化",
			Default:     "",
			Required:    false,
			Platform:    "All",
		},
		"nofishloot": {
			Name:        "nofishloot",
			Type:        "command_line",
			Category:    "gameplay",
			Description: "禁用钓鱼竿的非肉类鱼类战利品",
			Default:     "",
			Required:    false,
			Platform:    "All",
		},
		"noninlinesaveload": {
			Name:        "noninlinesaveload",
			Type:        "command_line",
			Category:    "performance",
			Description: "修复大型存档的存档损坏问题",
			Default:     "",
			Required:    false,
			Platform:    "All",
		},
		"nomansky": {
			Name:        "nomansky",
			Type:        "command_line",
			Category:    "graphics",
			Description: "禁用详细的天空功能",
			Default:     "",
			Required:    false,
			Platform:    "All",
		},
		"nomemorybias": {
			Name:        "nomemorybias",
			Type:        "command_line",
			Category:    "performance",
			Description: "减少客户端游戏内存使用",
			Default:     "",
			Required:    false,
			Platform:    "All",
		},
		"oldsaveformat": {
			Name:        "oldsaveformat",
			Type:        "command_line",
			Category:    "compatibility",
			Description: "使用旧的存档格式",
			Default:     "",
			Required:    false,
			Platform:    "All",
		},
		"PreventHibernation": {
			Name:        "PreventHibernation",
			Type:        "command_line",
			Category:    "performance",
			Description: "防止非活跃区域的生物进入休眠状态",
			Default:     "",
			Required:    false,
			Platform:    "All",
		},
		"server": {
			Name:        "server",
			Type:        "command_line",
			Category:    "basic",
			Description: "以服务器模式运行",
			Default:     "",
			Required:    true,
			Platform:    "All",
		},
		"servergamelog": {
			Name:        "servergamelog",
			Type:        "command_line",
			Category:    "logging",
			Description: "启用服务器管理员日志",
			Default:     "",
			Required:    false,
			Platform:    "All",
		},
		"servergamelogincludetribelogs": {
			Name:        "servergamelogincludetribelogs",
			Type:        "command_line",
			Category:    "logging",
			Description: "在服务器游戏日志中包含部落日志",
			Default:     "",
			Required:    false,
			Platform:    "All",
		},
		"ServerRCONOutputTribeLogs": {
			Name:        "ServerRCONOutputTribeLogs",
			Type:        "command_line",
			Category:    "logging",
			Description: "在RCON中显示部落聊天",
			Default:     "",
			Required:    false,
			Platform:    "All",
		},
		"sm4": {
			Name:        "sm4",
			Type:        "command_line",
			Category:    "graphics",
			Description: "强制使用Shader Model 4",
			Default:     "",
			Required:    false,
			Platform:    "Windows",
		},
		"StasisKeepControllers": {
			Name:        "StasisKeepControllers",
			Type:        "command_line",
			Category:    "performance",
			Description: "在静止状态下保持AI控制器",
			Default:     "",
			Required:    false,
			Platform:    "All",
		},
		"StructureDestructionTag": {
			Name:        "StructureDestructionTag",
			Type:        "command_line",
			Category:    "maintenance",
			Description: "沼泽和雪地区域的一次性自动结构拆除",
			Default:     "",
			Required:    false,
			Platform:    "All",
		},
		"USEALLAVAILABLECORES": {
			Name:        "USEALLAVAILABLECORES",
			Type:        "command_line",
			Category:    "performance",
			Description: "使用所有可用的CPU核心",
			Default:     "",
			Required:    false,
			Platform:    "All",
		},
		"usecache": {
			Name:        "usecache",
			Type:        "command_line",
			Category:    "performance",
			Description: "使用缓存提高加载速度",
			Default:     "",
			Required:    false,
			Platform:    "All",
		},
		"webalarm": {
			Name:        "webalarm",
			Type:        "command_line",
			Category:    "notifications",
			Description: "激活Web警报",
			Default:     "",
			Required:    false,
			Platform:    "All",
		},
		"exclusivejoin": {
			Name:        "exclusivejoin",
			Type:        "command_line",
			Category:    "security",
			Description: "激活白名单模式",
			Default:     "",
			Required:    false,
			Platform:    "All",
		},
		"ActiveEvent": {
			Name:        "ActiveEvent",
			Type:        "command_line",
			Category:    "events",
			Description: "启用指定事件",
			Default:     "",
			Required:    false,
			Platform:    "PC",
		},
		"crossplay": {
			Name:        "crossplay",
			Type:        "command_line",
			Category:    "compatibility",
			Description: "启用跨平台游戏（EPIC和Steam之间）",
			Default:     "",
			Required:    false,
			Platform:    "PC",
		},
		"epiconly": {
			Name:        "epiconly",
			Type:        "command_line",
			Category:    "compatibility",
			Description: "只允许Epic游戏商店玩家连接",
			Default:     "",
			Required:    false,
			Platform:    "PC",
		},
		"UseVivox": {
			Name:        "UseVivox",
			Type:        "command_line",
			Category:    "voice",
			Description: "在Steam服务器上启用Vivox",
			Default:     "",
			Required:    false,
			Platform:    "PC",
		},
		"PublicIPForEpic": {
			Name:        "PublicIPForEpic",
			Type:        "command_line",
			Category:    "network",
			Description: "EGS客户端尝试连接的公共IP",
			Default:     "",
			Required:    false,
			Platform:    "PC",
		},
		"high": {
			Name:        "high",
			Type:        "command_line",
			Category:    "performance",
			Description: "以高优先级运行游戏",
			Default:     "",
			Required:    false,
			Platform:    "PC",
		},
		"noaafonts": {
			Name:        "noaafonts",
			Type:        "command_line",
			Category:    "graphics",
			Description: "移除字体抗锯齿",
			Default:     "",
			Required:    false,
			Platform:    "PC",
		},
		"UseItemDupeCheck": {
			Name:        "UseItemDupeCheck",
			Type:        "command_line",
			Category:    "security",
			Description: "启用额外的复制保护",
			Default:     "",
			Required:    false,
			Platform:    "PC",
		},
		// 默认参数
		"structurememopts": {
			Name:        "structurememopts",
			Type:        "command_line",
			Category:    "performance",
			Description: "结构内存优化",
			Default:     "",
			Required:    false,
			Platform:    "All",
		},
		"UseStructureStasisGrid": {
			Name:        "UseStructureStasisGrid",
			Type:        "command_line",
			Category:    "performance",
			Description: "使用结构静止网格",
			Default:     "",
			Required:    false,
			Platform:    "All",
		},
		"SecureSendArKPayload": {
			Name:        "SecureSendArKPayload",
			Type:        "command_line",
			Category:    "security",
			Description: "安全发送ARK负载",
			Default:     "",
			Required:    false,
			Platform:    "All",
		},
		"UseSecureSpawnRules": {
			Name:        "UseSecureSpawnRules",
			Type:        "command_line",
			Category:    "security",
			Description: "使用安全重生规则",
			Default:     "",
			Required:    false,
			Platform:    "All",
		},
		"nosteamclient": {
			Name:        "nosteamclient",
			Type:        "command_line",
			Category:    "compatibility",
			Description: "不使用Steam客户端",
			Default:     "",
			Required:    false,
			Platform:    "All",
		},
		"game": {
			Name:        "game",
			Type:        "command_line",
			Category:    "basic",
			Description: "游戏模式",
			Default:     "",
			Required:    false,
			Platform:    "All",
		},
		"log": {
			Name:        "log",
			Type:        "command_line",
			Category:    "logging",
			Description: "启用日志记录",
			Default:     "",
			Required:    false,
			Platform:    "All",
		},
		"MinimumTimeBetweenInventoryRetrieval": {
			Name:        "MinimumTimeBetweenInventoryRetrieval",
			Type:        "command_line",
			Category:    "gameplay",
			Description: "物品检索之间的最小时间间隔（秒）",
			Default:     "3600",
			Required:    false,
			Platform:    "All",
		},
		"newsaveformat": {
			Name:        "newsaveformat",
			Type:        "command_line",
			Category:    "compatibility",
			Description: "使用新的存档格式",
			Default:     "",
			Required:    false,
			Platform:    "All",
		},
		"usestore": {
			Name:        "usestore",
			Type:        "command_line",
			Category:    "compatibility",
			Description: "使用商店功能",
			Default:     "",
			Required:    false,
			Platform:    "All",
		},
	}
}

// GetArgCategories 获取参数分类
func GetArgCategories() map[string]string {
	return map[string]string{
		"basic":         "基础设置",
		"gameplay":      "游戏玩法",
		"performance":   "性能优化",
		"security":      "安全设置",
		"network":       "网络设置",
		"graphics":      "图形设置",
		"logging":       "日志记录",
		"mods":          "模组设置",
		"events":        "事件设置",
		"pvp":           "PvP设置",
		"transfer":      "数据转移",
		"maintenance":   "维护设置",
		"compatibility": "兼容性",
		"localization":  "本地化",
		"notifications": "通知设置",
		"voice":         "语音设置",
		"advanced":      "高级设置",
	}
}
