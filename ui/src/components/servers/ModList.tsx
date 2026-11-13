"use client";

import { useEffect, useState } from 'react';
import { useLocale, useTranslations } from 'next-intl';
import axios from 'axios';
import { ExternalLink, AlertTriangle, ChevronUp, ChevronDown, GripVertical, X } from 'lucide-react';
import { formatBytes } from '@/lib/format';

export interface ModInfo {
  id: string;
  title: string;
  // false 表示没能从 Steam 拿到结果（例如服务器无法访问 Steam），此时 found 无意义
  checked: boolean;
  found: boolean;
  file_size?: number;
  time_updated?: number;
  preview_url?: string;
}

const MOD_ID_LIST = /^\d+(,\d+)*$/;
// 创意工坊链接里的 Mod ID，例如 https://steamcommunity.com/sharedfiles/filedetails/?id=731604991
const WORKSHOP_URL_ID = /[?&]id=(\d+)/;

export const WORKSHOP_BROWSE_URL = 'https://steamcommunity.com/app/346110/workshop/';

// 规范化 Mod ID：去掉空格和多余逗号，支持中文逗号；创意工坊链接换成其中的 ID；重复的只保留第一个
export function normalizeModIds(value: string): string {
  const ids = value
    .split(/[,，\s]+/)
    .filter(Boolean)
    .map((part) => part.match(WORKSHOP_URL_ID)?.[1] ?? part);
  return Array.from(new Set(ids)).join(',');
}

export function isValidModIds(value: string): boolean {
  const v = normalizeModIds(value);
  return v === '' || MOD_ID_LIST.test(v);
}

// Unix 秒转为本地日期
export function formatUnixDate(seconds: number | undefined, locale: string): string {
  if (!seconds) return '-';
  return new Date(seconds * 1000).toLocaleDateString(locale === 'zh' ? 'zh-CN' : 'en-US');
}

export function ModThumb({ url }: { url?: string }) {
  if (!url) return <span className="size-8 flex-shrink-0 rounded-sm border bg-muted" />;
  // 创意工坊预览图来自 Steam CDN，不经过 next/image 的域名白名单
  // eslint-disable-next-line @next/next/no-img-element
  return <img src={url} alt="" loading="lazy" className="size-8 flex-shrink-0 rounded-sm border object-cover" />;
}

// 查询创意工坊名称并列出 Mod；ids 变化后防抖 600ms 再请求。
// 传入 onChange 时可以拖动或用上下按钮调整加载顺序、移除 Mod
export function ModList({ ids, compact = false, onChange }: { ids: string; compact?: boolean; onChange?: (ids: string) => void }) {
  const t = useTranslations('mods');
  const locale = useLocale();
  const normalized = normalizeModIds(ids);
  const valid = MOD_ID_LIST.test(normalized);
  const [result, setResult] = useState<{ key: string; mods: ModInfo[] } | null>(null);
  const [dragIndex, setDragIndex] = useState<number | null>(null);

  useEffect(() => {
    if (!normalized || !valid) return;
    let cancelled = false;
    const timer = setTimeout(async () => {
      try {
        const res = await axios.get('/api/mods', { params: { ids: normalized } });
        if (!cancelled) setResult({ key: normalized, mods: res.data.data ?? [] });
      } catch {
        if (!cancelled) setResult({ key: normalized, mods: normalized.split(',').map((id) => ({ id, title: '', checked: false, found: false })) });
      }
    }, 600);
    return () => {
      cancelled = true;
      clearTimeout(timer);
    };
  }, [normalized, valid]);

  if (!normalized) return null;
  if (!valid) {
    return (
      <p className="flex items-center gap-1 text-xs text-red-600">
        <AlertTriangle className="h-3.5 w-3.5" />
        {t('invalid')}
      </p>
    );
  }

  const order = normalized.split(',');
  // 调整顺序时 ID 集合不变，已查到的信息按 ID 取，不用等重新查询
  const known = new Map((result?.mods ?? []).map((m) => [m.id, m]));
  const mods: ModInfo[] = order.map((id) => known.get(id) ?? { id, title: '', checked: false, found: false });
  const loading = (id: string) => !known.has(id);

  const move = (from: number, to: number) => {
    if (!onChange || from === to || to < 0 || to >= order.length) return;
    const next = [...order];
    const [item] = next.splice(from, 1);
    next.splice(to, 0, item);
    onChange(next.join(','));
  };
  const remove = (index: number) => onChange?.(order.filter((_, i) => i !== index).join(','));

  return (
    <div className="space-y-1">
      {onChange && order.length > 1 && <p className="text-xs text-muted-foreground">{t('orderHint')}</p>}
      <ul className={compact ? 'space-y-1' : 'space-y-1.5'}>
        {mods.map((m, i) => (
          <li
            key={m.id}
            draggable={!!onChange}
            onDragStart={() => setDragIndex(i)}
            onDragOver={(e) => {
              if (dragIndex === null) return;
              e.preventDefault();
            }}
            onDrop={(e) => {
              e.preventDefault();
              if (dragIndex !== null) move(dragIndex, i);
              setDragIndex(null);
            }}
            onDragEnd={() => setDragIndex(null)}
            className={`flex items-center justify-between gap-2 rounded-sm border bg-muted/50 px-2 py-1 text-xs ${
              dragIndex === i ? 'opacity-50' : ''
            } ${onChange ? 'cursor-grab' : ''}`}
          >
            <span className="flex min-w-0 items-center gap-2">
              {onChange && <GripVertical className="h-3.5 w-3.5 flex-shrink-0 text-gray-400" />}
              <span className="num text-gray-400">{i + 1}.</span>
              {!compact && <ModThumb url={m.preview_url} />}
              <span className="min-w-0">
                <span
                  className={`block truncate ${
                    loading(m.id) || !m.checked ? 'text-gray-500' : m.found ? 'text-gray-800' : 'text-red-600'
                  }`}
                >
                  {loading(m.id) ? t('loading') : !m.checked ? t('unavailable') : m.found ? m.title || m.id : t('notFound')}
                </span>
                {!compact && m.found && m.time_updated ? (
                  <span className="block text-[11px] text-gray-500">{t('updatedAt', { date: formatUnixDate(m.time_updated, locale) })}</span>
                ) : null}
              </span>
            </span>
            <span className="flex flex-shrink-0 items-center gap-2 text-gray-500">
              {!compact && m.file_size ? <span>{formatBytes(m.file_size)}</span> : null}
              <a
                href={`https://steamcommunity.com/sharedfiles/filedetails/?id=${m.id}`}
                target="_blank"
                rel="noreferrer"
                className="inline-flex items-center gap-0.5 font-mono hover:text-blue-600"
              >
                {m.id}
                <ExternalLink className="h-3 w-3" />
              </a>
              {onChange && (
                <span className="flex items-center">
                  <button
                    type="button"
                    className="rounded-sm p-0.5 hover:text-foreground disabled:opacity-30"
                    disabled={i === 0}
                    onClick={() => move(i, i - 1)}
                    aria-label={t('moveUp')}
                    title={t('moveUp')}
                  >
                    <ChevronUp className="h-3.5 w-3.5" />
                  </button>
                  <button
                    type="button"
                    className="rounded-sm p-0.5 hover:text-foreground disabled:opacity-30"
                    disabled={i === order.length - 1}
                    onClick={() => move(i, i + 1)}
                    aria-label={t('moveDown')}
                    title={t('moveDown')}
                  >
                    <ChevronDown className="h-3.5 w-3.5" />
                  </button>
                  <button
                    type="button"
                    className="rounded-sm p-0.5 hover:text-red-600"
                    onClick={() => remove(i)}
                    aria-label={t('remove')}
                    title={t('remove')}
                  >
                    <X className="h-3.5 w-3.5" />
                  </button>
                </span>
              )}
            </span>
          </li>
        ))}
      </ul>
    </div>
  );
}
