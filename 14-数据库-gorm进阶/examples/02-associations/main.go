// 02-associations: GORM 关联关系
//
// 📌 关联类型:
//   - BelongsTo: 属于（外键在当前表）
//   - HasOne: 一对一（外键在关联表）
//   - HasMany: 一对多（外键在关联表）
//   - Many2Many: 多对多（中间表）
package main

import (
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// User 用户
type User struct {
	gorm.Model
	Username string  `gorm:"uniqueIndex;size:50"`
	Profile  Profile // HasOne
	Orders   []Order // HasMany
}

// Profile 用户资料（一对一）
type Profile struct {
	gorm.Model
	UserID uint `gorm:"uniqueIndex"` // 外键
	Avatar string
	Bio    string
}

// Order 订单（多对一）
type Order struct {
	gorm.Model
	OrderNo string `gorm:"uniqueIndex;size:50"`
	UserID  uint   `gorm:"index"` // 外键
	User    User   // BelongsTo
	Amount  float64
}

func main() {
	db, _ := gorm.Open(sqlite.Open("test.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	// 自动迁移
	db.AutoMigrate(&User{}, &Profile{}, &Order{})

	// 清空测试数据
	db.Exec("DELETE FROM orders")
	db.Exec("DELETE FROM profiles")
	db.Exec("DELETE FROM users")

	// ==================== 创建关联数据 ====================
	fmt.Println("=== 创建关联数据 ===")

	// 创建用户和资料
	user := User{
		Username: "tom",
		Profile: Profile{
			Avatar: "/avatar/tom.png",
			Bio:    "Hello, I'm Tom!",
		},
		Orders: []Order{
			{OrderNo: "ORD001", Amount: 99.9},
			{OrderNo: "ORD002", Amount: 199.9},
		},
	}
	db.Create(&user)
	fmt.Printf("创建用户: ID=%d\n", user.ID)

	// ==================== 预加载查询 ====================
	fmt.Println("\n=== 预加载查询 ===")

	var foundUser User

	// Preload 预加载关联
	db.Preload("Profile").Preload("Orders").First(&foundUser, user.ID)
	fmt.Printf("用户: %s\n", foundUser.Username)
	fmt.Printf("头像: %s\n", foundUser.Profile.Avatar)
	fmt.Printf("订单数: %d\n", len(foundUser.Orders))

	// 条件预加载
	db.Preload("Orders", "amount > ?", 100).First(&foundUser, user.ID)
	fmt.Printf("大额订单数: %d\n", len(foundUser.Orders))

	// ==================== 关联查询 ====================
	fmt.Println("\n=== 关联查询 ===")

	// 通过订单找用户 (BelongsTo)
	var order Order
	db.Preload("User").First(&order, "order_no = ?", "ORD001")
	fmt.Printf("订单 %s 属于用户: %s\n", order.OrderNo, order.User.Username)

	// ==================== 关联操作 ====================
	fmt.Println("\n=== 关联操作 ===")

	// 添加关联
	newOrder := Order{OrderNo: "ORD003", Amount: 299.9}
	db.Model(&user).Association("Orders").Append(&newOrder)
	fmt.Printf("添加订单后: %d 条\n", db.Model(&user).Association("Orders").Count())

	// 替换关联
	// db.Model(&user).Association("Profile").Replace(&Profile{Avatar: "/new.png"})

	// 删除关联（只删除引用，不删除记录）
	// db.Model(&user).Association("Orders").Delete(&order)

	// 清空关联
	// db.Model(&user).Association("Orders").Clear()
}
