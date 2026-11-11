import { cn } from '@/lib/utils';

// 品牌字标：ARK 用 Saira 粗体，副名降一级
export function Wordmark({ className, size = 'md' }: { className?: string; size?: 'md' | 'lg' }) {
  return (
    <span className={cn('font-display inline-flex items-baseline gap-2 leading-none', className)}>
      <span className={cn('glow-text font-bold tracking-[0.12em] text-primary', size === 'lg' ? 'text-4xl' : 'text-xl')}>ARK</span>
      <span className={cn('font-medium tracking-[0.08em] uppercase opacity-80', size === 'lg' ? 'text-lg' : 'text-xs')}>Server Commander</span>
    </span>
  );
}
