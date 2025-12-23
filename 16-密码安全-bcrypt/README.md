# 16 - 密码安全 bcrypt

## 学习目标

掌握密码安全存储，使用 bcrypt 哈希算法。

---

## 🆚 Java 对比：密码处理

| 方面 | Spring Security | Go bcrypt |
|------|-----------------|-----------|
| API | `BCryptPasswordEncoder` | `bcrypt.GenerateFromPassword` |
| 默认强度 | cost=10 | `bcrypt.DefaultCost`=10 |
| 使用方式 | `@Bean` 注入 | 直接调用函数 |

> **洞察**：功能完全相同，但 Spring 把它包装成 Bean 注入，Go 直接暴露函数。Go 的方式更"原始"，但你更清楚发生了什么。

---

## 学习要点

### 1. 为什么不能明文存储

```
❌ 数据库泄露 → 所有密码暴露
❌ 用户在多个网站用同一密码 → 连锁反应
❌ 内部人员可以看到用户密码
```

### 2. 为什么不用 MD5/SHA

```
MD5("password") = "5f4dcc3b5aa765d61d8327deb882cf99"
```

- **彩虹表攻击**：预计算常见密码的哈希值
- **速度太快**：GPU 每秒可算数十亿次
- **无盐**：相同密码产生相同哈希

### 3. bcrypt 原理

```
$2a$10$N9qo8uLOickgx2ZMRZoMye.IjJ8.k0sE5z3T5n1P.CZH2h3KnZ/Cu
 ^   ^                                                      ^
算法 cost(2^10次迭代)         盐值(22字符) + 哈希(31字符)
```

- **自带盐**：每次生成不同
- **慢**：故意设计成慢，防暴力破解
- **可调节**：cost 越高越慢越安全

### 4. 密码哈希

```go
import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword(
        []byte(password),
        bcrypt.DefaultCost,  // cost=10
    )
    return string(bytes), err
}
```

> **🆚 Spring Security 对比**
> ```java
> @Bean
> public PasswordEncoder passwordEncoder() {
>     return new BCryptPasswordEncoder();
> }
>
> String hash = passwordEncoder.encode("password");
> ```
> 功能相同，Go 直接调用，Java 通过 Bean 注入。

### 5. 密码验证

```go
func CheckPassword(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}
```

> **🆚 Spring Security 对比**
> ```java
> boolean matches = passwordEncoder.matches(rawPassword, encodedPassword);
> ```

### 6. 注册流程

```go
func Register(username, password, email string) error {
    // 1. 检查用户是否存在
    if userExists(username) {
        return errors.New("用户名已存在")
    }

    // 2. 哈希密码
    hashedPassword, err := HashPassword(password)
    if err != nil {
        return err
    }

    // 3. 保存用户
    user := &User{
        Username: username,
        Password: hashedPassword,  // 存储哈希值
        Email:    email,
    }
    return db.Create(user).Error
}
```

### 7. 登录流程

```go
func Login(username, password string) (string, error) {
    // 1. 查找用户
    var user User
    if err := db.Where("username = ?", username).First(&user).Error; err != nil {
        return "", errors.New("用户不存在")
    }

    // 2. 验证密码
    if !CheckPassword(password, user.Password) {
        return "", errors.New("密码错误")
    }

    // 3. 生成 JWT
    return GenerateToken(user.ID, user.Username)
}
```

### 8. 安全最佳实践

| 实践 | 说明 |
|------|------|
| 使用 bcrypt | 不用 MD5/SHA |
| cost >= 10 | 生产环境建议 12 |
| 密码长度限制 | 6-72 字符（bcrypt 限制） |
| 不记录密码 | 日志中不要打印密码 |
| HTTPS | 传输加密 |
| 定期更换 JWT 密钥 | 降低泄露风险 |

---

## 示例代码

### examples/01-hash-password/
密码哈希示例

### examples/02-verify-password/
密码验证示例

### examples/03-register-login/
完整注册登录流程

---

## 作业任务

### 任务描述
将之前的注册登录接口改为使用 bcrypt。

### 要求
1. 注册时使用 bcrypt 加密密码
2. 登录时使用 bcrypt 验证密码
3. 数据库中不能存储明文密码

### 验收标准
```bash
# 注册
curl -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{"username":"tom","password":"123456","email":"tom@example.com"}'

# 查看数据库
sqlite3 data.db "SELECT password FROM users WHERE username='tom'"
# 输出类似: $2a$10$N9qo8uLOickgx2ZMRZoMye...（不是明文）

# 登录成功
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"username":"tom","password":"123456"}'
# 返回 token

# 密码错误
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"username":"tom","password":"wrong"}'
# 返回错误
```

---

## 参考资料
- [bcrypt 包文档](https://pkg.go.dev/golang.org/x/crypto/bcrypt)
- [OWASP 密码存储指南](https://cheatsheetseries.owasp.org/cheatsheets/Password_Storage_Cheat_Sheet.html)
