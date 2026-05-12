package router

import (
	v1 "admin/api/v1"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// API v1 路由组
	apiV1 := r.Group("/api/v1")
	{
		// 用户相关
		users := apiV1.Group("/users")
		{
			users.GET("", v1.ListUsers)
			users.GET("/:id", v1.GetUser)
			users.POST("", v1.CreateUser)
			users.PUT("/:id", v1.UpdateUser)
			users.DELETE("/:id", v1.DeleteUser)
		}
		users.Use() // 权限设置
		// 认证相关不需要权限控制
		auth := apiV1.Group("/auth")
		{
			auth.POST("/login", v1.Login)
			auth.POST("/register", v1.Register)
			auth.POST("/logout", v1.Logout)
		}
	}

	apiV1.Use() //跨域设置
	return r
}
