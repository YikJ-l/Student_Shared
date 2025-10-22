package api

import (
	"net/http"
	"student_shared/app/model"
	"student_shared/app/utils/database"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	// 引入统一的请求/响应包
	req "student_shared/app/model/req"
	resp "student_shared/app/model/resp"
)

// ListNoteComments 获取笔记评论列表（POST + JSON）
func ListNoteComments(c *gin.Context) {
	// 使用统一的请求结构体
	var params req.ListCommentsRequest
	params.Page = 1
	params.PageSize = 20
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数无效"})
		return
	}

	// 获取当前用户ID（可选）
	currentUserID, _ := c.Get("userID")

	// 分页参数
	page := params.Page
	pageSize := params.PageSize
	if page < 1 { page = 1 }
	if pageSize < 1 || pageSize > 100 { pageSize = 20 }

	// 查询顶级评论总数
	total, err := model.CountTopLevelComments(database.DB, params.NoteID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取评论数量失败"})
		return
	}

	// 查询顶级评论列表
	comments, err := model.FindTopLevelCommentsWithReplies(database.DB, params.NoteID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取评论列表失败"})
		return
	}

	// 转换为响应格式
	var commentResponses []resp.CommentResponse
	for _, comment := range comments {
		commentResponse := buildCommentResponse(database.DB, comment, currentUserID)
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
func buildCommentResponse(db *gorm.DB, comment model.Comment, currentUserID interface{}) resp.CommentResponse {
	// 检查当前用户是否点赞了该评论
	isLiked := false
	if currentUserID != nil {
		liked, _ := model.IsCommentLikedByUser(db, comment.ID, currentUserID.(uint))
		isLiked = liked
	}

	// 构建回复列表
	var replies []resp.CommentResponse
	for _, reply := range comment.Replies {
		replyResponse := resp.CommentResponse{
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
			liked, _ := model.IsCommentLikedByUser(db, reply.ID, currentUserID.(uint))
			replyResponse.IsLiked = liked
		}

		replies = append(replies, replyResponse)
	}

	return resp.CommentResponse{
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
	var params req.CommentRequest
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数无效"})
		return
	}

	// 检查笔记是否存在
	note, err := model.GetNoteByID(database.DB, params.NoteID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "笔记不存在"})
		return
	}

	// 验证评分范围
	if params.Rating < 0 || params.Rating > 5 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "评分必须在0-5之间"})
		return
	}

	// 如果是回复评论，验证父评论是否存在
	if params.ParentID != nil {
		parentComment, err := model.GetCommentByID(database.DB, *params.ParentID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "父评论不存在"})
			return
		}
		// 确保父评论属于同一笔记
		if parentComment.NoteID != params.NoteID {
			c.JSON(http.StatusBadRequest, gin.H{"error": "父评论不属于该笔记"})
			return
		}
	}

	// 如果指定了回复用户，验证用户是否存在
	if params.ReplyToUserID != nil {
		_, err := model.GetUserByID(database.DB, *params.ReplyToUserID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "回复的用户不存在"})
			return
		}
	}

	// 创建评论
	comment := model.Comment{
		Content:       params.Content,
		Rating:        params.Rating,
		UserID:        userID.(uint),
		NoteID:        note.ID,
		ParentID:      params.ParentID,
		ReplyToUserID: params.ReplyToUserID,
	}
	if err := model.CreateComment(database.DB, &comment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建评论失败"})
		return
	}

	// 重新查询评论以获取关联数据
	comment, _ = model.GetCommentWithPreloads(database.DB, comment.ID)

	// 构建响应
	response := buildCommentResponse(database.DB, comment, userID)

	c.JSON(http.StatusCreated, gin.H{
		"message": "评论成功",
		"comment": response,
	})
}

// UpdateComment 更新评论（POST + JSON）
func UpdateComment(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}
	var params req.CommentUpdateRequest
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数无效"})
		return
	}
	comment, err := model.GetCommentByID(database.DB, params.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "评论不存在"})
		return
	}
	if comment.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "您没有权限更新此评论"})
		return
	}
	if params.Rating < 0 || params.Rating > 5 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "评分必须在0-5之间"})
		return
	}
	comment.Content = params.Content
	comment.Rating = params.Rating
	if err := model.UpdateComment(database.DB, &comment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新评论失败"})
		return
	}
	comment, _ = model.GetCommentWithPreloads(database.DB, comment.ID)
	response := buildCommentResponse(database.DB, comment, userID)
	c.JSON(http.StatusOK, gin.H{"message": "更新评论成功", "comment": response})
}

// DeleteComment 删除评论（POST + JSON）
func DeleteComment(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}
	role, _ := c.Get("role")
	var params req.GetCommentRequest
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数无效"})
		return
	}
	comment, err := model.GetCommentByID(database.DB, params.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "评论不存在"})
		return
	}
	if comment.UserID != userID.(uint) && role.(string) != "admin" {
		note, err := model.GetNoteByID(database.DB, comment.NoteID)
		if err != nil || note.UserID != userID.(uint) {
			c.JSON(http.StatusForbidden, gin.H{"error": "您没有权限删除此评论"})
			return
		}
	}
	if err := model.DeleteCommentWithCascade(database.DB, comment.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除评论失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "删除评论成功"})
}

// LikeComment 点赞评论（POST + JSON）
func LikeComment(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}
	var params req.GetCommentRequest
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数无效"})
		return
	}
	_, err := model.GetCommentByID(database.DB, params.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "评论不存在"})
		return
	}
	newLikeCount, alreadyLiked, err := model.LikeCommentTx(database.DB, params.ID, userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "点赞失败"})
		return
	}
	if alreadyLiked {
		c.JSON(http.StatusConflict, gin.H{"error": "您已经点赞过此评论", "like_count": newLikeCount, "is_liked": true})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "点赞成功", "data": gin.H{"like_count": newLikeCount, "is_liked": true}})
}

// UnlikeComment 取消点赞评论（POST + JSON）
func UnlikeComment(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}
	var params req.GetCommentRequest
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数无效"})
		return
	}
	_, err := model.GetCommentByID(database.DB, params.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "评论不存在"})
		return
	}
	newLikeCount, hadLiked, err := model.UnlikeCommentTx(database.DB, params.ID, userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "取消点赞失败"})
		return
	}
	if !hadLiked {
		c.JSON(http.StatusNotFound, gin.H{"message": "您还没有点赞此评论", "like_count": newLikeCount, "is_liked": false})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "取消点赞成功", "data": gin.H{"like_count": newLikeCount, "is_liked": false}})
}
