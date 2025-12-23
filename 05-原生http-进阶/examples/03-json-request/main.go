// 03-json-request: JSON 请求体解析
//
// 📌 关键步骤:
//   1. 使用 json.NewDecoder(r.Body).Decode() 解析
//   2. 记得处理解析错误
//
// 📌 最佳实践:
//   - 限制请求体大小防止攻击
//   - 验证必填字段
//   - 使用指针字段区分"未传"和"空值"
//
// 运行: go run main.go
// 测试:
//   curl -X POST http://localhost:8080/users \
//     -H "Content-Type: application/json" \
//     -d '{"username":"tom","email":"tom@example.com"}'
package main

import (
	"encoding/json"
	"log"
	"net/http"
)

// CreateUserRequest 创建用户请求
type CreateUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Age      *int   `json:"age"` // 指针可区分 0 和未传
}

// Response 统一响应
type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /users", createUserHandler)

	addr := ":8080"
	log.Printf("服务器启动: http://localhost%s", addr)
	log.Println("测试命令:")
	log.Println(`  curl -X POST http://localhost:8080/users -H "Content-Type: application/json" -d '{"username":"tom","email":"tom@example.com"}'`)

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}

func createUserHandler(w http.ResponseWriter, r *http.Request) {
	// 1. 限制请求体大小（防止大请求攻击）
	r.Body = http.MaxBytesReader(w, r.Body, 1024*1024) // 1MB

	// 2. 解析 JSON
	var req CreateUserRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields() // 不允许未知字段（可选，更严格）

	if err := decoder.Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, Response{
			Code:    400,
			Message: "JSON 解析失败: " + err.Error(),
		})
		return
	}

	// 3. 验证必填字段
	if req.Username == "" {
		writeJSON(w, http.StatusBadRequest, Response{
			Code:    400,
			Message: "username 不能为空",
		})
		return
	}

	if req.Email == "" {
		writeJSON(w, http.StatusBadRequest, Response{
			Code:    400,
			Message: "email 不能为空",
		})
		return
	}

	// 4. 处理可选字段
	age := 0
	if req.Age != nil {
		age = *req.Age
	}

	// 5. 返回成功响应
	writeJSON(w, http.StatusCreated, Response{
		Code:    0,
		Message: "创建成功",
		Data: map[string]any{
			"id":       1,
			"username": req.Username,
			"email":    req.Email,
			"age":      age,
		},
	})
}

func writeJSON(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}
