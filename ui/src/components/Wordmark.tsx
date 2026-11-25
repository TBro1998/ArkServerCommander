import { cn } from '@/lib/utils';

// 品牌字标：ARKfgf 奖章图标（和 arkfgf.com、联机助手共用），ARK 用 Saira 粗体，副名降一级
export function Wordmark({ className, size = 'md' }: { className?: string; size?: 'md' | 'lg' }) {
  return (
    <span className={cn('font-display inline-flex items-center gap-2.5 leading-none', className)}>
      {/* eslint-disable-next-line @next/next/no-img-element -- 固定尺寸的小图，不需要 next/image 优化 */}
      <img src="/brand/arkfgf-icon-64.png" alt="" className={cn('shrink-0', size === 'lg' ? 'h-12 w-12' : 'h-8 w-8')} />
      <span className="inline-flex items-baseline gap-2">
        <span className={cn('glow-text font-bold tracking-[0.12em] text-primary', size === 'lg' ? 'text-4xl' : 'text-xl')}>ARK</span>
        <span className={cn('font-medium tracking-[0.08em] uppercase opacity-80', size === 'lg' ? 'text-lg' : 'text-xs')}>Server Commander</span>
      </span>
    </span>
  );
}
