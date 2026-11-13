"use client";

import { useCallback, useEffect, useState } from 'react';
import { useLocale, useTranslations } from 'next-intl';
import axios from 'axios';
import { Download, ExternalLink, Loader2, RefreshCw } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { ConfirmDialog } from '@/components/ui/confirm-dialog';
import { ModThumb, formatUnixDate, type ModInfo } from '@/components/servers/ModList';
import { useApiError } from '@/lib/api-error';
import { toast } from '@/lib/toast';
import { serversActions, isTransitioning, type Server } from '@/stores/servers';

interface ServerMod extends ModInfo {
  installed: boolean;
  installed_at?: number;
  has_update: boolean;
}

// 详情页的 Mod 卡片：安装状态、创意工坊是否有更新，可以只重新下载选中的 Mod
export function ServerModsPanel({ server }: { server: Server }) {
  const t = useTranslations('mods');
  const locale = useLocale();
  const apiError = useApiError();
  const [mods, setMods] = useState<ServerMod[] | null>(null);
  const [selected, setSelected] = useState<Set<string>>(new Set());
  const [refreshing, setRefreshing] = useState(false);
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
        <CardContent>
          <p className="text-sm text-muted-foreground">{t('none')}</p>
        </CardContent>
      </Card>
    );
  }

  const busy = isTransitioning(server.status);
  const outdated = mods?.filter((m) => m.has_update).length ?? 0;

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
      <CardContent>
        {!mods ? (
          <p className="text-sm text-muted-foreground">{t('loading')}</p>
        ) : (
          <ul className="space-y-1.5">
            {mods.map((m, i) => (
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
                    <span className={`block truncate ${!m.checked ? 'text-gray-500' : m.found ? 'text-gray-800' : 'text-red-600'}`}>
                      {!m.checked ? t('unavailable') : m.found ? m.title || m.id : t('notFound')}
                    </span>
                    <span className="flex flex-wrap items-center gap-x-2 text-[11px] text-gray-500">
                      <span>
                        {m.installed
                          ? m.installed_at
                            ? t('installedAt', { date: formatUnixDate(m.installed_at, locale) })
                            : t('installed')
                          : t('notInstalled')}
                      </span>
                      {m.found && m.time_updated ? <span>{t('updatedAt', { date: formatUnixDate(m.time_updated, locale) })}</span> : null}
                      {m.has_update && (
                        <span className="rounded-full border border-amber-200 bg-amber-50 px-1.5 text-amber-800">{t('hasUpdate')}</span>
                      )}
                    </span>
                  </span>
                </label>
                <a
                  href={`https://steamcommunity.com/sharedfiles/filedetails/?id=${m.id}`}
                  target="_blank"
                  rel="noreferrer"
                  className="inline-flex flex-shrink-0 items-center gap-0.5 font-mono text-gray-500 hover:text-blue-600"
                >
                  {m.id}
                  <ExternalLink className="h-3 w-3" />
                </a>
              </li>
            ))}
          </ul>
        )}
        <p className="mt-2 text-xs text-muted-foreground">{t('hint')}</p>
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
