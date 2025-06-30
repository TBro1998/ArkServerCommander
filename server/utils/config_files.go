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

	// 配置文件目录
	ConfigDirectory = "Saved/Config/WindowsServer"
)

// GetDefaultGameUserSettings 获取默认的GameUserSettings.ini配置
func GetDefaultGameUserSettings(serverName, mapName string, port, queryPort, rconPort, maxPlayers int, adminPassword string) string {
	return fmt.Sprintf(`[ServerSettings]
SessionName=%s
ServerPassword=
ServerAdminPassword=%s
Port=%d
QueryPort=%d
RCONEnabled=True
RCONPort=%d
MaxPlayers=%d

[SessionSettings]
SessionName=%s

[MessageOfTheDay]
Message=欢迎来到 %s ARK 服务器！

[/Script/ShooterGame.ShooterGameMode]
bUseSingleplayerSettings=False
bDisableStructurePlacementCollision=False
bAllowFlyerCarryPvE=True
bDisableStructureDecayPvE=False

[RCONSettings]
RCONEnabled=True
RCONPort=%d
`, serverName, adminPassword, port, queryPort, rconPort, maxPlayers, serverName, serverName, rconPort)
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

// EnsureConfigDirectory 确保配置文件目录存在
func EnsureConfigDirectory(serverID uint) error {
	configPath := GetServerConfigPath(serverID)
	return os.MkdirAll(configPath, 0755)
}

// ReadConfigFile 读取配置文件内容
func ReadConfigFile(serverID uint, fileName string) (string, error) {
	filePath := GetConfigFilePath(serverID, fileName)

	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return "", fmt.Errorf("配置文件不存在: %s", fileName)
	}

	// 读取文件内容
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("读取配置文件失败: %v", err)
	}

	return string(content), nil
}

// WriteConfigFile 写入配置文件内容
func WriteConfigFile(serverID uint, fileName, content string) error {
	// 确保配置目录存在
	if err := EnsureConfigDirectory(serverID); err != nil {
		return fmt.Errorf("创建配置目录失败: %v", err)
	}

	filePath := GetConfigFilePath(serverID, fileName)

	// 写入文件
	err := os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		return fmt.Errorf("写入配置文件失败: %v", err)
	}

	return nil
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

// CreateDefaultConfigFiles 为服务器创建默认配置文件
func CreateDefaultConfigFiles(serverID uint, serverName, mapName string, port, queryPort, rconPort, maxPlayers int, adminPassword string) error {
	// 确保配置目录存在
	if err := EnsureConfigDirectory(serverID); err != nil {
		return err
	}

	// 创建默认的 GameUserSettings.ini
	gameUserSettings := GetDefaultGameUserSettings(serverName, mapName, port, queryPort, rconPort, maxPlayers, adminPassword)
	if err := WriteConfigFile(serverID, GameUserSettingsFileName, gameUserSettings); err != nil {
		return fmt.Errorf("创建默认 %s 失败: %v", GameUserSettingsFileName, err)
	}

	// 创建默认的 Game.ini
	gameIni := GetDefaultGameIni()
	if err := WriteConfigFile(serverID, GameIniFileName, gameIni); err != nil {
		return fmt.Errorf("创建默认 %s 失败: %v", GameIniFileName, err)
	}

	return nil
}

// DeleteConfigFiles 删除服务器的配置文件
func DeleteConfigFiles(serverID uint) error {
	configPath := GetServerConfigPath(serverID)

	// 删除整个配置目录
	if err := os.RemoveAll(configPath); err != nil {
		return fmt.Errorf("删除配置文件目录失败: %v", err)
	}

	return nil
}

// ListConfigFiles 列出服务器配置目录下的所有文件
func ListConfigFiles(serverID uint) ([]models.ServerConfigFileInfo, error) {
	configPath := GetServerConfigPath(serverID)

	var files []models.ServerConfigFileInfo

	// 检查配置目录是否存在
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return files, nil // 返回空列表而不是错误
	}

	// 读取目录内容
	entries, err := os.ReadDir(configPath)
	if err != nil {
		return nil, fmt.Errorf("读取配置目录失败: %v", err)
	}

	// 处理每个文件
	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".ini" {
			fileInfo := GetConfigFileInfo(serverID, entry.Name())
			files = append(files, fileInfo)
		}
	}

	return files, nil
}
