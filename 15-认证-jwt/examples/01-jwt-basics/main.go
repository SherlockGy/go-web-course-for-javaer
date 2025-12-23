// 01-jwt-basics: JWT 基础用法
//
// 📌 JWT 最佳实践:
//   - 使用 golang-jwt/jwt/v5（官方维护版本）
//   - Token 有效期不宜过长（access: 15min-2h, refresh: 7-30d）
//   - 密钥至少 256 位（32 字节）
//   - 生产环境密钥从环境变量/配置中心读取
//   - 敏感信息不要放入 payload（可被 Base64 解码）
//
// 📌 与 Java 对比:
//   - Java: io.jsonwebtoken:jjwt（Builder 模式）
//   - Go: golang-jwt/jwt（函数式，更简洁）
package main

import (
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// 📌 密钥管理最佳实践: 生产环境应从配置或环境变量读取
var jwtSecret = []byte("your-256-bit-secret-key-here!!!")

// Claims 自定义声明
// 📌 嵌入 jwt.RegisteredClaims 获取标准字段
type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

func main() {
	// ==================== 生成 Token ====================
	fmt.Println("=== 生成 JWT Token ===")

	token, err := GenerateToken(1, "tom", "admin")
	if err != nil {
		log.Fatalf("生成 token 失败: %v", err)
	}
	fmt.Printf("Token: %s\n\n", token)

	// ==================== 解析 Token ====================
	fmt.Println("=== 解析 JWT Token ===")

	claims, err := ParseToken(token)
	if err != nil {
		log.Fatalf("解析 token 失败: %v", err)
	}
	fmt.Printf("UserID: %d\n", claims.UserID)
	fmt.Printf("Username: %s\n", claims.Username)
	fmt.Printf("Role: %s\n", claims.Role)
	fmt.Printf("过期时间: %v\n", claims.ExpiresAt.Time)

	// ==================== 验证过期 Token ====================
	fmt.Println("\n=== 测试过期 Token ===")

	expiredToken, _ := GenerateExpiredToken(1, "tom")
	_, err = ParseToken(expiredToken)
	if err != nil {
		fmt.Printf("过期 Token 验证失败（预期）: %v\n", err)
	}

	// ==================== 验证无效 Token ====================
	fmt.Println("\n=== 测试无效 Token ===")

	_, err = ParseToken("invalid.token.here")
	if err != nil {
		fmt.Printf("无效 Token 验证失败（预期）: %v\n", err)
	}
}

// GenerateToken 生成 JWT Token
// 📌 最佳实践: 设置合理的过期时间
func GenerateToken(userID uint, username, role string) (string, error) {
	claims := Claims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			// 签发者
			Issuer: "go-web-tutorial",
			// 主题
			Subject: username,
			// 过期时间（2小时）
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(2 * time.Hour)),
			// 签发时间
			IssuedAt: jwt.NewNumericDate(time.Now()),
			// 生效时间
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	// 📌 HS256 是对称加密，适合单服务
	// 📌 RS256 是非对称加密，适合微服务（公钥验证）
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// ParseToken 解析并验证 JWT Token
// 📌 最佳实践: 验证签名方法防止算法替换攻击
func ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// 📌 重要: 验证签名算法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

// GenerateExpiredToken 生成已过期的 Token（仅用于测试）
func GenerateExpiredToken(userID uint, username string) (string, error) {
	claims := Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)), // 已过期
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}
