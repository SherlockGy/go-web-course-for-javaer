// 01-variables: 变量、类型与零值深度解析
//
// 📌 最佳实践:
//   - 优先使用 := 短声明（函数内部）
//   - 包级变量使用 var
//   - 常量使用 const，编译期确定
//   - 零值是有意义的: string="", int=0, bool=false, pointer=nil
//   - 理解零值可用 vs 零值陷阱的区别
//
// 🆚 与 Java 对比:
//   - Java 的 null 需要防御性检查，Go 的零值让很多类型开箱即用
//   - 但 map、channel 的零值(nil)仍有陷阱，需要注意
//
// 📌 运行: go run main.go
package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

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

	// === 零值基础 ===
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

	// === 零值可用：strings.Builder ===
	// 📌 Java 中 StringBuilder sb = null; sb.append() 会 NPE
	// 🆚 Go 的 strings.Builder 零值可以直接使用！
	var sb strings.Builder // 无需 new 或 make
	sb.WriteString("Hello")
	sb.WriteString(", ")
	sb.WriteString("World!")
	fmt.Printf("strings.Builder 零值可用: %s\n", sb.String())

	// === 零值陷阱1：nil map ===
	// ⚠️ nil map 可以读取（返回零值），但写入会 panic
	var nilMap map[string]int
	fmt.Printf("nil map 读取: %d (不会panic)\n", nilMap["any"]) // 返回 int 零值 0

	// 下面这行会 panic，已注释
	// nilMap["key"] = 1  // 💥 panic: assignment to entry in nil map

	// 📌 正确做法：使用 make 初始化
	initMap := make(map[string]int)
	initMap["key"] = 1 // ✅ 可以写入
	fmt.Printf("初始化后的 map: %v\n", initMap)

	// === 零值陷阱2：nil slice vs 空 slice ===
	var nilSliceDemo []int          // nil slice
	emptySliceDemo := []int{}       // 空 slice（非 nil）
	makeSliceDemo := make([]int, 0) // 空 slice（非 nil）

	// 功能上几乎等价：都可以 append
	nilSliceDemo = append(nilSliceDemo, 1)
	emptySliceDemo = append(emptySliceDemo, 1)
	fmt.Printf("nil slice append 后: %v\n", nilSliceDemo)
	fmt.Printf("空 slice append 后: %v\n", emptySliceDemo)

	// ⚠️ 但 JSON 序列化结果不同！
	var forJSON []int
	emptyForJSON := []int{}

	nilJSON, _ := json.Marshal(forJSON)
	emptyJSON, _ := json.Marshal(emptyForJSON)
	fmt.Printf("nil slice JSON: %s (注意是 null)\n", nilJSON)
	fmt.Printf("空 slice JSON: %s (注意是 [])\n", emptyJSON)

	// 判断是否为 nil
	fmt.Printf("nilSlice == nil: %t\n", forJSON == nil)
	fmt.Printf("emptySlice == nil: %t\n", emptyForJSON == nil)
	_ = makeSliceDemo // 避免未使用变量警告

	// === 零值陷阱3：time.Time ===
	// ⚠️ time.Time 零值是 0001-01-01，不是 null
	var zeroTime time.Time
	fmt.Printf("time.Time 零值: %s\n", zeroTime.Format(time.DateTime))
	fmt.Printf("是否为零值: %t (使用 IsZero() 判断)\n", zeroTime.IsZero())

	// 📌 实际使用中，用 IsZero() 判断时间是否已设置
	if zeroTime.IsZero() {
		fmt.Println("时间未设置，使用当前时间")
		zeroTime = time.Now()
	}
	fmt.Printf("设置后的时间: %s\n", zeroTime.Format(time.DateTime))

	// === 零值实战：map 计数器 ===
	// 📌 利用 int 零值简化代码
	// 🆚 Java 需要 getOrDefault(key, 0) + 1
	words := []string{"apple", "banana", "apple", "cherry", "banana", "apple"}
	wordCount := make(map[string]int)
	for _, word := range words {
		wordCount[word]++ // 零值 0，直接 ++，无需判断 key 是否存在
	}
	fmt.Printf("词频统计: %v\n", wordCount)

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
