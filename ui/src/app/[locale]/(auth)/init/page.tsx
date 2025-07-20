"use client";

import { useEffect, useState } from 'react';
import { useAuthActions, useAuthIsLoading } from '@/stores/auth';
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { ClosableAlert } from "@/components/ui/closable-alert";
import { Settings, User, Lock, Info } from 'lucide-react';
import { useTranslations } from 'next-intl';
import { useRouter } from '@/navigation';

export default function InitPage() {
  const router = useRouter();
  const t = useTranslations('auth');
  const { init, checkInit } = useAuthActions();
  const isLoading = useAuthIsLoading();

  const [form, setForm] = useState({ username: '', password: '', confirmPassword: '' });
  const [error, setError] = useState('');
  const [success, setSuccess] = useState('');

  useEffect(() => {
    const initialize = async () => {
      const initialized = await checkInit();
      if (initialized) {
        router.push('/login');
      }
    };
    initialize();
  }, [checkInit, router]);

  const handleInit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setSuccess('');

    if (form.password !== form.confirmPassword) {
      setError(t('passwordMismatch'));
      return;
    }
    if (form.password.length < 6) {
      setError(t('passwordMinLengthError'));
      return;
    }

    const result = await init({ username: form.username, password: form.password });
    if (result.success) {
      setSuccess(result.message);
      setTimeout(() => router.push('/'), 1500);
    } else {
      setError(result.message);
    }
  };

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setForm({ ...form, [e.target.name]: e.target.value });
  };

  return (
    <div className="min-h-screen flex flex-col bg-gradient-to-br from-green-50 to-emerald-100">
      <div className="flex-1 flex items-center justify-center py-12 px-4 sm:px-6 lg:px-8">
        <div className="w-full max-w-md">
          <div className="text-center mb-8">
            <div className="mx-auto w-16 h-16 bg-green-600 rounded-full flex items-center justify-center mb-4">
              <Settings className="w-8 h-8 text-white" />
            </div>
            <h1 className="text-3xl font-bold text-gray-900 mb-2">
              {t('initTitle')}
            </h1>
            <p className="text-gray-600">
              {t('initSubtitle')}
            </p>
          </div>

          <Card className="shadow-xl">
            <CardHeader>
              <CardTitle>{t('adminUsername')}</CardTitle>
            </CardHeader>
            <CardContent>
              <form onSubmit={handleInit} className="space-y-6">
                {error && <ClosableAlert variant="destructive" title={t('loginError')} onClose={() => setError('')}>{error}</ClosableAlert>}
                {success && <ClosableAlert title={t('loginSuccess')} onClose={() => setSuccess('')}>{success}</ClosableAlert>}

                <div className="space-y-2">
                  <Label htmlFor="username">{t('adminUsername')}</Label>
                  <div className="relative">
                    <User className="absolute left-3 top-1/2 -translate-y-1/2 h-5 w-5 text-gray-400" />
                    <Input
                      id="username"
                      name="username"
                      value={form.username}
                      onChange={handleChange}
                      placeholder={t('enterAdminUsername')}
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
                      placeholder={t('passwordMinLength')}
                      required
                      className="pl-10"
                    />
                  </div>
                </div>
                <div className="space-y-2">
                  <Label htmlFor="confirmPassword">{t('confirmPassword')}</Label>
                  <div className="relative">
                    <Lock className="absolute left-3 top-1/2 -translate-y-1/2 h-5 w-5 text-gray-400" />
                    <Input
                      id="confirmPassword"
                      name="confirmPassword"
                      type="password"
                      value={form.confirmPassword}
                      onChange={handleChange}
                      placeholder={t('enterConfirmPassword')}
                      required
                      className="pl-10"
                    />
                  </div>
                </div>
                <Button type="submit" disabled={isLoading} className="w-full mt-8" variant="default">
                  {isLoading ? t('initLoading') : t('initButton')}
                </Button>
              </form>
            </CardContent>
          </Card>

          <div className="text-center mt-6">
            <div className="flex items-center justify-center space-x-2 text-sm text-gray-500">
                <Info className="w-4 h-4" />
                <span>{t('initTip')}</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}