package database

import (
	"fmt"
	"log"
	"time"

	"admin/config"
	"admin/internal/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

// InitDB 初始化数据库连接
func InitDB(cfg *config.Config) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.Database.Username,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.DBName,
	)
	//log.Printf("数据库连接DSN: %s", dsn) // 添加这行

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}

	// 获取底层 sql.DB 以配置连接池
	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatalf("获取数据库连接池失败: %v", err)
	}

	// 设置连接池参数
	sqlDB.SetMaxIdleConns(10)           // 最大空闲连接数
	sqlDB.SetMaxOpenConns(100)          // 最大打开连接数
	sqlDB.SetConnMaxLifetime(time.Hour) // 连接最大生命周期

	// 自动迁移表结构
	log.Println("正在执行数据库迁移...")
	if err := DB.AutoMigrate(&models.Admin{}); err != nil {
		log.Fatalf("数据库迁移失败: %v", err)
	}

	log.Println("数据库连接成功")

}
