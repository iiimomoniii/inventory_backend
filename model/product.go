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
	CategoryName string    `json:"categoryName"` // ← join จาก categories
	Price        float64   `json:"price"`
	Stock        int       `json:"stock"`
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

// ─── Custom Errors ─────────────────────────────────────────

type NotFoundError struct {
	ID int64
}

func (e *NotFoundError) Error() string { return "not found" }

type ValidationError struct {
	Code string
}

func (e *ValidationError) Error() string { return e.Code }
