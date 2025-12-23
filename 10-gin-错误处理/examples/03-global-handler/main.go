// 03-global-handler: 全局错误处理
//
// 📌 全局错误处理方案:
//   - 方案1: 中间件 + c.Error() 收集错误
//   - 方案2: 封装响应函数统一处理
//   - 方案3: 自定义 Handler 包装器
//
// 📌 最佳实践:
//   - 所有错误统一格式响应
//   - 记录错误日志
//   - 生产环境隐藏内部错误
package main

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ==================== 错误定义 ====================

type BizError struct {
	Code    int
	Message string
}

func (e *BizError) Error() string {
	return e.Message
}

var (
	ErrUserNotFound  = &BizError{Code: 1001, Message: "用户不存在"}
	ErrPasswordWrong = &BizError{Code: 1003, Message: "密码错误"}
)

// ==================== 方案1: 错误处理中间件 ====================

// ErrorHandler 全局错误处理中间件
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// 检查是否有错误
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err
			handleError(c, err)
		}
	}
}

func handleError(c *gin.Context, err error) {
	// 记录错误日志
	log.Printf("Error: %v", err)

	var bizErr *BizError
	if errors.As(err, &bizErr) {
		c.JSON(http.StatusOK, gin.H{
			"code":    bizErr.Code,
			"message": bizErr.Message,
		})
		return
	}

	// 未知错误
	c.JSON(http.StatusInternalServerError, gin.H{
		"code":    500,
		"message": "服务器内部错误",
	})
}

// ==================== 方案2: Handler 包装器 ====================

// HandlerFunc 自定义处理函数类型
type HandlerFunc func(c *gin.Context) error

// Wrap 包装 Handler，统一处理错误
func Wrap(fn HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := fn(c); err != nil {
			handleError(c, err)
		}
	}
}

// ==================== 使用示例 ====================

func main() {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// 方案1: 使用错误处理中间件
	r.Use(ErrorHandler())

	// 传统方式：使用 c.Error() 记录错误
	r.GET("/users/:id", func(c *gin.Context) {
		id := c.Param("id")
		if id == "999" {
			c.Error(ErrUserNotFound) // 记录错误
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"data": gin.H{"id": id, "username": "tom"},
		})
	})

	// 方案2: 使用 Wrap 包装
	r.GET("/v2/users/:id", Wrap(getUserV2))
	r.POST("/v2/login", Wrap(loginV2))

	log.Println("测试命令:")
	log.Println("  curl http://localhost:8080/users/1")
	log.Println("  curl http://localhost:8080/users/999")
	log.Println("  curl http://localhost:8080/v2/users/999")

	r.Run(":8080")
}

// 使用 Wrap 的 Handler，返回 error
func getUserV2(c *gin.Context) error {
	id := c.Param("id")
	if id == "999" {
		return ErrUserNotFound
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{"id": id, "username": "tom"},
	})
	return nil
}

func loginV2(c *gin.Context) error {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		return &BizError{Code: 400, Message: "参数错误: " + err.Error()}
	}

	if req.Username != "admin" {
		return ErrUserNotFound
	}
	if req.Password != "123456" {
		return ErrPasswordWrong
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "登录成功",
	})
	return nil
}
