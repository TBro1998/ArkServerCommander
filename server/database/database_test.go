package database

import (
	"path/filepath"
	"testing"

	"ark-server-commander/models"
	"ark-server-commander/utils"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

// 以前的版本允许同时开启发布和联机房：发布成功过的保留发布，其他的保留联机房，只开一个的不变
func TestFixPublishAndLink(t *testing.T) {
	utils.InitLogger()
	db, err := gorm.Open(sqlite.Open(filepath.Join(t.TempDir(), "test.db")), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.AutoMigrate(&models.Server{}); err != nil {
		t.Fatal(err)
	}
	old := DB
	DB = db
	t.Cleanup(func() {
		DB = old
		if sqlDB, err := db.DB(); err == nil {
			sqlDB.Close()
		}
	})

	servers := []models.Server{
		{SessionName: "published", Port: 7777, Publish: true, Link: true, Published: true},
		{SessionName: "unpublished", Port: 7779, Publish: true, Link: true},
		{SessionName: "publish only", Port: 7781, Publish: true, Published: true},
		{SessionName: "link only", Port: 7783, Link: true},
	}
	for i := range servers {
		if err := db.Create(&servers[i]).Error; err != nil {
			t.Fatal(err)
		}
		// default:true 之外的 bool 字段 GORM 照常写入；这里显式写回，避免默认值影响
		db.Model(&servers[i]).Updates(map[string]any{"publish": servers[i].Publish, "link": servers[i].Link, "published": servers[i].Published})
	}

	fixPublishAndLink()

	want := map[string][2]bool{
		"published":    {true, false},
		"unpublished":  {false, true},
		"publish only": {true, false},
		"link only":    {false, true},
	}
	var got []models.Server
	db.Find(&got)
	for _, s := range got {
		if w := want[s.SessionName]; s.Publish != w[0] || s.Link != w[1] {
			t.Errorf("%s: publish = %v, link = %v, want %v", s.SessionName, s.Publish, s.Link, w)
		}
	}
}
