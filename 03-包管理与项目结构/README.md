# 03 - 包管理与项目结构

## 学习目标

掌握 Go 模块管理和标准项目布局，理解包拆分最佳实践。

---

## 学习要点

### 1. go mod 深入

```bash
# 初始化模块
go mod init github.com/yourname/project

# 整理依赖（添加缺失、移除无用）
go mod tidy

# 下载依赖到本地缓存
go mod download

# 将依赖复制到 vendor 目录
go mod vendor

# 查看某个依赖为何被引入
go mod why github.com/gin-gonic/gin
```

### 2. 包的导入与导出

```go
// 大写开头 = 导出（public）
func PublicFunc() {}
type PublicType struct{}

// 小写开头 = 私有（package private）
func privateFunc() {}
type privateType struct{}
```

### 3. 标准项目布局

```
project/
├── cmd/                      # 可执行文件入口
│   ├── server/
│   │   └── main.go          # Web 服务器入口
│   └── cli/
│       └── main.go          # CLI 工具入口
│
├── internal/                 # 私有包（外部无法导入）
│   ├── handler/             # HTTP 处理器
│   ├── service/             # 业务逻辑
│   ├── repository/          # 数据访问
│   ├── model/               # 数据模型
│   └── config/              # 配置管理
│
├── pkg/                      # 公共包（可被外部导入）
│   ├── logger/              # 日志工具
│   ├── jwt/                 # JWT 工具
│   └── response/            # 响应工具
│
├── api/                      # API 定义
│   ├── proto/               # gRPC proto 文件
│   └── openapi/             # OpenAPI/Swagger 定义
│
├── configs/                  # 配置文件模板
│   └── config.example.yaml
│
├── scripts/                  # 脚本文件
│   └── build.sh
│
├── go.mod                    # 模块定义
├── go.sum                    # 依赖校验
└── README.md
```

### 4. internal 的特殊作用

```
myproject/
├── internal/
│   └── secret/
│       └── secret.go       # 只有 myproject 内部可以导入
└── pkg/
    └── public/
        └── public.go       # 任何项目都可以导入
```

**规则**：`internal` 目录下的包只能被其父目录的代码导入。

### 5. 包拆分原则

| 原则 | 说明 | 反例 |
|------|------|------|
| **单一职责** | 一个包只做一件事 | `utils` 包塞满杂七杂八的函数 |
| **按功能拆分** | 按业务功能组织 | `models/`, `controllers/` 按类型分 |
| **避免循环依赖** | A→B→A 是不允许的 | handler 导入 service，service 又导入 handler |
| **internal 隔离** | 业务代码放 internal | 把所有代码都放 pkg |
| **包名简洁** | 小写、单词 | `userService`, `user_service` |

### 6. 循环依赖解决方案

```
# 问题：handler → service → handler

# 方案1：提取公共接口到独立包
internal/
├── handler/     → 依赖 service
├── service/     → 依赖 repository, types
└── types/       ← 公共类型（handler 和 service 都可以用）

# 方案2：依赖注入
service 定义接口，handler 实现接口
```

### 7. go.mod 与 go.sum

```go
// go.mod
module github.com/yourname/project

go 1.21

require (
    github.com/gin-gonic/gin v1.9.1
    gorm.io/gorm v1.25.5
)

// go.sum（自动生成，用于校验依赖完整性）
github.com/gin-gonic/gin v1.9.1 h1:...
github.com/gin-gonic/gin v1.9.1/go.mod h1:...
```

---

## 示例代码

### examples/01-package-basics/
包的导入导出演示

### examples/02-project-layout/
标准项目布局示例

### examples/03-internal-pkg/
internal 与 pkg 的区别演示

### examples/04-dependency-direction/
依赖方向示例（避免循环依赖）

---

## 作业任务

### 任务描述
创建一个符合标准布局的项目结构。

### 要求
1. 使用 `go mod init` 初始化项目
2. 创建以下目录结构：
   ```
   myproject/
   ├── cmd/server/main.go
   ├── internal/
   │   ├── handler/user.go
   │   └── service/user.go
   ├── pkg/response/response.go
   └── go.mod
   ```
3. `handler` 调用 `service`
4. `handler` 和 `service` 都使用 `pkg/response`
5. 确保代码能编译通过

### 验收标准
- 运行 `go build ./...` 无报错
- 理解为什么 `internal` 包不能被外部项目导入

### 思考题
如果 `service` 包需要调用 `handler` 包的某个函数，应该怎么做？（提示：接口或提取公共包）

---

## 参考资料
- [Standard Go Project Layout](https://github.com/golang-standards/project-layout)
- [Go Modules Reference](https://go.dev/ref/mod)
