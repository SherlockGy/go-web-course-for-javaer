# internal vs pkg 的区别

## 目录结构
```
03-internal-pkg/
├── go.mod
├── main.go
├── internal/
│   └── secret/           # 私有包 - 只能本模块使用
│       └── secret.go
└── pkg/
    └── public/           # 公共包 - 可被外部导入
        └── public.go
```

## 关键区别

| 特性 | internal/ | pkg/ |
|------|-----------|------|
| 可被外部导入 | ❌ 不可以 | ✅ 可以 |
| 编译器保护 | ✅ 强制限制 | ❌ 无限制 |
| 适合内容 | 业务逻辑、内部实现 | 通用工具、可复用代码 |

## 演示

如果外部模块尝试导入 `internal/secret`，会得到编译错误：
```
use of internal package not allowed
```

## 运行
```bash
go run main.go
```
