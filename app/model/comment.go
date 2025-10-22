package model

import "gorm.io/gorm"

// Comment 评论模型
type Comment struct {
	BaseModel
	Content   string `gorm:"size:1000;not null" json:"content"`
	Rating    int    `gorm:"default:0" json:"rating"` // 0-5星评分
	LikeCount int    `gorm:"default:0" json:"like_count"` // 点赞数
	
	// 回复相关字段
	ParentID *uint `gorm:"index" json:"parent_id"` // 父评论ID，为空表示顶级评论
	ReplyToUserID *uint `json:"reply_to_user_id"` // 回复的用户ID
	
	// 外键
	UserID uint `json:"user_id"`
	NoteID uint `json:"note_id"`
	
	// 关联
	User User `json:"-"`
	Note Note `json:"-"`
	Parent *Comment `gorm:"foreignKey:ParentID" json:"-"`
	Replies []Comment `gorm:"foreignKey:ParentID" json:"-"`
	ReplyToUser *User `gorm:"foreignKey:ReplyToUserID" json:"-"`
}

// CommentLike 评论点赞模型
type CommentLike struct {
	BaseModel
	UserID    uint `gorm:"not null;index;uniqueIndex:uk_comment_likes_user_comment" json:"user_id"`
	CommentID uint `gorm:"not null;index;uniqueIndex:uk_comment_likes_user_comment" json:"comment_id"`
	IsActive  bool `gorm:"default:true" json:"is_active"` // 点赞状态，true为点赞，false为取消点赞
	
	// 关联
	User    User    `json:"-"`
	Comment Comment `json:"-"`
}

// GetCommentByID 根据ID获取评论
func GetCommentByID(db *gorm.DB, id uint) (Comment, error) {
	var c Comment
	err := db.First(&c, id).Error
	return c, err
}

// CountTopLevelComments 统计顶级评论数量
func CountTopLevelComments(db *gorm.DB, noteID uint) (int64, error) {
	var total int64
	err := db.Model(&Comment{}).Where("note_id = ? AND parent_id IS NULL", noteID).Count(&total).Error
	return total, err
}

// FindTopLevelCommentsWithReplies 获取顶级评论及其回复（包含关联用户）
func FindTopLevelCommentsWithReplies(db *gorm.DB, noteID uint, page, pageSize int) ([]Comment, error) {
	var comments []Comment
	offset := (page - 1) * pageSize
	res := db.Where("note_id = ? AND parent_id IS NULL", noteID).
		Preload("User").
		Preload("Replies", func(tx *gorm.DB) *gorm.DB { return tx.Preload("User").Preload("ReplyToUser").Order("created_at ASC") }).
		Preload("Replies.User").
		Preload("Replies.ReplyToUser").
		Order("created_at DESC").
		Limit(pageSize).Offset(offset).Find(&comments)
	return comments, res.Error
}

// IsCommentLikedByUser 判断用户是否点赞了指定评论
func IsCommentLikedByUser(db *gorm.DB, commentID, userID uint) (bool, error) {
	var count int64
	err := db.Model(&CommentLike{}).Where("comment_id = ? AND user_id = ? AND is_active = ?", commentID, userID, true).Count(&count).Error
	return count > 0, err
}

// CreateComment 创建评论
func CreateComment(db *gorm.DB, c *Comment) error {
	return db.Create(c).Error
}

// GetCommentWithPreloads 获取评论及其关联用户
func GetCommentWithPreloads(db *gorm.DB, id uint) (Comment, error) {
	var c Comment
	err := db.Preload("User").Preload("ReplyToUser").First(&c, id).Error
	return c, err
}

// UpdateComment 保存评论更新
func UpdateComment(db *gorm.DB, c *Comment) error {
	return db.Save(c).Error
}

// DeleteCommentCascade 递归删除评论及其所有子评论（包含点赞记录）
func DeleteCommentCascade(db *gorm.DB, commentID uint) error {
	// 查找子评论
	var children []Comment
	if err := db.Where("parent_id = ?", commentID).Find(&children).Error; err != nil { return err }
	// 递归删除子评论
	for _, child := range children {
		if err := DeleteCommentCascade(db, child.ID); err != nil { return err }
		if err := db.Where("comment_id = ?", child.ID).Delete(&CommentLike{}).Error; err != nil { return err }
		if err := db.Delete(&child).Error; err != nil { return err }
	}
	return nil
}

// DeleteCommentWithCascade 删除根评论及其子评论与点赞
func DeleteCommentWithCascade(db *gorm.DB, id uint) error {
	if err := DeleteCommentCascade(db, id); err != nil { return err }
	if err := db.Where("comment_id = ?", id).Delete(&CommentLike{}).Error; err != nil { return err }
	return db.Delete(&Comment{}, id).Error
}

// LikeCommentTx 点赞评论（事务处理），返回新点赞数以及是否已点赞
func LikeCommentTx(db *gorm.DB, commentID, userID uint) (int, bool, error) {
	tx := db.Begin()
	defer func() { if r := recover(); r != nil { tx.Rollback() } }()
	var comment Comment
	if err := tx.First(&comment, commentID).Error; err != nil { tx.Rollback(); return 0, false, err }
	var existing CommentLike
	res := tx.Where("comment_id = ? AND user_id = ?", commentID, userID).First(&existing)
	if res.RowsAffected > 0 {
		if existing.IsActive { tx.Rollback(); return comment.LikeCount, true, nil }
		if err := tx.Model(&existing).Update("is_active", true).Error; err != nil { tx.Rollback(); return 0, false, err }
	} else {
		like := CommentLike{CommentID: commentID, UserID: userID, IsActive: true}
		if err := tx.Create(&like).Error; err != nil { tx.Rollback(); return 0, false, err }
	}
	newCount := comment.LikeCount + 1
	if err := tx.Model(&comment).Update("like_count", newCount).Error; err != nil { tx.Rollback(); return 0, false, err }
	if err := tx.Commit().Error; err != nil { return 0, false, err }
	return newCount, false, nil
}

// UnlikeCommentTx 取消点赞评论（事务处理），返回新点赞数以及是否原本已点赞
func UnlikeCommentTx(db *gorm.DB, commentID, userID uint) (int, bool, error) {
	tx := db.Begin()
	defer func() { if r := recover(); r != nil { tx.Rollback() } }()
	var comment Comment
	if err := tx.First(&comment, commentID).Error; err != nil { tx.Rollback(); return 0, false, err }
	var like CommentLike
	res := tx.Where("comment_id = ? AND user_id = ?", commentID, userID).First(&like)
	if res.RowsAffected == 0 || !like.IsActive { tx.Rollback(); return comment.LikeCount, false, nil }
	if err := tx.Model(&like).Update("is_active", false).Error; err != nil { tx.Rollback(); return 0, false, err }
	newCount := comment.LikeCount
	if newCount > 0 { newCount-- }
	if err := tx.Model(&comment).Update("like_count", newCount).Error; err != nil { tx.Rollback(); return 0, false, err }
	if err := tx.Commit().Error; err != nil { return 0, false, err }
	return newCount, true, nil
}