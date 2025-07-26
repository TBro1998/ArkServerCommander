"use client";

import { useEffect, useState } from 'react';
import { useServers, serversActions, useServersIsLoading, useImageStatus } from '@/stores/servers';
import { ServerCard } from '@/components/servers/ServerCard';
import { ServerEditModal } from '@/components/servers/ServerEditModal';
import { Button } from '@/components/ui/button';
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert';
import { ClosableAlert } from '@/components/ui/closable-alert';
import { Plus, Loader2, Server as ServerIcon, AlertCircle } from 'lucide-react';
import { Server } from '@/stores/servers';
import { useTranslations } from 'next-intl';

export default function ServersPage() {
  const t = useTranslations('servers');
  const servers = useServers();
  const { fetchServers, getImageStatus, startServer, stopServer, createServer, updateServer, deleteServer, getServer } = serversActions;
  const isLoading = useServersIsLoading();
  const imageStatus = useImageStatus();

  const [showEditModal, setShowEditModal] = useState(false);
  const [editMode, setEditMode] = useState<'create' | 'edit'>('create');
  const [currentServer, setCurrentServer] = useState<Partial<Server> | null>(null);
  const [isModalLoading, setIsModalLoading] = useState(false);
  const [isSaving, setIsSaving] = useState(false);
  const [error, setError] = useState('');
  const [success, setSuccess] = useState('');

  useEffect(() => {
    fetchServers().catch(() => setError(t('getServerListFailed')));
    getImageStatus();
  }, [fetchServers, getImageStatus, t]);

  const handleAddServer = () => {
    setCurrentServer(null);
    setEditMode('create');
    setShowEditModal(true);
  };

  const handleEditServer = async (server: Server) => {
    setIsModalLoading(true);
    setCurrentServer(null);
    setEditMode('edit');
    setShowEditModal(true);
    try {
      const serverData = await getServer(server.id);
      setCurrentServer(serverData);
    } catch {
      setError(t('loadServerInfoFailed'));
      setShowEditModal(false);
    } finally {
      setIsModalLoading(false);
    }
  };

  const handleDeleteServer = async (server: Server) => {
    if (server.status === 'running') {
      setError(t('cannotDeleteRunning'));
      return;
    }
    if (confirm(t('card.confirmDeleteMessage', { identifier: server.name }))) {
      try {
        await deleteServer(server.id);
        setSuccess(t('serverDeleteSuccess'));
      } catch {
        setError(t('deleteFailed'));
      }
    }
  };

  const handleSaveServer = async (data: Partial<Server>) => {
    setIsSaving(true);
    try {
      if (editMode === 'create') {
        await createServer(data);
        setSuccess(t('serverCreateSuccess'));
      } else if (currentServer?.id) {
        await updateServer(currentServer.id, data);
        setSuccess(t('serverUpdateSuccess'));
      }
      setShowEditModal(false);
    } catch {
      setError(editMode === 'create' ? t('serverCreateError') : t('serverUpdateError'));
    } finally {
      setIsSaving(false);
    }
  };

  const handleStartServer = (server: Server) => startServer(server.id);
  const handleStopServer = (server: Server) => stopServer(server.id);

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

      {error && <ClosableAlert variant="destructive" className="mb-4" title={t('common.error')} onClose={() => setError('')}>{error}</ClosableAlert>}
      {success && <ClosableAlert className="mb-4" title={t('common.success')} onClose={() => setSuccess('')}>{success}</ClosableAlert>}

      {isLoading && servers.length === 0 ? (
        <div className="text-center py-12">
          <Loader2 className="w-8 h-8 animate-spin text-blue-600 mx-auto mb-4" />
          <p className="text-gray-600">{t('common.loading')}</p>
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
            />
          ))}
        </div>
      )}

      <ServerEditModal
        show={showEditModal}
        mode={editMode}
        server={currentServer}
        loading={isModalLoading}
        saving={isSaving}
        onClose={() => setShowEditModal(false)}
        onSave={handleSaveServer}
      />
    </div>
  );
}