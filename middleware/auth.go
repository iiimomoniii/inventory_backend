package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

func ExtractToken(authHeader string) string {
	parts := strings.Split(authHeader, " ")
	if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
		return parts[1]
	}
	return ""
}

// AuthMiddleware — validate JWT token แล้วเซ็ต Locals
// ไม่ verify signature เพราะ verify ที่ API Gateway แล้ว
func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			setEmptyLocals(c)
			return c.Next()
		}

		tokenString := ExtractToken(authHeader)
		if tokenString == "" {
			setEmptyLocals(c)
			return c.Next()
		}

		claims := jwt.MapClaims{}
		token, _, err := jwt.NewParser().ParseUnverified(tokenString, claims)
		if err != nil {
			setEmptyLocals(c)
			return c.Next()
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			setEmptyLocals(c)
			return c.Next()
		}

		// name
		if name, ok := claims["name"].(string); ok {
			c.Locals("name", name)
		} else {
			c.Locals("name", "")
		}

		// userID
		if userID, ok := claims["preferred_username"].(string); ok {
			c.Locals("userID", userID)
		} else {
			c.Locals("userID", "")
		}

		return c.Next()
	}
}

func setEmptyLocals(c *fiber.Ctx) {
	c.Locals("name", "")
	c.Locals("userID", "")
}
