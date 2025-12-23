# 17 - 分层架构

## 学习目标

掌握 Web 应用的分层架构设计，实现关注点分离。

---

## 🆚 Java 对比：分层命名

| 层级 | Spring MVC | Go 习惯 |
|------|------------|---------|
| 表示层 | `@Controller` | `handler/` |
| 业务层 | `@Service` | `service/` |
| 数据层 | `@Repository` | `repository/` |
| 实体 | `@Entity` | `model/` |

> **洞察**：Go 用目录名而非注解来表示职责。没有框架魔法，全靠约定。

---

## 学习要点

### 1. 三层架构思想

```
┌─────────────┐
│   Handler   │  ← HTTP 请求/响应处理
├─────────────┤
│   Service   │  ← 业务逻辑
├─────────────┤
│ Repository  │  ← 数据访问
└─────────────┘
```

### 2. 各层职责

| 层级 | 职责 | 不该做的事 |
|------|------|------------|
| **Handler** | 参数解析、验证、调用 Service、返回响应 | 写业务逻辑、直接操作数据库 |
| **Service** | 业务逻辑、事务控制、调用多个 Repository | 处理 HTTP、直接写 SQL |
| **Repository** | 数据库 CRUD、数据查询 | 写业务逻辑、处理 HTTP |

### 3. 目录结构

```
internal/
├── handler/
│   ├── user_handler.go      # HTTP 处理
│   └── auth_handler.go
├── service/
│   └── user_service.go      # 业务逻辑
├── repository/
│   └── user_repository.go   # 数据访问
└── model/
    ├── user.go              # 实体
    ├── request.go           # 请求 DTO
    └── response.go          # 响应 DTO
```

### 4. 层间数据传递

```go
// Entity（数据库模型）
type User struct {
    gorm.Model
    Username string
    Password string  // 不应暴露给前端
    Email    string
}

// DTO（对外传输）
type UserDTO struct {
    ID       uint   `json:"id"`
    Username string `json:"username"`
    Email    string `json:"email"`
}

// Entity → DTO 转换
func (u *User) ToDTO() *UserDTO {
    return &UserDTO{
        ID:       u.ID,
        Username: u.Username,
        Email:    u.Email,
    }
}
```

> **🆚 Java 对比**
> ```java
> // Java 常用 ModelMapper 或 MapStruct
> UserDTO dto = modelMapper.map(user, UserDTO.class);
>
> // Go 通常手写转换方法
> dto := user.ToDTO()
> ```
> Go 更显式，Java 更自动化。Go 的方式虽然代码多一点，但转换逻辑清晰可控。

### 5. 完整代码示例

**repository/user_repository.go**
```go
type UserRepository struct {
    db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
    return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *model.User) error {
    return r.db.Create(user).Error
}

func (r *UserRepository) FindByID(id uint) (*model.User, error) {
    var user model.User
    err := r.db.First(&user, id).Error
    return &user, err
}
```

**service/user_service.go**
```go
type UserService struct {
    repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
    return &UserService{repo: repo}
}

func (s *UserService) GetUser(id uint) (*model.UserDTO, error) {
    user, err := s.repo.FindByID(id)
    if err != nil {
        return nil, errors.New("用户不存在")
    }
    return user.ToDTO(), nil
}
```

**handler/user_handler.go**
```go
type UserHandler struct {
    svc *service.UserService
}

func NewUserHandler(svc *service.UserService) *UserHandler {
    return &UserHandler{svc: svc}
}

func (h *UserHandler) GetUser(c *gin.Context) {
    id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

    user, err := h.svc.GetUser(uint(id))
    if err != nil {
        c.JSON(404, gin.H{"error": err.Error()})
        return
    }

    c.JSON(200, gin.H{"data": user})
}
```

### 6. 组装依赖（main.go）

```go
func main() {
    db := initDB()

    // 依赖注入（从下往上）
    userRepo := repository.NewUserRepository(db)
    userSvc := service.NewUserService(userRepo)
    userHandler := handler.NewUserHandler(userSvc)

    // 路由
    r := gin.Default()
    r.GET("/users/:id", userHandler.GetUser)
    r.Run(":8080")
}
```

> **🆚 Spring 对比**
> ```java
> // Spring 用 @Autowired 自动注入
> @Service
> public class UserService {
>     @Autowired
>     private UserRepository repo;
> }
> ```
> Spring 的 DI 是"隐式"的（框架扫描注解），Go 的 DI 是"显式"的（手动 New）。

---

## 示例代码

### examples/01-layered-architecture/
完整的分层架构示例

---

## 作业任务

### 任务描述
将之前的 Todo API 重构为三层架构。

### 要求
1. 创建 `internal/handler/todo_handler.go`
2. 创建 `internal/service/todo_service.go`
3. 创建 `internal/repository/todo_repository.go`
4. 使用 GORM 持久化（替代内存存储）

### 目录结构
```
homework/
├── cmd/server/main.go
├── internal/
│   ├── handler/todo_handler.go
│   ├── service/todo_service.go
│   ├── repository/todo_repository.go
│   └── model/todo.go
└── go.mod
```

### 验收标准
- Handler 不直接操作数据库
- Service 不处理 HTTP
- Repository 只做数据访问
- 功能与之前相同

---

## 思考题

1. 为什么 Handler 不应该直接调用 Repository？
2. Service 层的价值是什么？（如果业务很简单，可以省略吗？）
3. DTO 和 Entity 分离的好处是什么？

---

## 参考资料
- [Standard Go Project Layout](https://github.com/golang-standards/project-layout)
