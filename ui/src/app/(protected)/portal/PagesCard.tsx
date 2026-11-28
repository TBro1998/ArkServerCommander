"use client";

import { useCallback, useEffect, useState } from 'react';
import { useTranslations } from 'next-intl';
import axios from 'axios';
import { ArrowDown, ArrowUp, BookOpen, Loader2, Pencil, Plus, Trash2 } from 'lucide-react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Textarea } from '@/components/ui/textarea';
import { PageLoader } from '@/components/ui/page-loader';
import { ConfirmDialog } from '@/components/ui/confirm-dialog';
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from '@/components/ui/dialog';
import { useApiError } from '@/lib/api-error';
import { toast } from '@/lib/toast';

interface PortalPage {
  id: number;
  title: string;
  body: string;
}

// 玩家站上服主写的页面（规则、新手指南等）：纯文本，「# 」开头是标题、「- 」开头是列表
export function PagesCard() {
  const t = useTranslations('portal');
  const apiError = useApiError();
  const [list, setList] = useState<PortalPage[] | null>(null);
  const [editing, setEditing] = useState<Partial<PortalPage> | null>(null);
  const [deleting, setDeleting] = useState<PortalPage | null>(null);
  const [busy, setBusy] = useState(false);

  const load = useCallback(async () => {
    try {
      const res = await axios.get<{ data: PortalPage[] }>('/api/portal/pages');
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
      const body = { title: editing.title ?? '', body: editing.body ?? '' };
      if (editing.id) await axios.put(`/api/portal/pages/${editing.id}`, body);
      else await axios.post('/api/portal/pages', body);
      toast.success(t('pageSaved'));
      setEditing(null);
      load();
    } catch (err) {
      toast.error(apiError(err, t('saveFailed')));
    } finally {
      setBusy(false);
    }
  };

  const move = async (p: PortalPage, delta: number) => {
    try {
      await axios.post(`/api/portal/pages/${p.id}/move`, { delta });
      load();
    } catch (err) {
      toast.error(apiError(err, t('saveFailed')));
    }
  };

  return (
    <Card>
      <CardHeader className="flex flex-row items-center justify-between space-y-0">
        <CardTitle className="flex items-center gap-2 text-lg">
          <BookOpen className="size-5 text-primary" />
          {t('pages')}
        </CardTitle>
        <Button size="sm" onClick={() => setEditing({ title: '', body: '' })}>
          <Plus />
          {t('newPage')}
        </Button>
      </CardHeader>
      <CardContent>
        <p className="mb-3 text-xs text-muted-foreground">{t('pagesHint')}</p>
        {list === null ? (
          <PageLoader />
        ) : list.length === 0 ? (
          <p className="py-6 text-center text-sm text-muted-foreground">{t('noPages')}</p>
        ) : (
          <ul className="divide-y divide-border">
            {list.map((p, i) => (
              <li key={p.id} className="flex items-center gap-2 py-2">
                <span className="min-w-0 flex-1 truncate text-sm">{p.title}</span>
                <Button variant="ghost" size="icon" aria-label={t('moveUp')} disabled={i === 0} onClick={() => move(p, -1)}>
                  <ArrowUp />
                </Button>
                <Button variant="ghost" size="icon" aria-label={t('moveDown')} disabled={i === list.length - 1} onClick={() => move(p, 1)}>
                  <ArrowDown />
                </Button>
                <Button variant="ghost" size="icon" aria-label={t('edit')} onClick={() => setEditing(p)}>
                  <Pencil />
                </Button>
                <Button variant="ghost" size="icon" aria-label={t('delete')} onClick={() => setDeleting(p)}>
                  <Trash2 />
                </Button>
              </li>
            ))}
          </ul>
        )}
      </CardContent>

      <Dialog open={!!editing} onOpenChange={(o) => !o && setEditing(null)}>
        <DialogContent className="sm:max-w-2xl">
          <DialogHeader>
            <DialogTitle>{editing?.id ? t('editPage') : t('newPage')}</DialogTitle>
            <DialogDescription>{t('pageFormatHint')}</DialogDescription>
          </DialogHeader>
          <div className="space-y-4">
            <div className="space-y-1.5">
              <Label>{t('pageTitle')}</Label>
              <Input value={editing?.title ?? ''} maxLength={64} onChange={(e) => setEditing((x) => ({ ...x, title: e.target.value }))} />
            </div>
            <div className="space-y-1.5">
              <Label>{t('pageBody')}</Label>
              <Textarea
                value={editing?.body ?? ''}
                maxLength={20000}
                rows={14}
                className="font-mono text-sm"
                placeholder={t('pageBodyPlaceholder')}
                onChange={(e) => setEditing((x) => ({ ...x, body: e.target.value }))}
              />
            </div>
          </div>
          <DialogFooter>
            <Button variant="ghost" onClick={() => setEditing(null)}>
              {t('cancel')}
            </Button>
            <Button onClick={submit} disabled={busy || !editing?.title?.trim()}>
              {busy && <Loader2 className="animate-spin" />}
              {t('save')}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      <ConfirmDialog
        open={!!deleting}
        onOpenChange={(o) => !o && setDeleting(null)}
        title={t('deletePageTitle')}
        description={t('deletePageDesc', { title: deleting?.title ?? '' })}
        confirmText={t('delete')}
        destructive
        onConfirm={async () => {
          if (!deleting) return;
          try {
            await axios.delete(`/api/portal/pages/${deleting.id}`);
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
