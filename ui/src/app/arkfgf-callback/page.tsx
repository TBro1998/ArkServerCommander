"use client";

import { useEffect, useRef, useState } from 'react';
import { useTranslations } from 'next-intl';
import axios from 'axios';
import { CircleCheck, Loader2 } from 'lucide-react';
import { Link } from '@/navigation';
import { AuthShell, FormError } from '@/components/AuthShell';
import { Button } from '@/components/ui/button';
import { useApiError } from '@/lib/api-error';

// Steam 登录后回到这里（设置页"arkfgf.com 账号"的回调地址）。公开页面，不需要登录管理端：
// 在 Steam 客户端里登录时，这一页是在客户端的内置浏览器里打开的，那里没有管理端的登录 cookie。
// 把 URL 上的全部参数（state、openid.*）交给后端（后端只接受自己开始的登录的 state），
// 后端转给 arkfgf.com 核验并保存账号。state 只能用一次，所以只提交一次（开发模式下 effect 会执行两次）
export default function ArkfgfCallbackPage() {
  const t = useTranslations('settings.arkfgfCallback');
  const apiError = useApiError();
  const submitted = useRef(false);
  const [result, setResult] = useState<{ ok: true } | { ok: false; error: string } | null>(null);

  useEffect(() => {
    if (submitted.current) return;
    submitted.current = true;
    const params: Record<string, string> = {};
    new URLSearchParams(window.location.search).forEach((value, key) => {
      params[key] = value;
    });
    // 参数用过就失效，从地址栏去掉，刷新页面不会再提交一次
    window.history.replaceState(null, '', window.location.pathname);
    axios
      .post('/api/arkfgf/account/login/finish', { params })
      .then(() => setResult({ ok: true }))
      .catch((err) => setResult({ ok: false, error: apiError(err, t('failed')) }));
  }, [apiError, t]);

  return (
    <AuthShell tagline={t('tagline')} title={t('title')}>
      {result === null && (
        <div role="status" aria-live="polite" className="flex items-center gap-3 text-sm text-muted-foreground">
          <Loader2 className="size-5 animate-spin text-primary" />
          {t('finishing')}
        </div>
      )}
      {result?.ok && (
        <div className="space-y-5">
          <div className="flex items-start gap-3 rounded-sm border border-green-200 bg-green-50 px-3 py-3 text-sm text-green-800">
            <CircleCheck className="mt-0.5 size-5 flex-shrink-0" />
            <div>
              <div className="font-medium">{t('success')}</div>
              <div className="mt-1 text-green-700">{t('successHint')}</div>
            </div>
          </div>
          <Button asChild variant="outline" className="w-full">
            <Link href="/settings">{t('openSettings')}</Link>
          </Button>
        </div>
      )}
      {result && !result.ok && (
        <div className="space-y-5">
          <FormError message={t('failedWith', { error: result.error })} />
          <p className="text-sm text-muted-foreground">{t('failedHint')}</p>
          <Button asChild variant="outline" className="w-full">
            <Link href="/settings">{t('openSettings')}</Link>
          </Button>
        </div>
      )}
    </AuthShell>
  );
}
