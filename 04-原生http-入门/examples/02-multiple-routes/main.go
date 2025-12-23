// 02-multiple-routes: 多路由处理
//
// 📌 Go 1.22 新特性:
//   - 支持 HTTP 方法匹配: "GET /users"
//   - 支持路径参数: "/users/{id}"
//   - 更精确的路由匹配
//
// 运行: go run main.go
// 测试:
//   curl http://localhost:8080/
//   curl http://localhost:8080/users
//   curl -X POST http://localhost:8080/users
package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	// Go 1.22+ 新语法：方法 + 路径
	mux.HandleFunc("GET /", homeHandler)
	mux.HandleFunc("GET /users", listUsersHandler)
	mux.HandleFunc("POST /users", createUserHandler)
	mux.HandleFunc("GET /health", healthHandler)

	addr := ":8080"
	log.Printf("服务器启动: http://localhost%s", addr)

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "首页 - Go 1.22 路由示例")
}

func listUsersHandler(w http.ResponseWriter, r *http.Request) {
	// 只会匹配 GET 请求
	fmt.Fprintln(w, "用户列表: [tom, jerry, alice]")
}

func createUserHandler(w http.ResponseWriter, r *http.Request) {
	// 只会匹配 POST 请求
	fmt.Fprintln(w, "创建用户成功")
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintln(w, `{"status":"ok"}`)
}
