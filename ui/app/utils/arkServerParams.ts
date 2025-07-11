// ARK服务器启动参数定义 - TypeScript版本支持i18n
// 参数类型定义
interface ServerParam {
  type: 'boolean' | 'number' | 'string' | 'select'
  category: string
  default?: any
  min?: number
  max?: number
  step?: number
  options?: string[]
  hidden?: boolean
}

// 查询参数（以?开头的参数）
export const queryParams: Record<string, ServerParam> = {
  // 基础参数（从服务器配置自动生成，不在UI中显示）
  listen: { type: 'boolean', category: 'basic', default: true, hidden: true },
  Port: { type: 'number', category: 'basic', default: 7777, hidden: true },
  QueryPort: { type: 'number', category: 'basic', default: 27015, hidden: true },
  RCONPort: { type: 'number', category: 'basic', default: 32330, hidden: true },
  ServerAdminPassword: { type: 'string', category: 'basic', default: '', hidden: true },
  SessionName: { type: 'string', category: 'basic', default: '', hidden: true },
  GameModIds: { type: 'string', category: 'basic', default: '', hidden: true },

  // 核心游戏设置
  OverrideOfficialDifficulty: { type: 'number', category: 'core', default: 4.0, min: 1.0, max: 10.0, step: 0.1 },
  AllowThirdPersonPlayer: { type: 'boolean', category: 'core', default: true },
  ServerCrosshair: { type: 'boolean', category: 'core', default: false },
  ShowFloatingDamageText: { type: 'boolean', category: 'core', default: false },
  MapPlayerLocation: { type: 'boolean', category: 'core', default: false },
  DisableWeatherFog: { type: 'boolean', category: 'core', default: false },

  // 恐龙管理
  MaxPersonalTamedDinos: { type: 'number', category: 'dinos', default: 500, min: 1, max: 10000 },
  MaxTamedDinos: { type: 'number', category: 'dinos', default: 4000, min: 1, max: 100000 },
  PersonalTamedDinosSaddleStructureCost: { type: 'number', category: 'dinos', default: 19, min: 1, max: 100 },
  MinimumDinoReuploadInterval: { type: 'number', category: 'dinos', default: 0, min: 0, max: 86400 },
  ServerAutoForceRespawnWildDinosInterval: { type: 'number', category: 'dinos', default: 0, min: 0, max: 31536000 },
  AllowAnyoneBabyImprintCuddle: { type: 'boolean', category: 'dinos', default: false },
  DisableImprintDinoBuff: { type: 'boolean', category: 'dinos', default: false },
  PreventMateBoost: { type: 'boolean', category: 'dinos', default: false },
  AllowFlyerCarryPvE: { type: 'boolean', category: 'dinos', default: false },
  AllowFlyingStaminaRecovery: { type: 'boolean', category: 'dinos', default: false },
  ForceFlyerExplosives: { type: 'boolean', category: 'dinos', default: false },
  AllowMultipleAttachedC4: { type: 'boolean', category: 'dinos', default: false },
  AutoDestroyDecayedDinos: { type: 'boolean', category: 'dinos', default: false },

  // 建筑系统
  TheMaxStructuresInRange: { type: 'number', category: 'structures', default: 10500, min: 1000, max: 50000 },
  ForceAllStructureLocking: { type: 'boolean', category: 'structures', default: false },
  EnableExtraStructurePreventionVolumes: { type: 'boolean', category: 'structures', default: false },
  OverrideStructurePlatformPrevention: { type: 'boolean', category: 'structures', default: false },
  DestroyUnconnectedWaterPipes: { type: 'boolean', category: 'structures', default: false },
  FastDecayUnsnappedCoreStructures: { type: 'boolean', category: 'structures', default: false },
  OnlyAutoDestroyCoreStructures: { type: 'boolean', category: 'structures', default: false },
  OnlyDecayUnsnappedCoreStructures: { type: 'boolean', category: 'structures', default: false },

  // PvP设置
  PreventOfflinePvP: { type: 'boolean', category: 'pvp', default: false },
  PreventOfflinePvPInterval: { type: 'number', category: 'pvp', default: 900, min: 0, max: 86400 },
  PreventTribeAlliances: { type: 'boolean', category: 'pvp', default: false },

  // 游戏机制
  AllowCrateSpawnsOnTopOfStructures: { type: 'boolean', category: 'mechanics', default: false },
  RandomSupplyCratePoints: { type: 'boolean', category: 'mechanics', default: false },
  NonPermanentDiseases: { type: 'boolean', category: 'mechanics', default: false },
  PreventDiseases: { type: 'boolean', category: 'mechanics', default: false },
  PreventSpawnAnimations: { type: 'boolean', category: 'mechanics', default: false },
  UseOptimizedHarvestingHealth: { type: 'boolean', category: 'mechanics', default: false },
  ClampItemSpoilingTimes: { type: 'boolean', category: 'mechanics', default: false },
  ClampItemStats: { type: 'boolean', category: 'mechanics', default: false },
  ExtinctionEventTimeInterval: { type: 'number', category: 'mechanics', default: 0, min: 0, max: 31536000 },
  OxygenSwimSpeedStatMultiplier: { type: 'number', category: 'mechanics', default: 1.0, min: 0.1, max: 10.0, step: 0.1 },

  // 数据传输控制
  PreventDownloadSurvivors: { type: 'boolean', category: 'transfer', default: false },
  PreventDownloadItems: { type: 'boolean', category: 'transfer', default: false },
  PreventDownloadDinos: { type: 'boolean', category: 'transfer', default: false },
  PreventUploadSurvivors: { type: 'boolean', category: 'transfer', default: false },
  PreventUploadItems: { type: 'boolean', category: 'transfer', default: false },
  PreventUploadDinos: { type: 'boolean', category: 'transfer', default: false },

  // 高级配置
  CustomLiveTuningUrl: { type: 'string', category: 'advanced', default: '' },
  customdynamicconfigurl: { type: 'string', category: 'advanced', default: '' },
  EnableFullDump: { type: 'boolean', category: 'advanced', default: false },
  TribeLogDestroyedEnemyStructures: { type: 'boolean', category: 'advanced', default: false }
}

// 命令行参数（以-开头的参数）
export const commandLineArgs: Record<string, ServerParam> = {
  // 基础服务器设置
  server: { type: 'boolean', category: 'basic', default: true, hidden: true },
  clusterid: { type: 'string', category: 'basic', default: '', hidden: true },

  // 性能优化
  allcores: { type: 'boolean', category: 'performance', default: false },
  USEALLAVAILABLECORES: { type: 'boolean', category: 'performance', default: false },
  lowmemory: { type: 'boolean', category: 'performance', default: false },
  nomemorybias: { type: 'boolean', category: 'performance', default: false },
  usecache: { type: 'boolean', category: 'performance', default: false },
  high: { type: 'boolean', category: 'performance', default: false },
  nocombineclientmoves: { type: 'boolean', category: 'performance', default: false },

  // 图形优化
  allowansel: { type: 'boolean', category: 'graphics', default: false },
  d3d10: { type: 'boolean', category: 'graphics', default: false },
  sm4: { type: 'boolean', category: 'graphics', default: false },
  nomansky: { type: 'boolean', category: 'graphics', default: false },
  noaafonts: { type: 'boolean', category: 'graphics', default: false },

  // 安全设置
  NoBattlEye: { type: 'boolean', category: 'security', default: true },
  noantispeedhack: { type: 'boolean', category: 'security', default: false },
  insecure: { type: 'boolean', category: 'security', default: false },
  UseItemDupeCheck: { type: 'boolean', category: 'security', default: true },
  exclusivejoin: { type: 'boolean', category: 'security', default: false },

  // 日志记录
  servergamelog: { type: 'boolean', category: 'logging', default: true },
  servergamelogincludetribelogs: { type: 'boolean', category: 'logging', default: false },
  ServerRCONOutputTribeLogs: { type: 'boolean', category: 'logging', default: false },
  gameplaylogging: { type: 'boolean', category: 'logging', default: false },
  webalarm: { type: 'boolean', category: 'logging', default: false },

  // 模组和平台
  automanagedmods: { type: 'boolean', category: 'mods', default: true },
  crossplay: { type: 'boolean', category: 'mods', default: false },
  epiconly: { type: 'boolean', category: 'mods', default: false },
  UseVivox: { type: 'boolean', category: 'mods', default: false },
  MapModID: { type: 'string', category: 'mods', default: '' },
  PublicIPForEpic: { type: 'string', category: 'mods', default: '' },

  // 游戏功能
  ForceAllowCaveFlyers: { type: 'boolean', category: 'features', default: false },
  ForceRespawnDinos: { type: 'boolean', category: 'features', default: false },
  NoBiomeWalls: { type: 'boolean', category: 'features', default: false },
  PreventHibernation: { type: 'boolean', category: 'features', default: false },
  StasisKeepControllers: { type: 'boolean', category: 'features', default: false },
  EnableIdlePlayerKick: { type: 'boolean', category: 'features', default: false },
  nofishloot: { type: 'boolean', category: 'features', default: false },
  ActiveEvent: { type: 'string', category: 'features', default: '' },

  // 维护工具
  ClearOldItems: { type: 'boolean', category: 'maintenance', default: false },
  noninlinesaveload: { type: 'boolean', category: 'maintenance', default: false },
  oldsaveformat: { type: 'boolean', category: 'maintenance', default: false },
  StructureDestructionTag: { type: 'string', category: 'maintenance', default: '' },
  culture: { 
    type: 'select', 
    category: 'maintenance', 
    default: '', 
    options: ['ca', 'cs', 'da', 'de', 'en', 'es', 'eu', 'fi', 'fr', 'hu', 'it', 'ja', 'ka', 'ko', 'nl', 'pl', 'pt_BR', 'ru', 'sv', 'th', 'tr', 'zh', 'zh-Hans-CN', 'zh-TW'] 
  }
}

// 参数分类定义
export const categories = {
  basic: 'basic',
  core: 'core',
  dinos: 'dinos',
  structures: 'structures',
  pvp: 'pvp',
  mechanics: 'mechanics',
  transfer: 'transfer',
  performance: 'performance',
  graphics: 'graphics',
  security: 'security',
  logging: 'logging',
  mods: 'mods',
  features: 'features',
  maintenance: 'maintenance',
  advanced: 'advanced',
  custom: 'custom'
} as const

export type CategoryKey = keyof typeof categories

// 获取服务器启动参数分类
export function getServerParamsByCategory() {
  const result: Record<CategoryKey, Array<{ key: string, param: ServerParam }>> = {
    basic: [],
    core: [],
    dinos: [],
    structures: [],
    pvp: [],
    mechanics: [],
    transfer: [],
    performance: [],
    graphics: [],
    security: [],
    logging: [],
    mods: [],
    features: [],
    maintenance: [],
    advanced: [],
    custom: []
  }
  
  // 分类查询参数
  Object.entries(queryParams).forEach(([key, param]) => {
    if (!param.hidden) {
      result[param.category as CategoryKey].push({ key, param })
    }
  })
  
  // 分类命令行参数
  Object.entries(commandLineArgs).forEach(([key, param]) => {
    if (!param.hidden) {
      result[param.category as CategoryKey].push({ key, param })
    }
  })
  
  return result
}

// 获取参数默认值
export function getDefaultValues() {
  const defaults = {
    query_params: {} as Record<string, string>,
    command_line_args: {} as Record<string, any>,
    custom_args: [] as string[]
  }
  
  Object.entries(queryParams).forEach(([key, param]) => {
    if (!param.hidden) {
      // 查询参数的值必须是字符串类型
      if (param.type === 'boolean') {
        defaults.query_params[key] = param.default ? 'True' : 'False'
      } else {
        defaults.query_params[key] = String(param.default || '')
      }
    }
  })
  
  Object.entries(commandLineArgs).forEach(([key, param]) => {
    if (!param.hidden) {
      // 命令行参数可以保持原始类型
      defaults.command_line_args[key] = param.default
    }
  })
  
  return defaults
}

// 导出类型
export type { ServerParam } 