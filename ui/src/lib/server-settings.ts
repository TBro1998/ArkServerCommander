// 服务器设置的统一字典。
//
// 用户看到的是按用途分组的设置（常规、倍率、恐龙、建筑……），不关心每一项存在哪里：
// GameUserSettings.ini、Game.ini、启动参数（?查询参数 / -命令行参数），或服务器记录本身的字段。
// 每一项用 source 记录存放位置，读写由这里的函数统一处理。
//
// 默认值是"没有设置时游戏的表现"：ini 里没有这个键、启动参数里没有这个参数。
// 恢复默认就是删掉这个键或参数。

import type { LucideIcon } from 'lucide-react';
import {
  Settings2,
  Network,
  Package,
  Gamepad2,
  Gauge,
  User,
  PawPrint,
  Egg,
  Castle,
  Swords,
  Wrench,
  ShieldAlert,
} from 'lucide-react';
import { getIniValue, setIniValue, type IniIndex } from '@/lib/ini';
import type { Server } from '@/stores/servers';
import type { Game } from '@/lib/games';

// field：服务器记录的字段（名称、端口、Mod 等），由表单单独渲染
// gus：GameUserSettings.ini；game：Game.ini；query：?Key=Value；arg：-Key / -Key=Value
export type SettingSource = 'field' | 'gus' | 'game' | 'query' | 'arg';
// flags：从 options 里多选，用 separator 连成一个值（如 -ServerPlatform=PC+XSX），全选等于默认
export type SettingType = 'boolean' | 'number' | 'text' | 'password' | 'select' | 'flags';
// x：倍率；s：秒；min：分钟；h：小时
export type SettingUnit = 'x' | 's' | 'min' | 'h';
export type SettingValue = boolean | number | string;

export interface SettingDef {
  key: string;
  source: SettingSource;
  type: SettingType;
  default?: SettingValue;
  unit?: SettingUnit;
  step?: number;
  options?: string[];
  // flags 的分隔符，缺省为 +
  separator?: string;
  // ini 的 section，缺省按 source 决定
  section?: string;
  // 占满一行
  wide?: boolean;
  // 只对这些游戏有效，缺省两个游戏都有。ASE 的存档、建筑优化、反作弊类参数在 ASA 里不存在，
  // ASA 也有自己的设置（跨平台、地堡、无线制作等），不适用的项不显示。依据是 ark.wiki.gg 的 Server configuration，
  // 那里对 ASA 标"未知"的项照常显示
  games?: Game[];
  // 两个游戏的默认值不同时（如封禁列表地址），按游戏覆盖 default，settingGroupsFor 里换好
  gameDefaults?: Partial<Record<Game, SettingValue>>;
}

export type PanelId = 'customArgs' | 'rawConfig';

export interface SettingSectionDef {
  id: string;
  settings: SettingDef[];
  // 不是单个设置项的特殊面板（自定义启动参数、直接编辑配置文件）
  panel?: PanelId;
}

export type GroupId =
  | 'general'
  | 'network'
  | 'mods'
  | 'gameplay'
  | 'rates'
  | 'players'
  | 'dinos'
  | 'breeding'
  | 'structures'
  | 'tribes'
  | 'operations'
  | 'advanced';

export interface SettingGroupDef {
  id: GroupId;
  icon: LucideIcon;
  sections: SettingSectionDef[];
}

export const GUS_SECTION = '[ServerSettings]';
export const MOTD_SECTION = '[MessageOfTheDay]';
export const GAME_SECTION = '[/script/shootergame.shootergamemode]';

export function sectionOf(def: SettingDef): string {
  return def.section ?? (def.source === 'game' ? GAME_SECTION : GUS_SECTION);
}

// ---------- 定义用的简写 ----------

type Src = Exclude<SettingSource, 'field'>;
type Extra = Partial<Omit<SettingDef, 'key' | 'source' | 'type' | 'default'>>;

const field = (key: string, extra?: Extra): SettingDef => ({ key, source: 'field', type: 'text', ...extra });
const bool = (source: Src, key: string, def = false, extra?: Extra): SettingDef => ({ key, source, type: 'boolean', default: def, ...extra });
const num = (source: Src, key: string, def?: number, extra?: Extra): SettingDef => ({ key, source, type: 'number', default: def, ...extra });
const mult = (source: Src, key: string, def = 1): SettingDef => ({ key, source, type: 'number', default: def, unit: 'x', step: 0.1 });
const secs = (source: Src, key: string, def?: number): SettingDef => ({ key, source, type: 'number', default: def, unit: 's' });
const text = (source: Src, key: string, def = '', extra?: Extra): SettingDef => ({ key, source, type: 'text', default: def, ...extra });

// 只在 ASE / ASA 有效
const ase = (def: SettingDef): SettingDef => ({ ...def, games: ['ase'] });
const asa = (def: SettingDef): SettingDef => ({ ...def, games: ['asa'] });
const hours = (source: Src, key: string, def?: number): SettingDef => ({ key, source, type: 'number', default: def, unit: 'h', step: 0.1 });

const G = 'gus' as const;
const I = 'game' as const;
const Q = 'query' as const;
const A = 'arg' as const;

// ---------- 分组 ----------

export const SETTING_GROUPS: SettingGroupDef[] = [
  {
    id: 'general',
    icon: Settings2,
    sections: [
      { id: 'server', settings: [field('session_name'), field('map'), field('max_players'), ase(field('game_branch'))] },
      {
        id: 'access',
        settings: [
          { key: 'ServerPassword', source: G, type: 'password', default: '' },
          { key: 'SpectatorPassword', source: G, type: 'password', default: '' },
          field('admin_password', { wide: true }),
        ],
      },
      {
        id: 'motd',
        settings: [
          text(G, 'Message', '', { section: MOTD_SECTION, wide: true }),
          num(G, 'Duration', 20, { section: MOTD_SECTION, unit: 's' }),
        ],
      },
    ],
  },
  {
    id: 'network',
    icon: Network,
    sections: [
      { id: 'ports', settings: [field('port'), ase(field('query_port')), field('rcon_port')] },
      {
        id: 'platform',
        settings: [
          ase(bool(A, 'crossplay')),
          ase(bool(A, 'epiconly')),
          ase(text(A, 'PublicIPForEpic')),
          // ASA：不设置时允许所有平台（跨平台）
          asa({ key: 'ServerPlatform', source: A, type: 'flags', options: ['PC', 'XSX', 'PS5', 'WINGDK'], wide: true }),
          text(A, 'MULTIHOME'),
        ],
      },
      {
        id: 'connections',
        settings: [bool(A, 'exclusivejoin'), ase(num(A, 'MaxConnectionsPerIP')), bool(G, 'AllowSharedConnections')],
      },
    ],
  },
  {
    id: 'mods',
    icon: Package,
    sections: [
      {
        id: 'mods',
        settings: [
          field('game_mod_ids', { wide: true }),
          ase(text(A, 'MapModID')),
          ase(bool(A, 'automanagedmods')),
          asa(text(A, 'passivemods', '', { wide: true })),
        ],
      },
      {
        id: 'cosmetics',
        settings: [
          asa(bool(A, 'DisableCustomCosmetics')),
          asa(bool(A, 'DoCustomCosmeticValidation')),
          asa(text(G, 'CosmeticWhitelistOverride', '', { wide: true })),
        ],
      },
      {
        id: 'cluster',
        settings: [
          field('cluster_id', { wide: true }),
          bool(A, 'NoTransferFromFiltering'),
          ase(bool(A, 'usestore')),
          ase(bool(A, 'BackupTransferPlayerDatas')),
          ase(bool(A, 'converttostore')),
        ],
      },
      {
        id: 'transfer',
        settings: [
          bool(G, 'NoTributeDownloads'),
          bool(G, 'CrossARKAllowForeignDinoDownloads'),
          bool(G, 'PreventDownloadSurvivors'),
          bool(G, 'PreventUploadSurvivors'),
          bool(G, 'PreventDownloadItems'),
          bool(G, 'PreventUploadItems'),
          bool(G, 'PreventDownloadDinos'),
          bool(G, 'PreventUploadDinos'),
          num(G, 'MaxTributeDinos', 20),
          num(G, 'MaxTributeItems', 50),
          secs(G, 'TributeCharacterExpirationSeconds', 86400),
          secs(G, 'TributeItemExpirationSeconds', 86400),
          secs(G, 'TributeDinoExpirationSeconds', 86400),
          secs(G, 'MinimumDinoReuploadInterval', 0),
          ase(secs(A, 'MinimumTimeBetweenInventoryRetrieval')),
          bool(A, 'DisableCustomFoldersInTributeInventories'),
        ],
      },
    ],
  },
  {
    id: 'gameplay',
    icon: Gamepad2,
    sections: [
      {
        id: 'mode',
        settings: [
          bool(G, 'serverPVE'),
          bool(G, 'serverHardcore'),
          num(G, 'DifficultyOffset', 0.2, { step: 0.1 }),
          num(G, 'OverrideOfficialDifficulty', 5, { step: 0.1 }),
          bool(I, 'bUseSingleplayerSettings'),
        ],
      },
      {
        id: 'autoPve',
        settings: [
          bool(I, 'bAutoPvETimer'),
          bool(I, 'bAutoPvEUseSystemTime'),
          secs(I, 'AutoPvEStartTimeSeconds', 0),
          secs(I, 'AutoPvEStopTimeSeconds', 0),
        ],
      },
      {
        id: 'hud',
        settings: [
          bool(G, 'ShowMapPlayerLocation'),
          bool(G, 'allowThirdPersonPlayer'),
          bool(G, 'ServerCrosshair', true),
          bool(G, 'serverForceNoHud'),
          bool(G, 'ShowFloatingDamageText'),
          bool(G, 'AllowHitMarkers', true),
          bool(G, 'EnablePvPGamma'),
          bool(G, 'DisablePvEGamma'),
          bool(G, 'DisableWeatherFog'),
          bool(I, 'bShowCreativeMode'),
          asa(bool(I, 'bDisablePhotoMode')),
          asa(num(I, 'PhotoModeRangeLimit', 3000)),
        ],
      },
      {
        id: 'world',
        settings: [
          mult(G, 'DayCycleSpeedScale'),
          mult(G, 'DayTimeSpeedScale'),
          mult(G, 'NightTimeSpeedScale'),
          bool(G, 'RandomSupplyCratePoints'),
          bool(I, 'bDisableLootCrates'),
          bool(G, 'PreventDiseases'),
          bool(G, 'NonPermanentDiseases'),
          bool(G, 'PreventSpawnAnimations'),
          bool(G, 'AllowTekSuitPowersInGenesis'),
          num(G, 'MaxHexagonsPerCharacter', 2000000000),
          asa(bool(G, 'UseAstraeosTraversalBuff', true)),
          asa(text(G, 'WorldBossKingKaijuSpawnTime', '19:00:00')),
        ],
      },
      {
        id: 'chat',
        settings: [
          bool(G, 'globalVoiceChat'),
          bool(G, 'proximityChat'),
          bool(G, 'alwaysNotifyPlayerJoined'),
          bool(G, 'alwaysNotifyPlayerLeft'),
          bool(G, 'DontAlwaysNotifyPlayerJoined'),
          bool(G, 'bFilterChat'),
          bool(G, 'bFilterCharacterNames'),
          bool(G, 'bFilterTribeNames'),
          ase(bool(A, 'UseVivox')),
          bool(A, 'AllowChatSpam'),
        ],
      },
      {
        id: 'events',
        settings: [
          // ASA 的活动由 Mod 提供，-ActiveEvent 已经不起作用
          ase(text(A, 'ActiveEvent')),
          ase(bool(A, 'NewYearEvent')),
          ase(num(Q, 'NewYear1UTC')),
          ase(num(Q, 'NewYear2UTC')),
          ase(num(Q, 'EventColorsChanceOverride', 0, { step: 0.01 })),
          asa(bool(A, 'ServerUseEventColors')),
          asa(bool(A, 'EasterColors')),
        ],
      },
      {
        // 失落殖民地（Lost Colony）加入的玩法
        id: 'lostColony',
        settings: [
          asa(num(G, 'MaxActiveOutposts')),
          asa(num(G, 'MaxActiveCityOutposts')),
          asa(num(G, 'MaxActiveResourceCaches')),
          asa(mult(G, 'OutpostSigilRewardMultiplier')),
          asa(num(G, 'BloodforgeReinforceExtraDurability', 0.3, { step: 0.05 })),
          asa(mult(G, 'BloodforgeReinforceResourceCostMultiplier', 3)),
          asa(mult(G, 'BloodforgeReinforceSpeedMultiplier', 0.1)),
        ],
      },
    ],
  },
  {
    id: 'rates',
    icon: Gauge,
    sections: [
      {
        id: 'core',
        settings: [
          mult(G, 'XPMultiplier'),
          mult(G, 'TamingSpeedMultiplier'),
          mult(G, 'HarvestAmountMultiplier'),
          mult(G, 'HarvestHealthMultiplier'),
          mult(G, 'ResourcesRespawnPeriodMultiplier'),
          mult(G, 'ItemStackSizeMultiplier'),
        ],
      },
      {
        id: 'items',
        settings: [
          mult(I, 'CropGrowthSpeedMultiplier'),
          mult(I, 'GlobalItemDecompositionTimeMultiplier'),
          mult(I, 'GlobalCorpseDecompositionTimeMultiplier'),
          bool(G, 'ClampItemSpoilingTimes'),
          bool(G, 'ClampResourceHarvestDamage'),
          bool(G, 'UseOptimizedHarvestingHealth', true),
        ],
      },
    ],
  },
  {
    id: 'players',
    icon: User,
    sections: [
      {
        id: 'stats',
        settings: [
          mult(G, 'PlayerDamageMultiplier'),
          mult(G, 'PlayerResistanceMultiplier'),
          mult(G, 'PlayerCharacterHealthRecoveryMultiplier'),
          mult(G, 'PlayerCharacterStaminaDrainMultiplier'),
          mult(G, 'PlayerCharacterFoodDrainMultiplier'),
          mult(G, 'PlayerCharacterWaterDrainMultiplier'),
          mult(G, 'OxygenSwimSpeedStatMultiplier'),
        ],
      },
      {
        id: 'progression',
        settings: [
          bool(I, 'bAllowUnlimitedRespecs'),
          bool(I, 'bAutoUnlockAllEngrams'),
          bool(I, 'bOnlyAllowSpecifiedEngrams'),
          bool(I, 'bAllowCustomRecipes', true),
          asa(bool(I, 'bAllowSpeedLeveling')),
          bool(I, 'bAllowFlyerSpeedLeveling'),
        ],
      },
      {
        id: 'misc',
        settings: [
          bool(A, 'EnableIdlePlayerKick'),
          secs(G, 'KickIdlePlayersPeriod', 3600),
          asa(secs(G, 'ImplantSuicideCD', 28800)),
          bool(I, 'bUseCorpseLocator', true),
        ],
      },
    ],
  },
  {
    id: 'dinos',
    icon: PawPrint,
    sections: [
      {
        id: 'wild',
        settings: [
          mult(G, 'DinoCountMultiplier'),
          mult(G, 'DinoDamageMultiplier'),
          mult(G, 'DinoResistanceMultiplier'),
          mult(I, 'WildDinoCharacterFoodDrainMultiplier'),
          mult(I, 'WildDinoTorporDrainMultiplier'),
          bool(A, 'NoDinos'),
          asa(bool(A, 'NoDinosExceptForcedSpawn')),
          asa(bool(A, 'NoDinosExceptStreamingSpawn')),
          asa(bool(A, 'NoDinosExceptManualSpawn')),
          asa(bool(A, 'NoDinosExceptWaterSpawn')),
          asa(bool(A, 'NoWildBabies')),
          bool(A, 'ForceRespawnDinos'),
        ],
      },
      {
        id: 'stats',
        settings: [
          mult(G, 'DinoCharacterHealthRecoveryMultiplier'),
          mult(G, 'DinoCharacterStaminaDrainMultiplier'),
          mult(G, 'DinoCharacterFoodDrainMultiplier'),
        ],
      },
      {
        id: 'tamed',
        settings: [
          mult(G, 'TamedDinoDamageMultiplier'),
          mult(G, 'TamedDinoResistanceMultiplier'),
          num(G, 'MaxTamedDinos', 4000),
          num(G, 'MaxPersonalTamedDinos', 0),
          asa(num(G, 'MaxTamedDinos_SoftTameLimit', 5000)),
          asa(bool(G, 'DestroyTamesOverTheSoftTameLimit')),
          asa(secs(G, 'MaxTamedDinos_SoftTameLimit_CountdownForDeletionDuration', 604800)),
          bool(I, 'bDisableDinoTaming'),
          bool(I, 'bDisableDinoRiding'),
          bool(G, 'AllowRaidDinoFeeding'),
          mult(G, 'RaidDinoCharacterFoodDrainMultiplier'),
          ase(bool(A, 'UseTameEffectivenessClamp')),
        ],
      },
      {
        id: 'flyers',
        settings: [
          bool(G, 'AllowFlyerCarryPvE'),
          bool(G, 'bForceCanRideFliers'),
          bool(G, 'AllowFlyingStaminaRecovery'),
          bool(G, 'ForceFlyerExplosives'),
          bool(A, 'ForceAllowCaveFlyers'),
        ],
      },
      {
        id: 'species',
        settings: [
          asa(secs(G, 'ArmadoggoDeathCooldown', 3600)),
          asa(secs(G, 'YoungIceFoxDeathCooldown', 3600)),
          asa(num(G, 'MaxCosmoWeaponAmmo', -1)),
          asa(num(G, 'CosmoWeaponAmmoReloadAmount', 1, { step: 0.1 })),
          asa(bool(G, 'ForceGachaUnhappyInCaves', true)),
          asa(bool(G, 'DisableBurrowDecayTimers')),
        ],
      },
      {
        id: 'decay',
        settings: [
          bool(G, 'DisableDinoDecayPvE'),
          mult(G, 'PvEDinoDecayPeriodMultiplier'),
          bool(G, 'PvPDinoDecay'),
          bool(G, 'AutoDestroyDecayedDinos'),
        ],
      },
      {
        id: 'cryopod',
        settings: [
          bool(G, 'EnableCryoSicknessPVE', true),
          bool(G, 'EnableCryopodNerf'),
          secs(G, 'CryopodNerfDuration', 10),
          num(G, 'CryopodNerfDamageMult', 0.01, { step: 0.01 }),
          num(G, 'CryopodNerfIncomingDamageMultPercent', 0.25, { step: 0.01 }),
          asa(bool(G, 'DisableCryopodEnemyCheck')),
          asa(bool(G, 'DisableCryopodFridgeRequirement')),
          asa(bool(G, 'AllowCryoFridgeOnSaddle')),
          asa(hours(G, 'CryoHospitalHoursToRegenHP', 1)),
          asa(hours(G, 'CryoHospitalHoursToRegenFood', 24)),
          asa(hours(G, 'CryoHospitalHoursToDrainTorpor', 1)),
          asa(mult(G, 'CryoHospitalMatingCooldownReduction', 2)),
        ],
      },
    ],
  },
  {
    id: 'breeding',
    icon: Egg,
    sections: [
      {
        id: 'speed',
        settings: [
          mult(I, 'MatingIntervalMultiplier'),
          mult(I, 'EggHatchSpeedMultiplier'),
          mult(I, 'BabyMatureSpeedMultiplier'),
          mult(I, 'BabyFoodConsumptionSpeedMultiplier'),
        ],
      },
      {
        id: 'imprint',
        settings: [
          mult(I, 'BabyCuddleIntervalMultiplier'),
          mult(I, 'BabyCuddleGracePeriodMultiplier'),
          mult(I, 'BabyCuddleLoseImprintQualitySpeedMultiplier'),
          mult(I, 'BabyImprintingStatScaleMultiplier'),
          bool(G, 'AllowAnyoneBabyImprintCuddle'),
          bool(G, 'DisableImprintDinoBuff'),
          num(A, 'imprintlimit', 101),
        ],
      },
    ],
  },
  {
    id: 'structures',
    icon: Castle,
    sections: [
      {
        id: 'general',
        settings: [
          mult(G, 'StructureDamageMultiplier'),
          mult(G, 'StructureResistanceMultiplier'),
          secs(I, 'StructureDamageRepairCooldown', 180),
          mult(G, 'StructurePreventResourceRadiusMultiplier'),
          bool(I, 'bDisableStructurePlacementCollision'),
          bool(G, 'AllowIntegratedSPlusStructures', true),
          ase(bool(G, 'DestroyUnconnectedWaterPipes')),
        ],
      },
      {
        id: 'limits',
        settings: [
          num(G, 'TheMaxStructuresInRange', 10500),
          num(G, 'NewMaxStructuresInRange', 6000),
          num(G, 'MaxStructuresInRange', 1300),
          bool(G, 'EnableExtraStructurePreventionVolumes'),
          bool(G, 'OverrideStructurePlatformPrevention'),
          bool(G, 'AllowCaveBuildingPvE'),
          bool(G, 'AllowCaveBuildingPvP'),
          bool(G, 'PvEAllowStructuresAtSupplyDrops'),
          bool(G, 'AllowCrateSpawnsOnTopOfStructures'),
          asa(num(I, 'LimitGeneratorsNum', 3)),
          asa(num(I, 'LimitGeneratorsRange', 15000)),
        ],
      },
      {
        id: 'pickup',
        settings: [
          secs(G, 'StructurePickupTimeAfterPlacement', 30),
          num(G, 'StructurePickupHoldDuration', 0.5, { unit: 's', step: 0.1 }),
          bool(G, 'AlwaysAllowStructurePickup'),
        ],
      },
      {
        id: 'decay',
        settings: [
          bool(G, 'DisableStructureDecayPvE'),
          mult(G, 'PvEStructureDecayPeriodMultiplier'),
          bool(G, 'PvPStructureDecay'),
          bool(G, 'OnlyAutoDestroyCoreStructures'),
          bool(G, 'OnlyDecayUnsnappedCoreStructures'),
          bool(G, 'FastDecayUnsnappedCoreStructures'),
        ],
      },
      {
        id: 'platforms',
        settings: [
          num(G, 'MaxPlatformSaddleStructureLimit', 75),
          mult(G, 'PerPlatformMaxStructuresMultiplier'),
          mult(G, 'PlatformSaddleBuildAreaBoundsMultiplier'),
          bool(I, 'bAllowPlatformSaddleMultiFloors'),
          num(G, 'MaxGateFrameOnSaddles', 0),
          asa(num(G, 'MaxTrainCars', 8)),
        ],
      },
      {
        // 地堡（ASA 70.0 加入）
        id: 'bunkers',
        settings: [
          asa(bool(G, 'LimitBunkersPerTribe', true)),
          asa(num(G, 'LimitBunkersPerTribeNum', 3)),
          asa(num(G, 'MinDistanceBetweenBunkers', 3000)),
          asa(bool(G, 'AllowBunkersInPreventionZones')),
          asa(bool(G, 'AllowBunkerModulesAboveGround')),
          asa(bool(G, 'AllowBunkerModulesInPreventionZones')),
          asa(bool(G, 'AllowRidingDinosInsideBunkers', true)),
          asa(bool(G, 'AllowDinoAIInsideBunkers', true)),
          asa(num(G, 'EnemyAccessBunkerHPThreshold', 0.25, { step: 0.05 })),
          asa(mult(G, 'BunkerUnderHPThresholdDmgMultiplier', 0.05)),
        ],
      },
      {
        // TEK 专用存储的无线制作
        id: 'crafting',
        settings: [
          asa(bool(I, 'bDisableWirelessCrafting')),
          asa(bool(I, 'bDisableWirelessCraftingForPlayers')),
          asa(bool(I, 'bDisableWirelessCraftingForDinos')),
          asa(bool(I, 'bDisableWirelessCraftingForStructures')),
          asa(num(I, 'WirelessCraftingRangeOverride', 3000)),
        ],
      },
    ],
  },
  {
    id: 'tribes',
    icon: Swords,
    sections: [
      {
        id: 'tribes',
        settings: [
          num(I, 'MaxNumberOfPlayersInTribe', 0),
          num(G, 'TribeNameChangeCooldown', 15, { unit: 'min' }),
          bool(G, 'PreventTribeAlliances'),
          num(I, 'MaxAlliancesPerTribe', 10),
          num(I, 'MaxTribesPerAlliance', 10),
          bool(I, 'bPvEAllowTribeWar'),
          bool(I, 'bPvEAllowTribeWarCancel'),
          asa(mult(I, 'TribeTowerBonusMultiplier', 2)),
        ],
      },
      {
        id: 'combat',
        settings: [
          bool(I, 'bDisableFriendlyFire'),
          bool(I, 'bPvEDisableFriendlyFire'),
          bool(I, 'bPassiveDefensesDamageRiderlessDinos'),
          bool(G, 'AllowMultipleAttachedC4'),
          ase(bool(A, 'DisableRailgunPVP')),
          asa(bool(G, 'IgnorePVPMountedWeaponryRestrictions')),
        ],
      },
      {
        id: 'respawn',
        settings: [
          bool(I, 'bIncreasePvPRespawnInterval'),
          secs(I, 'IncreasePvPRespawnIntervalCheckPeriod', 300),
          mult(I, 'IncreasePvPRespawnIntervalMultiplier', 2),
          secs(I, 'IncreasePvPRespawnIntervalBaseAmount', 60),
        ],
      },
      { id: 'offline', settings: [bool(G, 'PreventOfflinePvP'), secs(G, 'PreventOfflinePvPInterval', 0)] },
    ],
  },
  {
    id: 'operations',
    icon: Wrench,
    sections: [
      // 只放游戏本身的设置。管理端和 Docker 层面的开关不在这里：ArkApi 在详情页"插件"标签，
      // 自动重启、自动更新在"管理"标签，收录站和联机房在概览的 arkfgf.com 卡片
      {
        id: 'saves',
        settings: [
          num(G, 'AutoSavePeriodMinutes', 15, { unit: 'min' }),
          ase(num(A, 'MaxNumOfSaveBackups', 20)),
          ase(bool(A, 'newsaveformat')),
          text(Q, 'AltSaveDirectoryName'),
        ],
      },
      {
        id: 'logs',
        settings: [
          bool(G, 'AdminLogging'),
          bool(A, 'servergamelog'),
          bool(A, 'servergamelogincludetribelogs'),
          bool(A, 'ServerRCONOutputTribeLogs'),
          bool(A, 'NotifyAdminCommandsInChat'),
          bool(G, 'AllowHideDamageSourceFromLogs'),
          num(G, 'RCONServerGameLogBuffer', 600),
        ],
      },
      {
        id: 'misc',
        settings: [
          {
            key: 'culture',
            source: A,
            type: 'select',
            default: '',
            options: ['ca', 'cs', 'da', 'de', 'en', 'es', 'eu', 'fi', 'fr', 'hu', 'it', 'ja', 'ka', 'ko', 'nl', 'pl', 'pt_BR', 'ru', 'sv', 'th', 'tr', 'zh', 'zh-Hans-CN', 'zh-TW'],
          },
          ase(bool(A, 'NoHangDetection')),
          // 内存保护：占用超过这个值（GB）时游戏自己存档并重启，0 关闭。用文本框：数字参数为 0 时不会传给服务器
          asa(text(A, 'GBUsageToForceRestart', '35')),
          asa(bool(G, 'UseCharacterTracker')),
          text(G, 'BanListURL', 'http://arkdedicated.com/banlist.txt', {
            wide: true,
            gameDefaults: { asa: 'https://cdn2.arkdedicated.com/asa/BanList.txt' },
          }),
          asa(text(G, 'AdminListURL', '', { wide: true })),
          text(G, 'CustomDynamicConfigUrl', '', { wide: true }),
        ],
      },
    ],
  },
  {
    id: 'advanced',
    icon: ShieldAlert,
    sections: [
      {
        id: 'antiCheat',
        settings: [
          bool(A, 'NoBattlEye'),
          ase(bool(A, 'BattlEyeServerRecheck')),
          ase(bool(A, 'insecure')),
          ase(bool(A, 'noantispeedhack')),
          ase(num(A, 'speedhackbias', 1, { step: 0.1 })),
          ase(bool(A, 'noundermeshchecking')),
          ase(bool(A, 'noundermeshkilling')),
          ase(bool(A, 'forcedisablemeshchecking')),
          ase(bool(A, 'SecureSendArKPayload')),
          ase(bool(A, 'UseItemDupeCheck')),
          ase(bool(A, 'UseSecureSpawnRules')),
          ase(bool(A, 'EnableVictoryCoreDupeCheck')),
          ase(bool(A, 'ignoredupeditems')),
          ase(bool(A, 'ForceDupeLog')),
          ase(bool(A, 'DisableDupeLogDeletes')),
          ase(bool(A, 'UseServerNetSpeedCheck')),
          ase(bool(A, 'pauseonddos')),
        ],
      },
      {
        id: 'performance',
        settings: [
          ase(bool(A, 'structurememopts')),
          ase(bool(A, 'UseStructureStasisGrid')),
          ase(bool(A, 'StasisKeepControllers')),
          ase(bool(A, 'nocombineclientmoves')),
          ase(num(A, 'DormancyNetMultiplier', 1, { step: 0.1 })),
          ase(bool(A, 'nodormancythrottling')),
          ase(bool(A, 'dedihibernation')),
          ase(bool(A, 'nitradotest2')),
          ase(bool(A, 'ServerAllowAnsel')),
          asa(bool(A, 'disabledinonetrangescaling')),
          asa(bool(A, 'AlwaysTickDedicatedSkeletalMeshes')),
          mult(G, 'ListenServerTetherDistanceMultiplier'),
          num(G, 'NPCNetworkStasisRangeScalePlayerCountStart', 0),
          num(G, 'NPCNetworkStasisRangeScalePlayerCountEnd', 0),
          num(G, 'NPCNetworkStasisRangeScalePercentEnd', 0, { step: 0.1 }),
        ],
      },
      {
        id: 'otherArgs',
        settings: [
          ase(bool(A, 'webalarm')),
          ase(bool(A, 'EnableOfficialOnlyVersioningCode')),
          ase(bool(A, 'parseservertojson')),
          ase(bool(A, 'PreventTotalConversionSaveDir')),
          ase(bool(A, 'ReloadedForBackup')),
          ase(bool(A, 'UnstasisDinoObstructionCheck')),
          ase(bool(A, 'pveallowtribewar')),
          ase(bool(A, 'pvedisallowtribewar')),
          ase(text(A, 'CustomAdminCommandTrackingURL', '', { wide: true })),
          ase(text(A, 'CustomMerticsURL', '', { wide: true })),
          ase(text(A, 'CustomNotificationURL', '', { wide: true })),
        ],
      },
      { id: 'customArgs', settings: [], panel: 'customArgs' },
      { id: 'rawConfig', settings: [], panel: 'rawConfig' },
    ],
  },
];

// ---------- 读写 ----------

export type ServerArgsValue = NonNullable<Server['server_args']>;

export const EMPTY_SERVER_ARGS: ServerArgsValue = { query_params: {}, command_line_args: {}, custom_args: [] };

export interface SettingsDoc {
  gus: IniIndex;
  game: IniIndex;
  args: ServerArgsValue;
}

function parseScalar(def: SettingDef, raw: unknown): SettingValue | undefined {
  if (raw === undefined || raw === null) return undefined;
  if (def.type === 'boolean') {
    if (typeof raw === 'boolean') return raw;
    return /^true$/i.test(String(raw).trim());
  }
  if (def.type === 'number') {
    const n = typeof raw === 'number' ? raw : Number(String(raw).trim());
    return String(raw).trim() === '' || !Number.isFinite(n) ? undefined : n;
  }
  return String(raw);
}

/** 读取设置的当前值；undefined 表示没有设置（使用默认值） */
export function readSetting(def: SettingDef, doc: SettingsDoc): SettingValue | undefined {
  switch (def.source) {
    case 'gus':
      return parseScalar(def, getIniValue(doc.gus, sectionOf(def), def.key));
    case 'game':
      return parseScalar(def, getIniValue(doc.game, sectionOf(def), def.key));
    case 'query':
      return parseScalar(def, doc.args.query_params?.[def.key]);
    case 'arg':
      return parseScalar(def, doc.args.command_line_args?.[def.key]);
    default:
      return undefined;
  }
}

/** 已设置且和默认值不同 */
export function isModified(def: SettingDef, value: SettingValue | undefined): boolean {
  if (value === undefined) return false;
  if (def.default === undefined) return value !== '';
  return value !== def.default;
}

/** 写入设置，返回要合并进表单的字段。value 为 undefined 时删除（恢复默认） */
export function writeSetting(def: SettingDef, value: SettingValue | undefined, form: Partial<Server>): Partial<Server> {
  switch (def.source) {
    case 'gus':
      return { game_user_settings: setIniValue(form.game_user_settings ?? '', sectionOf(def), def.key, value) };
    case 'game':
      return { game_ini: setIniValue(form.game_ini ?? '', sectionOf(def), def.key, value) };
    case 'query': {
      const args = form.server_args ?? EMPTY_SERVER_ARGS;
      const query_params = { ...args.query_params };
      // 查询参数的值都是字符串，空值后端会跳过
      if (value === undefined || value === '') delete query_params[def.key];
      else query_params[def.key] = typeof value === 'boolean' ? (value ? 'True' : 'False') : String(value);
      return { server_args: { ...args, query_params } };
    }
    case 'arg': {
      const args = form.server_args ?? EMPTY_SERVER_ARGS;
      const command_line_args = { ...args.command_line_args };
      // 开关参数关闭就是不传
      if (value === undefined || value === '' || value === false) delete command_line_args[def.key];
      else command_line_args[def.key] = value;
      return { server_args: { ...args, command_line_args } };
    }
    default:
      return {};
  }
}

export function allSettings(group: SettingGroupDef): SettingDef[] {
  return group.sections.flatMap((s) => s.settings);
}

/** 设置项对这个游戏是否有效 */
export function settingAppliesTo(def: SettingDef, game: Game): boolean {
  return !def.games || def.games.includes(game);
}

function withGameDefault(def: SettingDef, game: Game): SettingDef {
  const value = def.gameDefaults?.[game];
  return value === undefined ? def : { ...def, default: value };
}

// 按游戏过滤后的分组：去掉不适用的设置项，以及因此变空的分区（特殊面板保留），换上这个游戏的默认值
function filterGroups(game: Game): SettingGroupDef[] {
  return SETTING_GROUPS.map((group) => ({
    ...group,
    sections: group.sections
      .map((section) => ({
        ...section,
        settings: section.settings.filter((def) => settingAppliesTo(def, game)).map((def) => withGameDefault(def, game)),
      }))
      .filter((section) => section.panel || section.settings.length > 0),
  }));
}

const GROUPS_BY_GAME: Record<Game, SettingGroupDef[]> = { ase: filterGroups('ase'), asa: filterGroups('asa') };

/** 这个游戏的服务器在设置界面上显示的分组 */
export function settingGroupsFor(game: Game): SettingGroupDef[] {
  return GROUPS_BY_GAME[game];
}
