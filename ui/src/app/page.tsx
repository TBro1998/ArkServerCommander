"use client";

import { useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { useIsAuthenticated, useAuthActions } from '@/stores/auth';

export default function HomePage() {
  const router = useRouter();
  const isAuthenticated = useIsAuthenticated();
  const { initFromStorage } = useAuthActions();

  useEffect(() => {
    initFromStorage();
  }, [initFromStorage]);
  
  useEffect(() => {
    if (isAuthenticated) {
      router.replace('/servers');
    } else {
      router.replace('/login');
    }
  }, [isAuthenticated, router]);

  return (
    <div className="min-h-screen flex items-center justify-center">
      <p>正在加载...</p>
    </div>
  );
}
