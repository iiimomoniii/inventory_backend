package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/iiimomoniii/inventory_backend/model"
	"github.com/iiimomoniii/inventory_backend/repository"
	"golang.org/x/crypto/bcrypt"
)

// ─── Interface ─────────────────────────────────────────────

type UserService interface {
	Create(ctx context.Context, req model.UserCreateRequest, createdBy string) (*model.UserResponse, error)
}

// ─── Implementation ────────────────────────────────────────

type UserServiceImpl struct {
	Repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserServiceImpl {
	return &UserServiceImpl{Repo: repo}
}

func (s *UserServiceImpl) Create(ctx context.Context, req model.UserCreateRequest, createdBy string) (*model.UserResponse, error) {
	req.Username = strings.TrimSpace(req.Username)
	req.Name = strings.TrimSpace(req.Name)
	req.Role = strings.TrimSpace(req.Role)

	if req.Username == "" {
		return nil, &ServiceError{Code: "USR001"}
	}
	if req.Password == "" {
		return nil, &ServiceError{Code: "USR002"}
	}
	if len(req.Password) < 6 {
		return nil, &ServiceError{Code: "USR003"}
	}
	if req.Name == "" {
		return nil, &ServiceError{Code: "USR004"}
	}
	if req.Role == "" {
		req.Role = "user"
	}

	existingUser, err := s.Repo.FindByUsername(ctx, req.Username)
	if err == nil && existingUser != nil {
		return nil, &ServiceError{Code: "USR005"}
	}

	var notFoundErr *model.NotFoundError
	if err != nil && !errors.As(err, &notFoundErr) {
		return nil, fmt.Errorf("check duplicate username failed: %w", err)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hash password failed: %w", err)
	}
	req.Password = string(hash)

	return s.Repo.Create(ctx, req, createdBy)
}
