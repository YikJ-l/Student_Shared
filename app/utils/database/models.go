package database

import (
	"gorm.io/gorm"
	"time"
)

// BaseModel 基础模型，包含共有字段
type BaseModel struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// User 用户模型
type User struct {
	BaseModel
	Username    string    `gorm:"size:50;not null;uniqueIndex" json:"username"`
	Password    string    `gorm:"size:100;not null" json:"-"`
	Email       string    `gorm:"size:100;not null;uniqueIndex" json:"email"`
	Nickname    string    `gorm:"size:50" json:"nickname"`
	Avatar      string    `gorm:"size:255" json:"avatar"`
	School      string    `gorm:"size:100" json:"school"`
	Department  string    `gorm:"size:100" json:"department"`
	Major       string    `gorm:"size:100" json:"major"`
	Introduction string    `gorm:"size:500" json:"introduction"`
	LastLogin   time.Time `json:"last_login"`
	Role        string    `gorm:"size:20;default:'student'" json:"role"` // student, teacher, admin
	
	// 关联
	Courses   []Course   `gorm:"many2many:user_courses;" json:"-"`
	Notes     []Note     `json:"-"`
	Comments  []Comment  `json:"-"`
	Favorites []Favorite `json:"-"`
}

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
	Notes []Note `json:"-"`
}

// Note 笔记模型
type Note struct {
	BaseModel
	Title       string `gorm:"size:200;not null" json:"title"`
	Content     string `gorm:"type:text" json:"content"`
	Description string `gorm:"size:500" json:"description"`
	FilePath    string `gorm:"size:255" json:"file_path"`
	FileSize    int64  `json:"file_size"`
	FileType    string `gorm:"size:50" json:"file_type"`
	ViewCount   int    `gorm:"default:0" json:"view_count"`
	DownloadCount int  `gorm:"default:0" json:"download_count"`
	LikeCount   int    `gorm:"default:0" json:"like_count"`
	Status      string `gorm:"size:20;default:'public'" json:"status"` // public, private
	
	// 外键
	UserID   uint   `json:"user_id"`
	CourseID uint   `json:"course_id"`
	
	// 关联
	User     User      `json:"-"`
	Course   Course    `json:"-"`
	Comments []Comment `json:"-"`
}

// Comment 评论模型
type Comment struct {
	BaseModel
	Content string `gorm:"size:1000;not null" json:"content"`
	Rating  int    `gorm:"default:0" json:"rating"` // 0-5星评分
	
	// 外键
	UserID uint `json:"user_id"`
	NoteID uint `json:"note_id"`
	
	// 关联
	User User `json:"-"`
	Note Note `json:"-"`
}

// Favorite 收藏模型
type Favorite struct {
	BaseModel
	
	// 外键
	UserID uint `json:"user_id"`
	NoteID uint `json:"note_id"`
	
	// 关联
	User User `json:"-"`
	Note Note `json:"-"`
}