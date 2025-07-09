// GameUserSettings.ini 配置参数定义
export const gameUserSettingsParams = {
  // 服务器基本设置
  serverBasic: {
    title: '服务器基本设置',
    params: {
      ServerPassword: {
        type: 'password',
        label: '服务器密码',
        default: '',
        description: '玩家加入服务器需要的密码，留空表示无密码'
      },
      SpectatorPassword: {
        type: 'password',
        label: '观察者密码',
        default: '',
        description: '非管理员观察模式密码'
      },
      AdminLogging: {
        type: 'boolean',
        label: '管理员日志',
        default: true,
        description: '将所有管理员命令记录到游戏聊天中'
      }
    }
  },

  // 游戏模式设置
  gameMode: {
    title: '游戏模式设置',
    params: {
      serverPVE: {
        type: 'boolean',
        label: 'PvE模式',
        default: false,
        description: '启用PvE模式，禁用玩家间伤害'
      },
      serverHardcore: {
        type: 'boolean',
        label: '硬核模式',
        default: false,
        description: '玩家死亡后重新开始，等级重置为1'
      },
      ShowMapPlayerLocation: {
        type: 'boolean',
        label: '显示玩家位置',
        default: false,
        description: '在地图上显示玩家自己的位置'
      },
      allowThirdPersonPlayer: {
        type: 'boolean',
        label: '允许第三人称',
        default: false,
        description: '允许玩家使用第三人称视角'
      },
      ServerCrosshair: {
        type: 'boolean',
        label: '显示准星',
        default: true,
        description: '显示武器瞄准准星'
      },
      EnablePvPGamma: {
        type: 'boolean',
        label: 'PvP伽马调节',
        default: false,
        description: '在PvP模式下允许使用gamma控制台命令'
      },
      DisablePvEGamma: {
        type: 'boolean',
        label: '禁用PvE伽马调节',
        default: false,
        description: '在PvE模式下禁用gamma控制台命令'
      },
      serverForceNoHud: {
        type: 'boolean',
        label: '强制隐藏HUD',
        default: false,
        description: '强制禁用游戏界面显示'
      },
      ShowFloatingDamageText: {
        type: 'boolean',
        label: '显示浮动伤害文字',
        default: false,
        description: '启用RPG风格的弹出伤害数字显示'
      },
      AllowHitMarkers: {
        type: 'boolean',
        label: '允许命中标记',
        default: true,
        description: '允许远程攻击的可选标记'
      }
    }
  },

  // 聊天和通讯设置
  communication: {
    title: '聊天和通讯设置',
    params: {
      globalVoiceChat: {
        type: 'boolean',
        label: '全局语音聊天',
        default: false,
        description: '启用全局语音聊天，所有玩家都能听到'
      },
      proximityChat: {
        type: 'boolean',
        label: '附近聊天',
        default: false,
        description: '只有附近的玩家才能看到彼此的聊天消息'
      },
      alwaysNotifyPlayerJoined: {
        type: 'boolean',
        label: '总是通知玩家加入',
        default: false,
        description: '玩家加入服务器时总是通知其他玩家'
      },
      alwaysNotifyPlayerLeft: {
        type: 'boolean',
        label: '总是通知玩家离开',
        default: false,
        description: '玩家离开服务器时总是通知其他玩家'
      },
      DontAlwaysNotifyPlayerJoined: {
        type: 'boolean',
        label: '禁用玩家加入通知',
        default: false,
        description: '全局禁用玩家加入通知'
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
      },
      OxygenSwimSpeedStatMultiplier: {
        type: 'number',
        label: '氧气游泳速度倍率',
        default: 1.0,
        min: 0.1,
        max: 10,
        step: 0.1,
        description: '设置游泳速度如何被氧气等级倍增'
      },
      ImplantSuicideCD: {
        type: 'number',
        label: '植入体自杀冷却',
        default: 28800,
        min: 0,
        max: 86400,
        description: '玩家使用植入体重生功能间隔时间(秒)'
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
        default: 4000,
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
      },
      DisableDinoDecayPvE: {
        type: 'boolean',
        label: '禁用PvE恐龙腐朽',
        default: false,
        description: '在PvE模式下禁用恐龙自动腐朽'
      },
      AutoDestroyDecayedDinos: {
        type: 'boolean',
        label: '自动销毁腐朽恐龙',
        default: false,
        description: '加载时自动销毁可认领的腐朽恐龙'
      },
      PvEDinoDecayPeriodMultiplier: {
        type: 'number',
        label: 'PvE恐龙腐朽倍率',
        default: 1.0,
        min: 0.1,
        max: 10,
        step: 0.1,
        description: 'PvE模式下恐龙腐朽期间倍率'
      },
      PvPDinoDecay: {
        type: 'boolean',
        label: 'PvP恐龙腐朽',
        default: false,
        description: '启用离线突袭保护时恐龙是否腐朽'
      },
      AllowRaidDinoFeeding: {
        type: 'boolean',
        label: '允许突袭恐龙喂食',
        default: false,
        description: '允许泰坦龙等突袭恐龙被永久驯服'
      },
      RaidDinoCharacterFoodDrainMultiplier: {
        type: 'number',
        label: '突袭恐龙饥饿倍率',
        default: 1.0,
        min: 0.1,
        max: 10,
        step: 0.1,
        description: '突袭恐龙(如泰坦龙)食物消耗速度倍率'
      },
      AllowFlyerCarryPvE: {
        type: 'boolean',
        label: 'PvE飞行载具抓取',
        default: false,
        description: 'PvE模式下允许飞行生物抓取野生恐龙和玩家'
      },
      bForceCanRideFliers: {
        type: 'boolean',
        label: '强制允许骑乘飞行生物',
        default: false,
        description: '允许在通常禁用飞行生物的地图上使用飞行生物'
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
        default: 0.2,
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
      },
      RandomSupplyCratePoints: {
        type: 'boolean',
        label: '随机补给箱位置',
        default: false,
        description: '补给箱在随机位置生成'
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
      NewMaxStructuresInRange: {
        type: 'number',
        label: '新范围内最大建筑数',
        default: 6000,
        min: 1000,
        max: 50000,
        description: '新的建筑数量限制系统'
      },
      MaxStructuresInRange: {
        type: 'number',
        label: '最大建筑范围数',
        default: 1300,
        min: 100,
        max: 10000,
        description: '建筑范围限制的最大数量'
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
      PvEStructureDecayDestructionPeriod: {
        type: 'number',
        label: 'PvE建筑腐朽销毁期',
        default: 0,
        min: 0,
        max: 1000,
        description: 'PvE模式下建筑腐朽销毁时间(天)'
      },
      PvPStructureDecay: {
        type: 'boolean',
        label: 'PvP建筑腐朽',
        default: false,
        description: '启用离线突袭保护时建筑是否腐朽'
      },
      StructurePickupTimeAfterPlacement: {
        type: 'number',
        label: '建筑快捷拾取时间',
        default: 30,
        min: 0,
        max: 3600,
        description: '建筑放置后可快捷拾取的时间(秒)'
      },
      StructurePickupHoldDuration: {
        type: 'number',
        label: '拾取按键持续时间',
        default: 0.5,
        min: 0,
        max: 10,
        step: 0.1,
        description: '快捷拾取的按键持续时间，0表示立即拾取'
      },
      AlwaysAllowStructurePickup: {
        type: 'boolean',
        label: '总是允许拾取建筑',
        default: false,
        description: '禁用建筑快捷拾取的时间限制'
      },
      OnlyAutoDestroyCoreStructures: {
        type: 'boolean',
        label: '仅自动销毁核心建筑',
        default: false,
        description: '防止非核心/非地基建筑自动销毁'
      },
      OnlyDecayUnsnappedCoreStructures: {
        type: 'boolean',
        label: '仅腐朽未连接核心建筑',
        default: false,
        description: '只有未连接的核心建筑会腐朽'
      },
      FastDecayUnsnappedCoreStructures: {
        type: 'boolean',
        label: '快速腐朽未连接建筑',
        default: false,
        description: '未连接地基/柱子以5倍速度腐朽'
      },
      DestroyUnconnectedWaterPipes: {
        type: 'boolean',
        label: '销毁未连接水管',
        default: false,
        description: '2天后自动销毁未连接的水管'
      },
      StructurePreventResourceRadiusMultiplier: {
        type: 'number',
        label: '建筑阻止资源半径倍率',
        default: 1.0,
        min: 0.1,
        max: 10,
        step: 0.1,
        description: '建筑阻止资源重生的半径倍率'
      },
      MaxPlatformSaddleStructureLimit: {
        type: 'number',
        label: '平台鞍最大建筑数',
        default: 25,
        min: 1,
        max: 1000,
        description: '单个平台鞍上允许的最大建筑数量'
      },
      PerPlatformMaxStructuresMultiplier: {
        type: 'number',
        label: '平台建筑倍率',
        default: 1.0,
        min: 0.1,
        max: 10,
        step: 0.1,
        description: '增加鞍座和木筏上可放置物品的最大数量'
      },
      PlatformSaddleBuildAreaBoundsMultiplier: {
        type: 'number',
        label: '平台建造区域倍率',
        default: 1.0,
        min: 0.1,
        max: 10,
        step: 0.1,
        description: '允许建筑在平台更远距离放置'
      },
      OverrideStructurePlatformPrevention: {
        type: 'boolean',
        label: '覆盖平台限制',
        default: false,
        description: '允许炮塔在平台鞍上建造和使用'
      },
      EnableExtraStructurePreventionVolumes: {
        type: 'boolean',
        label: '启用额外建筑限制区域',
        default: false,
        description: '在特定资源丰富区域禁用建筑'
      },
      AllowCaveBuildingPvE: {
        type: 'boolean',
        label: '允许PvE洞穴建造',
        default: false,
        description: 'PvE模式下允许在洞穴中建造'
      },
      AllowCaveBuildingPvP: {
        type: 'boolean',
        label: '允许PvP洞穴建造',
        default: false,
        description: 'PvP模式下允许在洞穴中建造'
      },
      PvEAllowStructuresAtSupplyDrops: {
        type: 'boolean',
        label: 'PvE允许补给点建造',
        default: false,
        description: 'PvE模式下允许在补给点附近建造'
      },
      AllowCrateSpawnsOnTopOfStructures: {
        type: 'boolean',
        label: '允许补给箱生成在建筑上',
        default: false,
        description: '允许空投补给箱出现在建筑物顶部'
      },
      bAllowPlatformSaddleMultiFloors: {
        type: 'boolean',
        label: '允许平台鞍多层',
        default: false,
        description: '允许平台鞍建造多层建筑'
      },
      MaxGateFrameOnSaddles: {
        type: 'number',
        label: '鞍座最大门框数',
        default: 0,
        min: 0,
        max: 100,
        description: '平台鞍上允许的最大门框数量'
      }
    }
  },

  // 部落和联盟设置
  tribeSettings: {
    title: '部落和联盟设置',
    params: {
      MaxNumberOfPlayersInTribe: {
        type: 'number',
        label: '部落最大人数',
        default: 0,
        min: 0,
        max: 500,
        description: '单个部落最大玩家数量，0表示无限制'
      },
      TribeNameChangeCooldown: {
        type: 'number',
        label: '部落改名冷却',
        default: 15,
        min: 0,
        max: 1440,
        description: '部落名称更改之间的冷却时间(分钟)'
      },
      PreventTribeAlliances: {
        type: 'boolean',
        label: '防止部落联盟',
        default: false,
        description: '阻止部落创建联盟'
      },
      MaxAlliancesPerTribe: {
        type: 'number',
        label: '每个部落最大联盟数',
        default: 10,
        min: 0,
        max: 100,
        description: '每个部落可以加入的最大联盟数量'
      },
      MaxTribesPerAlliance: {
        type: 'number',
        label: '每个联盟最大部落数',
        default: 10,
        min: 0,
        max: 100,
        description: '每个联盟可以包含的最大部落数量'
      }
    }
  },

  // 繁殖和印记设置
  breedingSettings: {
    title: '繁殖和印记设置',
    params: {
      AllowAnyoneBabyImprintCuddle: {
        type: 'boolean',
        label: '允许任何人照顾幼体',
        default: false,
        description: '允许任何人照顾幼体恐龙，而不仅仅是印记者'
      },
      DisableImprintDinoBuff: {
        type: 'boolean',
        label: '禁用印记恐龙加成',
        default: false,
        description: '禁用恐龙印记玩家的属性加成'
      },
      BabyImprintingStatScaleMultiplier: {
        type: 'number',
        label: '幼体印记属性倍率',
        default: 1.0,
        min: 0.1,
        max: 10,
        step: 0.1,
        description: '幼体印记后属性提升倍率'
      }
    }
  },

  // 物品和补给设置
  itemSettings: {
    title: '物品和补给设置',
    params: {
      ClampItemSpoilingTimes: {
        type: 'boolean',
        label: '限制物品腐坏时间',
        default: false,
        description: '将所有腐坏时间限制在物品最大腐坏时间内'
      },
      ClampResourceHarvestDamage: {
        type: 'boolean',
        label: '限制资源采集伤害',
        default: false,
        description: '基于资源剩余血量限制驯服生物对资源的伤害'
      },
      UseOptimizedHarvestingHealth: {
        type: 'boolean',
        label: '使用优化采集血量',
        default: true,
        description: '使用优化的资源采集血量计算'
      },
      BanListURL: {
        type: 'text',
        label: '封禁列表URL',
        default: 'http://arkdedicated.com/banlist.txt',
        description: '指向在线封禁列表的URL，必须用引号包围'
      }
    }
  },

  // 服务器性能设置
  performanceSettings: {
    title: '服务器性能设置',
    params: {
      AutoSavePeriodMinutes: {
        type: 'number',
        label: '自动保存间隔',
        default: 15.0,
        min: 1,
        max: 60,
        step: 0.1,
        description: '自动保存间隔时间(分钟)，设为0将持续保存'
      },
      KickIdlePlayersPeriod: {
        type: 'number',
        label: '踢出空闲玩家时间',
        default: 3600,
        min: 0,
        max: 86400,
        description: '踢出空闲玩家的时间(秒)，0禁用此功能'
      },
      ListenServerTetherDistanceMultiplier: {
        type: 'number',
        label: '监听服务器距离倍率',
        default: 1.0,
        min: 0.1,
        max: 10,
        step: 0.1,
        description: '非专用服务器的束缚距离倍率'
      },
      RCONServerGameLogBuffer: {
        type: 'number',
        label: 'RCON游戏日志缓冲',
        default: 600,
        min: 100,
        max: 10000,
        description: 'RCON服务器游戏日志缓冲区大小'
      },
      NPCNetworkStasisRangeScalePlayerCountStart: {
        type: 'number',
        label: 'NPC网络停滞范围玩家数起始',
        default: 0,
        min: 0,
        max: 1000,
        description: 'NPC网络停滞范围缩放的玩家数量起始值'
      },
      NPCNetworkStasisRangeScalePlayerCountEnd: {
        type: 'number',
        label: 'NPC网络停滞范围玩家数结束',
        default: 0,
        min: 0,
        max: 1000,
        description: 'NPC网络停滞范围缩放的玩家数量结束值'
      },
      NPCNetworkStasisRangeScalePercentEnd: {
        type: 'number',
        label: 'NPC网络停滞范围缩放百分比结束',
        default: 0,
        min: 0,
        max: 100,
        step: 0.1,
        description: 'NPC网络停滞范围缩放的百分比结束值'
      }
    }
  },

  // 疾病和状态设置
  diseaseSettings: {
    title: '疾病和状态设置',
    params: {
      PreventDiseases: {
        type: 'boolean',
        label: '防止疾病',
        default: false,
        description: '完全禁用服务器上的疾病'
      },
      NonPermanentDiseases: {
        type: 'boolean',
        label: '非永久性疾病',
        default: false,
        description: '疾病不永久存在，重生后会失去疾病'
      },
      PreventSpawnAnimations: {
        type: 'boolean',
        label: '防止生成动画',
        default: false,
        description: '玩家角色重生时不播放苏醒动画'
      }
    }
  },

  // 离线突袭保护设置
  offlineRaidSettings: {
    title: '离线突袭保护设置',
    params: {
      PreventOfflinePvP: {
        type: 'boolean',
        label: '防止离线PvP',
        default: false,
        description: '启用离线突袭保护(ORP)功能'
      },
      PreventOfflinePvPInterval: {
        type: 'number',
        label: '离线保护等待时间',
        default: 900,
        min: 0,
        max: 3600,
        description: '部落/玩家下线后ORP激活的等待时间(秒)'
      }
    }
  },

  // 跨服传输设置
  crossArkSettings: {
    title: '跨服传输设置',
    params: {
      NoTributeDownloads: {
        type: 'boolean',
        label: '禁用贡品下载',
        default: false,
        description: '禁用所有跨服数据传输'
      },
      PreventDownloadSurvivors: {
        type: 'boolean',
        label: '防止下载角色',
        default: false,
        description: '防止从其他服务器下载角色'
      },
      PreventDownloadItems: {
        type: 'boolean',
        label: '防止下载物品',
        default: false,
        description: '防止从其他服务器下载物品'
      },
      PreventDownloadDinos: {
        type: 'boolean',
        label: '防止下载恐龙',
        default: false,
        description: '防止从其他服务器下载恐龙'
      },
      PreventUploadSurvivors: {
        type: 'boolean',
        label: '防止上传角色',
        default: false,
        description: '防止角色上传到其他服务器'
      },
      PreventUploadItems: {
        type: 'boolean',
        label: '防止上传物品',
        default: false,
        description: '防止物品上传到其他服务器'
      },
      PreventUploadDinos: {
        type: 'boolean',
        label: '防止上传恐龙',
        default: false,
        description: '防止恐龙上传到其他服务器'
      },
      CrossARKAllowForeignDinoDownloads: {
        type: 'boolean',
        label: '允许外来恐龙下载',
        default: false,
        description: '允许下载外来地图的恐龙'
      },
      MaxTributeDinos: {
        type: 'number',
        label: '最大贡品恐龙数',
        default: 20,
        min: 1,
        max: 500,
        description: '上传恐龙的槽位数量'
      },
      MaxTributeItems: {
        type: 'number',
        label: '最大贡品物品数',
        default: 50,
        min: 1,
        max: 500,
        description: '上传物品和资源的槽位数量'
      },
      MinimumDinoReuploadInterval: {
        type: 'number',
        label: '最小恐龙重新上传间隔',
        default: 0,
        min: 0,
        max: 86400,
        description: '允许恐龙重新上传之间的冷却时间(秒)'
      },
      TributeItemExpirationSeconds: {
        type: 'number',
        label: '贡品物品过期时间',
        default: 86400,
        min: 0,
        max: 604800,
        description: '上传物品的过期时间(秒)'
      },
      TributeDinoExpirationSeconds: {
        type: 'number',
        label: '贡品恐龙过期时间',
        default: 86400,
        min: 0,
        max: 604800,
        description: '上传恐龙的过期时间(秒)'
      },
      TributeCharacterExpirationSeconds: {
        type: 'number',
        label: '贡品角色过期时间',
        default: 86400,
        min: 0,
        max: 604800,
        description: '上传角色的过期时间(秒)'
      }
    }
  },

  // 飞行载具设置
  flyerSettings: {
    title: '飞行载具设置',
    params: {
      AllowFlyingStaminaRecovery: {
        type: 'boolean',
        label: '飞行耐力恢复',
        default: true,
        description: '玩家站在飞行生物上时是否恢复耐力'
      },
      ForceFlyerExplosives: {
        type: 'boolean',
        label: '强制飞行器爆炸物',
        default: false,
        description: '允许携带C4的飞行生物继续飞行'
      }
    }
  },

  // 高级功能设置
  advancedSettings: {
    title: '高级功能设置',
    params: {
      AllowMultipleAttachedC4: {
        type: 'boolean',
        label: '允许多个C4',
        default: false,
        description: '允许在一个恐龙上附着多个C4'
      },
      AllowIntegratedSPlusStructures: {
        type: 'boolean',
        label: '允许集成S+建筑',
        default: true,
        description: '允许使用集成的Structures Plus建筑'
      },
      AllowHideDamageSourceFromLogs: {
        type: 'boolean',
        label: '隐藏伤害来源日志',
        default: false,
        description: '在部落日志中隐藏伤害来源'
      },
      AllowSharedConnections: {
        type: 'boolean',
        label: '允许共享连接',
        default: false,
        description: '允许共享网络连接，False表示禁用家庭共享'
      },
      bFilterTribeNames: {
        type: 'boolean',
        label: '过滤部落名称',
        default: false,
        description: '启用部落名称过滤'
      },
      bFilterCharacterNames: {
        type: 'boolean',
        label: '过滤角色名称',
        default: false,
        description: '启用角色名称过滤'
      },
      bFilterChat: {
        type: 'boolean',
        label: '过滤聊天',
        default: false,
        description: '启用聊天内容过滤'
      },
      EnableCryoSicknessPVE: {
        type: 'boolean',
        label: '启用PvE冷冻舱疾病',
        default: true,
        description: 'PvE模式下启用冷冻舱疾病机制'
      },
      EnableCryopodNerf: {
        type: 'boolean',
        label: '启用冷冻舱削弱',
        default: false,
        description: '启用从冷冻舱部署生物时的冷冻疾病'
      },
      CryopodNerfDuration: {
        type: 'number',
        label: '冷冻舱削弱持续时间',
        default: 10,
        min: 0,
        max: 3600,
        description: '从冷冻舱部署生物后冷冻疾病持续时间(秒)'
      },
      CryopodNerfDamageMult: {
        type: 'number',
        label: '冷冻舱削弱伤害倍率',
        default: 0.01,
        min: 0,
        max: 1,
        step: 0.01,
        description: '从冷冻舱部署的生物造成的伤害倍率'
      },
      CryopodNerfIncomingDamageMultPercent: {
        type: 'number',
        label: '冷冻舱削弱受到伤害倍率',
        default: 0.25,
        min: 0,
        max: 1,
        step: 0.01,
        description: '从冷冻舱部署的生物受到伤害的倍率'
      },
      DisableCryopodEnemyCheck: {
        type: 'boolean',
        label: '禁用冷冻舱敌人检查',
        default: false,
        description: '允许在敌人附近使用冷冻舱'
      },
      DisableCryopodFridgeRequirement: {
        type: 'boolean',
        label: '禁用冷冻舱冰箱需求',
        default: false,
        description: '允许不在冷冻冰箱范围内使用冷冻舱'
      },
      AllowCryoFridgeOnSaddle: {
        type: 'boolean',
        label: '允许鞍座冷冻冰箱',
        default: false,
        description: '允许在平台鞍座和木筏上建造冷冻冰箱'
      },
      MaxTrainCars: {
        type: 'number',
        label: '最大火车车厢数',
        default: 8,
        min: 1,
        max: 20,
        description: '火车可以拥有的最大车厢数量'
      },
      MaxHexagonsPerCharacter: {
        type: 'number',
        label: '每角色最大六边形数',
        default: 2000000000,
        min: 0,
        max: 2000000000,
        description: '每个角色可以拥有的最大六边形数量'
      },
      AllowTekSuitPowersInGenesis: {
        type: 'boolean',
        label: '允许创世纪TEK套装能力',
        default: false,
        description: '在创世纪第一部分启用或禁用TEK套装能力'
      },
      CustomDynamicConfigUrl: {
        type: 'text',
        label: '自定义动态配置URL',
        default: '',
        description: '配置文件的直接链接'
      }
    }
  }
} 