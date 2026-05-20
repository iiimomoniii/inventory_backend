package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
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
	Refresh(ctx context.Context, refreshToken string) (*model.TokenResponse, error)
	Logout(ctx context.Context, refreshToken string) error
}

// ─── Implementation ────────────────────────────────────────

type AuthServiceImpl struct {
	UserRepo         repository.UserRepository
	RefreshTokenRepo repository.RefreshTokenRepository
	SecretKey        string
	ExpiresIn        int // access token seconds
	RefreshExpiresIn int // refresh token days
}

func NewAuthService(
	userRepo repository.UserRepository,
	refreshTokenRepo repository.RefreshTokenRepository,
	secretKey string,
	expiresIn int,
) *AuthServiceImpl {
	return &AuthServiceImpl{
		UserRepo:         userRepo,
		RefreshTokenRepo: refreshTokenRepo,
		SecretKey:        secretKey,
		ExpiresIn:        expiresIn,
		RefreshExpiresIn: 7, // 7 วัน
	}
}

func (s *AuthServiceImpl) Login(ctx context.Context, username, password string) (*model.TokenResponse, error) {
	// 1. หา user จาก DB
	user, err := s.UserRepo.FindByUsername(ctx, username)
	if err != nil {
		var notFound *model.NotFoundError
		if errors.As(err, &notFound) {
			return nil, &ServiceError{Code: "GLB004"}
		}
		return nil, fmt.Errorf("login failed: %w", err)
	}

	// 2. เช็ค password (TODO: bcrypt)
	if user.Password != password {
		return nil, &ServiceError{Code: "GLB004"}
	}

	// 3. สร้าง access token
	accessToken, err := s.generateJWT(user)
	if err != nil {
		return nil, fmt.Errorf("generate access token failed: %w", err)
	}

	// 4. สร้าง refresh token
	refreshToken, expiresAt, err := s.generateRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("generate refresh token failed: %w", err)
	}

	// 5. เก็บ refresh token ใน DB
	if err := s.RefreshTokenRepo.Save(ctx, user.ID, refreshToken, expiresAt); err != nil {
		return nil, fmt.Errorf("save refresh token failed: %w", err)
	}

	return &model.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    s.ExpiresIn,
	}, nil
}

func (s *AuthServiceImpl) Refresh(ctx context.Context, refreshToken string) (*model.TokenResponse, error) {
	// 1. หา refresh token จาก DB
	rt, err := s.RefreshTokenRepo.FindByToken(ctx, refreshToken)
	if err != nil {
		return nil, &ServiceError{Code: "GLB004"}
	}

	// 2. เช็คว่า revoke แล้วหรือยัง
	if rt.RevokedAt != nil {
		return nil, &ServiceError{Code: "GLB004"}
	}

	// 3. เช็คหมดอายุ
	if time.Now().After(rt.ExpiresAt) {
		return nil, &ServiceError{Code: "GLB004"}
	}

	// 4. หา user จาก DB
	user, err := s.UserRepo.FindByID(ctx, rt.UserID)
	if err != nil {
		return nil, fmt.Errorf("find user failed: %w", err)
	}

	// 5. สร้าง access token ใหม่
	accessToken, err := s.generateJWT(user)
	if err != nil {
		return nil, fmt.Errorf("generate access token failed: %w", err)
	}

	return &model.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken, // คืน refresh token เดิม
		TokenType:    "Bearer",
		ExpiresIn:    s.ExpiresIn,
	}, nil
}

func (s *AuthServiceImpl) Logout(ctx context.Context, refreshToken string) error {
	if err := s.RefreshTokenRepo.Revoke(ctx, refreshToken); err != nil {
		var notFound *model.NotFoundError
		if errors.As(err, &notFound) {
			return &ServiceError{Code: "GLB004"}
		}
		return fmt.Errorf("logout failed: %w", err)
	}
	return nil
}

// ─── Helpers ───────────────────────────────────────────────

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

func (s *AuthServiceImpl) generateRefreshToken() (string, time.Time, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", time.Time{}, fmt.Errorf("generate random bytes failed: %w", err)
	}
	token := hex.EncodeToString(b)
	expiresAt := time.Now().AddDate(0, 0, s.RefreshExpiresIn)
	return token, expiresAt, nil
}
