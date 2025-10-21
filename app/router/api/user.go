package api

import (
	"fmt"
	"net/http"
	"student_shared/app/model"
	"student_shared/app/utils/database"
	"student_shared/app/utils/jwt"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// UserRegisterRequest 用户注册请求
type UserRegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6"`
	Email    string `json:"email" binding:"required,email"`
	Nickname string `json:"nickname"`
	School   string `json:"school"`
}

// UserLoginRequest 用户登录请求
type UserLoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// UserProfileResponse 用户资料响应
type UserProfileResponse struct {
	ID           uint      `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	Nickname     string    `json:"nickname"`
	Avatar       string    `json:"avatar"`
	School       string    `json:"school"`
	Department   string    `json:"department"`
	Major        string    `json:"major"`
	Introduction string    `json:"introduction"`
	Role         string    `json:"role"`
	LastLogin    time.Time `json:"last_login"`
	CreatedAt    time.Time `json:"created_at"`
}

// RegisterUser 用户注册
func RegisterUser(c *gin.Context) {
	var req UserRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数无效"})
		return
	}

	// 检查用户名是否已存在
	var existingUser model.User
	result := database.DB.Where("username = ?", req.Username).First(&existingUser)
	if result.RowsAffected > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "用户名已存在"})
		return
	}

	// 检查邮箱是否已存在
	result = database.DB.Where("email = ?", req.Email).First(&existingUser)
	if result.RowsAffected > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "邮箱已被注册"})
		return
	}

	// 密码加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "密码加密失败"})
		return
	}

	// 创建用户
	user := model.User{
		Username:  req.Username,
		Password:  string(hashedPassword),
		Email:     req.Email,
		Nickname:  req.Nickname,
		School:    req.School,
		LastLogin: time.Now(),
		Role:      "student",
	}

	result = database.DB.Create(&user)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建用户失败"})
		return
	}

	// 生成JWT令牌
	token, err := jwt.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成令牌失败"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "注册成功",
		"token":   token,
		"user": UserProfileResponse{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			Nickname:  user.Nickname,
			School:    user.School,
			Role:      user.Role,
			CreatedAt: user.CreatedAt,
		},
	})
}

// LoginUser 用户登录
func LoginUser(c *gin.Context) {
	var req UserLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数无效"})
		return
	}

	// 查找用户
	var user model.User
	result := database.DB.Where("username = ?", req.Username).First(&user)
	if result.RowsAffected == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	// 验证密码
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	// 更新最后登录时间
	user.LastLogin = time.Now()
	database.DB.Save(&user)

	// 生成JWT令牌
	token, err := jwt.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成令牌失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "登录成功",
		"token":   token,
		"user": UserProfileResponse{
			ID:           user.ID,
			Username:     user.Username,
			Email:        user.Email,
			Nickname:     user.Nickname,
			Avatar:       user.Avatar,
			School:       user.School,
			Department:   user.Department,
			Major:        user.Major,
			Introduction: user.Introduction,
			Role:         user.Role,
			LastLogin:    user.LastLogin,
			CreatedAt:    user.CreatedAt,
		},
	})
}

// GetUserProfile 获取用户资料
func GetUserProfile(c *gin.Context) {
	// 从JWT中获取用户ID
	userID, exists := c.Get("userID")
	fmt.Println(userID)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	// 查询用户信息
	var user model.User
	result := database.DB.First(&user, userID)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	c.JSON(http.StatusOK, UserProfileResponse{
		ID:           user.ID,
		Username:     user.Username,
		Email:        user.Email,
		Nickname:     user.Nickname,
		Avatar:       user.Avatar,
		School:       user.School,
		Department:   user.Department,
		Major:        user.Major,
		Introduction: user.Introduction,
		Role:         user.Role,
		LastLogin:    user.LastLogin,
		CreatedAt:    user.CreatedAt,
	})
}

// UpdateUserProfile 更新用户资料
func UpdateUserProfile(c *gin.Context) {
	// 从JWT中获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	// 查询用户信息
	var user model.User
	result := database.DB.First(&user, userID)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	// 绑定请求数据
	var updateData struct {
		Nickname     *string `json:"nickname"`
		Avatar       *string `json:"avatar"`
		School       *string `json:"school"`
		Department   *string `json:"department"`
		Major        *string `json:"major"`
		Introduction *string `json:"introduction"`
	}

	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数无效"})
		return
	}

	// 更新用户信息 - 使用指针判断字段是否被提供
	if updateData.Nickname != nil {
		user.Nickname = *updateData.Nickname
	}
	if updateData.Avatar != nil {
		user.Avatar = *updateData.Avatar
	}
	if updateData.School != nil {
		user.School = *updateData.School
	}
	if updateData.Department != nil {
		user.Department = *updateData.Department
	}
	if updateData.Major != nil {
		user.Major = *updateData.Major
	}
	if updateData.Introduction != nil {
		user.Introduction = *updateData.Introduction
	}

	result = database.DB.Save(&user)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新用户信息失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "更新成功",
		"user": UserProfileResponse{
			ID:           user.ID,
			Username:     user.Username,
			Email:        user.Email,
			Nickname:     user.Nickname,
			Avatar:       user.Avatar,
			School:       user.School,
			Department:   user.Department,
			Major:        user.Major,
			Introduction: user.Introduction,
			Role:         user.Role,
			LastLogin:    user.LastLogin,
			CreatedAt:    user.CreatedAt,
		},
	})
}
