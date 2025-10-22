package api

import (
	"net/http"
	"strings"
	"student_shared/app/model"
	"student_shared/app/utils/database"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	// 引入统一的请求/响应包
	req "student_shared/app/model/req"
	resp "student_shared/app/model/resp"
)


func ListCourses(c *gin.Context) {
	// 使用统一的请求结构体
	var params req.ListCoursesRequest
	params.Page = 1
	params.PageSize = 10
	_ = c.ShouldBindJSON(&params)
	page := params.Page
	pageSize := params.PageSize
	if page < 1 { page = 1 }
	if pageSize < 1 || pageSize > 100 { pageSize = 10 }
	// 构建查询条件（迁移到 model 层）
	filter := model.CourseListFilter{
		School:     params.School,
		Department: params.Department,
		Semester:   params.Semester,
		Status:     params.Status,
		Search:     params.Search,
	}
	courses, total, err := model.ListCourses(database.DB, filter, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询课程失败"})
		return
	}

	// 转换为响应格式
	var courseResponses []resp.CourseResponse
	for _, course := range courses {
		// 使用 model 层封装获取课程统计
		stats, sErr := model.GetCourseStats(database.DB, course.ID)
		if sErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "查询课程统计失败"})
			return
		}
		courseResponse := resp.CourseResponse{
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
			StudentCount: stats.StudentCount,
			NoteCount:    stats.NoteCount,
		}

		// 检查用户是否已加入课程（使用封装）
		if uid, exists := c.Get("userID"); exists {
			joined, jErr := model.IsUserJoinedCourse(database.DB, uid.(uint), course.ID)
			if jErr != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "查询加入状态失败"})
				return
			}
			courseResponse.IsJoined = joined
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
	var params req.GetCourseRequest
	if err := c.ShouldBindJSON(&params); err != nil || params.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "缺少有效的课程ID"})
		return
	}

	course, err := model.GetCourseByID(database.DB, params.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "课程不存在"})
		return
	}

	isJoined := false
	if uid, exists := c.Get("userID"); exists {
		joined, jErr := model.IsUserJoinedCourse(database.DB, uid.(uint), course.ID)
		if jErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "查询加入状态失败"})
			return
		}
		isJoined = joined
	}

	stats, sErr := model.GetCourseStats(database.DB, course.ID)
	if sErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询课程统计失败"})
		return
	}

	courseResponse := resp.CourseResponse{
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
		StudentCount: stats.StudentCount,
		NoteCount:    stats.NoteCount,
	}

	c.JSON(http.StatusOK, courseResponse)
}

func CreateCourse(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}
	role, _ := c.Get("role")
	if role.(string) != "teacher" && role.(string) != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "只有教师和管理员可以创建课程"})
		return
	}

	var params req.CourseRequest
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数无效"})
		return
	}

	course := model.Course{
		Code:        params.Code,
		Name:        params.Name,
		Description: params.Description,
		School:      params.School,
		Department:  params.Department,
		Teacher:     params.Teacher,
		Semester:    params.Semester,
		Cover:       params.Cover,
		Status:      "active",
	}
	if params.Status != "" {
		course.Status = params.Status
	}

	if err := model.CreateCourse(database.DB, &course); err != nil {
		if err.Error() == "课程代码已存在" {
			c.JSON(http.StatusConflict, gin.H{"error": "课程代码已存在"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "创建课程失败"})
		}
		return
	}

	if role.(string) == "teacher" {
		if err := model.AddUserToCourseWithRole(database.DB, userID.(uint), course.ID, "teacher"); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "创建者加入课程失败"})
			return
		}
	}

	c.JSON(http.StatusCreated, resp.CourseResponse{
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
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}
	role, _ := c.Get("role")

	var params req.UpdateCourseRequest
	if err := c.ShouldBindJSON(&params); err != nil || params.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数无效或缺少课程ID"})
		return
	}

	course, err := model.GetCourseByID(database.DB, params.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "课程不存在"})
		return
	}

	canUpdate := false
	if role.(string) == "admin" {
		canUpdate = true
	} else if role.(string) == "teacher" {
		ok, rErr := model.HasUserCourseRole(database.DB, userID.(uint), course.ID, "teacher")
		if rErr != nil { c.JSON(http.StatusInternalServerError, gin.H{"error": "权限校验失败"}); return }
		canUpdate = ok
	}
	if !canUpdate {
		c.JSON(http.StatusForbidden, gin.H{"error": "没有权限修改此课程"})
		return
	}

	if params.Code != "" && params.Code != course.Code {
		taken, tErr := model.IsCourseCodeTaken(database.DB, params.Code, course.ID)
		if tErr != nil { c.JSON(http.StatusInternalServerError, gin.H{"error": "代码校验失败"}); return }
		if taken { c.JSON(http.StatusConflict, gin.H{"error": "课程代码已存在"}); return }
	}

	updates := map[string]interface{}{
		"code":        params.Code,
		"name":        params.Name,
		"description": params.Description,
		"school":      params.School,
		"department":  params.Department,
		"teacher":     params.Teacher,
		"semester":    params.Semester,
		"cover":       params.Cover,
		"status":      params.Status,
	}
	if err := model.UpdateCourse(database.DB, course.ID, updates); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新课程信息失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "课程更新成功"})
}

func DeleteCourse(c *gin.Context) {
	// 获取用户ID和角色
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}
	role, _ := c.Get("role")

	var params req.GetCourseRequest
	if err := c.ShouldBindJSON(&params); err != nil || params.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "缺少有效的课程ID"})
		return
	}

	course, err := model.GetCourseByID(database.DB, params.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "课程不存在"})
		return
	}

	canDelete := false
	if role.(string) == "admin" {
		canDelete = true
	} else if role.(string) == "teacher" {
		ok, rErr := model.HasUserCourseRole(database.DB, userID.(uint), course.ID, "teacher")
		if rErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "权限校验失败"})
			return
		}
		canDelete = ok
	}
	if !canDelete {
		c.JSON(http.StatusForbidden, gin.H{"error": "没有权限删除此课程"})
		return
	}

	if err := model.SoftDeleteCourse(database.DB, course.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除课程失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "课程删除成功"})
}

func JoinCourse(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}
	var params req.GetCourseRequest
	if err := c.ShouldBindJSON(&params); err != nil || params.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "缺少有效的课程ID"})
		return
	}

	if err := model.JoinCourse(database.DB, userID.(uint), params.ID); err != nil {
		switch {
		case err.Error() == "课程不可用":
			c.JSON(http.StatusBadRequest, gin.H{"error": "课程不可用"})
		case err.Error() == "已加入该课程":
			c.JSON(http.StatusConflict, gin.H{"error": "已加入该课程"})
		case err == gorm.ErrRecordNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "课程不存在"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "加入课程失败"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "加入课程成功"})
}

func SearchCourses(c *gin.Context) {
	var params req.SearchCoursesRequest
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数无效"})
		return
	}
	keyword := strings.TrimSpace(params.Keyword)
	page := params.Page
	pageSize := params.PageSize
	if page < 1 { page = 1 }
	if pageSize < 1 || pageSize > 100 { pageSize = 10 }

	courses, total, err := model.SearchCourses(database.DB, keyword, params.Status, params.SortBy, params.Order, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "搜索课程失败"})
		return
	}

	var courseResponses []resp.CourseResponse
	for _, course := range courses {
		stats, sErr := model.GetCourseStats(database.DB, course.ID)
		if sErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "查询课程统计失败"})
			return
		}
		courseResponses = append(courseResponses, resp.CourseResponse{
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
			StudentCount: stats.StudentCount,
			NoteCount:    stats.NoteCount,
		})
	}

	totalPages := (int(total) + pageSize - 1) / pageSize
	c.JSON(http.StatusOK, gin.H{
		"courses":     courseResponses,
		"total":       total,
		"page":        page,
		"page_size":   pageSize,
		"total_pages": totalPages,
	})
}

func GetLatestCourses(c *gin.Context) {
	courses, err := model.GetLatestActiveCourses(database.DB, 3)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取最新课程失败"})
		return
	}

	var courseResponses []resp.CourseResponse
	for _, course := range courses {
		stats, sErr := model.GetCourseStats(database.DB, course.ID)
		if sErr != nil { c.JSON(http.StatusInternalServerError, gin.H{"error": "查询课程统计失败"}); return }
		cr := resp.CourseResponse{
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
			StudentCount: stats.StudentCount,
			NoteCount:    stats.NoteCount,
		}
		if uid, exists := c.Get("userID"); exists {
			joined, jErr := model.IsUserJoinedCourse(database.DB, uid.(uint), course.ID)
			if jErr != nil { c.JSON(http.StatusInternalServerError, gin.H{"error": "查询加入状态失败"}); return }
			cr.IsJoined = joined
		}
		courseResponses = append(courseResponses, cr)
	}

	c.JSON(http.StatusOK, courseResponses)
}

func GetMyCourses(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	var params req.ListMyCoursesRequest
	params.Page = 1
	params.PageSize = 10
	_ = c.ShouldBindJSON(&params)
	page := params.Page
	pageSize := params.PageSize
	if page < 1 { page = 1 }
	if pageSize < 1 || pageSize > 100 { pageSize = 10 }

	userCourses, total, err := model.ListUserCourses(database.DB, userID.(uint), page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询用户课程失败"})
		return
	}

	var myCourseResponses []resp.MyCourseResponse
	for _, uc := range userCourses {
		course, cErr := model.GetCourseByID(database.DB, uc.CourseID)
		if cErr != nil {
			continue
		}
		stats, sErr := model.GetCourseStats(database.DB, course.ID)
		if sErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "查询课程统计失败"})
			return
		}
		myCourseResponses = append(myCourseResponses, resp.MyCourseResponse{
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
			Role:         uc.Role,
			JoinedAt:     uc.JoinedAt,
			CreatedAt:    course.CreatedAt,
			UpdatedAt:    course.UpdatedAt,
			StudentCount: stats.StudentCount,
			NoteCount:    stats.NoteCount,
		})
	}

	totalPages := (int(total) + pageSize - 1) / pageSize
	c.JSON(http.StatusOK, gin.H{
		"data":        myCourseResponses,
		"total":       total,
		"page":        page,
		"page_size":   pageSize,
		"total_pages": totalPages,
	})
}
