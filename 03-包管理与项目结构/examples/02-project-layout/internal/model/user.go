// Package model 定义数据模型
//
// 📌 internal 包特性:
//   - 只能被同一模块内的代码导入
//   - 外部模块无法导入，编译器强制保护
//   - 适合放业务模型、内部逻辑
package model

import "time"

// User 用户模型
type User struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}
