// 01-first-program: 第一个 Go 程序
//
// 📌 运行方式:
//
//	go run main.go
//
// 📌 编译运行:
//
//	go build -o hello.exe && ./hello.exe      # Windows
//	go build -o hello && ./hello              # Linux/Mac
//
// 📌 与 Java 对比:
//   - Java: public class + public static void main(String[] args)
//   - Go: package main + func main()，更简洁
//   - Java 需要 JVM，Go 直接编译成原生可执行文件
package main

import "fmt"

// main 是程序入口
// 📌 Go 没有类的概念，函数直接定义在包级别
func main() {
	fmt.Println("Hello, Go!")

	// 📌 fmt.Println vs fmt.Printf
	// Println: 自动换行，参数用空格分隔
	// Printf: 格式化输出，需要手动 \n
	name := "Gopher"
	fmt.Println("Hello,", name)      // Hello, Gopher
	fmt.Printf("Hello, %s!\n", name) // Hello, Gopher!
}
