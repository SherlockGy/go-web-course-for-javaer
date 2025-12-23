# 02 - 基础语法速览

## 学习目标

快速掌握 Go 语言核心语法，为 Web 开发打基础。

---

## 学习要点

### 1. 变量声明
```go
// var 声明
var name string = "Tom"
var age int        // 零值：0

// 短声明（函数内部）
name := "Tom"
age := 18

// 常量
const Pi = 3.14159
```

### 2. 基本类型
| 类型 | 说明 | 零值 |
|------|------|------|
| `string` | 字符串 | `""` |
| `int`, `int64` | 整数 | `0` |
| `float64` | 浮点数 | `0.0` |
| `bool` | 布尔 | `false` |
| `[]T` | 切片 | `nil` |
| `map[K]V` | 映射 | `nil` |

### 3. 函数定义
```go
// 基本函数
func add(a, b int) int {
    return a + b
}

// 多返回值
func divide(a, b int) (int, error) {
    if b == 0 {
        return 0, errors.New("除数不能为零")
    }
    return a / b, nil
}

// 命名返回值
func split(sum int) (x, y int) {
    x = sum * 4 / 9
    y = sum - x
    return  // 裸返回
}
```

### 4. 结构体与方法
```go
// 定义结构体
type User struct {
    ID       int64
    Username string
    Email    string
}

// 值接收者方法
func (u User) GetName() string {
    return u.Username
}

// 指针接收者方法（可修改结构体）
func (u *User) SetName(name string) {
    u.Username = name
}
```

### 5. 接口（隐式实现）
```go
// 定义接口
type Stringer interface {
    String() string
}

// User 实现 Stringer（无需显式声明）
func (u User) String() string {
    return fmt.Sprintf("User{ID: %d, Name: %s}", u.ID, u.Username)
}
```

### 6. 错误处理
```go
// Go 没有异常，用 error 接口表示错误
func doSomething() error {
    // 成功返回 nil
    return nil
}

// 调用时检查错误
result, err := doSomething()
if err != nil {
    // 处理错误
    return err
}
```

> **🆚 Java 对比**
> ```java
> // Java: 异常会向上传播，调用者可能不知道会抛什么异常
> User user = userService.findById(id);  // 可能抛 RuntimeException
>
> // Go: 错误必须显式处理，代码即文档
> user, err := userService.FindByID(id)  // 签名告诉你：可能返回错误
> if err != nil { ... }
> ```
> **洞察**：Go 用编译器强制你思考"如果失败了怎么办"，Java 让你用 try-catch "以后再说"。

### 7. 指针基础
```go
x := 10
p := &x      // p 是指向 x 的指针
*p = 20      // 通过指针修改 x 的值
fmt.Println(x) // 输出 20
```

---

## 示例代码

### examples/01-variables/
变量声明与类型演示

### examples/02-functions/
函数定义与多返回值

### examples/03-structs/
结构体与方法

### examples/04-interfaces/
接口定义与隐式实现

### examples/05-errors/
错误处理模式

---

## 作业任务

### 任务描述
定义一个 `User` 结构体，实现 `Stringer` 接口，编写函数返回 `(User, error)`。

### 要求
1. 定义 `User` 结构体，包含 `ID`、`Name`、`Email` 字段
2. 实现 `String()` 方法，返回格式化的用户信息
3. 编写 `FindUser(id int64) (User, error)` 函数
   - 如果 id > 0，返回模拟用户
   - 如果 id <= 0，返回错误

### 验收标准
```go
user, err := FindUser(1)
if err != nil {
    fmt.Println("错误:", err)
} else {
    fmt.Println(user)  // 调用 String() 方法
}
```

### 对比思考
Go 的错误处理 vs Java 的异常：
- Go：显式检查，`if err != nil`
- Java：隐式抛出，`try-catch`

哪种方式更好？思考各自的优缺点。

---

## 参考资料
- [Effective Go](https://go.dev/doc/effective_go)
- [Go 语言圣经](https://gopl.io/)
