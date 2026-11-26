"use client";

import { useEffect, useMemo, useRef, useState } from 'react';
import { useTranslations } from 'next-intl';
import { Loader2, Dices, Save, Plus, RotateCw, ExternalLink, Search, X, FileUp } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Card, CardContent } from '@/components/ui/card';
import { MapSelector } from '@/components/servers/MapSelector';
import { ModIdInput, ModList, isValidModIds, mergeModIds, modBrowseUrl, normalizeModIds } from '@/components/servers/ModList';
import { SettingTile, SettingControl, PasswordInput } from '@/components/servers/settings/SettingTile';
import { CustomArgsPanel, RawConfigPanel } from '@/components/servers/settings/Panels';
import { ImportConfigDialog } from '@/components/servers/ServerConfigDialogs';
import { applyConfigToForm } from '@/lib/server-config';
import { toast } from '@/lib/toast';
import { randomPassword } from '@/lib/format';
import { indexIni } from '@/lib/ini';
import {
  SETTING_GROUPS,
  EMPTY_SERVER_ARGS,
  allSettings,
  isModified,
  readSetting,
  settingGroupsFor,
  writeSetting,
  type GroupId,
  type SettingDef,
  type SettingGroupDef,
  type SettingSectionDef,
  type SettingsDoc,
  type SettingValue,
} from '@/lib/server-settings';
import { cn } from '@/lib/utils';
import type { Server } from '@/stores/servers';
import { gameOf } from '@/lib/games';

const CLUSTER_ID = /^[A-Za-z0-9][A-Za-z0-9_.-]{0,63}$/;

// 可选的游戏版本（Steam 分支），public 对应后端的空字符串
const GAME_BRANCHES = ['public', 'preaquatica'] as const;

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

// 服务器记录本身的字段所在的分组，用于在出错时跳到对应分组
const FIELD_GROUP: Record<string, GroupId> = {};
for (const group of SETTING_GROUPS) {
  for (const def of allSettings(group)) if (def.source === 'field') FIELD_GROUP[def.key] = group.id;
}

export function ServerForm({ mode, initial, running, saving, onSubmit, onCancel, onDirtyChange }: ServerFormProps) {
  const t = useTranslations('serverForm');
  const tS = useTranslations('serverSettings');
  const tMods = useTranslations('mods');
  const tEdit = useTranslations('servers.edit');
  const tCommon = useTranslations('common');
  const tConfig = useTranslations('serverConfig');
  const [importOpen, setImportOpen] = useState(false);
  const [form, setForm] = useState<Partial<Server>>(initial);
  const [dirty, setDirty] = useState(false);
  const [errors, setErrors] = useState<Record<string, string>>({});
  const [activeGroup, setActiveGroup] = useState<GroupId>('general');
  const [query, setQuery] = useState('');
  // 生成的随机密码直接显示出来，方便复制
  const [generatedPassword, setGeneratedPassword] = useState('');
  const contentRef = useRef<HTMLDivElement>(null);
  // 游戏创建后不能修改；新建页切换游戏时会重新挂载表单。不适用于这个游戏的设置项不显示
  const game = gameOf(form);
  const groups = settingGroupsFor(game);
  // 父组件异步加载到建议端口、默认配置后会更新 initial；用户还没改过时跟着更新
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

  const markDirty = () => {
    if (!dirty) {
      setDirty(true);
      onDirtyChange?.(true);
    }
  };

  const update = (patch: Partial<Server> | ((prev: Partial<Server>) => Partial<Server>)) => {
    setForm((prev) => ({ ...prev, ...(typeof patch === 'function' ? patch(prev) : patch) }));
    markDirty();
  };

  const clearError = (key: string) => {
    if (errors[key]) setErrors((prev) => ({ ...prev, [key]: '' }));
  };

  const setField = (key: keyof Server, value: Server[keyof Server]) => {
    update({ [key]: value } as Partial<Server>);
    clearError(key);
  };

  const setSetting = (def: SettingDef, value: SettingValue | undefined) => {
    update((prev) => writeSetting(def, value, prev));
  };

  // 每次渲染从文本解析，设置项和"配置文件"面板始终一致
  const gusIndex = useMemo(() => indexIni(form.game_user_settings ?? ''), [form.game_user_settings]);
  const gameIndex = useMemo(() => indexIni(form.game_ini ?? ''), [form.game_ini]);
  const doc: SettingsDoc = { gus: gusIndex, game: gameIndex, args: form.server_args ?? EMPTY_SERVER_ARGS };

  const validate = (): Record<string, string> => {
    const next: Record<string, string> = {};
    if (mode === 'create' && !form.admin_password) next.admin_password = t('passwordRequired');
    if (form.admin_password && /[\r\n]/.test(form.admin_password)) next.admin_password = t('passwordNoNewline');
    // ASA 没有查询端口
    const portKeys = game === 'asa' ? (['port', 'rcon_port'] as const) : (['port', 'query_port', 'rcon_port'] as const);
    for (const key of portKeys) {
      const v = Number(form[key]);
      if (!Number.isInteger(v) || v < 1 || v > 65535) next[key] = t('portRange');
    }
    const ports = [Number(form.port), Number(form.port) + 1, ...(game === 'asa' ? [] : [Number(form.query_port)]), Number(form.rcon_port)];
    if (!next.port && !next.query_port && !next.rcon_port && new Set(ports).size !== ports.length) {
      next.port = game === 'asa' ? t('portDuplicateAsa') : t('portDuplicate');
    }
    const mp = Number(form.max_players);
    if (!Number.isInteger(mp) || mp < 1 || mp > 200) next.max_players = t('maxPlayersRange');
    if (form.cluster_id && !CLUSTER_ID.test(form.cluster_id.trim())) next.cluster_id = t('clusterInvalid');
    if (form.game_mod_ids && !isValidModIds(form.game_mod_ids, game)) next.game_mod_ids = t('modsInvalid');
    return next;
  };

  const showGroup = (id: GroupId) => {
    setActiveGroup(id);
    setQuery('');
    contentRef.current?.scrollIntoView({ block: 'nearest' });
  };

  const submit = async (restart: boolean) => {
    const next = validate();
    setErrors(next);
    const firstError = Object.keys(next)[0];
    if (firstError) {
      // 跳到第一个出错字段所在的分组
      showGroup(FIELD_GROUP[firstError] ?? 'general');
      return;
    }
    // 收录站、联机房、ArkApi、自动重启和自动更新不在这个表单里编辑（在详情页切换），编辑时不提交，
    // 免得覆盖别处的修改；只有导入配置改了它们时才提交。新建时照常提交（导入配置可能带了这些值）
    const rest = { ...form };
    delete rest.publish;
    delete rest.link;
    if (mode === 'edit') {
      for (const key of ['arkapi', 'auto_restart', 'auto_update'] as const) {
        if (form[key] === initial[key]) delete rest[key];
      }
    }
    const data: Partial<Server> = {
      ...rest,
      cluster_id: form.cluster_id?.trim() ?? '',
      game_mod_ids: normalizeModIds(form.game_mod_ids ?? '', game),
      server_args: form.server_args && {
        ...form.server_args,
        custom_args: form.server_args.custom_args.map((a) => a.trim()).filter(Boolean),
      },
    };
    // ASA 没有查询端口，后端也会忽略
    if (game === 'asa') delete data.query_port;
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

  // ---------- 文案 ----------

  const fieldLabel = (key: string): string => tS(`fields.${key}`);
  const fieldHint = (key: string): string | undefined => {
    if (key === 'admin_password') return mode === 'edit' ? t('passwordKeepHint') : t('passwordHint');
    // ASA 有不同说明的字段（如 Mod 是 CurseForge 项目 ID）
    if (game === 'asa' && tS.has(`fieldHintsAsa.${key}`)) return tS(`fieldHintsAsa.${key}`);
    return tS.has(`fieldHints.${key}`) ? tS(`fieldHints.${key}`) : undefined;
  };
  const labelOf = (def: SettingDef): string =>
    def.source === 'field' ? fieldLabel(def.key) : tS.has(`params.${def.key}`) ? tS(`params.${def.key}`) : def.key;
  const hintOf = (def: SettingDef): string | undefined =>
    def.source === 'field' ? fieldHint(def.key) : tS.has(`hints.${def.key}`) ? tS(`hints.${def.key}`) : undefined;
  const sectionTitle = (group: GroupId, section: string) => tS(`sections.${group}.${section}`);
  const sectionDesc = (group: GroupId, section: string) =>
    game === 'asa' && tS.has(`sectionHintsAsa.${group}.${section}`)
      ? tS(`sectionHintsAsa.${group}.${section}`)
      : tS.has(`sectionHints.${group}.${section}`)
        ? tS(`sectionHints.${group}.${section}`)
        : undefined;

  // ---------- 服务器字段 ----------

  const renderField = (def: SettingDef) => {
    const key = def.key;
    const common = { key, label: fieldLabel(key), htmlFor: `f-${key}`, hint: fieldHint(key), error: errors[key], wide: def.wide };

    switch (key) {
      case 'session_name':
      case 'cluster_id':
        return (
          <SettingTile {...common}>
            <Input
              id={common.htmlFor}
              value={form[key] || ''}
              onChange={(e) => setField(key, e.target.value)}
              placeholder={key === 'cluster_id' ? tEdit('clusterIdPlaceholder') : undefined}
            />
          </SettingTile>
        );
      case 'map':
        return (
          <SettingTile {...common}>
            <MapSelector value={form.map || ''} onChange={(value) => setField('map', value)} game={game} />
          </SettingTile>
        );
      case 'game_branch': {
        // Select 不能用空字符串作值，蜗牛海洋正式版（后端为空）在界面上用 public 表示
        const branch = form.game_branch || 'public';
        const switched = mode === 'edit' && (form.game_branch || '') !== (initial.game_branch || '');
        return (
          <SettingTile {...common}>
            <Select value={branch} onValueChange={(v) => setField('game_branch', v === 'public' ? '' : v)}>
              <SelectTrigger id={common.htmlFor} className="w-full">
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                {GAME_BRANCHES.map((b) => (
                  <SelectItem key={b} value={b}>
                    {t(`gameBranches.${b}`)}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
            {switched && <p className="mt-1.5 text-xs text-amber-800">{t('gameBranchSwitched')}</p>}
          </SettingTile>
        );
      }
      case 'max_players':
      case 'port':
      case 'query_port':
      case 'rcon_port':
        return (
          <SettingTile {...common}>
            <Input
              id={common.htmlFor}
              type="number"
              className="num"
              value={form[key] || ''}
              onChange={(e) => setField(key, e.target.value === '' ? 0 : Number(e.target.value))}
              aria-invalid={Boolean(errors[key]) || undefined}
            />
          </SettingTile>
        );
      case 'admin_password':
        return (
          <SettingTile {...common} label={mode === 'create' ? `${common.label} *` : common.label}>
            <div className="flex gap-2">
              <PasswordInput
                id={common.htmlFor}
                className="flex-1"
                value={form.admin_password || ''}
                onChange={(v) => setField('admin_password', v)}
                placeholder={mode === 'edit' ? t('passwordKeepPlaceholder') : ''}
              />
              <Button
                type="button"
                variant="outline"
                onClick={() => {
                  const pwd = randomPassword();
                  setGeneratedPassword(pwd);
                  setField('admin_password', pwd);
                }}
              >
                <Dices />
                {t('generate')}
              </Button>
            </div>
            {generatedPassword && form.admin_password === generatedPassword && (
              <p className="mt-1.5 font-mono text-xs text-gray-700">{form.admin_password}</p>
            )}
          </SettingTile>
        );
      case 'game_mod_ids': {
        const ids = form.game_mod_ids ?? '';
        return (
          <SettingTile {...common}>
            <div className="space-y-2">
              <ModIdInput
                id={common.htmlFor}
                existing={ids}
                game={game}
                onAdd={(added) => update({ game_mod_ids: mergeModIds(ids, added.join(',')) })}
              />
              {ids ? (
                <ModList ids={ids} game={game} onChange={(value) => update({ game_mod_ids: value })} />
              ) : (
                <p className="text-xs text-gray-500">{tMods('none')}</p>
              )}
              {/* 想一次粘贴一串 ID 或修格式时用原始输入框，平时直接在上面添加 */}
              <details>
                <summary className="cursor-pointer text-xs text-gray-500 select-none hover:text-gray-800">{tMods('pasteIds')}</summary>
                <Input
                  id={`${common.htmlFor}-raw`}
                  className="mt-1.5 font-mono"
                  value={ids}
                  onChange={(e) => setField('game_mod_ids', e.target.value)}
                  onBlur={() => {
                    // 粘贴的创意工坊 / CurseForge 链接换成 ID
                    if (isValidModIds(ids, game) && normalizeModIds(ids, game) !== ids) update({ game_mod_ids: normalizeModIds(ids, game) });
                  }}
                  placeholder={tEdit('modIdsPlaceholder')}
                />
              </details>
              <a
                href={modBrowseUrl(game)}
                target="_blank"
                rel="noreferrer"
                className="inline-flex items-center gap-1 text-xs text-primary hover:underline"
              >
                {game === 'asa' ? t('browseCurseforge') : t('browseWorkshop')}
                <ExternalLink className="size-3" />
              </a>
            </div>
          </SettingTile>
        );
      }
      default:
        return null;
    }
  };

  const renderSetting = (def: SettingDef) => {
    if (def.source === 'field') return renderField(def);
    const value = readSetting(def, doc);
    const modified = isModified(def, value);
    const id = `s-${def.source}-${def.key}`;
    return (
      <SettingTile
        key={id}
        label={labelOf(def)}
        htmlFor={id}
        name={def.source === 'arg' ? `-${def.key}` : def.source === 'query' ? `?${def.key}` : def.key}
        hint={hintOf(def)}
        modified={modified}
        onReset={() => setSetting(def, undefined)}
        inline={def.type === 'boolean'}
        wide={def.wide}
      >
        <SettingControl def={def} id={id} value={value} onChange={(v) => setSetting(def, v)} />
      </SettingTile>
    );
  };

  const renderPanel = (section: SettingSectionDef) => {
    if (section.panel === 'customArgs') {
      const args = form.server_args ?? EMPTY_SERVER_ARGS;
      return (
        <CustomArgsPanel
          value={args.custom_args ?? []}
          onChange={(custom_args) => update((prev) => ({ server_args: { ...(prev.server_args ?? EMPTY_SERVER_ARGS), custom_args } }))}
        />
      );
    }
    if (section.panel === 'rawConfig') {
      return <RawConfigPanel gameUserSettings={form.game_user_settings ?? ''} gameIni={form.game_ini ?? ''} onChange={update} />;
    }
    return null;
  };

  const renderSection = (group: GroupId, section: SettingSectionDef, settings: SettingDef[]) => {
    const desc = sectionDesc(group, section.id);
    return (
      <section key={`${group}-${section.id}`} className="space-y-3">
        <div>
          <h3 className="hud-caps text-xs font-semibold tracking-wider text-gray-700">{sectionTitle(group, section.id)}</h3>
          {desc && <p className="mt-1 text-xs text-gray-500">{desc}</p>}
        </div>
        <div className="grid grid-cols-1 gap-3 sm:grid-cols-2">{section.panel ? renderPanel(section) : settings.map(renderSetting)}</div>
      </section>
    );
  };

  // ---------- 分组统计 ----------

  const modifiedCount = (group: SettingGroupDef) =>
    allSettings(group).filter((def) => def.source !== 'field' && isModified(def, readSetting(def, doc))).length;
  const groupHasError = (group: SettingGroupDef) => allSettings(group).some((def) => def.source === 'field' && errors[def.key]);

  // ---------- 搜索 ----------

  const q = query.trim().toLowerCase();
  const searchResults = q
    ? groups.map((group) => ({
        group,
        sections: group.sections
          .map((section) => {
            if (sectionTitle(group.id, section.id).toLowerCase().includes(q)) return { section, settings: section.settings };
            const settings = section.settings.filter((def) =>
              [labelOf(def), def.key, hintOf(def) ?? ''].some((s) => s.toLowerCase().includes(q))
            );
            return { section, settings };
          })
          .filter(({ section, settings }) => settings.length > 0 || (section.panel && sectionTitle(group.id, section.id).toLowerCase().includes(q))),
      })).filter((r) => r.sections.length > 0)
    : [];

  const current = groups.find((g) => g.id === activeGroup) ?? groups[0];

  return (
    <Card className="py-0">
      <CardContent className="px-0">
        <div className="grid md:grid-cols-[14rem_minmax(0,1fr)]">
          {/* 分组导航：桌面端在左侧，窄屏时横向滚动 */}
          <aside className="border-b p-3 md:border-r md:border-b-0">
            <div className="md:sticky md:top-4 space-y-3">
              <div className="relative">
                <Search className="pointer-events-none absolute top-1/2 left-2.5 size-4 -translate-y-1/2 text-gray-500" />
                <Input
                  value={query}
                  onChange={(e) => setQuery(e.target.value)}
                  placeholder={tS('searchPlaceholder')}
                  className="pr-8 pl-8"
                  aria-label={tS('searchPlaceholder')}
                />
                {query && (
                  <button
                    type="button"
                    onClick={() => setQuery('')}
                    className="absolute inset-y-0 right-0 flex items-center pr-2.5 text-gray-500 hover:text-gray-900"
                    aria-label={tS('clearSearch')}
                  >
                    <X className="size-4" />
                  </button>
                )}
              </div>
              <nav className="-mx-3 flex gap-1 overflow-x-auto px-3 md:mx-0 md:flex-col md:overflow-visible md:px-0">
                {groups.map((group) => {
                  const Icon = group.icon;
                  const active = !q && group.id === activeGroup;
                  const count = modifiedCount(group);
                  return (
                    <button
                      key={group.id}
                      type="button"
                      onClick={() => showGroup(group.id)}
                      aria-current={active ? 'page' : undefined}
                      className={cn(
                        'flex shrink-0 items-center gap-2 rounded-md px-2.5 py-1.5 text-left text-sm whitespace-nowrap transition-colors',
                        active ? 'bg-primary/15 text-primary' : 'text-gray-700 hover:bg-gray-100 hover:text-gray-900'
                      )}
                    >
                      <Icon className="size-4 shrink-0" />
                      <span className="flex-1">{tS(`groups.${group.id}.title`)}</span>
                      {groupHasError(group) ? (
                        <span className="size-2 rounded-full bg-red-500" />
                      ) : count > 0 ? (
                        <span className="num rounded bg-gray-100 px-1.5 text-[11px] text-gray-600" title={tS('modifiedCount', { count })}>
                          {count}
                        </span>
                      ) : null}
                    </button>
                  );
                })}
              </nav>
            </div>
          </aside>

          <form
            onSubmit={(e) => {
              e.preventDefault();
              submit(false);
            }}
            className="min-w-0 scroll-mt-4 p-4 sm:p-6"
          >
            <div ref={contentRef} className="scroll-mt-20" />
            {q ? (
              searchResults.length === 0 ? (
                <p className="py-12 text-center text-sm text-gray-500">{tS('noResults')}</p>
              ) : (
                <div className="space-y-8">
                  {searchResults.map(({ group, sections }) => (
                    <div key={group.id} className="space-y-5">
                      <button
                        type="button"
                        onClick={() => showGroup(group.id)}
                        className="font-display text-base font-semibold hover:text-primary"
                      >
                        {tS(`groups.${group.id}.title`)}
                      </button>
                      {sections.map(({ section, settings }) => renderSection(group.id, section, settings))}
                    </div>
                  ))}
                </div>
              )
            ) : (
              <div className="space-y-7">
                <div>
                  <h2 className="font-display text-lg font-semibold">{tS(`groups.${current.id}.title`)}</h2>
                  <p className="mt-0.5 text-sm text-gray-500">{tS(`groups.${current.id}.desc`)}</p>
                </div>
                {current.sections.map((section) => renderSection(current.id, section, section.settings))}
              </div>
            )}
            {/* 回车提交 */}
            <button type="submit" className="hidden" />
          </form>
        </div>

        {/* 操作栏固定在视口底部，长表单不用滚回去找保存按钮 */}
        <div className="sticky bottom-0 z-10 flex flex-col-reverse gap-2 rounded-b-sm border-t bg-card/95 px-6 py-4 backdrop-blur sm:flex-row sm:items-center sm:justify-end">
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
          <Button variant="outline" onClick={() => setImportOpen(true)} disabled={saving} title={tConfig('importHint')}>
            <FileUp />
            {tConfig('import')}
          </Button>
          <Button variant="outline" onClick={cancel} disabled={saving}>
            {tCommon('cancel')}
          </Button>
          {mode === 'edit' && running && (
            <Button variant="outline" onClick={() => submit(true)} disabled={saving}>
              {saving ? <Loader2 className="animate-spin" /> : <RotateCw />}
              {t('saveAndRestart')}
            </Button>
          )}
          <Button onClick={() => submit(false)} disabled={saving} title={t('saveShortcut')}>
            {saving ? <Loader2 className="animate-spin" /> : mode === 'create' ? <Plus /> : <Save />}
            {mode === 'create' ? tEdit('createServer') : tEdit('saveChanges')}
          </Button>
        </div>
      </CardContent>
      <ImportConfigDialog
        open={importOpen}
        onOpenChange={setImportOpen}
        excludeServerId={initial.id}
        targetGame={game}
        onImport={(config, parts) => {
          update((prev) => applyConfigToForm(prev, config, parts));
          setErrors({});
          toast.success(tConfig('imported'));
        }}
      />
    </Card>
  );
}
