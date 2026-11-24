"use client";

import { useMemo, useState } from 'react';
import { useTranslations } from 'next-intl';
import { ChevronRight, Copy, Pencil, Plus, Search, Trash2, Undo2, X } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Switch } from '@/components/ui/switch';
import { Textarea } from '@/components/ui/textarea';
import {
  NUMBER_PATTERN,
  childCount,
  duplicateChild,
  filterJsonTree,
  hasKey,
  isContainer,
  newNode,
  removeChild,
  type JsonContainer,
  type JsonFilter,
  type JsonKind,
  type JsonNode,
  type JsonPath,
} from '@/lib/json-tree';
import { cn } from '@/lib/utils';

// 一个容器一次最多渲染这么多子项，ShopItems 这类上百项的列表点"显示全部"再展开
const CHILD_LIMIT = 100;

// 顶层容器子项不多时默认展开
const AUTO_OPEN_CHILDREN = 10;

type AddKind = Exclude<JsonKind, 'null'>;
const ADD_KINDS: AddKind[] = ['string', 'number', 'boolean', 'object', 'array'];

// mergeKey 相同的连续修改（同一个输入框里打字）只记一次撤销
export type TreeUpdate = (path: JsonPath, fn: (node: JsonNode) => JsonNode, mergeKey?: string) => void;

interface Ctx {
  filter: JsonFilter | null;
  update: TreeUpdate;
}

// 字段名、字符串值是否像在描述这一项（ArkShop 的 Description、Title 等），用于折叠时的摘要
const SUMMARY_KEYS = /^(description|title|name|displayname|blueprint)$/i;

function summaryText(node: JsonContainer): string {
  if (node.kind !== 'object') return '';
  const preferred = node.entries.find((e) => SUMMARY_KEYS.test(e.key) && e.value.kind === 'string');
  const entry = preferred ?? node.entries.find((e) => e.value.kind === 'string');
  return entry && entry.value.kind === 'string' ? entry.value.value : '';
}

// 设置项的按层级编辑：对象和数组可以展开，标量按类型显示输入框、开关
export function JsonTreeEditor({
  value,
  update,
  canUndo,
  onUndo,
}: {
  value: JsonContainer;
  update: TreeUpdate;
  canUndo: boolean;
  onUndo: () => void;
}) {
  const t = useTranslations('plugins');
  const [query, setQuery] = useState('');
  const filter = useMemo(() => filterJsonTree(value, query), [value, query]);
  const ctx: Ctx = { filter, update };
  const noMatch = filter !== null && filter.shown.size === 0;

  return (
    <div className="space-y-2">
      <div className="flex items-center gap-2">
        <div className="relative flex-1">
          <Search className="pointer-events-none absolute top-1/2 left-2.5 size-3.5 -translate-y-1/2 text-muted-foreground" />
          <Input
            value={query}
            onChange={(e) => setQuery(e.target.value)}
            placeholder={t('searchFields')}
            className="h-8 pr-8 pl-8 text-xs md:text-xs"
          />
          {query && (
            <button
              type="button"
              className="absolute top-1/2 right-2 -translate-y-1/2 text-muted-foreground hover:text-foreground"
              onClick={() => setQuery('')}
              aria-label={t('clearSearch')}
            >
              <X className="size-3.5" />
            </button>
          )}
        </div>
        <Button variant="outline" size="sm" onClick={onUndo} disabled={!canUndo}>
          <Undo2 className="mr-1.5 size-3.5" />
          {t('undo')}
        </Button>
      </div>
      <div className="max-h-[55vh] overflow-y-auto pr-1">
        {noMatch ? (
          <p className="py-6 text-center text-sm text-muted-foreground">{t('noMatch')}</p>
        ) : (
          <Children node={value} path={[]} depth={0} ctx={ctx} />
        )}
      </div>
    </div>
  );
}

function Children({ node, path, depth, ctx }: { node: JsonContainer; path: JsonPath; depth: number; ctx: Ctx }) {
  const t = useTranslations('plugins');
  const [limit, setLimit] = useState(CHILD_LIMIT);
  const all =
    node.kind === 'object'
      ? node.entries.map((e, index) => ({ key: e.key as string | null, value: e.value, index }))
      : node.items.map((value, index) => ({ key: null as string | null, value, index }));
  const visible = ctx.filter ? all.filter((c) => ctx.filter!.shown.has(c.value.id)) : all;
  const rendered = visible.slice(0, limit);

  return (
    <div className="space-y-0.5">
      {all.length === 0 && <p className="px-1 py-1 text-xs text-muted-foreground">{t('emptyContainer')}</p>}
      {rendered.map((c) => (
        <Row key={c.value.id} entryKey={c.key} index={c.index} node={c.value} parent={node} parentPath={path} depth={depth} ctx={ctx} />
      ))}
      {visible.length > rendered.length && (
        <Button variant="ghost" size="sm" className="text-xs" onClick={() => setLimit(visible.length)}>
          {t('showAll', { count: visible.length })}
        </Button>
      )}
      {!ctx.filter && <AddChild node={node} path={path} ctx={ctx} />}
    </div>
  );
}

function Row({
  entryKey,
  index,
  node,
  parent,
  parentPath,
  depth,
  ctx,
}: {
  entryKey: string | null;
  index: number;
  node: JsonNode;
  parent: JsonContainer;
  parentPath: JsonPath;
  depth: number;
  ctx: Ctx;
}) {
  const t = useTranslations('plugins');
  const path = [...parentPath, index];
  const container = isContainer(node) ? node : null;
  const [expanded, setExpanded] = useState(() => {
    if (!container) return false;
    const n = childCount(container);
    return n === 0 || (depth === 0 && n <= AUTO_OPEN_CHILDREN);
  });
  const [renaming, setRenaming] = useState(false);
  const open = container !== null && (expanded || (ctx.filter?.open.has(node.id) ?? false));
  const label = entryKey ?? `#${index + 1}`;

  const setValue = (next: JsonNode, merge = true) => ctx.update(path, () => next, merge ? `v${node.id}` : undefined);
  const updateParent = (fn: (p: JsonContainer) => JsonContainer) => ctx.update(parentPath, (p) => fn(p as JsonContainer));

  const keyLabel = renaming ? (
    <RenameInput
      initial={entryKey ?? ''}
      taken={(k) => hasKey(parent, k, index)}
      onDone={(k) => {
        setRenaming(false);
        if (k !== null && k !== entryKey) {
          updateParent((p) =>
            p.kind === 'object' ? { ...p, entries: p.entries.map((e, i) => (i === index ? { ...e, key: k } : e)) } : p,
          );
        }
      }}
    />
  ) : (
    <span
      className={cn('truncate font-mono text-xs', entryKey === null ? 'text-muted-foreground' : 'text-foreground/90')}
      title={label}
    >
      {label}
    </span>
  );

  const actions = (
    <div className="flex shrink-0 items-center opacity-100 transition-opacity group-focus-within:opacity-100 group-hover:opacity-100 [@media(hover:hover)]:opacity-0">
      {entryKey !== null && (
        <IconButton label={t('renameField')} onClick={() => setRenaming(true)}>
          <Pencil className="size-3.5" />
        </IconButton>
      )}
      <IconButton label={t('duplicateField')} onClick={() => updateParent((p) => duplicateChild(p, index))}>
        <Copy className="size-3.5" />
      </IconButton>
      <IconButton label={t('removeField')} onClick={() => updateParent((p) => removeChild(p, index))}>
        <Trash2 className="size-3.5 text-red-600" />
      </IconButton>
    </div>
  );

  if (container) {
    const n = childCount(container);
    const summary = summaryText(container);
    const chevron = (
      <ChevronRight className={cn('size-3.5 shrink-0 text-muted-foreground transition-transform', open && 'rotate-90')} />
    );
    const badge = (
      <span className="shrink-0 rounded-sm border px-1.5 text-[10px] leading-4 text-muted-foreground">
        {container.kind === 'object' ? t('objectCount', { count: n }) : t('arrayCount', { count: n })}
      </span>
    );
    return (
      <div>
        <div className="group flex min-h-8 items-center gap-1 rounded-sm px-1 hover:bg-muted/40">
          {renaming ? (
            // 改名时输入框不能放在展开按钮里
            <div className="flex min-w-0 flex-1 items-center gap-1.5 py-1">
              {chevron}
              <div className="w-44 shrink-0">{keyLabel}</div>
              {badge}
            </div>
          ) : (
            <button
              type="button"
              className="flex min-w-0 flex-1 items-center gap-1.5 py-1 text-left"
              onClick={() => setExpanded(!open)}
              aria-expanded={open}
            >
              {chevron}
              <span className="flex max-w-[45%] shrink-0 items-center">{keyLabel}</span>
              {badge}
              {summary && !open && <span className="truncate text-xs text-muted-foreground">{summary}</span>}
            </button>
          )}
          {actions}
        </div>
        {open && (
          <div className="ml-[0.7rem] border-l border-border/70 pl-3">
            <Children node={container} path={path} depth={depth + 1} ctx={ctx} />
          </div>
        )}
      </div>
    );
  }

  return (
    <div className="group flex min-h-8 items-center gap-2 rounded-sm px-1 hover:bg-muted/40">
      <span className="w-3.5 shrink-0" />
      <div className="flex w-40 shrink-0 items-center">{keyLabel}</div>
      <div className="flex min-w-0 flex-1 items-center gap-2">
        <ScalarEditor node={node} label={label} onChange={setValue} />
      </div>
      {actions}
    </div>
  );
}

function ScalarEditor({ node, label, onChange }: { node: JsonNode; label: string; onChange: (next: JsonNode, merge?: boolean) => void }) {
  if (node.kind === 'boolean') {
    return (
      <label className="flex items-center gap-2">
        <Switch checked={node.value} onCheckedChange={(value) => onChange({ ...node, value }, false)} aria-label={label} />
        <span className="font-mono text-xs text-muted-foreground">{node.value ? 'true' : 'false'}</span>
      </label>
    );
  }
  if (node.kind === 'number') {
    return <NumberInput raw={node.raw} label={label} onCommit={(raw) => onChange({ ...node, raw })} />;
  }
  if (node.kind === 'string') {
    // 输入框会吃掉换行，多行文本用文本框
    if (node.value.includes('\n')) {
      return (
        <Textarea
          value={node.value}
          onChange={(e) => onChange({ ...node, value: e.target.value })}
          aria-label={label}
          spellCheck={false}
          className="min-h-9 py-1.5 font-mono text-xs md:text-xs"
        />
      );
    }
    return (
      <Input
        value={node.value}
        onChange={(e) => onChange({ ...node, value: e.target.value })}
        aria-label={label}
        spellCheck={false}
        className="h-8 font-mono text-xs md:text-xs"
      />
    );
  }
  return <span className="font-mono text-xs text-muted-foreground">null</span>;
}

// 数字按原文编辑（大整数不丢精度），不合法时标红、不写回
function NumberInput({ raw, label, onCommit }: { raw: string; label: string; onCommit: (raw: string) => void }) {
  const t = useTranslations('plugins');
  const [draft, setDraft] = useState(raw);
  const [synced, setSynced] = useState(raw);
  // 撤销等外部修改时跟着变
  if (raw !== synced) {
    setSynced(raw);
    setDraft(raw);
  }
  const valid = NUMBER_PATTERN.test(draft.trim());
  return (
    <Input
      value={draft}
      inputMode="decimal"
      aria-label={label}
      aria-invalid={!valid}
      title={valid ? undefined : t('fieldNumberInvalid')}
      className="num h-8 max-w-56 font-mono text-xs md:text-xs"
      onChange={(e) => {
        const text = e.target.value;
        setDraft(text);
        const trimmed = text.trim();
        if (NUMBER_PATTERN.test(trimmed) && trimmed !== raw) {
          setSynced(trimmed);
          onCommit(trimmed);
        }
      }}
      onBlur={() => setDraft(raw)}
    />
  );
}

function RenameInput({ initial, taken, onDone }: { initial: string; taken: (key: string) => boolean; onDone: (key: string | null) => void }) {
  const t = useTranslations('plugins');
  const [value, setValue] = useState(initial);
  const invalid = value === '' || (value !== initial && taken(value));
  const commit = () => onDone(invalid ? null : value);
  return (
    <Input
      autoFocus
      value={value}
      onChange={(e) => setValue(e.target.value)}
      onKeyDown={(e) => {
        if (e.key === 'Enter') {
          e.preventDefault();
          if (!invalid) commit();
        } else if (e.key === 'Escape') {
          e.preventDefault();
          e.stopPropagation();
          onDone(null);
        }
      }}
      onBlur={commit}
      aria-invalid={invalid}
      title={invalid ? t('fieldKeyInvalid') : undefined}
      className="h-7 font-mono text-xs md:text-xs"
    />
  );
}

function AddChild({ node, path, ctx }: { node: JsonContainer; path: JsonPath; ctx: Ctx }) {
  const t = useTranslations('plugins');
  const [open, setOpen] = useState(false);
  const [key, setKey] = useState('');
  const lastKind = node.kind === 'array' ? node.items[node.items.length - 1]?.kind : undefined;
  const [kind, setKind] = useState<AddKind>(lastKind && lastKind !== 'null' ? lastKind : 'string');
  const isObject = node.kind === 'object';
  const keyInvalid = isObject && (key === '' || hasKey(node, key));

  if (!open) {
    return (
      <button
        type="button"
        className="flex items-center gap-1 rounded-sm px-1 py-1 text-xs text-muted-foreground hover:text-primary"
        onClick={() => setOpen(true)}
      >
        <Plus className="size-3.5" />
        {isObject ? t('addField') : t('addItem')}
      </button>
    );
  }

  const add = () => {
    if (keyInvalid) return;
    const value = newNode(kind);
    ctx.update(path, (p) => {
      const c = p as JsonContainer;
      return c.kind === 'object' ? { ...c, entries: [...c.entries, { key, value }] } : { ...c, items: [...c.items, value] };
    });
    setKey('');
    setOpen(false);
  };

  return (
    <div className="flex flex-wrap items-center gap-2 px-1 py-1">
      {isObject && (
        <Input
          autoFocus
          value={key}
          onChange={(e) => setKey(e.target.value)}
          onKeyDown={(e) => {
            if (e.key === 'Enter') add();
            else if (e.key === 'Escape') {
              e.stopPropagation();
              setOpen(false);
            }
          }}
          placeholder={t('newFieldKey')}
          aria-invalid={key !== '' && keyInvalid}
          title={key !== '' && keyInvalid ? t('fieldKeyInvalid') : undefined}
          className="h-8 w-44 font-mono text-xs md:text-xs"
        />
      )}
      <Select value={kind} onValueChange={(v) => setKind(v as AddKind)}>
        <SelectTrigger size="sm" className="w-28 text-xs">
          <SelectValue />
        </SelectTrigger>
        <SelectContent>
          {ADD_KINDS.map((k) => (
            <SelectItem key={k} value={k}>
              {t(`kind.${k}`)}
            </SelectItem>
          ))}
        </SelectContent>
      </Select>
      <Button size="sm" variant="outline" onClick={add} disabled={keyInvalid}>
        <Plus className="mr-1 size-3.5" />
        {isObject ? t('addField') : t('addItem')}
      </Button>
      <Button size="sm" variant="ghost" onClick={() => setOpen(false)}>
        <X className="size-3.5" />
      </Button>
    </div>
  );
}

function IconButton({ label, onClick, children }: { label: string; onClick: () => void; children: React.ReactNode }) {
  return (
    <Button variant="ghost" size="icon" className="h-7 w-7" onClick={onClick} title={label} aria-label={label}>
      {children}
    </Button>
  );
}
