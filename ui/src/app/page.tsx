"use client";

import { useEffect } from 'react';
import { useTranslations } from 'next-intl';
import { PageLoader } from '@/components/ui/page-loader';
import { useRouter } from '@/navigation';
import { useIsAuthenticated, useAuthActions, useAuthIsInitialized } from '@/stores/auth';

export default function HomePage() {
  const router = useRouter();
  const isAuthenticated = useIsAuthenticated();
  const isInitialized = useAuthIsInitialized();
  const { restoreSession } = useAuthActions();
  const t = useTranslations('layout');

  useEffect(() => {
    if (!isInitialized) {
      restoreSession();
    }
  }, [restoreSession, isInitialized]);

  useEffect(() => {
    // 只有在初始化完成后才进行跳转
    if (isInitialized) {
      if (isAuthenticated) {
        router.replace('/home');
      } else {
        router.replace('/login');
      }
    }
  }, [isAuthenticated, isInitialized, router]);

  return <PageLoader fullScreen text={isInitialized ? t('redirecting') : t('initializing')} />;
}
