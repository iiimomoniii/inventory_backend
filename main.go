package main

import "github.com/yourname/inventory-api/routes"

func main() {
	// Bootstrap wire ทุกอย่าง แล้วส่ง server กลับมา
	// เหมือนของบริษัทเลย — main.go ไม่รู้จัก detail ใดๆ
	app, wg := routes.Bootstrap()

	// Start server (non-blocking)
	go app.Start()

	// รอ Ctrl+C แล้ว graceful shutdown
	wg.Wait()
}
