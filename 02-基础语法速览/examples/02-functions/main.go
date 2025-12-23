// 02-functions: 函数与多返回值
//
// 📌 最佳实践:
//   - 函数返回 error 作为最后一个返回值
//   - 调用后立即检查 error
//   - 使用命名返回值提高可读性（复杂函数）
//   - 函数名用驼峰命名
package main

import (
	"errors"
	"fmt"
)

func main() {
	// === 基本函数调用 ===
	result := add(1, 2)
	fmt.Printf("1 + 2 = %d\n", result)

	// === 多返回值 ===
	sum, diff := sumAndDiff(10, 3)
	fmt.Printf("sum=%d, diff=%d\n", sum, diff)

	// === 返回 error（最重要的模式）===
	user, err := findUser(1)
	if err != nil {
		fmt.Printf("查找失败: %v\n", err)
	} else {
		fmt.Printf("找到用户: %s\n", user)
	}

	// 找不到的情况
	_, err = findUser(999)
	if err != nil {
		fmt.Printf("查找失败: %v\n", err)
	}

	// === 可变参数 ===
	total := sumAll(1, 2, 3, 4, 5)
	fmt.Printf("1+2+3+4+5 = %d\n", total)

	// 也可以传入切片（展开）
	nums := []int{10, 20, 30}
	total = sumAll(nums...)
	fmt.Printf("10+20+30 = %d\n", total)

	// === 函数作为值 ===
	op := multiply // 函数赋值给变量
	fmt.Printf("3 * 4 = %d\n", op(3, 4))

	// === 匿名函数（闭包）===
	counter := makeCounter()
	fmt.Printf("count: %d\n", counter()) // 1
	fmt.Printf("count: %d\n", counter()) // 2
	fmt.Printf("count: %d\n", counter()) // 3

	// === defer（延迟执行）===
	deferDemo()
}

// 基本函数
func add(a, b int) int {
	return a + b
}

// 多返回值
func sumAndDiff(a, b int) (int, int) {
	return a + b, a - b
}

// 命名返回值 - 复杂函数时提高可读性
func divide(a, b float64) (result float64, err error) {
	if b == 0 {
		err = errors.New("除数不能为零")
		return // 直接返回命名返回值
	}
	result = a / b
	return
}

// 返回 error 模式 - Go 最核心的错误处理
func findUser(id int) (string, error) {
	// 模拟数据库查询
	users := map[int]string{
		1: "Tom",
		2: "Jerry",
	}

	if user, ok := users[id]; ok {
		return user, nil
	}
	return "", fmt.Errorf("用户不存在: id=%d", id)
}

// 可变参数
func sumAll(nums ...int) int {
	total := 0
	for _, n := range nums {
		total += n
	}
	return total
}

// 函数类型
func multiply(a, b int) int {
	return a * b
}

// 返回函数（闭包）
func makeCounter() func() int {
	count := 0
	return func() int {
		count++
		return count
	}
}

// defer 演示
func deferDemo() {
	fmt.Println("=== defer 演示 ===")
	defer fmt.Println("3. 最后执行（defer）")
	fmt.Println("1. 先执行")
	fmt.Println("2. 再执行")
	// defer 常用于：关闭文件、释放锁、记录日志
}
