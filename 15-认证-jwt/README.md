# 15 - 认证 JWT

## 学习目标

掌握 JWT 认证机制，实现用户登录和接口鉴权。

---

## 🆚 Java 对比：认证实现

| 特性 | Spring Security | Go + JWT |
|------|-----------------|----------|
| 配置方式 | `SecurityConfig` 类 | 中间件函数 |
| 过滤链 | `FilterChain` | `gin.HandlerFunc` 链 |
| 注解鉴权 | `@PreAuthorize` | 手动检查或中间件 |
| JWT 支持 | 需要额外配置 | `golang-jwt` 包 |

> **洞察**：Spring Security 功能全面但学习曲线陡峭，Go 的 JWT 实现更"原始"但更透明——你清楚地知道每一步在做什么。

---

## 学习要点

### 1. JWT 结构

```
Header.Payload.Signature
eyJhbGciOiJIUzI1NiJ9.eyJ1c2VyX2lkIjoxfQ.signature
```

- **Header**：算法类型（HS256）
- **Payload**：用户数据（Claims）
- **Signature**：签名（防篡改）

### 2. 安装 JWT 库

```bash
go get -u github.com/golang-jwt/jwt/v5
```

### 3. 生成 Token

```go
import "github.com/golang-jwt/jwt/v5"

type Claims struct {
    UserID   int64  `json:"user_id"`
    Username string `json:"username"`
    jwt.RegisteredClaims
}

func GenerateToken(userID int64, username string) (string, error) {
    claims := Claims{
        UserID:   userID,
        Username: username,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte("your-secret-key"))
}
```

> **🆚 Java 对比**
> ```java
> String token = Jwts.builder()
>     .setSubject(userId.toString())
>     .setIssuedAt(new Date())
>     .setExpiration(new Date(System.currentTimeMillis() + 86400000))
>     .signWith(key, SignatureAlgorithm.HS256)
>     .compact();
> ```
> 写法类似，都是 Builder 模式。但 Go 用结构体定义 Claims，Java 用 Map。

### 4. 验证 Token

```go
func ValidateToken(tokenString string) (*Claims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
        return []byte("your-secret-key"), nil
    })

    if err != nil {
        return nil, err
    }

    if claims, ok := token.Claims.(*Claims); ok && token.Valid {
        return claims, nil
    }

    return nil, errors.New("无效的 token")
}
```

### 5. 认证中间件

```go
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 获取 Authorization header
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(401, gin.H{"error": "未提供认证令牌"})
            c.Abort()
            return
        }

        // 解析 Bearer token
        parts := strings.SplitN(authHeader, " ", 2)
        if len(parts) != 2 || parts[0] != "Bearer" {
            c.JSON(401, gin.H{"error": "认证格式错误"})
            c.Abort()
            return
        }

        // 验证 token
        claims, err := ValidateToken(parts[1])
        if err != nil {
            c.JSON(401, gin.H{"error": "无效的令牌"})
            c.Abort()
            return
        }

        // 将用户信息存入上下文
        c.Set("userID", claims.UserID)
        c.Set("username", claims.Username)
        c.Next()
    }
}
```

> **🆚 Spring Security 对比**
> ```java
> // Spring 用 Filter 链，配置在 SecurityConfig
> http.addFilterBefore(jwtFilter, UsernamePasswordAuthenticationFilter.class);
> ```
> Spring 的 Filter 链是"声明式"配置，Go 的中间件是"命令式"代码。Go 更直观，Spring 更"框架化"。

### 6. 使用中间件

```go
r := gin.Default()

// 公开路由
r.POST("/login", loginHandler)
r.POST("/register", registerHandler)

// 需要认证的路由
auth := r.Group("/api")
auth.Use(AuthMiddleware())
{
    auth.GET("/me", getCurrentUser)
    auth.GET("/users", listUsers)
}
```

### 7. 获取当前用户

```go
func getCurrentUser(c *gin.Context) {
    userID := c.GetInt64("userID")
    username := c.GetString("username")

    c.JSON(200, gin.H{
        "user_id":  userID,
        "username": username,
    })
}
```

---

## 示例代码

### examples/01-jwt-basics/
JWT 生成与验证

### examples/02-gin-auth/
Gin 认证中间件

### examples/03-protected-routes/
受保护路由示例

---

## 作业任务

### 任务描述
实现完整的 JWT 认证流程。

### 要求
1. `POST /login` - 登录接口，返回 JWT
2. `GET /api/me` - 获取当前用户（需认证）
3. 认证中间件处理 Bearer Token

### 验收标准
```bash
# 登录获取 token
TOKEN=$(curl -s -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"username":"tom","password":"123456"}' | jq -r '.token')

# 访问受保护接口
curl http://localhost:8080/api/me \
  -H "Authorization: Bearer $TOKEN"
# 返回: {"user_id":1,"username":"tom"}

# 无 token 访问
curl http://localhost:8080/api/me
# 返回: {"error":"未提供认证令牌"}
```

---

## 参考资料
- [golang-jwt](https://github.com/golang-jwt/jwt)
- [JWT.io](https://jwt.io/) - JWT 调试工具
