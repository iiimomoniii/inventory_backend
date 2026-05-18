package routes

import (
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

	// 2. รัน migration
	if err := db.RunMigrations(cfg.Database); err != nil {
		fmt.Printf("Migration warning: %v\n", err)
		// ไม่ panic เพื่อให้ยังรันได้กับ in-memory repo
	}

	// 3. สร้าง server
	apiServer := NewAPIServer(cfg)

	// 4. graceful shutdown
	var wg sync.WaitGroup
	wg.Add(1)
	addShutdownHook(&wg, func() {
		apiServer.Stop()
	})

	return apiServer, &wg
}

func addShutdownHook(wg *sync.WaitGroup, f func()) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		defer wg.Done()
		<-c
		fmt.Println("\nShutting down...")
		f()
	}()
}
