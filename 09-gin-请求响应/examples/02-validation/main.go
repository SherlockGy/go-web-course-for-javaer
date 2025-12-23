// 02-validation: 参数验证
//
// 📌 常用验证 tag:
//   - required: 必填
//   - min/max: 最小/最大值（数字）或长度（字符串）
//   - email: 邮箱格式
//   - url: URL 格式
//   - oneof: 枚举值
//   - gt/gte/lt/lte: 大于/大于等于/小于/小于等于
//
// 📌 自定义验证:
//   - 可以注册自定义验证函数
//   - 使用 validator/v10 包
package main

import (
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// RegisterRequest 注册请求
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=20,alphanum"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6,max=50"`
	Phone    string `json:"phone" binding:"required,phone"` // 自定义验证
	Gender   string `json:"gender" binding:"required,oneof=male female other"`
	Age      int    `json:"age" binding:"required,gte=18,lte=120"`
}

// SearchRequest 搜索请求
type SearchRequest struct {
	Keyword  string `form:"keyword" binding:"required,min=1"`
	Page     int    `form:"page" binding:"required,min=1"`
	PageSize int    `form:"page_size" binding:"required,min=1,max=100"`
	OrderBy  string `form:"order_by" binding:"omitempty,oneof=created_at updated_at name"`
}

func main() {
	r := gin.Default()

	// 注册自定义验证器
	registerCustomValidators()

	r.POST("/register", register)
	r.GET("/search", search)

	r.Run(":8080")
}

// 注册自定义验证器
func registerCustomValidators() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// 注册手机号验证
		v.RegisterValidation("phone", validatePhone)
	}
}

// 自定义手机号验证
func validatePhone(fl validator.FieldLevel) bool {
	phone := fl.Field().String()
	// 简单的中国手机号正则
	matched, _ := regexp.MatchString(`^1[3-9]\d{9}$`, phone)
	return matched
}

func register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// 解析验证错误
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数验证失败",
			"errors":  parseValidationError(err),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"code":    0,
		"message": "注册成功",
		"data": gin.H{
			"username": req.Username,
			"email":    req.Email,
			"phone":    req.Phone,
		},
	})
}

func search(c *gin.Context) {
	var req SearchRequest
	// ShouldBindQuery 绑定 URL 查询参数
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数验证失败",
			"errors":  parseValidationError(err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"keyword":   req.Keyword,
			"page":      req.Page,
			"page_size": req.PageSize,
			"order_by":  req.OrderBy,
		},
	})
}

// 解析验证错误为友好格式
func parseValidationError(err error) []gin.H {
	var errors []gin.H

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrors {
			errors = append(errors, gin.H{
				"field":   e.Field(),
				"tag":     e.Tag(),
				"value":   e.Value(),
				"message": getErrorMessage(e),
			})
		}
	} else {
		errors = append(errors, gin.H{
			"message": err.Error(),
		})
	}

	return errors
}

// 获取友好的错误消息
func getErrorMessage(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return e.Field() + " 是必填字段"
	case "email":
		return e.Field() + " 必须是有效的邮箱地址"
	case "min":
		return e.Field() + " 长度不能小于 " + e.Param()
	case "max":
		return e.Field() + " 长度不能大于 " + e.Param()
	case "gte":
		return e.Field() + " 必须大于等于 " + e.Param()
	case "lte":
		return e.Field() + " 必须小于等于 " + e.Param()
	case "oneof":
		return e.Field() + " 必须是以下值之一: " + e.Param()
	case "phone":
		return e.Field() + " 必须是有效的手机号"
	default:
		return e.Field() + " 验证失败"
	}
}
