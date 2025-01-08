package dal

import (
	"database/sql"
	"fmt"

	"github.com/InnoServe/blockSBOM/internal/config"
	"github.com/InnoServe/blockSBOM/internal/dal/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

// Init 初始化数据库连接
func Init(cfg *config.Config) error {
	// 首先尝试创建数据库
	if err := createDatabaseIfNotExists(cfg); err != nil {
		return fmt.Errorf("create database failed: %v", err)
	}

	// 连接到指定数据库
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.Database.Username,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.DBName,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("connect database failed: %v", err)
	}

	// 自动迁移表结构
	if err := autoMigrate(db); err != nil {
		return fmt.Errorf("auto migrate failed: %v", err)
	}

	// 初始化基础数据
	if err := initBaseData(db); err != nil {
		return fmt.Errorf("init base data failed: %v", err)
	}

	DB = db
	return nil
}

// createDatabaseIfNotExists 如果数据库不存在则创建
func createDatabaseIfNotExists(cfg *config.Config) error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/",
		cfg.Database.Username,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return err
	}
	defer db.Close()

	// 创建数据库
	createSQL := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;",
		cfg.Database.DBName)

	_, err = db.Exec(createSQL)
	return err
}

// autoMigrate 自动迁移表结构
func autoMigrate(db *gorm.DB) error {
	// 在这里添加需要自动迁移的模型
	return db.AutoMigrate(
		&model.User{},
		// 后续其他模型都在这里添加
	)
}

// initBaseData 初始化基础数据
func initBaseData(db *gorm.DB) error {
	// 检查是否已经有管理员用户
	var count int64
	db.Model(&model.User{}).Where("username = ?", "admin").Count(&count)
	if count == 0 {
		// 创建默认管理员账户
		admin := &model.User{
			Username: "admin",
			Password: "$2a$10$ZYkfwgz6hnxZ0T0TrYWPU.VyoG9w4xgAG2gjqfvhqrFX1mGGtNp.e", // 密码: admin123
			Email:    "admin@example.com",
			Status:   "active",
		}
		if err := db.Create(admin).Error; err != nil {
			return fmt.Errorf("create admin user failed: %v", err)
		}
	}
	return nil
}

// Close 关闭数据库连接
func Close() error {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}
