// 01-recovery: Recovery 中间件
//
// 📌 Recovery 中间件作用:
//   - 捕获 panic
//   - 返回 500 错误
//   - 记录错误日志
//   - 防止服务器崩溃
//
// 📌 最佳实践:
//   - 始终使用 Recovery 中间件
//   - 自定义 Recovery 返回统一格式
//   - panic 仅用于不可恢复的错误
package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.New()

	// 添加日志中间件
	r.Use(gin.Logger())

	// 自定义 Recovery 中间件
	r.Use(CustomRecovery())

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Hello"})
	})

	r.GET("/panic", func(c *gin.Context) {
		panic("发生了严重错误!")
	})

	r.GET("/nil-pointer", func(c *gin.Context) {
		var ptr *string
		log.Println(*ptr) // 空指针访问，会 panic
	})

	log.Println("测试命令:")
	log.Println("  curl http://localhost:8080/")
	log.Println("  curl http://localhost:8080/panic")
	log.Println("  curl http://localhost:8080/nil-pointer")

	r.Run(":8080")
}

// CustomRecovery 自定义恢复中间件
func CustomRecovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 记录错误日志
				log.Printf("Panic recovered: %v", err)

				// 返回统一格式的错误响应
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"code":    500,
					"message": "服务器内部错误",
					// 生产环境不要暴露具体错误
					// "error": fmt.Sprintf("%v", err),
				})
			}
		}()
		c.Next()
	}
}
