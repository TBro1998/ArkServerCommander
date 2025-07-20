"use client";

import { useState, useEffect, useCallback } from 'react';
import { useTranslations } from 'next-intl';
import { Textarea } from '@/components/ui/textarea';
import { Label } from '@/components/ui/label';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Alert, AlertDescription } from '@/components/ui/alert';
import { Info, CheckCircle, RefreshCw, Code, Eye, EyeOff } from 'lucide-react';
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

function getCategoryName(categoryKey: GameUserSettingsCategoryKey): string {
  const categoryNames: Record<GameUserSettingsCategoryKey, string> = {
    serverBasic: '服务器基本设置',
    gameMode: '游戏模式设置',
    communication: '聊天和通讯设置',
    gameMultipliers: '游戏倍率设置',
    characterSettings: '角色设置',
    dinoSettings: '恐龙设置',
    environmentSettings: '环境设置',
    structureSettings: '建筑设置',
    tribeSettings: '部落和联盟设置',
    breedingSettings: '繁殖和印记设置',
    itemSettings: '物品和补给设置',
    performanceSettings: '服务器性能设置',
    diseaseSettings: '疾病和状态设置',
    offlineRaidSettings: '离线突袭保护设置',
    crossArkSettings: '跨服传输设置',
    flyerSettings: '飞行载具设置',
    advancedSettings: '高级功能设置'
  };
  return categoryNames[categoryKey] || categoryKey;
}

function getParamDisplayName(paramKey: string): string {
  // Convert camelCase to readable format
  return paramKey
    .replace(/([A-Z])/g, ' $1')
    .replace(/^./, str => str.toUpperCase())
    .trim();
}

interface GameUserSettingsEditorProps {
  value?: string;
  onChange?: (value: string) => void;
}

export function GameUserSettingsEditor({ value, onChange }: GameUserSettingsEditorProps) {
  const t = useTranslations('servers.editor');
  const [editMode, setEditMode] = useState<'visual' | 'text'>('visual');
  const [isSyncing, setIsSyncing] = useState(false);
  const [textContent, setTextContent] = useState('');
  const [visualConfig, setVisualConfig] = useState<Record<string, any>>({});
  const [showPasswords, setShowPasswords] = useState<Record<string, boolean>>({});
  const [activeTab, setActiveTab] = useState<GameUserSettingsCategoryKey>('serverBasic');
  const [showTooltip, setShowTooltip] = useState<string | null>(null);

  useEffect(() => {
    setTextContent(value || '');
    parseTextToVisual(value || '');
  }, [value]);

  const parseTextToVisual = useCallback((content: string) => {
    if (!content) {
      initializeVisualConfig();
      return;
    }

    try {
      const values = extractConfigValues(content);
      if (values && Object.keys(values).length > 0) {
        initializeVisualConfig();
        Object.keys(values).forEach(key => {
          if (values[key] !== undefined && values[key] !== null) {
            setVisualConfig(prev => ({ ...prev, [key]: values[key] }));
          }
        });
      } else {
        initializeVisualConfig();
      }
    } catch (error) {
      console.error('解析配置文件失败:', error);
      initializeVisualConfig();
    }
  }, []);

  const initializeVisualConfig = useCallback(() => {
    const defaultConfig: Record<string, any> = {};
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

  const extractConfigValues = (content: string): Record<string, any> => {
    const values: Record<string, any> = {};
    const lines = content.split('\n');
    
    lines.forEach(line => {
      line = line.trim();
      if (line && !line.startsWith('[') && !line.startsWith(';') && line.includes('=')) {
        const [key, value] = line.split('=').map(s => s.trim());
        if (key && value !== undefined) {
          // Try to parse as boolean
          if (value.toLowerCase() === 'true') {
            values[key] = true;
          } else if (value.toLowerCase() === 'false') {
            values[key] = false;
          } else if (!isNaN(Number(value))) {
            values[key] = Number(value);
          } else {
            values[key] = value;
          }
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
    onChange?.(newContent);
    setIsSyncing(true);
    setTimeout(() => {
      parseTextToVisual(newContent);
      setIsSyncing(false);
    }, 800);
  };

  const handleVisualChange = (paramKey: string, value: any) => {
    setVisualConfig(prev => ({ ...prev, [paramKey]: value }));
    setIsSyncing(true);
    setTimeout(() => {
      syncVisualToText();
      setIsSyncing(false);
    }, 500);
  };

  const togglePasswordVisibility = (paramKey: string) => {
    setShowPasswords(prev => ({ ...prev, [paramKey]: !prev[paramKey] }));
  };

  const resetToDefault = () => {
    const defaultContent = `[ServerSettings]
ServerPassword=
ServerAdminPassword=
SpectatorPassword=
AdminLogging=True
serverPVE=False
serverHardcore=False
ShowMapPlayerLocation=False
allowThirdPersonPlayer=False
ServerCrosshair=True
EnablePvPGamma=False
DisablePvEGamma=False
serverForceNoHud=False
ShowFloatingDamageText=False
AllowHitMarkers=True
globalVoiceChat=False
proximityChat=False
alwaysNotifyPlayerJoined=False
alwaysNotifyPlayerLeft=False
DontAlwaysNotifyPlayerJoined=False
XPMultiplier=1.0
TamingSpeedMultiplier=1.0
HarvestAmountMultiplier=1.0
HarvestHealthMultiplier=1.0
ResourcesRespawnPeriodMultiplier=1.0
ItemStackSizeMultiplier=1.0
PlayerCharacterHealthRecoveryMultiplier=1.0
PlayerCharacterFoodDrainMultiplier=1.0
PlayerCharacterWaterDrainMultiplier=1.0
PlayerCharacterStaminaDrainMultiplier=1.0
PlayerDamageMultiplier=1.0
PlayerResistanceMultiplier=1.0
OxygenSwimSpeedStatMultiplier=1.0
ImplantSuicideCD=28800
DinoCountMultiplier=1.0
DinoCharacterHealthRecoveryMultiplier=1.0
DinoCharacterFoodDrainMultiplier=1.0
DinoCharacterStaminaDrainMultiplier=1.0
DinoDamageMultiplier=1.0
TamedDinoDamageMultiplier=1.0
DinoResistanceMultiplier=1.0
TamedDinoResistanceMultiplier=1.0
MaxTamedDinos=4000
MaxPersonalTamedDinos=0
DisableDinoDecayPvE=False
AutoDestroyDecayedDinos=False
PvEDinoDecayPeriodMultiplier=1.0
PvPDinoDecay=False
AllowRaidDinoFeeding=False
RaidDinoCharacterFoodDrainMultiplier=1.0
AllowFlyerCarryPvE=False
bForceCanRideFliers=False
DayCycleSpeedScale=1.0
DayTimeSpeedScale=1.0
NightTimeSpeedScale=1.0
DisableWeatherFog=False
DifficultyOffset=0.2
OverrideOfficialDifficulty=5.0
RandomSupplyCratePoints=False
StructureDamageMultiplier=1.0
StructureResistanceMultiplier=1.0
TheMaxStructuresInRange=10500
NewMaxStructuresInRange=6000
MaxStructuresInRange=1300
DisableStructureDecayPvE=False
PvEStructureDecayPeriodMultiplier=1.0
PvPStructureDecay=False
StructurePickupTimeAfterPlacement=30
StructurePickupHoldDuration=0.5
AlwaysAllowStructurePickup=False
OnlyAutoDestroyCoreStructures=False
OnlyDecayUnsnappedCoreStructures=False
FastDecayUnsnappedCoreStructures=False
DestroyUnconnectedWaterPipes=False
StructurePreventResourceRadiusMultiplier=1.0
MaxPlatformSaddleStructureLimit=25
PerPlatformMaxStructuresMultiplier=1.0
PlatformSaddleBuildAreaBoundsMultiplier=1.0
OverrideStructurePlatformPrevention=False
EnableExtraStructurePreventionVolumes=False
AllowCaveBuildingPvE=False
AllowCaveBuildingPvP=False
PvEAllowStructuresAtSupplyDrops=False
AllowCrateSpawnsOnTopOfStructures=False
bAllowPlatformSaddleMultiFloors=False
MaxGateFrameOnSaddles=0
MaxNumberOfPlayersInTribe=0
TribeNameChangeCooldown=15
PreventTribeAlliances=False
MaxAlliancesPerTribe=10
MaxTribesPerAlliance=10
AllowAnyoneBabyImprintCuddle=False
DisableImprintDinoBuff=False
BabyImprintingStatScaleMultiplier=1.0
ClampItemSpoilingTimes=False
ClampResourceHarvestDamage=False
UseOptimizedHarvestingHealth=True
BanListURL=http://arkdedicated.com/banlist.txt
AutoSavePeriodMinutes=15.0
KickIdlePlayersPeriod=3600
ListenServerTetherDistanceMultiplier=1.0
RCONServerGameLogBuffer=600
NPCNetworkStasisRangeScalePlayerCountStart=0
NPCNetworkStasisRangeScalePlayerCountEnd=0
NPCNetworkStasisRangeScalePercentEnd=0
PreventDiseases=False
NonPermanentDiseases=False
PreventSpawnAnimations=False
PreventOfflinePvP=False
PreventOfflinePvPInterval=900
NoTributeDownloads=False
PreventDownloadSurvivors=False
PreventDownloadItems=False
PreventDownloadDinos=False
PreventUploadSurvivors=False
PreventUploadItems=False
PreventUploadDinos=False
CrossARKAllowForeignDinoDownloads=False
MaxTributeDinos=20
MaxTributeItems=50
MinimumDinoReuploadInterval=0
TributeItemExpirationSeconds=86400
TributeDinoExpirationSeconds=86400
TributeCharacterExpirationSeconds=86400
AllowFlyingStaminaRecovery=True
ForceFlyerExplosives=False
AllowMultipleAttachedC4=False
AllowIntegratedSPlusStructures=True
AllowHideDamageSourceFromLogs=False
AllowSharedConnections=False
bFilterTribeNames=False
bFilterCharacterNames=False
bFilterChat=False
EnableCryoSicknessPVE=True
EnableCryopodNerf=False
CryopodNerfDuration=10
CryopodNerfDamageMult=0.01
CryopodNerfIncomingDamageMultPercent=0.25
DisableCryopodEnemyCheck=False
DisableCryopodFridgeRequirement=False
AllowCryoFridgeOnSaddle=False
MaxTrainCars=8
MaxHexagonsPerCharacter=2000000000
AllowTekSuitPowersInGenesis=False
CustomDynamicConfigUrl=

[/Script/Engine.GameSession]
MaxPlayers=70

[SessionSettings]
SessionName=My ARK Server

[MessageOfTheDay]
Message=欢迎来到ARK服务器！
Duration=30`;
    
    setTextContent(defaultContent);
    onChange?.(defaultContent);
  };

  const formatConfig = () => {
    // Simple formatting - ensure consistent spacing
    const lines = textContent.split('\n');
    const formattedLines: string[] = [];
    let currentSection = '';
    
    lines.forEach(line => {
      const trimmed = line.trim();
      if (trimmed.startsWith('[') && trimmed.endsWith(']')) {
        if (currentSection) formattedLines.push('');
        formattedLines.push(trimmed);
        currentSection = trimmed;
      } else if (trimmed && trimmed.includes('=')) {
        const [key, ...valueParts] = trimmed.split('=');
        const value = valueParts.join('=').trim();
        formattedLines.push(`${key.trim()}=${value}`);
      } else if (trimmed) {
        formattedLines.push(trimmed);
      }
    });
    
    const formatted = formattedLines.join('\n');
    setTextContent(formatted);
    onChange?.(formatted);
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
                              value={visualConfig[paramKey] || ''}
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
                                value={visualConfig[paramKey] || ''}
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
              <div className="flex gap-2">
                <Button
                  variant="ghost"
                  size="sm"
                  onClick={resetToDefault}
                  className="text-blue-600 hover:text-blue-800"
                >
                  {t('resetToDefault')}
                </Button>
                <Button
                  variant="ghost"
                  size="sm"
                  onClick={formatConfig}
                  className="text-green-600 hover:text-green-800"
                >
                  {t('format')}
                </Button>
              </div>
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
              className="min-h-[400px] font-mono text-sm"
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