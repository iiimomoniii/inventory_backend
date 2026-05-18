package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/iiimomoniii/inventory_backend/model"
)

func ExtractToken(authHeader string) string {
	parts := strings.Split(authHeader, " ")
	if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
		return parts[1]
	}
	return ""
}

// AuthMiddleware — validate JWT token
// Priority 1: ไม่มี token      → 401 ทันที
// Priority 2: token format ผิด → 401
// Priority 3: token parse ผิด  → 401
// Priority 4: claims ผิด       → 401
// Priority 5: ผ่านทั้งหมด      → เซ็ต Locals แล้วไปต่อ
func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {

		// Priority 1 — ไม่มี Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return unauthorized(c, "Authorization header is missing")
		}

		// Priority 2 — format ต้องเป็น "Bearer <token>"
		tokenString := ExtractToken(authHeader)
		if tokenString == "" {
			return unauthorized(c, "Invalid token format, expected: Bearer <token>")
		}

		// Priority 3 — parse token
		claims := jwt.MapClaims{}
		token, _, err := jwt.NewParser().ParseUnverified(tokenString, claims)
		if err != nil {
			return unauthorized(c, "Invalid token")
		}

		// Priority 4 — ดึง claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return unauthorized(c, "Invalid token claims")
		}

		// Priority 5 — เซ็ต Locals แล้วไปต่อ
		c.Locals("name", getStringClaim(claims, "name"))
		c.Locals("userID", getStringClaim(claims, "preferred_username"))

		return c.Next()
	}
}

// ─── Helpers ───────────────────────────────────────────────

func unauthorized(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusUnauthorized).JSON(model.Response{
		Status:  fiber.StatusUnauthorized,
		Message: message,
	})
}

func getStringClaim(claims jwt.MapClaims, key string) string {
	if val, ok := claims[key].(string); ok {
		return val
	}
	return ""
}
