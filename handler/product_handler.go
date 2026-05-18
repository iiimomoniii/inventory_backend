package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/iiimomoniii/inventory_backend/model"
	"github.com/iiimomoniii/inventory_backend/service"
)

type ProductHandler struct {
	Service service.ProductService
}

func NewProductHandler(svc service.ProductService) *ProductHandler {
	return &ProductHandler{Service: svc}
}

// ServeHTTP — router อย่างง่าย
// POST /products/search
// GET  /products/{id}
// POST /products
// PUT  /products/{id}
// DELETE /products/{id}
func (h *ProductHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	path := strings.TrimPrefix(r.URL.Path, "/products")
	path = strings.TrimSuffix(path, "/")

	switch {
	case r.Method == http.MethodPost && path == "/search":
		h.search(w, r)
	case r.Method == http.MethodPost && path == "":
		h.create(w, r)
	case r.Method == http.MethodGet && path != "":
		h.getByID(w, r, strings.TrimPrefix(path, "/"))
	case r.Method == http.MethodPut && path != "":
		h.update(w, r, strings.TrimPrefix(path, "/"))
	case r.Method == http.MethodDelete && path != "":
		h.delete(w, r, strings.TrimPrefix(path, "/"))
	default:
		writeJSON(w, http.StatusNotFound, model.Response{Message: "route not found"})
	}
}

func (h *ProductHandler) search(w http.ResponseWriter, r *http.Request) {
	var req model.ProductSearchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, model.Response{Message: "invalid request body"})
		return
	}

	result, err := h.Service.Search(context.Background(), req)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, model.Response{Message: "success", Data: result})
}

func (h *ProductHandler) getByID(w http.ResponseWriter, r *http.Request, idStr string) {
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, model.Response{Message: "invalid id"})
		return
	}

	product, err := h.Service.GetByID(context.Background(), id)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, model.Response{Message: "success", Data: product})
}

func (h *ProductHandler) create(w http.ResponseWriter, r *http.Request) {
	var req model.ProductCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, model.Response{Message: "invalid request body"})
		return
	}

	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		userID = "anonymous"
	}

	product, err := h.Service.Create(context.Background(), req, userID)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, model.Response{Message: "created", Data: product})
}

func (h *ProductHandler) update(w http.ResponseWriter, r *http.Request, idStr string) {
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, model.Response{Message: "invalid id"})
		return
	}

	var req model.ProductUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, model.Response{Message: "invalid request body"})
		return
	}

	userID := r.Header.Get("X-User-ID")
	product, err := h.Service.Update(context.Background(), id, req, userID)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, model.Response{Message: "updated", Data: product})
}

func (h *ProductHandler) delete(w http.ResponseWriter, r *http.Request, idStr string) {
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, model.Response{Message: "invalid id"})
		return
	}

	if err := h.Service.Delete(context.Background(), id); err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, model.Response{Message: "deleted"})
}

// ─── Helpers ───────────────────────────────────────────────

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, err error) {
	var notFound *model.NotFoundError
	if errors.As(err, &notFound) {
		writeJSON(w, http.StatusNotFound, model.Response{Message: "not found"})
		return
	}
	var valErr *model.ValidationError
	if errors.As(err, &valErr) {
		writeJSON(w, http.StatusBadRequest, model.Response{Message: valErr.Message})
		return
	}
	writeJSON(w, http.StatusInternalServerError, model.Response{Message: "internal server error"})
}
