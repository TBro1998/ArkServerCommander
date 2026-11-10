"use client";

import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import { useServers, serversActions, useServersIsLoading, useImageStatus } from '@/stores/servers';
import { ServerCard } from '@/components/servers/ServerCard';
import { Button } from '@/components/ui/button';
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert';
import { ClosableAlert } from '@/components/ui/closable-alert';
import { Plus, Loader2, Server as ServerIcon, AlertCircle } from 'lucide-react';
import { Server } from '@/stores/servers';
import { useTranslations } from 'next-intl';
import { getErrorMessage } from '@/lib/api-error';

// 有服务器处于启动中或停止中时，按这个间隔刷新列表以获取实际结果
const TRANSITION_POLL_INTERVAL = 3000;

export default function ServersPage() {
  const t = useTranslations('servers');
  const tCommon = useTranslations('common');
  const router = useRouter();
  const servers = useServers();
  const { fetchServers, getImageStatus, startServer, stopServer, deleteServer, recreateServer } = serversActions;
  const isLoading = useServersIsLoading();
  const imageStatus = useImageStatus();

  const [error, setError] = useState('');
  const [success, setSuccess] = useState('');

  useEffect(() => {
    fetchServers().catch(() => setError(t('getServerListFailed')));
    getImageStatus().catch(() => {});
  }, [fetchServers, getImageStatus, t]);

  const hasTransitioning = servers.some((s) => s.status === 'starting' || s.status === 'stopping');

  useEffect(() => {
    if (!hasTransitioning) return;
    const timer = setInterval(() => {
      fetchServers().catch(() => {});
    }, TRANSITION_POLL_INTERVAL);
    return () => clearInterval(timer);
  }, [hasTransitioning, fetchServers]);

  const handleAddServer = () => {
    router.push('/servers/new');
  };

  const handleEditServer = (server: Server) => {
    router.push(`/servers/${server.id}/edit`);
  };

  const handleDeleteServer = async (server: Server) => {
    if (server.status === 'running') {
      setError(t('cannotDeleteRunning'));
      return;
    }
    // 卡片上的弹出框已经确认过一次，这里直接删除
    try {
      await deleteServer(server.id);
      setSuccess(t('serverDeleteSuccess'));
    } catch (err) {
      setError(getErrorMessage(err, t('deleteFailed')));
    }
  };

  const handleStartServer = async (server: Server) => {
    try {
      await startServer(server.id);
    } catch (err) {
      setError(getErrorMessage(err, t('startServerFailed')));
    }
  };

  const handleRecreateServer = async (server: Server) => {
    try {
      await recreateServer(server.id);
      setSuccess(t('recreateStarted'));
    } catch (err) {
      setError(getErrorMessage(err, t('recreateFailed')));
    }
  };

  const handleStopServer = async (server: Server) => {
    try {
      await stopServer(server.id);
    } catch (err) {
      setError(getErrorMessage(err, t('stopServerFailed')));
    }
  };

  return (
    <div className="w-full max-w-none py-8">
      <div className="mb-6">
        <div className="flex justify-between items-center">
          <h1 className="text-2xl lg:text-3xl font-bold text-gray-900">{t('title')}</h1>
          <Button onClick={handleAddServer} disabled={!imageStatus?.can_create_server}>
            <Plus className="mr-2 h-4 w-4" />
            {t('addServer')}
          </Button>
        </div>
        {imageStatus && !imageStatus.can_create_server && (
          <Alert variant="destructive" className="mt-4">
            <AlertCircle className="h-4 w-4" />
            <AlertTitle>{imageStatus.any_pulling ? t('imageDownloading') : t('imageNotReady')}</AlertTitle>
            <AlertDescription>
              {imageStatus.any_pulling ? t('imageDownloadingDesc') : t('imageNotReadyDesc')}
            </AlertDescription>
          </Alert>
        )}
      </div>

      {error && <ClosableAlert variant="destructive" className="mb-4" title={tCommon('error')} onClose={() => setError('')}>{error}</ClosableAlert>}
      {success && <ClosableAlert className="mb-4" title={tCommon('success')} onClose={() => setSuccess('')}>{success}</ClosableAlert>}

      {isLoading && servers.length === 0 ? (
        <div className="text-center py-12">
          <Loader2 className="w-8 h-8 animate-spin text-blue-600 mx-auto mb-4" />
          <p className="text-gray-600">{tCommon('loading')}</p>
        </div>
      ) : servers.length === 0 ? (
        <div className="text-center py-16 px-4">
          <div className="mx-auto w-16 h-16 bg-gray-100 rounded-full flex items-center justify-center mb-6">
            <ServerIcon className="w-8 h-8 text-gray-400" />
          </div>
          <h3 className="text-xl font-medium text-gray-900 mb-3">{t('noServers')}</h3>
          <p className="text-gray-500 mb-8 max-w-md mx-auto">{t('noServersDesc')}</p>
          <Button onClick={handleAddServer} disabled={!imageStatus?.can_create_server}>
            <Plus className="mr-2 h-4 w-4" />
            {t('addServer')}
          </Button>
        </div>
      ) : (
        <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5 2xl:grid-cols-6 gap-4 lg:gap-6">
          {servers.map((server) => (
            <ServerCard
              key={server.id}
              server={server}
              canStartServer={imageStatus?.can_start_server ?? false}
              onStart={handleStartServer}
              onStop={handleStopServer}
              onEdit={handleEditServer}
              onDelete={handleDeleteServer}
              onRecreate={handleRecreateServer}
            />
          ))}
        </div>
      )}


    </div>
  );
}