"use client";

import { useState, useEffect, useCallback } from 'react';
import { useTranslations } from 'next-intl';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Switch } from '@/components/ui/switch';
import { Textarea } from '@/components/ui/textarea';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Info, RotateCcw, Eye, EyeOff } from 'lucide-react';
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
  description?: string;
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
    bUseSingleplayerSettings: { type: 'boolean', default: false, description: '使用单人游戏设置' },
    bDisableStructurePlacementCollision: { type: 'boolean', default: false, description: '禁用建筑放置碰撞检测' },
    bAllowFlyerCarryPvE: { type: 'boolean', default: true, description: 'PvE模式下允许飞行载具抓取' },
    bDisableStructureDecayPvE: { type: 'boolean', default: false, description: 'PvE模式下禁用建筑衰减' },
    bAllowUnlimitedRespecs: { type: 'boolean', default: true, description: '允许无限重置技能点' },
    bAllowPlatformSaddleMultiFloors: { type: 'boolean', default: true, description: '允许平台鞍座多层建筑' },
    bPassiveDefensesDamageRiderlessDinos: { type: 'boolean', default: true, description: '被动防御伤害无骑手恐龙' },
    bPvEDisableFriendlyFire: { type: 'boolean', default: false, description: 'PvE模式禁用友军伤害' },
    bDisableFriendlyFire: { type: 'boolean', default: false, description: '禁用友军伤害' },
    bEnablePvPGamma: { type: 'boolean', default: false, description: '启用PvP伽马调节' },
    DifficultyOffset: { type: 'number', default: 1.0, min: 0.0, max: 1.0, step: 0.1, description: '难度偏移值' },
    OverrideOfficialDifficulty: { type: 'number', default: 5.0, min: 1.0, max: 10.0, step: 0.1, description: '覆盖官方难度' }
  },
  experienceSettings: {
    XPMultiplier: { type: 'number', default: 1.0, min: 0.1, max: 10.0, step: 0.1, description: '经验值倍率' },
    PlayerCharacterWaterDrainMultiplier: { type: 'number', default: 1.0, min: 0.1, max: 10.0, step: 0.1, description: '玩家水分消耗倍率' },
    PlayerCharacterFoodDrainMultiplier: { type: 'number', default: 1.0, min: 0.1, max: 10.0, step: 0.1, description: '玩家食物消耗倍率' },
    PlayerCharacterStaminaDrainMultiplier: { type: 'number', default: 1.0, min: 0.1, max: 10.0, step: 0.1, description: '玩家耐力消耗倍率' },
    PlayerCharacterHealthRecoveryMultiplier: { type: 'number', default: 1.0, min: 0.1, max: 10.0, step: 0.1, description: '玩家生命恢复倍率' }
  },
  breedingSettings: {
    MatingIntervalMultiplier: { type: 'number', default: 1.0, min: 0.1, max: 10.0, step: 0.1, description: '交配间隔倍率' },
    EggHatchSpeedMultiplier: { type: 'number', default: 1.0, min: 0.1, max: 50.0, step: 0.1, description: '蛋孵化速度倍率' },
    BabyMatureSpeedMultiplier: { type: 'number', default: 1.0, min: 0.1, max: 50.0, step: 0.1, description: '幼体成长速度倍率' },
    BabyFoodConsumptionSpeedMultiplier: { type: 'number', default: 1.0, min: 0.1, max: 10.0, step: 0.1, description: '幼体食物消耗速度倍率' },
    BabyCuddleIntervalMultiplier: { type: 'number', default: 1.0, min: 0.1, max: 10.0, step: 0.1, description: '幼体印记间隔倍率' },
    BabyCuddleGracePeriodMultiplier: { type: 'number', default: 1.0, min: 0.1, max: 10.0, step: 0.1, description: '幼体印记宽限期倍率' },
    BabyCuddleLoseImprintQualitySpeedMultiplier: { type: 'number', default: 1.0, min: 0.1, max: 10.0, step: 0.1, description: '幼体印记质量丢失速度倍率' }
  },
  itemSettings: {
    HarvestAmountMultiplier: { type: 'number', default: 1.0, min: 0.1, max: 10.0, step: 0.1, description: '收获量倍率' },
    HarvestHealthMultiplier: { type: 'number', default: 1.0, min: 0.1, max: 10.0, step: 0.1, description: '收获生命值倍率' },
    ResourcesRespawnPeriodMultiplier: { type: 'number', default: 1.0, min: 0.1, max: 10.0, step: 0.1, description: '资源重生周期倍率' },
    ItemStackSizeMultiplier: { type: 'number', default: 1.0, min: 0.1, max: 10.0, step: 0.1, description: '物品堆叠大小倍率' },
    GlobalSpoilingTimeMultiplier: { type: 'number', default: 1.0, min: 0.1, max: 10.0, step: 0.1, description: '全局腐败时间倍率' },
    GlobalItemDecompositionTimeMultiplier: { type: 'number', default: 1.0, min: 0.1, max: 10.0, step: 0.1, description: '全局物品分解时间倍率' },
    GlobalCorpseDecompositionTimeMultiplier: { type: 'number', default: 1.0, min: 0.1, max: 10.0, step: 0.1, description: '全局尸体分解时间倍率' }
  },
  dinoSettings: {
    TamingSpeedMultiplier: { type: 'number', default: 1.0, min: 0.1, max: 50.0, step: 0.1, description: '驯服速度倍率' },
    DinoCharacterFoodDrainMultiplier: { type: 'number', default: 1.0, min: 0.1, max: 10.0, step: 0.1, description: '恐龙食物消耗倍率' },
    DinoCharacterStaminaDrainMultiplier: { type: 'number', default: 1.0, min: 0.1, max: 10.0, step: 0.1, description: '恐龙耐力消耗倍率' },
    DinoCharacterHealthRecoveryMultiplier: { type: 'number', default: 1.0, min: 0.1, max: 10.0, step: 0.1, description: '恐龙生命恢复倍率' },
    DinoCountMultiplier: { type: 'number', default: 1.0, min: 0.1, max: 10.0, step: 0.1, description: '恐龙数量倍率' },
    WildDinoCharacterFoodDrainMultiplier: { type: 'number', default: 1.0, min: 0.1, max: 10.0, step: 0.1, description: '野生恐龙食物消耗倍率' },
    WildDinoTorporDrainMultiplier: { type: 'number', default: 1.0, min: 0.1, max: 10.0, step: 0.1, description: '野生恐龙眩晕值消耗倍率' }
  },
  tribeSettings: {
    MaxNumberOfPlayersInTribe: { type: 'number', default: 0, min: 0, max: 100, step: 1, description: '部落最大玩家数量（0为无限制）' },
    TribeNameChangeCooldown: { type: 'number', default: 15, min: 0, max: 1440, step: 1, description: '部落名称更改冷却时间（分钟）' },
    bPvEAllowTribeWar: { type: 'boolean', default: false, description: 'PvE模式允许部落战争' },
    bPvEAllowTribeWarCancel: { type: 'boolean', default: false, description: 'PvE模式允许取消部落战争' }
  },
  pvpSettings: {
    bIncreasePvPRespawnInterval: { type: 'boolean', default: false, description: '增加PvP重生间隔' },
    IncreasePvPRespawnIntervalCheckPeriod: { type: 'number', default: 300, min: 60, max: 3600, step: 60, description: 'PvP重生间隔检查周期（秒）' },
    IncreasePvPRespawnIntervalMultiplier: { type: 'number', default: 1.0, min: 1.0, max: 5.0, step: 0.1, description: 'PvP重生间隔倍率' },
    IncreasePvPRespawnIntervalBaseAmount: { type: 'number', default: 60, min: 30, max: 600, step: 30, description: 'PvP重生间隔基础时间（秒）' }
  },
  structureSettings: {
    StructureDamageMultiplier: { type: 'number', default: 1.0, min: 0.1, max: 10.0, step: 0.1, description: '建筑伤害倍率' },
    StructureResistanceMultiplier: { type: 'number', default: 1.0, min: 0.1, max: 10.0, step: 0.1, description: '建筑抗性倍率' },
    StructureDamageRepairCooldown: { type: 'number', default: 180, min: 0, max: 3600, step: 60, description: '建筑伤害修复冷却时间（秒）' },
    PvEStructureDecayPeriodMultiplier: { type: 'number', default: 1.0, min: 0.1, max: 10.0, step: 0.1, description: 'PvE建筑衰减周期倍率' },
    MaxStructuresInRange: { type: 'number', default: 1300, min: 100, max: 10000, step: 100, description: '范围内最大建筑数量' }
  },
  advancedSettings: {
    bAutoPvETimer: { type: 'boolean', default: false, description: '自动PvE计时器' },
    bAutoPvEUseSystemTime: { type: 'boolean', default: false, description: '自动PvE使用系统时间' },
    AutoPvEStartTimeSeconds: { type: 'number', default: 0, min: 0, max: 86400, step: 3600, description: '自动PvE开始时间（秒）' },
    AutoPvEStopTimeSeconds: { type: 'number', default: 0, min: 0, max: 86400, step: 3600, description: '自动PvE结束时间（秒）' },
    bOnlyAllowSpecifiedEngrams: { type: 'boolean', default: false, description: '仅允许指定图纸' },
    bAutoUnlockAllEngrams: { type: 'boolean', default: false, description: '自动解锁所有图纸' },
    bShowCreativeMode: { type: 'boolean', default: false, description: '显示创造模式' },
    bUseCorpseLocator: { type: 'boolean', default: false, description: '使用尸体定位器' },
    bDisableLootCrates: { type: 'boolean', default: false, description: '禁用补给箱' },
    bDisableDinoRiding: { type: 'boolean', default: false, description: '禁用骑乘恐龙' },
    bDisableDinoTaming: { type: 'boolean', default: false, description: '禁用驯服恐龙' },
    bAllowCustomRecipes: { type: 'boolean', default: true, description: '允许自定义配方' }
  },
  customSettings: {
    DayCycleSpeedScale: { type: 'number', default: 1.0, min: 0.1, max: 10.0, step: 0.1, description: '昼夜循环速度倍率' },
    NightTimeSpeedScale: { type: 'number', default: 1.0, min: 0.1, max: 10.0, step: 0.1, description: '夜晚时间流逝速度倍率' },
    DayTimeSpeedScale: { type: 'number', default: 1.0, min: 0.1, max: 10.0, step: 0.1, description: '白天时间流逝速度倍率' }
  }
};

function getGameIniParamsByCategory(category: GameIniCategoryKey): Record<string, GameIniParam> {
  return gameIniParams[category] || {};
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
  const [isSyncing, setIsSyncing] = useState(false);
  const [textContent, setTextContent] = useState('');
  const [visualConfig, setVisualConfig] = useState<Record<string, string | number | boolean>>({});
  const [activeTab, setActiveTab] = useState<GameIniCategoryKey>('gameBasic');
  const [showTooltip, setShowTooltip] = useState<string | null>(null);
  const [isUserEditing, setIsUserEditing] = useState(false);
  const [lastUserEditTime, setLastUserEditTime] = useState(0);

  // Initialize with default values
  useEffect(() => {
    if (!value) {
      const defaultConfig: Record<string, string | number | boolean> = {};
      Object.keys(gameIniParams).forEach(categoryKey => {
        const category = categoryKey as GameIniCategoryKey;
        const params = getGameIniParamsByCategory(category);
        Object.keys(params).forEach(paramKey => {
          defaultConfig[paramKey] = params[paramKey].default;
        });
      });
      setVisualConfig(defaultConfig);
    } else {
      setTextContent(value);
      parseTextToVisual(value);
    }
  }, [value]);

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
      console.error('解析Game.ini文本失败:', error);
    }
  }, []);

  const syncVisualToText = useCallback(() => {
    try {
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
      
      setTextContent(iniContent);
      onChange?.(iniContent);
    } catch (error) {
      console.error('同步可视化配置到文本失败:', error);
    }
  }, [visualConfig, onChange]);

  // Sync visual to text when visual config changes
  useEffect(() => {
    if (editMode === 'visual' && !isUserEditing) {
      const now = Date.now();
      if (now - lastUserEditTime > 500) {
        syncVisualToText();
      }
    }
  }, [visualConfig, editMode, syncVisualToText, isUserEditing, lastUserEditTime]);

  const handleVisualChange = (key: string, value: string | number | boolean) => {
    setIsUserEditing(true);
    setLastUserEditTime(Date.now());
    setVisualConfig(prev => ({ ...prev, [key]: value }));
    setTimeout(() => setIsUserEditing(false), 1000);
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

  const resetToDefault = () => {
    const defaultConfig: Record<string, string | number | boolean> = {};
    Object.keys(gameIniParams).forEach(categoryKey => {
      const category = categoryKey as GameIniCategoryKey;
      const params = getGameIniParamsByCategory(category);
      Object.keys(params).forEach(paramKey => {
        defaultConfig[paramKey] = params[paramKey].default;
      });
    });
    setVisualConfig(defaultConfig);
  };

  const getParamDisplayName = (paramKey: string): string => {
    const translated = tParams(paramKey as any);
    return translated !== paramKey ? translated : paramKey;
  };

  const getCategoryDisplayName = (categoryKey: GameIniCategoryKey): string => {
    const translationKey = `servers.gameIniCategories.${categoryKey}`;
    const translated = tCategories(translationKey as any);
    return translated !== translationKey ? translated : categoryKey;
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
        
        {editMode === 'visual' && (
          <Button variant="outline" size="sm" onClick={resetToDefault}>
            <RotateCcw className="w-4 h-4 mr-1" />
            {t('resetToDefault')}
          </Button>
        )}
      </div>

      {/* Mode Description */}
      <div className="text-sm text-muted-foreground p-3 bg-muted/50 rounded-md">
        {editMode === 'visual' ? t('visualEditModeTip') : t('gameIniTextEditDesc')}
      </div>

      {editMode === 'visual' ? (
        <Tabs value={activeTab} onValueChange={(value) => setActiveTab(value as GameIniCategoryKey)}>
          <TabsList className="grid w-full grid-cols-5 lg:grid-cols-10">
            {Object.keys(gameIniParams).map((categoryKey) => (
              <TabsTrigger key={categoryKey} value={categoryKey} className="text-xs">
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
                      {Object.keys(params).length} 个参数
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
                                      {param.description && (
                                        <p className="text-sm mt-1">{param.description}</p>
                                      )}
                                      <p className="text-xs mt-1 text-muted-foreground">
                                        默认值: {String(param.default)}
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