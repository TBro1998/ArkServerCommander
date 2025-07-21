"use client";

import { useState, useEffect } from 'react';
import { useRouter } from '@/navigation';
import { useTranslations } from 'next-intl';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { useAuthActions, useIsAuthenticated } from '@/stores/auth';
import { LanguageSwitcher } from '@/components/LanguageSwitcher';

export default function LoginPage() {
  const t = useTranslations('auth');
  const router = useRouter();
  const { login, checkInit } = useAuthActions();
  const isAuthenticated = useIsAuthenticated();

  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState('');
  const [needsInit, setNeedsInit] = useState(false);
  const [checkingInit, setCheckingInit] = useState(true);

  useEffect(() => {
    if (isAuthenticated) {
      router.replace('/home');
      return;
    }

    const checkInitStatus = async () => {
      try {
        const isInitialized = await checkInit();
        setNeedsInit(!isInitialized); // 如果未初始化，则需要初始化
      } catch (error) {
        console.error('Failed to check init status:', error);
      } finally {
        setCheckingInit(false);
      }
    };

    checkInitStatus();
  }, [isAuthenticated, router, checkInit]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!username.trim()) {
      setError(t('enterUsername'));
      return;
    }

    if (!password.trim()) {
      setError(t('enterPassword'));
      return;
    }

    setIsLoading(true);
    setError('');

    try {
      const result = await login({ username, password });
      if (result.success) {
        router.replace('/home');
      } else {
        setError(result.message || t('loginError'));
      }
    } catch (error: unknown) {
      const errorMessage = error instanceof Error ? error.message : t('loginError');
      setError(errorMessage);
    } finally {
      setIsLoading(false);
    }
  };

  if (checkingInit) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="text-center">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-gray-900 mx-auto mb-4"></div>
          <p>{t('initCheck')}</p>
        </div>
      </div>
    );
  }

  if (needsInit) {
    router.replace('/init');
    return null;
  }

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
      <div className="max-w-md w-full space-y-8">
        <div className="flex justify-end">
          <LanguageSwitcher />
        </div>

        <Card>
          <CardHeader className="space-y-1">
            <CardTitle className="text-2xl text-center">{t('loginTitle')}</CardTitle>
            <CardDescription className="text-center">
              {t('loginSubtitle')}
            </CardDescription>
          </CardHeader>
          <CardContent>
            <form onSubmit={handleSubmit} className="space-y-4">
              <div className="space-y-2">
                <Label htmlFor="username">{t('username')}</Label>
                <Input
                  id="username"
                  type="text"
                  value={username}
                  onChange={(e) => setUsername(e.target.value)}
                  placeholder={t('enterUsername')}
                  disabled={isLoading}
                />
              </div>

              <div className="space-y-2">
                <Label htmlFor="password">{t('password')}</Label>
                <Input
                  id="password"
                  type="password"
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  placeholder={t('enterPassword')}
                  disabled={isLoading}
                />
              </div>

              {error && (
                <div className="text-red-600 text-sm text-center">
                  {error}
                </div>
              )}

              <Button
                type="submit"
                className="w-full"
                disabled={isLoading}
              >
                {isLoading ? t('loginLoading') : t('loginButton')}
              </Button>
            </form>

            <div className="mt-4 text-center text-sm text-gray-600">
              {t('firstTimeTip')}
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}