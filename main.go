package main

import (
	"flag"
	"fmt"
	"os"

	routes "github.com/iiimomoniii/inventory_backend/route"
)

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
