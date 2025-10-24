package resp

import (
	"time"
)

// -------- 用户 --------

type UserProfileResponse struct {
	ID           uint      `json:"id"`
	Username     string    `json:"username"`
	Nickname     string    `json:"nickname"`
	Email        string    `json:"email"`
	School       string    `json:"school"`
	Department   string    `json:"department"`
	Major        string    `json:"major"`
	Introduction string    `json:"introduction"`
	Avatar       string    `json:"avatar"`
	Role         string    `json:"role"`
	LastLogin    time.Time `json:"last_login"`
	CreatedAt    time.Time `json:"created_at"`
}

// -------- 评论 --------

type CommentResponse struct {
	ID              uint              `json:"id"`
	Content         string            `json:"content"`
	Rating          int               `json:"rating"`
	LikeCount       int               `json:"like_count"`
	UserID          uint              `json:"user_id"`
	NoteID          uint              `json:"note_id"`
	Username        string            `json:"username"`
	UserAvatar      string            `json:"user_avatar"`
	ParentID        *uint             `json:"parent_id"`
	ReplyToUserID   *uint             `json:"reply_to_user_id"`
	ReplyToUsername *string           `json:"reply_to_username"`
	Replies         []CommentResponse `json:"replies"`
	IsLiked         bool              `json:"is_liked"`
	CreatedAt       time.Time         `json:"created_at"`
	UpdatedAt       time.Time         `json:"updated_at"`
}

// -------- 笔记 --------

type NoteResponse struct {
	ID            uint      `json:"id"`
	Title         string    `json:"title"`
	Content       string    `json:"content"`
	Description   string    `json:"description"`
	FilePath      string    `json:"file_path,omitempty"`
	FileSize      int64     `json:"file_size"`
	FileType      string    `json:"file_type"`
	ViewCount     int       `json:"view_count"`
	DownloadCount int       `json:"download_count"`
	LikeCount     int       `json:"like_count"`
	CommentCount  int       `json:"comment_count"`
	Status        string    `json:"status"`
	UserID        uint      `json:"user_id"`
	CourseID      uint      `json:"course_id"`
	Username      string    `json:"username"`
	AuthorAvatar  string    `json:"author_avatar"`
	CourseName    string    `json:"course_name"`
	IsLiked       bool      `json:"is_liked"`
	IsFavorited   bool      `json:"is_favorited"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	// 语义检索可选字段
	Similarity    float64   `json:"similarity,omitempty"`
	Excerpt       string    `json:"excerpt,omitempty"`
	HighlightTitle       string    `json:"highlight_title,omitempty"`
	HighlightDescription string    `json:"highlight_description,omitempty"`
}

// -------- 管理员 --------

type AdminUserListResponse struct {
	Data  []UserProfileResponse `json:"data"`
	Total int64                 `json:"total"`
	Page  int                   `json:"page"`
	Limit int                   `json:"limit"`
}

// -------- 课程 --------

type CourseResponse struct {
	ID           uint      `json:"id"`
	Code         string    `json:"code"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	School       string    `json:"school"`
	Department   string    `json:"department"`
	Teacher      string    `json:"teacher"`
	Semester     string    `json:"semester"`
	Cover        string    `json:"cover"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	IsJoined     bool      `json:"is_joined,omitempty"`
	StudentCount int64     `json:"student_count,omitempty"`
	NoteCount    int64     `json:"note_count,omitempty"`
	// 语义检索可选字段
	Similarity   float64   `json:"similarity,omitempty"`
	HighlightName        string    `json:"highlight_name,omitempty"`
	HighlightDescription string    `json:"highlight_description,omitempty"`
}

type MyCourseResponse struct {
	ID           uint      `json:"id"`
	Code         string    `json:"code"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	School       string    `json:"school"`
	Department   string    `json:"department"`
	Teacher      string    `json:"teacher"`
	Semester     string    `json:"semester"`
	Cover        string    `json:"cover"`
	Status       string    `json:"status"`
	Role         string    `json:"role"`
	JoinedAt     time.Time `json:"joined_at"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	StudentCount int64     `json:"student_count,omitempty"`
	NoteCount    int64     `json:"note_count,omitempty"`
}

// -------- 上传 --------

type UploadResponse struct {
	URL  string `json:"url"`
	Size int64  `json:"size"`
}