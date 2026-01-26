# GEMINI.md — Инструкции для AI coding agents

## Обзор проекта
Telegram-бот на Go для получения погоды (go-weather-bot).
Основной файл: main.go
Ключевые компоненты: инициализация бота, обработка команд, интеграция с weather API (OpenWeatherMap или аналог), конфиг из env.

## Tech Stack
- Язык: Go (последняя стабильная версия)
- Зависимости: go.mod (tgbotapi, HTTP-клиент для weather API, dotenv или flags)
- Без фреймворков (чистый stdlib + внешние либы)

## Команды Build / Lint / Test
- Сборка: `go build -o weather-bot main.go`
- Запуск: `go run main.go` (или `./weather-bot` после сборки)
- Тесты все: `go test ./... -v`
- Тест одного файла: `go test ./path/to/file_test.go -v`
- Тест одного теста: `go test -run ^ИмяТеста$ ./... -v`
- Линт: `go vet ./...` или `golangci-lint run` (если установлен)
- Форматирование: `go fmt ./...` и `goimports -w .`

## Правила стиля кода
- Именование: camelCase для переменных/функций, PascalCase для экспортируемых
- Импорты: стандартная библиотека → сторонние → локальные. Пустая строка между группами
- Обработка ошибок: всегда проверять err != nil, возвращать с контекстом (fmt.Errorf("... : %w", err))
- Логирование: использовать log/slog из stdlib
- Конкурентность: context.Context для отмены, sync.WaitGroup для goroutines
- Тесты: table-driven где возможно, покрывать happy path + edge cases (плохие API-ответы, пустой город)
- Комментарии: на русском или английском последовательно

## Архитектурные замечания
- Инициализация бота в main.go
- Хендлеры команд Telegram в отдельном пакете/файлах
- Вызов weather API с обработкой ошибок и таймаутами
- Конфиг: API-ключи из environment variablesац

## Правила для агентов
- Всегда писать тесты для новых фич и изменений
- Не добавлять новые зависимости без подтверждения
- Следовать существующему стилю кода
- Перед изменениями — делать план (в Plan mode)
- Проверять edge cases: нет сети, invalid city, rate limit API

## Потенциальные проблемы (для ревью)
- Обработка ошибок Telegram updates
- Защита от flood/spam
- Rate limiting запросов к weather API
- Graceful shutdown бота
