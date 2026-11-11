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
      // 整行都可以点击复制，目标更大
      <button
        type="button"
        onClick={() => copy(gameAddress)}
        title={t('copyAddress')}
        className="group flex w-full items-center justify-between gap-2 rounded-md border bg-background px-2.5 py-1.5 text-left text-xs transition-colors hover:border-primary/40 hover:bg-accent"
      >
        <span className="text-muted-foreground">{t('address')}</span>
        <span className="flex min-w-0 items-center gap-1.5">
          <span className="truncate font-mono font-medium">{gameAddress}</span>
          <Copy className="h-3.5 w-3.5 flex-shrink-0 text-muted-foreground group-hover:text-primary" />
        </span>
      </button>
    );
  }

  return (
    <div className="space-y-2 text-sm">
      <div className="flex flex-wrap items-center justify-between gap-2">
        <span className="text-muted-foreground">{t('address')}</span>
        <span className="flex items-center gap-1">
          <span className="font-mono font-medium">{gameAddress}</span>
          <Button variant="ghost" size="sm" className="h-7 w-7 p-0" onClick={() => copy(gameAddress)} title={t('copyAddress')}>
            <Copy className="h-3.5 w-3.5" />
          </Button>
        </span>
      </div>
      <div className="flex flex-wrap items-center justify-between gap-2">
        <span className="text-muted-foreground">{t('steamConnect')}</span>
        <span className="flex items-center gap-1">
          <a href={steamUrl} className="inline-flex items-center gap-1 font-mono text-primary hover:underline">
            {host}:{server.query_port}
            <ExternalLink className="h-3.5 w-3.5" />
          </a>
          <Button variant="ghost" size="sm" className="h-7 w-7 p-0" onClick={() => copy(steamUrl)} title={t('copyLink')}>
            <Copy className="h-3.5 w-3.5" />
          </Button>
        </span>
      </div>
      <p className="text-xs text-muted-foreground">{t('hint')}</p>
    </div>
  );
}
