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

	// 连接SQLite数据库。WAL 模式下读写互不阻塞，busy_timeout 让并发写入排队等待，而不是立即报 database is locked
	dsn := config.DBPath + "?_pragma=busy_timeout(5000)&_pragma=journal_mode(WAL)"
	DB, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{
		// Info 级别会打印每条 SQL，其中包含保存服务器时的管理员密码明文
		Logger: logger.Default.LogMode(logger.Warn),
	})

	if err != nil {
		utils.Fatal("数据库连接失败", zap.Error(err))
	}

	// 自动迁移数据库结构
	err = DB.AutoMigrate(&models.User{}, &models.Server{}, &models.Schedule{})
	if err != nil {
		utils.Fatal("数据库迁移失败", zap.Error(err))
	}

	utils.Info("数据库初始化成功", zap.String("db_path", config.DBPath))
}
