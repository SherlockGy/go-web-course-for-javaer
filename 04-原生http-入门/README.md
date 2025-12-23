# 04 - 原生 HTTP 入门

## 学习目标

掌握 Go 标准库 `net/http` 的基本用法，理解 HTTP 服务器的工作原理。

---

## 学习要点

### 1. 最简 HTTP 服务器

```go
package main

import (
    "fmt"
    "net/http"
)

func main() {
    // 注册处理函数
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Hello, World!")
    })

    // 启动服务器
    fmt.Println("服务器启动: http://localhost:8080")
    http.ListenAndServe(":8080", nil)
}
```

### 2. http.Handler 接口

```go
// 标准接口定义
type Handler interface {
    ServeHTTP(ResponseWriter, *Request)
}

// 自定义 Handler
type MyHandler struct{}

func (h *MyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Hello from MyHandler"))
}

// 使用
http.Handle("/", &MyHandler{})
```

### 3. http.HandleFunc 便捷方式

```go
// HandleFunc 是 Handler 的便捷封装
http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Hello"))
})

// 等价于
http.Handle("/hello", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Hello"))
}))
```

### 4. ResponseWriter 和 Request

```go
func handler(w http.ResponseWriter, r *http.Request) {
    // 请求信息
    method := r.Method           // GET, POST, PUT, DELETE
    path := r.URL.Path           // /users/123
    query := r.URL.Query()       // URL 查询参数
    header := r.Header           // 请求头

    // 响应操作
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)  // 200
    w.Write([]byte(`{"message": "ok"}`))
}
```

### 5. Go 1.22 新路由语法

```go
// Go 1.22+ 支持方法和路径参数
mux := http.NewServeMux()

// 指定 HTTP 方法
mux.HandleFunc("GET /users", listUsers)
mux.HandleFunc("POST /users", createUser)

// 路径参数
mux.HandleFunc("GET /users/{id}", getUser)

func getUser(w http.ResponseWriter, r *http.Request) {
    id := r.PathValue("id")  // 获取路径参数
    fmt.Fprintf(w, "User ID: %s", id)
}
```

---

## 示例代码

### examples/01-simple-server/
最简 HTTP 服务器

### examples/02-multiple-routes/
多路由处理

### examples/03-path-params/
路径参数（Go 1.22+）

---

## 作业任务

### 任务描述
创建一个 HTTP 服务器，实现 `GET /hello/{name}` 接口。

### 要求
1. 使用 Go 1.22+ 的路由语法
2. 从路径中获取 `name` 参数
3. 返回 `Hello, {name}!`

### 验收标准
```bash
curl http://localhost:8080/hello/Tom
# 返回: Hello, Tom!

curl http://localhost:8080/hello/World
# 返回: Hello, World!
```

### 代码模板
```go
package main

import (
    "fmt"
    "net/http"
)

func main() {
    mux := http.NewServeMux()

    // TODO: 注册 GET /hello/{name} 路由

    fmt.Println("服务器启动: http://localhost:8080")
    http.ListenAndServe(":8080", mux)
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
    // TODO: 获取 name 参数并返回问候语
}
```

---

## 参考资料
- [net/http 包文档](https://pkg.go.dev/net/http)
- [Go 1.22 路由增强](https://go.dev/blog/routing-enhancements)
