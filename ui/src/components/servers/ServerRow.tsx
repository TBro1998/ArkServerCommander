"use client";

import { useTranslations } from 'next-intl';
import Link from 'next/link';
import { AlertTriangle, Copy } from 'lucide-react';
import { Progress } from '@/components/ui/progress';
import type { Server } from '@/stores/servers';
import { StatusBadge } from './StatusBadge';
import { ServerTags } from './ServerCard';
import { PowerButtons, ServerShortcuts, useServerControls } from './ServerActions';
import { copyText, formatBytes, formatPercent } from '@/lib/format';
import { mapColor } from '@/lib/maps';
import { useMapName } from '@/hooks/use-map-name';
import { toast } from '@/lib/toast';
import { cn } from '@/lib/utils';

// 列表视图的列宽，表头和每一行共用。每行是独立的 grid，列宽都要固定（操作列不能用 auto），否则各行对不齐
export const ROW_GRID = 'lg:grid lg:grid-cols-[minmax(0,1fr)_10.5rem_5rem_9rem_12rem_15rem] lg:items-center lg:gap-4';

/** 列表视图的表头，只在宽屏显示 */
export function ServerRowHeader() {
  const t = useTranslations('serverList.columns');
  return (
    <div className={cn('hidden px-4 pb-1.5 text-xs text-muted-foreground', ROW_GRID)}>
      <span>{t('server')}</span>
      <span>{t('status')}</span>
      <span>{t('players')}</span>
      <span>{t('resources')}</span>
      <span>{t('address')}</span>
      <span className="text-right">{t('actions')}</span>
    </div>
  );
}

/** 列表视图的一行：一台服务器的关键信息和常用操作，宽屏时按列对齐，窄屏时分行排列 */
export function ServerRow({ server, publicHost }: { server: Server; publicHost: string }) {
  const t = useTranslations('serverList.columns');
  const tCard = useTranslations('serverCard');
  const tConn = useTranslations('connection');
  const mapName = useMapName();
  const controls = useServerControls(server);

  const runtime = server.status === 'running' ? server.runtime : undefined;
  const ready = runtime?.phase === 'ready';
  const color = mapColor(server.map);
  const address = publicHost ? `${publicHost}:${server.port}` : '';

  const copyAddress = async () => {
    if (await copyText(address)) toast.success(tConn('copied', { text: address }));
    else toast.error(tConn('copyFailed'));
  };

  const progress =
    runtime?.phase === 'downloading'
      ? runtime.progress ?? 0
      : runtime?.phase === 'installing_mods' && runtime.mod_total
        ? ((runtime.mod_current ?? 0) / runtime.mod_total) * 100
        : null;

  return (
    <div
      className={cn(
        'relative flex flex-col gap-3 rounded-sm border bg-card px-4 py-3 backdrop-blur-sm transition-colors hover:border-primary/50',
        ROW_GRID,
      )}
      // 左边的地图色条，停止的服务器变灰
      style={{ boxShadow: `inset 3px 0 0 ${server.status === 'stopped' ? `${color}66` : color}` }}
    >
      {/* 名称、地图和标记 */}
      <div className="min-w-0">
        <div className="flex items-center gap-2">
          <Link
            href={`/servers/${server.id}`}
            className="truncate font-semibold hover:text-primary hover:underline"
            title={server.session_name}
          >
            {server.session_name}
          </Link>
          {server.last_error && (
            <Link href={`/servers/${server.id}?tab=logs`} title={`${tCard('lastError')}: ${server.last_error}`} className="shrink-0 text-red-600">
              <AlertTriangle className="size-4" />
            </Link>
          )}
        </div>
        <div className="mt-1 flex flex-wrap items-center gap-1.5 text-xs">
          <span className="font-medium" style={{ color }}>
            {mapName(server.map)}
          </span>
          <ServerTags server={server} />
        </div>
      </div>

      {/* 状态：窄屏时和人数、资源排在同一行 */}
      <div className="flex flex-wrap items-center gap-x-4 gap-y-2 lg:contents">
        <div className="min-w-0 space-y-1.5">
          <StatusBadge server={server} />
          {progress !== null && <Progress value={progress} className="h-1 w-32" />}
        </div>
        <div className="num text-sm" title={tCard('players')}>
          <span className="text-xs text-muted-foreground lg:hidden">{t('players')} </span>
          {ready ? (
            <>
              <span className="font-semibold text-primary">{runtime.players}</span>
              <span className="text-muted-foreground">/{runtime.max_players || server.max_players}</span>
            </>
          ) : (
            <span className="text-muted-foreground">-</span>
          )}
        </div>
        <div className={cn('num text-xs', !runtime && 'hidden lg:block')}>
          {runtime ? (
            <>
              CPU {formatPercent(runtime.cpu_percent, 0)}
              <span className="mx-1.5 text-gray-300">|</span>
              {formatBytes(runtime.memory_bytes)}
            </>
          ) : (
            <span className="text-muted-foreground">-</span>
          )}
        </div>
        <div className="min-w-0">
          {address && (
            <button
              type="button"
              onClick={copyAddress}
              title={tConn('copyAddress')}
              className="group flex max-w-full items-center gap-1.5 rounded px-1.5 py-0.5 font-mono text-xs transition-colors hover:bg-accent"
            >
              <span className="truncate">{address}</span>
              <Copy className="size-3.5 shrink-0 text-muted-foreground group-hover:text-primary" />
            </button>
          )}
        </div>
      </div>

      <div className="flex items-center justify-between gap-2 border-t pt-2.5 lg:justify-end lg:border-t-0 lg:pt-0">
        <PowerButtons controls={controls} compact />
        <ServerShortcuts server={server} />
      </div>
      {controls.dialogs}
    </div>
  );
}
