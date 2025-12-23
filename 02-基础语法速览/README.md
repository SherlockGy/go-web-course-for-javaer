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

Go 的变量永远有值，未初始化时为"零值"。这不是 bug，是特性——**让很多类型"开箱即用"，无需初始化**。

```go
var count int      // 0，可以直接 count++
var name string    // ""，可以直接 name += "hello"
var ok bool        // false，可以直接用于条件判断
```

**🆚 与 Java 对比：null 的噩梦 vs 零值的安全**

```java
// Java: null 是万恶之源
StringBuilder sb = null;
sb.append("hello");  // 💥 NullPointerException!

// Java 需要防御性编程
if (sb != null) {
    sb.append("hello");
}
```

```go
// Go: 很多类型的零值可以直接使用
var sb strings.Builder  // 零值可用，无需 new 或 make
sb.WriteString("hello") // ✅ 直接使用，不会 panic
sb.WriteString(" world")
fmt.Println(sb.String()) // "hello world"
```

> **洞察**：Java 的 `null` 表示"没有对象"，但程序员经常忘记检查。Go 的零值设计让很多类型在未初始化时就能安全使用，从根本上减少了空指针类问题。

---

#### 2.1 零值最佳实践（先记住这个）

**📌 一句话总结：只有 map、channel、指针 需要初始化，其他放心用零值。**

```go
// ✅ 这些零值直接用，不用想
var count int              // 0，直接 ++
var name string            // ""，直接拼接
var ok bool                // false，直接判断
var items []string         // nil，直接 append
var sb strings.Builder     // 空，直接 WriteString
var mu sync.Mutex          // 未锁定，直接 Lock

// ⚠️ 只有这三个需要 make 初始化
m := make(map[string]int)  // map 必须初始化才能写入
ch := make(chan int)       // channel 必须初始化
// 指针需要指向有效地址
```

这就是 Go 零值设计的**一致性**：
- **值类型**（int、string、bool、struct）→ 零值总是可用
- **引用类型**（slice）→ 零值可用（append 会自动分配）
- **需要底层资源的类型**（map、channel）→ 必须 make

**📌 nil slice 安全操作指南**

```go
var s []int  // nil slice

// ✅ 这些都安全
len(s)                    // 0
cap(s)                    // 0
s = append(s, 1)          // 正常工作
for _, v := range s { }   // 空的不执行，不会 panic

// ❌ 下标访问会 panic（因为长度是 0）
s[0]  // panic: index out of range
```

**一句话：nil slice 当空集合用，append 和 range 随便用，下标访问先查长度。**

```go
// 常见用法
var users []User
users = append(users, user1, user2)  // 直接 append

for _, u := range users {  // 安全遍历
    fmt.Println(u.Name)
}

if len(users) > 0 {  // 取值前检查长度
    first := users[0]
}
```

---

#### 2.2 如何判断一个类型的零值是否可用？

为什么 `strings.Builder` 零值可用，但 `map` 不行？看内部结构：

```go
// strings.Builder 内部
type Builder struct {
    buf []byte  // slice 零值是 nil，但 append(nil, ...) 可以工作
}

// 所以 Builder 零值可用——它的字段零值都能正常工作
```

**📌 判断规则**

| 判断方法 | 说明 |
|---------|------|
| **看文档** | 好的库会明确说明，如 `sync.Mutex`: "A Mutex must not be copied after first use. The zero value for a Mutex is an unlocked mutex." |
| **看是否有 `NewXxx()` 构造函数** | 如果库提供了构造函数，通常意味着零值不够用 |
| **看内部是否有 map/chan 字段** | 如果方法需要写入这些字段，零值会 panic |

**📌 标准库中"零值可用"的典型类型**

```go
var mu sync.Mutex          // ✅ 文档明确说明零值是未锁定的锁
var wg sync.WaitGroup      // ✅ 零值可用
var buf bytes.Buffer       // ✅ 零值是空缓冲区
var sb strings.Builder     // ✅ 零值是空构建器
var client http.Client     // ✅ 零值使用默认配置
var once sync.Once         // ✅ 零值可用
```

> **洞察**：Go 标准库的设计原则是尽量让零值可用。当你设计自己的结构体时，也应该遵循这个原则——让用户能 `var x MyType` 直接使用。

---

#### 2.3 利用零值简化代码（实战案例）

**案例 1：计数器无需初始化**

```go
// Go: 零值直接可用
counts := make(map[string]int)
for _, word := range words {
    counts[word]++  // 零值 0，直接 ++，无需判断 key 是否存在
}
```

```java
// Java: 需要 getOrDefault
Map<String, Integer> counts = new HashMap<>();
for (String word : words) {
    counts.put(word, counts.getOrDefault(word, 0) + 1);
}
```

**案例 2：布尔零值做默认配置**

```go
type Config struct {
    Debug    bool  // 零值 false = 默认不开启调试 ✅ 合理
    MaxRetry int   // 零值 0 = 不重试？需要考虑是否合理
}

// 零值配置可以直接使用
var cfg Config
if cfg.Debug {
    log.Println("debug mode")
}
```

**案例 3：strings.Builder 零值可用**

```go
func buildSQL(conditions []string) string {
    var sb strings.Builder  // 📌 无需 new，零值即可用
    sb.WriteString("SELECT * FROM users WHERE 1=1")
    for _, cond := range conditions {
        sb.WriteString(" AND ")
        sb.WriteString(cond)
    }
    return sb.String()
}
```

---

#### 2.4 零值陷阱详解（需要注意的边界情况）

前面说了"只有 map、channel、指针需要初始化"，这里详细解释为什么：

**陷阱 1：nil map 写入 panic**

```go
var m map[string]int  // nil map
fmt.Println(m["key"]) // ✅ 读取返回零值 0，不会 panic
m["key"] = 1          // 💥 panic: assignment to entry in nil map

// 正确做法
m = make(map[string]int)
m["key"] = 1          // ✅ 现在可以写入
```

> 这和 Java 的 NPE 类似，但更隐蔽——**读取不报错，写入才 panic**。

**陷阱 2：nil slice vs 空 slice 的 JSON 序列化**

```go
var nilSlice []int        // nil
emptySlice := []int{}     // 空但非 nil

// 功能上几乎等价
nilSlice = append(nilSlice, 1)    // ✅ 都可以
emptySlice = append(emptySlice, 1) // ✅

// 但 JSON 序列化结果不同！
json.Marshal(nilSlice)    // null  ← API 返回这个可能有问题
json.Marshal(emptySlice)  // []    ← 前端通常期望这个
```

| 场景 | 推荐 | 原因 |
|------|------|------|
| 函数内部使用 | `var s []T` | nil slice 可以直接 append |
| JSON API 返回 | `make([]T, 0)` 或 `[]T{}` | 避免返回 `null` |

**陷阱 3：time.Time 零值是 0001-01-01**

```go
var t time.Time
fmt.Println(t.Format(time.DateTime))  // 0001-01-01 00:00:00 😱

// 判断是否为零值
if t.IsZero() {
    fmt.Println("时间未设置")
}
```

🆚 Java 用 `null` 表示未设置，Go 用 `IsZero()` 或指针 `*time.Time`。

---

#### 2.5 时间处理（Go 1.20+ 最佳实践）

既然提到了 `time.Time` 零值，这里完整介绍时间处理：

```go
import "time"

now := time.Now()

// ✅ Go 1.20+ 推荐：使用预定义常量
fmt.Println(now.Format(time.DateTime))  // 2024-01-15 14:30:00
fmt.Println(now.Format(time.DateOnly))  // 2024-01-15
fmt.Println(now.Format(time.TimeOnly))  // 14:30:00

// ⚠️ 旧写法（仍然有效，但不推荐）
fmt.Println(now.Format("2006-01-02 15:04:05"))
```

**Go 1.20+ 预定义时间常量**：

| 常量 | 值 | 说明 |
|------|-----|------|
| `time.DateTime` | `"2006-01-02 15:04:05"` | 日期时间 |
| `time.DateOnly` | `"2006-01-02"` | 仅日期 |
| `time.TimeOnly` | `"15:04:05"` | 仅时间 |
| `time.RFC3339` | `"2006-01-02T15:04:05Z07:00"` | 标准格式（API 常用） |

> **📌 为什么是 2006-01-02 15:04:05？**
> 这是 Go 的独特设计：`01/02 03:04:05PM '06 -0700`（美式日期顺序 1-2-3-4-5-6-7），便于记忆。

**🆚 与 Java 对比**

```java
// Java: 使用 DateTimeFormatter
LocalDateTime now = LocalDateTime.now();
String formatted = now.format(DateTimeFormatter.ofPattern("yyyy-MM-dd HH:mm:ss"));
```

```go
// Go: 使用参考时间作为模板
now := time.Now()
formatted := now.Format(time.DateTime)  // 或 "2006-01-02 15:04:05"
```

> **洞察**：Java 用 `yyyy-MM-dd` 这种符号模式，Go 用具体的参考时间 `2006-01-02`。Go 的方式更直观——格式串本身就是输出示例。

---

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

本章包含 4 个作业，分别对应不同的知识点。建议按顺序完成。

| 作业 | 目录 | 知识点 | 难度 |
|-----|------|-------|------|
| 1. 词频统计 | `homework/01-word-count/` | 零值、map、slice、range | ⭐ |
| 2. 用户管理 | `homework/02-user-manager/` | 结构体、方法、接口 | ⭐⭐ |
| 3. 安全计算器 | `homework/03-calculator/` | 多返回值、错误处理 | ⭐⭐ |
| 4. 计数器 | `homework/04-counter/` | 指针、值传递 | ⭐⭐ |

---

### 作业 1：词频统计

**目标**：理解 map 的零值陷阱，利用 int 零值简化计数逻辑

```bash
cd homework/01-word-count && go run main.go
```

**要求**：
1. 实现 `CountWords(words []string) map[string]int` 统计词频
2. 实现 `TopWords(counts map[string]int, n int) []string` 返回 Top N

**提示**：
- `counts[word]++` 利用 int 零值直接自增，无需判断 key 是否存在

---

### 作业 2：用户管理

**目标**：掌握结构体、方法定义、接口的隐式实现

```bash
cd homework/02-user-manager && go run main.go
```

**要求**：
1. 定义 `User` 结构体和 `NewUser` 构造函数
2. 实现 `String()` 方法满足 `fmt.Stringer` 接口
3. 实现 `Activate()` 和 `Deactivate()` 方法

**思考**：为什么 `String()` 用值接收者，`Activate()` 用指针接收者？

---

### 作业 3：安全计算器

**目标**：掌握 Go 的多返回值和错误处理模式

```bash
cd homework/03-calculator && go run main.go
```

**要求**：
1. 定义哨兵错误 `ErrDivideByZero`、`ErrNegativeSqrt`
2. 实现 `SafeDivide`、`SafeSqrt`、`Calculate` 函数
3. 使用 `errors.Is` 检查具体错误类型

**对比 Java**：
```java
// Java: 抛异常
if (b == 0) throw new ArithmeticException("除数不能为零");
```
```go
// Go: 返回错误
if b == 0 { return 0, ErrDivideByZero }
```

---

### 作业 4：计数器

**目标**：理解 Go 的值传递机制和指针的必要性

```bash
cd homework/04-counter && go run main.go
```

**要求**：
1. 实现 `IncrementWrong()` —— 值接收者，故意写"错"
2. 实现 `IncrementRight()` —— 指针接收者，正确实现
3. 观察并理解两者的区别

**这是 Java 程序员最容易困惑的点**：
```java
// Java: c 是引用，方法总能修改原对象
Counter c = new Counter();
c.increment();  // 总是工作
```
```go
// Go: c 是值，值接收者收到的是拷贝
c := Counter{}
c.IncrementWrong()  // 不工作！修改的是拷贝
c.IncrementRight()  // 需要指针接收者
```

---

## 参考资料
- [Effective Go](https://go.dev/doc/effective_go)
- [Go 语言圣经](https://gopl.io/)
- [Go 错误处理最佳实践](https://go.dev/blog/go1.13-errors)
