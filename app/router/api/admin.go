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

// GetAllUsers 管理员获取所有用户列表（POST + JSON）
func GetAllUsers(c *gin.Context) {
	var params req.AdminUserListRequest
	params.Page = 1
	params.Limit = 10
	_ = c.ShouldBindJSON(&params)

	filter := model.AdminUserListFilter{ Search: params.Search, Role: params.Role }
	users, total, err := model.ListUsersAdmin(database.DB, filter, params.Page, params.Limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询用户失败"})
		return
	}

	var userResponses []resp.UserProfileResponse
	for _, user := range users {
		userResponses = append(userResponses, resp.UserProfileResponse{
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

	c.JSON(http.StatusOK, resp.AdminUserListResponse{ Data: userResponses, Total: total, Page: params.Page, Limit: params.Limit })
}

// AdminUpdateUser 管理员更新用户信息（POST + JSON）
func AdminUpdateUser(c *gin.Context) {
	var params struct {
		ID uint `json:"id" binding:"required"`
		req.AdminUpdateUserRequest
	}
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数无效"})
		return
	}

	fields := model.AdminUserUpdateFields{
		Nickname:     params.Nickname,
		Email:        params.Email,
		School:       params.School,
		Department:   params.Department,
		Major:        params.Major,
		Introduction: params.Introduction,
	}
	user, err := model.AdminUpdateUserFields(database.DB, params.ID, fields)
	if err != nil {
		// 邮箱唯一性冲突
		if err == gorm.ErrInvalidData {
			c.JSON(http.StatusConflict, gin.H{"error": "邮箱已被其他用户使用"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新用户信息失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "用户信息更新成功", "user": resp.UserProfileResponse{
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
	}})
}

// AdminUpdateUserRole 管理员更新用户角色（POST + JSON）
func AdminUpdateUserRole(c *gin.Context) {
	var params struct {
		ID uint `json:"id" binding:"required"`
		req.AdminUpdateUserRoleRequest
	}
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数无效"})
		return
	}

	user, err := model.AdminUpdateUserRole(database.DB, params.ID, params.Role)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "用户角色更新成功", "user": resp.UserProfileResponse{
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
	}})
}

// AdminDeleteUser 管理员删除用户（POST + JSON）
func AdminDeleteUser(c *gin.Context) {
	var params req.AdminDeleteUserRequest
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
		return
	}

	// 防止删除自己
	currentUserID, exists := c.Get("userID")
	if exists && currentUserID == params.ID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "不能删除自己的账户"})
		return
	}

	if err := model.AdminDeleteUserSoft(database.DB, params.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除用户失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "用户删除成功"})
}

// GetAdminStats 获取管理员统计数据
func GetAdminStats(c *gin.Context) {
	s, err := model.GetAdminStats(database.DB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取统计数据失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"total_users":     s.TotalUsers,
		"total_courses":   s.TotalCourses,
		"total_notes":     s.TotalNotes,
		"total_comments":  s.TotalComments,
		"today_new_users": s.TodayNewUsers,
		"online_users":    s.OnlineUsers,
	})
}
