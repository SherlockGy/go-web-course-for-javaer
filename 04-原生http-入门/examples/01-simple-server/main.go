// 01-simple-server: 最简单的 HTTP 服务器
//
// 📌 关键概念:
//   - http.ListenAndServe 启动服务器
//   - http.HandleFunc 注册路由处理函数
//   - http.ResponseWriter 写入响应
//   - *http.Request 读取请求
//
// 运行: go run main.go
// 测试: curl http://localhost:8080/
package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	// 注册路由处理函数
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/about", aboutHandler)

	// 启动服务器
	addr := ":8080"
	log.Printf("服务器启动: http://localhost%s", addr)

	// ListenAndServe 会阻塞，直到服务器停止
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}

// homeHandler 处理首页请求
func homeHandler(w http.ResponseWriter, r *http.Request) {
	// r.Method: 请求方法 (GET/POST/...)
	// r.URL.Path: 请求路径
	// r.Header: 请求头
	log.Printf("[%s] %s", r.Method, r.URL.Path)

	// 设置响应头
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	// 写入响应体
	fmt.Fprintln(w, "欢迎来到 Go Web 世界!")
}

// aboutHandler 处理关于页面
func aboutHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[%s] %s", r.Method, r.URL.Path)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintln(w, "<h1>关于我们</h1>")
	fmt.Fprintln(w, "<p>这是一个 Go Web 学习项目</p>")
}
