package handlers

import (
	"time"

	"github.com/InnoServe/blockSBOM/internal/dal"
	"github.com/InnoServe/blockSBOM/internal/dal/model"
	"github.com/InnoServe/blockSBOM/pkg/utils"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

type LoginRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6,max=32"`
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6,max=32"`
	Email    string `json:"email" binding:"required,email"`
}

type AuthResponse struct {
	Token string      `json:"token"`
	User  *model.User `json:"user"`
}

// Register 处理用户注册
func Register(c *app.RequestContext) {
	var req RegisterRequest
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(consts.StatusBadRequest, map[string]interface{}{
			"error": "无效的请求参数: " + err.Error(),
		})
		return
	}

	// 检查用户名是否已存在
	var existingUser model.User
	if err := dal.DB.Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
		c.JSON(consts.StatusConflict, map[string]interface{}{
			"error": "用户名已存在",
		})
		return
	}

	// 检查邮箱是否已存在
	if err := dal.DB.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		c.JSON(consts.StatusConflict, map[string]interface{}{
			"error": "邮箱已被注册",
		})
		return
	}

	// 加密密码
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, map[string]interface{}{
			"error": "密码加密失败",
		})
		return
	}

	// 创建新用户
	user := &model.User{
		Username: req.Username,
		Password: hashedPassword,
		Email:    req.Email,
		Status:   "active",
	}

	if err := dal.DB.Create(user).Error; err != nil {
		c.JSON(consts.StatusInternalServerError, map[string]interface{}{
			"error": "创建用户失败: " + err.Error(),
		})
		return
	}

	c.JSON(consts.StatusCreated, map[string]interface{}{
		"message": "注册成功",
	})
}

// Login 处理用户登录
func Login(c *app.RequestContext) {
	var req LoginRequest
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(consts.StatusBadRequest, map[string]interface{}{
			"error": "无效的请求参数: " + err.Error(),
		})
		return
	}

	// 查找用户
	var user model.User
	if err := dal.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		c.JSON(consts.StatusUnauthorized, map[string]interface{}{
			"error": "用户名或密码错误",
		})
		return
	}

	// 检查用户状态
	if user.Status != "active" {
		c.JSON(consts.StatusForbidden, map[string]interface{}{
			"error": "账号已被禁用",
		})
		return
	}

	// 验证密码
	if !utils.CheckPassword(req.Password, user.Password) {
		c.JSON(consts.StatusUnauthorized, map[string]interface{}{
			"error": "用户名或密码错误",
		})
		return
	}

	// 生成 JWT token
	token, err := utils.GenerateToken(user.ID, user.Username)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, map[string]interface{}{
			"error": "生成令牌失败",
		})
		return
	}

	// 更新最后登录时间
	now := time.Now()
	user.LastLoginAt = &now
	if err := dal.DB.Save(&user).Error; err != nil {
		// 仅记录日志，不影响登录
		println("更新最后登录时间失败:", err.Error())
	}

	// 清除密码后返回用户信息
	user.Password = ""

	c.JSON(consts.StatusOK, AuthResponse{
		Token: token,
		User:  &user,
	})
}

// GetUserInfo 获取当前用户信息
func GetUserInfo(c *app.RequestContext) {
	// 从上下文中获取用户ID（由认证中间件设置）
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(consts.StatusUnauthorized, map[string]interface{}{
			"error": "未登录",
		})
		return
	}

	var user model.User
	if err := dal.DB.First(&user, userID).Error; err != nil {
		c.JSON(consts.StatusNotFound, map[string]interface{}{
			"error": "用户不存在",
		})
		return
	}

	// 清除密码后返回
	user.Password = ""

	c.JSON(consts.StatusOK, map[string]interface{}{
		"user": user,
	})
}
