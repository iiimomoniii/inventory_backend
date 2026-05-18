package routes

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
)

// App interface — กำหนดว่า server ต้องทำอะไรได้บ้าง
// เหมือนของบริษัท ใครก็ implement ได้ถ้ามี Start() Stop()
type App interface {
	Start()
	Stop()
}

// Bootstrap — wire dependencies ทั้งหมด แล้วส่ง App กลับไป
// เหมือน main.go ที่ประกอบทุกอย่างเข้าด้วยกัน
func Bootstrap() (App, *sync.WaitGroup) {

	// 1. Load config จาก env
	cfg, err := loadConfig()
	if err != nil {
		fmt.Println("Error loading config:", err)
	}

	// 2. สร้าง server พร้อม wire dependencies ทั้งหมด
	apiServer := NewAPIServer(cfg)

	// 3. ตั้ง WaitGroup รอ graceful shutdown
	var wg sync.WaitGroup
	wg.Add(1)
	addShutdownHook(&wg, func() {
		apiServer.Stop()
	})

	return apiServer, &wg
}

// addShutdownHook — รอรับ signal Ctrl+C แล้วรัน f()
// เหมือนของบริษัทเลยครับ
func addShutdownHook(wg *sync.WaitGroup, f func()) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		defer wg.Done()
		<-c // รอ Ctrl+C
		fmt.Println("\nShutting down...")
		f()
	}()
}
