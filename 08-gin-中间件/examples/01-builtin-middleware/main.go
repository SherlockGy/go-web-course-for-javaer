// 01-builtin-middleware: 内置中间件使用
//
// 📌 Gin 内置中间件:
//   - gin.Logger()   - 请求日志
//   - gin.Recovery() - panic 恢复
//   - gin.BasicAuth() - HTTP Basic 认证
//
// 📌 gin.Default() vs gin.New():
//   - Default() = New() + Logger() + Recovery()
//   - 生产环境建议用 New() 自定义
package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	// 生产环境设置
	// gin.SetMode(gin.ReleaseMode)

	// 方式1: 使用 Default（包含 Logger 和 Recovery）
	// r := gin.Default()

	// 方式2: 使用 New 自定义中间件
	r := gin.New()

	// 添加 Logger 中间件（自定义输出）
	r.Use(gin.LoggerWithWriter(os.Stdout))

	// 添加 Recovery 中间件（自定义错误处理）
	r.Use(gin.CustomRecovery(func(c *gin.Context, recovered any) {
		if err, ok := recovered.(string); ok {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "服务器内部错误",
				"error":   err,
			})
		}
		c.AbortWithStatus(http.StatusInternalServerError)
	}))

	// 测试路由
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Hello"})
	})

	// 测试 panic（会被 Recovery 捕获）
	r.GET("/panic", func(c *gin.Context) {
		panic("测试 panic!")
	})

	// 使用 BasicAuth
	authorized := r.Group("/admin", gin.BasicAuth(gin.Accounts{
		"admin": "admin123",
		"user":  "user123",
	}))
	{
		authorized.GET("/dashboard", func(c *gin.Context) {
			user := c.MustGet(gin.AuthUserKey).(string)
			c.JSON(http.StatusOK, gin.H{
				"message": "欢迎 " + user,
			})
		})
	}

	log.Println("测试命令:")
	log.Println("  curl http://localhost:8080/")
	log.Println("  curl http://localhost:8080/panic")
	log.Println("  curl -u admin:admin123 http://localhost:8080/admin/dashboard")

	r.Run(":8080")
}
