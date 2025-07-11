// GameUserSettings.ini 配置参数定义 - TypeScript 版本
export interface GameUserSettingsParam {
  type: 'boolean' | 'number' | 'text' | 'password'
  default: boolean | number | string
  min?: number
  max?: number
  step?: number
}

export type GameUserSettingsCategoryKey = 
  | 'serverBasic'
  | 'gameMode'
  | 'communication'
  | 'gameMultipliers'
  | 'characterSettings'
  | 'dinoSettings'
  | 'environmentSettings'
  | 'structureSettings'
  | 'tribeSettings'
  | 'breedingSettings'
  | 'itemSettings'
  | 'performanceSettings'
  | 'diseaseSettings'
  | 'offlineRaidSettings'
  | 'crossArkSettings'
  | 'flyerSettings'
  | 'advancedSettings'

export const gameUserSettingsParams: Record<GameUserSettingsCategoryKey, Record<string, GameUserSettingsParam>> = {
  // 服务器基本设置
  serverBasic: {
    ServerPassword: {
      type: 'password',
      default: ''
    },
    SpectatorPassword: {
      type: 'password',
      default: ''
    },
    AdminLogging: {
      type: 'boolean',
      default: true
    }
  },

  // 游戏模式设置
  gameMode: {
    serverPVE: {
      type: 'boolean',
      default: false
    },
    serverHardcore: {
      type: 'boolean',
      default: false
    },
    ShowMapPlayerLocation: {
      type: 'boolean',
      default: false
    },
    allowThirdPersonPlayer: {
      type: 'boolean',
      default: false
    },
    ServerCrosshair: {
      type: 'boolean',
      default: true
    },
    EnablePvPGamma: {
      type: 'boolean',
      default: false
    },
    DisablePvEGamma: {
      type: 'boolean',
      default: false
    },
    serverForceNoHud: {
      type: 'boolean',
      default: false
    },
    ShowFloatingDamageText: {
      type: 'boolean',
      default: false
    },
    AllowHitMarkers: {
      type: 'boolean',
      default: true
    }
  },

  // 聊天和通讯设置
  communication: {
    globalVoiceChat: {
      type: 'boolean',
      default: false
    },
    proximityChat: {
      type: 'boolean',
      default: false
    },
    alwaysNotifyPlayerJoined: {
      type: 'boolean',
      default: false
    },
    alwaysNotifyPlayerLeft: {
      type: 'boolean',
      default: false
    },
    DontAlwaysNotifyPlayerJoined: {
      type: 'boolean',
      default: false
    }
  },

  // 游戏倍率设置
  gameMultipliers: {
    XPMultiplier: {
      type: 'number',
      default: 1.0,
      min: 0.1,
      max: 100,
      step: 0.1
    },
    TamingSpeedMultiplier: {
      type: 'number',
      default: 1.0,
      min: 0.1,
      max: 100,
      step: 0.1
    },
    HarvestAmountMultiplier: {
      type: 'number',
      default: 1.0,
      min: 0.1,
      max: 100,
      step: 0.1
    },
    HarvestHealthMultiplier: {
      type: 'number',
      default: 1.0,
      min: 0.1,
      max: 10,
      step: 0.1
    },
    ResourcesRespawnPeriodMultiplier: {
      type: 'number',
      default: 1.0,
      min: 0.1,
      max: 10,
      step: 0.1
    },
    ItemStackSizeMultiplier: {
      type: 'number',
      default: 1.0,
      min: 1,
      max: 100,
      step: 1
    }
  },

  // 角色设置
  characterSettings: {
    PlayerCharacterHealthRecoveryMultiplier: {
      type: 'number',
      default: 1.0,
      min: 0.1,
      max: 10,
      step: 0.1
    },
    PlayerCharacterFoodDrainMultiplier: {
      type: 'number',
      default: 1.0,
      min: 0.1,
      max: 10,
      step: 0.1
    },
    PlayerCharacterWaterDrainMultiplier: {
      type: 'number',
      default: 1.0,
      min: 0.1,
      max: 10,
      step: 0.1
    },
    PlayerCharacterStaminaDrainMultiplier: {
      type: 'number',
      default: 1.0,
      min: 0.1,
      max: 10,
      step: 0.1
    },
    PlayerDamageMultiplier: {
      type: 'number',
      default: 1.0,
      min: 0.1,
      max: 10,
      step: 0.1
    },
    PlayerResistanceMultiplier: {
      type: 'number',
      default: 1.0,
      min: 0.1,
      max: 10,
      step: 0.1
    },
    OxygenSwimSpeedStatMultiplier: {
      type: 'number',
      default: 1.0,
      min: 0.1,
      max: 10,
      step: 0.1
    },
    ImplantSuicideCD: {
      type: 'number',
      default: 28800,
      min: 0,
      max: 86400
    }
  },

  // 恐龙设置
  dinoSettings: {
    DinoCountMultiplier: {
      type: 'number',
      default: 1.0,
      min: 0.1,
      max: 10,
      step: 0.1
    },
    DinoCharacterHealthRecoveryMultiplier: {
      type: 'number',
      default: 1.0,
      min: 0.1,
      max: 10,
      step: 0.1
    },
    DinoCharacterFoodDrainMultiplier: {
      type: 'number',
      default: 1.0,
      min: 0.1,
      max: 10,
      step: 0.1
    },
    DinoCharacterStaminaDrainMultiplier: {
      type: 'number',
      default: 1.0,
      min: 0.1,
      max: 10,
      step: 0.1
    },
    DinoDamageMultiplier: {
      type: 'number',
      default: 1.0,
      min: 0.1,
      max: 10,
      step: 0.1
    },
    TamedDinoDamageMultiplier: {
      type: 'number',
      default: 1.0,
      min: 0.1,
      max: 10,
      step: 0.1
    },
    DinoResistanceMultiplier: {
      type: 'number',
      default: 1.0,
      min: 0.1,
      max: 10,
      step: 0.1
    },
    TamedDinoResistanceMultiplier: {
      type: 'number',
      default: 1.0,
      min: 0.1,
      max: 10,
      step: 0.1
    },
    MaxTamedDinos: {
      type: 'number',
      default: 4000,
      min: 0,
      max: 50000
    },
    MaxPersonalTamedDinos: {
      type: 'number',
      default: 0,
      min: 0,
      max: 5000
    },
    DisableDinoDecayPvE: {
      type: 'boolean',
      default: false
    },
    AutoDestroyDecayedDinos: {
      type: 'boolean',
      default: false
    },
    PvEDinoDecayPeriodMultiplier: {
      type: 'number',
      default: 1.0,
      min: 0.1,
      max: 10,
      step: 0.1
    },
    PvPDinoDecay: {
      type: 'boolean',
      default: false
    },
    AllowRaidDinoFeeding: {
      type: 'boolean',
      default: false
    },
    RaidDinoCharacterFoodDrainMultiplier: {
      type: 'number',
      default: 1.0,
      min: 0.1,
      max: 10,
      step: 0.1
    },
    AllowFlyerCarryPvE: {
      type: 'boolean',
      default: false
    },
    bForceCanRideFliers: {
      type: 'boolean',
      default: false
    }
  },

  // 环境设置
  environmentSettings: {
    DayCycleSpeedScale: {
      type: 'number',
      default: 1.0,
      min: 0.1,
      max: 10,
      step: 0.1
    },
    DayTimeSpeedScale: {
      type: 'number',
      default: 1.0,
      min: 0.1,
      max: 10,
      step: 0.1
    },
    NightTimeSpeedScale: {
      type: 'number',
      default: 1.0,
      min: 0.1,
      max: 10,
      step: 0.1
    },
    DisableWeatherFog: {
      type: 'boolean',
      default: false
    },
    DifficultyOffset: {
      type: 'number',
      default: 0.2,
      min: 0,
      max: 1,
      step: 0.1
    },
    OverrideOfficialDifficulty: {
      type: 'number',
      default: 5.0,
      min: 1,
      max: 10,
      step: 0.1
    },
    RandomSupplyCratePoints: {
      type: 'boolean',
      default: false
    }
  },

  // 建筑设置
  structureSettings: {
    StructureDamageMultiplier: {
      type: 'number',
      default: 1.0,
      min: 0.1,
      max: 10,
      step: 0.1
    },
    StructureResistanceMultiplier: {
      type: 'number',
      default: 1.0,
      min: 0.1,
      max: 10,
      step: 0.1
    },
    TheMaxStructuresInRange: {
      type: 'number',
      default: 10500,
      min: 1000,
      max: 50000
    },
    NewMaxStructuresInRange: {
      type: 'number',
      default: 6000,
      min: 1000,
      max: 50000
    },
    MaxStructuresInRange: {
      type: 'number',
      default: 1300,
      min: 100,
      max: 10000
    },
    DisableStructureDecayPvE: {
      type: 'boolean',
      default: false
    },
    PvEStructureDecayPeriodMultiplier: {
      type: 'number',
      default: 1.0,
      min: 0.1,
      max: 10,
      step: 0.1
    },
    PvEStructureDecayDestructionPeriod: {
      type: 'number',
      default: 0,
      min: 0,
      max: 1000
    },
    PvPStructureDecay: {
      type: 'boolean',
      default: false
    },
    StructurePickupTimeAfterPlacement: {
      type: 'number',
      default: 30,
      min: 0,
      max: 3600
    },
    StructurePickupHoldDuration: {
      type: 'number',
      default: 0.5,
      min: 0,
      max: 10,
      step: 0.1
    },
    AlwaysAllowStructurePickup: {
      type: 'boolean',
      default: false
    },
    OnlyAutoDestroyCoreStructures: {
      type: 'boolean',
      default: false
    },
    OnlyDecayUnsnappedCoreStructures: {
      type: 'boolean',
      default: false
    },
    FastDecayUnsnappedCoreStructures: {
      type: 'boolean',
      default: false
    },
    DestroyUnconnectedWaterPipes: {
      type: 'boolean',
      default: false
    },
    StructurePreventResourceRadiusMultiplier: {
      type: 'number',
      default: 1.0,
      min: 0.1,
      max: 10,
      step: 0.1
    },
    MaxPlatformSaddleStructureLimit: {
      type: 'number',
      default: 25,
      min: 1,
      max: 1000
    },
    PerPlatformMaxStructuresMultiplier: {
      type: 'number',
      default: 1.0,
      min: 0.1,
      max: 10,
      step: 0.1
    },
    PlatformSaddleBuildAreaBoundsMultiplier: {
      type: 'number',
      default: 1.0,
      min: 0.1,
      max: 10,
      step: 0.1
    },
    OverrideStructurePlatformPrevention: {
      type: 'boolean',
      default: false
    },
    EnableExtraStructurePreventionVolumes: {
      type: 'boolean',
      default: false
    },
    AllowCaveBuildingPvE: {
      type: 'boolean',
      default: false
    },
    AllowCaveBuildingPvP: {
      type: 'boolean',
      default: false
    },
    PvEAllowStructuresAtSupplyDrops: {
      type: 'boolean',
      default: false
    },
    AllowCrateSpawnsOnTopOfStructures: {
      type: 'boolean',
      default: false
    },
    bAllowPlatformSaddleMultiFloors: {
      type: 'boolean',
      default: false
    },
    MaxGateFrameOnSaddles: {
      type: 'number',
      default: 0,
      min: 0,
      max: 100
    }
  },

  // 部落和联盟设置
  tribeSettings: {
    MaxNumberOfPlayersInTribe: {
      type: 'number',
      default: 0,
      min: 0,
      max: 500
    },
    TribeNameChangeCooldown: {
      type: 'number',
      default: 15,
      min: 0,
      max: 1440
    },
    PreventTribeAlliances: {
      type: 'boolean',
      default: false
    },
    MaxAlliancesPerTribe: {
      type: 'number',
      default: 10,
      min: 0,
      max: 100
    },
    MaxTribesPerAlliance: {
      type: 'number',
      default: 10,
      min: 0,
      max: 100
    }
  },

  // 繁殖和印记设置
  breedingSettings: {
    AllowAnyoneBabyImprintCuddle: {
      type: 'boolean',
      default: false
    },
    DisableImprintDinoBuff: {
      type: 'boolean',
      default: false
    },
    BabyImprintingStatScaleMultiplier: {
      type: 'number',
      default: 1.0,
      min: 0.1,
      max: 10,
      step: 0.1
    }
  },

  // 物品和补给设置
  itemSettings: {
    ClampItemSpoilingTimes: {
      type: 'boolean',
      default: false
    },
    ClampResourceHarvestDamage: {
      type: 'boolean',
      default: false
    },
    UseOptimizedHarvestingHealth: {
      type: 'boolean',
      default: true
    },
    BanListURL: {
      type: 'text',
      default: 'http://arkdedicated.com/banlist.txt'
    }
  },

  // 服务器性能设置
  performanceSettings: {
    AutoSavePeriodMinutes: {
      type: 'number',
      default: 15.0,
      min: 1,
      max: 60,
      step: 0.1
    },
    KickIdlePlayersPeriod: {
      type: 'number',
      default: 3600,
      min: 0,
      max: 86400
    },
    ListenServerTetherDistanceMultiplier: {
      type: 'number',
      default: 1.0,
      min: 0.1,
      max: 10,
      step: 0.1
    },
    RCONServerGameLogBuffer: {
      type: 'number',
      default: 600,
      min: 100,
      max: 10000
    },
    NPCNetworkStasisRangeScalePlayerCountStart: {
      type: 'number',
      default: 0,
      min: 0,
      max: 1000
    },
    NPCNetworkStasisRangeScalePlayerCountEnd: {
      type: 'number',
      default: 0,
      min: 0,
      max: 1000
    },
    NPCNetworkStasisRangeScalePercentEnd: {
      type: 'number',
      default: 0,
      min: 0,
      max: 100,
      step: 0.1
    }
  },

  // 疾病和状态设置
  diseaseSettings: {
    PreventDiseases: {
      type: 'boolean',
      default: false
    },
    NonPermanentDiseases: {
      type: 'boolean',
      default: false
    },
    PreventSpawnAnimations: {
      type: 'boolean',
      default: false
    }
  },

  // 离线突袭保护设置
  offlineRaidSettings: {
    PreventOfflinePvP: {
      type: 'boolean',
      default: false
    },
    PreventOfflinePvPInterval: {
      type: 'number',
      default: 900,
      min: 0,
      max: 3600
    }
  },

  // 跨服传输设置
  crossArkSettings: {
    NoTributeDownloads: {
      type: 'boolean',
      default: false
    },
    PreventDownloadSurvivors: {
      type: 'boolean',
      default: false
    },
    PreventDownloadItems: {
      type: 'boolean',
      default: false
    },
    PreventDownloadDinos: {
      type: 'boolean',
      default: false
    },
    PreventUploadSurvivors: {
      type: 'boolean',
      default: false
    },
    PreventUploadItems: {
      type: 'boolean',
      default: false
    },
    PreventUploadDinos: {
      type: 'boolean',
      default: false
    },
    CrossARKAllowForeignDinoDownloads: {
      type: 'boolean',
      default: false
    },
    MaxTributeDinos: {
      type: 'number',
      default: 20,
      min: 1,
      max: 500
    },
    MaxTributeItems: {
      type: 'number',
      default: 50,
      min: 1,
      max: 500
    },
    MinimumDinoReuploadInterval: {
      type: 'number',
      default: 0,
      min: 0,
      max: 86400
    },
    TributeItemExpirationSeconds: {
      type: 'number',
      default: 86400,
      min: 0,
      max: 604800
    },
    TributeDinoExpirationSeconds: {
      type: 'number',
      default: 86400,
      min: 0,
      max: 604800
    },
    TributeCharacterExpirationSeconds: {
      type: 'number',
      default: 86400,
      min: 0,
      max: 604800
    }
  },

  // 飞行载具设置
  flyerSettings: {
    AllowFlyingStaminaRecovery: {
      type: 'boolean',
      default: true
    },
    ForceFlyerExplosives: {
      type: 'boolean',
      default: false
    }
  },

  // 高级功能设置
  advancedSettings: {
    AllowMultipleAttachedC4: {
      type: 'boolean',
      default: false
    },
    AllowIntegratedSPlusStructures: {
      type: 'boolean',
      default: true
    },
    AllowHideDamageSourceFromLogs: {
      type: 'boolean',
      default: false
    },
    AllowSharedConnections: {
      type: 'boolean',
      default: false
    },
    bFilterTribeNames: {
      type: 'boolean',
      default: false
    },
    bFilterCharacterNames: {
      type: 'boolean',
      default: false
    },
    bFilterChat: {
      type: 'boolean',
      default: false
    },
    EnableCryoSicknessPVE: {
      type: 'boolean',
      default: true
    },
    EnableCryopodNerf: {
      type: 'boolean',
      default: false
    },
    CryopodNerfDuration: {
      type: 'number',
      default: 10,
      min: 0,
      max: 3600
    },
    CryopodNerfDamageMult: {
      type: 'number',
      default: 0.01,
      min: 0,
      max: 1,
      step: 0.01
    },
    CryopodNerfIncomingDamageMultPercent: {
      type: 'number',
      default: 0.25,
      min: 0,
      max: 1,
      step: 0.01
    },
    DisableCryopodEnemyCheck: {
      type: 'boolean',
      default: false
    },
    DisableCryopodFridgeRequirement: {
      type: 'boolean',
      default: false
    },
    AllowCryoFridgeOnSaddle: {
      type: 'boolean',
      default: false
    },
    MaxTrainCars: {
      type: 'number',
      default: 8,
      min: 1,
      max: 20
    },
    MaxHexagonsPerCharacter: {
      type: 'number',
      default: 2000000000,
      min: 0,
      max: 2000000000
    },
    AllowTekSuitPowersInGenesis: {
      type: 'boolean',
      default: false
    },
    CustomDynamicConfigUrl: {
      type: 'text',
      default: ''
    }
  }
}

// GameUserSettings.ini 辅助函数
export function getGameUserSettingsParamsByCategory(categoryKey: GameUserSettingsCategoryKey): Record<string, GameUserSettingsParam> {
  return gameUserSettingsParams[categoryKey] || {}
}

export function getAllGameUserSettingsCategories(): GameUserSettingsCategoryKey[] {
  return Object.keys(gameUserSettingsParams) as GameUserSettingsCategoryKey[]
} 