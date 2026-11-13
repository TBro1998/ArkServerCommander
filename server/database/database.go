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

	// arkapi 列是后加的：在此之前创建的服务器都是 Proton 版并带 ArkApi，迁移后保持原样，新服务器默认用 Linux 版
	addingArkApi := DB.Migrator().HasTable(&models.Server{}) && !DB.Migrator().HasColumn(&models.Server{}, "arkapi")

	// 自动迁移数据库结构
	err = DB.AutoMigrate(&models.User{}, &models.Server{}, &models.Schedule{})
	if err != nil {
		utils.Fatal("数据库迁移失败", zap.Error(err))
	}

	if addingArkApi {
		res := DB.Exec("UPDATE servers SET arkapi = ?", true)
		if res.Error != nil {
			utils.Fatal("数据库迁移失败", zap.Error(res.Error))
		}
		utils.Info("已有服务器保持启用 ArkApi（Proton 镜像）", zap.Int64("servers", res.RowsAffected))
	}

	utils.Info("数据库初始化成功", zap.String("db_path", config.DBPath))
}
