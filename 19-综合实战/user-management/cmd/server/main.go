// cmd/server/main.go - 应用入口
//
// 📌 综合实战：用户管理系统
//
// 技术栈:
//   - Gin: Web 框架
//   - GORM: ORM 框架
//   - Viper: 配置管理
//   - Zap: 结构化日志
//   - JWT: 认证
//   - bcrypt: 密码哈希
//
// 项目结构:
//   ├── cmd/server/          # 入口
//   ├── config.yaml          # 配置文件
//   └── internal/            # 内部包
//       ├── config/          # 配置
//       ├── model/           # 模型
//       ├── repository/      # 数据访问
//       ├── service/         # 业务逻辑
//       ├── handler/         # HTTP 处理
//       └── middleware/      # 中间件
//
// API:
//   POST /api/register       - 注册
//   POST /api/login          - 登录
//   GET  /api/profile        - 获取个人信息 (需认证)
//   PUT  /api/profile        - 更新个人信息 (需认证)
//   PUT  /api/password       - 修改密码 (需认证)
//   GET  /api/admin/users    - 用户列表 (需管理员)
//   DELETE /api/admin/users/:id - 删除用户 (需管理员)
package main

import (
	"fmt"
	"log"
	"os"
	"user-management/internal/config"
	"user-management/internal/handler"
	"user-management/internal/middleware"
	"user-management/internal/model"
	"user-management/internal/repository"
	"user-management/internal/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	// 1. 加载配置
	cfg, err := config.Load("config.yaml")
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 2. 初始化日志
	logger := initLogger(cfg.Log.Level)
	defer logger.Sync()

	// 3. 初始化数据库
	db, err := initDB(cfg.Database.DSN)
	if err != nil {
		logger.Fatal("数据库初始化失败", zap.Error(err))
	}

	// 创建管理员账号
	createAdminUser(db)

	// 4. 依赖注入
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo, &cfg.JWT)
	userHandler := handler.NewUserHandler(userService)

	// 5. 设置 Gin
	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(middleware.RecoveryMiddleware(logger))
	r.Use(middleware.LoggerMiddleware(logger))

	// 6. 注册路由
	api := r.Group("/api")
	authMiddleware := middleware.AuthMiddleware(&cfg.JWT)
	adminMiddleware := middleware.AdminMiddleware()
	userHandler.RegisterRoutes(api, authMiddleware, adminMiddleware)

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// 7. 启动服务
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	logger.Info("服务器启动", zap.String("addr", addr))

	printUsage(cfg.Server.Port)

	if err := r.Run(addr); err != nil {
		logger.Fatal("服务器启动失败", zap.Error(err))
	}
}

func initLogger(level string) *zap.Logger {
	var zapLevel zapcore.Level
	switch level {
	case "debug":
		zapLevel = zapcore.DebugLevel
	case "info":
		zapLevel = zapcore.InfoLevel
	case "warn":
		zapLevel = zapcore.WarnLevel
	case "error":
		zapLevel = zapcore.ErrorLevel
	default:
		zapLevel = zapcore.InfoLevel
	}

	config := zap.Config{
		Level:            zap.NewAtomicLevelAt(zapLevel),
		Development:      true,
		Encoding:         "console",
		EncoderConfig:    zap.NewDevelopmentEncoderConfig(),
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	logger, _ := config.Build()
	return logger
}

func initDB(dsn string) (*gorm.DB, error) {
	// 确保数据目录存在
	os.MkdirAll("./data", 0755)

	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// 自动迁移
	if err := db.AutoMigrate(&model.User{}); err != nil {
		return nil, err
	}

	return db, nil
}

func createAdminUser(db *gorm.DB) {
	var count int64
	db.Model(&model.User{}).Where("username = ?", "admin").Count(&count)
	if count > 0 {
		return
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	admin := &model.User{
		Username: "admin",
		Email:    "admin@example.com",
		Password: string(hashedPassword),
		Role:     "admin",
	}
	db.Create(admin)
	fmt.Println("✓ 创建管理员账号: admin / admin123")
}

func printUsage(port int) {
	fmt.Println("\n========== 用户管理系统 API ==========")
	fmt.Printf("服务地址: http://localhost:%d\n\n", port)
	fmt.Println("测试命令:")
	fmt.Println("  # 注册")
	fmt.Printf("  curl -X POST http://localhost:%d/api/register -H \"Content-Type: application/json\" -d '{\"username\":\"tom\",\"email\":\"tom@example.com\",\"password\":\"123456\"}'\n\n", port)
	fmt.Println("  # 登录")
	fmt.Printf("  curl -X POST http://localhost:%d/api/login -H \"Content-Type: application/json\" -d '{\"username\":\"admin\",\"password\":\"admin123\"}'\n\n", port)
	fmt.Println("  # 获取个人信息 (需要 token)")
	fmt.Printf("  curl http://localhost:%d/api/profile -H \"Authorization: Bearer <token>\"\n\n", port)
	fmt.Println("  # 用户列表 (需要管理员 token)")
	fmt.Printf("  curl http://localhost:%d/api/admin/users -H \"Authorization: Bearer <admin_token>\"\n", port)
	fmt.Println("==========================================\n")
}
