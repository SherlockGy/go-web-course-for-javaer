// 01-advanced-query: GORM 高级查询
//
// 📌 查询方法:
//   - Where: 条件查询
//   - Select: 选择字段
//   - Order: 排序
//   - Limit/Offset: 分页
//   - Group/Having: 分组
package main

import (
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type User struct {
	gorm.Model
	Username string `gorm:"uniqueIndex;size:50"`
	Email    string `gorm:"size:100"`
	Age      int
	Status   int `gorm:"default:1"`
}

func main() {
	db, _ := gorm.Open(sqlite.Open("test.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	db.AutoMigrate(&User{})

	// 清空并插入测试数据
	db.Exec("DELETE FROM users")
	testUsers := []User{
		{Username: "tom", Email: "tom@example.com", Age: 25, Status: 1},
		{Username: "jerry", Email: "jerry@example.com", Age: 20, Status: 1},
		{Username: "alice", Email: "alice@example.com", Age: 30, Status: 0},
		{Username: "bob", Email: "bob@example.com", Age: 35, Status: 1},
		{Username: "carol", Email: "carol@example.com", Age: 28, Status: 1},
	}
	db.Create(&testUsers)

	// ==================== Where 条件查询 ====================
	fmt.Println("=== Where 条件查询 ===")

	var users []User

	// 简单条件
	db.Where("age > ?", 25).Find(&users)
	fmt.Printf("age > 25: %d 条\n", len(users))

	// 多条件 AND
	db.Where("age > ? AND status = ?", 20, 1).Find(&users)
	fmt.Printf("age > 20 AND status = 1: %d 条\n", len(users))

	// Or 条件
	db.Where("age < ?", 22).Or("age > ?", 30).Find(&users)
	fmt.Printf("age < 22 OR age > 30: %d 条\n", len(users))

	// IN 查询
	db.Where("username IN ?", []string{"tom", "jerry"}).Find(&users)
	fmt.Printf("IN ('tom', 'jerry'): %d 条\n", len(users))

	// LIKE 查询
	db.Where("email LIKE ?", "%example.com").Find(&users)
	fmt.Printf("LIKE %%example.com: %d 条\n", len(users))

	// 结构体条件（非零值）
	db.Where(&User{Status: 1}).Find(&users)
	fmt.Printf("status = 1: %d 条\n", len(users))

	// Map 条件
	db.Where(map[string]interface{}{"status": 1, "age": 25}).Find(&users)
	fmt.Printf("status=1 AND age=25: %d 条\n", len(users))

	// ==================== Select 选择字段 ====================
	fmt.Println("\n=== Select 选择字段 ===")

	var usernames []string
	db.Model(&User{}).Select("username").Find(&usernames)
	fmt.Printf("用户名列表: %v\n", usernames)

	type UserSimple struct {
		Username string
		Email    string
	}
	var simpleUsers []UserSimple
	db.Model(&User{}).Select("username", "email").Find(&simpleUsers)
	fmt.Printf("简化用户: %+v\n", simpleUsers[0])

	// ==================== Order 排序 ====================
	fmt.Println("\n=== Order 排序 ===")

	db.Order("age desc").Find(&users)
	fmt.Printf("按年龄降序: %s (age=%d)\n", users[0].Username, users[0].Age)

	db.Order("status desc, age asc").Find(&users)
	fmt.Printf("多字段排序: %s\n", users[0].Username)

	// ==================== Limit/Offset 分页 ====================
	fmt.Println("\n=== Limit/Offset 分页 ===")

	page := 1
	pageSize := 2
	offset := (page - 1) * pageSize

	db.Offset(offset).Limit(pageSize).Find(&users)
	fmt.Printf("第 %d 页 (每页 %d 条): %d 条\n", page, pageSize, len(users))

	// 获取总数
	var total int64
	db.Model(&User{}).Count(&total)
	fmt.Printf("总记录数: %d\n", total)

	// ==================== 原生 SQL ====================
	fmt.Println("\n=== 原生 SQL ===")

	db.Raw("SELECT * FROM users WHERE age > ?", 25).Scan(&users)
	fmt.Printf("原生查询: %d 条\n", len(users))
}
