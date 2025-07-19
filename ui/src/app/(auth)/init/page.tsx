"use client";

import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import { useAuthActions, useAuthIsLoading } from '@/stores/auth';
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { ClosableAlert } from "@/components/ui/closable-alert";
import { Settings, User, Lock, Info } from 'lucide-react';

export default function InitPage() {
  const router = useRouter();
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
      setError('密码不匹配');
      return;
    }
    if (form.password.length < 6) {
      setError('密码长度至少为6位');
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
              系统初始化
            </h1>
            <p className="text-gray-600">
              首次使用，请设置管理员账户
            </p>
          </div>

          <Card className="shadow-xl">
            <CardHeader>
              <CardTitle>创建管理员</CardTitle>
            </CardHeader>
            <CardContent>
              <form onSubmit={handleInit} className="space-y-6">
                {error && <ClosableAlert variant="destructive" title="错误" onClose={() => setError('')}>{error}</ClosableAlert>}
                {success && <ClosableAlert title="成功" onClose={() => setSuccess('')}>{success}</ClosableAlert>}

                <div className="space-y-2">
                  <Label htmlFor="username">管理员用户名</Label>
                  <div className="relative">
                    <User className="absolute left-3 top-1/2 -translate-y-1/2 h-5 w-5 text-gray-400" />
                    <Input
                      id="username"
                      name="username"
                      value={form.username}
                      onChange={handleChange}
                      placeholder="输入管理员用户名"
                      required
                      className="pl-10"
                    />
                  </div>
                </div>
                <div className="space-y-2">
                  <Label htmlFor="password">密码</Label>
                  <div className="relative">
                    <Lock className="absolute left-3 top-1/2 -translate-y-1/2 h-5 w-5 text-gray-400" />
                    <Input
                      id="password"
                      name="password"
                      type="password"
                      value={form.password}
                      onChange={handleChange}
                      placeholder="输入密码 (至少6位)"
                      required
                      className="pl-10"
                    />
                  </div>
                </div>
                <div className="space-y-2">
                  <Label htmlFor="confirmPassword">确认密码</Label>
                  <div className="relative">
                    <Lock className="absolute left-3 top-1/2 -translate-y-1/2 h-5 w-5 text-gray-400" />
                    <Input
                      id="confirmPassword"
                      name="confirmPassword"
                      type="password"
                      value={form.confirmPassword}
                      onChange={handleChange}
                      placeholder="再次输入密码"
                      required
                      className="pl-10"
                    />
                  </div>
                </div>
                <Button type="submit" disabled={isLoading} className="w-full mt-8" variant="default">
                  {isLoading ? '初始化中...' : '完成初始化'}
                </Button>
              </form>
            </CardContent>
          </Card>

          <div className="text-center mt-6">
            <div className="flex items-center justify-center space-x-2 text-sm text-gray-500">
                <Info className="w-4 h-4" />
                <span>此操作只需执行一次，用于创建第一个管理员帐户。</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}