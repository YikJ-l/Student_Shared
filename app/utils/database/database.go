package database

import (
	"fmt"
	"log"
	"os"
	"student_shared/app/model"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB 全局数据库连接实例
var DB *gorm.DB

// InitDB 初始化数据库连接
func InitDB() error {
	// 数据库配置
	dbUser := "root"           // 数据库用户名
	dbPass := "abc123456"      // 数据库密码
	dbHost := "localhost"      // 数据库主机
	dbPort := "3306"           // 数据库端口
	dbName := "shared_student" // 数据库名称

	// 构建DSN连接字符串
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPass, dbHost, dbPort, dbName)

	// 配置日志
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // 慢SQL阈值
			LogLevel:                  logger.Info, // 日志级别
			IgnoreRecordNotFoundError: true,        // 忽略记录未找到错误
			Colorful:                  true,        // 彩色输出
		},
	)

	// 打开数据库连接
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		return fmt.Errorf("连接数据库失败: %w", err)
	}

	// 获取底层SQL连接池
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("获取连接池失败: %w", err)
	}

	// 设置连接池参数
	sqlDB.SetMaxIdleConns(10)           // 最大空闲连接数
	sqlDB.SetMaxOpenConns(100)          // 最大打开连接数
	sqlDB.SetConnMaxLifetime(time.Hour) // 连接最大生命周期

	// 自动迁移数据库模型
	err = autoMigrate()
	if err != nil {
		return fmt.Errorf("自动迁移数据库模型失败: %w", err)
	}

	return nil
}

// autoMigrate 自动迁移数据库模型
func autoMigrate() error {
	// 自动迁移模型
	err := DB.AutoMigrate(
		&model.User{},
		&model.Course{},
		&model.UserCourse{},
		&model.Note{},
		&model.Comment{},
		&model.Favorite{},
		&model.NoteLike{},
		&model.CommentLike{},
	)
	return err
}
