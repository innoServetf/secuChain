package api

import (
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/yourusername/blockSBOM/internal/api/handlers"
	"github.com/yourusername/blockSBOM/internal/api/middleware"
)

func RegisterRoutes(h *server.Hertz) {
	// 全局中间件
	h.Use(middleware.Cors())
	h.Use(middleware.RequestID())
	h.Use(middleware.Logger())

	// API 路由组
	v1 := h.Group("/api/v1")
	{
		// SBOM 相关路由
		sbom := v1.Group("/sbom")
		{
			sbom.POST("", handlers.CreateSBOM)
			sbom.GET("/:id", handlers.GetSBOM)
			sbom.PUT("/:id", handlers.UpdateSBOM)
			sbom.GET("", handlers.ListSBOMs)
		}

		// DID 相关路由
		did := v1.Group("/did")
		{
			did.POST("", handlers.CreateDID)
			did.GET("/:id", handlers.GetDID)
			did.PUT("/:id", handlers.UpdateDID)
		}

		// 漏洞扫描相关路由
		vuln := v1.Group("/vulnerability")
		{
			vuln.POST("/scan", handlers.ScanVulnerability)
			vuln.GET("/report/:id", handlers.GetVulnerabilityReport)
		}
	}
}
