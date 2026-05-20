package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/iiimomoniii/inventory_backend/model"
)

// ─── Interface ─────────────────────────────────────────────

type ProductRepository interface {
	Search(ctx context.Context, req model.ProductSearchRequest) ([]model.ProductResponse, int, error)
	FindByID(ctx context.Context, id int64) (*model.ProductResponse, error)
	Create(ctx context.Context, req model.ProductCreateRequest, createdBy string) (*model.ProductResponse, error)
	Update(ctx context.Context, id int64, req model.ProductUpdateRequest, updatedBy string) (*model.ProductResponse, error)
	Delete(ctx context.Context, id int64) error
}

// ─── Implementation ────────────────────────────────────────

type ProductRepositoryImpl struct {
	DB *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepositoryImpl {
	return &ProductRepositoryImpl{DB: db}
}

// baseSelect — SELECT + JOIN ที่ใช้ซ้ำ รวม created_by, updated_by
const baseSelect = `
	SELECT p.id, p.name, p.category_id, c.name AS category_name,
	       p.price, p.stock,
	       COALESCE(p.created_by, ''), COALESCE(p.updated_by, ''),
	       p.created_at, p.updated_at
	FROM products p
	JOIN categories c ON c.id = p.category_id`

// scanProduct — scan row เป็น ProductResponse
func scanProduct(row interface{ Scan(...any) error }) (model.ProductResponse, error) {
	var p model.ProductResponse
	err := row.Scan(
		&p.ID, &p.Name, &p.CategoryID, &p.CategoryName,
		&p.Price, &p.Stock,
		&p.CreatedBy, &p.UpdatedBy,
		&p.CreatedAt, &p.UpdatedAt,
	)
	return p, err
}

func (r *ProductRepositoryImpl) Search(ctx context.Context, req model.ProductSearchRequest) ([]model.ProductResponse, int, error) {
	page := req.Page
	pageSize := req.PageSize
	if page < 0 {
		page = 0
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	offset := page * pageSize

	where := "WHERE p.deleted_at IS NULL"
	params := []interface{}{}
	n := 1

	if req.Name != nil && strings.TrimSpace(*req.Name) != "" {
		where += fmt.Sprintf(" AND p.name ILIKE $%d", n)
		params = append(params, "%"+*req.Name+"%")
		n++
	}
	if req.CategoryID != nil {
		where += fmt.Sprintf(" AND p.category_id = $%d", n)
		params = append(params, *req.CategoryID)
		n++
	}
	if req.MinPrice != nil {
		where += fmt.Sprintf(" AND p.price >= $%d", n)
		params = append(params, *req.MinPrice)
		n++
	}
	if req.MaxPrice != nil {
		where += fmt.Sprintf(" AND p.price <= $%d", n)
		params = append(params, *req.MaxPrice)
		n++
	}

	var total int
	countCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	if err := r.DB.QueryRowContext(countCtx,
		"SELECT COUNT(*) FROM products p "+where, params...,
	).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count query failed: %w", err)
	}

	query := fmt.Sprintf("%s %s ORDER BY p.created_at DESC LIMIT $%d OFFSET $%d",
		baseSelect, where, n, n+1,
	)
	params = append(params, pageSize, offset)

	rows, err := r.DB.QueryContext(ctx, query, params...)
	if err != nil {
		return nil, 0, fmt.Errorf("search query failed: %w", err)
	}
	defer rows.Close()

	var products []model.ProductResponse
	for rows.Next() {
		p, err := scanProduct(rows)
		if err != nil {
			return nil, 0, fmt.Errorf("scan failed: %w", err)
		}
		products = append(products, p)
	}
	return products, total, nil
}

func (r *ProductRepositoryImpl) FindByID(ctx context.Context, id int64) (*model.ProductResponse, error) {
	query := baseSelect + " WHERE p.id = $1 AND p.deleted_at IS NULL"

	p, err := scanProduct(r.DB.QueryRowContext(ctx, query, id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, &model.NotFoundError{ID: id}
		}
		return nil, fmt.Errorf("find by id failed: %w", err)
	}
	return &p, nil
}

func (r *ProductRepositoryImpl) Create(ctx context.Context, req model.ProductCreateRequest, createdBy string) (*model.ProductResponse, error) {
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("begin tx failed: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	var id int64
	err = tx.QueryRowContext(ctx, `
		INSERT INTO products (name, category_id, price, stock, created_by, updated_by)
		VALUES ($1, $2, $3, $4, $5, $5)
		RETURNING id`,
		req.Name, req.CategoryID, req.Price, req.Stock, createdBy,
	).Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("insert failed: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit failed: %w", err)
	}

	return r.FindByID(ctx, id)
}

func (r *ProductRepositoryImpl) Update(ctx context.Context, id int64, req model.ProductUpdateRequest, updatedBy string) (*model.ProductResponse, error) {
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("begin tx failed: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	setClauses := []string{"updated_by = $1", "updated_at = NOW()"}
	params := []interface{}{updatedBy}
	n := 2

	if req.Price != nil {
		setClauses = append(setClauses, fmt.Sprintf("price = $%d", n))
		params = append(params, *req.Price)
		n++
	}
	if req.Stock != nil {
		setClauses = append(setClauses, fmt.Sprintf("stock = $%d", n))
		params = append(params, *req.Stock)
		n++
	}

	params = append(params, id)
	result, err := tx.ExecContext(ctx, fmt.Sprintf(`
		UPDATE products SET %s WHERE id = $%d AND deleted_at IS NULL`,
		strings.Join(setClauses, ", "), n,
	), params...)
	if err != nil {
		return nil, fmt.Errorf("update failed: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		tx.Rollback()
		return nil, &model.NotFoundError{ID: id}
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit failed: %w", err)
	}

	return r.FindByID(ctx, id)
}

func (r *ProductRepositoryImpl) Delete(ctx context.Context, id int64) error {
	result, err := r.DB.ExecContext(ctx,
		"UPDATE products SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL", id,
	)
	if err != nil {
		return fmt.Errorf("delete failed: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return &model.NotFoundError{ID: id}
	}
	return nil
}
