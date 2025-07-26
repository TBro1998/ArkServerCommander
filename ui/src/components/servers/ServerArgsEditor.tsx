"use client";

import { useState } from 'react';
import { useTranslations } from 'next-intl';
import { ServerParam, getServerParamsByCategory, CategoryKey } from '@/lib/ark-settings';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Switch } from '@/components/ui/switch';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Button } from '@/components/ui/button';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { Card, CardContent } from '@/components/ui/card';
import { Trash2 } from 'lucide-react';

interface ServerArgsEditorProps {
  value: {
    query_params: Record<string, string>;
    command_line_args: Record<string, string | number | boolean>;
    custom_args: string[];
  };
  onChange: (value: ServerArgsEditorProps['value']) => void;
}

export function ServerArgsEditor({ value, onChange }: ServerArgsEditorProps) {
  const t = useTranslations('servers.argsEditor');
  const tParams = useTranslations('servers.paramCategories');
  const paramCategories = getServerParamsByCategory();
  const [activeTab, setActiveTab] = useState<CategoryKey>('basic');

  // 获取有参数的分类
  const availableCategories = Object.entries(paramCategories)
    .filter(([_, params]) => params.length > 0)
    .map(([category]) => category as CategoryKey);

  const handleParamChange = (
    type: 'query_params' | 'command_line_args' | 'custom_args',
    key: string,
    val: string | number | boolean
  ) => {
    onChange({
      ...value,
      [type]: {
        ...value[type],
        [key]: val,
      },
    });
  };

  const handleCustomArgChange = (index: number, val: string) => {
    const newCustomArgs = [...value.custom_args];
    newCustomArgs[index] = val;
    onChange({ ...value, custom_args: newCustomArgs });
  };

  const addCustomArg = () => {
    onChange({ ...value, custom_args: [...value.custom_args, ''] });
  };

  const removeCustomArg = (index: number) => {
    const newCustomArgs = value.custom_args.filter((_, i) => i !== index);
    onChange({ ...value, custom_args: newCustomArgs });
  };

  const renderParam = (type: 'query' | 'cmd', key: string, param: ServerParam) => {
    const id = `${type}-${key}`;
    const currentValue = type === 'query' ? value.query_params[key] : value.command_line_args[key];
    
    // 获取参数的翻译名称
    const tQueryParams = useTranslations('servers.queryParams');
    const tCommandLineArgs = useTranslations('servers.commandLineArgs');
    const paramDisplayName = type === 'query' 
      ? (tQueryParams.has(key) ? tQueryParams(key) : key)
      : (tCommandLineArgs.has(key) ? tCommandLineArgs(key) : key);

    switch (param.type) {
      case 'boolean':
        return (
          <div key={id} className="flex items-center space-x-2">
            <Switch
              id={id}
              checked={type === 'query' ? currentValue === 'True' : Boolean(currentValue)}
              onCheckedChange={(checked: boolean) => 
                handleParamChange(
                  type === 'query' ? 'query_params' : 'command_line_args', 
                  key, 
                  type === 'query' ? (checked ? 'True' : 'False') : checked
                )
              }
            />
            <Label htmlFor={id}>{paramDisplayName}</Label>
          </div>
        );
      case 'number':
        return (
          <div key={id}>
            <Label htmlFor={id}>{paramDisplayName}</Label>
            <Input
              id={id}
              type="number"
              value={String(currentValue || '')}
              onChange={(e) => 
                handleParamChange(
                  type === 'query' ? 'query_params' : 'command_line_args', 
                  key, 
                  e.target.value
                )
              }
              min={param.min}
              max={param.max}
              step={param.step}
            />
          </div>
        );
      case 'string':
        return (
          <div key={id}>
            <Label htmlFor={id}>{paramDisplayName}</Label>
            <Input
              id={id}
              type="text"
              value={String(currentValue || '')}
              onChange={(e) => 
                handleParamChange(
                  type === 'query' ? 'query_params' : 'command_line_args', 
                  key, 
                  e.target.value
                )
              }
            />
          </div>
        );
      case 'select':
        return (
          <div key={id}>
            <Label htmlFor={id}>{paramDisplayName}</Label>
            <Select 
              value={String(currentValue || '')} 
              onValueChange={(val: string) => 
                handleParamChange(
                  type === 'query' ? 'query_params' : 'command_line_args', 
                  key, 
                  val
                )
              }
            >
              <SelectTrigger>
                <SelectValue placeholder={t('pleaseSelect')} />
              </SelectTrigger>
              <SelectContent>
                {param.options?.map(opt => (
                  <SelectItem key={opt} value={opt}>{opt}</SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>
        );
      default:
        return null;
    }
  };

  return (
    <div className="space-y-6">
      <Tabs value={activeTab} onValueChange={(value) => setActiveTab(value as CategoryKey)}>
        <TabsList className="flex flex-wrap h-auto">
          {availableCategories.map((category) => (
            <TabsTrigger key={category} value={category} className="text-sm">
              {tParams(category)}
            </TabsTrigger>
          ))}
          <TabsTrigger value="custom" className="text-sm">
            {t('customArgs')}
          </TabsTrigger>
        </TabsList>

        {availableCategories.map((category) => (
          <TabsContent key={category} value={category} className="mt-4">
            <Card>
              <CardContent className="pt-6">
                <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
                  {paramCategories[category].map(({ key, param }) => {
                    const type = key in value.query_params ? 'query' : 'cmd';
                    return renderParam(type, key, param);
                  })}
                </div>
              </CardContent>
            </Card>
          </TabsContent>
        ))}

        <TabsContent value="custom" className="mt-4">
          <Card>
            <CardContent className="pt-6">
              <div className="space-y-4">
                <div className="space-y-2">
                  {value.custom_args.map((arg, index) => (
                    <div key={index} className="flex items-center gap-2">
                      <Input 
                        value={arg} 
                        onChange={(e) => handleCustomArgChange(index, e.target.value)}
                        placeholder={t('customArgPlaceholder')}
                      />
                      <Button variant="ghost" size="icon" onClick={() => removeCustomArg(index)}>
                        <Trash2 className="h-4 w-4" />
                      </Button>
                    </div>
                  ))}
                  <Button onClick={addCustomArg}>{t('addCustomArg')}</Button>
                </div>
              </div>
            </CardContent>
          </Card>
        </TabsContent>
      </Tabs>
    </div>
  );
}