package service

import (
	"context"
	"errors"
	"testing"

	"github.com/iiimomoniii/inventory_backend/model"
	"github.com/iiimomoniii/inventory_backend/repository"
)

// ─── Mock Repository ───────────────────────────────────────

type mockProductRepository struct {
	searchResult []model.ProductResponse
	searchTotal  int
	searchErr    error

	findResult *model.ProductResponse
	findErr    error

	createResult *model.ProductResponse
	createErr    error

	updateResult *model.ProductResponse
	updateErr    error

	deleteErr error

	capturedCreateReq model.ProductCreateRequest
	capturedUpdatedBy string
	capturedDeleteID  int64
}

func (m *mockProductRepository) Search(_ context.Context, _ model.ProductSearchRequest) ([]model.ProductResponse, int, error) {
	return m.searchResult, m.searchTotal, m.searchErr
}
func (m *mockProductRepository) FindByID(_ context.Context, _ int64) (*model.ProductResponse, error) {
	return m.findResult, m.findErr
}
func (m *mockProductRepository) Create(_ context.Context, req model.ProductCreateRequest, _ string) (*model.ProductResponse, error) {
	m.capturedCreateReq = req
	return m.createResult, m.createErr
}
func (m *mockProductRepository) Update(_ context.Context, _ int64, _ model.ProductUpdateRequest, updatedBy string) (*model.ProductResponse, error) {
	m.capturedUpdatedBy = updatedBy
	return m.updateResult, m.updateErr
}
func (m *mockProductRepository) Delete(_ context.Context, id int64) error {
	m.capturedDeleteID = id
	return m.deleteErr
}

// ─── Helpers ───────────────────────────────────────────────

// ใช้ repository.ProductRepository interface แทน *mockProductRepository
func newSvc(mock repository.ProductRepository) *ProductServiceImpl {
	return &ProductServiceImpl{Repo: mock}
}

func ptr[T any](v T) *T { return &v }

// ─── Search ────────────────────────────────────────────────

func TestSearch_Success(t *testing.T) {
	mock := &mockProductRepository{
		searchResult: []model.ProductResponse{
			{ID: 1, Name: "Apple"},
			{ID: 2, Name: "Banana"},
		},
		searchTotal: 2,
	}
	result, err := newSvc(mock).Search(context.Background(), model.ProductSearchRequest{Page: 0, PageSize: 10})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.TotalElements != 2 {
		t.Errorf("TotalElements = %d, want 2", result.TotalElements)
	}
	if len(result.Data) != 2 {
		t.Errorf("len(Data) = %d, want 2", len(result.Data))
	}
}

func TestSearch_InvalidPage(t *testing.T) {
	_, err := newSvc(&mockProductRepository{}).Search(context.Background(), model.ProductSearchRequest{Page: -1})
	assertValidationError(t, err, "PRD006")
}

func TestSearch_InvalidPriceRange(t *testing.T) {
	_, err := newSvc(&mockProductRepository{}).Search(context.Background(), model.ProductSearchRequest{
		Page: 0, PageSize: 10,
		MinPrice: ptr(100.0),
		MaxPrice: ptr(50.0),
	})
	assertValidationError(t, err, "PRD007")
}

func TestSearch_DefaultPageSize(t *testing.T) {
	result, err := newSvc(&mockProductRepository{}).Search(context.Background(), model.ProductSearchRequest{Page: 0, PageSize: 0})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.PageSize != 20 {
		t.Errorf("PageSize = %d, want 20", result.PageSize)
	}
}

func TestSearch_RepoError_Propagates(t *testing.T) {
	sentinel := errors.New("db down")
	_, err := newSvc(&mockProductRepository{searchErr: sentinel}).Search(
		context.Background(), model.ProductSearchRequest{Page: 0, PageSize: 10},
	)
	if !errors.Is(err, sentinel) {
		t.Errorf("expected sentinel error, got: %v", err)
	}
}

// ─── GetByID ───────────────────────────────────────────────

func TestGetByID_Success(t *testing.T) {
	expected := &model.ProductResponse{ID: 1, Name: "Apple"}
	result, err := newSvc(&mockProductRepository{findResult: expected}).GetByID(context.Background(), 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.ID != 1 {
		t.Errorf("ID = %d, want 1", result.ID)
	}
}

func TestGetByID_InvalidID(t *testing.T) {
	tests := []struct {
		name string
		id   int64
	}{
		{"zero", 0},
		{"negative", -5},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := newSvc(&mockProductRepository{}).GetByID(context.Background(), tc.id)
			assertValidationError(t, err, "PRD005")
		})
	}
}

func TestGetByID_NotFound(t *testing.T) {
	_, err := newSvc(&mockProductRepository{
		findErr: &model.NotFoundError{ID: 99},
	}).GetByID(context.Background(), 99)

	var notFound *model.NotFoundError
	if !errors.As(err, &notFound) {
		t.Fatalf("expected *NotFoundError, got %T", err)
	}
	if notFound.ID != 99 {
		t.Errorf("ID = %d, want 99", notFound.ID)
	}
}

// ─── Create ────────────────────────────────────────────────

func TestCreate_Success(t *testing.T) {
	expected := &model.ProductResponse{ID: 1, Name: "Apple", Price: 10.0}
	mock := &mockProductRepository{createResult: expected}

	result, err := newSvc(mock).Create(context.Background(), model.ProductCreateRequest{
		Name: "Apple", CategoryID: 1, Price: 10.0, Stock: 100,
	}, "admin")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.ID != expected.ID {
		t.Errorf("ID = %d, want %d", result.ID, expected.ID)
	}
	if mock.capturedCreateReq.Name != "Apple" {
		t.Errorf("capturedName = %q, want 'Apple'", mock.capturedCreateReq.Name)
	}
}

func TestCreate_ValidationErrors(t *testing.T) {
	tests := []struct {
		name     string
		req      model.ProductCreateRequest
		wantCode string
	}{
		{"empty name", model.ProductCreateRequest{Name: "", CategoryID: 1, Price: 10, Stock: 0}, "PRD001"},
		{"zero categoryId", model.ProductCreateRequest{Name: "Apple", CategoryID: 0, Price: 10, Stock: 0}, "PRD002"},
		{"zero price", model.ProductCreateRequest{Name: "Apple", CategoryID: 1, Price: 0, Stock: 0}, "PRD003"},
		{"negative price", model.ProductCreateRequest{Name: "Apple", CategoryID: 1, Price: -1, Stock: 0}, "PRD003"},
		{"negative stock", model.ProductCreateRequest{Name: "Apple", CategoryID: 1, Price: 10, Stock: -1}, "PRD004"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := newSvc(&mockProductRepository{}).Create(context.Background(), tc.req, "admin")
			assertValidationError(t, err, tc.wantCode)
		})
	}
}

// ─── Update ────────────────────────────────────────────────

func TestUpdate_Success(t *testing.T) {
	expected := &model.ProductResponse{ID: 1, Price: 99.0, Stock: 50}
	mock := &mockProductRepository{updateResult: expected}

	result, err := newSvc(mock).Update(context.Background(), 1, model.ProductUpdateRequest{
		Price: ptr(99.0), Stock: ptr(50),
	}, "admin")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Price != 99.0 {
		t.Errorf("Price = %f, want 99.0", result.Price)
	}
	if mock.capturedUpdatedBy != "admin" {
		t.Errorf("updatedBy = %q, want 'admin'", mock.capturedUpdatedBy)
	}
}

func TestUpdate_NoFields(t *testing.T) {
	_, err := newSvc(&mockProductRepository{}).Update(context.Background(), 1, model.ProductUpdateRequest{}, "admin")
	assertValidationError(t, err, "PRD008")
}

func TestUpdate_NotFound(t *testing.T) {
	_, err := newSvc(&mockProductRepository{updateErr: &model.NotFoundError{ID: 99}}).Update(
		context.Background(), 99, model.ProductUpdateRequest{Price: ptr(10.0)}, "admin",
	)
	var notFound *model.NotFoundError
	if !errors.As(err, &notFound) {
		t.Fatalf("expected *NotFoundError, got %T", err)
	}
}

// ─── Delete ────────────────────────────────────────────────

func TestDelete_Success(t *testing.T) {
	mock := &mockProductRepository{}
	if err := newSvc(mock).Delete(context.Background(), 1); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if mock.capturedDeleteID != 1 {
		t.Errorf("capturedDeleteID = %d, want 1", mock.capturedDeleteID)
	}
}

func TestDelete_InvalidID(t *testing.T) {
	err := newSvc(&mockProductRepository{}).Delete(context.Background(), 0)
	assertValidationError(t, err, "PRD005")
}

func TestDelete_NotFound(t *testing.T) {
	err := newSvc(&mockProductRepository{deleteErr: &model.NotFoundError{ID: 99}}).Delete(context.Background(), 99)
	var notFound *model.NotFoundError
	if !errors.As(err, &notFound) {
		t.Fatalf("expected *NotFoundError, got %T", err)
	}
}

func TestDelete_RepoError_Propagates(t *testing.T) {
	sentinel := errors.New("db error")
	err := newSvc(&mockProductRepository{deleteErr: sentinel}).Delete(context.Background(), 1)
	if !errors.Is(err, sentinel) {
		t.Errorf("expected sentinel error, got: %v", err)
	}
}

// ─── Helper ────────────────────────────────────────────────

func assertValidationError(t *testing.T, err error, wantCode string) {
	t.Helper()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var valErr *model.ValidationError
	if !errors.As(err, &valErr) {
		t.Fatalf("expected *ValidationError, got %T: %v", err, err)
	}
	if valErr.Code != wantCode {
		t.Errorf("code = %q, want %q", valErr.Code, wantCode)
	}
}
