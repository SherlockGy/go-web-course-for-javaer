// 02-custom-middleware: 自定义中间件
//
// 📌 中间件签名:
//   func() gin.HandlerFunc
//   func(c *gin.Context)
//
// 📌 关键方法:
//   - c.Next()  - 执行后续处理器
//   - c.Abort() - 终止后续处理
//   - c.Set() / c.Get() - 传递数据
//
// 📌 最佳实践:
//   - 单一职责：一个中间件只做一件事
//   - 注意执行顺序：先注册先执行
//   - 使用 defer 确保后处理执行
package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.New()

	// 全局中间件
	r.Use(RequestIDMiddleware())
	r.Use(LoggerMiddleware())
	r.Use(TimerMiddleware())
	r.Use(gin.Recovery())

	r.GET("/", func(c *gin.Context) {
		// 获取中间件设置的值
		requestID, _ := c.Get("request_id")
		c.JSON(http.StatusOK, gin.H{
			"message":    "Hello",
			"request_id": requestID,
		})
	})

	r.GET("/slow", func(c *gin.Context) {
		time.Sleep(100 * time.Millisecond)
		c.JSON(http.StatusOK, gin.H{"message": "慢接口"})
	})

	r.Run(":8080")
}

// RequestIDMiddleware 请求 ID 中间件
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 生成或获取请求 ID
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = generateRequestID()
		}

		// 存储到上下文
		c.Set("request_id", requestID)

		// 设置响应头
		c.Header("X-Request-ID", requestID)

		c.Next()
	}
}

// LoggerMiddleware 日志中间件
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 请求前
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method
		requestID, _ := c.Get("request_id")

		log.Printf("[%s] → %s %s", requestID, method, path)

		// 执行后续处理器
		c.Next()

		// 请求后（使用 defer 也可以）
		duration := time.Since(start)
		status := c.Writer.Status()

		log.Printf("[%s] ← %s %s %d %v",
			requestID, method, path, status, duration)
	}
}

// TimerMiddleware 计时中间件
func TimerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		// 如果请求耗时超过阈值，记录警告
		duration := time.Since(start)
		if duration > 50*time.Millisecond {
			log.Printf("⚠️ 慢请求: %s %s 耗时 %v",
				c.Request.Method, c.Request.URL.Path, duration)
		}
	}
}

// 生成简单的请求 ID
func generateRequestID() string {
	return time.Now().Format("20060102150405.000000")
}
