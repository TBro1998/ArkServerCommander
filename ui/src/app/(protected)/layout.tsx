"use client";

import { useEffect, useState } from 'react';
import { useTranslations } from 'next-intl';
import { Menu, LogOut } from 'lucide-react';
import { Link, usePathname, useRouter } from '@/navigation';
import { useIsAuthenticated, useAuthActions, useAuthIsInitialized, useAuthUser } from '@/stores/auth';
import { LanguageSwitcher } from '@/components/LanguageSwitcher';
import { Button } from '@/components/ui/button';
import { Sheet, SheetContent, SheetHeader, SheetTitle, SheetTrigger } from '@/components/ui/sheet';
import { cn } from '@/lib/utils';

const NAV_ITEMS = [
  { href: '/home', key: 'home', match: (p: string) => p === '/home' },
  { href: '/servers', key: 'servers', match: (p: string) => p.startsWith('/servers') },
  { href: '/settings', key: 'settings', match: (p: string) => p.startsWith('/settings') },
] as const;

function Spinner({ text }: { text: string }) {
  return (
    <div className="flex min-h-screen items-center justify-center">
      <div className="text-center">
        <div className="mx-auto mb-4 h-8 w-8 animate-spin rounded-full border-b-2 border-gray-900"></div>
        <p>{text}</p>
      </div>
    </div>
  );
}

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

  if (!isInitialized) return <Spinner text={tLayout('initializing')} />;
  // 初始化完成但未认证，很快会跳转
  if (!isAuthenticated) return <Spinner text={tLayout('verifying')} />;

  const navLink = (item: (typeof NAV_ITEMS)[number], mobile = false) => (
    <Link
      key={item.href}
      href={item.href}
      onClick={() => setMenuOpen(false)}
      className={cn(
        'rounded-md px-3 py-2 text-sm font-medium',
        mobile && 'block text-base',
        item.match(pathname) ? 'bg-gray-200 text-gray-900' : 'text-gray-500 hover:bg-gray-100 hover:text-gray-900',
      )}
    >
      {t(item.key)}
    </Link>
  );

  return (
    <div className="min-h-screen bg-gray-50">
      <header className="sticky top-0 z-40 border-b bg-white shadow-sm">
        <div className="w-full px-4 sm:px-6 lg:px-8">
          <div className="flex h-14 items-center justify-between sm:h-16">
            <div className="flex min-w-0 items-center">
              <Link href="/home" className="truncate text-lg font-semibold text-gray-900 sm:text-xl">
                ARK Server Commander
              </Link>
              <nav className="ml-8 hidden items-baseline space-x-2 md:flex">{NAV_ITEMS.map((item) => navLink(item))}</nav>
            </div>
            <div className="hidden items-center space-x-2 md:flex">
              {user && <span className="text-sm text-gray-500">{user.username}</span>}
              <LanguageSwitcher />
              <Button variant="ghost" size="sm" onClick={logout} className="text-gray-500 hover:text-gray-700">
                <LogOut className="mr-1 h-4 w-4" />
                {t('logout')}
              </Button>
            </div>

            {/* 手机端菜单 */}
            <div className="flex items-center gap-1 md:hidden">
              <LanguageSwitcher />
              <Sheet open={menuOpen} onOpenChange={setMenuOpen}>
                <SheetTrigger asChild>
                  <Button variant="ghost" size="icon" aria-label={tLayout('menu')}>
                    <Menu className="h-5 w-5" />
                  </Button>
                </SheetTrigger>
                <SheetContent side="right" className="w-64">
                  <SheetHeader>
                    <SheetTitle>{user?.username ?? 'ARK Server Commander'}</SheetTitle>
                  </SheetHeader>
                  <nav className="space-y-1 px-4">{NAV_ITEMS.map((item) => navLink(item, true))}</nav>
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
        </div>
      </header>
      <main className="w-full px-4 py-2 sm:px-6 sm:py-6 lg:px-8">{children}</main>
    </div>
  );
}

