// 01-middleware: 中间件模式
//
// 📌 中间件原理:
//   - 中间件是一个函数，接收 Handler 返回 Handler
//   - 形成洋葱模型: middleware1 → middleware2 → handler → middleware2 → middleware1
//
// 📌 最佳实践:
//   - 日志、认证、限流等通用逻辑放在中间件
//   - 中间件应该职责单一
//   - 注意中间件的执行顺序
//
// 运行: go run main.go
// 测试: curl http://localhost:8080/
package main

import (
	"log"
	"net/http"
	"time"
)

// Middleware 类型定义
type Middleware func(http.Handler) http.Handler

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", homeHandler)
	mux.HandleFunc("GET /users", usersHandler)

	// 应用中间件（从外到内）
	// 请求顺序: Logger → Timer → Handler
	// 响应顺序: Handler → Timer → Logger
	handler := Chain(mux, LoggerMiddleware, TimerMiddleware)

	addr := ":8080"
	log.Printf("服务器启动: http://localhost%s", addr)

	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}

// Chain 将多个中间件串联起来
func Chain(handler http.Handler, middlewares ...Middleware) http.Handler {
	// 从后向前包装，最后一个中间件最先执行
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}
	return handler
}

// LoggerMiddleware 日志中间件
func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("→ [%s] %s", r.Method, r.URL.Path)

		// 调用下一个处理器
		next.ServeHTTP(w, r)

		log.Printf("← [%s] %s 完成", r.Method, r.URL.Path)
	})
}

// TimerMiddleware 计时中间件
func TimerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// 调用下一个处理器
		next.ServeHTTP(w, r)

		duration := time.Since(start)
		log.Printf("⏱ %s %s 耗时: %v", r.Method, r.URL.Path, duration)
	})
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	time.Sleep(10 * time.Millisecond) // 模拟处理时间
	w.Write([]byte("首页"))
}

func usersHandler(w http.ResponseWriter, r *http.Request) {
	time.Sleep(50 * time.Millisecond) // 模拟处理时间
	w.Write([]byte("用户列表"))
}
