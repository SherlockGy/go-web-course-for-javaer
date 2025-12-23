// 01-json-binding: JSON 请求绑定
//
// 📌 绑定方法对比:
//   - ShouldBindJSON: 失败返回 error，不自动响应
//   - BindJSON: 失败自动响应 400，设置 Content-Type
//   - 推荐使用 ShouldBindXXX 系列，更灵活
//
// 📌 最佳实践:
//   - 使用 ShouldBindJSON 自行处理错误
//   - 合理使用 binding tag 进行验证
//   - 使用指针区分 "未传" 和 "零值"
package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// CreateUserRequest 创建用户请求
type CreateUserRequest struct {
	Username string `json:"username" binding:"required,min=3,max=20"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Age      *int   `json:"age" binding:"omitempty,min=0,max=150"` // 可选字段
}

// UpdateUserRequest 更新用户请求
type UpdateUserRequest struct {
	Username *string `json:"username" binding:"omitempty,min=3,max=20"` // 可选
	Email    *string `json:"email" binding:"omitempty,email"`           // 可选
	Age      *int    `json:"age" binding:"omitempty,min=0,max=150"`     // 可选
}

func main() {
	r := gin.Default()

	r.POST("/users", createUser)
	r.PUT("/users/:id", updateUser)

	r.Run(":8080")
}

func createUser(c *gin.Context) {
	var req CreateUserRequest

	// ShouldBindJSON: 推荐方式，自行处理错误
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数验证失败",
			"error":   err.Error(),
		})
		return
	}

	// 处理可选字段
	age := 0
	if req.Age != nil {
		age = *req.Age
	}

	c.JSON(http.StatusCreated, gin.H{
		"code":    0,
		"message": "创建成功",
		"data": gin.H{
			"id":       1,
			"username": req.Username,
			"email":    req.Email,
			"age":      age,
		},
	})
}

func updateUser(c *gin.Context) {
	id := c.Param("id")

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数验证失败",
			"error":   err.Error(),
		})
		return
	}

	// 只更新传入的字段
	updates := gin.H{"id": id}
	if req.Username != nil {
		updates["username"] = *req.Username
	}
	if req.Email != nil {
		updates["email"] = *req.Email
	}
	if req.Age != nil {
		updates["age"] = *req.Age
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "更新成功",
		"data":    updates,
	})
}
