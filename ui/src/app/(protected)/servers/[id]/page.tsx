"use client";

import { useCallback, useEffect, useState } from 'react';
import { useTranslations } from 'next-intl';
import { useParams, useRouter, useSearchParams } from 'next/navigation';
import Link from 'next/link';
import { ArrowLeft, AlertTriangle, Users, Cpu, MemoryStick, Loader2 } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { Progress } from '@/components/ui/progress';
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
import { ModList } from '@/components/servers/ModList';
import { usePolling } from '@/hooks/use-polling';
import { formatBytes, formatPercent } from '@/lib/format';

const TABS = ['overview', 'console', 'logs', 'backups', 'schedules'] as const;
type Tab = (typeof TABS)[number];

function Row({ label, children }: { label: string; children: React.ReactNode }) {
  return (
    <div className="flex items-start justify-between gap-4 py-1.5 text-sm">
      <span className="flex-shrink-0 text-gray-600">{label}</span>
      <span className="min-w-0 text-right font-medium break-all">{children}</span>
    </div>
  );
}

function Overview({ server, publicHost }: { server: Server; publicHost: string }) {
  const t = useTranslations('detail');
  const tMaps = useTranslations('servers.edit.maps');
  const tPhase = useTranslations('phase');
  const copying = server.status === 'starting' && server.runtime?.phase === 'copying';
  const runtime = server.status === 'running' || copying ? server.runtime : undefined;
  const mapName = tMaps.has(server.map) ? tMaps(server.map) : server.map;

  return (
    <div className="grid gap-4 lg:grid-cols-2">
      <Card>
        <CardHeader className="pb-2">
          <CardTitle className="text-base">{t('runtime')}</CardTitle>
        </CardHeader>
        <CardContent className="space-y-3">
          {!runtime ? (
            <p className="text-sm text-gray-500">{server.status === 'running' ? t('collecting') : t('notRunning')}</p>
          ) : (
            <>
              <div className="text-sm text-gray-700">{tPhase(`${runtime.phase || 'unknown'}Desc`)}</div>
              {runtime.phase === 'downloading' && <Progress value={runtime.progress ?? 0} className="h-2" />}
              {runtime.phase === 'installing_mods' && runtime.mod_total ? (
                <Progress value={((runtime.mod_current ?? 0) / runtime.mod_total) * 100} className="h-2" />
              ) : null}
              {!copying && (
              <div className="grid grid-cols-3 gap-3">
                <div className="rounded-lg bg-gray-50 p-3 text-center">
                  <Users className="mx-auto mb-1 h-4 w-4 text-green-600" />
                  <div className="text-lg font-semibold">
                    {runtime.phase === 'ready' ? `${runtime.players}/${runtime.max_players || server.max_players}` : '-'}
                  </div>
                  <div className="text-xs text-gray-500">{t('players')}</div>
                </div>
                <div className="rounded-lg bg-gray-50 p-3 text-center">
                  <Cpu className="mx-auto mb-1 h-4 w-4 text-blue-600" />
                  <div className="text-lg font-semibold">{formatPercent(runtime.cpu_percent)}</div>
                  <div className="text-xs text-gray-500">CPU</div>
                </div>
                <div className="rounded-lg bg-gray-50 p-3 text-center">
                  <MemoryStick className="mx-auto mb-1 h-4 w-4 text-purple-600" />
                  <div className="text-lg font-semibold">{formatBytes(runtime.memory_bytes)}</div>
                  <div className="text-xs text-gray-500">{t('memory')}</div>
                </div>
              </div>
              )}
              {runtime.server_name && <Row label={t('advertisedName')}>{runtime.server_name}</Row>}
            </>
          )}
        </CardContent>
      </Card>

      <Card>
        <CardHeader className="pb-2">
          <CardTitle className="text-base">{t('connection')}</CardTitle>
        </CardHeader>
        <CardContent>
          {publicHost ? <ConnectionInfo server={server} host={publicHost} /> : <p className="text-sm text-gray-500">{t('noPublicHost')}</p>}
          <div className="mt-3 border-t pt-2">
            <Row label={t('gamePort')}>{`${server.port} (+${server.port + 1})`}</Row>
            <Row label={t('queryPort')}>{server.query_port}</Row>
            <Row label={t('rconPort')}>{server.rcon_port}</Row>
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardHeader className="pb-2">
          <CardTitle className="text-base">{t('config')}</CardTitle>
        </CardHeader>
        <CardContent>
          <Row label={t('identifier')}>{server.identifier}</Row>
          <Row label={t('map')}>{mapName}</Row>
          <Row label={t('maxPlayers')}>{server.max_players}</Row>
          <Row label={t('cluster')}>{server.cluster_id || t('none')}</Row>
          <Row label={t('autoRestart')}>{server.auto_restart ? t('on') : t('off')}</Row>
          <Row label={t('autoUpdate')}>{server.auto_update ? t('on') : t('off')}</Row>
          <Row label={t('createdAt')}>{server.created_at}</Row>
        </CardContent>
      </Card>

      <Card>
        <CardHeader className="pb-2">
          <CardTitle className="text-base">{t('mods')}</CardTitle>
        </CardHeader>
        <CardContent>
          {server.game_mod_ids ? <ModList ids={server.game_mod_ids} /> : <p className="text-sm text-gray-500">{t('noMods')}</p>}
        </CardContent>
      </Card>
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

  if (!loaded) {
    return (
      <div className="py-12 text-center">
        <Loader2 className="mx-auto h-8 w-8 animate-spin text-blue-600" />
      </div>
    );
  }

  if (!server) {
    return (
      <div className="py-12 text-center">
        <p className="mb-4 text-gray-600">{t('notFound')}</p>
        <Button onClick={() => router.push('/servers')}>
          <ArrowLeft className="mr-2 h-4 w-4" />
          {t('backToList')}
        </Button>
      </div>
    );
  }

  return (
    <div className="w-full max-w-none py-4 sm:py-8">
      <div className="mb-6 space-y-3">
        <Link href="/servers" className="inline-flex items-center text-sm text-gray-600 hover:text-gray-900">
          <ArrowLeft className="mr-1 h-4 w-4" />
          {t('backToList')}
        </Link>
        <div className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
          <div className="min-w-0">
            <h1 className="truncate text-2xl font-bold text-gray-900 lg:text-3xl">{server.session_name}</h1>
            <div className="mt-2 flex flex-wrap items-center gap-2">
              <StatusBadge server={server} />
              {server.image_outdated && <span className="rounded-md bg-yellow-100 px-2 py-0.5 text-xs text-yellow-800">{t('imageOutdated')}</span>}
              {server.config_pending && <span className="rounded-md bg-blue-100 px-2 py-0.5 text-xs text-blue-800">{t('configPending')}</span>}
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
          <div className="flex items-start gap-2 rounded-lg bg-red-50 p-3 text-sm text-red-700">
            <AlertTriangle className="mt-0.5 h-4 w-4 flex-shrink-0" />
            <span className="break-all">
              <span className="font-medium">{t('lastError')}:</span> {server.last_error}
              <button type="button" className="ml-2 underline" onClick={() => changeTab('logs')}>
                {t('viewLogs')}
              </button>
            </span>
          </div>
        )}
      </div>

      <Tabs value={tab} onValueChange={changeTab}>
        <div className="-mx-4 overflow-x-auto px-4 sm:mx-0 sm:px-0">
          <TabsList>
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
      </Tabs>
    </div>
  );
}
