package req

import (
	"mime/multipart"
)

// -------- 评论 --------

// CommentRequest 评论创建请求
type CommentRequest struct {
	Content       string `json:"content" binding:"required"`
	NoteID        uint   `json:"note_id" binding:"required"`
	Rating        int    `json:"rating"`           // 0-5星评分
	ParentID      *uint  `json:"parent_id"`        // 父评论ID，用于回复
	ReplyToUserID *uint  `json:"reply_to_user_id"` // 回复的用户ID
}

// ListCommentsRequest 获取评论列表请求
type ListCommentsRequest struct {
	NoteID   uint `json:"note_id" binding:"required"`
	Page     int  `json:"page"`
	PageSize int  `json:"page_size"`
}

// CommentUpdateRequest 评论更新请求
type CommentUpdateRequest struct {
	ID      uint   `json:"id" binding:"required"`
	Content string `json:"content" binding:"required"`
	Rating  int    `json:"rating"`
}

// GetCommentRequest 获取/操作单个评论请求
type GetCommentRequest struct {
	ID uint `json:"id" binding:"required"`
}

// -------- 笔记 --------

// NoteRequest 笔记创建请求
type NoteRequest struct {
	Title       string `json:"title" binding:"required"`
	Content     string `json:"content"`
	Description string `json:"description"`
	CourseID    *uint  `json:"course_id"`
	Status      string `json:"status"`
}

// NoteUpdateRequest 笔记更新请求（部分更新）
type NoteUpdateRequest struct {
	Title       *string `json:"title"`
	Content     *string `json:"content"`
	Description *string `json:"description"`
	CourseID    *uint   `json:"course_id"`
	Status      *string `json:"status"`
}

// ListNotesRequest 获取笔记列表请求
type ListNotesRequest struct {
	CourseID *uint  `json:"course_id"`
	UserID   *uint  `json:"user_id"`
	Status   string `json:"status"`
	Page     int    `json:"page"`
	PageSize int    `json:"page_size"`
	SortBy   string `json:"sort_by"`
	Order    string `json:"order"`
}

// GetNoteRequest 获取单个笔记请求
type GetNoteRequest struct {
	ID uint `json:"id"`
}

// UpdateNoteRequest 更新笔记请求（包含ID）
type UpdateNoteRequest struct {
	ID uint `json:"id"`
	NoteUpdateRequest
}

// SearchNotesRequest 搜索笔记请求
type SearchNotesRequest struct {
	Keyword  string `json:"keyword"`
	Page     int    `json:"page"`
	PageSize int    `json:"page_size"`
	SortBy   string `json:"sort_by"`
	Order    string `json:"order"`
}

// -------- 课程 --------

// CourseRequest 课程创建/更新公共字段
type CourseRequest struct {
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description"`
	School      string `json:"school"`
	Department  string `json:"department"`
	Teacher     string `json:"teacher"`
	Semester    string `json:"semester"`
	Cover       string `json:"cover"`
	Status      string `json:"status"`
}

// ListCoursesRequest 获取课程列表
type ListCoursesRequest struct {
	School     string `json:"school"`
	Department string `json:"department"`
	Semester   string `json:"semester"`
	Search     string `json:"search"`
	Status     string `json:"status"`
	Page       int    `json:"page"`
	PageSize   int    `json:"page_size"`
}

// GetCourseRequest 获取课程详情
type GetCourseRequest struct {
	ID uint `json:"id"`
}

// UpdateCourseRequest 更新课程（包含ID）
type UpdateCourseRequest struct {
	ID uint `json:"id"`
	CourseRequest
}

// SearchCoursesRequest 搜索课程
type SearchCoursesRequest struct {
	Keyword  string `json:"keyword"`
	Page     int    `json:"page"`
	PageSize int    `json:"page_size"`
	SortBy   string `json:"sort_by"`
	Order    string `json:"order"`
	Status   string `json:"status"`
}

// ListMyCoursesRequest 我的课程列表分页
type ListMyCoursesRequest struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
}

// -------- 用户/认证 --------

// UserRegisterRequest 用户注册
type UserRegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Nickname string `json:"nickname"`
	School   string `json:"school"`
}

// UserLoginRequest 用户登录
type UserLoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// -------- 管理员 --------

// AdminUserListRequest 管理员用户列表
type AdminUserListRequest struct {
	Page   int    `json:"page" binding:"min=1"`
	Limit  int    `json:"limit" binding:"min=1,max=100"`
	Search string `json:"search"`
	Role   string `json:"role"`
}

// AdminUpdateUserRequest 管理员更新用户
type AdminUpdateUserRequest struct {
	Nickname     string `json:"nickname"`
	Email        string `json:"email"`
	School       string `json:"school"`
	Department   string `json:"department"`
	Major        string `json:"major"`
	Introduction string `json:"introduction"`
}

// AdminUpdateUserRoleRequest 管理员更新角色
type AdminUpdateUserRoleRequest struct {
	Role string `json:"role" binding:"required,oneof=student teacher admin"`
}

// AdminDeleteUserRequest 管理员删除用户
type AdminDeleteUserRequest struct {
	ID uint `json:"id" binding:"required"`
}

// -------- 上传 --------

// UploadAvatarRequest 头像上传（表单文件）
type UploadAvatarRequest struct {
	File *multipart.FileHeader `form:"file" binding:"required"`
}

// -------- AI --------

// SummarizeRequest AI 摘要请求（笔记ID或文本）
type SummarizeRequest struct {
	NoteID  *uint   `json:"note_id"`
	Content *string `json:"content"`
}
