# 18 - 依赖注入

## 学习目标

理解依赖注入的价值，掌握 Go 中的 DI 实现方式。

---

## 🆚 Java 对比：DI 哲学

| 方面 | Spring | Go |
|------|--------|-----|
| 注入方式 | `@Autowired`（自动） | 构造函数（手动） |
| 容器 | IoC Container | 无（或 Wire） |
| 配置 | 注解 + XML | 代码 |
| 运行时 | 反射扫描 | 编译期确定 |

> **洞察**：Spring 的 DI 是"魔法"——你不知道依赖怎么来的。Go 的 DI 是"显式"的——依赖关系在 main 里一目了然。

---

## 学习要点

### 1. 为什么需要依赖注入

**没有 DI（紧耦合）**：
```go
type UserService struct {
    repo *UserRepository  // 直接依赖具体实现
}

func NewUserService() *UserService {
    return &UserService{
        repo: NewUserRepository(),  // 内部创建依赖
    }
}
```

**问题**：
- 无法替换 `UserRepository`（比如测试时用 Mock）
- `UserService` 和 `UserRepository` 强绑定

### 2. 手动依赖注入

**有 DI（松耦合）**：
```go
type UserService struct {
    repo *UserRepository
}

// 依赖从外部传入
func NewUserService(repo *UserRepository) *UserService {
    return &UserService{repo: repo}
}
```

```go
// main.go 组装
repo := NewUserRepository(db)
service := NewUserService(repo)  // 注入依赖
handler := NewUserHandler(service)
```

> **🆚 Spring 对比**
> ```java
> @Service
> public class UserService {
>     private final UserRepository repo;
>
>     // Spring 推荐构造函数注入
>     public UserService(UserRepository repo) {
>         this.repo = repo;
>     }
> }
> ```
> 思想相同！Spring 只是用 IoC 容器自动完成了 `new` 的过程。

### 3. 接口解耦

```go
// 定义接口（在 service 包中）
type UserRepository interface {
    Create(user *User) error
    FindByID(id uint) (*User, error)
    FindByUsername(username string) (*User, error)
}

// Service 依赖接口
type UserService struct {
    repo UserRepository  // 接口，不是具体类型
}
```

```go
// 实现接口（在 repository 包中）
type userRepositoryImpl struct {
    db *gorm.DB
}

func (r *userRepositoryImpl) Create(user *User) error {
    return r.db.Create(user).Error
}

func (r *userRepositoryImpl) FindByID(id uint) (*User, error) {
    // ...
}
```

**好处**：
- 可以轻松替换实现（MySQL → PostgreSQL）
- 测试时可以注入 Mock

### 4. 测试时注入 Mock

```go
// Mock 实现
type mockUserRepository struct {
    users map[uint]*User
}

func (m *mockUserRepository) FindByID(id uint) (*User, error) {
    if user, ok := m.users[id]; ok {
        return user, nil
    }
    return nil, errors.New("not found")
}

// 测试
func TestUserService_GetUser(t *testing.T) {
    // 注入 Mock
    mockRepo := &mockUserRepository{
        users: map[uint]*User{1: {ID: 1, Username: "tom"}},
    }
    service := NewUserService(mockRepo)

    user, err := service.GetUser(1)
    assert.NoError(t, err)
    assert.Equal(t, "tom", user.Username)
}
```

> **🆚 Java 对比**
> ```java
> @Mock
> private UserRepository mockRepo;
>
> @InjectMocks
> private UserService service;
>
> @Test
> void testGetUser() {
>     when(mockRepo.findById(1L)).thenReturn(Optional.of(user));
>     // ...
> }
> ```
> Java 用 Mockito 框架，Go 手写 Mock 或用 gomock。思想相同。

### 5. Wire（可选）

当依赖很多时，手动组装很繁琐。Wire 是 Google 的编译期 DI 工具。

```go
// wire.go
//go:build wireinject

func InitializeApp(db *gorm.DB) *App {
    wire.Build(
        repository.NewUserRepository,
        service.NewUserService,
        handler.NewUserHandler,
        NewApp,
    )
    return nil
}
```

```bash
wire ./...  # 生成 wire_gen.go
```

> **洞察**：Wire 在编译期生成代码，不像 Spring 在运行时反射。性能更好，但需要额外工具。

---

## 示例代码

### examples/01-manual-di/
手动依赖注入

### examples/02-interface-decouple/
接口解耦示例

---

## 作业任务

### 任务描述
为 Repository 定义接口，实现依赖注入。

### 要求
1. 定义 `UserRepository` 接口
2. `UserService` 依赖接口而非具体实现
3. 在 main 中组装依赖
4. 编写测试，使用 Mock Repository

### 代码结构
```go
// internal/repository/interface.go
type UserRepository interface {
    Create(user *model.User) error
    FindByID(id uint) (*model.User, error)
    FindByUsername(username string) (*model.User, error)
    ExistsByUsername(username string) (bool, error)
}

// internal/repository/user_repository.go
type userRepository struct {
    db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
    return &userRepository{db: db}
}

// internal/service/user_service.go
type UserService struct {
    repo repository.UserRepository  // 依赖接口
}
```

### 验收标准
```go
// 可以轻松替换实现
realRepo := repository.NewUserRepository(db)
mockRepo := &MockUserRepository{}

// 两种都能工作
service1 := service.NewUserService(realRepo)
service2 := service.NewUserService(mockRepo)
```

---

## 思考题

1. Go 为什么没有像 Spring 那样的 IoC 容器？
2. "依赖接口而非实现"这个原则的价值是什么？
3. 什么时候应该用 Wire，什么时候手动注入就够了？

---

## 参考资料
- [Wire GitHub](https://github.com/google/wire)
- [Go 依赖注入实践](https://go.dev/blog/wire)
