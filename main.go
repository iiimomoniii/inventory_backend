package main

import (
	"flag"
	"fmt"
	"os"

	routes "github.com/iiimomoniii/inventory_backend/route"
)

// @title Inventory Backend API
// @version 1.0
// @description Inventory Backend API documentation
// @host localhost:8080
// @BasePath /
// @schemes http
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	// รับ -env flag จาก command line
	env := flag.String("env", "dev", "environment: dev, qa, uat, prod")
	flag.Parse()

	// เซ็ต APP_ENV ให้ config โหลดถูกไฟล์
	os.Setenv("APP_ENV", *env)
	fmt.Printf("[main] environment: %s\n", *env)

	app, wg := routes.Bootstrap()
	go app.Start()
	wg.Wait()
}
