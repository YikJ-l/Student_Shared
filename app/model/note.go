package model

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
	CourseID *uint  `json:"course_id"` // 使用指针类型，允许为nil表示无关联课程
	
	// 关联
	User     User      `json:"-"`
	Course   Course    `json:"-"`
}

// NoteLike 笔记点赞模型
type NoteLike struct {
	BaseModel
	UserID uint `gorm:"not null;index" json:"user_id"`
	NoteID uint `gorm:"not null;index" json:"note_id"`
	
	// 关联
	User User `json:"-"`
	Note Note `json:"-"`
}