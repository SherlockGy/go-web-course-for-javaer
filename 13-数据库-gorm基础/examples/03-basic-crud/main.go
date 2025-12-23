// 03-basic-crud: GORM 增删改查
//
// 📌 CRUD 操作:
//   - Create: 创建记录
//   - First/Find: 查询记录
//   - Save/Updates: 更新记录
//   - Delete: 删除记录
//
// 📌 最佳实践:
//   - 检查操作结果的 Error 和 RowsAffected
//   - 使用 First 返回单条，Find 返回多条
//   - 更新使用 Updates（只更新非零值）或 Select
package main

import (
	"errors"
	"fmt"
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type User struct {
	gorm.Model
	Username string `gorm:"uniqueIndex;size:50"`
	Email    string `gorm:"size:100"`
	Age      int    `gorm:"default:0"`
}

func main() {
	db, _ := gorm.Open(sqlite.Open("test.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	db.AutoMigrate(&User{})

	// 清空测试数据
	db.Exec("DELETE FROM users")

	// ==================== Create 创建 ====================
	fmt.Println("=== 创建记录 ===")

	// 创建单条
	user := User{Username: "tom", Email: "tom@example.com", Age: 25}
	result := db.Create(&user)
	if result.Error != nil {
		log.Fatalf("创建失败: %v", result.Error)
	}
	fmt.Printf("创建成功: ID=%d, 影响行数=%d\n", user.ID, result.RowsAffected)

	// 批量创建
	users := []User{
		{Username: "jerry", Email: "jerry@example.com", Age: 20},
		{Username: "alice", Email: "alice@example.com", Age: 30},
	}
	db.Create(&users)
	fmt.Printf("批量创建: %d 条记录\n", len(users))

	// ==================== Read 查询 ====================
	fmt.Println("\n=== 查询记录 ===")

	// First: 获取第一条
	var firstUser User
	db.First(&firstUser) // SELECT * FROM users ORDER BY id LIMIT 1
	fmt.Printf("First: %s\n", firstUser.Username)

	// First by ID
	var userByID User
	db.First(&userByID, 1) // SELECT * FROM users WHERE id = 1
	fmt.Printf("By ID: %s\n", userByID.Username)

	// First by condition
	var userByName User
	db.First(&userByName, "username = ?", "tom")
	fmt.Printf("By Name: %s\n", userByName.Username)

	// Find: 获取多条
	var allUsers []User
	db.Find(&allUsers)
	fmt.Printf("All Users: %d 条\n", len(allUsers))

	// Find with condition
	var adults []User
	db.Where("age >= ?", 25).Find(&adults)
	fmt.Printf("Adults (age>=25): %d 条\n", len(adults))

	// 处理记录不存在的情况
	var notFound User
	result = db.First(&notFound, "username = ?", "nobody")
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		fmt.Println("记录不存在")
	}

	// ==================== Update 更新 ====================
	fmt.Println("\n=== 更新记录 ===")

	// Save: 保存所有字段（包括零值）
	user.Age = 26
	db.Save(&user)
	fmt.Printf("Save 后: Age=%d\n", user.Age)

	// Updates: 更新指定字段（忽略零值）
	db.Model(&user).Updates(User{Age: 27, Email: "tom.new@example.com"})

	// Updates with map (可以更新为零值)
	db.Model(&user).Updates(map[string]interface{}{"age": 0})

	// 条件更新
	db.Model(&User{}).Where("age < ?", 25).Update("age", 25)

	// ==================== Delete 删除 ====================
	fmt.Println("\n=== 删除记录 ===")

	// 软删除（设置 deleted_at）
	db.Delete(&user) // user.DeletedAt 被设置
	fmt.Printf("软删除: ID=%d\n", user.ID)

	// 查询被软删除的记录
	var deletedUser User
	db.Unscoped().First(&deletedUser, user.ID)
	fmt.Printf("软删除记录仍存在: %s\n", deletedUser.Username)

	// 永久删除
	db.Unscoped().Delete(&deletedUser)
	fmt.Println("永久删除完成")

	// 条件删除
	db.Where("age < ?", 18).Delete(&User{})
}
