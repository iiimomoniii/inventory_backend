package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/iiimomoniii/inventory_backend/model"
	"github.com/iiimomoniii/inventory_backend/repository"
)

// ─── Interface ─────────────────────────────────────────────

type ProductService interface {
	Search(ctx context.Context, req model.ProductSearchRequest) (*model.ProductSearchResponse, error)
	GetByID(ctx context.Context, id int64) (*model.ProductResponse, error)
	Create(ctx context.Context, req model.ProductCreateRequest, createdBy string) (*model.ProductResponse, error)
	Update(ctx context.Context, id int64, req model.ProductUpdateRequest, updatedBy string) (*model.ProductResponse, error)
	Delete(ctx context.Context, id int64) error
}

// ─── Implementation ────────────────────────────────────────

type ProductServiceImpl struct {
	Repo repository.ProductRepository
}

func NewProductService(repo repository.ProductRepository) *ProductServiceImpl {
	return &ProductServiceImpl{Repo: repo}
}

func (s *ProductServiceImpl) Search(ctx context.Context, req model.ProductSearchRequest) (*model.ProductSearchResponse, error) {
	// Priority 1 — validate page
	if req.Page < 0 {
		return nil, &model.ValidationError{Code: "PRD006"}
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}
	// Priority 2 — validate price range
	if req.MinPrice != nil && req.MaxPrice != nil && *req.MinPrice > *req.MaxPrice {
		return nil, &model.ValidationError{Code: "PRD007"}
	}

	products, total, err := s.Repo.Search(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("search failed: %w", err)
	}

	totalPages := 0
	if total > 0 {
		totalPages = (total + req.PageSize - 1) / req.PageSize
	}

	return &model.ProductSearchResponse{
		Data:          products,
		TotalElements: total,
		TotalPages:    totalPages,
		Page:          req.Page,
		PageSize:      req.PageSize,
	}, nil
}

func (s *ProductServiceImpl) GetByID(ctx context.Context, id int64) (*model.ProductResponse, error) {
	// Priority 1 — validate id
	if id <= 0 {
		return nil, &model.ValidationError{Code: "PRD005"}
	}
	return s.Repo.FindByID(ctx, id)
}

func (s *ProductServiceImpl) Create(ctx context.Context, req model.ProductCreateRequest, createdBy string) (*model.ProductResponse, error) {
	req.Name = strings.TrimSpace(req.Name)
	req.Category = strings.TrimSpace(req.Category)

	// Priority 1 — name
	if req.Name == "" {
		return nil, &model.ValidationError{Code: "PRD001"}
	}
	// Priority 2 — category
	if req.Category == "" {
		return nil, &model.ValidationError{Code: "PRD002"}
	}
	// Priority 3 — price
	if req.Price <= 0 {
		return nil, &model.ValidationError{Code: "PRD003"}
	}
	// Priority 4 — stock
	if req.Stock < 0 {
		return nil, &model.ValidationError{Code: "PRD004"}
	}

	product, err := s.Repo.Create(ctx, req, createdBy)
	if err != nil {
		return nil, fmt.Errorf("create product failed: %w", err)
	}
	return product, nil
}

func (s *ProductServiceImpl) Update(ctx context.Context, id int64, req model.ProductUpdateRequest, updatedBy string) (*model.ProductResponse, error) {
	// Priority 1 — id
	if id <= 0 {
		return nil, &model.ValidationError{Code: "PRD005"}
	}
	// Priority 2 — ต้องมีอย่างน้อย 1 field
	if req.Price == nil && req.Stock == nil {
		return nil, &model.ValidationError{Code: "PRD008"}
	}
	// Priority 3 — price
	if req.Price != nil && *req.Price <= 0 {
		return nil, &model.ValidationError{Code: "PRD003"}
	}
	// Priority 4 — stock
	if req.Stock != nil && *req.Stock < 0 {
		return nil, &model.ValidationError{Code: "PRD004"}
	}

	return s.Repo.Update(ctx, id, req, updatedBy)
}

func (s *ProductServiceImpl) Delete(ctx context.Context, id int64) error {
	// Priority 1 — id
	if id <= 0 {
		return &model.ValidationError{Code: "PRD005"}
	}

	if err := s.Repo.Delete(ctx, id); err != nil {
		var notFound *model.NotFoundError
		if errors.As(err, &notFound) {
			return notFound
		}
		return fmt.Errorf("delete product failed: %w", err)
	}
	return nil
}
