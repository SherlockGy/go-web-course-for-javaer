# 01 - Hello World

## 学习目标

掌握 Go 语言开发环境搭建和第一个程序的编写。

---

## 学习要点

### 1. Go 环境安装与配置
- 下载安装 Go：https://go.dev/dl/
- 环境变量说明：
  - `GOROOT`：Go 安装目录
  - `GOPATH`：工作区目录（Go 1.11 后可选）
  - `GOPROXY`：模块代理（国内推荐 `https://goproxy.cn,direct`）

### 2. 第一个 Go 程序结构
```go
package main          // 声明包名，main 包是可执行程序入口

import "fmt"          // 导入标准库

func main() {         // main 函数是程序入口
    fmt.Println("Hello, World!")
}
```

### 3. go mod 模块化管理
```bash
go mod init <module-name>   # 初始化模块
go mod tidy                 # 整理依赖
go mod download             # 下载依赖
```

### 4. 常用命令
| 命令 | 说明 |
|------|------|
| `go run main.go` | 编译并运行 |
| `go build` | 编译生成可执行文件 |
| `go fmt ./...` | 格式化代码 |
| `go vet ./...` | 静态检查 |

---

## 示例代码

### examples/01-first-program/
最简单的 Hello World 程序

### examples/02-project-structure/
标准项目结构示例

---

## 作业任务

### 任务描述
创建一个 Go 模块，编写程序输出当前时间。

### 要求
1. 使用 `go mod init` 初始化项目
2. 使用 `time` 包获取当前时间
3. 格式化输出：`当前时间: 2024-01-01 12:00:00`

### 验收标准
- `go run main.go` 能正确运行
- 输出格式正确

### 提示
```go
import "time"

time.Now()                              // 获取当前时间
time.Now().Format("2006-01-02 15:04:05") // 格式化时间
```
> 注意：Go 的时间格式化使用固定的参考时间 `2006-01-02 15:04:05`

---

## 参考资料
- [Go 官方文档](https://go.dev/doc/)
- [Go by Example](https://gobyexample.com/)
