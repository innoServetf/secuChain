package middleware

import (
	"context"
	"time"

	//"github.com/InnoServe/blockSBOM/pkg/logger"
	"github.com/cloudwego/hertz/pkg/app"
	"gorm.io/gorm/logger"
)

// Logger 日志中间件
func Logger() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		start := time.Now()
		path := string(c.Request.URI().Path())
		method := string(c.Request.Method())

		c.Next(ctx)

		latency := time.Since(start)
		statusCode := c.Response.StatusCode()
		clientIP := c.ClientIP()

		logger.Info("request",
			"method", method,
			"path", path,
			"status_code", statusCode,
			"latency", latency,
			"client_ip", clientIP,
			"user_id", c.GetUint("user_id"),
		)
	}
}
