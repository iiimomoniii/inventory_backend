package config

import (
	"bytes"
	_ "embed"
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

// ─── Loader ────────────────────────────────────────────────

func LoadConfig() (AppConfig, error) {
	gotenv.Load()

	viper.SetConfigType("yaml")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "__", "-", "_"))

	if err := viper.ReadConfig(bytes.NewBuffer(configFile)); err != nil {
		return AppConfig{}, err
	}

	var cfg AppConfig
	if err := viper.Unmarshal(&cfg); err != nil {
		return AppConfig{}, err
	}

	overrideWithEnv(&cfg)
	return cfg, nil
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
}
