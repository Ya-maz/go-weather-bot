# Project Development Plan (PLANS.md)

As a Senior Golang Developer, I recommend the following improvements to transition this project from a prototype to a production-ready application.

## Phase 1: Stability & Reliability (Critical)

### 1. Resource Management
- **Issue**: HTTP responses in `openweather.go` are not being closed.
- **Action**: Add `defer resp.Body.Close()` after successful HTTP calls.

### 2. Concurrency & Performance
- **Issue**: Sequential processing of Telegram updates blocks the entire bot.
- **Action**: Implement a worker pool or process each update in a separate goroutine (with a limit/semaphore).
- **Action**: Switch from `pgx.Conn` to `pgxpool.Pool` for managed database connections.

### 3. Context & Timeouts
- **Issue**: Lacking context support throughout the call chain.
- **Action**: Thread `context.Context` through all client and repository methods.
- **Action**: Set reasonable timeouts for external API calls and database queries.

## Phase 2: Architecture & Clean Code

### 4. Configuration Management
- **Issue**: Direct `os.Getenv` usage in `main.go`.
- **Action**: Create a `config` package to load and validate all environment variables into a structured object.

### 5. Graceful Shutdown
- **Issue**: Bot stops abruptly on SIGINT/SIGTERM.
- **Action**: Listen for OS signals and ensure the DB connection closes and the bot finishes processing current updates before exiting.

### 6. Dependency Injection & Interfaces
- **Issue**: Handler depends on a concrete `OpenWeatherClient`.
- **Action**: Define an interface for the weather service to allow easier testing and flexibility.

## Phase 3: Quality & Security

### 7. Automated Testing
- **Issue**: No tests in the project.
- **Action**: Add unit tests for `handler`, `repo`, and `clients` using mocks.
- **Action**: Implement table-driven tests for weather logic.

### 8. Input Validation
- **Issue**: City name is saved without validation.
- **Action**: Add basic validation or normalization for city names before saving them to the database.

### 9. Linting & Formatting
- **Action**: Integrated `golangci-lint` into the development workflow to ensure consistent code quality.

## Future Features
- **Daily Notifications**: Send weather updates to users at a scheduled time.
- **Multiple Locations**: Allow users to save and check weather for multiple cities.
- **Extended Forecast**: Provide 3-day or 7-day weather forecasts.
