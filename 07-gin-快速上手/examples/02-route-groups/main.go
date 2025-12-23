// 02-route-groups: 路由分组
//
// 📌 路由分组的好处:
//   - 统一 URL 前缀
//   - 分组应用中间件
//   - 代码组织更清晰
//
// 📌 最佳实践:
//   - 按版本分组: /api/v1, /api/v2
//   - 按功能分组: /auth, /users, /admin
//   - 嵌套分组实现层级
package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// 公开路由
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// API v1 分组
	v1 := r.Group("/api/v1")
	{
		// /api/v1/users
		v1.GET("/users", listUsersV1)
		v1.POST("/users", createUserV1)
	}

	// API v2 分组（可能有不同的响应格式）
	v2 := r.Group("/api/v2")
	{
		// /api/v2/users
		v2.GET("/users", listUsersV2)
	}

	// 认证相关路由
	auth := r.Group("/auth")
	{
		auth.POST("/login", loginHandler)
		auth.POST("/register", registerHandler)
		auth.POST("/logout", logoutHandler)
	}

	// 嵌套分组：需要认证的 API
	api := r.Group("/api")
	{
		// 用户模块
		users := api.Group("/users")
		{
			users.GET("", listUsersV1)        // GET /api/users
			users.GET("/:id", getUserHandler) // GET /api/users/:id
			users.POST("", createUserV1)      // POST /api/users
		}

		// 订单模块
		orders := api.Group("/orders")
		{
			orders.GET("", listOrdersHandler)   // GET /api/orders
			orders.GET("/:id", getOrderHandler) // GET /api/orders/:id
		}
	}

	r.Run(":8080")
}

// ==================== Handler 函数 ====================

func listUsersV1(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"version": "v1",
		"users":   []string{"tom", "jerry"},
	})
}

func listUsersV2(c *gin.Context) {
	// V2 版本可能有不同的响应格式
	c.JSON(http.StatusOK, gin.H{
		"version": "v2",
		"data": gin.H{
			"items": []gin.H{
				{"id": 1, "name": "tom"},
				{"id": 2, "name": "jerry"},
			},
			"total": 2,
		},
	})
}

func createUserV1(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{"message": "用户创建成功"})
}

func getUserHandler(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{"user_id": id})
}

func loginHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "登录成功"})
}

func registerHandler(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{"message": "注册成功"})
}

func logoutHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "登出成功"})
}

func listOrdersHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"orders": []string{}})
}

func getOrderHandler(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{"order_id": id})
}
