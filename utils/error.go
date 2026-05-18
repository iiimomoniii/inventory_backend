package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/iiimomoniii/inventory_backend/model"
)

// ─── Locale Store ──────────────────────────────────────────

type ErrorEntry struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	MessageEN string `json:"messageEN"`
	MessageTH string `json:"messageTH"`
}

var (
	localeStore = map[string]ErrorEntry{}
	localeMu    sync.RWMutex
)

func LoadLocaleFile(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read locale file %s: %w", path, err)
	}
	var entries map[string]ErrorEntry
	if err := json.Unmarshal(data, &entries); err != nil {
		return fmt.Errorf("failed to parse locale file %s: %w", path, err)
	}
	localeMu.Lock()
	defer localeMu.Unlock()
	for k, v := range entries {
		localeStore[k] = v
	}
	return nil
}

func GetMessage(code string, lang string) string {
	localeMu.RLock()
	defer localeMu.RUnlock()
	entry, ok := localeStore[code]
	if !ok {
		return code
	}
	if lang == "th" {
		return entry.MessageTH
	}
	return entry.MessageEN
}

func GetName(code string) string {
	localeMu.RLock()
	defer localeMu.RUnlock()
	entry, ok := localeStore[code]
	if !ok {
		return code
	}
	return entry.Name
}

// ─── Single Error Response ─────────────────────────────────

func ErrorResponse(c *fiber.Ctx, status int, code string) error {
	messageEN := GetMessage(code, "en")
	messageTH := GetMessage(code, "th")
	name := GetName(code)

	return c.Status(status).JSON(model.Response{
		Status:  status,
		Message: messageEN,
		Data: model.ErrorData{
			Timestamp: time.Now().Format(time.RFC3339),
			Status:    status,
			ErrorCode: code,
			Error:     name,
			MessageEN: messageEN,
			MessageTH: messageTH,
			Path:      c.Path(),
		},
	})
}

// ─── Items Response ────────────────────────────────────────

func ItemsResponse(c *fiber.Ctx, results []model.ItemResult) error {
	totalSuccess, totalError := 0, 0
	for _, r := range results {
		if r.Status == "success" {
			totalSuccess++
		} else {
			totalError++
		}
	}

	status := fiber.StatusOK
	message := "success"
	if totalError > 0 && totalSuccess > 0 {
		status = fiber.StatusMultiStatus
		message = GetMessage("PRD009", "en")
	} else if totalError > 0 && totalSuccess == 0 {
		status = fiber.StatusBadRequest
		message = "All items failed"
	}

	return c.Status(status).JSON(model.Response{
		Status:  status,
		Message: message,
		Data: model.ItemsErrorData{
			Timestamp:    time.Now().Format(time.RFC3339),
			Status:       status,
			TotalSuccess: totalSuccess,
			TotalError:   totalError,
			Results:      results,
			Path:         c.Path(),
		},
	})
}

func BuildSuccessItem(index int, data any) model.ItemResult {
	return model.ItemResult{Index: index, Status: "success", Data: data}
}

func BuildErrorItem(index int, code string) model.ItemResult {
	return model.ItemResult{
		Index:     index,
		Status:    "error",
		ErrorCode: code,
		Error:     GetName(code),
		MessageEN: GetMessage(code, "en"),
		MessageTH: GetMessage(code, "th"),
	}
}

// ─── Global Errors (GLB) ───────────────────────────────────

func NotFound(c *fiber.Ctx) error { return ErrorResponse(c, fiber.StatusNotFound, "GLB001") }
func InternalError(c *fiber.Ctx) error {
	return ErrorResponse(c, fiber.StatusInternalServerError, "GLB002")
}
func InvalidBody(c *fiber.Ctx) error  { return ErrorResponse(c, fiber.StatusBadRequest, "GLB003") }
func Unauthorized(c *fiber.Ctx) error { return ErrorResponse(c, fiber.StatusUnauthorized, "GLB004") }

// ─── Product Errors (PRD) ──────────────────────────────────

func BadRequest(c *fiber.Ctx, code string) error {
	return ErrorResponse(c, fiber.StatusBadRequest, code)
}
