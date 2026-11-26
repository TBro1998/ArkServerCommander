"use client";

import { useState, type ReactNode } from 'react';
import Link from 'next/link';
import { useTranslations } from 'next-intl';
import { Play, Square, Loader2, RotateCw, SquareTerminal, ScrollText, Settings2 } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Tooltip, TooltipContent, TooltipTrigger } from '@/components/ui/tooltip';
import { ConfirmDialog } from '@/components/ui/confirm-dialog';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Switch } from '@/components/ui/switch';
import { serversActions, isTransitioning, type Server } from '@/stores/servers';
import { useApiError } from '@/lib/api-error';
import { toast } from '@/lib/toast';
import { cn } from '@/lib/utils';
import { gameOf } from '@/lib/games';

export type ServerDialog = 'stop' | 'restart' | 'update' | 'recreate' | 'delete';

/**
 * 服务器的启停、重启、更新、重建、删除操作和对应的确认弹窗。
 * 调用方渲染 controls.dialogs，再用 PowerButtons、管理面板等触发操作，同一台服务器的弹窗只有一份
 */
export function useServerControls(server: Server, options: { onDeleted?: () => void } = {}) {
  const t = useTranslations('serverActions');
  const apiError = useApiError();
  const [dialog, setDialog] = useState<ServerDialog | null>(null);
  const [pending, setPending] = useState(false);
  const [updateGame, setUpdateGame] = useState(true);
  const [updateMods, setUpdateMods] = useState(true);
  const [deleteData, setDeleteData] = useState(false);
  const [confirmName, setConfirmName] = useState('');

  const busy = isTransitioning(server.status) || pending;
  const running = server.status === 'running';
  const stopped = server.status === 'stopped';
  const players = server.runtime?.players ?? 0;
  const hasMods = !!server.game_mod_ids;
  // ASA 的 Mod 由服务端每次启动时自己从 CurseForge 下载和更新，"更新"只需要更新游戏本体
  const asa = gameOf(server) === 'asa';
  const canUpdateMods = hasMods && !asa;

  // 执行操作：失败时弹出错误并抛出，让确认框保持打开
  const run = async (action: () => Promise<void>, success: string, failure: string) => {
    setPending(true);
    try {
      await action();
      toast.success(success);
    } catch (err) {
      toast.error(apiError(err, failure));
      throw err;
    } finally {
      setPending(false);
    }
  };

  const start = () => run(() => serversActions.startServer(server.id), t('startSent'), t('startFailed')).catch(() => {});
  const stop = () => run(() => serversActions.stopServer(server.id), t('stopSent'), t('stopFailed'));
  const restart = () => run(() => serversActions.restartServer(server.id), t('restartSent'), t('restartFailed'));
  const recreate = () => run(() => serversActions.recreateServer(server.id), t('recreateSent'), t('recreateFailed'));
  const update = () =>
    run(() => serversActions.updateGame(server.id, updateGame, updateMods && canUpdateMods), t('updateSent'), t('updateFailed'));
  const remove = () =>
    run(
      async () => {
        await serversActions.deleteServer(server.id, deleteData);
        options.onDeleted?.();
      },
      t('deleteSuccess'),
      t('deleteFailed'),
    );

  const open = (kind: ServerDialog) => {
    setConfirmName('');
    setDeleteData(false);
    setUpdateGame(true);
    setUpdateMods(true);
    setDialog(kind);
  };

  // 没有玩家在线时直接停止，有玩家时先确认
  const requestStop = (alwaysConfirm = false) => {
    if (players > 0 || alwaysConfirm) open('stop');
    else stop().catch(() => {});
  };

  const close = (o: boolean) => !o && setDialog(null);

  const dialogs = (
    <>
      <ConfirmDialog
        open={dialog === 'stop'}
        onOpenChange={close}
        title={t('stopTitle', { name: server.session_name })}
        description={players > 0 ? t('stopWithPlayers', { count: players }) : t('stopDesc')}
        confirmText={t('stop')}
        destructive
        onConfirm={stop}
      />

      <ConfirmDialog
        open={dialog === 'restart'}
        onOpenChange={close}
        title={t('restartTitle', { name: server.session_name })}
        description={
          <>
            <p>{running ? t('restartDesc') : t('restartStoppedDesc')}</p>
            {players > 0 && <p className="mt-2 font-medium text-amber-700">{t('playersOnline', { count: players })}</p>}
          </>
        }
        confirmText={t('restart')}
        onConfirm={restart}
      />

      <ConfirmDialog
        open={dialog === 'update'}
        onOpenChange={close}
        title={t('updateTitle', { name: server.session_name })}
        description={running ? t('updateDescRunning') : t('updateDesc')}
        confirmText={t('updateConfirm')}
        confirmDisabled={!updateGame && !(updateMods && canUpdateMods)}
        onConfirm={update}
      >
        <div className="space-y-3 text-sm">
          <label className="flex items-start justify-between gap-4">
            <span>
              <span className="font-medium">{t('updateGameOption')}</span>
              <span className="block text-xs text-gray-500">{t('updateGameOptionDesc')}</span>
            </span>
            <Switch checked={updateGame} onCheckedChange={setUpdateGame} />
          </label>
          {asa ? (
            <p className="text-xs text-gray-500">{t('updateModsAsa')}</p>
          ) : (
            <label className="flex items-start justify-between gap-4">
              <span>
                <span className="font-medium">{t('updateModsOption')}</span>
                <span className="block text-xs text-gray-500">{hasMods ? t('updateModsOptionDesc') : t('noMods')}</span>
              </span>
              <Switch checked={updateMods && hasMods} onCheckedChange={setUpdateMods} disabled={!hasMods} />
            </label>
          )}
          {players > 0 && <p className="font-medium text-amber-700">{t('playersOnline', { count: players })}</p>}
        </div>
      </ConfirmDialog>

      <ConfirmDialog
        open={dialog === 'recreate'}
        onOpenChange={close}
        title={t('recreateTitle', { name: server.session_name })}
        description={t('recreateDesc')}
        confirmText={t('recreate')}
        onConfirm={recreate}
      />

      <ConfirmDialog
        open={dialog === 'delete'}
        onOpenChange={close}
        title={t('deleteTitle', { name: server.session_name })}
        description={t('deleteDesc')}
        confirmText={t('delete')}
        destructive
        confirmDisabled={deleteData && confirmName !== server.session_name}
        onConfirm={remove}
      >
        <div className="space-y-3 text-sm">
          <label className="flex items-start justify-between gap-4">
            <span>
              <span className="font-medium">{t('deleteData')}</span>
              <span className="block text-xs text-gray-500">{deleteData ? t('deleteDataOn') : t('deleteDataOff')}</span>
            </span>
            <Switch checked={deleteData} onCheckedChange={setDeleteData} />
          </label>
          {deleteData && (
            <div className="space-y-1">
              <Label htmlFor={`confirm-${server.id}`}>{t('typeNameToConfirm', { name: server.session_name })}</Label>
              <Input
                id={`confirm-${server.id}`}
                value={confirmName}
                onChange={(e) => setConfirmName(e.target.value)}
                autoComplete="off"
              />
            </div>
          )}
        </div>
      </ConfirmDialog>
    </>
  );

  return { server, busy, running, stopped, players, hasMods, start, requestStop, open, dialogs };
}

export type ServerControls = ReturnType<typeof useServerControls>;

// ---------- 按钮 ----------

/** 只有图标的按钮，悬停显示说明。禁用时浏览器不触发悬停，提示放在外层 span 上 */
export function IconAction({
  label,
  icon,
  onClick,
  href,
  disabled,
  className,
}: {
  label: string;
  icon: ReactNode;
  onClick?: () => void;
  href?: string;
  disabled?: boolean;
  className?: string;
}) {
  const button = href ? (
    <Button asChild variant="ghost" size="icon" className={cn('size-8 text-gray-700 hover:text-primary', className)} aria-label={label}>
      <Link href={href}>{icon}</Link>
    </Button>
  ) : (
    <Button
      variant="ghost"
      size="icon"
      className={cn('size-8 text-gray-700 hover:text-primary', className)}
      onClick={onClick}
      disabled={disabled}
      aria-label={label}
    >
      {icon}
    </Button>
  );
  return (
    <Tooltip>
      <TooltipTrigger asChild>{disabled ? <span className="inline-flex">{button}</span> : button}</TooltipTrigger>
      <TooltipContent sideOffset={4}>{label}</TooltipContent>
    </Tooltip>
  );
}

/**
 * 启动 / 停止和重启。compact：卡片、列表、首页用，重启只显示图标；否则带文字（详情页）。
 * 已停止时不显示重启（与启动相同）
 */
export function PowerButtons({ controls, compact = false }: { controls: ServerControls; compact?: boolean }) {
  const t = useTranslations('serverActions');
  const { busy, running, stopped } = controls;

  let primary;
  if (busy) {
    primary = (
      <Button variant="outline" size="sm" disabled className="min-w-[5.5rem]">
        <Loader2 className="animate-spin" />
        {t('working')}
      </Button>
    );
  } else if (running) {
    primary = (
      <Button
        variant="outline"
        size="sm"
        className="min-w-[5.5rem] hover:border-red-200 hover:bg-red-50 hover:text-red-800"
        onClick={() => controls.requestStop(!compact)}
      >
        <Square className="size-2.5 fill-red-600 text-red-600" />
        {t('stop')}
      </Button>
    );
  } else {
    primary = (
      <Button size="sm" className="min-w-[5.5rem]" onClick={controls.start} disabled={!stopped}>
        <Play className="size-3.5 fill-current" />
        {t('start')}
      </Button>
    );
  }

  const showRestart = !stopped;
  return (
    <div className="flex items-center gap-1.5">
      {primary}
      {showRestart &&
        (compact ? (
          <IconAction label={t('restart')} icon={<RotateCw />} onClick={() => controls.open('restart')} disabled={busy} />
        ) : (
          <Button variant="outline" size="sm" onClick={() => controls.open('restart')} disabled={busy}>
            <RotateCw />
            {t('restart')}
          </Button>
        ))}
    </div>
  );
}

/** 常用页面的快捷入口：控制台、日志、编辑设置 */
export function ServerShortcuts({ server }: { server: Server }) {
  const t = useTranslations('serverActions');
  return (
    <div className="flex items-center gap-0.5">
      <IconAction label={t('console')} icon={<SquareTerminal />} href={`/servers/${server.id}?tab=console`} />
      <IconAction label={t('logs')} icon={<ScrollText />} href={`/servers/${server.id}?tab=logs`} />
      <IconAction label={t('settings')} icon={<Settings2 />} href={`/servers/${server.id}/edit`} />
    </div>
  );
}

/** 只有启停按钮的独立组件（首页名单） */
export function ServerPowerControls({ server }: { server: Server }) {
  const controls = useServerControls(server);
  return (
    <>
      <PowerButtons controls={controls} compact />
      {controls.dialogs}
    </>
  );
}
