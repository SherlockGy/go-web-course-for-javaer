// Package repository 数据访问层
//
// 依赖: model
// 被依赖: service
package repository

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"dependency-demo/internal/model"
)

// UserRepository 用户数据仓库
type UserRepository struct {
	mu     sync.RWMutex
	users  map[string]*model.User
	nextID int
}

// NewUserRepository 创建仓库
func NewUserRepository() *UserRepository {
	return &UserRepository{
		users:  make(map[string]*model.User),
		nextID: 1,
	}
}

// Create 创建用户
func (r *UserRepository) Create(user *model.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	user.ID = fmt.Sprintf("%d", r.nextID)
	user.CreatedAt = time.Now()
	r.nextID++

	r.users[user.ID] = user
	return nil
}

// FindByID 按 ID 查找
func (r *UserRepository) FindByID(id string) (*model.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, ok := r.users[id]
	if !ok {
		return nil, errors.New("用户不存在")
	}
	return user, nil
}
