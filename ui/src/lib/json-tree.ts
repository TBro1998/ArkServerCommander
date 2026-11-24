// config.json 的树形编辑：自己解析 JSON，保留键的顺序和数字的原文。
// JSON.parse 会把 "1"、"2" 这类整数键排到其他键前面，大整数（如 SteamID）会丢精度，
// 所以改写文件时不经过 JS 对象，而是在这棵树上改，再按原文件的缩进写回去。

export type JsonNode =
  | { id: number; kind: 'object'; entries: JsonEntry[] }
  | { id: number; kind: 'array'; items: JsonNode[] }
  | { id: number; kind: 'string'; value: string }
  | { id: number; kind: 'number'; raw: string }
  | { id: number; kind: 'boolean'; value: boolean }
  | { id: number; kind: 'null' };

export type JsonKind = JsonNode['kind'];
export type JsonContainer = Extract<JsonNode, { kind: 'object' | 'array' }>;

export interface JsonEntry {
  key: string;
  value: JsonNode;
}

// 节点 ID 只用作 React 的 key 和展开状态，复制出来的节点要换新 ID
let nextId = 1;
const newId = () => nextId++;

export const NUMBER_PATTERN = /^-?(?:0|[1-9]\d*)(?:\.\d+)?(?:[eE][+-]?\d+)?$/;
const NUMBER_AT = /-?(?:0|[1-9]\d*)(?:\.\d+)?(?:[eE][+-]?\d+)?/y;

export class JsonTreeError extends Error {
  readonly position: number;
  constructor(message: string, position: number) {
    super(message);
    this.position = position;
  }
}

// 解析 JSON 文本（可带 BOM），格式错误时抛出 JsonTreeError
export function parseJsonTree(text: string): JsonNode {
  let i = text.charCodeAt(0) === 0xfeff ? 1 : 0;

  const fail = (message: string): never => {
    throw new JsonTreeError(`${message} (position ${i})`, i);
  };

  const skipWs = () => {
    while (i < text.length) {
      const c = text[i];
      if (c === ' ' || c === '\t' || c === '\n' || c === '\r') i++;
      else break;
    }
  };

  const parseString = (): string => {
    const start = i;
    i++;
    while (i < text.length) {
      const c = text[i];
      if (c === '"') {
        i++;
        try {
          return JSON.parse(text.slice(start, i)) as string;
        } catch {
          i = start;
          return fail('Invalid string');
        }
      }
      i += c === '\\' ? 2 : 1;
    }
    i = start;
    return fail('Unterminated string');
  };

  const expectWord = (word: string) => {
    if (text.startsWith(word, i)) i += word.length;
    else fail('Unexpected token');
  };

  const parseValue = (): JsonNode => {
    skipWs();
    const c = text[i];
    if (c === '{') {
      i++;
      const entries: JsonEntry[] = [];
      skipWs();
      if (text[i] === '}') {
        i++;
        return { id: newId(), kind: 'object', entries };
      }
      for (;;) {
        skipWs();
        if (text[i] !== '"') fail('Expected property name');
        const key = parseString();
        skipWs();
        if (text[i] !== ':') fail("Expected ':'");
        i++;
        entries.push({ key, value: parseValue() });
        skipWs();
        if (text[i] === ',') {
          i++;
          continue;
        }
        if (text[i] === '}') {
          i++;
          return { id: newId(), kind: 'object', entries };
        }
        fail("Expected ',' or '}'");
      }
    }
    if (c === '[') {
      i++;
      const items: JsonNode[] = [];
      skipWs();
      if (text[i] === ']') {
        i++;
        return { id: newId(), kind: 'array', items };
      }
      for (;;) {
        items.push(parseValue());
        skipWs();
        if (text[i] === ',') {
          i++;
          continue;
        }
        if (text[i] === ']') {
          i++;
          return { id: newId(), kind: 'array', items };
        }
        fail("Expected ',' or ']'");
      }
    }
    if (c === '"') return { id: newId(), kind: 'string', value: parseString() };
    if (c === 't') {
      expectWord('true');
      return { id: newId(), kind: 'boolean', value: true };
    }
    if (c === 'f') {
      expectWord('false');
      return { id: newId(), kind: 'boolean', value: false };
    }
    if (c === 'n') {
      expectWord('null');
      return { id: newId(), kind: 'null' };
    }
    NUMBER_AT.lastIndex = i;
    const m = NUMBER_AT.exec(text);
    if (m) {
      i += m[0].length;
      return { id: newId(), kind: 'number', raw: m[0] };
    }
    return fail(i >= text.length ? 'Unexpected end of input' : 'Unexpected token');
  };

  const root = parseValue();
  skipWs();
  if (i < text.length) fail('Unexpected content after JSON');
  return root;
}

export interface JsonFormat {
  indent: string;
  eol: string;
  finalNewline: boolean;
}

// 按原文件的缩进（制表符或几个空格）和换行符写回，减少无关的改动
export function detectJsonFormat(text: string): JsonFormat {
  const m = /\n([ \t]+)\S/.exec(text);
  let indent = '  ';
  if (m) indent = m[1].startsWith('\t') ? '\t' : ' '.repeat(Math.min(m[1].length, 8));
  return {
    indent,
    eol: text.includes('\r\n') ? '\r\n' : '\n',
    finalNewline: /\n\s*$/.test(text),
  };
}

export function serializeJsonTree(root: JsonNode, format: JsonFormat = { indent: '  ', eol: '\n', finalNewline: false }): string {
  const { indent, eol } = format;
  const write = (node: JsonNode, pad: string): string => {
    switch (node.kind) {
      case 'object': {
        if (node.entries.length === 0) return '{}';
        const inner = pad + indent;
        const body = node.entries.map((e) => `${inner}${JSON.stringify(e.key)}: ${write(e.value, inner)}`).join(`,${eol}`);
        return `{${eol}${body}${eol}${pad}}`;
      }
      case 'array': {
        if (node.items.length === 0) return '[]';
        const inner = pad + indent;
        const body = node.items.map((item) => inner + write(item, inner)).join(`,${eol}`);
        return `[${eol}${body}${eol}${pad}]`;
      }
      case 'string':
        return JSON.stringify(node.value);
      case 'number':
        return node.raw;
      case 'boolean':
        return node.value ? 'true' : 'false';
      case 'null':
        return 'null';
    }
  };
  return write(root, '') + (format.finalNewline ? eol : '');
}

export function isContainer(node: JsonNode): node is JsonContainer {
  return node.kind === 'object' || node.kind === 'array';
}

export function childCount(node: JsonContainer): number {
  return node.kind === 'object' ? node.entries.length : node.items.length;
}

export function newNode(kind: JsonKind): JsonNode {
  switch (kind) {
    case 'object':
      return { id: newId(), kind, entries: [] };
    case 'array':
      return { id: newId(), kind, items: [] };
    case 'string':
      return { id: newId(), kind, value: '' };
    case 'number':
      return { id: newId(), kind, raw: '0' };
    case 'boolean':
      return { id: newId(), kind, value: false };
    case 'null':
      return { id: newId(), kind };
  }
}

// 深拷贝并换上新 ID
export function cloneNode(node: JsonNode): JsonNode {
  switch (node.kind) {
    case 'object':
      return { id: newId(), kind: 'object', entries: node.entries.map((e) => ({ key: e.key, value: cloneNode(e.value) })) };
    case 'array':
      return { id: newId(), kind: 'array', items: node.items.map(cloneNode) };
    default:
      return { ...node, id: newId() };
  }
}

// 路径是每一层在 entries / items 里的下标
export type JsonPath = readonly number[];

export function updateAt(root: JsonNode, path: JsonPath, fn: (node: JsonNode) => JsonNode): JsonNode {
  if (path.length === 0) return fn(root);
  const [index, ...rest] = path;
  if (root.kind === 'object') {
    const entries = root.entries.slice();
    entries[index] = { ...entries[index], value: updateAt(entries[index].value, rest, fn) };
    return { ...root, entries };
  }
  if (root.kind === 'array') {
    const items = root.items.slice();
    items[index] = updateAt(items[index], rest, fn);
    return { ...root, items };
  }
  return root;
}

export function hasKey(node: JsonContainer, key: string, exceptIndex = -1): boolean {
  return node.kind === 'object' && node.entries.some((e, i) => i !== exceptIndex && e.key === key);
}

// 复制一个字段时给它起不重名的键：结尾是数字的往上加（A001 → A002，保留补零），否则加 _2、_3
export function nextFreeKey(node: JsonContainer, base: string): string {
  const m = /^(.*?)(\d+)$/.exec(base);
  if (m) {
    const [, prefix, digits] = m;
    for (let n = Number(digits) + 1; ; n++) {
      const key = prefix + String(n).padStart(digits.length, '0');
      if (!hasKey(node, key)) return key;
    }
  }
  for (let n = 2; ; n++) {
    const key = `${base}_${n}`;
    if (!hasKey(node, key)) return key;
  }
}

// 在容器的 index 后面插入一份拷贝
export function duplicateChild(node: JsonContainer, index: number): JsonContainer {
  if (node.kind === 'object') {
    const src = node.entries[index];
    const entries = node.entries.slice();
    entries.splice(index + 1, 0, { key: nextFreeKey(node, src.key), value: cloneNode(src.value) });
    return { ...node, entries };
  }
  const items = node.items.slice();
  items.splice(index + 1, 0, cloneNode(node.items[index]));
  return { ...node, items };
}

export function removeChild(node: JsonContainer, index: number): JsonContainer {
  if (node.kind === 'object') return { ...node, entries: node.entries.filter((_, i) => i !== index) };
  return { ...node, items: node.items.filter((_, i) => i !== index) };
}

export interface JsonFilter {
  // 要显示的节点
  shown: Set<number>;
  // 子孙里有匹配、搜索时要展开的容器
  open: Set<number>;
}

// 按字段名和值搜索（不区分大小写）。字段名匹配时整棵子树都显示，值匹配时只显示到它的路径
export function filterJsonTree(root: JsonNode, query: string): JsonFilter | null {
  const q = query.trim().toLowerCase();
  if (!q) return null;
  const shown = new Set<number>();
  const open = new Set<number>();

  const showAll = (node: JsonNode) => {
    shown.add(node.id);
    if (node.kind === 'object') node.entries.forEach((e) => showAll(e.value));
    else if (node.kind === 'array') node.items.forEach(showAll);
  };

  const visit = (node: JsonNode, key: string | null): boolean => {
    if (key !== null && key.toLowerCase().includes(q)) {
      showAll(node);
      return true;
    }
    let hit = false;
    if (node.kind === 'object') {
      for (const e of node.entries) if (visit(e.value, e.key)) hit = true;
    } else if (node.kind === 'array') {
      for (const item of node.items) if (visit(item, null)) hit = true;
    } else {
      const text = node.kind === 'string' ? node.value : node.kind === 'number' ? node.raw : String(node.kind === 'boolean' ? node.value : 'null');
      hit = text.toLowerCase().includes(q);
    }
    if (hit) {
      shown.add(node.id);
      if (isContainer(node)) open.add(node.id);
    }
    return hit;
  };

  visit(root, null);
  return { shown, open };
}
