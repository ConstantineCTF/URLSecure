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
COPY --from=builder /app/urlsecure .
COPY --from=builder /app/migrations ./migrations
COPY --from=builder /app/.env.example .env
EXPOSE 8080
CMD ["./urlsecure"]
