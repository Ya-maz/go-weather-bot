package main

import (
	"context"
	"log"
	"os"
	"study/weatherbot/clients/openweather"
	"study/weatherbot/handler"
	"study/weatherbot/repo"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

    err = conn.Ping(context.Background())
    if err != nil {
        log.Fatal("Error ping db")
    }

	bot, err := tgbotapi.NewBotAPI(os.Getenv("BOT_TOKEN"))
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	owClient := openweather.New(os.Getenv("OPEN_WEATHER_API_KEY"))

    userRepo := repo.New(conn)

	botHandler := handler.New(bot, owClient, userRepo)

	botHandler.Start()
}
