// 作业4：计数器 - 理解指针与值传递
//
// 📌 学习目标：
//   - 理解 Go 的值传递机制
//   - 理解为什么需要指针
//   - 区分值接收者和指针接收者的效果
//   - 理解 Java 程序员常见的困惑点
//
// 📌 要求：
//   1. 定义 Counter 结构体，包含 value int 字段
//   2. 实现 IncrementWrong() 方法（值接收者）—— 故意写"错"
//   3. 实现 IncrementRight() 方法（指针接收者）—— 正确实现
//   4. 实现 Value() 方法返回当前值
//   5. 运行并观察两种方法的区别，理解为什么一个"不工作"
//
// 📌 思考题：
//   - 为什么 IncrementWrong 不能修改原对象？
//   - Java 中为什么没有这个问题？（提示：Java 对象变量本身就是引用）
//   - 什么时候用值接收者，什么时候用指针接收者？
//
// 🆚 与 Java 对比：
//   Java:
//     Counter c = new Counter();  // c 是引用（隐式指针）
//     c.increment();              // 总是能修改原对象
//
//   Go:
//     c := Counter{}              // c 是值
//     c.IncrementWrong()          // 方法收到的是 c 的拷贝！
//     c.IncrementRight()          // 需要指针接收者才能修改
//
// 📌 运行：go run main.go
//
// 📌 预期输出：
//   === 值接收者（错误示范）===
//   调用前: 0
//   调用 IncrementWrong() 3 次...
//   调用后: 0  <-- 没变！因为修改的是拷贝
//
//   === 指针接收者（正确做法）===
//   调用前: 0
//   调用 IncrementRight() 3 次...
//   调用后: 3  <-- 正确修改
package main

import "fmt"

// TODO: 1. 定义 Counter 结构体
// type Counter struct {
//     value int
// }

// TODO: 2. 实现 IncrementWrong（值接收者）
// 这个方法"不工作"，用于演示值传递的问题
// func (c Counter) IncrementWrong() {
//     c.value++  // 这里修改的是 c 的拷贝！
// }

// TODO: 3. 实现 IncrementRight（指针接收者）
// 这个方法正确工作
// func (c *Counter) IncrementRight() {
//     c.value++  // 这里修改的是原对象
// }

// TODO: 4. 实现 Value() 方法
// func (c Counter) Value() int {
//     return c.value
// }

func main() {
	fmt.Println("=== 值接收者（错误示范）===")
	// TODO: 5. 演示 IncrementWrong 的问题
	// c1 := Counter{}
	// fmt.Printf("调用前: %d\n", c1.Value())
	// c1.IncrementWrong()
	// c1.IncrementWrong()
	// c1.IncrementWrong()
	// fmt.Printf("调用后: %d  <-- 没变！\n", c1.Value())

	fmt.Println("\n=== 指针接收者（正确做法）===")
	// TODO: 6. 演示 IncrementRight 的效果
	// c2 := Counter{}
	// fmt.Printf("调用前: %d\n", c2.Value())
	// c2.IncrementRight()
	// c2.IncrementRight()
	// c2.IncrementRight()
	// fmt.Printf("调用后: %d  <-- 正确！\n", c2.Value())

	fmt.Println("\n=== 额外练习：函数参数的值传递 ===")
	// TODO: 7. 实现两个函数来演示函数参数的值传递
	// func doubleWrong(n int) { n = n * 2 }       // 不工作
	// func doubleRight(n *int) { *n = *n * 2 }    // 工作
	//
	// x := 10
	// doubleWrong(x)
	// fmt.Printf("doubleWrong 后: %d\n", x)  // 仍然是 10
	//
	// doubleRight(&x)
	// fmt.Printf("doubleRight 后: %d\n", x)  // 变成 20

	// 以下是占位代码，完成后删除
	fmt.Println("请完成作业")
}
