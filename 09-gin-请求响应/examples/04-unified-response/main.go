// 04-unified-response: 统一响应封装
//
// 📌 统一响应格式:
//   {
//     "code": 0,       // 业务状态码
//     "message": "",   // 消息
//     "data": {}       // 数据
//   }
//
// 📌 最佳实践:
//   - 封装响应函数，减少重复代码
//   - HTTP 状态码和业务状态码分离
//   - 成功统一返回 200，错误通过 code 区分
package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 统一响应结构
type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// PageData 分页数据
type PageData struct {
	Items    any   `json:"items"`
	Total    int64 `json:"total"`
	Page     int   `json:"page"`
	PageSize int   `json:"page_size"`
}

// 响应辅助函数
func Success(c *gin.Context, data any) {
	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "success",
		Data:    data,
	})
}

func SuccessMessage(c *gin.Context, message string) {
	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: message,
	})
}

func SuccessPage(c *gin.Context, items any, total int64, page, pageSize int) {
	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "success",
		Data: PageData{
			Items:    items,
			Total:    total,
			Page:     page,
			PageSize: pageSize,
		},
	})
}

func Created(c *gin.Context, data any) {
	c.JSON(http.StatusCreated, Response{
		Code:    0,
		Message: "created",
		Data:    data,
	})
}

func Error(c *gin.Context, httpCode int, bizCode int, message string) {
	c.JSON(httpCode, Response{
		Code:    bizCode,
		Message: message,
	})
}

func BadRequest(c *gin.Context, message string) {
	Error(c, http.StatusBadRequest, 400, message)
}

func Unauthorized(c *gin.Context, message string) {
	Error(c, http.StatusUnauthorized, 401, message)
}

func Forbidden(c *gin.Context, message string) {
	Error(c, http.StatusForbidden, 403, message)
}

func NotFound(c *gin.Context, message string) {
	Error(c, http.StatusNotFound, 404, message)
}

func ServerError(c *gin.Context, message string) {
	Error(c, http.StatusInternalServerError, 500, message)
}

// 业务错误码
const (
	ErrCodeUserNotFound  = 1001
	ErrCodeUserExists    = 1002
	ErrCodePasswordWrong = 1003
)

func BusinessError(c *gin.Context, code int, message string) {
	c.JSON(http.StatusOK, Response{
		Code:    code,
		Message: message,
	})
}

// ==================== 使用示例 ====================

func main() {
	r := gin.Default()

	r.GET("/users", listUsers)
	r.GET("/users/:id", getUser)
	r.POST("/users", createUser)
	r.DELETE("/users/:id", deleteUser)

	r.Run(":8080")
}

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

func listUsers(c *gin.Context) {
	users := []User{
		{ID: 1, Username: "tom", Email: "tom@example.com"},
		{ID: 2, Username: "jerry", Email: "jerry@example.com"},
	}

	// 分页响应
	SuccessPage(c, users, 100, 1, 10)
}

func getUser(c *gin.Context) {
	id := c.Param("id")

	if id == "999" {
		// 业务错误：使用业务状态码
		BusinessError(c, ErrCodeUserNotFound, "用户不存在")
		return
	}

	user := User{ID: 1, Username: "tom", Email: "tom@example.com"}
	Success(c, user)
}

func createUser(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		BadRequest(c, "参数错误: "+err.Error())
		return
	}

	user.ID = 1
	Created(c, user)
}

func deleteUser(c *gin.Context) {
	id := c.Param("id")
	if id == "999" {
		NotFound(c, "用户不存在")
		return
	}

	SuccessMessage(c, "删除成功")
}
