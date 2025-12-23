# 10 - Gin 错误处理

## 学习目标

掌握 Gin 的错误处理机制，实现统一的异常处理。

---

## 学习要点

### 1. Recovery 中间件原理

```go
// Recovery 捕获 panic，防止程序崩溃
func Recovery() gin.HandlerFunc {
    return func(c *gin.Context) {
        defer func() {
            if err := recover(); err != nil {
                // 记录错误
                log.Printf("panic recovered: %v", err)
                // 返回 500 错误
                c.AbortWithStatusJSON(500, gin.H{
                    "error": "内部服务器错误",
                })
            }
        }()
        c.Next()
    }
}
```

### 2. 自定义错误类型

```go
// 业务错误定义
type AppError struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
}

func (e *AppError) Error() string {
    return e.Message
}

// 预定义错误
var (
    ErrUserNotFound    = &AppError{Code: 404, Message: "用户不存在"}
    ErrUserExists      = &AppError{Code: 400, Message: "用户名已存在"}
    ErrInvalidPassword = &AppError{Code: 401, Message: "密码错误"}
    ErrUnauthorized    = &AppError{Code: 401, Message: "未授权"}
    ErrInternal        = &AppError{Code: 500, Message: "服务器内部错误"}
)
```

### 3. 错误响应工具

```go
func HandleError(c *gin.Context, err error) {
    if appErr, ok := err.(*AppError); ok {
        // 业务错误
        c.JSON(appErr.Code, gin.H{
            "code":    appErr.Code,
            "message": appErr.Message,
        })
    } else {
        // 未知错误
        log.Printf("未知错误: %v", err)
        c.JSON(500, gin.H{
            "code":    500,
            "message": "服务器内部错误",
        })
    }
}

// 使用
func getUser(c *gin.Context) {
    user, err := userService.GetByID(id)
    if err != nil {
        HandleError(c, err)
        return
    }
    c.JSON(200, user)
}
```

### 4. 全局错误处理中间件

```go
func ErrorHandler() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Next()

        // 检查是否有错误
        if len(c.Errors) > 0 {
            err := c.Errors.Last().Err
            HandleError(c, err)
        }
    }
}

// 在 handler 中添加错误
func handler(c *gin.Context) {
    if err := doSomething(); err != nil {
        c.Error(err)  // 添加到错误列表
        return
    }
}
```

### 5. panic vs error 的选择

```go
// 使用 error：可恢复的业务错误
func GetUser(id int64) (*User, error) {
    user := db.Find(id)
    if user == nil {
        return nil, ErrUserNotFound  // 返回错误
    }
    return user, nil
}

// 使用 panic：不可恢复的严重错误（少用）
func MustInit() {
    if err := initDB(); err != nil {
        panic("数据库初始化失败: " + err.Error())  // 程序无法继续
    }
}
```

### 6. 业务错误码设计

```go
const (
    // 成功
    CodeSuccess = 200

    // 客户端错误 4xx
    CodeBadRequest   = 400  // 参数错误
    CodeUnauthorized = 401  // 未授权
    CodeForbidden    = 403  // 禁止访问
    CodeNotFound     = 404  // 资源不存在

    // 业务错误 1xxx
    CodeUserNotFound    = 1001  // 用户不存在
    CodeUserExists      = 1002  // 用户已存在
    CodeInvalidPassword = 1003  // 密码错误
    CodeEmailExists     = 1004  // 邮箱已存在

    // 服务器错误 5xx
    CodeInternal = 500  // 内部错误
)
```

### 7. 完整的错误处理流程

```go
// errors.go
type AppError struct {
    HTTPCode int    `json:"-"`
    Code     int    `json:"code"`
    Message  string `json:"message"`
}

func NewAppError(httpCode, code int, message string) *AppError {
    return &AppError{
        HTTPCode: httpCode,
        Code:     code,
        Message:  message,
    }
}

// handler.go
func createUser(c *gin.Context) {
    var req CreateUserRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"code": 400, "message": "参数错误: " + err.Error()})
        return
    }

    user, err := userService.Create(req)
    if err != nil {
        if appErr, ok := err.(*AppError); ok {
            c.JSON(appErr.HTTPCode, appErr)
        } else {
            c.JSON(500, gin.H{"code": 500, "message": "服务器错误"})
        }
        return
    }

    c.JSON(201, gin.H{"code": 200, "message": "success", "data": user})
}
```

---

## 示例代码

### examples/01-recovery/
Recovery 中间件演示

### examples/02-custom-errors/
自定义错误类型

### examples/03-global-handler/
全局错误处理中间件

---

## 作业任务

### 任务描述
设计业务错误码系统，实现全局错误处理中间件。

### 要求
1. 定义业务错误码：
   - 1001: 用户不存在
   - 1002: 用户名已存在
   - 1003: 密码错误
   - 1004: 邮箱已存在
   - 1005: 无效的令牌

2. 实现统一响应格式：
   ```json
   {
       "code": 1001,
       "message": "用户不存在",
       "data": null
   }
   ```

3. 实现全局错误处理中间件

### 验收标准
```bash
# 用户不存在
curl http://localhost:8080/api/users/999
# 返回: {"code":1001,"message":"用户不存在","data":null}

# 参数错误
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{}'
# 返回: {"code":400,"message":"参数错误: ...","data":null}

# 成功
curl http://localhost:8080/api/users/1
# 返回: {"code":200,"message":"success","data":{...}}
```

### 代码结构建议
```
homework/
├── main.go
├── errors/
│   └── errors.go     # 错误定义
├── middleware/
│   └── error.go      # 错误处理中间件
└── handler/
    └── user.go       # 用户处理器
```

---

## 思考题

1. 为什么要区分 HTTP 状态码和业务错误码？
2. 什么情况下应该使用 `panic`？
3. 如何做到对客户端友好（隐藏内部错误）又方便调试（记录详细日志）？

---

## 参考资料
- [Gin 自定义错误处理](https://gin-gonic.com/docs/examples/custom-http-config/)
- [Go 错误处理最佳实践](https://go.dev/blog/error-handling-and-go)
