// ARK 服务器配置参数定义
export const arkConfigParams = {
  // GameUserSettings.ini 参数配置
  gameUserSettings: {
    // 服务器基本设置
    serverBasic: {
      title: '服务器基本设置',
      params: {
        SessionName: {
          type: 'text',
          label: '服务器名称',
          default: 'ARK Server',
          description: '在服务器列表中显示的服务器名称'
        },
        ServerPassword: {
          type: 'password',
          label: '服务器密码',
          default: '',
          description: '玩家加入服务器需要的密码，留空表示无密码'
        },
        ServerAdminPassword: {
          type: 'password',
          label: '管理员密码',
          default: '',
          description: '获得管理员权限所需的密码'
        },
        Port: {
          type: 'number',
          label: '游戏端口',
          default: 7777,
          min: 1024,
          max: 65535,
          description: '玩家连接到服务器使用的端口'
        },
        QueryPort: {
          type: 'number',
          label: '查询端口',
          default: 27015,
          min: 1024,
          max: 65535,
          description: 'Steam 查询端口，用于在服务器列表中显示'
        },
        RCONPort: {
          type: 'number',
          label: 'RCON端口',
          default: 32330,
          min: 1024,
          max: 65535,
          description: '远程控制台访问端口'
        },
        MaxPlayers: {
          type: 'number',
          label: '最大玩家数',
          default: 70,
          min: 1,
          max: 999,
          description: '服务器同时允许的最大玩家数量'
        }
      }
    },

    // 游戏模式设置
    gameMode: {
      title: '游戏模式设置',
      params: {
        ServerPVE: {
          type: 'boolean',
          label: 'PvE模式',
          default: false,
          description: '启用PvE模式，禁用玩家间伤害'
        },
        ServerHardcore: {
          type: 'boolean',
          label: '硬核模式',
          default: false,
          description: '玩家死亡后重新开始，等级重置为1'
        },
        ShowMapPlayerLocation: {
          type: 'boolean',
          label: '显示玩家位置',
          default: true,
          description: '在地图上显示玩家自己的位置'
        },
        AllowThirdPersonPlayer: {
          type: 'boolean',
          label: '允许第三人称',
          default: true,
          description: '允许玩家使用第三人称视角'
        },
        ServerCrosshair: {
          type: 'boolean',
          label: '显示准星',
          default: true,
          description: '显示武器瞄准准星'
        },
        EnablePVPGamma: {
          type: 'boolean',
          label: 'PvP伽马调节',
          default: false,
          description: '在PvP模式下允许使用gamma控制台命令'
        }
      }
    },

    // 游戏倍率设置
    gameMultipliers: {
      title: '游戏倍率设置',
      params: {
        XPMultiplier: {
          type: 'number',
          label: '经验倍率',
          default: 1.0,
          min: 0.1,
          max: 100,
          step: 0.1,
          description: '玩家、部落和恐龙获得经验的倍率'
        },
        TamingSpeedMultiplier: {
          type: 'number',
          label: '驯服速度倍率',
          default: 1.0,
          min: 0.1,
          max: 100,
          step: 0.1,
          description: '恐龙驯服速度倍率，值越高驯服越快'
        },
        HarvestAmountMultiplier: {
          type: 'number',
          label: '采集倍率',
          default: 1.0,
          min: 0.1,
          max: 100,
          step: 0.1,
          description: '所有采集活动的产量倍率'
        },
        HarvestHealthMultiplier: {
          type: 'number',
          label: '资源血量倍率',
          default: 1.0,
          min: 0.1,
          max: 10,
          step: 0.1,
          description: '可采集物品(树木、岩石等)的血量倍率'
        },
        ResourcesRespawnPeriodMultiplier: {
          type: 'number',
          label: '资源重生倍率',
          default: 1.0,
          min: 0.1,
          max: 10,
          step: 0.1,
          description: '资源重生速度倍率，值越低重生越快'
        },
        ItemStackSizeMultiplier: {
          type: 'number',
          label: '物品堆叠倍率',
          default: 1.0,
          min: 1,
          max: 100,
          step: 1,
          description: '物品最大堆叠数量倍率'
        }
      }
    },

    // 角色设置
    characterSettings: {
      title: '角色设置',
      params: {
        PlayerCharacterHealthRecoveryMultiplier: {
          type: 'number',
          label: '玩家回血倍率',
          default: 1.0,
          min: 0.1,
          max: 10,
          step: 0.1,
          description: '玩家角色血量恢复速度倍率'
        },
        PlayerCharacterFoodDrainMultiplier: {
          type: 'number',
          label: '玩家饥饿倍率',
          default: 1.0,
          min: 0.1,
          max: 10,
          step: 0.1,
          description: '玩家角色食物消耗速度倍率'
        },
        PlayerCharacterWaterDrainMultiplier: {
          type: 'number',
          label: '玩家口渴倍率',
          default: 1.0,
          min: 0.1,
          max: 10,
          step: 0.1,
          description: '玩家角色水分消耗速度倍率'
        },
        PlayerCharacterStaminaDrainMultiplier: {
          type: 'number',
          label: '玩家耐力倍率',
          default: 1.0,
          min: 0.1,
          max: 10,
          step: 0.1,
          description: '玩家角色耐力消耗速度倍率'
        },
        PlayerDamageMultiplier: {
          type: 'number',
          label: '玩家伤害倍率',
          default: 1.0,
          min: 0.1,
          max: 10,
          step: 0.1,
          description: '玩家造成伤害的倍率'
        },
        PlayerResistanceMultiplier: {
          type: 'number',
          label: '玩家抗性倍率',
          default: 1.0,
          min: 0.1,
          max: 10,
          step: 0.1,
          description: '玩家受到伤害的抗性倍率，值越高受到伤害越多'
        }
      }
    },

    // 恐龙设置
    dinoSettings: {
      title: '恐龙设置',
      params: {
        DinoCountMultiplier: {
          type: 'number',
          label: '恐龙数量倍率',
          default: 1.0,
          min: 0.1,
          max: 10,
          step: 0.1,
          description: '恐龙生成数量倍率'
        },
        DinoCharacterHealthRecoveryMultiplier: {
          type: 'number',
          label: '恐龙回血倍率',
          default: 1.0,
          min: 0.1,
          max: 10,
          step: 0.1,
          description: '恐龙血量恢复速度倍率'
        },
        DinoCharacterFoodDrainMultiplier: {
          type: 'number',
          label: '恐龙饥饿倍率',
          default: 1.0,
          min: 0.1,
          max: 10,
          step: 0.1,
          description: '恐龙食物消耗速度倍率'
        },
        DinoCharacterStaminaDrainMultiplier: {
          type: 'number',
          label: '恐龙耐力倍率',
          default: 1.0,
          min: 0.1,
          max: 10,
          step: 0.1,
          description: '恐龙耐力消耗速度倍率'
        },
        DinoDamageMultiplier: {
          type: 'number',
          label: '恐龙伤害倍率',
          default: 1.0,
          min: 0.1,
          max: 10,
          step: 0.1,
          description: '野生恐龙造成伤害的倍率'
        },
        TamedDinoDamageMultiplier: {
          type: 'number',
          label: '驯服恐龙伤害倍率',
          default: 1.0,
          min: 0.1,
          max: 10,
          step: 0.1,
          description: '驯服恐龙造成伤害的倍率'
        },
        DinoResistanceMultiplier: {
          type: 'number',
          label: '恐龙抗性倍率',
          default: 1.0,
          min: 0.1,
          max: 10,
          step: 0.1,
          description: '野生恐龙受到伤害的抗性倍率'
        },
        TamedDinoResistanceMultiplier: {
          type: 'number',
          label: '驯服恐龙抗性倍率',
          default: 1.0,
          min: 0.1,
          max: 10,
          step: 0.1,
          description: '驯服恐龙受到伤害的抗性倍率'
        },
        MaxTamedDinos: {
          type: 'number',
          label: '最大驯服恐龙数',
          default: 5000,
          min: 0,
          max: 50000,
          description: '服务器全局驯服恐龙数量上限'
        },
        MaxPersonalTamedDinos: {
          type: 'number',
          label: '个人最大驯服数',
          default: 0,
          min: 0,
          max: 5000,
          description: '每个部落可驯服的恐龙数量上限，0表示无限制'
        }
      }
    },

    // 环境设置
    environmentSettings: {
      title: '环境设置',
      params: {
        DayCycleSpeedScale: {
          type: 'number',
          label: '昼夜循环速度',
          default: 1.0,
          min: 0.1,
          max: 10,
          step: 0.1,
          description: '昼夜循环速度倍率，值越高时间流逝越快'
        },
        DayTimeSpeedScale: {
          type: 'number',
          label: '白天时长倍率',
          default: 1.0,
          min: 0.1,
          max: 10,
          step: 0.1,
          description: '白天持续时间倍率，值越低白天越长'
        },
        NightTimeSpeedScale: {
          type: 'number',
          label: '夜晚时长倍率',
          default: 1.0,
          min: 0.1,
          max: 10,
          step: 0.1,
          description: '夜晚持续时间倍率，值越低夜晚越长'
        },
        DisableWeatherFog: {
          type: 'boolean',
          label: '禁用雾效',
          default: false,
          description: '禁用游戏中的雾效效果'
        },
        DifficultyOffset: {
          type: 'number',
          label: '难度偏移',
          default: 1.0,
          min: 0,
          max: 1,
          step: 0.1,
          description: '游戏难度偏移，影响恐龙等级'
        },
        OverrideOfficialDifficulty: {
          type: 'number',
          label: '覆盖官方难度',
          default: 5.0,
          min: 1,
          max: 10,
          step: 0.1,
          description: '覆盖默认难度级别，5.0匹配官方服务器'
        }
      }
    },

    // 建筑设置
    structureSettings: {
      title: '建筑设置',
      params: {
        StructureDamageMultiplier: {
          type: 'number',
          label: '建筑伤害倍率',
          default: 1.0,
          min: 0.1,
          max: 10,
          step: 0.1,
          description: '建筑物造成伤害的倍率'
        },
        StructureResistanceMultiplier: {
          type: 'number',
          label: '建筑抗性倍率',
          default: 1.0,
          min: 0.1,
          max: 10,
          step: 0.1,
          description: '建筑物受到伤害的抗性倍率'
        },
        TheMaxStructuresInRange: {
          type: 'number',
          label: '范围内最大建筑数',
          default: 10500,
          min: 1000,
          max: 50000,
          description: '特定范围内可建造的最大建筑数量'
        },
        DisableStructureDecayPvE: {
          type: 'boolean',
          label: '禁用PvE建筑腐朽',
          default: false,
          description: '在PvE模式下禁用建筑物自动腐朽'
        },
        PvEStructureDecayPeriodMultiplier: {
          type: 'number',
          label: 'PvE建筑腐朽倍率',
          default: 1.0,
          min: 0.1,
          max: 10,
          step: 0.1,
          description: 'PvE模式下建筑腐朽速度倍率'
        },
        StructurePickupTimeAfterPlacement: {
          type: 'number',
          label: '建筑快捷拾取时间',
          default: 30,
          min: 0,
          max: 3600,
          description: '建筑放置后可快捷拾取的时间(秒)'
        },
        AlwaysAllowStructurePickup: {
          type: 'boolean',
          label: '总是允许拾取建筑',
          default: false,
          description: '禁用建筑快捷拾取的时间限制'
        }
      }
    },

    // RCON 和管理设置
    adminSettings: {
      title: 'RCON和管理设置',
      params: {
        RCONEnabled: {
          type: 'boolean',
          label: '启用RCON',
          default: true,
          description: '启用远程控制台访问'
        },
        AdminLogging: {
          type: 'boolean',
          label: '管理员日志',
          default: false,
          description: '将所有管理员命令记录到游戏聊天中'
        },
        SpectatorPassword: {
          type: 'password',
          label: '观察者密码',
          default: '',
          description: '非管理员观察模式密码'
        }
      }
    }
  },

  // Game.ini 参数配置
  gameIni: {
    // 基础游戏设置
    gameBasic: {
      title: '基础游戏设置',
      params: {
        bUseSingleplayerSettings: {
          type: 'boolean',
          label: '使用单人设置',
          default: false,
          description: '使用单人游戏的平衡设置'
        },
        bDisableStructurePlacementCollision: {
          type: 'boolean',
          label: '禁用建筑碰撞',
          default: false,
          description: '允许建筑物重叠放置'
        },
        bAllowFlyerCarryPvE: {
          type: 'boolean',
          label: 'PvE飞行载具',
          default: true,
          description: 'PvE模式下允许飞行生物抓取野生恐龙'
        },
        bAllowUnlimitedRespecs: {
          type: 'boolean',
          label: '无限重置技能',
          default: true,
          description: '允许无限制重置角色技能点'
        },
        bPassiveDefensesDamageRiderlessDinos: {
          type: 'boolean',
          label: '防御伤害无骑手恐龙',
          default: true,
          description: '尖刺墙等防御设施对无骑手恐龙造成伤害'
        }
      }
    },

    // 繁殖设置
    breedingSettings: {
      title: '繁殖设置',
      params: {
        MatingIntervalMultiplier: {
          type: 'number',
          label: '交配间隔倍率',
          default: 1.0,
          min: 0.1,
          max: 10,
          step: 0.1,
          description: '恐龙交配冷却时间倍率，值越高冷却越长'
        },
        EggHatchSpeedMultiplier: {
          type: 'number',
          label: '蛋孵化速度倍率',
          default: 1.0,
          min: 0.1,
          max: 100,
          step: 0.1,
          description: '恐龙蛋孵化速度倍率，值越高孵化越快'
        },
        BabyMatureSpeedMultiplier: {
          type: 'number',
          label: '幼体成长速度倍率',
          default: 1.0,
          min: 0.1,
          max: 100,
          step: 0.1,
          description: '幼体恐龙成长速度倍率，值越高成长越快'
        },
        BabyCuddleIntervalMultiplier: {
          type: 'number',
          label: '幼体关爱间隔倍率',
          default: 1.0,
          min: 0.1,
          max: 10,
          step: 0.1,
          description: '幼体恐龙关爱间隔时间倍率'
        },
        BabyFoodConsumptionSpeedMultiplier: {
          type: 'number',
          label: '幼体食物消耗倍率',
          default: 1.0,
          min: 0.1,
          max: 10,
          step: 0.1,
          description: '幼体恐龙食物消耗速度倍率'
        }
      }
    },

    // 物品设置
    itemSettings: {
      title: '物品设置',
      params: {
        CropGrowthSpeedMultiplier: {
          type: 'number',
          label: '作物生长速度倍率',
          default: 1.0,
          min: 0.1,
          max: 100,
          step: 0.1,
          description: '农作物生长速度倍率'
        },
        GlobalSpoilingTimeMultiplier: {
          type: 'number',
          label: '全局腐坏时间倍率',
          default: 1.0,
          min: 0.1,
          max: 100,
          step: 0.1,
          description: '所有物品腐坏时间倍率，值越高腐坏越慢'
        },
        GlobalItemDecompositionTimeMultiplier: {
          type: 'number',
          label: '物品消失时间倍率',
          default: 1.0,
          min: 0.1,
          max: 100,
          step: 0.1,
          description: '掉落物品消失时间倍率'
        },
        GlobalCorpseDecompositionTimeMultiplier: {
          type: 'number',
          label: '尸体消失时间倍率',
          default: 1.0,
          min: 0.1,
          max: 100,
          step: 0.1,
          description: '玩家和恐龙尸体消失时间倍率'
        }
      }
    },

    // 高级设置
    advancedSettings: {
      title: '高级设置',
      params: {
        SupplyCrateLootQualityMultiplier: {
          type: 'number',
          label: '补给箱品质倍率',
          default: 1.0,
          min: 0.1,
          max: 10,
          step: 0.1,
          description: '补给箱物品品质倍率'
        },
        FishingLootQualityMultiplier: {
          type: 'number',
          label: '钓鱼品质倍率',
          default: 1.0,
          min: 0.1,
          max: 10,
          step: 0.1,
          description: '钓鱼获得物品品质倍率'
        },
        MaxNumberOfPlayersInTribe: {
          type: 'number',
          label: '部落最大人数',
          default: 0,
          min: 0,
          max: 500,
          description: '单个部落最大玩家数量，0表示无限制'
        },
        bAllowCustomRecipes: {
          type: 'boolean',
          label: '允许自定义配方',
          default: true,
          description: '允许玩家创建自定义烹饪配方'
        },
        UseCorpseLifeSpanMultiplier: {
          type: 'number',
          label: '尸体生命周期倍率',
          default: 1.0,
          min: 0.1,
          max: 10,
          step: 0.1,
          description: '尸体存在时间倍率'
        }
      }
    }
  }
} 