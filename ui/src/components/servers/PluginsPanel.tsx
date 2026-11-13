"use client";

import { useCallback, useEffect, useRef, useState } from 'react';
import { useTranslations } from 'next-intl';
import axios from 'axios';
import { AlertTriangle, Download, FileCog, Loader2, Puzzle, RefreshCw, Settings2, Trash2, Upload } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { ConfirmDialog } from '@/components/ui/confirm-dialog';
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from '@/components/ui/dialog';
import { Input } from '@/components/ui/input';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Switch } from '@/components/ui/switch';
import { Textarea } from '@/components/ui/textarea';
import { useApiError } from '@/lib/api-error';
import { formatBytes } from '@/lib/format';
import { toast } from '@/lib/toast';
import { serversActions, type Server } from '@/stores/servers';

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
  size?: number;
  published_at?: string;
}

interface ArkApiVersions {
  versions: ArkApiVersion[];
  latest?: string;
  github_error?: string;
}

// 按点分隔的数字比较版本号，与后端 compareVersions 一致
function compareVersions(a: string, b: string): number {
  const pa = a.split('.');
  const pb = b.split('.');
  for (let i = 0; i < Math.max(pa.length, pb.length); i++) {
    const x = pa[i] ?? '';
    const y = pb[i] ?? '';
    const nx = /^\d+$/.test(x) ? Number(x) : NaN;
    const ny = /^\d+$/.test(y) ? Number(y) : NaN;
    if (!Number.isNaN(nx) && !Number.isNaN(ny)) {
      if (nx !== ny) return nx < ny ? -1 : 1;
    } else if (x !== y) {
      return x < y ? -1 : 1;
    }
  }
  return 0;
}

// 正在编辑的配置：name 为 null 表示 ArkApi 本体
interface ConfigTarget {
  name: string | null;
  content: string;
}

function configUrl(serverId: number, name: string | null) {
  return name === null ? `/api/servers/${serverId}/arkapi/config` : `/api/servers/${serverId}/plugins/${encodeURIComponent(name)}/config`;
}

// 详情页"插件"标签：ArkApi 插件的安装、启停、删除和 config.json 编辑
export function PluginsPanel({ server, active }: { server: Server; active: boolean }) {
  const t = useTranslations('plugins');
  const tCommon = useTranslations('common');
  const apiError = useApiError();
  const [data, setData] = useState<PluginsOverview | null>(null);
  const [uploading, setUploading] = useState(false);
  const [pluginName, setPluginName] = useState('');
  const [deleteTarget, setDeleteTarget] = useState<Plugin | null>(null);
  const [config, setConfig] = useState<ConfigTarget | null>(null);
  const [configError, setConfigError] = useState('');
  const [saving, setSaving] = useState(false);
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

  const stopped = server.status === 'stopped';
  const arkApiEnabled = server.arkapi ?? false;

  const load = useCallback(async () => {
    try {
      const res = await axios.get(`/api/servers/${server.id}/plugins`);
      setData(res.data.data);
    } catch (err) {
      toast.error(apiError(err, t('loadFailed')));
    }
  }, [server.id, apiError, t]);

  const loadVersions = useCallback(
    async (refresh = false) => {
      setLoadingVersions(true);
      try {
        const res = await axios.get('/api/arkapi/versions', { params: refresh ? { refresh: true } : undefined });
        setVersions(res.data.data);
      } catch (err) {
        toast.error(apiError(err, t('versionsLoadFailed')));
      } finally {
        setLoadingVersions(false);
      }
    },
    [apiError, t],
  );

  useEffect(() => {
    if (!active) return;
    const timer = setTimeout(() => {
      load();
    }, 0);
    return () => clearTimeout(timer);
  }, [active, load, server.status]);

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
      toast.success(enabled ? t('arkapiEnabledToast') : t('arkapiDisabledToast'));
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
      toast.success(t('arkapiInstalledToast', { version: res.data.data ?? chosen }));
      load();
      // 只在 GitHub 上的版本安装时会下载进版本库
      if (chosenInfo && !chosenInfo.local) loadVersions();
    } catch (err) {
      toast.error(apiError(err, t('arkapiInstallFailed')));
    } finally {
      setInstallingApi(false);
    }
  };

  const uploadArkApi = async (file: File) => {
    setUploadingApi(true);
    const form = new FormData();
    form.append('file', file);
    if (apiVersionInput.trim()) form.append('version', apiVersionInput.trim());
    try {
      const res = await axios.post('/api/arkapi/versions', form);
      const version = res.data.data as string;
      toast.success(t('arkapiUploaded', { version }));
      setApiVersionInput('');
      setSelected(version);
      loadVersions();
    } catch (err) {
      toast.error(apiError(err, t('arkapiUploadFailed')));
    } finally {
      setUploadingApi(false);
      if (apiFileRef.current) apiFileRef.current.value = '';
    }
  };

  const removeVersion = async () => {
    if (!deleteVersion) return;
    try {
      await axios.delete(`/api/arkapi/versions/${encodeURIComponent(deleteVersion)}`);
      toast.success(t('versionDeleted', { version: deleteVersion }));
      if (selected === deleteVersion) setSelected('');
      loadVersions();
    } catch (err) {
      toast.error(apiError(err, t('versionDeleteFailed')));
      throw err;
    }
  };

  const versionTags = (v: ArkApiVersion) =>
    [v.builtin ? t('builtinTag') : v.local ? t('localTag') : t('githubTag'), v.size ? formatBytes(v.size) : ''].filter(Boolean).join(' · ');

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

  const openConfig = async (name: string | null) => {
    try {
      const res = await axios.get(configUrl(server.id, name));
      setConfigError('');
      setConfig({ name, content: res.data.data?.content ?? '' });
    } catch (err) {
      toast.error(apiError(err, t('configLoadFailed')));
    }
  };

  const formatConfig = () => {
    if (!config) return;
    try {
      setConfig({ ...config, content: JSON.stringify(JSON.parse(config.content), null, 2) });
      setConfigError('');
    } catch (e) {
      setConfigError(t('invalidJson', { message: (e as Error).message }));
    }
  };

  const saveConfig = async () => {
    if (!config) return;
    try {
      JSON.parse(config.content);
    } catch (e) {
      setConfigError(t('invalidJson', { message: (e as Error).message }));
      return;
    }
    setSaving(true);
    try {
      await axios.put(configUrl(server.id, config.name), { content: config.content });
      toast.success(server.status === 'stopped' ? t('configSaved') : t('configSavedRunning'));
      setConfig(null);
      load();
    } catch (err) {
      toast.error(apiError(err, t('configSaveFailed')));
    } finally {
      setSaving(false);
    }
  };

  return (
    <div className="space-y-4">
      <Card>
        <CardHeader className="flex flex-row flex-wrap items-center justify-between gap-2 space-y-0 pb-2">
          <CardTitle className="font-display flex flex-wrap items-center gap-2 text-base">
            <Puzzle className="h-4 w-4 text-primary" />
            ArkApi
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
              {t('arkapiEnable')}
              <Switch checked={arkApiEnabled} onCheckedChange={toggleArkApi} disabled={togglingApi} aria-label={t('arkapiEnable')} />
            </label>
            <Button variant="outline" size="sm" onClick={() => openConfig(null)} disabled={!data?.arkapi_installed}>
              <FileCog className="mr-1.5" />
              {t('arkapiConfig')}
            </Button>
          </div>
        </CardHeader>
        <CardContent className="space-y-3 text-sm">
          {!arkApiEnabled && (
            <p className="flex items-start gap-1.5 rounded-sm border border-amber-200 bg-amber-50 px-3 py-2 text-xs text-amber-800">
              <AlertTriangle className="mt-0.5 size-3.5 flex-shrink-0" />
              {t('arkapiDisabled')}
            </p>
          )}
          <p className="text-muted-foreground">{data && !data.arkapi_installed ? t('arkapiNotInstalled') : t('arkapiHint')}</p>

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
              {t('uploadArkApi')}
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
          <p className="text-xs text-muted-foreground">{t('uploadArkApiHint')}</p>
        </CardContent>
      </Card>

      <Card>
        <CardHeader className="pb-2">
          <CardTitle className="font-display text-base">{t('title')}</CardTitle>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="flex flex-col gap-2 sm:flex-row sm:items-center">
            <Input
              value={pluginName}
              onChange={(e) => setPluginName(e.target.value)}
              placeholder={t('namePlaceholder')}
              className="sm:max-w-xs"
              disabled={!stopped || uploading}
            />
            <span title={stopped ? undefined : t('stopFirst')}>
              <Button onClick={() => fileRef.current?.click()} disabled={!stopped || uploading}>
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
          <p className="text-xs text-muted-foreground">{t('uploadHint')}</p>
          {!stopped && (
            <p className="flex items-start gap-1.5 rounded-sm border border-amber-200 bg-amber-50 px-3 py-2 text-xs text-amber-800">
              <AlertTriangle className="mt-0.5 size-3.5 flex-shrink-0" />
              {t('runningHint')}
            </p>
          )}

          {!data ? (
            <p className="text-sm text-muted-foreground">{tCommon('loading')}</p>
          ) : data.plugins.length === 0 ? (
            <p className="text-sm text-muted-foreground">{t('empty')}</p>
          ) : (
            <ul className="divide-y rounded-sm border">
              {data.plugins.map((p) => (
                <li key={p.name} className={`flex flex-col gap-2 px-3 py-2.5 sm:flex-row sm:items-center sm:justify-between ${p.enabled ? '' : 'opacity-60'}`}>
                  <div className="min-w-0">
                    <div className="flex flex-wrap items-center gap-2 text-sm">
                      <span className="font-medium">{p.full_name || p.name}</span>
                      {p.full_name && p.full_name !== p.name && <span className="font-mono text-xs text-muted-foreground">{p.name}</span>}
                      {p.version && <span className="num text-xs text-muted-foreground">v{p.version}</span>}
                      {!p.enabled && <span className="rounded-full border px-1.5 text-xs text-muted-foreground">{t('disabled')}</span>}
                      {!p.has_dll && (
                        <span className="rounded-full border border-red-200 bg-red-50 px-1.5 text-xs text-red-800" title={t('missingDllTip', { name: p.name })}>
                          {t('missingDll')}
                        </span>
                      )}
                    </div>
                    {p.description && <p className="mt-0.5 text-xs text-muted-foreground">{p.description}</p>}
                    <p className="mt-0.5 text-xs text-muted-foreground">
                      {formatBytes(p.size)}
                      {p.min_api_version && ` · ${t('minApi', { version: p.min_api_version })}`}
                      {p.dependencies?.length ? ` · ${t('dependencies', { names: p.dependencies.join(', ') })}` : ''}
                    </p>
                  </div>
                  <div className="flex flex-shrink-0 items-center gap-2">
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
              ))}
            </ul>
          )}
        </CardContent>
      </Card>

      <ConfirmDialog
        open={!!deleteVersion}
        onOpenChange={(open) => !open && setDeleteVersion(null)}
        title={t('deleteVersionTitle')}
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

      <Dialog open={!!config} onOpenChange={(open) => !open && setConfig(null)}>
        <DialogContent className="sm:max-w-3xl">
          <DialogHeader>
            <DialogTitle>{config?.name === null ? t('arkapiConfig') : t('configTitle', { name: config?.name ?? '' })}</DialogTitle>
            <DialogDescription>{server.status === 'stopped' ? t('configDesc') : t('configDescRunning')}</DialogDescription>
          </DialogHeader>
          <Textarea
            value={config?.content ?? ''}
            onChange={(e) => {
              if (config) setConfig({ ...config, content: e.target.value });
              if (configError) setConfigError('');
            }}
            spellCheck={false}
            className="h-[55vh] [field-sizing:fixed] resize-none bg-terminal text-terminal-foreground font-mono text-xs"
            placeholder="{}"
          />
          {configError && <p className="text-xs text-red-600">{configError}</p>}
          <DialogFooter>
            <Button variant="outline" onClick={formatConfig}>
              {t('format')}
            </Button>
            <Button variant="outline" onClick={() => setConfig(null)}>
              {tCommon('cancel')}
            </Button>
            <Button onClick={saveConfig} disabled={saving}>
              {saving && <Loader2 className="mr-1.5 animate-spin" />}
              {tCommon('save')}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  );
}
