package model

import (
	"time"

	"gorm.io/gorm"
)

// BaseModel 基础模型，包含共有字段
type BaseModel struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// User 用户模型
type User struct {
	BaseModel
	Username     string    `gorm:"size:50;not null;uniqueIndex" json:"username"`
	Password     string    `gorm:"size:100;not null" json:"-"`
	Email        string    `gorm:"size:100;not null;uniqueIndex" json:"email"`
	Nickname     string    `gorm:"size:50" json:"nickname"`
	Avatar       string    `gorm:"size:255" json:"avatar"`
	School       string    `gorm:"size:100" json:"school"`
	Department   string    `gorm:"size:100" json:"department"`
	Major        string    `gorm:"size:100" json:"major"`
	Introduction string    `gorm:"size:500" json:"introduction"`
	LastLogin    time.Time `json:"last_login"`
	Role         string    `gorm:"size:20;default:'student'" json:"role"` // student, teacher, admin
}

// --- 用户相关数据库操作封装 ---

// UserExistsByUsername 检查用户名是否存在
func UserExistsByUsername(db *gorm.DB, username string) (bool, error) {
	var u User
	res := db.Where("username = ?", username).First(&u)
	if res.Error != nil && res.Error != gorm.ErrRecordNotFound {
		return false, res.Error
	}
	return res.RowsAffected > 0, nil
}

// UserExistsByEmail 检查邮箱是否存在
func UserExistsByEmail(db *gorm.DB, email string) (bool, error) {
	var u User
	res := db.Where("email = ?", email).First(&u)
	if res.Error != nil && res.Error != gorm.ErrRecordNotFound {
		return false, res.Error
	}
	return res.RowsAffected > 0, nil
}

// CreateUser 创建用户
func CreateUser(db *gorm.DB, user *User) error {
	return db.Create(user).Error
}

// GetUserByUsername 根据用户名获取用户
func GetUserByUsername(db *gorm.DB, username string) (*User, error) {
	var u User
	res := db.Where("username = ?", username).First(&u)
	if res.Error != nil {
		return nil, res.Error
	}
	return &u, nil
}

// GetUserByID 根据ID获取用户
func GetUserByID(db *gorm.DB, id uint) (*User, error) {
	var u User
	res := db.First(&u, id)
	if res.Error != nil {
		return nil, res.Error
	}
	return &u, nil
}

// UpdateLastLogin 更新最后登录时间
func UpdateLastLogin(db *gorm.DB, id uint, t time.Time) error {
	return db.Model(&User{}).Where("id = ?", id).Update("last_login", t).Error
}

// UpdateUserFields 部分更新用户字段
func UpdateUserFields(db *gorm.DB, id uint, updates map[string]interface{}) (*User, error) {
	var u User
	if err := db.First(&u, id).Error; err != nil {
		return nil, err
	}
	if err := db.Model(&u).Updates(updates).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

// -------- 管理员相关封装 --------

// AdminUserListFilter 管理员用户列表筛选
type AdminUserListFilter struct {
	Search string
	Role   string
}

// ListUsersAdmin 管理员分页查询用户
func ListUsersAdmin(db *gorm.DB, filter AdminUserListFilter, page, limit int) ([]User, int64, error) {
	if page < 1 { page = 1 }
	if limit < 1 || limit > 100 { limit = 10 }
	offset := (page - 1) * limit

	query := db.Model(&User{})
	if filter.Search != "" {
		like := "%" + filter.Search + "%"
		query = query.Where("username LIKE ? OR email LIKE ? OR nickname LIKE ?", like, like, like)
	}
	if filter.Role != "" {
		query = query.Where("role = ?", filter.Role)
	}
	var total int64
	if err := query.Count(&total).Error; err != nil { return nil, 0, err }
	var users []User
	if err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&users).Error; err != nil {
		return nil, 0, err
	}
	return users, total, nil
}

// IsEmailTakenByOther 检查邮箱是否被其他用户占用
func IsEmailTakenByOther(db *gorm.DB, email string, excludeID uint) (bool, error) {
	if email == "" { return false, nil }
	var u User
	res := db.Where("email = ? AND id != ?", email, excludeID).First(&u)
	if res.Error != nil && res.Error != gorm.ErrRecordNotFound {
		return false, res.Error
	}
	return res.RowsAffected > 0, nil
}

// AdminUserUpdateFields 管理员更新用户字段
type AdminUserUpdateFields struct {
	Nickname     string
	Email        string
	School       string
	Department   string
	Major        string
	Introduction string
}

// AdminUpdateUserFields 管理员更新用户信息（包含邮箱唯一性检查）
func AdminUpdateUserFields(db *gorm.DB, id uint, fields AdminUserUpdateFields) (*User, error) {
	var u User
	if err := db.First(&u, id).Error; err != nil { return nil, err }
	if fields.Email != "" {
		taken, err := IsEmailTakenByOther(db, fields.Email, id)
		if err != nil { return nil, err }
		if taken { return nil, gorm.ErrInvalidData }
	}
	updates := map[string]interface{}{}
	if fields.Nickname != "" { updates["nickname"] = fields.Nickname }
	if fields.Email != "" { updates["email"] = fields.Email }
	if fields.School != "" { updates["school"] = fields.School }
	if fields.Department != "" { updates["department"] = fields.Department }
	if fields.Major != "" { updates["major"] = fields.Major }
	if fields.Introduction != "" { updates["introduction"] = fields.Introduction }
	if len(updates) == 0 { return &u, nil }
	if err := db.Model(&u).Updates(updates).Error; err != nil { return nil, err }
	return &u, nil
}

// AdminUpdateUserRole 管理员更新用户角色
func AdminUpdateUserRole(db *gorm.DB, id uint, role string) (*User, error) {
	var u User
	if err := db.First(&u, id).Error; err != nil { return nil, err }
	u.Role = role
	if err := db.Save(&u).Error; err != nil { return nil, err }
	return &u, nil
}

// AdminDeleteUserSoft 管理员软删除用户
func AdminDeleteUserSoft(db *gorm.DB, id uint) error {
	var u User
	if err := db.First(&u, id).Error; err != nil { return err }
	return db.Delete(&u).Error
}