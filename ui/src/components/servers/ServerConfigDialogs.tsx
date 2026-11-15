"use client";

import { useEffect, useRef, useState } from 'react';
import { useTranslations } from 'next-intl';
import axios from 'axios';
import { FileUp, Loader2, Server as ServerIcon } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Switch } from '@/components/ui/switch';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog';
import { serversActions, useServers, type Server } from '@/stores/servers';
import { useApiError, useErrorText } from '@/lib/api-error';
import { toast } from '@/lib/toast';
import {
  CONFIG_PARTS,
  DEFAULT_CONFIG_PARTS,
  InvalidConfigFileError,
  fetchServerConfig,
  parseServerConfig,
  type ConfigPart,
  type ServerConfig,
} from '@/lib/server-config';
import { cn } from '@/lib/utils';

// ---------- 选择要导入 / 复制的部分 ----------

function ConfigPartsPicker({ value, onChange }: { value: ConfigPart[]; onChange: (parts: ConfigPart[]) => void }) {
  const t = useTranslations('serverConfig');
  return (
    <div className="space-y-2.5 text-sm">
      <div className="font-medium">{t('parts')}</div>
      {CONFIG_PARTS.map((part) => (
        <label key={part} className="flex items-start justify-between gap-4">
          <span>
            <span className="font-medium">{t(`part.${part}`)}</span>
            <span className="block text-xs text-gray-500">{t(`partDesc.${part}`)}</span>
          </span>
          <Switch
            checked={value.includes(part)}
            onCheckedChange={(on) => onChange(on ? CONFIG_PARTS.filter((p) => p === part || value.includes(p)) : value.filter((p) => p !== part))}
          />
        </label>
      ))}
    </div>
  );
}

// 弹窗打开时刷新服务器列表（编辑页、新建页可能还没加载过）
function useServersOnOpen(open: boolean) {
  useEffect(() => {
    if (open) serversActions.fetchServers().catch(() => {});
  }, [open]);
  return useServers();
}

// ---------- 导入到表单 ----------

interface ImportConfigDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  // 编辑页的服务器自己，不出现在来源列表里
  excludeServerId?: number;
  onImport: (config: ServerConfig, parts: ConfigPart[]) => void;
}

/** 从文件或另一台服务器读取配置，选中的部分写进表单，保存后才生效 */
export function ImportConfigDialog({ open, onOpenChange, excludeServerId, onImport }: ImportConfigDialogProps) {
  const t = useTranslations('serverConfig');
  const tCommon = useTranslations('common');
  const apiError = useApiError();
  const servers = useServersOnOpen(open).filter((s) => s.id !== excludeServerId);
  const fileInput = useRef<HTMLInputElement>(null);

  const [source, setSource] = useState<'file' | 'server'>('file');
  const [fileConfig, setFileConfig] = useState<ServerConfig | null>(null);
  const [fileName, setFileName] = useState('');
  const [fileError, setFileError] = useState('');
  const [sourceId, setSourceId] = useState('');
  const [parts, setParts] = useState<ConfigPart[]>(DEFAULT_CONFIG_PARTS);
  const [loading, setLoading] = useState(false);

  // 每次打开时清空上次的选择
  const [prevOpen, setPrevOpen] = useState(open);
  if (open !== prevOpen) {
    setPrevOpen(open);
    if (open) {
      setFileConfig(null);
      setFileName('');
      setFileError('');
      setSourceId('');
      setParts(DEFAULT_CONFIG_PARTS);
    }
  }

  const readFile = async (file: File | undefined) => {
    if (!file) return;
    setFileName(file.name);
    setFileConfig(null);
    setFileError('');
    try {
      setFileConfig(parseServerConfig(await file.text()));
    } catch (err) {
      setFileError(err instanceof InvalidConfigFileError ? t('invalidFile') : t('readFileFailed'));
    }
  };

  const ready = parts.length > 0 && (source === 'file' ? fileConfig !== null : sourceId !== '');

  const confirm = async () => {
    let cfg = fileConfig;
    if (source === 'server') {
      setLoading(true);
      try {
        cfg = await fetchServerConfig(Number(sourceId));
      } catch (err) {
        toast.error(apiError(err, t('loadFailed')));
        return;
      } finally {
        setLoading(false);
      }
    }
    if (!cfg) return;
    onImport(cfg, parts);
    onOpenChange(false);
  };

  return (
    <Dialog open={open} onOpenChange={(next) => !loading && onOpenChange(next)}>
      <DialogContent className="max-w-lg">
        <DialogHeader>
          <DialogTitle>{t('importTitle')}</DialogTitle>
          <DialogDescription>{t('importDesc')}</DialogDescription>
        </DialogHeader>

        <div className="space-y-4">
          <div className="grid grid-cols-2 gap-2">
            {(['file', 'server'] as const).map((s) => (
              <button
                key={s}
                type="button"
                onClick={() => setSource(s)}
                className={cn(
                  'flex items-center justify-center gap-2 rounded-md border px-3 py-2 text-sm transition-colors',
                  source === s ? 'border-primary bg-primary/15 text-primary' : 'text-gray-700 hover:bg-gray-100'
                )}
              >
                {s === 'file' ? <FileUp className="size-4" /> : <ServerIcon className="size-4" />}
                {t(s === 'file' ? 'fromFile' : 'fromServer')}
              </button>
            ))}
          </div>

          {source === 'file' ? (
            <div className="space-y-2 text-sm">
              <input
                ref={fileInput}
                type="file"
                accept=".json,application/json"
                className="hidden"
                onChange={(e) => {
                  readFile(e.target.files?.[0]);
                  e.target.value = '';
                }}
              />
              <Button type="button" variant="outline" className="w-full" onClick={() => fileInput.current?.click()}>
                <FileUp />
                {fileName || t('chooseFile')}
              </Button>
              {fileError && <p className="text-xs text-red-700">{fileError}</p>}
              {fileConfig && (
                <p className="text-xs text-gray-600">
                  {fileConfig.session_name ? t('fileFrom', { name: fileConfig.session_name }) : t('fileLoaded')}
                </p>
              )}
            </div>
          ) : servers.length === 0 ? (
            <p className="text-sm text-gray-500">{t('noOtherServers')}</p>
          ) : (
            <Select value={sourceId} onValueChange={setSourceId}>
              <SelectTrigger className="w-full">
                <SelectValue placeholder={t('chooseServer')} />
              </SelectTrigger>
              <SelectContent>
                {servers.map((s) => (
                  <SelectItem key={s.id} value={String(s.id)}>
                    {s.session_name}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          )}

          <ConfigPartsPicker value={parts} onChange={setParts} />
          <p className="text-xs text-gray-500">{t('keepIdentity')}</p>
        </div>

        <DialogFooter>
          <Button variant="outline" onClick={() => onOpenChange(false)} disabled={loading}>
            {tCommon('cancel')}
          </Button>
          <Button onClick={confirm} disabled={!ready || loading}>
            {loading && <Loader2 className="animate-spin" />}
            {t('importConfirm')}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}

// ---------- 复制到其他服务器 ----------

interface CopyResult {
  server_id: number;
  session_name?: string;
  success: boolean;
  error?: string;
  code?: string;
  config_pending?: boolean;
}

interface CopyConfigDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  server: Server;
}

/** 把一台服务器的设置复制到选中的其他服务器，直接保存 */
export function CopyConfigDialog({ open, onOpenChange, server }: CopyConfigDialogProps) {
  const t = useTranslations('serverConfig');
  const tCommon = useTranslations('common');
  const apiError = useApiError();
  const errorText = useErrorText();
  const targets = useServersOnOpen(open).filter((s) => s.id !== server.id);

  const [selected, setSelected] = useState<number[]>([]);
  const [parts, setParts] = useState<ConfigPart[]>(DEFAULT_CONFIG_PARTS);
  const [submitting, setSubmitting] = useState(false);

  const [prevOpen, setPrevOpen] = useState(open);
  if (open !== prevOpen) {
    setPrevOpen(open);
    if (open) {
      setSelected([]);
      setParts(DEFAULT_CONFIG_PARTS);
    }
  }

  const toggle = (id: number, on: boolean) => setSelected((prev) => (on ? [...prev, id] : prev.filter((x) => x !== id)));
  const allSelected = targets.length > 0 && targets.every((s) => selected.includes(s.id));
  const anyNotStopped = targets.some((s) => selected.includes(s.id) && s.status !== 'stopped');

  const confirm = async () => {
    setSubmitting(true);
    try {
      const res = await axios.post(`/api/servers/${server.id}/config/copy`, { target_ids: selected, parts });
      const results: CopyResult[] = res.data.data ?? [];
      const ok = results.filter((r) => r.success);
      const failed = results.filter((r) => !r.success);
      if (ok.length > 0) toast.success(t('copySuccess', { count: ok.length }));
      for (const r of failed) {
        toast.error(t('copyFailedOne', { name: r.session_name || `#${r.server_id}`, error: errorText(r.error, r.code, t('copyFailed')) }));
      }
      serversActions.fetchServers().catch(() => {});
      if (failed.length === 0) onOpenChange(false);
      else setSelected(failed.map((r) => r.server_id));
    } catch (err) {
      toast.error(apiError(err, t('copyFailed')));
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <Dialog open={open} onOpenChange={(next) => !submitting && onOpenChange(next)}>
      <DialogContent className="max-w-lg">
        <DialogHeader>
          <DialogTitle>{t('copyTitle', { name: server.session_name })}</DialogTitle>
          <DialogDescription>{t('copyDesc')}</DialogDescription>
        </DialogHeader>

        <div className="space-y-4">
          <div className="space-y-2 text-sm">
            <div className="flex items-center justify-between">
              <span className="font-medium">{t('targets')}</span>
              {targets.length > 1 && (
                <button
                  type="button"
                  className="text-xs text-primary hover:underline"
                  onClick={() => setSelected(allSelected ? [] : targets.map((s) => s.id))}
                >
                  {allSelected ? t('selectNone') : t('selectAll')}
                </button>
              )}
            </div>
            {targets.length === 0 ? (
              <p className="text-gray-500">{t('noOtherServers')}</p>
            ) : (
              <ul className="max-h-48 space-y-1 overflow-y-auto">
                {targets.map((s) => (
                  <li key={s.id}>
                    <label className="flex cursor-pointer items-center gap-3 rounded border px-3 py-1.5 hover:bg-gray-50">
                      <input
                        type="checkbox"
                        className="size-4 accent-[var(--primary)]"
                        checked={selected.includes(s.id)}
                        onChange={(e) => toggle(s.id, e.target.checked)}
                      />
                      <span className="flex-1 truncate">{s.session_name}</span>
                      {s.status !== 'stopped' && <span className="text-xs text-amber-800">{t('notStopped')}</span>}
                    </label>
                  </li>
                ))}
              </ul>
            )}
          </div>

          <ConfigPartsPicker value={parts} onChange={setParts} />
          <p className="text-xs text-gray-500">{t('keepIdentity')}</p>
          {anyNotStopped && <p className="text-xs text-amber-800">{t('appliesOnRestart')}</p>}
        </div>

        <DialogFooter>
          <Button variant="outline" onClick={() => onOpenChange(false)} disabled={submitting}>
            {tCommon('cancel')}
          </Button>
          <Button onClick={confirm} disabled={submitting || selected.length === 0 || parts.length === 0}>
            {submitting && <Loader2 className="animate-spin" />}
            {t('copyConfirm', { count: selected.length })}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
