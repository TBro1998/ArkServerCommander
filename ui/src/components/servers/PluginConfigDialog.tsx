"use client";

import { useState } from 'react';
import { useTranslations } from 'next-intl';
import axios from 'axios';
import { Loader2 } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from '@/components/ui/dialog';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { Textarea } from '@/components/ui/textarea';
import { gameOf, pluginApiName } from '@/lib/games';
import { useApiError } from '@/lib/api-error';
import {
  detectJsonFormat,
  isContainer,
  parseJsonTree,
  serializeJsonTree,
  updateAt,
  type JsonContainer,
} from '@/lib/json-tree';
import { toast } from '@/lib/toast';
import type { Server } from '@/stores/servers';
import { JsonTreeEditor, type TreeUpdate } from './JsonTreeEditor';

// 正在编辑的配置：name 为 null 表示插件框架本体（ArkApi / AsaApi）
export interface ConfigTarget {
  name: string | null;
  content: string;
}

type ConfigMode = 'fields' | 'source';

// 撤销最多保留这么多步
const HISTORY_LIMIT = 200;

export function configUrl(serverId: number, name: string | null) {
  return name === null ? `/api/servers/${serverId}/arkapi/config` : `/api/servers/${serverId}/plugins/${encodeURIComponent(name)}/config`;
}

// 能按字段编辑时返回树，否则返回错误信息
function parseContainer(content: string): { tree: JsonContainer } | { error: 'notContainer' | string } {
  try {
    const tree = parseJsonTree(content);
    return isContainer(tree) ? { tree } : { error: 'notContainer' };
  } catch (e) {
    return { error: (e as Error).message };
  }
}

// config.json 编辑弹窗：按字段（树形）或源码编辑。每次打开都重新挂载，状态从 target 初始化
export function PluginConfigDialog({ server, target, onClose, onSaved }: { server: Server; target: ConfigTarget; onClose: () => void; onSaved: () => void }) {
  const t = useTranslations('plugins');
  const tCommon = useTranslations('common');
  const apiError = useApiError();
  const [format] = useState(() => detectJsonFormat(target.content));
  const [initial] = useState(() => parseContainer(target.content));
  const [mode, setMode] = useState<ConfigMode>('tree' in initial ? 'fields' : 'source');
  // 源码模式下的文本；按字段编辑时以 tree 为准，切回源码或保存时才写回文本
  const [content, setContent] = useState(target.content);
  const [tree, setTree] = useState<JsonContainer | null>('tree' in initial ? initial.tree : null);
  const [treeDirty, setTreeDirty] = useState(false);
  const [history, setHistory] = useState<{ past: JsonContainer[]; mergeKey?: string }>({ past: [] });
  const [error, setError] = useState('');
  const [saving, setSaving] = useState(false);

  const errorText = (e: string) => (e === 'notContainer' ? t('rootNotContainer') : t('invalidJson', { message: e }));

  const update: TreeUpdate = (path, fn, mergeKey) => {
    if (!tree) return;
    const next = updateAt(tree, path, fn) as JsonContainer;
    if (next === tree) return;
    setHistory((h) =>
      mergeKey !== undefined && h.mergeKey === mergeKey ? h : { past: [...h.past.slice(-(HISTORY_LIMIT - 1)), tree], mergeKey },
    );
    setTree(next);
    setTreeDirty(true);
    if (error) setError('');
  };

  const undo = () => {
    const prev = history.past[history.past.length - 1];
    if (!prev) return;
    setHistory({ past: history.past.slice(0, -1) });
    setTree(prev);
    setTreeDirty(true);
  };

  // 当前要保存的文本：按字段编辑过就从树重新生成，否则保留原文
  const currentText = (): string => (mode === 'fields' && tree && treeDirty ? serializeJsonTree(tree, format) : content);

  const switchMode = (next: ConfigMode) => {
    if (next === mode) return;
    if (next === 'source') {
      setContent(currentText());
      setTreeDirty(false);
      setError('');
      setMode('source');
      return;
    }
    const parsed = parseContainer(content);
    if ('error' in parsed) {
      setError(errorText(parsed.error));
      return;
    }
    setTree(parsed.tree);
    setTreeDirty(false);
    setHistory({ past: [] });
    setError('');
    setMode('fields');
  };

  const formatSource = () => {
    try {
      setContent(serializeJsonTree(parseJsonTree(content), format));
      setError('');
    } catch (e) {
      setError(errorText((e as Error).message));
    }
  };

  const save = async () => {
    const text = currentText();
    try {
      parseJsonTree(text);
    } catch (e) {
      setError(errorText((e as Error).message));
      return;
    }
    setSaving(true);
    try {
      await axios.put(configUrl(server.id, target.name), { content: text });
      toast.success(server.status === 'stopped' ? t('configSaved') : t('configSavedRunning'));
      onSaved();
    } catch (err) {
      toast.error(apiError(err, t('configSaveFailed')));
    } finally {
      setSaving(false);
    }
  };

  return (
    <Dialog open onOpenChange={(open) => !open && onClose()}>
      <DialogContent className="sm:max-w-4xl">
        <DialogHeader>
          <DialogTitle>{target.name === null ? t('arkapiConfig', { api: pluginApiName(gameOf(server)) }) : t('configTitle', { name: target.name })}</DialogTitle>
          <DialogDescription>{server.status === 'stopped' ? t('configDesc') : t('configDescRunning')}</DialogDescription>
        </DialogHeader>

        <Tabs value={mode} onValueChange={(value) => switchMode(value as ConfigMode)} className="min-w-0">
          <TabsList variant="line">
            <TabsTrigger value="fields">{t('fieldsMode')}</TabsTrigger>
            <TabsTrigger value="source">{t('sourceMode')}</TabsTrigger>
          </TabsList>

          <TabsContent value="fields" className="space-y-2">
            <p className="text-xs text-muted-foreground">{t('fieldsHint')}</p>
            {tree && <JsonTreeEditor value={tree} update={update} canUndo={history.past.length > 0} onUndo={undo} />}
          </TabsContent>

          <TabsContent value="source">
            <Textarea
              value={content}
              onChange={(e) => {
                setContent(e.target.value);
                if (error) setError('');
              }}
              spellCheck={false}
              className="h-[60vh] [field-sizing:fixed] resize-none bg-terminal text-terminal-foreground font-mono text-xs"
              placeholder="{}"
            />
          </TabsContent>
        </Tabs>

        {error && <p className="text-xs text-red-600">{error}</p>}
        <DialogFooter>
          {mode === 'source' && (
            <Button variant="outline" onClick={formatSource}>
              {t('format')}
            </Button>
          )}
          <Button variant="outline" onClick={onClose}>
            {tCommon('cancel')}
          </Button>
          <Button onClick={save} disabled={saving}>
            {saving && <Loader2 className="mr-1.5 animate-spin" />}
            {tCommon('save')}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
