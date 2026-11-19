"use client";

import { useCallback, useEffect, useState, useSyncExternalStore } from 'react';
import { useLocale, useTranslations } from 'next-intl';
import axios from 'axios';
import { ExternalLink, Globe, Link2, Loader2, LogIn, Unlink, UserRound, X } from 'lucide-react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { ConfirmDialog } from '@/components/ui/confirm-dialog';
import { FormError } from '@/components/AuthShell';
import { usePolling } from '@/hooks/use-polling';
import { useApiError, useErrorText } from '@/lib/api-error';
import { formatDateTime } from '@/lib/format';
import { toast } from '@/lib/toast';
import type { ArkfgfAccountStatus } from '@/stores/servers';

const DEFAULT_NAME = 'ArkServerCommander';
// Steam 登录后回到的公开页面（app/arkfgf-callback），不需要登录管理端
const ARKFGF_CALLBACK_PATH = '/arkfgf-callback';

// 备用方式：在当前浏览器里跳到 Steam 登录页（浏览器里没登录 Steam 时可以用 Steam 手机 App 扫码）
function BrowserFallback({ starting, onClick }: { starting: 'client' | 'browser' | null; onClick: () => void }) {
  const t = useTranslations('settings.arkfgf');
  return (
    <button
      type="button"
      disabled={starting !== null}
      onClick={onClick}
      className="inline-flex items-center gap-1.5 text-left text-xs text-primary hover:underline disabled:opacity-50"
    >
      {starting === 'browser' ? <Loader2 className="size-3.5 animate-spin" /> : <Globe className="size-3.5" />}
      {t('loginBrowser')}
    </button>
  );
}

// 浏览器当前访问管理端的来源（Steam 登录后回到这里）；服务端渲染时为空
const subscribeNoop = () => () => {};
const useOrigin = () =>
  useSyncExternalStore(
    subscribeNoop,
    () => window.location.origin,
    () => '',
  );

// 网站没拿到 Steam 资料时用户名是占位的 Player，这时只显示 SteamID
const hasNickname = (name?: string) => !!name && name !== 'Player';

function Avatar({ url }: { url?: string }) {
  if (!url) {
    return (
      <span className="flex size-12 flex-shrink-0 items-center justify-center rounded-full bg-gray-100 text-gray-500">
        <UserRound className="size-6" />
      </span>
    );
  }
  // Steam 头像是外链图片，不用 next/image
  // eslint-disable-next-line @next/next/no-img-element
  return <img src={url} alt="" width={48} height={48} className="size-12 flex-shrink-0 rounded-full ring-1 ring-primary/30" />;
}

// 设置页"arkfgf.com 账号"：发布到收录站和 P2P 联机房都需要先绑定。
// 未绑定时用 Steam 登录：主要方式是在 Steam 客户端里打开登录页（steam://openurl/），备用方式是在浏览器里跳转，
// 登录后都回到公开页面 /arkfgf-callback 完成绑定；已绑定时显示账号和解除绑定
export function ArkfgfAccountCard() {
  const t = useTranslations('settings.arkfgf');
  const locale = useLocale();
  const apiError = useApiError();
  const origin = useOrigin();

  const [status, setStatus] = useState<ArkfgfAccountStatus | null>(null);
  const [name, setName] = useState(DEFAULT_NAME);
  const errorText = useErrorText();
  const [starting, setStarting] = useState<'client' | 'browser' | null>(null);
  const [waiting, setWaiting] = useState(false);
  const [waitError, setWaitError] = useState<string | null>(null);
  const [confirmUnbind, setConfirmUnbind] = useState(false);
  const [forceError, setForceError] = useState<string | null>(null);

  const load = useCallback(
    async (refresh = false) => {
      try {
        const res = await axios.get('/api/arkfgf/account', { params: refresh ? { refresh: true } : undefined });
        setStatus(res.data.data);
      } catch (err) {
        toast.error(apiError(err, t('loadFailed')));
      }
    },
    [apiError, t],
  );

  // 打开设置页时向网站核对一次账号（服主可能在网站后台移除了这个管理端）
  useEffect(() => {
    const timer = setTimeout(() => {
      load(true);
    }, 0);
    return () => clearTimeout(timer);
  }, [load]);

  // 开始登录：client 在 Steam 客户端里打开登录页（客户端内置浏览器共享客户端的登录状态，服主只要点"登录"），
  // 回调页在客户端里打开，这里轮询状态等结果；browser 在当前页面跳到 Steam 登录页（可以用 Steam 手机 App 扫码），
  // 登录后回调页在当前浏览器里显示结果
  const login = async (mode: 'client' | 'browser') => {
    setStarting(mode);
    setWaitError(null);
    try {
      const res = await axios.post('/api/arkfgf/account/login', {
        name: name.trim(),
        return_to: window.location.origin + ARKFGF_CALLBACK_PATH,
      });
      const steamURL: string = res.data.data.steam_url;
      if (mode === 'browser') {
        // 按钮保持加载状态直到页面离开
        window.location.href = steamURL;
        return;
      }
      window.location.href = 'steam://openurl/' + steamURL;
      setWaiting(true);
    } catch (err) {
      toast.error(apiError(err, t('loginFailed')));
    }
    setStarting(null);
  };

  // 等待 Steam 客户端里完成登录：绑定成功、完成登录失败（回调页在客户端里报的错）或登录过期时停止
  const poll = useCallback(async () => {
    const res = await axios.get('/api/arkfgf/account');
    const next: ArkfgfAccountStatus = res.data.data;
    setStatus(next);
    if (next.bound) {
      setWaiting(false);
      toast.success(t('bound'));
    } else if (next.login?.error_code) {
      setWaiting(false);
      setWaitError(errorText(next.login.error, next.login.error_code, t('loginFailed')));
    } else if (!next.login?.pending) {
      setWaiting(false);
      setWaitError(t('loginExpired'));
    }
  }, [errorText, t]);
  usePolling(poll, 2000, waiting);

  const unbind = async (force: boolean) => {
    try {
      await axios.delete('/api/arkfgf/account', { params: force ? { force: true } : undefined });
      toast.success(t('unbound'));
      setForceError(null);
      load();
    } catch (err) {
      // 连不上网站：让服主选择只清本地
      if (!force && axios.isAxiosError(err) && err.response?.data?.code === 'arkfgf_unavailable') {
        setForceError(apiError(err, t('unbindFailed')));
        return;
      }
      toast.error(apiError(err, t('unbindFailed')));
      throw err;
    }
  };

  const account = status?.bound ? status.account : undefined;

  return (
    <Card>
      <CardHeader>
        <CardTitle className="flex items-center gap-2 text-lg">
          <Link2 className="h-5 w-5" />
          {t('title')}
        </CardTitle>
      </CardHeader>
      <CardContent className="space-y-4 text-sm">
        <p className="text-xs text-muted-foreground">{t('description')}</p>

        {status === null && <Loader2 className="size-5 animate-spin text-muted-foreground" />}

        {status && !status.enabled && (
          <p className="rounded-md border border-amber-300 bg-amber-50 px-3 py-2 text-amber-900">{t('disabled')}</p>
        )}

        {status?.enabled && !status.bound && !waiting && (
          <form
            onSubmit={(e) => {
              e.preventDefault();
              login('client');
            }}
            className="space-y-3"
          >
            <div className="space-y-1">
              <Label htmlFor="arkfgf-name">{t('name')}</Label>
              <Input id="arkfgf-name" value={name} maxLength={64} onChange={(e) => setName(e.target.value)} placeholder={DEFAULT_NAME} />
              <p className="text-xs text-muted-foreground">{t('nameHint')}</p>
            </div>
            {waitError && <FormError message={waitError} />}
            <Button type="submit" disabled={starting !== null} className="hud-caps">
              {starting === 'client' ? <Loader2 className="animate-spin" /> : <LogIn />}
              {t('loginClient')}
            </Button>
            <p className="text-xs text-muted-foreground">{t('loginClientHint')}</p>
            <BrowserFallback starting={starting} onClick={() => login('browser')} />
            <p className="text-xs text-muted-foreground">{t('loginHint', { origin: origin || '-' })}</p>
          </form>
        )}

        {status?.enabled && !status.bound && waiting && (
          <div className="hud relative space-y-3 rounded-sm border border-amber-200 bg-amber-50 px-4 py-3">
            <div className="flex items-start gap-3 text-amber-900">
              <Loader2 className="mt-0.5 size-5 flex-shrink-0 animate-spin" />
              <div className="space-y-1">
                <div className="font-medium">{t('waitingTitle')}</div>
                <div className="text-xs text-amber-800">{t('waitingHint')}</div>
              </div>
            </div>
            <div className="flex flex-wrap items-center gap-2">
              <Button variant="outline" size="sm" onClick={() => setWaiting(false)}>
                <X className="size-3.5" />
                {t('cancelWaiting')}
              </Button>
            </div>
            <BrowserFallback starting={starting} onClick={() => login('browser')} />
          </div>
        )}

        {account && (
          <div className="space-y-4">
            <div className="flex items-center gap-4">
              <Avatar url={account.avatar_url} />
              <dl className="grid min-w-0 flex-1 grid-cols-[auto_1fr] gap-x-4 gap-y-1">
                <dt className="text-muted-foreground">{t('steamName')}</dt>
                <dd className="min-w-0 truncate">
                  <a
                    href={`https://steamcommunity.com/profiles/${account.steam_id}`}
                    target="_blank"
                    rel="noreferrer"
                    className="num font-medium text-primary hover:underline"
                  >
                    {account.steam_id}
                  </a>
                  {hasNickname(account.user_name) && <span className="ml-2 text-muted-foreground">{account.user_name}</span>}
                </dd>
                <dt className="text-muted-foreground">{t('commanderName')}</dt>
                <dd className="truncate">{account.commander_name}</dd>
                <dt className="text-muted-foreground">{t('boundAt')}</dt>
                <dd className="num">{formatDateTime(account.bound_at, locale)}</dd>
              </dl>
            </div>
            <div className="flex flex-wrap items-center gap-2">
              {status?.dashboard_url && (
                <Button variant="outline" size="sm" asChild>
                  <a href={status.dashboard_url} target="_blank" rel="noreferrer">
                    <ExternalLink className="size-3.5" />
                    {t('manage')}
                  </a>
                </Button>
              )}
              <Button variant="outline" size="sm" className="text-red-700 hover:text-red-800" onClick={() => setConfirmUnbind(true)}>
                <Unlink className="size-3.5" />
                {t('unbind')}
              </Button>
            </div>
          </div>
        )}
      </CardContent>

      <ConfirmDialog
        open={confirmUnbind}
        onOpenChange={setConfirmUnbind}
        title={t('unbindTitle')}
        description={t('unbindDesc')}
        confirmText={t('unbind')}
        destructive
        onConfirm={() => unbind(false)}
      />
      <ConfirmDialog
        open={forceError !== null}
        onOpenChange={(open) => !open && setForceError(null)}
        title={t('forceTitle')}
        description={t('forceDesc', { error: forceError ?? '' })}
        confirmText={t('force')}
        destructive
        onConfirm={() => unbind(true)}
      />
    </Card>
  );
}
