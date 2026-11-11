"use client";

import { useTranslations } from 'next-intl';
import { Loader2 } from 'lucide-react';
import { cn } from '@/lib/utils';
import type { Server } from '@/stores/servers';

const STATUS_STYLES: Record<string, string> = {
  running: 'bg-green-100 text-green-800 border-green-200',
  stopped: 'bg-gray-100 text-gray-700 border-gray-200',
  starting: 'bg-blue-100 text-blue-800 border-blue-200',
  stopping: 'bg-amber-100 text-amber-800 border-amber-200',
};

const PHASE_STYLES: Record<string, string> = {
  copying: 'bg-cyan-100 text-cyan-800 border-cyan-200',
  downloading: 'bg-sky-100 text-sky-800 border-sky-200',
  installing_mods: 'bg-indigo-100 text-indigo-800 border-indigo-200',
  launching: 'bg-blue-100 text-blue-800 border-blue-200',
  ready: 'bg-green-100 text-green-800 border-green-200',
};

// 状态徽章：运行中的服务器按阶段显示（下载服务端 / 安装 Mod / 加载中 / 可连接）
export function StatusBadge({ server, className }: { server: Server; className?: string }) {
  const t = useTranslations('status');
  const tPhase = useTranslations('phase');
  const runtime = server.runtime;

  let label: string;
  let style: string;
  let busy = false;

  // 运行中按阶段显示；启动中正在复制服务端文件时也显示阶段
  if (runtime?.phase && (server.status === 'running' || (server.status === 'starting' && runtime.phase === 'copying'))) {
    const phase = runtime.phase;
    style = PHASE_STYLES[phase] ?? STATUS_STYLES.running;
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
    style = STATUS_STYLES[server.status] ?? 'bg-gray-100 text-gray-600 border-gray-200';
    busy = server.status === 'starting' || server.status === 'stopping' || server.status === 'restarting';
  }

  return (
    <span
      className={cn('inline-flex items-center gap-1 rounded-md border px-2 py-0.5 text-xs font-medium whitespace-nowrap', style, className)}
    >
      {busy && <Loader2 className="h-3 w-3 animate-spin" />}
      {!busy && (
        <span
          className={cn(
            'h-1.5 w-1.5 rounded-full',
            server.status === 'running' ? 'bg-green-500' : server.status === 'stopped' ? 'bg-gray-400' : 'bg-current',
          )}
        />
      )}
      {label}
    </span>
  );
}
