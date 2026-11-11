"use client";

import { useEffect, useState } from 'react';
import { useTranslations } from 'next-intl';
import axios from 'axios';
import { ExternalLink, AlertTriangle } from 'lucide-react';
import { formatBytes } from '@/lib/format';

export interface ModInfo {
  id: string;
  title: string;
  // false 表示没能从 Steam 拿到结果（例如服务器无法访问 Steam），此时 found 无意义
  checked: boolean;
  found: boolean;
  file_size?: number;
  time_updated?: number;
}

const MOD_ID_LIST = /^\d+(,\d+)*$/;

// 规范化 Mod ID：去掉空格和多余逗号，支持中文逗号
export function normalizeModIds(value: string): string {
  return value
    .split(/[,，\s]+/)
    .filter(Boolean)
    .join(',');
}

export function isValidModIds(value: string): boolean {
  const v = normalizeModIds(value);
  return v === '' || MOD_ID_LIST.test(v);
}

// 查询创意工坊名称并列出 Mod；ids 变化后防抖 600ms 再请求
export function ModList({ ids, compact = false }: { ids: string; compact?: boolean }) {
  const t = useTranslations('mods');
  const normalized = normalizeModIds(ids);
  const valid = MOD_ID_LIST.test(normalized);
  const [result, setResult] = useState<{ key: string; mods: ModInfo[] } | null>(null);

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

  // 结果对应的不是当前输入时，先按 ID 显示
  const mods: ModInfo[] =
    result?.key === normalized ? result.mods : normalized.split(',').map((id) => ({ id, title: '', checked: false, found: false }));
  const loading = result?.key !== normalized;

  return (
    <ul className={compact ? 'space-y-1' : 'space-y-1.5'}>
      {mods.map((m, i) => (
        <li key={`${m.id}-${i}`} className="flex items-center justify-between gap-2 rounded-sm border bg-muted/50 px-2 py-1 text-xs">
          <span className="flex min-w-0 items-center gap-2">
            <span className="text-gray-400">{i + 1}.</span>
            <span
              className={`truncate ${
                loading || !m.checked ? 'text-gray-500' : m.found ? 'text-gray-800' : 'text-red-600'
              }`}
            >
              {loading ? t('loading') : !m.checked ? t('unavailable') : m.found ? m.title || m.id : t('notFound')}
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
          </span>
        </li>
      ))}
    </ul>
  );
}
