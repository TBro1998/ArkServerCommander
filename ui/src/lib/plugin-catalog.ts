// 插件目录（arkfgf.com 整理的 ark-server-api.com 上的插件，ASE 的 ArkApi 插件和 ASA 的 AsaApi 插件各一份）和已安装插件的对应关系

export interface CatalogPluginInfo {
  name: string;
  full_name?: string;
  version?: string;
  min_api_version?: string;
  dependencies?: string[];
}

export interface CatalogPlugin {
  id: number;
  title: string;
  tagline?: string;
  author?: string;
  category: 'official' | 'free' | 'paid' | string;
  version?: string;
  downloads: number;
  rating?: number;
  updated_at?: string;
  page_url: string;
  download_url?: string;
  file_name?: string;
  file_size?: number;
  plugins?: CatalogPluginInfo[];
}

export interface PluginCatalog {
  game: 'ase' | 'asa';
  source: string;
  source_url: string;
  generated_at: string;
  plugins: CatalogPlugin[];
}

export interface PluginSource {
  source: string;
  resource_id: number;
  version: string;
  installed_at: string;
}

// 已安装插件里用于对应目录的字段
export interface InstalledPluginRef {
  name: string;
  version?: string;
  source?: PluginSource;
}

// 按点分隔的数字比较版本号，与后端 compareVersions 一致
export function compareVersions(a: string, b: string): number {
  const pa = a.split('.');
  const pb = b.split('.');
  for (let i = 0; i < Math.max(pa.length, pb.length); i++) {
    const x = pa[i] ?? '';
    const y = pb[i] ?? '';
    const nx = /^\d+$/.test(x) ? Number(x) : NaN;
    const ny = /^\d+$/.test(y) ? Number(y) : NaN;
    if (!Number.isNaN(nx) && !Number.isNaN(ny)) {
      if (nx !== ny) return nx < ny ? -1 : 1;
    } else if (x !== y) {
      return x < y ? -1 : 1;
    }
  }
  return 0;
}

export interface CatalogIndex {
  byId: Map<number, CatalogPlugin>;
  // 插件目录名（小写）-> 包含它的资源，只算能一键安装的
  byPluginName: Map<string, CatalogPlugin>;
}

export function indexCatalog(catalog: PluginCatalog | null): CatalogIndex {
  const byId = new Map<number, CatalogPlugin>();
  const byPluginName = new Map<string, CatalogPlugin>();
  for (const entry of catalog?.plugins ?? []) {
    byId.set(entry.id, entry);
    if (!entry.download_url) continue;
    for (const p of entry.plugins ?? []) {
      const key = p.name.toLowerCase();
      // 同名插件出现在多个资源里时，官方的和下载多的优先（目录已经按这个顺序排好）
      if (!byPluginName.has(key)) byPluginName.set(key, entry);
    }
  }
  return { byId, byPluginName };
}

// 已安装插件对应的目录资源：从目录安装的按资源 ID，其他的按插件目录名
export function catalogEntryFor(plugin: InstalledPluginRef, index: CatalogIndex): CatalogPlugin | undefined {
  if (plugin.source) return index.byId.get(plugin.source.resource_id);
  return index.byPluginName.get(plugin.name.toLowerCase());
}

// 有更新时返回目录里的资源版本：从目录安装的比较资源版本，其他的比较 PluginInfo.json 的版本
export function availableUpdate(plugin: InstalledPluginRef, entry: CatalogPlugin | undefined): string | null {
  if (!entry?.download_url || !entry.version) return null;
  if (plugin.source) {
    return compareVersions(entry.version, plugin.source.version) > 0 ? entry.version : null;
  }
  const info = entry.plugins?.find((p) => p.name.toLowerCase() === plugin.name.toLowerCase());
  if (info?.version && plugin.version && compareVersions(info.version, plugin.version) > 0) return entry.version;
  return null;
}
