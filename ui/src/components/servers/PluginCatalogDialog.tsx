"use client";

import { useMemo, useState } from 'react';
import { useLocale, useTranslations } from 'next-intl';
import { Check, Download, ExternalLink, Loader2, RefreshCw, Search } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Dialog, DialogContent, DialogDescription, DialogHeader, DialogTitle } from '@/components/ui/dialog';
import { Input } from '@/components/ui/input';
import { Tabs, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { formatBytes, formatDateTime } from '@/lib/format';
import {
  availableUpdate,
  catalogEntryFor,
  type CatalogIndex,
  type CatalogPlugin,
  type InstalledPluginRef,
  type PluginCatalog,
} from '@/lib/plugin-catalog';
import { pluginApiName, type Game } from '@/lib/games';

type CategoryFilter = 'all' | 'official' | 'free' | 'paid';

const categoryStyle: Record<string, string> = {
  official: 'border-amber-200 bg-amber-50 text-amber-800',
  free: 'border-green-200 bg-green-50 text-green-800',
  paid: 'border-purple-200 bg-purple-50 text-purple-800',
};

interface Props {
  game: Game;
  open: boolean;
  onOpenChange: (open: boolean) => void;
  catalog: PluginCatalog | null;
  index: CatalogIndex;
  loading: boolean;
  error: string;
  onRefresh: () => void;
  installed: InstalledPluginRef[];
  stopped: boolean;
  installingId: number | null;
  onInstall: (entry: CatalogPlugin) => void;
}

// 插件页的"浏览插件"：arkfgf.com 提供的 ark-server-api.com 插件目录（服务器所属游戏的那一份），能直接下载的一键安装
export function PluginCatalogDialog({ game, open, onOpenChange, catalog, index, loading, error, onRefresh, installed, stopped, installingId, onInstall }: Props) {
  const t = useTranslations('plugins');
  const locale = useLocale();
  const [query, setQuery] = useState('');
  const [category, setCategory] = useState<CategoryFilter>('all');

  const installedNames = useMemo(() => new Set(installed.map((p) => p.name.toLowerCase())), [installed]);

  // 每个资源对应的已安装插件，以及其中有没有可以更新的
  const status = useMemo(() => {
    const map = new Map<number, { installed: boolean; update: string | null }>();
    for (const plugin of installed) {
      const entry = catalogEntryFor(plugin, index);
      if (!entry) continue;
      const prev = map.get(entry.id);
      const update = availableUpdate(plugin, entry);
      map.set(entry.id, { installed: true, update: prev?.update ?? update });
    }
    return map;
  }, [installed, index]);

  const list = useMemo(() => {
    const q = query.trim().toLowerCase();
    return (catalog?.plugins ?? []).filter((p) => {
      if (category !== 'all' && p.category !== category) return false;
      if (!q) return true;
      return [p.title, p.tagline, p.author, ...(p.plugins ?? []).map((i) => i.name)].some((v) => v?.toLowerCase().includes(q));
    });
  }, [catalog, query, category]);

  const counts = useMemo(() => {
    const c: Record<string, number> = { all: 0, official: 0, free: 0, paid: 0 };
    for (const p of catalog?.plugins ?? []) {
      c.all++;
      c[p.category] = (c[p.category] ?? 0) + 1;
    }
    return c;
  }, [catalog]);

  const action = (entry: CatalogPlugin) => {
    const s = status.get(entry.id);
    const busy = installingId === entry.id;
    if (!entry.download_url) {
      return (
        <Button variant="outline" size="sm" asChild>
          <a href={entry.page_url} target="_blank" rel="noopener noreferrer">
            <ExternalLink className="mr-1.5" />
            {entry.category === 'paid' ? t('catalogBuy') : t('catalogOpenPage')}
          </a>
        </Button>
      );
    }
    if (s?.installed && !s.update) {
      return (
        <Button variant="outline" size="sm" disabled>
          <Check className="mr-1.5" />
          {t('catalogInstalled')}
        </Button>
      );
    }
    return (
      <span title={stopped ? undefined : t('stopFirst')} className="inline-flex">
        <Button size="sm" onClick={() => onInstall(entry)} disabled={!stopped || installingId !== null}>
          {busy ? <Loader2 className="mr-1.5 animate-spin" /> : <Download className="mr-1.5" />}
          {s?.update ? t('updateTo', { version: s.update }) : t('catalogInstall')}
        </Button>
      </span>
    );
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-4xl">
        <DialogHeader>
          <DialogTitle>{t('catalogTitle')}</DialogTitle>
          <DialogDescription>
            {t('catalogDesc', { game: game.toUpperCase(), api: pluginApiName(game) })}
            {catalog?.generated_at && ` ${t('catalogUpdatedAt', { time: formatDateTime(catalog.generated_at, locale) })}`}
          </DialogDescription>
        </DialogHeader>

        <div className="flex flex-wrap items-center gap-2">
          <div className="relative min-w-48 flex-1">
            <Search className="pointer-events-none absolute left-2 top-1/2 size-3.5 -translate-y-1/2 text-muted-foreground" />
            <Input value={query} onChange={(e) => setQuery(e.target.value)} placeholder={t('catalogSearch')} className="pl-7" aria-label={t('catalogSearch')} />
          </div>
          <Tabs value={category} onValueChange={(v) => setCategory(v as CategoryFilter)}>
            <TabsList>
              {(['all', 'official', 'free', 'paid'] as const).map((c) => (
                <TabsTrigger key={c} value={c}>
                  {t(`category_${c}`)}
                  <span className="num ml-1 text-xs text-muted-foreground">{counts[c] ?? 0}</span>
                </TabsTrigger>
              ))}
            </TabsList>
          </Tabs>
          <Button variant="ghost" size="icon" className="h-8 w-8" onClick={onRefresh} disabled={loading} title={t('catalogRefresh')} aria-label={t('catalogRefresh')}>
            <RefreshCw className={loading ? 'animate-spin' : ''} />
          </Button>
        </div>

        {!stopped && <p className="text-xs text-amber-800">{t('catalogStopFirst')}</p>}

        <div className="max-h-[60vh] overflow-y-auto rounded-sm border">
          {!catalog ? (
            <p className="p-4 text-sm text-muted-foreground">{loading ? t('catalogLoading') : error || t('catalogUnavailable')}</p>
          ) : list.length === 0 ? (
            <p className="p-4 text-sm text-muted-foreground">{t('searchEmpty')}</p>
          ) : (
            <ul className="divide-y">
              {list.map((entry) => {
                const deps = Array.from(new Set((entry.plugins ?? []).flatMap((p) => p.dependencies ?? []))).filter(
                  (d) => !(entry.plugins ?? []).some((p) => p.name.toLowerCase() === d.toLowerCase()),
                );
                return (
                  <li key={entry.id} className="flex flex-col gap-2 px-3 py-2.5 sm:flex-row sm:items-center sm:justify-between">
                    <div className="min-w-0">
                      <div className="flex flex-wrap items-center gap-2 text-sm">
                        <a href={entry.page_url} target="_blank" rel="noopener noreferrer" className="font-medium hover:underline">
                          {entry.title}
                        </a>
                        {entry.version && <span className="num text-xs text-muted-foreground">v{entry.version}</span>}
                        <span className={`rounded-full border px-1.5 text-xs ${categoryStyle[entry.category] ?? ''}`}>{t(`category_${entry.category}`)}</span>
                        {entry.category !== 'paid' && !entry.download_url && (
                          <span className="rounded-full border px-1.5 text-xs text-muted-foreground" title={t('catalogExternalTip')}>
                            {t('catalogExternal')}
                          </span>
                        )}
                      </div>
                      {entry.tagline && <p className="mt-0.5 text-xs text-muted-foreground">{entry.tagline}</p>}
                      <p className="mt-0.5 text-xs text-muted-foreground">
                        {[
                          entry.author,
                          t('catalogDownloads', { count: entry.downloads }),
                          entry.updated_at ? t('catalogUpdated', { time: formatDateTime(entry.updated_at, locale) }) : '',
                          entry.file_size ? formatBytes(entry.file_size) : '',
                        ]
                          .filter(Boolean)
                          .join(' · ')}
                      </p>
                      {deps.length > 0 && (
                        <p className="mt-1 flex flex-wrap items-center gap-1 text-xs text-muted-foreground">
                          <span>{t('dependsOn')}</span>
                          {deps.map((d) => {
                            const missing = !installedNames.has(d.toLowerCase());
                            return (
                              <span
                                key={d}
                                className={missing ? 'rounded-full border border-amber-200 bg-amber-50 px-1.5 text-amber-800' : 'rounded-full border px-1.5'}
                                title={missing ? t('dependencyMissingTip', { name: d }) : undefined}
                              >
                                {d}
                                {missing ? ` · ${t('dependencyMissing')}` : ''}
                              </span>
                            );
                          })}
                        </p>
                      )}
                    </div>
                    <div className="flex-shrink-0">{action(entry)}</div>
                  </li>
                );
              })}
            </ul>
          )}
        </div>
        <p className="text-xs text-muted-foreground">{t('catalogDisclaimer')}</p>
      </DialogContent>
    </Dialog>
  );
}
