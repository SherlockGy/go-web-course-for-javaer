# REST API 示例

使用原生 net/http 实现完整的用户 CRUD API。

## 运行

```bash
go run main.go
```

## 测试命令

```bash
# 获取用户列表
curl http://localhost:8080/users

# 获取单个用户
curl http://localhost:8080/users/1

# 创建用户
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{"username":"alice","email":"alice@example.com"}'

# 更新用户
curl -X PUT http://localhost:8080/users/1 \
  -H "Content-Type: application/json" \
  -d '{"username":"tom_updated"}'

# 删除用户
curl -X DELETE http://localhost:8080/users/2

# 健康检查
curl http://localhost:8080/health
```

## 响应格式

所有 API 返回统一格式：

```json
{
  "code": 0,
  "message": "success",
  "data": { ... }
}
```

## HTTP 状态码

| 状态码 | 含义 |
|--------|------|
| 200 | 成功 |
| 201 | 创建成功 |
| 400 | 请求参数错误 |
| 404 | 资源不存在 |
| 500 | 服务器内部错误 |
