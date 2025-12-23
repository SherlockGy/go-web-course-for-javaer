// 04-gin-integration: Zap 与 Gin 集成
//
// 📌 集成要点:
//   - 替换 Gin 默认日志中间件
//   - 记录请求方法、路径、状态码、耗时
//   - 添加 request_id 追踪
package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var logger *zap.Logger

func main() {
	// 初始化日志
	logger, _ = zap.NewProduction()
	defer logger.Sync()

	// 禁用 Gin 默认日志
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()

	// 使用自定义中间件
	r.Use(ZapLogger(logger))
	r.Use(ZapRecovery(logger))

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Hello"})
	})

	r.GET("/users/:id", func(c *gin.Context) {
		id := c.Param("id")
		c.JSON(http.StatusOK, gin.H{"user_id": id})
	})

	r.GET("/error", func(c *gin.Context) {
		panic("测试 panic")
	})

	logger.Info("服务器启动", zap.String("addr", ":8080"))
	r.Run(":8080")
}

// ZapLogger Gin 日志中间件
func ZapLogger(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		// 生成请求 ID
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = generateRequestID()
		}
		c.Set("request_id", requestID)
		c.Header("X-Request-ID", requestID)

		// 处理请求
		c.Next()

		// 记录日志
		latency := time.Since(start)
		status := c.Writer.Status()

		fields := []zap.Field{
			zap.String("request_id", requestID),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.Int("status", status),
			zap.Duration("latency", latency),
			zap.String("client_ip", c.ClientIP()),
			zap.String("user_agent", c.Request.UserAgent()),
		}

		if len(c.Errors) > 0 {
			fields = append(fields, zap.String("errors", c.Errors.String()))
		}

		// 根据状态码选择日志级别
		switch {
		case status >= 500:
			logger.Error("服务器错误", fields...)
		case status >= 400:
			logger.Warn("客户端错误", fields...)
		default:
			logger.Info("请求完成", fields...)
		}
	}
}

// ZapRecovery Gin Recovery 中间件
func ZapRecovery(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				requestID, _ := c.Get("request_id")

				logger.Error("Panic 恢复",
					zap.Any("request_id", requestID),
					zap.Any("error", err),
					zap.String("path", c.Request.URL.Path),
					zap.Stack("stacktrace"),
				)

				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"code":    500,
					"message": "服务器内部错误",
				})
			}
		}()
		c.Next()
	}
}

func generateRequestID() string {
	return time.Now().Format("20060102150405.000000")
}
