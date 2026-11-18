"use client";

import { useEffect, useState } from 'react';
import { useTranslations } from 'next-intl';
import { useRouter } from 'next/navigation';
import axios from 'axios';
import { ArrowLeft } from 'lucide-react';
import Link from 'next/link';
import { PageHeader } from '@/components/PageHeader';
import { PageLoader } from '@/components/ui/page-loader';
import { ServerForm } from '@/components/servers/ServerForm';
import { serversActions, type Server } from '@/stores/servers';
import { useApiError } from '@/lib/api-error';
import { toast } from '@/lib/toast';
import { indexIni, getIniValue, replaceIniValues, setIniValue } from '@/lib/ini';
import { MOTD_SECTION } from '@/lib/server-settings';

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
  publish: false,
  link: false,
};

// 默认配置文件里写的是默认名称和人数；创建时换成表单里填的值。
// 游戏以启动参数里的名称和人数为准，这里只是让文件内容保持一致
function applyIdentity(data: Partial<Server>, defaults: Partial<Server>): Partial<Server> {
  let gus = data.game_user_settings;
  if (!gus) return data;
  if (data.session_name) gus = replaceIniValues(gus, 'SessionName', data.session_name);
  if (data.max_players) gus = replaceIniValues(gus, 'MaxPlayers', data.max_players);
  // 欢迎消息没改过时，把里面的默认名称换成新名称
  const motd = getIniValue(indexIni(gus), MOTD_SECTION, 'Message');
  const defaultName = defaults.session_name;
  if (motd && defaultName && data.session_name && motd === getIniValue(indexIni(defaults.game_user_settings ?? ''), MOTD_SECTION, 'Message')) {
    gus = setIniValue(gus, MOTD_SECTION, 'Message', motd.split(defaultName).join(data.session_name));
  }
  return { ...data, game_user_settings: gus };
}

export default function ServerNewPage() {
  const t = useTranslations('servers.edit');
  const tForm = useTranslations('serverForm');
  const tList = useTranslations('detail');
  const tCommon = useTranslations('common');
  const router = useRouter();
  const apiError = useApiError();
  const [initial, setInitial] = useState<Partial<Server>>(DEFAULTS);
  // 默认配置文件和启动参数：设置页以它为基础编辑，显示的值就是创建后实际写入的值
  const [defaults, setDefaults] = useState<Partial<Server> | null>(null);
  const [saving, setSaving] = useState(false);

  useEffect(() => {
    let cancelled = false;
    axios
      .get('/api/system/server-defaults')
      .then((res) => {
        if (cancelled) return;
        const d = res.data.data ?? {};
        const next = { session_name: DEFAULTS.session_name, game_user_settings: d.game_user_settings, game_ini: d.game_ini, server_args: d.server_args };
        setDefaults(next);
        setInitial((prev) => ({ ...prev, ...next }));
      })
      // 取不到时不传配置文件和启动参数，后端使用同样的默认值
      .catch(() => {
        if (!cancelled) setDefaults({});
      });
    return () => {
      cancelled = true;
    };
  }, []);

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
      const created = await serversActions.createServer(applyIdentity(data, defaults ?? {}));
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
      {defaults === null ? (
        <PageLoader text={tCommon('loading')} />
      ) : (
        <ServerForm mode="create" initial={initial} saving={saving} onSubmit={handleSubmit} onCancel={() => router.push('/servers')} />
      )}
    </div>
  );
}
