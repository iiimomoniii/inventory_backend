package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/iiimomoniii/inventory_backend/model"
)

// ─── Interface ─────────────────────────────────────────────

type CategoryRepository interface {
	FindAll(ctx context.Context) ([]model.Category, error)
	FindByID(ctx context.Context, id int64) (*model.Category, error)
}

// ─── Implementation ────────────────────────────────────────

type CategoryRepositoryImpl struct {
	DB *sql.DB
}

func NewCategoryRepository(db *sql.DB) *CategoryRepositoryImpl {
	return &CategoryRepositoryImpl{DB: db}
}

func (r *CategoryRepositoryImpl) FindAll(ctx context.Context) ([]model.Category, error) {
	query := `SELECT id, name, created_at, updated_at
	          FROM categories
	          WHERE deleted_at IS NULL
	          ORDER BY id ASC`

	rows, err := r.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("find all categories failed: %w", err)
	}
	defer rows.Close()

	var categories []model.Category
	for rows.Next() {
		var c model.Category
		if err := rows.Scan(&c.ID, &c.Name, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan failed: %w", err)
		}
		categories = append(categories, c)
	}
	return categories, nil
}

func (r *CategoryRepositoryImpl) FindByID(ctx context.Context, id int64) (*model.Category, error) {
	query := `SELECT id, name, created_at, updated_at
	          FROM categories
	          WHERE id = $1 AND deleted_at IS NULL`

	var c model.Category
	err := r.DB.QueryRowContext(ctx, query, id).Scan(
		&c.ID, &c.Name, &c.CreatedAt, &c.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, &model.NotFoundError{ID: id}
		}
		return nil, fmt.Errorf("find category by id failed: %w", err)
	}
	return &c, nil
}
