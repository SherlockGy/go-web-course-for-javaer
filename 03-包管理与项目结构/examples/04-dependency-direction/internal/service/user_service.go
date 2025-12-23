// Package service 业务逻辑层
//
// 依赖: model, repository
// 被依赖: handler
package service

import (
	"dependency-demo/internal/model"
	"dependency-demo/internal/repository"
)

// UserService 用户服务
type UserService struct {
	repo *repository.UserRepository
}

// NewUserService 创建服务
func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

// CreateUser 创建用户（业务逻辑）
func (s *UserService) CreateUser(username, email string) (*model.User, error) {
	// 这里可以添加业务规则
	// 比如：用户名格式验证、邮箱唯一性检查等

	user := &model.User{
		Username: username,
		Email:    email,
	}

	if err := s.repo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

// GetUser 获取用户
func (s *UserService) GetUser(id string) (*model.User, error) {
	return s.repo.FindByID(id)
}
