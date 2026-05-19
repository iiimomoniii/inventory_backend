package routes

import (
	"database/sql"
	"fmt"
	"os"
	"os/signal"
	"sync"

	"github.com/iiimomoniii/inventory_backend/config"
	"github.com/iiimomoniii/inventory_backend/db"
)

type App interface {
	Start()
	Stop()
}

func Bootstrap() (App, *sync.WaitGroup) {

	// 1. โหลด config
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Println("Error loading config:", err)
	}

	// 2. รัน migration ด้วย connection แยก (ไม่ใช้ sqlDB ของ app)
	if err := db.RunMigrations(cfg.Database); err != nil {
		fmt.Printf("[bootstrap] Migration warning: %v\n", err)
	}

	// 3. Connect database สำหรับ app
	sqlDB, err := db.NewDatabase(cfg.Database)
	if err != nil {
		fmt.Printf("[bootstrap] ❌ DB connect failed: %v\n", err)
	} else {
		fmt.Println("[bootstrap] ✅ DB ready")
	}

	// 4. สร้าง server
	apiServer := NewAPIServer(cfg, sqlDB)

	// 5. graceful shutdown
	var wg sync.WaitGroup
	wg.Add(1)
	addShutdownHook(&wg, sqlDB, func() {
		apiServer.Stop()
	})

	return apiServer, &wg
}

func addShutdownHook(wg *sync.WaitGroup, sqlDB *sql.DB, f func()) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		defer wg.Done()
		<-c
		fmt.Println("\nShutting down...")
		f()
		if sqlDB != nil {
			sqlDB.Close()
			fmt.Println("[bootstrap] DB closed")
		}
	}()
}
