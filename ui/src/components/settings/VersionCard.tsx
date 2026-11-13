"use client";

import { useCallback, useEffect, useState } from 'react';
import { useLocale, useTranslations } from 'next-intl';
import axios from 'axios';
import { Copy, ExternalLink, Info, Loader2, RefreshCw } from 'lucide-react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { useApiError } from '@/lib/api-error';
import { copyText, formatDateTime } from '@/lib/format';
import { toast } from '@/lib/toast';

interface VersionInfo {
  current: string;
  latest?: string;
  has_update: boolean;
  release_url?: string;
  published_at?: string;
  checked_at?: string;
  error?: string;
}

const UPGRADE_COMMAND = 'docker compose pull && docker compose up -d';

// 设置页"关于"：当前版本、最新发布版本和升级方法。管理器不自我更新（要替换正在运行的自身容器），只给出命令
export function VersionCard() {
  const t = useTranslations('settings.version');
  const locale = useLocale();
  const apiError = useApiError();
  const [info, setInfo] = useState<VersionInfo | null>(null);
  const [checking, setChecking] = useState(false);

  const load = useCallback(
    async (refresh = false) => {
      try {
        const res = await axios.get('/api/system/version', { params: refresh ? { refresh: true } : undefined });
        setInfo(res.data.data);
      } catch (err) {
        toast.error(apiError(err, t('loadFailed')));
      }
    },
    [apiError, t],
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

  const dev = info?.current === 'dev';

  return (
    <Card>
      <CardHeader className="flex flex-row flex-wrap items-center justify-between gap-2 space-y-0">
        <CardTitle className="flex items-center gap-2 text-lg">
          <Info className="h-5 w-5" />
          {t('title')}
        </CardTitle>
        <Button variant="outline" size="sm" onClick={check} disabled={checking}>
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

        {info?.error && <p className="text-xs text-red-600">{t('checkFailed', { error: info.error })}</p>}

        {info?.has_update ? (
          <div className="space-y-2 rounded-sm border border-amber-200 bg-amber-50 p-3 text-amber-800">
            <p>{t('updateAvailable', { version: info.latest ?? '' })}</p>
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
          </div>
        ) : info && !info.error && info.latest ? (
          <p className="text-xs text-muted-foreground">{dev ? t('devHint') : t('upToDate')}</p>
        ) : null}
      </CardContent>
    </Card>
  );
}
