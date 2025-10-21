package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"runtime/debug"
)

// ErrorResponse 定义统一的错误响应格式
type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// ErrorHandler 全局错误处理中间件
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 使用defer捕获可能的panic
		defer func() {
			if err := recover(); err != nil {
				// 记录堆栈信息
				debug.PrintStack()
				
				// 返回500错误
				c.JSON(http.StatusInternalServerError, ErrorResponse{
					Code:    http.StatusInternalServerError,
					Message: "服务器内部错误",
				})
				
				// 终止后续中间件
				c.Abort()
			}
		}()
		
		// 继续执行后续中间件
		c.Next()
		
		// 处理请求过程中设置的错误
		if len(c.Errors) > 0 {
			// 获取最后一个错误
			err := c.Errors.Last()
			
			// 如果响应已经发送，则不再处理
			if c.Writer.Written() {
				return
			}
			
			// 返回400错误
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
			})
		}
	}
}