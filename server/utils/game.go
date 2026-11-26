package utils

// 游戏：ARK: Survival Evolved（ASE）和 ARK: Survival Ascended（ASA）。服务器创建时选定，之后不能修改：
// 两个游戏的服务端文件、存档、配置、Mod（ASE 是 Steam 创意工坊，ASA 是 CurseForge）和插件都不通用
const (
	GameASE = "ase"
	GameASA = "asa"
)

// Games 可选的游戏
var Games = []string{GameASE, GameASA}

// NormalizeGame 空值（加这个字段之前创建的服务器、旧版请求）视为 ASE
func NormalizeGame(game string) string {
	if game == "" {
		return GameASE
	}
	return game
}

// ValidGame 是否是可选的游戏（空值视为 ASE）
func ValidGame(game string) bool {
	game = NormalizeGame(game)
	return game == GameASE || game == GameASA
}
