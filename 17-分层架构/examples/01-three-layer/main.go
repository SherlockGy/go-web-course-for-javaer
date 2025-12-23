// 01-three-layer: 三层架构示例
//
// 📌 架构分层:
//   Handler (表现层) → Service (业务层) → Repository (数据层)
//
// 📌 依赖方向:
//   - 上层依赖下层接口
//   - 下层不依赖上层
//
// 📌 与 Java Spring 对比:
//   - Java: @Controller → @Service → @Repository
//   - Go: 手动构造，依赖注入更显式
//
// 📌 项目结构:
//   ├── main.go          # 启动入口，组装依赖
//   ├── model/           # 数据模型、DTO
//   ├── handler/         # HTTP 处理器
//   ├── service/         # 业务逻辑
//   └── repository/      # 数据访问
package main

import (
	"fmt"
	"log"
	"three-layer/handler"
	"three-layer/model"
	"three-layer/repository"
	"three-layer/service"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	// 1. 初始化数据库
	db, err := initDB()
	if err != nil {
		log.Fatalf("数据库初始化失败: %v", err)
	}

	// 2. 组装依赖（依赖注入）
	// 📌 与 Java Spring 的 @Autowired 类似，但更显式
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	// 3. 设置路由
	r := gin.Default()

	api := r.Group("/api")
	userHandler.RegisterRoutes(api)

	// 4. 启动服务
	fmt.Println("服务器运行在 http://localhost:8080")
	fmt.Println("\n测试命令:")
	fmt.Println(`创建: curl -X POST http://localhost:8080/api/users -H "Content-Type: application/json" -d '{"username":"tom","email":"tom@example.com","password":"123456"}'`)
	fmt.Println(`列表: curl http://localhost:8080/api/users`)
	fmt.Println(`详情: curl http://localhost:8080/api/users/1`)
	fmt.Println(`更新: curl -X PUT http://localhost:8080/api/users/1 -H "Content-Type: application/json" -d '{"username":"tom2"}'`)
	fmt.Println(`删除: curl -X DELETE http://localhost:8080/api/users/1`)

	r.Run(":8080")
}

func initDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// 自动迁移
	if err := db.AutoMigrate(&model.User{}); err != nil {
		return nil, err
	}

	return db, nil
}
