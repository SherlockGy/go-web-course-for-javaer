// 01-hello-gin: Gin 框架入门
//
// 📌 Gin vs 原生 net/http:
//   - 更简洁的路由定义
//   - 内置参数绑定和验证
//   - 中间件机制更完善
//   - 更好的性能
//
// 📌 gin.Default() vs gin.New():
//   - Default() = New() + Logger + Recovery
//   - 生产环境建议用 New() 自定义中间件
//
// 安装: go get github.com/gin-gonic/gin
// 运行: go run main.go
package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// 创建 Gin 引擎
	// gin.Default() 包含 Logger 和 Recovery 中间件
	r := gin.Default()

	// 基本路由
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello, Gin!")
	})

	// JSON 响应
	r.GET("/json", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "success",
			"data": gin.H{
				"name":    "Gin",
				"version": "1.9+",
			},
		})
	})

	// 路径参数
	r.GET("/users/:id", func(c *gin.Context) {
		id := c.Param("id")
		c.JSON(http.StatusOK, gin.H{
			"user_id": id,
		})
	})

	// 查询参数
	r.GET("/search", func(c *gin.Context) {
		q := c.Query("q")                   // 获取参数，不存在返回空
		page := c.DefaultQuery("page", "1") // 带默认值
		c.JSON(http.StatusOK, gin.H{
			"query": q,
			"page":  page,
		})
	})

	// 启动服务器
	r.Run(":8080") // 默认 :8080
}
