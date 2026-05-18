package middleware

import (
	"github.com/gofiber/fiber/v2"
	utils "github.com/iiimomoniii/inventory_backend/utils"
)

// InitI18n — โหลด locale files ทั้งหมดเข้า memory
func InitI18n() {
	// ─── Global ────────────────────────────────────────────
	if err := utils.LoadLocaleFile("locales/global_errors.json"); err != nil {
		panic(err)
	}
	// ─── Product ───────────────────────────────────────────
	if err := utils.LoadLocaleFile("locales/product_errors.json"); err != nil {
		panic(err)
	}
	// ─── เพิ่ม module ใหม่ตรงนี้ ───────────────────────────
	// utils.LoadLocaleFile("locales/user.json")
	// utils.LoadLocaleFile("locales/order.json")
}

// I18nMiddleware — detect ภาษาจาก Accept-Language header
func I18nMiddleware(c *fiber.Ctx) error {
	lang := c.Get("Accept-Language", "en")
	if len(lang) >= 2 && lang[:2] == "th" {
		lang = "th"
	} else {
		lang = "en"
	}
	c.Locals("lang", lang)
	return c.Next()
}

// GetLang — ดึงภาษาจาก context
func GetLang(c *fiber.Ctx) string {
	if lang, ok := c.Locals("lang").(string); ok {
		return lang
	}
	return "en"
}
