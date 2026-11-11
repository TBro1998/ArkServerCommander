import { AlertCircle } from 'lucide-react';
import { LanguageSwitcher } from '@/components/LanguageSwitcher';
import { Wordmark } from '@/components/Wordmark';
import { mapColor } from '@/lib/maps';

// 登录页左侧的色带：与服务器卡片上的地图颜色相同，作为品牌标识
const STRATA = ['TheIsland', 'ScorchedEarth_P', 'Aberration_P', 'Extinction', 'Ragnarok', 'Genesis', 'CrystalIsles', 'Fjordur'];

// 登录 / 初始化页面的外框：宽屏左侧品牌面板，右侧表单
export function AuthShell({ tagline, title, subtitle, children }: { tagline: string; title: string; subtitle?: string; children: React.ReactNode }) {
  return (
    <div className="grid min-h-screen lg:grid-cols-[minmax(0,5fr)_minmax(0,6fr)]">
      <aside className="relative hidden flex-col justify-between overflow-hidden border-r border-primary/25 bg-chrome p-10 text-chrome-foreground lg:flex">
        <div className="space-y-6">
          <Wordmark size="lg" />
          <p className="max-w-sm text-lg leading-relaxed text-chrome-muted">{tagline}</p>
        </div>
        <div aria-hidden className="flex h-2 overflow-hidden rounded-full">
          {STRATA.map((m) => (
            <span key={m} className="flex-1" style={{ backgroundColor: mapColor(m) }} />
          ))}
        </div>
      </aside>

      <main className="flex flex-col px-4 py-6 sm:px-8">
        <div className="flex items-center justify-between lg:justify-end">
          <Wordmark className="lg:hidden" />
          <LanguageSwitcher />
        </div>
        <div className="mx-auto flex w-full max-w-sm flex-1 flex-col justify-center py-10">
          <div className="hud relative rounded-sm border bg-card p-7 backdrop-blur-sm">
            <h1 className="font-display hud-caps text-2xl font-semibold">{title}</h1>
            {subtitle && <p className="mt-1.5 text-sm text-muted-foreground">{subtitle}</p>}
            <div className="mt-7">{children}</div>
          </div>
        </div>
      </main>
    </div>
  );
}

export function FormError({ message }: { message: string }) {
  if (!message) return null;
  return (
    <div role="alert" className="flex items-start gap-2 rounded-sm border border-red-200 bg-red-50 px-3 py-2 text-sm text-red-800">
      <AlertCircle className="mt-0.5 h-4 w-4 flex-shrink-0" />
      <span>{message}</span>
    </div>
  );
}
