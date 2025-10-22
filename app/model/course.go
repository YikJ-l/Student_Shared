package model

import (
	"strings"
	"time"
	"fmt"
	"errors"

	"gorm.io/gorm"
)

// Course 课程模型
type Course struct {
	BaseModel
	Code        string `gorm:"size:50;not null;uniqueIndex" json:"code"`
	Name        string `gorm:"size:100;not null" json:"name"`
	Description string `gorm:"size:1000" json:"description"`
	School      string `gorm:"size:100" json:"school"`
	Department  string `gorm:"size:100" json:"department"`
	Teacher     string `gorm:"size:100" json:"teacher"`
	Semester    string `gorm:"size:50" json:"semester"`
	Cover       string `gorm:"size:255" json:"cover"`
	Status      string `gorm:"size:20;default:'active'" json:"status"` // active, inactive
	Embedding   string `gorm:"type:text" json:"-"` // 保存嵌入向量的JSON字符串
	
	// 关联
	Users []User `gorm:"many2many:user_courses;" json:"-"`
}

// UserCourse 用户-课程关联表
type UserCourse struct {
	UserID   uint `gorm:"primaryKey"`
	CourseID uint `gorm:"primaryKey"`
	Role     string `gorm:"size:20;default:'student'"` // student, teacher
	JoinedAt time.Time
}

// CourseListFilter 课程列表查询过滤条件
type CourseListFilter struct {
	School     string
	Department string
	Semester   string
	Status     string
	Search     string
}

// BuildCourseQuery 根据过滤条件构建课程查询
func BuildCourseQuery(db *gorm.DB, f CourseListFilter) *gorm.DB {
	query := db.Model(&Course{})
	if strings.TrimSpace(f.School) != "" {
		query = query.Where("school = ?", f.School)
	}
	if strings.TrimSpace(f.Department) != "" {
		query = query.Where("department = ?", f.Department)
	}
	if strings.TrimSpace(f.Semester) != "" {
		query = query.Where("semester = ?", f.Semester)
	}
	if strings.TrimSpace(f.Status) != "" {
		query = query.Where("status = ?", f.Status)
	}
	if s := strings.TrimSpace(f.Search); s != "" {
		like := "%" + s + "%"
		query = query.Where("name LIKE ? OR code LIKE ? OR teacher LIKE ?", like, like, like)
	}
	return query
}

// ListCourses 按分页返回课程及总数
func ListCourses(db *gorm.DB, f CourseListFilter, page, pageSize int) (courses []Course, total int64, err error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}
	query := BuildCourseQuery(db, f)
	query.Count(&total)
	offset := (page - 1) * pageSize
	result := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&courses)
	return courses, total, result.Error
}

// GetLatestActiveCourses 获取最新活跃课程（按创建时间倒序）
func GetLatestActiveCourses(db *gorm.DB, limit int) (courses []Course, err error) {
	if limit <= 0 {
		limit = 3
	}
	result := db.Model(&Course{}).Where("status = ?", "active").Order("created_at DESC").Limit(limit).Find(&courses)
	return courses, result.Error
}

// GetCourseByID 根据ID获取课程
func GetCourseByID(db *gorm.DB, id uint) (Course, error) {
	var course Course
	result := db.First(&course, id)
	return course, result.Error
}

// SoftDeleteCourse 软删除课程（设置为 inactive）
func SoftDeleteCourse(db *gorm.DB, id uint) error {
	var course Course
	result := db.First(&course, id)
	if result.Error != nil {
		return result.Error
	}
	course.Status = "inactive"
	return db.Save(&course).Error
}

// JoinCourse 用户加入课程（含状态与重复校验）
func JoinCourse(db *gorm.DB, userID uint, courseID uint) error {
	var course Course
	result := db.First(&course, courseID)
	if result.Error != nil {
		return result.Error
	}
	if course.Status != "active" {
		return fmt.Errorf("课程不可用")
	}
	var existing UserCourse
	result = db.Where("user_id = ? AND course_id = ?", userID, courseID).First(&existing)
	if result.RowsAffected > 0 {
		return fmt.Errorf("已加入该课程")
	}
	userCourse := UserCourse{UserID: userID, CourseID: courseID, Role: "student", JoinedAt: time.Now()}
	return db.Create(&userCourse).Error
}

// ListUserCourses 获取用户加入课程（分页）
func ListUserCourses(db *gorm.DB, userID uint, page, pageSize int) (userCourses []UserCourse, total int64, err error) {
	if page < 1 { page = 1 }
	if pageSize < 1 || pageSize > 100 { pageSize = 10 }
	offset := (page - 1) * pageSize
	result := db.Where("user_id = ?", userID).Offset(offset).Limit(pageSize).Order("joined_at DESC").Find(&userCourses)
	if result.Error != nil { return nil, 0, result.Error }
	db.Model(&UserCourse{}).Where("user_id = ?", userID).Count(&total)
	return userCourses, total, nil
}

// SearchCourses 搜索课程（关键字、状态、排序、分页）
func SearchCourses(db *gorm.DB, keyword, status, sortBy, order string, page, pageSize int) (courses []Course, total int64, err error) {
	if page < 1 { page = 1 }
	if pageSize < 1 || pageSize > 100 { pageSize = 10 }
	searchPattern := "%" + strings.TrimSpace(keyword) + "%"
	query := db.Model(&Course{})
	if strings.TrimSpace(status) != "" {
		query = query.Where("status = ?", status)
	} else {
		query = query.Where("status = ?", "active")
	}
	query = query.Where("(name LIKE ? OR description LIKE ? OR teacher LIKE ? OR code LIKE ?)", searchPattern, searchPattern, searchPattern, searchPattern)
	// 排序
	orderUpper := strings.ToUpper(order)
	if orderUpper != "ASC" { orderUpper = "DESC" }
	orderClause := "created_at DESC"
	if sortBy == "created_at" {
		orderClause = "created_at " + orderUpper
	}
	// 统计与分页
	query.Count(&total)
	offset := (page - 1) * pageSize
	result := query.Offset(offset).Limit(pageSize).Order(orderClause).Find(&courses)
	return courses, total, result.Error
}

// CourseStats 课程统计数据
type CourseStats struct {
	StudentCount int64
	NoteCount    int64
}

// IsUserJoinedCourse 判断用户是否已加入课程
func IsUserJoinedCourse(db *gorm.DB, userID uint, courseID uint) (bool, error) {
	var uc UserCourse
	result := db.Where("user_id = ? AND course_id = ?", userID, courseID).First(&uc)
	if result.Error == nil {
		return true, nil
	}
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return false, nil
	}
	return false, result.Error
}

// HasUserCourseRole 判断用户在课程中的角色是否匹配
func HasUserCourseRole(db *gorm.DB, userID uint, courseID uint, role string) (bool, error) {
	var uc UserCourse
	result := db.Where("user_id = ? AND course_id = ? AND role = ?", userID, courseID, role).First(&uc)
	if result.Error == nil {
		return true, nil
	}
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return false, nil
	}
	return false, result.Error
}

// GetCourseStats 获取课程的学生数与笔记数
func GetCourseStats(db *gorm.DB, courseID uint) (CourseStats, error) {
	var students int64
	if err := db.Model(&UserCourse{}).Where("course_id = ?", courseID).Count(&students).Error; err != nil {
		return CourseStats{}, err
	}
	var notes int64
	if err := db.Model(&Note{}).Where("course_id = ?", courseID).Count(&notes).Error; err != nil {
		return CourseStats{}, err
	}
	return CourseStats{StudentCount: students, NoteCount: notes}, nil
}

// CreateCourse 创建课程并校验代码唯一性
func CreateCourse(db *gorm.DB, course *Course) error {
	var existing Course
	res := db.Where("code = ?", course.Code).First(&existing)
	if res.Error == nil {
		return fmt.Errorf("课程代码已存在")
	}
	if res.Error != nil && !errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return res.Error
	}
	return db.Create(course).Error
}

// IsCourseCodeTaken 检查课程代码是否被其他课程占用
func IsCourseCodeTaken(db *gorm.DB, code string, excludeID uint) (bool, error) {
	var c Course
	res := db.Where("code = ? AND id != ?", code, excludeID).First(&c)
	if res.Error == nil {
		return true, nil
	}
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return false, nil
	}
	return false, res.Error
}

// UpdateCourse 更新课程信息
func UpdateCourse(db *gorm.DB, courseID uint, updates map[string]interface{}) error {
	return db.Model(&Course{}).Where("id = ?", courseID).Updates(updates).Error
}

// AddUserToCourseWithRole 加入课程并指定角色（若已加入则忽略）
func AddUserToCourseWithRole(db *gorm.DB, userID, courseID uint, role string) error {
	var uc UserCourse
	res := db.Where("user_id = ? AND course_id = ?", userID, courseID).First(&uc)
	if res.Error == nil {
		// 已加入，若角色不同可根据需求更新，这里保持不变
		return nil
	}
	if res.Error != nil && !errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return res.Error
	}
	return db.Create(&UserCourse{UserID: userID, CourseID: courseID, Role: role, JoinedAt: time.Now()}).Error
}