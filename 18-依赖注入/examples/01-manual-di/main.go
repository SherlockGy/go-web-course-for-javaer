// 01-manual-di: 手动依赖注入
//
// 📌 Go 依赖注入理念:
//   - 优先使用手动构造函数注入
//   - 简单、显式、无魔法
//   - 编译时检查依赖关系
//
// 📌 与 Java Spring 对比:
//   - Java: @Autowired 自动注入，运行时反射
//   - Go: 构造函数显式传入，编译时检查
//
// 📌 优点:
//   - 依赖关系一目了然
//   - 无需学习 DI 框架
//   - IDE 跳转和重构友好
package main

import (
	"fmt"
	"time"
)

// ==================== 接口定义 ====================

// Logger 日志接口
type Logger interface {
	Info(msg string)
	Error(msg string)
}

// UserRepository 用户仓储接口
type UserRepository interface {
	FindByID(id int) (*User, error)
	Save(user *User) error
}

// EmailService 邮件服务接口
type EmailService interface {
	Send(to, subject, body string) error
}

// ==================== 实体 ====================

type User struct {
	ID       int
	Name     string
	Email    string
	CreateAt time.Time
}

// ==================== 实现 ====================

// ConsoleLogger 控制台日志实现
type ConsoleLogger struct{}

func (l *ConsoleLogger) Info(msg string) {
	fmt.Printf("[INFO] %s\n", msg)
}

func (l *ConsoleLogger) Error(msg string) {
	fmt.Printf("[ERROR] %s\n", msg)
}

// MemoryUserRepository 内存用户仓储实现
type MemoryUserRepository struct {
	users map[int]*User
}

func NewMemoryUserRepository() *MemoryUserRepository {
	return &MemoryUserRepository{
		users: make(map[int]*User),
	}
}

func (r *MemoryUserRepository) FindByID(id int) (*User, error) {
	if user, ok := r.users[id]; ok {
		return user, nil
	}
	return nil, fmt.Errorf("user not found: %d", id)
}

func (r *MemoryUserRepository) Save(user *User) error {
	r.users[user.ID] = user
	return nil
}

// SMTPEmailService SMTP 邮件服务实现
type SMTPEmailService struct {
	host string
	port int
}

func NewSMTPEmailService(host string, port int) *SMTPEmailService {
	return &SMTPEmailService{host: host, port: port}
}

func (s *SMTPEmailService) Send(to, subject, body string) error {
	fmt.Printf("发送邮件: to=%s, subject=%s\n", to, subject)
	return nil
}

// ==================== 业务服务 ====================

// UserService 用户服务
// 📌 通过构造函数接收依赖（接口类型）
type UserService struct {
	repo   UserRepository
	email  EmailService
	logger Logger
}

// NewUserService 构造函数注入
// 📌 与 Java Spring @Autowired 构造函数注入类似
func NewUserService(repo UserRepository, email EmailService, logger Logger) *UserService {
	return &UserService{
		repo:   repo,
		email:  email,
		logger: logger,
	}
}

func (s *UserService) Register(id int, name, email string) error {
	s.logger.Info(fmt.Sprintf("注册用户: %s", name))

	user := &User{
		ID:       id,
		Name:     name,
		Email:    email,
		CreateAt: time.Now(),
	}

	if err := s.repo.Save(user); err != nil {
		s.logger.Error(fmt.Sprintf("保存用户失败: %v", err))
		return err
	}

	if err := s.email.Send(email, "欢迎注册", "感谢您的注册!"); err != nil {
		s.logger.Error(fmt.Sprintf("发送邮件失败: %v", err))
		// 邮件失败不影响注册
	}

	return nil
}

func (s *UserService) GetUser(id int) (*User, error) {
	return s.repo.FindByID(id)
}

// ==================== 主函数：组装依赖 ====================

func main() {
	// 📌 在入口处组装所有依赖
	// 这种模式称为 "Composition Root"

	// 1. 创建基础设施组件
	logger := &ConsoleLogger{}
	userRepo := NewMemoryUserRepository()
	emailService := NewSMTPEmailService("smtp.example.com", 587)

	// 2. 创建业务服务（注入依赖）
	userService := NewUserService(userRepo, emailService, logger)

	// 3. 使用服务
	fmt.Println("=== 手动依赖注入示例 ===\n")

	// 注册用户
	if err := userService.Register(1, "Tom", "tom@example.com"); err != nil {
		fmt.Printf("注册失败: %v\n", err)
		return
	}

	// 查询用户
	user, err := userService.GetUser(1)
	if err != nil {
		fmt.Printf("查询失败: %v\n", err)
		return
	}
	fmt.Printf("\n查询结果: ID=%d, Name=%s, Email=%s\n", user.ID, user.Name, user.Email)
}
