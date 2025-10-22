package api

import (
	"net/http"
	"student_shared/app/model"
	"student_shared/app/utils/database"

	"github.com/gin-gonic/gin"
)

// GetPublicStats 获取公开统计数据
func GetPublicStats(c *gin.Context) {
	stats, err := model.GetPublicStats(database.DB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取统计数据失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"users":    stats.Users,
		"courses":  stats.Courses,
		"notes":    stats.Notes,
		"comments": stats.Comments,
	})
}