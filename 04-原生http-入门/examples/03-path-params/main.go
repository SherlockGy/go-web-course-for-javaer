// 03-path-params: 路径参数（Go 1.22+）
//
// 📌 Go 1.22 重大更新:
//   - 路径参数语法: /users/{id}
//   - 使用 r.PathValue("id") 获取参数值
//   - 通配符: /files/{path...} 匹配剩余路径
//
// 这是 Go 1.22 之前需要第三方库才能实现的功能！
//
// 运行: go run main.go
// 测试:
//   curl http://localhost:8080/users/123
//   curl http://localhost:8080/users/456/orders
//   curl http://localhost:8080/files/docs/readme.md
package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	// 路径参数 {id}
	mux.HandleFunc("GET /users/{id}", getUserHandler)

	// 多个路径参数
	mux.HandleFunc("GET /users/{userId}/orders/{orderId}", getOrderHandler)

	// 通配符 {path...} 匹配剩余所有路径
	mux.HandleFunc("GET /files/{path...}", getFileHandler)

	// 精确匹配优先
	mux.HandleFunc("GET /users/me", getCurrentUserHandler) // 优先于 /users/{id}

	addr := ":8080"
	log.Printf("服务器启动: http://localhost%s", addr)
	log.Println("测试命令:")
	log.Println("  curl http://localhost:8080/users/123")
	log.Println("  curl http://localhost:8080/users/me")
	log.Println("  curl http://localhost:8080/users/1/orders/100")
	log.Println("  curl http://localhost:8080/files/docs/readme.md")

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}

func getUserHandler(w http.ResponseWriter, r *http.Request) {
	// Go 1.22+ 新方法: PathValue
	id := r.PathValue("id")
	fmt.Fprintf(w, "获取用户: id=%s\n", id)
}

func getOrderHandler(w http.ResponseWriter, r *http.Request) {
	userId := r.PathValue("userId")
	orderId := r.PathValue("orderId")
	fmt.Fprintf(w, "获取订单: userId=%s, orderId=%s\n", userId, orderId)
}

func getFileHandler(w http.ResponseWriter, r *http.Request) {
	// {path...} 匹配剩余所有路径
	path := r.PathValue("path")
	fmt.Fprintf(w, "获取文件: path=%s\n", path)
}

func getCurrentUserHandler(w http.ResponseWriter, r *http.Request) {
	// 精确路由优先于参数路由
	fmt.Fprintln(w, "获取当前登录用户")
}
