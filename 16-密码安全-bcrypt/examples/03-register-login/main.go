// 03-register-login: 完整的注册登录示例
//
// 📌 安全最佳实践:
//   - 密码强度验证（长度、复杂度）
//   - bcrypt 哈希存储
//   - 登录失败锁定
//   - 统一错误消息（防止用户枚举）
//
// 📌 与 Java Spring Security 对比:
//   - Java: PasswordEncoder + UserDetailsService
//   - Go: 手动实现，更灵活但需自行处理
package main

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// User 用户模型
type User struct {
	ID           uint      `json:"id"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"-"` // JSON 序列化时忽略
	Email        string    `json:"email"`
	CreatedAt    time.Time `json:"created_at"`
}

// 模拟数据库
var (
	users   = make(map[string]*User)
	usersMu sync.RWMutex
	nextID  uint = 1
)

// 登录失败计数器
var (
	loginAttempts   = make(map[string]int)
	loginAttemptsMu sync.Mutex
)

func main() {
	r := gin.Default()

	r.POST("/register", registerHandler)
	r.POST("/login", loginHandler)

	fmt.Println("服务器运行在 http://localhost:8080")
	fmt.Println("测试命令:")
	fmt.Println(`  curl -X POST http://localhost:8080/register -H "Content-Type: application/json" -d '{"username":"tom","password":"Test@123456","email":"tom@example.com"}'`)
	fmt.Println(`  curl -X POST http://localhost:8080/login -H "Content-Type: application/json" -d '{"username":"tom","password":"Test@123456"}'`)

	r.Run(":8080")
}

// RegisterRequest 注册请求
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=20"`
	Password string `json:"password" binding:"required,min=8"`
	Email    string `json:"email" binding:"required,email"`
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func registerHandler(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	// 1. 验证密码强度
	if err := ValidatePassword(req.Password); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 2. 检查用户名是否存在
	usersMu.RLock()
	if _, exists := users[req.Username]; exists {
		usersMu.RUnlock()
		c.JSON(http.StatusConflict, gin.H{"error": "用户名已存在"})
		return
	}
	usersMu.RUnlock()

	// 3. 哈希密码
	hash, err := HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器错误"})
		return
	}

	// 4. 保存用户
	usersMu.Lock()
	user := &User{
		ID:           nextID,
		Username:     req.Username,
		PasswordHash: hash,
		Email:        req.Email,
		CreatedAt:    time.Now(),
	}
	users[req.Username] = user
	nextID++
	usersMu.Unlock()

	c.JSON(http.StatusCreated, gin.H{
		"message": "注册成功",
		"user_id": user.ID,
	})
}

func loginHandler(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	// 1. 检查登录失败次数
	if isLocked(req.Username) {
		c.JSON(http.StatusTooManyRequests, gin.H{
			"error": "登录尝试次数过多，请稍后再试",
		})
		return
	}

	// 2. 查找用户
	usersMu.RLock()
	user, exists := users[req.Username]
	usersMu.RUnlock()

	// 📌 安全最佳实践: 统一错误消息，防止用户枚举
	if !exists {
		recordFailedLogin(req.Username)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	// 3. 验证密码
	if !CheckPassword(req.Password, user.PasswordHash) {
		recordFailedLogin(req.Username)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	// 4. 登录成功，清除失败计数
	clearFailedLogins(req.Username)

	c.JSON(http.StatusOK, gin.H{
		"message": "登录成功",
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
		},
	})
}

// ==================== 密码工具函数 ====================

// HashPassword 哈希密码
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}

// CheckPassword 验证密码
func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// ValidatePassword 验证密码强度
// 📌 最佳实践: 至少8位，包含大小写字母、数字、特殊字符
func ValidatePassword(password string) error {
	if len(password) < 8 {
		return errors.New("密码至少8位")
	}
	if len(password) > 72 { // bcrypt 限制
		return errors.New("密码不能超过72位")
	}

	var (
		hasUpper   = regexp.MustCompile(`[A-Z]`).MatchString(password)
		hasLower   = regexp.MustCompile(`[a-z]`).MatchString(password)
		hasNumber  = regexp.MustCompile(`[0-9]`).MatchString(password)
		hasSpecial = regexp.MustCompile(`[!@#$%^&*(),.?":{}|<>]`).MatchString(password)
	)

	if !hasUpper {
		return errors.New("密码需包含大写字母")
	}
	if !hasLower {
		return errors.New("密码需包含小写字母")
	}
	if !hasNumber {
		return errors.New("密码需包含数字")
	}
	if !hasSpecial {
		return errors.New("密码需包含特殊字符")
	}

	return nil
}

// ==================== 登录失败锁定 ====================

const maxLoginAttempts = 5

// isLocked 检查是否被锁定
func isLocked(username string) bool {
	loginAttemptsMu.Lock()
	defer loginAttemptsMu.Unlock()
	return loginAttempts[username] >= maxLoginAttempts
}

// recordFailedLogin 记录失败登录
func recordFailedLogin(username string) {
	loginAttemptsMu.Lock()
	defer loginAttemptsMu.Unlock()
	loginAttempts[username]++
}

// clearFailedLogins 清除失败记录
func clearFailedLogins(username string) {
	loginAttemptsMu.Lock()
	defer loginAttemptsMu.Unlock()
	delete(loginAttempts, username)
}
