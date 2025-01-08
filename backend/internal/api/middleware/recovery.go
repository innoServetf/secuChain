package middleware

import (
	"context"
	"fmt"
	"runtime/debug"

	//"github.com/InnoServe/blockSBOM/pkg/logger"
	"github.com/InnoServe/blockSBOM/pkg/response"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"gorm.io/gorm/logger"
)

// Recovery 恢复中间件
func Recovery() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		defer func() {
			if err := recover(); err != nil {
				stack := string(debug.Stack())
				logger.Error("panic recovered",
					"error", err,
					"stack", stack,
				)

				response.Error(c, consts.StatusInternalServerError,
					fmt.Sprintf("Internal Server Error: %v", err),
					nil,
				)
				c.Abort()
			}
		}()
		c.Next(ctx)
	}
}
