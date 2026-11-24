"use client";

import { useCallback, useEffect, useRef, useState } from 'react';
import { useLocale, useTranslations } from 'next-intl';
import axios from 'axios';
import { ArrowDown, Download, FileText, RefreshCw } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { useApiError } from '@/lib/api-error';
import { formatBytes } from '@/lib/format';
import { toast } from '@/lib/toast';
import { cn } from '@/lib/utils';
import type { Server } from '@/stores/servers';
import { lineClass } from './LogViewer';

const LOG_TAIL = 2000;
const DOWNLOAD_TAIL = 5000;
const REFRESH_INTERVAL = 5000;

interface PluginLogFile {
  name: string;
  size: number;
  modified_at: string;
}

interface PluginLogFiles {
  files: PluginLogFile[];
  total: number;
}

function formatTime(value: string, locale: string): string {
  const d = new Date(value);
  if (Number.isNaN(d.getTime())) return '';
  return d.toLocaleString(locale === 'zh' ? 'zh-CN' : 'en-US', { hour12: false });
}

// 插件页里的日志：ArkApi 和插件写在 ShooterGame/Binaries/Win64/logs 下的 .log 文件。
// 默认显示最新的文件，运行中每 5 秒刷新；停止后仍能看上次运行的日志
export function PluginLogsCard({ server, active }: { server: Server; active: boolean }) {
  const t = useTranslations('pluginLogs');
  const locale = useLocale();
  const apiError = useApiError();
  const [list, setList] = useState<PluginLogFiles | null>(null);
  // 空字符串表示跟随最新的文件（重启后会换成新的日志文件）
  const [picked, setPicked] = useState('');
  const [lines, setLines] = useState<string[]>([]);
  const [loadingLines, setLoadingLines] = useState(false);
  const [filter, setFilter] = useState('');
  const [autoScroll, setAutoScroll] = useState(true);
  const [session, setSession] = useState(0);
  const containerRef = useRef<HTMLDivElement>(null);

  const running = server.status === 'running';
  const files = list?.files ?? [];
  const current = picked && files.some((f) => f.name === picked) ? picked : (files[0]?.name ?? '');

  const loadList = useCallback(async () => {
    try {
      const res = await axios.get(`/api/servers/${server.id}/arkapi/logs`);
      setList(res.data.data);
    } catch (err) {
      toast.error(apiError(err, t('listFailed')));
    }
  }, [server.id, apiError, t]);

  // 文件列表：打开标签、服务器状态变化、手动刷新时读取；运行中跟着日志内容一起刷新
  useEffect(() => {
    if (!active) return;
    const first = setTimeout(loadList, 0);
    const timer = running ? setInterval(loadList, REFRESH_INTERVAL) : undefined;
    return () => {
      clearTimeout(first);
      if (timer) clearInterval(timer);
    };
  }, [active, running, loadList, session]);

  // 日志内容：读取所选文件的末尾，整体替换
  useEffect(() => {
    if (!active || !current) return;
    let cancelled = false;
    const load = async (first: boolean) => {
      if (first) setLoadingLines(true);
      try {
        const res = await axios.get(`/api/servers/${server.id}/arkapi/logs/${encodeURIComponent(current)}`, {
          params: { tail: LOG_TAIL },
        });
        if (!cancelled) setLines(res.data.data?.lines ?? []);
      } catch (err) {
        if (!cancelled && first) toast.error(apiError(err, t('loadFailed')));
      } finally {
        if (!cancelled && first) setLoadingLines(false);
      }
    };
    const first = setTimeout(() => load(true), 0);
    const timer = running ? setInterval(() => load(false), REFRESH_INTERVAL) : undefined;
    return () => {
      cancelled = true;
      clearTimeout(first);
      if (timer) clearInterval(timer);
    };
  }, [active, current, running, server.id, apiError, t, session]);

  // 新内容到达时滚动到底部（用户向上翻看时不打扰）
  useEffect(() => {
    if (!autoScroll) return;
    const el = containerRef.current;
    if (el) el.scrollTop = el.scrollHeight;
  }, [lines, autoScroll]);

  const onScroll = () => {
    const el = containerRef.current;
    if (!el) return;
    setAutoScroll(el.scrollHeight - el.scrollTop - el.clientHeight < 40);
  };

  const choose = (name: string) => {
    // 选中最新的文件等于跟随最新
    setPicked(name === files[0]?.name ? '' : name);
    setLines([]);
    setAutoScroll(true);
  };

  const download = async () => {
    if (!current) return;
    try {
      const res = await axios.get(`/api/servers/${server.id}/arkapi/logs/${encodeURIComponent(current)}`, {
        params: { tail: DOWNLOAD_TAIL },
      });
      const data: string[] = res.data.data?.lines ?? [];
      const blob = new Blob([data.join('\n')], { type: 'text/plain;charset=utf-8' });
      const url = URL.createObjectURL(blob);
      const a = document.createElement('a');
      a.href = url;
      a.download = current;
      a.click();
      URL.revokeObjectURL(url);
    } catch (err) {
      toast.error(apiError(err, t('downloadFailed')));
    }
  };

  const keyword = filter.trim().toLowerCase();
  const visible = keyword ? lines.filter((l) => l.toLowerCase().includes(keyword)) : lines;

  return (
    <Card>
      <CardHeader className="flex flex-row flex-wrap items-center justify-between gap-2 space-y-0 pb-2">
        <CardTitle className="font-display flex items-center gap-2 text-base">
          <FileText className="h-4 w-4 text-primary" />
          {t('title')}
          {list && list.total > 0 && <span className="num text-sm font-normal text-muted-foreground">{list.total}</span>}
        </CardTitle>
        <Button
          variant="ghost"
          size="icon"
          className="h-8 w-8"
          onClick={() => setSession((s) => s + 1)}
          title={t('refresh')}
          aria-label={t('refresh')}
        >
          <RefreshCw className={loadingLines ? 'animate-spin' : ''} />
        </Button>
      </CardHeader>
      <CardContent className="space-y-2 text-sm">
        {!list ? null : files.length === 0 ? (
          <p className="text-muted-foreground">{t('empty')}</p>
        ) : (
          <>
            <div className="flex flex-wrap items-center gap-2">
              <Select value={current} onValueChange={choose}>
                <SelectTrigger className="w-full font-mono text-xs sm:w-96">
                  <SelectValue placeholder={t('selectFile')} />
                </SelectTrigger>
                <SelectContent>
                  {files.map((f, i) => (
                    <SelectItem key={f.name} value={f.name} className="font-mono text-xs">
                      {f.name}
                      <span className="ml-2 font-sans text-muted-foreground">
                        {[i === 0 ? t('latest') : '', formatBytes(f.size), formatTime(f.modified_at, locale)].filter(Boolean).join(' · ')}
                      </span>
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
              <Input value={filter} onChange={(e) => setFilter(e.target.value)} placeholder={t('filter')} className="h-8 w-full sm:w-56" />
              <div className="ml-auto flex items-center gap-1">
                {running && (
                  <span className="inline-flex items-center gap-1.5 rounded-md bg-green-100 px-2 py-1 text-xs text-green-800">
                    <span className="h-2 w-2 animate-pulse rounded-full bg-green-500" />
                    {t('autoRefresh')}
                  </span>
                )}
                <Button variant="outline" size="sm" onClick={download} title={t('download')} aria-label={t('download')}>
                  <Download className="size-3.5" />
                </Button>
              </div>
            </div>

            <div className="relative">
              <div
                ref={containerRef}
                onScroll={onScroll}
                className="h-[45vh] min-h-[260px] overflow-auto rounded-sm border bg-terminal p-3 font-mono text-xs leading-5"
              >
                {visible.length === 0 ? (
                  <div className="text-gray-500">{loadingLines ? t('loading') : keyword ? t('filterEmpty') : t('fileEmpty')}</div>
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
                  <ArrowDown className="mr-1 size-3.5" />
                  {t('scrollToBottom')}
                </Button>
              )}
            </div>
            {list.total > files.length && (
              <p className="text-xs text-muted-foreground">{t('truncated', { shown: files.length, total: list.total })}</p>
            )}
          </>
        )}
        <p className="text-xs text-muted-foreground">{t('hint', { tail: LOG_TAIL })}</p>
      </CardContent>
    </Card>
  );
}
