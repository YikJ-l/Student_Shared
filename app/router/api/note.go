package api

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"student_shared/app/model"
	"student_shared/app/utils/database"

	"student_shared/app/model/req"
	"student_shared/app/model/resp"

	"github.com/gin-gonic/gin"
)

func ListNotes(c *gin.Context) {
	// 读取 JSON 请求体
	var params req.ListNotesRequest
	params.Status = "public"
	params.Page = 1
	params.PageSize = 10
	params.SortBy = "created_at"
	params.Order = "desc"
	_ = c.ShouldBindJSON(&params)

	// 获取当前用户ID（如果已登录），用于私有笔记可见性判断
	currentUserID, exists := c.Get("userID")
	canViewPrivate := false
	if exists {
		if params.UserID != nil && currentUserID.(uint) == *params.UserID {
			canViewPrivate = true
		}
		// 可在此添加管理员权限判断：如果是管理员，也允许查看私有
	}

	// 调用 model 层进行查询
	filter := model.NoteListFilter{CourseID: params.CourseID, UserID: params.UserID, Status: params.Status, CanViewPrivate: canViewPrivate}
	notes, total, err := model.ListNotes(database.DB, filter, params.Page, params.PageSize, params.SortBy, params.Order)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取笔记列表失败"})
		return
	}

	// 构建响应
	var noteResponses []resp.NoteResponse
	for _, note := range notes {
		commentCount, _ := model.CountCommentsForNote(database.DB, note.ID)
		var courseID uint
		var courseName string
		if note.CourseID != nil {
			courseID = *note.CourseID
			courseName = note.Course.Name
		}
		noteResponses = append(noteResponses, resp.NoteResponse{
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

	page := params.Page
	pageSize := params.PageSize
	c.JSON(http.StatusOK, gin.H{
		"total":       total,
		"page":        page,
		"page_size":   pageSize,
		"total_pages": (int(total) + pageSize - 1) / pageSize,
		"notes":       noteResponses,
	})
}

func buildNoteResponse(note model.Note) resp.NoteResponse {
	commentCount, _ := model.CountCommentsForNote(database.DB, note.ID)
	var courseID uint
	var courseName string
	if note.CourseID != nil {
		courseID = *note.CourseID
		if note.Course.ID != 0 {
			courseName = note.Course.Name
		}
	}
	return resp.NoteResponse{
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
		IsLiked:       false,
		IsFavorited:   false,
		CreatedAt:     note.CreatedAt,
		UpdatedAt:     note.UpdatedAt,
	}
}

func GetNote(c *gin.Context) {
	var params req.GetNoteRequest
	if err := c.ShouldBindJSON(&params); err != nil || params.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的笔记ID"})
		return
	}
	// 获取笔记（预加载）
	note, err := model.GetNoteWithPreloads(database.DB, params.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "笔记不存在"})
		return
	}
	// 权限检查
	if note.Status == "private" {
		userID, exists := c.Get("userID")
		if !exists || userID.(uint) != note.UserID {
			c.JSON(http.StatusForbidden, gin.H{"error": "您没有权限查看此笔记"})
			return
		}
	}
	// 增加查看次数
	_ = model.IncrementViewCount(database.DB, note.ID)
	// 检查点赞与收藏状态
	var isLiked, isFavorited bool
	if currentUserID, exists := c.Get("userID"); exists {
		isLiked, _ = model.IsNoteLikedByUser(database.DB, note.ID, currentUserID.(uint))
		isFavorited, _ = model.IsNoteFavoritedByUser(database.DB, note.ID, currentUserID.(uint))
	}
	// 课程字段
	var courseID uint
	var courseName string
	if note.CourseID != nil {
		courseID = *note.CourseID
		if note.Course.ID != 0 {
			courseName = note.Course.Name
		}
	}
	c.JSON(http.StatusOK, resp.NoteResponse{
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
		Username:      note.User.Username,
		AuthorAvatar:  note.User.Avatar,
		CourseName:    courseName,
		IsLiked:       isLiked,
		IsFavorited:   isFavorited,
		CreatedAt:     note.CreatedAt,
		UpdatedAt:     note.UpdatedAt,
	})
}

func CreateNote(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}
	type NoteCreatePayload struct {
		Title       string `json:"title" binding:"required"`
		Content     string `json:"content"`
		Description string `json:"description"`
		CourseID    uint   `json:"course_id"`
		Status      string `json:"status"`
		Visibility  string `json:"visibility"`
	}
	var req NoteCreatePayload
	if err := c.ShouldBindJSON(&req); err != nil || strings.TrimSpace(req.Title) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "标题不能为空或请求体无效"})
		return
	}
	title := strings.TrimSpace(req.Title)
	content := req.Content
	description := req.Description
	status := strings.TrimSpace(req.Status)
	if status == "" {
		vis := strings.TrimSpace(req.Visibility)
		if vis != "" {
			status = vis
		} else {
			status = "public"
		}
	}
	if status != "public" && status != "private" {
		status = "public"
	}
	courseID := req.CourseID
	if courseID != 0 {
		if _, err := model.GetCourseByID(database.DB, courseID); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "课程不存在"})
			return
		}
		joined, err := model.IsUserJoinedCourse(database.DB, userID.(uint), courseID)
		if err != nil || !joined {
			c.JSON(http.StatusForbidden, gin.H{"error": "您需要先加入课程才能创建笔记"})
			return
		}
	}
	note := model.Note{Title: title, Content: content, Description: description, Status: status, UserID: userID.(uint)}
	if courseID != 0 {
		note.CourseID = &courseID
	}
	if err := model.CreateNote(database.DB, &note); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建笔记失败"})
		return
	}
	var courseName string
	var cid uint
	if note.CourseID != nil {
		cid = *note.CourseID
		course, _ := model.GetCourseByID(database.DB, cid)
		courseName = course.Name
	}
	c.JSON(http.StatusCreated, resp.NoteResponse{
		ID:           note.ID,
		Title:        note.Title,
		Content:      note.Content,
		Description:  note.Description,
		Status:       note.Status,
		UserID:       note.UserID,
		CourseID:     cid,
		Username:     getUsernameFromNote(note),
		AuthorAvatar: getAvatarFromNote(note),
		CourseName:   courseName,
		CreatedAt:    note.CreatedAt,
		UpdatedAt:    note.UpdatedAt,
	})
}

func UpdateNote(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}
	role, _ := c.Get("role")
	var params req.UpdateNoteRequest
	if err := c.ShouldBindJSON(&params); err != nil || params.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数无效或缺少笔记ID"})
		return
	}
	note, err := model.GetNoteByID(database.DB, params.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "笔记不存在"})
		return
	}
	if note.UserID != userID.(uint) && role.(string) != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "您没有权限更新此笔记"})
		return
	}
	if params.Title != nil {
		note.Title = *params.Title
	}
	if params.Content != nil {
		note.Content = *params.Content
	}
	if params.Description != nil {
		note.Description = *params.Description
	}
	if params.Status != nil {
		note.Status = *params.Status
	}
	if params.CourseID != nil {
		needUpdate := false
		if note.CourseID == nil && *params.CourseID != 0 {
			needUpdate = true
		}
		if note.CourseID != nil && *params.CourseID != *note.CourseID {
			needUpdate = true
		}
		if note.CourseID != nil && *params.CourseID == 0 {
			needUpdate = true
		}
		if needUpdate {
			if *params.CourseID == 0 {
				note.CourseID = nil
			} else {
				if _, err := model.GetCourseByID(database.DB, *params.CourseID); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": "目标课程不存在"})
					return
				}
				note.CourseID = params.CourseID
			}
		}
	}
	if err := model.UpdateNote(database.DB, note); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新笔记失败"})
		return
	}
	user, _ := model.GetUserByID(database.DB, note.UserID)
	var courseName string
	var courseID uint
	if note.CourseID != nil {
		courseID = *note.CourseID
		course, _ := model.GetCourseByID(database.DB, courseID)
		courseName = course.Name
	}
	c.JSON(http.StatusOK, gin.H{"message": "更新笔记成功", "note": resp.NoteResponse{
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
		AuthorAvatar:  user.Avatar,
		CourseName:    courseName,
		CreatedAt:     note.CreatedAt,
		UpdatedAt:     note.UpdatedAt,
	}})
}

func DeleteNote(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}
	role, _ := c.Get("role")
	var req struct {
		ID uint `json:"id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "缺少有效的笔记ID"})
		return
	}
	note, err := model.GetNoteByID(database.DB, req.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "笔记不存在"})
		return
	}
	if note.UserID != userID.(uint) && role.(string) != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "您没有权限删除此笔记"})
		return
	}
	if note.FilePath != "" {
		_ = os.Remove(note.FilePath)
	}
	if err := model.DeleteNoteByID(database.DB, note.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除笔记失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "删除笔记成功"})
}

func DownloadNote(c *gin.Context) {
	var req struct {
		ID uint `json:"id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "缺少有效的笔记ID"})
		return
	}
	note, err := model.GetNoteByID(database.DB, req.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "笔记不存在"})
		return
	}
	if note.Status == "private" {
		userID, exists := c.Get("userID")
		if !exists || userID.(uint) != note.UserID {
			c.JSON(http.StatusForbidden, gin.H{"error": "您没有权限下载此笔记"})
			return
		}
	}
	if note.FilePath == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "此笔记没有关联文件"})
		return
	}
	if _, statErr := os.Stat(note.FilePath); os.IsNotExist(statErr) {
		c.JSON(http.StatusNotFound, gin.H{"error": "文件不存在"})
		return
	}
	_ = model.IncrementDownloadCount(database.DB, note.ID)
	fileName := filepath.Base(note.FilePath)
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
	c.Header("Content-Type", "application/octet-stream")
	c.File(note.FilePath)
}

func LikeNote(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}
	var req struct {
		ID uint `json:"id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "缺少有效的笔记ID"})
		return
	}
	alreadyLiked, likeCount, err := model.LikeNoteTx(database.DB, req.ID, userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "点赞失败"})
		return
	}
	if alreadyLiked {
		c.JSON(http.StatusConflict, gin.H{"error": "您已经点赞过此笔记", "like_count": likeCount, "is_liked": true})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "点赞成功", "like_count": likeCount, "is_liked": true})
}

func UnlikeNote(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}
	var req struct {
		ID uint `json:"id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "缺少有效的笔记ID"})
		return
	}
	notLiked, likeCount, err := model.UnlikeNoteTx(database.DB, req.ID, userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "取消点赞失败"})
		return
	}
	if notLiked {
		c.JSON(http.StatusNotFound, gin.H{"error": "您还没有点赞此笔记", "like_count": likeCount, "is_liked": false})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "取消点赞成功", "like_count": likeCount, "is_liked": false})
}

func FavoriteNote(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}
	var req struct {
		ID uint `json:"id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的笔记ID"})
		return
	}
	alreadyFavorited, err := model.FavoriteNote(database.DB, req.ID, userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "收藏失败"})
		return
	}
	if alreadyFavorited {
		c.JSON(http.StatusBadRequest, gin.H{"error": "已收藏该笔记"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "收藏成功", "status": "favorited"})
}

func UnfavoriteNote(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}
	var req struct {
		ID uint `json:"id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的笔记ID"})
		return
	}
	notFavorited, err := model.UnfavoriteNote(database.DB, req.ID, userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "取消收藏失败"})
		return
	}
	if notFavorited {
		c.JSON(http.StatusBadRequest, gin.H{"error": "未收藏该笔记"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "取消收藏成功", "status": "unfavorited"})
}

func GetMyFavorites(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}
	var req struct {
		Page  int `json:"page"`
		Limit int `json:"limit"`
	}
	req.Page = 1
	req.Limit = 10
	_ = c.ShouldBindJSON(&req)
	if req.Page < 1 {
		req.Page = 1
	}
	if req.Limit < 1 || req.Limit > 100 {
		req.Limit = 10
	}
	favorites, err := model.GetFavoritesByUser(database.DB, userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询收藏失败"})
		return
	}
	notes, total, err := model.GetFavoritesNotesForUser(database.DB, userID.(uint), req.Page, req.Limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询笔记失败"})
		return
	}
	favMap := make(map[uint]model.Favorite)
	for _, f := range favorites {
		favMap[f.NoteID] = f
	}
	var responseNotes []gin.H
	for _, note := range notes {
		commentCount, _ := model.CountCommentsForNote(database.DB, note.ID)
		var courseID uint
		var courseName string
		if note.CourseID != nil {
			courseID = *note.CourseID
			if note.Course.ID != 0 {
				courseName = note.Course.Name
			}
		}
		responseNotes = append(responseNotes, gin.H{
			"id":            note.ID,
			"title":         note.Title,
			"content":       note.Content,
			"description":   note.Description,
			"course_id":     courseID,
			"course_name":   courseName,
			"user_id":       note.UserID,
			"username":      note.User.Username,
			"status":        note.Status,
			"view_count":    note.ViewCount,
			"like_count":    note.LikeCount,
			"comment_count": int(commentCount),
			"created_at":    note.CreatedAt,
			"updated_at":    note.UpdatedAt,
			"favorited_at":  favMap[note.ID].CreatedAt,
		})
	}
	c.JSON(http.StatusOK, gin.H{"data": responseNotes, "total": total, "page": req.Page, "limit": req.Limit})
}

func GetMyLikes(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}
	var req struct {
		Page  int `json:"page"`
		Limit int `json:"limit"`
	}
	req.Page = 1
	req.Limit = 10
	_ = c.ShouldBindJSON(&req)
	if req.Page < 1 {
		req.Page = 1
	}
	if req.Limit < 1 || req.Limit > 100 {
		req.Limit = 10
	}
	likes, err := model.GetNoteLikesByUser(database.DB, userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询点赞失败"})
		return
	}
	notes, total, err := model.GetLikedNotesForUser(database.DB, userID.(uint), req.Page, req.Limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询笔记失败"})
		return
	}
	likeMap := make(map[uint]model.NoteLike)
	for _, l := range likes {
		likeMap[l.NoteID] = l
	}
	var responseNotes []gin.H
	for _, note := range notes {
		commentCount, _ := model.CountCommentsForNote(database.DB, note.ID)
		var courseID uint
		var courseName string
		if note.CourseID != nil {
			courseID = *note.CourseID
			if note.Course.ID != 0 {
				courseName = note.Course.Name
			}
		}
		responseNotes = append(responseNotes, gin.H{
			"id":            note.ID,
			"title":         note.Title,
			"content":       note.Content,
			"description":   note.Description,
			"course_id":     courseID,
			"course_name":   courseName,
			"user_id":       note.UserID,
			"username":      note.User.Username,
			"status":        note.Status,
			"view_count":    note.ViewCount,
			"like_count":    note.LikeCount,
			"comment_count": int(commentCount),
			"created_at":    note.CreatedAt,
			"updated_at":    note.UpdatedAt,
			"liked_at":      likeMap[note.ID].CreatedAt,
		})
	}
	c.JSON(http.StatusOK, gin.H{"data": responseNotes, "total": total, "page": req.Page, "limit": req.Limit})
}

func SearchNotes(c *gin.Context) {
	var params req.SearchNotesRequest
	params.Page = 1
	params.PageSize = 10
	params.SortBy = "created_at"
	params.Order = "desc"
	_ = c.ShouldBindJSON(&params)
	kw := strings.TrimSpace(params.Keyword)
	if kw == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "搜索关键词不能为空"})
		return
	}
	notes, total, err := model.SearchPublicNotes(database.DB, kw, params.SortBy, params.Order, params.Page, params.PageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "搜索笔记失败"})
		return
	}
	var noteResponses []resp.NoteResponse
	for _, note := range notes {
		commentCount, _ := model.CountCommentsForNote(database.DB, note.ID)
		var courseID uint
		var courseName string
		if note.CourseID != nil {
			courseID = *note.CourseID
			courseName = note.Course.Name
		}
		noteResponses = append(noteResponses, resp.NoteResponse{
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
	page := params.Page
	pageSize := params.PageSize
	c.JSON(http.StatusOK, gin.H{"total": total, "page": page, "page_size": pageSize, "total_pages": (int(total) + pageSize - 1) / pageSize, "notes": noteResponses})
}

func GetPopularNotes(c *gin.Context) {
	notes, err := model.GetPopularPublicNotes(database.DB, 3)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取热门笔记失败"})
		return
	}
	var noteResponses []resp.NoteResponse
	for _, note := range notes {
		commentCount, _ := model.CountCommentsForNote(database.DB, note.ID)
		var courseID uint
		var courseName string
		if note.CourseID != nil {
			courseID = *note.CourseID
			if note.Course.ID != 0 {
				courseName = note.Course.Name
			}
		}
		noteResponses = append(noteResponses, resp.NoteResponse{
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
	c.JSON(http.StatusOK, gin.H{"notes": noteResponses})
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
