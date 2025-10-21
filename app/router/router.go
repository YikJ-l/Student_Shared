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
		userRoutes.GET("/profile", middleware.AuthMiddleware(), api.GetUserProfile)
		userRoutes.PUT("/profile", middleware.AuthMiddleware(), api.UpdateUserProfile)
	}

	// 课程路由
	courseRoutes := v1.Group("/courses")
	{
		courseRoutes.GET("", api.ListCourses)       // 获取课程列表
		courseRoutes.GET("/:id", api.GetCourse)     // 获取课程详情
		courseRoutes.GET("/search", api.SearchCourses) // 搜索课程
	}

	// 需要认证的课程路由
	authCourseRoutes := v1.Group("/courses")
	authCourseRoutes.Use(middleware.AuthMiddleware())
	{
		authCourseRoutes.GET("/my", api.GetMyCourses)      // 获取我的课程
		authCourseRoutes.POST("", api.CreateCourse)        // 创建课程
		authCourseRoutes.PUT("/:id", api.UpdateCourse)     // 更新课程
		authCourseRoutes.DELETE("/:id", api.DeleteCourse)  // 删除课程
		authCourseRoutes.POST("/:id/join", api.JoinCourse) // 加入课程
	}

	// 笔记相关路由
	noteRoutes := v1.Group("/notes")
	{
		// 公开接口 - 浏览公开笔记不需要认证，但支持可选认证以显示私有笔记
		noteRoutes.GET("", middleware.OptionalAuthMiddleware(), api.ListNotes)
		noteRoutes.GET("/:id", middleware.OptionalAuthMiddleware(), api.GetNote)

		// 需要认证的接口
		noteRoutes.POST("", middleware.AuthMiddleware(), api.CreateNote)
		noteRoutes.PUT("/:id", middleware.AuthMiddleware(), api.UpdateNote)
		noteRoutes.DELETE("/:id", middleware.AuthMiddleware(), api.DeleteNote)
		noteRoutes.POST("/:id/like", middleware.AuthMiddleware(), api.LikeNote)
		noteRoutes.DELETE("/:id/like", middleware.AuthMiddleware(), api.UnlikeNote)
		noteRoutes.POST("/:id/favorite", middleware.AuthMiddleware(), api.FavoriteNote)
		noteRoutes.DELETE("/:id/favorite", middleware.AuthMiddleware(), api.UnfavoriteNote)
		noteRoutes.GET("/favorites", middleware.AuthMiddleware(), api.GetMyFavorites)
		noteRoutes.GET("/likes", middleware.AuthMiddleware(), api.GetMyLikes)
		noteRoutes.GET("/:id/download", middleware.AuthMiddleware(), api.DownloadNote)
	}

	// 评论相关路由
	commentGroup := v1.Group("/comment")
	{
		commentGroup.GET("/note/:noteId", middleware.OptionalAuthMiddleware(), api.ListNoteComments) // 获取笔记评论列表
		commentGroup.POST("/", middleware.AuthMiddleware(), api.CreateComment)                       // 创建评论
		commentGroup.PUT("/:id", middleware.AuthMiddleware(), api.UpdateComment)                     // 更新评论
		commentGroup.DELETE("/:id", middleware.AuthMiddleware(), api.DeleteComment)                  // 删除评论
		commentGroup.POST("/:id/like", middleware.AuthMiddleware(), api.LikeComment)                 // 点赞评论
		commentGroup.DELETE("/:id/like", middleware.AuthMiddleware(), api.UnlikeComment)             // 取消点赞评论
	}

	// 搜索路由 - 公开接口
	searchRoutes := v1.Group("/search")
	{
		searchRoutes.GET("/courses", middleware.OptionalAuthMiddleware(), api.SearchCourses)
		searchRoutes.GET("/notes", middleware.OptionalAuthMiddleware(), api.SearchNotes)
	}

	// 首页相关路由 - 无需认证
	homeRoutes := v1.Group("/home")
	{
		homeRoutes.GET("/popular-notes", api.GetPopularNotes)   // 获取热门笔记（前3个点赞最多的）
		homeRoutes.GET("/latest-courses", api.GetLatestCourses) // 获取最新课程（前3个最新创建的）
		homeRoutes.GET("/stats", api.GetPublicStats)           // 获取公开统计数据
	}

	// 文件上传路由 - 需要认证
	uploadRoutes := v1.Group("/upload")
	uploadRoutes.Use(middleware.AuthMiddleware())
	{
		uploadRoutes.POST("/avatar", api.UploadAvatar)
		uploadRoutes.DELETE("/avatar", api.DeleteAvatar)
	}

	// 静态文件服务
	r.Static("/uploads", "./uploads")

	// 管理员路由 - 需要管理员权限
	adminRoutes := v1.Group("/admin")
	adminRoutes.Use(middleware.AuthMiddleware())
	adminRoutes.Use(api.RequireAdminRole())
	{
		// 统计数据
		adminRoutes.GET("/stats", api.GetAdminStats)

		// 用户管理
		adminRoutes.GET("/users", api.GetAllUsers)
		adminRoutes.PUT("/users/:id", api.AdminUpdateUser)
		adminRoutes.PUT("/users/:id/role", api.AdminUpdateUserRole)
		adminRoutes.DELETE("/users/:id", api.AdminDeleteUser)
	}
}
