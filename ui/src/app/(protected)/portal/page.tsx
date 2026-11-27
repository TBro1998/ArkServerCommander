"use client";

import { useCallback, useEffect, useState } from 'react';
import { useLocale, useTranslations } from 'next-intl';
import axios from 'axios';
import { ExternalLink, Loader2, Megaphone, Pencil, Pin, Plus, Store, Trash2, X } from 'lucide-react';
import { PageHeader } from '@/components/PageHeader';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Switch } from '@/components/ui/switch';
import { Textarea } from '@/components/ui/textarea';
import { Badge } from '@/components/ui/badge';
import { PageLoader } from '@/components/ui/page-loader';
import { ConfirmDialog } from '@/components/ui/confirm-dialog';
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from '@/components/ui/dialog';
import { GameBadge } from '@/components/servers/GameBadge';
import type { Game } from '@/lib/games';
import { usePolling } from '@/hooks/use-polling';
import { useApiError } from '@/lib/api-error';
import { formatDateTime } from '@/lib/format';
import { toast } from '@/lib/toast';

interface PortalLink {
  title: string;
  url: string;
}

interface PortalInfo {
  settings: {
    enabled: boolean;
    port: number;
    public_url: string;
    name: string;
    description: string;
    qq_group: string;
    discord: string;
  };
  links: PortalLink[];
  status: { state: 'disabled' | 'starting' | 'running' | 'error'; url?: string; error_code?: string; error?: string };
  servers: { id: number; name: string; game: Game; hidden: boolean; shop: boolean; shop_error?: string }[];
  image: string;
}

interface Announcement {
  id: number;
  title: string;
  body: string;
  pinned: boolean;
  published_at: string;
  updated_at: string;
}

type Form = PortalInfo['settings'] & { links: PortalLink[]; hidden: number[] };

function toForm(info: PortalInfo): Form {
  return { ...info.settings, links: info.links, hidden: info.servers.filter((s) => s.hidden).map((s) => s.id) };
}

export default function PortalPage() {
  const t = useTranslations('portal');
  const apiError = useApiError();
  const [info, setInfo] = useState<PortalInfo | null>(null);
  const [form, setForm] = useState<Form | null>(null);
  const [dirty, setDirty] = useState(false);
  const [saving, setSaving] = useState(false);
  const [loadError, setLoadError] = useState('');

  const load = useCallback(
    async (resetForm: boolean) => {
      try {
        const res = await axios.get<{ data: PortalInfo }>('/api/portal');
        setInfo(res.data.data);
        if (resetForm) setForm(toForm(res.data.data));
        setLoadError('');
      } catch (err) {
        setLoadError(apiError(err, t('loadFailed')));
      }
    },
    [apiError, t],
  );

  useEffect(() => {
    const timer = setTimeout(() => load(true), 0);
    return () => clearTimeout(timer);
  }, [load]);
  // 状态（容器启动、镜像下载）每 5 秒刷新，不覆盖正在编辑的表单
  usePolling(() => load(!dirty), 5000);

  const update = (patch: Partial<Form>) => {
    setForm((f) => (f ? { ...f, ...patch } : f));
    setDirty(true);
  };

  const save = async (next: Form) => {
    setSaving(true);
    try {
      const res = await axios.put<{ data: PortalInfo }>('/api/portal', {
        enabled: next.enabled,
        port: Number(next.port) || 0,
        public_url: next.public_url,
        name: next.name,
        description: next.description,
        qq_group: next.qq_group,
        discord: next.discord,
        links: next.links,
        hidden_servers: next.hidden,
      });
      setInfo(res.data.data);
      setForm(toForm(res.data.data));
      setDirty(false);
      toast.success(t('saved'));
      return true;
    } catch (err) {
      toast.error(apiError(err, t('saveFailed')));
      return false;
    } finally {
      setSaving(false);
    }
  };

  if (!info || !form) {
    return loadError ? <div className="text-sm text-red-700">{loadError}</div> : <PageLoader />;
  }

  return (
    <div className="w-full max-w-5xl space-y-6">
      <PageHeader title={t('title')} description={t('description')} />

      <StatusCard info={info} form={form} saving={saving} dirty={dirty} onToggle={(enabled) => save({ ...form, enabled })} />

      <div className="grid gap-6 lg:grid-cols-[minmax(0,1fr)_minmax(0,1fr)]">
        <Card>
          <CardHeader>
            <CardTitle className="text-lg">{t('siteTitle')}</CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            <Field label={t('name')} hint={t('nameHint')}>
              <Input value={form.name} maxLength={64} onChange={(e) => update({ name: e.target.value })} />
            </Field>
            <Field label={t('siteDescription')}>
              <Textarea value={form.description} maxLength={2000} rows={4} onChange={(e) => update({ description: e.target.value })} />
            </Field>
            <div className="grid gap-4 sm:grid-cols-2">
              <Field label={t('qqGroup')}>
                <Input value={form.qq_group} inputMode="numeric" onChange={(e) => update({ qq_group: e.target.value })} />
              </Field>
              <Field label="Discord">
                <Input value={form.discord} placeholder="https://discord.gg/..." onChange={(e) => update({ discord: e.target.value })} />
              </Field>
            </div>
            <div className="space-y-2">
              <Label>{t('links')}</Label>
              {form.links.map((l, i) => (
                <div key={i} className="flex gap-2">
                  <Input
                    className="w-1/3"
                    value={l.title}
                    placeholder={t('linkTitle')}
                    onChange={(e) => update({ links: form.links.map((x, j) => (j === i ? { ...x, title: e.target.value } : x)) })}
                  />
                  <Input
                    value={l.url}
                    placeholder="https://"
                    onChange={(e) => update({ links: form.links.map((x, j) => (j === i ? { ...x, url: e.target.value } : x)) })}
                  />
                  <Button variant="ghost" size="icon" aria-label={t('removeLink')} onClick={() => update({ links: form.links.filter((_, j) => j !== i) })}>
                    <X />
                  </Button>
                </div>
              ))}
              {form.links.length < 10 && (
                <Button variant="outline" size="sm" onClick={() => update({ links: [...form.links, { title: '', url: '' }] })}>
                  <Plus />
                  {t('addLink')}
                </Button>
              )}
            </div>
          </CardContent>
        </Card>

        <div className="space-y-6">
          <Card>
            <CardHeader>
              <CardTitle className="text-lg">{t('access')}</CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              <Field label={t('port')} hint={t('portHint')}>
                <Input type="number" min={1024} max={65535} value={form.port} onChange={(e) => update({ port: Number(e.target.value) })} className="w-32" />
              </Field>
              <Field label={t('publicUrl')} hint={t('publicUrlHint')}>
                <Input value={form.public_url} placeholder="https://ark.example.com" onChange={(e) => update({ public_url: e.target.value })} />
              </Field>
            </CardContent>
          </Card>

          <Card>
            <CardHeader>
              <CardTitle className="text-lg">{t('servers')}</CardTitle>
            </CardHeader>
            <CardContent>
              {info.servers.length === 0 ? (
                <p className="text-sm text-muted-foreground">{t('noServers')}</p>
              ) : (
                <ul className="divide-y divide-border">
                  {info.servers.map((s) => {
                    const shown = !form.hidden.includes(s.id);
                    return (
                      <li key={s.id} className="flex items-center gap-3 py-2.5">
                        <GameBadge game={s.game} />
                        <div className="min-w-0 flex-1">
                          <div className="truncate text-sm">{s.name}</div>
                          {s.shop ? (
                            <div className="mt-0.5 flex items-center gap-1 text-xs text-green-800">
                              <Store className="size-3" />
                              {t('hasShop')}
                            </div>
                          ) : s.shop_error ? (
                            <div className="mt-0.5 text-xs text-red-700">{s.shop_error}</div>
                          ) : null}
                        </div>
                        <Switch
                          checked={shown}
                          aria-label={t('showServer')}
                          onCheckedChange={(v) => update({ hidden: v ? form.hidden.filter((id) => id !== s.id) : [...form.hidden, s.id] })}
                        />
                      </li>
                    );
                  })}
                </ul>
              )}
              <p className="mt-3 text-xs text-muted-foreground">{t('shopHint')}</p>
            </CardContent>
          </Card>
        </div>
      </div>

      {dirty && (
        <div className="sticky bottom-4 z-10 flex justify-end gap-2 rounded-lg border border-primary/30 bg-card p-3 shadow-lg">
          <span className="mr-auto self-center text-sm text-muted-foreground">{t('unsaved')}</span>
          <Button
            variant="ghost"
            onClick={() => {
              setForm(toForm(info));
              setDirty(false);
            }}
          >
            {t('discard')}
          </Button>
          <Button onClick={() => save(form)} disabled={saving}>
            {saving && <Loader2 className="animate-spin" />}
            {t('save')}
          </Button>
        </div>
      )}

      <AnnouncementsCard />
    </div>
  );
}

function Field({ label, hint, children }: { label: string; hint?: string; children: React.ReactNode }) {
  return (
    <div className="space-y-1.5">
      <Label>{label}</Label>
      {children}
      {hint && <p className="text-xs text-muted-foreground">{hint}</p>}
    </div>
  );
}

function StatusCard({
  info,
  form,
  saving,
  dirty,
  onToggle,
}: {
  info: PortalInfo;
  form: Form;
  saving: boolean;
  dirty: boolean;
  onToggle: (enabled: boolean) => void;
}) {
  const t = useTranslations('portal');
  const st = info.status;
  const enabled = info.settings.enabled;
  const tone =
    st.state === 'running' ? 'bg-green-100 text-green-800' : st.state === 'error' ? 'bg-red-100 text-red-800' : st.state === 'starting' ? 'bg-amber-100 text-amber-800' : 'bg-stone-100 text-stone-700';
  return (
    <Card className="hud">
      <CardContent className="flex flex-wrap items-center gap-4 py-5">
        <Switch checked={enabled} disabled={saving || dirty} onCheckedChange={onToggle} aria-label={t('enable')} />
        <div className="min-w-0 flex-1">
          <div className="flex flex-wrap items-center gap-2">
            <span className="font-medium">{t('enable')}</span>
            <Badge className={tone}>{t(`state.${st.state}`)}</Badge>
          </div>
          <p className="mt-1 text-sm text-muted-foreground">
            {dirty ? t('saveFirst') : enabled ? t('enabledHint') : t('disabledHint', { port: form.port })}
          </p>
          {st.state === 'error' && st.error && <p className="mt-1 text-sm text-red-700">{st.error}</p>}
        </div>
        {st.url && enabled && (
          <Button variant="outline" asChild>
            <a href={st.url} target="_blank" rel="noreferrer noopener">
              <ExternalLink />
              {st.url}
            </a>
          </Button>
        )}
      </CardContent>
    </Card>
  );
}

function AnnouncementsCard() {
  const t = useTranslations('portal');
  const locale = useLocale();
  const apiError = useApiError();
  const [list, setList] = useState<Announcement[] | null>(null);
  const [editing, setEditing] = useState<Partial<Announcement> | null>(null);
  const [deleting, setDeleting] = useState<Announcement | null>(null);
  const [busy, setBusy] = useState(false);

  const load = useCallback(async () => {
    try {
      const res = await axios.get<{ data: Announcement[] }>('/api/portal/announcements');
      setList(res.data.data);
    } catch (err) {
      toast.error(apiError(err, t('loadFailed')));
    }
  }, [apiError, t]);

  useEffect(() => {
    const timer = setTimeout(load, 0);
    return () => clearTimeout(timer);
  }, [load]);

  const submit = async () => {
    if (!editing) return;
    setBusy(true);
    try {
      const body = { title: editing.title ?? '', body: editing.body ?? '', pinned: !!editing.pinned };
      if (editing.id) await axios.put(`/api/portal/announcements/${editing.id}`, body);
      else await axios.post('/api/portal/announcements', body);
      toast.success(editing.id ? t('announcementUpdated') : t('announcementPublished'));
      setEditing(null);
      load();
    } catch (err) {
      toast.error(apiError(err, t('saveFailed')));
    } finally {
      setBusy(false);
    }
  };

  const sorted = [...(list ?? [])].sort((a, b) => Number(b.pinned) - Number(a.pinned) || b.published_at.localeCompare(a.published_at));

  return (
    <Card>
      <CardHeader className="flex flex-row items-center justify-between space-y-0">
        <CardTitle className="flex items-center gap-2 text-lg">
          <Megaphone className="size-5 text-primary" />
          {t('announcements')}
        </CardTitle>
        <Button size="sm" onClick={() => setEditing({ title: '', body: '', pinned: false })}>
          <Plus />
          {t('newAnnouncement')}
        </Button>
      </CardHeader>
      <CardContent>
        {list === null ? (
          <PageLoader />
        ) : sorted.length === 0 ? (
          <p className="py-6 text-center text-sm text-muted-foreground">{t('noAnnouncements')}</p>
        ) : (
          <ul className="space-y-3">
            {sorted.map((a) => (
              <li key={a.id} className="rounded-md border border-border p-3">
                <div className="flex items-start gap-2">
                  {a.pinned && <Pin className="mt-0.5 size-4 shrink-0 text-primary" aria-label={t('pinned')} />}
                  <div className="min-w-0 flex-1">
                    <div className="font-medium break-words">{a.title}</div>
                    <div className="text-xs text-muted-foreground">{formatDateTime(a.published_at, locale)}</div>
                  </div>
                  <Button variant="ghost" size="icon" aria-label={t('edit')} onClick={() => setEditing(a)}>
                    <Pencil />
                  </Button>
                  <Button variant="ghost" size="icon" aria-label={t('delete')} onClick={() => setDeleting(a)}>
                    <Trash2 />
                  </Button>
                </div>
                {a.body && <p className="mt-2 line-clamp-3 whitespace-pre-line text-sm text-muted-foreground">{a.body}</p>}
              </li>
            ))}
          </ul>
        )}
      </CardContent>

      <Dialog open={!!editing} onOpenChange={(o) => !o && setEditing(null)}>
        <DialogContent className="sm:max-w-xl">
          <DialogHeader>
            <DialogTitle>{editing?.id ? t('editAnnouncement') : t('newAnnouncement')}</DialogTitle>
            <DialogDescription>{t('announcementHint')}</DialogDescription>
          </DialogHeader>
          <div className="space-y-4">
            <Field label={t('announcementTitle')}>
              <Input value={editing?.title ?? ''} maxLength={100} onChange={(e) => setEditing((x) => ({ ...x, title: e.target.value }))} />
            </Field>
            <Field label={t('announcementBody')}>
              <Textarea value={editing?.body ?? ''} maxLength={10000} rows={8} onChange={(e) => setEditing((x) => ({ ...x, body: e.target.value }))} />
            </Field>
            <label className="flex items-center gap-2 text-sm">
              <Switch checked={!!editing?.pinned} onCheckedChange={(v) => setEditing((x) => ({ ...x, pinned: v }))} />
              {t('pinned')}
            </label>
          </div>
          <DialogFooter>
            <Button variant="ghost" onClick={() => setEditing(null)}>
              {t('cancel')}
            </Button>
            <Button onClick={submit} disabled={busy || !editing?.title?.trim()}>
              {busy && <Loader2 className="animate-spin" />}
              {editing?.id ? t('save') : t('publish')}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      <ConfirmDialog
        open={!!deleting}
        onOpenChange={(o) => !o && setDeleting(null)}
        title={t('deleteTitle')}
        description={t('deleteDesc', { title: deleting?.title ?? '' })}
        confirmText={t('delete')}
        destructive
        onConfirm={async () => {
          if (!deleting) return;
          try {
            await axios.delete(`/api/portal/announcements/${deleting.id}`);
            setDeleting(null);
            load();
          } catch (err) {
            toast.error(apiError(err, t('saveFailed')));
          }
        }}
      />
    </Card>
  );
}
