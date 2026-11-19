"use client";

import { useCallback, useEffect, useState } from 'react';
import { useLocale, useTranslations } from 'next-intl';
import axios from 'axios';
import { Ban, ExternalLink, Loader2, Pencil, RefreshCw, UserMinus, UserRound } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { ConfirmDialog } from '@/components/ui/confirm-dialog';
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from '@/components/ui/dialog';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Textarea } from '@/components/ui/textarea';
import { useApiError } from '@/lib/api-error';
import { formatDateTime } from '@/lib/format';
import { toast } from '@/lib/toast';
import type { LinkRoomInfo, Server } from '@/stores/servers';

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

// 详情页的"P2P 联机房"卡片（开启联机房后显示）。每个管理端一个房间，开启了联机房的服务器都在里面，
// 所以这张卡片显示的是整个房间：房间资料（可编辑）、房间里的服务器、玩家加入方式、成员（踢出、封禁）和封禁列表（解封）
export function LinkCard({ server }: { server: Server }) {
  const t = useTranslations('linkRoom');
  const locale = useLocale();
  const apiError = useApiError();
  const status = server.link_status;
  const online = status?.state === 'online';
  // 没绑定 arkfgf.com 账号时房间没有登记，不请求成员
  const loginRequired = status?.state === 'login_required';

  // data 为 null 表示还没有登记；undefined 表示加载中
  const [data, setData] = useState<LinkMembers | null | undefined>(undefined);
  const [error, setError] = useState('');
  const [room, setRoom] = useState<LinkRoomInfo | null>(null);
  const [roomError, setRoomError] = useState('');
  const [loading, setLoading] = useState(false);
  const [action, setAction] = useState<Action | null>(null);
  const [editing, setEditing] = useState(false);

  const load = useCallback(async () => {
    setLoading(true);
    const [members, roomRes] = await Promise.allSettled([
      loginRequired ? Promise.reject(null) : axios.get('/api/link/members'),
      axios.get('/api/link/room'),
    ]);
    if (members.status === 'fulfilled') {
      setData(members.value.data.data);
      setError('');
    } else if (loginRequired || (axios.isAxiosError(members.reason) && members.reason.response?.status === 409)) {
      setData(null);
      setError('');
    } else {
      setError(apiError(members.reason, t('loadFailed')));
    }
    if (roomRes.status === 'fulfilled') {
      setRoom(roomRes.value.data.data);
      setRoomError('');
    } else {
      setRoomError(apiError(roomRes.reason, t('roomLoadFailed')));
    }
    setLoading(false);
  }, [loginRequired, apiError, t]);

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
        await axios.delete(`/api/link/bans/${encodeURIComponent(action.steamId)}`);
        toast.success(t('unbanned'));
      } else {
        await axios.post('/api/link/kick', { steam_id: action.steamId, ban: action.kind === 'ban' });
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

        {/* 房间资料和房间里的服务器 */}
        {roomError && <p className="text-sm text-red-700">{roomError}</p>}
        {room && (
          <div className="grid gap-6 md:grid-cols-2">
            <div className="space-y-2 text-sm">
              <div className="flex items-center justify-between gap-2">
                <span className="text-sm font-medium">{t('room')}</span>
                <Button variant="outline" size="sm" onClick={() => setEditing(true)}>
                  <Pencil className="mr-1.5 size-3.5" />
                  {t('editRoom')}
                </Button>
              </div>
              <div>
                <span className="text-xs text-muted-foreground">{t('roomName')}：</span>
                {room.name ? <span className="font-medium">{room.name}</span> : <span className="text-muted-foreground">{t('roomNameFallback')}</span>}
              </div>
              {room.listing.description || room.listing.tags.length || room.listing.qq_group || room.listing.discord ? (
                <>
                  {room.listing.tags.length > 0 && (
                    <div className="flex flex-wrap gap-1">
                      {room.listing.tags.map((tag) => (
                        <span key={tag} className="rounded border px-1.5 py-0.5 text-xs">
                          {tag}
                        </span>
                      ))}
                    </div>
                  )}
                  {room.listing.description && (
                    <p className="line-clamp-3 whitespace-pre-line text-muted-foreground">{room.listing.description}</p>
                  )}
                </>
              ) : (
                <p className="text-xs text-muted-foreground">{t('profileFallback')}</p>
              )}
              <p className="text-xs text-muted-foreground">{t('roomHint')}</p>
            </div>
            <div>
              <div className="mb-1.5 text-sm font-medium">
                {t('servers')} <span className="num text-muted-foreground">{room.servers.length}</span>
              </div>
              <ul className="divide-y rounded-md border">
                {room.servers.map((s) => (
                  <li key={s.id} className="flex items-center gap-2.5 px-3 py-2 text-sm">
                    <span className={`size-2 flex-shrink-0 rounded-full ${s.online ? 'bg-green-500' : 'bg-stone-500'}`} />
                    <div className="min-w-0 flex-1">
                      <div className="truncate font-medium">
                        {s.session_name}
                        {s.id === server.id && <span className="ml-1.5 text-xs text-muted-foreground">({t('thisServer')})</span>}
                      </div>
                      <div className="num truncate text-xs text-muted-foreground">
                        {s.port} / {s.query_port}
                      </div>
                    </div>
                    <span className="text-xs text-muted-foreground">{s.online ? t('serverOnline') : t('serverOffline')}</span>
                  </li>
                ))}
              </ul>
            </div>
          </div>
        )}

        {data === undefined && !error && (
          <p className="flex items-center gap-2 text-sm text-muted-foreground">
            <Loader2 className="size-4 animate-spin" />
            {t('loading')}
          </p>
        )}
        {error && <p className="text-sm text-red-700">{error}</p>}
        {data === null && (
          <p className="text-sm text-muted-foreground">{status?.state === 'login_required' ? t('loginRequired') : t('notRegistered')}</p>
        )}

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
      {room && (
        <LinkRoomDialog
          room={room}
          open={editing}
          onClose={() => setEditing(false)}
          onSaved={(info) => setRoom(info)}
        />
      )}
    </Card>
  );
}

const DESCRIPTION_MAX = 500;

// 编辑联机房的名称和收录资料（PUT /api/link/room），字段和收录站的收录资料一样
function LinkRoomDialog({
  room,
  open,
  onClose,
  onSaved,
}: {
  room: LinkRoomInfo;
  open: boolean;
  onClose: () => void;
  onSaved: (info: LinkRoomInfo) => void;
}) {
  const t = useTranslations('linkRoom');
  const tl = useTranslations('listing');
  const apiError = useApiError();
  const [form, setForm] = useState({ name: '', description: '', tags: '', qq_group: '', discord: '' });
  const [saving, setSaving] = useState(false);

  // 每次打开时用当前资料填表（在渲染期间按 props 调整状态，不在 effect 里 setState）
  const [openedFor, setOpenedFor] = useState(false);
  if (open !== openedFor) {
    setOpenedFor(open);
    if (open) {
      setForm({
        name: room.name,
        description: room.listing.description,
        tags: room.listing.tags.join(', '),
        qq_group: room.listing.qq_group,
        discord: room.listing.discord,
      });
    }
  }

  const save = async () => {
    setSaving(true);
    try {
      const res = await axios.put('/api/link/room', {
        name: form.name.trim(),
        listing: {
          description: form.description,
          tags: form.tags.split(/[,，]/).map((s) => s.trim()).filter(Boolean),
          qq_group: form.qq_group.trim(),
          discord: form.discord.trim(),
        },
      });
      onSaved(res.data.data);
      toast.success(t('roomSaved'));
      onClose();
    } catch (err) {
      toast.error(apiError(err, tl('saveFailed')));
    } finally {
      setSaving(false);
    }
  };

  return (
    <Dialog open={open} onOpenChange={(o) => !o && !saving && onClose()}>
      <DialogContent className="max-w-lg">
        <DialogHeader>
          <DialogTitle>{t('editRoomTitle')}</DialogTitle>
          <DialogDescription>{t('editRoomDesc')}</DialogDescription>
        </DialogHeader>
        <div className="space-y-4">
          <div className="space-y-1">
            <Label htmlFor="link-room-name">{t('roomName')}</Label>
            <Input
              id="link-room-name"
              value={form.name}
              maxLength={64}
              onChange={(e) => setForm((f) => ({ ...f, name: e.target.value }))}
              placeholder={t('roomNamePlaceholder')}
            />
          </div>
          <div className="space-y-1">
            <Label htmlFor="link-room-description">{tl('description')}</Label>
            <Textarea
              id="link-room-description"
              value={form.description}
              maxLength={DESCRIPTION_MAX}
              rows={4}
              onChange={(e) => setForm((f) => ({ ...f, description: e.target.value }))}
              placeholder={tl('descriptionPlaceholder')}
            />
            <p className="text-right text-xs text-muted-foreground">
              {form.description.length}/{DESCRIPTION_MAX}
            </p>
          </div>
          <div className="space-y-1">
            <Label htmlFor="link-room-tags">{tl('tags')}</Label>
            <Input
              id="link-room-tags"
              value={form.tags}
              onChange={(e) => setForm((f) => ({ ...f, tags: e.target.value }))}
              placeholder={tl('tagsPlaceholder')}
            />
            <p className="text-xs text-muted-foreground">{tl('tagsHint')}</p>
          </div>
          <div className="grid gap-4 sm:grid-cols-2">
            <div className="space-y-1">
              <Label htmlFor="link-room-qq">{tl('qq')}</Label>
              <Input
                id="link-room-qq"
                inputMode="numeric"
                value={form.qq_group}
                onChange={(e) => setForm((f) => ({ ...f, qq_group: e.target.value }))}
                placeholder="123456789"
              />
            </div>
            <div className="space-y-1">
              <Label htmlFor="link-room-discord">Discord</Label>
              <Input
                id="link-room-discord"
                value={form.discord}
                onChange={(e) => setForm((f) => ({ ...f, discord: e.target.value }))}
                placeholder="https://discord.gg/xxxx"
              />
            </div>
          </div>
          <p className="text-xs text-muted-foreground">{t('profileFallback')}</p>
        </div>
        <DialogFooter>
          <Button variant="outline" onClick={onClose} disabled={saving}>
            {tl('cancel')}
          </Button>
          <Button onClick={save} disabled={saving}>
            {saving && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
            {tl('save')}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
