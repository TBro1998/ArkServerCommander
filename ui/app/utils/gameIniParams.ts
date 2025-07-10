// Game.ini 配置参数定义 - TypeScript 版本

// 参数接口类型定义
export interface GameIniParam {
  type: 'boolean' | 'number' | 'text' | 'select'
  default?: any
  min?: number
  max?: number
  step?: number
  options?: { value: string | number, label: string }[]
  description?: string
}

// 分类键类型
export type GameIniCategoryKey = 
  | 'gameBasic'
  | 'experienceSettings'
  | 'breedingSettings'
  | 'itemSettings'
  | 'dinoSettings'
  | 'tribeSettings'
  | 'pvpSettings'
  | 'structureSettings'
  | 'advancedSettings'
  | 'customSettings'

// Game.ini 配置参数定义
export const gameIniParams: Record<GameIniCategoryKey, Record<string, GameIniParam>> = {
  // 基础游戏设置
  gameBasic: {
    bUseSingleplayerSettings: {
      type: 'boolean',
      default: false
    },
    bDisableStructurePlacementCollision: {
      type: 'boolean',
      default: false
    },
    bAllowFlyerCarryPvE: {
      type: 'boolean',
      default: true
    },
    bAllowUnlimitedRespecs: {
      type: 'boolean',
      default: true
    },
    bPassiveDefensesDamageRiderlessDinos: {
      type: 'boolean',
      default: true
    },
    bOnlyAllowSpecifiedEngrams: {
      type: 'boolean',
      default: false
    },
    bAutoUnlockAllEngrams: {
      type: 'boolean',
      default: false
    },
    bShowCreativeMode: {
      type: 'boolean',
      default: false
    },
    bUseCorpseLocator: {
      type: 'boolean',
      default: false
    },
    bDisableLootCrates: {
      type: 'boolean',
      default: false
    },
    bDisableDinoRiding: {
      type: 'boolean',
      default: false
    },
    bDisableDinoTaming: {
      type: 'boolean',
      default: false
    },
    bAllowCustomRecipes: {
      type: 'boolean',
      default: true
    },
    bHardLimitTurretsInRange: {
      type: 'boolean',
      default: false
    },
    bPvEAllowTribeWar: {
      type: 'boolean',
      default: false
    },
    bPvEAllowTribeWarCancel: {
      type: 'boolean',
      default: false
    },
    bPvEDisableFriendlyFire: {
      type: 'boolean',
      default: false
    },
    bDisableFriendlyFire: {
      type: 'boolean',
      default: false
    },
    bFlyerPlatformAllowUnalignedDinoBasing: {
      type: 'boolean',
      default: false
    },
    bIncreasePvPRespawnInterval: {
      type: 'boolean',
      default: false
    },
    bAutoPvETimer: {
      type: 'boolean',
      default: false
    },
    bAutoPvEUseSystemTime: {
      type: 'boolean',
      default: false
    },
    bIgnoreStructuresPreventionVolumes: {
      type: 'boolean',
      default: false
    },
    bGenesisUseStructuresPreventionVolumes: {
      type: 'boolean',
      default: false
    },
    bAllowFlyerSpeedLeveling: {
      type: 'boolean',
      default: false
    },
    bUseTameLimitForStructuresOnly: {
      type: 'boolean',
      default: false
    }
  },

  // 经验值和等级设置
  experienceSettings: {
    OverrideMaxExperiencePointsPlayer: {
      type: 'number',
      default: 0,
      min: 0,
      max: 999999999
    },
    OverrideMaxExperiencePointsDino: {
      type: 'number',
      default: 0,
      min: 0,
      max: 999999999
    },
    OverridePlayerLevelEngramPoints: {
      type: 'number',
      default: 0,
      min: 0,
      max: 1000
    },
    KillXPMultiplier: {
      type: 'number',
      default: 1.0,
      min: 0.1,
      max: 100,
      step: 0.1
    },
    HarvestXPMultiplier: {
      type: 'number',
      default: 1.0,
      min: 0.1,
      max: 100,
      step: 0.1
    },
    CraftXPMultiplier: {
      type: 'number',
      default: 1.0,
      min: 0.1,
      max: 100,
      step: 0.1
    },
    GenericXPMultiplier: {
      type: 'number',
      default: 1.0,
      min: 0.1,
      max: 100,
      step: 0.1
    },
    SpecialXPMultiplier: {
      type: 'number',
      default: 1.0,
      min: 0.1,
      max: 100,
      step: 0.1
    }
  },

  // 繁殖设置
  breedingSettings: {
    MatingIntervalMultiplier: {
      type: 'number',
      default: 1.0,
      min: 0.1,
      max: 10,
      step: 0.1
    },
    MatingSpeedMultiplier: {
      type: 'number',
      default: 1.0,
      min: 0.1,
      max: 10,
      step: 0.1
    },
    EggHatchSpeedMultiplier: {
      type: 'number',
      default: 1.0,
      min: 0.1,
      max: 100,
      step: 0.1
    },
    BabyMatureSpeedMultiplier: {
      type: 'number',
      default: 1.0,
      min: 0.1,
      max: 100,
      step: 0.1
    },
    BabyCuddleIntervalMultiplier: {
      type: 'number',
      default: 1.0,
      min: 0.1,
      max: 10,
      step: 0.1
    },
    BabyFoodConsumptionSpeedMultiplier: {
      type: 'number',
      default: 1.0,
      min: 0.1,
      max: 10,
      step: 0.1
    },
    BabyImprintingStatScaleMultiplier: {
      type: 'number',
      default: 1.0,
      min: 0.1,
      max: 10,
      step: 0.1
    },
    BabyImprintAmountMultiplier: {
      type: 'number',
      default: 1.0,
      min: 0.1,
      max: 10,
      step: 0.1
    },
    BabyCuddleGracePeriodMultiplier: {
      type: 'number',
      default: 1.0,
      min: 0.1,
      max: 10,
      step: 0.1
    },
    BabyCuddleLoseImprintQualitySpeedMultiplier: {
      type: 'number',
      default: 1.0,
      min: 0.1,
      max: 10,
      step: 0.1
    }
  },

  // 物品和资源设置
  itemSettings: {
    CropGrowthSpeedMultiplier: {
      type: 'number',
      default: 1.0,
      min: 0.1,
      max: 100,
      step: 0.1
    },
    CropDecaySpeedMultiplier: {
      type: 'number',
      default: 1.0,
      min: 0.1,
      max: 10,
      step: 0.1
    },
    GlobalSpoilingTimeMultiplier: {
      type: 'number',
      default: 1.0,
      min: 0.1,
      max: 100,
      step: 0.1
    },
    GlobalItemDecompositionTimeMultiplier: {
      type: 'number',
      default: 1.0,
      min: 0.1,
      max: 100,
      step: 0.1
    },
    GlobalCorpseDecompositionTimeMultiplier: {
      type: 'number',
      default: 1.0,
      min: 0.1,
      max: 100,
      step: 0.1
    },
    UseCorpseLifeSpanMultiplier: {
      type: 'number',
      default: 1.0,
      min: 0.1,
      max: 10,
      step: 0.1
    },
    SupplyCrateLootQualityMultiplier: {
      type: 'number',
      default: 1.0,
      min: 0.1,
      max: 10,
      step: 0.1
    },
    FishingLootQualityMultiplier: {
      type: 'number',
      default: 1.0,
      min: 0.1,
      max: 10,
      step: 0.1
    },
    CustomRecipeEffectivenessMultiplier: {
      type: 'number',
      default: 1.0,
      min: 0.1,
      max: 10,
      step: 0.1
    },
    CustomRecipeSkillMultiplier: {
      type: 'number',
      default: 1.0,
      min: 0.1,
      max: 10,
      step: 0.1
    },
    CraftingSkillBonusMultiplier: {
      type: 'number',
      default: 1.0,
      min: 0.1,
      max: 10,
      step: 0.1
    },
    ResourceNoReplenishRadiusPlayers: {
      type: 'number',
      default: 1.0,
      min: 0.1,
      max: 1000,
      step: 0.1
    },
    ResourceNoReplenishRadiusStructures: {
      type: 'number',
      default: 1.0,
      min: 0.1,
      max: 1000,
      step: 0.1
    },
    HarvestResourceItemAmountClassMultipliers: {
      type: 'text',
      default: ''
    },
    DinoHarvestingDamageMultiplier: {
      type: 'number',
      default: 1.0,
      min: 0.1,
      max: 10,
      step: 0.1
    },
    PlayerHarvestingDamageMultiplier: {
      type: 'number',
      default: 1.0,
      min: 0.1,
      max: 10,
      step: 0.1
    }
  },

  // 恐龙设置
  dinoSettings: {
    TamedDinoCharacterFoodDrainMultiplier: {
      type: 'number',
      default: 1.0,
      min: 0.1,
      max: 10,
      step: 0.1
    },
    WildDinoCharacterFoodDrainMultiplier: {
      type: 'number',
      default: 1.0,
      min: 0.1,
      max: 10,
      step: 0.1
    },
    WildDinoTorporDrainMultiplier: {
      type: 'number',
      default: 1.0,
      min: 0.1,
      max: 10,
      step: 0.1
    },
    TamedDinoTorporDrainMultiplier: {
      type: 'number',
      default: 1.0,
      min: 0.1,
      max: 10,
      step: 0.1
    },
    PassiveTameIntervalMultiplier: {
      type: 'number',
      default: 1.0,
      min: 0.1,
      max: 10,
      step: 0.1
    },
    DinoTurretDamageMultiplier: {
      type: 'number',
      default: 1.0,
      min: 0.1,
      max: 10,
      step: 0.1
    },
    PreventDinoTameClassNames: {
      type: 'text',
      default: ''
    },
    PreventTransferForClassName: {
      type: 'text',
      default: ''
    },
    DinoClassDamageMultipliers: {
      type: 'text',
      default: ''
    },
    TamedDinoClassDamageMultipliers: {
      type: 'text',
      default: ''
    },
    DinoClassResistanceMultipliers: {
      type: 'text',
      default: ''
    },
    TamedDinoClassResistanceMultipliers: {
      type: 'text',
      default: ''
    },
    DestroyTamesOverLevelClamp: {
      type: 'number',
      default: 0,
      min: 0,
      max: 1000
    }
  },

  // 部落和玩家设置
  tribeSettings: {
    MaxNumberOfPlayersInTribe: {
      type: 'number',
      default: 0,
      min: 0,
      max: 500
    },
    MaxAlliancesPerTribe: {
      type: 'number',
      default: 0,
      min: 0,
      max: 100
    },
    MaxTribesPerAlliance: {
      type: 'number',
      default: 0,
      min: 0,
      max: 100
    },
    TribeSlotReuseCooldown: {
      type: 'number',
      default: 0,
      min: 0,
      max: 86400
    },
    MaxTribeLogs: {
      type: 'number',
      default: 100,
      min: 1,
      max: 1000
    },
    KickIdlePlayersPeriod: {
      type: 'number',
      default: 0,
      min: 0,
      max: 86400
    }
  },

  // PvP设置
  pvpSettings: {
    IncreasePvPRespawnIntervalCheckPeriod: {
      type: 'number',
      default: 300,
      min: 1,
      max: 3600
    },
    IncreasePvPRespawnIntervalMultiplier: {
      type: 'number',
      default: 2.0,
      min: 1.0,
      max: 10.0,
      step: 0.1
    },
    IncreasePvPRespawnIntervalBaseAmount: {
      type: 'number',
      default: 60,
      min: 1,
      max: 3600
    },
    AutoPvEStartTimeSeconds: {
      type: 'number',
      default: 0,
      min: 0,
      max: 86400
    },
    AutoPvEStopTimeSeconds: {
      type: 'number',
      default: 0,
      min: 0,
      max: 86400
    },
    PvPZoneStructureDamageMultiplier: {
      type: 'number',
      default: 1.0,
      min: 0.1,
      max: 10,
      step: 0.1
    }
  },

  // 建筑和结构设置
  structureSettings: {
    StructureDamageRepairCooldown: {
      type: 'number',
      default: 180,
      min: 0,
      max: 3600
    },
    StructureResistanceMultiplier: {
      type: 'number',
      default: 1.0,
      min: 0.1,
      max: 10,
      step: 0.1
    },
    StructureDamageMultiplier: {
      type: 'number',
      default: 1.0,
      min: 0.1,
      max: 10,
      step: 0.1
    },
    PvEStructureDecayPeriodMultiplier: {
      type: 'number',
      default: 1.0,
      min: 0.1,
      max: 100,
      step: 0.1
    },
    FastDecayInterval: {
      type: 'number',
      default: 0,
      min: 0,
      max: 86400
    }
  },

  // 高级设置
  advancedSettings: {
    HairGrowthSpeedMultiplier: {
      type: 'number',
      default: 1.0,
      min: 0.1,
      max: 10,
      step: 0.1
    },
    PoopIntervalMultiplier: {
      type: 'number',
      default: 1.0,
      min: 0.1,
      max: 10,
      step: 0.1
    },
    LayEggIntervalMultiplier: {
      type: 'number',
      default: 1.0,
      min: 0.1,
      max: 10,
      step: 0.1
    },
    FuelConsumptionIntervalMultiplier: {
      type: 'number',
      default: 1.0,
      min: 0.1,
      max: 10,
      step: 0.1
    },
    GlobalPoweredBatteryDurabilityDecreasePerSecond: {
      type: 'number',
      default: 0,
      min: 0,
      max: 100,
      step: 0.1
    },
    LimitNonPlayerDroppedItemsRange: {
      type: 'number',
      default: 0,
      min: 0,
      max: 1000
    },
    LimitNonPlayerDroppedItemsCount: {
      type: 'number',
      default: 0,
      min: 0,
      max: 10000
    },
    MaxFallSpeedMultiplier: {
      type: 'number',
      default: 1.0,
      min: 0.1,
      max: 10,
      step: 0.1
    },
    PreventOfflinePvPConnectionInvincibleInterval: {
      type: 'number',
      default: 0,
      min: 0,
      max: 3600
    }
  },

  // 自定义配置
  customSettings: {
    ConfigOverrideItemMaxQuantity: {
      type: 'text',
      default: ''
    },
    ConfigOverrideItemCraftingCosts: {
      type: 'text',
      default: ''
    },
    ConfigOverrideSupplyCrateItems: {
      type: 'text',
      default: ''
    },
    ExcludeItemIndices: {
      type: 'text',
      default: ''
    },
    LevelExperienceRampOverrides: {
      type: 'text',
      default: ''
    },
    EngramEntryAutoUnlocks: {
      type: 'text',
      default: ''
    },
    ModIDS: {
      type: 'text',
      default: ''
    }
  }
}

// 根据分类获取参数的辅助函数
export function getParamsByCategory(category: GameIniCategoryKey): { key: string, param: GameIniParam }[] {
  const categoryParams = gameIniParams[category]
  if (!categoryParams) return []
  
  return Object.entries(categoryParams).map(([key, param]) => ({
    key,
    param
  }))
}

// 获取所有参数分类的辅助函数
export function getAllCategories(): { key: GameIniCategoryKey, params: { key: string, param: GameIniParam }[] }[] {
  return Object.keys(gameIniParams).map(key => ({
    key: key as GameIniCategoryKey,
    params: getParamsByCategory(key as GameIniCategoryKey)
  }))
} 