package api

import (
	"fmt"
	"net/http"
	"time"

	"student_shared/app/model"
	"student_shared/app/utils/database"
	"student_shared/app/utils/jwt"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	req "student_shared/app/model/req"
	resp "student_shared/app/model/resp"
)

// RegisterUser 用户注册
func RegisterUser(c *gin.Context) {
	var req req.UserRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数无效"})
		return
	}

	// 检查用户名是否已存在
	existsUsername, err := model.UserExistsByUsername(database.DB, req.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "检查用户名失败"})
		return
	}
	if existsUsername {
		c.JSON(http.StatusConflict, gin.H{"error": "用户名已存在"})
		return
	}

	// 检查邮箱是否已存在
	existsEmail, err := model.UserExistsByEmail(database.DB, req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "检查邮箱失败"})
		return
	}
	if existsEmail {
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

	if err := model.CreateUser(database.DB, &user); err != nil {
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
		"user": resp.UserProfileResponse{
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
	var req req.UserLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数无效"})
		return
	}

	// 查找用户
	user, err := model.GetUserByUsername(database.DB, req.Username)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	// 更新最后登录时间
	loginTime := time.Now()
	if err := model.UpdateLastLogin(database.DB, user.ID, loginTime); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新登录时间失败"})
		return
	}
	user.LastLogin = loginTime

	// 生成JWT令牌
	token, err := jwt.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成令牌失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "登录成功",
		"token":   token,
		"user": resp.UserProfileResponse{
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
	user, err := model.GetUserByID(database.DB, userID.(uint))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	c.JSON(http.StatusOK, resp.UserProfileResponse{
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
	user, err := model.GetUserByID(database.DB, userID.(uint))
	if err != nil {
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

	updates := map[string]interface{}{}
	if updateData.Nickname != nil { updates["nickname"] = *updateData.Nickname }
	if updateData.Avatar != nil { updates["avatar"] = *updateData.Avatar }
	if updateData.School != nil { updates["school"] = *updateData.School }
	if updateData.Department != nil { updates["department"] = *updateData.Department }
	if updateData.Major != nil { updates["major"] = *updateData.Major }
	if updateData.Introduction != nil { updates["introduction"] = *updateData.Introduction }

	updated, err := model.UpdateUserFields(database.DB, user.ID, updates)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新用户信息失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "用户信息更新成功",
		"user": resp.UserProfileResponse{
			ID:           updated.ID,
			Username:     updated.Username,
			Email:        updated.Email,
			Nickname:     updated.Nickname,
			Avatar:       updated.Avatar,
			School:       updated.School,
			Department:   updated.Department,
			Major:        updated.Major,
			Introduction: updated.Introduction,
			Role:         updated.Role,
			LastLogin:    updated.LastLogin,
			CreatedAt:    updated.CreatedAt,
		},
	})
}
