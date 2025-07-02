package database

import (
	"log"

	"ark-server-manager/config"
	"ark-server-manager/models"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB() {
	var err error

	// 连接SQLite数据库
	DB, err = gorm.Open(sqlite.Open(config.DBPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatal("数据库连接失败:", err)
	}

	// 自动迁移数据库结构
	err = DB.AutoMigrate(&models.User{}, &models.Server{})
	if err != nil {
		log.Fatal("数据库迁移失败:", err)
	}

	log.Println("数据库初始化成功")
}

func GetDB() *gorm.DB {
	return DB
}
