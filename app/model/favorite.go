package model

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