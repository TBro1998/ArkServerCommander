"use client";

import { CheckCircle2, AlertCircle, Info, X } from 'lucide-react';
import { useToastStore } from '@/lib/toast';
import { cn } from '@/lib/utils';

const ICONS = {
  success: CheckCircle2,
  error: AlertCircle,
  info: Info,
};

const STYLES = {
  success: 'border-green-200 bg-green-50 text-green-900',
  error: 'border-red-200 bg-red-50 text-red-900',
  info: 'border-blue-200 bg-blue-50 text-blue-900',
};

// 全局通知：固定在右下角（手机上为底部），不随页面滚动
export function Toaster() {
  const toasts = useToastStore((s) => s.toasts);
  const dismiss = useToastStore((s) => s.dismiss);

  return (
    <div
      aria-live="polite"
      className="pointer-events-none fixed inset-x-0 bottom-0 z-[100] flex flex-col items-center gap-2 p-4 sm:inset-x-auto sm:right-0 sm:items-end"
    >
      {toasts.map((t) => {
        const Icon = ICONS[t.variant];
        return (
          <div
            key={t.id}
            role={t.variant === 'error' ? 'alert' : 'status'}
            className={cn(
              'pointer-events-auto flex w-full max-w-sm items-start gap-3 rounded-lg border p-3 shadow-lg animate-in fade-in slide-in-from-bottom-2',
              STYLES[t.variant],
            )}
          >
            <Icon className="mt-0.5 h-4 w-4 flex-shrink-0" />
            <div className="min-w-0 flex-1 text-sm">
              {t.title && <div className="font-medium">{t.title}</div>}
              <div className="break-words">{t.message}</div>
            </div>
            <button
              type="button"
              onClick={() => dismiss(t.id)}
              className="-m-1 rounded p-1 opacity-60 hover:opacity-100"
              aria-label="close"
            >
              <X className="h-4 w-4" />
            </button>
          </div>
        );
      })}
    </div>
  );
}
