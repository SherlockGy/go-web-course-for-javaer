// 04-query-params: URL 查询参数
//
// 📌 获取 Query 参数:
//   - r.URL.Query() 返回 url.Values (map[string][]string)
//   - query.Get("key") 获取单个值
//   - query["key"] 获取多个值（数组参数）
//
// 运行: go run main.go
// 测试:
//   curl "http://localhost:8080/search?q=golang&page=2&limit=10"
//   curl "http://localhost:8080/filter?tags=go&tags=web&tags=api"
package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /search", searchHandler)
	mux.HandleFunc("GET /filter", filterHandler)

	addr := ":8080"
	log.Printf("服务器启动: http://localhost%s", addr)
	log.Println("测试命令:")
	log.Println(`  curl "http://localhost:8080/search?q=golang&page=2&limit=10"`)
	log.Println(`  curl "http://localhost:8080/filter?tags=go&tags=web&tags=api"`)

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	// 获取单个参数
	q := query.Get("q")
	if q == "" {
		q = "默认搜索词"
	}

	// 获取整数参数（带默认值）
	page := getIntParam(query.Get("page"), 1)
	limit := getIntParam(query.Get("limit"), 20)

	// 限制最大值
	if limit > 100 {
		limit = 100
	}

	writeJSON(w, http.StatusOK, Response{
		Code:    0,
		Message: "success",
		Data: map[string]any{
			"query": q,
			"page":  page,
			"limit": limit,
			"results": []string{
				"结果1: " + q,
				"结果2: " + q,
			},
		},
	})
}

func filterHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	// 获取数组参数（同名多个参数）
	// URL: /filter?tags=go&tags=web&tags=api
	tags := query["tags"] // []string{"go", "web", "api"}

	if len(tags) == 0 {
		tags = []string{"default"}
	}

	writeJSON(w, http.StatusOK, Response{
		Code:    0,
		Message: "success",
		Data: map[string]any{
			"tags":  tags,
			"count": len(tags),
		},
	})
}

// getIntParam 获取整数参数，转换失败返回默认值
func getIntParam(s string, defaultVal int) int {
	if s == "" {
		return defaultVal
	}
	val, err := strconv.Atoi(s)
	if err != nil {
		return defaultVal
	}
	return val
}

func writeJSON(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}
