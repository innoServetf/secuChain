package handlers

import (
	"context"

	"github.com/InnoServe/blockSBOM/internal/service/user"
	"github.com/InnoServe/blockSBOM/pkg/response"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

type UserHandler struct {
	userService *user.Service
}

func NewUserHandler() *UserHandler {
	return &UserHandler{
		userService: user.NewService(),
	}
}

// Register 用户注册
func (h *UserHandler) Register(ctx context.Context, c *app.RequestContext) {
	var req user.RegisterRequest
	if err := c.BindAndValidate(&req); err != nil {
		response.Error(c, consts.StatusBadRequest, "invalid request", err)
		return
	}

	if err := h.userService.Register(ctx, &req); err != nil {
		switch err {
		case user.ErrUserExists:
			response.Error(c, consts.StatusConflict, "user already exists", err)
		default:
			response.Error(c, consts.StatusInternalServerError, "registration failed", err)
		}
		return
	}

	response.Success(c, "registration successful, please check your email for verification")
}

// Login 用户登录
func (h *UserHandler) Login(ctx context.Context, c *app.RequestContext) {
	var req user.LoginRequest
	if err := c.BindAndValidate(&req); err != nil {
		response.Error(c, consts.StatusBadRequest, "invalid request", err)
		return
	}

	token, err := h.userService.Login(ctx, &req)
	if err != nil {
		switch err {
		case user.ErrInvalidCredential:
			response.Error(c, consts.StatusUnauthorized, "invalid credentials", err)
		default:
			response.Error(c, consts.StatusInternalServerError, "login failed", err)
		}
		return
	}

	response.Success(c, token)
}

// VerifyEmail 验证邮箱
func (h *UserHandler) VerifyEmail(ctx context.Context, c *app.RequestContext) {
	var req user.EmailVerificationRequest
	if err := c.BindAndValidate(&req); err != nil {
		response.Error(c, consts.StatusBadRequest, "invalid request", err)
		return
	}

	if err := h.userService.VerifyEmail(ctx, &req); err != nil {
		switch err {
		case user.ErrInvalidCode:
			response.Error(c, consts.StatusBadRequest, "invalid verification code", err)
		case user.ErrCodeExpired:
			response.Error(c, consts.StatusBadRequest, "verification code expired", err)
		default:
			response.Error(c, consts.StatusInternalServerError, "verification failed", err)
		}
		return
	}

	response.Success(c, "email verified successfully")
}

// ResetPassword 重置密码
func (h *UserHandler) ResetPassword(ctx context.Context, c *app.RequestContext) {
	var req user.ResetPasswordRequest
	if err := c.BindAndValidate(&req); err != nil {
		response.Error(c, consts.StatusBadRequest, "invalid request", err)
		return
	}

	if err := h.userService.ResetPassword(ctx, &req); err != nil {
		switch err {
		case user.ErrInvalidCode:
			response.Error(c, consts.StatusBadRequest, "invalid reset code", err)
		case user.ErrCodeExpired:
			response.Error(c, consts.StatusBadRequest, "reset code expired", err)
		default:
			response.Error(c, consts.StatusInternalServerError, "password reset failed", err)
		}
		return
	}

	response.Success(c, "password reset successfully")
}

// GetUserInfo 获取用户信息
func (h *UserHandler) GetUserInfo(ctx context.Context, c *app.RequestContext) {
	userID := c.GetUint("user_id") // 从上下文中获取用户ID（由认证中间件设置）

	userInfo, err := h.userService.GetUserInfo(ctx, userID)
	if err != nil {
		switch err {
		case user.ErrUserNotFound:
			response.Error(c, consts.StatusNotFound, "user not found", err)
		default:
			response.Error(c, consts.StatusInternalServerError, "failed to get user info", err)
		}
		return
	}

	response.Success(c, userInfo)
}
