"use client";

import { useEffect } from 'react';
import { useRouter } from '@/navigation';
import { useIsAuthenticated, useAuthActions } from '@/stores/auth';
import ProtectedHomePage from './(protected)/page';

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

  if (isAuthenticated === null) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <p>正在加载...</p>
      </div>
    );
  }

  if (isAuthenticated === true) {
    return <ProtectedHomePage />;
  }

  return null;
}