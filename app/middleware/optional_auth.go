package middleware

import (
	"student_shared/app/utils/jwt"

	"github.com/gin-gonic/gin"
)

// OptionalAuthMiddleware 可选认证中间件
// 如果提供了有效的token，则设置用户信息到上下文中
// 如果没有token或token无效，则继续处理请求但不设置用户信息
func OptionalAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头获取token
		token := c.GetHeader("token")
		if token == "" {
			// 没有token，继续处理请求
			c.Next()
			return
		}

		// 尝试解析令牌
		claims, err := jwt.ParseToken(token)
		if err != nil {
			// token无效，继续处理请求但不设置用户信息
			c.Next()
			return
		}

		// 将用户信息存储到上下文中
		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)

		c.Next()
	}
}