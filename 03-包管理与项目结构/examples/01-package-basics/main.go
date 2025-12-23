// 01-package-basics: 包的导入与导出
//
// 📌 最佳实践:
//   - 大写开头 = 导出（public）
//   - 小写开头 = 私有（package-private）
//   - 包名应简短、小写、单词
//   - 避免使用下划线或驼峰命名包
//
// 运行: go run .
package main

import (
	"fmt"

	"package-basics/greeting" // 导入子包
)

func main() {
	// 调用导出的函数
	msg := greeting.Hello("Tom")
	fmt.Println(msg)

	// 调用导出的函数（使用包内私有函数）
	formal := greeting.FormalGreeting("Dr.", "Smith")
	fmt.Println(formal)

	// 访问导出的常量
	fmt.Printf("最大名字长度: %d\n", greeting.MaxNameLength)

	// 访问导出的变量
	fmt.Printf("默认语言: %s\n", greeting.DefaultLanguage)

	// 下面这行会编译错误，因为 formatName 是私有的:
	// greeting.formatName("test")  // 错误！
}
