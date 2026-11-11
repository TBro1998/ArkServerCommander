import { cn } from '@/lib/utils';

// 品牌字标：ARK 用 Saira 粗体，副名降一级
export function Wordmark({ className, size = 'md' }: { className?: string; size?: 'md' | 'lg' }) {
  return (
    <span className={cn('font-display inline-flex items-baseline gap-2 leading-none', className)}>
      <span className={cn('font-bold tracking-[0.06em]', size === 'lg' ? 'text-4xl' : 'text-xl')}>ARK</span>
      <span className={cn('font-medium opacity-75', size === 'lg' ? 'text-xl' : 'text-sm')}>Server Commander</span>
    </span>
  );
}
