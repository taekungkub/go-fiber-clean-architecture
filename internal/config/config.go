package config

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	APP_ENV    string `envconfig:"APP_ENV"`
	DB_URI     string `envconfig:"DB_URI"`
	JWT_SECRET string `envconfig:"JWT_SECRET" default:"your-secret-key-change-this-in-production"`
	JWT_EXPIRY int    `envconfig:"JWT_EXPIRY" default:"24"` // hours
}

func LoadConfig() *Config {
	// 1. หา env (development / production)
	_ = godotenv.Load(".env")

	// 2. ตั้งค่า APP_ENV
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "dev"
	}

	// 3. โหลด env-specific file
	_ = godotenv.Load(".env." + env)

	// 4. map env → struct
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return nil
	}

	return &cfg
}
