// INI 配置文件合并。
//
// 可视化编辑器只认识内置参数字典里的键。直接按字典重新生成整个文件会把用户手写的
// 配置整体丢掉（未知键、自定义 section、注释、重复的 ConfigOverride* 键）。
// 这里改为在原文本上就地修改：命中的键只替换值，其余行原样保留，
// 字典里有但文件里没有的键追加到对应 section 末尾。

export type IniValue = string | number | boolean

export interface IniMergeOptions {
  // 该键应当写入的 section（文件里没有这个 section 时新建）
  sectionForKey: (key: string) => string
  // 文件为空时的初始内容（通常是默认模板）
  initialContent?: string
}

interface IniBlock {
  // section 名（含方括号）；文件开头的无 section 内容为 null
  header: string | null
  lines: string[]
}

function parseBlocks(content: string): IniBlock[] {
  const blocks: IniBlock[] = [{ header: null, lines: [] }]
  for (const line of content.split(/\r?\n/)) {
    const trimmed = line.trim()
    if (trimmed.startsWith('[') && trimmed.endsWith(']')) {
      blocks.push({ header: trimmed, lines: [] })
      continue
    }
    blocks[blocks.length - 1].lines.push(line)
  }
  return blocks
}

function serializeBlocks(blocks: IniBlock[]): string {
  const out: string[] = []
  for (const block of blocks) {
    if (block.header !== null) out.push(block.header)
    out.push(...block.lines)
  }
  // 去掉末尾多余空行，保留一个换行
  while (out.length > 0 && out[out.length - 1].trim() === '') out.pop()
  return out.join('\n') + '\n'
}

function formatValue(value: IniValue): string {
  if (typeof value === 'boolean') return value ? 'True' : 'False'
  return String(value)
}

// splitKeyValue 拆分 "Key=Value"，保留缩进；注释行和非键值行返回 null
function splitKeyValue(line: string): { indent: string; key: string; value: string } | null {
  const trimmed = line.trimStart()
  if (trimmed === '' || trimmed.startsWith(';') || trimmed.startsWith('#')) return null
  const eq = trimmed.indexOf('=')
  if (eq <= 0) return null
  const key = trimmed.slice(0, eq).trim()
  if (key === '') return null
  return {
    indent: line.slice(0, line.length - trimmed.length),
    key,
    value: trimmed.slice(eq + 1).trim(),
  }
}

function sameSection(a: string, b: string): boolean {
  return a.toLowerCase() === b.toLowerCase()
}

/**
 * 把 updates 合并进 content，返回新的文件内容。
 * 已存在的键就地替换（同一 section 或不同 section 中的同名键都会更新），
 * 不存在的键追加到 sectionForKey 指定的 section；其它行保持不变。
 */
export function mergeIni(content: string, updates: Record<string, IniValue>, options: IniMergeOptions): string {
  const base = content.trim() === '' ? (options.initialContent ?? '') : content
  const blocks = parseBlocks(base)

  // INI 键不区分大小写，按小写匹配
  const pending = new Map<string, { key: string; value: IniValue }>()
  for (const [key, value] of Object.entries(updates)) {
    if (value === undefined || value === null) continue
    pending.set(key.toLowerCase(), { key, value })
  }

  for (const block of blocks) {
    block.lines = block.lines.map((line) => {
      const parsed = splitKeyValue(line)
      if (!parsed) return line
      const update = pending.get(parsed.key.toLowerCase())
      if (!update) return line
      pending.delete(parsed.key.toLowerCase())
      return `${parsed.indent}${parsed.key}=${formatValue(update.value)}`
    })
  }

  // 剩下的键是文件里还没有的，按目标 section 分组后追加
  const additions = new Map<string, string[]>()
  for (const { key, value } of pending.values()) {
    const section = options.sectionForKey(key)
    const lines = additions.get(section) ?? []
    lines.push(`${key}=${formatValue(value)}`)
    additions.set(section, lines)
  }

  for (const [section, lines] of additions) {
    let target = blocks.find((b) => b.header !== null && sameSection(b.header, section))
    if (!target) {
      target = { header: section, lines: [] }
      blocks.push(target)
    }
    target.lines.push(...lines)
  }

  return serializeBlocks(blocks)
}
