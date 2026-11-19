"use client";

import { useState } from 'react';
import { useLocale, useTranslations } from 'next-intl';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Switch } from '@/components/ui/switch';
import { useApiError } from '@/lib/api-error';
import { toast } from '@/lib/toast';
import { serversActions, type Server } from '@/stores/servers';

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
    else if (status.error_code && t.has(key)) detail = t(key, { port: server.query_port });
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
// 它们是网站提供的服务，不属于游戏设置，所以不在"编辑设置"里；开关立即生效，不用重启
export function ArkfgfCard({ server }: { server: Server }) {
  const t = useTranslations('detail');
  const apiError = useApiError();
  const [saving, setSaving] = useState<'publish' | 'link' | null>(null);

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
    { key: 'publish' as const, label: t('publishSwitch'), hint: t('publishHint'), state: <PublishState server={server} /> },
    { key: 'link' as const, label: t('linkSwitch'), hint: t('linkHint'), state: <LinkState server={server} /> },
  ];

  return (
    <Card>
      <CardHeader className="pb-2">
        <CardTitle className="font-display text-base">{t('external')}</CardTitle>
        <p className="text-xs text-muted-foreground">{t('externalDesc')}</p>
      </CardHeader>
      <CardContent className="divide-y">
        {items.map((item) => {
          const on = server[item.key] ?? false;
          const id = `arkfgf-${item.key}`;
          return (
            <div key={item.key} className="space-y-1.5 py-3 first:pt-0 last:pb-0">
              <div className="flex items-center justify-between gap-4">
                <label htmlFor={id} className="text-sm font-medium">
                  {item.label}
                </label>
                <Switch id={id} checked={on} disabled={saving !== null} onCheckedChange={(v) => toggle(item.key, v)} />
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
