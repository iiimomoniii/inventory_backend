package config

import (
	"bytes"
	_ "embed"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
	"github.com/subosito/gotenv"
)

//go:embed config.yaml
var configFile []byte

// ─── Struct ────────────────────────────────────────────────

type AppConfig struct {
	App      AppInfo        `mapstructure:"app"`
	Fiber    FiberConfig    `mapstructure:"fiber"`
	Database DatabaseConfig `mapstructure:"database"`
	JWT      JWTConfig      `mapstructure:"jwt"`
}

type AppInfo struct {
	Name        string `mapstructure:"name"`
	Environment string `mapstructure:"environment"`
}

type FiberConfig struct {
	Address        string `mapstructure:"address"`
	ReadTimeout    int    `mapstructure:"readTimeout"`
	WriteTimeout   int    `mapstructure:"writeTimeout"`
	IdleTimeout    int    `mapstructure:"idleTimeout"`
	ReadBufferSize int    `mapstructure:"readBufferSize"`
	BodyLimitSize  int    `mapstructure:"bodyLimitSize"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	SSLMode  string `mapstructure:"sslmode"`
}

type JWTConfig struct {
	Secret    string `mapstructure:"secret"`
	ExpiresIn int    `mapstructure:"expiresIn"`
}

// ─── Loader ────────────────────────────────────────────────

func LoadConfig() (AppConfig, error) {
	// โหลด .env ตาม APP_ENV
	env := os.Getenv("APP_ENV")
	envFile := loadEnvFile(env)

	viper.SetConfigType("yaml")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "__", "-", "_"))

	// โหลด config.yaml ก่อน
	if err := viper.ReadConfig(bytes.NewBuffer(configFile)); err != nil {
		return AppConfig{}, err
	}

	var cfg AppConfig
	if err := viper.Unmarshal(&cfg); err != nil {
		return AppConfig{}, err
	}

	// Override ด้วย .env
	overrideWithEnv(&cfg)

	fmt.Printf("[config] env=%s file=%s db=%s:%d/%s\n",
		env, envFile, cfg.Database.Host, cfg.Database.Port, cfg.Database.DBName)

	return cfg, nil
}

// loadEnvFile — โหลด .env file ตาม environment
func loadEnvFile(env string) string {
	envFiles := map[string]string{
		"dev":  ".env.dev",
		"qa":   ".env.qa",
		"uat":  ".env.uat",
		"prod": ".env.prod",
	}

	file, ok := envFiles[env]
	if !ok {
		file = ".env.dev" // default
	}

	if err := gotenv.Load(file); err != nil {
		fmt.Printf("[config] warning: %s not found, using default\n", file)
		gotenv.Load(".env.dev") // fallback
	}

	return file
}

func overrideWithEnv(cfg *AppConfig) {
	// ─── App ───────────────────────────────────────────────
	if v := viper.GetString("APP__NAME"); v != "" {
		cfg.App.Name = v
	}
	if v := viper.GetString("APP__ENVIRONMENT"); v != "" {
		cfg.App.Environment = v
	}

	// ─── Fiber ─────────────────────────────────────────────
	if v := viper.GetString("FIBER__ADDRESS"); v != "" {
		cfg.Fiber.Address = v
	}
	if v := viper.GetInt("FIBER__READ_TIMEOUT"); v != 0 {
		cfg.Fiber.ReadTimeout = v
	}
	if v := viper.GetInt("FIBER__WRITE_TIMEOUT"); v != 0 {
		cfg.Fiber.WriteTimeout = v
	}
	if v := viper.GetInt("FIBER__IDLE_TIMEOUT"); v != 0 {
		cfg.Fiber.IdleTimeout = v
	}
	if v := viper.GetInt("FIBER__READ_BUFFER_SIZE"); v != 0 {
		cfg.Fiber.ReadBufferSize = v
	}
	if v := viper.GetInt("FIBER__BODY_LIMIT_SIZE"); v != 0 {
		cfg.Fiber.BodyLimitSize = v
	}

	// ─── Database ──────────────────────────────────────────
	if v := viper.GetString("DATABASE__HOST"); v != "" {
		cfg.Database.Host = v
	}
	if v := viper.GetInt("DATABASE__PORT"); v != 0 {
		cfg.Database.Port = v
	}
	if v := viper.GetString("DATABASE__USER"); v != "" {
		cfg.Database.User = v
	}
	if v := viper.GetString("DATABASE__PASSWORD"); v != "" {
		cfg.Database.Password = v
	}
	if v := viper.GetString("DATABASE__DBNAME"); v != "" {
		cfg.Database.DBName = v
	}
	if v := viper.GetString("DATABASE__SSLMODE"); v != "" {
		cfg.Database.SSLMode = v
	}

	// ─── JWT ───────────────────────────────────────────────
	if v := viper.GetString("JWT__SECRET"); v != "" {
		cfg.JWT.Secret = v
	}
	if v := viper.GetInt("JWT__EXPIRES_IN"); v != 0 {
		cfg.JWT.ExpiresIn = v
	}
}
