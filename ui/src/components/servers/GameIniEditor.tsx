"use client";

import { useState, useCallback } from 'react';
import { useTranslations } from 'next-intl';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Switch } from '@/components/ui/switch';
import { Textarea } from '@/components/ui/textarea';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { mergeIni } from '@/lib/ini';

import { Info, Eye, EyeOff } from 'lucide-react';
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from '@/components/ui/tooltip';

// Game.ini configuration parameters
interface GameIniParam {
  type: 'boolean' | 'number' | 'text';
  default: boolean | number | string;
  min?: number;
  max?: number;
  step?: number;
}

type GameIniCategoryKey =
  | 'gameBasic'
  | 'experienceSettings'
  | 'breedingSettings'
  | 'itemSettings'
  | 'dinoSettings'
  | 'tribeSettings'
  | 'pvpSettings'
  | 'structureSettings'
  | 'advancedSettings'
  | 'customSettings';

const gameIniParams: Record<GameIniCategoryKey, Record<string, GameIniParam>> = {
  gameBasic: {
    bUseSingleplayerSettings: { type: 'boolean', default: false },
    bDisableStructurePlacementCollision: { type: 'boolean', default: false },
    bAllowFlyerCarryPvE: { type: 'boolean', default: true },
    bDisableStructureDecayPvE: { type: 'boolean', default: false },
    bAllowUnlimitedRespecs: { type: 'boolean', default: true },
    bAllowPlatformSaddleMultiFloors: { type: 'boolean', default: true },
    bPassiveDefensesDamageRiderlessDinos: { type: 'boolean', default: true },
    bPvEDisableFriendlyFire: { type: 'boolean', default: false },
    bDisableFriendlyFire: { type: 'boolean', default: false },
    bEnablePvPGamma: { type: 'boolean', default: false },
    DifficultyOffset: { type: 'number', default: 1.0, min: 0.0, max: 1.0, step: 0.1 },
    OverrideOfficialDifficulty: { type: 'number', default: 5.0, min: 1.0, max: 10.0, step: 0.1 }
  },
  experienceSettings: {
    XPMultiplier: { type: 'number', default: 1.0, min: 0.1, max: 10.0, step: 0.1 },
    PlayerCharacterWaterDrainMultiplier: { type: 'number', default: 1.0, min: 0.1, max: 10.0, step: 0.1 },
    PlayerCharacterFoodDrainMultiplier: { type: 'number', default: 1.0, min: 0.1, max: 10.0, step: 0.1 },
    PlayerCharacterStaminaDrainMultiplier: { type: 'number', default: 1.0, min: 0.1, max: 10.0, step: 0.1 },
    PlayerCharacterHealthRecoveryMultiplier: { type: 'number', default: 1.0, min: 0.1, max: 10.0, step: 0.1 }
  },
  breedingSettings: {
    MatingIntervalMultiplier: { type: 'number', default: 1.0, min: 0.1, max: 10.0, step: 0.1 },
    EggHatchSpeedMultiplier: { type: 'number', default: 1.0, min: 0.1, max: 50.0, step: 0.1 },
    BabyMatureSpeedMultiplier: { type: 'number', default: 1.0, min: 0.1, max: 50.0, step: 0.1 },
    BabyFoodConsumptionSpeedMultiplier: { type: 'number', default: 1.0, min: 0.1, max: 10.0, step: 0.1 },
    BabyCuddleIntervalMultiplier: { type: 'number', default: 1.0, min: 0.1, max: 10.0, step: 0.1 },
    BabyCuddleGracePeriodMultiplier: { type: 'number', default: 1.0, min: 0.1, max: 10.0, step: 0.1 },
    BabyCuddleLoseImprintQualitySpeedMultiplier: { type: 'number', default: 1.0, min: 0.1, max: 10.0, step: 0.1 }
  },
  itemSettings: {
    HarvestAmountMultiplier: { type: 'number', default: 1.0, min: 0.1, max: 10.0, step: 0.1 },
    HarvestHealthMultiplier: { type: 'number', default: 1.0, min: 0.1, max: 10.0, step: 0.1 },
    ResourcesRespawnPeriodMultiplier: { type: 'number', default: 1.0, min: 0.1, max: 10.0, step: 0.1 },
    ItemStackSizeMultiplier: { type: 'number', default: 1.0, min: 0.1, max: 10.0, step: 0.1 },
    CropGrowthSpeedMultiplier: { type: 'number', default: 1.0, min: 0.1, max: 10.0, step: 0.1 },
    GlobalItemDecompositionTimeMultiplier: { type: 'number', default: 1.0, min: 0.1, max: 10.0, step: 0.1 },
    GlobalCorpseDecompositionTimeMultiplier: { type: 'number', default: 1.0, min: 0.1, max: 10.0, step: 0.1 }
  },
  dinoSettings: {
    TamingSpeedMultiplier: { type: 'number', default: 1.0, min: 0.1, max: 50.0, step: 0.1 },
    DinoCharacterFoodDrainMultiplier: { type: 'number', default: 1.0, min: 0.1, max: 10.0, step: 0.1 },
    DinoCharacterStaminaDrainMultiplier: { type: 'number', default: 1.0, min: 0.1, max: 10.0, step: 0.1 },
    DinoCharacterHealthRecoveryMultiplier: { type: 'number', default: 1.0, min: 0.1, max: 10.0, step: 0.1 },
    DinoCountMultiplier: { type: 'number', default: 1.0, min: 0.1, max: 10.0, step: 0.1 },
    WildDinoCharacterFoodDrainMultiplier: { type: 'number', default: 1.0, min: 0.1, max: 10.0, step: 0.1 },
    WildDinoTorporDrainMultiplier: { type: 'number', default: 1.0, min: 0.1, max: 10.0, step: 0.1 }
  },
  tribeSettings: {
    MaxNumberOfPlayersInTribe: { type: 'number', default: 0, min: 0, max: 100, step: 1 },
    TribeNameChangeCooldown: { type: 'number', default: 15, min: 0, max: 1440, step: 1 },
    bPvEAllowTribeWar: { type: 'boolean', default: false },
    bPvEAllowTribeWarCancel: { type: 'boolean', default: false }
  },
  pvpSettings: {
    bIncreasePvPRespawnInterval: { type: 'boolean', default: false },
    IncreasePvPRespawnIntervalCheckPeriod: { type: 'number', default: 300, min: 60, max: 3600, step: 60 },
    IncreasePvPRespawnIntervalMultiplier: { type: 'number', default: 1.0, min: 1.0, max: 5.0, step: 0.1 },
    IncreasePvPRespawnIntervalBaseAmount: { type: 'number', default: 60, min: 30, max: 600, step: 30 }
  },
  structureSettings: {
    StructureDamageMultiplier: { type: 'number', default: 1.0, min: 0.1, max: 10.0, step: 0.1 },
    StructureResistanceMultiplier: { type: 'number', default: 1.0, min: 0.1, max: 10.0, step: 0.1 },
    StructureDamageRepairCooldown: { type: 'number', default: 180, min: 0, max: 3600, step: 60 },
    PvEStructureDecayPeriodMultiplier: { type: 'number', default: 1.0, min: 0.1, max: 10.0, step: 0.1 },
    MaxStructuresInRange: { type: 'number', default: 1300, min: 100, max: 10000, step: 100 }
  },
  advancedSettings: {
    bAutoPvETimer: { type: 'boolean', default: false },
    bAutoPvEUseSystemTime: { type: 'boolean', default: false },
    AutoPvEStartTimeSeconds: { type: 'number', default: 0, min: 0, max: 86400, step: 3600 },
    AutoPvEStopTimeSeconds: { type: 'number', default: 0, min: 0, max: 86400, step: 3600 },
    bOnlyAllowSpecifiedEngrams: { type: 'boolean', default: false },
    bAutoUnlockAllEngrams: { type: 'boolean', default: false },
    bShowCreativeMode: { type: 'boolean', default: false },
    bUseCorpseLocator: { type: 'boolean', default: false },
    bDisableLootCrates: { type: 'boolean', default: false },
    bDisableDinoRiding: { type: 'boolean', default: false },
    bDisableDinoTaming: { type: 'boolean', default: false },
    bAllowCustomRecipes: { type: 'boolean', default: true }
  },
  customSettings: {
    DayCycleSpeedScale: { type: 'number', default: 1.0, min: 0.1, max: 10.0, step: 0.1 },
    NightTimeSpeedScale: { type: 'number', default: 1.0, min: 0.1, max: 10.0, step: 0.1 },
    DayTimeSpeedScale: { type: 'number', default: 1.0, min: 0.1, max: 10.0, step: 0.1 }
  }
};

function getGameIniParamsByCategory(category: GameIniCategoryKey): Record<string, GameIniParam> {
  return gameIniParams[category] || {};
}



function buildDefaultConfig(): Record<string, string | number | boolean> {
  const defaultConfig: Record<string, string | number | boolean> = {};
  Object.keys(gameIniParams).forEach(categoryKey => {
    const params = getGameIniParamsByCategory(categoryKey as GameIniCategoryKey);
    Object.keys(params).forEach(paramKey => {
      defaultConfig[paramKey] = params[paramKey].default;
    });
  });
  return defaultConfig;
}

// 生成默认的 Game.ini 模板，只在文件为空（新建服务器）时使用
function buildDefaultGameIniText(visualConfig: Record<string, string | number | boolean>): string {
  let iniContent = '';

  // Add sections based on categories
  iniContent += '[/script/shootergame.shootergamemode]\n';

  // Basic game settings
  const gameBasicParams = getGameIniParamsByCategory('gameBasic');
  Object.keys(gameBasicParams).forEach(key => {
    const value = visualConfig[key];
    if (value !== undefined) {
      iniContent += `${key}=${value}\n`;
    }
  });

  // Tribe settings
  const tribeParams = getGameIniParamsByCategory('tribeSettings');
  Object.keys(tribeParams).forEach(key => {
    const value = visualConfig[key];
    if (value !== undefined) {
      iniContent += `${key}=${value}\n`;
    }
  });

  iniContent += '\n[/Script/ShooterGame.ShooterGameMode]\n';

  // Experience settings
  const expParams = getGameIniParamsByCategory('experienceSettings');
  Object.keys(expParams).forEach(key => {
    const value = visualConfig[key];
    if (value !== undefined) {
      iniContent += `${key}=${value}\n`;
    }
  });

  // Item settings
  const itemParams = getGameIniParamsByCategory('itemSettings');
  Object.keys(itemParams).forEach(key => {
    const value = visualConfig[key];
    if (value !== undefined) {
      iniContent += `${key}=${value}\n`;
    }
  });

  // Dino settings
  const dinoParams = getGameIniParamsByCategory('dinoSettings');
  Object.keys(dinoParams).forEach(key => {
    const value = visualConfig[key];
    if (value !== undefined) {
      iniContent += `${key}=${value}\n`;
    }
  });

  // Breeding settings
  const breedingParams = getGameIniParamsByCategory('breedingSettings');
  Object.keys(breedingParams).forEach(key => {
    const value = visualConfig[key];
    if (value !== undefined) {
      iniContent += `${key}=${value}\n`;
    }
  });

  // PvP settings
  const pvpParams = getGameIniParamsByCategory('pvpSettings');
  Object.keys(pvpParams).forEach(key => {
    const value = visualConfig[key];
    if (value !== undefined) {
      iniContent += `${key}=${value}\n`;
    }
  });

  // Structure settings
  const structureParams = getGameIniParamsByCategory('structureSettings');
  Object.keys(structureParams).forEach(key => {
    const value = visualConfig[key];
    if (value !== undefined) {
      iniContent += `${key}=${value}\n`;
    }
  });

  // Advanced settings
  const advancedParams = getGameIniParamsByCategory('advancedSettings');
  Object.keys(advancedParams).forEach(key => {
    const value = visualConfig[key];
    if (value !== undefined) {
      iniContent += `${key}=${value}\n`;
    }
  });

  // Custom settings
  const customParams = getGameIniParamsByCategory('customSettings');
  Object.keys(customParams).forEach(key => {
    const value = visualConfig[key];
    if (value !== undefined) {
      iniContent += `${key}=${value}\n`;
    }
  });
  return iniContent;
}

interface GameIniEditorProps {
  value?: string;
  onChange?: (value: string) => void;
}

export function GameIniEditor({ value, onChange }: GameIniEditorProps) {
  const t = useTranslations('servers.editor');
  const tCategories = useTranslations('servers.gameIniCategories');
  const tParams = useTranslations('servers.gameIniParams');
  const [editMode, setEditMode] = useState<'visual' | 'text'>('visual');
  const [textContent, setTextContent] = useState('');
  const [visualConfig, setVisualConfig] = useState<Record<string, string | number | boolean>>({});
  const [activeTab, setActiveTab] = useState<GameIniCategoryKey>('gameBasic');
  // 在可视化界面里改动过的键，保存时只合并这些键
  const [dirtyKeys, setDirtyKeys] = useState<Set<string>>(() => new Set());

  // Get parameter display name from translations
  const getParamDisplayName = (paramKey: string): string => {
    try {
      return tParams(paramKey);
    } catch {
      return paramKey; // Fallback to parameter key if translation not found
    }
  };

  const parseTextToVisual = useCallback((text: string) => {
    try {
      const config: Record<string, string | number | boolean> = {};
      const lines = text.split('\n');

      for (const line of lines) {
        const trimmedLine = line.trim();
        if (trimmedLine && !trimmedLine.startsWith(';') && !trimmedLine.startsWith('#') && !trimmedLine.startsWith('[')) {
          const [key, ...valueParts] = trimmedLine.split('=');
          if (key && valueParts.length > 0) {
            const value = valueParts.join('=').trim();

            // Find parameter definition
            let paramDef: GameIniParam | undefined;
            for (const categoryKey of Object.keys(gameIniParams)) {
              const category = categoryKey as GameIniCategoryKey;
              const params = getGameIniParamsByCategory(category);
              if (params[key.trim()]) {
                paramDef = params[key.trim()];
                break;
              }
            }

            if (paramDef) {
              if (paramDef.type === 'boolean') {
                config[key.trim()] = value.toLowerCase() === 'true';
              } else if (paramDef.type === 'number') {
                const numValue = parseFloat(value);
                config[key.trim()] = isNaN(numValue) ? paramDef.default : numValue;
              } else {
                config[key.trim()] = value;
              }
            }
          }
        }
      }

      setVisualConfig(config);
    } catch (error) {
      console.error(t('parseGameIniError') + ':', error);
    }
  }, [t]);

  // 首次渲染以及外部 value 变化时同步文本和可视化配置。
  // 在渲染期间调整状态，而不是在 effect 里同步 setState
  const [syncedValue, setSyncedValue] = useState<string | null>(null);
  const incomingValue = value || '';
  if (syncedValue !== incomingValue) {
    setSyncedValue(incomingValue);
    if (!value) {
      setVisualConfig(buildDefaultConfig());
    } else {
      setTextContent(value);
      parseTextToVisual(value);
    }
  }

  // 只把用户改过的键合并进现有文本，文件里的其它内容（未知参数、自定义 section、注释）原样保留
  const syncVisualToTextWithConfig = useCallback((config: Record<string, string | number | boolean>) => {
    try {
      const updates: Record<string, string | number | boolean> = {};
      dirtyKeys.forEach((key) => {
        const value = config[key];
        if (value !== undefined) updates[key] = value;
      });

      const base = textContent.trim() === '' ? buildDefaultGameIniText(config) : textContent;
      const iniContent = mergeIni(base, updates, {
        // 基础设置和部落设置属于 ShooterGameMode，其余沿用 shootergamemode
        sectionForKey: (key) =>
          key in gameIniParams.gameBasic || key in gameIniParams.tribeSettings
            ? '[/Script/ShooterGame.ShooterGameMode]'
            : '[/script/shootergame.shootergamemode]',
      });

      setTextContent(iniContent);
      onChange?.(iniContent);
    } catch (error) {
      console.error(t('syncVisualToTextError') + ':', error);
    }
  }, [dirtyKeys, onChange, t, textContent]);

  const syncVisualToText = useCallback(() => {
    syncVisualToTextWithConfig(visualConfig);
  }, [visualConfig, syncVisualToTextWithConfig]);

  // 可视化编辑在事件里直接同步到文本并通知父组件，不再依赖 effect 和计时器
  const handleVisualChange = (key: string, value: string | number | boolean) => {
    setDirtyKeys(prev => {
      if (prev.has(key)) return prev;
      const next = new Set(prev);
      next.add(key);
      return next;
    });
    const next = { ...visualConfig, [key]: value };
    setVisualConfig(next);
    syncVisualToTextWithConfig(next);
  };

  const handleTextChange = (newText: string) => {
    setTextContent(newText);
    onChange?.(newText);

    // Parse text to visual with debounce
    setTimeout(() => {
      parseTextToVisual(newText);
    }, 500);
  };

  const handleModeSwitch = (mode: 'visual' | 'text') => {
    if (mode === 'text' && editMode === 'visual') {
      syncVisualToText();
    } else if (mode === 'visual' && editMode === 'text') {
      parseTextToVisual(textContent);
    }
    setEditMode(mode);
  };



  const getCategoryDisplayName = (categoryKey: GameIniCategoryKey): string => {
    const translated = tCategories(categoryKey as string);
    return translated !== categoryKey ? translated : categoryKey;
  };

  const renderParamControl = (paramKey: string, param: GameIniParam) => {
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
              step={param.step}
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

      default:
        return null;
    }
  };

  return (
    <div className="space-y-4">
      {/* Mode Switch */}
      <div className="flex items-center justify-between">
        <div className="flex items-center space-x-2">
          <Button
            variant={editMode === 'visual' ? 'default' : 'outline'}
            size="sm"
            onClick={() => handleModeSwitch('visual')}
          >
            <Eye className="w-4 h-4 mr-1" />
            {t('visualEditMode')}
          </Button>
          <Button
            variant={editMode === 'text' ? 'default' : 'outline'}
            size="sm"
            onClick={() => handleModeSwitch('text')}
          >
            <EyeOff className="w-4 h-4 mr-1" />
            {t('textEditMode')}
          </Button>
        </div>


      </div>

      {editMode === 'visual' ? (
        <Tabs value={activeTab} onValueChange={(value) => setActiveTab(value as GameIniCategoryKey)}>
          <TabsList className="grid w-full grid-cols-3 sm:grid-cols-4 md:grid-cols-6 lg:grid-cols-9 gap-2 h-auto p-2 bg-muted/50 border rounded-lg">
            {Object.keys(gameIniParams).map((categoryKey) => (
              <TabsTrigger 
                key={categoryKey} 
                value={categoryKey} 
                className="text-xs px-3 py-2 flex-shrink-0 border-2 border-transparent data-[state=active]:border-primary data-[state=active]:bg-primary data-[state=active]:text-primary-foreground data-[state=active]:shadow-md hover:border-muted-foreground/50 hover:bg-muted transition-all duration-200 font-medium rounded-md"
              >
                {getCategoryDisplayName(categoryKey as GameIniCategoryKey)}
              </TabsTrigger>
            ))}
          </TabsList>

          {Object.keys(gameIniParams).map((categoryKey) => {
            const category = categoryKey as GameIniCategoryKey;
            const params = getGameIniParamsByCategory(category);

            return (
              <TabsContent key={categoryKey} value={categoryKey} className="space-y-4">
                <Card>
                  <CardHeader>
                    <CardTitle className="text-lg">
                      {getCategoryDisplayName(category)}
                    </CardTitle>
                    <CardDescription>
                      {Object.keys(params).length}{t('parametersCount')}
                    </CardDescription>
                  </CardHeader>
                  <CardContent className="space-y-4">
                    <div className="grid gap-4">
                      {Object.keys(params).map((paramKey) => {
                        const param = params[paramKey];
                        return (
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
                        );
                      })}
                    </div>
                  </CardContent>
                </Card>
              </TabsContent>
            );
          })}
        </Tabs>
      ) : (
        <div className="space-y-2">
          <Label htmlFor="game-ini-text">{t('textEditMode')}</Label>
          <Textarea
            id="game-ini-text"
            value={textContent}
            onChange={(e) => handleTextChange(e.target.value)}
            className="min-h-[400px] font-mono text-sm"
            placeholder={t('placeholders.gameIni')}
          />
        </div>
      )}
    </div>
  );
}