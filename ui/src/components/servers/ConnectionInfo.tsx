"use client";

import { useTranslations } from 'next-intl';
import { Copy, ExternalLink } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { copyText } from '@/lib/format';
import { toast } from '@/lib/toast';
import type { Server } from '@/stores/servers';

// 连接地址：游戏内"加入"用 IP:游戏端口，Steam 服务器浏览器 / steam://connect 用 IP:查询端口
export function ConnectionInfo({ server, host, compact = false }: { server: Server; host: string; compact?: boolean }) {
  const t = useTranslations('connection');
  if (!host) return null;

  const gameAddress = `${host}:${server.port}`;
  const steamUrl = `steam://connect/${host}:${server.query_port}`;

  const copy = async (text: string) => {
    if (await copyText(text)) toast.success(t('copied', { text }));
    else toast.error(t('copyFailed'));
  };

  if (compact) {
    return (
      <div className="flex items-center justify-between gap-2 rounded-lg bg-blue-50 px-2 py-1.5 text-xs">
        <span className="text-gray-600">{t('address')}</span>
        <span className="flex min-w-0 items-center gap-1">
          <span className="truncate font-mono font-medium">{gameAddress}</span>
          <Button variant="ghost" size="sm" className="h-6 w-6 p-0" onClick={() => copy(gameAddress)} title={t('copyAddress')}>
            <Copy className="h-3 w-3" />
          </Button>
        </span>
      </div>
    );
  }

  return (
    <div className="space-y-2 text-sm">
      <div className="flex flex-wrap items-center justify-between gap-2">
        <span className="text-gray-600">{t('address')}</span>
        <span className="flex items-center gap-1">
          <span className="font-mono font-medium">{gameAddress}</span>
          <Button variant="ghost" size="sm" className="h-7 w-7 p-0" onClick={() => copy(gameAddress)} title={t('copyAddress')}>
            <Copy className="h-3.5 w-3.5" />
          </Button>
        </span>
      </div>
      <div className="flex flex-wrap items-center justify-between gap-2">
        <span className="text-gray-600">{t('steamConnect')}</span>
        <span className="flex items-center gap-1">
          <a href={steamUrl} className="inline-flex items-center gap-1 font-mono text-blue-600 hover:underline">
            {host}:{server.query_port}
            <ExternalLink className="h-3.5 w-3.5" />
          </a>
          <Button variant="ghost" size="sm" className="h-7 w-7 p-0" onClick={() => copy(steamUrl)} title={t('copyLink')}>
            <Copy className="h-3.5 w-3.5" />
          </Button>
        </span>
      </div>
      <p className="text-xs text-gray-500">{t('hint')}</p>
    </div>
  );
}
