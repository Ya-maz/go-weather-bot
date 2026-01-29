package handler

import (
	"context"
	"fmt"
	"log"
	"math"
	"study/weatherbot/clients/openweather"
	"study/weatherbot/models"
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type userRepository interface {
	GetUserCity(ctx context.Context, userID int64) (string, error)
	CreateUser(ctx context.Context, userID int64) error
	UpdateUserCity(ctx context.Context, userID int64, city string) error
	GetUser(ctx context.Context, userID int64) (*models.User, error)
}

type weatherProvider interface {
	Coordinates(ctx context.Context, city string) (openweather.Coordinate, error)
	Weather(ctx context.Context, lat float64, lon float64) (openweather.Weather, error)
}

type Handler struct {
	bot        *tgbotapi.BotAPI
	owProvider weatherProvider
	userRepo   userRepository
}

func New(bot *tgbotapi.BotAPI, owProvider weatherProvider, userRepo userRepository) *Handler {
	return &Handler{
		bot:        bot,
		owProvider: owProvider,
		userRepo:   userRepo,
	}
}

func (h *Handler) handleUpdate(ctx context.Context, update tgbotapi.Update) {
	if update.Message == nil {
		return
	}

	if update.Message.IsCommand() {
		err := h.ensureUser(ctx, update)
		if err != nil {
			log.Println("error h.ensureUser: ", err)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Произошлла ошибка")
			msg.ReplyToMessageID = update.Message.MessageID
			h.bot.Send(msg)
			return
		}
		switch update.Message.Command() {
		case "city":
			h.handleSetCity(ctx, update)
			return
		case "weather":
			h.handleSendWeather(ctx, update)
			return
		default:
			h.handleUnknownCommand(update)
			return
		}
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Воcпользуйтесь доступными командами")
	msg.ReplyToMessageID = update.Message.MessageID
	h.bot.Send(msg)
}

func (h *Handler) Start(ctx context.Context) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := h.bot.GetUpdatesChan(u)

	var wg sync.WaitGroup

	log.Println("Bot handler started...")

	for {
		select {
		case <-ctx.Done():
			log.Println("Stopping bot handler...")
			h.bot.StopReceivingUpdates()
			wg.Wait()
			log.Println("Bot handler stopped gracefully.")
			return
		case update, ok := <-updates:
			if !ok {
				wg.Wait()
				return
			}
			wg.Add(1)
			go func(upd tgbotapi.Update) {
				defer wg.Done()
				h.handleUpdate(ctx, upd)
			}(update)
		}
	}
}

func (h *Handler) handleSetCity(ctx context.Context, update tgbotapi.Update) {
	city := update.Message.CommandArguments()
	err := h.userRepo.UpdateUserCity(ctx, update.Message.From.ID, city)
	if err != nil {
		log.Println("error userRepo.updateUserCity: ", err)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Ошибка при попытке сохранения города - %s", city))
		msg.ReplyToMessageID = update.Message.MessageID
		h.bot.Send(msg)
	}
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Город %s сохранен", city))
	msg.ReplyToMessageID = update.Message.MessageID
	h.bot.Send(msg)
}

func (h *Handler) handleSendWeather(ctx context.Context, update tgbotapi.Update) {
	city, err := h.userRepo.GetUserCity(ctx, update.Message.From.ID)
	if err != nil {
		log.Println("error userRepo.updateUserCity: ", err)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Произошлла ошибка")
		msg.ReplyToMessageID = update.Message.MessageID
		h.bot.Send(msg)
		return
	}

	if city == "" {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Сначала сохраните ваш город - /city <your city>")
		msg.ReplyToMessageID = update.Message.MessageID
		h.bot.Send(msg)
		return
	}

	log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

	weatherCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	coordinate, err := h.owProvider.Coordinates(weatherCtx, city)
	if err != nil {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Не смогли получить координаты")
		msg.ReplyToMessageID = update.Message.MessageID
		h.bot.Send(msg)
		return
	}

	weather, err := h.owProvider.Weather(weatherCtx, coordinate.Lat, coordinate.Lon)
	if err != nil {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Не смогли получить погоду в этой местности")
		msg.ReplyToMessageID = update.Message.MessageID
		h.bot.Send(msg)
		return
	}

	msg := tgbotapi.NewMessage(
		update.Message.Chat.ID,
		fmt.Sprintf("Температура в вашем городе \n%s: %d°C", city, int(math.Round(weather.Temp))),
	)
	msg.ReplyToMessageID = update.Message.MessageID
	h.bot.Send(msg)
}

func (h *Handler) handleUnknownCommand(update tgbotapi.Update) {
	log.Printf("Unknown command - [%s] %s", update.Message.From.UserName, update.Message.Text)
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Такая команда не доступна")
	msg.ReplyToMessageID = update.Message.MessageID
	h.bot.Send(msg)
}

func (h *Handler) ensureUser(ctx context.Context, update tgbotapi.Update) error {
	user, err := h.userRepo.GetUser(ctx, update.Message.From.ID)
	if err != nil {
		return fmt.Errorf("error userRepo.GetUser: %w", err)
	}

	if user == nil {
		err = h.userRepo.CreateUser(ctx, update.Message.From.ID)

		if err != nil {
			return fmt.Errorf("error userRepo.CreateUser: %w", err)
		}

	}
	return nil
}
