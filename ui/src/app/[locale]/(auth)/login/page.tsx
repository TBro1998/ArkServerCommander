"use client";

import { useEffect, useState } from 'react';
import { useAuthActions, useAuthIsLoading, useIsAuthenticated } from '@/stores/auth';
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";
import { Server, User, Lock, LogIn, AlertCircle } from 'lucide-react';
import { useTranslations } from 'next-intl';
import { useRouter } from '@/navigation';

export default function LoginPage() {
  const router = useRouter();
  const t = useTranslations('auth');
  const { login, checkInit } = useAuthActions();
  const isAuthenticated = useIsAuthenticated();
  const isLoading = useAuthIsLoading();

  const [form, setForm] = useState({ username: '', password: '' });
  const [error, setError] = useState('');

  useEffect(() => {
    const initialize = async () => {
      const initialized = await checkInit();
      if (!initialized) {
        router.push('/init');
      } else if (isAuthenticated) {
        router.push('/');
      }
    };
    initialize();
  }, [checkInit, isAuthenticated, router]);

  const handleLogin = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    const result = await login(form);
    if (result.success) {
      router.push('/');
    } else {
      setError(result.message);
    }
  };

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setForm({ ...form, [e.target.name]: e.target.value });
  };

  return (
    <div className="min-h-screen flex flex-col bg-gradient-to-br from-blue-50 to-indigo-100">
      <div className="flex-1 flex items-center justify-center py-12 px-4 sm:px-6 lg:px-8">
        <div className="w-full max-w-md">
          <div className="text-center mb-8">
            <div className="mx-auto w-16 h-16 bg-blue-600 rounded-full flex items-center justify-center mb-4">
              <Server className="w-8 h-8 text-white" />
            </div>
            <h1 className="text-3xl font-bold text-gray-900 mb-2">
              {t('loginTitle')}
            </h1>
            <p className="text-gray-600">
              {t('loginSubtitle')}
            </p>
          </div>

          <Card className="shadow-xl">
            <CardHeader>
              <CardTitle>{t('login')}</CardTitle>
              <CardDescription>{t('enterUsername')}</CardDescription>
            </CardHeader>
            <CardContent>
              <form onSubmit={handleLogin} className="space-y-6">
                {error && (
                  <Alert variant="destructive">
                    <AlertCircle className="h-4 w-4" />
                    <AlertTitle>{t('loginError')}</AlertTitle>
                    <AlertDescription>{error}</AlertDescription>
                  </Alert>
                )}
                <div className="space-y-2">
                  <Label htmlFor="username">{t('username')}</Label>
                  <div className="relative">
                    <User className="absolute left-3 top-1/2 -translate-y-1/2 h-5 w-5 text-gray-400" />
                    <Input
                      id="username"
                      name="username"
                      value={form.username}
                      onChange={handleChange}
                      placeholder={t('enterUsername')}
                      required
                      className="pl-10"
                    />
                  </div>
                </div>
                <div className="space-y-2">
                  <Label htmlFor="password">{t('password')}</Label>
                  <div className="relative">
                    <Lock className="absolute left-3 top-1/2 -translate-y-1/2 h-5 w-5 text-gray-400" />
                    <Input
                      id="password"
                      name="password"
                      type="password"
                      value={form.password}
                      onChange={handleChange}
                      placeholder={t('enterPassword')}
                      required
                      className="pl-10"
                    />
                  </div>
                </div>
                <Button type="submit" disabled={isLoading} className="w-full mt-8">
                  {isLoading ? (
                    t('loginLoading')
                  ) : (
                    <>
                      <LogIn className="w-4 h-4 mr-2" />
                      {t('loginButton')}
                    </>
                  )}
                </Button>
              </form>
            </CardContent>
          </Card>

          <div className="text-center mt-6">
            <p className="text-sm text-gray-500">
              {t('firstTimeTip')}
            </p>
          </div>
        </div>
      </div>
    </div>
  );
}