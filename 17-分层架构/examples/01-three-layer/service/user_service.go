// service/user_service.go - 业务逻辑层
//
// 📌 Service 层职责:
//   - 实现业务逻辑
//   - 调用 Repository 层
//   - 事务管理
//   - 数据转换
//
// 📌 与 Java 对比:
//   - Java: @Service + @Transactional
//   - Go: 显式事务控制，更清晰
//
// 📌 最佳实践:
//   - Service 依赖 Repository 接口
//   - 返回业务错误，由 Handler 转换为 HTTP 错误
package service

import (
	"errors"
	"three-layer/model"
	"three-layer/repository"

	"golang.org/x/crypto/bcrypt"
)

// 业务错误
var (
	ErrInvalidCredentials = errors.New("用户名或密码错误")
	ErrPasswordTooWeak    = errors.New("密码强度不足")
)

// UserService 用户服务接口
type UserService interface {
	CreateUser(req *model.CreateUserRequest) (*model.UserResponse, error)
	GetUser(id uint) (*model.UserResponse, error)
	GetUsers(page, pageSize int) ([]*model.UserResponse, int64, error)
	UpdateUser(id uint, req *model.UpdateUserRequest) (*model.UserResponse, error)
	DeleteUser(id uint) error
}

// userService 实现
type userService struct {
	repo repository.UserRepository
}

// NewUserService 构造函数
// 📌 依赖注入: 接收接口而非具体实现
func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) CreateUser(req *model.CreateUserRequest) (*model.UserResponse, error) {
	// 1. 业务逻辑：密码哈希
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// 2. 构造实体
	user := &model.User{
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	// 3. 调用 Repository
	if err := s.repo.Create(user); err != nil {
		return nil, err
	}

	// 4. 返回 DTO
	return user.ToResponse(), nil
}

func (s *userService) GetUser(id uint) (*model.UserResponse, error) {
	user, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	return user.ToResponse(), nil
}

func (s *userService) GetUsers(page, pageSize int) ([]*model.UserResponse, int64, error) {
	// 计算偏移量
	offset := (page - 1) * pageSize

	users, total, err := s.repo.FindAll(offset, pageSize)
	if err != nil {
		return nil, 0, err
	}

	// 转换为响应 DTO
	responses := make([]*model.UserResponse, len(users))
	for i, user := range users {
		responses[i] = user.ToResponse()
	}

	return responses, total, nil
}

func (s *userService) UpdateUser(id uint, req *model.UpdateUserRequest) (*model.UserResponse, error) {
	// 1. 查找用户
	user, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// 2. 更新字段（只更新非空字段）
	if req.Username != "" {
		user.Username = req.Username
	}
	if req.Email != "" {
		user.Email = req.Email
	}

	// 3. 保存
	if err := s.repo.Update(user); err != nil {
		return nil, err
	}

	return user.ToResponse(), nil
}

func (s *userService) DeleteUser(id uint) error {
	return s.repo.Delete(id)
}
