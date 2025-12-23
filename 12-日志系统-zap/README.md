# 12 - 日志系统 Zap

## 学习目标

掌握 Zap 结构化日志，替代 `fmt.Println` 进行生产级日志记录。

---

## 🆚 Java 对比：日志设计哲学

| 特性 | SLF4J + Logback | Go Zap |
|------|-----------------|--------|
| API 风格 | `log.info("User {} logged in", name)` | `log.Info("User logged in", zap.String("name", name))` |
| 配置方式 | `logback.xml` | 代码配置 |
| 性能 | 中等 | 极高（零分配） |
| 结构化 | MDC（侵入式） | 原生支持 |

> **洞察**：Java 日志习惯用字符串模板 `{}`，Go 用强类型字段。Go 的方式对日志分析更友好（JSON 结构化），但写起来稍繁琐。

---

## 学习要点

### 1. Zap 基础

```bash
go get go.uber.org/zap
```

```go
// 快速开发用 Sugar（像 printf）
logger, _ := zap.NewDevelopment()
sugar := logger.Sugar()
sugar.Infof("User %s logged in", "tom")

// 高性能用 Logger（强类型）
logger.Info("User logged in",
    zap.String("username", "tom"),
    zap.Int("userID", 123),
)
```

> **🆚 Java 对比**
> ```java
> // SLF4J
> log.info("User {} logged in, ID: {}", username, userId);
>
> // Go Zap
> logger.Info("User logged in",
>     zap.String("username", username),
>     zap.Int("userID", userId))
> ```
> Zap 更啰嗦，但类型安全，且生成的 JSON 日志更易于 ELK 等工具分析。

### 2. 日志级别

```go
logger.Debug("调试信息")  // 开发时用
logger.Info("一般信息")   // 正常运行
logger.Warn("警告信息")   // 需要注意
logger.Error("错误信息")  // 出错但能继续
logger.Fatal("致命错误")  // 出错后退出程序
```

### 3. 结构化日志字段

```go
logger.Info("请求完成",
    zap.String("method", "GET"),
    zap.String("path", "/api/users"),
    zap.Int("status", 200),
    zap.Duration("latency", time.Millisecond*50),
    zap.Error(err),  // 错误字段
)
```

**输出 JSON**：
```json
{
  "level": "info",
  "msg": "请求完成",
  "method": "GET",
  "path": "/api/users",
  "status": 200,
  "latency": "50ms"
}
```

### 4. 日志输出配置

```go
// 同时输出到控制台和文件
core := zapcore.NewTee(
    // 控制台：人类可读
    zapcore.NewCore(
        zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig()),
        zapcore.AddSync(os.Stdout),
        zapcore.DebugLevel,
    ),
    // 文件：JSON 格式
    zapcore.NewCore(
        zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
        zapcore.AddSync(logFile),
        zapcore.InfoLevel,
    ),
)
logger := zap.New(core)
```

### 5. 日志轮转（lumberjack）

```go
import "gopkg.in/natefinch/lumberjack.v2"

writer := &lumberjack.Logger{
    Filename:   "./logs/app.log",
    MaxSize:    10,    // MB
    MaxBackups: 5,     // 保留旧文件数
    MaxAge:     30,    // 保留天数
    Compress:   true,  // 压缩旧文件
}
```

> **🆚 Logback 对比**
> ```xml
> <appender class="ch.qos.logback.core.rolling.RollingFileAppender">
>     <rollingPolicy class="TimeBasedRollingPolicy">
>         <maxHistory>30</maxHistory>
>     </rollingPolicy>
> </appender>
> ```
> Java 用 XML 配置，Go 用代码配置。各有优劣：XML 不用重新编译，代码更灵活。

### 6. 与 Gin 集成

```go
func GinLogger(logger *zap.Logger) gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        c.Next()

        logger.Info("HTTP 请求",
            zap.String("method", c.Request.Method),
            zap.String("path", c.Request.URL.Path),
            zap.Int("status", c.Writer.Status()),
            zap.Duration("latency", time.Since(start)),
        )
    }
}
```

---

## 示例代码

### examples/01-basic-zap/
Zap 基础使用

### examples/02-structured-log/
结构化日志示例

### examples/03-file-output/
输出到文件 + 轮转

### examples/04-gin-integration/
与 Gin 集成

---

## 作业任务

### 任务描述
创建生产级日志系统，同时输出到控制台和文件。

### 要求
1. 控制台：彩色人类可读格式
2. 文件：JSON 格式，便于 ELK 分析
3. 支持日志轮转（10MB/文件，保留 5 个）
4. 封装为 Gin 中间件

### 验收标准
- 控制台输出彩色日志
- `logs/app.log` 有 JSON 日志
- 日志包含：时间、级别、方法、路径、状态码、耗时

---

## 参考资料
- [Zap GitHub](https://github.com/uber-go/zap)
- [lumberjack](https://github.com/natefinch/lumberjack)
