package app

import (
	"fmt"
	"log"
	"student_shared/app/middleware"
	"student_shared/app/router"
	"student_shared/app/utils/config"
	"student_shared/app/utils/database"

	"github.com/gin-gonic/gin"
)

// Start 启动应用程序
func Start() {
	// 初始化数据库连接
	err := database.InitDB()
	if err != nil {
		log.Fatalf("数据库初始化失败: %v", err)
	}

	// 创建Gin引擎
	r := gin.Default()

	// 注册中间件
	r.Use(middleware.Cors())
	r.Use(middleware.ErrorHandler())

	// 注册路由
	router.RegisterRoutes(r)

	// 读取统一配置中的服务端口
	c := config.Load()
	port := c.Server.Port
	if port == "" {
		port = "8080"
	}
	fmt.Printf("服务器启动在 http://localhost:%s\n", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}
