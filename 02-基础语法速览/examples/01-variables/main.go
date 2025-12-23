// 01-variables: 变量与类型
//
// 📌 最佳实践:
//   - 优先使用 := 短声明（函数内部）
//   - 包级变量使用 var
//   - 常量使用 const，编译期确定
//   - 零值是有意义的: string="", int=0, bool=false, pointer=nil
package main

import "fmt"

// 包级变量必须用 var
var globalVar = "我是包级变量"

// 常量 - 编译期确定，不可修改
const (
	MaxRetries = 3
	AppName    = "MyApp"
)

func main() {
	// === 变量声明 ===

	// 方式1: var 声明（显式类型）
	var name string = "Tom"

	// 方式2: var 声明（类型推断）
	var age = 25

	// 方式3: 短声明 := （推荐，仅函数内可用）
	email := "tom@example.com"

	// 方式4: 多变量声明
	var x, y int = 1, 2
	a, b := "hello", true

	fmt.Printf("name=%s, age=%d, email=%s\n", name, age, email)
	fmt.Printf("x=%d, y=%d, a=%s, b=%t\n", x, y, a, b)

	// === 基本类型 ===

	// 整数
	var i int = 42      // 平台相关（32/64位）
	var i64 int64 = 100 // 明确64位
	var u uint = 10     // 无符号

	// 浮点
	var f float64 = 3.14

	// 布尔
	var ok bool = true

	// 字符串（不可变）
	var s string = "Hello, 世界"

	fmt.Printf("i=%d, i64=%d, u=%d, f=%.2f, ok=%t, s=%s\n", i, i64, u, f, ok, s)

	// === 零值 ===
	// Go 的变量总有值，未初始化时为"零值"
	var (
		zeroInt    int            // 0
		zeroFloat  float64        // 0.0
		zeroBool   bool           // false
		zeroString string         // ""
		zeroSlice  []int          // nil
		zeroMap    map[string]int // nil
	)
	fmt.Printf("零值: int=%d, float=%.1f, bool=%t, string=%q, slice=%v, map=%v\n",
		zeroInt, zeroFloat, zeroBool, zeroString, zeroSlice, zeroMap)

	// === Slice（切片）===
	// 动态数组，最常用的集合类型
	nums := []int{1, 2, 3}
	nums = append(nums, 4, 5) // 追加元素
	fmt.Printf("slice: %v, len=%d, cap=%d\n", nums, len(nums), cap(nums))

	// 切片操作
	sub := nums[1:3] // [2, 3] - 左闭右开
	fmt.Printf("sub slice: %v\n", sub)

	// === Map（映射）===
	// 键值对集合
	scores := map[string]int{
		"Tom":   90,
		"Jerry": 85,
	}
	scores["Alice"] = 95 // 添加

	// 检查键是否存在（重要！）
	if score, exists := scores["Tom"]; exists {
		fmt.Printf("Tom's score: %d\n", score)
	}

	// 删除
	delete(scores, "Jerry")
	fmt.Printf("scores: %v\n", scores)
}
