"use client";

import { useEffect, useState } from 'react';
import { useLocale, useTranslations } from 'next-intl';
import axios from 'axios';
import { Link2 } from 'lucide-react';
import { Link } from '@/navigation';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Switch } from '@/components/ui/switch';
import { useApiError } from '@/lib/api-error';
import { toast } from '@/lib/toast';
import { serversActions, type ArkfgfAccountStatus, type Server } from '@/stores/servers';
import { gameOf } from '@/lib/games';

// 状态文字和出错原因（中文界面用后端的详情，其他语言按错误码翻译）
function StateText({ color, text, url, linkText, detail }: { color: string; text: string; url?: string; linkText: string; detail: string }) {
  return (
    <span className="inline-flex flex-col gap-0.5">
      <span className={color}>
        {text}
        {url && (
          <a href={url} target="_blank" rel="noreferrer" className="ml-2 text-primary hover:underline">
            {linkText}
          </a>
        )}
      </span>
      {detail && <span className="text-xs text-muted-foreground">{detail}</span>}
    </span>
  );
}

function PublishState({ server }: { server: Server }) {
  const t = useTranslations('detail');
  const locale = useLocale();
  const status = server.publish_status;
  if (!status) return null;
  const color = status.state === 'published' ? 'text-green-700' : status.state === 'error' ? 'text-red-700' : 'text-amber-800';
  let detail = '';
  if (status.state === 'error') {
    const key = `publishErrors.${status.error_code}`;
    if (locale === 'zh' && status.error) detail = status.error;
    // ASA 没有查询端口，提示里用游戏端口
    else if (status.error_code && t.has(key)) detail = t(key, { port: status.error_code === 'asa_not_listed' || gameOf(server) === 'asa' ? server.port : server.query_port });
    else detail = status.error ?? '';
  }
  return (
    <StateText
      color={color}
      text={t(`publishStates.${status.state}`)}
      url={status.state === 'published' ? status.url : undefined}
      linkText={t('publishView')}
      detail={detail}
    />
  );
}

function LinkState({ server }: { server: Server }) {
  const t = useTranslations('detail');
  const locale = useLocale();
  const status = server.link_status;
  if (!status) return null;
  const color =
    status.state === 'online'
      ? 'text-green-700'
      : status.state === 'error'
        ? 'text-red-700'
        : status.state === 'offline' || status.state === 'disabled'
          ? 'text-muted-foreground'
          : 'text-amber-800';
  let detail = '';
  if (status.state === 'error') {
    const key = `linkErrors.${status.error_code}`;
    if (locale === 'zh' && status.error) detail = status.error;
    else if (status.error_code && t.has(key)) detail = t(key);
    else detail = status.error ?? '';
  }
  return <StateText color={color} text={t(`linkStates.${status.state}`)} url={status.room_url} linkText={t('linkView')} detail={detail} />;
}

// 详情页概览的 arkfgf.com 卡片：发布到收录站、P2P 联机房的开关和状态。
// 它们是网站提供的服务，不属于游戏设置，所以不在"编辑设置"里；开关立即生效，不用重启。
// 都需要先在设置页绑定 arkfgf.com 账号：未绑定时不能打开开关（已打开的可以关闭）
export function ArkfgfCard({ server }: { server: Server }) {
  const t = useTranslations('detail');
  const apiError = useApiError();
  const [saving, setSaving] = useState<'publish' | 'link' | null>(null);
  const [account, setAccount] = useState<ArkfgfAccountStatus | null>(null);

  useEffect(() => {
    let cancelled = false;
    axios
      .get('/api/arkfgf/account')
      .then((res) => {
        if (!cancelled) setAccount(res.data.data);
      })
      .catch(() => {
        // 读不到时不限制开关，后端的状态（login_required）仍会显示
      });
    return () => {
      cancelled = true;
    };
  }, []);
  // ARKFGF_API 为空时整体关闭，开关照旧（状态显示 disabled）
  const needLogin = account !== null && account.enabled && !account.bound;

  const toggle = async (key: 'publish' | 'link', enabled: boolean) => {
    setSaving(key);
    try {
      await serversActions.updateServer(server.id, key === 'publish' ? { publish: enabled } : { link: enabled });
      toast.success(t(`${key}${enabled ? 'On' : 'Off'}Toast`));
    } catch (err) {
      toast.error(apiError(err, t('toggleFailed')));
    } finally {
      setSaving(null);
    }
  };

  const items = [
    { key: 'publish' as const, label: t('publishSwitch'), hint: gameOf(server) === 'asa' ? t('publishHintAsa') : t('publishHint'), state: <PublishState server={server} /> },
    { key: 'link' as const, label: t('linkSwitch'), hint: t('linkHint'), state: <LinkState server={server} /> },
  ];

  return (
    <Card>
      <CardHeader className="pb-2">
        <CardTitle className="font-display text-base">{t('external')}</CardTitle>
        <p className="text-xs text-muted-foreground">{t('externalDesc')}</p>
      </CardHeader>
      <CardContent className="divide-y">
        {needLogin && (
          <div className="mb-3 flex items-start gap-2 rounded-md border border-amber-300 bg-amber-50 px-3 py-2 text-sm text-amber-900">
            <Link2 className="mt-0.5 size-4 flex-shrink-0" />
            <span>
              {t('loginRequiredNotice')}
              <Link href="/settings" className="ml-1 font-medium text-primary hover:underline">
                {t('loginRequiredLink')}
              </Link>
            </span>
          </div>
        )}
        {items.map((item) => {
          const on = server[item.key] ?? false;
          const id = `arkfgf-${item.key}`;
          // 未绑定账号时不能打开，已打开的可以关闭
          const blocked = needLogin && !on;
          return (
            <div key={item.key} className="space-y-1.5 py-3 first:pt-0 last:pb-0">
              <div className="flex items-center justify-between gap-4">
                <label htmlFor={id} className="text-sm font-medium">
                  {item.label}
                </label>
                <span title={blocked ? t('loginRequiredSwitch') : undefined}>
                  <Switch id={id} checked={on} disabled={saving !== null || blocked} onCheckedChange={(v) => toggle(item.key, v)} />
                </span>
              </div>
              {on && <div className="text-sm">{item.state}</div>}
              <p className="text-xs leading-relaxed text-muted-foreground">{item.hint}</p>
            </div>
          );
        })}
      </CardContent>
    </Card>
  );
}
