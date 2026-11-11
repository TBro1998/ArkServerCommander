"use client";

import { useCallback, useEffect, useRef, useState } from 'react';
import { useTranslations } from 'next-intl';
import axios from 'axios';
import { Send, Loader2, RefreshCw, UserX, Ban, Megaphone, Save, Copy } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { ConfirmDialog } from '@/components/ui/confirm-dialog';
import { useApiError } from '@/lib/api-error';
import { copyText } from '@/lib/format';
import { toast } from '@/lib/toast';
import type { Server } from '@/stores/servers';

interface Player {
  index: number;
  name: string;
  steam_id: string;
}

interface HistoryEntry {
  command: string;
  output: string;
  ok: boolean;
  time: string;
}

// 常用命令：点击后填入输入框，带参数的命令需要补全
const QUICK_COMMANDS = ['SaveWorld', 'ListPlayers', 'GetChat', 'DestroyWildDinos', 'SetTimeOfDay 08:00', 'Broadcast '];

export function ConsolePanel({ server, active }: { server: Server; active: boolean }) {
  const t = useTranslations('console');
  const apiError = useApiError();
  const running = server.status === 'running';
  const ready = running && server.runtime?.phase === 'ready';

  const [command, setCommand] = useState('');
  const [sending, setSending] = useState(false);
  const [history, setHistory] = useState<HistoryEntry[]>([]);
  const [historyIndex, setHistoryIndex] = useState(-1);
  const [broadcastText, setBroadcastText] = useState('');
  const [players, setPlayers] = useState<Player[] | null>(null);
  const [loadingPlayers, setLoadingPlayers] = useState(false);
  const [target, setTarget] = useState<{ player: Player; action: 'kick' | 'ban' } | null>(null);
  const outputRef = useRef<HTMLDivElement>(null);

  const exec = useCallback(
    async (cmd: string): Promise<string> => {
      const res = await axios.post(`/api/servers/${server.id}/rcon/command`, { command: cmd });
      return res.data.data?.output ?? '';
    },
    [server.id],
  );

  const send = async (cmd: string) => {
    const trimmed = cmd.trim();
    if (!trimmed || sending) return;
    setSending(true);
    const time = new Date().toLocaleTimeString();
    try {
      const output = await exec(trimmed);
      setHistory((h) => [...h, { command: trimmed, output: output || t('noOutput'), ok: true, time }]);
      setCommand('');
    } catch (err) {
      setHistory((h) => [...h, { command: trimmed, output: apiError(err, t('commandFailed')), ok: false, time }]);
    } finally {
      setSending(false);
      setHistoryIndex(-1);
      requestAnimationFrame(() => {
        const el = outputRef.current;
        if (el) el.scrollTop = el.scrollHeight;
      });
    }
  };

  // 上下方向键浏览历史命令
  const onKeyDown = (e: React.KeyboardEvent<HTMLInputElement>) => {
    const commands = history.map((h) => h.command);
    if (e.key === 'ArrowUp' && commands.length > 0) {
      e.preventDefault();
      const next = historyIndex < 0 ? commands.length - 1 : Math.max(0, historyIndex - 1);
      setHistoryIndex(next);
      setCommand(commands[next]);
    } else if (e.key === 'ArrowDown' && historyIndex >= 0) {
      e.preventDefault();
      const next = historyIndex + 1;
      if (next >= commands.length) {
        setHistoryIndex(-1);
        setCommand('');
      } else {
        setHistoryIndex(next);
        setCommand(commands[next]);
      }
    }
  };

  const loadPlayers = useCallback(async () => {
    setLoadingPlayers(true);
    try {
      const res = await axios.get(`/api/servers/${server.id}/players`);
      setPlayers(res.data.data ?? []);
    } catch (err) {
      setPlayers(null);
      toast.error(apiError(err, t('playersFailed')));
    } finally {
      setLoadingPlayers(false);
    }
  }, [server.id, apiError, t]);

  // 切到控制台时，如果服务器可以连接，自动加载一次玩家列表
  useEffect(() => {
    if (!active || !ready) return;
    const timer = setTimeout(() => {
      loadPlayers();
    }, 0);
    return () => clearTimeout(timer);
  }, [active, ready, loadPlayers]);

  const doBroadcast = async () => {
    const text = broadcastText.trim();
    if (!text) return;
    try {
      await exec(`Broadcast ${text}`);
      toast.success(t('broadcastSent'));
      setBroadcastText('');
    } catch (err) {
      toast.error(apiError(err, t('commandFailed')));
    }
  };

  const doSave = async () => {
    try {
      await exec('SaveWorld');
      toast.success(t('saved'));
    } catch (err) {
      toast.error(apiError(err, t('commandFailed')));
    }
  };

  const doPlayerAction = async () => {
    if (!target) return;
    const cmd = target.action === 'kick' ? `KickPlayer ${target.player.steam_id}` : `BanPlayer ${target.player.steam_id}`;
    try {
      await exec(cmd);
      toast.success(target.action === 'kick' ? t('kicked', { name: target.player.name }) : t('banned', { name: target.player.name }));
      loadPlayers();
    } catch (err) {
      toast.error(apiError(err, t('commandFailed')));
      throw err;
    }
  };

  if (!running) {
    return <p className="py-8 text-center text-sm text-gray-500">{t('notRunning')}</p>;
  }

  return (
    <div className="grid gap-4 lg:grid-cols-3">
      <div className="space-y-4 lg:col-span-2">
        {!ready && <p className="rounded-lg bg-amber-50 p-3 text-sm text-amber-800">{t('notReady')}</p>}

        <Card>
          <CardHeader className="pb-3">
            <CardTitle className="text-base">{t('quickActions')}</CardTitle>
          </CardHeader>
          <CardContent className="space-y-3">
            <div className="flex flex-col gap-2 sm:flex-row">
              <Input
                value={broadcastText}
                onChange={(e) => setBroadcastText(e.target.value)}
                onKeyDown={(e) => e.key === 'Enter' && doBroadcast()}
                placeholder={t('broadcastPlaceholder')}
              />
              <Button onClick={doBroadcast} disabled={!broadcastText.trim()}>
                <Megaphone className="mr-2 h-4 w-4" />
                {t('broadcast')}
              </Button>
            </div>
            <Button variant="outline" onClick={doSave}>
              <Save className="mr-2 h-4 w-4" />
              {t('saveWorld')}
            </Button>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="pb-3">
            <CardTitle className="text-base">{t('rconConsole')}</CardTitle>
          </CardHeader>
          <CardContent className="space-y-3">
            <div ref={outputRef} className="h-72 overflow-auto rounded-lg bg-gray-950 p-3 font-mono text-xs leading-5">
              {history.length === 0 ? (
                <div className="text-gray-500">{t('consoleHint')}</div>
              ) : (
                history.map((h, i) => (
                  <div key={i} className="mb-2">
                    <div className="text-sky-300">
                      <span className="text-gray-500">[{h.time}]</span> &gt; {h.command}
                    </div>
                    <div className={`whitespace-pre-wrap break-all ${h.ok ? 'text-gray-200' : 'text-red-400'}`}>{h.output}</div>
                  </div>
                ))
              )}
            </div>
            <form
              className="flex gap-2"
              onSubmit={(e) => {
                e.preventDefault();
                send(command);
              }}
            >
              <Input
                value={command}
                onChange={(e) => setCommand(e.target.value)}
                onKeyDown={onKeyDown}
                placeholder={t('commandPlaceholder')}
                className="font-mono"
                autoComplete="off"
              />
              <Button type="submit" disabled={sending || !command.trim()}>
                {sending ? <Loader2 className="h-4 w-4 animate-spin" /> : <Send className="h-4 w-4" />}
              </Button>
            </form>
            <div className="flex flex-wrap gap-1">
              {QUICK_COMMANDS.map((c) => (
                <button
                  key={c}
                  type="button"
                  className="rounded border bg-gray-50 px-2 py-0.5 font-mono text-xs hover:bg-gray-100"
                  onClick={() => setCommand(c)}
                >
                  {c.trim()}
                </button>
              ))}
            </div>
          </CardContent>
        </Card>
      </div>

      <Card className="h-fit">
        <CardHeader className="flex flex-row items-center justify-between pb-3">
          <CardTitle className="text-base">
            {t('players')}
            {players && <span className="ml-1 text-sm font-normal text-gray-500">({players.length})</span>}
          </CardTitle>
          <Button variant="ghost" size="sm" onClick={loadPlayers} disabled={loadingPlayers || !ready}>
            <RefreshCw className={`h-4 w-4 ${loadingPlayers ? 'animate-spin' : ''}`} />
          </Button>
        </CardHeader>
        <CardContent>
          {!ready ? (
            <p className="text-sm text-gray-500">{t('playersNotReady')}</p>
          ) : players === null ? (
            <p className="text-sm text-gray-500">{loadingPlayers ? t('loading') : t('playersUnknown')}</p>
          ) : players.length === 0 ? (
            <p className="text-sm text-gray-500">{t('noPlayers')}</p>
          ) : (
            <ul className="divide-y">
              {players.map((p) => (
                <li key={p.steam_id} className="flex items-center justify-between gap-2 py-2">
                  <div className="min-w-0">
                    <div className="truncate text-sm font-medium">{p.name || '-'}</div>
                    <button
                      type="button"
                      className="flex items-center gap-1 font-mono text-xs text-gray-500 hover:text-gray-800"
                      onClick={async () => {
                        if (await copyText(p.steam_id)) toast.success(t('steamIdCopied'));
                      }}
                    >
                      {p.steam_id}
                      <Copy className="h-3 w-3" />
                    </button>
                  </div>
                  <div className="flex flex-shrink-0 gap-1">
                    <Button variant="ghost" size="sm" className="h-7 w-7 p-0" title={t('kick')} onClick={() => setTarget({ player: p, action: 'kick' })}>
                      <UserX className="h-4 w-4" />
                    </Button>
                    <Button
                      variant="ghost"
                      size="sm"
                      className="h-7 w-7 p-0 text-red-600"
                      title={t('ban')}
                      onClick={() => setTarget({ player: p, action: 'ban' })}
                    >
                      <Ban className="h-4 w-4" />
                    </Button>
                  </div>
                </li>
              ))}
            </ul>
          )}
        </CardContent>
      </Card>

      <ConfirmDialog
        open={target !== null}
        onOpenChange={(o) => !o && setTarget(null)}
        title={target?.action === 'ban' ? t('banTitle', { name: target?.player.name ?? '' }) : t('kickTitle', { name: target?.player.name ?? '' })}
        description={target?.action === 'ban' ? t('banDesc') : t('kickDesc')}
        confirmText={target?.action === 'ban' ? t('ban') : t('kick')}
        destructive={target?.action === 'ban'}
        onConfirm={doPlayerAction}
      />
    </div>
  );
}
