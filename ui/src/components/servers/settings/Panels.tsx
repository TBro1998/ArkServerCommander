"use client";

import { useTranslations } from 'next-intl';
import { Plus, Trash2 } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Textarea } from '@/components/ui/textarea';

// 自定义启动参数：字典里没有的参数，原样追加到启动命令
export function CustomArgsPanel({ value, onChange }: { value: string[]; onChange: (value: string[]) => void }) {
  const t = useTranslations('serverSettings.panels');
  return (
    <div className="space-y-2 sm:col-span-2">
      {value.map((arg, index) => (
        <div key={index} className="flex items-center gap-2">
          <Input
            value={arg}
            onChange={(e) => onChange(value.map((v, i) => (i === index ? e.target.value : v)))}
            placeholder={t('customArgPlaceholder')}
            className="font-mono"
          />
          <Button type="button" variant="ghost" size="icon" onClick={() => onChange(value.filter((_, i) => i !== index))} aria-label={t('removeCustomArg')}>
            <Trash2 />
          </Button>
        </div>
      ))}
      <Button type="button" variant="outline" size="sm" onClick={() => onChange([...value, ''])}>
        <Plus />
        {t('addCustomArg')}
      </Button>
    </div>
  );
}

// 直接编辑配置文件：给需要字典里没有的键（如 ConfigOverride*、重复键）的用户。
// 和上面的设置项读写的是同一份内容，在这里改的值切回设置页就能看到
export function RawConfigPanel({
  gameUserSettings,
  gameIni,
  onChange,
}: {
  gameUserSettings: string;
  gameIni: string;
  onChange: (patch: { game_user_settings?: string; game_ini?: string }) => void;
}) {
  const t = useTranslations('serverSettings.panels');
  return (
    <div className="space-y-4 sm:col-span-2">
      <div className="space-y-1.5">
        <label htmlFor="raw-gus" className="font-mono text-xs text-gray-700">
          GameUserSettings.ini
        </label>
        <Textarea
          id="raw-gus"
          value={gameUserSettings}
          onChange={(e) => onChange({ game_user_settings: e.target.value })}
          className="min-h-64 font-mono text-xs"
          spellCheck={false}
          placeholder={t('rawEmpty')}
        />
      </div>
      <div className="space-y-1.5">
        <label htmlFor="raw-game" className="font-mono text-xs text-gray-700">
          Game.ini
        </label>
        <Textarea
          id="raw-game"
          value={gameIni}
          onChange={(e) => onChange({ game_ini: e.target.value })}
          className="min-h-48 font-mono text-xs"
          spellCheck={false}
          placeholder={t('rawEmpty')}
        />
      </div>
    </div>
  );
}
