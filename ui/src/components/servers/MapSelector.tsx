"use client";

import { useState, useEffect } from 'react';
import { useTranslations } from 'next-intl';
import { Label } from '@/components/ui/label';
import { Button } from '@/components/ui/button';
import { Popover, PopoverContent, PopoverTrigger } from '@/components/ui/popover';
import { Command, CommandEmpty, CommandGroup, CommandInput, CommandItem, CommandList } from '@/components/ui/command';
import { ChevronDown, X, Check, Map } from 'lucide-react';
import { cn } from '@/lib/utils';

interface MapSelectorProps {
  value: string;
  onChange: (value: string) => void;
  label?: string;
}

export function MapSelector({ value, onChange, label }: MapSelectorProps) {
  const t = useTranslations('servers.edit');
  const [open, setOpen] = useState(false);
  const [searchValue, setSearchValue] = useState('');

  // 预定义的地图选项
  const predefinedMaps = {
    'TheIsland': t('maps.TheIsland'),
    'TheCenter': t('maps.TheCenter'),
    'ScorchedEarth_P': t('maps.ScorchedEarth_P'),
    'Aberration_P': t('maps.Aberration_P'),
    'Extinction': t('maps.Extinction'),
    'Valguero_P': t('maps.Valguero_P'),
    'Genesis': t('maps.Genesis'),
    'Genesis2': t('maps.Genesis2'),
    'CrystalIsles': t('maps.CrystalIsles'),
    'LostIsland': t('maps.LostIsland'),
    'Fjordur': t('maps.Fjordur')
  };

  // 同步外部value变化
  useEffect(() => {
    // 当外部value变化时，重置搜索值
    if (!open) {
      setSearchValue('');
    }
  }, [value, open]);

  // 过滤地图选项
  const filteredMaps = Object.entries(predefinedMaps).filter(([key, displayName]) => {
    if (!searchValue) return true;
    const searchLower = searchValue.toLowerCase();
    return key.toLowerCase().includes(searchLower) ||
      displayName.toLowerCase().includes(searchLower);
  });

  // 处理地图选择
  const handleMapSelect = (mapKey: string) => {
    onChange(mapKey);
    setOpen(false);
    setSearchValue('');
  };

  // 清空选择
  const handleClear = () => {
    onChange('');
    setSearchValue('');
  };

  // 获取当前选中地图的显示名称
  const getDisplayName = (mapKey: string) => {
    return predefinedMaps[mapKey as keyof typeof predefinedMaps] || mapKey;
  };

  return (
    <div className="space-y-2">
      {label && <Label>{label}</Label>}

      <Popover open={open} onOpenChange={setOpen}>
        <PopoverTrigger asChild>
          <Button
            variant="outline"
            role="combobox"
            aria-expanded={open}
            className="w-full justify-between"
          >
            <div className="flex items-center gap-2">
              <Map className="h-4 w-4 opacity-50" />
              <span className={cn(
                "truncate",
                !value && "text-muted-foreground"
              )}>
                {value ? getDisplayName(value) : t('selectMapPlaceholder')}
              </span>
            </div>
            <div className="flex items-center gap-1">
              {value && (
                <Button
                  variant="ghost"
                  size="sm"
                  className="h-4 w-4 p-0 hover:bg-transparent"
                  onClick={(e) => {
                    e.stopPropagation();
                    handleClear();
                  }}
                >
                  <X className="h-3 w-3" />
                </Button>
              )}
              <ChevronDown className={cn(
                "h-4 w-4 shrink-0 opacity-50 transition-transform duration-200",
                open && "rotate-180"
              )} />
            </div>
          </Button>
        </PopoverTrigger>

        <PopoverContent className="w-[--radix-popover-trigger-width] p-0" align="start">
          <Command>
            <CommandInput
              placeholder={t('searchMaps') || 'Search maps...'}
              value={searchValue}
              onValueChange={setSearchValue}
            />
            <CommandList>
              <CommandEmpty>{t('noMatchingMaps') || 'No maps found.'}</CommandEmpty>

              <CommandGroup heading={t('officialMaps') || 'Official Maps'}>
                {filteredMaps.map(([key, displayName]) => (
                  <CommandItem
                    key={key}
                    value={key}
                    onSelect={() => handleMapSelect(key)}
                    className="flex items-center justify-between"
                  >
                    <div className="flex flex-col">
                      <span className="font-medium">{displayName}</span>
                      {displayName !== key && (
                        <span className="text-xs text-muted-foreground">{key}</span>
                      )}
                    </div>
                    {value === key && (
                      <Check className="h-4 w-4" />
                    )}
                  </CommandItem>
                ))}
              </CommandGroup>

              {/* 自定义地图选项 */}
              {searchValue && !Object.keys(predefinedMaps).includes(searchValue) && (
                <CommandGroup heading={t('customMap') || 'Custom Map'}>
                  <CommandItem
                    value={searchValue}
                    onSelect={() => handleMapSelect(searchValue)}
                    className="flex items-center justify-between"
                  >
                    <div className="flex flex-col">
                      <span className="font-medium">{searchValue}</span>
                      <span className="text-xs text-muted-foreground">
                        {t('customMapPlaceholder') || 'Custom map'}
                      </span>
                    </div>
                  </CommandItem>
                </CommandGroup>
              )}
            </CommandList>
          </Command>
        </PopoverContent>
      </Popover>

      {/* 显示当前选中地图的友好名称 */}
      {value && Object.keys(predefinedMaps).includes(value) && value !== getDisplayName(value) && (
        <div className="text-xs text-muted-foreground">
          {t('mapId')}: {value}
        </div>
      )}
    </div>
  );
}