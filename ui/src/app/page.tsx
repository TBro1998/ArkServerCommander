"use client";

import { useEffect } from 'react';
import { useRouter } from '@/navigation';
import { useIsAuthenticated, useAuthActions, useAuthIsInitialized } from '@/stores/auth';

export default function HomePage() {
  const router = useRouter();
  const isAuthenticated = useIsAuthenticated();
  const isInitialized = useAuthIsInitialized();
  const { initFromStorage } = useAuthActions();

  useEffect(() => {
    if (!isInitialized) {
      initFromStorage();
    }
  }, [initFromStorage, isInitialized]);

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

  // 显示加载状态
  return (
    <div className="min-h-screen flex items-center justify-center">
      <div className="text-center">
        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-gray-900 mx-auto mb-4"></div>
        <p>{isInitialized ? '正在跳转...' : '正在加载...'}</p>
      </div>
    </div>
  );
}