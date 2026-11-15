// 服务器配置的导出、导入和复制。格式与后端 models.ServerConfig 一致：
// 不含端口和管理员密码，导入或复制时目标服务器保留自己的名称、端口、人数和密码
import axios from 'axios';
import { keepIniValues, replaceIniValues } from '@/lib/ini';
import type { Server } from '@/stores/servers';

export const SERVER_CONFIG_FORMAT = 'ark-server-commander/server-config';
export const SERVER_CONFIG_VERSION = 1;

export interface ServerConfig {
  format: string;
  version: number;
  exported_at?: string;
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

// 默认选中的部分：游戏规则相关的内容；地图、集群和选项通常每台服务器不同
export const DEFAULT_CONFIG_PARTS: ConfigPart[] = ['config_files', 'server_args', 'mods'];

// 每台服务器自己的键，导入配置文件时保留原值（与后端 identityIniKeys 一致）
const IDENTITY_INI_KEYS = ['SessionName', 'ServerAdminPassword', 'Port', 'QueryPort', 'RCONPort', 'MaxPlayers'];

export class InvalidConfigFileError extends Error {}

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

/** 把配置里选中的部分写进表单（导入到新建或编辑页），规则与后端 applyConfig 相同 */
export function applyConfigToForm(form: Partial<Server>, cfg: ServerConfig, parts: ConfigPart[]): Partial<Server> {
  const has = (p: ConfigPart) => parts.includes(p);
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
