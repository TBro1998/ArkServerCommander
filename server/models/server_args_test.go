package models

import (
	"testing"
)

func TestGenerateArgsString(t *testing.T) {
	// 创建测试服务器
	server := Server{
		ID:            1,
		Identifier:    "test-server",
		Port:          7777,
		QueryPort:     27015,
		RCONPort:      32330,
		AdminPassword: "testpassword",
		Map:           "TheIsland",
		GameModIds:    "123456,789012",
	}

	// 创建启动参数
	serverArgs := NewServerArgs()
	serverArgs.QueryParams["AllowThirdPersonPlayer"] = "True"
	serverArgs.CommandLineArgs["NoBattlEye"] = true
	serverArgs.CommandLineArgs["servergamelog"] = true

	// 生成启动参数字符串
	result := serverArgs.GenerateArgsString(server)

	// 验证结果包含基础参数
	expectedParts := []string{
		"TheIsland",
		"?listen",
		"?Port=7777",
		"?QueryPort=27015",
		"?MaxPlayers=70",
		"?RCONEnabled=True",
		"?RCONPort=32330",
		"?ServerAdminPassword=testpassword",
		"?GameModIds=123456,789012",
		"?AllowThirdPersonPlayer=True",
		"-NoBattlEye",
		"-servergamelog",
	}

	for _, part := range expectedParts {
		if !contains(result, part) {
			t.Errorf("生成的启动参数字符串缺少: %s", part)
		}
	}

	t.Logf("生成的启动参数字符串: %s", result)
}

func TestFromServer(t *testing.T) {
	// 创建测试服务器
	server := Server{
		ID:            1,
		Identifier:    "test-server",
		Port:          7777,
		QueryPort:     27015,
		RCONPort:      32330,
		AdminPassword: "testpassword",
		Map:           "TheIsland",
		GameModIds:    "123456,789012",
	}

	// 从服务器创建启动参数
	serverArgs := FromServer(server)

	// 验证默认参数设置
	if serverArgs.QueryParams["listen"] != "" {
		t.Error("listen参数应该设置为空字符串")
	}

	// 验证默认命令行参数
	expectedCommandLineArgs := map[string]interface{}{
		"NoBattlEye":    true,
		"servergamelog": true,
		"game":          true,
		"server":        true,
	}

	for key, expectedValue := range expectedCommandLineArgs {
		if value, exists := serverArgs.CommandLineArgs[key]; !exists {
			t.Errorf("缺少默认命令行参数: %s", key)
		} else if value != expectedValue {
			t.Errorf("命令行参数 %s 的值不正确，期望: %v，实际: %v", key, expectedValue, value)
		}
	}

	// 验证不包含基础参数
	basicParams := []string{"MaxPlayers", "RCONEnabled", "ServerAdminPassword"}
	for _, param := range basicParams {
		if _, exists := serverArgs.QueryParams[param]; exists {
			t.Errorf("不应该包含基础参数: %s", param)
		}
	}
}

// 辅助函数：检查字符串是否包含子字符串
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) &&
		(s[:len(substr)] == substr || s[len(s)-len(substr):] == substr ||
			containsSubstring(s, substr)))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
