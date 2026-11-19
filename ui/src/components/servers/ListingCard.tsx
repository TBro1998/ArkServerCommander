"use client";

import { useCallback, useEffect, useState } from 'react';
import { useTranslations } from 'next-intl';
import axios from 'axios';
import { Copy, ExternalLink, Loader2, Pencil } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from '@/components/ui/dialog';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Textarea } from '@/components/ui/textarea';
import { useApiError } from '@/lib/api-error';
import { toast } from '@/lib/toast';
import { serversActions, type ListingProfile, type Server } from '@/stores/servers';

interface StatTotals {
  views: number;
  connects: number;
  copies: number;
}

interface PublishStats {
  id: number;
  url: string;
  badge: string;
  week: StatTotals;
  month: StatTotals;
}

const DESCRIPTION_MAX = 500;
const EMPTY_LISTING: ListingProfile = { description: '', tags: [], qq_group: '', discord: '' };

// 详情页的"收录站"卡片（开启发布后显示）：收录站上的访问统计、状态徽章的引用代码、编辑收录资料
export function ListingCard({ server }: { server: Server }) {
  const t = useTranslations('listing');
  const apiError = useApiError();
  const listing = server.listing ?? EMPTY_LISTING;
  const published = server.publish_status?.state === 'published';
  // 没绑定 arkfgf.com 账号时没有统计可查，不发请求
  const loginRequired = server.publish_status?.state === 'login_required';

  // stats 为 null 表示还没有发布过；undefined 表示加载中
  const [stats, setStats] = useState<PublishStats | null | undefined>(undefined);
  const [statsError, setStatsError] = useState('');
  const [editing, setEditing] = useState(false);

  const load = useCallback(async () => {
    if (loginRequired) {
      setStats(null);
      return;
    }
    try {
      const res = await axios.get(`/api/servers/${server.id}/publish/stats`);
      setStats(res.data.data);
      setStatsError('');
    } catch (err) {
      if (axios.isAxiosError(err) && err.response?.status === 409) {
        setStats(null);
      } else {
        setStatsError(apiError(err, t('statsFailed')));
      }
    }
  }, [server.id, loginRequired, apiError, t]);

  // 发布成功后（第一次发布才有令牌和统计）重新读取
  useEffect(() => {
    const timer = setTimeout(() => {
      load();
    }, 0);
    return () => clearTimeout(timer);
  }, [load, published]);

  const copy = async (text: string) => {
    try {
      await navigator.clipboard.writeText(text);
      toast.success(t('copied'));
    } catch {
      toast.error(t('copyFailed'));
    }
  };

  const snippets = stats
    ? [
        { key: 'image', code: stats.badge },
        { key: 'html', code: `<a href="${stats.url}"><img src="${stats.badge}" alt="${escapeHtml(server.session_name)}"></a>` },
        { key: 'markdown', code: `[![${server.session_name.replace(/[[\]]/g, '')}](${stats.badge})](${stats.url})` },
        { key: 'bbcode', code: `[url=${stats.url}][img]${stats.badge}[/img][/url]` },
      ]
    : [];

  return (
    <Card className="lg:col-span-2">
      <CardHeader className="flex flex-row items-center justify-between gap-2 pb-2">
        <CardTitle className="font-display text-base">{t('title')}</CardTitle>
        {stats?.url && (
          <a href={stats.url} target="_blank" rel="noreferrer" className="inline-flex items-center gap-1 text-sm text-primary hover:underline">
            {t('viewPage')}
            <ExternalLink className="size-3.5" />
          </a>
        )}
      </CardHeader>
      <CardContent className="grid gap-6 md:grid-cols-2">
        {/* 访问统计和徽章 */}
        <div className="space-y-4">
          {stats === undefined && !statsError && <p className="text-sm text-muted-foreground">{t('loading')}</p>}
          {statsError && <p className="text-sm text-red-700">{statsError}</p>}
          {stats === null && (
            <p className="text-sm text-muted-foreground">
              {server.publish_status?.state === 'login_required' ? t('loginRequired') : t('notPublished')}
            </p>
          )}
          {stats && (
            <>
              <div>
                <div className="mb-1.5 text-xs text-muted-foreground">{t('week')}</div>
                <div className="grid grid-cols-3 gap-px overflow-hidden rounded-md border bg-border">
                  {(['views', 'connects', 'copies'] as const).map((k) => (
                    <div key={k} className="bg-card px-3 py-2.5">
                      <div className="text-xs text-muted-foreground">{t(`stats.${k}`)}</div>
                      <div className="font-display text-xl font-semibold">{stats.week[k]}</div>
                    </div>
                  ))}
                </div>
                <p className="mt-1.5 text-xs text-muted-foreground">
                  {t('month', { views: stats.month.views, connects: stats.month.connects, copies: stats.month.copies })}
                </p>
              </div>
              <div>
                <div className="mb-1.5 text-xs text-muted-foreground">{t('badge')}</div>
                {/* 收录站生成的图片，带缓存，不用 next/image */}
                {/* eslint-disable-next-line @next/next/no-img-element */}
                <img src={stats.badge} alt={t('badgeAlt')} width={260} height={48} className="mb-2 rounded" />
                <div className="flex flex-wrap gap-1.5">
                  {snippets.map((s) => (
                    <Button key={s.key} variant="outline" size="sm" onClick={() => copy(s.code)}>
                      <Copy className="mr-1.5 size-3.5" />
                      {t(`snippets.${s.key}`)}
                    </Button>
                  ))}
                </div>
                <p className="mt-1.5 text-xs text-muted-foreground">{t('badgeHint')}</p>
              </div>
            </>
          )}
        </div>

        {/* 收录资料 */}
        <div className="space-y-2 text-sm">
          <div className="flex items-center justify-between gap-2">
            <span className="text-xs text-muted-foreground">{t('profile')}</span>
            <Button variant="outline" size="sm" onClick={() => setEditing(true)}>
              <Pencil className="mr-1.5 size-3.5" />
              {t('edit')}
            </Button>
          </div>
          {listing.description || listing.tags.length || listing.qq_group || listing.discord ? (
            <>
              {listing.tags.length > 0 && (
                <div className="flex flex-wrap gap-1">
                  {listing.tags.map((tag) => (
                    <span key={tag} className="rounded border px-1.5 py-0.5 text-xs">
                      {tag}
                    </span>
                  ))}
                </div>
              )}
              {listing.description && <p className="line-clamp-4 whitespace-pre-line text-muted-foreground">{listing.description}</p>}
              {listing.qq_group && <p>{t('qqValue', { qq: listing.qq_group })}</p>}
              {listing.discord && <p className="break-all">{listing.discord}</p>}
            </>
          ) : (
            <p className="text-muted-foreground">{t('profileEmpty')}</p>
          )}
        </div>
      </CardContent>
      <ListingDialog server={server} open={editing} onClose={() => setEditing(false)} />
    </Card>
  );
}

function escapeHtml(s: string) {
  return s.replace(/[&<>"]/g, (c) => ({ '&': '&amp;', '<': '&lt;', '>': '&gt;', '"': '&quot;' })[c]!);
}

function ListingDialog({ server, open, onClose }: { server: Server; open: boolean; onClose: () => void }) {
  const t = useTranslations('listing');
  const apiError = useApiError();
  const current = server.listing ?? EMPTY_LISTING;
  const [form, setForm] = useState({ description: '', tags: '', qq_group: '', discord: '' });
  const [saving, setSaving] = useState(false);

  // 每次打开时用当前资料填表（在渲染期间按 props 调整状态，不在 effect 里 setState）
  const [openedFor, setOpenedFor] = useState(false);
  if (open !== openedFor) {
    setOpenedFor(open);
    if (open) {
      setForm({
        description: current.description,
        tags: current.tags.join(', '),
        qq_group: current.qq_group,
        discord: current.discord,
      });
    }
  }

  const save = async () => {
    setSaving(true);
    try {
      const listing: ListingProfile = {
        description: form.description,
        tags: form.tags.split(/[,，]/).map((s) => s.trim()).filter(Boolean),
        qq_group: form.qq_group.trim(),
        discord: form.discord.trim(),
      };
      await serversActions.updateServer(server.id, { listing });
      toast.success(server.publish_status?.state === 'published' ? t('savedRepublish') : t('saved'));
      onClose();
    } catch (err) {
      toast.error(apiError(err, t('saveFailed')));
    } finally {
      setSaving(false);
    }
  };

  return (
    <Dialog open={open} onOpenChange={(o) => !o && !saving && onClose()}>
      <DialogContent className="max-w-lg">
        <DialogHeader>
          <DialogTitle>{t('editTitle')}</DialogTitle>
          <DialogDescription>{t('editDesc')}</DialogDescription>
        </DialogHeader>
        <div className="space-y-4">
          <div className="space-y-1">
            <Label htmlFor="listing-description">{t('description')}</Label>
            <Textarea
              id="listing-description"
              value={form.description}
              maxLength={DESCRIPTION_MAX}
              rows={5}
              onChange={(e) => setForm((f) => ({ ...f, description: e.target.value }))}
              placeholder={t('descriptionPlaceholder')}
            />
            <p className="text-right text-xs text-muted-foreground">
              {form.description.length}/{DESCRIPTION_MAX}
            </p>
          </div>
          <div className="space-y-1">
            <Label htmlFor="listing-tags">{t('tags')}</Label>
            <Input
              id="listing-tags"
              value={form.tags}
              onChange={(e) => setForm((f) => ({ ...f, tags: e.target.value }))}
              placeholder={t('tagsPlaceholder')}
            />
            <p className="text-xs text-muted-foreground">{t('tagsHint')}</p>
          </div>
          <div className="grid gap-4 sm:grid-cols-2">
            <div className="space-y-1">
              <Label htmlFor="listing-qq">{t('qq')}</Label>
              <Input
                id="listing-qq"
                inputMode="numeric"
                value={form.qq_group}
                onChange={(e) => setForm((f) => ({ ...f, qq_group: e.target.value }))}
                placeholder="123456789"
              />
            </div>
            <div className="space-y-1">
              <Label htmlFor="listing-discord">Discord</Label>
              <Input
                id="listing-discord"
                value={form.discord}
                onChange={(e) => setForm((f) => ({ ...f, discord: e.target.value }))}
                placeholder="https://discord.gg/xxxx"
              />
            </div>
          </div>
        </div>
        <DialogFooter>
          <Button variant="outline" onClick={onClose} disabled={saving}>
            {t('cancel')}
          </Button>
          <Button onClick={save} disabled={saving}>
            {saving && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
            {t('save')}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
