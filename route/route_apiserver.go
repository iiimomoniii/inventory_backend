package routes

import (
	"database/sql"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/iiimomoniii/inventory_backend/config"
	"github.com/iiimomoniii/inventory_backend/handler"
	"github.com/iiimomoniii/inventory_backend/middleware"
	"github.com/iiimomoniii/inventory_backend/repository"
	"github.com/iiimomoniii/inventory_backend/service"
)

type APIServer struct {
	server *fiber.App
	cfg    config.AppConfig
}

func NewAPIServer(cfg config.AppConfig, sqlDB *sql.DB) App {

	app := fiber.New(fiber.Config{
		AppName:        cfg.App.Name,
		ReadTimeout:    time.Millisecond * time.Duration(cfg.Fiber.ReadTimeout),
		WriteTimeout:   time.Millisecond * time.Duration(cfg.Fiber.WriteTimeout),
		IdleTimeout:    time.Millisecond * time.Duration(cfg.Fiber.IdleTimeout),
		ReadBufferSize: cfg.Fiber.ReadBufferSize,
		BodyLimit:      cfg.Fiber.BodyLimitSize,
	})

	middleware.InitI18n()
	app.Use(middleware.CorsMiddleware())

	// ─── Repositories ──────────────────────────────────────
	productRepo := repository.NewPostgresProductRepository(sqlDB)
	categoryRepo := repository.NewCategoryRepository(sqlDB)
	userRepo := repository.NewUserRepository(sqlDB)

	// ─── Services ──────────────────────────────────────────
	productSvc := service.NewProductService(productRepo)
	authSvc := service.NewAuthService(userRepo, cfg.JWT.Secret, cfg.JWT.ExpiresIn)

	// ─── Handlers ──────────────────────────────────────────
	productHandler := handler.NewProductHandler(productSvc)
	categoryHandler := handler.NewCategoryHandler(categoryRepo)
	authHandler := handler.NewAuthHandler(authSvc)

	registerRoutes(app, productHandler, categoryHandler, authHandler)

	return &APIServer{server: app, cfg: cfg}
}

func registerRoutes(app *fiber.App, productHandler *handler.ProductHandler, categoryHandler *handler.CategoryHandler, authHandler *handler.AuthHandler) {

	// ─── Public ────────────────────────────────────────────
	app.Get("/live", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})
	app.Post("/auth/token", authHandler.GenerateToken)

	// ─── Auth middleware ────────────────────────────────────
	app.Use(middleware.AuthMiddleware())
	app.Use(middleware.I18nMiddleware)

	// ─── v1 ────────────────────────────────────────────────
	v1 := app.Group("/v1")

	// Categories
	v1.Get("/categories", categoryHandler.GetAll)
	v1.Get("/categories/:id", categoryHandler.GetByID)

	// Products
	v1.Post("/products/search", productHandler.Search)
	v1.Post("/products/create/items", productHandler.CreateItems)
	v1.Get("/products/:id", productHandler.GetByID)
	v1.Post("/products/create", productHandler.Create)
	v1.Put("/products/update/:id", productHandler.Update)
	v1.Delete("/products/delete/:id", productHandler.Delete)
}

func (s *APIServer) Start() {
	if err := s.server.Listen(s.cfg.Fiber.Address); err != nil {
		panic(err)
	}
}

func (s *APIServer) Stop() {
	if err := s.server.Shutdown(); err != nil {
		panic(err)
	}
}
