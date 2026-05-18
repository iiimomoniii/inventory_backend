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
	App   AppInfo     `mapstructure:"app"`
	Fiber FiberConfig `mapstructure:"fiber"`
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

// ─── Loader ────────────────────────────────────────────────

func LoadConfig() (AppConfig, error) {
	// โหลด .env ก่อน
	gotenv.Load()

	viper.SetConfigType("yaml")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "__", "-", "_"))

	// โหลดจาก config.yaml ที่ embed ไว้
	if err := viper.ReadConfig(bytes.NewBuffer(configFile)); err != nil {
		return AppConfig{}, err
	}

	var cfg AppConfig
	if err := viper.Unmarshal(&cfg); err != nil {
		return AppConfig{}, err
	}

	// Override ด้วย .env หรือ ENV variable
	overrideWithEnv(&cfg)

	return cfg, nil
}

// overrideWithEnv — ค่าจาก .env ทับค่าจาก yaml
func overrideWithEnv(cfg *AppConfig) {

	// ─── App ───────────────────────────────────────────────
	if name := viper.GetString("APP__NAME"); name != "" {
		cfg.App.Name = name
	}
	if env := viper.GetString("APP__ENVIRONMENT"); env != "" {
		cfg.App.Environment = env
	}

	// ─── Fiber ─────────────────────────────────────────────
	if addr := viper.GetString("FIBER__ADDRESS"); addr != "" {
		cfg.Fiber.Address = addr
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
}
