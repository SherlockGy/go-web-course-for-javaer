# 06 - 原生 HTTP 实战

## 学习目标

用原生 `net/http` 实现完整的 RESTful API，不使用任何框架。

---

## 学习要点

### 1. RESTful API 设计

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /todos | 获取所有待办事项 |
| POST | /todos | 创建待办事项 |
| GET | /todos/{id} | 获取单个待办事项 |
| PUT | /todos/{id} | 更新待办事项 |
| DELETE | /todos/{id} | 删除待办事项 |

### 2. 内存数据存储

```go
import "sync"

type Todo struct {
    ID        int64  `json:"id"`
    Title     string `json:"title"`
    Completed bool   `json:"completed"`
}

// 线程安全的存储
var (
    todos   = make(map[int64]*Todo)
    todosMu sync.RWMutex
    nextID  int64 = 1
)

// 读取（加读锁）
func getTodo(id int64) (*Todo, bool) {
    todosMu.RLock()
    defer todosMu.RUnlock()
    todo, ok := todos[id]
    return todo, ok
}

// 写入（加写锁）
func addTodo(todo *Todo) {
    todosMu.Lock()
    defer todosMu.Unlock()
    todo.ID = nextID
    nextID++
    todos[todo.ID] = todo
}
```

### 3. 统一响应格式

```go
type ApiResponse struct {
    Code    int         `json:"code"`
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`
}

func success(w http.ResponseWriter, data interface{}) {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(ApiResponse{
        Code:    200,
        Message: "success",
        Data:    data,
    })
}

func errorResponse(w http.ResponseWriter, code int, message string) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(code)
    json.NewEncoder(w).Encode(ApiResponse{
        Code:    code,
        Message: message,
    })
}
```

### 4. 完整处理器示例

```go
// GET /todos - 列表
func listTodos(w http.ResponseWriter, r *http.Request) {
    todosMu.RLock()
    defer todosMu.RUnlock()

    list := make([]*Todo, 0, len(todos))
    for _, todo := range todos {
        list = append(list, todo)
    }
    success(w, list)
}

// POST /todos - 创建
func createTodo(w http.ResponseWriter, r *http.Request) {
    var todo Todo
    if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
        errorResponse(w, 400, "无效的请求体")
        return
    }
    addTodo(&todo)
    w.WriteHeader(http.StatusCreated)
    success(w, todo)
}

// GET /todos/{id} - 获取单个
func getTodoHandler(w http.ResponseWriter, r *http.Request) {
    id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
    if err != nil {
        errorResponse(w, 400, "无效的 ID")
        return
    }

    todo, ok := getTodo(id)
    if !ok {
        errorResponse(w, 404, "待办事项不存在")
        return
    }
    success(w, todo)
}
```

---

## 示例代码

### examples/01-rest-api/
完整的用户 CRUD API（无框架）

---

## 作业任务

### 任务描述
用原生 `net/http` 实现 Todo List API。

### 要求
1. 实现以下 5 个接口：
   - `GET /todos` - 获取所有待办事项
   - `POST /todos` - 创建待办事项
   - `GET /todos/{id}` - 获取单个待办事项
   - `PUT /todos/{id}` - 更新待办事项
   - `DELETE /todos/{id}` - 删除待办事项

2. 使用内存存储（map + sync.RWMutex）

3. 统一响应格式：
   ```json
   {
     "code": 200,
     "message": "success",
     "data": {}
   }
   ```

4. 添加日志中间件

### 数据结构
```go
type Todo struct {
    ID        int64  `json:"id"`
    Title     string `json:"title"`
    Completed bool   `json:"completed"`
}
```

### 验收标准
```bash
# 创建
curl -X POST http://localhost:8080/todos \
  -H "Content-Type: application/json" \
  -d '{"title": "学习 Go", "completed": false}'
# 返回: {"code":200,"message":"success","data":{"id":1,"title":"学习 Go","completed":false}}

# 列表
curl http://localhost:8080/todos
# 返回: {"code":200,"message":"success","data":[{"id":1,"title":"学习 Go","completed":false}]}

# 获取单个
curl http://localhost:8080/todos/1
# 返回: {"code":200,"message":"success","data":{"id":1,"title":"学习 Go","completed":false}}

# 更新
curl -X PUT http://localhost:8080/todos/1 \
  -H "Content-Type: application/json" \
  -d '{"title": "学习 Go", "completed": true}'
# 返回: {"code":200,"message":"success","data":{"id":1,"title":"学习 Go","completed":true}}

# 删除
curl -X DELETE http://localhost:8080/todos/1
# 返回: {"code":200,"message":"success"}

# 验证删除
curl http://localhost:8080/todos/1
# 返回: {"code":404,"message":"待办事项不存在"}
```

### 代码结构建议
```
homework/
├── main.go           # 入口 + 路由注册
├── handler.go        # HTTP 处理器
├── store.go          # 内存存储
└── response.go       # 统一响应
```

---

## 思考题

1. 为什么需要 `sync.RWMutex`？不加锁会怎样？
2. 原生实现和框架的主要差距在哪里？
3. 如果要持久化数据，应该改哪些代码？

---

## 参考资料
- [RESTful API 设计指南](https://restfulapi.net/)
- [sync.RWMutex 文档](https://pkg.go.dev/sync#RWMutex)
