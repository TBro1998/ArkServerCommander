"use client";

import { useState } from 'react';
import { useTranslations } from 'next-intl';
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
  DialogFooter,
} from '@/components/ui/dialog';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { AlertTriangle, Loader2 } from 'lucide-react';
import { Server } from '@/stores/servers';

interface ImageUpdateConfirmModalProps {
  open: boolean;
  imageName: string;
  imageLabel: string;
  affectedServers: Server[];
  submitting: boolean;
  onCancel: () => void;
  onConfirm: (restartRunning: boolean) => void;
}

// 父组件用 key 在每次打开时重新挂载，选项会回到默认的"仅更新镜像"
export function ImageUpdateConfirmModal({
  open,
  imageName,
  imageLabel,
  affectedServers,
  submitting,
  onCancel,
  onConfirm,
}: ImageUpdateConfirmModalProps) {
  const t = useTranslations('servers.dockerImages');
  const tCommon = useTranslations('common');
  const [restartRunning, setRestartRunning] = useState(false);

  const runningCount = affectedServers.filter((s) => s.status === 'running').length;

  const options = [
    { value: false, title: t('updateImageOnly'), desc: t('updateImageOnlyDesc') },
    { value: true, title: t('updateAndRecreate'), desc: t('updateAndRecreateDesc') },
  ];

  return (
    <Dialog open={open} onOpenChange={(next) => !next && !submitting && onCancel()}>
      <DialogContent className="max-w-lg">
        <DialogHeader>
          <DialogTitle>{t('updateConfirm')}</DialogTitle>
          <DialogDescription>
            {imageLabel} <span className="font-mono">({imageName})</span>
          </DialogDescription>
        </DialogHeader>

        {/* 受影响的服务器 */}
        <div className="space-y-2 text-sm">
          <div className="font-medium">
            {t('affectedServers')}
            {runningCount > 0 && (
              <span className="ml-2 text-xs text-yellow-700">{t('runningCount', { count: runningCount })}</span>
            )}
          </div>
          {affectedServers.length === 0 ? (
            <p className="text-gray-500">{t('noAffectedServers')}</p>
          ) : (
            <ul className="max-h-40 overflow-y-auto space-y-1">
              {affectedServers.map((server) => (
                <li key={server.id} className="flex items-center justify-between rounded border px-3 py-1.5">
                  <span className="truncate">{server.session_name}</span>
                  <Badge variant={server.status === 'running' ? 'default' : 'secondary'} className="text-xs">
                    {server.status}
                  </Badge>
                </li>
              ))}
            </ul>
          )}
        </div>

        {/* 更新须知 */}
        <div className="rounded-lg border border-yellow-200 bg-yellow-50 p-3 text-sm text-yellow-800 space-y-1">
          <div className="flex items-center font-medium">
            <AlertTriangle className="w-4 h-4 mr-1.5" />
            {t('updateWarning')}
          </div>
          <ul className="list-disc pl-5 space-y-0.5">
            <li>{t('warningDownloadTime')}</li>
            <li>{t('warningContainerRecreate')}</li>
            <li>{t('warningDataSafety')}</li>
          </ul>
        </div>

        {/* 更新选项：只有存在运行中的服务器时才需要选择 */}
        {runningCount > 0 && (
          <div className="space-y-2 text-sm">
            <div className="font-medium">{t('updateOptions')}</div>
            {options.map((option) => (
              <label
                key={String(option.value)}
                className={`flex gap-3 rounded-lg border p-3 cursor-pointer ${
                  restartRunning === option.value ? 'border-blue-500 bg-blue-50' : 'hover:bg-gray-50'
                }`}
              >
                <input
                  type="radio"
                  name="update-option"
                  className="mt-1"
                  checked={restartRunning === option.value}
                  onChange={() => setRestartRunning(option.value)}
                />
                <div>
                  <div className="font-medium">{option.title}</div>
                  <div className="text-gray-600">{option.desc}</div>
                </div>
              </label>
            ))}
          </div>
        )}

        <DialogFooter>
          <Button variant="outline" onClick={onCancel} disabled={submitting}>
            {tCommon('cancel')}
          </Button>
          <Button onClick={() => onConfirm(runningCount > 0 && restartRunning)} disabled={submitting}>
            {submitting && <Loader2 className="w-4 h-4 mr-2 animate-spin" />}
            {t('confirmUpdate')}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
