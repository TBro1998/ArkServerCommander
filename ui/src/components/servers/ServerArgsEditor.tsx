"use client";

import { useTranslations } from 'next-intl';
import { ServerParam, getServerParamsByCategory } from '@/lib/ark-settings';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Switch } from '@/components/ui/switch';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Button } from '@/components/ui/button';
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
  const tParams = useTranslations('servers'); // 添加这行
  const paramCategories = getServerParamsByCategory();

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

    switch (param.type) {
      case 'boolean':
        return (
          <div className="flex items-center space-x-2">
            <Switch
              id={id}
              checked={type === 'query' ? currentValue === 'True' : Boolean(currentValue)}
              onCheckedChange={(checked: boolean) => handleParamChange(type === 'query' ? 'query_params' : 'command_line_args', key, type === 'query' ? (checked ? 'True' : 'False') : checked)}
            />
            <Label htmlFor={id}>{tParams(`queryParams.${key}`) || tParams(`commandLineArgs.${key}`)}</Label>
          </div>
        );
      case 'number':
        return (
          <div>
            <Label htmlFor={id}>{tParams(`queryParams.${key}`) || tParams(`commandLineArgs.${key}`)}</Label>
            <Input
              id={id}
              type="number"
              value={String(currentValue || '')}
              onChange={(e) => handleParamChange(type === 'query' ? 'query_params' : 'command_line_args', key, e.target.value)}
              min={param.min}
              max={param.max}
              step={param.step}
            />
          </div>
        );
      case 'string':
        return (
          <div>
            <Label htmlFor={id}>{tParams(`queryParams.${key}`) || tParams(`commandLineArgs.${key}`)}</Label>
            <Input
              id={id}
              type="text"
              value={String(currentValue || '')}
              onChange={(e) => handleParamChange(type === 'query' ? 'query_params' : 'command_line_args', key, e.target.value)}
            />
          </div>
        );
      case 'select':
        return (
          <div>
            <Label htmlFor={id}>{tParams(`queryParams.${key}`) || tParams(`commandLineArgs.${key}`)}</Label>
            <Select value={String(currentValue || '')} onValueChange={(val: string) => handleParamChange(type === 'query' ? 'query_params' : 'command_line_args', key, val)}>
              <SelectTrigger>
                <SelectValue placeholder={t('pleaseSelect')} />
              </SelectTrigger>
              <SelectContent>
                {param.options?.map(opt => <SelectItem key={opt} value={opt}>{opt}</SelectItem>)}
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
      {Object.entries(paramCategories).map(([category, params]) => (
        <div key={category}>
          <h3 className="text-lg font-semibold mb-2">{t(`paramCategories.${category}`)}</h3>
          <div className="space-y-4">
            {params.map(({ key, param }) => {
              const type = key in value.query_params ? 'query' : 'cmd';
              return renderParam(type, key, param)
            })}
          </div>
        </div>
      ))}
      <div>
        <h3 className="text-lg font-semibold mb-2">{t('customArgs')}</h3>
        <div className="space-y-2">
          {value.custom_args.map((arg, index) => (
            <div key={index} className="flex items-center gap-2">
              <Input value={arg} onChange={(e) => handleCustomArgChange(index, e.target.value)} />
              <Button variant="ghost" size="icon" onClick={() => removeCustomArg(index)}>
                <Trash2 className="h-4 w-4" />
              </Button>
            </div>
          ))}
          <Button onClick={addCustomArg}>{t('addCustomArg')}</Button>
        </div>
      </div>
    </div>
  );
}