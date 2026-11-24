"use client";

import { useCallback, useEffect, useState } from 'react';
import { useLocale, useTranslations } from 'next-intl';
import axios from 'axios';
import { Copy, Download, ExternalLink, Info, Loader2, RefreshCw } from 'lucide-react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Progress } from '@/components/ui/progress';
import { ConfirmDialog } from '@/components/ui/confirm-dialog';
import { useApiError, useErrorText } from '@/lib/api-error';
import { copyText, formatDateTime } from '@/lib/format';
import { toast } from '@/lib/toast';
import { usePolling } from '@/hooks/use-polling';

// 最近一次一键更新的记录（数据目录的 .self-update.json）
interface SelfUpdateRecord {
  state: 'started' | 'success' | 'rolled_back' | 'failed' | 'interrupted';
  from_version: string;
  to_version?: string;
  error?: string;
  started_at: string;
  finished_at?: string;
}

interface SelfUpdateStatus {
  supported: boolean;
  reason?: 'dev' | 'not_container' | 'custom_image' | 'no_socket' | 'unavailable';
  state: 'idle' | 'pulling' | 'restarting' | 'error';
  progress?: number;
  error_code?: string;
  error?: string;
  last?: SelfUpdateRecord;
}

interface VersionInfo {
  current: string;
  latest?: string;
  has_update: boolean;
  release_url?: string;
  published_at?: string;
  checked_at?: string;
  error?: string;
  self_update?: SelfUpdateStatus;
}

// 正在跟踪的一次更新：from 是更新前的版本；ignore 是更新开始前已有的那条结果记录的 started_at，不能当成这次的结果
interface Tracking {
  from: string;
  ignore?: string;
}

const UPGRADE_COMMAND = 'docker compose pull && docker compose up -d';
const QQ_GROUP = '953242676';
// 失败的更新结果在卡片上显示多久
const FAILURE_VISIBLE_MS = 3 * 24 * 3600 * 1000;

function isUpdating(su?: SelfUpdateStatus) {
  return su?.state === 'pulling' || su?.state === 'restarting';
}

function recentFailure(su?: SelfUpdateStatus): SelfUpdateRecord | null {
  const last = su?.last;
  if (!last || last.state === 'success' || last.state === 'started') return null;
  const at = Date.parse(last.finished_at ?? last.started_at);
  return Number.isFinite(at) && Date.now() - at < FAILURE_VISIBLE_MS ? last : null;
}

// 设置页"关于"：当前版本、最新发布版本和升级。用发布的 latest 镜像运行时可以一键更新（由临时的更新容器替换管理端容器），
// 否则只给出命令
export function VersionCard() {
  const t = useTranslations('settings.version');
  const locale = useLocale();
  const apiError = useApiError();
  const errorText = useErrorText();
  const [info, setInfo] = useState<VersionInfo | null>(null);
  const [checking, setChecking] = useState(false);
  const [confirmOpen, setConfirmOpen] = useState(false);
  const [tracking, setTracking] = useState<Tracking | null>(null);
  // 更新期间管理端在重启，请求失败是正常的
  const [offline, setOffline] = useState(false);
  const [failure, setFailure] = useState<SelfUpdateRecord | null>(null);

  const apply = useCallback((data: VersionInfo) => {
    setInfo(data);
    setFailure(recentFailure(data.self_update));
  }, []);

  const load = useCallback(
    async (refresh = false) => {
      try {
        const res = await axios.get('/api/system/version', { params: refresh ? { refresh: true } : undefined });
        const data: VersionInfo = res.data.data;
        apply(data);
        // 打开页面时更新已经在进行（例如刷新了页面）
        if (isUpdating(data.self_update)) {
          const last = data.self_update?.last;
          setTracking((cur) => cur ?? { from: data.current, ignore: last?.state === 'started' ? undefined : last?.started_at });
        }
      } catch (err) {
        toast.error(apiError(err, t('loadFailed')));
      }
    },
    [apiError, apply, t],
  );

  useEffect(() => {
    const timer = setTimeout(() => {
      load();
    }, 0);
    return () => clearTimeout(timer);
  }, [load]);

  const check = async () => {
    setChecking(true);
    await load(true);
    setChecking(false);
  };

  const startUpdate = async () => {
    if (!info) return;
    try {
      await axios.post('/api/system/update', { version: info.latest ?? '' });
    } catch (err) {
      toast.error(apiError(err, t('updateStartFailed')));
      throw err;
    }
    const last = info.self_update?.last;
    setFailure(null);
    setTracking({ from: info.current, ignore: last?.started_at });
    setInfo({ ...info, self_update: { ...(info.self_update ?? { supported: true }), state: 'pulling', progress: 0 } });
  };

  const poll = useCallback(async () => {
    if (!tracking) return;
    let data: VersionInfo;
    try {
      const res = await axios.get('/api/system/version');
      data = res.data.data;
    } catch {
      setOffline(true);
      return;
    }
    setOffline(false);
    apply(data);
    const su = data.self_update;
    const last = su?.last;

    const succeeded = () => {
      setTracking(null);
      toast.success(t('updated', { version: data.current }));
      // 新版本的前端资源和旧页面不兼容，重新加载
      setTimeout(() => window.location.reload(), 1500);
    };
    if (data.current !== tracking.from) {
      succeeded();
      return;
    }
    if (last && last.state !== 'started' && last.started_at !== tracking.ignore) {
      if (last.state === 'success') {
        succeeded();
      } else {
        setTracking(null);
        toast.error(t('updateRolledBack'));
      }
      return;
    }
    if (su?.state === 'error') {
      setTracking(null);
      toast.error(errorText(su.error, su.error_code, t('updateStartFailed')));
    }
  }, [apply, errorText, t, tracking]);

  usePolling(poll, 2000, tracking !== null);

  const dev = info?.current === 'dev';
  const su = info?.self_update;
  const updating = tracking !== null || isUpdating(su);
  const canSelfUpdate = !!su?.supported && !updating;
  const unsupportedHint = su && !su.supported && su.reason && su.reason !== 'dev' ? t(`unsupported.${su.reason}`) : null;

  return (
    <Card>
      <CardHeader className="flex flex-row flex-wrap items-center justify-between gap-2 space-y-0">
        <CardTitle className="flex items-center gap-2 text-lg">
          <Info className="h-5 w-5" />
          {t('title')}
        </CardTitle>
        <Button variant="outline" size="sm" onClick={check} disabled={checking || updating}>
          {checking ? <Loader2 className="mr-1.5 animate-spin" /> : <RefreshCw className="mr-1.5" />}
          {t('check')}
        </Button>
      </CardHeader>
      <CardContent className="space-y-3 text-sm">
        <div className="flex justify-between py-1">
          <span className="text-muted-foreground">{t('current')}</span>
          <span className="num font-medium">{!info ? '-' : dev ? t('dev') : `v${info.current}`}</span>
        </div>
        <div className="flex justify-between py-1">
          <span className="text-muted-foreground">{t('latest')}</span>
          <span className="num flex items-center gap-2 font-medium">
            {info?.latest ? (
              <a href={info.release_url} target="_blank" rel="noreferrer" className="inline-flex items-center gap-1 hover:text-primary hover:underline">
                v{info.latest}
                <ExternalLink className="size-3" />
              </a>
            ) : (
              '-'
            )}
            {info?.published_at && <span className="text-xs font-normal text-muted-foreground">{formatDateTime(info.published_at, locale)}</span>}
          </span>
        </div>

        <div className="flex items-center justify-between py-1">
          <span className="text-muted-foreground">{t('qqGroup')}</span>
          <span className="flex items-center gap-1">
            <span className="num font-medium select-all">{QQ_GROUP}</span>
            <Button
              variant="ghost"
              size="icon"
              className="size-7"
              onClick={async () => {
                if (await copyText(QQ_GROUP)) toast.success(t('qqCopied'));
                else toast.error(t('copyFailed'));
              }}
              aria-label={t('copyQQ')}
              title={t('copyQQ')}
            >
              <Copy />
            </Button>
          </span>
        </div>

        {info?.error && <p className="text-xs text-red-600">{t('checkFailed', { error: info.error })}</p>}

        {updating && (
          <div className="space-y-2 rounded-sm border border-amber-200 bg-amber-50 p-3 text-amber-800">
            <p className="flex items-center gap-2 font-medium">
              <Loader2 className="size-4 animate-spin" />
              {su?.state === 'pulling' && !offline
                ? t('pulling', { progress: su.progress ?? 0 })
                : t('restarting')}
            </p>
            {su?.state === 'pulling' && !offline && <Progress value={su.progress ?? 0} />}
            <p className="text-xs">{t('updatingHint')}</p>
          </div>
        )}

        {!updating && failure && (
          <div className="space-y-1 rounded-sm border border-red-200 bg-red-50 p-3 text-red-800">
            <p className="font-medium">
              {failure.state === 'failed' ? t('lastFailed') : failure.state === 'interrupted' ? t('lastInterrupted') : t('lastRolledBack')}
            </p>
            {failure.error && <p className="break-all font-mono text-xs">{failure.error}</p>}
            <p className="text-xs text-muted-foreground">{formatDateTime(failure.finished_at ?? failure.started_at, locale)}</p>
          </div>
        )}

        {info?.has_update && !updating ? (
          <div className="space-y-2 rounded-sm border border-amber-200 bg-amber-50 p-3 text-amber-800">
            <p>{t('updateAvailable', { version: info.latest ?? '' })}</p>
            {canSelfUpdate && (
              <Button size="sm" onClick={() => setConfirmOpen(true)}>
                <Download className="mr-1.5" />
                {t('updateNow')}
              </Button>
            )}
            <p className="text-xs">{canSelfUpdate ? t('manualCommand') : t('runCommand')}</p>
            <div className="flex items-center gap-2">
              <code className="min-w-0 flex-1 truncate rounded-sm bg-terminal px-2 py-1 font-mono text-xs text-terminal-foreground">{UPGRADE_COMMAND}</code>
              <Button
                variant="outline"
                size="sm"
                onClick={async () => {
                  if (await copyText(UPGRADE_COMMAND)) toast.success(t('copied'));
                  else toast.error(t('copyFailed'));
                }}
                aria-label={t('copy')}
              >
                <Copy />
              </Button>
            </div>
            <p className="text-xs">{t('upgradeHint')}</p>
            {unsupportedHint && <p className="text-xs">{unsupportedHint}</p>}
          </div>
        ) : info && !info.error && info.latest && !updating ? (
          <p className="text-xs text-muted-foreground">{dev ? t('devHint') : t('upToDate')}</p>
        ) : null}
      </CardContent>

      <ConfirmDialog
        open={confirmOpen}
        onOpenChange={setConfirmOpen}
        title={t('confirmTitle', { version: info?.latest ?? '' })}
        description={
          <span className="block space-y-2">
            <span className="block">{t('confirmSteps')}</span>
            <span className="block">{t('confirmSafe')}</span>
            <span className="block">{t('confirmRollback')}</span>
          </span>
        }
        confirmText={t('updateNow')}
        onConfirm={startUpdate}
      />
    </Card>
  );
}
