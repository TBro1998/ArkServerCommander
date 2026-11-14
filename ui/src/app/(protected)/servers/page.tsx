"use client";

import { useCallback, useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import { useTranslations } from 'next-intl';
import { Plus, Server as ServerIcon, RefreshCw } from 'lucide-react';
import {
  useServers,
  useServersLoaded,
  useServersIsLoading,
  usePublicHost,
  serversActions,
  isTransitioning,
  type Server,
} from '@/stores/servers';
import { ServerCard } from '@/components/servers/ServerCard';
import { PageHeader } from '@/components/PageHeader';
import { Button } from '@/components/ui/button';
import { PageLoader } from '@/components/ui/page-loader';
import { usePolling } from '@/hooks/use-polling';
import { useApiError } from '@/lib/api-error';
import { toast } from '@/lib/toast';
import { cn } from '@/lib/utils';

// 有服务器在启停时快速刷新，平时也定期刷新，以便发现崩溃、阶段变化和在线人数
const TRANSITION_POLL_INTERVAL = 3000;
const IDLE_POLL_INTERVAL = 10000;

type Filter = 'all' | 'active' | 'stopped';
const FILTERS: { key: Filter; match: (s: Server) => boolean }[] = [
  { key: 'all', match: () => true },
  // "运行中"包括启动中、停止中，这些服务器都占着资源
  { key: 'active', match: (s) => s.status !== 'stopped' },
  { key: 'stopped', match: (s) => s.status === 'stopped' },
];

export default function ServersPage() {
  const t = useTranslations('servers');
  const tList = useTranslations('serverList');
  const router = useRouter();
  const servers = useServers();
  const loaded = useServersLoaded();
  const isLoading = useServersIsLoading();
  const publicHost = usePublicHost();
  const apiError = useApiError();
  const { fetchServers, fetchPublicHost } = serversActions;
  const [filter, setFilter] = useState<Filter>('all');

  const refresh = useCallback(() => fetchServers(), [fetchServers]);

  useEffect(() => {
    fetchServers().catch((err) => toast.error(apiError(err, t('getServerListFailed'))));
    fetchPublicHost();
  }, [fetchServers, fetchPublicHost, apiError, t]);

  const hasTransitioning = servers.some((s) => isTransitioning(s.status));
  // 运行中但还没有可以连接（下载、安装 Mod、加载地图）时也快速刷新
  const hasPreparing = servers.some((s) => s.status === 'running' && s.runtime?.phase !== 'ready');
  usePolling(refresh, hasTransitioning || hasPreparing ? TRANSITION_POLL_INTERVAL : IDLE_POLL_INTERVAL);

  const handleAddServer = () => router.push('/servers/new');

  const runningCount = servers.filter((s) => s.status === 'running').length;
  const playerCount = servers.reduce((sum, s) => sum + (s.status === 'running' ? s.runtime?.players ?? 0 : 0), 0);
  const visible = servers.filter(FILTERS.find((f) => f.key === filter)!.match);

  const addButton = (
    <Button onClick={handleAddServer}>
      <Plus className="h-4 w-4" />
      {t('addServer')}
    </Button>
  );

  return (
    <div className="w-full">
      <PageHeader
        title={t('title')}
        description={servers.length > 0 ? tList('summary', { total: servers.length, running: runningCount, players: playerCount }) : undefined}
        actions={
          <>
            <Button
              variant="outline"
              size="icon"
              onClick={() => refresh().catch(() => {})}
              disabled={isLoading}
              title={tList('refresh')}
              aria-label={tList('refresh')}
            >
              <RefreshCw className={cn('h-4 w-4', isLoading && 'animate-spin')} />
            </Button>
            {addButton}
          </>
        }
      />

      {!loaded ? (
        <PageLoader />
      ) : servers.length === 0 ? (
        <div className="rounded-lg border border-dashed bg-card/60 px-6 py-16 text-center">
          <div className="mx-auto mb-5 flex h-14 w-14 items-center justify-center rounded-full bg-accent">
            <ServerIcon className="h-7 w-7 text-primary" />
          </div>
          <h2 className="font-display mb-2 text-xl font-semibold">{t('noServers')}</h2>
          <p className="mx-auto mb-6 max-w-md text-sm text-muted-foreground">{t('noServersDesc')}</p>
          {addButton}
        </div>
      ) : (
        <>
          {/* 状态筛选：服务器多时快速找到运行中或已停止的 */}
          <div role="group" aria-label={tList('filter')} className="mb-4 inline-flex rounded-lg border bg-card p-0.5">
            {FILTERS.map((f) => {
              const count = servers.filter(f.match).length;
              const selected = filter === f.key;
              return (
                <button
                  key={f.key}
                  type="button"
                  aria-pressed={selected}
                  onClick={() => setFilter(f.key)}
                  className={cn(
                    'flex items-center gap-1.5 rounded-md px-3 py-1.5 text-sm font-medium transition-colors',
                    selected ? 'bg-primary text-primary-foreground' : 'text-muted-foreground hover:bg-accent hover:text-foreground',
                  )}
                >
                  {tList(`filters.${f.key}`)}
                  <span className={cn('num text-xs', selected ? 'opacity-80' : 'opacity-60')}>{count}</span>
                </button>
              );
            })}
          </div>

          {visible.length === 0 ? (
            <p className="rounded-lg border border-dashed px-6 py-10 text-center text-sm text-muted-foreground">{tList('filterEmpty')}</p>
          ) : (
            <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3 2xl:grid-cols-4">
              {visible.map((server) => (
                <ServerCard key={server.id} server={server} publicHost={publicHost} />
              ))}
            </div>
          )}
        </>
      )}
    </div>
  );
}
