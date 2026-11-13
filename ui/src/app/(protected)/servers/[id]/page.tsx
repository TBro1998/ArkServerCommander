"use client";

import { useCallback, useEffect, useState } from 'react';
import { useTranslations } from 'next-intl';
import { useParams, useRouter, useSearchParams } from 'next/navigation';
import Link from 'next/link';
import { ArrowLeft, AlertTriangle, Users } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { Progress } from '@/components/ui/progress';
import { PageLoader } from '@/components/ui/page-loader';
import {
  useServers,
  useServersLoaded,
  useImageStatus,
  usePublicHost,
  serversActions,
  isTransitioning,
  type Server,
} from '@/stores/servers';
import { StatusBadge } from '@/components/servers/StatusBadge';
import { ServerActions } from '@/components/servers/ServerActions';
import { ConnectionInfo } from '@/components/servers/ConnectionInfo';
import { LogViewer } from '@/components/servers/LogViewer';
import { ConsolePanel } from '@/components/servers/ConsolePanel';
import { BackupsPanel } from '@/components/servers/BackupsPanel';
import { SchedulesPanel } from '@/components/servers/SchedulesPanel';
import { PluginsPanel } from '@/components/servers/PluginsPanel';
import { ServerModsPanel } from '@/components/servers/ServerModsPanel';
import { MapBand } from '@/components/servers/ServerCard';
import { useMapName } from '@/hooks/use-map-name';
import { usePolling } from '@/hooks/use-polling';
import { formatBytes, formatPercent } from '@/lib/format';

const TABS = ['overview', 'console', 'logs', 'backups', 'schedules', 'plugins'] as const;
type Tab = (typeof TABS)[number];

function Row({ label, children }: { label: string; children: React.ReactNode }) {
  return (
    <div className="flex items-start justify-between gap-4 py-1.5 text-sm">
      <span className="flex-shrink-0 text-muted-foreground">{label}</span>
      <span className="min-w-0 text-right font-medium break-all">{children}</span>
    </div>
  );
}

function Overview({ server, publicHost }: { server: Server; publicHost: string }) {
  const t = useTranslations('detail');
  const mapName = useMapName();
  const tPhase = useTranslations('phase');
  const copying = server.status === 'starting' && server.runtime?.phase === 'copying';
  const runtime = server.status === 'running' || copying ? server.runtime : undefined;

  return (
    <div className="grid gap-4 lg:grid-cols-2">
      <Card>
        <CardHeader className="pb-2">
          <CardTitle className="font-display text-base">{t('runtime')}</CardTitle>
        </CardHeader>
        <CardContent className="space-y-3">
          {!runtime ? (
            <p className="text-sm text-muted-foreground">{server.status === 'running' ? t('collecting') : t('notRunning')}</p>
          ) : (
            <>
              <div className="text-sm">{tPhase(`${runtime.phase || 'unknown'}Desc`)}</div>
              {runtime.phase === 'downloading' && <Progress value={runtime.progress ?? 0} className="h-2" />}
              {runtime.phase === 'installing_mods' && runtime.mod_total ? (
                <Progress value={((runtime.mod_current ?? 0) / runtime.mod_total) * 100} className="h-2" />
              ) : null}
              {!copying && (
                <div className="grid grid-cols-3 gap-px overflow-hidden rounded-md border bg-border">
                  {[
                    {
                      label: t('players'),
                      value: runtime.phase === 'ready' ? `${runtime.players}/${runtime.max_players || server.max_players}` : '-',
                    },
                    { label: 'CPU', value: formatPercent(runtime.cpu_percent) },
                    { label: t('memory'), value: formatBytes(runtime.memory_bytes) },
                  ].map((item) => (
                    <div key={item.label} className="bg-card px-3 py-2.5">
                      <div className="text-xs text-muted-foreground">{item.label}</div>
                      <div className="font-display text-xl font-semibold">{item.value}</div>
                    </div>
                  ))}
                </div>
              )}
              {runtime.server_name && <Row label={t('advertisedName')}>{runtime.server_name}</Row>}
            </>
          )}
        </CardContent>
      </Card>

      <Card>
        <CardHeader className="pb-2">
          <CardTitle className="font-display text-base">{t('connection')}</CardTitle>
        </CardHeader>
        <CardContent>
          {publicHost ? <ConnectionInfo server={server} host={publicHost} /> : <p className="text-sm text-muted-foreground">{t('noPublicHost')}</p>}
          <div className="mt-3 border-t pt-2">
            <Row label={t('gamePort')}>{`${server.port} (+${server.port + 1})`}</Row>
            <Row label={t('queryPort')}>{server.query_port}</Row>
            <Row label={t('rconPort')}>{server.rcon_port}</Row>
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardHeader className="pb-2">
          <CardTitle className="font-display text-base">{t('config')}</CardTitle>
        </CardHeader>
        <CardContent>
          <Row label={t('identifier')}>{server.identifier}</Row>
          <Row label={t('map')}>{mapName(server.map)}</Row>
          <Row label={t('maxPlayers')}>{server.max_players}</Row>
          <Row label={t('cluster')}>{server.cluster_id || t('none')}</Row>
          <Row label={t('autoRestart')}>{server.auto_restart ? t('on') : t('off')}</Row>
          <Row label={t('autoUpdate')}>{server.auto_update ? t('on') : t('off')}</Row>
          <Row label={t('createdAt')}>{server.created_at}</Row>
        </CardContent>
      </Card>

      <ServerModsPanel server={server} />
    </div>
  );
}

export default function ServerDetailPage() {
  const t = useTranslations('detail');
  const params = useParams();
  const router = useRouter();
  const searchParams = useSearchParams();
  const serverId = Number(params.id);
  const servers = useServers();
  const loaded = useServersLoaded();
  const imageStatus = useImageStatus();
  const publicHost = usePublicHost();
  const server = servers.find((s) => s.id === serverId);

  const initialTab = searchParams.get('tab');
  const [tab, setTab] = useState<Tab>(TABS.includes(initialTab as Tab) ? (initialTab as Tab) : 'overview');

  const { fetchServers, getImageStatus, fetchPublicHost } = serversActions;
  // 详情页用列表接口刷新（单个服务器接口会读取配置文件，开销大）
  const refresh = useCallback(() => fetchServers(), [fetchServers]);

  useEffect(() => {
    fetchServers().catch(() => {});
    getImageStatus().catch(() => {});
    fetchPublicHost();
  }, [fetchServers, getImageStatus, fetchPublicHost]);

  const fast = server ? isTransitioning(server.status) || (server.status === 'running' && server.runtime?.phase !== 'ready') : false;
  usePolling(refresh, fast ? 3000 : 10000);

  const changeTab = (value: string) => {
    setTab(value as Tab);
    const url = new URL(window.location.href);
    if (value === 'overview') url.searchParams.delete('tab');
    else url.searchParams.set('tab', value);
    window.history.replaceState(null, '', url.toString());
  };

  if (!loaded) return <PageLoader />;

  if (!server) {
    return (
      <div className="py-12 text-center">
        <p className="mb-4 text-muted-foreground">{t('notFound')}</p>
        <Button onClick={() => router.push('/servers')}>
          <ArrowLeft className="mr-2 h-4 w-4" />
          {t('backToList')}
        </Button>
      </div>
    );
  }

  return (
    <div className="w-full">
      <Link href="/servers" className="mb-3 inline-flex items-center text-sm text-muted-foreground hover:text-foreground">
        <ArrowLeft className="mr-1 h-4 w-4" />
        {t('backToList')}
      </Link>

      {/* 顶部：地图色带 + 名称、状态、操作 */}
      <div className="hud relative mb-6 rounded-sm border bg-card backdrop-blur-sm">
        <MapBand server={server} className="h-10 text-sm">
          {server.status === 'running' && server.runtime?.phase === 'ready' && (
            <span className="font-display flex items-center gap-1.5 text-base font-semibold text-primary [text-shadow:0_0_10px_rgb(57_195_236/0.6)]">
              <Users className="h-4 w-4 opacity-80" />
              {server.runtime.players}
              <span className="text-sm font-normal opacity-75">/{server.runtime.max_players || server.max_players}</span>
            </span>
          )}
        </MapBand>
        <div className="flex flex-col gap-4 p-4 sm:flex-row sm:items-center sm:justify-between sm:p-5">
          <div className="min-w-0">
            <h1 className="font-display truncate text-2xl font-semibold lg:text-[1.75rem]">{server.session_name}</h1>
            <div className="mt-2 flex flex-wrap items-center gap-1.5">
              <StatusBadge server={server} />
              {server.image_outdated && (
                <span className="rounded-full border border-amber-200 bg-amber-50 px-2 py-0.5 text-xs text-amber-800">{t('imageOutdated')}</span>
              )}
              {server.config_pending && (
                <span className="rounded-full border border-blue-200 bg-blue-50 px-2 py-0.5 text-xs text-blue-800">{t('configPending')}</span>
              )}
              {!!server.mods_outdated?.length && (
                <span className="rounded-full border border-amber-200 bg-amber-50 px-2 py-0.5 text-xs text-amber-800">
                  {t('modsOutdated', { count: server.mods_outdated.length })}
                </span>
              )}
            </div>
          </div>
          <ServerActions
            server={server}
            canStart={imageStatus?.can_start_server ?? false}
            variant="detail"
            onDeleted={() => router.push('/servers')}
          />
        </div>
        {server.last_error && (
          <div className="flex items-start gap-2 border-t border-red-200 bg-red-50 px-4 py-3 text-sm text-red-800 sm:px-5">
            <AlertTriangle className="mt-0.5 h-4 w-4 flex-shrink-0" />
            <span className="break-all">
              <span className="font-medium">{t('lastError')}:</span>{' '}
              {server.last_error}
              <button type="button" className="ml-2 font-medium underline underline-offset-2" onClick={() => changeTab('logs')}>
                {t('viewLogs')}
              </button>
            </span>
          </div>
        )}
      </div>

      <Tabs value={tab} onValueChange={changeTab}>
        <div className="-mx-4 overflow-x-auto px-4 sm:mx-0 sm:px-0">
          <TabsList variant="line">
            {TABS.map((v) => (
              <TabsTrigger key={v} value={v}>
                {t(`tabs.${v}`)}
              </TabsTrigger>
            ))}
          </TabsList>
        </div>
        <TabsContent value="overview" className="mt-4">
          <Overview server={server} publicHost={publicHost} />
        </TabsContent>
        <TabsContent value="console" className="mt-4">
          <ConsolePanel server={server} active={tab === 'console'} />
        </TabsContent>
        <TabsContent value="logs" className="mt-4">
          <LogViewer serverId={server.id} active={tab === 'logs'} />
        </TabsContent>
        <TabsContent value="backups" className="mt-4">
          <BackupsPanel server={server} active={tab === 'backups'} onRestored={() => fetchServers().catch(() => {})} />
        </TabsContent>
        <TabsContent value="schedules" className="mt-4">
          <SchedulesPanel server={server} active={tab === 'schedules'} />
        </TabsContent>
        <TabsContent value="plugins" className="mt-4">
          <PluginsPanel server={server} active={tab === 'plugins'} />
        </TabsContent>
      </Tabs>
    </div>
  );
}
