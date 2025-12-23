// 02-model-define: GORM 模型定义
//
// 📌 gorm.Model 包含:
//   - ID        uint           `gorm:"primarykey"`
//   - CreatedAt time.Time
//   - UpdatedAt time.Time
//   - DeletedAt gorm.DeletedAt `gorm:"index"`
//
// 📌 常用 Tag:
//   - column: 指定列名
//   - type: 指定列类型
//   - size: 指定大小
//   - primaryKey: 主键
//   - unique: 唯一
//   - uniqueIndex: 唯一索引
//   - index: 索引
//   - not null: 非空
//   - default: 默认值
//   - "-": 忽略该字段
package main

import (
	"log"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	gorm.Model        // 嵌入 ID, CreatedAt, UpdatedAt, DeletedAt
	Username   string `gorm:"uniqueIndex;size:50;not null"` // 唯一索引
	Password   string `gorm:"size:100;not null"`
	Email      string `gorm:"uniqueIndex;size:100"`
	Nickname   string `gorm:"size:50"`
	Age        int    `gorm:"default:0"`
	Status     int    `gorm:"default:1;comment:1-正常 0-禁用"`
}

// Product 商品模型（自定义主键）
type Product struct {
	ID          uint    `gorm:"primaryKey;autoIncrement"`
	Code        string  `gorm:"uniqueIndex;size:50"`
	Name        string  `gorm:"size:200;not null"`
	Price       float64 `gorm:"type:decimal(10,2)"`
	Stock       int     `gorm:"default:0"`
	Description string  `gorm:"type:text"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Order 订单模型（自定义表名）
type Order struct {
	ID        uint    `gorm:"primaryKey"`
	OrderNo   string  `gorm:"uniqueIndex;size:50"`
	UserID    uint    `gorm:"index"`
	Amount    float64 `gorm:"type:decimal(10,2)"`
	Status    int     `gorm:"default:0;comment:0-待支付 1-已支付 2-已完成"`
	CreatedAt time.Time
}

// TableName 自定义表名
func (Order) TableName() string {
	return "t_orders"
}

func main() {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("连接数据库失败: %v", err)
	}

	// 自动迁移（创建表）
	err = db.AutoMigrate(&User{}, &Product{}, &Order{})
	if err != nil {
		log.Fatalf("自动迁移失败: %v", err)
	}

	log.Println("自动迁移成功!")
	log.Println("创建的表: users, products, t_orders")
}
