// 作业2：用户管理
//
// 📌 学习目标：
//   - 定义结构体和构造函数
//   - 实现值接收者和指针接收者方法
//   - 理解 Go 的隐式接口实现
//   - 使用 fmt.Stringer 接口
//
// 📌 要求：
//   1. 定义 User 结构体：ID (int64)、Name (string)、Email (string)、Active (bool)
//   2. 实现 NewUser 构造函数，返回 *User
//   3. 实现 String() 方法（值接收者），满足 fmt.Stringer 接口
//   4. 实现 Deactivate() 方法（指针接收者），将 Active 设为 false
//   5. 实现 Activate() 方法（指针接收者），将 Active 设为 true
//   6. 思考：String() 为什么用值接收者，Activate() 为什么用指针接收者？
//
// 📌 提示：
//   - 构造函数惯例：func NewUser(...) *User { return &User{...} }
//   - 值接收者：func (u User) Method()
//   - 指针接收者：func (u *User) Method()
//   - fmt.Stringer 接口只要求实现 String() string 方法
//
// 📌 运行：go run main.go
//
// 📌 预期输出示例：
//   新用户: User{ID: 1, Name: Tom, Email: tom@example.com, Active: true}
//   停用后: User{ID: 1, Name: Tom, Email: tom@example.com, Active: false}
//   激活后: User{ID: 1, Name: Tom, Email: tom@example.com, Active: true}
package main

import "fmt"

// TODO: 1. 定义 User 结构体
// type User struct {
//     ...
// }

// TODO: 2. 实现 NewUser 构造函数
// 提示：新用户默认 Active 为 true
// func NewUser(id int64, name, email string) *User {
//     ...
// }

// TODO: 3. 实现 String() 方法（值接收者）
// 🆚 Java: 类似于 @Override public String toString()
// func (u User) String() string {
//     ...
// }

// TODO: 4. 实现 Deactivate() 方法（指针接收者）
// func (u *User) Deactivate() {
//     ...
// }

// TODO: 5. 实现 Activate() 方法（指针接收者）
// func (u *User) Activate() {
//     ...
// }

func main() {
	// TODO: 6. 创建用户并演示各方法
	fmt.Println("=== 创建用户 ===")
	// user := NewUser(1, "Tom", "tom@example.com")
	// fmt.Printf("新用户: %s\n", user)

	fmt.Println("\n=== 停用用户 ===")
	// user.Deactivate()
	// fmt.Printf("停用后: %s\n", user)

	fmt.Println("\n=== 激活用户 ===")
	// user.Activate()
	// fmt.Printf("激活后: %s\n", user)

	// TODO: 7. 验证接口实现
	fmt.Println("\n=== 接口验证 ===")
	// var s fmt.Stringer = user  // 如果编译通过，说明 User 实现了 Stringer
	// fmt.Printf("作为 Stringer: %s\n", s.String())

	// 以下是占位代码，完成后删除
	fmt.Println("请完成作业")
}
