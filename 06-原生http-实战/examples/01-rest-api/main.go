// 01-rest-api: 完整的用户 CRUD API（无框架）
//
// 📌 RESTful 设计:
//   GET    /users      - 列表
//   GET    /users/{id} - 详情
//   POST   /users      - 创建
//   PUT    /users/{id} - 更新
//   DELETE /users/{id} - 删除
//
// 📌 最佳实践:
//   - 使用 sync.RWMutex 保证并发安全
//   - 统一响应格式
//   - 合理的 HTTP 状态码
//
// 运行: go run main.go
// 测试见 README.md
package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

// ==================== 模型定义 ====================

type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

type UpdateUserRequest struct {
	Username string `json:"username,omitempty"`
	Email    string `json:"email,omitempty"`
}

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// ==================== 存储层 ====================

type UserStore struct {
	mu     sync.RWMutex
	users  map[int]*User
	nextID int
}

func NewUserStore() *UserStore {
	return &UserStore{
		users:  make(map[int]*User),
		nextID: 1,
	}
}

func (s *UserStore) List() []*User {
	s.mu.RLock()
	defer s.mu.RUnlock()

	users := make([]*User, 0, len(s.users))
	for _, u := range s.users {
		users = append(users, u)
	}
	return users
}

func (s *UserStore) Get(id int) (*User, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	user, ok := s.users[id]
	return user, ok
}

func (s *UserStore) Create(username, email string) *User {
	s.mu.Lock()
	defer s.mu.Unlock()

	user := &User{
		ID:        s.nextID,
		Username:  username,
		Email:     email,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	s.users[s.nextID] = user
	s.nextID++
	return user
}

func (s *UserStore) Update(id int, username, email string) (*User, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	user, ok := s.users[id]
	if !ok {
		return nil, false
	}

	if username != "" {
		user.Username = username
	}
	if email != "" {
		user.Email = email
	}
	user.UpdatedAt = time.Now()
	return user, true
}

func (s *UserStore) Delete(id int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.users[id]; !ok {
		return false
	}
	delete(s.users, id)
	return true
}

// ==================== Handler ====================

type UserHandler struct {
	store *UserStore
}

func NewUserHandler(store *UserStore) *UserHandler {
	return &UserHandler{store: store}
}

// GET /users - 用户列表
func (h *UserHandler) List(w http.ResponseWriter, r *http.Request) {
	users := h.store.List()
	writeJSON(w, http.StatusOK, Response{
		Code:    0,
		Message: "success",
		Data:    users,
	})
}

// GET /users/{id} - 用户详情
func (h *UserHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, Response{
			Code:    400,
			Message: "无效的用户 ID",
		})
		return
	}

	user, ok := h.store.Get(id)
	if !ok {
		writeJSON(w, http.StatusNotFound, Response{
			Code:    404,
			Message: "用户不存在",
		})
		return
	}

	writeJSON(w, http.StatusOK, Response{
		Code:    0,
		Message: "success",
		Data:    user,
	})
}

// POST /users - 创建用户
func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, Response{
			Code:    400,
			Message: "无效的请求体",
		})
		return
	}

	if req.Username == "" || req.Email == "" {
		writeJSON(w, http.StatusBadRequest, Response{
			Code:    400,
			Message: "用户名和邮箱不能为空",
		})
		return
	}

	user := h.store.Create(req.Username, req.Email)
	writeJSON(w, http.StatusCreated, Response{
		Code:    0,
		Message: "创建成功",
		Data:    user,
	})
}

// PUT /users/{id} - 更新用户
func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, Response{
			Code:    400,
			Message: "无效的用户 ID",
		})
		return
	}

	var req UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, Response{
			Code:    400,
			Message: "无效的请求体",
		})
		return
	}

	user, ok := h.store.Update(id, req.Username, req.Email)
	if !ok {
		writeJSON(w, http.StatusNotFound, Response{
			Code:    404,
			Message: "用户不存在",
		})
		return
	}

	writeJSON(w, http.StatusOK, Response{
		Code:    0,
		Message: "更新成功",
		Data:    user,
	})
}

// DELETE /users/{id} - 删除用户
func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, Response{
			Code:    400,
			Message: "无效的用户 ID",
		})
		return
	}

	if !h.store.Delete(id) {
		writeJSON(w, http.StatusNotFound, Response{
			Code:    404,
			Message: "用户不存在",
		})
		return
	}

	writeJSON(w, http.StatusOK, Response{
		Code:    0,
		Message: "删除成功",
	})
}

// ==================== 辅助函数 ====================

func writeJSON(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

// ==================== 主函数 ====================

func main() {
	store := NewUserStore()
	handler := NewUserHandler(store)

	// 预置一些数据
	store.Create("tom", "tom@example.com")
	store.Create("jerry", "jerry@example.com")

	mux := http.NewServeMux()

	// Go 1.22+ 路由语法
	mux.HandleFunc("GET /users", handler.List)
	mux.HandleFunc("GET /users/{id}", handler.Get)
	mux.HandleFunc("POST /users", handler.Create)
	mux.HandleFunc("PUT /users/{id}", handler.Update)
	mux.HandleFunc("DELETE /users/{id}", handler.Delete)

	// 健康检查
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, Response{Code: 0, Message: "ok"})
	})

	addr := ":8080"
	log.Printf("服务器启动: http://localhost%s", addr)
	log.Println("API 列表:")
	log.Println("  GET    /users      - 用户列表")
	log.Println("  GET    /users/{id} - 用户详情")
	log.Println("  POST   /users      - 创建用户")
	log.Println("  PUT    /users/{id} - 更新用户")
	log.Println("  DELETE /users/{id} - 删除用户")

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}
