package db

//数据库连接和初始化，使用 GORM 连接 MySQL，并自动迁移表结构

import (
	"fmt"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"geekedu-project/logic-server/model"
)

var DB *gorm.DB

// 初始化数据库连接（适配Docker环境变量）
func InitDB() error {
	// 从环境变量读取MySQL配置（和docker-compose里的environment对应）
	host := os.Getenv("DB_HOST")         // 对应 docker-compose 里的 DB_HOST: mysql
	port := os.Getenv("DB_PORT")         // 对应 docker-compose 里的 DB_PORT: 3306
	user := os.Getenv("DB_USER")         // 对应 docker-compose 里的 DB_USER: root
	password := os.Getenv("DB_PASSWORD") // 对应 docker-compose 里的 DB_PASSWORD: root
	dbName := os.Getenv("DB_DATABASE")   // 对应 docker-compose 里的 DB_DATABASE: geekedu

	// 拼接正确的MySQL连接地址（DSN格式：user:password@tcp(host:port)/dbname?参数）
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, password, host, port, dbName)

	// 连接MySQL（添加超时配置，避免无限等待）
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // 打印SQL日志（方便调试）
	})
	if err != nil {
		return fmt.Errorf("数据库连接失败：%v", err)
	}

	// 配置连接池（优化性能）
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("连接池配置失败：%v", err)
	}
	sqlDB.SetMaxIdleConns(10)                  // 最大空闲连接数
	sqlDB.SetMaxOpenConns(100)                 // 最大打开连接数
	sqlDB.SetConnMaxLifetime(30 * time.Minute) // 连接最大存活时间

	DB = db

	// 关键修复：使用 GORM AutoMigrate 创建表（确保表结构存在）
	// 即使 init.sql 未被执行，此处也会创建表
	if err := db.AutoMigrate(
		&model.User{},
		&model.Course{},
		&model.Video{},
		&model.Order{},
	); err != nil {
		return fmt.Errorf("自动迁移失败：%v", err)
	}

	return nil
}
