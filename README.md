# Go Web 快速上手实战教程 (For Java Developers)

> 🚀 **这是一个由 Claude AI 生成，专为 Java 开发者打造的 Go 语言快速转型指南。**
> 本教程通过对比 Java (Spring Boot) 与 Go 的核心概念，帮助你利用现有的经验快速掌握 Go Web 开发的最佳实践。

## 📚 项目简介

作为一个 Java 开发者，你可能习惯了注解、反射、庞大的依赖注入容器和各种 "Magic"。转到 Go 语言时，最大的挑战往往不是语法，而是**思维模式的转变**。

本项目旨在：
1.  **对比式学习**：建立 `Java Concept` -> `Go Concept` 的映射。
2.  **最佳实践**：直接教授符合 Go 社区规范（Standard Go Project Layout）的写代码方式，避免写出 "Java 味的 Go"。
3.  **实战驱动**：从原生 `net/http` 到 `Gin` 框架，再到 `Gorm`、`Wire`、`JWT` 等全栈生态整合。

---

## 🗺️ Java vs Go 概念极速映射

在开始之前，先建立心理模型：

| Java / Spring 生态 | Go 生态 | 备注 |
| :--- | :--- | :--- |
| **Maven / Gradle** | `go.mod` | 依赖管理，不再需要繁琐的 XML |
| **class** | `struct` | Go 没有类，只有结构体 |
| **interface (implements)** | `interface` (Duck Typing) | 不需要显式 `implements`，实现了方法即实现了接口 |
| **try-catch-finally** | `if err != nil` + `defer` | Go 视错误为普通值，不通过异常流控制逻辑 |
| **Thread** | `Goroutine` | 极轻量级线程，启动只需 `go func()` |
| **Spring Boot** | `Gin` / `Echo` / `Kratos` | Go 倾向于组合库，而不是全家桶框架 |
| **Spring IoC (Autowired)** | `Wire` / 手动注入 | 编译期注入，无反射，代码更清晰 |
| **MyBatis / Hibernate** | `GORM` / `Ent` | ORM 选择类似，但更轻量 |
| **Lombok** | (无) | Go 崇尚显式代码，简单的 getter/setter 通常直接访问字段 |

---

## 📂 课程目录与进度

每个章节都包含 `examples` (示例代码) 和 `homework` (课后作业)。

### 第一阶段：语言基础与工程化
- **[01-hello-world](./01-hello-world)**: 环境搭建，运行第一个 Go 程序。
- **[02-基础语法速览](./02-基础语法速览)**: 变量、函数、结构体、接口、错误处理（对比 Java 差异）。
- **[03-包管理与项目结构](./03-包管理与项目结构)**: 理解 `go.mod`，`internal` 目录的特殊性，Standard Go Layout。

### 第二阶段：Web 开发核心 (HTTP)
- **[04-原生http-入门](./04-原生http-入门)**: 不依赖框架，理解 Go 标准库 `net/http` 的强大。
- **[05-原生http-进阶](./05-原生http-进阶)**: 中间件原理（类似于 Java Filter）、JSON 处理。
- **[06-原生http-实战](./06-原生http-实战)**: 手写一个纯原生 REST API。

### 第三阶段：Gin 框架实战 (The "Spring Boot" way)
- **[07-gin-快速上手](./07-gin-快速上手)**: 路由、参数解析。
- **[08-gin-中间件](./08-gin-中间件)**: 鉴权、日志、CORS。
- **[09-gin-请求响应](./09-gin-请求响应)**: 统一响应封装、结构体绑定与校验（Validator）。
- **[10-gin-错误处理](./10-gin-错误处理)**: 全局异常捕获（Recovery）、自定义错误体系。

### 第四阶段：全栈生态整合
- **[11-配置管理-viper](./11-配置管理-viper)**: 替代 Spring Cloud Config/Properties，支持热加载。
- **[12-日志系统-zap](./12-日志系统-zap)**: 高性能日志，替代 Logback/Log4j2。
- **[13-数据库-gorm基础](./13-数据库-gorm基础)**: CRUD，连接池配置。
- **[14-数据库-gorm进阶](./14-数据库-gorm进阶)**: 事务、钩子、预加载、复杂查询。
- **[15-认证-jwt](./15-认证-jwt)**: 实现无状态登录。
- **[16-密码安全-bcrypt](./16-密码安全-bcrypt)**: 用户密码加密存储最佳实践。

### 第五阶段：架构与工程化 (进阶)
- **[17-分层架构](./17-分层架构)**: `Controller` -> `Service` -> `Repository` (DAO) 的 Go 语言实现。
- **[18-依赖注入](./18-依赖注入)**: 使用 `Google Wire` 实现依赖注入，解耦代码。
- **[19-综合实战](./19-综合实战)**: 综合以上所有知识，构建一个完整的 Web 后端项目。

---

## 🛠️ 环境准备

确保你已经安装了 Go (建议 1.21+)。

```bash
# 1. 检查版本
go version

# 2. 配置国内代理 (加速下载)
go env -w GOPROXY=https://goproxy.cn,direct
```

## 📝 学习建议 (致 Java 开发者)

1.  **忘掉继承**：Go 没有继承，只有组合。不要试图创建复杂的父类子类层级。
2.  **拥抱显式错误**：不要觉得 `if err != nil` 烦。这是 Go 鲁棒性的基石，它强迫你处理每一个可能的错误，而不是把它们抛给顶层。
3.  **少用反射**：Java 中反射满天飞，但在 Go 中反射性能较差且不安全。尽量使用代码生成（如 Wire, EasyJSON）。
4.  **关注数据布局**：Go 是“按值传递”的语言，理解指针 (`*Struct`) 和值 (`Struct`) 的区别至关重要。

## 🤝 贡献与反馈
如果你发现任何问题，或者有更好的对比案例，欢迎提交 PR。
Happy Coding in Go! 🐹
