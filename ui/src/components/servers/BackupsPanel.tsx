"use client";

import { useCallback, useEffect, useRef, useState } from 'react';
import { useLocale, useTranslations } from 'next-intl';
import axios from 'axios';
import { Archive, Download, Loader2, RotateCcw, Trash2, Upload } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { ConfirmDialog } from '@/components/ui/confirm-dialog';
import { useApiError } from '@/lib/api-error';
import { formatBytes, formatDateTime } from '@/lib/format';
import { toast } from '@/lib/toast';
import type { Server } from '@/stores/servers';

interface Backup {
  name: string;
  kind: 'manual' | 'scheduled' | 'pre-restore' | 'uploaded';
  size: number;
  created_at: string;
}

const KIND_STYLES: Record<Backup['kind'], string> = {
  manual: 'bg-blue-100 text-blue-800',
  scheduled: 'bg-green-100 text-green-800',
  'pre-restore': 'bg-amber-100 text-amber-800',
  uploaded: 'bg-purple-100 text-purple-800',
};

export function BackupsPanel({ server, active, onRestored }: { server: Server; active: boolean; onRestored?: () => void }) {
  const t = useTranslations('backups');
  const locale = useLocale();
  const apiError = useApiError();
  const [backups, setBackups] = useState<Backup[] | null>(null);
  const [creating, setCreating] = useState(false);
  const [uploading, setUploading] = useState(false);
  const [restoreTarget, setRestoreTarget] = useState<Backup | null>(null);
  const [deleteTarget, setDeleteTarget] = useState<Backup | null>(null);
  const fileRef = useRef<HTMLInputElement>(null);

  const load = useCallback(async () => {
    try {
      const res = await axios.get(`/api/servers/${server.id}/backups`);
      setBackups(res.data.data ?? []);
    } catch (err) {
      toast.error(apiError(err, t('loadFailed')));
    }
  }, [server.id, apiError, t]);

  useEffect(() => {
    if (!active) return;
    const timer = setTimeout(() => {
      load();
    }, 0);
    return () => clearTimeout(timer);
  }, [active, load]);

  const create = async () => {
    setCreating(true);
    try {
      await axios.post(`/api/servers/${server.id}/backups`);
      toast.success(t('created'));
      load();
    } catch (err) {
      toast.error(apiError(err, t('createFailed')));
    } finally {
      setCreating(false);
    }
  };

  const upload = async (file: File) => {
    setUploading(true);
    const form = new FormData();
    form.append('file', file);
    try {
      await axios.post(`/api/servers/${server.id}/backups/upload`, form);
      toast.success(t('uploaded'));
      load();
    } catch (err) {
      toast.error(apiError(err, t('uploadFailed')));
    } finally {
      setUploading(false);
      if (fileRef.current) fileRef.current.value = '';
    }
  };

  const restore = async () => {
    if (!restoreTarget) return;
    try {
      await axios.post(`/api/servers/${server.id}/backups/${encodeURIComponent(restoreTarget.name)}/restore`);
      toast.success(t('restored'));
      load();
      onRestored?.();
    } catch (err) {
      toast.error(apiError(err, t('restoreFailed')));
      throw err;
    }
  };

  const remove = async () => {
    if (!deleteTarget) return;
    try {
      await axios.delete(`/api/servers/${server.id}/backups/${encodeURIComponent(deleteTarget.name)}`);
      toast.success(t('deleted'));
      load();
    } catch (err) {
      toast.error(apiError(err, t('deleteFailed')));
      throw err;
    }
  };

  const stopped = server.status === 'stopped';
  const busy = server.status === 'starting' || server.status === 'stopping';
  const totalSize = (backups ?? []).reduce((sum, b) => sum + b.size, 0);

  return (
    <Card>
      <CardHeader className="flex flex-col gap-3 pb-3 sm:flex-row sm:items-center sm:justify-between">
        <div>
          <CardTitle className="text-base">{t('title')}</CardTitle>
          <p className="mt-1 text-xs text-gray-500">
            {t('summary', { count: backups?.length ?? 0, size: formatBytes(totalSize) })}
          </p>
        </div>
        <div className="flex gap-2">
          <input
            ref={fileRef}
            type="file"
            accept=".gz,.tar.gz,application/gzip"
            className="hidden"
            onChange={(e) => {
              const file = e.target.files?.[0];
              if (file) upload(file);
            }}
          />
          <Button variant="outline" onClick={() => fileRef.current?.click()} disabled={uploading}>
            {uploading ? <Loader2 className="mr-2 h-4 w-4 animate-spin" /> : <Upload className="mr-2 h-4 w-4" />}
            {t('upload')}
          </Button>
          <Button onClick={create} disabled={creating || busy}>
            {creating ? <Loader2 className="mr-2 h-4 w-4 animate-spin" /> : <Archive className="mr-2 h-4 w-4" />}
            {t('create')}
          </Button>
        </div>
      </CardHeader>
      <CardContent>
        <p className="mb-3 rounded-lg bg-gray-50 p-3 text-xs text-gray-600">{t('hint')}</p>
        {backups === null ? (
          <div className="py-6 text-center">
            <Loader2 className="mx-auto h-6 w-6 animate-spin text-gray-400" />
          </div>
        ) : backups.length === 0 ? (
          <p className="py-6 text-center text-sm text-gray-500">{t('empty')}</p>
        ) : (
          <div className="overflow-x-auto">
            <table className="w-full text-sm">
              <thead>
                <tr className="border-b text-left text-xs text-gray-500">
                  <th className="py-2 pr-3 font-medium">{t('time')}</th>
                  <th className="py-2 pr-3 font-medium">{t('kind')}</th>
                  <th className="py-2 pr-3 font-medium">{t('size')}</th>
                  <th className="py-2 text-right font-medium">{t('actions')}</th>
                </tr>
              </thead>
              <tbody>
                {backups.map((b) => (
                  <tr key={b.name} className="border-b last:border-0">
                    <td className="py-2 pr-3 whitespace-nowrap">{formatDateTime(b.created_at, locale)}</td>
                    <td className="py-2 pr-3">
                      <span className={`rounded px-1.5 py-0.5 text-xs ${KIND_STYLES[b.kind]}`}>{t(`kinds.${b.kind}`)}</span>
                    </td>
                    <td className="py-2 pr-3 whitespace-nowrap">{formatBytes(b.size)}</td>
                    <td className="py-2 text-right whitespace-nowrap">
                      <Button variant="ghost" size="sm" className="h-7 w-7 p-0" asChild title={t('download')}>
                        <a href={`/api/servers/${server.id}/backups/${encodeURIComponent(b.name)}`} download>
                          <Download className="h-4 w-4" />
                        </a>
                      </Button>
                      <Button
                        variant="ghost"
                        size="sm"
                        className="h-7 w-7 p-0"
                        title={stopped ? t('restore') : t('restoreNeedStop')}
                        disabled={!stopped}
                        onClick={() => setRestoreTarget(b)}
                      >
                        <RotateCcw className="h-4 w-4" />
                      </Button>
                      <Button
                        variant="ghost"
                        size="sm"
                        className="h-7 w-7 p-0 text-red-600"
                        title={t('delete')}
                        onClick={() => setDeleteTarget(b)}
                      >
                        <Trash2 className="h-4 w-4" />
                      </Button>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </CardContent>

      <ConfirmDialog
        open={restoreTarget !== null}
        onOpenChange={(o) => !o && setRestoreTarget(null)}
        title={t('restoreTitle')}
        description={t('restoreDesc', { time: restoreTarget ? formatDateTime(restoreTarget.created_at, locale) : '' })}
        confirmText={t('restore')}
        destructive
        onConfirm={restore}
      />
      <ConfirmDialog
        open={deleteTarget !== null}
        onOpenChange={(o) => !o && setDeleteTarget(null)}
        title={t('deleteTitle')}
        description={t('deleteDesc', { time: deleteTarget ? formatDateTime(deleteTarget.created_at, locale) : '' })}
        confirmText={t('delete')}
        destructive
        onConfirm={remove}
      />
    </Card>
  );
}
