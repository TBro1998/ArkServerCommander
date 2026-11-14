package apitoken

import (
	"errors"
	"path/filepath"
	"strings"
	"testing"

	"ark-server-commander/database"
	"ark-server-commander/models"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(filepath.Join(t.TempDir(), "test.db")), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.AutoMigrate(&models.User{}, &models.APIToken{}); err != nil {
		t.Fatal(err)
	}
	if err := db.Create(&models.User{Username: "admin", Password: "x"}).Error; err != nil {
		t.Fatal(err)
	}

	old := database.DB
	database.DB = db
	t.Cleanup(func() {
		database.DB = old
		if sqlDB, err := db.DB(); err == nil {
			sqlDB.Close()
		}
	})
}

func TestTokenLifecycle(t *testing.T) {
	setupTestDB(t)

	created, err := Create(1, "  claude  ", true)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(created.Token, TokenPrefix) || !strings.HasPrefix(created.Token, created.Prefix) {
		t.Fatalf("令牌格式不对: token=%q prefix=%q", created.Token, created.Prefix)
	}
	if created.Name != "claude" || !created.ReadOnly {
		t.Fatalf("名称或只读标记不对: %+v", created.APIToken)
	}

	// 数据库里不能有明文
	var stored models.APIToken
	database.DB.First(&stored, created.ID)
	if stored.TokenHash == created.Token || strings.Contains(stored.TokenHash, created.Token) {
		t.Fatal("数据库中保存了明文令牌")
	}

	got, err := Verify(created.Token)
	if err != nil {
		t.Fatal(err)
	}
	if got.ID != created.ID || got.UserID != 1 || !got.ReadOnly || got.LastUsedAt == nil {
		t.Fatalf("校验结果不对: %+v", got)
	}

	for _, bad := range []string{"", "asc_", created.Token + "x", strings.TrimPrefix(created.Token, TokenPrefix)} {
		if _, err := Verify(bad); !errors.Is(err, ErrInvalidToken) {
			t.Errorf("Verify(%q) 应返回 ErrInvalidToken，实际 %v", bad, err)
		}
	}

	list, err := List(1)
	if err != nil || len(list) != 1 {
		t.Fatalf("List = %v, %v", list, err)
	}

	if err := Delete(2, created.ID); !errors.Is(err, ErrTokenNotFound) {
		t.Fatalf("不能删除其他用户的令牌，实际 %v", err)
	}
	if err := Delete(1, created.ID); err != nil {
		t.Fatal(err)
	}
	if _, err := Verify(created.Token); !errors.Is(err, ErrInvalidToken) {
		t.Fatalf("删除后令牌应失效，实际 %v", err)
	}
}

func TestTokenRejectsInvalidName(t *testing.T) {
	setupTestDB(t)

	for _, name := range []string{"", "   ", strings.Repeat("长", 65)} {
		if _, err := Create(1, name, false); !errors.Is(err, ErrInvalidTokenName) {
			t.Errorf("Create(%q) 应返回 ErrInvalidTokenName，实际 %v", name, err)
		}
	}
}

func TestTokenOfDeletedUserIsInvalid(t *testing.T) {
	setupTestDB(t)

	created, err := Create(1, "agent", false)
	if err != nil {
		t.Fatal(err)
	}
	database.DB.Delete(&models.User{}, 1)
	if _, err := Verify(created.Token); !errors.Is(err, ErrInvalidToken) {
		t.Fatalf("用户删除后令牌应失效，实际 %v", err)
	}
}
