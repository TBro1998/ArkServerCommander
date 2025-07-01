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