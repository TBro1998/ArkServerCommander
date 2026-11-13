"use client";

import { useCallback, useEffect, useState } from 'react';
import { useTranslations } from 'next-intl';
import axios from 'axios';
import Link from 'next/link';
import { RefreshCw, AlertTriangle, Copy, ChevronRight, Plus, Trash2, Users } from 'lucide-react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { PageLoader } from '@/components/ui/page-loader';
import { ConfirmDialog } from '@/components/ui/confirm-dialog';
import { PageHeader } from '@/components/PageHeader';
import {
  useImageStatus,
  useServers,
  useServersLoaded,
  usePublicHost,
  serversActions,
  isTransitioning,
  type Server as ServerInfo,
} from '@/stores/servers';
import { ImageStatus } from '@/components/docker/ImageStatus';
import { ImageUpdateConfirmModal } from '@/components/docker/ImageUpdateConfirmModal';
import { StatusBadge } from '@/components/servers/StatusBadge';
import { ServerActions } from '@/components/servers/ServerActions';
import { usePolling } from '@/hooks/use-polling';
import { useMapName } from '@/hooks/use-map-name';
import { useApiError } from '@/lib/api-error';
import { copyText, formatBytes, formatPercent } from '@/lib/format';
import { mapColor } from '@/lib/maps';
import { toast } from '@/lib/toast';
import { cn } from '@/lib/utils';

interface DiskUsage {
  server_id: number;
  session_name: string;
  files: number;
  data: number;
  plugins: number;
  backups: number;
  total: number;
}

interface Overview {
  public_host: string;
  total_servers: number;
  running_servers: number;
  ready_servers: number;
  online_players: number;
  cpu_percent: number;
  memory_bytes: number;
  disk: DiskUsage[];
  cluster_volumes: Record<string, number>;
  unused_clusters: { id: string; size: number }[];
  disk_total: number;
  disk_updated_at?: string;
  disk_error?: string;
  recent_errors: { server_id: number; session_name: string; error: string }[];
}

interface UpdateTarget {
  imageName: string;
  servers: ServerInfo[];
}

// 仪表条里的一格：标签在上，读数用 Saira 大号数字
function Readout({ label, value, sub, live }: { label: string; value: string; sub?: string; live?: boolean }) {
  return (
    <div className="min-w-0 bg-[#0d1b28] px-4 py-3.5">
      <div className="hud-caps truncate text-xs text-muted-foreground">{label}</div>
      <div className={cn('font-display mt-0.5 truncate text-[1.75rem] leading-tight font-semibold', live && 'glow-text text-primary')}>{value}</div>
      <div className="h-4 truncate text-xs text-muted-foreground">{sub}</div>
    </div>
  );
}

// 服务器名单的一行：地图色条、名称、状态、人数和启停按钮
function RosterRow({ server, canStart }: { server: ServerInfo; canStart: boolean }) {
  const mapName = useMapName();
  const runtime = server.status === 'running' && server.runtime?.phase === 'ready' ? server.runtime : undefined;
  return (
    <li className="flex items-center gap-3 px-4 py-2.5">
      <span
        aria-hidden
        className={cn('h-9 w-1.5 flex-shrink-0 rounded-full', server.status === 'stopped' && 'saturate-[.3]')}
        style={{ backgroundColor: mapColor(server.map) }}
      />
      <div className="min-w-0 flex-1">
        <Link href={`/servers/${server.id}`} className="block truncate font-medium hover:text-primary hover:underline">
          {server.session_name}
        </Link>
        <div className="hidden truncate text-xs text-muted-foreground sm:block">{mapName(server.map)}</div>
        {/* 窄屏放不下状态列，状态挪到名称下面 */}
        <div className="mt-1 sm:hidden">
          <StatusBadge server={server} />
        </div>
      </div>
      <div className="hidden sm:block">
        <StatusBadge server={server} />
      </div>
      <div className="num hidden w-16 items-center justify-end gap-1 text-sm md:flex">
        <Users className="h-3.5 w-3.5 text-muted-foreground" />
        {runtime ? (
          <span>
            {runtime.players}
            <span className="text-muted-foreground">/{runtime.max_players || server.max_players}</span>
          </span>
        ) : (
          <span className="text-muted-foreground">-</span>
        )}
      </div>
      <ServerActions server={server} canStart={canStart} />
    </li>
  );
}

export default function HomePage() {
  const t = useTranslations('dashboard');
  const tImages = useTranslations('servers.dockerImages');
  const imageStatus = useImageStatus();
  const servers = useServers();
  const serversLoaded = useServersLoaded();
  const publicHost = usePublicHost();
  const apiError = useApiError();
  const { getImageStatus, fetchServers, fetchPublicHost } = serversActions;
  const [overview, setOverview] = useState<Overview | null>(null);
  const [refreshingDisk, setRefreshingDisk] = useState(false);
  const [checkingUpdates, setCheckingUpdates] = useState(false);
  const [updateTarget, setUpdateTarget] = useState<UpdateTarget | null>(null);
  const [submittingUpdate, setSubmittingUpdate] = useState(false);
  const [clusterTarget, setClusterTarget] = useState<{ id: string; size: number } | null>(null);

  const loadOverview = useCallback(async (refreshDisk = false) => {
    const res = await axios.get('/api/system/overview', { params: refreshDisk ? { refresh: true } : undefined });
    setOverview(res.data.data);
  }, []);

  const refreshImageStatus = useCallback(() => getImageStatus().catch(() => {}), [getImageStatus]);
  const refreshServers = useCallback(() => fetchServers(), [fetchServers]);

  useEffect(() => {
    let cancelled = false;
    axios
      .get('/api/system/overview')
      .then((res) => {
        if (!cancelled) setOverview(res.data.data);
      })
      .catch(() => {});
    refreshImageStatus();
    fetchServers().catch(() => {});
    fetchPublicHost();
    return () => {
      cancelled = true;
    };
  }, [refreshImageStatus, fetchServers, fetchPublicHost]);

  // 概览每 15 秒刷新；有镜像在下载时每 2 秒刷新进度（刷新页面后也会继续）
  usePolling(() => loadOverview(), 15000);
  usePolling(refreshImageStatus, 2000, imageStatus?.any_pulling ?? false);
  // 名单里的状态：有服务器在启停或准备时快速刷新
  const fast = servers.some((s) => isTransitioning(s.status) || (s.status === 'running' && s.runtime?.phase !== 'ready'));
  usePolling(refreshServers, fast ? 3000 : 10000);

  const getImageLabel = (imageName: string): string => {
    switch (imageName) {
      case 'tbro98/ase-server:latest':
        return tImages('arkServer');
      case 'alpine:latest':
        return tImages('alpineSystem');
      default:
        return imageName;
    }
  };

  const handleDownloadImage = async (imageName: string) => {
    try {
      await axios.post('/api/images/pull', { image_name: imageName });
      toast.success(t('downloadStarted', { name: getImageLabel(imageName) }));
      await refreshImageStatus();
    } catch (err) {
      toast.error(apiError(err, t('downloadFailed')));
    }
  };

  const handleManualDownload = async () => {
    if (!imageStatus?.images) return;
    // 依次下载所有未就绪的镜像
    const notReady = Object.entries(imageStatus.images).filter(([, s]) => !s.ready && !s.pulling);
    for (const [name] of notReady) {
      await handleDownloadImage(name);
    }
  };

  const handleUpdateImage = async (imageName: string) => {
    try {
      const res = await axios.get('/api/images/affected', { params: { image_name: imageName } });
      setUpdateTarget({ imageName, servers: res.data.data?.servers || [] });
    } catch (err) {
      toast.error(apiError(err, tImages('loadAffectedFailed')));
    }
  };

  const handleConfirmUpdate = async (restartRunning: boolean) => {
    if (!updateTarget) return;
    setSubmittingUpdate(true);
    try {
      await axios.post('/api/images/update', { image_name: updateTarget.imageName, restart_running: restartRunning });
      toast.success(tImages('updateStarted'));
      setUpdateTarget(null);
      await refreshImageStatus();
    } catch (err) {
      toast.error(apiError(err, tImages('updateFailed')));
    } finally {
      setSubmittingUpdate(false);
    }
  };

  const handleCheckUpdates = async () => {
    setCheckingUpdates(true);
    try {
      const res = await axios.get('/api/images/check-updates');
      const results: Record<string, boolean> = res.data.data || {};
      const anyUpdate = Object.values(results).some(Boolean);
      toast.success(anyUpdate ? tImages('updatesFound') : tImages('upToDate'));
      await refreshImageStatus();
    } catch (err) {
      toast.error(apiError(err, tImages('checkUpdatesFailed')));
    } finally {
      setCheckingUpdates(false);
    }
  };

  const refreshDisk = async () => {
    setRefreshingDisk(true);
    try {
      await loadOverview(true);
    } catch (err) {
      toast.error(apiError(err, t('diskFailed')));
    } finally {
      setRefreshingDisk(false);
    }
  };

  const deleteCluster = async () => {
    if (!clusterTarget) return;
    try {
      await axios.delete(`/api/system/clusters/${encodeURIComponent(clusterTarget.id)}`);
      toast.success(t('clusterDeleted', { name: clusterTarget.id }));
      await loadOverview();
    } catch (err) {
      toast.error(apiError(err, t('clusterDeleteFailed')));
      throw err;
    }
  };

  const ov = overview;
  const host = ov?.public_host || publicHost;
  const canStart = imageStatus?.can_start_server ?? false;

  const copyHost = async () => {
    if (host && (await copyText(host))) toast.success(t('hostCopied'));
  };

  return (
    <div className="w-full space-y-6">
      <PageHeader
        title={t('title')}
        description={
          host ? (
            <button type="button" onClick={copyHost} title={t('copyHost')} className="group inline-flex items-center gap-1.5 rounded hover:text-foreground">
              {t('publicHost')}
              <span className="font-mono text-foreground">{host}</span>
              <Copy className="h-3.5 w-3.5 opacity-60 group-hover:opacity-100" />
            </button>
          ) : undefined
        }
        actions={
          <Button asChild variant="outline">
            <Link href="/servers">
              {t('manageServers')}
              <ChevronRight className="h-4 w-4" />
            </Link>
          </Button>
        }
      />

      {/* 仪表条：一行读完全局状态 */}
      <section className="hud relative grid grid-cols-2 gap-px rounded-sm border bg-border sm:grid-cols-3 xl:grid-cols-6">
        <Readout label={t('servers')} value={String(ov?.total_servers ?? '-')} />
        <Readout
          label={t('running')}
          value={String(ov?.running_servers ?? '-')}
          sub={ov ? t('readyCount', { count: ov.ready_servers }) : undefined}
          live={!!ov && ov.running_servers > 0}
        />
        <Readout label={t('players')} value={String(ov?.online_players ?? '-')} live={!!ov && ov.online_players > 0} />
        <Readout label="CPU" value={ov ? formatPercent(ov.cpu_percent, 0) : '-'} sub={t('cpuHint')} />
        <Readout label={t('memory')} value={ov ? formatBytes(ov.memory_bytes) : '-'} />
        <Readout label={t('disk')} value={ov ? formatBytes(ov.disk_total) : '-'} />
      </section>

      {/* 最近的错误 */}
      {ov && ov.recent_errors.length > 0 && (
        <section className="rounded-sm border border-l-4 border-red-200 border-l-red-500 bg-red-50/80 px-4 py-3">
          <h2 className="mb-2 flex items-center gap-2 text-sm font-semibold text-red-800">
            <AlertTriangle className="h-4 w-4" />
            {t('recentErrors')}
          </h2>
          <ul className="space-y-1.5 text-sm">
            {ov.recent_errors.map((e) => (
              <li key={e.server_id} className="flex flex-col gap-0.5 sm:flex-row sm:gap-3">
                <Link href={`/servers/${e.server_id}?tab=logs`} className="flex-shrink-0 font-medium hover:underline">
                  {e.session_name}
                </Link>
                <span className="break-all text-red-800">{e.error}</span>
              </li>
            ))}
          </ul>
        </section>
      )}

      <div className="grid grid-cols-1 items-start gap-6 xl:grid-cols-[minmax(0,3fr)_minmax(0,2fr)]">
        {/* 服务器名单：在首页就能启停 */}
        <Card className="gap-0 py-0">
          <CardHeader className="flex flex-row items-center justify-between gap-3 border-b py-3">
            <CardTitle className="font-display text-lg">{t('servers')}</CardTitle>
            {servers.length > 0 && (
              <Button asChild variant="ghost" size="sm" className="text-muted-foreground">
                <Link href="/servers">
                  {t('viewAll')}
                  <ChevronRight className="h-4 w-4" />
                </Link>
              </Button>
            )}
          </CardHeader>
          {!serversLoaded ? (
            <PageLoader />
          ) : servers.length === 0 ? (
            <div className="px-6 py-12 text-center">
              <p className="mb-4 text-sm text-muted-foreground">{t('noServers')}</p>
              <Button asChild>
                <Link href="/servers/new">
                  <Plus className="h-4 w-4" />
                  {t('createFirst')}
                </Link>
              </Button>
            </div>
          ) : (
            <ul className="divide-y">
              {servers.map((s) => (
                <RosterRow key={s.id} server={s} canStart={canStart} />
              ))}
            </ul>
          )}
        </Card>

        <div className="space-y-6">
          {/* 镜像管理 */}
          <Card className="gap-4">
            <CardHeader>
              <CardTitle className="font-display text-lg">{t('imageManagement')}</CardTitle>
            </CardHeader>
            <CardContent>
              {imageStatus ? (
                <ImageStatus
                  imageStatus={imageStatus}
                  checkingUpdates={checkingUpdates}
                  onRefresh={refreshImageStatus}
                  onManualDownload={handleManualDownload}
                  onCheckUpdates={handleCheckUpdates}
                  onDownloadImage={handleDownloadImage}
                  onUpdateImage={handleUpdateImage}
                />
              ) : (
                <PageLoader text={tImages('loadingStatus')} className="py-8" />
              )}
              {updateTarget && (
                <ImageUpdateConfirmModal
                  key={updateTarget.imageName}
                  open
                  imageName={updateTarget.imageName}
                  imageLabel={getImageLabel(updateTarget.imageName)}
                  affectedServers={updateTarget.servers}
                  submitting={submittingUpdate}
                  onCancel={() => setUpdateTarget(null)}
                  onConfirm={handleConfirmUpdate}
                />
              )}
            </CardContent>
          </Card>

          {/* 磁盘占用 */}
          <Card className="gap-4">
            <CardHeader className="flex flex-row items-start justify-between gap-3">
              <div>
                <CardTitle className="font-display text-lg">{t('diskUsage')}</CardTitle>
                <p className="mt-1 text-xs text-muted-foreground">{t('diskHint')}</p>
              </div>
              <Button variant="outline" size="sm" onClick={refreshDisk} disabled={refreshingDisk}>
                <RefreshCw className={cn('h-3.5 w-3.5', refreshingDisk && 'animate-spin')} />
                {t('recalculate')}
              </Button>
            </CardHeader>
            <CardContent>
              {ov?.disk_error && <p className="mb-2 text-xs text-red-700">{ov.disk_error}</p>}
              {!ov ? null : ov.disk.length === 0 && !ov.unused_clusters?.length ? (
                <p className="text-sm text-muted-foreground">{t('noServers')}</p>
              ) : (
                <div className="overflow-x-auto">
                  <table className="num w-full text-sm">
                    <thead>
                      <tr className="border-b text-left text-xs text-muted-foreground">
                        <th className="py-2 pr-3 font-medium">{t('server')}</th>
                        <th className="py-2 pr-3 text-right font-medium">{t('files')}</th>
                        <th className="py-2 pr-3 text-right font-medium">{t('saves')}</th>
                        <th className="py-2 pr-3 text-right font-medium">{t('backups')}</th>
                        <th className="py-2 text-right font-medium">{t('total')}</th>
                      </tr>
                    </thead>
                    <tbody>
                      {ov.disk.map((d) => (
                        <tr key={d.server_id} className="border-b last:border-0">
                          <td className="max-w-[10rem] truncate py-2 pr-3">
                            <Link href={`/servers/${d.server_id}`} className="hover:underline">
                              {d.session_name}
                            </Link>
                          </td>
                          <td className="py-2 pr-3 text-right whitespace-nowrap text-muted-foreground">{formatBytes(d.files)}</td>
                          <td className="py-2 pr-3 text-right whitespace-nowrap text-muted-foreground">{formatBytes(d.data + d.plugins)}</td>
                          <td className="py-2 pr-3 text-right whitespace-nowrap text-muted-foreground">{formatBytes(d.backups)}</td>
                          <td className="py-2 text-right font-medium whitespace-nowrap">{formatBytes(d.total)}</td>
                        </tr>
                      ))}
                      {Object.entries(ov.cluster_volumes).map(([name, size]) => (
                        <tr key={name} className="border-b text-muted-foreground last:border-0">
                          <td className="py-2 pr-3" colSpan={4}>
                            {t('clusterVolume', { name })}
                          </td>
                          <td className="py-2 text-right whitespace-nowrap">{formatBytes(size)}</td>
                        </tr>
                      ))}
                      {(ov.unused_clusters ?? []).map((c) => (
                        <tr key={c.id} className="border-b text-muted-foreground last:border-0">
                          <td className="py-2 pr-3" colSpan={4}>
                            <span className="inline-flex items-center gap-2">
                              {t('unusedCluster', { name: c.id })}
                              <Button variant="ghost" size="sm" onClick={() => setClusterTarget(c)} aria-label={t('deleteCluster')} title={t('deleteCluster')}>
                                <Trash2 className="text-red-600" />
                              </Button>
                            </span>
                          </td>
                          <td className="py-2 text-right whitespace-nowrap">{formatBytes(c.size)}</td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>
              )}
            </CardContent>
          </Card>
        </div>
      </div>

      <ConfirmDialog
        open={!!clusterTarget}
        onOpenChange={(open) => !open && setClusterTarget(null)}
        title={t('deleteClusterTitle')}
        description={t('deleteClusterDesc', { name: clusterTarget?.id ?? '', size: formatBytes(clusterTarget?.size ?? 0) })}
        confirmText={t('deleteCluster')}
        destructive
        onConfirm={deleteCluster}
      />
    </div>
  );
}
