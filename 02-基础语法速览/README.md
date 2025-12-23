# 02 - 基础语法速览

## 学习目标

快速掌握 Go 语言核心语法，理解 Go 与 Java 在设计理念上的差异。

---

## 学习要点

### 1. 变量声明

```go
// var 声明（显式类型）
var name string = "Tom"
var age int        // 零值：0

// var 声明（类型推断）
var count = 10     // 推断为 int

// 短声明 :=（推荐，仅函数内可用）
name := "Tom"
age := 18

// 常量
const Pi = 3.14159
const MaxRetries = 3
```

**📌 var vs := 如何选择？**

| 场景 | 推荐 | 原因 |
|------|------|------|
| 函数内部声明 | `:=` | 简洁，Go 社区惯例 |
| 包级变量 | `var` | `:=` 语法不允许在包级使用 |
| 显式指定类型 | `var x int64 = 10` | 类型推断默认 int，需要 int64 时显式声明 |
| 声明零值变量 | `var s string` | 比 `s := ""` 更清晰地表达"我要零值" |

**🆚 与 Java 对比**

```java
// Java: 必须声明类型，或用 var（Java 10+）
String name = "Tom";
var age = 18;  // Java 10+ 局部变量类型推断
```

```go
// Go: := 更简洁，且从 Go 1.0 就支持
name := "Tom"
age := 18
```

> **洞察**：Go 的 `:=` 让代码更简洁，但仅限函数内使用。这是刻意的设计——包级变量应该更明确，需要思考后才声明。

---

### 2. 基本类型与零值

| Go 类型 | Java 等价 | 零值 | 说明 |
|---------|-----------|------|------|
| `string` | `String` | `""` | 不可变，UTF-8 编码 |
| `int` | `int/long` | `0` | 平台相关（32/64位），不确定时用 `int64` |
| `int64` | `long` | `0` | 明确 64 位 |
| `float64` | `double` | `0.0` | 默认浮点类型 |
| `bool` | `boolean` | `false` | |
| `[]T` (slice) | `List<T>` | `nil` | 动态数组，最常用集合 |
| `map[K]V` | `Map<K,V>` | `nil` | 哈希表 |
| `*T` (pointer) | 引用 | `nil` | Go 有显式指针，Java 没有 |

**📌 零值设计理念**

Go 的变量永远有值，未初始化时为"零值"。这不是 bug，是特性：

```go
var count int      // 0，可以直接 count++
var name string    // ""，可以直接 name += "hello"
var ok bool        // false，可以直接用于条件判断

// Java 中 null 引发的问题，Go 通过零值避免了
var sb strings.Builder  // 零值可用，无需 new
sb.WriteString("hello") // 直接使用，不会 NPE
```

**📌 int vs int64 如何选择？**

```go
// 一般情况：用 int
for i := 0; i < 100; i++ { }  // 循环变量用 int
len(slice)  // 返回 int

// 明确需要 64 位时：用 int64
var userID int64 = 1234567890123  // 数据库 ID、时间戳
var fileSize int64                 // 文件大小
```

---

### 3. 函数定义

```go
// 基本函数
func add(a, b int) int {
    return a + b
}

// 多返回值（Go 特色，Java 需要包装类）
func divide(a, b int) (int, error) {
    if b == 0 {
        return 0, errors.New("除数不能为零")
    }
    return a / b, nil
}

// 调用多返回值
result, err := divide(10, 2)
if err != nil {
    // 处理错误
}
```

**🆚 与 Java 对比**

```java
// Java: 需要定义 Result 类或使用 Optional
public class DivideResult {
    int value;
    String error;
}

// 或者抛异常
public int divide(int a, int b) throws ArithmeticException {
    if (b == 0) throw new ArithmeticException();
    return a / b;
}
```

```go
// Go: 多返回值，简洁直接
func divide(a, b int) (int, error) {
    if b == 0 {
        return 0, errors.New("除数不能为零")
    }
    return a / b, nil
}
```

**📌 命名返回值（谨慎使用）**

```go
// 命名返回值 + 裸返回
func split(sum int) (x, y int) {
    x = sum * 4 / 9
    y = sum - x
    return  // 裸返回，返回 x 和 y
}
```

> **⚠️ 最佳实践**：命名返回值在短函数中可用于文档目的，但**裸返回（naked return）在长函数中会降低可读性**，Go 官方也建议避免在长函数中使用。

---

### 4. 结构体与方法

```go
// 定义结构体（相当于 Java 的 class，但没有继承）
type User struct {
    ID       int64   // 公开字段（大写开头）
    Username string
    Email    string
    password string  // 私有字段（小写开头，仅包内可见）
}

// 构造函数（Go 惯例：NewXxx）
func NewUser(id int64, username, email string) *User {
    return &User{
        ID:       id,
        Username: username,
        Email:    email,
    }
}

// 值接收者方法（不修改原对象）
func (u User) DisplayName() string {
    return u.Username
}

// 指针接收者方法（可以修改原对象）
func (u *User) SetEmail(email string) {
    u.Email = email
}
```

**📌 值接收者 vs 指针接收者**

| 场景 | 推荐 | 原因 |
|------|------|------|
| 需要修改结构体 | `*T` 指针 | 值传递会拷贝，修改不影响原对象 |
| 结构体较大 | `*T` 指针 | 避免拷贝开销 |
| 需要一致性 | 统一用 `*T` | 如果有一个方法用指针，建议全部用指针 |
| 小结构体只读 | `T` 值 | 如 `time.Time`，不可变语义更清晰 |

**🆚 与 Java 对比**

```java
// Java: 类 + 方法
public class User {
    private String username;

    public String getUsername() {
        return username;
    }

    public void setUsername(String username) {
        this.username = username;  // Java 对象总是"引用"
    }
}
```

```go
// Go: 结构体 + 方法
type User struct {
    Username string
}

func (u User) GetUsername() string {
    return u.Username
}

func (u *User) SetUsername(name string) {
    u.Username = name  // 必须用指针接收者才能修改
}
```

> **洞察**：Java 的对象变量本质是引用（指针），所以 `user.setName()` 总能修改。Go 的方法接收者默认是值传递，必须显式用 `*User` 才能修改——这让代码意图更清晰。

---

### 5. 接口（隐式实现）

```go
// 定义接口
type Stringer interface {
    String() string
}

// User 实现 Stringer（无需 implements 关键字！）
func (u User) String() string {
    return fmt.Sprintf("User{ID: %d, Name: %s}", u.ID, u.Username)
}

// 使用
var s Stringer = User{ID: 1, Username: "tom"}
fmt.Println(s.String())
```

**🆚 与 Java 对比**

```java
// Java: 显式 implements
public interface Stringer {
    String toString();
}

public class User implements Stringer {  // 必须声明
    @Override
    public String toString() { ... }
}
```

```go
// Go: 隐式实现，只要方法签名匹配就实现了接口
type Stringer interface {
    String() string
}

func (u User) String() string { ... }  // 自动实现 Stringer
```

**📌 隐式实现的好处**

1. **解耦**：实现者不需要知道接口的存在，可以后期添加接口
2. **组合**：容易让已有类型实现新接口（Java 需要修改类定义或用适配器）
3. **小接口**：鼓励定义小而精的接口（1-2 个方法），而非大而全

> **Go 谚语**："Accept interfaces, return structs"（接受接口，返回结构体）

---

### 6. 错误处理（Go 1.13+ 最佳实践）

```go
// Go 没有异常，用 error 接口表示错误
func doSomething() error {
    return nil  // 成功返回 nil
}

// 调用时检查错误
result, err := doSomething()
if err != nil {
    return err  // 向上传递
}
```

**📌 Go 1.13+ 错误包装与检查**

```go
import "errors"

// 预定义错误（哨兵错误）
var ErrNotFound = errors.New("资源不存在")

// 包装错误（添加上下文）
func findUser(id int) (*User, error) {
    user, err := db.Query(id)
    if err != nil {
        return nil, fmt.Errorf("查找用户 %d 失败: %w", id, err)  // %w 包装
    }
    return user, nil
}

// 检查错误链
if errors.Is(err, ErrNotFound) {
    // 处理"不存在"错误
}

// 提取特定错误类型
var validErr *ValidationError
if errors.As(err, &validErr) {
    fmt.Println(validErr.Field)  // 访问错误详情
}
```

**🆚 与 Java 对比**

| 对比项 | Java | Go |
|--------|------|-----|
| 机制 | 异常（try-catch） | 返回值（error） |
| 传播 | 自动向上抛出 | 显式 return err |
| 检查 | instanceof / catch 类型 | errors.Is / errors.As |
| 包装 | new Exception(cause) | fmt.Errorf("%w", err) |

```java
// Java: 异常可能被忽略
User user = userService.findById(id);  // 可能抛 RuntimeException
```

```go
// Go: 错误必须处理
user, err := userService.FindByID(id)
if err != nil {
    // 必须处理，否则编译器警告（unused variable）
}
```

> **洞察**：Go 强制你在每个调用点思考"失败了怎么办"，代码更健壮但也更啰嗦。Java 的异常让代码更简洁，但容易忽略错误处理。

---

### 7. 指针基础

```go
x := 10
p := &x      // p 是指向 x 的指针，类型为 *int
*p = 20      // 通过指针修改 x 的值
fmt.Println(x) // 输出 20
```

**📌 为什么 Go 需要指针？**

Go 是值传递语言，函数参数都是拷贝：

```go
func double(n int) {
    n = n * 2  // 修改的是拷贝，原值不变
}

x := 10
double(x)
fmt.Println(x)  // 仍然是 10

// 用指针才能修改
func doublePtr(n *int) {
    *n = *n * 2
}
doublePtr(&x)
fmt.Println(x)  // 现在是 20
```

**🆚 与 Java 对比**

```java
// Java: 对象变量本身就是引用（隐式指针）
User user = new User();
modifyUser(user);  // 可以修改 user 的字段

// Java 基本类型是值传递
void doubleValue(int n) {
    n = n * 2;  // 不影响原值
}
```

```go
// Go: 所有类型默认值传递，需要指针才能修改
func modifyUser(u *User) {
    u.Name = "new"  // 用指针才能修改
}
```

> **洞察**：Java 程序员容易困惑——为什么 Go 的方法参数不能修改原对象？因为 Java 对象变量隐含了指针语义，而 Go 把这层抽象暴露出来了。显式比隐式更清晰。

---

## 示例代码

| 目录 | 内容 | 重点 |
|------|------|------|
| `examples/01-variables/` | 变量声明与类型 | 零值、slice、map |
| `examples/02-functions/` | 函数定义 | 多返回值 |
| `examples/03-structs/` | 结构体与方法 | 构造函数、嵌入 |
| `examples/04-interfaces/` | 接口 | 隐式实现 |
| `examples/05-errors/` | 错误处理 | errors.Is/As、Go 1.20+ errors.Join |

---

## 作业任务

### 任务描述
完成 `homework/main.go`，实现一个简单的用户查找功能。

### 要求
1. 定义 `User` 结构体，包含 `ID`、`Name`、`Email` 字段
2. 实现 `String()` 方法（满足 `fmt.Stringer` 接口）
3. 编写 `FindUser(id int64) (*User, error)` 函数：
   - id > 0：返回模拟用户
   - id <= 0：返回 `ErrInvalidID` 错误
4. 使用 `errors.Is` 检查错误

### 验收标准
```bash
cd homework && go run main.go
```
输出应包含：成功查找用户、错误处理演示。

---

## 参考资料
- [Effective Go](https://go.dev/doc/effective_go)
- [Go 语言圣经](https://gopl.io/)
- [Go 错误处理最佳实践](https://go.dev/blog/go1.13-errors)
