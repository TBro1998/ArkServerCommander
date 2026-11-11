"use client";

import { useEffect, useState } from 'react';
import { useTranslations } from 'next-intl';
import { Menu, LogOut, CircleUser } from 'lucide-react';
import { Link, usePathname, useRouter } from '@/navigation';
import { useIsAuthenticated, useAuthActions, useAuthIsInitialized, useAuthUser } from '@/stores/auth';
import { LanguageSwitcher } from '@/components/LanguageSwitcher';
import { Wordmark } from '@/components/Wordmark';
import { Button } from '@/components/ui/button';
import { PageLoader } from '@/components/ui/page-loader';
import { Sheet, SheetContent, SheetHeader, SheetTitle, SheetTrigger } from '@/components/ui/sheet';
import { cn } from '@/lib/utils';

const NAV_ITEMS = [
  { href: '/home', key: 'home', match: (p: string) => p === '/home' },
  { href: '/servers', key: 'servers', match: (p: string) => p.startsWith('/servers') },
  { href: '/settings', key: 'settings', match: (p: string) => p.startsWith('/settings') },
] as const;

export default function ProtectedLayout({ children }: { children: React.ReactNode }) {
  const router = useRouter();
  const isAuthenticated = useIsAuthenticated();
  const isInitialized = useAuthIsInitialized();
  const user = useAuthUser();
  const { restoreSession, logout } = useAuthActions();
  const t = useTranslations('navigation');
  const tLayout = useTranslations('layout');
  const pathname = usePathname();
  const [menuOpen, setMenuOpen] = useState(false);

  useEffect(() => {
    if (!isInitialized) {
      restoreSession();
    }
  }, [restoreSession, isInitialized]);

  useEffect(() => {
    // 只有在初始化完成且未认证时才跳转
    if (isInitialized && !isAuthenticated) {
      router.replace('/login');
    }
  }, [isAuthenticated, isInitialized, router]);

  if (!isInitialized) return <PageLoader fullScreen text={tLayout('initializing')} />;
  // 初始化完成但未认证，很快会跳转
  if (!isAuthenticated) return <PageLoader fullScreen text={tLayout('verifying')} />;

  const headerButton = 'text-chrome-muted hover:bg-white/10 hover:text-chrome-foreground';

  return (
    <div className="min-h-screen bg-background">
      <a
        href="#main"
        className="sr-only focus:not-sr-only focus:fixed focus:top-2 focus:left-2 focus:z-50 focus:rounded-md focus:bg-card focus:px-3 focus:py-2 focus:text-sm"
      >
        {tLayout('skipToContent')}
      </a>
      <header className="sticky top-0 z-40 border-b border-primary/25 bg-chrome text-chrome-foreground shadow-[0_1px_16px_-6px] shadow-primary/40 backdrop-blur">
        <div className="mx-auto flex h-14 max-w-[1600px] items-center justify-between px-4 sm:px-6 lg:px-8">
          <div className="flex h-full min-w-0 items-center gap-8">
            <Link href="/home" className="truncate text-chrome-foreground">
              <Wordmark />
            </Link>
            <nav className="hidden h-full items-stretch gap-1 md:flex" aria-label={tLayout('menu')}>
              {NAV_ITEMS.map((item) => {
                const active = item.match(pathname);
                return (
                  <Link
                    key={item.href}
                    href={item.href}
                    aria-current={active ? 'page' : undefined}
                    className={cn(
                      'hud-caps relative flex items-center px-3 text-sm font-medium transition-colors',
                      active ? 'glow-text text-primary' : 'text-chrome-muted hover:text-chrome-foreground',
                    )}
                  >
                    {t(item.key)}
                    {active && <span className="absolute inset-x-3 -bottom-px h-0.5 bg-primary shadow-[0_0_8px] shadow-primary" />}
                  </Link>
                );
              })}
            </nav>
          </div>
          <div className="hidden items-center gap-1 md:flex">
            {user && (
              <span className="mr-2 flex items-center gap-1.5 text-sm text-chrome-muted">
                <CircleUser className="h-4 w-4" />
                {user.username}
              </span>
            )}
            <LanguageSwitcher className={headerButton} />
            <Button variant="ghost" size="sm" onClick={logout} className={headerButton}>
              <LogOut className="h-4 w-4" />
              {t('logout')}
            </Button>
          </div>

          {/* 手机端菜单 */}
          <div className="flex items-center gap-1 md:hidden">
            <LanguageSwitcher className={headerButton} />
            <Sheet open={menuOpen} onOpenChange={setMenuOpen}>
              <SheetTrigger asChild>
                <Button variant="ghost" size="icon" aria-label={tLayout('menu')} className={headerButton}>
                  <Menu className="h-5 w-5" />
                </Button>
              </SheetTrigger>
              <SheetContent side="right" className="w-64">
                <SheetHeader>
                  <SheetTitle className="flex items-center gap-2">
                    <CircleUser className="h-4 w-4" />
                    {user?.username ?? 'ARK Server Commander'}
                  </SheetTitle>
                </SheetHeader>
                <nav className="space-y-1 px-4">
                  {NAV_ITEMS.map((item) => {
                    const active = item.match(pathname);
                    return (
                      <Link
                        key={item.href}
                        href={item.href}
                        onClick={() => setMenuOpen(false)}
                        aria-current={active ? 'page' : undefined}
                        className={cn(
                          'block rounded-md border-l-[3px] px-3 py-2.5 text-base font-medium',
                          active ? 'border-primary bg-accent text-primary' : 'border-transparent text-muted-foreground hover:bg-accent',
                        )}
                      >
                        {t(item.key)}
                      </Link>
                    );
                  })}
                </nav>
                <div className="mt-4 border-t px-4 pt-4">
                  <Button
                    variant="ghost"
                    className="w-full justify-start"
                    onClick={() => {
                      setMenuOpen(false);
                      logout();
                    }}
                  >
                    <LogOut className="mr-2 h-4 w-4" />
                    {t('logout')}
                  </Button>
                </div>
              </SheetContent>
            </Sheet>
          </div>
        </div>
      </header>
      <main id="main" className="mx-auto w-full max-w-[1600px] px-4 pt-6 pb-16 sm:px-6 sm:pt-8 lg:px-8">
        {children}
      </main>
    </div>
  );
}
