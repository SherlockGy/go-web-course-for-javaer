// 作业：用户查找功能
//
// 📌 要求：
//   1. 定义 User 结构体，包含 ID (int64)、Name (string)、Email (string) 字段
//   2. 实现 String() 方法，满足 fmt.Stringer 接口
//   3. 定义 ErrInvalidID 哨兵错误
//   4. 编写 FindUser(id int64) (*User, error) 函数：
//      - id > 0：返回模拟用户（可以用固定数据）
//      - id <= 0：返回 ErrInvalidID 错误
//   5. 在 main 函数中演示：成功查找 和 错误处理
//
// 📌 提示：
//   - errors.New("错误信息") 创建错误
//   - errors.Is(err, target) 检查错误
//   - fmt.Sprintf() 格式化字符串
//   - 结构体方法：func (u User) String() string { ... }
//   - 指针返回：return &User{...}, nil
//
// 📌 运行：go run main.go
//
// 📌 预期输出示例：
//   === 成功查找 ===
//   User{ID: 1, Name: xxx, Email: xxx}
//
//   === 错误处理 ===
//   错误: 无效的用户 ID
package main

import (
	"errors"
	"fmt"
)

// TODO: 1. 定义 ErrInvalidID 哨兵错误
// 提示：var ErrInvalidID = errors.New("...")

// TODO: 2. 定义 User 结构体
// 提示：type User struct { ... }

// TODO: 3. 为 User 实现 String() 方法
// 提示：func (u User) String() string { ... }
// 🆚 Java: 类似于重写 toString() 方法

// TODO: 4. 实现 FindUser 函数
// 提示：func FindUser(id int64) (*User, error) { ... }

func main() {
	fmt.Println("=== 成功查找 ===")
	// TODO: 5. 调用 FindUser(1)，打印返回的用户

	fmt.Println("\n=== 错误处理 ===")
	// TODO: 6. 调用 FindUser(-1)，使用 errors.Is 检查错误并打印

	// 以下是占位代码，请删除并替换为你的实现
	fmt.Println("请完成作业")
	_ = errors.New // 删除这行
}
