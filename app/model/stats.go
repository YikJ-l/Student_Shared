package model

import (
	"time"
	"gorm.io/gorm"
)

// PublicStats 公开统计数据
type PublicStats struct {
	Users    int64
	Courses  int64
	Notes    int64
	Comments int64
}

// GetPublicStats 查询公开统计数据
func GetPublicStats(db *gorm.DB) (PublicStats, error) {
	var s PublicStats
	if err := db.Model(&User{}).Count(&s.Users).Error; err != nil { return s, err }
	if err := db.Model(&Course{}).Count(&s.Courses).Error; err != nil { return s, err }
	if err := db.Model(&Note{}).Count(&s.Notes).Error; err != nil { return s, err }
	if err := db.Model(&Comment{}).Count(&s.Comments).Error; err != nil { return s, err }
	return s, nil
}

// AdminStats 管理员统计数据
type AdminStats struct {
	TotalUsers     int64
	TotalCourses   int64
	TotalNotes     int64
	TotalComments  int64
	TodayNewUsers  int64
	OnlineUsers    int64
}

// GetAdminStats 获取管理员统计数据
func GetAdminStats(db *gorm.DB) (AdminStats, error) {
	var s AdminStats
	if err := db.Model(&User{}).Count(&s.TotalUsers).Error; err != nil { return s, err }
	if err := db.Model(&Course{}).Count(&s.TotalCourses).Error; err != nil { return s, err }
	if err := db.Model(&Note{}).Count(&s.TotalNotes).Error; err != nil { return s, err }
	if err := db.Model(&Comment{}).Count(&s.TotalComments).Error; err != nil { return s, err }
	// 今日新增用户（按日期）
	today := time.Now().Format("2006-01-02")
	if err := db.Model(&User{}).Where("DATE(created_at) = ?", today).Count(&s.TodayNewUsers).Error; err != nil { return s, err }
	// 在线用户暂不统计，默认 0（可接入会话/心跳数据时更新）
	s.OnlineUsers = 0
	return s, nil
}