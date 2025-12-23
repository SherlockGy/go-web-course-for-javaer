// Package handler HTTP 处理层（最上层）
//
// 依赖: service
// 被依赖: main
package handler

import (
	"fmt"

	"dependency-demo/internal/service"
)

// UserHandler 用户处理器
type UserHandler struct {
	svc *service.UserService
}

// NewUserHandler 创建处理器
func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

// HandleCreateUser 处理创建用户请求
func (h *UserHandler) HandleCreateUser(username, email string) string {
	user, err := h.svc.CreateUser(username, email)
	if err != nil {
		return fmt.Sprintf("创建失败: %v", err)
	}
	return fmt.Sprintf("创建成功: ID=%s, Username=%s", user.ID, user.Username)
}

// HandleGetUser 处理获取用户请求
func (h *UserHandler) HandleGetUser(id string) string {
	user, err := h.svc.GetUser(id)
	if err != nil {
		return fmt.Sprintf("获取失败: %v", err)
	}
	return fmt.Sprintf("用户信息: ID=%s, Username=%s, Email=%s",
		user.ID, user.Username, user.Email)
}
