package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/iiimomoniii/inventory_backend/model"
)

// ─── Interface ─────────────────────────────────────────────

type RefreshTokenRepository interface {
	Save(ctx context.Context, userID int64, token string, expiresAt time.Time) error
	FindByToken(ctx context.Context, token string) (*model.RefreshToken, error)
	Revoke(ctx context.Context, token string) error
	RevokeAllByUserID(ctx context.Context, userID int64) error
}

// ─── Implementation ────────────────────────────────────────

type RefreshTokenRepositoryImpl struct {
	DB *sql.DB
}

func NewRefreshTokenRepository(db *sql.DB) *RefreshTokenRepositoryImpl {
	return &RefreshTokenRepositoryImpl{DB: db}
}

func (r *RefreshTokenRepositoryImpl) Save(ctx context.Context, userID int64, token string, expiresAt time.Time) error {
	_, err := r.DB.ExecContext(ctx, `
		INSERT INTO refresh_tokens (user_id, token, expires_at)
		VALUES ($1, $2, $3)`,
		userID, token, expiresAt,
	)
	if err != nil {
		return fmt.Errorf("save refresh token failed: %w", err)
	}
	return nil
}

func (r *RefreshTokenRepositoryImpl) FindByToken(ctx context.Context, token string) (*model.RefreshToken, error) {
	var t model.RefreshToken
	err := r.DB.QueryRowContext(ctx, `
		SELECT id, user_id, token, expires_at, created_at, revoked_at
		FROM refresh_tokens
		WHERE token = $1`, token,
	).Scan(&t.ID, &t.UserID, &t.Token, &t.ExpiresAt, &t.CreatedAt, &t.RevokedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, &model.NotFoundError{}
		}
		return nil, fmt.Errorf("find refresh token failed: %w", err)
	}
	return &t, nil
}

func (r *RefreshTokenRepositoryImpl) Revoke(ctx context.Context, token string) error {
	result, err := r.DB.ExecContext(ctx, `
		UPDATE refresh_tokens SET revoked_at = NOW()
		WHERE token = $1 AND revoked_at IS NULL`, token,
	)
	if err != nil {
		return fmt.Errorf("revoke refresh token failed: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return &model.NotFoundError{}
	}
	return nil
}

func (r *RefreshTokenRepositoryImpl) RevokeAllByUserID(ctx context.Context, userID int64) error {
	_, err := r.DB.ExecContext(ctx, `
		UPDATE refresh_tokens SET revoked_at = NOW()
		WHERE user_id = $1 AND revoked_at IS NULL`, userID,
	)
	if err != nil {
		return fmt.Errorf("revoke all refresh tokens failed: %w", err)
	}
	return nil
}
