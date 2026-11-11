"use client";

import { useEffect, useRef, useState } from 'react';
import { useTranslations } from 'next-intl';
import { Loader2, Eye, EyeOff, Dices, Save, Plus, RotateCw } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Switch } from '@/components/ui/switch';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { Card, CardContent } from '@/components/ui/card';
import { GameUserSettingsEditor } from '@/components/servers/GameUserSettingsEditor';
import { GameIniEditor } from '@/components/servers/GameIniEditor';
import { ServerArgsEditor } from '@/components/servers/ServerArgsEditor';
import { MapSelector } from '@/components/servers/MapSelector';
import { ModList, isValidModIds, normalizeModIds } from '@/components/servers/ModList';
import { randomPassword } from '@/lib/format';
import type { Server } from '@/stores/servers';

const CLUSTER_ID = /^[A-Za-z0-9][A-Za-z0-9_.-]{0,63}$/;

export interface ServerFormProps {
  mode: 'create' | 'edit';
  initial: Partial<Server>;
  // 编辑时服务器在运行，显示"保存并重启"
  running?: boolean;
  saving: boolean;
  onSubmit: (data: Partial<Server>, options: { restart: boolean }) => Promise<boolean>;
  onCancel: () => void;
  onDirtyChange?: (dirty: boolean) => void;
}

function Field({ label, htmlFor, hint, error, children }: { label: string; htmlFor?: string; hint?: string; error?: string; children: React.ReactNode }) {
  return (
    <div className="space-y-1">
      <Label htmlFor={htmlFor}>{label}</Label>
      {children}
      {error ? <p className="text-xs text-red-600">{error}</p> : hint ? <p className="text-xs text-gray-500">{hint}</p> : null}
    </div>
  );
}

export function ServerForm({ mode, initial, running, saving, onSubmit, onCancel, onDirtyChange }: ServerFormProps) {
  const t = useTranslations('serverForm');
  const tEdit = useTranslations('servers.edit');
  const tCommon = useTranslations('common');
  const [form, setForm] = useState<Partial<Server>>(initial);
  const [dirty, setDirty] = useState(false);
  const [showPassword, setShowPassword] = useState(false);
  const [errors, setErrors] = useState<Record<string, string>>({});
  // 父组件异步加载到建议端口后会更新 initial；用户还没改过时跟着更新
  const [prevInitial, setPrevInitial] = useState(initial);
  if (initial !== prevInitial) {
    setPrevInitial(initial);
    if (!dirty) setForm(initial);
  }

  // 有未保存的修改时，关闭或刷新页面前提醒
  useEffect(() => {
    if (!dirty) return;
    const handler = (e: BeforeUnloadEvent) => {
      e.preventDefault();
    };
    window.addEventListener('beforeunload', handler);
    return () => window.removeEventListener('beforeunload', handler);
  }, [dirty]);

  const update = (patch: Partial<Server>) => {
    setForm((prev) => ({ ...prev, ...patch }));
    if (!dirty) {
      setDirty(true);
      onDirtyChange?.(true);
    }
  };

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value, type } = e.target;
    update({ [name]: type === 'number' ? (value === '' ? 0 : Number(value)) : value } as Partial<Server>);
    if (errors[name]) setErrors((prev) => ({ ...prev, [name]: '' }));
  };

  const validate = (): boolean => {
    const next: Record<string, string> = {};
    if (!form.identifier?.trim()) next.identifier = t('identifierRequired');
    if (mode === 'create' && !form.admin_password) next.admin_password = t('passwordRequired');
    if (form.admin_password && /[\r\n]/.test(form.admin_password)) next.admin_password = t('passwordNoNewline');
    for (const key of ['port', 'query_port', 'rcon_port'] as const) {
      const v = Number(form[key]);
      if (!Number.isInteger(v) || v < 1 || v > 65535) next[key] = t('portRange');
    }
    const ports = [Number(form.port), Number(form.port) + 1, Number(form.query_port), Number(form.rcon_port)];
    if (!next.port && !next.query_port && !next.rcon_port && new Set(ports).size !== ports.length) next.port = t('portDuplicate');
    const mp = Number(form.max_players);
    if (!Number.isInteger(mp) || mp < 1 || mp > 200) next.max_players = t('maxPlayersRange');
    if (form.cluster_id && !CLUSTER_ID.test(form.cluster_id.trim())) next.cluster_id = t('clusterInvalid');
    if (form.game_mod_ids && !isValidModIds(form.game_mod_ids)) next.game_mod_ids = t('modsInvalid');
    setErrors(next);
    return Object.keys(next).length === 0;
  };

  const submit = async (restart: boolean) => {
    if (!validate()) return;
    const data: Partial<Server> = {
      ...form,
      identifier: form.identifier?.trim(),
      cluster_id: form.cluster_id?.trim() ?? '',
      game_mod_ids: normalizeModIds(form.game_mod_ids ?? ''),
    };
    const ok = await onSubmit(data, { restart });
    if (ok) {
      setDirty(false);
      onDirtyChange?.(false);
    }
  };

  // Ctrl+S / Cmd+S 保存（用 ref 拿到最新的 submit，监听只注册一次）
  const submitRef = useRef(submit);
  useEffect(() => {
    submitRef.current = submit;
  });
  useEffect(() => {
    const onKey = (e: KeyboardEvent) => {
      if ((e.ctrlKey || e.metaKey) && e.key.toLowerCase() === 's') {
        e.preventDefault();
        submitRef.current(false);
      }
    };
    window.addEventListener('keydown', onKey);
    return () => window.removeEventListener('keydown', onKey);
  }, []);

  const cancel = () => {
    if (dirty && !window.confirm(t('discardConfirm'))) return;
    onCancel();
  };

  const hasErrors = Object.values(errors).some(Boolean);

  return (
    <Card className="pt-3">
      <CardContent>
        <Tabs defaultValue="basic">
          <div className="-mx-2 overflow-x-auto px-2">
            <TabsList variant="line">
              <TabsTrigger value="basic">
                {tEdit('basicParams')}
                {hasErrors && <span className="ml-1 h-2 w-2 rounded-full bg-red-500" />}
              </TabsTrigger>
              <TabsTrigger value="game_user_settings">{tEdit('gameUserSettings')}</TabsTrigger>
              <TabsTrigger value="game_ini">{tEdit('gameIni')}</TabsTrigger>
              <TabsTrigger value="server_args">{tEdit('serverArgs')}</TabsTrigger>
            </TabsList>
          </div>

          <TabsContent value="basic">
            <form
              onSubmit={(e) => {
                e.preventDefault();
                submit(false);
              }}
              className="space-y-6 py-4"
            >
              <section className="space-y-4">
                <h3 className="font-display text-base font-semibold">{t('sectionBasic')}</h3>
                <div className="grid grid-cols-1 gap-4 md:grid-cols-2">
                  <Field label={tEdit('serverName')} htmlFor="session_name" hint={t('sessionNameHint')}>
                    <Input id="session_name" name="session_name" value={form.session_name || ''} onChange={handleChange} />
                  </Field>
                  <Field label={`${tEdit('serverIdentifier')} *`} htmlFor="identifier" hint={t('identifierHint')} error={errors.identifier}>
                    <Input
                      id="identifier"
                      name="identifier"
                      value={form.identifier || ''}
                      onChange={handleChange}
                      placeholder={tEdit('serverIdentifierPlaceholder')}
                      aria-invalid={!!errors.identifier}
                    />
                  </Field>
                  <div className="space-y-1">
                    <MapSelector value={form.map || ''} onChange={(value) => update({ map: value })} label={tEdit('map')} />
                  </div>
                  <Field label={tEdit('maxPlayers')} htmlFor="max_players" hint={tEdit('maxPlayersDesc')} error={errors.max_players}>
                    <Input id="max_players" name="max_players" type="number" min={1} max={200} value={form.max_players || ''} onChange={handleChange} />
                  </Field>
                </div>
              </section>

              <section className="space-y-4">
                <h3 className="font-display text-base font-semibold">{t('sectionNetwork')}</h3>
                <div className="grid grid-cols-1 gap-4 sm:grid-cols-3">
                  <Field label={tEdit('gamePort')} htmlFor="port" hint={t('gamePortHint')} error={errors.port}>
                    <Input id="port" name="port" type="number" value={form.port || ''} onChange={handleChange} />
                  </Field>
                  <Field label={tEdit('queryPort')} htmlFor="query_port" hint={t('queryPortHint')} error={errors.query_port}>
                    <Input id="query_port" name="query_port" type="number" value={form.query_port || ''} onChange={handleChange} />
                  </Field>
                  <Field label={tEdit('rconPort')} htmlFor="rcon_port" hint={t('rconPortHint')} error={errors.rcon_port}>
                    <Input id="rcon_port" name="rcon_port" type="number" value={form.rcon_port || ''} onChange={handleChange} />
                  </Field>
                </div>
                <p className="text-xs text-gray-500">{t('portsFirewallHint')}</p>

                <Field
                  label={mode === 'create' ? `${tEdit('adminPassword')} *` : tEdit('adminPassword')}
                  htmlFor="admin_password"
                  hint={mode === 'edit' ? t('passwordKeepHint') : t('passwordHint')}
                  error={errors.admin_password}
                >
                  <div className="flex gap-2">
                    <div className="relative flex-1">
                      <Input
                        id="admin_password"
                        name="admin_password"
                        type={showPassword ? 'text' : 'password'}
                        value={form.admin_password || ''}
                        onChange={handleChange}
                        placeholder={mode === 'edit' ? t('passwordKeepPlaceholder') : ''}
                        autoComplete="new-password"
                        className="pr-10"
                      />
                      <button
                        type="button"
                        onClick={() => setShowPassword(!showPassword)}
                        className="absolute inset-y-0 right-0 flex items-center pr-3"
                        aria-label={showPassword ? tEdit('hidePassword') : tEdit('showPassword')}
                      >
                        {showPassword ? <EyeOff className="h-4 w-4 text-gray-500" /> : <Eye className="h-4 w-4 text-gray-500" />}
                      </button>
                    </div>
                    <Button
                      type="button"
                      variant="outline"
                      onClick={() => {
                        update({ admin_password: randomPassword() });
                        setShowPassword(true);
                        setErrors((prev) => ({ ...prev, admin_password: '' }));
                      }}
                    >
                      <Dices className="mr-2 h-4 w-4" />
                      {t('generate')}
                    </Button>
                  </div>
                </Field>
              </section>

              <section className="space-y-4">
                <h3 className="font-display text-base font-semibold">{t('sectionMods')}</h3>
                <Field label={tEdit('modIds')} htmlFor="game_mod_ids" hint={t('modsHint')} error={errors.game_mod_ids}>
                  <Input
                    id="game_mod_ids"
                    name="game_mod_ids"
                    value={form.game_mod_ids || ''}
                    onChange={handleChange}
                    placeholder={tEdit('modIdsPlaceholder')}
                  />
                </Field>
                {form.game_mod_ids && <ModList ids={form.game_mod_ids} />}

                <Field label={tEdit('clusterId')} htmlFor="cluster_id" hint={t('clusterHint')} error={errors.cluster_id}>
                  <Input
                    id="cluster_id"
                    name="cluster_id"
                    value={form.cluster_id || ''}
                    onChange={handleChange}
                    placeholder={tEdit('clusterIdPlaceholder')}
                  />
                </Field>
              </section>

              <section className="space-y-3">
                <h3 className="font-display text-base font-semibold">{t('sectionMaintenance')}</h3>
                <label className="flex items-start justify-between gap-4 rounded-lg border p-3">
                  <span>
                    <span className="text-sm font-medium">{t('autoRestart')}</span>
                    <span className="block text-xs text-gray-500">{t('autoRestartHint')}</span>
                  </span>
                  <Switch checked={form.auto_restart ?? true} onCheckedChange={(v) => update({ auto_restart: v })} />
                </label>
                <label className="flex items-start justify-between gap-4 rounded-lg border p-3">
                  <span>
                    <span className="text-sm font-medium">{t('autoUpdate')}</span>
                    <span className="block text-xs text-gray-500">{t('autoUpdateHint')}</span>
                  </span>
                  <Switch checked={form.auto_update ?? false} onCheckedChange={(v) => update({ auto_update: v })} />
                </label>
              </section>
              {/* 回车提交 */}
              <button type="submit" className="hidden" />
            </form>
          </TabsContent>

          <TabsContent value="game_user_settings">
            <GameUserSettingsEditor
              value={form.game_user_settings}
              maxPlayers={form.max_players}
              onChange={(v) => update({ game_user_settings: v })}
            />
          </TabsContent>

          <TabsContent value="game_ini">
            <GameIniEditor value={form.game_ini} onChange={(v) => update({ game_ini: v })} />
          </TabsContent>

          <TabsContent value="server_args">
            <ServerArgsEditor value={form.server_args} onChange={(v) => update({ server_args: v })} />
          </TabsContent>
        </Tabs>

        {/* 操作栏固定在视口底部，长表单不用滚回去找保存按钮 */}
        <div className="sticky bottom-0 z-10 -mx-6 mt-6 -mb-6 flex flex-col-reverse gap-2 rounded-b-lg border-t bg-card/95 px-6 py-4 backdrop-blur sm:flex-row sm:items-center sm:justify-end">
          <div className="text-sm sm:mr-auto" aria-live="polite">
            {hasErrors ? (
              <span className="text-red-700">{t('fixErrors')}</span>
            ) : dirty ? (
              <span className="flex items-center gap-2 text-amber-800">
                <span className="h-2 w-2 rounded-full bg-amber-500" />
                {t('unsaved')}
              </span>
            ) : null}
          </div>
          <Button variant="outline" onClick={cancel} disabled={saving}>
            {tCommon('cancel')}
          </Button>
          {mode === 'edit' && running && (
            <Button variant="outline" onClick={() => submit(true)} disabled={saving}>
              {saving ? <Loader2 className="mr-2 h-4 w-4 animate-spin" /> : <RotateCw className="mr-2 h-4 w-4" />}
              {t('saveAndRestart')}
            </Button>
          )}
          <Button onClick={() => submit(false)} disabled={saving} title={t('saveShortcut')}>
            {saving ? <Loader2 className="mr-2 h-4 w-4 animate-spin" /> : mode === 'create' ? <Plus className="mr-2 h-4 w-4" /> : <Save className="mr-2 h-4 w-4" />}
            {mode === 'create' ? tEdit('createServer') : tEdit('saveChanges')}
          </Button>
        </div>
      </CardContent>
    </Card>
  );
}
