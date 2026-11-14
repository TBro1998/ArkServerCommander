"use client";

import { useEffect, useState } from 'react';
import { useTranslations } from 'next-intl';
import { useRouter } from 'next/navigation';
import axios from 'axios';
import { ArrowLeft } from 'lucide-react';
import Link from 'next/link';
import { PageHeader } from '@/components/PageHeader';
import { ServerForm } from '@/components/servers/ServerForm';
import { serversActions, type Server } from '@/stores/servers';
import { useApiError } from '@/lib/api-error';
import { toast } from '@/lib/toast';

const DEFAULTS: Partial<Server> = {
  session_name: 'ARK Server',
  cluster_id: '',
  port: 7777,
  query_port: 27015,
  rcon_port: 32330,
  map: 'TheIsland',
  max_players: 70,
  game_mod_ids: '',
  auto_restart: true,
  auto_update: false,
  server_args: {
    query_params: {},
    command_line_args: {},
    custom_args: [],
  },
};

export default function ServerNewPage() {
  const t = useTranslations('servers.edit');
  const tForm = useTranslations('serverForm');
  const tList = useTranslations('detail');
  const router = useRouter();
  const apiError = useApiError();
  const [initial, setInitial] = useState<Partial<Server>>(DEFAULTS);
  const [saving, setSaving] = useState(false);

  // 默认端口已被其他服务器占用时，换成第一组空闲端口
  useEffect(() => {
    let cancelled = false;
    axios
      .get('/api/system/suggest-ports')
      .then((res) => {
        const p = res.data.data;
        if (!cancelled && p) setInitial((prev) => ({ ...prev, port: p.port, query_port: p.query_port, rcon_port: p.rcon_port }));
      })
      .catch(() => {});
    return () => {
      cancelled = true;
    };
  }, []);

  const handleSubmit = async (data: Partial<Server>) => {
    setSaving(true);
    try {
      const created = await serversActions.createServer(data);
      toast.success(tForm('created'));
      router.push(`/servers/${created.id}`);
      return true;
    } catch (err) {
      toast.error(apiError(err, t('serverCreateError')));
      return false;
    } finally {
      setSaving(false);
    }
  };

  return (
    <div className="w-full">
      <Link href="/servers" className="mb-3 inline-flex items-center text-sm text-muted-foreground hover:text-foreground">
        <ArrowLeft className="mr-1 h-4 w-4" />
        {tList('backToList')}
      </Link>
      <PageHeader title={t('createTitle')} description={t('createServerDesc')} />
      <ServerForm mode="create" initial={initial} saving={saving} onSubmit={handleSubmit} onCancel={() => router.push('/servers')} />
    </div>
  );
}
