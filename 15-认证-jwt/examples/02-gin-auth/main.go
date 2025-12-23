// 02-gin-auth: Gin + JWT 完整认证示例
//
// 📌 认证流程最佳实践:
//   1. 用户登录 -> 返回 access_token + refresh_token
//   2. 请求携带 Authorization: Bearer <token>
//   3. 中间件验证 token -> 设置用户信息到 Context
//   4. access_token 过期 -> 用 refresh_token 刷新
//
// 📌 与 Java 对比:
//   - Java Spring: JwtAuthenticationFilter extends OncePerRequestFilter
//   - Go Gin: 中间件函数，更轻量直接
package main

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("your-256-bit-secret-key-here!!!")

// Claims 自定义声明
type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// TokenResponse 登录响应
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"` // 秒
}

func main() {
	r := gin.Default()

	// 公开路由
	r.POST("/login", loginHandler)
	r.POST("/refresh", refreshHandler)

	// 受保护路由
	protected := r.Group("/api")
	protected.Use(JWTAuthMiddleware())
	{
		protected.GET("/profile", profileHandler)
		protected.GET("/admin", RoleMiddleware("admin"), adminHandler)
	}

	r.Run(":8080")
}

// ==================== Handlers ====================

func loginHandler(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	// 📌 实际应查询数据库并验证密码（bcrypt）
	if req.Username != "admin" || req.Password != "123456" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	// 生成 Token
	accessToken, _ := generateToken(1, req.Username, "admin", 2*time.Hour)
	refreshToken, _ := generateToken(1, req.Username, "admin", 7*24*time.Hour)

	c.JSON(http.StatusOK, TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    7200, // 2小时
	})
}

func refreshHandler(c *gin.Context) {
	refreshToken := c.GetHeader("X-Refresh-Token")
	if refreshToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "缺少 refresh token"})
		return
	}

	claims, err := parseToken(refreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "refresh token 无效"})
		return
	}

	// 生成新的 access token
	newAccessToken, _ := generateToken(claims.UserID, claims.Username, claims.Role, 2*time.Hour)

	c.JSON(http.StatusOK, gin.H{
		"access_token": newAccessToken,
		"expires_in":   7200,
	})
}

func profileHandler(c *gin.Context) {
	// 📌 从 Context 获取用户信息（中间件已设置）
	userID, _ := c.Get("userID")
	username, _ := c.Get("username")
	role, _ := c.Get("role")

	c.JSON(http.StatusOK, gin.H{
		"user_id":  userID,
		"username": username,
		"role":     role,
		"message":  "个人信息获取成功",
	})
}

func adminHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "欢迎管理员!",
	})
}

// ==================== Middleware ====================

// JWTAuthMiddleware JWT 认证中间件
// 📌 最佳实践: 统一处理认证逻辑，业务代码无需关心
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 📌 标准格式: Authorization: Bearer <token>
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "缺少认证头"})
			c.Abort()
			return
		}

		// 解析 Bearer Token
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "认证头格式错误"})
			c.Abort()
			return
		}

		// 验证 Token
		claims, err := parseToken(parts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token 无效或已过期"})
			c.Abort()
			return
		}

		// 📌 将用户信息存入 Context，后续 handler 可用
		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)

		c.Next()
	}
}

// RoleMiddleware 角色验证中间件
// 📌 与 Java Spring @PreAuthorize("hasRole('ADMIN')") 类似
func RoleMiddleware(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists || role != requiredRole {
			c.JSON(http.StatusForbidden, gin.H{"error": "权限不足"})
			c.Abort()
			return
		}
		c.Next()
	}
}

// ==================== JWT Utils ====================

func generateToken(userID uint, username, role string, duration time.Duration) (string, error) {
	claims := Claims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "go-web-tutorial",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func parseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, jwt.ErrSignatureInvalid
}
