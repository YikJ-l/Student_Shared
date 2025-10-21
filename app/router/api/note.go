package api

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"student_shared/app/model"
	"student_shared/app/utils/database"
	"time"

	"github.com/gin-gonic/gin"
)

// NoteRequest 笔记请求
type NoteRequest struct {
	Title       string `json:"title" binding:"required"`
	Content     string `json:"content"`
	Description string `json:"description"`
	CourseID    *uint  `json:"course_id"` // 课程ID，可以为nil表示不属于任何课程
	Status      string `json:"status"`    // public, private
}

// NoteUpdateRequest 笔记更新请求（支持部分更新）
type NoteUpdateRequest struct {
	Title       *string `json:"title"`
	Content     *string `json:"content"`
	Description *string `json:"description"`
	CourseID    *uint   `json:"course_id"`
	Status      *string `json:"status"` // public, private, draft
}

// NoteResponse 笔记响应
type NoteResponse struct {
	ID            uint      `json:"id"`
	Title         string    `json:"title"`
	Content       string    `json:"content"`
	Description   string    `json:"description"`
	FilePath      string    `json:"file_path,omitempty"`
	FileSize      int64     `json:"file_size"`
	FileType      string    `json:"file_type"`
	ViewCount     int       `json:"view_count"`
	DownloadCount int       `json:"download_count"`
	LikeCount     int       `json:"like_count"`
	CommentCount  int       `json:"comment_count"`
	Status        string    `json:"status"`
	UserID        uint      `json:"user_id"`
	CourseID      uint      `json:"course_id"`
	Username      string    `json:"username"`
	AuthorAvatar  string    `json:"author_avatar"`
	CourseName    string    `json:"course_name"`
	IsLiked       bool      `json:"is_liked"`
	IsFavorited   bool      `json:"is_favorited"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// ListNotes 获取笔记列表
func ListNotes(c *gin.Context) {
	// 获取查询参数
	courseID := c.Query("course_id")
	userID := c.Query("user_id")
	status := c.DefaultQuery("status", "public")

	// 分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	// 构建查询
	query := database.DB.Model(&model.Note{})

	// 应用过滤条件
	if courseID != "" {
		if courseID == "0" {
			// 查询没有关联课程的笔记
			query = query.Where("course_id IS NULL")
		} else {
			query = query.Where("course_id = ?", courseID)
		}
	}
	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}
	if status != "all" {
		query = query.Where("status = ?", status)
	}

	// 获取当前用户ID（如果已登录）
	currentUserID, exists := c.Get("userID")

	// 权限控制：只有在以下情况下才能看到私有笔记
	// 1. 查询自己的笔记（userID 参数等于当前用户ID）
	// 2. 管理员查询所有笔记
	canViewPrivate := false
	if exists {
		// 检查是否查询自己的笔记
		if userID != "" && fmt.Sprintf("%v", currentUserID) == userID {
			canViewPrivate = true
		}
		// 检查是否为管理员（可以根据需要添加管理员权限检查）
		// 这里可以添加管理员权限检查逻辑
	}

	// 如果不能查看私有笔记，则只显示公开笔记
	if !canViewPrivate {
		query = query.Where("status = ?", "public")
	}

	// 获取总数
	var total int64
	query.Count(&total)

	// 排序参数
	sortBy := c.DefaultQuery("sort_by", "created_at")
	order := c.DefaultQuery("order", "desc")

	var orderClause string
	switch sortBy {
	case "like_count":
		orderClause = "like_count " + strings.ToUpper(order)
	case "view_count":
		orderClause = "view_count " + strings.ToUpper(order)
	case "created_at":
		orderClause = "created_at " + strings.ToUpper(order)
	default:
		orderClause = "created_at DESC"
	}

	// 获取分页数据，并预加载用户和课程信息
	var notes []model.Note
	result := query.Preload("User").Preload("Course").Limit(pageSize).Offset(offset).Order(orderClause).Find(&notes)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取笔记列表失败"})
		return
	}

	// 转换为响应格式
	var noteResponses []NoteResponse
	for _, note := range notes {
		// 计算评论数量
		var commentCount int64
		database.DB.Model(&model.Comment{}).Where("note_id = ?", note.ID).Count(&commentCount)

		// 处理CourseID和CourseName
		var courseID uint
		var courseName string
		if note.CourseID != nil {
			courseID = *note.CourseID
			courseName = note.Course.Name
		} else {
			courseID = 0
			courseName = ""
		}

		noteResponses = append(noteResponses, NoteResponse{
			ID:            note.ID,
			Title:         note.Title,
			Description:   note.Description,
			FileSize:      note.FileSize,
			FileType:      note.FileType,
			ViewCount:     note.ViewCount,
			DownloadCount: note.DownloadCount,
			LikeCount:     note.LikeCount,
			CommentCount:  int(commentCount),
			Status:        note.Status,
			UserID:        note.UserID,
			CourseID:      courseID,
			Username:      note.User.Username,
			AuthorAvatar:  note.User.Avatar,
			CourseName:    courseName,
			CreatedAt:     note.CreatedAt,
			UpdatedAt:     note.UpdatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"total":       total,
		"page":        page,
		"page_size":   pageSize,
		"total_pages": (int(total) + pageSize - 1) / pageSize,
		"notes":       noteResponses,
	})
}

// buildNoteResponse 构建笔记响应对象
func buildNoteResponse(note model.Note) NoteResponse {
	// 计算评论数量
	var commentCount int64
	database.DB.Model(&model.Comment{}).Where("note_id = ?", note.ID).Count(&commentCount)

	// 处理CourseID和CourseName
	var courseID uint
	var courseName string
	if note.CourseID != nil {
		courseID = *note.CourseID
		if note.Course.ID != 0 {
			courseName = note.Course.Name
		}
	} else {
		courseID = 0
		courseName = ""
	}

	return NoteResponse{
		ID:            note.ID,
		Title:         note.Title,
		Content:       note.Content,
		Description:   note.Description,
		FilePath:      note.FilePath,
		FileSize:      note.FileSize,
		FileType:      note.FileType,
		ViewCount:     note.ViewCount,
		DownloadCount: note.DownloadCount,
		LikeCount:     note.LikeCount,
		CommentCount:  int(commentCount),
		Status:        note.Status,
		UserID:        note.UserID,
		CourseID:      courseID,
		Username:      getUsernameFromNote(note),
		AuthorAvatar:  getAvatarFromNote(note),
		CourseName:    courseName,
		CreatedAt:     note.CreatedAt,
		UpdatedAt:     note.UpdatedAt,
	}
}

// GetNote 获取笔记详情
func GetNote(c *gin.Context) {
	// 获取笔记ID
	id := c.Param("id")
	
	// 验证ID是否为有效数字
	noteID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的笔记ID"})
		return
	}

	// 查询笔记
	var note model.Note
	result := database.DB.Preload("User").Preload("Course").First(&note, noteID)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "笔记不存在"})
		return
	}

	// 检查权限
	if note.Status == "private" {
		userID, exists := c.Get("userID")
		if !exists || userID.(uint) != note.UserID {
			c.JSON(http.StatusForbidden, gin.H{"error": "您没有权限查看此笔记"})
			return
		}
	}

	// 增加查看次数
	note.ViewCount++
	database.DB.Save(&note)

	// 检查用户是否已点赞和收藏（如果用户已登录）
	var isLiked, isFavorited bool
	currentUserID, exists := c.Get("userID")
	if exists {
		// 检查是否已点赞
		var noteLike model.NoteLike
		result := database.DB.Where("note_id = ? AND user_id = ?", note.ID, currentUserID).First(&noteLike)
		isLiked = result.RowsAffected > 0

		// 检查是否已收藏
		var favorite model.Favorite
		result = database.DB.Where("note_id = ? AND user_id = ?", note.ID, currentUserID).First(&favorite)
		isFavorited = result.RowsAffected > 0
	}

	// 处理CourseID和CourseName
	var courseID uint
	var courseName string
	if note.CourseID != nil {
		courseID = *note.CourseID
		if note.Course.ID != 0 {
			courseName = note.Course.Name
		}
	} else {
		courseID = 0
		courseName = ""
	}

	// 返回笔记信息
	c.JSON(http.StatusOK, NoteResponse{
		ID:            note.ID,
		Title:         note.Title,
		Content:       note.Content,
		Description:   note.Description,
		FilePath:      note.FilePath,
		FileSize:      note.FileSize,
		FileType:      note.FileType,
		ViewCount:     note.ViewCount,
		DownloadCount: note.DownloadCount,
		LikeCount:     note.LikeCount,
		Status:        note.Status,

		UserID:       note.UserID,
		CourseID:     courseID,
		Username:     note.User.Username,
		AuthorAvatar: note.User.Avatar,
		CourseName:   courseName,
		IsLiked:      isLiked,
		IsFavorited:  isFavorited,
		CreatedAt:    note.CreatedAt,
		UpdatedAt:    note.UpdatedAt,
	})
}

// CreateNote 创建笔记
func CreateNote(c *gin.Context) {
	// 获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	// 处理表单数据
	title := c.PostForm("title")
	content := c.PostForm("content")
	description := c.PostForm("description")
	courseIDStr := c.PostForm("course_id")
	status := c.DefaultPostForm("status", "public")

	if title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "标题不能为空"})
		return
	}

	// 处理课程ID，允许为空或0（表示无关联课程）
	var courseID uint = 0
	var err error
	if courseIDStr != "" {
		courseIDUint64, parseErr := strconv.ParseUint(courseIDStr, 10, 32)
		if parseErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无效的课程ID"})
			return
		}
		courseID = uint(courseIDUint64)
	}

	// 如果课程ID不为0，检查课程是否存在和用户权限
	if courseID != 0 {
		// 检查课程是否存在
		var course model.Course
		result := database.DB.First(&course, courseID)
		if result.Error != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "课程不存在"})
			return
		}

		// 检查用户是否已加入课程
		var userCourse model.UserCourse
		result = database.DB.Where("user_id = ? AND course_id = ?", userID, courseID).First(&userCourse)
		if result.Error != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "您需要先加入课程才能创建笔记"})
			return
		}
	}

	// 创建笔记对象
	note := model.Note{
		Title:       title,
		Content:     content,
		Description: description,
		UserID:      userID.(uint),
		Status:      status,
	}

	// 设置CourseID（使用指针类型）
	if courseID != 0 {
		note.CourseID = &courseID
	} else {
		note.CourseID = nil // 无所属课程
	}

	// 处理上传文件
	file, err := c.FormFile("file")
	if err == nil {
		// 创建上传目录
		uploadDir := "uploads/notes"
		if err := os.MkdirAll(uploadDir, 0755); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "创建上传目录失败"})
			return
		}

		// 生成文件名
		fileExt := filepath.Ext(file.Filename)
		fileName := fmt.Sprintf("%d_%d%s", userID, time.Now().Unix(), fileExt)
		filePath := filepath.Join(uploadDir, fileName)

		// 保存文件
		if err := c.SaveUploadedFile(file, filePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "保存文件失败"})
			return
		}

		// 更新笔记信息
		note.FilePath = filePath
		note.FileSize = file.Size
		note.FileType = fileExt[1:] // 去掉点号
	}

	// 保存笔记
	result := database.DB.Create(&note)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建笔记失败"})
		return
	}

	// 查询用户信息
	var user model.User
	database.DB.First(&user, userID)

	// 查询课程信息（如果课程ID不为nil）
	var courseName string
	if note.CourseID != nil {
		var course model.Course
		database.DB.First(&course, *note.CourseID)
		courseName = course.Name
		courseID = *note.CourseID
	} else {
		courseName = "" // 无所属课程
		courseID = 0
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "创建笔记成功",
		"note": NoteResponse{
			ID:            note.ID,
			Title:         note.Title,
			Content:       note.Content,
			Description:   note.Description,
			FilePath:      note.FilePath,
			FileSize:      note.FileSize,
			FileType:      note.FileType,
			ViewCount:     note.ViewCount,
			DownloadCount: note.DownloadCount,
			LikeCount:     note.LikeCount,
			Status:        note.Status,
			UserID:        note.UserID,
			CourseID:      courseID,
			Username:      user.Username,
			CourseName:    courseName,
			CreatedAt:     note.CreatedAt,
			UpdatedAt:     note.UpdatedAt,
		},
	})
}

// UpdateNote 更新笔记
func UpdateNote(c *gin.Context) {
	// 获取用户ID和角色
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	role, _ := c.Get("role")

	// 获取笔记ID
	id := c.Param("id")

	// 查询笔记
	var note model.Note
	result := database.DB.First(&note, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "笔记不存在"})
		return
	}

	// 检查权限
	if note.UserID != userID.(uint) && role.(string) != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "您没有权限更新此笔记"})
		return
	}

	// 绑定请求数据
	var req NoteUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数无效"})
		return
	}

	// 更新笔记（只更新提供的字段）
	if req.Title != nil {
		note.Title = *req.Title
	}
	if req.Content != nil {
		note.Content = *req.Content
	}
	if req.Description != nil {
		note.Description = *req.Description
	}
	if req.Status != nil {
		note.Status = *req.Status
	}

	// 如果更改了课程，检查新课程是否存在
	if req.CourseID != nil {
		// 检查是否真的需要更新CourseID
		needUpdate := false
		if note.CourseID == nil && *req.CourseID != 0 {
			needUpdate = true
		} else if note.CourseID != nil && *req.CourseID != *note.CourseID {
			needUpdate = true
		} else if note.CourseID != nil && *req.CourseID == 0 {
			needUpdate = true
		}

		if needUpdate {
			// 如果课程ID为0，表示笔记不属于任何课程
			if *req.CourseID == 0 {
				note.CourseID = nil
			} else {
				// 检查课程是否存在
				var course model.Course
				result = database.DB.First(&course, *req.CourseID)
				if result.Error != nil {
					c.JSON(http.StatusNotFound, gin.H{"error": "课程不存在"})
					return
				}

				// 检查用户是否已加入课程
				var userCourse model.UserCourse
				result = database.DB.Where("user_id = ? AND course_id = ?", userID, *req.CourseID).First(&userCourse)
				if result.Error != nil {
					c.JSON(http.StatusForbidden, gin.H{"error": "您需要先加入课程才能将笔记关联到该课程"})
					return
				}

				courseIDPtr := *req.CourseID
				note.CourseID = &courseIDPtr
			}
		}
	}

	result = database.DB.Save(&note)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新笔记失败"})
		return
	}

	// 查询用户信息
	var user model.User
	database.DB.First(&user, userID)

	// 查询课程信息（如果课程ID不为nil）
	var courseName string
	var courseID uint
	if note.CourseID != nil {
		var course model.Course
		database.DB.First(&course, *note.CourseID)
		courseName = course.Name
		courseID = *note.CourseID
	} else {
		courseName = "" // 无所属课程
		courseID = 0
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "更新笔记成功",
		"note": NoteResponse{
			ID:            note.ID,
			Title:         note.Title,
			Content:       note.Content,
			Description:   note.Description,
			FilePath:      note.FilePath,
			FileSize:      note.FileSize,
			FileType:      note.FileType,
			ViewCount:     note.ViewCount,
			DownloadCount: note.DownloadCount,
			LikeCount:     note.LikeCount,
			Status:        note.Status,
			UserID:        note.UserID,
			CourseID:      courseID,
			Username:      user.Username,
			CourseName:    courseName,
			CreatedAt:     note.CreatedAt,
			UpdatedAt:     note.UpdatedAt,
		},
	})
}

// DeleteNote 删除笔记
func DeleteNote(c *gin.Context) {
	// 获取用户ID和角色
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	role, _ := c.Get("role")

	// 获取笔记ID
	id := c.Param("id")

	// 查询笔记
	var note model.Note
	result := database.DB.First(&note, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "笔记不存在"})
		return
	}

	// 检查权限
	if note.UserID != userID.(uint) && role.(string) != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "您没有权限删除此笔记"})
		return
	}

	// 删除文件（如果存在）
	if note.FilePath != "" {
		_ = os.Remove(note.FilePath) // 忽略错误
	}

	// 删除笔记
	result = database.DB.Delete(&note)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除笔记失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "删除笔记成功",
	})
}

// DownloadNote 下载笔记
func DownloadNote(c *gin.Context) {
	// 获取笔记ID
	id := c.Param("id")

	// 查询笔记
	var note model.Note
	result := database.DB.First(&note, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "笔记不存在"})
		return
	}

	// 检查权限
	if note.Status == "private" {
		userID, exists := c.Get("userID")
		if !exists || userID.(uint) != note.UserID {
			c.JSON(http.StatusForbidden, gin.H{"error": "您没有权限下载此笔记"})
			return
		}
	}

	// 检查文件是否存在
	if note.FilePath == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "此笔记没有关联文件"})
		return
	}

	// 检查文件是否存在
	if _, err := os.Stat(note.FilePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "文件不存在"})
		return
	}

	// 增加下载次数
	note.DownloadCount++
	database.DB.Save(&note)

	// 获取文件名
	fileName := filepath.Base(note.FilePath)

	// 设置响应头
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
	c.Header("Content-Type", "application/octet-stream")
	c.File(note.FilePath)
}

// LikeNote 点赞笔记
func LikeNote(c *gin.Context) {
	// 获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	// 获取笔记ID
	id := c.Param("id")

	// 查询笔记是否存在
	var note model.Note
	result := database.DB.First(&note, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "笔记不存在"})
		return
	}

	// 检查是否已经点赞
	var existingLike model.NoteLike
	result = database.DB.Where("note_id = ? AND user_id = ?", note.ID, userID).First(&existingLike)
	if result.RowsAffected > 0 {
		// 如果已经点赞过，返回409状态码表示冲突
		c.JSON(http.StatusConflict, gin.H{
			"error":      "您已经点赞过此笔记",
			"like_count": note.LikeCount,
			"is_liked":   true,
		})
		return
	}

	// 创建点赞记录
	like := model.NoteLike{
		NoteID: note.ID,
		UserID: userID.(uint),
	}

	result = database.DB.Create(&like)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "点赞失败"})
		return
	}

	// 更新笔记点赞数
	note.LikeCount++
	database.DB.Save(&note)

	c.JSON(http.StatusOK, gin.H{
		"message":    "点赞成功",
		"like_count": note.LikeCount,
		"is_liked":   true,
	})
}

// UnlikeNote 取消点赞笔记
func UnlikeNote(c *gin.Context) {
	// 获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	// 获取笔记ID
	id := c.Param("id")

	// 查询笔记是否存在
	var note model.Note
	result := database.DB.First(&note, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "笔记不存在"})
		return
	}

	// 查询点赞记录
	var like model.NoteLike
	result = database.DB.Where("note_id = ? AND user_id = ?", note.ID, userID).First(&like)
	if result.RowsAffected == 0 {
		// 如果还没有点赞过，返回404状态码
		c.JSON(http.StatusNotFound, gin.H{
			"error":      "您还没有点赞此笔记",
			"like_count": note.LikeCount,
			"is_liked":   false,
		})
		return
	}

	// 删除点赞记录
	result = database.DB.Delete(&like)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "取消点赞失败"})
		return
	}

	// 更新笔记点赞数
	if note.LikeCount > 0 {
		note.LikeCount--
		database.DB.Save(&note)
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "取消点赞成功",
		"like_count": note.LikeCount,
		"is_liked":   false,
	})
}

// FavoriteNote 收藏笔记
func FavoriteNote(c *gin.Context) {
	// 获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	// 获取笔记ID
	id := c.Param("id")
	noteID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的笔记ID"})
		return
	}

	// 查询笔记
	var note model.Note
	result := database.DB.First(&note, noteID)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "笔记不存在"})
		return
	}

	// 检查是否已收藏
	var favorite model.Favorite
	result = database.DB.Where("user_id = ? AND note_id = ?", userID, noteID).First(&favorite)
	if result.RowsAffected > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "已收藏该笔记"})
		return
	}

	// 添加收藏
	favorite = model.Favorite{
		UserID: userID.(uint),
		NoteID: uint(noteID),
	}

	result = database.DB.Create(&favorite)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "收藏失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "收藏成功",
		"status":  "favorited",
	})
}

// UnfavoriteNote 取消收藏笔记
func UnfavoriteNote(c *gin.Context) {
	// 获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	// 获取笔记ID
	id := c.Param("id")
	noteID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的笔记ID"})
		return
	}

	// 查询笔记
	var note model.Note
	result := database.DB.First(&note, noteID)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "笔记不存在"})
		return
	}

	// 查找收藏记录
	var favorite model.Favorite
	result = database.DB.Where("user_id = ? AND note_id = ?", userID, noteID).First(&favorite)
	if result.RowsAffected == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "未收藏该笔记"})
		return
	}

	// 删除收藏记录
	result = database.DB.Delete(&favorite)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "取消收藏失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "取消收藏成功",
		"status":  "unfavorited",
	})
}

// GetMyFavorites 获取我的收藏笔记
func GetMyFavorites(c *gin.Context) {
	// 获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	// 分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit

	// 查询用户收藏的笔记
	var favorites []model.Favorite
	result := database.DB.Where("user_id = ?", userID).Find(&favorites)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询收藏失败"})
		return
	}

	// 获取笔记ID列表
	noteIDs := make([]uint, len(favorites))
	favoriteMap := make(map[uint]model.Favorite)
	for i, favorite := range favorites {
		noteIDs[i] = favorite.NoteID
		favoriteMap[favorite.NoteID] = favorite
	}

	if len(noteIDs) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"data":  []interface{}{},
			"total": 0,
			"page":  page,
			"limit": limit,
		})
		return
	}

	// 查询笔记详情
	var notes []model.Note
	query := database.DB.Where("id IN ?", noteIDs).Preload("User").Preload("Course")
	result = query.Limit(limit).Offset(offset).Order("created_at DESC").Find(&notes)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询笔记失败"})
		return
	}

	// 构建响应数据
	var responseNotes []gin.H
	for _, note := range notes {
		favorite := favoriteMap[note.ID]

		// 计算评论数量
		var commentCount int64
		database.DB.Model(&model.Comment{}).Where("note_id = ?", note.ID).Count(&commentCount)

		responseNote := gin.H{
			"id":            note.ID,
			"title":         note.Title,
			"description":   note.Description,
			"content":       note.Content,
			"course_id":     note.CourseID,
			"course_name":   note.Course.Name,
			"user_id":       note.UserID,
			"username":      note.User.Username,
			"author_avatar": note.User.Avatar,
			"view_count":    note.ViewCount,
			"like_count":    note.LikeCount,
			"comment_count": int(commentCount),
			"created_at":    note.CreatedAt,
			"updated_at":    note.UpdatedAt,
			"favorited_at":  favorite.CreatedAt,
		}
		responseNotes = append(responseNotes, responseNote)
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  responseNotes,
		"total": len(favorites),
		"page":  page,
		"limit": limit,
	})
}

// GetMyLikes 获取我的点赞笔记
func GetMyLikes(c *gin.Context) {
	// 获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	// 分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit

	// 查询用户点赞的笔记
	var likes []model.NoteLike
	result := database.DB.Where("user_id = ?", userID).Find(&likes)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询点赞失败"})
		return
	}

	// 获取笔记ID列表
	noteIDs := make([]uint, len(likes))
	likeMap := make(map[uint]model.NoteLike)
	for i, like := range likes {
		noteIDs[i] = like.NoteID
		likeMap[like.NoteID] = like
	}

	if len(noteIDs) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"data":  []interface{}{},
			"total": 0,
			"page":  page,
			"limit": limit,
		})
		return
	}

	// 查询笔记详情
	var notes []model.Note
	query := database.DB.Where("id IN ?", noteIDs).Preload("User").Preload("Course")
	result = query.Limit(limit).Offset(offset).Order("created_at DESC").Find(&notes)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询笔记失败"})
		return
	}

	// 构建响应数据
	var responseNotes []gin.H
	for _, note := range notes {
		like := likeMap[note.ID]

		// 计算评论数量
		var commentCount int64
		database.DB.Model(&model.Comment{}).Where("note_id = ?", note.ID).Count(&commentCount)

		responseNote := gin.H{
			"id":            note.ID,
			"title":         note.Title,
			"description":   note.Description,
			"content":       note.Content,
			"course_id":     note.CourseID,
			"course_name":   note.Course.Name,
			"user_id":       note.UserID,
			"username":      note.User.Username,
			"author_avatar": note.User.Avatar,
			"view_count":    note.ViewCount,
			"like_count":    note.LikeCount,
			"comment_count": int(commentCount),
			"created_at":    note.CreatedAt,
			"updated_at":    note.UpdatedAt,
			"liked_at":      like.CreatedAt,
		}
		responseNotes = append(responseNotes, responseNote)
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  responseNotes,
		"total": len(likes),
		"page":  page,
		"limit": limit,
	})
}

// SearchNotes 搜索笔记
func SearchNotes(c *gin.Context) {
	// 获取搜索关键词
	keyword := c.Query("keyword")
	if keyword == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "搜索关键词不能为空"})
		return
	}

	// 分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	// 构建查询
	query := database.DB.Model(&model.Note{})

	// 应用搜索条件
	searchTerm := "%" + keyword + "%"
	query = query.Where(
		"title LIKE ? OR content LIKE ? OR description LIKE ?",
		searchTerm, searchTerm, searchTerm,
	)

	// 只查询公开笔记
	query = query.Where("status = ?", "public")

	// 获取总数
	var total int64
	query.Count(&total)

	// 排序参数
	sortBy := c.DefaultQuery("sort_by", "created_at")
	order := c.DefaultQuery("order", "desc")

	var orderClause string
	switch sortBy {
	case "like_count":
		orderClause = "like_count " + strings.ToUpper(order)
	case "view_count":
		orderClause = "view_count " + strings.ToUpper(order)
	case "created_at":
		orderClause = "created_at " + strings.ToUpper(order)
	default:
		orderClause = "created_at DESC"
	}

	// 获取分页数据，并预加载用户和课程信息
	var notes []model.Note
	result := query.Preload("User").Preload("Course").Limit(pageSize).Offset(offset).Order(orderClause).Find(&notes)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "搜索笔记失败"})
		return
	}

	// 转换为响应格式
	var noteResponses []NoteResponse
	for _, note := range notes {
		// 计算评论数量
		var commentCount int64
		database.DB.Model(&model.Comment{}).Where("note_id = ?", note.ID).Count(&commentCount)

		// 处理CourseID和CourseName
		var courseID uint
		var courseName string
		if note.CourseID != nil {
			courseID = *note.CourseID
			courseName = note.Course.Name
		} else {
			courseID = 0
			courseName = ""
		}

		noteResponses = append(noteResponses, NoteResponse{
			ID:            note.ID,
			Title:         note.Title,
			Description:   note.Description,
			FileSize:      note.FileSize,
			FileType:      note.FileType,
			ViewCount:     note.ViewCount,
			DownloadCount: note.DownloadCount,
			LikeCount:     note.LikeCount,
			CommentCount:  int(commentCount),
			Status:        note.Status,
			UserID:        note.UserID,
			CourseID:      courseID,
			Username:      note.User.Username,
			AuthorAvatar:  note.User.Avatar,
			CourseName:    courseName,
			CreatedAt:     note.CreatedAt,
			UpdatedAt:     note.UpdatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"total":       total,
		"page":        page,
		"page_size":   pageSize,
		"total_pages": (int(total) + pageSize - 1) / pageSize,
		"notes":       noteResponses,
	})
}

// GetPopularNotes 获取热门笔记（点赞量最高的前3个）
func GetPopularNotes(c *gin.Context) {
	// 构建查询，只获取公开笔记
	query := database.DB.Model(&model.Note{}).Where("status = ?", "public")

	// 获取前3个点赞最多的笔记
	var notes []model.Note
	result := query.Preload("User").Preload("Course").Order("like_count DESC").Limit(3).Find(&notes)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取热门笔记失败"})
		return
	}

	// 转换为响应格式
	var noteResponses []NoteResponse
	for _, note := range notes {
		// 计算评论数量
		var commentCount int64
		database.DB.Model(&model.Comment{}).Where("note_id = ?", note.ID).Count(&commentCount)

		// 处理CourseID和CourseName
		var courseID uint
		var courseName string
		if note.CourseID != nil {
			courseID = *note.CourseID
			if note.Course.ID != 0 {
				courseName = note.Course.Name
			}
		}

		noteResponses = append(noteResponses, NoteResponse{
			ID:           note.ID,
			Title:        note.Title,
			Content:      note.Content,
			Description:  note.Description,
			CourseID:     courseID,
			CourseName:   courseName,
			UserID:       note.UserID,
			Username:     note.User.Username,
			Status:       note.Status,
			ViewCount:    note.ViewCount,
			LikeCount:    note.LikeCount,
			CommentCount: int(commentCount),
			CreatedAt:    note.CreatedAt,
			UpdatedAt:    note.UpdatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"notes": noteResponses,
	})
}

// getUsernameFromNote 安全地获取用户名
func getUsernameFromNote(note model.Note) string {
	if note.User.ID != 0 && note.User.Username != "" {
		return note.User.Username
	}
	return ""
}

// getAvatarFromNote 安全地获取用户头像
func getAvatarFromNote(note model.Note) string {
	if note.User.ID != 0 && note.User.Avatar != "" {
		return note.User.Avatar
	}
	return ""
}
