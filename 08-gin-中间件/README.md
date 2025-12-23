# 08 - Gin 中间件

## 学习目标

掌握 Gin 中间件的使用和自定义开发。

---

## 学习要点

### 1. 内置中间件

```go
// gin.Default() 包含以下中间件
r := gin.New()
r.Use(gin.Logger())    // 请求日志
r.Use(gin.Recovery())  // panic 恢复
```

### 2. 中间件执行顺序

```go
r.Use(middleware1)  // 先执行
r.Use(middleware2)  // 后执行

// 请求流程：
// middleware1 前置 → middleware2 前置 → handler → middleware2 后置 → middleware1 后置
```

### 3. 自定义中间件

```go
func LoggerMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()

        // 前置处理
        fmt.Printf("[开始] %s %s\n", c.Request.Method, c.Request.URL.Path)

        // 执行下一个处理器
        c.Next()

        // 后置处理
        duration := time.Since(start)
        status := c.Writer.Status()
        fmt.Printf("[完成] %s %s %d %v\n", c.Request.Method, c.Request.URL.Path, status, duration)
    }
}

// 使用
r.Use(LoggerMiddleware())
```

### 4. c.Next() 和 c.Abort()

```go
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        token := c.GetHeader("Authorization")

        if token == "" {
            c.JSON(401, gin.H{"error": "未授权"})
            c.Abort()  // 终止后续处理
            return
        }

        c.Next()  // 继续执行
    }
}
```

### 5. 中间件传值

```go
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 设置值
        c.Set("userID", 123)
        c.Set("username", "tom")
        c.Next()
    }
}

func handler(c *gin.Context) {
    // 获取值
    userID := c.GetInt("userID")
    username := c.GetString("username")

    // 或使用 MustGet（不存在会 panic）
    // userID := c.MustGet("userID").(int)
}
```

### 6. 不同级别的中间件

```go
// 全局中间件
r.Use(GlobalMiddleware())

// 分组中间件
api := r.Group("/api")
api.Use(APIMiddleware())
{
    api.GET("/users", listUsers)
}

// 路由级中间件
r.GET("/admin", AdminMiddleware(), adminHandler)
```

### 7. 常用中间件示例

**CORS 跨域**：
```go
func CORSMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Header("Access-Control-Allow-Origin", "*")
        c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }
        c.Next()
    }
}
```

**请求 ID**：
```go
func RequestIDMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        requestID := uuid.New().String()
        c.Set("requestID", requestID)
        c.Header("X-Request-ID", requestID)
        c.Next()
    }
}
```

---

## 示例代码

### examples/01-builtin-middleware/
内置中间件使用

### examples/02-custom-middleware/
自定义日志中间件

### examples/03-auth-middleware/
简单认证中间件

### examples/04-middleware-chain/
中间件链执行顺序演示

---

## 作业任务

### 任务 1：请求耗时中间件
实现一个中间件，记录每个请求的耗时。

**要求**：
- 打印格式：`[2024-01-01 12:00:00] GET /api/users 200 2.5ms`
- 包含：时间、方法、路径、状态码、耗时

### 任务 2：API Key 认证中间件
实现一个简单的 API Key 认证中间件。

**要求**：
- 从请求头 `X-API-Key` 获取 Key
- 有效 Key 列表：`["key1", "key2", "key3"]`
- 无效 Key 返回 401

### 验收标准
```bash
# 无 API Key
curl http://localhost:8080/api/users
# 返回: {"error": "未提供 API Key"}

# 错误的 API Key
curl -H "X-API-Key: wrong" http://localhost:8080/api/users
# 返回: {"error": "无效的 API Key"}

# 正确的 API Key
curl -H "X-API-Key: key1" http://localhost:8080/api/users
# 返回: 正常数据
```

### 代码模板
```go
// 耗时中间件
func TimingMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // TODO: 实现
    }
}

// API Key 认证中间件
func APIKeyAuthMiddleware(validKeys []string) gin.HandlerFunc {
    return func(c *gin.Context) {
        // TODO: 实现
    }
}

func main() {
    r := gin.New()

    // 全局中间件
    r.Use(TimingMiddleware())
    r.Use(gin.Recovery())

    // 需要认证的路由
    api := r.Group("/api")
    api.Use(APIKeyAuthMiddleware([]string{"key1", "key2", "key3"}))
    {
        api.GET("/users", listUsers)
    }

    r.Run(":8080")
}
```

---

## 参考资料
- [Gin 中间件文档](https://gin-gonic.com/docs/examples/custom-middleware/)
- [常用 Gin 中间件集合](https://github.com/gin-contrib)
