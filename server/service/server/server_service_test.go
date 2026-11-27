package server

import (
	"errors"
	"fmt"
	"path/filepath"
	"testing"

	"ark-server-commander/database"
	"ark-server-commander/models"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

// setupTestDB 把全局 database.DB 换成临时的 SQLite，测试结束后还原
func setupTestDB(t *testing.T) {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(filepath.Join(t.TempDir(), "test.db")), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.AutoMigrate(&models.User{}, &models.Server{}, &models.Schedule{}, &models.LinkRoom{}, &models.ArkfgfAccount{}); err != nil {
		t.Fatal(err)
	}
	if err := db.Create(&models.User{Username: "admin", Password: "x"}).Error; err != nil {
		t.Fatal(err)
	}

	old := database.DB
	database.DB = db
	// 账号状态缓存在内存里，换数据库后重新读取
	accountManager.forget()
	t.Cleanup(func() {
		database.DB = old
		accountManager.forget()
		// 先关闭连接，Windows 上打开着的数据库文件无法被 TempDir 清理
		if sqlDB, err := db.DB(); err == nil {
			sqlDB.Close()
		}
	})
}

func createTestServer(t *testing.T, port, queryPort, rconPort int) models.Server {
	t.Helper()

	s := models.Server{
		SessionName:   "test",
		Port:          port,
		QueryPort:     queryPort,
		RCONPort:      rconPort,
		AdminPassword: "pw",
		MaxPlayers:    70,
		Status:        "stopped",
		UserID:        1,
	}
	if err := database.DB.Create(&s).Error; err != nil {
		t.Fatal(err)
	}
	return s
}

func TestCheckPorts(t *testing.T) {
	setupTestDB(t)
	a := createTestServer(t, 7777, 27015, 32330)

	tests := []struct {
		name              string
		excludeID         uint
		port, query, rcon int
		want              error
	}{
		{"完全相同", 0, 7777, 27015, 32330, ErrPortConflict},
		{"游戏端口+1 撞上其他服务器的游戏端口", 0, 7776, 27016, 32331, ErrPortConflict},
		{"查询端口相同", 0, 7780, 27015, 32331, ErrPortConflict},
		{"RCON 端口相同", 0, 7780, 27016, 32330, ErrPortConflict},
		{"互不冲突", 0, 7780, 27016, 32331, nil},
		{"排除自己后不冲突", a.ID, 7777, 27015, 32330, nil},
		{"本服务器内部重复", 0, 7780, 27016, 7781, ErrPortConflict},
		{"游戏端口+1 超出范围", 0, 65535, 27016, 32331, ErrInvalidPort},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := checkPorts(tt.excludeID, tt.port, tt.query, tt.rcon)
			if !errors.Is(err, tt.want) {
				t.Fatalf("期望 %v，实际 %v", tt.want, err)
			}
		})
	}
}

// 已删除的服务器容器已经移除，不应继续占用端口
func TestCheckPortsIgnoresDeletedServers(t *testing.T) {
	setupTestDB(t)
	a := createTestServer(t, 7777, 27015, 32330)

	if err := checkPorts(0, 7777, 27015, 32330); !errors.Is(err, ErrPortConflict) {
		t.Fatalf("删除前应该冲突，实际 %v", err)
	}
	if err := database.DB.Delete(&a).Error; err != nil {
		t.Fatal(err)
	}
	if err := checkPorts(0, 7777, 27015, 32330); err != nil {
		t.Fatalf("删除后应该可用，实际 %v", err)
	}
}

// UpdateServer 保存时不能把读取之后被启停协程改掉的状态写回旧值
func TestUpdateSaveDoesNotClobberStatus(t *testing.T) {
	setupTestDB(t)
	s := createTestServer(t, 7777, 27015, 32330)

	var stale models.Server
	if err := database.DB.First(&stale, s.ID).Error; err != nil {
		t.Fatal(err)
	}

	// 模拟读取之后，启动协程把状态改成了 starting
	finishOperation(s.ID, "starting", "")

	stale.SessionName = "renamed"
	if err := database.DB.Omit("Status", "LastError").Save(&stale).Error; err != nil {
		t.Fatal(err)
	}

	var got models.Server
	if err := database.DB.First(&got, s.ID).Error; err != nil {
		t.Fatal(err)
	}
	if got.Status != "starting" {
		t.Errorf("状态被覆盖了: %q", got.Status)
	}
	if got.SessionName != "renamed" {
		t.Errorf("名称没有保存: %q", got.SessionName)
	}
}

// GORM 创建记录时会把 bool 零值换成 default:true，关闭的开关需要显式写回
func TestCreateScheduleKeepsDisabled(t *testing.T) {
	setupTestDB(t)
	s := createTestServer(t, 7777, 27015, 32330)

	enabled := false
	sch, err := NewServerService().CreateSchedule(1, fmt.Sprint(s.ID), models.ScheduleRequest{
		Action: models.ScheduleBackup, Cron: "@daily", Enabled: &enabled,
	})
	if err != nil {
		t.Fatal(err)
	}
	if sch.Enabled {
		t.Fatal("response should report the schedule as disabled")
	}
	var stored models.Schedule
	database.DB.First(&stored, sch.ID)
	if stored.Enabled {
		t.Fatal("schedule should be stored as disabled")
	}
}

func TestGetUserServerChecksOwner(t *testing.T) {
	setupTestDB(t)
	s := createTestServer(t, 7777, 27015, 32330)

	if _, err := getUserServer(1, fmt.Sprint(s.ID)); err != nil {
		t.Fatalf("owner should get the server: %v", err)
	}
	if _, err := getUserServer(2, fmt.Sprint(s.ID)); !errors.Is(err, ErrServerNotFound) {
		t.Fatalf("other users should get ErrServerNotFound, got %v", err)
	}
	if _, err := getUserServer(1, "abc"); !errors.Is(err, ErrInvalidServerID) {
		t.Fatalf("expected ErrInvalidServerID, got %v", err)
	}
}

func TestSuggestPortsSkipsUsedPorts(t *testing.T) {
	setupTestDB(t)
	createTestServer(t, 7777, 27015, 32330)

	p := NewServerService().SuggestPorts()
	if p.Port == 7777 || p.Port == 7778 || p.QueryPort == 27015 || p.RCONPort == 32330 {
		t.Fatalf("suggested ports overlap the existing server: %+v", p)
	}
	if err := checkPorts(0, p.Port, p.QueryPort, p.RCONPort); err != nil {
		t.Fatalf("suggested ports should be free: %v", err)
	}
}

// 发布到收录站和联机房只能二选一：创建、修改时两个都开返回 ErrPublishLinkExclusive，在同一个请求里关掉另一个可以切换
func TestPublishLinkExclusive(t *testing.T) {
	setupTestDB(t)
	svc := &ServerService{}
	on, off := true, false
	if _, err := svc.CreateServer(1, models.ServerRequest{Port: 7790, QueryPort: 27090, RCONPort: 27091, AdminPassword: "pw", MaxPlayers: 10,
		Publish: &on, Link: &on}); !errors.Is(err, ErrPublishLinkExclusive) {
		t.Fatalf("create with both: %v", err)
	}

	s := createTestServer(t, 7777, 27015, 32330)
	id := fmt.Sprint(s.ID)
	if _, _, err := svc.UpdateServer(1, id, models.ServerUpdateRequest{Link: &on}); err != nil {
		t.Fatal(err)
	}
	if _, _, err := svc.UpdateServer(1, id, models.ServerUpdateRequest{Publish: &on}); !errors.Is(err, ErrPublishLinkExclusive) {
		t.Fatalf("publish while linked: %v", err)
	}
	if _, _, err := svc.UpdateServer(1, id, models.ServerUpdateRequest{Publish: &on, Link: &off}); err != nil {
		t.Fatalf("switch to publish: %v", err)
	}
	var got models.Server
	database.DB.First(&got, s.ID)
	if !got.Publish || got.Link {
		t.Fatalf("publish = %v, link = %v", got.Publish, got.Link)
	}
}
