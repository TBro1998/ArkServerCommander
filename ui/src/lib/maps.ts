import type { Game } from '@/lib/games';

// 官方地图（地图选择器里的顺序）。ASA 的地图名都带 _WP 后缀
export const ASE_MAPS = [
  'TheIsland',
  'TheCenter',
  'ScorchedEarth_P',
  'Aberration_P',
  'Extinction',
  'Ragnarok',
  'Valguero_P',
  'Genesis',
  'Genesis2',
  'CrystalIsles',
  'LostIsland',
  'Fjordur',
];

export const ASA_MAPS = [
  'TheIsland_WP',
  'TheCenter_WP',
  'ScorchedEarth_WP',
  'Aberration_WP',
  'Extinction_WP',
  'Ragnarok_WP',
  'Valguero_WP',
  'LostColony_WP',
  'Astraeos_WP',
  'Genesis_WP',
  // Club ARK：需要 CurseForge Mod 1005639
  'BobsMissions_WP',
];

export function officialMaps(game: Game): string[] {
  return game === 'asa' ? ASA_MAPS : ASE_MAPS;
}

/** 新建服务器的默认地图，与后端 DefaultMapFor 一致 */
export function defaultMap(game: Game): string {
  return game === 'asa' ? 'TheIsland_WP' : 'TheIsland';
}

// 每张地图一个固定的颜色，取自地图的地貌：服务器卡片和详情页顶部的色带用它区分"这是哪台服务器"。
// ASA 的同名地图用同一个颜色
const MAP_COLORS: Record<string, string> = {
  TheIsland: '#3d7a4c', // 丛林
  TheCenter: '#5a7d3a', // 浮岛草地
  ScorchedEarth_P: '#a86d27', // 沙漠
  Aberration_P: '#6a4c9c', // 地下荧光
  Extinction: '#8a4636', // 废墟
  Ragnarok: '#2f6c78', // 北欧峡湾
  Valguero_P: '#4c7658', // 河谷
  Genesis: '#23758a', // 模拟器
  Genesis2: '#3d5a9c', // 飞船
  CrystalIsles: '#9a4882', // 水晶
  LostIsland: '#4b7a6a', // 失落岛
  Fjordur: '#3f78a0', // 冰川
  TheIsland_WP: '#3d7a4c',
  TheCenter_WP: '#5a7d3a',
  ScorchedEarth_WP: '#a86d27',
  Aberration_WP: '#6a4c9c',
  Extinction_WP: '#8a4636',
  Ragnarok_WP: '#2f6c78',
  Valguero_WP: '#4c7658',
  LostColony_WP: '#7a6a3d', // 殖民地遗迹
  Astraeos_WP: '#4a5f9a', // 神话群岛
  Genesis_WP: '#23758a',
  BobsMissions_WP: '#b0703a', // 俱乐部
};

const FALLBACK = '#56645b';

export function mapColor(map: string | undefined): string {
  return (map && MAP_COLORS[map]) || FALLBACK;
}

// 没有翻译的自定义地图：去掉 _P / _WP 后缀，并在驼峰处断开（ScorchedEarth_P → Scorched Earth）
export function prettyMapName(map: string): string {
  return map.replace(/_W?P$/, '').replace(/([a-z])([A-Z0-9])/g, '$1 $2');
}
