"use client";

import { useCallback, useEffect, useMemo, useRef, useState } from 'react';
import { useLocale, useTranslations } from 'next-intl';
import axios from 'axios';
import {
  AlertTriangle,
  Download,
  FileCog,
  Globe,
  Loader2,
  Puzzle,
  RefreshCw,
  Search,
  Settings2,
  Store,
  Trash2,
  Upload,
} from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Checkbox } from '@/components/ui/checkbox';
import { ConfirmDialog } from '@/components/ui/confirm-dialog';
import { Input } from '@/components/ui/input';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Switch } from '@/components/ui/switch';
import { useApiError } from '@/lib/api-error';
import { formatBytes } from '@/lib/format';
import {
  availableUpdate,
  catalogEntryFor,
  compareVersions,
  indexCatalog,
  type CatalogPlugin,
  type PluginCatalog,
  type PluginSource,
} from '@/lib/plugin-catalog';
import { toast } from '@/lib/toast';
import { cn } from '@/lib/utils';
import { serversActions, type Server } from '@/stores/servers';
import { gameOf, pluginApiName } from '@/lib/games';
import { PluginCatalogDialog } from './PluginCatalogDialog';
import { PluginConfigDialog, configUrl, type ConfigTarget } from './PluginConfigDialog';

interface Plugin {
  name: string;
  enabled: boolean;
  full_name?: string;
  description?: string;
  version?: string;
  min_api_version?: string;
  dependencies?: string[];
  has_dll: boolean;
  has_config: boolean;
  size: number;
  source?: PluginSource;
}

interface PluginsOverview {
  arkapi_installed: boolean;
  arkapi_version?: string;
  plugins: Plugin[];
}

// ArkApi 版本库（所有服务器共用）：内置、已下载或上传的版本，以及 GitHub 上的版本
interface ArkApiVersion {
  version: string;
  builtin?: boolean;
  local: boolean;
  deletable?: boolean;
  on_github?: boolean;
  on_catalog?: boolean;
  size?: number;
  published_at?: string;
}

interface ArkApiVersions {
  // ase（ArkApi）或 asa（AsaApi）
  game?: string;
  name?: string;
  versions: ArkApiVersion[];
  latest?: string;
  github_error?: string;
}

// ISO 时间只显示日期，用于版本列表
function formatDate(value: string | undefined, locale: string): string {
  if (!value) return '';
  const d = new Date(value);
  if (Number.isNaN(d.getTime())) return '';
  return d.toLocaleDateString(locale === 'zh' ? 'zh-CN' : 'en-US');
}

// 详情页"插件"标签：插件框架（ASE 是 ArkApi，ASA 是 AsaApi）和插件的安装、启停、删除和 config.json 编辑。
// 插件目录（ark-server-api.com 的 ASE 插件）只对 ASE 服务器显示
export function PluginsPanel({ server, active }: { server: Server; active: boolean }) {
  const t = useTranslations('plugins');
  const tCommon = useTranslations('common');
  const locale = useLocale();
  const apiError = useApiError();
  const [data, setData] = useState<PluginsOverview | null>(null);
  const [uploading, setUploading] = useState(false);
  const [pluginName, setPluginName] = useState('');
  const [deleteTarget, setDeleteTarget] = useState<Plugin | null>(null);
  const [config, setConfig] = useState<ConfigTarget | null>(null);
  const fileRef = useRef<HTMLInputElement>(null);
  const [versions, setVersions] = useState<ArkApiVersions | null>(null);
  const [loadingVersions, setLoadingVersions] = useState(false);
  const [selected, setSelected] = useState('');
  const [installingApi, setInstallingApi] = useState(false);
  const [uploadingApi, setUploadingApi] = useState(false);
  const [apiVersionInput, setApiVersionInput] = useState('');
  const [deleteVersion, setDeleteVersion] = useState<string | null>(null);
  const [togglingApi, setTogglingApi] = useState(false);
  const apiFileRef = useRef<HTMLInputElement>(null);

  // 插件列表：搜索、多选、批量操作、从地址安装、拖拽上传
  const [query, setQuery] = useState('');
  const [picked, setPicked] = useState<string[]>([]);
  const [bulkBusy, setBulkBusy] = useState(false);
  const [bulkDeleteOpen, setBulkDeleteOpen] = useState(false);
  const [installUrl, setInstallUrl] = useState('');
  const [fetching, setFetching] = useState(false);
  const [dragging, setDragging] = useState(false);
  const dragDepth = useRef(0);

  // 插件目录（arkfgf.com 提供）
  const [catalog, setCatalog] = useState<PluginCatalog | null>(null);
  const [catalogLoading, setCatalogLoading] = useState(false);
  const [catalogError, setCatalogError] = useState('');
  const [catalogOpen, setCatalogOpen] = useState(false);
  const [installingId, setInstallingId] = useState<number | null>(null);

  const stopped = server.status === 'stopped';
  const arkApiEnabled = server.arkapi ?? false;
  const game = gameOf(server);
  const asa = game === 'asa';
  // 插件框架的名称，用在这台服务器的文案里
  const api = pluginApiName(game);

  const load = useCallback(async () => {
    try {
      const res = await axios.get(`/api/servers/${server.id}/plugins`);
      const overview = res.data.data as PluginsOverview;
      setData(overview);
      // 插件可能已经被删除，清掉选不中的
      const names = new Set(overview.plugins.map((p) => p.name));
      setPicked((prev) => prev.filter((name) => names.has(name)));
    } catch (err) {
      toast.error(apiError(err, t('loadFailed')));
    }
  }, [server.id, apiError, t]);

  const loadVersions = useCallback(
    async (refresh = false) => {
      setLoadingVersions(true);
      try {
        const res = await axios.get('/api/arkapi/versions', { params: { ...(asa ? { game: 'asa' } : {}), ...(refresh ? { refresh: true } : {}) } });
        setVersions(res.data.data);
      } catch (err) {
        toast.error(apiError(err, t('versionsLoadFailed', { api })));
      } finally {
        setLoadingVersions(false);
      }
    },
    [apiError, t, asa, api],
  );

  useEffect(() => {
    if (!active) return;
    const timer = setTimeout(() => {
      load();
    }, 0);
    return () => clearTimeout(timer);
  }, [active, load, server.status]);

  const loadCatalog = useCallback(
    async (refresh = false) => {
      setCatalogLoading(true);
      try {
        const res = await axios.get('/api/plugins/catalog', { params: { game, ...(refresh ? { refresh: true } : {}) } });
        setCatalog(res.data.data);
        setCatalogError('');
      } catch (err) {
        // 打开标签时静默失败，在目录弹窗里显示原因
        setCatalogError(apiError(err, t('catalogUnavailable')));
      } finally {
        setCatalogLoading(false);
      }
    },
    [apiError, t, game],
  );

  // 目录用来提示更新和缺失的依赖，标签第一次打开时读取一次。按服务器的游戏读取（ASE 是 ArkApi 插件，ASA 是 AsaApi 插件）
  const catalogRequested = useRef(false);
  useEffect(() => {
    if (!active || catalogRequested.current) return;
    catalogRequested.current = true;
    const timer = setTimeout(() => {
      loadCatalog();
    }, 0);
    return () => clearTimeout(timer);
  }, [active, loadCatalog]);

  // 版本库变化不频繁，标签第一次打开时读取一次
  const versionsRequested = useRef(false);
  useEffect(() => {
    if (!active || versionsRequested.current) return;
    versionsRequested.current = true;
    const timer = setTimeout(() => {
      loadVersions();
    }, 0);
    return () => clearTimeout(timer);
  }, [active, loadVersions]);

  const installedVersion = data?.arkapi_installed ? data.arkapi_version : undefined;
  const latest = versions?.latest;
  const hasUpdate = !!(installedVersion && latest && compareVersions(latest, installedVersion) > 0);
  // 没有手动选择时默认最新版本
  const chosen = selected && versions?.versions.some((v) => v.version === selected) ? selected : (latest ?? '');
  const chosenInfo = versions?.versions.find((v) => v.version === chosen);

  const toggleArkApi = async (enabled: boolean) => {
    setTogglingApi(true);
    try {
      await serversActions.updateServer(server.id, { arkapi: enabled });
      toast.success(enabled ? t('arkapiEnabledToast', { api }) : t('arkapiDisabledToast', { api }));
    } catch (err) {
      toast.error(apiError(err, t('arkapiToggleFailed')));
    } finally {
      setTogglingApi(false);
    }
  };

  const installArkApi = async () => {
    if (!chosen) return;
    setInstallingApi(true);
    try {
      const res = await axios.post(`/api/servers/${server.id}/arkapi/install`, { version: chosen });
      toast.success(t('arkapiInstalledToast', { api, version: res.data.data ?? chosen }));
      load();
      // 只在 GitHub 上的版本安装时会下载进版本库
      if (chosenInfo && !chosenInfo.local) loadVersions();
    } catch (err) {
      toast.error(apiError(err, t('arkapiInstallFailed', { api })));
    } finally {
      setInstallingApi(false);
    }
  };

  const uploadArkApi = async (file: File) => {
    setUploadingApi(true);
    const form = new FormData();
    form.append('file', file);
    if (apiVersionInput.trim()) form.append('version', apiVersionInput.trim());
    if (asa) form.append('game', game);
    try {
      const res = await axios.post('/api/arkapi/versions', form);
      const version = res.data.data as string;
      toast.success(t('arkapiUploaded', { api, version }));
      setApiVersionInput('');
      setSelected(version);
      loadVersions();
    } catch (err) {
      toast.error(apiError(err, t('arkapiUploadFailed', { api })));
    } finally {
      setUploadingApi(false);
      if (apiFileRef.current) apiFileRef.current.value = '';
    }
  };

  const removeVersion = async () => {
    if (!deleteVersion) return;
    try {
      await axios.delete(`/api/arkapi/versions/${encodeURIComponent(deleteVersion)}`, { params: asa ? { game } : undefined });
      toast.success(t('versionDeleted', { version: deleteVersion }));
      if (selected === deleteVersion) setSelected('');
      loadVersions();
    } catch (err) {
      toast.error(apiError(err, t('versionDeleteFailed')));
      throw err;
    }
  };

  const versionTags = (v: ArkApiVersion) =>
    [
      v.builtin ? t('builtinTag') : v.local ? t('localTag') : v.on_github ? t('githubTag') : v.on_catalog ? t('catalogTag') : '',
      v.size ? formatBytes(v.size) : '',
      formatDate(v.published_at, locale),
    ]
      .filter(Boolean)
      .join(' · ');

  const upload = async (file: File) => {
    setUploading(true);
    const form = new FormData();
    form.append('file', file);
    if (pluginName.trim()) form.append('name', pluginName.trim());
    try {
      const res = await axios.post(`/api/servers/${server.id}/plugins`, form);
      toast.success(t('installed', { names: (res.data.data ?? []).join(', ') }));
      setPluginName('');
      load();
    } catch (err) {
      toast.error(apiError(err, t('installFailed')));
    } finally {
      setUploading(false);
      if (fileRef.current) fileRef.current.value = '';
    }
  };

  const installFromUrl = async () => {
    const url = installUrl.trim();
    if (!/^https?:\/\/\S+$/i.test(url)) {
      toast.error(t('urlInvalid'));
      return;
    }
    setFetching(true);
    try {
      const res = await axios.post(`/api/servers/${server.id}/plugins/fetch`, { url, name: pluginName.trim() });
      toast.success(t('installed', { names: (res.data.data ?? []).join(', ') }));
      setInstallUrl('');
      setPluginName('');
      load();
    } catch (err) {
      toast.error(apiError(err, t('fetchFailed')));
    } finally {
      setFetching(false);
    }
  };

  const installCatalog = async (entry: CatalogPlugin) => {
    setInstallingId(entry.id);
    try {
      const res = await axios.post(`/api/servers/${server.id}/plugins/catalog`, { resource_id: entry.id });
      toast.success(t('installed', { names: (res.data.data ?? []).join(', ') }));
      load();
    } catch (err) {
      toast.error(apiError(err, t('installFailed')));
    } finally {
      setInstallingId(null);
    }
  };

  const toggle = async (plugin: Plugin, enabled: boolean) => {
    try {
      await axios.put(`/api/servers/${server.id}/plugins/${encodeURIComponent(plugin.name)}`, { enabled });
      toast.success(enabled ? t('enabledToast', { name: plugin.name }) : t('disabledToast', { name: plugin.name }));
      load();
    } catch (err) {
      toast.error(apiError(err, t('toggleFailed')));
    }
  };

  const remove = async () => {
    if (!deleteTarget) return;
    try {
      await axios.delete(`/api/servers/${server.id}/plugins/${encodeURIComponent(deleteTarget.name)}`);
      toast.success(t('deleted'));
      load();
    } catch (err) {
      toast.error(apiError(err, t('deleteFailed')));
      throw err;
    }
  };

  // ---- 插件列表：搜索、多选、批量 ----

  const plugins = data?.plugins;
  const filtered = useMemo(() => {
    const list = plugins ?? [];
    const q = query.trim().toLowerCase();
    if (!q) return list;
    return list.filter((p) => [p.name, p.full_name, p.description].some((v) => v?.toLowerCase().includes(q)));
  }, [plugins, query]);

  const pickedSet = useMemo(() => new Set(picked), [picked]);
  const allPicked = filtered.length > 0 && filtered.every((p) => pickedSet.has(p.name));

  const installedNames = useMemo(() => new Set((plugins ?? []).map((p) => p.name.toLowerCase())), [plugins]);
  const catalogIndex = useMemo(() => indexCatalog(catalog), [catalog]);
  // 依赖里没装上的（不算插件自己），在列表里标红
  const missingDeps = (plugin: Plugin) =>
    (plugin.dependencies ?? []).filter((d) => {
      const name = d.toLowerCase();
      return name !== plugin.name.toLowerCase() && !installedNames.has(name);
    });

  const togglePick = (name: string, checked: boolean) => {
    setPicked((prev) => (checked ? [...prev, name] : prev.filter((n) => n !== name)));
  };

  const pickAll = (checked: boolean) => {
    const names = filtered.map((p) => p.name);
    setPicked((prev) => (checked ? Array.from(new Set([...prev, ...names])) : prev.filter((n) => !names.includes(n))));
  };

  const bulkSetEnabled = async (enabled: boolean) => {
    const targets = (plugins ?? []).filter((p) => pickedSet.has(p.name) && p.enabled !== enabled);
    if (targets.length === 0) {
      toast.success(t('bulkUpToDate'));
      return;
    }
    setBulkBusy(true);
    let failed = 0;
    for (const plugin of targets) {
      try {
        await axios.put(`/api/servers/${server.id}/plugins/${encodeURIComponent(plugin.name)}`, { enabled });
      } catch {
        failed++;
      }
    }
    setBulkBusy(false);
    if (failed > 0) toast.error(t('bulkPartial', { count: failed }));
    else toast.success(enabled ? t('bulkEnabled', { count: targets.length }) : t('bulkDisabled', { count: targets.length }));
    setPicked([]);
    load();
  };

  const bulkDelete = async () => {
    setBulkBusy(true);
    let failed = 0;
    for (const name of picked) {
      try {
        await axios.delete(`/api/servers/${server.id}/plugins/${encodeURIComponent(name)}`);
      } catch {
        failed++;
      }
    }
    setBulkBusy(false);
    if (failed > 0) toast.error(t('bulkPartial', { count: failed }));
    else toast.success(t('bulkDeleted', { count: picked.length }));
    setPicked([]);
    load();
  };

  // ---- config.json 编辑 ----

  const openConfig = async (name: string | null) => {
    try {
      const res = await axios.get(configUrl(server.id, name));
      const content: string = res.data.data?.content ?? '';
      // 新建配置时从空对象开始，按字段编辑可以直接用
      setConfig({ name, content: content.trim() ? content : '{}' });
    } catch (err) {
      toast.error(apiError(err, t('configLoadFailed')));
    }
  };

  const handleDrop = (e: React.DragEvent) => {
    e.preventDefault();
    dragDepth.current = 0;
    setDragging(false);
    if (!stopped || uploading) return;
    const file = e.dataTransfer.files?.[0];
    if (file) upload(file);
  };

  return (
    <div className="space-y-4">
      <Card>
        <CardHeader className="flex flex-row flex-wrap items-center justify-between gap-2 space-y-0 pb-2">
          <CardTitle className="font-display flex flex-wrap items-center gap-2 text-base">
            <Puzzle className="h-4 w-4 text-primary" />
            {api}
            <span className="text-sm font-normal text-muted-foreground">
              {data?.arkapi_installed ? (installedVersion ? `v${installedVersion}` : t('versionUnknown')) : data ? t('notInstalled') : ''}
            </span>
            {hasUpdate && (
              <span className="rounded-full border border-amber-200 bg-amber-50 px-2 py-0.5 text-xs font-normal text-amber-800">
                {t('updateAvailable', { version: latest ?? '' })}
              </span>
            )}
          </CardTitle>
          <div className="flex items-center gap-3">
            <label className="flex items-center gap-2 text-sm">
              {t('arkapiEnable', { api })}
              <Switch checked={arkApiEnabled} onCheckedChange={toggleArkApi} disabled={togglingApi} aria-label={t('arkapiEnable', { api })} />
            </label>
            <Button variant="outline" size="sm" onClick={() => openConfig(null)} disabled={!data?.arkapi_installed}>
              <FileCog className="mr-1.5" />
              {t('arkapiConfig', { api })}
            </Button>
          </div>
        </CardHeader>
        <CardContent className="space-y-3 text-sm">
          {!arkApiEnabled && (
            <p className="flex items-start gap-1.5 rounded-sm border border-amber-200 bg-amber-50 px-3 py-2 text-xs text-amber-800">
              <AlertTriangle className="mt-0.5 size-3.5 flex-shrink-0" />
              {asa ? t('asaapiDisabled') : t('arkapiDisabled')}
            </p>
          )}
          <p className="text-muted-foreground">{data && !data.arkapi_installed ? t('arkapiNotInstalled', { api }) : t('arkapiHint', { api })}</p>

          <div className="flex flex-col gap-2 sm:flex-row sm:items-center">
            <Select value={chosen} onValueChange={setSelected} disabled={!versions?.versions.length}>
              <SelectTrigger className="sm:w-72">
                <SelectValue placeholder={t('selectVersion')} />
              </SelectTrigger>
              <SelectContent>
                {versions?.versions.map((v) => (
                  <SelectItem key={v.version} value={v.version}>
                    v{v.version}
                    <span className="ml-2 text-xs text-muted-foreground">{versionTags(v)}</span>
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
            <span title={stopped ? undefined : t('stopFirst')}>
              <Button onClick={installArkApi} disabled={!stopped || !chosen || installingApi}>
                {installingApi ? <Loader2 className="mr-1.5 animate-spin" /> : <Download className="mr-1.5" />}
                {chosen && chosen === installedVersion ? t('reinstallArkApi') : hasUpdate && chosen === latest ? t('updateArkApi') : t('installArkApi')}
              </Button>
            </span>
            {chosenInfo?.deletable && (
              <Button variant="outline" size="sm" onClick={() => setDeleteVersion(chosenInfo.version)} aria-label={t('deleteVersion')} title={t('deleteVersion')}>
                <Trash2 className="text-red-600" />
              </Button>
            )}
            <Button
              variant="ghost"
              size="icon"
              className="h-8 w-8"
              onClick={() => loadVersions(true)}
              disabled={loadingVersions}
              title={t('refreshVersions')}
              aria-label={t('refreshVersions')}
            >
              <RefreshCw className={loadingVersions ? 'animate-spin' : ''} />
            </Button>
          </div>
          {versions?.github_error && <p className="text-xs text-muted-foreground">{t('githubUnavailable')}</p>}

          <div className="flex flex-col gap-2 border-t pt-3 sm:flex-row sm:items-center">
            <Input
              value={apiVersionInput}
              onChange={(e) => setApiVersionInput(e.target.value)}
              placeholder={t('arkapiVersionPlaceholder')}
              className="sm:max-w-48"
              disabled={uploadingApi}
            />
            <Button variant="outline" onClick={() => apiFileRef.current?.click()} disabled={uploadingApi}>
              {uploadingApi ? <Loader2 className="mr-1.5 animate-spin" /> : <Upload className="mr-1.5" />}
              {t('uploadArkApi', { api })}
            </Button>
            <input
              ref={apiFileRef}
              type="file"
              accept=".zip,application/zip"
              className="hidden"
              onChange={(e) => {
                const file = e.target.files?.[0];
                if (file) uploadArkApi(file);
              }}
            />
          </div>
          <p className="text-xs text-muted-foreground">{asa ? t('uploadAsaApiHint') : t('uploadArkApiHint')}</p>
        </CardContent>
      </Card>

      <Card>
        <CardHeader className="flex flex-row flex-wrap items-center justify-between gap-2 space-y-0 pb-2">
          <CardTitle className="font-display flex items-center gap-2 text-base">
            {t('title')}
            {plugins && plugins.length > 0 && <span className="num text-sm font-normal text-muted-foreground">{plugins.length}</span>}
          </CardTitle>
          <Button
            variant="outline"
            size="sm"
            onClick={() => {
              setCatalogOpen(true);
              if (!catalog && !catalogLoading) loadCatalog();
            }}
          >
            <Store className="mr-1.5" />
            {t('browseCatalog')}
          </Button>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="flex flex-col gap-2 sm:flex-row sm:items-center">
            <Input
              value={pluginName}
              onChange={(e) => setPluginName(e.target.value)}
              placeholder={t('namePlaceholder')}
              className="sm:max-w-xs"
              disabled={!stopped || uploading || fetching}
            />
            <span title={stopped ? undefined : t('stopFirst')}>
              <Button onClick={() => fileRef.current?.click()} disabled={!stopped || uploading || fetching}>
                {uploading ? <Loader2 className="mr-1.5 animate-spin" /> : <Upload className="mr-1.5" />}
                {t('upload')}
              </Button>
            </span>
            <input
              ref={fileRef}
              type="file"
              accept=".zip,application/zip"
              className="hidden"
              onChange={(e) => {
                const file = e.target.files?.[0];
                if (file) upload(file);
              }}
            />
          </div>

          <div
            onDragEnter={(e) => {
              e.preventDefault();
              if (!stopped || uploading || fetching) return;
              dragDepth.current += 1;
              setDragging(true);
            }}
            onDragOver={(e) => e.preventDefault()}
            onDragLeave={() => {
              dragDepth.current = Math.max(0, dragDepth.current - 1);
              if (dragDepth.current === 0) setDragging(false);
            }}
            onDrop={handleDrop}
            onClick={() => {
              if (stopped && !uploading && !fetching) fileRef.current?.click();
            }}
            className={cn(
              'flex cursor-pointer flex-col items-center gap-1 rounded-sm border border-dashed px-3 py-4 text-center text-xs text-muted-foreground transition-colors',
              dragging ? 'border-primary bg-primary/10 text-foreground' : 'border-border',
              (!stopped || uploading || fetching) && 'cursor-not-allowed opacity-60',
            )}
          >
            <Upload className="size-4" />
            <span>{stopped ? t('dropHint') : t('uploadHint')}</span>
          </div>

          <div className="flex flex-col gap-2 sm:flex-row sm:items-center">
            <Input
              value={installUrl}
              onChange={(e) => setInstallUrl(e.target.value)}
              placeholder={t('urlPlaceholder')}
              className="flex-1 font-mono text-xs"
              disabled={!stopped || fetching}
              onKeyDown={(e) => {
                if (e.key === 'Enter') installFromUrl();
              }}
            />
            <span title={stopped ? undefined : t('stopFirst')}>
              <Button variant="outline" onClick={installFromUrl} disabled={!stopped || fetching || !installUrl.trim()}>
                {fetching ? <Loader2 className="mr-1.5 animate-spin" /> : <Globe className="mr-1.5" />}
                {fetching ? t('fetching') : t('installFromUrl')}
              </Button>
            </span>
          </div>
          <p className="text-xs text-muted-foreground">{t('urlHint')}</p>

          {!stopped && (
            <p className="flex items-start gap-1.5 rounded-sm border border-amber-200 bg-amber-50 px-3 py-2 text-xs text-amber-800">
              <AlertTriangle className="mt-0.5 size-3.5 flex-shrink-0" />
              {t('runningHint')}
            </p>
          )}

          {!data ? (
            <p className="text-sm text-muted-foreground">{tCommon('loading')}</p>
          ) : plugins && plugins.length === 0 ? (
            <p className="text-sm text-muted-foreground">{t('empty')}</p>
          ) : (
            <>
              <div className="flex flex-wrap items-center gap-3">
                <div className="relative min-w-48 flex-1 sm:max-w-xs">
                  <Search className="pointer-events-none absolute left-2 top-1/2 size-3.5 -translate-y-1/2 text-muted-foreground" />
                  <Input
                    value={query}
                    onChange={(e) => setQuery(e.target.value)}
                    placeholder={t('searchPlaceholder')}
                    className="pl-7"
                    aria-label={t('searchPlaceholder')}
                  />
                </div>
                <label className="flex items-center gap-1.5 text-xs text-muted-foreground">
                  <Checkbox checked={allPicked} onChange={(e) => pickAll(e.target.checked)} disabled={filtered.length === 0} aria-label={t('selectAll')} />
                  {t('selectAll')}
                </label>
                <span className="num text-xs text-muted-foreground">
                  {t('count', { shown: filtered.length, total: plugins?.length ?? 0 })}
                </span>
              </div>

              {picked.length > 0 && (
                <div className="flex flex-wrap items-center gap-2 rounded-sm border bg-muted/40 px-3 py-2">
                  <span className="text-xs text-muted-foreground">{t('selectedCount', { count: picked.length })}</span>
                  <span title={stopped ? undefined : t('stopFirst')} className="inline-flex">
                    <Button variant="outline" size="sm" onClick={() => bulkSetEnabled(true)} disabled={!stopped || bulkBusy}>
                      {bulkBusy && <Loader2 className="mr-1.5 animate-spin" />}
                      {t('enableSelected')}
                    </Button>
                  </span>
                  <span title={stopped ? undefined : t('stopFirst')} className="inline-flex">
                    <Button variant="outline" size="sm" onClick={() => bulkSetEnabled(false)} disabled={!stopped || bulkBusy}>
                      {t('disableSelected')}
                    </Button>
                  </span>
                  <span title={stopped ? undefined : t('stopFirst')} className="inline-flex">
                    <Button variant="outline" size="sm" onClick={() => setBulkDeleteOpen(true)} disabled={!stopped || bulkBusy}>
                      <Trash2 className="mr-1.5 size-3.5 text-red-600" />
                      {t('deleteSelected')}
                    </Button>
                  </span>
                  <Button variant="ghost" size="sm" onClick={() => setPicked([])} disabled={bulkBusy}>
                    {t('clearSelection')}
                  </Button>
                </div>
              )}

              {filtered.length === 0 ? (
                <p className="text-sm text-muted-foreground">{t('searchEmpty')}</p>
              ) : (
                <ul className="divide-y rounded-sm border">
                  {filtered.map((p) => {
                    const missing = missingDeps(p);
                    const entry = catalogEntryFor(p, catalogIndex);
                    const update = availableUpdate(p, entry);
                    return (
                      <li
                        key={p.name}
                        className={`flex flex-col gap-2 px-3 py-2.5 sm:flex-row sm:items-center sm:justify-between ${p.enabled ? '' : 'opacity-60'}`}
                      >
                        <div className="flex min-w-0 flex-1 items-start gap-2">
                          <Checkbox
                            className="mt-1"
                            checked={pickedSet.has(p.name)}
                            onChange={(e) => togglePick(p.name, e.target.checked)}
                            aria-label={t('selectPlugin', { name: p.name })}
                          />
                          <div className="min-w-0">
                            <div className="flex flex-wrap items-center gap-2 text-sm">
                              <span className="font-medium">{p.full_name || p.name}</span>
                              {p.full_name && p.full_name !== p.name && <span className="font-mono text-xs text-muted-foreground">{p.name}</span>}
                              {p.version && <span className="num text-xs text-muted-foreground">v{p.version}</span>}
                              {!p.enabled && <span className="rounded-full border px-1.5 text-xs text-muted-foreground">{t('disabled')}</span>}
                              {p.source && (
                                <a
                                  href={entry?.page_url}
                                  target="_blank"
                                  rel="noopener noreferrer"
                                  className="rounded-full border px-1.5 text-xs text-muted-foreground hover:underline"
                                  title={t('fromCatalogTip', { version: p.source.version })}
                                >
                                  {p.source.source}
                                </a>
                              )}
                              {update && (
                                <span className="rounded-full border border-amber-200 bg-amber-50 px-1.5 text-xs text-amber-800">
                                  {t('updateAvailable', { version: update })}
                                </span>
                              )}
                              {!p.has_dll && (
                                <span className="rounded-full border border-red-200 bg-red-50 px-1.5 text-xs text-red-800" title={t('missingDllTip', { api, name: p.name })}>
                                  {t('missingDll')}
                                </span>
                              )}
                            </div>
                            {p.description && <p className="mt-0.5 text-xs text-muted-foreground">{p.description}</p>}
                            <p className="mt-0.5 text-xs text-muted-foreground">
                              {formatBytes(p.size)}
                              {p.min_api_version && ` · ${t('minApi', { api, version: p.min_api_version })}`}
                            </p>
                            {!!p.dependencies?.length && (
                              <p className="mt-1 flex flex-wrap items-center gap-1 text-xs text-muted-foreground">
                                <span>{t('dependsOn')}</span>
                                {p.dependencies.map((dep) => {
                                  const isMissing = missing.some((m) => m.toLowerCase() === dep.toLowerCase());
                                  const depEntry = isMissing ? catalogIndex.byPluginName.get(dep.toLowerCase()) : undefined;
                                  if (depEntry) {
                                    return (
                                      <button
                                        key={dep}
                                        type="button"
                                        onClick={() => installCatalog(depEntry)}
                                        disabled={!stopped || installingId !== null}
                                        className="inline-flex items-center gap-1 rounded-full border border-red-200 bg-red-50 px-1.5 text-red-800 hover:bg-red-100 disabled:cursor-not-allowed disabled:opacity-60"
                                        title={stopped ? t('installDependencyTip', { name: dep, title: depEntry.title }) : t('stopFirst')}
                                      >
                                        {installingId === depEntry.id ? <Loader2 className="size-3 animate-spin" /> : <Download className="size-3" />}
                                        {dep} · {t('installDependency')}
                                      </button>
                                    );
                                  }
                                  return (
                                    <span
                                      key={dep}
                                      className={
                                        isMissing
                                          ? 'rounded-full border border-red-200 bg-red-50 px-1.5 text-red-800'
                                          : 'rounded-full border px-1.5'
                                      }
                                      title={isMissing ? t('dependencyMissingTip', { name: dep }) : undefined}
                                    >
                                      {dep}
                                      {isMissing ? ` · ${t('dependencyMissing')}` : ''}
                                    </span>
                                  );
                                })}
                              </p>
                            )}
                          </div>
                        </div>
                        <div className="flex flex-shrink-0 items-center gap-2">
                          {update && entry && (
                            <span title={stopped ? undefined : t('stopFirst')} className="inline-flex">
                              <Button size="sm" onClick={() => installCatalog(entry)} disabled={!stopped || installingId !== null}>
                                {installingId === entry.id ? <Loader2 className="mr-1.5 animate-spin" /> : <Download className="mr-1.5" />}
                                {t('update')}
                              </Button>
                            </span>
                          )}
                          <span title={stopped ? (p.enabled ? t('disable') : t('enable')) : t('stopFirst')} className="inline-flex">
                            <Switch checked={p.enabled} onCheckedChange={(v) => toggle(p, v)} disabled={!stopped} aria-label={p.enabled ? t('disable') : t('enable')} />
                          </span>
                          <Button variant="outline" size="sm" onClick={() => openConfig(p.name)}>
                            <Settings2 className="mr-1.5" />
                            {p.has_config ? t('config') : t('createConfig')}
                          </Button>
                          <span title={stopped ? undefined : t('stopFirst')}>
                            <Button variant="outline" size="sm" onClick={() => setDeleteTarget(p)} disabled={!stopped} aria-label={t('delete')}>
                              <Trash2 className="text-red-600" />
                            </Button>
                          </span>
                        </div>
                      </li>
                    );
                  })}
                </ul>
              )}
            </>
          )}
        </CardContent>
      </Card>

      <PluginCatalogDialog
        open={catalogOpen}
        onOpenChange={setCatalogOpen}
        game={game}
        catalog={catalog}
        index={catalogIndex}
        loading={catalogLoading}
        error={catalogError}
        onRefresh={() => loadCatalog(true)}
        installed={plugins ?? []}
        stopped={stopped}
        installingId={installingId}
        onInstall={installCatalog}
      />

      <ConfirmDialog
        open={!!deleteVersion}
        onOpenChange={(open) => !open && setDeleteVersion(null)}
        title={t('deleteVersionTitle', { api })}
        description={t('deleteVersionDesc', { version: deleteVersion ?? '' })}
        confirmText={t('delete')}
        destructive
        onConfirm={removeVersion}
      />

      <ConfirmDialog
        open={!!deleteTarget}
        onOpenChange={(open) => !open && setDeleteTarget(null)}
        title={t('deleteTitle')}
        description={t('deleteDesc', { name: deleteTarget?.name ?? '' })}
        confirmText={t('delete')}
        destructive
        onConfirm={remove}
      />

      <ConfirmDialog
        open={bulkDeleteOpen}
        onOpenChange={setBulkDeleteOpen}
        title={t('deleteSelectedTitle')}
        description={t('deleteSelectedDesc', { count: picked.length })}
        confirmText={t('delete')}
        destructive
        onConfirm={bulkDelete}
      />

      {config && (
        <PluginConfigDialog
          server={server}
          target={config}
          onClose={() => setConfig(null)}
          onSaved={() => {
            setConfig(null);
            load();
          }}
        />
      )}
    </div>
  );
}
