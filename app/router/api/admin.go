package api

import (
	"net/http"
	"strconv"
	"student_shared/app/model"
	"student_shared/app/utils/database"
	"time"

	"github.com/gin-gonic/gin"
)

// AdminUserListRequest 管理员获取用户列表请求
type AdminUserListRequest struct {
	Page   int    `form:"page" binding:"min=1"`
	Limit  int    `form:"limit" binding:"min=1,max=100"`
	Search string `form:"search"`
	Role   string `form:"role"`
}

// AdminUserListResponse 管理员用户列表响应
type AdminUserListResponse struct {
	Data  []UserProfileResponse `json:"data"`
	Total int64                 `json:"total"`
	Page  int                   `json:"page"`
	Limit int                   `json:"limit"`
}

// AdminUpdateUserRequest 管理员更新用户请求
type AdminUpdateUserRequest struct {
	Nickname     string `json:"nickname"`
	Email        string `json:"email"`
	School       string `json:"school"`
	Department   string `json:"department"`
	Major        string `json:"major"`
	Introduction string `json:"introduction"`
}

// AdminUpdateUserRoleRequest 管理员更新用户角色请求
type AdminUpdateUserRoleRequest struct {
	Role string `json:"role" binding:"required,oneof=student teacher admin"`
}

// 检查管理员权限的中间件
func RequireAdminRole() gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("role")
		if !exists || userRole != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "需要管理员权限"})
			c.Abort()
			return
		}
		c.Next()
	}
}

// GetAllUsers 管理员获取所有用户列表
func GetAllUsers(c *gin.Context) {
	var req AdminUserListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数无效"})
		return
	}

	// 设置默认值
	if req.Page == 0 {
		req.Page = 1
	}
	if req.Limit == 0 {
		req.Limit = 10
	}

	// 构建查询
	query := database.DB.Model(&model.User{})

	// 搜索条件
	if req.Search != "" {
		query = query.Where("username LIKE ? OR email LIKE ? OR nickname LIKE ?",
			"%"+req.Search+"%", "%"+req.Search+"%", "%"+req.Search+"%")
	}

	// 角色筛选
	if req.Role != "" {
		query = query.Where("role = ?", req.Role)
	}

	// 获取总数
	var total int64
	query.Count(&total)

	// 分页查询
	var users []model.User
	offset := (req.Page - 1) * req.Limit
	result := query.Offset(offset).Limit(req.Limit).Order("created_at DESC").Find(&users)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询用户失败"})
		return
	}

	// 转换为响应格式
	var userResponses []UserProfileResponse
	for _, user := range users {
		userResponses = append(userResponses, UserProfileResponse{
			ID:           user.ID,
			Username:     user.Username,
			Email:        user.Email,
			Nickname:     user.Nickname,
			Avatar:       user.Avatar,
			School:       user.School,
			Department:   user.Department,
			Major:        user.Major,
			Introduction: user.Introduction,
			Role:         user.Role,
			LastLogin:    user.LastLogin,
			CreatedAt:    user.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, AdminUserListResponse{
		Data:  userResponses,
		Total: total,
		Page:  req.Page,
		Limit: req.Limit,
	})
}

// AdminUpdateUser 管理员更新用户信息
func AdminUpdateUser(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
		return
	}

	// 查询用户
	var user model.User
	result := database.DB.First(&user, userID)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	// 绑定请求数据
	var req AdminUpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数无效"})
		return
	}

	// 更新用户信息
	if req.Nickname != "" {
		user.Nickname = req.Nickname
	}
	if req.Email != "" {
		// 检查邮箱是否已被其他用户使用
		var existingUser model.User
		result := database.DB.Where("email = ? AND id != ?", req.Email, userID).First(&existingUser)
		if result.RowsAffected > 0 {
			c.JSON(http.StatusConflict, gin.H{"error": "邮箱已被其他用户使用"})
			return
		}
		user.Email = req.Email
	}
	if req.School != "" {
		user.School = req.School
	}
	if req.Department != "" {
		user.Department = req.Department
	}
	if req.Major != "" {
		user.Major = req.Major
	}
	if req.Introduction != "" {
		user.Introduction = req.Introduction
	}

	// 保存更新
	result = database.DB.Save(&user)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新用户信息失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "用户信息更新成功",
		"user": UserProfileResponse{
			ID:           user.ID,
			Username:     user.Username,
			Email:        user.Email,
			Nickname:     user.Nickname,
			Avatar:       user.Avatar,
			School:       user.School,
			Department:   user.Department,
			Major:        user.Major,
			Introduction: user.Introduction,
			Role:         user.Role,
			LastLogin:    user.LastLogin,
			CreatedAt:    user.CreatedAt,
		},
	})
}

// AdminUpdateUserRole 管理员更新用户角色
func AdminUpdateUserRole(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
		return
	}

	// 查询用户
	var user model.User
	result := database.DB.First(&user, userID)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	// 绑定请求数据
	var req AdminUpdateUserRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数无效"})
		return
	}

	// 更新用户角色
	user.Role = req.Role

	// 保存更新
	result = database.DB.Save(&user)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新用户角色失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "用户角色更新成功",
		"user": UserProfileResponse{
			ID:           user.ID,
			Username:     user.Username,
			Email:        user.Email,
			Nickname:     user.Nickname,
			Avatar:       user.Avatar,
			School:       user.School,
			Department:   user.Department,
			Major:        user.Major,
			Introduction: user.Introduction,
			Role:         user.Role,
			LastLogin:    user.LastLogin,
			CreatedAt:    user.CreatedAt,
		},
	})
}

// AdminDeleteUser 管理员删除用户
func AdminDeleteUser(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
		return
	}

	// 获取当前管理员ID，防止删除自己
	currentUserID, exists := c.Get("userID")
	if exists && currentUserID == uint(userID) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "不能删除自己的账户"})
		return
	}

	// 查询用户
	var user model.User
	result := database.DB.First(&user, userID)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	// 软删除用户
	result = database.DB.Delete(&user)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除用户失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "用户删除成功"})
}

// GetAdminStats 获取管理员统计数据
func GetAdminStats(c *gin.Context) {
	// 统计用户数量
	var userCount int64
	database.DB.Model(&model.User{}).Count(&userCount)

	// 统计课程数量
	var courseCount int64
	database.DB.Model(&model.Course{}).Count(&courseCount)

	// 统计笔记数量
	var noteCount int64
	database.DB.Model(&model.Note{}).Count(&noteCount)

	// 统计评论数量
	var commentCount int64
	database.DB.Model(&model.Comment{}).Count(&commentCount)

	// 统计今日新增用户
	var todayUserCount int64
	today := time.Now().Format("2006-01-02")
	database.DB.Model(&model.User{}).Where("DATE(created_at) = ?", today).Count(&todayUserCount)

	c.JSON(http.StatusOK, gin.H{
		"total_users":     userCount,
		"total_courses":   courseCount,
		"total_notes":     noteCount,
		"total_comments":  commentCount,
		"today_new_users": todayUserCount,
		"online_users":    0, // 这里可以实现在线用户统计
	})
}
