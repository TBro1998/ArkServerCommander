"use client";

import { useCallback, useEffect, useState } from 'react';
import { useLocale, useTranslations } from 'next-intl';
import axios from 'axios';
import { AlertTriangle, CloudOff, Download, ExternalLink, Loader2, RefreshCw, X } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { ConfirmDialog } from '@/components/ui/confirm-dialog';
import {
  ModIdInput,
  ModThumb,
  formatUnixDate,
  mergeModIds,
  modBrowseUrl,
  modStateOf,
  modUrl,
  normalizeModIds,
  type ModInfo,
} from '@/components/servers/ModList';
import { useApiError } from '@/lib/api-error';
import { toast } from '@/lib/toast';
import { serversActions, isTransitioning, type Server } from '@/stores/servers';
import { gameOf } from '@/lib/games';

interface ServerMod extends ModInfo {
  installed: boolean;
  installed_at?: number;
  has_update: boolean;
}

// 详情页的 Mod 卡片。ASE：安装状态、创意工坊是否有更新，可以添加 / 移除 Mod，只重新下载选中的 Mod；
// ASA：Mod 是 CurseForge 项目 ID，由服务端每次启动时自己下载和更新，只能添加 / 移除
export function ServerModsPanel({ server }: { server: Server }) {
  return gameOf(server) === 'asa' ? <AsaModsPanel server={server} /> : <AseModsPanel server={server} />;
}

// 添加和移除只改服务器记录里的 Mod 列表，下次启动才下载或停止加载
function useSaveModIds(server: Server) {
  const t = useTranslations('mods');
  const apiError = useApiError();
  const [saving, setSaving] = useState(false);
  const save = async (next: string, message: string) => {
    setSaving(true);
    try {
      await serversActions.updateServer(server.id, { game_mod_ids: next });
      toast.success(message);
    } catch (err) {
      toast.error(apiError(err, t('modChangeFailed')));
    } finally {
      setSaving(false);
    }
  };
  return { saving, save };
}

function AsaModsPanel({ server }: { server: Server }) {
  const t = useTranslations('mods');
  const { saving, save } = useSaveModIds(server);
  const ids = normalizeModIds(server.game_mod_ids ?? '', 'asa');
  const list = ids ? ids.split(',') : [];
  const busy = isTransitioning(server.status);

  const addMods = (added: string[]) => save(mergeModIds(ids, added.join(',')), t('modAdded'));
  const removeMod = (id: string) => save(list.filter((x) => x !== id).join(','), t('modRemoved'));

  return (
    <Card className={list.length > 0 ? 'lg:col-span-2' : undefined}>
      <CardHeader className="flex flex-row flex-wrap items-center justify-between gap-2 space-y-0 pb-2">
        <CardTitle className="font-display text-base">{t('title')}</CardTitle>
        <a
          href={modBrowseUrl('asa')}
          target="_blank"
          rel="noreferrer"
          className="inline-flex items-center gap-1 text-xs text-primary hover:underline"
        >
          {t('browseCurseforge')}
          <ExternalLink className="size-3" />
        </a>
      </CardHeader>
      <CardContent className="space-y-2">
        <ModIdInput existing={ids} onAdd={addMods} disabled={saving} game="asa" />
        {list.length === 0 ? (
          <p className="text-sm text-muted-foreground">{t('none')}</p>
        ) : (
          <ul className="space-y-1.5">
            {list.map((id, i) => (
              <li key={id} className="flex items-center justify-between gap-3 rounded-sm border bg-muted/50 px-2 py-1.5 text-xs">
                <span className="flex min-w-0 items-center gap-2">
                  <span className="num w-5 flex-shrink-0 text-right text-gray-400">{i + 1}.</span>
                  <span className="truncate text-gray-800">{t('curseforgeMod')}</span>
                </span>
                <span className="flex flex-shrink-0 items-center gap-1">
                  <a
                    href={modUrl(id, 'asa')}
                    target="_blank"
                    rel="noreferrer"
                    title={t('curseforgeSearch')}
                    className="inline-flex items-center gap-0.5 font-mono text-gray-500 hover:text-blue-600"
                  >
                    {id}
                    <ExternalLink className="h-3 w-3" />
                  </a>
                  <button
                    type="button"
                    onClick={() => removeMod(id)}
                    disabled={saving || busy}
                    aria-label={t('removeMod')}
                    title={t('removeMod')}
                    className="rounded-sm p-0.5 text-gray-500 hover:text-red-600 disabled:opacity-30"
                  >
                    <X className="h-3.5 w-3.5" />
                  </button>
                </span>
              </li>
            ))}
          </ul>
        )}
        <p className="text-xs text-muted-foreground">{t('hintAsa')}</p>
        {server.status === 'running' && <p className="text-xs text-amber-800">{t('changesApplyOnRestart')}</p>}
      </CardContent>
    </Card>
  );
}

function AseModsPanel({ server }: { server: Server }) {
  const t = useTranslations('mods');
  const locale = useLocale();
  const apiError = useApiError();
  const [mods, setMods] = useState<ServerMod[] | null>(null);
  const [selected, setSelected] = useState<Set<string>>(new Set());
  const [refreshing, setRefreshing] = useState(false);
  const [saving, setSaving] = useState(false);
  const [confirmOpen, setConfirmOpen] = useState(false);

  const load = useCallback(
    async (refresh = false) => {
      try {
        const res = await axios.get(`/api/servers/${server.id}/mods`, { params: refresh ? { refresh: true } : undefined });
        const list: ServerMod[] = res.data.data ?? [];
        setMods(list);
        // 默认选中有更新的 Mod
        setSelected(new Set(list.filter((m) => m.has_update).map((m) => m.id)));
      } catch (err) {
        toast.error(apiError(err, t('loadFailed')));
      }
    },
    [server.id, apiError, t],
  );

  // Mod 配置变化或启停完成后重新读取安装状态
  useEffect(() => {
    if (!server.game_mod_ids) return;
    const timer = setTimeout(() => {
      load();
    }, 0);
    return () => clearTimeout(timer);
  }, [server.game_mod_ids, server.status, load]);

  const checkUpdates = async () => {
    setRefreshing(true);
    await load(true);
    setRefreshing(false);
  };

  const toggle = (id: string) => {
    setSelected((prev) => {
      const next = new Set(prev);
      if (next.has(id)) next.delete(id);
      else next.add(id);
      return next;
    });
  };

  // 添加和移除只改服务器记录里的 Mod 列表，下次启动才下载或停止加载
  const saveModIds = async (next: string, message: string) => {
    setSaving(true);
    try {
      await serversActions.updateServer(server.id, { game_mod_ids: next });
      toast.success(message);
    } catch (err) {
      toast.error(apiError(err, t('modChangeFailed')));
    } finally {
      setSaving(false);
    }
  };

  const addMods = (ids: string[]) => {
    saveModIds(mergeModIds(server.game_mod_ids ?? '', ids.join(',')), t('modAdded'));
  };

  const removeMod = (id: string) => {
    const next = normalizeModIds(server.game_mod_ids ?? '')
      .split(',')
      .filter((x) => x && x !== id)
      .join(',');
    saveModIds(next, t('modRemoved'));
  };

  const updateSelected = async () => {
    try {
      await serversActions.updateGame(server.id, false, false, Array.from(selected));
      toast.success(t('updateSent'));
    } catch (err) {
      toast.error(apiError(err, t('updateFailed')));
      throw err;
    }
  };

  if (!server.game_mod_ids) {
    return (
      <Card>
        <CardHeader className="pb-2">
          <CardTitle className="font-display text-base">{t('title')}</CardTitle>
        </CardHeader>
        <CardContent className="space-y-2">
          <ModIdInput existing="" onAdd={addMods} disabled={saving} />
          <p className="text-sm text-muted-foreground">{t('none')}</p>
        </CardContent>
      </Card>
    );
  }

  const busy = isTransitioning(server.status);
  const outdated = mods?.filter((m) => m.has_update).length ?? 0;
  // 运行中改 Mod 列表要等重启才生效，先说明白，免得服主以为刚移除的 Mod 立刻就不加载了
  const running = server.status === 'running';

  return (
    <Card className="lg:col-span-2">
      <CardHeader className="flex flex-row flex-wrap items-center justify-between gap-2 space-y-0 pb-2">
        <CardTitle className="font-display text-base">
          {t('title')}
          {outdated > 0 && (
            <span className="ml-2 rounded-full border border-amber-200 bg-amber-50 px-2 py-0.5 align-middle text-xs font-normal text-amber-800">
              {t('outdatedCount', { count: outdated })}
            </span>
          )}
        </CardTitle>
        <div className="flex flex-wrap items-center gap-2">
          <Button variant="outline" size="sm" onClick={checkUpdates} disabled={refreshing}>
            {refreshing ? <Loader2 className="mr-1.5 animate-spin" /> : <RefreshCw className="mr-1.5" />}
            {t('checkUpdates')}
          </Button>
          <span title={busy ? t('busy') : undefined}>
            <Button size="sm" onClick={() => setConfirmOpen(true)} disabled={selected.size === 0 || busy}>
              <Download className="mr-1.5" />
              {t('updateSelected', { count: selected.size })}
            </Button>
          </span>
        </div>
      </CardHeader>
      <CardContent className="space-y-2">
        <ModIdInput existing={server.game_mod_ids} onAdd={addMods} disabled={saving} />
        {!mods ? (
          <p className="text-sm text-muted-foreground">{t('loading')}</p>
        ) : (
          <ul className="space-y-1.5">
            {mods.map((m, i) => {
              const state = modStateOf(m);
              return (
                <li key={m.id} className="flex items-center justify-between gap-3 rounded-sm border bg-muted/50 px-2 py-1.5 text-xs">
                  <label className="flex min-w-0 flex-1 cursor-pointer items-center gap-2">
                    <input
                      type="checkbox"
                      className="size-3.5 flex-shrink-0 accent-primary"
                      checked={selected.has(m.id)}
                      disabled={!m.installed}
                      onChange={() => toggle(m.id)}
                      aria-label={t('select', { name: m.title || m.id })}
                    />
                    <span className="num w-5 flex-shrink-0 text-right text-gray-400">{i + 1}.</span>
                    <ModThumb url={m.preview_url} />
                    <span className="min-w-0">
                      <span
                        className={`flex items-center gap-1 truncate ${
                          state === 'notFound' ? 'text-red-600' : state === 'ok' ? 'text-gray-800' : 'text-gray-500 italic'
                        }`}
                      >
                        {state === 'notFound' && <AlertTriangle className="h-3 w-3 flex-shrink-0" />}
                        {state === 'unavailable' && <CloudOff className="h-3 w-3 flex-shrink-0" />}
                        <span className="truncate">
                          {state === 'unavailable' ? t('unavailable') : state === 'notFound' ? t('notFound') : m.title || m.id}
                        </span>
                      </span>
                      <span className="flex flex-wrap items-center gap-x-2 text-[11px] text-gray-500">
                        <span>
                          {m.installed
                            ? m.installed_at
                              ? t('installedAt', { date: formatUnixDate(m.installed_at, locale) })
                              : t('installed')
                            : t('notInstalled')}
                        </span>
                        {state === 'ok' && m.time_updated ? <span>{t('updatedAt', { date: formatUnixDate(m.time_updated, locale) })}</span> : null}
                        {state === 'unavailable' ? <span>{t('unavailableHint')}</span> : null}
                        {m.has_update && (
                          <span className="rounded-full border border-amber-200 bg-amber-50 px-1.5 text-amber-800">{t('hasUpdate')}</span>
                        )}
                      </span>
                    </span>
                  </label>
                  <span className="flex flex-shrink-0 items-center gap-1">
                    <a
                      href={`https://steamcommunity.com/sharedfiles/filedetails/?id=${m.id}`}
                      target="_blank"
                      rel="noreferrer"
                      className="inline-flex items-center gap-0.5 font-mono text-gray-500 hover:text-blue-600"
                    >
                      {m.id}
                      <ExternalLink className="h-3 w-3" />
                    </a>
                    <button
                      type="button"
                      onClick={() => removeMod(m.id)}
                      disabled={saving || busy}
                      aria-label={t('removeMod')}
                      title={t('removeMod')}
                      className="rounded-sm p-0.5 text-gray-500 hover:text-red-600 disabled:opacity-30"
                    >
                      <X className="h-3.5 w-3.5" />
                    </button>
                  </span>
                </li>
              );
            })}
          </ul>
        )}
        <p className="text-xs text-muted-foreground">{t('hint')}</p>
        {running && <p className="text-xs text-amber-800">{t('changesApplyOnRestart')}</p>}
      </CardContent>

      <ConfirmDialog
        open={confirmOpen}
        onOpenChange={setConfirmOpen}
        title={t('updateConfirmTitle')}
        description={server.status === 'running' ? t('updateConfirmRunning', { count: selected.size }) : t('updateConfirmStopped', { count: selected.size })}
        confirmText={t('updateConfirm')}
        onConfirm={updateSelected}
      />
    </Card>
  );
}
