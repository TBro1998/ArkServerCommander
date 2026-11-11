import { Loader2 } from 'lucide-react';
import { cn } from '@/lib/utils';

// 统一的加载状态：全屏（布局、跳转）或区块内（列表、面板）
export function PageLoader({ text, fullScreen = false, className }: { text?: string; fullScreen?: boolean; className?: string }) {
  return (
    <div
      role="status"
      aria-live="polite"
      className={cn('flex items-center justify-center', fullScreen ? 'min-h-screen bg-background' : 'py-16', className)}
    >
      <div className="flex items-center gap-3 text-sm text-muted-foreground">
        <Loader2 className="h-5 w-5 animate-spin text-primary" />
        {text && <span>{text}</span>}
      </div>
    </div>
  );
}
