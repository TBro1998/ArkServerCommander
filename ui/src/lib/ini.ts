// INI 配置文件的读取和就地修改。
//
// 设置界面只认识内置字典里的键，修改时只改动对应的那一行：命中的键替换值，
// 没有的键追加到所在 section 末尾，其余内容（未知键、自定义 section、注释、
// 重复的 ConfigOverride* 键）原样保留。键和 section 名都不区分大小写。

export type IniValue = string | number | boolean

interface IniBlock {
  // section 名（含方括号）；文件开头的无 section 内容为 null
  header: string | null
  lines: string[]
}

function parseBlocks(content: string): IniBlock[] {
  const blocks: IniBlock[] = [{ header: null, lines: [] }]
  if (content === '') return blocks
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
  // 去掉开头和末尾多余空行，保留一个换行
  while (out.length > 0 && out[out.length - 1].trim() === '') out.pop()
  while (out.length > 0 && out[0].trim() === '') out.shift()
  return out.length === 0 ? '' : out.join('\n') + '\n'
}

export function formatIniValue(value: IniValue): string {
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

const lower = (s: string) => s.toLowerCase()

// IniIndex：section（小写，含方括号）→ 键（小写）→ 值。同一 section 里重复的键取第一个
export type IniIndex = Map<string, Map<string, string>>

export function indexIni(content: string): IniIndex {
  const index: IniIndex = new Map()
  for (const block of parseBlocks(content)) {
    if (block.header === null) continue
    const name = lower(block.header)
    const entries = index.get(name) ?? new Map<string, string>()
    for (const line of block.lines) {
      const kv = splitKeyValue(line)
      if (kv && !entries.has(lower(kv.key))) entries.set(lower(kv.key), kv.value)
    }
    index.set(name, entries)
  }
  return index
}

export function getIniValue(index: IniIndex, section: string, key: string): string | undefined {
  return index.get(lower(section))?.get(lower(key))
}

/**
 * 在 section 里设置一个键，返回新的文件内容。value 为 undefined 时删除这个键（恢复游戏默认值）。
 * 只改动目标 section：其它 section 里的同名键不受影响（游戏只读取所属 section 里的键）。
 */
export function setIniValue(content: string, section: string, key: string, value: IniValue | undefined): string {
  const blocks = parseBlocks(content)
  let target = blocks.find((b) => b.header !== null && lower(b.header) === lower(section))

  if (value === undefined) {
    if (!target) return content
    const before = target.lines.length
    target.lines = target.lines.filter((line) => lower(splitKeyValue(line)?.key ?? '') !== lower(key))
    return target.lines.length === before ? content : serializeBlocks(blocks)
  }

  const formatted = formatIniValue(value)
  if (!target) {
    // 新 section 前空一行
    const last = blocks[blocks.length - 1]
    if (last.lines.length > 0 && last.lines[last.lines.length - 1].trim() !== '') last.lines.push('')
    target = { header: section, lines: [] }
    blocks.push(target)
  }

  let replaced = false
  target.lines = target.lines.map((line) => {
    if (replaced) return line
    const kv = splitKeyValue(line)
    if (!kv || lower(kv.key) !== lower(key)) return line
    replaced = true
    return `${kv.indent}${kv.key}=${formatted}`
  })
  if (!replaced) {
    // 插在 section 最后一个非空行之后，保留 section 之间的空行
    let at = target.lines.length
    while (at > 0 && target.lines[at - 1].trim() === '') at--
    target.lines.splice(at, 0, `${key}=${formatted}`)
  }
  return serializeBlocks(blocks)
}

/** 只替换已经存在的键（任何 section 里的都替换），不新增。用于同步分散在多个 section 里的名称、人数 */
export function replaceIniValues(content: string, key: string, value: IniValue): string {
  let changed = false
  const blocks = parseBlocks(content)
  for (const block of blocks) {
    block.lines = block.lines.map((line) => {
      const kv = splitKeyValue(line)
      if (!kv || lower(kv.key) !== lower(key)) return line
      changed = true
      return `${kv.indent}${kv.key}=${formatIniValue(value)}`
    })
  }
  return changed ? serializeBlocks(blocks) : content
}
