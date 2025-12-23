// model/user.go - 数据模型层
//
// 📌 模型层职责:
//   - 定义数据结构（对应数据库表）
//   - 定义 DTO（数据传输对象）
//   - 不包含业务逻辑
//
// 📌 与 Java 对比:
//   - Java: Entity + DTO 分离，@Entity 注解
//   - Go: struct tag，更简洁
package model

import "time"

// User 用户实体（对应数据库表）
// 📌 与 Java @Entity 类似
type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Username  string    `json:"username" gorm:"uniqueIndex;size:50"`
	Email     string    `json:"email" gorm:"uniqueIndex;size:100"`
	Password  string    `json:"-" gorm:"size:100"` // json:"-" 序列化时忽略
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateUserRequest 创建用户请求 DTO
// 📌 DTO 用于接收请求参数，与实体分离
type CreateUserRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// UpdateUserRequest 更新用户请求 DTO
type UpdateUserRequest struct {
	Username string `json:"username" binding:"omitempty,min=3,max=50"`
	Email    string `json:"email" binding:"omitempty,email"`
}

// UserResponse 用户响应 DTO
// 📌 控制返回给前端的字段
type UserResponse struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

// ToResponse 转换为响应 DTO
func (u *User) ToResponse() *UserResponse {
	return &UserResponse{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
	}
}
