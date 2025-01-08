package api

import (
	"github.com/InnoServe/blockSBOM/internal/api/handlers"
	"github.com/InnoServe/blockSBOM/internal/api/middleware"
	"github.com/cloudwego/hertz/pkg/app/server"
)

func RegisterRoutes(h *server.Hertz) {
	// 全局中间件
	h.Use(middleware.Recovery())
	h.Use(middleware.Logger())
	h.Use(middleware.Cors())

	// 用户认证相关路由
	userHandler := handlers.NewUserHandler()
	auth := h.Group("/api/v1/auth")
	{
		auth.POST("/register", userHandler.Register)
		auth.POST("/login", userHandler.Login)
		auth.POST("/verify-email", userHandler.VerifyEmail)
		auth.POST("/reset-password", userHandler.ResetPassword)
	}

	// 需要认证的路由
	api := h.Group("/api/v1")
	api.Use(middleware.AuthMiddleware())
	{
		// 用户相关
		api.GET("/user/info", userHandler.GetUserInfo)

		// 管理员路由
		admin := api.Group("/admin")
		admin.Use(middleware.RoleAuth("admin"))
		{
			// TODO: 添加管理员路由
		}
	}
}
