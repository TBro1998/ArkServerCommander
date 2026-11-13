package utils

import (
	"fmt"
	"strings"
)

const (
	// 配置文件名常量
	GameUserSettingsFileName = "GameUserSettings.ini"
	GameIniFileName          = "Game.ini"

	// 配置文件目录 - 适配新的ASE服务器镜像路径
	ConfigDirectory = "Config/WindowsServer"
)

// GetDefaultGameUserSettings 获取默认的GameUserSettings.ini配置。
// 倍率、难度和昼夜速度等服务器选项属于 GameUserSettings.ini 的 [ServerSettings]，写在 Game.ini 里不会生效
// （见 docs/ASE_GameUserSetting.md）；ARK 停服时会重写这个文件，所以这里不写注释
func GetDefaultGameUserSettings(serverName, mapName string, maxPlayers int) string {
	return fmt.Sprintf(`[ServerSettings]
SessionName=%s
ServerPassword=
MaxPlayers=%d
DifficultyOffset=1.0
OverrideOfficialDifficulty=5.0
XPMultiplier=1.0
TamingSpeedMultiplier=1.0
HarvestAmountMultiplier=1.0
HarvestHealthMultiplier=1.0
ResourcesRespawnPeriodMultiplier=1.0
DinoCountMultiplier=1.0
DinoCharacterFoodDrainMultiplier=1.0
DinoCharacterStaminaDrainMultiplier=1.0
DinoCharacterHealthRecoveryMultiplier=1.0
PlayerCharacterWaterDrainMultiplier=1.0
PlayerCharacterFoodDrainMultiplier=1.0
PlayerCharacterStaminaDrainMultiplier=1.0
PlayerCharacterHealthRecoveryMultiplier=1.0
DayCycleSpeedScale=1.0
NightTimeSpeedScale=1.0
StructureResistanceMultiplier=1.0
StructureDamageMultiplier=1.0
PvEStructureDecayPeriodMultiplier=1.0
AllowFlyerCarryPvE=True
DisableStructureDecayPvE=False
EnablePvPGamma=False

[SessionSettings]
SessionName=%s

[MessageOfTheDay]
Message=欢迎来到 %s ARK 服务器！

[/Script/Engine.GameSession]
MaxPlayers=%d
`, serverName, maxPlayers, serverName, serverName, maxPlayers)
}

// GetDefaultGameIni 获取默认的Game.ini配置。
// 只放 Game.ini 的选项，全部在 [/script/shootergame.shootergamemode] 下（见 docs/ASE_Game.ini.md）
func GetDefaultGameIni() string {
	return `[/script/shootergame.shootergamemode]
bUseSingleplayerSettings=False
bDisableStructurePlacementCollision=False
bAllowUnlimitedRespecs=True
bAllowPlatformSaddleMultiFloors=True
bPassiveDefensesDamageRiderlessDinos=True
MaxNumberOfPlayersInTribe=0
StructureDamageRepairCooldown=180
bPvEDisableFriendlyFire=False
bDisableFriendlyFire=False
`
}

// ValidateINIContent 验证INI内容的基本格式
func ValidateINIContent(content string) error {
	if content == "" {
		return nil // 空内容也是有效的
	}

	lines := strings.Split(content, "\n")

	for i, line := range lines {
		line = strings.TrimSpace(line)

		// 跳过空行和注释
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, ";") {
			continue
		}

		// 检查section格式
		if strings.HasPrefix(line, "[") {
			if !strings.HasSuffix(line, "]") {
				return fmt.Errorf("第%d行：section格式错误，缺少闭合括号", i+1)
			}
			continue
		}

		// 检查键值对格式
		if strings.Contains(line, "=") {
			parts := strings.SplitN(line, "=", 2)
			if len(parts) != 2 || strings.TrimSpace(parts[0]) == "" {
				return fmt.Errorf("第%d行：键值对格式错误", i+1)
			}
			continue
		}
	}

	return nil
}
