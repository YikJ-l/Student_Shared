package model

import "time"

// Course 课程模型
type Course struct {
	BaseModel
	Code        string `gorm:"size:50;not null;uniqueIndex" json:"code"`
	Name        string `gorm:"size:100;not null" json:"name"`
	Description string `gorm:"size:1000" json:"description"`
	School      string `gorm:"size:100" json:"school"`
	Department  string `gorm:"size:100" json:"department"`
	Teacher     string `gorm:"size:100" json:"teacher"`
	Semester    string `gorm:"size:50" json:"semester"`
	Cover       string `gorm:"size:255" json:"cover"`
	Status      string `gorm:"size:20;default:'active'" json:"status"` // active, inactive
	
	// 关联
	Users []User `gorm:"many2many:user_courses;" json:"-"`
}

// UserCourse 用户-课程关联表
type UserCourse struct {
	UserID   uint `gorm:"primaryKey"`
	CourseID uint `gorm:"primaryKey"`
	Role     string `gorm:"size:20;default:'student'"` // student, teacher
	JoinedAt time.Time
}