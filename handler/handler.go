package handler

import (
	"context"
	"fmt"
	"log"
	"math"
	"study/weatherbot/clients/openweather"
	"study/weatherbot/models"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type userRepository interface {
	GetUserCity(ctx context.Context, userID int64) (string, error)
	CreateUser(ctx context.Context, userID int64) error
	UpdateUserCity(ctx context.Context, userID int64, city string) error
	GetUser(ctx context.Context, userID int64) (*models.User, error)
}

type Handler struct {
	bot      *tgbotapi.BotAPI
	owClient *openweather.OpenWeatherClient
	userRepo userRepository
}

func New(bot *tgbotapi.BotAPI, owClient *openweather.OpenWeatherClient, userRepo userRepository) *Handler {
	return &Handler{
		bot:      bot,
		owClient: owClient,
		userRepo: userRepo,
	}
}

func (h *Handler) handleUpdate(update tgbotapi.Update) {
	if update.Message == nil {
		return
	}

	ctx := context.Background()

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

func (h *Handler) Start() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := h.bot.GetUpdatesChan(u)

	for update := range updates {
		h.handleUpdate(update)
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
	coordinate, err := h.owClient.Coordinates(city)
	if err != nil {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Не смогли получить координаты")
		msg.ReplyToMessageID = update.Message.MessageID
		h.bot.Send(msg)
		return
	}

	weather, err := h.owClient.Weather(coordinate.Lat, coordinate.Lon)
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
