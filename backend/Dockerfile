# Stage 1: build stage using Go 1.25 alpine image
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy modules files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy entire backend source code
COPY . .

# Build the Go binary statically for Linux
RUN CGO_ENABLED=0 GOOS=linux go build -o urlsecure cmd/shortener/main.go

# Stage 2: lightweight runtime image
FROM alpine:latest

WORKDIR /root/

# Copy the compiled binary from builder stage
COPY --from=builder /app/urlsecure .

# Copy database migrations
COPY --from=builder /app/migrations ./migrations

# Copy frontend static assets
COPY --from=builder /app/public ./public

# Command to start the backend server
CMD ["./urlsecure"]
