package repository

import (
	"context"

	"github.com/iiimomoniii/inventory_backend/model"
)

// Interface เท่านั้น — implementation จริงต่อ DB แยกไฟล์
// ทำให้ service และ test ไม่พึ่ง DB ตรงๆ

type ProductRepository interface {
	Search(ctx context.Context, req model.ProductSearchRequest) ([]model.ProductResponse, int, error)
	FindByID(ctx context.Context, id int64) (*model.ProductResponse, error)
	Create(ctx context.Context, req model.ProductCreateRequest, createdBy string) (*model.ProductResponse, error)
	Update(ctx context.Context, id int64, req model.ProductUpdateRequest, updatedBy string) (*model.ProductResponse, error)
	Delete(ctx context.Context, id int64) error
}
