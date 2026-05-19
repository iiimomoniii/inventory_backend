package db

import (
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/iiimomoniii/inventory_backend/config"
)

func buildDSN(cfg config.DatabaseConfig) string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.SSLMode,
	)
}

// RunMigrations — สร้าง connection แยกสำหรับ migration
// ไม่รับ *sql.DB เพื่อป้องกัน m.Close() ปิด connection ของ app
func RunMigrations(cfg config.DatabaseConfig) error {
	dsn := buildDSN(cfg)

	m, err := migrate.New("file://migrations", dsn)
	if err != nil {
		return fmt.Errorf("failed to create migrator: %w", err)
	}
	defer m.Close() // ปิดแค่ connection ของ migration เอง

	version, dirty, _ := m.Version()
	fmt.Printf("[migration] current version: %d, dirty: %v\n", version, dirty)

	if err := m.Up(); err != nil {
		if err == migrate.ErrNoChange {
			fmt.Println("[migration] no new migrations to run ✅")
			return nil
		}
		return fmt.Errorf("migration failed: %w", err)
	}

	newVersion, _, _ := m.Version()
	fmt.Printf("[migration] migrated successfully → version: %d ✅\n", newVersion)
	return nil
}

// RollbackMigration — rollback 1 step
func RollbackMigration(cfg config.DatabaseConfig) error {
	dsn := buildDSN(cfg)

	m, err := migrate.New("file://migrations", dsn)
	if err != nil {
		return fmt.Errorf("failed to create migrator: %w", err)
	}
	defer m.Close()

	version, _, _ := m.Version()
	fmt.Printf("[migration] rolling back from version: %d\n", version)

	if err := m.Steps(-1); err != nil {
		return fmt.Errorf("rollback failed: %w", err)
	}

	newVersion, _, _ := m.Version()
	fmt.Printf("[migration] rolled back → version: %d ✅\n", newVersion)
	return nil
}
