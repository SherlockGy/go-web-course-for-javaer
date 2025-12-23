// internal/handler/user_handler.go - 用户处理器
package handler

import (
	"errors"
	"net/http"
	"strconv"
	"user-management/internal/model"
	"user-management/internal/repository"
	"user-management/internal/service"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service service.UserService
}

func NewUserHandler(service service.UserService) *UserHandler {
	return &UserHandler{service: service}
}

// RegisterRoutes 注册路由
func (h *UserHandler) RegisterRoutes(r *gin.RouterGroup, authMiddleware, adminMiddleware gin.HandlerFunc) {
	// 公开路由
	r.POST("/register", h.Register)
	r.POST("/login", h.Login)

	// 需要认证的路由
	auth := r.Group("")
	auth.Use(authMiddleware)
	{
		auth.GET("/profile", h.GetProfile)
		auth.PUT("/profile", h.UpdateProfile)
		auth.PUT("/password", h.ChangePassword)
	}

	// 管理员路由
	admin := r.Group("/admin")
	admin.Use(authMiddleware, adminMiddleware)
	{
		admin.GET("/users", h.GetUsers)
		admin.DELETE("/users/:id", h.DeleteUser)
	}
}

// Register 用户注册
// @Summary 用户注册
// @Tags 用户
// @Accept json
// @Produce json
// @Param request body model.RegisterRequest true "注册信息"
// @Success 201 {object} Response
// @Router /api/register [post]
func (h *UserHandler) Register(c *gin.Context) {
	var req model.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{Code: 400, Message: "参数错误: " + err.Error()})
		return
	}

	user, err := h.service.Register(&req)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, Response{Code: 0, Message: "注册成功", Data: user})
}

// Login 用户登录
func (h *UserHandler) Login(c *gin.Context) {
	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{Code: 400, Message: "参数错误"})
		return
	}

	resp, err := h.service.Login(&req)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, Response{Code: 0, Message: "登录成功", Data: resp})
}

// GetProfile 获取个人信息
func (h *UserHandler) GetProfile(c *gin.Context) {
	userID := c.GetUint("userID")

	user, err := h.service.GetProfile(userID)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, Response{Code: 0, Message: "success", Data: user})
}

// UpdateProfile 更新个人信息
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID := c.GetUint("userID")

	var req model.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{Code: 400, Message: "参数错误"})
		return
	}

	user, err := h.service.UpdateProfile(userID, &req)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, Response{Code: 0, Message: "更新成功", Data: user})
}

// ChangePassword 修改密码
func (h *UserHandler) ChangePassword(c *gin.Context) {
	userID := c.GetUint("userID")

	var req model.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{Code: 400, Message: "参数错误"})
		return
	}

	if err := h.service.ChangePassword(userID, &req); err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, Response{Code: 0, Message: "密码修改成功"})
}

// GetUsers 获取用户列表（管理员）
func (h *UserHandler) GetUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

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
		Data: PageData{
			List:     users,
			Total:    total,
			Page:     page,
			PageSize: pageSize,
		},
	})
}

// DeleteUser 删除用户（管理员）
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

	c.JSON(http.StatusOK, Response{Code: 0, Message: "删除成功"})
}

func (h *UserHandler) handleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, repository.ErrUserNotFound):
		c.JSON(http.StatusNotFound, Response{Code: 404, Message: "用户不存在"})
	case errors.Is(err, service.ErrInvalidCredentials):
		c.JSON(http.StatusUnauthorized, Response{Code: 401, Message: "用户名或密码错误"})
	case errors.Is(err, service.ErrUsernameExists):
		c.JSON(http.StatusConflict, Response{Code: 409, Message: "用户名已存在"})
	case errors.Is(err, service.ErrEmailExists):
		c.JSON(http.StatusConflict, Response{Code: 409, Message: "邮箱已存在"})
	case errors.Is(err, service.ErrWrongPassword):
		c.JSON(http.StatusBadRequest, Response{Code: 400, Message: "原密码错误"})
	default:
		c.JSON(http.StatusInternalServerError, Response{Code: 500, Message: "服务器错误"})
	}
}

// Response 统一响应
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// PageData 分页数据
type PageData struct {
	List     interface{} `json:"list"`
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"page_size"`
}
