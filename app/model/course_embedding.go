package model

// CourseEmbedding 独立存储课程的嵌入向量，避免重复计算
// 一门课程仅维护最新版本的嵌入（uniqueIndex）
type CourseEmbedding struct {
	BaseModel

	CourseID  uint   `gorm:"uniqueIndex" json:"course_id"`
	Embedding string `gorm:"type:text" json:"-"` // 嵌入向量的JSON字符串

	// 关联
	Course Course `json:"-"`
}