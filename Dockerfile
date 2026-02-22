# Build stage
FROM golang:1.24-alpine AS builder

RUN apk add --no-cache git ca-certificates && update-ca-certificates

WORKDIR /app

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source and build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/weatherbot main.go

# Install goose for migrations
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

# Final stage
FROM alpine:latest

RUN apk add --no-cache ca-certificates tzdata

# Create a non-root user
RUN adduser -D -u 1000 appuser
WORKDIR /app

# Copy necessary files from builder
COPY --from=builder /app/weatherbot .
COPY --from=builder /app/migrations ./migrations
COPY --from=builder /go/bin/goose /usr/local/bin/goose

RUN chown -R appuser:appuser /app
USER appuser

# Run migrations and start the bot
ENTRYPOINT ["/bin/sh", "-c", "goose -dir ./migrations postgres \"$DATABASE_URL\" up && exec ./weatherbot"]
