// 04-dependency-direction: 依赖方向示例
//
// 📌 依赖方向最佳实践:
//   - 依赖应该单向流动: handler → service → repository
//   - 避免循环依赖（Go 编译器会报错）
//   - 使用接口解耦，依赖接口而非实现
//
// 正确的依赖方向:
//   main.go
//      ↓
//   handler (HTTP 处理)
//      ↓
//   service (业务逻辑)
//      ↓
//   repository (数据访问)
//      ↓
//   model (数据模型) ← 被所有层依赖
package main

import (
	"fmt"
	"log"

	"dependency-demo/internal/handler"
	"dependency-demo/internal/repository"
	"dependency-demo/internal/service"
)

func main() {
	// 依赖注入：从下往上构建
	// 1. 创建 Repository（最底层）
	userRepo := repository.NewUserRepository()

	// 2. 创建 Service（依赖 Repository）
	userSvc := service.NewUserService(userRepo)

	// 3. 创建 Handler（依赖 Service）
	userHandler := handler.NewUserHandler(userSvc)

	// 模拟 HTTP 请求处理
	fmt.Println("=== 模拟创建用户请求 ===")
	result := userHandler.HandleCreateUser("tom", "tom@example.com")
	fmt.Println(result)

	fmt.Println("\n=== 模拟获取用户请求 ===")
	result = userHandler.HandleGetUser("1")
	fmt.Println(result)
}

// 如果尝试循环依赖（如 repository 导入 service），
// Go 编译器会报错:
//   import cycle not allowed
//
// 这是 Go 强制良好架构的方式！
func demoCircularDependency() {
	log.Println("Go 不允许循环依赖，强制你思考正确的依赖方向")
}
