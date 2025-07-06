// ARK服务器启动参数定义
// 查询参数（以?开头的参数）
export const queryParams = {
  // 基础参数（从服务器配置自动生成，不在UI中显示）
  listen: { name: 'listen', type: 'boolean', category: 'basic', description: '监听模式', default: true, hidden: true },
  Port: { name: 'Port', type: 'number', category: 'basic', description: '游戏端口', default: 7777, hidden: true },
  QueryPort: { name: 'QueryPort', type: 'number', category: 'basic', description: '查询端口', default: 27015, hidden: true },
  RCONPort: { name: 'RCONPort', type: 'number', category: 'basic', description: 'RCON端口', default: 32330, hidden: true },
  ServerAdminPassword: { name: 'ServerAdminPassword', type: 'string', category: 'basic', description: '管理员密码', default: '', hidden: true },
  SessionName: { name: 'SessionName', type: 'string', category: 'basic', description: '服务器名称', default: '', hidden: true },
  GameModIds: { name: 'GameModIds', type: 'string', category: 'basic', description: '模组ID列表', default: '', hidden: true },

  // 核心游戏设置
  OverrideOfficialDifficulty: { name: 'OverrideOfficialDifficulty', type: 'number', category: 'core', description: '覆盖官方难度等级', default: 4.0, min: 1.0, max: 10.0, step: 0.1 },
  AllowThirdPersonPlayer: { name: 'AllowThirdPersonPlayer', type: 'boolean', category: 'core', description: '启用第三人称视角', default: true },
  ServerCrosshair: { name: 'ServerCrosshair', type: 'boolean', category: 'core', description: '禁用准星', default: false },
  ShowFloatingDamageText: { name: 'ShowFloatingDamageText', type: 'boolean', category: 'core', description: '启用RPG风格伤害弹出文本', default: false },
  MapPlayerLocation: { name: 'MapPlayerLocation', type: 'boolean', category: 'core', description: '在地图上显示玩家精确位置', default: false },
  DisableWeatherFog: { name: 'DisableWeatherFog', type: 'boolean', category: 'core', description: '禁用雾天', default: false },

  // 恐龙管理
  MaxPersonalTamedDinos: { name: 'MaxPersonalTamedDinos', type: 'number', category: 'dinos', description: '每个部落的恐龙驯服上限', default: 500, min: 1, max: 10000 },
  MaxTamedDinos: { name: 'MaxTamedDinos', type: 'number', category: 'dinos', description: '服务器驯服恐龙总数上限', default: 4000, min: 1, max: 100000 },
  PersonalTamedDinosSaddleStructureCost: { name: 'PersonalTamedDinosSaddleStructureCost', type: 'number', category: 'dinos', description: '平台鞍占用恐龙槽位数量', default: 19, min: 1, max: 100 },
  MinimumDinoReuploadInterval: { name: 'MinimumDinoReuploadInterval', type: 'number', category: 'dinos', description: '恐龙重新上传冷却时间（秒）', default: 0, min: 0, max: 86400 },
  ServerAutoForceRespawnWildDinosInterval: { name: 'ServerAutoForceRespawnWildDinosInterval', type: 'number', category: 'dinos', description: '强制野生恐龙重生间隔', default: 0, min: 0, max: 31536000 },
  AllowAnyoneBabyImprintCuddle: { name: 'AllowAnyoneBabyImprintCuddle', type: 'boolean', category: 'dinos', description: '允许任何人照顾幼龙', default: false },
  DisableImprintDinoBuff: { name: 'DisableImprintDinoBuff', type: 'boolean', category: 'dinos', description: '禁用恐龙印记玩家状态加成', default: false },
  PreventMateBoost: { name: 'PreventMateBoost', type: 'boolean', category: 'dinos', description: '禁用恐龙配偶加成', default: false },
  AllowFlyerCarryPvE: { name: 'AllowFlyerCarryPvE', type: 'boolean', category: 'dinos', description: 'PvE服务器允许飞行生物携带野生恐龙', default: false },
  AllowFlyingStaminaRecovery: { name: 'AllowFlyingStaminaRecovery', type: 'boolean', category: 'dinos', description: '飞行生物在玩家站立时恢复耐力', default: false },
  ForceFlyerExplosives: { name: 'ForceFlyerExplosives', type: 'boolean', category: 'dinos', description: '允许飞行生物携带C4飞行', default: false },
  AllowMultipleAttachedC4: { name: 'AllowMultipleAttachedC4', type: 'boolean', category: 'dinos', description: '允许在恐龙上附加多个C4', default: false },
  AutoDestroyDecayedDinos: { name: 'AutoDestroyDecayedDinos', type: 'boolean', category: 'dinos', description: '自动销毁已腐烂的可认领恐龙', default: false },

  // 建筑系统
  TheMaxStructuresInRange: { name: 'TheMaxStructuresInRange', type: 'number', category: 'structures', description: '服务器最大建筑数量', default: 10500, min: 1000, max: 50000 },
  ForceAllStructureLocking: { name: 'ForceAllStructureLocking', type: 'boolean', category: 'structures', description: '默认锁定所有建筑', default: false },
  EnableExtraStructurePreventionVolumes: { name: 'EnableExtraStructurePreventionVolumes', type: 'boolean', category: 'structures', description: '在特定资源丰富区域禁用建筑', default: false },
  OverrideStructurePlatformPrevention: { name: 'OverrideStructurePlatformPrevention', type: 'boolean', category: 'structures', description: '允许在平台鞍上建造炮塔', default: false },
  DestroyUnconnectedWaterPipes: { name: 'DestroyUnconnectedWaterPipes', type: 'boolean', category: 'structures', description: '自动销毁未连接的水管', default: false },
  FastDecayUnsnappedCoreStructures: { name: 'FastDecayUnsnappedCoreStructures', type: 'boolean', category: 'structures', description: '未连接核心建筑快速腐烂', default: false },
  OnlyAutoDestroyCoreStructures: { name: 'OnlyAutoDestroyCoreStructures', type: 'boolean', category: 'structures', description: '仅自动销毁核心建筑', default: false },
  OnlyDecayUnsnappedCoreStructures: { name: 'OnlyDecayUnsnappedCoreStructures', type: 'boolean', category: 'structures', description: '仅未连接核心建筑腐烂', default: false },

  // PvP设置
  PreventOfflinePvP: { name: 'PreventOfflinePvP', type: 'boolean', category: 'pvp', description: '启用离线PvP保护', default: false },
  PreventOfflinePvPInterval: { name: 'PreventOfflinePvPInterval', type: 'number', category: 'pvp', description: '离线PvP保护等待时间（秒）', default: 900, min: 0, max: 86400 },
  PreventTribeAlliances: { name: 'PreventTribeAlliances', type: 'boolean', category: 'pvp', description: '禁用部落联盟', default: false },

  // 游戏机制
  AllowCrateSpawnsOnTopOfStructures: { name: 'AllowCrateSpawnsOnTopOfStructures', type: 'boolean', category: 'mechanics', description: '允许空投箱出现在建筑上', default: false },
  RandomSupplyCratePoints: { name: 'RandomSupplyCratePoints', type: 'boolean', category: 'mechanics', description: '空投箱随机位置', default: false },
  NonPermanentDiseases: { name: 'NonPermanentDiseases', type: 'boolean', category: 'mechanics', description: '疾病非永久性（重生后消失）', default: false },
  PreventDiseases: { name: 'PreventDiseases', type: 'boolean', category: 'mechanics', description: '完全禁用疾病', default: false },
  PreventSpawnAnimations: { name: 'PreventSpawnAnimations', type: 'boolean', category: 'mechanics', description: '禁用重生动画', default: false },
  UseOptimizedHarvestingHealth: { name: 'UseOptimizedHarvestingHealth', type: 'boolean', category: 'mechanics', description: '优化收获健康值', default: false },
  ClampItemSpoilingTimes: { name: 'ClampItemSpoilingTimes', type: 'boolean', category: 'mechanics', description: '限制物品腐烂时间到最大值', default: false },
  ClampItemStats: { name: 'ClampItemStats', type: 'boolean', category: 'mechanics', description: '启用物品属性限制', default: false },
  ExtinctionEventTimeInterval: { name: 'ExtinctionEventTimeInterval', type: 'number', category: 'mechanics', description: '灭绝事件时间间隔（秒）', default: 0, min: 0, max: 31536000 },
  OxygenSwimSpeedStatMultiplier: { name: 'OxygenSwimSpeedStatMultiplier', type: 'number', category: 'mechanics', description: '氧气游泳速度倍数', default: 1.0, min: 0.1, max: 10.0, step: 0.1 },

  // 数据传输控制
  PreventDownloadSurvivors: { name: 'PreventDownloadSurvivors', type: 'boolean', category: 'transfer', description: '禁止下载幸存者数据', default: false },
  PreventDownloadItems: { name: 'PreventDownloadItems', type: 'boolean', category: 'transfer', description: '禁止下载物品数据', default: false },
  PreventDownloadDinos: { name: 'PreventDownloadDinos', type: 'boolean', category: 'transfer', description: '禁止下载恐龙数据', default: false },
  PreventUploadSurvivors: { name: 'PreventUploadSurvivors', type: 'boolean', category: 'transfer', description: '禁止上传幸存者数据', default: false },
  PreventUploadItems: { name: 'PreventUploadItems', type: 'boolean', category: 'transfer', description: '禁止上传物品数据', default: false },
  PreventUploadDinos: { name: 'PreventUploadDinos', type: 'boolean', category: 'transfer', description: '禁止上传恐龙数据', default: false },

  // 高级配置
  CustomLiveTuningUrl: { name: 'CustomLiveTuningUrl', type: 'string', category: 'advanced', description: '自定义实时调优文件URL', default: '' },
  customdynamicconfigurl: { name: 'customdynamicconfigurl', type: 'string', category: 'advanced', description: '自定义动态配置文件URL', default: '' },
  EnableFullDump: { name: 'EnableFullDump', type: 'boolean', category: 'advanced', description: '启用完整服务器日志转储', default: false },
  TribeLogDestroyedEnemyStructures: { name: 'TribeLogDestroyedEnemyStructures', type: 'boolean', category: 'advanced', description: '在部落日志中显示敌方建筑破坏', default: false }
}

// 命令行参数（以-开头的参数）
export const commandLineArgs = {
  // 基础服务器设置
  server: { name: 'server', type: 'boolean', category: 'basic', description: '服务器模式', default: true, hidden: true },
  clusterid: { name: 'clusterid', type: 'string', category: 'basic', description: '集群ID', default: '', hidden: true },

  // 性能优化
  allcores: { name: 'allcores', type: 'boolean', category: 'performance', description: '使用所有CPU核心', default: false },
  USEALLAVAILABLECORES: { name: 'USEALLAVAILABLECORES', type: 'boolean', category: 'performance', description: '强制使用所有可用CPU核心', default: false },
  lowmemory: { name: 'lowmemory', type: 'boolean', category: 'performance', description: '低内存模式（节省800MB RAM）', default: false },
  nomemorybias: { name: 'nomemorybias', type: 'boolean', category: 'performance', description: '减少内存使用（约600MB）', default: false },
  usecache: { name: 'usecache', type: 'boolean', category: 'performance', description: '使用缓存（提升70%加载速度）', default: false },
  high: { name: 'high', type: 'boolean', category: 'performance', description: '高优先级进程运行', default: false },
  nocombineclientmoves: { name: 'nocombineclientmoves', type: 'boolean', category: 'performance', description: '禁用客户端移动优化', default: false },

  // 图形优化
  allowansel: { name: 'allowansel', type: 'boolean', category: 'graphics', description: '启用NVIDIA Ansel支持', default: false },
  d3d10: { name: 'd3d10', type: 'boolean', category: 'graphics', description: '强制使用DX10（仅Windows）', default: false },
  sm4: { name: 'sm4', type: 'boolean', category: 'graphics', description: '强制使用Shader Model 4（仅Windows）', default: false },
  nomansky: { name: 'nomansky', type: 'boolean', category: 'graphics', description: '禁用天空细节（云层、星空）', default: false },
  noaafonts: { name: 'noaafonts', type: 'boolean', category: 'graphics', description: '移除字体抗锯齿', default: false },

  // 安全设置
  NoBattlEye: { name: 'NoBattlEye', type: 'boolean', category: 'security', description: '禁用BattleEye反作弊', default: true },
  noantispeedhack: { name: 'noantispeedhack', type: 'boolean', category: 'security', description: '禁用反速度黑客检测', default: false },
  insecure: { name: 'insecure', type: 'boolean', category: 'security', description: '禁用VAC系统', default: false },
  UseItemDupeCheck: { name: 'UseItemDupeCheck', type: 'boolean', category: 'security', description: '启用额外复制保护', default: true },
  exclusivejoin: { name: 'exclusivejoin', type: 'boolean', category: 'security', description: '白名单模式', default: false },

  // 日志记录
  servergamelog: { name: 'servergamelog', type: 'boolean', category: 'logging', description: '启用服务器游戏日志', default: true },
  servergamelogincludetribelogs: { name: 'servergamelogincludetribelogs', type: 'boolean', category: 'logging', description: '在游戏日志中包含部落日志', default: false },
  ServerRCONOutputTribeLogs: { name: 'ServerRCONOutputTribeLogs', type: 'boolean', category: 'logging', description: '在RCON中显示部落聊天', default: false },
  gameplaylogging: { name: 'gameplaylogging', type: 'boolean', category: 'logging', description: '游戏玩法日志记录', default: false },
  webalarm: { name: 'webalarm', type: 'boolean', category: 'logging', description: '启用Web警报', default: false },

  // 模组和平台
  automanagedmods: { name: 'automanagedmods', type: 'boolean', category: 'mods', description: '自动管理MOD下载/安装/更新', default: true },
  crossplay: { name: 'crossplay', type: 'boolean', category: 'mods', description: '启用跨平台（EPIC和Steam）', default: false },
  epiconly: { name: 'epiconly', type: 'boolean', category: 'mods', description: '仅允许Epic玩家连接', default: false },
  UseVivox: { name: 'UseVivox', type: 'boolean', category: 'mods', description: '启用Vivox语音', default: false },
  MapModID: { name: 'MapModID', type: 'string', category: 'mods', description: '自定义地图MOD ID', default: '' },
  PublicIPForEpic: { name: 'PublicIPForEpic', type: 'string', category: 'mods', description: 'Epic客户端的公共IP地址', default: '' },

  // 游戏功能
  ForceAllowCaveFlyers: { name: 'ForceAllowCaveFlyers', type: 'boolean', category: 'features', description: '强制允许飞行生物进入洞穴', default: false },
  ForceRespawnDinos: { name: 'ForceRespawnDinos', type: 'boolean', category: 'features', description: '启动时销毁所有野生生物', default: false },
  NoBiomeWalls: { name: 'NoBiomeWalls', type: 'boolean', category: 'features', description: '消除生物群落变化区域墙效果', default: false },
  PreventHibernation: { name: 'PreventHibernation', type: 'boolean', category: 'features', description: '防止休眠（消耗性能）', default: false },
  StasisKeepControllers: { name: 'StasisKeepControllers', type: 'boolean', category: 'features', description: '在静止状态保持AI控制器', default: false },
  EnableIdlePlayerKick: { name: 'EnableIdlePlayerKick', type: 'boolean', category: 'features', description: '踢出空闲玩家', default: false },
  nofishloot: { name: 'nofishloot', type: 'boolean', category: 'features', description: '禁用钓鱼非肉类战利品', default: false },
  ActiveEvent: { name: 'ActiveEvent', type: 'string', category: 'features', description: '启用指定事件', default: '' },

  // 维护工具
  ClearOldItems: { name: 'ClearOldItems', type: 'boolean', category: 'maintenance', description: '清除旧物品（一次性）', default: false },
  noninlinesaveload: { name: 'noninlinesaveload', type: 'boolean', category: 'maintenance', description: '修复大型存档损坏问题', default: false },
  oldsaveformat: { name: 'oldsaveformat', type: 'boolean', category: 'maintenance', description: '使用旧存档格式', default: false },
  StructureDestructionTag: { name: 'StructureDestructionTag', type: 'string', category: 'maintenance', description: '结构销毁标签', default: '' },
  culture: { name: 'culture', type: 'string', category: 'maintenance', description: '语言代码覆盖', default: '', options: ['ca', 'cs', 'da', 'de', 'en', 'es', 'eu', 'fi', 'fr', 'hu', 'it', 'ja', 'ka', 'ko', 'nl', 'pl', 'pt_BR', 'ru', 'sv', 'th', 'tr', 'zh', 'zh-Hans-CN', 'zh-TW'] }
}

// 参数分类 - 重新设计为更合理的分类
export const categories = {
  basic: { name: '基础设置', description: '服务器基础配置参数' },
  core: { name: '核心游戏', description: '游戏核心功能和显示设置' },
  dinos: { name: '恐龙管理', description: '恐龙驯服、限制、行为相关参数' },
  structures: { name: '建筑系统', description: '建筑限制、腐烂、保护相关参数' },
  pvp: { name: 'PvP设置', description: 'PvP保护、联盟、战斗相关参数' },
  mechanics: { name: '游戏机制', description: '游戏玩法、事件、机制相关参数' },
  transfer: { name: '数据传输', description: '角色、物品、恐龙上传下载控制' },
  performance: { name: '性能优化', description: 'CPU、内存、网络性能优化参数' },
  graphics: { name: '图形优化', description: '图形渲染、显示效果优化参数' },
  security: { name: '安全设置', description: '反作弊、安全控制相关参数' },
  logging: { name: '日志记录', description: '日志记录、调试、监控相关参数' },
  mods: { name: '模组平台', description: '模组管理、跨平台、语音相关参数' },
  features: { name: '游戏功能', description: '特殊游戏功能和行为控制' },
  maintenance: { name: '维护工具', description: '服务器维护、存档、清理工具' },
  advanced: { name: '高级配置', description: '高级自定义配置参数' },
  custom: { name: '自定义参数', description: '用户自定义启动参数' }
}

// 获取参数分类
export function getParamsByCategory() {
  const result = {}
  
  // 初始化分类
  Object.keys(categories).forEach(category => {
    result[category] = []
  })
  
  // 分类查询参数
  Object.values(queryParams).forEach(param => {
    if (!param.hidden) {
      result[param.category].push(param)
    }
  })
  
  // 分类命令行参数
  Object.values(commandLineArgs).forEach(param => {
    if (!param.hidden) {
      result[param.category].push(param)
    }
  })
  
  // 自定义参数分类不包含预定义参数，由前端动态管理
  
  return result
}

// 获取参数默认值
export function getDefaultValues() {
  const defaults = {
    query_params: {},
    command_line_args: {},
    custom_args: []
  }
  
  Object.values(queryParams).forEach(param => {
    if (!param.hidden) {
      // 查询参数的值必须是字符串类型
      if (param.type === 'boolean') {
        defaults.query_params[param.name] = param.default ? 'True' : 'False'
      } else {
        defaults.query_params[param.name] = String(param.default || '')
      }
    }
  })
  
  Object.values(commandLineArgs).forEach(param => {
    if (!param.hidden) {
      // 命令行参数可以保持原始类型
      defaults.command_line_args[param.name] = param.default
    }
  })
  
  return defaults
}