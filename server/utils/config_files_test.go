package utils

import (
	"strings"
	"testing"
)

// 这些是 GameUserSettings.ini [ServerSettings] 的选项，写在 Game.ini 里不会生效
var serverSettingsOnly = []string{
	"XPMultiplier", "TamingSpeedMultiplier", "HarvestAmountMultiplier", "DifficultyOffset",
	"OverrideOfficialDifficulty", "DayCycleSpeedScale", "AllowFlyerCarryPvE", "DisableStructureDecayPvE", "EnablePvPGamma",
}

func iniKeys(content string) (keys map[string]string, sections []string) {
	keys = map[string]string{}
	section := ""
	for _, line := range strings.Split(content, "\n") {
		line = strings.TrimSpace(line)
		switch {
		case line == "":
		case strings.HasPrefix(line, "["):
			section = line
			sections = append(sections, line)
		default:
			k, _, _ := strings.Cut(line, "=")
			keys[k] = section
		}
	}
	return keys, sections
}

func TestDefaultConfigSections(t *testing.T) {
	gameKeys, sections := iniKeys(GetDefaultGameIni())
	if len(sections) != 1 || sections[0] != "[/script/shootergame.shootergamemode]" {
		t.Fatalf("Game.ini should only have the shootergamemode section, got %v", sections)
	}
	gusKeys, _ := iniKeys(GetDefaultGameUserSettings(GameASE, "test", "TheIsland", 70))
	for _, k := range serverSettingsOnly {
		for gk := range gameKeys {
			if strings.EqualFold(strings.TrimPrefix(gk, "b"), k) {
				t.Errorf("%s belongs to GameUserSettings.ini, found %s in Game.ini", k, gk)
			}
		}
		if gusKeys[k] != "[ServerSettings]" {
			t.Errorf("%s should be under [ServerSettings], got %q", k, gusKeys[k])
		}
	}
}

// ASA 的人数上限只看 -WinLiveMaxPlayers，默认模板不写 GameSession；ASE 照常写
func TestDefaultGameUserSettingsGameSession(t *testing.T) {
	_, aseSections := iniKeys(GetDefaultGameUserSettings(GameASE, "test", "TheIsland", 70))
	_, asaSections := iniKeys(GetDefaultGameUserSettings(GameASA, "test", "TheIsland_WP", 70))
	has := func(sections []string) bool {
		for _, s := range sections {
			if s == "[/Script/Engine.GameSession]" {
				return true
			}
		}
		return false
	}
	if !has(aseSections) {
		t.Error("ASE template should have [/Script/Engine.GameSession]")
	}
	if has(asaSections) {
		t.Error("ASA template should not have [/Script/Engine.GameSession]")
	}
}
