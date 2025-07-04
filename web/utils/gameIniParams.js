// Game.ini 配置参数定义
export const gameIniParams = {
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
      },
      bOnlyAllowSpecifiedEngrams: {
        type: 'boolean',
        label: '仅允许指定图纸',
        default: false,
        description: '如果为true，未明确指定的图纸将被隐藏'
      },
      bAutoUnlockAllEngrams: {
        type: 'boolean',
        label: '自动解锁所有图纸',
        default: false,
        description: '解锁所有可用的图纸'
      },
      bShowCreativeMode: {
        type: 'boolean',
        label: '显示创造模式',
        default: false,
        description: '启用创造模式'
      },
      bUseCorpseLocator: {
        type: 'boolean',
        label: '使用尸体定位器',
        default: false,
        description: '如果设置为true，你将在死亡位置看到绿色光束'
      },
      bDisableLootCrates: {
        type: 'boolean',
        label: '禁用补给箱',
        default: false,
        description: 'True阻止补给箱生成'
      },
      bDisableDinoRiding: {
        type: 'boolean',
        label: '禁用骑乘恐龙',
        default: false,
        description: '禁用恐龙骑乘功能'
      },
      bDisableDinoTaming: {
        type: 'boolean',
        label: '禁用驯服恐龙',
        default: false,
        description: '禁用恐龙驯服功能'
      },
      bAllowCustomRecipes: {
        type: 'boolean',
        label: '允许自定义配方',
        default: true,
        description: '允许玩家创建自定义烹饪配方'
      },
      bHardLimitTurretsInRange: {
        type: 'boolean',
        label: '硬限制范围内炮塔',
        default: false,
        description: '硬限制范围内的炮塔数量'
      },
      bPvEAllowTribeWar: {
        type: 'boolean',
        label: 'PvE允许部落战争',
        default: false,
        description: 'False禁用部落正式宣战的能力'
      },
      bPvEAllowTribeWarCancel: {
        type: 'boolean',
        label: 'PvE允许取消部落战争',
        default: false,
        description: 'True允许在战争实际开始前取消已同意的战争'
      },
      bPvEDisableFriendlyFire: {
        type: 'boolean',
        label: 'PvE禁用友军伤害',
        default: false,
        description: 'PvE服务器的友军伤害防护选项'
      },
      bDisableFriendlyFire: {
        type: 'boolean',
        label: '禁用友军伤害',
        default: false,
        description: 'PvP服务器的友军伤害防护选项'
      },
      bFlyerPlatformAllowUnalignedDinoBasing: {
        type: 'boolean',
        label: '飞行平台允许未对齐恐龙基地',
        default: false,
        description: '飞行时风神平台不允许任何非盟友恐龙在其上建立基地'
      },
      bIncreasePvPRespawnInterval: {
        type: 'boolean',
        label: '增加PvP重生间隔',
        default: false,
        description: 'PvP服务器有可选的+1分钟额外重生时间，如果在5分钟内被同一队伍杀死，每次都会翻倍'
      },
      bAutoPvETimer: {
        type: 'boolean',
        label: '自动PvE计时器',
        default: false,
        description: '允许在预设的游戏时间或现实时间从PvE切换到PvP模式'
      },
      bAutoPvEUseSystemTime: {
        type: 'boolean',
        label: '自动PvE使用系统时间',
        default: false,
        description: '使用系统时间而不是游戏时间进行PvE/PvP切换'
      },
      bIgnoreStructuresPreventionVolumes: {
        type: 'boolean',
        label: '忽略建筑防护体积',
        default: false,
        description: '在创世纪第一部分的任务体积中启用建筑'
      },
      bGenesisUseStructuresPreventionVolumes: {
        type: 'boolean',
        label: '创世纪使用建筑防护体积',
        default: false,
        description: '如果设置为true，在创世纪第一部分的任务区域禁用建筑'
      },
      bAllowFlyerSpeedLeveling: {
        type: 'boolean',
        label: '允许飞行生物速度升级',
        default: false,
        description: '指定飞行生物是否可以升级移动速度'
      },
      bUseTameLimitForStructuresOnly: {
        type: 'boolean',
        label: '仅对建筑使用驯服限制',
        default: false,
        description: '如果为true，驯服单位将仅应用于带有建筑和木筏的平台'
      }
    }
  },

  // 经验值和等级设置
  experienceSettings: {
    title: '经验值和等级设置',
    params: {
      OverrideMaxExperiencePointsPlayer: {
        type: 'number',
        label: '覆盖玩家最大经验值',
        default: 0,
        min: 0,
        max: 999999999,
        description: '覆盖玩家的最大经验值上限'
      },
      OverrideMaxExperiencePointsDino: {
        type: 'number',
        label: '覆盖恐龙最大经验值',
        default: 0,
        min: 0,
        max: 999999999,
        description: '覆盖恐龙的最大经验值上限'
      },
      OverridePlayerLevelEngramPoints: {
        type: 'number',
        label: '覆盖玩家等级图纸点数',
        default: 0,
        min: 0,
        max: 1000,
        description: '配置玩家每升一级获得的图纸点数'
      },
      KillXPMultiplier: {
        type: 'number',
        label: '击杀经验倍率',
        default: 1.0,
        min: 0.1,
        max: 100,
        step: 0.1,
        description: '击杀获得经验值的倍率'
      },
      HarvestXPMultiplier: {
        type: 'number',
        label: '采集经验倍率',
        default: 1.0,
        min: 0.1,
        max: 100,
        step: 0.1,
        description: '采集获得经验值的倍率'
      },
      CraftXPMultiplier: {
        type: 'number',
        label: '制作经验倍率',
        default: 1.0,
        min: 0.1,
        max: 100,
        step: 0.1,
        description: '制作获得经验值的倍率'
      },
      GenericXPMultiplier: {
        type: 'number',
        label: '通用经验倍率',
        default: 1.0,
        min: 0.1,
        max: 100,
        step: 0.1,
        description: '通用经验值（随时间自动获得）的倍率'
      },
      SpecialXPMultiplier: {
        type: 'number',
        label: '特殊事件经验倍率',
        default: 1.0,
        min: 0.1,
        max: 100,
        step: 0.1,
        description: '特殊事件获得经验值的倍率'
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
      MatingSpeedMultiplier: {
        type: 'number',
        label: '交配速度倍率',
        default: 1.0,
        min: 0.1,
        max: 10,
        step: 0.1,
        description: '恐龙交配速度倍率，值越高交配越快'
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
      },
      BabyImprintingStatScaleMultiplier: {
        type: 'number',
        label: '幼体印记属性缩放倍率',
        default: 1.0,
        min: 0.1,
        max: 10,
        step: 0.1,
        description: '印记质量对属性影响的程度'
      },
      BabyImprintAmountMultiplier: {
        type: 'number',
        label: '幼体印记数量倍率',
        default: 1.0,
        min: 0.1,
        max: 10,
        step: 0.1,
        description: '应用于每次印记提供的百分比的倍率'
      },
      BabyCuddleGracePeriodMultiplier: {
        type: 'number',
        label: '幼体关爱宽限期倍率',
        default: 1.0,
        min: 0.1,
        max: 10,
        step: 0.1,
        description: '延迟关爱幼体后印记质量开始下降前的宽限期倍率'
      },
      BabyCuddleLoseImprintQualitySpeedMultiplier: {
        type: 'number',
        label: '幼体关爱失去印记质量速度倍率',
        default: 1.0,
        min: 0.1,
        max: 10,
        step: 0.1,
        description: '宽限期后印记质量下降速度的倍率'
      }
    }
  },

  // 物品和资源设置
  itemSettings: {
    title: '物品和资源设置',
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
      CropDecaySpeedMultiplier: {
        type: 'number',
        label: '作物腐烂速度倍率',
        default: 1.0,
        min: 0.1,
        max: 10,
        step: 0.1,
        description: '农田中作物腐烂速度倍率，值越高腐烂越快'
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
      },
      UseCorpseLifeSpanMultiplier: {
        type: 'number',
        label: '尸体生命周期倍率',
        default: 1.0,
        min: 0.1,
        max: 10,
        step: 0.1,
        description: '尸体存在时间倍率'
      },
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
      CustomRecipeEffectivenessMultiplier: {
        type: 'number',
        label: '自定义配方效果倍率',
        default: 1.0,
        min: 0.1,
        max: 10,
        step: 0.1,
        description: '自定义配方效果的倍率，值越高效果越好'
      },
      CustomRecipeSkillMultiplier: {
        type: 'number',
        label: '自定义配方技能倍率',
        default: 1.0,
        min: 0.1,
        max: 10,
        step: 0.1,
        description: '玩家制作速度等级对自定义配方的影响倍率'
      },
      CraftingSkillBonusMultiplier: {
        type: 'number',
        label: '制作技能奖励倍率',
        default: 1.0,
        min: 0.1,
        max: 10,
        step: 0.1,
        description: '升级制作技能获得的奖励倍率'
      },
      ResourceNoReplenishRadiusPlayers: {
        type: 'number',
        label: '玩家资源不重生半径',
        default: 1.0,
        min: 0.1,
        max: 1000,
        step: 0.1,
        description: '允许资源在距离玩家更近或更远的地方重生'
      },
      ResourceNoReplenishRadiusStructures: {
        type: 'number',
        label: '建筑资源不重生半径',
        default: 1.0,
        min: 0.1,
        max: 1000,
        step: 0.1,
        description: '允许资源在距离建筑更近或更远的地方重生'
      },
      HarvestResourceItemAmountClassMultipliers: {
        type: 'text',
        label: '资源采集数量类别倍率',
        default: '',
        description: '按资源类型缩放采集的资源数量，格式：(ClassName="classname", Multiplier=value)'
      },
      DinoHarvestingDamageMultiplier: {
        type: 'number',
        label: '恐龙采集伤害倍率',
        default: 1.0,
        min: 0.1,
        max: 10,
        step: 0.1,
        description: '恐龙对可采集物品/实体造成的伤害倍率'
      },
      PlayerHarvestingDamageMultiplier: {
        type: 'number',
        label: '玩家采集伤害倍率',
        default: 1.0,
        min: 0.1,
        max: 10,
        step: 0.1,
        description: '玩家对可采集物品/实体造成的伤害倍率'
      }
    }
  },

  // 恐龙设置
  dinoSettings: {
    title: '恐龙设置',
    params: {
      TamedDinoCharacterFoodDrainMultiplier: {
        type: 'number',
        label: '驯服恐龙食物消耗倍率',
        default: 1.0,
        min: 0.1,
        max: 10,
        step: 0.1,
        description: '驯服恐龙消耗食物的速度倍率'
      },
      WildDinoCharacterFoodDrainMultiplier: {
        type: 'number',
        label: '野生恐龙食物消耗倍率',
        default: 1.0,
        min: 0.1,
        max: 10,
        step: 0.1,
        description: '野生恐龙消耗食物的速度倍率'
      },
      WildDinoTorporDrainMultiplier: {
        type: 'number',
        label: '野生恐龙昏迷流失倍率',
        default: 1.0,
        min: 0.1,
        max: 10,
        step: 0.1,
        description: '野生恐龙失去昏迷值的速度倍率'
      },
      TamedDinoTorporDrainMultiplier: {
        type: 'number',
        label: '驯服恐龙昏迷流失倍率',
        default: 1.0,
        min: 0.1,
        max: 10,
        step: 0.1,
        description: '驯服恐龙失去昏迷值的速度倍率'
      },
      PassiveTameIntervalMultiplier: {
        type: 'number',
        label: '被动驯服间隔倍率',
        default: 1.0,
        min: 0.1,
        max: 10,
        step: 0.1,
        description: '被动驯服恐龙获得驯服请求的频率倍率'
      },
      DinoTurretDamageMultiplier: {
        type: 'number',
        label: '恐龙炮塔伤害倍率',
        default: 1.0,
        min: 0.1,
        max: 10,
        step: 0.1,
        description: '炮塔对恐龙造成伤害的倍率'
      },
      PreventDinoTameClassNames: {
        type: 'text',
        label: '禁止驯服恐龙类别',
        default: '',
        description: '通过类别名称禁止驯服特定恐龙'
      },
      PreventTransferForClassName: {
        type: 'text',
        label: '禁止转移类别名称',
        default: '',
        description: '通过类别名称禁止转移特定恐龙'
      },
      DinoClassDamageMultipliers: {
        type: 'text',
        label: '恐龙类别伤害倍率',
        default: '',
        description: '通过类别名称乘以特定恐龙的伤害，格式：(ClassName="classname", Multiplier=value)'
      },
      TamedDinoClassDamageMultipliers: {
        type: 'text',
        label: '驯服恐龙类别伤害倍率',
        default: '',
        description: '通过类别名称乘以特定驯服恐龙的伤害，格式：(ClassName="classname", Multiplier=value)'
      },
      DinoClassResistanceMultipliers: {
        type: 'text',
        label: '恐龙类别抗性倍率',
        default: '',
        description: '通过类别名称乘以特定恐龙的抗性，格式：(ClassName="classname", Multiplier=value)'
      },
      TamedDinoClassResistanceMultipliers: {
        type: 'text',
        label: '驯服恐龙类别抗性倍率',
        default: '',
        description: '通过类别名称乘以特定驯服恐龙的抗性，格式：(ClassName="classname", Multiplier=value)'
      },
      DestroyTamesOverLevelClamp: {
        type: 'number',
        label: '销毁超过等级限制的驯服',
        default: 0,
        min: 0,
        max: 1000,
        description: '超过此等级的驯服将在服务器启动时被删除'
      }
    }
  },

  // 部落和玩家设置
  tribeSettings: {
    title: '部落和玩家设置',
    params: {
      MaxNumberOfPlayersInTribe: {
        type: 'number',
        label: '部落最大人数',
        default: 0,
        min: 0,
        max: 500,
        description: '单个部落最大玩家数量，0表示无限制'
      },
      MaxAlliancesPerTribe: {
        type: 'number',
        label: '每个部落最大联盟数',
        default: 0,
        min: 0,
        max: 100,
        description: '定义部落可以形成或参与的最大联盟数量'
      },
      MaxTribesPerAlliance: {
        type: 'number',
        label: '每个联盟最大部落数',
        default: 0,
        min: 0,
        max: 100,
        description: '定义联盟中的最大部落数量'
      },
      TribeSlotReuseCooldown: {
        type: 'number',
        label: '部落槽位重用冷却',
        default: 0,
        min: 0,
        max: 86400,
        description: '某人离开部落后锁定部落槽位的指定持续时间计时器'
      },
      MaxTribeLogs: {
        type: 'number',
        label: '最大部落日志',
        default: 100,
        min: 1,
        max: 1000,
        description: '每个部落显示的部落日志数量'
      },
      KickIdlePlayersPeriod: {
        type: 'number',
        label: '踢出空闲玩家周期',
        default: 0,
        min: 0,
        max: 86400,
        description: '未移动或交互的角色将被踢出的时间'
      }
    }
  },

  // PvP设置
  pvpSettings: {
    title: 'PvP设置',
    params: {
      IncreasePvPRespawnIntervalCheckPeriod: {
        type: 'number',
        label: 'PvP重生间隔检查周期',
        default: 300,
        min: 1,
        max: 3600,
        description: '检查PvP重生间隔的时间周期'
      },
      IncreasePvPRespawnIntervalMultiplier: {
        type: 'number',
        label: 'PvP重生间隔倍率',
        default: 2.0,
        min: 1.0,
        max: 10.0,
        step: 0.1,
        description: 'PvP重生间隔的倍率'
      },
      IncreasePvPRespawnIntervalBaseAmount: {
        type: 'number',
        label: 'PvP重生间隔基础数量',
        default: 60,
        min: 1,
        max: 3600,
        description: 'PvP重生间隔的基础数量'
      },
      AutoPvEStartTimeSeconds: {
        type: 'number',
        label: '自动PvE开始时间（秒）',
        default: 0,
        min: 0,
        max: 86400,
        description: 'PvE模式的开始时间'
      },
      AutoPvEStopTimeSeconds: {
        type: 'number',
        label: '自动PvE停止时间（秒）',
        default: 0,
        min: 0,
        max: 86400,
        description: 'PvE模式的停止时间'
      },
      PvPZoneStructureDamageMultiplier: {
        type: 'number',
        label: 'PvP区域建筑伤害倍率',
        default: 1.0,
        min: 0.1,
        max: 10,
        step: 0.1,
        description: '指定洞穴内建筑所受伤害的缩放因子'
      }
    }
  },

  // 建筑和结构设置
  structureSettings: {
    title: '建筑和结构设置',
    params: {
      StructureDamageRepairCooldown: {
        type: 'number',
        label: '建筑伤害修复冷却',
        default: 180,
        min: 0,
        max: 3600,
        description: '建筑从上次受损后修复的冷却期'
      },
      StructureResistanceMultiplier: {
        type: 'number',
        label: '建筑抗性倍率',
        default: 1.0,
        min: 0.1,
        max: 10,
        step: 0.1,
        description: '建筑抗性倍率'
      },
      StructureDamageMultiplier: {
        type: 'number',
        label: '建筑伤害倍率',
        default: 1.0,
        min: 0.1,
        max: 10,
        step: 0.1,
        description: '建筑伤害倍率'
      },
      PvEStructureDecayPeriodMultiplier: {
        type: 'number',
        label: 'PvE建筑腐烂周期倍率',
        default: 1.0,
        min: 0.1,
        max: 100,
        step: 0.1,
        description: 'PvE模式下建筑腐烂周期倍率'
      },
      FastDecayInterval: {
        type: 'number',
        label: '快速腐烂间隔',
        default: 0,
        min: 0,
        max: 86400,
        description: '"快速腐烂"结构的固定常数腐烂周期'
      }
    }
  },

  // 高级设置
  advancedSettings: {
    title: '高级设置',
    params: {
      HairGrowthSpeedMultiplier: {
        type: 'number',
        label: '头发生长速度倍率',
        default: 1.0,
        min: 0.1,
        max: 10,
        step: 0.1,
        description: '头发生长速度倍率'
      },
      PoopIntervalMultiplier: {
        type: 'number',
        label: '排便间隔倍率',
        default: 1.0,
        min: 0.1,
        max: 10,
        step: 0.1,
        description: '排便频率倍率，值越高排便越少'
      },
      LayEggIntervalMultiplier: {
        type: 'number',
        label: '产蛋间隔倍率',
        default: 1.0,
        min: 0.1,
        max: 10,
        step: 0.1,
        description: '蛋生成/产下的时间间隔倍率，值越高间隔越长'
      },
      FuelConsumptionIntervalMultiplier: {
        type: 'number',
        label: '燃料消耗间隔倍率',
        default: 1.0,
        min: 0.1,
        max: 10,
        step: 0.1,
        description: '定义燃料消耗的间隔'
      },
      GlobalPoweredBatteryDurabilityDecreasePerSecond: {
        type: 'number',
        label: '全局供电电池耐久度每秒减少',
        default: 0,
        min: 0,
        max: 100,
        step: 0.1,
        description: '供电电池耐久度每秒减少量'
      },
      LimitNonPlayerDroppedItemsRange: {
        type: 'number',
        label: '限制非玩家掉落物品范围',
        default: 0,
        min: 0,
        max: 1000,
        description: '限制区域内掉落物品的数量'
      },
      LimitNonPlayerDroppedItemsCount: {
        type: 'number',
        label: '限制非玩家掉落物品数量',
        default: 0,
        min: 0,
        max: 10000,
        description: '限制区域内掉落物品的数量'
      },
      MaxFallSpeedMultiplier: {
        type: 'number',
        label: '最大下落速度倍率',
        default: 1.0,
        min: 0.1,
        max: 10,
        step: 0.1,
        description: '定义下落速度倍率'
      },
      PreventOfflinePvPConnectionInvincibleInterval: {
        type: 'number',
        label: '离线PvP连接无敌间隔',
        default: 0,
        min: 0,
        max: 3600,
        description: '离线PvP连接无敌时间间隔'
      }
    }
  },

  // 自定义配置
  customSettings: {
    title: '自定义配置',
    params: {
      ConfigOverrideItemMaxQuantity: {
        type: 'text',
        label: '配置覆盖物品最大数量',
        default: '',
        description: '允许手动覆盖每个物品的堆叠大小，格式：(ItemClassString="Class Name", Quantity=(MaxItemQuantity=n, bIgnoreMultiplier=value))'
      },
      ConfigOverrideItemCraftingCosts: {
        type: 'text',
        label: '配置覆盖物品制作成本',
        default: '',
        description: '自定义物品制作成本配置'
      },
      ConfigOverrideSupplyCrateItems: {
        type: 'text',
        label: '配置覆盖补给箱物品',
        default: '',
        description: '自定义补给箱物品配置'
      },
      ExcludeItemIndices: {
        type: 'text',
        label: '排除物品索引',
        default: '',
        description: '从补给箱中排除物品'
      },
      LevelExperienceRampOverrides: {
        type: 'text',
        label: '等级经验斜坡覆盖',
        default: '',
        description: '配置玩家和恐龙可用的总等级数以及达到每个等级所需的经验值，格式：(ExperiencePointsForLevel[n]=points...)'
      },
      EngramEntryAutoUnlocks: {
        type: 'text',
        label: '图纸条目自动解锁',
        default: '',
        description: '达到指定等级时自动解锁指定图纸，格式：(EngramClassName="index", LevelToAutoUnlock=value)'
      },
      ModIDS: {
        type: 'text',
        label: '模组ID',
        default: '',
        description: '指定要下载/安装/更新的额外Steam工坊模组/地图/TC ID的手动列表'
      }
    }
  }
} 