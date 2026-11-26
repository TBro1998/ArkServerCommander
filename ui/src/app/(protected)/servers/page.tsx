"use client";

import { useCallback, useEffect, useState } from 'react';
import Link from 'next/link';
import { useRouter } from 'next/navigation';
import { useTranslations } from 'next-intl';
import { Plus, Server as ServerIcon, RefreshCw, LayoutGrid, List } from 'lucide-react';
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
import { ServerRow, ServerRowHeader } from '@/components/servers/ServerRow';
import { useLocalStorage } from '@/hooks/use-local-storage';
import { PageHeader } from '@/components/PageHeader';
import { Button } from '@/components/ui/button';
import { PageLoader } from '@/components/ui/page-loader';
import { usePolling } from '@/hooks/use-polling';
import { useApiError } from '@/lib/api-error';
import { toast } from '@/lib/toast';
import { cn } from '@/lib/utils';
import { GAMES, GAME_NAMES, gameOf, type Game } from '@/lib/games';

// 有服务器在启停时快速刷新，平时也定期刷新，以便发现崩溃、阶段变化和在线人数
const TRANSITION_POLL_INTERVAL = 3000;
const IDLE_POLL_INTERVAL = 10000;

type Filter = 'all' | 'active' | 'stopped';
// 显示方式：卡片信息更全，列表一屏能看到更多服务器
const VIEWS = ['cards', 'list'] as const;
type View = (typeof VIEWS)[number];
// 每个游戏一个分区的强调色：ASA 橙色（和 GameBadge 一致），ASE 青色
const GAME_ACCENT: Record<Game, string> = {
  ase: 'bg-primary',
  asa: 'bg-orange-500',
};
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
  const [view, setView] = useLocalStorage<View>('servers-view', 'cards', VIEWS);

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
  // ASE 和 ASA 的服务器分区显示：地图、Mod、插件都不通用，混在一起不好找。没有服务器的游戏不显示分区
  const sections = GAMES.map((game) => ({
    game,
    all: servers.filter((s) => gameOf(s) === game),
    shown: visible.filter((s) => gameOf(s) === game),
  })).filter((g) => g.all.length > 0);

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
          <div className="mb-4 flex flex-wrap items-center justify-between gap-2">
            {/* 状态筛选：服务器多时快速找到运行中或已停止的 */}
            <div role="group" aria-label={tList('filter')} className="inline-flex rounded-lg border bg-card p-0.5">
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

            <div role="group" aria-label={tList('view')} className="inline-flex rounded-lg border bg-card p-0.5">
              {VIEWS.map((v) => {
                const Icon = v === 'cards' ? LayoutGrid : List;
                const selected = view === v;
                return (
                  <button
                    key={v}
                    type="button"
                    aria-pressed={selected}
                    onClick={() => setView(v)}
                    title={tList(`views.${v}`)}
                    className={cn(
                      'flex items-center gap-1.5 rounded-md px-2.5 py-1.5 text-sm font-medium transition-colors',
                      selected ? 'bg-primary text-primary-foreground' : 'text-muted-foreground hover:bg-accent hover:text-foreground',
                    )}
                  >
                    <Icon className="size-4" />
                    <span className="hidden sm:inline">{tList(`views.${v}`)}</span>
                  </button>
                );
              })}
            </div>
          </div>

          <div className="space-y-8">
            {sections.map(({ game, all, shown }) => {
              const running = all.filter((s) => s.status === 'running').length;
              const players = all.reduce((sum, s) => sum + (s.status === 'running' ? s.runtime?.players ?? 0 : 0), 0);
              return (
                <section key={game} aria-labelledby={`game-${game}`}>
                  <div className="mb-3 flex flex-wrap items-center justify-between gap-2 border-b pb-2">
                    <div className="flex min-w-0 items-center gap-3">
                      <span aria-hidden className={cn('h-7 w-1 flex-shrink-0 rounded-full', GAME_ACCENT[game])} />
                      <div className="min-w-0">
                        <h2 id={`game-${game}`} className="font-display flex items-baseline gap-2 text-lg font-semibold">
                          <span className="hud-caps">{game.toUpperCase()}</span>
                          <span className="truncate text-sm font-normal text-muted-foreground">{GAME_NAMES[game]}</span>
                        </h2>
                        <p className="text-xs text-muted-foreground">{tList('gameSummary', { total: all.length, running, players })}</p>
                      </div>
                    </div>
                    <Button asChild variant="ghost" size="sm" className="text-muted-foreground">
                      <Link href={`/servers/new?game=${game}`}>
                        <Plus className="h-4 w-4" />
                        {tList('addGameServer', { game: game.toUpperCase() })}
                      </Link>
                    </Button>
                  </div>
                  {shown.length === 0 ? (
                    <p className="rounded-lg border border-dashed px-6 py-6 text-center text-sm text-muted-foreground">{tList('filterEmpty')}</p>
                  ) : view === 'list' ? (
                    <div className="space-y-2">
                      <ServerRowHeader />
                      {shown.map((server) => (
                        <ServerRow key={server.id} server={server} publicHost={publicHost} />
                      ))}
                    </div>
                  ) : (
                    <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3 2xl:grid-cols-4">
                      {shown.map((server) => (
                        <ServerCard key={server.id} server={server} publicHost={publicHost} />
                      ))}
                    </div>
                  )}
                </section>
              );
            })}
          </div>
        </>
      )}
    </div>
  );
}
