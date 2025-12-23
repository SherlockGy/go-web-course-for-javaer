// 04-middleware-chain: 中间件链执行顺序
//
// 📌 执行顺序（洋葱模型）:
//   请求 → M1(前) → M2(前) → M3(前) → Handler → M3(后) → M2(后) → M1(后) → 响应
//
// 📌 关键理解:
//   - c.Next() 之前的代码：请求处理前执行
//   - c.Next() 之后的代码：请求处理后执行
//   - c.Abort() 会终止后续处理，但当前中间件的后续代码仍执行
package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.New()

	// 中间件按注册顺序执行
	r.Use(Middleware1())
	r.Use(Middleware2())
	r.Use(Middleware3())

	r.GET("/", func(c *gin.Context) {
		fmt.Println(">>> Handler 执行")
		c.JSON(http.StatusOK, gin.H{"message": "OK"})
	})

	// 测试 Abort
	r.GET("/abort", func(c *gin.Context) {
		fmt.Println(">>> Handler 执行（不应该看到这条）")
		c.JSON(http.StatusOK, gin.H{"message": "OK"})
	})

	fmt.Println("测试命令:")
	fmt.Println("  curl http://localhost:8080/")
	fmt.Println("  curl http://localhost:8080/abort?abort=true")

	r.Run(":8080")
}

func Middleware1() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("→ M1 进入")

		c.Next() // 执行下一个中间件

		fmt.Println("← M1 退出")
	}
}

func Middleware2() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("  → M2 进入")

		// 检查是否需要终止
		if c.Query("abort") == "true" {
			fmt.Println("  ✗ M2 终止请求")
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "请求被终止",
			})
			// 注意：Abort 后当前中间件的后续代码仍会执行
			fmt.Println("  ← M2 退出（Abort 后）")
			return
		}

		c.Next()

		fmt.Println("  ← M2 退出")
	}
}

func Middleware3() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("    → M3 进入")

		c.Next()

		fmt.Println("    ← M3 退出")
	}
}

/*
正常请求输出:
→ M1 进入
  → M2 进入
    → M3 进入
>>> Handler 执行
    ← M3 退出
  ← M2 退出
← M1 退出

Abort 请求输出:
→ M1 进入
  → M2 进入
  ✗ M2 终止请求
  ← M2 退出（Abort 后）
← M1 退出

注意：M3 和 Handler 都没有执行！
*/
