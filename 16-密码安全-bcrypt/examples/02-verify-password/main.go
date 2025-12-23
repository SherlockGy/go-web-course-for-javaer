// 02-verify-password: bcrypt 密码验证
//
// 📌 验证最佳实践:
//   - 使用 bcrypt.CompareHashAndPassword
//   - 时间恒定比较，防止时序攻击
//   - 不要自己实现比较逻辑
//
// 📌 错误处理:
//   - bcrypt.ErrMismatchedHashAndPassword: 密码错误
//   - bcrypt.ErrHashTooShort: 哈希格式错误
package main

import (
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	// 模拟数据库存储的哈希
	storedHash := "$2a$10$N9qo8uLOickgx2ZMRZoMye3VtB3/BfJ.bSNKVGrk1Ie9Oa3Ghj8K6"
	correctPassword := "password123"
	wrongPassword := "wrongpassword"

	// ==================== 正确密码验证 ====================
	fmt.Println("=== 密码验证 ===")

	if CheckPassword(correctPassword, storedHash) {
		fmt.Println("✓ 正确密码验证通过")
	}

	if !CheckPassword(wrongPassword, storedHash) {
		fmt.Println("✗ 错误密码验证失败（预期）")
	}

	// ==================== 详细错误处理 ====================
	fmt.Println("\n=== 详细错误处理 ===")

	// 密码错误
	err := VerifyPassword(wrongPassword, storedHash)
	if err != nil {
		fmt.Printf("验证失败: %v\n", err)
	}

	// 哈希格式错误
	err = VerifyPassword("test", "invalid-hash")
	if err != nil {
		fmt.Printf("哈希无效: %v\n", err)
	}

	// ==================== 完整流程演示 ====================
	fmt.Println("\n=== 完整流程演示 ===")

	password := "MySecurePass@123"

	// 注册时：生成哈希存储
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	fmt.Printf("注册时存储: %s\n", hash)

	// 登录时：验证密码
	err = bcrypt.CompareHashAndPassword(hash, []byte(password))
	if err == nil {
		fmt.Println("登录成功!")
	}

	// ==================== 时序攻击防护说明 ====================
	fmt.Println("\n=== 安全说明 ===")
	fmt.Println("bcrypt.CompareHashAndPassword 使用恒定时间比较")
	fmt.Println("无论密码在哪个字符出错，比较时间都相同")
	fmt.Println("这防止了通过响应时间猜测密码的攻击")
}

// CheckPassword 验证密码（简化版）
// 📌 返回 bool，适合简单场景
func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// VerifyPassword 验证密码（带错误信息）
// 📌 返回具体错误，便于调试和日志
func VerifyPassword(password, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return errors.New("密码错误")
		}
		return fmt.Errorf("验证失败: %w", err)
	}
	return nil
}

// PasswordMatcher 密码验证器接口
// 📌 面向接口编程，便于测试和替换实现
type PasswordMatcher interface {
	Hash(password string) (string, error)
	Verify(password, hash string) error
}

// BcryptMatcher bcrypt 实现
type BcryptMatcher struct {
	cost int
}

func NewBcryptMatcher(cost int) *BcryptMatcher {
	if cost < bcrypt.MinCost || cost > bcrypt.MaxCost {
		cost = bcrypt.DefaultCost
	}
	return &BcryptMatcher{cost: cost}
}

func (m *BcryptMatcher) Hash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), m.cost)
	return string(hash), err
}

func (m *BcryptMatcher) Verify(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
