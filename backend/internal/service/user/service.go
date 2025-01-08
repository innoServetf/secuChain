package user

import (
	"context"
	"errors"
	"time"

	"github.com/InnoServe/blockSBOM/internal/dal"
	"github.com/InnoServe/blockSBOM/internal/dal/model"
	"github.com/InnoServe/blockSBOM/pkg/utils"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	ErrUserExists        = errors.New("user already exists")
	ErrInvalidCredential = errors.New("invalid username or password")
	ErrInvalidCode       = errors.New("invalid verification code")
	ErrCodeExpired       = errors.New("verification code expired")
	ErrUserNotFound      = errors.New("user not found")
)

type Service struct {
	db *gorm.DB
}

func NewService() *Service {
	return &Service{
		db: dal.DB,
	}
}

// Register 用户注册
func (s *Service) Register(ctx context.Context, req *RegisterRequest) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		// 检查用户是否已存在
		var exists bool
		err := tx.Model(&model.User{}).
			Where("username = ? OR email = ?", req.Username, req.Email).
			Select("1").
			Scan(&exists).Error
		if err != nil {
			return err
		}
		if exists {
			return ErrUserExists
		}

		// 密码加密
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		// 创建用户
		user := &model.User{
			Username:     req.Username,
			Email:        req.Email,
			Password:     string(hashedPassword),
			Organization: req.Organization,
			Role:         "user",
			Status:       "pending",
		}

		if err := tx.Create(user).Error; err != nil {
			return err
		}

		// 生成验证码
		code := utils.GenerateVerificationCode()
		verification := &model.UserVerification{
			UserID:    user.ID,
			Code:      code,
			Type:      "register",
			ExpiredAt: time.Now().Add(15 * time.Minute),
		}

		if err := tx.Create(verification).Error; err != nil {
			return err
		}

		// 发送验证邮件
		return s.sendVerificationEmail(user.Email, code)
	})
}

// Login 用户登录
func (s *Service) Login(ctx context.Context, req *LoginRequest) (*TokenResponse, error) {
	var user model.User
	err := s.db.Where("username = ? AND status = ?", req.Username, "active").First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidCredential
		}
		return nil, err
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, ErrInvalidCredential
	}

	// 生成Token
	accessToken, err := utils.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		return nil, err
	}

	refreshToken, err := utils.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, err
	}

	// 更新最后登录时间
	now := time.Now()
	if err := s.db.Model(&user).Update("last_login_at", &now).Error; err != nil {
		return nil, err
	}

	return &TokenResponse{
		AccessToken:  accessToken,
		TokenType:    "Bearer",
		ExpiresIn:    utils.AccessTokenExpiry,
		RefreshToken: refreshToken,
	}, nil
}

// VerifyEmail 验证邮箱
func (s *Service) VerifyEmail(ctx context.Context, req *EmailVerificationRequest) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		var verification model.UserVerification
		err := tx.Where("code = ? AND type = ? AND used_at IS NULL",
			req.Code, "register").First(&verification).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrInvalidCode
			}
			return err
		}

		if time.Now().After(verification.ExpiredAt) {
			return ErrCodeExpired
		}

		// 更新验证状态
		now := time.Now()
		if err := tx.Model(&verification).Update("used_at", &now).Error; err != nil {
			return err
		}

		// 激活用户
		return tx.Model(&model.User{}).
			Where("id = ?", verification.UserID).
			Updates(map[string]interface{}{
				"status":         "active",
				"email_verified": true,
			}).Error
	})
}

// ResetPassword 重置密码
func (s *Service) ResetPassword(ctx context.Context, req *ResetPasswordRequest) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		var verification model.UserVerification
		err := tx.Where("code = ? AND type = ? AND used_at IS NULL",
			req.Code, "reset_password").First(&verification).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrInvalidCode
			}
			return err
		}

		if time.Now().After(verification.ExpiredAt) {
			return ErrCodeExpired
		}

		// 更新密码
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		if err := tx.Model(&model.User{}).
			Where("id = ?", verification.UserID).
			Update("password", string(hashedPassword)).Error; err != nil {
			return err
		}

		// 标记验证码已使用
		now := time.Now()
		return tx.Model(&verification).Update("used_at", &now).Error
	})
}

// GetUserInfo 获取用户信息
func (s *Service) GetUserInfo(ctx context.Context, userID uint) (*UserResponse, error) {
	var user model.User
	if err := s.db.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return &UserResponse{
		ID:            user.ID,
		Username:      user.Username,
		Email:         user.Email,
		Organization:  user.Organization,
		Role:          user.Role,
		Status:        user.Status,
		EmailVerified: user.EmailVerified,
	}, nil
}

// sendVerificationEmail 发送验证邮件
func (s *Service) sendVerificationEmail(email, code string) error {
	// TODO: 实现邮件发送逻辑
	return nil
}
