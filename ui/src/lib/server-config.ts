// 服务器配置的导出、导入和复制。格式与后端 models.ServerConfig 一致：
// 不含端口和管理员密码，导入或复制时目标服务器保留自己的名称、端口、人数和密码
import axios from 'axios';
import { getIniValue, indexIni, keepIniValues, replaceIniValues, setIniValue } from '@/lib/ini';
import type { Server } from '@/stores/servers';
import { gameOf, type Game } from '@/lib/games';

export const SERVER_CONFIG_FORMAT = 'ark-server-commander/server-config';
export const SERVER_CONFIG_VERSION = 1;

export interface ServerConfig {
  format: string;
  version: number;
  exported_at?: string;
  // 来源服务器的游戏，旧版导出的没有，视为 ASE
  game?: Game;
  // 来源服务器名称，只用于显示
  session_name?: string;
  map?: string;
  max_players?: number;
  cluster_id?: string;
  game_mod_ids?: string;
  arkapi?: boolean;
  auto_restart?: boolean;
  auto_update?: boolean;
  game_user_settings?: string;
  game_ini?: string;
  server_args?: Server['server_args'];
}

// 可以选择导入或复制的部分，顺序即界面上的顺序
export const CONFIG_PARTS = ['config_files', 'server_args', 'mods', 'map', 'cluster', 'options'] as const;
export type ConfigPart = (typeof CONFIG_PARTS)[number];

// 跨游戏（ASE 和 ASA 之间）导入或复制时不通用、不会导入的部分（与后端 applyConfig 一致）：
// 地图名、Mod（创意工坊 / CurseForge）和启动参数都不一样
export const CROSS_GAME_PARTS: ConfigPart[] = ['map', 'mods', 'server_args'];

// 默认选中的部分：游戏规则相关的内容；地图、集群和选项通常每台服务器不同
export const DEFAULT_CONFIG_PARTS: ConfigPart[] = ['config_files', 'server_args', 'mods'];

// 每台服务器自己的键，导入配置文件时保留原值（与后端 identityIniKeys 一致）
const IDENTITY_INI_KEYS = ['SessionName', 'ServerAdminPassword', 'Port', 'QueryPort', 'RCONPort', 'MaxPlayers'];

// invalid：不是本程序导出的 JSON 或版本不支持；iniUnknown：认不出是哪个 ini；
// iniDuplicate：同一种 ini 选了两份；mixed：JSON 和 ini 混选或选了多个 JSON；empty：ini 文件是空的
export type InvalidConfigFileReason = 'invalid' | 'iniUnknown' | 'iniDuplicate' | 'mixed' | 'empty';

export class InvalidConfigFileError extends Error {
  reason: InvalidConfigFileReason;
  fileName?: string;
  constructor(message: string, reason: InvalidConfigFileReason = 'invalid', fileName?: string) {
    super(message);
    this.reason = reason;
    this.fileName = fileName;
  }
}

export type IniFileKind = 'game_user_settings' | 'game_ini';

// 各自独有的 section（小写），用来认出是哪个文件
const GUS_SECTIONS = [
  '[serversettings]',
  '[sessionsettings]',
  '[/script/shootergame.shootergameusersettings]',
  '[/script/engine.gamesession]',
  '[messageoftheday]',
  '[scalabilitygroups]',
];
const GAME_INI_SECTIONS = ['[/script/shootergame.shootergamemode]'];

/** 判断 ini 是 GameUserSettings.ini 还是 Game.ini：先看内容里的 section，认不出再看文件名 */
export function detectIniKind(fileName: string, text: string): IniFileKind | null {
  const sections = [...indexIni(text).keys()];
  const isGus = sections.some((s) => GUS_SECTIONS.includes(s));
  const isGame = sections.some((s) => GAME_INI_SECTIONS.includes(s));
  if (isGus !== isGame) return isGus ? 'game_user_settings' : 'game_ini';
  const name = fileName.toLowerCase();
  if (name.includes('gameusersettings')) return 'game_user_settings';
  if (/^game\b/.test(name)) return 'game_ini';
  return null;
}

/** 按 BOM 选择编码读取文本：服务端写入含中文的 ini 时会用 UTF-16 */
async function readTextFile(file: File): Promise<string> {
  const buf = new Uint8Array(await file.arrayBuffer());
  let encoding = 'utf-8';
  if (buf[0] === 0xff && buf[1] === 0xfe) encoding = 'utf-16le';
  else if (buf[0] === 0xfe && buf[1] === 0xff) encoding = 'utf-16be';
  return new TextDecoder(encoding).decode(buf);
}

const isIniFile = (file: File) => file.name.toLowerCase().endsWith('.ini');

export interface ImportedConfigFile {
  config: ServerConfig;
  // ini：只含配置文件，只能导入「游戏设置」部分（以及 ActiveMods 里的 Mod）；iniKinds 为读到的文件
  kind: 'json' | 'ini';
  iniKinds: IniFileKind[];
}

const SERVER_SETTINGS_SECTION = '[ServerSettings]';

/**
 * 取出 GameUserSettings.ini 里 [ServerSettings] 的 ActiveMods，并从文本里删掉这一行。
 * 管理端用 Mod 列表（启动参数 ?GameModIds=）管理 Mod，入口脚本也只下载列表里的 Mod；
 * 留着 ActiveMods 时，Mod 列表为空的服务器会去加载没下载的 Mod
 */
export function extractActiveMods(gus: string): { gus: string; modIds: string } {
  const value = getIniValue(indexIni(gus), SERVER_SETTINGS_SECTION, 'ActiveMods');
  if (value === undefined) return { gus, modIds: '' };
  const ids = [...new Set(value.split(',').map((id) => id.trim()).filter((id) => /^\d+$/.test(id)))];
  return { gus: setIniValue(gus, SERVER_SETTINGS_SECTION, 'ActiveMods', undefined), modIds: ids.join(',') };
}

/**
 * 读取导入时选择的文件：一个导出的 JSON，或者 GameUserSettings.ini / Game.ini（可同时选两个）。
 * ini 被包装成只含配置文件的 ServerConfig，导入时只用 config_files 部分，没选的那个文件保持不变；
 * GameUserSettings.ini 里的 ActiveMods 去掉，Mod ID 放进 game_mod_ids，由用户决定是否同时导入 mods 部分
 */
export async function readImportFiles(files: File[]): Promise<ImportedConfigFile> {
  const inis = files.filter(isIniFile);
  if (inis.length === 0) {
    if (files.length !== 1) throw new InvalidConfigFileError('multiple json', 'mixed');
    return { config: parseServerConfig(await readTextFile(files[0])), kind: 'json', iniKinds: [] };
  }
  if (inis.length !== files.length) throw new InvalidConfigFileError('mixed', 'mixed');

  const config: ServerConfig = { format: SERVER_CONFIG_FORMAT, version: SERVER_CONFIG_VERSION };
  for (const file of inis) {
    const text = await readTextFile(file);
    if (text.trim() === '') throw new InvalidConfigFileError('empty', 'empty', file.name);
    const kind = detectIniKind(file.name, text);
    if (!kind) throw new InvalidConfigFileError('unknown ini', 'iniUnknown', file.name);
    if (config[kind] !== undefined) throw new InvalidConfigFileError('duplicate ini', 'iniDuplicate', file.name);
    if (kind === 'game_user_settings') {
      const { gus, modIds } = extractActiveMods(text);
      config.game_user_settings = gus;
      if (modIds) config.game_mod_ids = modIds;
    } else {
      config.game_ini = text;
    }
  }
  const iniKinds = (['game_user_settings', 'game_ini'] as const).filter((k) => config[k] !== undefined);
  return { config, kind: 'ini', iniKinds };
}

/** 解析导入的文件，不是本程序导出的配置时抛出 InvalidConfigFileError */
export function parseServerConfig(text: string): ServerConfig {
  let data: unknown;
  try {
    data = JSON.parse(text);
  } catch {
    throw new InvalidConfigFileError('invalid json');
  }
  const cfg = data as ServerConfig;
  if (!cfg || typeof cfg !== 'object' || cfg.format !== SERVER_CONFIG_FORMAT) throw new InvalidConfigFileError('format');
  if (typeof cfg.version !== 'number' || cfg.version < 1 || cfg.version > SERVER_CONFIG_VERSION) throw new InvalidConfigFileError('version');
  return cfg;
}

export async function fetchServerConfig(serverId: number): Promise<ServerConfig> {
  const res = await axios.get(`/api/servers/${serverId}/config/export`);
  return res.data.data;
}

/** 导出配置并下载为 JSON 文件 */
export async function downloadServerConfig(server: Pick<Server, 'id' | 'session_name'>): Promise<void> {
  const cfg = await fetchServerConfig(server.id);
  const blob = new Blob([JSON.stringify(cfg, null, 2) + '\n'], { type: 'application/json;charset=utf-8' });
  const url = URL.createObjectURL(blob);
  const a = document.createElement('a');
  a.href = url;
  // 文件名里去掉系统不允许的字符
  const name = server.session_name.replace(/[\\/:*?"<>|\s]+/g, '_').replace(/^_+|_+$/g, '') || `server-${server.id}`;
  a.download = `${name}.ark-config.json`;
  a.click();
  URL.revokeObjectURL(url);
}

/** 把配置里选中的部分写进表单（导入到新建或编辑页），规则与后端 applyConfig 相同：跨游戏时跳过地图、Mod 和启动参数 */
export function applyConfigToForm(form: Partial<Server>, cfg: ServerConfig, parts: ConfigPart[]): Partial<Server> {
  const crossGame = gameOf(cfg) !== gameOf(form);
  const has = (p: ConfigPart) => parts.includes(p) && !(crossGame && CROSS_GAME_PARTS.includes(p));
  const next: Partial<Server> = { ...form };

  if (has('options')) {
    if (cfg.max_players) next.max_players = cfg.max_players;
    if (cfg.arkapi !== undefined) next.arkapi = cfg.arkapi;
    if (cfg.auto_restart !== undefined) next.auto_restart = cfg.auto_restart;
    if (cfg.auto_update !== undefined) next.auto_update = cfg.auto_update;
  }
  if (has('map') && cfg.map) next.map = cfg.map;
  if (has('mods')) next.game_mod_ids = cfg.game_mod_ids ?? '';
  if (has('cluster')) next.cluster_id = cfg.cluster_id ?? '';
  if (has('server_args') && cfg.server_args) {
    next.server_args = {
      query_params: cfg.server_args.query_params ?? {},
      command_line_args: cfg.server_args.command_line_args ?? {},
      custom_args: cfg.server_args.custom_args ?? [],
    };
  }
  if (has('config_files')) {
    if (cfg.game_user_settings) {
      let gus = keepIniValues(cfg.game_user_settings, form.game_user_settings ?? '', IDENTITY_INI_KEYS);
      if (has('options') && cfg.max_players) gus = replaceIniValues(gus, 'MaxPlayers', cfg.max_players);
      next.game_user_settings = gus;
    }
    if (cfg.game_ini) next.game_ini = cfg.game_ini;
  }
  return next;
}
