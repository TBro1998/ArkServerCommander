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

// field：服务器记录的字段（名称、端口、Mod 等），由表单单独渲染
// gus：GameUserSettings.ini；game：Game.ini；query：?Key=Value；arg：-Key / -Key=Value
export type SettingSource = 'field' | 'gus' | 'game' | 'query' | 'arg';
export type SettingType = 'boolean' | 'number' | 'text' | 'password' | 'select';
// x：倍率；s：秒；min：分钟
export type SettingUnit = 'x' | 's' | 'min';
export type SettingValue = boolean | number | string;

export interface SettingDef {
  key: string;
  source: SettingSource;
  type: SettingType;
  default?: SettingValue;
  unit?: SettingUnit;
  step?: number;
  options?: string[];
  // ini 的 section，缺省按 source 决定
  section?: string;
  // 占满一行
  wide?: boolean;
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
      { id: 'server', settings: [field('session_name'), field('map'), field('max_players'), field('game_branch')] },
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
      { id: 'ports', settings: [field('port'), field('query_port'), field('rcon_port')] },
      {
        id: 'platform',
        settings: [bool(A, 'crossplay'), bool(A, 'epiconly'), text(A, 'PublicIPForEpic'), text(A, 'MULTIHOME')],
      },
      {
        id: 'connections',
        settings: [bool(A, 'exclusivejoin'), num(A, 'MaxConnectionsPerIP'), bool(G, 'AllowSharedConnections')],
      },
    ],
  },
  {
    id: 'mods',
    icon: Package,
    sections: [
      { id: 'mods', settings: [field('game_mod_ids', { wide: true }), text(A, 'MapModID'), bool(A, 'automanagedmods')] },
      {
        id: 'cluster',
        settings: [
          field('cluster_id', { wide: true }),
          bool(A, 'NoTransferFromFiltering'),
          bool(A, 'usestore'),
          bool(A, 'BackupTransferPlayerDatas'),
          bool(A, 'converttostore'),
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
          secs(A, 'MinimumTimeBetweenInventoryRetrieval'),
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
          bool(A, 'UseVivox'),
          bool(A, 'AllowChatSpam'),
        ],
      },
      {
        id: 'events',
        settings: [
          text(A, 'ActiveEvent'),
          bool(A, 'NewYearEvent'),
          num(Q, 'NewYear1UTC'),
          num(Q, 'NewYear2UTC'),
          num(Q, 'EventColorsChanceOverride', 0, { step: 0.01 }),
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
        ],
      },
      {
        id: 'misc',
        settings: [
          bool(A, 'EnableIdlePlayerKick'),
          secs(G, 'KickIdlePlayersPeriod', 3600),
          secs(G, 'ImplantSuicideCD', 28800),
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
          bool(I, 'bDisableDinoTaming'),
          bool(I, 'bDisableDinoRiding'),
          bool(G, 'AllowRaidDinoFeeding'),
          mult(G, 'RaidDinoCharacterFoodDrainMultiplier'),
          bool(A, 'UseTameEffectivenessClamp'),
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
          bool(G, 'DisableCryopodEnemyCheck'),
          bool(G, 'DisableCryopodFridgeRequirement'),
          bool(G, 'AllowCryoFridgeOnSaddle'),
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
          bool(G, 'DestroyUnconnectedWaterPipes'),
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
          num(G, 'MaxPlatformSaddleStructureLimit', 25),
          mult(G, 'PerPlatformMaxStructuresMultiplier'),
          mult(G, 'PlatformSaddleBuildAreaBoundsMultiplier'),
          bool(I, 'bAllowPlatformSaddleMultiFloors'),
          num(G, 'MaxGateFrameOnSaddles', 0),
          num(G, 'MaxTrainCars', 8),
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
        ],
      },
      {
        id: 'combat',
        settings: [
          bool(I, 'bDisableFriendlyFire'),
          bool(I, 'bPvEDisableFriendlyFire'),
          bool(I, 'bPassiveDefensesDamageRiderlessDinos'),
          bool(G, 'AllowMultipleAttachedC4'),
          bool(A, 'DisableRailgunPVP'),
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
          num(A, 'MaxNumOfSaveBackups', 20),
          bool(A, 'newsaveformat'),
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
          bool(A, 'NoHangDetection'),
          text(G, 'BanListURL', 'http://arkdedicated.com/banlist.txt', { wide: true }),
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
          bool(A, 'BattlEyeServerRecheck'),
          bool(A, 'insecure'),
          bool(A, 'noantispeedhack'),
          num(A, 'speedhackbias', 1, { step: 0.1 }),
          bool(A, 'noundermeshchecking'),
          bool(A, 'noundermeshkilling'),
          bool(A, 'forcedisablemeshchecking'),
          bool(A, 'SecureSendArKPayload'),
          bool(A, 'UseItemDupeCheck'),
          bool(A, 'UseSecureSpawnRules'),
          bool(A, 'EnableVictoryCoreDupeCheck'),
          bool(A, 'ignoredupeditems'),
          bool(A, 'ForceDupeLog'),
          bool(A, 'DisableDupeLogDeletes'),
          bool(A, 'UseServerNetSpeedCheck'),
          bool(A, 'pauseonddos'),
        ],
      },
      {
        id: 'performance',
        settings: [
          bool(A, 'structurememopts'),
          bool(A, 'UseStructureStasisGrid'),
          bool(A, 'StasisKeepControllers'),
          bool(A, 'nocombineclientmoves'),
          num(A, 'DormancyNetMultiplier', 1, { step: 0.1 }),
          bool(A, 'nodormancythrottling'),
          bool(A, 'dedihibernation'),
          bool(A, 'nitradotest2'),
          bool(A, 'ServerAllowAnsel'),
          mult(G, 'ListenServerTetherDistanceMultiplier'),
          num(G, 'NPCNetworkStasisRangeScalePlayerCountStart', 0),
          num(G, 'NPCNetworkStasisRangeScalePlayerCountEnd', 0),
          num(G, 'NPCNetworkStasisRangeScalePercentEnd', 0, { step: 0.1 }),
        ],
      },
      {
        id: 'otherArgs',
        settings: [
          bool(A, 'webalarm'),
          bool(A, 'EnableOfficialOnlyVersioningCode'),
          bool(A, 'parseservertojson'),
          bool(A, 'PreventTotalConversionSaveDir'),
          bool(A, 'ReloadedForBackup'),
          bool(A, 'UnstasisDinoObstructionCheck'),
          bool(A, 'pveallowtribewar'),
          bool(A, 'pvedisallowtribewar'),
          text(A, 'CustomAdminCommandTrackingURL', '', { wide: true }),
          text(A, 'CustomMerticsURL', '', { wide: true }),
          text(A, 'CustomNotificationURL', '', { wide: true }),
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
