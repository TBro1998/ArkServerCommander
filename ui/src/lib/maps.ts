// 每张地图一个固定的颜色，取自地图的地貌：服务器卡片和详情页顶部的色带用它区分"这是哪台服务器"
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
};

const FALLBACK = '#56645b';

export function mapColor(map: string | undefined): string {
  return (map && MAP_COLORS[map]) || FALLBACK;
}

// 没有翻译的自定义地图：去掉 _P 后缀，并在驼峰处断开（ScorchedEarth_P → Scorched Earth）
export function prettyMapName(map: string): string {
  return map.replace(/_P$/, '').replace(/([a-z])([A-Z0-9])/g, '$1 $2');
}
