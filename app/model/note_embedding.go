package model

// NoteEmbedding 独立存储笔记的嵌入向量，避免重复计算
// 一条笔记仅维护最新版本的嵌入（uniqueIndex）
type NoteEmbedding struct {
	BaseModel

	NoteID    uint   `gorm:"uniqueIndex" json:"note_id"`
	Embedding string `gorm:"type:text" json:"-"` // 嵌入向量的JSON字符串

	// 关联
	Note Note `json:"-"`
}