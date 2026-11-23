package server

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"ark-server-commander/database"
	"ark-server-commander/models"
	"ark-server-commander/service/docker_manager"
	"ark-server-commander/utils"
)

func TestWithAdminPassword(t *testing.T) {
	got, err := withAdminPassword("[ServerSettings]\nSessionName=a\n", "pw1")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(got, "ServerAdminPassword=pw1") {
		t.Fatalf("没有补上密码: %q", got)
	}

	// 已有不同的密码时替换，不重复
	got, err = withAdminPassword(got, "pw2")
	if err != nil {
		t.Fatal(err)
	}
	if strings.Count(got, "ServerAdminPassword=") != 1 || !strings.Contains(got, "ServerAdminPassword=pw2") {
		t.Fatalf("密码没有被替换: %q", got)
	}

	if _, err := withAdminPassword("", "a\nb"); err == nil {
		t.Fatal("密码含换行时应报错")
	}
}

// 编辑器、暂存配置、保存设置写 GameUserSettings.ini 时，即使内容里没有密码，写出的文件也要带上，
// 否则容器被 Docker 重启策略直接拉起时 ARK 不会开启 RCON
func TestConfigWritesKeepAdminPassword(t *testing.T) {
	setupTestDB(t)
	useTempDataDir(t)

	dm, err := docker_manager.GetDockerManager()
	if err != nil {
		t.Fatal(err)
	}
	srv := createTestServer(t, 7777, 27015, 32330)
	if err := os.MkdirAll(docker_manager.ServerSavedDir(srv.ID), 0o755); err != nil {
		t.Fatal(err)
	}
	read := func() string {
		t.Helper()
		content, err := dm.ReadConfigFile(srv.ID, utils.GameUserSettingsFileName)
		if err != nil {
			t.Fatal(err)
		}
		return content
	}

	// 直接写入（创建服务器走这里）
	if err := writeGameUserSettings(dm, srv.ID, srv.AdminPassword, "[ServerSettings]\nSessionName=x\n"); err != nil {
		t.Fatal(err)
	}
	if got := read(); !strings.Contains(got, "ServerAdminPassword=pw") || !strings.Contains(got, "SessionName=x") {
		t.Fatalf("直接写入: %q", got)
	}

	// 应用暂存的配置
	database.DB.Model(&srv).Update("pending_game_user_settings", "[ServerSettings]\nSessionName=y\n")
	if err := applyPendingConfig(srv.ID, dm); err != nil {
		t.Fatal(err)
	}
	if got := read(); !strings.Contains(got, "ServerAdminPassword=pw") || !strings.Contains(got, "SessionName=y") {
		t.Fatalf("应用暂存配置: %q", got)
	}

	// 已停止的服务器保存配置（编辑器、导入、复制、set_config_values 都经过 UpdateServer）
	if _, _, err := NewServerService().UpdateServer(1, fmt.Sprint(srv.ID), models.ServerUpdateRequest{
		GameUserSettings: "[ServerSettings]\nSessionName=z\n",
	}); err != nil {
		t.Fatal(err)
	}
	if got := read(); !strings.Contains(got, "ServerAdminPassword=pw") || !strings.Contains(got, "SessionName=z") {
		t.Fatalf("保存配置: %q", got)
	}
}
