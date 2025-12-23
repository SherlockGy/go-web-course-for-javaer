// repository/user_repository.go - 数据访问层
//
// 📌 Repository 层职责:
//   - 封装数据库操作
//   - 提供 CRUD 方法
//   - 不包含业务逻辑
//
// 📌 与 Java 对比:
//   - Java: JpaRepository<User, Long> 接口 + 方法命名约定
//   - Go: 手动实现，更灵活
//
// 📌 最佳实践:
//   - 定义接口，便于 mock 测试
//   - 返回领域错误而非数据库错误
package repository

import (
	"errors"
	"three-layer/model"

	"gorm.io/gorm"
)

// 定义领域错误
var (
	ErrUserNotFound      = errors.New("用户不存在")
	ErrUserAlreadyExists = errors.New("用户已存在")
)

// UserRepository 用户仓储接口
// 📌 面向接口编程，便于测试时 mock
type UserRepository interface {
	Create(user *model.User) error
	FindByID(id uint) (*model.User, error)
	FindByUsername(username string) (*model.User, error)
	FindAll(offset, limit int) ([]*model.User, int64, error)
	Update(user *model.User) error
	Delete(id uint) error
}

// userRepository 实现
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository 构造函数
// 📌 与 Java @Repository + 构造函数注入类似
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *model.User) error {
	// 检查是否已存在
	var count int64
	r.db.Model(&model.User{}).Where("username = ? OR email = ?", user.Username, user.Email).Count(&count)
	if count > 0 {
		return ErrUserAlreadyExists
	}

	return r.db.Create(user).Error
}

func (r *userRepository) FindByID(id uint) (*model.User, error) {
	var user model.User
	err := r.db.First(&user, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrUserNotFound
	}
	return &user, err
}

func (r *userRepository) FindByUsername(username string) (*model.User, error) {
	var user model.User
	err := r.db.Where("username = ?", username).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrUserNotFound
	}
	return &user, err
}

func (r *userRepository) FindAll(offset, limit int) ([]*model.User, int64, error) {
	var users []*model.User
	var total int64

	r.db.Model(&model.User{}).Count(&total)
	err := r.db.Offset(offset).Limit(limit).Find(&users).Error

	return users, total, err
}

func (r *userRepository) Update(user *model.User) error {
	result := r.db.Save(user)
	if result.RowsAffected == 0 {
		return ErrUserNotFound
	}
	return result.Error
}

func (r *userRepository) Delete(id uint) error {
	result := r.db.Delete(&model.User{}, id)
	if result.RowsAffected == 0 {
		return ErrUserNotFound
	}
	return result.Error
}
