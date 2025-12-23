// Package service 业务逻辑层
package service

import (
	"errors"
	"sync"
	"time"

	"project-layout/internal/model"
	"project-layout/pkg/utils"
)

// UserService 用户服务
type UserService struct {
	mu    sync.RWMutex
	users map[string]*model.User
}

// NewUserService 创建用户服务
func NewUserService() *UserService {
	return &UserService{
		users: make(map[string]*model.User),
	}
}

// CreateUser 创建用户
func (s *UserService) CreateUser(username, email string) (*model.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	user := &model.User{
		ID:        utils.GenerateID(), // 使用 pkg 中的工具
		Username:  username,
		Email:     email,
		CreatedAt: time.Now(),
	}

	s.users[user.ID] = user
	return user, nil
}

// GetUser 获取用户
func (s *UserService) GetUser(id string) (*model.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	user, ok := s.users[id]
	if !ok {
		return nil, errors.New("用户不存在")
	}
	return user, nil
}
