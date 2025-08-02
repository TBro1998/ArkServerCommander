"use client";

import { useState, useRef, useEffect } from 'react';
import { useTranslations } from 'next-intl';
import { Label } from '@/components/ui/label';
import { Input } from '@/components/ui/input';
import { Button } from '@/components/ui/button';
import { ChevronDown, X, Check } from 'lucide-react';
import { cn } from '@/lib/utils';

interface MapSelectorProps {
  value: string;
  onChange: (value: string) => void;
  label?: string;
}

export function MapSelector({ value, onChange, label }: MapSelectorProps) {
  const t = useTranslations('servers.edit');
  const [isOpen, setIsOpen] = useState(false);
  const [searchTerm, setSearchTerm] = useState('');
  const [inputValue, setInputValue] = useState(value || '');
  const inputRef = useRef<HTMLInputElement>(null);
  const dropdownRef = useRef<HTMLDivElement>(null);
  
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

  // 同步外部value变化到内部inputValue
  useEffect(() => {
    setInputValue(value || '');
  }, [value]);

  // 过滤地图选项
  const filteredMaps = Object.entries(predefinedMaps).filter(([key, displayName]) => {
    const searchLower = searchTerm.toLowerCase();
    return key.toLowerCase().includes(searchLower) || 
           displayName.toLowerCase().includes(searchLower);
  });

  // 点击外部关闭下拉框
  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (dropdownRef.current && !dropdownRef.current.contains(event.target as Node)) {
        setIsOpen(false);
        setSearchTerm('');
      }
    };

    document.addEventListener('mousedown', handleClickOutside);
    return () => document.removeEventListener('mousedown', handleClickOutside);
  }, []);

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const newValue = e.target.value;
    setInputValue(newValue);
    setSearchTerm(newValue);
    onChange(newValue);
    
    // 如果用户在输入，自动打开下拉框
    if (!isOpen && newValue) {
      setIsOpen(true);
    }
  };

  const handleMapSelect = (mapKey: string) => {
    setInputValue(mapKey);
    onChange(mapKey);
    setIsOpen(false);
    setSearchTerm('');
  };

  const handleInputFocus = () => {
    setSearchTerm(inputValue);
    setIsOpen(true);
  };

  const handleClearInput = () => {
    setInputValue('');
    onChange('');
    setSearchTerm('');
    setIsOpen(true); // 清空后打开下拉框显示所有选项
    inputRef.current?.focus();
  };

  const toggleDropdown = () => {
    setIsOpen(!isOpen);
    if (!isOpen) {
      setSearchTerm(inputValue);
      inputRef.current?.focus();
    } else {
      setSearchTerm('');
    }
  };

  // 获取当前选中地图的显示名称
  const getDisplayName = (mapKey: string) => {
    return predefinedMaps[mapKey as keyof typeof predefinedMaps] || mapKey;
  };

  return (
    <div className="relative" ref={dropdownRef}>
      {label && <Label htmlFor="map">{label}</Label>}
      <div className="relative">
        <Input
          ref={inputRef}
          id="map"
          name="map"
          value={inputValue}
          onChange={handleInputChange}
          onFocus={handleInputFocus}
          placeholder={t('selectMapPlaceholder')}
          className="pr-20"
          autoComplete="off"
        />
        
        {/* 右侧按钮组 */}
        <div className="absolute right-1 top-1/2 -translate-y-1/2 flex items-center gap-1">
          {inputValue && (
            <Button
              type="button"
              variant="ghost"
              size="sm"
              onClick={handleClearInput}
              className="h-6 w-6 p-0 hover:bg-gray-100"
            >
              <X className="h-3 w-3" />
            </Button>
          )}
          <Button
            type="button"
            variant="ghost"
            size="sm"
            onClick={toggleDropdown}
            className="h-6 w-6 p-0 hover:bg-gray-100"
          >
            <ChevronDown className={cn(
              "h-3 w-3 transition-transform duration-200",
              isOpen && "rotate-180"
            )} />
          </Button>
        </div>

        {/* 下拉选项 */}
        {isOpen && (
          <div className="absolute top-full left-0 right-0 z-50 mt-1 max-h-60 overflow-auto rounded-md border bg-white shadow-lg">
            {filteredMaps.length > 0 ? (
              <>
                {filteredMaps.map(([key, displayName]) => (
                  <div
                    key={key}
                    className={cn(
                      "flex items-center justify-between px-3 py-2 text-sm cursor-pointer hover:bg-gray-50",
                      inputValue === key && "bg-blue-50 text-blue-600"
                    )}
                    onClick={() => handleMapSelect(key)}
                  >
                    <div className="flex flex-col">
                      <span className="font-medium">{displayName}</span>
                      {displayName !== key && (
                        <span className="text-xs text-gray-500">{key}</span>
                      )}
                    </div>
                    {inputValue === key && (
                      <Check className="h-4 w-4 text-blue-600" />
                    )}
                  </div>
                ))}
                
                {/* 如果输入的不是预定义地图，显示自定义选项 */}
                {searchTerm && !Object.keys(predefinedMaps).includes(searchTerm) && (
                  <div className="border-t">
                    <div
                      className="flex items-center justify-between px-3 py-2 text-sm cursor-pointer hover:bg-gray-50 text-blue-600"
                      onClick={() => handleMapSelect(searchTerm)}
                    >
                      <div className="flex flex-col">
                        <span className="font-medium">{t('customMap')}: {searchTerm}</span>
                        <span className="text-xs text-gray-500">{t('customMapPlaceholder')}</span>
                      </div>
                    </div>
                  </div>
                )}
              </>
            ) : (
              <div className="px-3 py-2 text-sm text-gray-500">
                {searchTerm ? t('noMatchingMaps') : t('selectMapPlaceholder')}
              </div>
            )}
          </div>
        )}
      </div>
      
      {/* 显示当前选中地图的友好名称 */}
      {inputValue && Object.keys(predefinedMaps).includes(inputValue) && (
        <div className="mt-1 text-xs text-gray-600">
          {getDisplayName(inputValue)}
        </div>
      )}
    </div>
  );
}