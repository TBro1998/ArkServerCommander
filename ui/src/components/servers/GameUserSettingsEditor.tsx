"use client";

import { useState, useEffect, useCallback } from 'react';
import { useTranslations } from 'next-intl';
import { Textarea } from '@/components/ui/textarea';
import { Label } from '@/components/ui/label';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Alert, AlertDescription } from '@/components/ui/alert';
import { Info, CheckCircle, RefreshCw, Eye, EyeOff } from 'lucide-react';
import { Switch } from '@/components/ui/switch';
import { Input } from '@/components/ui/input';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { Button } from '@/components/ui/button';
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
  const tCategories = useTranslations('servers.gameUserSettingsCategories');
  const tParams = useTranslations('servers.gameUserSettingsParams');
  const [editMode, setEditMode] = useState<'visual' | 'text'>('visual');
  const [isSyncing, setIsSyncing] = useState(false);
  const [textContent, setTextContent] = useState('');
  const [visualConfig, setVisualConfig] = useState<Record<string, string | number | boolean>>({});
  const [showPasswords, setShowPasswords] = useState<Record<string, boolean>>({});
  const [activeTab, setActiveTab] = useState<GameUserSettingsCategoryKey>('serverBasic');
  const [showTooltip, setShowTooltip] = useState<string | null>(null);
  const [isUserEditing, setIsUserEditing] = useState(false);
  const [lastUserEditTime, setLastUserEditTime] = useState(0);

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

  const initializeVisualConfig = useCallback(() => {
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
    defaultConfig.SessionName = 'My ARK Server';
    defaultConfig.Message = '欢迎来到ARK服务器！';
    defaultConfig.Duration = 30;

    setVisualConfig(prev => ({ ...prev, ...defaultConfig }));
  }, []);

  const parseTextToVisual = useCallback((content: string) => {
    if (!content) {
      initializeVisualConfig();
      return;
    }

    try {
      const values = extractConfigValues(content);
      if (values && Object.keys(values).length > 0) {
        // 只更新实际发生变化的值，避免覆盖用户输入
        setVisualConfig(prev => {
          const newConfig = { ...prev };
          let hasChanges = false;

          Object.entries(values).forEach(([key, value]) => {
            if (value !== undefined && value !== null && newConfig[key] !== value) {
              newConfig[key] = value;
              hasChanges = true;
            }
          });

          return hasChanges ? newConfig : prev;
        });
      }
    } catch (error) {
      console.error('解析配置文件失败:', error);
    }
  }, [initializeVisualConfig]);

  useEffect(() => {
    setTextContent(value || '');
    // 只有在用户没有正在编辑时才更新可视化配置
    const now = Date.now();
    if (now - lastUserEditTime > 1000 && !isUserEditing) {
      parseTextToVisual(value || '');
    }
  }, [value, lastUserEditTime, isUserEditing, parseTextToVisual]);

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

  const syncVisualToText = useCallback(() => {
    try {
      // Build INI content matching GameUserSettings.ini format
      let iniContent = '';

      // ServerSettings section
      iniContent += '[ServerSettings]\n';
      const serverBasicParams = getGameUserSettingsParamsByCategory('serverBasic');
      Object.keys(serverBasicParams).forEach(key => {
        const value = visualConfig[key];
        if (value !== undefined && value !== '') {
          iniContent += `${key}=${value}\n`;
        }
      });

      // Game mode settings
      const gameModeParams = getGameUserSettingsParamsByCategory('gameMode');
      Object.keys(gameModeParams).forEach(key => {
        const value = visualConfig[key];
        if (value !== undefined) {
          iniContent += `${key}=${value}\n`;
        }
      });

      // Communication settings
      const communicationParams = getGameUserSettingsParamsByCategory('communication');
      Object.keys(communicationParams).forEach(key => {
        const value = visualConfig[key];
        if (value !== undefined) {
          iniContent += `${key}=${value}\n`;
        }
      });

      // Game multipliers
      const gameMultipliersParams = getGameUserSettingsParamsByCategory('gameMultipliers');
      Object.keys(gameMultipliersParams).forEach(key => {
        const value = visualConfig[key];
        if (value !== undefined) {
          iniContent += `${key}=${value}\n`;
        }
      });

      // Character settings
      const characterSettingsParams = getGameUserSettingsParamsByCategory('characterSettings');
      Object.keys(characterSettingsParams).forEach(key => {
        const value = visualConfig[key];
        if (value !== undefined) {
          iniContent += `${key}=${value}\n`;
        }
      });

      // Dino settings
      const dinoSettingsParams = getGameUserSettingsParamsByCategory('dinoSettings');
      Object.keys(dinoSettingsParams).forEach(key => {
        const value = visualConfig[key];
        if (value !== undefined) {
          iniContent += `${key}=${value}\n`;
        }
      });

      // Environment settings
      const environmentSettingsParams = getGameUserSettingsParamsByCategory('environmentSettings');
      Object.keys(environmentSettingsParams).forEach(key => {
        const value = visualConfig[key];
        if (value !== undefined) {
          iniContent += `${key}=${value}\n`;
        }
      });

      // Structure settings
      const structureSettingsParams = getGameUserSettingsParamsByCategory('structureSettings');
      Object.keys(structureSettingsParams).forEach(key => {
        const value = visualConfig[key];
        if (value !== undefined) {
          iniContent += `${key}=${value}\n`;
        }
      });

      // Tribe settings
      const tribeSettingsParams = getGameUserSettingsParamsByCategory('tribeSettings');
      Object.keys(tribeSettingsParams).forEach(key => {
        const value = visualConfig[key];
        if (value !== undefined) {
          iniContent += `${key}=${value}\n`;
        }
      });

      // Breeding settings
      const breedingSettingsParams = getGameUserSettingsParamsByCategory('breedingSettings');
      Object.keys(breedingSettingsParams).forEach(key => {
        const value = visualConfig[key];
        if (value !== undefined) {
          iniContent += `${key}=${value}\n`;
        }
      });

      // Item settings
      const itemSettingsParams = getGameUserSettingsParamsByCategory('itemSettings');
      Object.keys(itemSettingsParams).forEach(key => {
        const value = visualConfig[key];
        if (value !== undefined) {
          iniContent += `${key}=${value}\n`;
        }
      });

      // Performance settings
      const performanceSettingsParams = getGameUserSettingsParamsByCategory('performanceSettings');
      Object.keys(performanceSettingsParams).forEach(key => {
        const value = visualConfig[key];
        if (value !== undefined) {
          iniContent += `${key}=${value}\n`;
        }
      });

      // Disease settings
      const diseaseSettingsParams = getGameUserSettingsParamsByCategory('diseaseSettings');
      Object.keys(diseaseSettingsParams).forEach(key => {
        const value = visualConfig[key];
        if (value !== undefined) {
          iniContent += `${key}=${value}\n`;
        }
      });

      // Offline raid settings
      const offlineRaidSettingsParams = getGameUserSettingsParamsByCategory('offlineRaidSettings');
      Object.keys(offlineRaidSettingsParams).forEach(key => {
        const value = visualConfig[key];
        if (value !== undefined) {
          iniContent += `${key}=${value}\n`;
        }
      });

      // Cross ARK settings
      const crossArkSettingsParams = getGameUserSettingsParamsByCategory('crossArkSettings');
      Object.keys(crossArkSettingsParams).forEach(key => {
        const value = visualConfig[key];
        if (value !== undefined) {
          iniContent += `${key}=${value}\n`;
        }
      });

      // Flyer settings
      const flyerSettingsParams = getGameUserSettingsParamsByCategory('flyerSettings');
      Object.keys(flyerSettingsParams).forEach(key => {
        const value = visualConfig[key];
        if (value !== undefined) {
          iniContent += `${key}=${value}\n`;
        }
      });

      // Advanced settings
      const advancedSettingsParams = getGameUserSettingsParamsByCategory('advancedSettings');
      Object.keys(advancedSettingsParams).forEach(key => {
        const value = visualConfig[key];
        if (value !== undefined) {
          iniContent += `${key}=${value}\n`;
        }
      });

      // Add other required sections
      iniContent += '\n[/Script/Engine.GameSession]\n';
      iniContent += 'MaxPlayers=70\n\n';

      iniContent += '[SessionSettings]\n';
      iniContent += `SessionName=${visualConfig.SessionName || 'My ARK Server'}\n\n`;

      iniContent += '[MessageOfTheDay]\n';
      iniContent += `Message=${visualConfig.Message || '欢迎来到ARK服务器！'}\n`;
      iniContent += `Duration=${visualConfig.Duration || 30}\n`;

      setTextContent(iniContent);
      onChange?.(iniContent);
    } catch (error) {
      console.error('同步可视化配置到文本失败:', error);
    }
  }, [visualConfig, onChange]);

  const handleTextChange = (newContent: string) => {
    setTextContent(newContent);
    setIsUserEditing(true);
    setLastUserEditTime(Date.now());
    onChange?.(newContent);
    setIsSyncing(true);

    // 使用防抖机制避免快速连续修改
    if ((window as Window & { __textSyncTimeout?: NodeJS.Timeout }).__textSyncTimeout) {
      clearTimeout((window as Window & { __textSyncTimeout?: NodeJS.Timeout }).__textSyncTimeout);
    }
    (window as Window & { __textSyncTimeout?: NodeJS.Timeout }).__textSyncTimeout = setTimeout(() => {
      parseTextToVisual(newContent);
      setIsSyncing(false);
      setIsUserEditing(false);
    }, 500);
  };

  const handleVisualChange = (paramKey: string, value: string | number | boolean) => {
    setVisualConfig(prev => ({ ...prev, [paramKey]: value }));
    setIsUserEditing(true);
    setLastUserEditTime(Date.now());
    setIsSyncing(true);

    // 使用防抖机制避免快速连续修改
    if ((window as Window & { __visualSyncTimeout?: NodeJS.Timeout }).__visualSyncTimeout) {
      clearTimeout((window as Window & { __visualSyncTimeout?: NodeJS.Timeout }).__visualSyncTimeout);
    }
    (window as Window & { __visualSyncTimeout?: NodeJS.Timeout }).__visualSyncTimeout = setTimeout(() => {
      syncVisualToText();
      setIsSyncing(false);
      setIsUserEditing(false);
    }, 500);
  };

  const togglePasswordVisibility = (paramKey: string) => {
    setShowPasswords(prev => ({ ...prev, [paramKey]: !prev[paramKey] }));
  };


  const getAllCategories = () => getAllGameUserSettingsCategories();

  return (
    <div className="space-y-4">
      {/* Header with mode switch */}
      <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
        <div className="flex flex-wrap gap-2 sm:gap-4">
          <Button
            onClick={() => setEditMode('visual')}
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
            onClick={() => setEditMode('text')}
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
            ) : (
              <div className="flex items-center text-green-600">
                <CheckCircle className="h-4 w-4 mr-1" />
                {t('synced')}
              </div>
            )}
          </div>
        </div>
      </div>

      {/* Visual Edit Mode */}
      {editMode === 'visual' && (
        <div className="space-y-6">
          <Tabs value={activeTab} onValueChange={(value) => setActiveTab(value as GameUserSettingsCategoryKey)}>
            <TabsList className="flex flex-wrap h-auto">
              {getAllCategories().map((categoryKey) => (
                <TabsTrigger key={categoryKey} value={categoryKey} className="text-sm">
                  {getCategoryName(categoryKey)}
                </TabsTrigger>
              ))}
            </TabsList>

            {getAllCategories().map((categoryKey) => (
              <TabsContent key={categoryKey} value={categoryKey} className="mt-4">
                <Card>
                  <CardContent className="pt-6">
                    <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
                      {Object.entries(getGameUserSettingsParamsByCategory(categoryKey)).map(([paramKey, param]) => (
                        <div key={paramKey} className="space-y-2">
                          <Label className="text-sm font-medium">
                            {getParamDisplayName(paramKey)}
                            <div className="relative inline-block ml-1">
                              <Info
                                className="h-3 w-3 text-gray-400 cursor-help"
                                onMouseEnter={() => setShowTooltip(paramKey)}
                                onMouseLeave={() => setShowTooltip(null)}
                              />
                              {showTooltip === paramKey && (
                                <div className="absolute z-30 bg-gray-900 text-white text-xs rounded py-2 px-3 max-w-xs whitespace-normal shadow-lg"
                                  style={{ bottom: '100%', left: '50%', transform: 'translateX(-50%)', marginBottom: '8px' }}>
                                  {getParamDisplayName(paramKey)}
                                </div>
                              )}
                            </div>
                          </Label>

                          {param.type === 'boolean' && (
                            <div className="flex items-center space-x-2">
                              <Switch
                                checked={Boolean(visualConfig[paramKey])}
                                onCheckedChange={(checked) => handleVisualChange(paramKey, checked)}
                              />
                              <span className="text-sm">
                                {visualConfig[paramKey] ? t('enabled') : t('disabled')}
                              </span>
                            </div>
                          )}

                          {param.type === 'number' && (
                            <Input
                              type="number"
                              value={(visualConfig[paramKey]?.toString() || '')}
                              onChange={(e) => handleVisualChange(paramKey, parseFloat(e.target.value) || 0)}
                              min={param.min}
                              max={param.max}
                              step={param.step || 1}
                              className="w-full"
                            />
                          )}


                          {(param.type === 'text' || param.type === 'password') && (
                            <div className="relative">
                              <Input
                                type={param.type === 'password' && !showPasswords[paramKey] ? 'password' : 'text'}
                                value={visualConfig[paramKey]?.toString() || ''}
                                onChange={(e) => handleVisualChange(paramKey, e.target.value)}
                                className="w-full"
                              />
                              {param.type === 'password' && (
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
                              )}
                            </div>
                          )}
                        </div>
                      ))}
                    </div>
                  </CardContent>
                </Card>
              </TabsContent>
            ))}
          </Tabs>
        </div>
      )}

      {/* Text Edit Mode */}
      {editMode === 'text' && (
        <Card>
          <CardHeader>
            <CardTitle className="flex justify-between items-center">
              <span>{t('gameUserSettingsContent')}</span>
            </CardTitle>
          </CardHeader>
          <CardContent>
            <Textarea
              value={textContent}
              onChange={(e: React.ChangeEvent<HTMLTextAreaElement>) => handleTextChange(e.target.value)}
              placeholder="[ServerSettings]
SessionName=My ARK Server
ServerPassword=
MaxPlayers=70

[SessionSettings]
SessionName=My ARK Server

[MessageOfTheDay]
Message=欢迎来到ARK服务器！

[/Script/Engine.GameSession]
MaxPlayers=70"
              className="h-96 font-mono text-sm overflow-y-auto"
            />
            <p className="text-xs text-gray-500 mt-2">
              {t('gameUserSettingsDescription')}
            </p>
          </CardContent>
        </Card>
      )}

      {/* Info Alert */}
      {editMode === 'visual' && (
        <Alert>
          <Info className="h-4 w-4" />
          <AlertDescription>
            <p className="font-medium">{t('visualEditModeDesc')}</p>
            <p className="text-sm">{t('visualEditModeTip')}</p>
          </AlertDescription>
        </Alert>
      )}
    </div>
  );
}