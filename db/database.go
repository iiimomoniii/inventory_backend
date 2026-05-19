package db

import (
	"database/sql"
	"fmt"

	"github.com/iiimomoniii/inventory_backend/config"
	_ "github.com/lib/pq" // postgres driver
)

// NewDatabase — สร้าง native sql.DB connection
func NewDatabase(cfg config.DatabaseConfig) (*sql.DB, error) {
	dsn := buildDSN(cfg)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// ทดสอบ connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(10)

	fmt.Printf("[database] connected to %s:%d/%s ✅\n", cfg.Host, cfg.Port, cfg.DBName)
	return db, nil
}
