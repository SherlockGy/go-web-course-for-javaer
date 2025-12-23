# 11 - 配置管理 Viper

## 学习目标

掌握 Viper 配置管理，实现灵活的配置加载。

---

## 🆚 Java 对比：配置管理哲学

| 特性 | Spring Boot | Go Viper |
|------|-------------|----------|
| 配置文件 | `application.yml` | `config.yaml` |
| 环境覆盖 | `application-{profile}.yml` | 环境变量/命令行参数 |
| 配置注入 | `@Value("${key}")` | `viper.GetString("key")` |
| 热加载 | 需要 actuator | `viper.WatchConfig()` |

> **洞察**：Spring Boot 的配置是"声明式"的（注解注入），Go 是"命令式"的（主动读取）。Go 更显式，但也更灵活。

---

## 学习要点

### 1. Viper 基础

```bash
go get github.com/spf13/viper
```

```go
import "github.com/spf13/viper"

func main() {
    viper.SetConfigName("config")  // 配置文件名（无扩展名）
    viper.SetConfigType("yaml")    // 配置类型
    viper.AddConfigPath(".")       // 查找路径

    if err := viper.ReadInConfig(); err != nil {
        panic("读取配置失败: " + err.Error())
    }

    port := viper.GetInt("server.port")
    fmt.Println("端口:", port)
}
```

### 2. 配置结构体映射

```yaml
# config.yaml
server:
  host: "0.0.0.0"
  port: 8080

database:
  driver: "sqlite"
  dsn: "./data.db"
```

```go
type Config struct {
    Server   ServerConfig   `mapstructure:"server"`
    Database DatabaseConfig `mapstructure:"database"`
}

type ServerConfig struct {
    Host string `mapstructure:"host"`
    Port int    `mapstructure:"port"`
}

// 加载到结构体
var cfg Config
viper.Unmarshal(&cfg)
```

> **🆚 Spring Boot 对比**
> ```java
> @ConfigurationProperties(prefix = "server")
> public class ServerConfig {
>     private String host;
>     private int port;
> }
> ```
> Spring 用注解 + 自动装配，Go 用 `mapstructure` tag + 显式调用。

### 3. 环境变量覆盖

```go
// 自动绑定环境变量
viper.AutomaticEnv()
viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

// 配置: server.port
// 环境变量: SERVER_PORT=9090
```

```bash
SERVER_PORT=9090 go run main.go
```

> **🆚 Spring Boot 对比**
> ```bash
> # Spring Boot 也支持，但命名规则不同
> SERVER_PORT=9090 java -jar app.jar
> ```
> 思路相同，都是环境变量覆盖配置文件。

### 4. 配置热加载

```go
viper.WatchConfig()
viper.OnConfigChange(func(e fsnotify.Event) {
    fmt.Println("配置已更新:", e.Name)
    // 重新加载到结构体
    viper.Unmarshal(&cfg)
})
```

> **洞察**：Go 的热加载需要手动处理"配置变更后怎么办"，而 Spring Cloud Config 可以自动刷新 Bean。显式 vs 隐式的又一体现。

---

## 示例代码

### examples/01-basic-config/
基础配置读取

### examples/02-struct-mapping/
配置映射到结构体

### examples/03-env-override/
环境变量覆盖配置

---

## 作业任务

### 任务描述
创建配置系统，支持从文件和环境变量加载配置。

### config.yaml
```yaml
server:
  host: "127.0.0.1"
  port: 8080

database:
  dsn: "data.db"

jwt:
  secret: "dev-secret"
  expiration: 86400
```

### 要求
1. 定义对应的配置结构体
2. 实现 `LoadConfig() (*Config, error)` 函数
3. 支持环境变量覆盖（如 `SERVER_PORT`）

### 验收标准
```bash
# 默认配置
go run main.go
# 输出: Server: 127.0.0.1:8080

# 环境变量覆盖
SERVER_PORT=9090 go run main.go
# 输出: Server: 127.0.0.1:9090
```

---

## 参考资料
- [Viper GitHub](https://github.com/spf13/viper)
