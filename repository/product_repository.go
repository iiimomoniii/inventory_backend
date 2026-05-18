package repository

import (
	"context"

	"github.com/iiimomoniii/inventory_backend/model"
)

// ─── Interface ─────────────────────────────────────────────
// ใครก็ implement ได้ถ้ามี method ครบ
// ตอนนี้ใช้ InMemoryProductRepository
// อนาคตเปลี่ยนเป็น PostgresProductRepository โดยไม่แตะ service

type ProductRepository interface {
	Search(ctx context.Context, req model.ProductSearchRequest) ([]model.ProductResponse, int, error)
	FindByID(ctx context.Context, id int64) (*model.ProductResponse, error)
	Create(ctx context.Context, req model.ProductCreateRequest, createdBy string) (*model.ProductResponse, error)
	Update(ctx context.Context, id int64, req model.ProductUpdateRequest, updatedBy string) (*model.ProductResponse, error)
	Delete(ctx context.Context, id int64) error
}
