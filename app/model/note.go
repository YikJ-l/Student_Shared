package model

import (
	"errors"
	"strings"

	"gorm.io/gorm"
)

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
	Embedding   string `gorm:"type:text" json:"-"` // 保存嵌入向量的JSON字符串
	
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


// --- Note 相关数据库操作封装 ---

func GetNoteByID(db *gorm.DB, id uint) (*Note, error) {
	var n Note
	res := db.First(&n, id)
	if res.Error != nil {
		return nil, res.Error
	}
	return &n, nil
}

// GetNoteWithPreloads 获取笔记并预加载用户和课程
func GetNoteWithPreloads(db *gorm.DB, id uint) (*Note, error) {
	var n Note
	res := db.Preload("User").Preload("Course").First(&n, id)
	if res.Error != nil {
		return nil, res.Error
	}
	return &n, nil
}

// CountCommentsForNote 统计笔记的评论数量
func CountCommentsForNote(db *gorm.DB, noteID uint) (int64, error) {
	var cnt int64
	if err := db.Model(&Comment{}).Where("note_id = ?", noteID).Count(&cnt).Error; err != nil {
		return 0, err
	}
	return cnt, nil
}

// IncrementViewCount 查看次数 +1
func IncrementViewCount(db *gorm.DB, noteID uint) error {
	var n Note
	res := db.First(&n, noteID)
	if res.Error != nil { return res.Error }
	n.ViewCount++
	return db.Save(&n).Error
}

// IncrementDownloadCount 下载次数 +1
func IncrementDownloadCount(db *gorm.DB, noteID uint) error {
	var n Note
	res := db.First(&n, noteID)
	if res.Error != nil { return res.Error }
	n.DownloadCount++
	return db.Save(&n).Error
}

// IsNoteLikedByUser 用户是否点赞了该笔记
func IsNoteLikedByUser(db *gorm.DB, noteID, userID uint) (bool, error) {
	var nl NoteLike
	res := db.Where("note_id = ? AND user_id = ?", noteID, userID).First(&nl)
	if res.Error == nil { return true, nil }
	if errors.Is(res.Error, gorm.ErrRecordNotFound) { return false, nil }
	return false, res.Error
}

// IsNoteFavoritedByUser 用户是否收藏了该笔记
func IsNoteFavoritedByUser(db *gorm.DB, noteID, userID uint) (bool, error) {
	var f Favorite
	res := db.Where("note_id = ? AND user_id = ?", noteID, userID).First(&f)
	if res.Error == nil { return true, nil }
	if errors.Is(res.Error, gorm.ErrRecordNotFound) { return false, nil }
	return false, res.Error
}

// CreateNote 创建笔记
func CreateNote(db *gorm.DB, note *Note) error {
	return db.Create(note).Error
}

// UpdateNote 保存笔记更新
func UpdateNote(db *gorm.DB, note *Note) error {
	return db.Save(note).Error
}

// DeleteNoteByID 删除笔记
func DeleteNoteByID(db *gorm.DB, id uint) error {
	return db.Delete(&Note{}, id).Error
}

// NoteListFilter 列表过滤条件
type NoteListFilter struct {
	CourseID     *uint
	UserID       *uint
	Status       string // "public" | "private" | "all"
	CanViewPrivate bool // 是否允许查看私有笔记（本人或管理员）
}

// BuildNoteQuery 根据过滤条件构建查询
func BuildNoteQuery(db *gorm.DB, f NoteListFilter) *gorm.DB {
	q := db.Model(&Note{})
	if f.CourseID != nil {
		if *f.CourseID == 0 {
			q = q.Where("course_id IS NULL")
		} else {
			q = q.Where("course_id = ?", *f.CourseID)
		}
	}
	if f.UserID != nil {
		q = q.Where("user_id = ?", *f.UserID)
	}
	if strings.TrimSpace(f.Status) != "all" {
		q = q.Where("status = ?", strings.TrimSpace(f.Status))
	}
	if !f.CanViewPrivate {
		q = q.Where("status = ?", "public")
	}
	return q
}

// ListNotes 查询笔记列表（分页与排序）
func ListNotes(db *gorm.DB, f NoteListFilter, page, pageSize int, sortBy, order string) (notes []Note, total int64, err error) {
	if page < 1 { page = 1 }
	if pageSize < 1 || pageSize > 100 { pageSize = 10 }
	q := BuildNoteQuery(db, f)
	if err = q.Count(&total).Error; err != nil { return nil, 0, err }
	orderUpper := strings.ToUpper(order)
	if orderUpper != "ASC" { orderUpper = "DESC" }
	orderClause := "created_at DESC"
	switch sortBy {
	case "like_count":
		orderClause = "like_count " + orderUpper
	case "view_count":
		orderClause = "view_count " + orderUpper
	case "created_at":
		orderClause = "created_at " + orderUpper
	}
	offset := (page - 1) * pageSize
	res := q.Preload("User").Preload("Course").Offset(offset).Limit(pageSize).Order(orderClause).Find(&notes)
	return notes, total, res.Error
}

// LikeNoteTx 点赞（事务）
// 返回值：alreadyLiked 表示之前已点赞；likeCount 更新后的点赞数
func LikeNoteTx(db *gorm.DB, noteID, userID uint) (alreadyLiked bool, likeCount int, err error) {
	tx := db.Begin()
	defer func() { if r := recover(); r != nil { tx.Rollback() } }()
	var note Note
	if err = tx.First(&note, noteID).Error; err != nil { tx.Rollback(); return false, 0, err }
	var existing NoteLike
	res := tx.Where("note_id = ? AND user_id = ?", noteID, userID).First(&existing)
	if res.Error == nil { // 已存在点赞
		alreadyLiked = true
		likeCount = note.LikeCount
		_ = tx.Commit().Error
		return
	}
	if res.Error != nil && !errors.Is(res.Error, gorm.ErrRecordNotFound) { tx.Rollback(); return false, 0, res.Error }
	if err = tx.Create(&NoteLike{NoteID: noteID, UserID: userID}).Error; err != nil { tx.Rollback(); return false, 0, err }
	note.LikeCount++
	if err = tx.Save(&note).Error; err != nil { tx.Rollback(); return false, 0, err }
	likeCount = note.LikeCount
	return false, likeCount, tx.Commit().Error
}

// UnlikeNoteTx 取消点赞（事务）
// 返回值：notLiked 表示之前未点赞；likeCount 更新后的点赞数
func UnlikeNoteTx(db *gorm.DB, noteID, userID uint) (notLiked bool, likeCount int, err error) {
	tx := db.Begin()
	defer func() { if r := recover(); r != nil { tx.Rollback() } }()
	var note Note
	if err = tx.First(&note, noteID).Error; err != nil { tx.Rollback(); return false, 0, err }
	var existing NoteLike
	res := tx.Where("note_id = ? AND user_id = ?", noteID, userID).First(&existing)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			notLiked = true
			likeCount = note.LikeCount
			_ = tx.Commit().Error
			return
		}
		tx.Rollback(); return false, 0, res.Error
	}
	if err = tx.Delete(&existing).Error; err != nil { tx.Rollback(); return false, 0, err }
	if note.LikeCount > 0 { note.LikeCount-- }
	if err = tx.Save(&note).Error; err != nil { tx.Rollback(); return false, 0, err }
	likeCount = note.LikeCount
	return false, likeCount, tx.Commit().Error
}

// FavoriteNote 添加收藏
func FavoriteNote(db *gorm.DB, noteID, userID uint) (alreadyFavorited bool, err error) {
	var fav Favorite
	res := db.Where("note_id = ? AND user_id = ?", noteID, userID).First(&fav)
	if res.Error == nil { return true, nil }
	if res.Error != nil && !errors.Is(res.Error, gorm.ErrRecordNotFound) { return false, res.Error }
	return false, db.Create(&Favorite{NoteID: noteID, UserID: userID}).Error
}

// UnfavoriteNote 取消收藏
func UnfavoriteNote(db *gorm.DB, noteID, userID uint) (notFavorited bool, err error) {
	var fav Favorite
	res := db.Where("note_id = ? AND user_id = ?", noteID, userID).First(&fav)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) { return true, nil }
		return false, res.Error
	}
	return false, db.Delete(&fav).Error
}

// GetFavoritesNotesForUser 获取用户收藏的笔记（分页，预加载）
func GetFavoritesNotesForUser(db *gorm.DB, userID uint, page, limit int) (notes []Note, total int64, err error) {
	if page < 1 { page = 1 }
	if limit < 1 || limit > 100 { limit = 10 }
	var favorites []Favorite
	if err = db.Where("user_id = ?", userID).Find(&favorites).Error; err != nil { return nil, 0, err }
	noteIDs := make([]uint, 0, len(favorites))
	for _, f := range favorites { noteIDs = append(noteIDs, f.NoteID) }
	if len(noteIDs) == 0 { return []Note{}, 0, nil }
	if err = db.Model(&Note{}).Where("id IN ?", noteIDs).Count(&total).Error; err != nil { return nil, 0, err }
	offset := (page - 1) * limit
	res := db.Where("id IN ?", noteIDs).Preload("User").Preload("Course").Limit(limit).Offset(offset).Order("created_at DESC").Find(&notes)
	return notes, total, res.Error
}

// GetLikedNotesForUser 获取用户点赞的笔记（分页，预加载）
func GetLikedNotesForUser(db *gorm.DB, userID uint, page, limit int) (notes []Note, total int64, err error) {
	if page < 1 { page = 1 }
	if limit < 1 || limit > 100 { limit = 10 }
	var likes []NoteLike
	if err = db.Where("user_id = ?", userID).Find(&likes).Error; err != nil { return nil, 0, err }
	noteIDs := make([]uint, 0, len(likes))
	for _, l := range likes { noteIDs = append(noteIDs, l.NoteID) }
	if len(noteIDs) == 0 { return []Note{}, 0, nil }
	if err = db.Model(&Note{}).Where("id IN ?", noteIDs).Count(&total).Error; err != nil { return nil, 0, err }
	offset := (page - 1) * limit
	res := db.Where("id IN ?", noteIDs).Preload("User").Preload("Course").Limit(limit).Offset(offset).Order("created_at DESC").Find(&notes)
	return notes, total, res.Error
}

// SearchPublicNotes 搜索公开笔记（关键字、排序、分页）
func SearchPublicNotes(db *gorm.DB, keyword, sortBy, order string, page, pageSize int) (notes []Note, total int64, err error) {
	if page < 1 { page = 1 }
	if pageSize < 1 || pageSize > 100 { pageSize = 10 }
	kw := strings.TrimSpace(keyword)
	pattern := "%" + kw + "%"
	q := db.Model(&Note{}).Where("status = ?", "public")
	if kw != "" {
		q = q.Where("title LIKE ? OR description LIKE ?", pattern, pattern)
	}
	if err = q.Count(&total).Error; err != nil { return nil, 0, err }
	orderUpper := strings.ToUpper(order)
	if orderUpper != "ASC" { orderUpper = "DESC" }
	orderClause := "created_at DESC"
	switch sortBy {
	case "like_count":
		orderClause = "like_count " + orderUpper
	case "view_count":
		orderClause = "view_count " + orderUpper
	case "created_at":
		orderClause = "created_at " + orderUpper
	}
	offset := (page - 1) * pageSize
	res := q.Preload("User").Preload("Course").Offset(offset).Limit(pageSize).Order(orderClause).Find(&notes)
	return notes, total, res.Error
}

// GetPopularPublicNotes 获取按点赞数排序的热门公开笔记
func GetPopularPublicNotes(db *gorm.DB, limit int) (notes []Note, err error) {
	if limit <= 0 { limit = 10 }
	res := db.Model(&Note{}).Where("status = ?", "public").Preload("User").Preload("Course").Order("like_count DESC, view_count DESC").Limit(limit).Find(&notes)
	return notes, res.Error
}


// GetFavoritesByUser 获取用户的收藏记录
func GetFavoritesByUser(db *gorm.DB, userID uint) ([]Favorite, error) {
	var favorites []Favorite
	res := db.Where("user_id = ?", userID).Find(&favorites)
	return favorites, res.Error
}

// GetNoteLikesByUser 获取用户的笔记点赞记录
func GetNoteLikesByUser(db *gorm.DB, userID uint) ([]NoteLike, error) {
	var likes []NoteLike
	res := db.Where("user_id = ?", userID).Find(&likes)
	return likes, res.Error
}