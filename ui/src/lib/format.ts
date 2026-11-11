// 字节数转为易读的大小，例如 1536 -> "1.5 KB"
export function formatBytes(bytes: number | undefined | null, digits = 1): string {
  if (!bytes || bytes <= 0) return '0 B';
  const units = ['B', 'KB', 'MB', 'GB', 'TB'];
  const i = Math.min(Math.floor(Math.log(bytes) / Math.log(1024)), units.length - 1);
  const value = bytes / Math.pow(1024, i);
  return `${value.toFixed(i === 0 ? 0 : digits)} ${units[i]}`;
}

export function formatPercent(value: number | undefined | null, digits = 1): string {
  return `${(value ?? 0).toFixed(digits)}%`;
}

// ISO 时间转为本地时间字符串
export function formatDateTime(value: string | undefined | null, locale: string): string {
  if (!value) return '-';
  const d = new Date(value);
  if (Number.isNaN(d.getTime())) return value;
  return d.toLocaleString(locale === 'zh' ? 'zh-CN' : 'en-US', { hour12: false });
}

// 生成随机密码（字母和数字，避免 ini 和命令行里容易出问题的符号）
export function randomPassword(length = 16): string {
  const chars = 'ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz23456789';
  const values = new Uint32Array(length);
  crypto.getRandomValues(values);
  return Array.from(values, (v) => chars[v % chars.length]).join('');
}

// 复制文本：优先用 Clipboard API，非 HTTPS 环境下回退到 execCommand
export async function copyText(text: string): Promise<boolean> {
  try {
    if (navigator.clipboard && window.isSecureContext) {
      await navigator.clipboard.writeText(text);
      return true;
    }
  } catch {
    // 回退
  }
  try {
    const el = document.createElement('textarea');
    el.value = text;
    el.style.position = 'fixed';
    el.style.opacity = '0';
    document.body.appendChild(el);
    el.select();
    const ok = document.execCommand('copy');
    document.body.removeChild(el);
    return ok;
  } catch {
    return false;
  }
}
