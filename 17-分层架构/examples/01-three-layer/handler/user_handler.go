// handler/user_handler.go - 表现层（控制器）
//
// 📌 Handler 层职责:
//   - 接收 HTTP 请求
//   - 参数验证和绑定
//   - 调用 Service 层
//   - 返回 HTTP 响应
//
// 📌 与 Java 对比:
//   - Java: @RestController + @RequestMapping
//   - Go: 更轻量，无注解魔法
//
// 📌 最佳实践:
//   - Handler 只做参数处理和响应转换
//   - 业务逻辑放在 Service 层
//   - 统一响应格式
package handler

import (
	"errors"
	"net/http"
	"strconv"
	"three-layer/model"
	"three-layer/repository"
	"three-layer/service"

	"github.com/gin-gonic/gin"
)

// UserHandler 用户处理器
type UserHandler struct {
	service service.UserService
}

// NewUserHandler 构造函数
func NewUserHandler(service service.UserService) *UserHandler {
	return &UserHandler{service: service}
}

// Response 统一响应结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// PageResponse 分页响应
type PageResponse struct {
	List     interface{} `json:"list"`
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"page_size"`
}

// RegisterRoutes 注册路由
// 📌 与 Java @RequestMapping("/users") 类似
func (h *UserHandler) RegisterRoutes(r *gin.RouterGroup) {
	users := r.Group("/users")
	{
		users.POST("", h.CreateUser)
		users.GET("", h.GetUsers)
		users.GET("/:id", h.GetUser)
		users.PUT("/:id", h.UpdateUser)
		users.DELETE("/:id", h.DeleteUser)
	}
}

// CreateUser 创建用户
// POST /api/users
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req model.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "参数错误: " + err.Error(),
		})
		return
	}

	user, err := h.service.CreateUser(&req)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, Response{
		Code:    0,
		Message: "创建成功",
		Data:    user,
	})
}

// GetUser 获取用户
// GET /api/users/:id
func (h *UserHandler) GetUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{Code: 400, Message: "无效的ID"})
		return
	}

	user, err := h.service.GetUser(uint(id))
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "success",
		Data:    user,
	})
}

// GetUsers 获取用户列表
// GET /api/users?page=1&page_size=10
func (h *UserHandler) GetUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	// 参数校验
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	users, total, err := h.service.GetUsers(page, pageSize)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "success",
		Data: PageResponse{
			List:     users,
			Total:    total,
			Page:     page,
			PageSize: pageSize,
		},
	})
}

// UpdateUser 更新用户
// PUT /api/users/:id
func (h *UserHandler) UpdateUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{Code: 400, Message: "无效的ID"})
		return
	}

	var req model.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{Code: 400, Message: "参数错误"})
		return
	}

	user, err := h.service.UpdateUser(uint(id), &req)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "更新成功",
		Data:    user,
	})
}

// DeleteUser 删除用户
// DELETE /api/users/:id
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{Code: 400, Message: "无效的ID"})
		return
	}

	if err := h.service.DeleteUser(uint(id)); err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "删除成功",
	})
}

// handleError 统一错误处理
// 📌 将业务错误转换为 HTTP 响应
func (h *UserHandler) handleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, repository.ErrUserNotFound):
		c.JSON(http.StatusNotFound, Response{Code: 404, Message: "用户不存在"})
	case errors.Is(err, repository.ErrUserAlreadyExists):
		c.JSON(http.StatusConflict, Response{Code: 409, Message: "用户已存在"})
	default:
		c.JSON(http.StatusInternalServerError, Response{Code: 500, Message: "服务器错误"})
	}
}
