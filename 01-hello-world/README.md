# 01 - Hello World

## 学习目标

搭建 Go 开发环境，运行第一个程序，理解 Go 与 Java 的核心差异。

---

## 学习要点

### 1. Go 环境安装

**下载安装**：https://go.dev/dl/

**验证安装**：
```bash
go version    # 输出: go version go1.22.x ...
```

**环境变量说明**：

| 变量 | 说明 | 是否必须                                                                     |
|------|------|--------------------------------------------------------------------------|
| `GOROOT` | Go 安装目录 | 自动设置，无需配置                                                                |
| `GOPATH` | ~~工作区目录~~ | **已过时，现代 Go 不需要**                                                        |
| `GOPROXY` | 模块代理 | 默认值：`https://proxy.golang.org,direct`<br/>国内必须：`https://goproxy.cn,direct` |

> **📌 重要**：Go 1.16+ 默认启用 Module 模式，**完全不需要配置 GOPATH**。你可以在任意目录创建项目。这与 Java 必须配置 `JAVA_HOME` 不同。

**国内代理配置**（必须）：
```bash
go env -w GOPROXY=https://goproxy.cn,direct
```

---

### 2. 第一个程序：Go vs Java 对比

<table>
<tr><th>Go</th><th>Java</th></tr>
<tr>
<td>

```go
package main

import "fmt"

func main() {
    fmt.Println("Hello, Go!")
}
```

</td>
<td>

```java
package com.example;

public class Main {
    public static void main(String[] args) {
        System.out.println("Hello, Java!");
    }
}
```

</td>
</tr>
<tr>
<td>

- `package main` 是可执行程序入口
- 直接 `func main()` 无需类包装
- 运行：`go run main.go`

</td>
<td>

- 需要 `public class` 包装
- `main` 必须是 `public static void`
- 编译运行分离：`javac` → `java`

</td>
</tr>
</table>

**📌 核心差异**：
- Go 没有类，`func` 直接定义在包级别
- Go 编译快：无注解处理、依赖模型简单（Java 项目构建慢主要是 Maven/Gradle 流程复杂，而非 javac 本身）
- Go 生成单一可执行文件，无需 JVM

---

### 3. go mod 模块化管理

#### 初始化项目
```bash
mkdir myproject && cd myproject
go mod init myproject          # 本地项目
# 或
go mod init github.com/yourname/myproject  # 开源项目（推荐）
```

#### 核心命令区别

| 命令 | 作用 | 使用场景 |
|------|------|----------|
| `go mod init <name>` | 初始化模块，创建 go.mod | 新项目第一步 |
| `go mod tidy` | **分析代码**，自动添加/删除依赖 | 增删 import 后执行 |
| `go mod download` | 下载 go.mod 中声明的所有依赖 | CI/CD 预下载缓存 |

**📌 最佳实践**：
- 日常开发用 `go mod tidy`，它会分析你的代码自动管理依赖
- `go mod download` 主要用于 CI 场景（预热缓存）
- **go.sum 文件必须提交到 Git**（记录依赖的校验和，确保安全）

```bash
# 典型工作流
go mod init myproject     # 1. 初始化
# ... 编写代码，添加 import ...
go mod tidy               # 2. 自动下载并整理依赖
```

---

### 4. 常用命令

#### 基础命令
```bash
go run main.go       # 编译并运行（开发用）
go build             # 编译生成可执行文件
go fmt ./...         # 格式化代码（强制统一风格）
go vet ./...         # 静态检查（发现潜在 bug）
```

#### go build 真实场景

```bash
# 基本编译
go build                        # 输出与目录同名的可执行文件
go build -o myapp.exe           # 指定输出文件名

# 跨平台编译（Go 的杀手级特性）
# 在 Windows 上编译 Linux 可执行文件：
set GOOS=linux
set GOARCH=amd64
go build -o myapp-linux

# 在 Mac/Linux 上编译 Windows 可执行文件：
GOOS=windows GOARCH=amd64 go build -o myapp.exe

# 静态编译（无外部依赖，适合 Docker）
CGO_ENABLED=0 go build -o myapp

# 减小体积（去除调试信息）
go build -ldflags="-s -w" -o myapp
```

> **📌 与 Java 对比**：Java 需要 "Write Once, Run Anywhere" 依赖 JVM；Go 是 "Build Once, Run Anywhere" —— 直接编译成目标平台的原生可执行文件。

---

### 5. GoLand / VS Code 快捷输入（拓展）

高效编写 Go 代码的 Live Templates / Snippets：

**常用快捷键（GoLand & VS Code 通用）**：

| 输入     | 展开结果                         | 说明        |
|--------|------------------------------|-----------|
| `fp`   | `fmt.Println()`              | 打印输出      |
| `main` | `func main() { }`            | main 函数骨架 |
| `fori` | `for i := 0; i < 10; i++ {}` | for 循环    |
| `forr` | `for _, := range { }`        | range 循环  |

> **📌 使用方法**：输入缩写后按 `Tab` 键展开

**🆚 与 Java 习惯对比**：

| 场景 | Java (IntelliJ) | Go (GoLand/VS Code) |
|-----|------|-----|
| 打印输出 | `sout` → `System.out.println()` | `fp` → `fmt.Println()` |

---

## 示例代码

### examples/01-first-program/
最简单的 Hello World 程序，展示 Go 程序的基本结构。

---

## 作业任务

### 任务描述
完成 `homework/main.go`，使用 `fmt` 包输出个人学习计划。

### 要求
1. 使用 `fmt.Println` 输出至少 3 行信息
2. 使用 `fmt.Printf` 进行格式化输出（练习 `%s`、`%d` 占位符）
3. 内容包含：姓名、学习目标、预计学习天数

### 预期输出示例
```
=== Go 学习计划 ===
姓名: 张三
学习目标: 掌握 Go Web 开发
预计学习天数: 30 天
=== 开始学习！===
```

### 验收标准
- `cd homework && go run main.go` 能正确运行
- 使用 `fmt.Println` 和 `fmt.Printf` 两种方式

### 提示
```go
fmt.Println("Hello")           // 直接输出，自动换行
fmt.Printf("姓名: %s\n", name) // 格式化输出，%s 是字符串占位符
fmt.Printf("天数: %d\n", days) // %d 是整数占位符
```

> 📌 IDE 快捷键：`fp` + Tab 生成 `fmt.Println()`（GoLand/VS Code 通用）

---

## 参考资料
- [Go 官方文档](https://go.dev/doc/)
- [Go by Example](https://gobyexample.com/)
- [Go 1.20 Release Notes - time 常量](https://go.dev/doc/go1.20#time)
