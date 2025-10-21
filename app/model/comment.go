package model

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