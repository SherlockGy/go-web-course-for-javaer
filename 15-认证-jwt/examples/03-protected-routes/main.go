// 03-protected-routes: 完整的路由保护示例
//
// 📌 路由设计最佳实践:
//   - 公开路由: /login, /register, /health
//   - 认证路由: /api/* (需要登录)
//   - 权限路由: /admin/* (需要特定角色)
//
// 📌 与 Java Spring Security 对比:
//   - Java: SecurityFilterChain + antMatchers().permitAll()
//   - Go: 路由分组 + 中间件，更直观
package main

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("your-256-bit-secret-key-here!!!")

type Claims struct {
	UserID      uint     `json:"user_id"`
	Username    string   `json:"username"`
	Roles       []string `json:"roles"`       // 支持多角色
	Permissions []string `json:"permissions"` // 细粒度权限
	jwt.RegisteredClaims
}

func main() {
	r := gin.Default()

	// ==================== 公开路由 ====================
	// 无需认证
	public := r.Group("/")
	{
		public.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		})
		public.POST("/login", loginHandler)
		public.POST("/register", registerHandler)
	}

	// ==================== 认证路由 ====================
	// 需要登录，任意角色
	api := r.Group("/api")
	api.Use(AuthRequired())
	{
		api.GET("/profile", profileHandler)
		api.PUT("/profile", updateProfileHandler)
		api.GET("/orders", ordersHandler)
	}

	// ==================== 管理员路由 ====================
	// 需要 admin 角色
	admin := r.Group("/admin")
	admin.Use(AuthRequired(), RequireRole("admin"))
	{
		admin.GET("/users", listUsersHandler)
		admin.DELETE("/users/:id", deleteUserHandler)
		admin.GET("/stats", statsHandler)
	}

	// ==================== 细粒度权限路由 ====================
	// 需要特定权限
	products := r.Group("/products")
	products.Use(AuthRequired())
	{
		// 所有登录用户可读
		products.GET("", listProductsHandler)
		products.GET("/:id", getProductHandler)

		// 需要特定权限
		products.POST("", RequirePermission("product:create"), createProductHandler)
		products.PUT("/:id", RequirePermission("product:update"), updateProductHandler)
		products.DELETE("/:id", RequirePermission("product:delete"), deleteProductHandler)
	}

	r.Run(":8080")
}

// ==================== Handlers ====================

func loginHandler(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	// 模拟不同用户
	var userID uint
	var roles []string
	var permissions []string

	switch req.Username {
	case "admin":
		userID = 1
		roles = []string{"admin", "user"}
		permissions = []string{"product:create", "product:update", "product:delete"}
	case "editor":
		userID = 2
		roles = []string{"editor", "user"}
		permissions = []string{"product:create", "product:update"}
	default:
		userID = 3
		roles = []string{"user"}
		permissions = []string{}
	}

	token, _ := generateToken(userID, req.Username, roles, permissions)
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func registerHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "注册成功"})
}

func profileHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"user_id":     c.GetUint("userID"),
		"username":    c.GetString("username"),
		"roles":       c.GetStringSlice("roles"),
		"permissions": c.GetStringSlice("permissions"),
	})
}

func updateProfileHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "更新成功"})
}

func ordersHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"orders": []string{"ORD001", "ORD002"}})
}

func listUsersHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"users": []string{"admin", "tom", "jerry"}})
}

func deleteUserHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "用户已删除: " + c.Param("id")})
}

func statsHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"total_users":  100,
		"total_orders": 500,
	})
}

func listProductsHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"products": []string{"iPhone", "MacBook"}})
}

func getProductHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"product": c.Param("id")})
}

func createProductHandler(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{"message": "产品创建成功"})
}

func updateProductHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "产品更新成功"})
}

func deleteProductHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "产品删除成功"})
}

// ==================== Middleware ====================

// AuthRequired 认证中间件
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "需要登录"})
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "认证格式错误"})
			c.Abort()
			return
		}

		claims, err := parseToken(parts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token 无效"})
			c.Abort()
			return
		}

		// 存入 Context
		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("roles", claims.Roles)
		c.Set("permissions", claims.Permissions)

		c.Next()
	}
}

// RequireRole 角色验证中间件
// 📌 支持多角色: RequireRole("admin", "superadmin")
func RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRoles := c.GetStringSlice("roles")

		for _, required := range roles {
			for _, userRole := range userRoles {
				if userRole == required {
					c.Next()
					return
				}
			}
		}

		c.JSON(http.StatusForbidden, gin.H{"error": "角色权限不足"})
		c.Abort()
	}
}

// RequirePermission 权限验证中间件
// 📌 与 Java @PreAuthorize("hasAuthority('product:create')") 类似
func RequirePermission(required string) gin.HandlerFunc {
	return func(c *gin.Context) {
		permissions := c.GetStringSlice("permissions")

		for _, perm := range permissions {
			if perm == required {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{
			"error":    "权限不足",
			"required": required,
		})
		c.Abort()
	}
}

// ==================== JWT Utils ====================

func generateToken(userID uint, username string, roles, permissions []string) (string, error) {
	claims := Claims{
		UserID:      userID,
		Username:    username,
		Roles:       roles,
		Permissions: permissions,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
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
