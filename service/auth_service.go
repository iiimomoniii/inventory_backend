package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/iiimomoniii/inventory_backend/model"
	"github.com/iiimomoniii/inventory_backend/repository"
)

// ─── Interface ─────────────────────────────────────────────

type AuthService interface {
	Login(ctx context.Context, username, password string) (*model.TokenResponse, error)
}

// ─── Implementation ────────────────────────────────────────

type AuthServiceImpl struct {
	UserRepo  repository.UserRepository
	SecretKey string
	ExpiresIn int
}

func NewAuthService(userRepo repository.UserRepository, secretKey string, expiresIn int) *AuthServiceImpl {
	return &AuthServiceImpl{
		UserRepo:  userRepo,
		SecretKey: secretKey,
		ExpiresIn: expiresIn,
	}
}

func (s *AuthServiceImpl) Login(ctx context.Context, username, password string) (*model.TokenResponse, error) {
	// 1. หา user จาก DB
	user, err := s.UserRepo.FindByUsername(ctx, username)
	if err != nil {
		var notFound *model.NotFoundError
		if errors.As(err, &notFound) {
			return nil, &model.ValidationError{Code: "GLB004"} // Unauthorized
		}
		return nil, fmt.Errorf("login failed: %w", err)
	}

	// 2. เช็ค password
	// ในของจริงใช้ bcrypt.CompareHashAndPassword
	if user.Password != password {
		return nil, &model.ValidationError{Code: "GLB004"} // Unauthorized
	}

	// 3. สร้าง JWT token
	token, err := s.generateJWT(user)
	if err != nil {
		return nil, fmt.Errorf("generate token failed: %w", err)
	}

	return &model.TokenResponse{
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresIn:   s.ExpiresIn,
	}, nil
}

func (s *AuthServiceImpl) generateJWT(user *model.User) (string, error) {
	claims := jwt.MapClaims{
		"preferred_username": user.Username,
		"name":               user.Name,
		"role":               user.Role,
		"exp":                time.Now().Add(time.Duration(s.ExpiresIn) * time.Second).Unix(),
		"iat":                time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.SecretKey))
}
