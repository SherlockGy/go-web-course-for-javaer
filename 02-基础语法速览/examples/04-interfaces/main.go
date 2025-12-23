// 04-interfaces: 接口定义与实现
//
// 📌 最佳实践:
//   - 接口应该小而专一（1-3 个方法）
//   - 在使用方定义接口，而非实现方
//   - 空接口 any (interface{}) 谨慎使用
//   - 接口命名：动词+er（Reader, Writer, Stringer）
//
// 🆚 Java 对比:
//   Java: class Dog implements Animal { ... }  // 显式声明
//   Go:   只要实现了方法，就自动实现接口        // 隐式实现
package main

import "fmt"

// Stringer 接口 - fmt 包的标准接口
// 只要实现 String() 方法，fmt.Println 就会调用它
type Stringer interface {
	String() string
}

// Speaker 接口 - 小而专一
type Speaker interface {
	Speak() string
}

// Mover 接口
type Mover interface {
	Move() string
}

// Animal 组合接口
type Animal interface {
	Speaker
	Mover
}

// Dog 结构体
type Dog struct {
	Name string
}

// Dog 实现 Speaker 接口（隐式）
func (d Dog) Speak() string {
	return fmt.Sprintf("%s: 汪汪!", d.Name)
}

// Dog 实现 Mover 接口（隐式）
func (d Dog) Move() string {
	return fmt.Sprintf("%s 在跑", d.Name)
}

// Dog 实现 Stringer 接口
func (d Dog) String() string {
	return fmt.Sprintf("Dog{Name: %s}", d.Name)
}

// Cat 结构体
type Cat struct {
	Name string
}

func (c Cat) Speak() string {
	return fmt.Sprintf("%s: 喵~", c.Name)
}

func (c Cat) Move() string {
	return fmt.Sprintf("%s 在走", c.Name)
}

func main() {
	// === 接口使用 ===
	dog := Dog{Name: "旺财"}
	cat := Cat{Name: "咪咪"}

	// 接口变量可以持有任何实现该接口的类型
	var speaker Speaker
	speaker = dog
	fmt.Println(speaker.Speak())

	speaker = cat
	fmt.Println(speaker.Speak())

	// === 接口切片 ===
	animals := []Animal{dog, cat}
	for _, a := range animals {
		fmt.Printf("%s, %s\n", a.Speak(), a.Move())
	}

	// === 多态函数 ===
	makeSpeak(dog)
	makeSpeak(cat)

	// === Stringer 接口 ===
	// fmt.Println 会自动调用 String() 方法
	fmt.Println(dog) // 输出: Dog{Name: 旺财}

	// === 类型断言 ===
	var animal Animal = dog

	// 方式1: 直接断言（可能 panic）
	d := animal.(Dog)
	fmt.Printf("断言成功: %s\n", d.Name)

	// 方式2: 安全断言（推荐）
	if d, ok := animal.(Dog); ok {
		fmt.Printf("是 Dog: %s\n", d.Name)
	}

	// === 类型开关 ===
	checkType(dog)
	checkType(cat)
	checkType("hello")

	// === 空接口 any ===
	// any 是 interface{} 的别名（Go 1.18+）
	var anything any
	anything = 42
	anything = "hello"
	anything = dog
	fmt.Printf("any 可以持有任何类型: %v\n", anything)
}

// 多态函数 - 接受接口类型参数
func makeSpeak(s Speaker) {
	fmt.Println(s.Speak())
}

// 类型开关
func checkType(v any) {
	switch t := v.(type) {
	case Dog:
		fmt.Printf("这是一只狗: %s\n", t.Name)
	case Cat:
		fmt.Printf("这是一只猫: %s\n", t.Name)
	case string:
		fmt.Printf("这是字符串: %s\n", t)
	default:
		fmt.Printf("未知类型: %T\n", t)
	}
}
