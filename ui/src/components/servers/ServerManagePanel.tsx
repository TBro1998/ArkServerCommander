"use client";

import { useState, type ReactNode } from 'react';
import { useTranslations } from 'next-intl';
import axios from 'axios';
import { Copy, Download, FileDown, FileUp, RotateCcw, Trash2 } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { serversActions } from '@/stores/servers';
import { useApiError } from '@/lib/api-error';
import { toast } from '@/lib/toast';
import { downloadServerConfig, type ConfigPart, type ServerConfig } from '@/lib/server-config';
import { CopyConfigDialog, ImportConfigDialog } from '@/components/servers/ServerConfigDialogs';
import type { ServerControls } from '@/components/servers/ServerActions';
import { cn } from '@/lib/utils';

// 一项操作：左边说明，右边按钮；按钮禁用时把原因放在 title 上
function ActionRow({
  title,
  desc,
  note,
  children,
}: {
  title: string;
  desc: string;
  note?: ReactNode;
  children: ReactNode;
}) {
  return (
    <div className="flex flex-col gap-3 py-3.5 first:pt-0 last:pb-0 sm:flex-row sm:items-center sm:justify-between">
      <div className="min-w-0">
        <div className="text-sm font-medium">{title}</div>
        <p className="mt-0.5 text-xs text-gray-500">{desc}</p>
        {note}
      </div>
      <div className="flex shrink-0 flex-wrap gap-2">{children}</div>
    </div>
  );
}

/** 详情页"管理"标签：不常用或有风险的操作，每项都带说明 */
export function ServerManagePanel({ controls }: { controls: ServerControls }) {
  const t = useTranslations('manage');
  const tActions = useTranslations('serverActions');
  const tConfig = useTranslations('serverConfig');
  const apiError = useApiError();
  const { server, busy, stopped } = controls;
  const [importOpen, setImportOpen] = useState(false);
  const [copyOpen, setCopyOpen] = useState(false);

  const exportConfig = async () => {
    try {
      await downloadServerConfig(server);
    } catch (err) {
      toast.error(apiError(err, tConfig('exportFailed')));
    }
  };

  const importConfig = async (config: ServerConfig, parts: ConfigPart[]) => {
    try {
      const res = await axios.post(`/api/servers/${server.id}/config/import`, { config, parts });
      toast.success(res.data.data?.status === 'stopped' ? t('importSaved') : t('importPending'));
      serversActions.fetchServers().catch(() => {});
    } catch (err) {
      toast.error(apiError(err, t('importFailed')));
      throw err;
    }
  };

  const outdated = server.mods_outdated?.length ?? 0;
  const busyTitle = busy ? t('busy') : undefined;

  return (
    <div className="grid gap-4 lg:grid-cols-2">
      <Card>
        <CardHeader className="pb-2">
          <CardTitle className="font-display text-base">{t('maintenance')}</CardTitle>
        </CardHeader>
        <CardContent className="divide-y">
          <ActionRow
            title={tActions('updateGame')}
            desc={t('updateDesc')}
            note={outdated > 0 && <p className="mt-1 text-xs text-amber-800">{t('modsOutdated', { count: outdated })}</p>}
          >
            <span title={busyTitle}>
              <Button variant="outline" size="sm" onClick={() => controls.open('update')} disabled={busy}>
                <Download />
                {t('update')}
              </Button>
            </span>
          </ActionRow>
          <ActionRow title={tActions('recreate')} desc={t('recreateDesc')}>
            <span title={busyTitle}>
              <Button variant="outline" size="sm" onClick={() => controls.open('recreate')} disabled={busy}>
                <RotateCcw />
                {t('recreate')}
              </Button>
            </span>
          </ActionRow>
        </CardContent>
      </Card>

      <Card>
        <CardHeader className="pb-2">
          <CardTitle className="font-display text-base">{t('config')}</CardTitle>
        </CardHeader>
        <CardContent className="divide-y">
          <ActionRow title={tConfig('export')} desc={t('exportDesc')}>
            <Button variant="outline" size="sm" onClick={exportConfig}>
              <FileDown />
              {t('export')}
            </Button>
          </ActionRow>
          <ActionRow title={tConfig('import')} desc={t('importDesc')}>
            <Button variant="outline" size="sm" onClick={() => setImportOpen(true)}>
              <FileUp />
              {t('import')}
            </Button>
          </ActionRow>
          <ActionRow title={tConfig('copyTo')} desc={t('copyDesc')}>
            <Button variant="outline" size="sm" onClick={() => setCopyOpen(true)}>
              <Copy />
              {t('copy')}
            </Button>
          </ActionRow>
        </CardContent>
      </Card>

      <Card className="border-red-200 lg:col-span-2">
        <CardHeader className="pb-2">
          <CardTitle className="font-display text-base text-red-700">{t('danger')}</CardTitle>
        </CardHeader>
        <CardContent>
          <ActionRow title={tActions('delete')} desc={t('deleteDesc')}>
            <span title={stopped ? undefined : t('deleteNeedStop')}>
              <Button
                variant="outline"
                size="sm"
                className={cn('border-red-200 text-red-700 hover:bg-red-50 hover:text-red-800')}
                onClick={() => controls.open('delete')}
                disabled={!stopped || busy}
              >
                <Trash2 />
                {t('delete')}
              </Button>
            </span>
          </ActionRow>
        </CardContent>
      </Card>

      <ImportConfigDialog open={importOpen} onOpenChange={setImportOpen} excludeServerId={server.id} mode="apply" onImport={importConfig} />
      <CopyConfigDialog open={copyOpen} onOpenChange={setCopyOpen} server={server} />
    </div>
  );
}
