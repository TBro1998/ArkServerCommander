"use client";

import { useState, type ReactNode } from 'react';
import { useTranslations } from 'next-intl';
import { Loader2 } from 'lucide-react';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog';
import { Button } from '@/components/ui/button';

interface ConfirmDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  title: string;
  description?: ReactNode;
  children?: ReactNode;
  confirmText?: string;
  destructive?: boolean;
  confirmDisabled?: boolean;
  // 返回 Promise 时，执行期间按钮显示加载中；成功后自动关闭，抛错时保持打开
  onConfirm: () => unknown;
}

// 通用确认弹窗
export function ConfirmDialog({
  open,
  onOpenChange,
  title,
  description,
  children,
  confirmText,
  destructive,
  confirmDisabled,
  onConfirm,
}: ConfirmDialogProps) {
  const tCommon = useTranslations('common');
  const [busy, setBusy] = useState(false);

  const handleConfirm = async () => {
    setBusy(true);
    try {
      await onConfirm();
      onOpenChange(false);
    } catch {
      // 调用方负责提示错误
    } finally {
      setBusy(false);
    }
  };

  return (
    <Dialog open={open} onOpenChange={(next) => !busy && onOpenChange(next)}>
      <DialogContent className="max-w-md">
        <DialogHeader>
          <DialogTitle>{title}</DialogTitle>
          {description && <DialogDescription asChild><div>{description}</div></DialogDescription>}
        </DialogHeader>
        {children}
        <DialogFooter>
          <Button variant="outline" onClick={() => onOpenChange(false)} disabled={busy}>
            {tCommon('cancel')}
          </Button>
          <Button variant={destructive ? 'destructive' : 'default'} onClick={handleConfirm} disabled={busy || confirmDisabled}>
            {busy && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
            {confirmText ?? tCommon('confirm')}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
