// 03-file-upload: 文件上传
//
// 📌 单文件上传:
//   c.FormFile("file")
//
// 📌 多文件上传:
//   form.MultipartForm
//
// 📌 最佳实践:
//   - 限制文件大小
//   - 验证文件类型
//   - 使用安全的文件名
package main

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// 设置最大上传大小（8MB）
	r.MaxMultipartMemory = 8 << 20

	r.POST("/upload", uploadSingle)
	r.POST("/upload/multiple", uploadMultiple)

	r.Run(":8080")
}

// 单文件上传
func uploadSingle(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "获取文件失败: " + err.Error(),
		})
		return
	}

	// 验证文件类型
	if !isAllowedExt(file.Filename) {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "不支持的文件类型",
		})
		return
	}

	// 验证文件大小（5MB）
	if file.Size > 5*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "文件大小不能超过 5MB",
		})
		return
	}

	// 生成安全的文件名
	ext := filepath.Ext(file.Filename)
	newFilename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)

	// 保存文件
	dst := filepath.Join("./uploads", newFilename)
	if err := c.SaveUploadedFile(file, dst); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "保存文件失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "上传成功",
		"data": gin.H{
			"filename":      newFilename,
			"original_name": file.Filename,
			"size":          file.Size,
			"url":           "/files/" + newFilename,
		},
	})
}

// 多文件上传
func uploadMultiple(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "获取表单失败",
		})
		return
	}

	files := form.File["files"]
	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "没有上传文件",
		})
		return
	}

	var uploaded []gin.H
	for _, file := range files {
		// 验证
		if !isAllowedExt(file.Filename) {
			continue
		}
		if file.Size > 5*1024*1024 {
			continue
		}

		// 保存
		ext := filepath.Ext(file.Filename)
		newFilename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
		dst := filepath.Join("./uploads", newFilename)

		if err := c.SaveUploadedFile(file, dst); err != nil {
			continue
		}

		uploaded = append(uploaded, gin.H{
			"filename":      newFilename,
			"original_name": file.Filename,
			"size":          file.Size,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": fmt.Sprintf("成功上传 %d 个文件", len(uploaded)),
		"data":    uploaded,
	})
}

// 检查文件扩展名
func isAllowedExt(filename string) bool {
	allowed := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".pdf":  true,
	}
	ext := strings.ToLower(filepath.Ext(filename))
	return allowed[ext]
}
