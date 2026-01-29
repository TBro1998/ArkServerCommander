package database

import (
	"ark-server-commander/config"
	"ark-server-commander/models"
	"ark-server-commander/utils"

	"github.com/glebarez/sqlite"
	"go.uber.org/zap"
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
		utils.Fatal("数据库连接失败", zap.Error(err))
	}

	// 自动迁移数据库结构
	err = DB.AutoMigrate(&models.User{}, &models.Server{})
	if err != nil {
		utils.Fatal("数据库迁移失败", zap.Error(err))
	}

	utils.Info("数据库初始化成功", zap.String("db_path", config.DBPath))
}

func GetDB() *gorm.DB {
	return DB
}
