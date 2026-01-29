package config

import (
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	// Save original env vars and restore after test
	origBotToken := os.Getenv("BOT_TOKEN")
	origWeatherKey := os.Getenv("OPEN_WEATHER_API_KEY")
	origDBURL := os.Getenv("DATABASE_URL")
	defer func() {
		os.Setenv("BOT_TOKEN", origBotToken)
		os.Setenv("OPEN_WEATHER_API_KEY", origWeatherKey)
		os.Setenv("DATABASE_URL", origDBURL)
	}()

	tests := []struct {
		name    string
		envs    map[string]string
		wantErr bool
	}{
		{
			name: "Success",
			envs: map[string]string{
				"BOT_TOKEN":            "test_token",
				"OPEN_WEATHER_API_KEY": "test_key",
				"DATABASE_URL":         "postgres://localhost:5432/test",
			},
			wantErr: false,
		},
		{
			name: "Missing Bot Token",
			envs: map[string]string{
				"OPEN_WEATHER_API_KEY": "test_key",
				"DATABASE_URL":         "postgres://localhost:5432/test",
			},
			wantErr: true,
		},
		{
			name: "Missing API Key",
			envs: map[string]string{
				"BOT_TOKEN":    "test_token",
				"DATABASE_URL": "postgres://localhost:5432/test",
			},
			wantErr: true,
		},
		{
			name: "Missing DB URL",
			envs: map[string]string{
				"BOT_TOKEN":            "test_token",
				"OPEN_WEATHER_API_KEY": "test_key",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear relevant env vars
			os.Unsetenv("BOT_TOKEN")
			os.Unsetenv("OPEN_WEATHER_API_KEY")
			os.Unsetenv("DATABASE_URL")

			// Set test env vars
			for k, v := range tt.envs {
				os.Setenv(k, v)
			}

			cfg, err := Load()
			if (err != nil) != tt.wantErr {
				t.Errorf("Load() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if cfg.BotToken != tt.envs["BOT_TOKEN"] {
					t.Errorf("got BotToken %v, want %v", cfg.BotToken, tt.envs["BOT_TOKEN"])
				}
			}
		})
	}
}
