"use client";

import { useEffect, useState } from 'react';
import { useTranslations } from 'next-intl';
import { useRouter, useParams } from 'next/navigation';
import Link from 'next/link';
import { ArrowLeft } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Alert, AlertDescription } from '@/components/ui/alert';
import { PageLoader } from '@/components/ui/page-loader';
import { PageHeader } from '@/components/PageHeader';
import { ServerForm } from '@/components/servers/ServerForm';
import { serversActions, type Server } from '@/stores/servers';
import { useApiError } from '@/lib/api-error';
import { toast } from '@/lib/toast';

export default function ServerEditPage() {
  const tServers = useTranslations('servers');
  const tEdit = useTranslations('servers.edit');
  const tForm = useTranslations('serverForm');
  const tCommon = useTranslations('common');
  const router = useRouter();
  const params = useParams();
  const serverId = Number(params.id);
  const apiError = useApiError();

  const [server, setServer] = useState<Server | null>(null);
  const [loadError, setLoadError] = useState('');
  const [saving, setSaving] = useState(false);

  useEffect(() => {
    if (!serverId) return;
    let cancelled = false;
    serversActions
      .getServer(serverId)
      .then((data) => {
        if (!cancelled) setServer(data);
      })
      .catch((err) => {
        if (!cancelled) setLoadError(apiError(err, tEdit('loadServerInfoFailed')));
      });
    return () => {
      cancelled = true;
    };
  }, [serverId, apiError, tEdit]);

  const backToDetail = () => router.push(`/servers/${serverId}`);

  const handleSubmit = async (form: Partial<Server>, { restart }: { restart: boolean }) => {
    if (!server) return false;
    // 配置文件只在内容有变化时提交，避免每次保存都重写数据卷里的文件。
    // 服务器未停止时后端会先暂存，下次启动时再写入
    const payload: Partial<Server> = { ...form };
    if (payload.game_user_settings === server.game_user_settings) delete payload.game_user_settings;
    if (payload.game_ini === server.game_ini) delete payload.game_ini;
    // 密码留空表示不修改
    if (!payload.admin_password) delete payload.admin_password;

    setSaving(true);
    try {
      const updated = await serversActions.updateServer(serverId, payload);
      setServer({ ...updated, admin_password: undefined });
      if (restart) {
        await serversActions.restartServer(serverId);
        toast.success(tForm('savedAndRestarting'));
        backToDetail();
      } else {
        toast.success(updated.status === 'stopped' ? tServers('serverUpdateSuccess') : tForm('savedPending'));
      }
      return true;
    } catch (err) {
      toast.error(apiError(err, tServers('serverUpdateError')));
      return false;
    } finally {
      setSaving(false);
    }
  };

  if (loadError) {
    return (
      <div className="w-full">
        <Alert variant="destructive">
          <AlertDescription>{loadError}</AlertDescription>
        </Alert>
        <Button onClick={() => router.push('/servers')} className="mt-4">
          <ArrowLeft className="h-4 w-4" />
          {tCommon('back')}
        </Button>
      </div>
    );
  }

  if (!server) {
    return <PageLoader text={tCommon('loading')} />;
  }

  return (
    <div className="w-full">
      <Link href={`/servers/${serverId}`} className="mb-3 inline-flex items-center text-sm text-muted-foreground hover:text-foreground">
        <ArrowLeft className="mr-1 h-4 w-4" />
        {server.session_name}
      </Link>
      <PageHeader title={tEdit('editTitle')} />

      {(server.status !== 'stopped' || server.config_pending) && (
        <Alert className="mb-4 border-amber-300 bg-amber-50 text-amber-900">
          <AlertDescription>{server.config_pending ? tEdit('configPending') : tEdit('configAppliesOnRestart')}</AlertDescription>
        </Alert>
      )}

      <ServerForm
        mode="edit"
        initial={server}
        running={server.status === 'running'}
        saving={saving}
        onSubmit={handleSubmit}
        onCancel={backToDetail}
      />
    </div>
  );
}
