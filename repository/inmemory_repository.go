package repository

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/iiimomoniii/inventory_backend/model"
)

// InMemoryProductRepository — ใช้สำหรับ dev/demo
// ไม่ต้องต่อ DB จริง
// เปลี่ยนเป็น PostgresProductRepository ได้โดยไม่แตะ service เลย

type InMemoryProductRepository struct {
	mu       sync.RWMutex
	products map[int64]model.ProductResponse
	nextID   int64
}

func NewInMemoryProductRepository() *InMemoryProductRepository {
	r := &InMemoryProductRepository{
		products: make(map[int64]model.ProductResponse),
		nextID:   1,
	}
	r.seed()
	return r
}

// seed — ข้อมูลตัวอย่างเริ่มต้น
func (r *InMemoryProductRepository) seed() {
	items := []model.ProductCreateRequest{
		{Name: "Apple", Category: "Fruit", Price: 10.0, Stock: 100},
		{Name: "Banana", Category: "Fruit", Price: 5.0, Stock: 200},
		{Name: "Carrot", Category: "Vegetable", Price: 8.0, Stock: 5},
		{Name: "Milk", Category: "Dairy", Price: 35.0, Stock: 50},
		{Name: "Pork", Category: "Meat", Price: 120.0, Stock: 30},
	}
	for _, item := range items {
		p := model.ProductResponse{
			ID:        r.nextID,
			Name:      item.Name,
			Category:  item.Category,
			Price:     item.Price,
			Stock:     item.Stock,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		r.products[r.nextID] = p
		r.nextID++
	}
}

func (r *InMemoryProductRepository) Search(_ context.Context, req model.ProductSearchRequest) ([]model.ProductResponse, int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []model.ProductResponse
	for _, p := range r.products {
		if req.Name != nil && *req.Name != "" {
			if !strings.Contains(strings.ToLower(p.Name), strings.ToLower(*req.Name)) {
				continue
			}
		}
		if req.Category != nil && *req.Category != "" && p.Category != *req.Category {
			continue
		}
		if req.MinPrice != nil && p.Price < *req.MinPrice {
			continue
		}
		if req.MaxPrice != nil && p.Price > *req.MaxPrice {
			continue
		}
		result = append(result, p)
	}

	total := len(result)
	page := req.Page
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}
	start := page * pageSize
	if start >= total {
		return []model.ProductResponse{}, total, nil
	}
	end := start + pageSize
	if end > total {
		end = total
	}
	return result[start:end], total, nil
}

func (r *InMemoryProductRepository) FindByID(_ context.Context, id int64) (*model.ProductResponse, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	p, ok := r.products[id]
	if !ok {
		return nil, &model.NotFoundError{ID: id}
	}
	return &p, nil
}

func (r *InMemoryProductRepository) Create(_ context.Context, req model.ProductCreateRequest, _ string) (*model.ProductResponse, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	p := model.ProductResponse{
		ID:        r.nextID,
		Name:      req.Name,
		Category:  req.Category,
		Price:     req.Price,
		Stock:     req.Stock,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	r.products[r.nextID] = p
	r.nextID++
	return &p, nil
}

func (r *InMemoryProductRepository) Update(_ context.Context, id int64, req model.ProductUpdateRequest, _ string) (*model.ProductResponse, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	p, ok := r.products[id]
	if !ok {
		return nil, &model.NotFoundError{ID: id}
	}
	if req.Price != nil {
		p.Price = *req.Price
	}
	if req.Stock != nil {
		p.Stock = *req.Stock
	}
	p.UpdatedAt = time.Now()
	r.products[id] = p
	return &p, nil
}

func (r *InMemoryProductRepository) Delete(_ context.Context, id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.products[id]; !ok {
		return &model.NotFoundError{ID: id}
	}
	delete(r.products, id)
	return nil
}
