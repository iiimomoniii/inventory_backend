package routes

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/iiimomoniii/inventory_backend/config"
	"github.com/iiimomoniii/inventory_backend/handler"
	"github.com/iiimomoniii/inventory_backend/middleware"
	"github.com/iiimomoniii/inventory_backend/repository"
	"github.com/iiimomoniii/inventory_backend/service"
)

// APIServer — implement App interface
type APIServer struct {
	server *fiber.App
	cfg    config.AppConfig
}

// NewAPIServer — wire dependencies ทั้งหมด
// Repository → Service → Handler (Layered Architecture)
func NewAPIServer(cfg config.AppConfig) App {

	// ─── Fiber App ─────────────────────────────────────────
	app := fiber.New(fiber.Config{
		AppName:        cfg.App.Name,                                            // ← cfg.App.Name
		ReadTimeout:    time.Millisecond * time.Duration(cfg.Fiber.ReadTimeout), // ← จาก yaml
		WriteTimeout:   time.Millisecond * time.Duration(cfg.Fiber.WriteTimeout),
		IdleTimeout:    time.Millisecond * time.Duration(cfg.Fiber.IdleTimeout),
		ReadBufferSize: cfg.Fiber.ReadBufferSize, // ← 8192
		BodyLimit:      cfg.Fiber.BodyLimitSize,
	})

	// ─── Global Middleware ─────────────────────────────────
	middleware.InitI18n()
	app.Use(middleware.CorsMiddleware())

	// ─── Repositories ──────────────────────────────────────
	productRepo := repository.NewInMemoryProductRepository()

	// ─── Services ──────────────────────────────────────────
	productSvc := service.NewProductService(productRepo)

	// ─── Handlers ──────────────────────────────────────────
	productHandler := handler.NewProductHandler(productSvc)

	// ─── Routes ────────────────────────────────────────────
	registerRoutes(app, productHandler)

	return &APIServer{server: app, cfg: cfg}
}

func registerRoutes(app *fiber.App, productHandler *handler.ProductHandler) {

	// ─── Public (no auth required) ─────────────────────────
	app.Get("/live", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	// ─── Auth middleware ────────────────────────────────────
	app.Use(middleware.AuthMiddleware())
	app.Use(middleware.I18nMiddleware)

	// ─── Products ──────────────────────────────────────────
	// ⚠️ specific routes ต้องอยู่ก่อน wildcard /:id เสมอ
	app.Post("/products/search", productHandler.Search)
	app.Post("/products/create/items", productHandler.CreateItems) // ← ย้ายขึ้นก่อน /:id
	app.Get("/products/:id", productHandler.GetByID)
	app.Post("/products", productHandler.Create)
	app.Put("/products/:id", productHandler.Update)
	app.Delete("/products/:id", productHandler.Delete)
	app.Post("/products/create/items", productHandler.CreateItems)
}

// Start — implement App interface
func (s *APIServer) Start() {
	if err := s.server.Listen(s.cfg.Fiber.Address); err != nil { // ← cfg.Fiber.Address
		panic(err)
	}
}

// Stop — implement App interface
func (s *APIServer) Stop() {
	if err := s.server.Shutdown(); err != nil {
		panic(err)
	}
}
