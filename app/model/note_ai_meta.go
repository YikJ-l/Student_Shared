package model

import (
	"time"

	"gorm.io/gorm"
)

// NoteAIMeta 存储笔记的AI相关元数据（摘要、关键词、抄袭分数等）
type NoteAIMeta struct {
	BaseModel

	NoteID          uint       `gorm:"uniqueIndex" json:"note_id"`
	Summary         string     `gorm:"type:text" json:"summary"`
	Keywords        string     `gorm:"size:500" json:"keywords"` // 逗号分隔的关键词
	PlagiarismScore float64    `gorm:"default:0" json:"plagiarism_score"`
	Flags           string     `gorm:"size:100" json:"flags"` // 例如：is_plagiarized
	LastReviewedAt  *time.Time `json:"last_reviewed_at"`

	// 关联
	Note Note `json:"-"`
}

// --- NoteAIMeta 相关数据库操作封装 ---

// UpsertNoteAIMeta 根据 note_id 创建或更新 AI 元数据
func UpsertNoteAIMeta(db *gorm.DB, noteID uint, summary string, keywords string, reviewedAt time.Time) error {
	var meta NoteAIMeta
	res := db.Where("note_id = ?", noteID).First(&meta)
	if res.RowsAffected == 0 || res.Error == gorm.ErrRecordNotFound {
		meta = NoteAIMeta{
			NoteID:         noteID,
			Summary:        summary,
			Keywords:       keywords,
			PlagiarismScore: 0,
			Flags:          "",
			LastReviewedAt: &reviewedAt,
		}
		return db.Create(&meta).Error
	}
	if res.Error != nil {
		return res.Error
	}
	meta.Summary = summary
	meta.Keywords = keywords
	meta.LastReviewedAt = &reviewedAt
	return db.Save(&meta).Error
}

// GetNoteAIMetaByNoteID 获取某个笔记的AI元数据
func GetNoteAIMetaByNoteID(db *gorm.DB, noteID uint) (*NoteAIMeta, error) {
	var meta NoteAIMeta
	res := db.Where("note_id = ?", noteID).First(&meta)
	if res.Error != nil {
		return nil, res.Error
	}
	return &meta, nil
}