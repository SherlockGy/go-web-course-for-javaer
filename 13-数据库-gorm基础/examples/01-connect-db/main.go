// 01-connect-db: GORM 数据库连接
//
// 📌 GORM 支持的数据库:
//   - SQLite (gorm.io/driver/sqlite)
//   - MySQL (gorm.io/driver/mysql)
//   - PostgreSQL (gorm.io/driver/postgres)
//   - SQL Server (gorm.io/driver/sqlserver)
//
// 📌 最佳实践:
//   - 使用连接池配置
//   - 设置日志级别
//   - 生产环境禁用 Debug 模式
package main

import (
	"log"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	// ==================== SQLite 连接 ====================
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // 打印 SQL
	})
	if err != nil {
		log.Fatalf("连接数据库失败: %v", err)
	}

	log.Println("数据库连接成功")

	// 获取底层 *sql.DB 以配置连接池
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("获取 sql.DB 失败: %v", err)
	}

	// 连接池配置
	sqlDB.SetMaxIdleConns(10)           // 最大空闲连接数
	sqlDB.SetMaxOpenConns(100)          // 最大打开连接数
	sqlDB.SetConnMaxLifetime(time.Hour) // 连接最大生命周期

	log.Println("连接池配置完成")

	// 关闭连接（通常在程序结束时）
	defer sqlDB.Close()

	// ==================== MySQL 连接示例 ====================
	/*
		dsn := "user:password@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
		db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	*/

	// ==================== PostgreSQL 连接示例 ====================
	/*
		dsn := "host=localhost user=gorm password=gorm dbname=gorm port=5432 sslmode=disable"
		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	*/

	log.Println("数据库操作完成")
}
