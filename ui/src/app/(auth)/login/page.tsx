"use client";

import { useState, useEffect } from 'react';
import { useRouter } from '@/navigation';
import { useLocale, useTranslations } from 'next-intl';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { PageLoader } from '@/components/ui/page-loader';
import { AuthShell, FormError } from '@/components/AuthShell';
import { useAuthActions, useIsAuthenticated } from '@/stores/auth';
import { Loader2 } from 'lucide-react';

export default function LoginPage() {
  const t = useTranslations('auth');
  const tApi = useTranslations('apiErrors');
  const locale = useLocale();
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
        // 中文界面显示后端原文；其他语言按错误码显示翻译
        setError(result.code && locale !== 'zh' && tApi.has(result.code) ? tApi(result.code) : result.message || t('loginError'));
      }
    } catch (error: unknown) {
      const errorMessage = error instanceof Error ? error.message : t('loginError');
      setError(errorMessage);
    } finally {
      setIsLoading(false);
    }
  };

  if (checkingInit) return <PageLoader fullScreen text={t('initCheck')} />;

  if (needsInit) {
    router.replace('/init');
    return null;
  }

  return (
    <AuthShell tagline={t('tagline')} title={t('login')} subtitle={t('loginSubtitle')}>
      <form onSubmit={handleSubmit} className="space-y-5" noValidate>
        <div className="space-y-2">
          <Label htmlFor="username">{t('username')}</Label>
          <Input
            id="username"
            type="text"
            autoComplete="username"
            autoFocus
            value={username}
            onChange={(e) => setUsername(e.target.value)}
            disabled={isLoading}
            className="h-10"
          />
        </div>

        <div className="space-y-2">
          <Label htmlFor="password">{t('password')}</Label>
          <Input
            id="password"
            type="password"
            autoComplete="current-password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            disabled={isLoading}
            className="h-10"
          />
        </div>

        <FormError message={error} />

        <Button type="submit" className="h-10 w-full" disabled={isLoading}>
          {isLoading && <Loader2 className="h-4 w-4 animate-spin" />}
          {isLoading ? t('loginLoading') : t('loginButton')}
        </Button>
      </form>
    </AuthShell>
  );
}