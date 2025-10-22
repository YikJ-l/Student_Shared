package api

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	req "student_shared/app/model/req"
	resp "student_shared/app/model/resp"
)


// UploadAvatar 上传头像
func UploadAvatar(c *gin.Context) {
	// 检查用户是否已登录
	_, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	// 绑定上传的文件
	var form req.UploadAvatarRequest
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请选择要上传的文件"})
		return
	}
	file := form.File

	// 验证文件类型
	if !isValidImageType(file.Filename) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "只支持 JPG、JPEG、PNG、GIF 格式的图片"})
		return
	}

	// 验证文件大小 (2MB)
	if file.Size > 2*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "文件大小不能超过 2MB"})
		return
	}

	// 创建上传目录
	uploadDir := "uploads/avatars"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建上传目录失败"})
		return
	}

	// 生成唯一文件名
	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%s_%d%s", uuid.New().String(), time.Now().Unix(), ext)
	filePath := filepath.Join(uploadDir, filename)

	// 保存文件
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存文件失败"})
		return
	}

	// 生成访问URL
	fileURL := fmt.Sprintf("/uploads/avatars/%s", filename)

	// 返回成功响应
	c.JSON(http.StatusOK, resp.UploadResponse{
		URL:  fileURL,
		Size: file.Size,
	})
}

// isValidImageType 验证图片类型
func isValidImageType(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	validExts := []string{".jpg", ".jpeg", ".png", ".gif"}

	for _, validExt := range validExts {
		if ext == validExt {
			return true
		}
	}
	return false
}

// DeleteAvatar 删除头像文件
func DeleteAvatar(c *gin.Context) {
	// 检查用户是否已登录
	_, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	// 获取要删除的文件路径
	var req struct {
		Path string `json:"path" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数无效"})
		return
	}

	// 验证文件路径安全性
	if !strings.HasPrefix(req.Path, "uploads/avatars/") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的文件路径"})
		return
	}

	// 删除文件
	if err := os.Remove(req.Path); err != nil {
		// 文件不存在也认为删除成功
		if !os.IsNotExist(err) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "删除文件失败"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

// GetAvatarFile 获取头像文件
func GetAvatarFile(c *gin.Context) {
	filename := c.Param("filename")
	if filename == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "文件名不能为空"})
		return
	}

	// 构建文件路径
	filePath := filepath.Join("uploads/avatars", filename)

	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "文件不存在"})
		return
	}

	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "打开文件失败"})
		return
	}
	defer file.Close()

	// 设置响应头
	ext := filepath.Ext(filename)
	contentType := getContentType(ext)
	c.Header("Content-Type", contentType)
	c.Header("Cache-Control", "public, max-age=31536000") // 缓存一年

	// 返回文件内容
	io.Copy(c.Writer, file)
}

// getContentType 根据文件扩展名获取Content-Type
func getContentType(ext string) string {
	switch strings.ToLower(ext) {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	default:
		return "application/octet-stream"
	}
}
