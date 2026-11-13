"use client";

import { useCallback, useEffect, useRef, useState } from 'react';
import { useTranslations } from 'next-intl';
import axios from 'axios';
import { Pause, Play, Download, Trash2, RefreshCw, ArrowDown } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { cn } from '@/lib/utils';
import { toast } from '@/lib/toast';

const MAX_LINES = 5000;
const INITIAL_TAIL = 500;
const GAME_LOG_TAIL = 2000;
const GAME_LOG_INTERVAL = 5000;

type StreamState = 'connecting' | 'live' | 'ended' | 'error';
// container：容器输出（入口脚本、steamcmd），SSE 实时推送；game：存档目录里的 ShooterGame.log，定时读取
type LogSource = 'container' | 'game';

// steamcmd 下载进度行，下载服务端时会刷屏
const PROGRESS_LINE = /Update state|progress:/i;

// 按内容给日志行着色：错误、警告、脚本的阶段提示
function lineClass(line: string): string {
  if (/\[!\]|error|failed|fatal/i.test(line)) return 'text-red-700';
  if (/warn/i.test(line)) return 'text-amber-700';
  if (/^\s*\[\*\]|^\s*\[i\]/.test(line)) return 'text-primary';
  if (PROGRESS_LINE.test(line)) return 'text-muted-foreground';
  return 'text-terminal-foreground';
}

export function LogViewer({ serverId, active }: { serverId: number; active: boolean }) {
  const t = useTranslations('logs');
  const [lines, setLines] = useState<string[]>([]);
  const [state, setState] = useState<StreamState>('connecting');
  const [paused, setPaused] = useState(false);
  const [filter, setFilter] = useState('');
  const [autoScroll, setAutoScroll] = useState(true);
  const [session, setSession] = useState(0);
  const [source, setSource] = useState<LogSource>('container');
  const [hideProgress, setHideProgress] = useState(false);
  const [gameLogMissing, setGameLogMissing] = useState(false);
  const containerRef = useRef<HTMLDivElement>(null);
  // 暂停期间收到的行先放在这里，恢复时一次性追加
  const bufferRef = useRef<string[]>([]);
  const pausedRef = useRef(false);

  useEffect(() => {
    pausedRef.current = paused;
  }, [paused]);

  const flush = useCallback(() => {
    if (bufferRef.current.length === 0) return;
    const incoming = bufferRef.current;
    bufferRef.current = [];
    setLines((prev) => {
      const next = prev.concat(incoming);
      return next.length > MAX_LINES ? next.slice(next.length - MAX_LINES) : next;
    });
  }, []);

  // 游戏日志：定时读取文件末尾，整体替换；暂停时不刷新
  useEffect(() => {
    if (!active || source !== 'game') return;
    let cancelled = false;
    const load = async () => {
      if (pausedRef.current) return;
      try {
        const res = await axios.get(`/api/servers/${serverId}/logs/game`, { params: { tail: GAME_LOG_TAIL } });
        if (cancelled) return;
        setLines(res.data.data?.lines ?? []);
        setGameLogMissing(!res.data.data?.exists);
        setState('live');
      } catch {
        if (!cancelled) setState('error');
      }
    };
    const first = setTimeout(load, 0);
    const timer = setInterval(load, GAME_LOG_INTERVAL);
    return () => {
      cancelled = true;
      clearTimeout(first);
      clearInterval(timer);
    };
  }, [serverId, active, source, session]);

  // 打开日志页时建立 SSE 连接，离开或重连时关闭
  useEffect(() => {
    if (!active || source !== 'container') return;
    const es = new EventSource(`/api/servers/${serverId}/logs/stream?tail=${INITIAL_TAIL}`);
    let gotData = false;
    // 合并 200ms 内收到的行再渲染，日志刷屏时不会卡
    const timer = setInterval(() => {
      if (!pausedRef.current) flush();
    }, 200);

    es.onopen = () => setState('live');
    es.onmessage = (e) => {
      gotData = true;
      bufferRef.current.push(e.data);
      if (bufferRef.current.length > MAX_LINES) bufferRef.current.splice(0, bufferRef.current.length - MAX_LINES);
    };
    es.addEventListener('end', () => {
      es.close();
      flush();
      setState('ended');
    });
    es.onerror = () => {
      // EventSource 默认会无限重连；这里交给用户手动重连
      es.close();
      flush();
      setState(gotData ? 'ended' : 'error');
    };

    return () => {
      clearInterval(timer);
      es.close();
    };
  }, [serverId, active, source, session, flush]);

  // 新日志到达时滚动到底部（用户向上翻看时不打扰）
  useEffect(() => {
    if (!autoScroll) return;
    const el = containerRef.current;
    if (el) el.scrollTop = el.scrollHeight;
  }, [lines, autoScroll]);

  const onScroll = () => {
    const el = containerRef.current;
    if (!el) return;
    const atBottom = el.scrollHeight - el.scrollTop - el.clientHeight < 40;
    setAutoScroll(atBottom);
  };

  const reconnect = () => {
    bufferRef.current = [];
    setLines([]);
    setState('connecting');
    setSession((s) => s + 1);
  };

  const changeSource = (next: LogSource) => {
    if (next === source) return;
    bufferRef.current = [];
    setLines([]);
    setPaused(false);
    setGameLogMissing(false);
    setState('connecting');
    setSource(next);
  };

  const togglePause = () => {
    if (paused) flush();
    setPaused(!paused);
  };

  const download = async () => {
    try {
      const data: string[] =
        source === 'game'
          ? ((await axios.get(`/api/servers/${serverId}/logs/game`, { params: { tail: MAX_LINES } })).data.data?.lines ?? [])
          : (await axios.get(`/api/servers/${serverId}/logs`, { params: { tail: MAX_LINES } })).data.data;
      const blob = new Blob([data.join('\n')], { type: 'text/plain;charset=utf-8' });
      const url = URL.createObjectURL(blob);
      const a = document.createElement('a');
      a.href = url;
      const prefix = source === 'game' ? 'ShooterGame' : 'server';
      a.download = `${prefix}-${serverId}-${new Date().toISOString().replace(/[:.]/g, '-')}.log`;
      a.click();
      URL.revokeObjectURL(url);
    } catch {
      toast.error(t('downloadFailed'));
    }
  };

  const keyword = filter.trim().toLowerCase();
  const hidingProgress = hideProgress && source === 'container';
  const visible =
    keyword || hidingProgress
      ? lines.filter((l) => (!keyword || l.toLowerCase().includes(keyword)) && (!hidingProgress || !PROGRESS_LINE.test(l)))
      : lines;

  const stateLabel = {
    connecting: t('connecting'),
    live: paused ? t('paused') : t('live'),
    ended: t('ended'),
    error: t('error'),
  }[state];

  return (
    <div className="space-y-2">
      <div className="flex flex-wrap items-center gap-2">
        <div className="inline-flex rounded-sm border p-0.5" role="tablist" aria-label={t('source')}>
          {(['container', 'game'] as const).map((s) => (
            <button
              key={s}
              type="button"
              role="tab"
              aria-selected={source === s}
              onClick={() => changeSource(s)}
              className={cn(
                'rounded-[3px] px-2.5 py-1 text-xs transition-colors',
                source === s ? 'bg-primary/15 text-primary' : 'text-muted-foreground hover:text-foreground',
              )}
            >
              {s === 'container' ? t('sourceContainer') : t('sourceGame')}
            </button>
          ))}
        </div>
        <span
          className={cn(
            'inline-flex items-center gap-1.5 rounded-md px-2 py-1 text-xs',
            state === 'live' && !paused ? 'bg-green-100 text-green-800' : 'bg-gray-100 text-gray-700',
          )}
        >
          <span className={cn('h-2 w-2 rounded-full', state === 'live' && !paused ? 'animate-pulse bg-green-500' : 'bg-gray-400')} />
          {stateLabel}
        </span>
        <Input
          value={filter}
          onChange={(e) => setFilter(e.target.value)}
          placeholder={t('filter')}
          className="h-8 w-full sm:w-56"
        />
        {source === 'container' && (
          <label className="flex cursor-pointer items-center gap-1.5 text-xs text-muted-foreground">
            <input
              type="checkbox"
              className="size-3.5 accent-primary"
              checked={hideProgress}
              onChange={(e) => setHideProgress(e.target.checked)}
            />
            {t('hideProgress')}
          </label>
        )}
        <div className="ml-auto flex items-center gap-1">
          {state === 'live' && (
            <Button variant="outline" size="sm" onClick={togglePause}>
              {paused ? <Play className="mr-1 h-3.5 w-3.5" /> : <Pause className="mr-1 h-3.5 w-3.5" />}
              {paused ? t('resume') : t('pause')}
            </Button>
          )}
          {(state === 'ended' || state === 'error') && (
            <Button variant="outline" size="sm" onClick={reconnect}>
              <RefreshCw className="mr-1 h-3.5 w-3.5" />
              {t('reconnect')}
            </Button>
          )}
          <Button variant="outline" size="sm" onClick={() => setLines([])} title={t('clear')}>
            <Trash2 className="h-3.5 w-3.5" />
          </Button>
          <Button variant="outline" size="sm" onClick={download} title={t('download')}>
            <Download className="h-3.5 w-3.5" />
          </Button>
        </div>
      </div>

      <div className="relative">
        <div
          ref={containerRef}
          onScroll={onScroll}
          className="h-[55vh] min-h-[320px] overflow-auto rounded-sm border bg-terminal p-3 font-mono text-xs leading-5"
        >
          {visible.length === 0 ? (
            <div className="text-gray-500">
              {state === 'connecting' ? t('connecting') : source === 'game' ? (gameLogMissing ? t('gameEmpty') : t('empty')) : t('empty')}
            </div>
          ) : (
            visible.map((line, i) => (
              <div key={i} className={cn('whitespace-pre-wrap break-all', lineClass(line))}>
                {line}
              </div>
            ))
          )}
        </div>
        {!autoScroll && (
          <Button
            size="sm"
            className="absolute bottom-3 right-5 shadow"
            onClick={() => {
              setAutoScroll(true);
              const el = containerRef.current;
              if (el) el.scrollTop = el.scrollHeight;
            }}
          >
            <ArrowDown className="mr-1 h-3.5 w-3.5" />
            {t('scrollToBottom')}
          </Button>
        )}
      </div>
      <p className="text-xs text-gray-500">{source === 'game' ? t('gameHint') : t('hint')}</p>
    </div>
  );
}
