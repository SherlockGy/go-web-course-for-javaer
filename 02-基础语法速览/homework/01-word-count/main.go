// 作业1：词频统计
//
// 📌 学习目标：
//   - 理解 map 的零值陷阱（必须 make 初始化）
//   - 利用 int 零值简化计数逻辑
//   - 使用 range 遍历 slice 和 map
//   - 理解 nil slice vs 空 slice
//
// 📌 要求：
//   1. 实现 CountWords(words []string) map[string]int
//      - 统计每个单词出现的次数
//      - 利用 int 零值特性（不需要判断 key 是否存在）
//   2. 实现 TopWords(counts map[string]int, n int) []string
//      - 返回出现次数最多的 n 个单词
//      - 如果单词数不足 n 个，返回全部
//   3. 思考题：TopWords 返回 nil 还是空 slice？为什么？
//
// 📌 提示：
//   - make(map[string]int) 初始化 map
//   - counts[word]++ 利用零值直接自增
//   - for _, word := range words 遍历 slice
//   - for word, count := range counts 遍历 map
//
// 📌 运行：go run main.go
//
// 📌 预期输出示例：
//   词频统计: map[apple:3 banana:2 cherry:1 dog:1 elephant:1]
//   Top 3: [apple banana cherry] 或 [apple banana dog]（顺序可能不同）
package main

import "fmt"

// TODO: 1. 实现 CountWords 函数
// 统计单词出现次数
// func CountWords(words []string) map[string]int {
//     ...
// }

// TODO: 2. 实现 TopWords 函数
// 返回出现次数最多的 n 个单词
// 提示：可以简单实现，不需要严格排序，找出 top n 即可
// func TopWords(counts map[string]int, n int) []string {
//     ...
// }

func main() {
	words := []string{
		"apple", "banana", "apple", "cherry",
		"banana", "apple", "dog", "elephant",
	}

	// TODO: 3. 调用 CountWords 并打印结果
	fmt.Println("=== 词频统计 ===")
	// counts := CountWords(words)
	// fmt.Printf("词频统计: %v\n", counts)

	// TODO: 4. 调用 TopWords 并打印结果
	fmt.Println("\n=== Top 3 单词 ===")
	// top3 := TopWords(counts, 3)
	// fmt.Printf("Top 3: %v\n", top3)

	// TODO: 5. 测试边界情况：空 slice
	fmt.Println("\n=== 边界测试 ===")
	// emptyResult := CountWords([]string{})
	// fmt.Printf("空输入结果: %v\n", emptyResult)
	// 思考：emptyResult 是 nil 还是空 map？

	// 以下是占位代码，完成后删除
	fmt.Println("请完成作业")
	_ = words
}
