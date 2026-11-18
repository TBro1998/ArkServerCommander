"use client";

import { useCallback, useEffect, useState } from 'react';
import { useLocale, useTranslations } from 'next-intl';
import axios from 'axios';
import { Ban, ExternalLink, Loader2, RefreshCw, UserMinus, UserRound } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { ConfirmDialog } from '@/components/ui/confirm-dialog';
import { useApiError } from '@/lib/api-error';
import { formatDateTime } from '@/lib/format';
import { toast } from '@/lib/toast';
import type { Server } from '@/stores/servers';

// 网站返回的时间可能是 Unix 秒数或时间字符串
type SiteTime = number | string | null | undefined;

interface LinkMember {
  steam_id: string;
  name: string;
  avatar_url: string;
  credential_id: string;
  joined_at: SiteTime;
}

interface LinkBan {
  steam_id: string;
  name: string;
  avatar_url: string;
  banned_at: SiteTime;
}

interface LinkMembers {
  members: LinkMember[];
  bans: LinkBan[];
}

type Action = { kind: 'kick' | 'ban' | 'unban'; steamId: string; name: string };

function siteTime(value: SiteTime, locale: string) {
  if (value === null || value === undefined || value === '') return '-';
  if (typeof value === 'number') return formatDateTime(new Date(value * 1000).toISOString(), locale);
  return formatDateTime(value, locale);
}

function Avatar({ url }: { url: string }) {
  if (!url) {
    return (
      <span className="flex size-7 flex-shrink-0 items-center justify-center rounded-full bg-gray-100 text-gray-500">
        <UserRound className="size-4" />
      </span>
    );
  }
  // Steam 头像是外链图片，不用 next/image
  // eslint-disable-next-line @next/next/no-img-element
  return <img src={url} alt="" width={28} height={28} className="size-7 flex-shrink-0 rounded-full" />;
}

// 详情页的"P2P 联机房"卡片（开启联机房后显示）：房间链接、玩家加入方式、成员（踢出、封禁）和封禁列表（解封）
export function LinkCard({ server }: { server: Server }) {
  const t = useTranslations('linkRoom');
  const locale = useLocale();
  const apiError = useApiError();
  const status = server.link_status;
  const online = status?.state === 'online';

  // data 为 null 表示还没有登记；undefined 表示加载中
  const [data, setData] = useState<LinkMembers | null | undefined>(undefined);
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);
  const [action, setAction] = useState<Action | null>(null);

  const load = useCallback(async () => {
    setLoading(true);
    try {
      const res = await axios.get(`/api/servers/${server.id}/link/members`);
      setData(res.data.data);
      setError('');
    } catch (err) {
      if (axios.isAxiosError(err) && err.response?.status === 409) {
        setData(null);
        setError('');
      } else {
        setError(apiError(err, t('loadFailed')));
      }
    } finally {
      setLoading(false);
    }
  }, [server.id, apiError, t]);

  // 登记完成、上线后重新读取
  useEffect(() => {
    const timer = setTimeout(() => {
      load();
    }, 0);
    return () => clearTimeout(timer);
  }, [load, online, status?.room_url]);

  const run = async () => {
    if (!action) return;
    try {
      if (action.kind === 'unban') {
        await axios.delete(`/api/servers/${server.id}/link/bans/${encodeURIComponent(action.steamId)}`);
        toast.success(t('unbanned'));
      } else {
        await axios.post(`/api/servers/${server.id}/link/kick`, { steam_id: action.steamId, ban: action.kind === 'ban' });
        toast.success(action.kind === 'ban' ? t('banned') : t('kicked'));
      }
      load();
    } catch (err) {
      toast.error(apiError(err, t('actionFailed')));
      throw err;
    }
  };

  const displayName = (name: string, steamId: string) => name || steamId;

  return (
    <Card className="lg:col-span-2">
      <CardHeader className="flex flex-row items-center justify-between gap-2 pb-2">
        <CardTitle className="font-display text-base">{t('title')}</CardTitle>
        <div className="flex items-center gap-3">
          {status?.room_url && (
            <a href={status.room_url} target="_blank" rel="noreferrer" className="inline-flex items-center gap-1 text-sm text-primary hover:underline">
              {t('viewRoom')}
              <ExternalLink className="size-3.5" />
            </a>
          )}
          <Button variant="ghost" size="sm" className="h-7 w-7 p-0" title={t('refresh')} onClick={() => load()} disabled={loading}>
            <RefreshCw className={loading ? 'size-3.5 animate-spin' : 'size-3.5'} />
          </Button>
        </div>
      </CardHeader>
      <CardContent className="space-y-4">
        <p className="rounded-lg bg-gray-50 p-3 text-xs text-gray-600">{t('howTo')}</p>

        {data === undefined && !error && (
          <p className="flex items-center gap-2 text-sm text-muted-foreground">
            <Loader2 className="size-4 animate-spin" />
            {t('loading')}
          </p>
        )}
        {error && <p className="text-sm text-red-700">{error}</p>}
        {data === null && <p className="text-sm text-muted-foreground">{t('notRegistered')}</p>}

        {data && (
          <div className="grid gap-6 md:grid-cols-2">
            {/* 成员 */}
            <div>
              <div className="mb-1.5 flex items-baseline justify-between gap-2">
                <span className="text-sm font-medium">
                  {t('members')} <span className="num text-muted-foreground">{data.members.length}</span>
                </span>
                <span className="text-xs text-muted-foreground">{t('membersHint')}</span>
              </div>
              {data.members.length === 0 ? (
                <p className="py-4 text-center text-sm text-muted-foreground">{t('noMembers')}</p>
              ) : (
                <ul className="divide-y rounded-md border">
                  {data.members.map((m) => (
                    <li key={m.steam_id} className="flex items-center gap-2.5 px-3 py-2 text-sm">
                      <Avatar url={m.avatar_url} />
                      <div className="min-w-0 flex-1">
                        <div className="truncate font-medium">{displayName(m.name, m.steam_id)}</div>
                        <div className="truncate text-xs text-muted-foreground">
                          <span className="num">{m.steam_id}</span> · {t('joinedAt')} {siteTime(m.joined_at, locale)}
                        </div>
                      </div>
                      <Button
                        variant="ghost"
                        size="sm"
                        className="h-7 px-2 text-amber-800"
                        onClick={() => setAction({ kind: 'kick', steamId: m.steam_id, name: displayName(m.name, m.steam_id) })}
                      >
                        <UserMinus className="mr-1 size-3.5" />
                        {t('kick')}
                      </Button>
                      <Button
                        variant="ghost"
                        size="sm"
                        className="h-7 px-2 text-red-700"
                        onClick={() => setAction({ kind: 'ban', steamId: m.steam_id, name: displayName(m.name, m.steam_id) })}
                      >
                        <Ban className="mr-1 size-3.5" />
                        {t('ban')}
                      </Button>
                    </li>
                  ))}
                </ul>
              )}
            </div>

            {/* 封禁列表 */}
            <div>
              <div className="mb-1.5 text-sm font-medium">
                {t('bans')} <span className="num text-muted-foreground">{data.bans.length}</span>
              </div>
              {data.bans.length === 0 ? (
                <p className="py-4 text-center text-sm text-muted-foreground">{t('noBans')}</p>
              ) : (
                <ul className="divide-y rounded-md border">
                  {data.bans.map((b) => (
                    <li key={b.steam_id} className="flex items-center gap-2.5 px-3 py-2 text-sm">
                      <Avatar url={b.avatar_url} />
                      <div className="min-w-0 flex-1">
                        <div className="truncate font-medium">{displayName(b.name, b.steam_id)}</div>
                        <div className="truncate text-xs text-muted-foreground">
                          <span className="num">{b.steam_id}</span> · {t('bannedAt')} {siteTime(b.banned_at, locale)}
                        </div>
                      </div>
                      <Button
                        variant="outline"
                        size="sm"
                        className="h-7 px-2"
                        onClick={() => setAction({ kind: 'unban', steamId: b.steam_id, name: displayName(b.name, b.steam_id) })}
                      >
                        {t('unban')}
                      </Button>
                    </li>
                  ))}
                </ul>
              )}
            </div>
          </div>
        )}
      </CardContent>

      <ConfirmDialog
        open={action !== null}
        onOpenChange={(open) => !open && setAction(null)}
        title={action ? t(`${action.kind}Title`, { name: action.name }) : ''}
        description={action ? t(`${action.kind}Desc`) : undefined}
        confirmText={action ? t(action.kind) : undefined}
        destructive={action?.kind !== 'unban'}
        onConfirm={run}
      />
    </Card>
  );
}
