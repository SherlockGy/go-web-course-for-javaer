# 07 - Gin 快速上手

## 学习目标

掌握 Gin 框架的基本使用，对比原生 `net/http` 的优势。

---

## 学习要点

### 1. 安装 Gin

```bash
go get -u github.com/gin-gonic/gin
```

### 2. Hello Gin

```go
package main

import "github.com/gin-gonic/gin"

func main() {
    r := gin.Default()  // 创建路由器（带 Logger 和 Recovery）

    r.GET("/hello", func(c *gin.Context) {
        c.String(200, "Hello, Gin!")
    })

    r.Run(":8080")  // 启动服务器
}
```

### 3. gin.Default() vs gin.New()

```go
// Default = New + Logger + Recovery
r := gin.Default()

// New = 空白路由器
r := gin.New()
r.Use(gin.Logger())    // 手动添加日志
r.Use(gin.Recovery())  // 手动添加恢复
```

### 4. 路由注册

```go
// HTTP 方法
r.GET("/users", listUsers)
r.POST("/users", createUser)
r.PUT("/users/:id", updateUser)
r.DELETE("/users/:id", deleteUser)

// 任意方法
r.Any("/ping", pingHandler)

// 匹配多个方法
r.Match([]string{"GET", "POST"}, "/data", dataHandler)
```

### 5. 路由分组

```go
// API 分组
api := r.Group("/api")
{
    // /api/users
    api.GET("/users", listUsers)
    api.POST("/users", createUser)

    // 嵌套分组 /api/v1
    v1 := api.Group("/v1")
    {
        v1.GET("/users", listUsersV1)
    }
}
```

### 6. 路径参数

```go
// 路径参数 :name
r.GET("/users/:id", func(c *gin.Context) {
    id := c.Param("id")  // 获取参数
    c.String(200, "User ID: %s", id)
})

// 通配符 *path
r.GET("/files/*path", func(c *gin.Context) {
    path := c.Param("path")  // 例如: /static/css/style.css
    c.String(200, "File: %s", path)
})
```

### 7. 查询参数

```go
// GET /search?keyword=go&page=1
r.GET("/search", func(c *gin.Context) {
    keyword := c.Query("keyword")           // "go"
    page := c.DefaultQuery("page", "1")     // 默认值 "1"
    limit := c.Query("limit")               // "" (不存在)

    c.JSON(200, gin.H{
        "keyword": keyword,
        "page":    page,
    })
})
```

---

## Gin vs 原生对比

| 功能 | 原生 net/http | Gin |
|------|---------------|-----|
| 路由注册 | `mux.HandleFunc("GET /users/{id}", h)` | `r.GET("/users/:id", h)` |
| 获取路径参数 | `r.PathValue("id")` | `c.Param("id")` |
| 获取查询参数 | `r.URL.Query().Get("key")` | `c.Query("key")` |
| JSON 响应 | `json.NewEncoder(w).Encode(data)` | `c.JSON(200, data)` |
| 设置状态码 | `w.WriteHeader(404)` | `c.Status(404)` |

---

## 示例代码

### examples/01-hello-gin/
Gin 版 Hello World

### examples/02-route-groups/
路由分组示例

### examples/03-params/
参数获取（路径参数、查询参数）

---

## 作业任务

### 任务描述
用 Gin 重写第 6 章的 Todo API。

### 要求
1. 实现相同的 5 个接口
2. 使用 Gin 的路由分组
3. 对比代码量差异

### 路由结构
```go
r := gin.Default()

api := r.Group("/api")
{
    api.GET("/todos", listTodos)
    api.POST("/todos", createTodo)
    api.GET("/todos/:id", getTodo)
    api.PUT("/todos/:id", updateTodo)
    api.DELETE("/todos/:id", deleteTodo)
}
```

### 验收标准
- 功能与第 6 章完全相同
- 代码比原生版本更简洁
- 统计两个版本的代码行数差异

### 代码对比示例

**原生版本**：
```go
func createTodo(w http.ResponseWriter, r *http.Request) {
    var todo Todo
    if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(ApiResponse{Code: 400, Message: "无效的请求体"})
        return
    }
    addTodo(&todo)
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(ApiResponse{Code: 200, Message: "success", Data: todo})
}
```

**Gin 版本**：
```go
func createTodo(c *gin.Context) {
    var todo Todo
    if err := c.ShouldBindJSON(&todo); err != nil {
        c.JSON(400, gin.H{"code": 400, "message": "无效的请求体"})
        return
    }
    addTodo(&todo)
    c.JSON(201, gin.H{"code": 200, "message": "success", "data": todo})
}
```

---

## 参考资料
- [Gin 官方文档](https://gin-gonic.com/docs/)
- [Gin GitHub](https://github.com/gin-gonic/gin)
