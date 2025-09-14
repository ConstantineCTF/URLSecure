# Stage 1: build
FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o urlsecure cmd/shortener/main.go

# Stage 2: runtime
FROM alpine:latest
WORKDIR /root/

# Copy binary
COPY --from=builder /app/urlsecure .

# Copy migrations
COPY --from=builder /app/migrations ./migrations

# Copy frontend static files
COPY --from=builder /app/public ./public

# Copy environment example as .env
COPY --from=builder /app/.env.example .env

CMD ["./urlsecure"]
