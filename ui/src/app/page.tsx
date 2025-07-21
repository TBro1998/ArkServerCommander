"use client";

import { useEffect } from 'react';
import { useRouter } from '@/navigation';
import { useIsAuthenticated, useAuthActions } from '@/stores/auth';

export default function HomePage() {
  const router = useRouter();
  const isAuthenticated = useIsAuthenticated();
  const { initFromStorage } = useAuthActions();

  useEffect(() => {
    initFromStorage();
  }, [initFromStorage]);

  useEffect(() => {
    if (isAuthenticated === false) {
      router.replace('/login');
    }
  }, [isAuthenticated, router]);

  // 如果已认证，直接显示受保护的首页内容，但不包含布局
  // 这样用户会看到首页内容，但没有导航菜单
  // 我们需要修改这个逻辑来包含导航菜单
  if (isAuthenticated === true) {
    // 重定向到 /home
    router.replace('/home');
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="text-center">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-gray-900 mx-auto mb-4"></div>
          <p>正在跳转...</p>
        </div>
      </div>
    );
  }

  // 显示加载状态
  return (
    <div className="min-h-screen flex items-center justify-center">
      <div className="text-center">
        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-gray-900 mx-auto mb-4"></div>
        <p>正在加载...</p>
      </div>
    </div>
  );
}