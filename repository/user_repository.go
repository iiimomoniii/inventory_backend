package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/iiimomoniii/inventory_backend/model"
)

// ─── Interface ─────────────────────────────────────────────

type UserRepository interface {
	FindByUsername(ctx context.Context, username string) (*model.User, error)
	FindByID(ctx context.Context, id int64) (*model.User, error)
	Create(ctx context.Context, req model.UserCreateRequest, createdBy string) (*model.UserResponse, error)
}

// ─── Implementation ────────────────────────────────────────

type UserRepositoryImpl struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepositoryImpl {
	return &UserRepositoryImpl{DB: db}
}

func (r *UserRepositoryImpl) FindByUsername(ctx context.Context, username string) (*model.User, error) {
	query := `
		SELECT id, username, password, name, role,
		       created_by, updated_by, created_at, updated_at
		FROM users
		WHERE username = $1 AND deleted_at IS NULL`

	var u model.User
	err := r.DB.QueryRowContext(ctx, query, username).Scan(
		&u.ID, &u.Username, &u.Password, &u.Name, &u.Role,
		&u.CreatedBy, &u.UpdatedBy, &u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, &model.NotFoundError{}
		}
		return nil, fmt.Errorf("find user failed: %w", err)
	}
	return &u, nil
}

func (r *UserRepositoryImpl) FindByID(ctx context.Context, id int64) (*model.User, error) {
	query := `
		SELECT id, username, password, name, role,
		       created_by, updated_by, created_at, updated_at
		FROM users
		WHERE id = $1 AND deleted_at IS NULL`

	var u model.User
	err := r.DB.QueryRowContext(ctx, query, id).Scan(
		&u.ID, &u.Username, &u.Password, &u.Name, &u.Role,
		&u.CreatedBy, &u.UpdatedBy, &u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, &model.NotFoundError{}
		}
		return nil, fmt.Errorf("find user by id failed: %w", err)
	}
	return &u, nil
}

func (r *UserRepositoryImpl) Create(ctx context.Context, req model.UserCreateRequest, createdBy string) (*model.UserResponse, error) {
	var u model.UserResponse
	err := r.DB.QueryRowContext(ctx, `
		INSERT INTO users (username, password, name, role, created_by, updated_by)
		VALUES ($1, $2, $3, $4, $5, $5)
		RETURNING id, username, name, role, created_by, updated_by`,
		req.Username, req.Password, req.Name, req.Role, createdBy,
	).Scan(&u.ID, &u.Username, &u.Name, &u.Role, &u.CreatedBy, &u.UpdatedBy)
	if err != nil {
		return nil, fmt.Errorf("create user failed: %w", err)
	}
	return &u, nil
}
