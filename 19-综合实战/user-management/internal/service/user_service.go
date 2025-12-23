// internal/service/user_service.go - 用户业务逻辑层
package service

import (
	"errors"
	"time"
	"user-management/internal/config"
	"user-management/internal/model"
	"user-management/internal/repository"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("用户名或密码错误")
	ErrUsernameExists     = errors.New("用户名已存在")
	ErrEmailExists        = errors.New("邮箱已存在")
	ErrWrongPassword      = errors.New("原密码错误")
)

// Claims JWT 声明
type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// UserService 用户服务接口
type UserService interface {
	Register(req *model.RegisterRequest) (*model.UserResponse, error)
	Login(req *model.LoginRequest) (*model.LoginResponse, error)
	GetProfile(userID uint) (*model.UserResponse, error)
	UpdateProfile(userID uint, req *model.UpdateProfileRequest) (*model.UserResponse, error)
	ChangePassword(userID uint, req *model.ChangePasswordRequest) error
	GetUsers(page, pageSize int) ([]*model.UserResponse, int64, error)
	DeleteUser(id uint) error
}

type userService struct {
	repo      repository.UserRepository
	jwtConfig *config.JWTConfig
}

func NewUserService(repo repository.UserRepository, jwtConfig *config.JWTConfig) UserService {
	return &userService{
		repo:      repo,
		jwtConfig: jwtConfig,
	}
}

func (s *userService) Register(req *model.RegisterRequest) (*model.UserResponse, error) {
	// 检查用户名是否存在
	if s.repo.ExistsByUsername(req.Username) {
		return nil, ErrUsernameExists
	}

	// 检查邮箱是否存在
	if s.repo.ExistsByEmail(req.Email) {
		return nil, ErrEmailExists
	}

	// 哈希密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashedPassword),
		Role:     "user",
	}

	if err := s.repo.Create(user); err != nil {
		return nil, err
	}

	return user.ToResponse(), nil
}

func (s *userService) Login(req *model.LoginRequest) (*model.LoginResponse, error) {
	user, err := s.repo.FindByUsername(req.Username)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	// 生成 JWT
	token, err := s.generateToken(user)
	if err != nil {
		return nil, err
	}

	return &model.LoginResponse{
		Token: token,
		User:  user.ToResponse(),
	}, nil
}

func (s *userService) GetProfile(userID uint) (*model.UserResponse, error) {
	user, err := s.repo.FindByID(userID)
	if err != nil {
		return nil, err
	}
	return user.ToResponse(), nil
}

func (s *userService) UpdateProfile(userID uint, req *model.UpdateProfileRequest) (*model.UserResponse, error) {
	user, err := s.repo.FindByID(userID)
	if err != nil {
		return nil, err
	}

	if req.Email != "" && req.Email != user.Email {
		if s.repo.ExistsByEmail(req.Email) {
			return nil, ErrEmailExists
		}
		user.Email = req.Email
	}

	if err := s.repo.Update(user); err != nil {
		return nil, err
	}

	return user.ToResponse(), nil
}

func (s *userService) ChangePassword(userID uint, req *model.ChangePasswordRequest) error {
	user, err := s.repo.FindByID(userID)
	if err != nil {
		return err
	}

	// 验证原密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword)); err != nil {
		return ErrWrongPassword
	}

	// 哈希新密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)
	return s.repo.Update(user)
}

func (s *userService) GetUsers(page, pageSize int) ([]*model.UserResponse, int64, error) {
	offset := (page - 1) * pageSize
	users, total, err := s.repo.FindAll(offset, pageSize)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]*model.UserResponse, len(users))
	for i, user := range users {
		responses[i] = user.ToResponse()
	}

	return responses, total, nil
}

func (s *userService) DeleteUser(id uint) error {
	return s.repo.Delete(id)
}

func (s *userService) generateToken(user *model.User) (string, error) {
	claims := Claims{
		UserID:   user.ID,
		Username: user.Username,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(s.jwtConfig.ExpireHours) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "user-management",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtConfig.Secret))
}
