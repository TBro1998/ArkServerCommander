"use client";

import { useState } from 'react';
import { useTranslations } from 'next-intl';
import { useRouter } from 'next/navigation';
import {
  Play,
  Square,
  Loader2,
  MoreVertical,
  RotateCw,
  RotateCcw,
  Download,
  Edit,
  Trash2,
  ScrollText,
  LayoutDashboard,
} from 'lucide-react';
import { Button } from '@/components/ui/button';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import { ConfirmDialog } from '@/components/ui/confirm-dialog';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Switch } from '@/components/ui/switch';
import { serversActions, isTransitioning, type Server } from '@/stores/servers';
import { useApiError } from '@/lib/api-error';
import { toast } from '@/lib/toast';

type DialogKind = 'stop' | 'restart' | 'update' | 'recreate' | 'delete' | null;

interface ServerActionsProps {
  server: Server;
  // card：图标按钮；detail：带文字的按钮
  variant?: 'card' | 'detail';
  onDeleted?: () => void;
}

export function ServerActions({ server, variant = 'card', onDeleted }: ServerActionsProps) {
  const t = useTranslations('serverActions');
  const router = useRouter();
  const apiError = useApiError();
  const [dialog, setDialog] = useState<DialogKind>(null);
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
    run(() => serversActions.updateGame(server.id, updateGame, updateMods && hasMods), t('updateSent'), t('updateFailed'));
  const remove = () =>
    run(
      async () => {
        await serversActions.deleteServer(server.id, deleteData);
        onDeleted?.();
      },
      t('deleteSuccess'),
      t('deleteFailed'),
    );

  const openDialog = (kind: DialogKind) => {
    setConfirmName('');
    setDeleteData(false);
    setUpdateGame(true);
    setUpdateMods(true);
    setDialog(kind);
  };

  const iconButton = 'h-8 w-8 p-0';
  const isCard = variant === 'card';

  // 主操作始终带文字：卡片和详情页都能一眼看出点下去会发生什么
  let primary;
  if (busy) {
    primary = (
      <Button variant="outline" size="sm" disabled className="min-w-[5.5rem]">
        <Loader2 className="h-4 w-4 animate-spin" />
        {t('working')}
      </Button>
    );
  } else if (running) {
    primary = (
      <Button
        variant="outline"
        size="sm"
        className="min-w-[5.5rem] hover:border-red-200 hover:bg-red-50 hover:text-red-800"
        onClick={() => (players > 0 || !isCard ? openDialog('stop') : stop().catch(() => {}))}
      >
        <Square className="size-2.5 fill-red-600 text-red-600" />
        {t('stop')}
      </Button>
    );
  } else {
    primary = (
      <Button size="sm" className="min-w-[5.5rem]" onClick={start} disabled={!stopped}>
        <Play className="size-3.5 fill-current" />
        {t('start')}
      </Button>
    );
  }

  return (
    <>
      <div className="flex items-center gap-1.5">
        {primary}
        {!isCard && (
          <Button variant="outline" size="sm" onClick={() => openDialog('restart')} disabled={busy}>
            <RotateCw className="h-4 w-4" />
            {t('restart')}
          </Button>
        )}
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <Button variant={isCard ? 'ghost' : 'outline'} size="sm" className={iconButton} aria-label={t('more')} title={t('more')}>
              <MoreVertical className="h-4 w-4" />
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="end" className="w-52">
            {isCard && (
              <>
                <DropdownMenuItem onClick={() => router.push(`/servers/${server.id}`)}>
                  <LayoutDashboard className="mr-2 h-4 w-4" />
                  {t('details')}
                </DropdownMenuItem>
                <DropdownMenuItem onClick={() => router.push(`/servers/${server.id}?tab=logs`)}>
                  <ScrollText className="mr-2 h-4 w-4" />
                  {t('logs')}
                </DropdownMenuItem>
              </>
            )}
            <DropdownMenuItem onClick={() => router.push(`/servers/${server.id}/edit`)}>
              <Edit className="mr-2 h-4 w-4" />
              {t('edit')}
            </DropdownMenuItem>
            <DropdownMenuSeparator />
            {isCard && (
              <DropdownMenuItem onClick={() => openDialog('restart')} disabled={busy}>
                <RotateCw className="mr-2 h-4 w-4" />
                {t('restart')}
              </DropdownMenuItem>
            )}
            <DropdownMenuItem onClick={() => openDialog('update')} disabled={busy}>
              <Download className="mr-2 h-4 w-4" />
              {t('updateGame')}
            </DropdownMenuItem>
            <DropdownMenuItem onClick={() => openDialog('recreate')} disabled={busy}>
              <RotateCcw className="mr-2 h-4 w-4" />
              {t('recreate')}
            </DropdownMenuItem>
            <DropdownMenuSeparator />
            <DropdownMenuItem
              onClick={() => openDialog('delete')}
              disabled={!stopped || pending}
              className="text-red-600 focus:text-red-700"
            >
              <Trash2 className="mr-2 h-4 w-4" />
              {stopped ? t('delete') : t('deleteNeedStop')}
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      </div>

      <ConfirmDialog
        open={dialog === 'stop'}
        onOpenChange={(o) => !o && setDialog(null)}
        title={t('stopTitle', { name: server.session_name })}
        description={players > 0 ? t('stopWithPlayers', { count: players }) : t('stopDesc')}
        confirmText={t('stop')}
        destructive
        onConfirm={stop}
      />

      <ConfirmDialog
        open={dialog === 'restart'}
        onOpenChange={(o) => !o && setDialog(null)}
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
        onOpenChange={(o) => !o && setDialog(null)}
        title={t('updateTitle', { name: server.session_name })}
        description={running ? t('updateDescRunning') : t('updateDesc')}
        confirmText={t('updateConfirm')}
        confirmDisabled={!updateGame && !(updateMods && hasMods)}
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
          <label className="flex items-start justify-between gap-4">
            <span>
              <span className="font-medium">{t('updateModsOption')}</span>
              <span className="block text-xs text-gray-500">{hasMods ? t('updateModsOptionDesc') : t('noMods')}</span>
            </span>
            <Switch checked={updateMods && hasMods} onCheckedChange={setUpdateMods} disabled={!hasMods} />
          </label>
          {players > 0 && <p className="font-medium text-amber-700">{t('playersOnline', { count: players })}</p>}
        </div>
      </ConfirmDialog>

      <ConfirmDialog
        open={dialog === 'recreate'}
        onOpenChange={(o) => !o && setDialog(null)}
        title={t('recreateTitle', { name: server.session_name })}
        description={t('recreateDesc')}
        confirmText={t('recreate')}
        onConfirm={recreate}
      />

      <ConfirmDialog
        open={dialog === 'delete'}
        onOpenChange={(o) => !o && setDialog(null)}
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
}
