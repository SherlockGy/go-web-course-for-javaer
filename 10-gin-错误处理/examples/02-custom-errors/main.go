// 02-custom-errors: 自定义错误类型
//
// 📌 自定义错误的好处:
//   - 携带更多信息（错误码、字段等）
//   - 统一错误处理
//   - 便于错误分类和国际化
//
// 📌 最佳实践:
//   - 定义业务错误码
//   - 实现 error 接口
//   - 使用 errors.Is/As 判断类型
package main

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ==================== 错误定义 ====================

// AppError 应用错误
type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Err     error  `json:"-"` // 原始错误，不序列化
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return e.Message + ": " + e.Err.Error()
	}
	return e.Message
}

// Unwrap 实现错误链
func (e *AppError) Unwrap() error {
	return e.Err
}

// 预定义错误
var (
	ErrNotFound     = &AppError{Code: 404, Message: "资源不存在"}
	ErrUnauthorized = &AppError{Code: 401, Message: "未授权"}
	ErrForbidden    = &AppError{Code: 403, Message: "禁止访问"}
	ErrBadRequest   = &AppError{Code: 400, Message: "请求参数错误"}
	ErrInternal     = &AppError{Code: 500, Message: "服务器内部错误"}
)

// 业务错误
var (
	ErrUserNotFound  = &AppError{Code: 1001, Message: "用户不存在"}
	ErrUserExists    = &AppError{Code: 1002, Message: "用户已存在"}
	ErrPasswordWrong = &AppError{Code: 1003, Message: "密码错误"}
	ErrTokenInvalid  = &AppError{Code: 1004, Message: "令牌无效"}
	ErrTokenExpired  = &AppError{Code: 1005, Message: "令牌已过期"}
)

// NewAppError 创建新错误
func NewAppError(code int, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// WrapError 包装错误
func WrapError(base *AppError, err error) *AppError {
	return &AppError{
		Code:    base.Code,
		Message: base.Message,
		Err:     err,
	}
}

// ==================== 使用示例 ====================

func main() {
	r := gin.Default()

	r.GET("/users/:id", getUser)
	r.POST("/login", login)

	r.Run(":8080")
}

func getUser(c *gin.Context) {
	id := c.Param("id")

	user, err := findUserByID(id)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    user,
	})
}

func login(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		handleError(c, WrapError(ErrBadRequest, err))
		return
	}

	if err := authenticate(req.Username, req.Password); err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "登录成功",
	})
}

// 统一错误处理
func handleError(c *gin.Context, err error) {
	var appErr *AppError
	if errors.As(err, &appErr) {
		// 确定 HTTP 状态码
		httpCode := http.StatusOK // 业务错误返回 200
		if appErr.Code >= 400 && appErr.Code < 600 {
			httpCode = appErr.Code // HTTP 错误用对应状态码
		}

		c.JSON(httpCode, gin.H{
			"code":    appErr.Code,
			"message": appErr.Message,
		})
		return
	}

	// 未知错误
	c.JSON(http.StatusInternalServerError, gin.H{
		"code":    500,
		"message": "服务器内部错误",
	})
}

// 模拟业务逻辑
func findUserByID(id string) (gin.H, error) {
	if id == "999" {
		return nil, ErrUserNotFound
	}
	return gin.H{"id": id, "username": "tom"}, nil
}

func authenticate(username, password string) error {
	if username != "admin" {
		return ErrUserNotFound
	}
	if password != "123456" {
		return ErrPasswordWrong
	}
	return nil
}
