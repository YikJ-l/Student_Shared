package api

import (
	"net/http"
	"strconv"
	"strings"
	"student_shared/app/model"
	"student_shared/app/utils/database"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CommentRequest 评论请求
type CommentRequest struct {
	Content       string `json:"content" binding:"required"`
	NoteID        uint   `json:"note_id" binding:"required"`
	Rating        int    `json:"rating"`           // 0-5星评分
	ParentID      *uint  `json:"parent_id"`        // 父评论ID，用于回复
	ReplyToUserID *uint  `json:"reply_to_user_id"` // 回复的用户ID
}

// CommentResponse 评论响应
type CommentResponse struct {
	ID              uint              `json:"id"`
	Content         string            `json:"content"`
	Rating          int               `json:"rating"`
	LikeCount       int               `json:"like_count"`
	UserID          uint              `json:"user_id"`
	NoteID          uint              `json:"note_id"`
	Username        string            `json:"username"`
	UserAvatar      string            `json:"user_avatar"`
	ParentID        *uint             `json:"parent_id"`
	ReplyToUserID   *uint             `json:"reply_to_user_id"`
	ReplyToUsername *string           `json:"reply_to_username"`
	Replies         []CommentResponse `json:"replies"`
	IsLiked         bool              `json:"is_liked"` // 当前用户是否已点赞
	CreatedAt       time.Time         `json:"created_at"`
	UpdatedAt       time.Time         `json:"updated_at"`
}

// ListNoteComments 获取笔记评论列表
func ListNoteComments(c *gin.Context) {
	// 获取笔记ID
	noteID := c.Param("noteId")

	// 获取当前用户ID（可选）
	currentUserID, _ := c.Get("userID")

	// 分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	// 查询顶级评论总数
	var total int64
	database.DB.Model(&model.Comment{}).Where("note_id = ? AND parent_id IS NULL", noteID).Count(&total)

	// 查询顶级评论列表
	var comments []model.Comment
	result := database.DB.Where("note_id = ? AND parent_id IS NULL", noteID).
		Preload("User").
		Preload("Replies", func(db *gorm.DB) *gorm.DB {
			return db.Preload("User").Preload("ReplyToUser").Order("created_at ASC")
		}).
		Preload("Replies.User").
		Preload("Replies.ReplyToUser").
		Order("created_at DESC").
		Limit(pageSize).Offset(offset).Find(&comments)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取评论列表失败"})
		return
	}

	// 转换为响应格式
	var commentResponses []CommentResponse
	for _, comment := range comments {
		commentResponse := buildCommentResponse(comment, currentUserID)
		commentResponses = append(commentResponses, commentResponse)
	}

	c.JSON(http.StatusOK, gin.H{
		"total":       total,
		"page":        page,
		"page_size":   pageSize,
		"total_pages": (int(total) + pageSize - 1) / pageSize,
		"comments":    commentResponses,
	})
}

// buildCommentResponse 构建评论响应
func buildCommentResponse(comment model.Comment, currentUserID interface{}) CommentResponse {
	// 检查当前用户是否点赞了该评论
	isLiked := false
	if currentUserID != nil {
		var likeCount int64
		database.DB.Model(&model.CommentLike{}).Where("comment_id = ? AND user_id = ? AND is_active = ?", comment.ID, currentUserID, true).Count(&likeCount)
		isLiked = likeCount > 0
	}

	// 构建回复列表
	var replies []CommentResponse
	for _, reply := range comment.Replies {
		replyResponse := CommentResponse{
			ID:            reply.ID,
			Content:       reply.Content,
			Rating:        reply.Rating,
			LikeCount:     reply.LikeCount,
			UserID:        reply.UserID,
			NoteID:        reply.NoteID,
			Username:      reply.User.Username,
			UserAvatar:    reply.User.Avatar,
			ParentID:      reply.ParentID,
			ReplyToUserID: reply.ReplyToUserID,
			CreatedAt:     reply.CreatedAt,
			UpdatedAt:     reply.UpdatedAt,
		}

		// 设置回复目标用户名
		if reply.ReplyToUser != nil {
			replyResponse.ReplyToUsername = &reply.ReplyToUser.Username
		}

		// 检查当前用户是否点赞了该回复
		if currentUserID != nil {
			var replyLikeCount int64
			database.DB.Model(&model.CommentLike{}).Where("comment_id = ? AND user_id = ? AND is_active = ?", reply.ID, currentUserID, true).Count(&replyLikeCount)
			replyResponse.IsLiked = replyLikeCount > 0
		}

		replies = append(replies, replyResponse)
	}

	return CommentResponse{
		ID:            comment.ID,
		Content:       comment.Content,
		Rating:        comment.Rating,
		LikeCount:     comment.LikeCount,
		UserID:        comment.UserID,
		NoteID:        comment.NoteID,
		Username:      comment.User.Username,
		UserAvatar:    comment.User.Avatar,
		ParentID:      comment.ParentID,
		ReplyToUserID: comment.ReplyToUserID,
		Replies:       replies,
		IsLiked:       isLiked,
		CreatedAt:     comment.CreatedAt,
		UpdatedAt:     comment.UpdatedAt,
	}
}

// CreateComment 创建评论
func CreateComment(c *gin.Context) {
	// 获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	// 绑定请求数据
	var req CommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数无效"})
		return
	}

	// 检查笔记是否存在
	var note model.Note
	result := database.DB.First(&note, req.NoteID)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "笔记不存在"})
		return
	}

	// 验证评分范围
	if req.Rating < 0 || req.Rating > 5 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "评分必须在0-5之间"})
		return
	}

	// 如果是回复评论，验证父评论是否存在
	if req.ParentID != nil {
		var parentComment model.Comment
		result := database.DB.First(&parentComment, *req.ParentID)
		if result.Error != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "父评论不存在"})
			return
		}
		// 确保父评论属于同一笔记
		if parentComment.NoteID != req.NoteID {
			c.JSON(http.StatusBadRequest, gin.H{"error": "父评论不属于该笔记"})
			return
		}
	}

	// 如果指定了回复用户，验证用户是否存在
	if req.ReplyToUserID != nil {
		var replyToUser model.User
		result := database.DB.First(&replyToUser, *req.ReplyToUserID)
		if result.Error != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "回复的用户不存在"})
			return
		}
	}

	// 创建评论
	comment := model.Comment{
		Content:       req.Content,
		Rating:        req.Rating,
		UserID:        userID.(uint),
		NoteID:        req.NoteID,
		ParentID:      req.ParentID,
		ReplyToUserID: req.ReplyToUserID,
	}

	result = database.DB.Create(&comment)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建评论失败"})
		return
	}

	// 重新查询评论以获取关联数据
	database.DB.Preload("User").Preload("ReplyToUser").First(&comment, comment.ID)

	// 构建响应
	response := buildCommentResponse(comment, userID)

	c.JSON(http.StatusCreated, gin.H{
		"message": "评论成功",
		"comment": response,
	})
}

// UpdateComment 更新评论
func UpdateComment(c *gin.Context) {
	// 获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	// 获取评论ID
	id := c.Param("id")

	// 查询评论
	var comment model.Comment
	result := database.DB.First(&comment, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "评论不存在"})
		return
	}

	// 检查权限
	if comment.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "您没有权限更新此评论"})
		return
	}

	// 绑定请求数据
	var req struct {
		Content string `json:"content" binding:"required"`
		Rating  int    `json:"rating"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数无效"})
		return
	}

	// 验证评分范围
	if req.Rating < 0 || req.Rating > 5 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "评分必须在0-5之间"})
		return
	}

	// 更新评论
	comment.Content = req.Content
	comment.Rating = req.Rating

	result = database.DB.Save(&comment)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新评论失败"})
		return
	}

	// 重新查询评论以获取关联数据
	database.DB.Preload("User").Preload("ReplyToUser").First(&comment, comment.ID)

	// 构建响应
	response := buildCommentResponse(comment, userID)

	c.JSON(http.StatusOK, gin.H{
		"message": "更新评论成功",
		"comment": response,
	})
}

// DeleteComment 删除评论
func DeleteComment(c *gin.Context) {
	// 获取用户ID和角色
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	role, _ := c.Get("role")

	// 获取评论ID
	id := c.Param("id")

	// 查询评论
	var comment model.Comment
	result := database.DB.First(&comment, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "评论不存在"})
		return
	}

	// 检查权限
	if comment.UserID != userID.(uint) && role.(string) != "admin" {
		// 检查是否是笔记作者
		var note model.Note
		result = database.DB.First(&note, comment.NoteID)
		if result.Error != nil || note.UserID != userID.(uint) {
			c.JSON(http.StatusForbidden, gin.H{"error": "您没有权限删除此评论"})
			return
		}
	}

	// 递归删除子评论
	deleteCommentAndChildren(comment.ID)

	// 删除当前评论的点赞记录
	database.DB.Where("comment_id = ?", comment.ID).Delete(&model.CommentLike{})

	// 删除当前评论
	result = database.DB.Delete(&comment)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除评论失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "删除评论成功",
	})
}

// LikeComment 点赞评论
func LikeComment(c *gin.Context) {
	// 获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	// 获取评论ID
	id := c.Param("id")

	// 查询评论是否存在
	var comment model.Comment
	result := database.DB.First(&comment, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "评论不存在"})
		return
	}

	// 使用事务确保数据一致性
	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 在事务内重新查询评论以获取最新数据
	var commentInTx model.Comment
	if err := tx.First(&commentInTx, comment.ID).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{"error": "评论不存在"})
		return
	}

	// 检查是否已经存在点赞记录
	var existingLike model.CommentLike
	result = tx.Where("comment_id = ? AND user_id = ?", comment.ID, userID).First(&existingLike)
	
	if result.RowsAffected > 0 {
		// 如果记录存在且已激活，说明已经点赞过
		if existingLike.IsActive {
			tx.Rollback()
			c.JSON(http.StatusConflict, gin.H{
				"error": "您已经点赞过此评论",
				"like_count": commentInTx.LikeCount,
				"is_liked": true,
			})
			return
		}
		// 如果记录存在但未激活，重新激活点赞
		if err := tx.Model(&existingLike).Update("is_active", true).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "点赞失败"})
			return
		}
	} else {
		// 创建新的点赞记录
		like := model.CommentLike{
			CommentID: comment.ID,
			UserID:    userID.(uint),
			IsActive:  true,
		}
		
		if err := tx.Create(&like).Error; err != nil {
			tx.Rollback()
			// 检查是否是唯一约束冲突（重复点赞）
			if strings.Contains(err.Error(), "Duplicate entry") || strings.Contains(err.Error(), "uk_comment_likes_user_comment") {
				c.JSON(http.StatusConflict, gin.H{
					"error": "您已经点赞过此评论",
					"like_count": commentInTx.LikeCount,
					"is_liked": true,
				})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "点赞失败"})
			return
		}
	}

	// 更新评论点赞数
	newLikeCount := commentInTx.LikeCount + 1
	if err := tx.Model(&commentInTx).Update("like_count", newLikeCount).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新点赞数失败"})
		return
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "提交事务失败"})
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "点赞成功",
		"data": gin.H{
			"like_count": newLikeCount,
			"is_liked": true,
		},
	})
}

// UnlikeComment 取消点赞评论
func UnlikeComment(c *gin.Context) {
	// 获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	// 获取评论ID
	id := c.Param("id")

	// 查询评论是否存在
	var comment model.Comment
	result := database.DB.First(&comment, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "评论不存在"})
		return
	}

	// 使用事务确保数据一致性
	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 在事务内重新查询评论以获取最新数据
	var commentInTx model.Comment
	if err := tx.First(&commentInTx, comment.ID).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{"error": "评论不存在"})
		return
	}

	// 查询点赞记录
	var like model.CommentLike
	result = tx.Where("comment_id = ? AND user_id = ?", comment.ID, userID).First(&like)
	if result.RowsAffected == 0 || !like.IsActive {
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{
			"message": "您还没有点赞此评论",
			"like_count": commentInTx.LikeCount,
			"is_liked": false,
		})
		return
	}

	// 将点赞记录标记为非激活状态
	if err := tx.Model(&like).Update("is_active", false).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "取消点赞失败"})
		return
	}

	// 更新评论点赞数（确保不会小于0）
	newLikeCount := commentInTx.LikeCount
	if newLikeCount > 0 {
		newLikeCount--
	}
	if err := tx.Model(&commentInTx).Update("like_count", newLikeCount).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新点赞数失败"})
		return
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "提交事务失败"})
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "取消点赞成功",
		"data": gin.H{
			"like_count": newLikeCount,
			"is_liked": false,
		},
	})
}

// deleteCommentAndChildren 递归删除评论及其所有子评论
func deleteCommentAndChildren(commentID uint) {
	// 查找所有子评论
	var childComments []model.Comment
	database.DB.Where("parent_id = ?", commentID).Find(&childComments)

	// 递归删除每个子评论
	for _, child := range childComments {
		// 先递归删除子评论的子评论
		deleteCommentAndChildren(child.ID)
		
		// 删除子评论的点赞记录
		database.DB.Where("comment_id = ?", child.ID).Delete(&model.CommentLike{})
		
		// 删除子评论
		database.DB.Delete(&child)
	}
}
