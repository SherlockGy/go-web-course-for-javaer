// 03-transactions: GORM 事务操作
//
// 📌 事务使用场景:
//   - 转账：扣款和入账必须同时成功
//   - 订单：创建订单和扣减库存必须原子
//   - 注册：创建用户和初始化配置必须同时
//
// 📌 事务方法:
//   - db.Transaction(func(tx *gorm.DB) error {})
//   - db.Begin() / tx.Commit() / tx.Rollback()
package main

import (
	"errors"
	"fmt"
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Account struct {
	gorm.Model
	UserID  uint    `gorm:"uniqueIndex"`
	Balance float64 `gorm:"type:decimal(10,2)"`
}

type TransferLog struct {
	gorm.Model
	FromUserID uint
	ToUserID   uint
	Amount     float64
	Status     string // success / failed
}

func main() {
	db, _ := gorm.Open(sqlite.Open("test.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	db.AutoMigrate(&Account{}, &TransferLog{})

	// 初始化测试数据
	db.Exec("DELETE FROM accounts")
	db.Exec("DELETE FROM transfer_logs")
	db.Create(&Account{UserID: 1, Balance: 1000})
	db.Create(&Account{UserID: 2, Balance: 500})

	// ==================== 方式1: Transaction 闭包（推荐）====================
	fmt.Println("=== Transaction 闭包方式 ===")

	err := Transfer(db, 1, 2, 100)
	if err != nil {
		fmt.Printf("转账失败: %v\n", err)
	} else {
		fmt.Println("转账成功!")
	}

	// 验证结果
	var acc1, acc2 Account
	db.First(&acc1, "user_id = ?", 1)
	db.First(&acc2, "user_id = ?", 2)
	fmt.Printf("用户1余额: %.2f, 用户2余额: %.2f\n", acc1.Balance, acc2.Balance)

	// 测试失败场景
	err = Transfer(db, 1, 2, 10000) // 余额不足
	if err != nil {
		fmt.Printf("转账失败（预期）: %v\n", err)
	}

	// ==================== 方式2: 手动控制 ====================
	fmt.Println("\n=== 手动控制方式 ===")

	tx := db.Begin()
	if tx.Error != nil {
		log.Fatalf("开启事务失败: %v", tx.Error)
	}

	// 执行操作
	if err := tx.Model(&Account{}).Where("user_id = ?", 1).
		Update("balance", gorm.Expr("balance - ?", 50)).Error; err != nil {
		tx.Rollback()
		fmt.Printf("扣款失败，回滚: %v\n", err)
		return
	}

	if err := tx.Model(&Account{}).Where("user_id = ?", 2).
		Update("balance", gorm.Expr("balance + ?", 50)).Error; err != nil {
		tx.Rollback()
		fmt.Printf("入账失败，回滚: %v\n", err)
		return
	}

	// 提交事务
	tx.Commit()
	fmt.Println("手动事务提交成功!")
}

// Transfer 转账（使用事务）
func Transfer(db *gorm.DB, fromUserID, toUserID uint, amount float64) error {
	return db.Transaction(func(tx *gorm.DB) error {
		// 1. 检查源账户余额
		var fromAccount Account
		if err := tx.First(&fromAccount, "user_id = ?", fromUserID).Error; err != nil {
			return fmt.Errorf("源账户不存在: %w", err)
		}

		if fromAccount.Balance < amount {
			return errors.New("余额不足")
		}

		// 2. 扣减源账户
		if err := tx.Model(&fromAccount).
			Update("balance", gorm.Expr("balance - ?", amount)).Error; err != nil {
			return err
		}

		// 3. 增加目标账户
		result := tx.Model(&Account{}).Where("user_id = ?", toUserID).
			Update("balance", gorm.Expr("balance + ?", amount))
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return errors.New("目标账户不存在")
		}

		// 4. 记录日志
		log := TransferLog{
			FromUserID: fromUserID,
			ToUserID:   toUserID,
			Amount:     amount,
			Status:     "success",
		}
		if err := tx.Create(&log).Error; err != nil {
			return err
		}

		// 返回 nil 提交事务
		return nil
	})
}
