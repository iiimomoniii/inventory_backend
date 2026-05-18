package config

import "os"

type AppConfig struct {
	AppPort string
	AppName string
}

func LoadConfig() (AppConfig, error) {
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = ":8080"
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
