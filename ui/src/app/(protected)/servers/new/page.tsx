"use client";

import { Suspense, useEffect, useState } from 'react';
import { useTranslations } from 'next-intl';
import { useRouter, useSearchParams } from 'next/navigation';
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
import { GAMES, GAME_NAMES, type Game } from '@/lib/games';
import { defaultMap } from '@/lib/maps';
import { Card, CardContent } from '@/components/ui/card';
import { cn } from '@/lib/utils';

const DEFAULTS: Partial<Server> = {
  game: 'ase',
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

// useSearchParams 要放在 Suspense 里，否则构建时整个页面不能预渲染
export default function ServerNewPage() {
  return (
    <Suspense fallback={<PageLoader />}>
      <ServerNew />
    </Suspense>
  );
}

function ServerNew() {
  const t = useTranslations('servers.edit');
  const tForm = useTranslations('serverForm');
  const tList = useTranslations('detail');
  const tCommon = useTranslations('common');
  const router = useRouter();
  const apiError = useApiError();
  // 服务器列表的分区里"新建 ASA 服务器"带 ?game=asa
  const initialGame: Game = useSearchParams().get('game') === 'asa' ? 'asa' : 'ase';
  const [initial, setInitial] = useState<Partial<Server>>(() =>
    initialGame === 'ase' ? DEFAULTS : { ...DEFAULTS, game: initialGame, map: defaultMap(initialGame) },
  );
  // 默认配置文件和启动参数：设置页以它为基础编辑，显示的值就是创建后实际写入的值
  const [defaults, setDefaults] = useState<Partial<Server> | null>(null);
  const [saving, setSaving] = useState(false);
  // 游戏：创建后不能修改。切换时按这个游戏的默认配置重新开始填写
  const [game, setGame] = useState<Game>(initialGame);
  const [dirty, setDirty] = useState(false);

  const switchGame = (next: Game) => {
    if (next === game) return;
    if (dirty && !window.confirm(t('switchGameConfirm'))) return;
    setGame(next);
    setDirty(false);
    // 另一个游戏的配置文件和启动参数不通用，等新的默认值加载完再显示表单
    setDefaults(null);
    setInitial((prev) => ({
      ...prev,
      game: next,
      map: defaultMap(next),
      game_mod_ids: '',
      game_user_settings: undefined,
      game_ini: undefined,
      server_args: undefined,
    }));
  };

  useEffect(() => {
    let cancelled = false;
    axios
      .get('/api/system/server-defaults', { params: { game } })
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
  }, [game]);

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
      <Card className="mb-4 py-0">
        <CardContent className="space-y-2 p-4">
          <div className="flex flex-wrap items-baseline justify-between gap-2">
            <span className="text-sm font-medium">{tForm('game')}</span>
            <span className="text-xs text-muted-foreground">{tForm('gameFixed')}</span>
          </div>
          <div className="grid gap-2 sm:grid-cols-2" role="radiogroup" aria-label={tForm('game')}>
            {GAMES.map((g) => (
              <button
                key={g}
                type="button"
                role="radio"
                aria-checked={game === g}
                onClick={() => switchGame(g)}
                disabled={saving}
                className={cn(
                  'rounded-md border px-3 py-2.5 text-left transition-colors',
                  game === g ? 'border-primary bg-primary/15' : 'hover:bg-gray-100',
                )}
              >
                <span className={cn('block text-sm font-semibold', game === g ? 'text-primary' : 'text-gray-800')}>
                  {GAME_NAMES[g]} <span className="num text-xs font-normal text-muted-foreground">{g.toUpperCase()}</span>
                </span>
                <span className="mt-0.5 block text-xs text-gray-500">{tForm(`gameDesc.${g}`)}</span>
              </button>
            ))}
          </div>
        </CardContent>
      </Card>
      {defaults === null ? (
        <PageLoader text={tCommon('loading')} />
      ) : (
        // 切换游戏时重新挂载，表单从这个游戏的默认值开始
        <ServerForm
          key={game}
          mode="create"
          initial={initial}
          saving={saving}
          onSubmit={handleSubmit}
          onCancel={() => router.push('/servers')}
          onDirtyChange={setDirty}
        />
      )}
    </div>
  );
}
