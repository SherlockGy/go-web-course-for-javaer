# 19 - 综合实战

## 学习目标

整合前 18 章所学，从零实现一个完整的用户管理系统。

---

## 🆚 Java 对比：技术栈映射

| 功能 | Spring Boot | Go 实现 |
|------|-------------|---------|
| Web 框架 | Spring MVC | Gin |
| ORM | Spring Data JPA | GORM |
| 配置 | application.yml | Viper |
| 日志 | Logback | Zap |
| 认证 | Spring Security + JWT | golang-jwt |
| 密码 | BCryptPasswordEncoder | bcrypt |
| 依赖注入 | @Autowired | 手动构造函数 |

> **洞察**：Go 项目代码量可能比 Spring Boot 多 20%，但没有注解魔法，每一行都是你理解的逻辑。

---

## 最终项目结构

```
user-management/
├── cmd/
│   └── server/
│       └── main.go           # 入口，依赖组装
├── internal/
│   ├── config/
│   │   └── config.go         # 配置加载
│   ├── handler/
│   │   ├── auth_handler.go   # 认证处理
│   │   └── user_handler.go   # 用户处理
│   ├── middleware/
│   │   ├── auth.go           # JWT 认证
│   │   ├── logger.go         # 请求日志
│   │   └── cors.go           # 跨域
│   ├── model/
│   │   ├── user.go           # 用户实体
│   │   ├── request.go        # 请求 DTO
│   │   └── response.go       # 响应 DTO
│   ├── repository/
│   │   └── user_repository.go
│   └── service/
│       └── user_service.go
├── pkg/
│   ├── jwt/
│   │   └── jwt.go            # JWT 工具
│   └── logger/
│       └── logger.go         # 日志初始化
├── config.yaml               # 配置文件
├── go.mod
└── README.md
```

---

## API 设计

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| POST | /api/auth/register | 用户注册 | 否 |
| POST | /api/auth/login | 用户登录 | 否 |
| GET | /api/users | 用户列表 | 是 |
| GET | /api/users/:id | 用户详情 | 是 |
| PUT | /api/users/:id | 更新用户 | 是 |
| DELETE | /api/users/:id | 删除用户 | 是 |
| GET | /health | 健康检查 | 否 |

---

## 分步实现指南

### Step 1: 初始化项目

```bash
mkdir user-management && cd user-management
go mod init user-management
```

安装依赖：
```bash
go get github.com/gin-gonic/gin
go get gorm.io/gorm
go get gorm.io/driver/sqlite
go get github.com/spf13/viper
go get go.uber.org/zap
go get github.com/golang-jwt/jwt/v5
go get golang.org/x/crypto/bcrypt
```

### Step 2: 配置管理

**config.yaml**
```yaml
server:
  host: "0.0.0.0"
  port: 8080

database:
  dsn: "./data.db"

jwt:
  secret: "your-256-bit-secret-key-here"
  expiration: 86400

logging:
  level: "debug"
  file: "./logs/app.log"
```

### Step 3: 数据模型

**internal/model/user.go**
```go
type User struct {
    gorm.Model
    Username string `gorm:"uniqueIndex;size:50"`
    Password string `gorm:"size:100" json:"-"`
    Email    string `gorm:"uniqueIndex"`
    Nickname string `gorm:"size:50"`
}

func (u *User) ToDTO() *UserDTO {
    return &UserDTO{
        ID:       u.ID,
        Username: u.Username,
        Email:    u.Email,
        Nickname: u.Nickname,
    }
}
```

### Step 4: 实现各层

按照第 17 章的分层架构：
1. Repository → 数据访问
2. Service → 业务逻辑（含密码加密、JWT 生成）
3. Handler → HTTP 处理

### Step 5: 组装启动

**cmd/server/main.go**
```go
func main() {
    // 1. 加载配置
    cfg := config.Load()

    // 2. 初始化日志
    logger := logger.Init(cfg.Logging)

    // 3. 连接数据库
    db := initDB(cfg.Database)

    // 4. 依赖注入
    userRepo := repository.NewUserRepository(db)
    jwtUtil := jwt.NewJWTUtil(cfg.JWT)
    userSvc := service.NewUserService(userRepo, jwtUtil)
    authHandler := handler.NewAuthHandler(userSvc)
    userHandler := handler.NewUserHandler(userSvc)

    // 5. 路由配置
    r := gin.New()
    r.Use(middleware.Logger(logger))
    r.Use(middleware.Recovery())
    r.Use(middleware.CORS())

    // 公开路由
    r.GET("/health", healthHandler)
    auth := r.Group("/api/auth")
    {
        auth.POST("/register", authHandler.Register)
        auth.POST("/login", authHandler.Login)
    }

    // 受保护路由
    api := r.Group("/api")
    api.Use(middleware.JWT(jwtUtil))
    {
        api.GET("/users", userHandler.List)
        api.GET("/users/:id", userHandler.Get)
        api.PUT("/users/:id", userHandler.Update)
        api.DELETE("/users/:id", userHandler.Delete)
    }

    // 6. 启动服务
    addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
    logger.Info("服务器启动", zap.String("addr", addr))
    r.Run(addr)
}
```

---

## 功能检查清单

- [ ] **配置管理**：从 config.yaml 加载配置
- [ ] **日志系统**：请求日志、结构化日志、文件输出
- [ ] **数据库**：GORM + SQLite，自动迁移
- [ ] **用户注册**：用户名/邮箱唯一性检查，bcrypt 加密
- [ ] **用户登录**：密码验证，返回 JWT
- [ ] **JWT 认证**：中间件验证 Bearer Token
- [ ] **用户 CRUD**：列表、详情、更新、删除
- [ ] **统一响应**：`{code, message, data}` 格式
- [ ] **错误处理**：业务错误码，全局异常捕获
- [ ] **分层架构**：Handler → Service → Repository

---

## 验收测试

```bash
# 启动服务
go run cmd/server/main.go

# 健康检查
curl http://localhost:8080/health

# 注册
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"tom","password":"123456","email":"tom@example.com"}'

# 登录
TOKEN=$(curl -s -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"tom","password":"123456"}' | jq -r '.data.token')

echo "Token: $TOKEN"

# 获取用户列表（需认证）
curl http://localhost:8080/api/users \
  -H "Authorization: Bearer $TOKEN"

# 获取当前用户
curl http://localhost:8080/api/users/1 \
  -H "Authorization: Bearer $TOKEN"

# 更新用户
curl -X PUT http://localhost:8080/api/users/1 \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"nickname":"Tommy"}'

# 无认证访问（应返回 401）
curl http://localhost:8080/api/users
```

---

## 与参考项目对比

完成后，对比你的实现与 `java-go-comparison/go/web-demo/` 项目：

1. **目录结构**是否一致？
2. **代码风格**是否符合 Go 习惯？
3. **功能完整性**是否相同？

---

## 延伸挑战

完成基础功能后，可以尝试：

1. **添加分页**：用户列表支持 `page` 和 `pageSize` 参数
2. **添加搜索**：按用户名/邮箱模糊搜索
3. **添加角色**：用户角色（admin/user），权限控制
4. **添加测试**：单元测试、集成测试
5. **Docker 部署**：编写 Dockerfile，容器化部署

---

## 恭喜！

完成这个实战项目，你已经掌握了 Go Web 开发的核心技能：

- ✅ Gin 框架使用
- ✅ GORM 数据库操作
- ✅ JWT 认证
- ✅ 配置管理
- ✅ 日志系统
- ✅ 分层架构
- ✅ 依赖注入

接下来可以：
- 阅读开源项目代码
- 学习微服务架构（gRPC、服务发现）
- 深入 Go 并发编程

---

## 参考资料
- [java-go-comparison/go/web-demo](../../../java-go-comparison/go/web-demo/) - 完整参考实现
- [Gin 官方文档](https://gin-gonic.com/docs/)
- [GORM 官方文档](https://gorm.io/docs/)
