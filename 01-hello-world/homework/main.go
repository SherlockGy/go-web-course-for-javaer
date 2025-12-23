// 作业：输出当前时间
//
// 📌 要求：
//  1. 使用 time.DateTime 常量格式化时间（Go 1.20+ 新特性）
//  2. 输出格式：当前时间: 2024-01-15 14:30:00
//  3. 不要使用硬编码的时间格式字符串 "2006-01-02 15:04:05"
//
// 📌 提示：
//   - now := time.Now() 获取当前时间
//   - now.Format(layout) 格式化时间（now 是 time.Time 类型）
//   - time.DateTime 是 Go 1.20+ 预定义的时间格式常量
//   - fmt.Printf 可以格式化输出
//
// 📌 运行：go run main.go
package main

import (
	"fmt"
	"time"
)

func main() {
	// TODO: 1. 获取当前时间
	// 提示：使用 time.Now()

	// TODO: 2. 使用 time.DateTime 格式化并输出
	// 提示：使用 fmt.Printf("当前时间: %s\n", ...)

	// 以下是占位输出，请替换为你的实现
	fmt.Println("请完成作业")
	_ = time.Now      // 删除这行，这只是为了避免 import 报错
	_ = time.DateTime // 删除这行，这只是为了避免 import 报错
}
