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
	if req.Page < 0 {
		return nil, &model.ValidationError{Code: "INV001", Message: fmt.Sprintf("invalid page: %d", req.Page)}
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}
	if req.MinPrice != nil && req.MaxPrice != nil && *req.MinPrice > *req.MaxPrice {
		return nil, &model.ValidationError{Code: "INV002", Message: "minPrice cannot be greater than maxPrice"}
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
	if id <= 0 {
		return nil, &model.ValidationError{Code: "INV003", Message: "id must be greater than 0"}
	}
	return s.Repo.FindByID(ctx, id)
}

func (s *ProductServiceImpl) Create(ctx context.Context, req model.ProductCreateRequest, createdBy string) (*model.ProductResponse, error) {
	req.Name = strings.TrimSpace(req.Name)
	req.Category = strings.TrimSpace(req.Category)

	if req.Name == "" {
		return nil, &model.ValidationError{Code: "INV004", Message: "name is required"}
	}
	if req.Category == "" {
		return nil, &model.ValidationError{Code: "INV005", Message: "category is required"}
	}
	if req.Price <= 0 {
		return nil, &model.ValidationError{Code: "INV006", Message: "price must be greater than 0"}
	}
	if req.Stock < 0 {
		return nil, &model.ValidationError{Code: "INV007", Message: "stock must be >= 0"}
	}

	product, err := s.Repo.Create(ctx, req, createdBy)
	if err != nil {
		return nil, fmt.Errorf("create product failed: %w", err)
	}
	return product, nil
}

func (s *ProductServiceImpl) Update(ctx context.Context, id int64, req model.ProductUpdateRequest, updatedBy string) (*model.ProductResponse, error) {
	if id <= 0 {
		return nil, &model.ValidationError{Code: "INV003", Message: "id must be greater than 0"}
	}
	if req.Price == nil && req.Stock == nil {
		return nil, &model.ValidationError{Code: "INV008", Message: "at least one field (price or stock) must be provided"}
	}
	if req.Price != nil && *req.Price <= 0 {
		return nil, &model.ValidationError{Code: "INV006", Message: "price must be greater than 0"}
	}
	if req.Stock != nil && *req.Stock < 0 {
		return nil, &model.ValidationError{Code: "INV007", Message: "stock must be >= 0"}
	}

	return s.Repo.Update(ctx, id, req, updatedBy)
}

func (s *ProductServiceImpl) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return &model.ValidationError{Code: "INV003", Message: "id must be greater than 0"}
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
