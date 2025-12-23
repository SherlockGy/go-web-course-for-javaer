// 01-hash-password: bcrypt 密码哈希
//
// 📌 密码存储最佳实践:
//   - 永远不要明文存储密码
//   - 使用 bcrypt/argon2/scrypt 等慢哈希算法
//   - bcrypt 自动处理盐值，无需额外存储
//   - cost 参数决定计算强度（推荐 10-14）
//
// 📌 与 Java 对比:
//   - Java: BCryptPasswordEncoder (Spring Security)
//   - Go: golang.org/x/crypto/bcrypt（标准扩展库）
//
// 📌 为什么选择 bcrypt:
//   - 内置盐值：防止彩虹表攻击
//   - 可调 cost：随硬件升级提高安全性
//   - 故意慢：防止暴力破解
package main

import (
	"fmt"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	password := "MySecurePassword123!"

	// ==================== 基本哈希 ====================
	fmt.Println("=== 基本哈希 ===")

	hash, err := HashPassword(password)
	if err != nil {
		log.Fatalf("哈希失败: %v", err)
	}
	fmt.Printf("原始密码: %s\n", password)
	fmt.Printf("哈希结果: %s\n", hash)
	fmt.Printf("哈希长度: %d\n", len(hash))

	// ==================== 每次哈希结果不同 ====================
	fmt.Println("\n=== 每次哈希结果不同（因为盐值不同）===")

	for i := range 3 {
		h, _ := HashPassword(password)
		fmt.Printf("第%d次: %s\n", i+1, h)
	}

	// ==================== Cost 参数影响 ====================
	fmt.Println("\n=== Cost 参数对性能的影响 ===")

	costs := []int{10, 12, 14}
	for _, cost := range costs {
		start := time.Now()
		_, _ = bcrypt.GenerateFromPassword([]byte(password), cost)
		duration := time.Since(start)
		fmt.Printf("Cost=%d: %v\n", cost, duration)
	}

	// ==================== 哈希结构解析 ====================
	fmt.Println("\n=== bcrypt 哈希结构 ===")
	fmt.Println("格式: $2a$cost$salt(22字符)hash(31字符)")
	fmt.Printf("示例: %s\n", hash)
	fmt.Println("  - $2a$ : 版本标识")
	fmt.Println("  - 10   : cost 参数")
	fmt.Println("  - 前22位: 盐值 (Base64)")
	fmt.Println("  - 后31位: 哈希值 (Base64)")
}

// HashPassword 哈希密码
// 📌 最佳实践: cost 推荐 10-14，根据服务器性能调整
func HashPassword(password string) (string, error) {
	// bcrypt.DefaultCost = 10
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// HashPasswordWithCost 使用指定 cost 哈希密码
// 📌 cost 越高越安全，但耗时越长
// 📌 建议: 登录验证耗时控制在 100ms-500ms
func HashPasswordWithCost(password string, cost int) (string, error) {
	// cost 范围: 4-31，推荐 10-14
	if cost < bcrypt.MinCost || cost > bcrypt.MaxCost {
		cost = bcrypt.DefaultCost
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}
