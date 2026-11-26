// 游戏：ARK: Survival Evolved（ASE）和 ARK: Survival Ascended（ASA），与后端 utils/game.go 一致。
// 服务器创建时选定，之后不能修改：两个游戏的地图、Mod（ASE 是 Steam 创意工坊，ASA 是 CurseForge）、
// 启动参数和插件框架（ArkApi / AsaApi）都不通用
export type Game = 'ase' | 'asa';

export const GAMES: Game[] = ['ase', 'asa'];

export const GAME_NAMES: Record<Game, string> = {
  ase: 'ARK: Survival Evolved',
  asa: 'ARK: Survival Ascended',
};

/** 服务器或配置的游戏：空值（旧版服务器、旧版导出的配置）视为 ASE */
export function gameOf(value: { game?: string } | null | undefined): Game {
  return value?.game === 'asa' ? 'asa' : 'ase';
}

/** 插件框架的名称：ASE 是 ArkApi，ASA 是 AsaApi（服务器字段 arkapi 对 ASA 表示启用 AsaApi） */
export function pluginApiName(game: Game): 'ArkApi' | 'AsaApi' {
  return game === 'asa' ? 'AsaApi' : 'ArkApi';
}
