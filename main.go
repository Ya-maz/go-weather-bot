package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"study/weatherbot/clients/openweather"
	"study/weatherbot/config"
	"study/weatherbot/handler"
	"study/weatherbot/repo"
	"syscall"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v\n", err)
	}
	defer pool.Close()

	err = pool.Ping(ctx)
	if err != nil {
		log.Fatal("Error ping pool")
	}

	bot, err := tgbotapi.NewBotAPI(cfg.BotToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	owClient := openweather.New(cfg.OpenWeatherAPIKey)

	userRepo := repo.New(pool)

	botHandler := handler.New(bot, owClient, userRepo)

	botHandler.Start(ctx)
}
