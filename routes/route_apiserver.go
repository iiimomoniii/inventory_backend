package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/iiimomoniii/inventory_backend/config"
	"github.com/iiimomoniii/inventory_backend/handler"
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
		AppName: cfg.AppName,
	})

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

// registerRoutes — ลงทะเบียน route ทั้งหมด
// แยกออกมาเพื่อให้อ่านง่ายและเพิ่ม route ใหม่ได้สะดวก
func registerRoutes(app *fiber.App, productHandler *handler.ProductHandler) {

	// ─── Public (no auth required) ─────────────────────────
	app.Get("/live", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	// ─── Products ──────────────────────────────────────────
	app.Post("/products/search", productHandler.Search)
	app.Get("/products/:id", productHandler.GetByID)
	app.Post("/products", productHandler.Create)
	app.Put("/products/:id", productHandler.Update)
	app.Delete("/products/:id", productHandler.Delete)
}

// Start — implement App interface
func (s *APIServer) Start() {
	if err := s.server.Listen(s.cfg.AppPort); err != nil {
		panic(err)
	}
}

// Stop — implement App interface
func (s *APIServer) Stop() {
	if err := s.server.Shutdown(); err != nil {
		panic(err)
	}
}
