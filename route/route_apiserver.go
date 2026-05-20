package routes

import (
	"database/sql"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/iiimomoniii/inventory_backend/config"
	"github.com/iiimomoniii/inventory_backend/handler"
	"github.com/iiimomoniii/inventory_backend/middleware"
	"github.com/iiimomoniii/inventory_backend/repository"
	"github.com/iiimomoniii/inventory_backend/service"
	"github.com/iiimomoniii/inventory_backend/utils"
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
	productRepo := repository.NewProductRepository(sqlDB)
	categoryRepo := repository.NewCategoryRepository(sqlDB)
	userRepo := repository.NewUserRepository(sqlDB)
	refreshTokenRepo := repository.NewRefreshTokenRepository(sqlDB)

	// ─── Services ──────────────────────────────────────────
	productSvc := service.NewProductService(productRepo)
	authSvc := service.NewAuthService(userRepo, refreshTokenRepo, cfg.JWT.Secret, cfg.JWT.ExpiresIn)
	userSvc := service.NewUserService(userRepo)

	// ─── Handlers ──────────────────────────────────────────
	productHandler := handler.NewProductHandler(productSvc)
	categoryHandler := handler.NewCategoryHandler(categoryRepo)
	authHandler := handler.NewAuthHandler(authSvc)
	userHandler := handler.NewUserHandler(userSvc)

	registerRoutes(app, productHandler, categoryHandler, authHandler, userHandler)

	return &APIServer{server: app, cfg: cfg}
}

func registerRoutes(
	app *fiber.App,
	productHandler *handler.ProductHandler,
	categoryHandler *handler.CategoryHandler,
	authHandler *handler.AuthHandler,
	userHandler *handler.UserHandler,
) {
	// ─── Public ────────────────────────────────────────────
	app.Get("/live", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	app.Post("/auth/token",
		limiter.New(limiter.Config{
			Max:        5,
			Expiration: 1 * time.Minute,
			LimitReached: func(c *fiber.Ctx) error {
				return utils.TooManyRequests(c)
			},
		}),
		authHandler.GenerateToken,
	)

	app.Post("/auth/refresh", authHandler.RefreshToken)
	app.Post("/auth/logout", authHandler.Logout)

	// Public create user
	app.Post("/v1/users/create", userHandler.Create)

	// ─── Protected Routes ──────────────────────────────────
	protected := app.Group("/v1", middleware.AuthMiddleware(), middleware.I18nMiddleware)

	// Categories
	protected.Get("/categories", categoryHandler.GetAll)
	protected.Get("/categories/:id", categoryHandler.GetByID)

	// Products
	protected.Post("/products/search", productHandler.Search)
	protected.Post("/products/create/items", productHandler.CreateItems)
	protected.Post("/products/create", productHandler.Create)
	protected.Get("/products/:id", productHandler.GetByID)
	protected.Put("/products/update/:id", productHandler.Update)
	protected.Delete("/products/delete/:id", productHandler.Delete)
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
