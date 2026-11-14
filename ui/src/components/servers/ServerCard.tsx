"use client";

import { useState } from 'react';
import { useTranslations } from 'next-intl';
import Link from 'next/link';
import { Button } from '@/components/ui/button';
import { Progress } from '@/components/ui/progress';
import { Eye, EyeOff, Copy, AlertTriangle, Users, ChevronRight } from 'lucide-react';
import { serversActions, type Server } from '@/stores/servers';
import { StatusBadge } from './StatusBadge';
import { ServerActions } from './ServerActions';
import { ConnectionInfo } from './ConnectionInfo';
import { copyText, formatBytes, formatPercent } from '@/lib/format';
import { mapColor } from '@/lib/maps';
import { useMapName } from '@/hooks/use-map-name';
import { toast } from '@/lib/toast';
import { useApiError } from '@/lib/api-error';
import { cn } from '@/lib/utils';

interface ServerCardProps {
  server: Server;
  publicHost: string;
}

// 地图色带：以地图区分服务器；停止的服务器色带变灰，运行中的更醒目
export function MapBand({ server, className, children }: { server: Server; className?: string; children?: React.ReactNode }) {
  const mapName = useMapName();
  const live = server.status !== 'stopped';
  const color = mapColor(server.map);
  return (
    <div
      className={cn('relative flex items-center justify-between gap-3 border-b px-4 text-white', !live && 'opacity-60 saturate-[.25]', className)}
      style={{
        background: `linear-gradient(90deg, ${color}d9 0%, ${color}59 55%, transparent 100%)`,
        boxShadow: `inset 3px 0 0 ${color}`,
        borderBottomColor: `${color}80`,
      }}
    >
      <span className="font-display hud-caps relative truncate text-lg font-semibold [text-shadow:0_1px_6px_rgb(0_0_0/0.6)]">
        {mapName(server.map)}
      </span>
      {children && <span className="relative flex-shrink-0">{children}</span>}
    </div>
  );
}

export function ServerCard({ server, publicHost }: ServerCardProps) {
  const t = useTranslations('serverCard');
  const tActions = useTranslations('serverActions');
  const apiError = useApiError();
  const [password, setPassword] = useState<string | null>(null);
  const [loadingPassword, setLoadingPassword] = useState(false);

  const runtime = server.status === 'running' ? server.runtime : undefined;
  const ready = runtime?.phase === 'ready';

  // 密码不随列表返回，显示或复制时按需请求，且不做缓存
  const fetchPassword = async (): Promise<string | null> => {
    setLoadingPassword(true);
    try {
      return await serversActions.getServerPassword(server.id);
    } catch (err) {
      toast.error(apiError(err, t('passwordFailed')));
      return null;
    } finally {
      setLoadingPassword(false);
    }
  };

  const togglePassword = async () => {
    if (password !== null) {
      setPassword(null);
      return;
    }
    setPassword(await fetchPassword());
  };

  const copyPassword = async () => {
    const value = await fetchPassword();
    if (value === null) return;
    if (await copyText(value)) toast.success(t('passwordCopied'));
    else toast.error(t('copyFailed'));
  };

  return (
    <article className="hud relative flex h-full flex-col rounded-sm border bg-card backdrop-blur-sm transition-colors hover:border-primary/50">
      <MapBand server={server} className="h-12">
        {ready && (
          <span className="font-display flex items-center gap-1.5 text-lg font-semibold text-primary [text-shadow:0_0_10px_rgb(57_195_236/0.6)]" title={t('players')}>
            <Users className="h-4 w-4 opacity-80" />
            {runtime.players}
            <span className="text-sm font-normal opacity-75">/{runtime.max_players || server.max_players}</span>
          </span>
        )}
      </MapBand>

      <div className="flex flex-1 flex-col gap-3 p-4">
        <div className="space-y-2">
          <h3 className="min-w-0 text-base leading-snug font-semibold">
            <Link href={`/servers/${server.id}`} className="block truncate hover:text-primary hover:underline" title={server.session_name}>
              {server.session_name}
            </Link>
          </h3>
          <div className="flex flex-wrap items-center gap-1.5">
            <StatusBadge server={server} />
            {server.arkapi && (
              <span className="rounded-full border border-purple-200 bg-purple-50 px-2 py-0.5 text-xs text-purple-800" title={t('arkapiTip')}>
                {t('arkapi')}
              </span>
            )}
            {server.image_outdated && (
              <span className="rounded-full border border-amber-200 bg-amber-50 px-2 py-0.5 text-xs text-amber-800" title={t('imageOutdatedTip')}>
                {t('imageOutdated')}
              </span>
            )}
            {server.config_pending && (
              <span className="rounded-full border border-blue-200 bg-blue-50 px-2 py-0.5 text-xs text-blue-800" title={t('configPendingTip')}>
                {t('configPending')}
              </span>
            )}
            {!!server.mods_outdated?.length && (
              <Link
                href={`/servers/${server.id}`}
                className="rounded-full border border-amber-200 bg-amber-50 px-2 py-0.5 text-xs text-amber-800 hover:underline"
                title={t('modsOutdatedTip')}
              >
                {t('modsOutdated', { count: server.mods_outdated.length })}
              </Link>
            )}
            {server.cluster_id && (
              <span className="rounded-full border px-2 py-0.5 text-xs text-muted-foreground" title={t('clusterTip')}>
                {t('cluster', { id: server.cluster_id })}
              </span>
            )}
          </div>
        </div>

        {/* 最近一次启停失败的原因，下次启停时清空 */}
        {server.last_error && (
          <div className="flex items-start gap-2 rounded-md border border-red-200 bg-red-50 p-2 text-xs text-red-800">
            <AlertTriangle className="mt-0.5 h-3.5 w-3.5 flex-shrink-0" />
            <span className="min-w-0 break-words">
              <span className="font-medium">{t('lastError')}:</span>{' '}
              {server.last_error}{' '}
              <Link href={`/servers/${server.id}?tab=logs`} className="font-medium underline underline-offset-2">
                {t('viewLogs')}
              </Link>
            </span>
          </div>
        )}

        {/* 准备阶段：进度和说明 */}
        {runtime?.phase === 'downloading' && (
          <div className="space-y-1.5">
            <Progress value={runtime.progress ?? 0} className="h-1.5" />
            <p className="text-xs text-muted-foreground">{t('downloadHint')}</p>
          </div>
        )}
        {runtime?.phase === 'installing_mods' && runtime.mod_total ? (
          <Progress value={((runtime.mod_current ?? 0) / runtime.mod_total) * 100} className="h-1.5" />
        ) : null}
        {runtime?.phase === 'launching' && <p className="text-xs text-muted-foreground">{t('launchingHint')}</p>}
        {server.status === 'starting' && server.runtime?.phase === 'copying' && (
          <p className="text-xs text-muted-foreground">{t('copyingHint', { source: server.runtime.copy_source ?? '' })}</p>
        )}
        {server.status === 'starting' && server.runtime?.phase === 'pulling_image' && (
          <p className="text-xs text-muted-foreground">{t('pullingImageHint', { image: server.runtime.image ?? '' })}</p>
        )}

        <dl className="grid grid-cols-[auto_1fr] items-center gap-x-4 gap-y-1 text-xs">
          {runtime && (
            <>
              <dt className="text-muted-foreground">{t('resources')}</dt>
              <dd className="num text-right font-medium">
                CPU {formatPercent(runtime.cpu_percent, 0)}
                <span className="mx-1.5 text-gray-300">|</span>
                {formatBytes(runtime.memory_bytes)}
              </dd>
            </>
          )}
          <dt className="text-muted-foreground">{t('ports')}</dt>
          <dd className="num text-right font-medium">
            {server.port} / {server.query_port} / {server.rcon_port}
          </dd>
          <dt className="text-muted-foreground">{t('adminPassword')}</dt>
          <dd className="flex min-w-0 items-center justify-end gap-0.5">
            <span className="truncate font-mono">{password ?? '••••••••'}</span>
            <Button
              variant="ghost"
              size="sm"
              className="h-7 w-7 p-0"
              onClick={togglePassword}
              disabled={loadingPassword}
              aria-label={password !== null ? t('hidePassword') : t('showPassword')}
              title={password !== null ? t('hidePassword') : t('showPassword')}
            >
              {password !== null ? <EyeOff className="h-3.5 w-3.5" /> : <Eye className="h-3.5 w-3.5" />}
            </Button>
            <Button
              variant="ghost"
              size="sm"
              className="h-7 w-7 p-0"
              onClick={copyPassword}
              disabled={loadingPassword}
              aria-label={t('copyPassword')}
              title={t('copyPassword')}
            >
              <Copy className="h-3.5 w-3.5" />
            </Button>
          </dd>
        </dl>

        <div className="mt-auto">
          <ConnectionInfo server={server} host={publicHost} compact />
        </div>
      </div>

      <footer className="flex items-center justify-between gap-2 border-t bg-black/20 px-4 py-2.5">
        <ServerActions server={server} />
        <Button asChild variant="ghost" size="sm" className="text-muted-foreground">
          <Link href={`/servers/${server.id}`}>
            {tActions('details')}
            <ChevronRight className="h-4 w-4" />
          </Link>
        </Button>
      </footer>
    </article>
  );
}
