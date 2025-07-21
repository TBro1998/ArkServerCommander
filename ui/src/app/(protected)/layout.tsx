"use client";

import { useEffect } from 'react';
import { useTranslations } from 'next-intl';
import { Link, usePathname, useRouter } from '@/navigation';
import { useIsAuthenticated, useAuthActions } from '@/stores/auth';
import { LanguageSwitcher } from '@/components/LanguageSwitcher';

export default function ProtectedLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const router = useRouter();
  const isAuthenticated = useIsAuthenticated();
  const { initFromStorage, logout } = useAuthActions();
  const t = useTranslations('navigation');
  const pathname = usePathname();

  useEffect(() => {
    initFromStorage();
  }, [initFromStorage]);

  useEffect(() => {
    if (!isAuthenticated) {
      router.replace('/login');
    }
  }, [isAuthenticated, router]);

  if (!isAuthenticated) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="text-center">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-gray-900 mx-auto mb-4"></div>
          <p>验证身份中...</p>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50">
      <header className="bg-white shadow-sm border-b">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between items-center h-16">
            <div className="flex items-center">
              <h1 className="text-xl font-semibold text-gray-900">
                ARK Server Manager
              </h1>
              <nav className="ml-10 flex items-baseline space-x-4">
                <Link
                  href="/home"
                  className={`px-3 py-2 rounded-md text-sm font-medium ${pathname === '/home'
                      ? 'bg-gray-200 text-gray-900'
                      : 'text-gray-500 hover:bg-gray-100 hover:text-gray-900'
                    }`}
                >
                  {t('home')}
                </Link>
                <Link
                  href="/servers"
                  className={`px-3 py-2 rounded-md text-sm font-medium ${pathname.startsWith('/servers')
                      ? 'bg-gray-200 text-gray-900'
                      : 'text-gray-500 hover:bg-gray-100 hover:text-gray-900'
                    }`}
                >
                  {t('servers')}
                </Link>
              </nav>
            </div>
            <div className="flex items-center space-x-4">
              <LanguageSwitcher />
              <button
                onClick={logout}
                className="text-gray-500 hover:text-gray-700 px-3 py-2 rounded-md text-sm font-medium"
              >
                退出登录
              </button>
            </div>
          </div>
        </div>
      </header>
      <main className="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
        {children}
      </main>
    </div>
  );
}