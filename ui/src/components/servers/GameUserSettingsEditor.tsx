"use client";

import { useState, useEffect, useCallback } from 'react';
import { useTranslations } from 'next-intl';
import { Textarea } from '@/components/ui/textarea';
import { Label } from '@/components/ui/label';
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '@/components/ui/card';
import { Info, CheckCircle, RefreshCw, Eye, EyeOff } from 'lucide-react';
import { Switch } from '@/components/ui/switch';
import { Input } from '@/components/ui/input';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { Button } from '@/components/ui/button';
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from '@/components/ui/tooltip';
import { cn } from '@/lib/utils';

// GameUserSettings configuration parameters
interface GameUserSettingsParam {
  type: 'boolean' | 'number' | 'text' | 'password';
  default: boolean | number | string;
  min?: number;
  max?: number;
  step?: number;
}

type GameUserSettingsCategoryKey =
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
  | 'advancedSettings';

const gameUserSettingsParams: Record<GameUserSettingsCategoryKey, Record<string, GameUserSettingsParam>> = {
  serverBasic: {
    ServerPassword: { type: 'password', default: '' },
    ServerAdminPassword: { type: 'password', default: '' },
    SpectatorPassword: { type: 'password', default: '' },
    AdminLogging: { type: 'boolean', default: true }
  },
  gameMode: {
    serverPVE: { type: 'boolean', default: false },
    serverHardcore: { type: 'boolean', default: false },
    ShowMapPlayerLocation: { type: 'boolean', default: false },
    allowThirdPersonPlayer: { type: 'boolean', default: false },
    ServerCrosshair: { type: 'boolean', default: true },
    EnablePvPGamma: { type: 'boolean', default: false },
    DisablePvEGamma: { type: 'boolean', default: false },
    serverForceNoHud: { type: 'boolean', default: false },
    ShowFloatingDamageText: { type: 'boolean', default: false },
    AllowHitMarkers: { type: 'boolean', default: true }
  },
  communication: {
    globalVoiceChat: { type: 'boolean', default: false },
    proximityChat: { type: 'boolean', default: false },
    alwaysNotifyPlayerJoined: { type: 'boolean', default: false },
    alwaysNotifyPlayerLeft: { type: 'boolean', default: false },
    DontAlwaysNotifyPlayerJoined: { type: 'boolean', default: false }
  },
  gameMultipliers: {
    XPMultiplier: { type: 'number', default: 1.0, min: 0.1, max: 100, step: 0.1 },
    TamingSpeedMultiplier: { type: 'number', default: 1.0, min: 0.1, max: 100, step: 0.1 },
    HarvestAmountMultiplier: { type: 'number', default: 1.0, min: 0.1, max: 100, step: 0.1 },
    HarvestHealthMultiplier: { type: 'number', default: 1.0, min: 0.1, max: 10, step: 0.1 },
    ResourcesRespawnPeriodMultiplier: { type: 'number', default: 1.0, min: 0.1, max: 10, step: 0.1 },
    ItemStackSizeMultiplier: { type: 'number', default: 1.0, min: 1, max: 100, step: 1 }
  },
  characterSettings: {
    PlayerCharacterHealthRecoveryMultiplier: { type: 'number', default: 1.0, min: 0.1, max: 10, step: 0.1 },
    PlayerCharacterFoodDrainMultiplier: { type: 'number', default: 1.0, min: 0.1, max: 10, step: 0.1 },
    PlayerCharacterWaterDrainMultiplier: { type: 'number', default: 1.0, min: 0.1, max: 10, step: 0.1 },
    PlayerCharacterStaminaDrainMultiplier: { type: 'number', default: 1.0, min: 0.1, max: 10, step: 0.1 },
    PlayerDamageMultiplier: { type: 'number', default: 1.0, min: 0.1, max: 10, step: 0.1 },
    PlayerResistanceMultiplier: { type: 'number', default: 1.0, min: 0.1, max: 10, step: 0.1 },
    OxygenSwimSpeedStatMultiplier: { type: 'number', default: 1.0, min: 0.1, max: 10, step: 0.1 },
    ImplantSuicideCD: { type: 'number', default: 28800, min: 0, max: 86400 }
  },
  dinoSettings: {
    DinoCountMultiplier: { type: 'number', default: 1.0, min: 0.1, max: 10, step: 0.1 },
    DinoCharacterHealthRecoveryMultiplier: { type: 'number', default: 1.0, min: 0.1, max: 10, step: 0.1 },
    DinoCharacterFoodDrainMultiplier: { type: 'number', default: 1.0, min: 0.1, max: 10, step: 0.1 },
    DinoCharacterStaminaDrainMultiplier: { type: 'number', default: 1.0, min: 0.1, max: 10, step: 0.1 },
    DinoDamageMultiplier: { type: 'number', default: 1.0, min: 0.1, max: 10, step: 0.1 },
    TamedDinoDamageMultiplier: { type: 'number', default: 1.0, min: 0.1, max: 10, step: 0.1 },
    DinoResistanceMultiplier: { type: 'number', default: 1.0, min: 0.1, max: 10, step: 0.1 },
    TamedDinoResistanceMultiplier: { type: 'number', default: 1.0, min: 0.1, max: 10, step: 0.1 },
    MaxTamedDinos: { type: 'number', default: 4000, min: 0, max: 50000 },
    MaxPersonalTamedDinos: { type: 'number', default: 0, min: 0, max: 5000 },
    DisableDinoDecayPvE: { type: 'boolean', default: false },
    AutoDestroyDecayedDinos: { type: 'boolean', default: false },
    PvEDinoDecayPeriodMultiplier: { type: 'number', default: 1.0, min: 0.1, max: 10, step: 0.1 },
    PvPDinoDecay: { type: 'boolean', default: false },
    AllowRaidDinoFeeding: { type: 'boolean', default: false },
    RaidDinoCharacterFoodDrainMultiplier: { type: 'number', default: 1.0, min: 0.1, max: 10, step: 0.1 },
    AllowFlyerCarryPvE: { type: 'boolean', default: false },
    bForceCanRideFliers: { type: 'boolean', default: false }
  },
  environmentSettings: {
    DayCycleSpeedScale: { type: 'number', default: 1.0, min: 0.1, max: 10, step: 0.1 },
    DayTimeSpeedScale: { type: 'number', default: 1.0, min: 0.1, max: 10, step: 0.1 },
    NightTimeSpeedScale: { type: 'number', default: 1.0, min: 0.1, max: 10, step: 0.1 },
    DisableWeatherFog: { type: 'boolean', default: false },
    DifficultyOffset: { type: 'number', default: 0.2, min: 0, max: 1, step: 0.1 },
    OverrideOfficialDifficulty: { type: 'number', default: 5.0, min: 1, max: 10, step: 0.1 },
    RandomSupplyCratePoints: { type: 'boolean', default: false }
  },
  structureSettings: {
    StructureDamageMultiplier: { type: 'number', default: 1.0, min: 0.1, max: 10, step: 0.1 },
    StructureResistanceMultiplier: { type: 'number', default: 1.0, min: 0.1, max: 10, step: 0.1 },
    TheMaxStructuresInRange: { type: 'number', default: 10500, min: 1000, max: 50000 },
    NewMaxStructuresInRange: { type: 'number', default: 6000, min: 1000, max: 50000 },
    MaxStructuresInRange: { type: 'number', default: 1300, min: 100, max: 10000 },
    DisableStructureDecayPvE: { type: 'boolean', default: false },
    PvEStructureDecayPeriodMultiplier: { type: 'number', default: 1.0, min: 0.1, max: 10, step: 0.1 },
    PvPStructureDecay: { type: 'boolean', default: false },
    StructurePickupTimeAfterPlacement: { type: 'number', default: 30, min: 0, max: 3600 },
    StructurePickupHoldDuration: { type: 'number', default: 0.5, min: 0, max: 10, step: 0.1 },
    AlwaysAllowStructurePickup: { type: 'boolean', default: false },
    OnlyAutoDestroyCoreStructures: { type: 'boolean', default: false },
    OnlyDecayUnsnappedCoreStructures: { type: 'boolean', default: false },
    FastDecayUnsnappedCoreStructures: { type: 'boolean', default: false },
    DestroyUnconnectedWaterPipes: { type: 'boolean', default: false },
    StructurePreventResourceRadiusMultiplier: { type: 'number', default: 1.0, min: 0.1, max: 10, step: 0.1 },
    MaxPlatformSaddleStructureLimit: { type: 'number', default: 25, min: 1, max: 1000 },
    PerPlatformMaxStructuresMultiplier: { type: 'number', default: 1.0, min: 0.1, max: 10, step: 0.1 },
    PlatformSaddleBuildAreaBoundsMultiplier: { type: 'number', default: 1.0, min: 0.1, max: 10, step: 0.1 },
    OverrideStructurePlatformPrevention: { type: 'boolean', default: false },
    EnableExtraStructurePreventionVolumes: { type: 'boolean', default: false },
    AllowCaveBuildingPvE: { type: 'boolean', default: false },
    AllowCaveBuildingPvP: { type: 'boolean', default: false },
    PvEAllowStructuresAtSupplyDrops: { type: 'boolean', default: false },
    AllowCrateSpawnsOnTopOfStructures: { type: 'boolean', default: false },
    bAllowPlatformSaddleMultiFloors: { type: 'boolean', default: false },
    MaxGateFrameOnSaddles: { type: 'number', default: 0, min: 0, max: 100 }
  },
  tribeSettings: {
    MaxNumberOfPlayersInTribe: { type: 'number', default: 0, min: 0, max: 500 },
    TribeNameChangeCooldown: { type: 'number', default: 15, min: 0, max: 1440 },
    PreventTribeAlliances: { type: 'boolean', default: false },
    MaxAlliancesPerTribe: { type: 'number', default: 10, min: 0, max: 100 },
    MaxTribesPerAlliance: { type: 'number', default: 10, min: 0, max: 100 }
  },
  breedingSettings: {
    AllowAnyoneBabyImprintCuddle: { type: 'boolean', default: false },
    DisableImprintDinoBuff: { type: 'boolean', default: false },
    BabyImprintingStatScaleMultiplier: { type: 'number', default: 1.0, min: 0.1, max: 10, step: 0.1 }
  },
  itemSettings: {
    ClampItemSpoilingTimes: { type: 'boolean', default: false },
    ClampResourceHarvestDamage: { type: 'boolean', default: false },
    UseOptimizedHarvestingHealth: { type: 'boolean', default: true },
    BanListURL: { type: 'text', default: 'http://arkdedicated.com/banlist.txt' }
  },
  performanceSettings: {
    AutoSavePeriodMinutes: { type: 'number', default: 15.0, min: 1, max: 60, step: 0.1 },
    KickIdlePlayersPeriod: { type: 'number', default: 3600, min: 0, max: 86400 },
    ListenServerTetherDistanceMultiplier: { type: 'number', default: 1.0, min: 0.1, max: 10, step: 0.1 },
    RCONServerGameLogBuffer: { type: 'number', default: 600, min: 100, max: 10000 },
    NPCNetworkStasisRangeScalePlayerCountStart: { type: 'number', default: 0, min: 0, max: 1000 },
    NPCNetworkStasisRangeScalePlayerCountEnd: { type: 'number', default: 0, min: 0, max: 1000 },
    NPCNetworkStasisRangeScalePercentEnd: { type: 'number', default: 0, min: 0, max: 100, step: 0.1 }
  },
  diseaseSettings: {
    PreventDiseases: { type: 'boolean', default: false },
    NonPermanentDiseases: { type: 'boolean', default: false },
    PreventSpawnAnimations: { type: 'boolean', default: false }
  },
  offlineRaidSettings: {
    PreventOfflinePvP: { type: 'boolean', default: false },
    PreventOfflinePvPInterval: { type: 'number', default: 900, min: 0, max: 3600 }
  },
  crossArkSettings: {
    NoTributeDownloads: { type: 'boolean', default: false },
    PreventDownloadSurvivors: { type: 'boolean', default: false },
    PreventDownloadItems: { type: 'boolean', default: false },
    PreventDownloadDinos: { type: 'boolean', default: false },
    PreventUploadSurvivors: { type: 'boolean', default: false },
    PreventUploadItems: { type: 'boolean', default: false },
    PreventUploadDinos: { type: 'boolean', default: false },
    CrossARKAllowForeignDinoDownloads: { type: 'boolean', default: false },
    MaxTributeDinos: { type: 'number', default: 20, min: 1, max: 500 },
    MaxTributeItems: { type: 'number', default: 50, min: 1, max: 500 },
    MinimumDinoReuploadInterval: { type: 'number', default: 0, min: 0, max: 86400 },
    TributeItemExpirationSeconds: { type: 'number', default: 86400, min: 0, max: 604800 },
    TributeDinoExpirationSeconds: { type: 'number', default: 86400, min: 0, max: 604800 },
    TributeCharacterExpirationSeconds: { type: 'number', default: 86400, min: 0, max: 604800 }
  },
  flyerSettings: {
    AllowFlyingStaminaRecovery: { type: 'boolean', default: true },
    ForceFlyerExplosives: { type: 'boolean', default: false }
  },
  advancedSettings: {
    AllowMultipleAttachedC4: { type: 'boolean', default: false },
    AllowIntegratedSPlusStructures: { type: 'boolean', default: true },
    AllowHideDamageSourceFromLogs: { type: 'boolean', default: false },
    AllowSharedConnections: { type: 'boolean', default: false },
    bFilterTribeNames: { type: 'boolean', default: false },
    bFilterCharacterNames: { type: 'boolean', default: false },
    bFilterChat: { type: 'boolean', default: false },
    EnableCryoSicknessPVE: { type: 'boolean', default: true },
    EnableCryopodNerf: { type: 'boolean', default: false },
    CryopodNerfDuration: { type: 'number', default: 10, min: 0, max: 3600 },
    CryopodNerfDamageMult: { type: 'number', default: 0.01, min: 0, max: 1, step: 0.01 },
    CryopodNerfIncomingDamageMultPercent: { type: 'number', default: 0.25, min: 0, max: 1, step: 0.01 },
    DisableCryopodEnemyCheck: { type: 'boolean', default: false },
    DisableCryopodFridgeRequirement: { type: 'boolean', default: false },
    AllowCryoFridgeOnSaddle: { type: 'boolean', default: false },
    MaxTrainCars: { type: 'number', default: 8, min: 1, max: 20 },
    MaxHexagonsPerCharacter: { type: 'number', default: 2000000000, min: 0, max: 2000000000 },
    AllowTekSuitPowersInGenesis: { type: 'boolean', default: false },
    CustomDynamicConfigUrl: { type: 'text', default: '' }
  }
};

// Helper functions
function getGameUserSettingsParamsByCategory(categoryKey: GameUserSettingsCategoryKey): Record<string, GameUserSettingsParam> {
  return gameUserSettingsParams[categoryKey] || {};
}

function getAllGameUserSettingsCategories(): GameUserSettingsCategoryKey[] {
  return Object.keys(gameUserSettingsParams) as GameUserSettingsCategoryKey[];
}

// Helper functions moved inside component to access translations

interface GameUserSettingsEditorProps {
  value?: string;
  onChange?: (value: string) => void;
}

export function GameUserSettingsEditor({ value, onChange }: GameUserSettingsEditorProps) {
  const t = useTranslations('servers.editor');
  const tDefaultValues = useTranslations('servers.defaultValues');
  const tCategories = useTranslations('servers.gameUserSettingsCategories');
  const tParams = useTranslations('servers.gameUserSettingsParams');
  const [editMode, setEditMode] = useState<'visual' | 'text'>('visual');
  const [isSyncing, setIsSyncing] = useState(false);
  const [textContent, setTextContent] = useState('');
  const [visualConfig, setVisualConfig] = useState<Record<string, string | number | boolean>>({});
  const [showPasswords, setShowPasswords] = useState<Record<string, boolean>>({});
  const [activeTab, setActiveTab] = useState<GameUserSettingsCategoryKey>('serverBasic');
  const [pendingSync, setPendingSync] = useState(false);

  // Helper functions with i18n support
  const getCategoryName = (categoryKey: GameUserSettingsCategoryKey): string => {
    try {
      return tCategories(categoryKey);
    } catch {
      return categoryKey;
    }
  };

  const getParamDisplayName = (paramKey: string): string => {
    try {
      return tParams(paramKey);
    } catch {
      // Fallback to converting camelCase to readable format
      return paramKey
        .replace(/([A-Z])/g, ' $1')
        .replace(/^./, str => str.toUpperCase())
        .trim();
    }
  };



  const parseTextToVisual = useCallback((content: string) => {
    // 首先初始化默认配置
    const defaultConfig: Record<string, string | number | boolean> = {};
    getAllGameUserSettingsCategories().forEach(categoryKey => {
      const section = gameUserSettingsParams[categoryKey];
      if (section) {
        Object.keys(section).forEach(paramKey => {
          const param = section[paramKey];
          if (param && param.default !== undefined) {
            defaultConfig[paramKey] = param.default;
          }
        });
      }
    });

    // Add SessionSettings and MessageOfTheDay defaults
    defaultConfig.SessionName = tDefaultValues('sessionName');
    defaultConfig.Message = tDefaultValues('message');
    defaultConfig.Duration = 30;

    if (!content) {
      // 如果没有内容，直接使用默认配置
      setVisualConfig(defaultConfig);
      return;
    }

    try {
      const values = extractConfigValues(content);
      // 合并默认配置和解析出的配置
      const mergedConfig = { ...defaultConfig, ...values };
      setVisualConfig(mergedConfig);
    } catch (error) {
      console.error('解析配置文件失败:', error);
      // 解析失败时使用默认配置
      setVisualConfig(defaultConfig);
    }
  }, [tDefaultValues]);

  // 组件初始化时设置默认配置
  useEffect(() => {
    parseTextToVisual('');
  }, [parseTextToVisual]);

  useEffect(() => {
    setTextContent(value || '');
    // 解析后端传来的配置数据到可视化配置
    if (value) {
      parseTextToVisual(value);
    }
  }, [value, parseTextToVisual]);

  const extractConfigValues = (content: string): Record<string, string | number | boolean> => {
    const values: Record<string, string | number | boolean> = {};
    const lines = content.split('\n');

    lines.forEach(line => {
      line = line.trim();

      // 处理section
      if (line.startsWith('[') && line.endsWith(']')) {
        return;
      }

      // 跳过注释和空行
      if (!line || line.startsWith(';') || !line.includes('=')) {
        return;
      }

      const [key, ...valueParts] = line.split('=');
      const value = valueParts.join('=').trim();
      const keyName = key.trim();

      if (keyName && value !== undefined) {
        // 尝试解析为布尔值
        if (value.toLowerCase() === 'true') {
          values[keyName] = true;
        } else if (value.toLowerCase() === 'false') {
          values[keyName] = false;
        } else if (!isNaN(Number(value)) && value !== '') {
          values[keyName] = Number(value);
        } else {
          values[keyName] = value;
        }
      }
    });

    return values;
  };

  const syncVisualToTextWithConfig = useCallback((config: Record<string, string | number | boolean>) => {
    try {
      // Build INI content matching GameUserSettings.ini format
      let iniContent = '';

      // ServerSettings section
      iniContent += '[ServerSettings]\n';
      const serverBasicParams = getGameUserSettingsParamsByCategory('serverBasic');
      Object.keys(serverBasicParams).forEach(key => {
        const value = config[key];
        if (value !== undefined && value !== '') {
          iniContent += `${key}=${value}\n`;
        }
      });

      // Game mode settings
      const gameModeParams = getGameUserSettingsParamsByCategory('gameMode');
      Object.keys(gameModeParams).forEach(key => {
        const value = config[key];
        if (value !== undefined) {
          iniContent += `${key}=${value}\n`;
        }
      });

      // Communication settings
      const communicationParams = getGameUserSettingsParamsByCategory('communication');
      Object.keys(communicationParams).forEach(key => {
        const value = config[key];
        if (value !== undefined) {
          iniContent += `${key}=${value}\n`;
        }
      });

      // Game multipliers
      const gameMultipliersParams = getGameUserSettingsParamsByCategory('gameMultipliers');
      Object.keys(gameMultipliersParams).forEach(key => {
        const value = config[key];
        if (value !== undefined) {
          iniContent += `${key}=${value}\n`;
        }
      });

      // Character settings
      const characterSettingsParams = getGameUserSettingsParamsByCategory('characterSettings');
      Object.keys(characterSettingsParams).forEach(key => {
        const value = config[key];
        if (value !== undefined) {
          iniContent += `${key}=${value}\n`;
        }
      });

      // Dino settings
      const dinoSettingsParams = getGameUserSettingsParamsByCategory('dinoSettings');
      Object.keys(dinoSettingsParams).forEach(key => {
        const value = config[key];
        if (value !== undefined) {
          iniContent += `${key}=${value}\n`;
        }
      });

      // Environment settings
      const environmentSettingsParams = getGameUserSettingsParamsByCategory('environmentSettings');
      Object.keys(environmentSettingsParams).forEach(key => {
        const value = config[key];
        if (value !== undefined) {
          iniContent += `${key}=${value}\n`;
        }
      });

      // Structure settings
      const structureSettingsParams = getGameUserSettingsParamsByCategory('structureSettings');
      Object.keys(structureSettingsParams).forEach(key => {
        const value = config[key];
        if (value !== undefined) {
          iniContent += `${key}=${value}\n`;
        }
      });

      // Tribe settings
      const tribeSettingsParams = getGameUserSettingsParamsByCategory('tribeSettings');
      Object.keys(tribeSettingsParams).forEach(key => {
        const value = config[key];
        if (value !== undefined) {
          iniContent += `${key}=${value}\n`;
        }
      });

      // Breeding settings
      const breedingSettingsParams = getGameUserSettingsParamsByCategory('breedingSettings');
      Object.keys(breedingSettingsParams).forEach(key => {
        const value = config[key];
        if (value !== undefined) {
          iniContent += `${key}=${value}\n`;
        }
      });

      // Item settings
      const itemSettingsParams = getGameUserSettingsParamsByCategory('itemSettings');
      Object.keys(itemSettingsParams).forEach(key => {
        const value = config[key];
        if (value !== undefined) {
          iniContent += `${key}=${value}\n`;
        }
      });

      // Performance settings
      const performanceSettingsParams = getGameUserSettingsParamsByCategory('performanceSettings');
      Object.keys(performanceSettingsParams).forEach(key => {
        const value = config[key];
        if (value !== undefined) {
          iniContent += `${key}=${value}\n`;
        }
      });

      // Disease settings
      const diseaseSettingsParams = getGameUserSettingsParamsByCategory('diseaseSettings');
      Object.keys(diseaseSettingsParams).forEach(key => {
        const value = config[key];
        if (value !== undefined) {
          iniContent += `${key}=${value}\n`;
        }
      });

      // Offline raid settings
      const offlineRaidSettingsParams = getGameUserSettingsParamsByCategory('offlineRaidSettings');
      Object.keys(offlineRaidSettingsParams).forEach(key => {
        const value = config[key];
        if (value !== undefined) {
          iniContent += `${key}=${value}\n`;
        }
      });

      // Cross ARK settings
      const crossArkSettingsParams = getGameUserSettingsParamsByCategory('crossArkSettings');
      Object.keys(crossArkSettingsParams).forEach(key => {
        const value = config[key];
        if (value !== undefined) {
          iniContent += `${key}=${value}\n`;
        }
      });

      // Flyer settings
      const flyerSettingsParams = getGameUserSettingsParamsByCategory('flyerSettings');
      Object.keys(flyerSettingsParams).forEach(key => {
        const value = config[key];
        if (value !== undefined) {
          iniContent += `${key}=${value}\n`;
        }
      });

      // Advanced settings
      const advancedSettingsParams = getGameUserSettingsParamsByCategory('advancedSettings');
      Object.keys(advancedSettingsParams).forEach(key => {
        const value = config[key];
        if (value !== undefined) {
          iniContent += `${key}=${value}\n`;
        }
      });

      // Add other required sections
      iniContent += '\n[/Script/Engine.GameSession]\n';
      iniContent += 'MaxPlayers=70\n\n';

      iniContent += '[SessionSettings]\n';
      iniContent += `SessionName=${config.SessionName || tDefaultValues('sessionName')}\n\n`;

      iniContent += '[MessageOfTheDay]\n';
      iniContent += `Message=${config.Message || tDefaultValues('message')}\n`;
      iniContent += `Duration=${config.Duration || 30}\n`;

      setTextContent(iniContent);
      onChange?.(iniContent);
    } catch (error) {
      console.error('同步可视化配置到文本失败:', error);
    }
  }, [onChange, tDefaultValues]);

  const syncVisualToText = useCallback(() => {
    syncVisualToTextWithConfig(visualConfig);
  }, [visualConfig, syncVisualToTextWithConfig]);



  const handleTextChange = (newContent: string) => {
    setTextContent(newContent);
    onChange?.(newContent);
    setPendingSync(true);
  };

  const handleModeSwitch = async (mode: 'visual' | 'text') => {
    if (mode === editMode) return;

    setIsSyncing(true);
    
    try {
      if (mode === 'text' && editMode === 'visual') {
        // 从可视化模式切换到文本模式：将可视化配置同步到文本
        syncVisualToText();
      } else if (mode === 'visual' && editMode === 'text') {
        // 从文本模式切换到可视化模式：将文本解析到可视化配置
        parseTextToVisual(textContent);
      }
      
      setEditMode(mode);
      setPendingSync(false);
    } catch (error) {
      console.error('模式切换同步失败:', error);
    } finally {
      setIsSyncing(false);
    }
  };



  const handleVisualChange = (paramKey: string, value: string | number | boolean) => {
    setVisualConfig(prev => {
      const newConfig = { ...prev, [paramKey]: value };
      // 立即同步到文本内容并通知父组件
      setTimeout(() => {
        syncVisualToTextWithConfig(newConfig);
      }, 0);
      return newConfig;
    });
    setPendingSync(true);
  };

  const togglePasswordVisibility = (paramKey: string) => {
    setShowPasswords(prev => ({ ...prev, [paramKey]: !prev[paramKey] }));
  };

  const renderParamControl = (paramKey: string, param: GameUserSettingsParam) => {
    const currentValue = visualConfig[paramKey] ?? param.default;

    switch (param.type) {
      case 'boolean':
        return (
          <div className="flex items-center space-x-2">
            <Switch
              checked={Boolean(currentValue)}
              onCheckedChange={(checked) => handleVisualChange(paramKey, checked)}
            />
            <span className="text-sm text-muted-foreground">
              {Boolean(currentValue) ? t('enabled') : t('disabled')}
            </span>
          </div>
        );

      case 'number':
        return (
          <div className="space-y-1">
            <Input
              type="number"
              value={String(currentValue)}
              onChange={(e) => {
                const value = parseFloat(e.target.value);
                if (!isNaN(value)) {
                  handleVisualChange(paramKey, value);
                }
              }}
              min={param.min}
              max={param.max}
              step={param.step || 1}
              className="w-full"
            />
            {(param.min !== undefined || param.max !== undefined) && (
              <div className="text-xs text-muted-foreground">
                {t('range')}: {param.min ?? '∞'} - {param.max ?? '∞'}
              </div>
            )}
          </div>
        );

      case 'text':
        return (
          <Input
            type="text"
            value={String(currentValue)}
            onChange={(e) => handleVisualChange(paramKey, e.target.value)}
            className="w-full"
          />
        );

      case 'password':
        return (
          <div className="relative">
            <Input
              type={!showPasswords[paramKey] ? 'password' : 'text'}
              value={String(currentValue)}
              onChange={(e) => handleVisualChange(paramKey, e.target.value)}
              className="w-full pr-10"
            />
            <button
              type="button"
              onClick={() => togglePasswordVisibility(paramKey)}
              className="absolute inset-y-0 right-0 pr-3 flex items-center"
            >
              {showPasswords[paramKey] ? (
                <EyeOff className="h-4 w-4 text-gray-500" />
              ) : (
                <Eye className="h-4 w-4 text-gray-500" />
              )}
            </button>
          </div>
        );

      default:
        return null;
    }
  };

  const getAllCategories = () => getAllGameUserSettingsCategories();

  return (
    <div className="space-y-4">
      {/* Header with mode switch */}
      <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
        <div className="flex flex-wrap gap-2 sm:gap-4">
          <Button
            onClick={() => handleModeSwitch('visual')}
            variant={editMode === 'visual' ? 'default' : 'outline'}
            size="sm"
            className={cn(
              editMode === 'visual' && 'bg-blue-600 hover:bg-blue-700 text-white'
            )}
          >
            <i className="fas fa-sliders-h mr-2"></i>
            {t('visualEdit')}
          </Button>
          <Button
            onClick={() => handleModeSwitch('text')}
            variant={editMode === 'text' ? 'default' : 'outline'}
            size="sm"
            className={cn(
              editMode === 'text' && 'bg-blue-600 hover:bg-blue-700 text-white'
            )}
          >
            <i className="fas fa-code mr-2"></i>
            {t('textEdit')}
          </Button>

        </div>

        <div className="flex flex-col sm:flex-row items-start sm:items-center gap-2 sm:gap-4">
          <div className="text-sm text-gray-500">
            GameUserSettings.ini - {editMode === 'visual' ? t('visualEditMode') : t('textEditMode')}
          </div>
          <div className="flex items-center text-sm">
            {isSyncing ? (
              <div className="flex items-center text-blue-600">
                <RefreshCw className="h-4 w-4 animate-spin mr-2" />
                {t('syncing')}
              </div>
            ) : pendingSync ? (
              <div className="flex items-center text-yellow-600">
                <Info className="h-4 w-4 mr-1" />
                {t('pendingSync')}
              </div>
            ) : (
              <div className="flex items-center text-green-600">
                <CheckCircle className="h-4 w-4 mr-1" />
                {t('synced')}
              </div>
            )}
          </div>
        </div>
      </div>

      {/* Mode Description */}
      <div className="text-sm text-muted-foreground p-3 bg-muted/50 rounded-md">
        {editMode === 'visual' ? (
          <div>
            <p>{t('visualEditModeTip')}</p>
            <p className="mt-1 text-xs">
              {t('syncTip.visual')}
            </p>
          </div>
        ) : (
          <div>
            <p>{t('gameUserSettingsTextEditDesc')}</p>
            <p className="mt-1 text-xs">
              {t('syncTip.text')}
            </p>
          </div>
        )}
      </div>

      {/* Visual Edit Mode */}
      {editMode === 'visual' && (
        <Tabs value={activeTab} onValueChange={(value) => setActiveTab(value as GameUserSettingsCategoryKey)}>
          <TabsList className="grid w-full grid-cols-5 lg:grid-cols-9">
            {getAllCategories().map((categoryKey) => (
              <TabsTrigger key={categoryKey} value={categoryKey} className="text-xs">
                {getCategoryName(categoryKey)}
              </TabsTrigger>
            ))}
          </TabsList>

          {getAllCategories().map((categoryKey) => {
            const params = getGameUserSettingsParamsByCategory(categoryKey);
            
            return (
              <TabsContent key={categoryKey} value={categoryKey} className="space-y-4">
                <Card>
                  <CardHeader>
                    <CardTitle className="text-lg">
                      {getCategoryName(categoryKey)}
                    </CardTitle>
                    <CardDescription>
                      {Object.keys(params).length}{t('parametersCount')}
                    </CardDescription>
                  </CardHeader>
                  <CardContent className="space-y-4">
                    <div className="grid gap-4">
                      {Object.entries(params).map(([paramKey, param]) => (
                        <div key={paramKey} className="space-y-2">
                          <Label className="text-sm font-medium">
                            {getParamDisplayName(paramKey)}
                            <TooltipProvider>
                              <Tooltip>
                                <TooltipTrigger asChild>
                                  <Info className="w-4 h-4 ml-1 inline-block cursor-help" />
                                </TooltipTrigger>
                                <TooltipContent>
                                  <div className="max-w-xs">
                                    <p className="font-medium">{paramKey}</p>
                                    <p className="text-sm mt-1">{getParamDisplayName(paramKey)}</p>
                                    <p className="text-xs mt-1 text-muted-foreground">
                                      {t('defaultValue')}: {String(param.default)}
                                    </p>
                                  </div>
                                </TooltipContent>
                              </Tooltip>
                            </TooltipProvider>
                          </Label>
                          {renderParamControl(paramKey, param)}
                        </div>
                      ))}
                    </div>
                  </CardContent>
                </Card>
              </TabsContent>
            );
          })}
        </Tabs>
      )}

      {/* Text Edit Mode */}
      {editMode === 'text' && (
        <div className="space-y-2">
          <Label htmlFor="gameusersettings-text">{t('textEditMode')}</Label>
          <Textarea
            id="gameusersettings-text"
            value={textContent}
            onChange={(e) => handleTextChange(e.target.value)}
            className="min-h-[400px] font-mono text-sm"
            placeholder={t('placeholders.gameUserSettings')}
          />
        </div>
      )}
    </div>
  );
}