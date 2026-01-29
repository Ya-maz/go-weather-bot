package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	BotToken          string
	OpenWeatherAPIKey string
	DatabaseURL       string
}

func Load() (*Config, error) {
	// Load .env file if it exists
	_ = godotenv.Load()

	cfg := &Config{
		BotToken:          os.Getenv("BOT_TOKEN"),
		OpenWeatherAPIKey: os.Getenv("OPEN_WEATHER_API_KEY"),
		DatabaseURL:       os.Getenv("DATABASE_URL"),
	}

	if cfg.BotToken == "" {
		return nil, fmt.Errorf("BOT_TOKEN is required")
	}
	if cfg.OpenWeatherAPIKey == "" {
		return nil, fmt.Errorf("OPEN_WEATHER_API_KEY is required")
	}
	if cfg.DatabaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}

	return cfg, nil
}
