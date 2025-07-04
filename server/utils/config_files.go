package utils

import (
	"ark-server-manager/models"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	// 配置文件名常量
	GameUserSettingsFileName = "GameUserSettings.ini"
	GameIniFileName          = "Game.ini"

	// 配置文件目录 - 适配新的ASE服务器镜像路径
	ConfigDirectory = "Config/WindowsServer"
)

// GetDefaultGameUserSettings 获取默认的GameUserSettings.ini配置
func GetDefaultGameUserSettings(serverName, mapName string, maxPlayers int) string {
	return fmt.Sprintf(`[ServerSettings]
SessionName=%s
ServerPassword=
MaxPlayers=%d

[SessionSettings]
SessionName=%s

[MessageOfTheDay]
Message=欢迎来到 %s ARK 服务器！

[/Script/Engine.GameSession]
MaxPlayers=%d
`, serverName, maxPlayers, serverName, serverName, maxPlayers)
}

// GetDefaultGameIni 获取默认的Game.ini配置
func GetDefaultGameIni() string {
	return `[/script/shootergame.shootergamemode]
bUseSingleplayerSettings=false
bDisableStructurePlacementCollision=false
bAllowFlyerCarryPvE=true
bDisableStructureDecayPvE=false
bAllowUnlimitedRespecs=true
bAllowPlatformSaddleMultiFloors=true
bPassiveDefensesDamageRiderlessDinos=true
MaxNumberOfPlayersInTribe=0

[/script/engine.gamesession]
MaxPlayers=70

[/Script/ShooterGame.ShooterGameMode]
DifficultyOffset=1.0
OverrideOfficialDifficulty=5.0

# 资源重生倍率
ResourcesRespawnPeriodMultiplier=1.0

# 驯服相关设置
TamingSpeedMultiplier=1.0
DinoCharacterFoodDrainMultiplier=1.0
DinoCharacterStaminaDrainMultiplier=1.0
DinoCharacterHealthRecoveryMultiplier=1.0
DinoCountMultiplier=1.0

# 经验值倍率
XPMultiplier=1.0
PlayerCharacterWaterDrainMultiplier=1.0
PlayerCharacterFoodDrainMultiplier=1.0
PlayerCharacterStaminaDrainMultiplier=1.0
PlayerCharacterHealthRecoveryMultiplier=1.0

# 收集倍率
HarvestAmountMultiplier=1.0
HarvestHealthMultiplier=1.0

# 白天/夜晚时间流逝速度
DayCycleSpeedScale=1.0
NightTimeSpeedScale=1.0

# 结构相关
StructureResistanceMultiplier=1.0
StructureDamageMultiplier=1.0
StructureDamageRepairCooldown=180
PvEStructureDecayPeriodMultiplier=1.0

# PvP相关设置
bPvEDisableFriendlyFire=False
bEnablePvPGamma=False
bDisableFriendlyFire=False
bAllowFlyerCarryPvE=True
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

// GetServerConfigPath 获取服务器配置文件目录路径
func GetServerConfigPath(serverID uint) string {
	serverFolder := GetServerFolderPath(serverID)
	return filepath.Join(serverFolder, ConfigDirectory)
}

// GetConfigFilePath 获取配置文件的完整路径
func GetConfigFilePath(serverID uint, fileName string) string {
	configPath := GetServerConfigPath(serverID)
	return filepath.Join(configPath, fileName)
}

// GetConfigFileInfo 获取配置文件信息
func GetConfigFileInfo(serverID uint, fileName string) models.ServerConfigFileInfo {
	filePath := GetConfigFilePath(serverID, fileName)

	info := models.ServerConfigFileInfo{
		FileName: fileName,
		FilePath: filePath,
		Exists:   false,
	}

	// 获取文件信息
	fileInfo, err := os.Stat(filePath)
	if err == nil {
		info.Exists = true
		info.FileSize = fileInfo.Size()
		info.LastModified = fileInfo.ModTime().Format("2006-01-02 15:04:05")
	}

	return info
}

// GenerateServerArgs 生成ARK服务器的启动参数
// 已废弃，统一用ServerArgs.GenerateArgsString
func GenerateServerArgs(server models.Server) string {
	return ""
}

// GenerateGameModIdsEnv 生成GameModIds环境变量
// 如果GameModIds为空，返回空字符串
func GenerateGameModIdsEnv(gameModIds string) string {
	if gameModIds == "" {
		return ""
	}
	return gameModIds
}
