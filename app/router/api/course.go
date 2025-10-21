package api

import (
	"net/http"
	"strconv"
	"student_shared/app/model"
	"student_shared/app/utils/database"
	"time"

	"github.com/gin-gonic/gin"
)

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

func ListCourses(c *gin.Context) {
	// 获取查询参数
	school := c.Query("school")
	department := c.Query("department")
	semester := c.Query("semester")
	search := c.Query("search") // 添加搜索参数
	status := c.Query("status") // 添加状态参数

	// 分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	// 构建查询
	query := database.DB.Model(&model.Course{})
	
	// 如果指定了状态，则按状态筛选；否则默认只显示活跃课程
	if status != "" {
		query = query.Where("status = ?", status)
	} else {
		// 对于普通用户，默认只显示活跃课程
		// 对于管理员，如果没有指定状态参数，显示所有课程
		userRole, exists := c.Get("role")
		if !exists || userRole != "admin" {
			query = query.Where("status = ?", "active")
		}
	}

	// 添加筛选条件
	if school != "" {
		query = query.Where("school LIKE ?", "%"+school+"%")
	}
	if department != "" {
		query = query.Where("department LIKE ?", "%"+department+"%")
	}
	if semester != "" {
		query = query.Where("semester = ?", semester)
	}
	if search != "" {
		// 在课程名称、描述、教师名称中搜索
		searchPattern := "%" + search + "%"
		query = query.Where("name LIKE ? OR description LIKE ? OR teacher LIKE ?", searchPattern, searchPattern, searchPattern)
	}

	// 获取总数
	var total int64
	query.Count(&total)

	// 分页查询
	var courses []model.Course
	offset := (page - 1) * pageSize
	result := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&courses)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询课程失败"})
		return
	}

	// 转换为响应格式
	var courseResponses []CourseResponse
	for _, course := range courses {
		// 统计学生数量
		var studentCount int64
		database.DB.Model(&model.UserCourse{}).Where("course_id = ?", course.ID).Count(&studentCount)

		// 统计笔记数量
		var noteCount int64
		database.DB.Model(&model.Note{}).Where("course_id = ?", course.ID).Count(&noteCount)

		courseResponse := CourseResponse{
			ID:           course.ID,
			Code:         course.Code,
			Name:         course.Name,
			Description:  course.Description,
			School:       course.School,
			Department:   course.Department,
			Teacher:      course.Teacher,
			Semester:     course.Semester,
			Cover:        course.Cover,
			Status:       course.Status,
			CreatedAt:    course.CreatedAt,
			UpdatedAt:    course.UpdatedAt,
			StudentCount: studentCount,
			NoteCount:    noteCount,
		}

		// 检查用户是否已加入课程
		userID, exists := c.Get("userID")
		if exists {
			var userCourse model.UserCourse
			err := database.DB.Where("user_id = ? AND course_id = ?", userID, course.ID).First(&userCourse).Error
			courseResponse.IsJoined = err == nil
		}

		courseResponses = append(courseResponses, courseResponse)
	}

	// 计算分页信息
	totalPages := (int(total) + pageSize - 1) / pageSize

	c.JSON(http.StatusOK, gin.H{
		"courses":     courseResponses,
		"total":       total,
		"page":        page,
		"page_size":   pageSize,
		"total_pages": totalPages,
	})
}

func GetCourse(c *gin.Context) {
	// 获取课程ID
	id := c.Param("id")

	// 查询课程
	var course model.Course
	result := database.DB.First(&course, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "课程不存在"})
		return
	}

	// 检查用户是否已加入课程
	isJoined := false
	userID, exists := c.Get("userID")
	if exists {
		var userCourse model.UserCourse
		err := database.DB.Where("user_id = ? AND course_id = ?", userID, course.ID).First(&userCourse).Error
		isJoined = err == nil
	}

	// 统计学生数量
	var studentCount int64
	database.DB.Model(&model.UserCourse{}).Where("course_id = ?", course.ID).Count(&studentCount)

	// 统计笔记数量
	var noteCount int64
	database.DB.Model(&model.Note{}).Where("course_id = ?", course.ID).Count(&noteCount)

	courseResponse := CourseResponse{
		ID:           course.ID,
		Code:         course.Code,
		Name:         course.Name,
		Description:  course.Description,
		School:       course.School,
		Department:   course.Department,
		Teacher:      course.Teacher,
		Semester:     course.Semester,
		Cover:        course.Cover,
		Status:       course.Status,
		CreatedAt:    course.CreatedAt,
		UpdatedAt:    course.UpdatedAt,
		IsJoined:     isJoined,
		StudentCount: studentCount,
		NoteCount:    noteCount,
	}

	c.JSON(http.StatusOK, courseResponse)
}

func CreateCourse(c *gin.Context) {
	// 获取用户ID和角色
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	role, _ := c.Get("role")
	// 只有教师和管理员可以创建课程
	if role.(string) != "teacher" && role.(string) != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "只有教师和管理员可以创建课程"})
		return
	}

	// 绑定请求数据
	var req CourseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数无效"})
		return
	}

	// 检查课程代码是否已存在
	var existingCourse model.Course
	result := database.DB.Where("code = ?", req.Code).First(&existingCourse)
	if result.RowsAffected > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "课程代码已存在"})
		return
	}

	// 创建课程
	course := model.Course{
		Code:        req.Code,
		Name:        req.Name,
		Description: req.Description,
		School:      req.School,
		Department:  req.Department,
		Teacher:     req.Teacher,
		Semester:    req.Semester,
		Cover:       req.Cover,
		Status:      "active",
	}

	if req.Status != "" {
		course.Status = req.Status
	}

	result = database.DB.Create(&course)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建课程失败"})
		return
	}

	// 如果创建者是教师，自动加入课程并设置为教师角色
	if role.(string) == "teacher" {
		userCourse := model.UserCourse{
			UserID:   userID.(uint),
			CourseID: course.ID,
			Role:     "teacher",
			JoinedAt: time.Now(),
		}
		database.DB.Create(&userCourse)
	}

	c.JSON(http.StatusCreated, CourseResponse{
		ID:          course.ID,
		Code:        course.Code,
		Name:        course.Name,
		Description: course.Description,
		School:      course.School,
		Department:  course.Department,
		Teacher:     course.Teacher,
		Semester:    course.Semester,
		Cover:       course.Cover,
		Status:      course.Status,
		CreatedAt:   course.CreatedAt,
		UpdatedAt:   course.UpdatedAt,
	})
}

func UpdateCourse(c *gin.Context) {
	// 获取用户ID和角色
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	role, _ := c.Get("role")

	// 获取课程ID
	id := c.Param("id")

	// 查询课程
	var course model.Course
	result := database.DB.First(&course, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "课程不存在"})
		return
	}

	// 检查权限
	canUpdate := false
	if role.(string) == "admin" {
		canUpdate = true
	} else if role.(string) == "teacher" {
		// 检查是否是课程的教师
		var userCourse model.UserCourse
		err := database.DB.Where("user_id = ? AND course_id = ? AND role = ?", userID, course.ID, "teacher").First(&userCourse).Error
		canUpdate = err == nil
	}

	if !canUpdate {
		c.JSON(http.StatusForbidden, gin.H{"error": "没有权限修改此课程"})
		return
	}

	// 绑定请求数据
	var req CourseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数无效"})
		return
	}

	// 检查课程代码是否与其他课程冲突
	if req.Code != course.Code {
		var existingCourse model.Course
		result := database.DB.Where("code = ? AND id != ?", req.Code, course.ID).First(&existingCourse)
		if result.RowsAffected > 0 {
			c.JSON(http.StatusConflict, gin.H{"error": "课程代码已存在"})
			return
		}
	}

	// 更新课程信息 - 使用Select避免更新CreatedAt
	updateData := map[string]interface{}{
		"code":        req.Code,
		"name":        req.Name,
		"description": req.Description,
		"school":      req.School,
		"department":  req.Department,
		"teacher":     req.Teacher,
		"semester":    req.Semester,
		"cover":       req.Cover,
	}
	
	if req.Status != "" {
		updateData["status"] = req.Status
	}

	result = database.DB.Model(&course).Updates(updateData)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新课程失败"})
		return
	}

	// 重新查询更新后的课程信息
	var updatedCourse model.Course
	if err := database.DB.First(&updatedCourse, course.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取更新后的课程信息失败"})
		return
	}

	c.JSON(http.StatusOK, CourseResponse{
		ID:          updatedCourse.ID,
		Code:        updatedCourse.Code,
		Name:        updatedCourse.Name,
		Description: updatedCourse.Description,
		School:      updatedCourse.School,
		Department:  updatedCourse.Department,
		Teacher:     updatedCourse.Teacher,
		Semester:    updatedCourse.Semester,
		Cover:       updatedCourse.Cover,
		Status:      updatedCourse.Status,
		CreatedAt:   updatedCourse.CreatedAt,
		UpdatedAt:   updatedCourse.UpdatedAt,
	})
}

func DeleteCourse(c *gin.Context) {
	// 获取用户ID和角色
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	role, _ := c.Get("role")

	// 获取课程ID
	id := c.Param("id")

	// 查询课程
	var course model.Course
	result := database.DB.First(&course, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "课程不存在"})
		return
	}

	// 检查权限
	canDelete := false
	if role.(string) == "admin" {
		canDelete = true
	} else if role.(string) == "teacher" {
		// 检查是否是课程的教师
		var userCourse model.UserCourse
		err := database.DB.Where("user_id = ? AND course_id = ? AND role = ?", userID, course.ID, "teacher").First(&userCourse).Error
		canDelete = err == nil
	}

	if !canDelete {
		c.JSON(http.StatusForbidden, gin.H{"error": "没有权限删除此课程"})
		return
	}

	// 软删除课程（设置状态为inactive）
	course.Status = "inactive"
	result = database.DB.Save(&course)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除课程失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "课程删除成功"})
}

func JoinCourse(c *gin.Context) {
	// 获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	// 获取课程ID
	id := c.Param("id")

	// 查询课程
	var course model.Course
	result := database.DB.First(&course, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "课程不存在"})
		return
	}

	// 检查课程状态
	if course.Status != "active" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "课程不可用"})
		return
	}

	// 检查用户是否已加入课程
	var existingUserCourse model.UserCourse
	result = database.DB.Where("user_id = ? AND course_id = ?", userID, course.ID).First(&existingUserCourse)
	if result.RowsAffected > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "已加入该课程"})
		return
	}

	// 创建用户课程关联
	userCourse := model.UserCourse{
		UserID:   userID.(uint),
		CourseID: course.ID,
		Role:     "student",
		JoinedAt: time.Now(),
	}

	result = database.DB.Create(&userCourse)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "加入课程失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "加入课程成功"})
}

func SearchCourses(c *gin.Context) {
	// 获取搜索关键词
	keyword := c.Query("keyword")
	if keyword == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "搜索关键词不能为空"})
		return
	}

	// 分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	// 构建搜索查询
	searchPattern := "%" + keyword + "%"
	query := database.DB.Model(&model.Course{}).Where(
		"status = ? AND (name LIKE ? OR description LIKE ? OR teacher LIKE ? OR code LIKE ?)",
		"active", searchPattern, searchPattern, searchPattern, searchPattern,
	)

	// 获取总数
	var total int64
	query.Count(&total)

	// 分页查询
	var courses []model.Course
	offset := (page - 1) * pageSize
	result := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&courses)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "搜索课程失败"})
		return
	}

	// 转换为响应格式
	var courseResponses []CourseResponse
	for _, course := range courses {
		// 统计学生数量
		var studentCount int64
		database.DB.Model(&model.UserCourse{}).Where("course_id = ?", course.ID).Count(&studentCount)

		// 统计笔记数量
		var noteCount int64
		database.DB.Model(&model.Note{}).Where("course_id = ?", course.ID).Count(&noteCount)

		courseResponse := CourseResponse{
			ID:           course.ID,
			Code:         course.Code,
			Name:         course.Name,
			Description:  course.Description,
			School:       course.School,
			Department:   course.Department,
			Teacher:      course.Teacher,
			Semester:     course.Semester,
			Cover:        course.Cover,
			Status:       course.Status,
			CreatedAt:    course.CreatedAt,
			UpdatedAt:    course.UpdatedAt,
			StudentCount: studentCount,
			NoteCount:    noteCount,
		}

		// 检查用户是否已加入课程
		userID, exists := c.Get("userID")
		if exists {
			var userCourse model.UserCourse
			err := database.DB.Where("user_id = ? AND course_id = ?", userID, course.ID).First(&userCourse).Error
			courseResponse.IsJoined = err == nil
		}

		courseResponses = append(courseResponses, courseResponse)
	}

	// 计算分页信息
	totalPages := (int(total) + pageSize - 1) / pageSize

	c.JSON(http.StatusOK, gin.H{
		"courses":     courseResponses,
		"total":       total,
		"page":        page,
		"page_size":   pageSize,
		"total_pages": totalPages,
		"keyword":     keyword,
	})
}

func GetLatestCourses(c *gin.Context) {
	// 构建查询，只获取活跃课程
	query := database.DB.Model(&model.Course{}).Where("status = ?", "active")

	// 获取前3个最新创建的课程
	var courses []model.Course
	result := query.Order("created_at DESC").Limit(3).Find(&courses)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取最新课程失败"})
		return
	}

	// 转换为响应格式
	var courseResponses []CourseResponse
	for _, course := range courses {
		// 统计学生数量
		var studentCount int64
		database.DB.Model(&model.UserCourse{}).Where("course_id = ?", course.ID).Count(&studentCount)

		// 统计笔记数量
		var noteCount int64
		database.DB.Model(&model.Note{}).Where("course_id = ?", course.ID).Count(&noteCount)

		courseResponse := CourseResponse{
			ID:           course.ID,
			Code:         course.Code,
			Name:         course.Name,
			Description:  course.Description,
			School:       course.School,
			Department:   course.Department,
			Teacher:      course.Teacher,
			Semester:     course.Semester,
			Cover:        course.Cover,
			Status:       course.Status,
			CreatedAt:    course.CreatedAt,
			UpdatedAt:    course.UpdatedAt,
			StudentCount: studentCount,
			NoteCount:    noteCount,
		}

		// 检查用户是否已加入课程
		userID, exists := c.Get("userID")
		if exists {
			var userCourse model.UserCourse
			err := database.DB.Where("user_id = ? AND course_id = ?", userID, course.ID).First(&userCourse).Error
			courseResponse.IsJoined = err == nil
		}

		courseResponses = append(courseResponses, courseResponse)
	}

	c.JSON(http.StatusOK, courseResponses)
}

// GetMyCourses 获取用户加入的课程列表
func GetMyCourses(c *gin.Context) {
	// 获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	// 分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	// 查询用户加入的课程
	var userCourses []model.UserCourse
	offset := (page - 1) * pageSize
	result := database.DB.Where("user_id = ?", userID).Offset(offset).Limit(pageSize).Order("joined_at DESC").Find(&userCourses)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询用户课程失败"})
		return
	}

	// 获取总数
	var total int64
	database.DB.Model(&model.UserCourse{}).Where("user_id = ?", userID).Count(&total)

	// 获取课程详细信息
	var myCourseResponses []MyCourseResponse
	for _, userCourse := range userCourses {
		var course model.Course
		result := database.DB.First(&course, userCourse.CourseID)
		if result.Error != nil {
			continue // 跳过不存在的课程
		}

		// 统计学生数量
		var studentCount int64
		database.DB.Model(&model.UserCourse{}).Where("course_id = ?", course.ID).Count(&studentCount)

		// 统计笔记数量
		var noteCount int64
		database.DB.Model(&model.Note{}).Where("course_id = ?", course.ID).Count(&noteCount)

		myCourseResponse := MyCourseResponse{
			ID:           course.ID,
			Code:         course.Code,
			Name:         course.Name,
			Description:  course.Description,
			School:       course.School,
			Department:   course.Department,
			Teacher:      course.Teacher,
			Semester:     course.Semester,
			Cover:        course.Cover,
			Status:       course.Status,
			Role:         userCourse.Role,
			JoinedAt:     userCourse.JoinedAt,
			CreatedAt:    course.CreatedAt,
			UpdatedAt:    course.UpdatedAt,
			StudentCount: studentCount,
			NoteCount:    noteCount,
		}

		myCourseResponses = append(myCourseResponses, myCourseResponse)
	}

	// 计算分页信息
	totalPages := (int(total) + pageSize - 1) / pageSize

	c.JSON(http.StatusOK, gin.H{
		"courses":     myCourseResponses,
		"total":       total,
		"page":        page,
		"page_size":   pageSize,
		"total_pages": totalPages,
	})
}
