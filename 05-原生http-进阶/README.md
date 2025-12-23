# 05 - 原生 HTTP 进阶

## 学习目标

掌握中间件模式、JSON 处理和请求参数解析。

---

## 学习要点

### 1. 中间件模式

```go
// 中间件是一个包装 Handler 的函数
func loggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()

        // 调用下一个处理器
        next.ServeHTTP(w, r)

        // 记录日志
        fmt.Printf("%s %s %v\n", r.Method, r.URL.Path, time.Since(start))
    })
}

// 使用中间件
mux := http.NewServeMux()
mux.HandleFunc("/", homeHandler)

// 包装
handler := loggingMiddleware(mux)
http.ListenAndServe(":8080", handler)
```

### 2. 中间件链

```go
// 多个中间件串联
func chain(h http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
    for i := len(middlewares) - 1; i >= 0; i-- {
        h = middlewares[i](h)
    }
    return h
}

// 使用
handler := chain(mux, loggingMiddleware, authMiddleware, corsMiddleware)
```

### 3. JSON 响应

```go
import "encoding/json"

type User struct {
    ID   int64  `json:"id"`
    Name string `json:"name"`
}

func jsonHandler(w http.ResponseWriter, r *http.Request) {
    user := User{ID: 1, Name: "Tom"}

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(user)
}
```

### 4. JSON 请求体解析

```go
type CreateUserRequest struct {
    Name  string `json:"name"`
    Email string `json:"email"`
}

func createUser(w http.ResponseWriter, r *http.Request) {
    var req CreateUserRequest

    // 解析 JSON 请求体
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "无效的 JSON", http.StatusBadRequest)
        return
    }

    fmt.Printf("创建用户: %+v\n", req)
    w.WriteHeader(http.StatusCreated)
}
```

### 5. Query 参数获取

```go
// GET /search?keyword=go&page=1
func searchHandler(w http.ResponseWriter, r *http.Request) {
    query := r.URL.Query()

    keyword := query.Get("keyword")  // "go"
    page := query.Get("page")        // "1"（字符串）

    // 类型转换
    pageNum, _ := strconv.Atoi(page)

    fmt.Fprintf(w, "搜索: %s, 页码: %d", keyword, pageNum)
}
```

### 6. 表单数据解析

```go
func formHandler(w http.ResponseWriter, r *http.Request) {
    // 解析表单
    if err := r.ParseForm(); err != nil {
        http.Error(w, "解析表单失败", http.StatusBadRequest)
        return
    }

    username := r.FormValue("username")
    password := r.FormValue("password")

    fmt.Fprintf(w, "用户名: %s", username)
}
```

---

## 示例代码

### examples/01-middleware/
日志中间件实现

### examples/02-json-response/
JSON 响应示例

### examples/03-json-request/
JSON 请求体解析

### examples/04-query-params/
URL 查询参数处理

---

## 作业任务

### 任务 1：日志中间件
实现一个日志中间件，打印请求方法、路径、响应时间。

**输出格式**：
```
[2024-01-01 12:00:00] GET /users 2.5ms
```

### 任务 2：JSON API
实现 `POST /users` 接口：
- 接收 JSON 请求体：`{"name": "Tom", "email": "tom@example.com"}`
- 返回 JSON 响应：`{"id": 1, "name": "Tom", "email": "tom@example.com"}`

### 验收标准
```bash
# 测试日志中间件
curl http://localhost:8080/users
# 控制台输出: [2024-01-01 12:00:00] GET /users 1.2ms

# 测试 JSON API
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{"name": "Tom", "email": "tom@example.com"}'
# 返回: {"id":1,"name":"Tom","email":"tom@example.com"}
```

### 代码模板
```go
package main

import (
    "encoding/json"
    "fmt"
    "net/http"
    "time"
)

// TODO: 实现日志中间件
func loggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // 记录开始时间
        // 调用 next.ServeHTTP(w, r)
        // 打印日志
    })
}

// TODO: 实现创建用户处理器
func createUserHandler(w http.ResponseWriter, r *http.Request) {
    // 解析 JSON
    // 返回 JSON
}

func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("POST /users", createUserHandler)

    handler := loggingMiddleware(mux)

    fmt.Println("服务器启动: http://localhost:8080")
    http.ListenAndServe(":8080", handler)
}
```

---

## 参考资料
- [encoding/json 包文档](https://pkg.go.dev/encoding/json)
- [Writing Web Applications](https://go.dev/doc/articles/wiki/)
