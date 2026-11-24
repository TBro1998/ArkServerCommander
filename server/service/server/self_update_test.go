package server

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"ark-server-commander/config"
	"ark-server-commander/database"
	"ark-server-commander/service/docker_manager"
	"ark-server-commander/utils"
)

func TestPullPercent(t *testing.T) {
	if p := pullPercent(nil); p != 0 {
		t.Errorf("nil = %d", p)
	}
	status := &docker_manager.ImageStatus{Layers: map[string]*docker_manager.LayerStatus{
		"a": {Status: "exists", Size: 1000}, // 本地已有，不计
		"b": {Status: "downloading", Size: 100, Progress: 50},
		"c": {Status: "extracting", Size: 100, Progress: 10}, // 已下载完
		"d": {Status: "pending"},
	}}
	if p := pullPercent(status); p != 75 {
		t.Errorf("进度 = %d，应为 75", p)
	}
	// 全部下载完也不显示 100，留给替换容器
	status = &docker_manager.ImageStatus{Layers: map[string]*docker_manager.LayerStatus{
		"a": {Status: "downloaded", Size: 10},
	}}
	if p := pullPercent(status); p != 99 {
		t.Errorf("进度 = %d，应为 99", p)
	}
}

func TestStartSelfUpdateRejectsDev(t *testing.T) {
	old := utils.Version
	utils.Version = "dev"
	t.Cleanup(func() { utils.Version = old })

	err := NewServerService().StartSelfUpdate("")
	if !errors.Is(err, ErrSelfUpdateUnsupported) {
		t.Fatalf("开发版本应返回 ErrSelfUpdateUnsupported，得到 %v", err)
	}
	if ErrorCode(err) != "self_update_unsupported" {
		t.Errorf("错误码 = %q", ErrorCode(err))
	}
}

func TestCheckSelfUpdateIdle(t *testing.T) {
	setupTestDB(t)
	if err := checkSelfUpdateIdle(); err != nil {
		t.Fatalf("没有进行中的操作时应允许更新: %v", err)
	}

	s := createTestServer(t, 7777, 27015, 32330)
	database.DB.Model(&s).Update("status", "starting")
	if err := checkSelfUpdateIdle(); !errors.Is(err, ErrSelfUpdateBusy) {
		t.Errorf("有服务器正在启动时应拒绝，得到 %v", err)
	}
	database.DB.Model(&s).Update("status", "running")
	if err := checkSelfUpdateIdle(); err != nil {
		t.Errorf("运行中的服务器不影响更新: %v", err)
	}

	unlock, err := lockBackup(s.ID)
	if err != nil {
		t.Fatal(err)
	}
	if err := checkSelfUpdateIdle(); !errors.Is(err, ErrSelfUpdateBusy) {
		t.Errorf("正在备份时应拒绝，得到 %v", err)
	}
	unlock()
	if err := checkSelfUpdateIdle(); err != nil {
		t.Errorf("备份结束后应允许更新: %v", err)
	}
}

func TestBackupDatabaseForUpdate(t *testing.T) {
	setupTestDB(t)
	oldPath := config.DBPath
	config.DBPath = filepath.Join(t.TempDir(), "ark_server.db")
	t.Cleanup(func() { config.DBPath = oldPath })

	// 第二次覆盖上一次的备份
	for i := 0; i < 2; i++ {
		if err := backupDatabaseForUpdate(); err != nil {
			t.Fatal(err)
		}
	}
	info, err := os.Stat(config.DBPath + selfUpdateBackupSuffix)
	if err != nil || info.Size() == 0 {
		t.Fatalf("没有生成数据库备份: %v", err)
	}
}

func TestSelfUpdateStatusLastRecord(t *testing.T) {
	oldDir := config.DataDir
	config.DataDir = t.TempDir()
	t.Cleanup(func() { config.DataDir = oldDir })

	st := NewServerService().GetSelfUpdateStatus()
	if st.State != "idle" || st.Last != nil {
		t.Errorf("没有记录时: %+v", st)
	}

	rec := docker_manager.SelfUpdateRecord{State: docker_manager.SelfUpdateRolledBack, FromVersion: "0.0.5", Error: "新容器退出了"}
	if err := docker_manager.WriteSelfUpdateRecord(config.DataDir, rec); err != nil {
		t.Fatal(err)
	}
	st = NewServerService().GetSelfUpdateStatus()
	if st.Last == nil || st.Last.State != docker_manager.SelfUpdateRolledBack || st.Last.Error != rec.Error {
		t.Errorf("应返回最近一次的记录: %+v", st.Last)
	}

	// 辅助容器没写结果就没了
	rec.State = docker_manager.SelfUpdateStarted
	if err := docker_manager.WriteSelfUpdateRecord(config.DataDir, rec); err != nil {
		t.Fatal(err)
	}
	st = NewServerService().GetSelfUpdateStatus()
	if st.Last == nil || st.Last.State != "interrupted" {
		t.Errorf("没有结果、也没有在更新时应为 interrupted: %+v", st.Last)
	}
}
