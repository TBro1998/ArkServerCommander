"use client";

import { useTranslations } from 'next-intl';
import { Loader2 } from 'lucide-react';
import { cn } from '@/lib/utils';
import type { Server } from '@/stores/servers';

// 三种颜色对应三种含义：绿色可用，琥珀色进行中，石色已停止
const OK = 'bg-green-100 text-green-800 border-green-200';
const BUSY = 'bg-amber-50 text-amber-800 border-amber-200';
const OFF = 'bg-gray-100 text-gray-700 border-gray-200';

const STATUS_STYLES: Record<string, string> = {
  running: OK,
  stopped: OFF,
  starting: BUSY,
  stopping: BUSY,
  restarting: BUSY,
};

// 状态徽章：运行中的服务器按阶段显示（下载服务端 / 安装 Mod / 加载中 / 可连接）
export function StatusBadge({ server, className }: { server: Server; className?: string }) {
  const t = useTranslations('status');
  const tPhase = useTranslations('phase');
  const runtime = server.runtime;

  let label: string;
  let style: string;
  let busy = false;

  // 运行中按阶段显示；启动中正在下载镜像、复制服务端文件时也显示阶段
  if (runtime?.phase && (server.status === 'running' || (server.status === 'starting' && (runtime.phase === 'copying' || runtime.phase === 'pulling_image')))) {
    const phase = runtime.phase;
    style = phase === 'ready' ? OK : BUSY;
    busy = phase !== 'ready';
    if (phase === 'downloading') {
      label = runtime.progress ? tPhase('downloadingProgress', { progress: runtime.progress.toFixed(1) }) : tPhase('downloading');
    } else if (phase === 'installing_mods' && runtime.mod_total) {
      label = tPhase('installingModsProgress', { current: runtime.mod_current ?? 0, total: runtime.mod_total });
    } else {
      label = tPhase(phase);
    }
  } else {
    const known = ['running', 'stopped', 'starting', 'stopping', 'restarting'];
    label = known.includes(server.status) ? t(server.status) : t('unknown');
    style = STATUS_STYLES[server.status] ?? OFF;
    busy = server.status === 'starting' || server.status === 'stopping' || server.status === 'restarting';
  }

  return (
    <span
      className={cn('num inline-flex items-center gap-1.5 rounded-full border px-2.5 py-0.5 text-xs font-medium whitespace-nowrap', style, className)}
    >
      {busy && <Loader2 className="h-3 w-3 animate-spin" />}
      {!busy && (
        <span
          className={cn(
            'h-2 w-2 rounded-full',
            server.status === 'running' ? 'bg-green-500 shadow-[0_0_0_3px] shadow-green-500/20' : server.status === 'stopped' ? 'bg-gray-400' : 'bg-current',
          )}
        />
      )}
      {label}
    </span>
  );
}
