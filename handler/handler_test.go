package handler

import (
	"context"
	"study/weatherbot/clients/openweather"
	"study/weatherbot/models"
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type mockUserRepo struct {
	city string
	err  error
	user *models.User
}

func (m *mockUserRepo) GetUserCity(ctx context.Context, userID int64) (string, error) {
	return m.city, m.err
}
func (m *mockUserRepo) CreateUser(ctx context.Context, userID int64) error {
	return m.err
}
func (m *mockUserRepo) UpdateUserCity(ctx context.Context, userID int64, city string) error {
	m.city = city
	return m.err
}
func (m *mockUserRepo) GetUser(ctx context.Context, userID int64) (*models.User, error) {
	return m.user, m.err
}

type mockWeatherProvider struct {
	coord   openweather.Coordinate
	weather openweather.Weather
	err     error
}

func (m *mockWeatherProvider) Coordinates(ctx context.Context, city string) (openweather.Coordinate, error) {
	return m.coord, m.err
}
func (m *mockWeatherProvider) Weather(ctx context.Context, lat float64, lon float64) (openweather.Weather, error) {
	return m.weather, m.err
}

type mockBotAPI struct {
	sent []tgbotapi.Chattable
}

func (m *mockBotAPI) Send(c tgbotapi.Chattable) (tgbotapi.Message, error) {
	m.sent = append(m.sent, c)
	return tgbotapi.Message{}, nil
}
func (m *mockBotAPI) GetUpdatesChan(config tgbotapi.UpdateConfig) tgbotapi.UpdatesChannel {
	return nil
}
func (m *mockBotAPI) StopReceivingUpdates() {}

func TestHandler_HandleUpdate_SetCity(t *testing.T) {
	repo := &mockUserRepo{user: &models.User{ID: 1}}
	weather := &mockWeatherProvider{
		coord: openweather.Coordinate{Name: "Moscow", Lat: 55, Lon: 37},
	}
	bot := &mockBotAPI{}

	h := New(bot, weather, repo)

	update := tgbotapi.Update{
		Message: &tgbotapi.Message{
			From:     &tgbotapi.User{ID: 1},
			Chat:     &tgbotapi.Chat{ID: 1},
			Text:     "/city Moscow",
			Entities: []tgbotapi.MessageEntity{{Type: "bot_command", Offset: 0, Length: 5}},
		},
	}

	h.handleUpdate(context.Background(), update)

	if repo.city != "Moscow" {
		t.Errorf("got city %v, want Moscow", repo.city)
	}
	if len(bot.sent) != 1 {
		t.Errorf("got %d messages sent, want 1", len(bot.sent))
	}
}

func TestHandler_HandleSendWeather(t *testing.T) {
	repo := &mockUserRepo{user: &models.User{ID: 1}, city: "Moscow"}
	weather := &mockWeatherProvider{
		coord:   openweather.Coordinate{Lat: 55, Lon: 37},
		weather: openweather.Weather{Temp: 10.5},
	}
	bot := &mockBotAPI{}

	h := New(bot, weather, repo)

	update := tgbotapi.Update{
		Message: &tgbotapi.Message{
			From:     &tgbotapi.User{ID: 1},
			Chat:     &tgbotapi.Chat{ID: 1},
			Text:     "/weather",
			Entities: []tgbotapi.MessageEntity{{Type: "bot_command", Offset: 0, Length: 8}},
		},
	}

	h.handleSendWeather(context.Background(), update)

	if len(bot.sent) != 1 {
		t.Errorf("got %d messages sent, want 1", len(bot.sent))
	}
	msg := bot.sent[0].(tgbotapi.MessageConfig)
	if msg.Text != "Температура в вашем городе \nMoscow: 11°C" {
		t.Errorf("unexpected message text: %v", msg.Text)
	}
}

func TestHandler_HandleUpdate_EmptyCity(t *testing.T) {
	repo := &mockUserRepo{user: &models.User{ID: 1}}
	weather := &mockWeatherProvider{}
	bot := &mockBotAPI{}

	h := New(bot, weather, repo)

	update := tgbotapi.Update{
		Message: &tgbotapi.Message{
			From:     &tgbotapi.User{ID: 1},
			Chat:     &tgbotapi.Chat{ID: 1},
			Text:     "/city  ", // Empty arguments
			Entities: []tgbotapi.MessageEntity{{Type: "bot_command", Offset: 0, Length: 5}},
		},
	}

	h.handleUpdate(context.Background(), update)

	if repo.city != "" {
		t.Errorf("repo city should be empty, got %v", repo.city)
	}
	if len(bot.sent) != 1 {
		t.Errorf("got %d messages sent, want 1", len(bot.sent))
	}
	msg := bot.sent[0].(tgbotapi.MessageConfig)
	if msg.Text != "Пожалуйста, укажите город: /city <название>" {
		t.Errorf("unexpected message text: %v", msg.Text)
	}
}
