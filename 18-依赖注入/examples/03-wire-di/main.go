// 03-wire-di: Google Wire 编译时依赖注入
//
// 📌 Wire 特点:
//   - 编译时代码生成（无反射）
//   - 类型安全
//   - 自动解析依赖图
//
// 📌 与 Java Spring 对比:
//   - Java Spring: 运行时反射注入
//   - Go Wire: 编译时生成初始化代码
//
// 📌 使用步骤:
//   1. 定义 Provider 函数（构造函数）
//   2. 定义 Injector 函数（wire.Build）
//   3. 运行 wire 生成代码
//
// 📌 本示例展示 Wire 的概念，实际使用需要安装 wire 工具
//    go install github.com/google/wire/cmd/wire@latest
package main

import (
	"fmt"
)

// ==================== 配置 ====================

type Config struct {
	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	ServerPort int
}

// NewConfig Provider: 创建配置
func NewConfig() *Config {
	return &Config{
		DBHost:     "localhost",
		DBPort:     3306,
		DBUser:     "root",
		DBPassword: "password",
		ServerPort: 8080,
	}
}

// ==================== 数据库 ====================

type Database struct {
	host     string
	port     int
	user     string
	password string
}

// NewDatabase Provider: 创建数据库连接
// 📌 依赖 Config
func NewDatabase(config *Config) (*Database, error) {
	fmt.Printf("连接数据库: %s:%d\n", config.DBHost, config.DBPort)
	return &Database{
		host:     config.DBHost,
		port:     config.DBPort,
		user:     config.DBUser,
		password: config.DBPassword,
	}, nil
}

func (db *Database) Query(sql string) string {
	return fmt.Sprintf("执行SQL: %s", sql)
}

// ==================== Repository ====================

type UserRepository struct {
	db *Database
}

// NewUserRepository Provider: 创建用户仓储
// 📌 依赖 Database
func NewUserRepository(db *Database) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) FindByID(id int) string {
	return r.db.Query(fmt.Sprintf("SELECT * FROM users WHERE id = %d", id))
}

// ==================== Service ====================

type UserService struct {
	repo *UserRepository
}

// NewUserService Provider: 创建用户服务
// 📌 依赖 UserRepository
func NewUserService(repo *UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetUser(id int) string {
	return s.repo.FindByID(id)
}

// ==================== Handler ====================

type UserHandler struct {
	service *UserService
}

// NewUserHandler Provider: 创建用户处理器
// 📌 依赖 UserService
func NewUserHandler(service *UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) HandleGetUser(id int) string {
	return h.service.GetUser(id)
}

// ==================== Server ====================

type Server struct {
	config      *Config
	userHandler *UserHandler
}

// NewServer Provider: 创建服务器
// 📌 依赖 Config 和 UserHandler
func NewServer(config *Config, userHandler *UserHandler) *Server {
	return &Server{
		config:      config,
		userHandler: userHandler,
	}
}

func (s *Server) Start() {
	fmt.Printf("服务器启动在端口: %d\n", s.config.ServerPort)
}

// ==================== Wire Injector (概念展示) ====================
/*
// wire.go - 这个文件会被 wire 工具处理
// +build wireinject

package main

import "github.com/google/wire"

// InitializeServer 定义如何组装 Server
// Wire 会根据这个函数生成实际的初始化代码
func InitializeServer() (*Server, error) {
    wire.Build(
        NewConfig,
        NewDatabase,
        NewUserRepository,
        NewUserService,
        NewUserHandler,
        NewServer,
    )
    return nil, nil // 这行代码会被 wire 替换
}
*/

// ==================== 手动组装（等价于 Wire 生成的代码）====================

// InitializeServer 手动实现的依赖注入
// 📌 这就是 Wire 会自动生成的代码
func InitializeServer() (*Server, error) {
	config := NewConfig()

	database, err := NewDatabase(config)
	if err != nil {
		return nil, err
	}

	userRepository := NewUserRepository(database)
	userService := NewUserService(userRepository)
	userHandler := NewUserHandler(userService)

	server := NewServer(config, userHandler)
	return server, nil
}

// ==================== 主函数 ====================

func main() {
	fmt.Println("=== Google Wire 编译时依赖注入 ===\n")

	// 依赖关系图:
	// Config
	//   ├── Database
	//   │     └── UserRepository
	//   │           └── UserService
	//   │                 └── UserHandler
	//   │                       └── Server
	//   └── Server

	fmt.Println("依赖关系:")
	fmt.Println("Config → Database → UserRepository → UserService → UserHandler → Server")
	fmt.Println()

	// 初始化（Wire 会自动生成这部分代码）
	server, err := InitializeServer()
	if err != nil {
		fmt.Printf("初始化失败: %v\n", err)
		return
	}

	// 启动服务器
	server.Start()

	// 测试请求
	fmt.Println()
	result := server.userHandler.HandleGetUser(1)
	fmt.Println(result)

	fmt.Println("\n📌 Wire 的优势:")
	fmt.Println("1. 自动解析依赖图，无需手动排序初始化顺序")
	fmt.Println("2. 编译时检查，确保所有依赖都能被满足")
	fmt.Println("3. 生成普通 Go 代码，无运行时反射开销")
	fmt.Println("4. 依赖关系变更时，只需重新运行 wire 即可")
}
