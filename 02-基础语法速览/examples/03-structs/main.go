// 03-structs: 结构体与方法
//
// 📌 最佳实践:
//   - 结构体名大写开头 = 公开（可导出）
//   - 字段名大写开头 = 公开
//   - 方法接收者：修改状态用指针，只读用值
//   - 推荐统一使用指针接收者（一致性）
package main

import "fmt"

// User 结构体定义
// 大写开头 = 可导出（public）
type User struct {
	ID       int // 公开字段
	Username string
	Email    string
	password string // 小写 = 私有（仅包内可见）
}

// NewUser 构造函数（Go 惯例：New + 类型名）
func NewUser(id int, username, email, password string) *User {
	return &User{
		ID:       id,
		Username: username,
		Email:    email,
		password: password,
	}
}

// 值接收者方法 - 不修改原对象
func (u User) DisplayName() string {
	return fmt.Sprintf("%s <%s>", u.Username, u.Email)
}

// 指针接收者方法 - 可以修改原对象
func (u *User) UpdateEmail(newEmail string) {
	u.Email = newEmail
}

// 指针接收者 - 避免拷贝大对象
func (u *User) CheckPassword(pwd string) bool {
	return u.password == pwd
}

func main() {
	// === 创建结构体实例 ===

	// 方式1: 字面量
	user1 := User{
		ID:       1,
		Username: "tom",
		Email:    "tom@example.com",
	}
	fmt.Printf("user1: %+v\n", user1)

	// 方式2: 构造函数（推荐）
	user2 := NewUser(2, "jerry", "jerry@example.com", "secret123")
	fmt.Printf("user2: %+v\n", user2)

	// 方式3: new() - 返回指针，所有字段为零值
	user3 := new(User)
	user3.ID = 3
	user3.Username = "alice"
	fmt.Printf("user3: %+v\n", user3)

	// === 调用方法 ===
	fmt.Printf("DisplayName: %s\n", user2.DisplayName())

	// 修改字段
	user2.UpdateEmail("jerry.new@example.com")
	fmt.Printf("新邮箱: %s\n", user2.Email)

	// 密码检查
	fmt.Printf("密码正确: %t\n", user2.CheckPassword("secret123"))
	fmt.Printf("密码错误: %t\n", user2.CheckPassword("wrong"))

	// === 结构体嵌入（组合）===
	admin := Admin{
		User: User{
			ID:       100,
			Username: "admin",
			Email:    "admin@example.com",
		},
		Role: "super_admin",
	}

	// 可以直接访问嵌入结构体的字段和方法
	fmt.Printf("Admin username: %s\n", admin.Username) // 不需要 admin.User.Username
	fmt.Printf("Admin display: %s\n", admin.DisplayName())
	fmt.Printf("Admin role: %s\n", admin.Role)

	// === 匿名结构体（临时使用）===
	response := struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}{
		Code:    200,
		Message: "success",
	}
	fmt.Printf("response: %+v\n", response)
}

// Admin 嵌入 User（组合优于继承）
type Admin struct {
	User // 匿名嵌入
	Role string
}
