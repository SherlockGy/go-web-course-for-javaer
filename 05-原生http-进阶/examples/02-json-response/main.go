// 02-json-response: JSON 响应
//
// 📌 关键步骤:
//   1. 设置 Content-Type: application/json
//   2. 使用 json.NewEncoder(w).Encode() 写入响应
//
// 📌 最佳实践:
//   - 统一响应格式 {code, message, data}
//   - 使用结构体 tag 控制 JSON 字段名
//   - 敏感字段使用 json:"-" 跳过
//
// 运行: go run main.go
// 测试: curl http://localhost:8080/users
package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

// User 用户模型
type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"-"` // 不序列化到 JSON
	CreatedAt time.Time `json:"created_at"`
}

// Response 统一响应格式
type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"` // omitempty: 空值不输出
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /users", listUsersHandler)
	mux.HandleFunc("GET /users/{id}", getUserHandler)
	mux.HandleFunc("GET /error", errorHandler)

	addr := ":8080"
	log.Printf("服务器启动: http://localhost%s", addr)

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}

// 返回用户列表
func listUsersHandler(w http.ResponseWriter, r *http.Request) {
	users := []User{
		{ID: 1, Username: "tom", Email: "tom@example.com", Password: "secret", CreatedAt: time.Now()},
		{ID: 2, Username: "jerry", Email: "jerry@example.com", Password: "secret", CreatedAt: time.Now()},
	}

	// 使用统一响应格式
	writeJSON(w, http.StatusOK, Response{
		Code:    0,
		Message: "success",
		Data:    users,
	})
}

// 返回单个用户
func getUserHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	// 模拟数据
	user := User{
		ID:        1,
		Username:  "tom",
		Email:     "tom@example.com",
		Password:  "secret123", // 不会出现在 JSON 中
		CreatedAt: time.Now(),
	}

	writeJSON(w, http.StatusOK, Response{
		Code:    0,
		Message: "success",
		Data: map[string]any{
			"user":       user,
			"request_id": id,
		},
	})
}

// 返回错误响应
func errorHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusNotFound, Response{
		Code:    404,
		Message: "资源不存在",
		// Data 为空时不会输出（omitempty）
	})
}

// writeJSON 写入 JSON 响应（通用函数）
func writeJSON(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("JSON 编码失败: %v", err)
	}
}
