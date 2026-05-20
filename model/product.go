package model

import "time"

// ─── Request ───────────────────────────────────────────────

type ProductSearchRequest struct {
	Name       *string  `json:"name,omitempty"`
	CategoryID *int64   `json:"categoryId,omitempty"`
	MinPrice   *float64 `json:"minPrice,omitempty"`
	MaxPrice   *float64 `json:"maxPrice,omitempty"`
	Page       int      `json:"page"`
	PageSize   int      `json:"pageSize"`
}

type ProductCreateRequest struct {
	Name       string  `json:"name"`
	CategoryID int64   `json:"categoryId"`
	Price      float64 `json:"price"`
	Stock      int     `json:"stock"`
}

type ProductUpdateRequest struct {
	Price *float64 `json:"price,omitempty"`
	Stock *int     `json:"stock,omitempty"`
}

// ─── Response ──────────────────────────────────────────────

type ProductResponse struct {
	ID           int64     `json:"id"`
	Name         string    `json:"name"`
	CategoryID   int64     `json:"categoryId"`
	CategoryName string    `json:"categoryName"`
	Price        float64   `json:"price"`
	Stock        int       `json:"stock"`
	CreatedBy    string    `json:"createdBy"`
	UpdatedBy    string    `json:"updatedBy"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

type ProductSearchResponse struct {
	Data          []ProductResponse `json:"data"`
	TotalElements int               `json:"totalElements"`
	TotalPages    int               `json:"totalPages"`
	Page          int               `json:"page"`
	PageSize      int               `json:"pageSize"`
}
