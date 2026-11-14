"use client";

import { useCallback, useEffect, useState, useSyncExternalStore } from 'react';
import { useLocale, useTranslations } from 'next-intl';
import axios from 'axios';
import { Bot, Copy, Loader2, Plus, Trash2 } from 'lucide-react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Switch } from '@/components/ui/switch';
import { ConfirmDialog } from '@/components/ui/confirm-dialog';
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from '@/components/ui/dialog';
import { useApiError } from '@/lib/api-error';
import { copyText, formatDateTime } from '@/lib/format';
import { toast } from '@/lib/toast';

interface ApiToken {
  id: number;
  name: string;
  prefix: string;
  read_only: boolean;
  last_used_at?: string | null;
  created_at: string;
}

interface CreatedToken extends ApiToken {
  token: string;
}

const SERVER_NAME = 'ark-commander';

// MCP 地址取浏览器当前访问的来源；服务端渲染时为空
const subscribeNoop = () => () => {};
const useMcpUrl = () =>
  useSyncExternalStore(
    subscribeNoop,
    () => `${window.location.origin}/mcp`,
    () => '',
  );

function CopyField({ value, label, multiline }: { value: string; label: string; multiline?: boolean }) {
  const t = useTranslations('settings.mcp');
  return (
    <div className="flex items-start gap-2">
      {multiline ? (
        <pre className="min-w-0 flex-1 overflow-x-auto rounded-sm bg-terminal px-2 py-1.5 font-mono text-xs whitespace-pre text-terminal-foreground">{value}</pre>
      ) : (
        <code className="min-w-0 flex-1 truncate rounded-sm bg-terminal px-2 py-1.5 font-mono text-xs text-terminal-foreground">{value}</code>
      )}
      <Button
        variant="outline"
        size="sm"
        aria-label={label}
        title={label}
        onClick={async () => {
          if (await copyText(value)) toast.success(t('copied'));
          else toast.error(t('copyFailed'));
        }}
      >
        <Copy />
      </Button>
    </div>
  );
}

// 设置页"Agent 接入（MCP）"：显示 MCP 地址，管理 API 令牌
export function McpCard() {
  const t = useTranslations('settings.mcp');
  const locale = useLocale();
  const apiError = useApiError();
  const mcpUrl = useMcpUrl();

  const [tokens, setTokens] = useState<ApiToken[] | null>(null);
  const [name, setName] = useState('');
  const [readOnly, setReadOnly] = useState(false);
  const [creating, setCreating] = useState(false);
  const [created, setCreated] = useState<CreatedToken | null>(null);
  const [deleteTarget, setDeleteTarget] = useState<ApiToken | null>(null);

  const load = useCallback(async () => {
    try {
      const res = await axios.get('/api/tokens');
      setTokens(res.data.data ?? []);
    } catch (err) {
      toast.error(apiError(err, t('loadFailed')));
    }
  }, [apiError, t]);

  useEffect(() => {
    const timer = setTimeout(() => {
      load();
    }, 0);
    return () => clearTimeout(timer);
  }, [load]);

  const create = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!name.trim()) return;
    setCreating(true);
    try {
      const res = await axios.post('/api/tokens', { name: name.trim(), read_only: readOnly });
      setCreated(res.data.data);
      setName('');
      setReadOnly(false);
      load();
    } catch (err) {
      toast.error(apiError(err, t('createFailed')));
    } finally {
      setCreating(false);
    }
  };

  const remove = async () => {
    if (!deleteTarget) return;
    try {
      await axios.delete(`/api/tokens/${deleteTarget.id}`);
      toast.success(t('deleted'));
      load();
    } catch (err) {
      toast.error(apiError(err, t('deleteFailed')));
      throw err;
    }
  };

  const claudeCommand = created
    ? `claude mcp add --transport http ${SERVER_NAME} ${mcpUrl} --header "Authorization: Bearer ${created.token}"`
    : '';
  const jsonConfig = created
    ? JSON.stringify(
        { mcpServers: { [SERVER_NAME]: { type: 'http', url: mcpUrl, headers: { Authorization: `Bearer ${created.token}` } } } },
        null,
        2,
      )
    : '';

  return (
    <Card>
      <CardHeader>
        <CardTitle className="flex items-center gap-2 text-lg">
          <Bot className="h-5 w-5" />
          {t('title')}
        </CardTitle>
      </CardHeader>
      <CardContent className="space-y-5 text-sm">
        <p className="text-xs text-muted-foreground">{t('description')}</p>

        <div className="space-y-1">
          <Label>{t('endpoint')}</Label>
          <CopyField value={mcpUrl} label={t('copy')} />
          <p className="text-xs text-muted-foreground">{t('endpointHint')}</p>
        </div>

        <form onSubmit={create} className="space-y-3">
          <div className="space-y-1">
            <Label htmlFor="token_name">{t('newToken')}</Label>
            <div className="flex gap-2">
              <Input
                id="token_name"
                value={name}
                maxLength={64}
                placeholder={t('namePlaceholder')}
                onChange={(e) => setName(e.target.value)}
              />
              <Button type="submit" disabled={creating || !name.trim()}>
                {creating ? <Loader2 className="mr-1.5 animate-spin" /> : <Plus className="mr-1.5" />}
                {t('create')}
              </Button>
            </div>
          </div>
          <div className="flex items-center gap-2">
            <Switch id="token_read_only" checked={readOnly} onCheckedChange={setReadOnly} />
            <Label htmlFor="token_read_only" className="font-normal">
              {t('readOnly')}
            </Label>
          </div>
          <p className="text-xs text-muted-foreground">{t('readOnlyHint')}</p>
        </form>

        <div className="space-y-2">
          <Label>{t('tokens')}</Label>
          {tokens === null ? (
            <div className="flex justify-center py-4">
              <Loader2 className="animate-spin text-muted-foreground" />
            </div>
          ) : tokens.length === 0 ? (
            <p className="py-2 text-xs text-muted-foreground">{t('empty')}</p>
          ) : (
            <ul className="divide-y divide-border rounded-sm border border-border">
              {tokens.map((tok) => (
                <li key={tok.id} className="flex items-center gap-3 px-3 py-2">
                  <div className="min-w-0 flex-1">
                    <div className="flex items-center gap-2">
                      <span className="truncate font-medium">{tok.name}</span>
                      <span
                        className={
                          tok.read_only
                            ? 'rounded-sm bg-gray-100 px-1.5 py-0.5 text-xs text-gray-700'
                            : 'rounded-sm bg-amber-50 px-1.5 py-0.5 text-xs text-amber-800'
                        }
                      >
                        {tok.read_only ? t('scopeReadOnly') : t('scopeFull')}
                      </span>
                    </div>
                    <div className="mt-0.5 flex flex-wrap gap-x-3 text-xs text-muted-foreground">
                      <span className="font-mono">{tok.prefix}…</span>
                      <span>{t('createdAt', { time: formatDateTime(tok.created_at, locale) })}</span>
                      <span>
                        {tok.last_used_at ? t('lastUsed', { time: formatDateTime(tok.last_used_at, locale) }) : t('neverUsed')}
                      </span>
                    </div>
                  </div>
                  <Button variant="outline" size="sm" aria-label={t('delete')} title={t('delete')} onClick={() => setDeleteTarget(tok)}>
                    <Trash2 />
                  </Button>
                </li>
              ))}
            </ul>
          )}
        </div>
      </CardContent>

      <Dialog open={created !== null} onOpenChange={(o) => !o && setCreated(null)}>
        <DialogContent className="max-w-2xl">
          <DialogHeader>
            <DialogTitle>{t('createdTitle')}</DialogTitle>
            <DialogDescription>{t('createdDesc')}</DialogDescription>
          </DialogHeader>
          {created && (
            <div className="min-w-0 space-y-4 text-sm">
              <div className="space-y-1">
                <Label>{t('token')}</Label>
                <CopyField value={created.token} label={t('copy')} />
              </div>
              <div className="space-y-1">
                <Label>{t('claudeCode')}</Label>
                <CopyField value={claudeCommand} label={t('copy')} multiline />
              </div>
              <div className="space-y-1">
                <Label>{t('jsonConfig')}</Label>
                <CopyField value={jsonConfig} label={t('copy')} multiline />
              </div>
            </div>
          )}
          <DialogFooter>
            <Button onClick={() => setCreated(null)}>{t('done')}</Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      <ConfirmDialog
        open={deleteTarget !== null}
        onOpenChange={(o) => !o && setDeleteTarget(null)}
        title={t('deleteTitle')}
        description={t('deleteDesc', { name: deleteTarget?.name ?? '' })}
        confirmText={t('delete')}
        destructive
        onConfirm={remove}
      />
    </Card>
  );
}
