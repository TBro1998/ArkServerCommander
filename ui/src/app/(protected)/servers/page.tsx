"use client";

import { useCallback, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { useTranslations } from 'next-intl';
import { Plus, Loader2, Server as ServerIcon, AlertCircle, RefreshCw } from 'lucide-react';
import {
  useServers,
  useServersLoaded,
  useServersIsLoading,
  useImageStatus,
  usePublicHost,
  serversActions,
  isTransitioning,
} from '@/stores/servers';
import { ServerCard } from '@/components/servers/ServerCard';
import { Button } from '@/components/ui/button';
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert';
import { usePolling } from '@/hooks/use-polling';
import { useApiError } from '@/lib/api-error';
import { toast } from '@/lib/toast';

// 有服务器在启停时快速刷新，平时也定期刷新，以便发现崩溃、阶段变化和在线人数
const TRANSITION_POLL_INTERVAL = 3000;
const IDLE_POLL_INTERVAL = 10000;

export default function ServersPage() {
  const t = useTranslations('servers');
  const tList = useTranslations('serverList');
  const router = useRouter();
  const servers = useServers();
  const loaded = useServersLoaded();
  const isLoading = useServersIsLoading();
  const imageStatus = useImageStatus();
  const publicHost = usePublicHost();
  const apiError = useApiError();
  const { fetchServers, getImageStatus, fetchPublicHost } = serversActions;

  const refresh = useCallback(() => fetchServers(), [fetchServers]);

  useEffect(() => {
    fetchServers().catch((err) => toast.error(apiError(err, t('getServerListFailed'))));
    getImageStatus().catch(() => {});
    fetchPublicHost();
  }, [fetchServers, getImageStatus, fetchPublicHost, apiError, t]);

  const hasTransitioning = servers.some((s) => isTransitioning(s.status));
  // 运行中但还没有可以连接（下载、安装 Mod、加载地图）时也快速刷新
  const hasPreparing = servers.some((s) => s.status === 'running' && s.runtime?.phase !== 'ready');
  usePolling(refresh, hasTransitioning || hasPreparing ? TRANSITION_POLL_INTERVAL : IDLE_POLL_INTERVAL);

  const handleAddServer = () => router.push('/servers/new');
  const canCreate = imageStatus?.can_create_server ?? false;

  const runningCount = servers.filter((s) => s.status === 'running').length;
  const playerCount = servers.reduce((sum, s) => sum + (s.status === 'running' ? s.runtime?.players ?? 0 : 0), 0);

  return (
    <div className="w-full max-w-none py-4 sm:py-8">
      <div className="mb-6">
        <div className="flex flex-wrap items-center justify-between gap-3">
          <div>
            <h1 className="text-2xl font-bold text-gray-900 lg:text-3xl">{t('title')}</h1>
            {servers.length > 0 && (
              <p className="mt-1 text-sm text-gray-500">
                {tList('summary', { total: servers.length, running: runningCount, players: playerCount })}
              </p>
            )}
          </div>
          <div className="flex items-center gap-2">
            <Button variant="outline" size="icon" onClick={() => refresh().catch(() => {})} disabled={isLoading} title={tList('refresh')}>
              <RefreshCw className={`h-4 w-4 ${isLoading ? 'animate-spin' : ''}`} />
            </Button>
            <Button onClick={handleAddServer} disabled={!canCreate}>
              <Plus className="mr-2 h-4 w-4" />
              {t('addServer')}
            </Button>
          </div>
        </div>
        {imageStatus && !imageStatus.can_start_server && (
          <Alert variant="destructive" className="mt-4">
            <AlertCircle className="h-4 w-4" />
            <AlertTitle>{imageStatus.any_pulling ? t('imageDownloading') : t('imageNotReady')}</AlertTitle>
            <AlertDescription>
              {imageStatus.any_pulling ? t('imageDownloadingDesc') : tList('imageNotReadyGoHome')}
            </AlertDescription>
          </Alert>
        )}
      </div>

      {!loaded ? (
        <div className="py-12 text-center">
          <Loader2 className="mx-auto mb-4 h-8 w-8 animate-spin text-blue-600" />
        </div>
      ) : servers.length === 0 ? (
        <div className="px-4 py-16 text-center">
          <div className="mx-auto mb-6 flex h-16 w-16 items-center justify-center rounded-full bg-gray-100">
            <ServerIcon className="h-8 w-8 text-gray-400" />
          </div>
          <h3 className="mb-3 text-xl font-medium text-gray-900">{t('noServers')}</h3>
          <p className="mx-auto mb-8 max-w-md text-gray-500">{t('noServersDesc')}</p>
          <Button onClick={handleAddServer} disabled={!canCreate}>
            <Plus className="mr-2 h-4 w-4" />
            {t('addServer')}
          </Button>
        </div>
      ) : (
        <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3 lg:gap-6 xl:grid-cols-4 2xl:grid-cols-5">
          {servers.map((server) => (
            <ServerCard
              key={server.id}
              server={server}
              canStartServer={imageStatus?.can_start_server ?? false}
              publicHost={publicHost}
            />
          ))}
        </div>
      )}
    </div>
  );
}
