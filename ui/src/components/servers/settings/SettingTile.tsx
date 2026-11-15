"use client";

import { useState } from 'react';
import { useTranslations } from 'next-intl';
import { Eye, EyeOff, RotateCcw } from 'lucide-react';
import { Input } from '@/components/ui/input';
import { Switch } from '@/components/ui/switch';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { cn } from '@/lib/utils';
import type { SettingDef, SettingUnit, SettingValue } from '@/lib/server-settings';

interface SettingTileProps {
  label: string;
  htmlFor?: string;
  // 游戏里的键名，方便对照文档和搜索
  name?: string;
  hint?: string;
  error?: string;
  modified?: boolean;
  onReset?: () => void;
  // 控件和标题在同一行（开关）
  inline?: boolean;
  wide?: boolean;
  children: React.ReactNode;
}

export function SettingTile({ label, htmlFor, name, hint, error, modified, onReset, inline, wide, children }: SettingTileProps) {
  const t = useTranslations('serverSettings');
  return (
    <div
      className={cn(
        'rounded-lg border bg-background/40 px-3 py-2.5 transition-colors',
        modified && 'border-primary/40',
        error && 'border-red-500/70',
        wide && 'sm:col-span-2'
      )}
    >
      <div className={cn('flex justify-between gap-3', inline ? 'items-center' : 'items-start')}>
        <div className="min-w-0">
          <label htmlFor={htmlFor} className="flex items-center gap-1.5 text-sm font-medium">
            {label}
            {modified && <span className="size-1.5 shrink-0 rounded-full bg-primary" title={t('modified')} />}
          </label>
          {name && <div className="truncate font-mono text-[11px] text-gray-500">{name}</div>}
        </div>
        <div className="flex shrink-0 items-center gap-2">
          {modified && onReset && (
            <button
              type="button"
              onClick={onReset}
              className="rounded p-1 text-gray-500 hover:bg-gray-100 hover:text-gray-900"
              title={t('reset')}
              aria-label={t('reset')}
            >
              <RotateCcw className="size-3.5" />
            </button>
          )}
          {inline && children}
        </div>
      </div>
      {!inline && <div className="mt-2">{children}</div>}
      {error ? (
        <p className="mt-1.5 text-xs text-red-600">{error}</p>
      ) : hint ? (
        <p className="mt-1.5 text-xs text-gray-500">{hint}</p>
      ) : null}
    </div>
  );
}

// 数字输入框：输入过程中允许暂时不合法（如 "1."、"-"），合法时立即提交；清空表示恢复默认
export function NumberInput({
  id,
  value,
  placeholder,
  unit,
  onCommit,
  invalid,
}: {
  id?: string;
  value: number | undefined;
  placeholder?: string;
  unit?: SettingUnit;
  onCommit: (value: number | undefined) => void;
  invalid?: boolean;
}) {
  const t = useTranslations('serverSettings.units');
  const format = (v: number | undefined) => (v === undefined ? '' : String(v));
  const [draft, setDraft] = useState(format(value));
  const [prevValue, setPrevValue] = useState(value);
  // 外部值变化（恢复默认、在配置文件里修改）时同步；和草稿表示同一个数时保留草稿
  if (value !== prevValue) {
    setPrevValue(value);
    const parsed = draft.trim() === '' ? undefined : Number(draft);
    if (parsed !== value) setDraft(format(value));
  }

  return (
    <div className="relative">
      <Input
        id={id}
        inputMode="decimal"
        value={draft}
        placeholder={placeholder}
        aria-invalid={invalid || undefined}
        onChange={(e) => {
          const next = e.target.value;
          setDraft(next);
          if (next.trim() === '') onCommit(undefined);
          else if (Number.isFinite(Number(next))) onCommit(Number(next));
        }}
        onBlur={() => setDraft(format(value))}
        className={cn('num', unit && 'pr-12')}
      />
      {unit && (
        <span className="pointer-events-none absolute inset-y-0 right-3 flex items-center text-xs text-gray-500">
          {unit === 'x' ? '×' : t(unit)}
        </span>
      )}
    </div>
  );
}

export function PasswordInput({ id, value, onChange, placeholder, className }: { id?: string; value: string; onChange: (value: string) => void; placeholder?: string; className?: string }) {
  const t = useTranslations('serverSettings');
  const [show, setShow] = useState(false);
  return (
    <div className={cn('relative', className)}>
      <Input
        id={id}
        type={show ? 'text' : 'password'}
        value={value}
        onChange={(e) => onChange(e.target.value)}
        placeholder={placeholder}
        autoComplete="new-password"
        className="pr-10"
      />
      <button
        type="button"
        onClick={() => setShow(!show)}
        className="absolute inset-y-0 right-0 flex items-center pr-3"
        aria-label={show ? t('hidePassword') : t('showPassword')}
      >
        {show ? <EyeOff className="size-4 text-gray-500" /> : <Eye className="size-4 text-gray-500" />}
      </button>
    </div>
  );
}

const DEFAULT_OPTION = '__default__';

/** 按设置类型渲染控件。value 为 undefined 时显示默认值 */
export function SettingControl({
  def,
  id,
  value,
  onChange,
}: {
  def: SettingDef;
  id: string;
  value: SettingValue | undefined;
  onChange: (value: SettingValue | undefined) => void;
}) {
  const t = useTranslations('serverSettings');

  switch (def.type) {
    case 'boolean':
      return <Switch id={id} checked={Boolean(value ?? def.default)} onCheckedChange={(checked) => onChange(checked)} />;
    case 'number':
      return (
        <NumberInput
          id={id}
          value={typeof value === 'number' ? value : undefined}
          placeholder={def.default !== undefined ? String(def.default) : t('notSet')}
          unit={def.unit}
          onCommit={onChange}
        />
      );
    case 'password':
      return <PasswordInput id={id} value={String(value ?? '')} onChange={onChange} />;
    case 'select':
      return (
        <Select value={value ? String(value) : DEFAULT_OPTION} onValueChange={(v) => onChange(v === DEFAULT_OPTION ? undefined : v)}>
          <SelectTrigger id={id} className="w-full">
            <SelectValue />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value={DEFAULT_OPTION}>{t('defaultOption')}</SelectItem>
            {def.options?.map((opt) => (
              <SelectItem key={opt} value={opt}>
                {opt}
              </SelectItem>
            ))}
          </SelectContent>
        </Select>
      );
    default:
      return (
        <Input
          id={id}
          value={String(value ?? '')}
          placeholder={def.default ? String(def.default) : undefined}
          onChange={(e) => onChange(e.target.value)}
        />
      );
  }
}
