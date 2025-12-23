// 03-auth-middleware: 认证中间件
//
// 📌 认证中间件职责:
//   - 验证 Token
//   - 解析用户信息
//   - 存入上下文供后续使用
//   - 未认证时终止请求
//
// 📌 关键方法:
//   - c.Abort() - 终止后续处理
//   - c.AbortWithStatusJSON() - 终止并返回 JSON
package main

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// 公开路由
	r.POST("/login", loginHandler)

	// 受保护路由
	api := r.Group("/api")
	api.Use(AuthMiddleware()) // 应用认证中间件
	{
		api.GET("/me", getMeHandler)
		api.GET("/users", listUsersHandler)
	}

	// 管理员路由（双重检查）
	admin := r.Group("/admin")
	admin.Use(AuthMiddleware())
	admin.Use(AdminMiddleware())
	{
		admin.GET("/dashboard", adminDashboard)
	}

	r.Run(":8080")
}

// AuthMiddleware 认证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取 Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "未提供认证令牌",
			})
			return
		}

		// 解析 Bearer Token
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "认证格式错误，应为: Bearer <token>",
			})
			return
		}

		token := parts[1]

		// 验证 Token（简化示例）
		user, err := validateToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "无效的令牌",
			})
			return
		}

		// 将用户信息存入上下文
		c.Set("user_id", user.ID)
		c.Set("username", user.Username)
		c.Set("role", user.Role)

		c.Next()
	}
}

// AdminMiddleware 管理员权限中间件
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists || role != "admin" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"code":    403,
				"message": "需要管理员权限",
			})
			return
		}
		c.Next()
	}
}

// ==================== Handler ====================

func loginHandler(c *gin.Context) {
	// 简化示例：返回固定 Token
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "登录成功",
		"data": gin.H{
			"token": "valid-user-token",
		},
	})
}

func getMeHandler(c *gin.Context) {
	// 从上下文获取用户信息
	userID, _ := c.Get("user_id")
	username, _ := c.Get("username")
	role, _ := c.Get("role")

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"user_id":  userID,
			"username": username,
			"role":     role,
		},
	})
}

func listUsersHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    []string{"tom", "jerry"},
	})
}

func adminDashboard(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"dashboard": "管理员仪表盘",
		},
	})
}

// ==================== 辅助函数 ====================

type User struct {
	ID       int
	Username string
	Role     string
}

// validateToken 验证 Token（简化示例）
func validateToken(token string) (*User, error) {
	// 实际项目中应该使用 JWT 验证
	validTokens := map[string]*User{
		"valid-user-token":  {ID: 1, Username: "tom", Role: "user"},
		"valid-admin-token": {ID: 2, Username: "admin", Role: "admin"},
	}

	if user, ok := validTokens[token]; ok {
		return user, nil
	}
	return nil, gin.Error{Err: nil, Type: gin.ErrorTypePrivate}
}
