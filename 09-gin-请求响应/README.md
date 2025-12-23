# 09 - Gin 请求响应

## 学习目标

掌握 Gin 的请求绑定、参数验证和响应处理。

---

## 学习要点

### 1. JSON 绑定

```go
type CreateUserRequest struct {
    Name  string `json:"name"`
    Email string `json:"email"`
}

func createUser(c *gin.Context) {
    var req CreateUserRequest

    // ShouldBindJSON - 绑定失败返回错误，不终止请求
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }

    c.JSON(200, gin.H{"data": req})
}
```

### 2. 表单绑定

```go
type LoginForm struct {
    Username string `form:"username"`
    Password string `form:"password"`
}

func login(c *gin.Context) {
    var form LoginForm

    // ShouldBind 自动识别 Content-Type
    if err := c.ShouldBind(&form); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }

    c.JSON(200, gin.H{"username": form.Username})
}
```

### 3. 参数验证（binding tag）

```go
type RegisterRequest struct {
    // required: 必填
    // min/max: 长度限制
    // email: 邮箱格式
    Username string `json:"username" binding:"required,min=3,max=20"`
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=6,max=50"`
    Age      int    `json:"age" binding:"gte=0,lte=120"`  // 数值范围
}

func register(c *gin.Context) {
    var req RegisterRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    // 验证通过...
}
```

### 4. 常用验证规则

| 规则 | 说明 | 示例 |
|------|------|------|
| `required` | 必填 | `binding:"required"` |
| `min` | 最小长度/值 | `binding:"min=3"` |
| `max` | 最大长度/值 | `binding:"max=20"` |
| `len` | 精确长度 | `binding:"len=11"` |
| `email` | 邮箱格式 | `binding:"email"` |
| `url` | URL 格式 | `binding:"url"` |
| `gte` | 大于等于 | `binding:"gte=0"` |
| `lte` | 小于等于 | `binding:"lte=100"` |
| `oneof` | 枚举值 | `binding:"oneof=male female"` |

### 5. 文件上传

```go
// 单文件
func uploadFile(c *gin.Context) {
    file, err := c.FormFile("file")
    if err != nil {
        c.JSON(400, gin.H{"error": "文件上传失败"})
        return
    }

    // 保存文件
    dst := "./uploads/" + file.Filename
    c.SaveUploadedFile(file, dst)

    c.JSON(200, gin.H{"filename": file.Filename})
}

// 多文件
func uploadFiles(c *gin.Context) {
    form, _ := c.MultipartForm()
    files := form.File["files"]

    for _, file := range files {
        c.SaveUploadedFile(file, "./uploads/"+file.Filename)
    }

    c.JSON(200, gin.H{"count": len(files)})
}
```

### 6. 多种响应格式

```go
// JSON
c.JSON(200, gin.H{"message": "ok"})

// XML
c.XML(200, gin.H{"message": "ok"})

// String
c.String(200, "Hello, %s", name)

// HTML
c.HTML(200, "index.html", gin.H{"title": "首页"})

// 文件
c.File("./files/report.pdf")

// 重定向
c.Redirect(302, "/login")
```

### 7. 统一响应结构

```go
type Response struct {
    Code    int         `json:"code"`
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`
}

func Success(c *gin.Context, data interface{}) {
    c.JSON(200, Response{
        Code:    200,
        Message: "success",
        Data:    data,
    })
}

func Error(c *gin.Context, code int, message string) {
    c.JSON(code, Response{
        Code:    code,
        Message: message,
    })
}

// 使用
func handler(c *gin.Context) {
    user := getUser()
    Success(c, user)
}
```

---

## 示例代码

### examples/01-json-binding/
JSON 请求绑定

### examples/02-validation/
参数验证规则

### examples/03-file-upload/
文件上传处理

### examples/04-unified-response/
统一响应封装

---

## 作业任务

### 任务描述
实现用户注册接口，包含完整的参数验证。

### 请求格式
```json
POST /api/register
{
    "username": "tom",
    "email": "tom@example.com",
    "password": "123456",
    "confirm_password": "123456"
}
```

### 验证规则
- `username`: 必填，3-20 字符
- `email`: 必填，合法邮箱格式
- `password`: 必填，6-50 字符
- `confirm_password`: 必填，与 password 相同

### 响应格式
**成功**：
```json
{
    "code": 200,
    "message": "注册成功",
    "data": {
        "id": 1,
        "username": "tom",
        "email": "tom@example.com"
    }
}
```

**失败**：
```json
{
    "code": 400,
    "message": "用户名长度必须在 3-20 字符之间"
}
```

### 验收标准
```bash
# 成功注册
curl -X POST http://localhost:8080/api/register \
  -H "Content-Type: application/json" \
  -d '{"username":"tom","email":"tom@example.com","password":"123456","confirm_password":"123456"}'
# 返回: 200

# 用户名太短
curl -X POST http://localhost:8080/api/register \
  -H "Content-Type: application/json" \
  -d '{"username":"ab","email":"tom@example.com","password":"123456","confirm_password":"123456"}'
# 返回: 400

# 邮箱格式错误
curl -X POST http://localhost:8080/api/register \
  -H "Content-Type: application/json" \
  -d '{"username":"tom","email":"invalid","password":"123456","confirm_password":"123456"}'
# 返回: 400
```

### 提示
密码确认可以用自定义验证器：
```go
// 自定义验证
if req.Password != req.ConfirmPassword {
    Error(c, 400, "两次密码不一致")
    return
}
```

---

## 参考资料
- [Gin 模型绑定与验证](https://gin-gonic.com/docs/examples/binding-and-validation/)
- [validator 包文档](https://pkg.go.dev/github.com/go-playground/validator/v10)
