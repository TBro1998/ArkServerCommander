"use client";

import { useEffect, useState } from 'react';
import { useServers, useServersActions, useServersIsLoading, useImageStatus } from '@/stores/servers';
import { ServerCard } from '@/components/servers/ServerCard';
import { ServerEditModal } from '@/components/servers/ServerEditModal';
import { Button } from '@/components/ui/button';
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert';
import { ClosableAlert } from '@/components/ui/closable-alert';
import { Plus, Loader2, Server as ServerIcon, AlertCircle } from 'lucide-react';
import { Server } from '@/stores/servers';

export default function ServersPage() {
  const servers = useServers();
  const { fetchServers, getImageStatus, startServer, stopServer, createServer, updateServer, deleteServer, getServer } = useServersActions();
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
    fetchServers().catch((err) => setError('获取服务器列表失败'));
    const statusInterval = setInterval(() => getImageStatus(), 2000);
    return () => clearInterval(statusInterval);
  }, [fetchServers, getImageStatus]);

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
        setError('加载服务器信息失败');
        setShowEditModal(false);
    } finally {
        setIsModalLoading(false);
    }
  };

  const handleDeleteServer = async (server: Server) => {
    if (server.status === 'running') {
      setError('无法删除运行中的服务器');
      return;
    }
    if (confirm(`确定要删除服务器 ${server.name} 吗？`)) {
      try {
        await deleteServer(server.id);
        setSuccess('服务器删除成功');
      } catch {
        setError('删除服务器失败');
      }
    }
  };

  const handleSaveServer = async (data: Partial<Server>) => {
    setIsSaving(true);
    try {
      if (editMode === 'create') {
        await createServer(data);
        setSuccess('服务器创建成功');
      } else if (currentServer?.id) {
        await updateServer(currentServer.id, data);
        setSuccess('服务器更新成功');
      }
      setShowEditModal(false);
    } catch {
      setError(editMode === 'create' ? '创建服务器失败' : '更新服务器失败');
    } finally {
      setIsSaving(false);
    }
  };

  const handleStartServer = (server: Server) => startServer(server.id);
  const handleStopServer = (server: Server) => stopServer(server.id);

  return (
    <div className="w-full max-w-none py-8">
      <div className="mb-6">
        <div className="flex flex-col lg:flex-row justify-between items-start lg:items-center gap-4">
          <div className="flex-1">
            <h1 className="text-2xl lg:text-3xl font-bold text-gray-900 mb-2">服务器管理</h1>
            <p className="text-gray-600 text-sm lg:text-base">在这里管理您的 ARK 服务器</p>
          </div>
          <Button onClick={handleAddServer} disabled={!imageStatus?.can_create_server}>
            <Plus className="mr-2 h-4 w-4" />
            添加服务器
          </Button>
        </div>
        {imageStatus && !imageStatus.can_create_server && (
          <Alert variant="destructive" className="mt-4">
            <AlertCircle className="h-4 w-4" />
            <AlertTitle>{imageStatus.any_pulling ? '镜像下载中' : '镜像未就绪'}</AlertTitle>
            <AlertDescription>
              {imageStatus.any_pulling ? '请等待镜像下载完成。' : '创建服务器需要的基础镜像尚未准备好。'}
            </AlertDescription>
          </Alert>
        )}
      </div>

      {error && <ClosableAlert variant="destructive" className="mb-4" title="错误" onClose={() => setError('')}>{error}</ClosableAlert>}
      {success && <ClosableAlert className="mb-4" title="成功" onClose={() => setSuccess('')}>{success}</ClosableAlert>}

      {isLoading && servers.length === 0 ? (
        <div className="text-center py-12">
          <Loader2 className="w-8 h-8 animate-spin text-blue-600 mx-auto mb-4" />
          <p className="text-gray-600">正在加载服务器...</p>
        </div>
      ) : servers.length === 0 ? (
        <div className="text-center py-16 px-4">
          <div className="mx-auto w-16 h-16 bg-gray-100 rounded-full flex items-center justify-center mb-6">
            <ServerIcon className="w-8 h-8 text-gray-400" />
          </div>
          <h3 className="text-xl font-medium text-gray-900 mb-3">没有找到服务器</h3>
          <p className="text-gray-500 mb-8 max-w-md mx-auto">点击“添加服务器”来创建您的第一个 ARK 服务器。</p>
          <Button onClick={handleAddServer} disabled={!imageStatus?.can_create_server}>
            <Plus className="mr-2 h-4 w-4" />
            添加服务器
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