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
  CROSS_GAME_PARTS,
  DEFAULT_CONFIG_PARTS,
  InvalidConfigFileError,
  fetchServerConfig,
  readImportFiles,
  type ConfigPart,
  type ImportedConfigFile,
  type ServerConfig,
} from '@/lib/server-config';
import { cn } from '@/lib/utils';
import { GAME_NAMES, gameOf, type Game } from '@/lib/games';
import { GameBadge } from './GameBadge';

// ---------- 选择要导入 / 复制的部分 ----------

// disabled：不能导入的部分（跨游戏时的地图、Mod 和启动参数），显示为关闭
function ConfigPartsPicker({
  value,
  onChange,
  disabled = [],
}: {
  value: ConfigPart[];
  onChange: (parts: ConfigPart[]) => void;
  disabled?: ConfigPart[];
}) {
  const t = useTranslations('serverConfig');
  return (
    <div className="space-y-2.5 text-sm">
      <div className="font-medium">{t('parts')}</div>
      {CONFIG_PARTS.map((part) => {
        const off = disabled.includes(part);
        return (
          <label key={part} className={cn('flex items-start justify-between gap-4', off && 'opacity-50')}>
            <span>
              <span className="font-medium">{t(`part.${part}`)}</span>
              <span className="block text-xs text-gray-500">{t(`partDesc.${part}`)}</span>
            </span>
            <Switch
              checked={!off && value.includes(part)}
              disabled={off}
              onCheckedChange={(on) => onChange(on ? CONFIG_PARTS.filter((p) => p === part || value.includes(p)) : value.filter((p) => p !== part))}
            />
          </label>
        );
      })}
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
  // 导入到哪个游戏的服务器：来源是另一个游戏时不导入地图、Mod 和启动参数
  targetGame?: Game;
  // form：填进表单，保存后才生效（新建、编辑页）；apply：直接保存到服务器（详情页的管理标签）
  mode?: 'form' | 'apply';
  // 返回 Promise 时等待完成，抛错时保持弹窗打开（由调用方提示错误）
  onImport: (config: ServerConfig, parts: ConfigPart[]) => void | Promise<void>;
}

/** 从文件或另一台服务器读取配置，把选中的部分填进表单或直接保存 */
export function ImportConfigDialog({ open, onOpenChange, excludeServerId, targetGame = 'ase', mode = 'form', onImport }: ImportConfigDialogProps) {
  const t = useTranslations('serverConfig');
  const tCommon = useTranslations('common');
  const apiError = useApiError();
  const servers = useServersOnOpen(open).filter((s) => s.id !== excludeServerId);
  const fileInput = useRef<HTMLInputElement>(null);

  const [source, setSource] = useState<'file' | 'server'>('file');
  const [imported, setImported] = useState<ImportedConfigFile | null>(null);
  const [fileName, setFileName] = useState('');
  const [fileError, setFileError] = useState('');
  const [sourceId, setSourceId] = useState('');
  const [parts, setParts] = useState<ConfigPart[]>(DEFAULT_CONFIG_PARTS);
  // 导入 ini 时，GameUserSettings.ini 的 ActiveMods 是否同时导入为 Mod 列表
  const [iniMods, setIniMods] = useState(true);
  const [loading, setLoading] = useState(false);

  // 每次打开时清空上次的选择
  const [prevOpen, setPrevOpen] = useState(open);
  if (open !== prevOpen) {
    setPrevOpen(open);
    if (open) {
      setImported(null);
      setFileName('');
      setFileError('');
      setSourceId('');
      setParts(DEFAULT_CONFIG_PARTS);
      setIniMods(true);
    }
  }

  const readFiles = async (files: File[]) => {
    if (files.length === 0) return;
    setFileName(files.map((f) => f.name).join('、'));
    setImported(null);
    setFileError('');
    setIniMods(true);
    try {
      setImported(await readImportFiles(files));
    } catch (err) {
      if (!(err instanceof InvalidConfigFileError)) setFileError(t('readFileFailed'));
      else if (err.reason === 'invalid') setFileError(t('invalidFile'));
      else setFileError(t(`fileError.${err.reason}`, { name: err.fileName ?? '' }));
    }
  };

  // 导入 ini 文件时只有「游戏设置」部分，不显示部分选择；ini 里有 ActiveMods 时可以同时导入 Mod 列表
  const iniImport = source === 'file' && imported?.kind === 'ini';
  const iniModCount = iniImport ? (imported.config.game_mod_ids?.split(',').length ?? 0) : 0;
  // 来源的游戏：导出的文件里记着（旧版没有，视为 ASE），ini 文件两个游戏通用
  const sourceServer = source === 'server' ? servers.find((s) => String(s.id) === sourceId) : undefined;
  const sourceGame: Game | null =
    source === 'server' ? (sourceServer ? gameOf(sourceServer) : null) : imported?.kind === 'json' ? gameOf(imported.config) : null;
  const crossGame = sourceGame !== null && sourceGame !== targetGame;
  const effectiveParts: ConfigPart[] = iniImport
    ? iniModCount > 0 && iniMods
      ? ['config_files', 'mods']
      : ['config_files']
    : parts.filter((p) => !(crossGame && CROSS_GAME_PARTS.includes(p)));
  const ready = effectiveParts.length > 0 && (source === 'file' ? imported !== null : sourceId !== '');

  const confirm = async () => {
    setLoading(true);
    try {
      let cfg = imported?.config ?? null;
      if (source === 'server') {
        try {
          cfg = await fetchServerConfig(Number(sourceId));
        } catch (err) {
          toast.error(apiError(err, t('loadFailed')));
          return;
        }
      }
      if (!cfg) return;
      // ini 文件不属于哪个游戏，按导入到的服务器处理（ActiveMods 里的 Mod 也按这个游戏导入）
      if (iniImport) cfg = { ...cfg, game: targetGame };
      await onImport(cfg, effectiveParts);
      onOpenChange(false);
    } catch {
      // 调用方负责提示错误
    } finally {
      setLoading(false);
    }
  };

  return (
    <Dialog open={open} onOpenChange={(next) => !loading && onOpenChange(next)}>
      <DialogContent className="max-w-lg">
        <DialogHeader>
          <DialogTitle>{t('importTitle')}</DialogTitle>
          <DialogDescription>{mode === 'apply' ? t('importApplyDesc') : t('importDesc')}</DialogDescription>
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
                accept=".json,.ini,application/json"
                multiple
                className="hidden"
                onChange={(e) => {
                  readFiles(Array.from(e.target.files ?? []));
                  e.target.value = '';
                }}
              />
              <Button type="button" variant="outline" className="w-full" onClick={() => fileInput.current?.click()}>
                <FileUp />
                {fileName || t('chooseFile')}
              </Button>
              {fileError && <p className="text-xs text-red-700">{fileError}</p>}
              {!imported && !fileError && <p className="text-xs text-gray-500">{t('chooseFileHint')}</p>}
              {imported?.kind === 'json' && (
                <p className="text-xs text-gray-600">
                  {imported.config.session_name ? t('fileFrom', { name: imported.config.session_name }) : t('fileLoaded')}
                </p>
              )}
              {imported?.kind === 'ini' && (
                <div className="space-y-1 text-xs text-gray-600">
                  <p>
                    {t('iniLoaded', {
                      files: imported.iniKinds.map((k) => (k === 'game_user_settings' ? 'GameUserSettings.ini' : 'Game.ini')).join('、'),
                    })}
                  </p>
                  {imported.iniKinds.length === 1 && (
                    <p className="text-gray-500">
                      {t('iniOtherUnchanged', { file: imported.iniKinds[0] === 'game_user_settings' ? 'Game.ini' : 'GameUserSettings.ini' })}
                    </p>
                  )}
                </div>
              )}
              {iniModCount > 0 && (
                <label className="flex items-start justify-between gap-4 pt-1">
                  <span>
                    <span className="font-medium">{t('iniMods', { count: iniModCount })}</span>
                    <span className="block text-xs text-gray-500">{t('iniModsDesc')}</span>
                  </span>
                  <Switch checked={iniMods} onCheckedChange={setIniMods} />
                </label>
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
                    <GameBadge game={gameOf(s)} className="px-1.5 py-0 text-[10px]" />
                    {s.session_name}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          )}

          {!iniImport && <ConfigPartsPicker value={parts} onChange={setParts} disabled={crossGame ? CROSS_GAME_PARTS : []} />}
          {crossGame && sourceGame && (
            <p className="text-xs text-amber-800">{t('crossGameNote', { source: GAME_NAMES[sourceGame], target: GAME_NAMES[targetGame] })}</p>
          )}
          <p className="text-xs text-gray-500">{t('keepIdentity')}</p>
        </div>

        <DialogFooter>
          <Button variant="outline" onClick={() => onOpenChange(false)} disabled={loading}>
            {tCommon('cancel')}
          </Button>
          <Button onClick={confirm} disabled={!ready || loading}>
            {loading && <Loader2 className="animate-spin" />}
            {mode === 'apply' ? t('importApplyConfirm') : t('importConfirm')}
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
  // 另一个游戏的目标服务器不复制地图、Mod 和启动参数（后端跳过）
  const sourceGame = gameOf(server);
  const anyCrossGame = targets.some((s) => selected.includes(s.id) && gameOf(s) !== sourceGame);

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
                      <GameBadge game={gameOf(s)} className="px-1.5 py-0 text-[10px]" />
                      <span className="flex-1 truncate">{s.session_name}</span>
                      {s.status !== 'stopped' && <span className="text-xs text-amber-800">{t('notStopped')}</span>}
                    </label>
                  </li>
                ))}
              </ul>
            )}
          </div>

          <ConfigPartsPicker value={parts} onChange={setParts} />
          {anyCrossGame && <p className="text-xs text-amber-800">{t('crossGameCopyNote', { game: GAME_NAMES[sourceGame] })}</p>}
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
