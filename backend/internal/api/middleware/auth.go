package middleware

import (
	"context"
	"strings"

	"github.com/InnoServe/blockSBOM/pkg/response"
	"github.com/InnoServe/blockSBOM/pkg/utils"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

// Auth JWT认证中间件
func Auth() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		authHeader := string(c.Request.Header.Get("Authorization"))
		if authHeader == "" {
			c.JSON(consts.StatusUnauthorized, map[string]interface{}{
				"error": "未授权访问",
			})
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(consts.StatusUnauthorized, map[string]interface{}{
				"error": "无效的授权头",
			})
			c.Abort()
			return
		}

		claims, err := utils.ParseToken(parts[1])
		if err != nil {
			c.JSON(consts.StatusUnauthorized, map[string]interface{}{
				"error": "无效的token",
			})
			c.Abort()
			return
		}

		// 将用户信息存储在上下文中
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)

		c.Next(ctx)
	}
}

// RoleAuth 角色验证中间件
func RoleAuth(roles ...string) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		userRole := c.GetString("role")

		for _, role := range roles {
			if role == userRole {
				c.Next(ctx)
				return
			}
		}

		response.Error(c, consts.StatusForbidden, "permission denied", nil)
		c.Abort()
	}
}
