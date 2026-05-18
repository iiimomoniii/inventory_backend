package routes

import (
	"fmt"
	"net/http"
	"os"

	"github.com/yourname/inventory-api/handler"
	"github.com/yourname/inventory-api/repository"
	"github.com/yourname/inventory-api/service"
)

// AppConfig — config ที่โหลดจาก env
type AppConfig struct {
	AppPort string
	AppName string
}

// APIServer — implement App interface
type APIServer struct {
	server *http.Server
	cfg    AppConfig
}

// NewAPIServer — wire dependencies ทั้งหมด
// Repository → Service → Handler (Layered Architecture)
func NewAPIServer(cfg AppConfig) App {

	// ─── Repositories ──────────────────────────────────────
	// ในของจริงส่ง DB เข้ามา ตอนนี้ใช้ in-memory
	productRepo := repository.NewInMemoryProductRepository()

	// ─── Services ──────────────────────────────────────────
	productSvc := service.NewProductService(productRepo)

	// ─── Handlers ──────────────────────────────────────────
	productHandler := handler.NewProductHandler(productSvc)

	// ─── Routes ────────────────────────────────────────────
	mux := http.NewServeMux()

	// Public
	mux.HandleFunc("/live", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	// Product routes
	mux.Handle("/products", productHandler)
	mux.Handle("/products/", productHandler)

	// ─── Server ────────────────────────────────────────────
	srv := &http.Server{
		Addr:    cfg.AppPort,
		Handler: mux,
	}

	return &APIServer{
		server: srv,
		cfg:    cfg,
	}
}

// Start — implement App interface
func (s *APIServer) Start() {
	fmt.Printf("[%s] Starting server on %s\n", s.cfg.AppName, s.cfg.AppPort)
	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		fmt.Printf("Server error: %v\n", err)
	}
}

// Stop — implement App interface
func (s *APIServer) Stop() {
	fmt.Printf("[%s] Stopping server...\n", s.cfg.AppName)
	if err := s.server.Close(); err != nil {
		fmt.Printf("Error stopping server: %v\n", err)
		return
	}
	fmt.Printf("[%s] Server stopped\n", s.cfg.AppName)
}

// loadConfig — โหลดจาก environment variable
func loadConfig() (AppConfig, error) {
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = ":8080" // default
	}

	name := os.Getenv("APP_NAME")
	if name == "" {
		name = "inventory-api"
	}

	return AppConfig{
		AppPort: port,
		AppName: name,
	}, nil
}
