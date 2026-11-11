"use client";

import { useCallback, useEffect, useState } from 'react';
import { useTranslations } from 'next-intl';
import axios from 'axios';
import Link from 'next/link';
import { Server, Play, Users, Cpu, MemoryStick, HardDrive, RefreshCw, AlertTriangle, Globe, ArrowRight } from 'lucide-react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { useImageStatus, serversActions, type Server as ServerInfo } from '@/stores/servers';
import { ImageStatus } from '@/components/docker/ImageStatus';
import { ImageUpdateConfirmModal } from '@/components/docker/ImageUpdateConfirmModal';
import { usePolling } from '@/hooks/use-polling';
import { useApiError } from '@/lib/api-error';
import { formatBytes, formatPercent } from '@/lib/format';
import { toast } from '@/lib/toast';

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
  disk_total: number;
  disk_updated_at?: string;
  disk_error?: string;
  recent_errors: { server_id: number; session_name: string; error: string }[];
}

interface UpdateTarget {
  imageName: string;
  servers: ServerInfo[];
}

function Stat({ icon: Icon, label, value, sub, color }: { icon: React.ElementType; label: string; value: string; sub?: string; color: string }) {
  return (
    <Card className="border-0 shadow-sm">
      <CardContent className="flex items-center gap-3 p-4">
        <div className={`flex h-10 w-10 flex-shrink-0 items-center justify-center rounded-full ${color}`}>
          <Icon className="h-5 w-5" />
        </div>
        <div className="min-w-0">
          <div className="text-xs text-gray-500">{label}</div>
          <div className="truncate text-xl font-semibold text-gray-900">{value}</div>
          {sub && <div className="truncate text-xs text-gray-500">{sub}</div>}
        </div>
      </CardContent>
    </Card>
  );
}

export default function HomePage() {
  const t = useTranslations('dashboard');
  const tImages = useTranslations('servers.dockerImages');
  const imageStatus = useImageStatus();
  const apiError = useApiError();
  const { getImageStatus } = serversActions;
  const [overview, setOverview] = useState<Overview | null>(null);
  const [refreshingDisk, setRefreshingDisk] = useState(false);
  const [checkingUpdates, setCheckingUpdates] = useState(false);
  const [updateTarget, setUpdateTarget] = useState<UpdateTarget | null>(null);
  const [submittingUpdate, setSubmittingUpdate] = useState(false);

  const loadOverview = useCallback(async (refreshDisk = false) => {
    const res = await axios.get('/api/system/overview', { params: refreshDisk ? { refresh: true } : undefined });
    setOverview(res.data.data);
  }, []);

  const refreshImageStatus = useCallback(() => getImageStatus().catch(() => {}), [getImageStatus]);

  useEffect(() => {
    let cancelled = false;
    axios
      .get('/api/system/overview')
      .then((res) => {
        if (!cancelled) setOverview(res.data.data);
      })
      .catch(() => {});
    refreshImageStatus();
    return () => {
      cancelled = true;
    };
  }, [refreshImageStatus]);

  // 概览每 15 秒刷新；有镜像在下载时每 2 秒刷新进度（刷新页面后也会继续）
  usePolling(() => loadOverview(), 15000);
  usePolling(refreshImageStatus, 2000, imageStatus?.any_pulling ?? false);

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

  const ov = overview;

  return (
    <div className="w-full max-w-none space-y-6 py-4 sm:py-8">
      <div className="flex flex-wrap items-end justify-between gap-3">
        <div>
          <h1 className="text-2xl font-bold text-gray-900 lg:text-3xl">{t('title')}</h1>
          {ov?.public_host && (
            <p className="mt-1 flex items-center gap-1 text-sm text-gray-500">
              <Globe className="h-4 w-4" />
              {t('publicHost')}: <span className="font-mono">{ov.public_host}</span>
            </p>
          )}
        </div>
        <Button asChild>
          <Link href="/servers">
            {t('manageServers')}
            <ArrowRight className="ml-2 h-4 w-4" />
          </Link>
        </Button>
      </div>

      {/* 统计 */}
      <div className="grid grid-cols-2 gap-3 md:grid-cols-3 xl:grid-cols-6">
        <Stat icon={Server} label={t('servers')} value={String(ov?.total_servers ?? '-')} color="bg-gray-100 text-gray-700" />
        <Stat
          icon={Play}
          label={t('running')}
          value={String(ov?.running_servers ?? '-')}
          sub={ov ? t('readyCount', { count: ov.ready_servers }) : undefined}
          color="bg-green-100 text-green-700"
        />
        <Stat icon={Users} label={t('players')} value={String(ov?.online_players ?? '-')} color="bg-blue-100 text-blue-700" />
        <Stat icon={Cpu} label="CPU" value={ov ? formatPercent(ov.cpu_percent) : '-'} sub={t('cpuHint')} color="bg-orange-100 text-orange-700" />
        <Stat icon={MemoryStick} label={t('memory')} value={ov ? formatBytes(ov.memory_bytes) : '-'} color="bg-purple-100 text-purple-700" />
        <Stat icon={HardDrive} label={t('disk')} value={ov ? formatBytes(ov.disk_total) : '-'} color="bg-amber-100 text-amber-700" />
      </div>

      {/* 最近的错误 */}
      {ov && ov.recent_errors.length > 0 && (
        <Card className="border-red-200">
          <CardHeader className="pb-2">
            <CardTitle className="flex items-center gap-2 text-base text-red-700">
              <AlertTriangle className="h-4 w-4" />
              {t('recentErrors')}
            </CardTitle>
          </CardHeader>
          <CardContent>
            <ul className="space-y-2 text-sm">
              {ov.recent_errors.map((e) => (
                <li key={e.server_id} className="flex flex-col gap-0.5 sm:flex-row sm:gap-2">
                  <Link href={`/servers/${e.server_id}?tab=logs`} className="flex-shrink-0 font-medium text-gray-900 hover:underline">
                    {e.session_name}
                  </Link>
                  <span className="break-all text-red-700">{e.error}</span>
                </li>
              ))}
            </ul>
          </CardContent>
        </Card>
      )}

      <div className="grid gap-6 xl:grid-cols-2">
        {/* 镜像管理 */}
        <Card>
          <CardHeader>
            <CardTitle className="text-lg">{t('imageManagement')}</CardTitle>
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
              <div className="py-8 text-center">
                <div className="mx-auto h-8 w-8 animate-spin rounded-full border-b-2 border-blue-600"></div>
                <p className="mt-2 text-gray-600">{tImages('loadingStatus')}</p>
              </div>
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
        <Card>
          <CardHeader className="flex flex-row items-center justify-between">
            <div>
              <CardTitle className="text-lg">{t('diskUsage')}</CardTitle>
              <p className="mt-1 text-xs text-gray-500">{t('diskHint')}</p>
            </div>
            <Button variant="outline" size="sm" onClick={refreshDisk} disabled={refreshingDisk}>
              <RefreshCw className={`mr-1 h-3.5 w-3.5 ${refreshingDisk ? 'animate-spin' : ''}`} />
              {t('recalculate')}
            </Button>
          </CardHeader>
          <CardContent>
            {ov?.disk_error && <p className="mb-2 text-xs text-red-600">{ov.disk_error}</p>}
            {!ov ? null : ov.disk.length === 0 ? (
              <p className="text-sm text-gray-500">{t('noServers')}</p>
            ) : (
              <div className="overflow-x-auto">
                <table className="w-full text-sm">
                  <thead>
                    <tr className="border-b text-left text-xs text-gray-500">
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
                        <td className="py-2 pr-3 text-right whitespace-nowrap">{formatBytes(d.files)}</td>
                        <td className="py-2 pr-3 text-right whitespace-nowrap">{formatBytes(d.data + d.plugins)}</td>
                        <td className="py-2 pr-3 text-right whitespace-nowrap">{formatBytes(d.backups)}</td>
                        <td className="py-2 text-right font-medium whitespace-nowrap">{formatBytes(d.total)}</td>
                      </tr>
                    ))}
                    {Object.entries(ov.cluster_volumes).map(([name, size]) => (
                      <tr key={name} className="border-b text-gray-600 last:border-0">
                        <td className="py-2 pr-3" colSpan={4}>
                          {t('clusterVolume', { name })}
                        </td>
                        <td className="py-2 text-right whitespace-nowrap">{formatBytes(size)}</td>
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
  );
}
