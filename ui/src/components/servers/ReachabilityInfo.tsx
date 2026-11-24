"use client";

import { useState } from 'react';
import { useLocale, useTranslations } from 'next-intl';
import axios from 'axios';
import { AlertTriangle, CheckCircle2, Loader2, RefreshCw, ShieldQuestion } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { useApiError } from '@/lib/api-error';
import { toast } from '@/lib/toast';
import { serversActions, type Server } from '@/stores/servers';

// 公网可达性：管理端用玩家连接的地址（PUBLIC_HOST 或自动探测的公网 IP）查询这台服务器的查询端口的结果
// （服务器每次可以连接后自动检测一次）。容器里的 A2S 只说明游戏起来了，有没有发布到公网要看这里
export function ReachabilityInfo({ server }: { server: Server }) {
  const t = useTranslations('reachability');
  const locale = useLocale();
  const apiError = useApiError();
  const [busy, setBusy] = useState(false);
  const r = server.reachability;
  if (!r) return null;

  const recheck = async () => {
    setBusy(true);
    try {
      await axios.post(`/api/servers/${server.id}/reachability`);
      await serversActions.fetchServers();
    } catch (err) {
      toast.error(apiError(err, t('recheckFailed')));
    } finally {
      setBusy(false);
    }
  };

  const checking = r.state === 'checking' || busy;
  const addr = { ip: r.ip ?? '', port: server.query_port };
  const button = (
    <Button variant="outline" size="sm" className="h-7 flex-shrink-0 px-2 text-xs" onClick={recheck} disabled={checking}>
      {checking ? <Loader2 className="size-3.5 animate-spin" /> : <RefreshCw className="size-3.5" />}
      {t('recheck')}
    </Button>
  );
  if (checking) {
    return (
      <div className="flex items-center gap-2 rounded-md border px-3 py-2 text-sm text-amber-800">
        <Loader2 className="size-4 flex-shrink-0 animate-spin" />
        {t('checking', { port: server.query_port })}
      </div>
    );
  }

  if (r.state === 'reachable') {
    return (
      <div className="space-y-1 rounded-md border px-3 py-2 text-sm">
        <div className="flex items-center justify-between gap-2">
          <span className="flex items-center gap-2 text-green-700">
            <CheckCircle2 className="size-4 flex-shrink-0" />
            {t('reachable', addr)}
          </span>
          {button}
        </div>
        <p className="text-xs text-muted-foreground">{t('reachableHint', { port: server.port, port1: server.port + 1 })}</p>
      </div>
    );
  }

  if (r.state === 'not_public') {
    return (
      <div className="space-y-1 rounded-md border border-amber-300 bg-amber-50 px-3 py-2 text-sm text-amber-900">
        <div className="flex items-start justify-between gap-2">
          <span className="flex items-start gap-2 font-medium">
            <AlertTriangle className="mt-0.5 size-4 flex-shrink-0" />
            {t('notPublic', addr)}
          </span>
          {button}
        </div>
        <p className="text-xs">{t('notPublicHint')}</p>
      </div>
    );
  }

  if (r.state === 'error') {
    const detail = locale === 'zh' && r.error ? r.error : r.error_code && t.has(`errors.${r.error_code}`) ? t(`errors.${r.error_code}`) : (r.error ?? '');
    return (
      <div className="flex items-start justify-between gap-2 rounded-md border px-3 py-2 text-sm">
        <span className="flex items-start gap-2 text-muted-foreground">
          <ShieldQuestion className="mt-0.5 size-4 flex-shrink-0" />
          <span>
            {t('error')}
            {detail && <span className="block text-xs">{detail}</span>}
          </span>
        </span>
        {button}
      </div>
    );
  }

  // unreachable / mismatch
  const mismatch = r.state === 'mismatch';
  return (
    <div className="space-y-2 rounded-md border border-red-200 bg-red-50 px-3 py-2 text-sm text-red-800">
      <div className="flex items-start justify-between gap-2">
        <span className="flex items-start gap-2 font-medium">
          <AlertTriangle className="mt-0.5 size-4 flex-shrink-0" />
          {mismatch ? t('mismatch', addr) : t('unreachable', addr)}
        </span>
        {button}
      </div>
      {mismatch ? (
        <p className="text-xs">{r.name ? t('mismatchHint', { name: r.name }) : t('mismatchHintNoName')}</p>
      ) : (
        <ul className="list-disc space-y-0.5 pl-5 text-xs">
          <li>{t('causeNoPublicIp')}</li>
          <li>{t('causeForward', { port: server.port, port1: server.port + 1, query: server.query_port })}</li>
          <li>{t('causeFirewall')}</li>
        </ul>
      )}
      {!mismatch && <p className="text-xs">{t('loopbackNote')}</p>}
      {!mismatch && <p className="text-xs">{t('linkSuggestion')}</p>}
    </div>
  );
}
