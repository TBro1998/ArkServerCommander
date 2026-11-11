"use client";

import { useCallback, useEffect, useState } from 'react';
import { useLocale, useTranslations } from 'next-intl';
import axios from 'axios';
import { CalendarClock, Edit, Loader2, Play, Plus, Trash2 } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Switch } from '@/components/ui/switch';
import { Dialog, DialogContent, DialogFooter, DialogHeader, DialogTitle } from '@/components/ui/dialog';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { ConfirmDialog } from '@/components/ui/confirm-dialog';
import { useApiError } from '@/lib/api-error';
import { formatDateTime } from '@/lib/format';
import { toast } from '@/lib/toast';
import type { Server } from '@/stores/servers';

type Action = 'restart' | 'backup' | 'update' | 'command' | 'start' | 'stop';
const ACTIONS: Action[] = ['restart', 'backup', 'update', 'command', 'start', 'stop'];
const WARN_ACTIONS: Action[] = ['restart', 'update', 'stop'];

interface Schedule {
  id: number;
  action: Action;
  cron: string;
  enabled: boolean;
  warn_minutes: number;
  message: string;
  command: string;
  keep: number;
  last_run_at?: string | null;
  last_result?: string;
  next_run_at?: string | null;
}

type Mode = 'daily' | 'weekly' | 'interval' | 'custom';

interface FormState {
  action: Action;
  mode: Mode;
  time: string; // HH:MM
  weekday: string; // 0-6
  hours: number;
  cron: string;
  enabled: boolean;
  warn_minutes: number;
  message: string;
  command: string;
  keep: number;
}

const DEFAULT_FORM: FormState = {
  action: 'restart',
  mode: 'daily',
  time: '04:00',
  weekday: '1',
  hours: 6,
  cron: '0 4 * * *',
  enabled: true,
  warn_minutes: 5,
  message: '',
  command: '',
  keep: 7,
};

const pad = (n: number) => String(n).padStart(2, '0');

// 把表单里的频率转成 Cron 表达式
function buildCron(f: FormState): string {
  const [h, m] = f.time.split(':').map((v) => parseInt(v, 10) || 0);
  switch (f.mode) {
    case 'daily':
      return `${m} ${h} * * *`;
    case 'weekly':
      return `${m} ${h} * * ${f.weekday}`;
    case 'interval':
      return `0 */${Math.max(1, Math.min(23, f.hours))} * * *`;
    default:
      return f.cron.trim();
  }
}

// 识别常见的 Cron 写法，便于编辑时回填到简单模式
function parseCron(cron: string): Pick<FormState, 'mode' | 'time' | 'weekday' | 'hours' | 'cron'> {
  const base = { mode: 'custom' as Mode, time: '04:00', weekday: '1', hours: 6, cron };
  let m = cron.match(/^(\d{1,2}) (\d{1,2}) \* \* \*$/);
  if (m) return { ...base, mode: 'daily', time: `${pad(+m[2])}:${pad(+m[1])}` };
  m = cron.match(/^(\d{1,2}) (\d{1,2}) \* \* ([0-6])$/);
  if (m) return { ...base, mode: 'weekly', time: `${pad(+m[2])}:${pad(+m[1])}`, weekday: m[3] };
  m = cron.match(/^0 \*\/(\d{1,2}) \* \* \*$/);
  if (m) return { ...base, mode: 'interval', hours: +m[1] };
  return base;
}

export function SchedulesPanel({ server, active }: { server: Server; active: boolean }) {
  const t = useTranslations('schedules');
  const locale = useLocale();
  const apiError = useApiError();
  const [list, setList] = useState<Schedule[] | null>(null);
  const [editing, setEditing] = useState<Schedule | 'new' | null>(null);
  const [form, setForm] = useState<FormState>(DEFAULT_FORM);
  const [saving, setSaving] = useState(false);
  const [deleteTarget, setDeleteTarget] = useState<Schedule | null>(null);

  const load = useCallback(async () => {
    try {
      const res = await axios.get(`/api/servers/${server.id}/schedules`);
      setList(res.data.data ?? []);
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

  const describe = (s: Schedule): string => {
    const p = parseCron(s.cron);
    const weekdays = t.raw('weekdays') as string[];
    switch (p.mode) {
      case 'daily':
        return t('everyDayAt', { time: p.time });
      case 'weekly':
        return t('everyWeekAt', { day: weekdays[+p.weekday] ?? p.weekday, time: p.time });
      case 'interval':
        return t('everyHours', { hours: p.hours });
      default:
        return `Cron: ${s.cron}`;
    }
  };

  const openNew = () => {
    setForm(DEFAULT_FORM);
    setEditing('new');
  };

  const openEdit = (s: Schedule) => {
    setForm({
      ...DEFAULT_FORM,
      ...parseCron(s.cron),
      action: s.action,
      enabled: s.enabled,
      warn_minutes: s.warn_minutes,
      message: s.message,
      command: s.command,
      keep: s.keep,
    });
    setEditing(s);
  };

  const save = async () => {
    const payload = {
      action: form.action,
      cron: buildCron(form),
      enabled: form.enabled,
      warn_minutes: WARN_ACTIONS.includes(form.action) ? form.warn_minutes : 0,
      message: form.message,
      command: form.command,
      keep: form.action === 'backup' ? form.keep : 0,
    };
    setSaving(true);
    try {
      if (editing === 'new') await axios.post(`/api/servers/${server.id}/schedules`, payload);
      else if (editing) await axios.put(`/api/servers/${server.id}/schedules/${editing.id}`, payload);
      toast.success(t('saved'));
      setEditing(null);
      load();
    } catch (err) {
      toast.error(apiError(err, t('saveFailed')));
    } finally {
      setSaving(false);
    }
  };

  const toggle = async (s: Schedule, enabled: boolean) => {
    try {
      await axios.put(`/api/servers/${server.id}/schedules/${s.id}`, { ...s, enabled });
      load();
    } catch (err) {
      toast.error(apiError(err, t('saveFailed')));
    }
  };

  const runNow = async (s: Schedule) => {
    try {
      await axios.post(`/api/servers/${server.id}/schedules/${s.id}/run`);
      toast.success(t('runStarted'));
      setTimeout(load, 3000);
    } catch (err) {
      toast.error(apiError(err, t('runFailed')));
    }
  };

  const remove = async () => {
    if (!deleteTarget) return;
    try {
      await axios.delete(`/api/servers/${server.id}/schedules/${deleteTarget.id}`);
      toast.success(t('deleted'));
      load();
    } catch (err) {
      toast.error(apiError(err, t('deleteFailed')));
      throw err;
    }
  };

  const set = <K extends keyof FormState>(key: K, value: FormState[K]) => setForm((f) => ({ ...f, [key]: value }));
  const weekdays = t.raw('weekdays') as string[];

  return (
    <Card>
      <CardHeader className="flex flex-row items-center justify-between pb-3">
        <div>
          <CardTitle className="text-base">{t('title')}</CardTitle>
          <p className="mt-1 text-xs text-gray-500">{t('timezoneHint')}</p>
        </div>
        <Button onClick={openNew}>
          <Plus className="mr-2 h-4 w-4" />
          {t('add')}
        </Button>
      </CardHeader>
      <CardContent>
        {list === null ? (
          <div className="py-6 text-center">
            <Loader2 className="mx-auto h-6 w-6 animate-spin text-gray-400" />
          </div>
        ) : list.length === 0 ? (
          <div className="py-8 text-center text-sm text-gray-500">
            <CalendarClock className="mx-auto mb-2 h-8 w-8 text-gray-300" />
            {t('empty')}
          </div>
        ) : (
          <ul className="divide-y">
            {list.map((s) => (
              <li key={s.id} className="flex flex-col gap-2 py-3 sm:flex-row sm:items-center sm:justify-between">
                <div className="min-w-0 space-y-0.5">
                  <div className="flex flex-wrap items-center gap-2">
                    <span className="font-medium">{t(`actions.${s.action}`)}</span>
                    <span className="text-sm text-gray-600">{describe(s)}</span>
                    {s.action === 'command' && <code className="rounded bg-gray-100 px-1 text-xs">{s.command}</code>}
                  </div>
                  <div className="text-xs text-gray-500">
                    {s.enabled && s.next_run_at && <span>{t('nextRun', { time: formatDateTime(s.next_run_at, locale) })} · </span>}
                    {s.last_run_at ? (
                      <span className={s.last_result ? 'text-red-600' : ''}>
                        {t('lastRun', { time: formatDateTime(s.last_run_at, locale) })}
                        {s.last_result ? ` · ${s.last_result}` : ` · ${t('ok')}`}
                      </span>
                    ) : (
                      <span>{t('neverRun')}</span>
                    )}
                  </div>
                </div>
                <div className="flex flex-shrink-0 items-center gap-1">
                  <Switch checked={s.enabled} onCheckedChange={(v) => toggle(s, v)} title={t('enabled')} />
                  <Button variant="ghost" size="sm" className="h-8 w-8 p-0" title={t('runNow')} onClick={() => runNow(s)}>
                    <Play className="h-4 w-4" />
                  </Button>
                  <Button variant="ghost" size="sm" className="h-8 w-8 p-0" title={t('edit')} onClick={() => openEdit(s)}>
                    <Edit className="h-4 w-4" />
                  </Button>
                  <Button variant="ghost" size="sm" className="h-8 w-8 p-0 text-red-600" title={t('delete')} onClick={() => setDeleteTarget(s)}>
                    <Trash2 className="h-4 w-4" />
                  </Button>
                </div>
              </li>
            ))}
          </ul>
        )}
      </CardContent>

      <Dialog open={editing !== null} onOpenChange={(o) => !o && !saving && setEditing(null)}>
        <DialogContent className="max-w-lg">
          <DialogHeader>
            <DialogTitle>{editing === 'new' ? t('add') : t('edit')}</DialogTitle>
          </DialogHeader>
          <div className="space-y-4">
            <div className="space-y-1">
              <Label>{t('action')}</Label>
              <Select value={form.action} onValueChange={(v) => set('action', v as Action)}>
                <SelectTrigger className="w-full">
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  {ACTIONS.map((a) => (
                    <SelectItem key={a} value={a}>
                      {t(`actions.${a}`)}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
              <p className="text-xs text-gray-500">{t(`actionDesc.${form.action}`)}</p>
            </div>

            <div className="space-y-1">
              <Label>{t('frequency')}</Label>
              <Select value={form.mode} onValueChange={(v) => set('mode', v as Mode)}>
                <SelectTrigger className="w-full">
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  {(['daily', 'weekly', 'interval', 'custom'] as Mode[]).map((m) => (
                    <SelectItem key={m} value={m}>
                      {t(`modes.${m}`)}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>

            {(form.mode === 'daily' || form.mode === 'weekly') && (
              <div className="grid grid-cols-2 gap-3">
                {form.mode === 'weekly' && (
                  <div className="space-y-1">
                    <Label>{t('weekday')}</Label>
                    <Select value={form.weekday} onValueChange={(v) => set('weekday', v)}>
                      <SelectTrigger className="w-full">
                        <SelectValue />
                      </SelectTrigger>
                      <SelectContent>
                        {weekdays.map((d, i) => (
                          <SelectItem key={i} value={String(i)}>
                            {d}
                          </SelectItem>
                        ))}
                      </SelectContent>
                    </Select>
                  </div>
                )}
                <div className="space-y-1">
                  <Label>{t('time')}</Label>
                  <Input type="time" value={form.time} onChange={(e) => set('time', e.target.value)} />
                </div>
              </div>
            )}
            {form.mode === 'interval' && (
              <div className="space-y-1">
                <Label>{t('intervalHours')}</Label>
                <Input type="number" min={1} max={23} value={form.hours} onChange={(e) => set('hours', Number(e.target.value))} />
              </div>
            )}
            {form.mode === 'custom' && (
              <div className="space-y-1">
                <Label>Cron</Label>
                <Input value={form.cron} onChange={(e) => set('cron', e.target.value)} placeholder="0 4 * * *" className="font-mono" />
                <p className="text-xs text-gray-500">{t('cronHint')}</p>
              </div>
            )}

            {WARN_ACTIONS.includes(form.action) && (
              <>
                <div className="space-y-1">
                  <Label>{t('warnMinutes')}</Label>
                  <Input
                    type="number"
                    min={0}
                    max={60}
                    value={form.warn_minutes}
                    onChange={(e) => set('warn_minutes', Number(e.target.value))}
                  />
                  <p className="text-xs text-gray-500">{t('warnMinutesHint')}</p>
                </div>
                {form.warn_minutes > 0 && (
                  <div className="space-y-1">
                    <Label>{t('message')}</Label>
                    <Input value={form.message} onChange={(e) => set('message', e.target.value)} placeholder={t('messagePlaceholder', { minutes: '{minutes}' })} />
                  </div>
                )}
              </>
            )}
            {form.action === 'command' && (
              <div className="space-y-1">
                <Label>{t('command')}</Label>
                <Input value={form.command} onChange={(e) => set('command', e.target.value)} placeholder="SaveWorld" className="font-mono" />
              </div>
            )}
            {form.action === 'backup' && (
              <div className="space-y-1">
                <Label>{t('keep')}</Label>
                <Input type="number" min={0} max={1000} value={form.keep} onChange={(e) => set('keep', Number(e.target.value))} />
                <p className="text-xs text-gray-500">{t('keepHint')}</p>
              </div>
            )}

            <label className="flex items-center justify-between">
              <span className="text-sm font-medium">{t('enabled')}</span>
              <Switch checked={form.enabled} onCheckedChange={(v) => set('enabled', v)} />
            </label>
          </div>
          <DialogFooter>
            <Button variant="outline" onClick={() => setEditing(null)} disabled={saving}>
              {t('cancel')}
            </Button>
            <Button onClick={save} disabled={saving}>
              {saving && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
              {t('save')}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      <ConfirmDialog
        open={deleteTarget !== null}
        onOpenChange={(o) => !o && setDeleteTarget(null)}
        title={t('deleteTitle')}
        description={deleteTarget ? `${t(`actions.${deleteTarget.action}`)} · ${describe(deleteTarget)}` : ''}
        confirmText={t('delete')}
        destructive
        onConfirm={remove}
      />
    </Card>
  );
}
