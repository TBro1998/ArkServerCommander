"use client";

import { useState } from 'react';
import { useTranslations } from 'next-intl';
import Link from 'next/link';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Progress } from '@/components/ui/progress';
import { Lock, Eye, EyeOff, Copy, AlertTriangle, Users, Cpu, MemoryStick, Map as MapIcon } from 'lucide-react';
import { serversActions, type Server } from '@/stores/servers';
import { StatusBadge } from './StatusBadge';
import { ServerActions } from './ServerActions';
import { ConnectionInfo } from './ConnectionInfo';
import { copyText, formatBytes, formatPercent } from '@/lib/format';
import { toast } from '@/lib/toast';
import { useApiError } from '@/lib/api-error';

interface ServerCardProps {
  server: Server;
  canStartServer: boolean;
  publicHost: string;
}

export function ServerCard({ server, canStartServer, publicHost }: ServerCardProps) {
  const t = useTranslations('serverCard');
  const tMaps = useTranslations('servers.edit.maps');
  const apiError = useApiError();
  const [password, setPassword] = useState<string | null>(null);
  const [loadingPassword, setLoadingPassword] = useState(false);

  const mapName = tMaps.has(server.map) ? tMaps(server.map) : server.map;
  const runtime = server.status === 'running' ? server.runtime : undefined;

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
    <Card className="flex h-full flex-col border-0 shadow-sm transition-all duration-200 hover:shadow-lg">
      <CardHeader className="pb-3">
        <div className="mb-2 flex items-start justify-between gap-2">
          <CardTitle className="min-w-0 text-base font-semibold">
            <Link href={`/servers/${server.id}`} className="block truncate hover:text-blue-700 hover:underline" title={server.session_name}>
              {server.session_name}
            </Link>
          </CardTitle>
          <ServerActions server={server} canStart={canStartServer} />
        </div>
        <div className="flex flex-wrap items-center gap-1">
          <StatusBadge server={server} />
          {server.image_outdated && (
            <span className="rounded-md bg-yellow-100 px-2 py-0.5 text-xs text-yellow-800" title={t('imageOutdatedTip')}>
              {t('imageOutdated')}
            </span>
          )}
          {server.config_pending && (
            <span className="rounded-md bg-blue-100 px-2 py-0.5 text-xs text-blue-800" title={t('configPendingTip')}>
              {t('configPending')}
            </span>
          )}
          {server.cluster_id && (
            <span className="rounded-md bg-purple-100 px-2 py-0.5 text-xs text-purple-800" title={t('clusterTip')}>
              {t('cluster', { id: server.cluster_id })}
            </span>
          )}
        </div>
      </CardHeader>
      <CardContent className="flex flex-1 flex-col space-y-3 pt-0">
        {/* 最近一次启停失败的原因，下次启停时清空 */}
        {server.last_error && (
          <div className="flex items-start gap-1.5 rounded-lg bg-red-50 p-2 text-xs text-red-700">
            <AlertTriangle className="mt-0.5 h-3.5 w-3.5 flex-shrink-0" />
            <span className="break-all">
              <span className="font-medium">{t('lastError')}:</span> {server.last_error}
            </span>
          </div>
        )}

        {/* 下载服务端时显示进度条 */}
        {runtime?.phase === 'downloading' && (
          <div className="space-y-1">
            <Progress value={runtime.progress ?? 0} className="h-1.5" />
            <p className="text-xs text-gray-500">{t('downloadHint')}</p>
          </div>
        )}
        {runtime?.phase === 'launching' && <p className="text-xs text-gray-500">{t('launchingHint')}</p>}
        {server.status === 'starting' && server.runtime?.phase === 'copying' && (
          <p className="text-xs text-gray-500">{t('copyingHint', { source: server.runtime.copy_source ?? '' })}</p>
        )}

        {/* 运行中：人数和资源 */}
        {runtime && (
          <div className="grid grid-cols-3 gap-2 rounded-lg bg-gray-50 p-2 text-xs">
            <div className="text-center" title={t('players')}>
              <Users className="mx-auto mb-0.5 h-3.5 w-3.5 text-green-600" />
              <span className="font-semibold">
                {runtime.phase === 'ready' ? `${runtime.players}/${runtime.max_players || server.max_players}` : '-'}
              </span>
            </div>
            <div className="text-center" title="CPU">
              <Cpu className="mx-auto mb-0.5 h-3.5 w-3.5 text-blue-600" />
              <span className="font-semibold">{formatPercent(runtime.cpu_percent, 0)}</span>
            </div>
            <div className="text-center" title={t('memory')}>
              <MemoryStick className="mx-auto mb-0.5 h-3.5 w-3.5 text-purple-600" />
              <span className="font-semibold">{formatBytes(runtime.memory_bytes)}</span>
            </div>
          </div>
        )}

        <div className="space-y-1.5 text-sm">
          <div className="flex items-center justify-between">
            <span className="flex items-center text-gray-600">
              <MapIcon className="mr-1.5 h-3.5 w-3.5 text-green-600" />
              {t('map')}
            </span>
            <span className="ml-2 truncate font-medium">{mapName}</span>
          </div>
          <div className="flex items-center justify-between text-xs text-gray-500">
            <span>{t('ports')}</span>
            <span className="font-mono">
              {server.port} / {server.query_port} / {server.rcon_port}
            </span>
          </div>
        </div>

        <div className="mt-auto space-y-2">
          <ConnectionInfo server={server} host={publicHost} compact />

          {/* 管理员密码 */}
          <div className="flex items-center justify-between rounded-lg bg-amber-50 px-2 py-1.5">
            <span className="flex items-center text-xs font-medium text-gray-700">
              <Lock className="mr-1.5 h-3.5 w-3.5 text-amber-600" />
              {t('adminPassword')}
            </span>
            <span className="flex min-w-0 items-center gap-1">
              <span className="truncate font-mono text-xs">{password ?? '••••••••'}</span>
              <Button
                variant="ghost"
                size="sm"
                className="h-6 w-6 p-0"
                onClick={togglePassword}
                disabled={loadingPassword}
                title={password !== null ? t('hidePassword') : t('showPassword')}
              >
                {password !== null ? <EyeOff className="h-3 w-3" /> : <Eye className="h-3 w-3" />}
              </Button>
              <Button variant="ghost" size="sm" className="h-6 w-6 p-0" onClick={copyPassword} disabled={loadingPassword} title={t('copyPassword')}>
                <Copy className="h-3 w-3" />
              </Button>
            </span>
          </div>
        </div>
      </CardContent>
    </Card>
  );
}
