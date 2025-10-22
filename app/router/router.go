package router

import (
	"student_shared/app/middleware"
	"student_shared/app/router/api"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册所有路由
func RegisterRoutes(r *gin.Engine) {
	// API版本前缀
	v1 := r.Group("/api/v1")

	// 用户相关路由
	userRoutes := v1.Group("/users")
	{
		// 公开接口 - 不需要认证
		userRoutes.POST("/register", api.RegisterUser)
		userRoutes.POST("/login", api.LoginUser)

		// 需要认证的接口
		userRoutes.POST("/profile", middleware.AuthMiddleware(), api.GetUserProfile)
		userRoutes.POST("/profile/update", middleware.AuthMiddleware(), api.UpdateUserProfile)
	}

	// 课程路由
	courseRoutes := v1.Group("/courses")
	{
		courseRoutes.POST("", api.ListCourses)          // 获取课程列表
		// 允许携带令牌以判定 is_joined（可选认证）
		courseRoutes.POST("/detail", middleware.OptionalAuthMiddleware(), api.GetCourse)     // 获取课程详情（JSON传id）
		courseRoutes.POST("/search", api.SearchCourses) // 搜索课程
	}

	// 需要认证的课程路由
	authCourseRoutes := v1.Group("/courses")
	authCourseRoutes.Use(middleware.AuthMiddleware())
	{
		authCourseRoutes.POST("/my", api.GetMyCourses)     // 获取我的课程
		authCourseRoutes.POST("/create", api.CreateCourse) // 创建课程
		authCourseRoutes.POST("/update", api.UpdateCourse) // 更新课程
		authCourseRoutes.POST("/delete", api.DeleteCourse) // 删除课程
		authCourseRoutes.POST("/join", api.JoinCourse)     // 加入课程
	}

	// 笔记相关路由
	noteRoutes := v1.Group("/notes")
	{
		// 公开接口 - 浏览公开笔记不需要认证，但支持可选认证以显示私有笔记
		noteRoutes.POST("", middleware.OptionalAuthMiddleware(), api.ListNotes)
		noteRoutes.POST("/detail", middleware.OptionalAuthMiddleware(), api.GetNote)

		// 需要认证的接口
		noteRoutes.POST("/create", middleware.AuthMiddleware(), api.CreateNote)
		noteRoutes.POST("/update", middleware.AuthMiddleware(), api.UpdateNote)
		noteRoutes.POST("/delete", middleware.AuthMiddleware(), api.DeleteNote)
		noteRoutes.POST("/like", middleware.AuthMiddleware(), api.LikeNote)
		noteRoutes.POST("/unlike", middleware.AuthMiddleware(), api.UnlikeNote)
		noteRoutes.POST("/favorite", middleware.AuthMiddleware(), api.FavoriteNote)
		noteRoutes.POST("/unfavorite", middleware.AuthMiddleware(), api.UnfavoriteNote)
		noteRoutes.POST("/favorites", middleware.AuthMiddleware(), api.GetMyFavorites)
		noteRoutes.POST("/likes", middleware.AuthMiddleware(), api.GetMyLikes)
		noteRoutes.POST("/download", middleware.AuthMiddleware(), api.DownloadNote)
	}

	// 评论相关路由
	commentGroup := v1.Group("/comment")
	{
		commentGroup.POST("/list", middleware.OptionalAuthMiddleware(), api.ListNoteComments) // 获取笔记评论列表（JSON传 noteId）
		commentGroup.POST("/create", middleware.AuthMiddleware(), api.CreateComment)          // 创建评论
		commentGroup.POST("/update", middleware.AuthMiddleware(), api.UpdateComment)          // 更新评论
		commentGroup.POST("/delete", middleware.AuthMiddleware(), api.DeleteComment)          // 删除评论
		commentGroup.POST("/like", middleware.AuthMiddleware(), api.LikeComment)              // 点赞评论
		commentGroup.POST("/unlike", middleware.AuthMiddleware(), api.UnlikeComment)          // 取消点赞评论
	}

	// 搜索路由 - 公开接口
	searchRoutes := v1.Group("/search")
	{
		searchRoutes.POST("/courses", middleware.OptionalAuthMiddleware(), api.SearchCourses)
		searchRoutes.POST("/notes", middleware.OptionalAuthMiddleware(), api.SearchNotes)
	}

	// 首页相关路由 - 无需认证
	homeRoutes := v1.Group("/home")
	{
		homeRoutes.POST("/popular-notes", api.GetPopularNotes)   // 获取热门笔记（前3个点赞最多的）
		homeRoutes.POST("/latest-courses", api.GetLatestCourses) // 获取最新课程（前3个最新创建的）
		homeRoutes.POST("/stats", api.GetPublicStats)            // 获取公开统计数据
	}

	// 文件上传路由 - 需要认证
	uploadRoutes := v1.Group("/upload")
	uploadRoutes.Use(middleware.AuthMiddleware())
	{
		uploadRoutes.POST("/avatar", api.UploadAvatar)
		uploadRoutes.POST("/avatar/delete", api.DeleteAvatar)
	}

	// AI 路由
	haiRoutes := v1.Group("/ai")
	{
		// 生成摘要：需要认证
		haiRoutes.POST("/summarize", middleware.AuthMiddleware(), api.SummarizeText)
		// 获取某笔记的AI元数据：公开笔记允许匿名，私有笔记需要鉴权
		haiRoutes.POST("/notes/meta", middleware.OptionalAuthMiddleware(), api.GetNoteAIMeta)
	}

	// 静态文件服务
	r.Static("/uploads", "./uploads")

	// 管理员路由 - 需要管理员权限
	adminRoutes := v1.Group("/admin")
	adminRoutes.Use(middleware.AuthMiddleware())
	adminRoutes.Use(api.RequireAdminRole())
	{
		// 统计数据
		adminRoutes.POST("/stats", api.GetAdminStats)
		adminRoutes.POST("/users/list", api.GetAllUsers)
		adminRoutes.POST("/users/update", api.AdminUpdateUser)
		adminRoutes.POST("/users/update-role", api.AdminUpdateUserRole)
		adminRoutes.POST("/users/delete", api.AdminDeleteUser)
	}
}
