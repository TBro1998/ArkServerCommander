"use client";

import { useState, useEffect } from 'react';
import { useRouter } from '@/navigation';
import { useLocale, useTranslations } from 'next-intl';
import { Button } from '@/components/ui/button';
import { Loader2 } from 'lucide-react';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { PageLoader } from '@/components/ui/page-loader';
import { AuthShell, FormError } from '@/components/AuthShell';
import { useAuthActions, useIsAuthenticated } from '@/stores/auth';

export default function InitPage() {
  const t = useTranslations('auth');
  const tApi = useTranslations('apiErrors');
  const locale = useLocale();
  const router = useRouter();
  const { init, checkInit } = useAuthActions();
  const isAuthenticated = useIsAuthenticated();

  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState('');
  const [checkingInit, setCheckingInit] = useState(true);

  useEffect(() => {
    if (isAuthenticated) {
      router.replace('/home');
      return;
    }

    const checkInitStatus = async () => {
      try {
        const isInitialized = await checkInit();
        if (isInitialized) { // 如果已经初始化，跳转到登录页面
          router.replace('/login');
          return;
        }
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
      setError(t('enterAdminUsername'));
      return;
    }

    if (!password.trim()) {
      setError(t('enterPassword'));
      return;
    }

    if (password.length < 8) {
      setError(t('passwordMinLengthError'));
      return;
    }

    // 后端限制 72 字节（bcrypt 只处理前 72 字节），中文等多字节字符要按字节数计算
    if (new TextEncoder().encode(password).length > 72) {
      setError(t('passwordMaxLengthError'));
      return;
    }

    if (password !== confirmPassword) {
      setError(t('passwordMismatch'));
      return;
    }

    setIsLoading(true);
    setError('');

    try {
      const result = await init({ username, password });
      if (result.success) {
        router.replace('/home');
      } else {
        // 中文界面显示后端原文；其他语言按错误码显示翻译
        setError(result.code && locale !== 'zh' && tApi.has(result.code) ? tApi(result.code) : result.message || t('initError'));
      }
    } catch (error: unknown) {
      const errorMessage = error instanceof Error ? error.message : t('initError');
      setError(errorMessage);
    } finally {
      setIsLoading(false);
    }
  };

  if (checkingInit) return <PageLoader fullScreen text={t('initCheck')} />;

  return (
    <AuthShell tagline={t('tagline')} title={t('initTitle')} subtitle={t('initSubtitle')}>
      <form onSubmit={handleSubmit} className="space-y-5" noValidate>
        <div className="space-y-2">
          <Label htmlFor="username">{t('adminUsername')}</Label>
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
            autoComplete="new-password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            disabled={isLoading}
            aria-describedby="password-hint"
            className="h-10"
          />
          <p id="password-hint" className="text-xs text-muted-foreground">{t('passwordMinLength')}</p>
        </div>

        <div className="space-y-2">
          <Label htmlFor="confirmPassword">{t('confirmPassword')}</Label>
          <Input
            id="confirmPassword"
            type="password"
            autoComplete="new-password"
            value={confirmPassword}
            onChange={(e) => setConfirmPassword(e.target.value)}
            disabled={isLoading}
            className="h-10"
          />
        </div>

        <FormError message={error} />

        <Button type="submit" className="h-10 w-full" disabled={isLoading}>
          {isLoading && <Loader2 className="h-4 w-4 animate-spin" />}
          {isLoading ? t('initLoading') : t('initButton')}
        </Button>
        <p className="text-center text-xs text-muted-foreground">{t('initTip')}</p>
      </form>
    </AuthShell>
  );
}