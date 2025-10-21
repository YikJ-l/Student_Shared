package api

import (
	"net/http"
	"student_shared/app/model"
	"student_shared/app/utils/database"

	"github.com/gin-gonic/gin"
)

// GetPublicStats 获取公开统计数据
func GetPublicStats(c *gin.Context) {
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

	c.JSON(http.StatusOK, gin.H{
		"users":    userCount,
		"courses":  courseCount,
		"notes":    noteCount,
		"comments": commentCount,
	})
}