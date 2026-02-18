# --- Stage 1: Base & Build Dependencies ---
FROM golang:1.23-alpine AS base
WORKDIR /app
RUN apk add --no-cache git ca-certificates tzdata
COPY go.mod go.sum ./
RUN go mod download

# --- Stage 2: Development (Hot-reload with Air) ---
FROM base AS development
# Cài đặt Air cho hot-reloading
RUN go install github.com/air-verse/air@latest
COPY . .
# Chạy Air để tự động restart khi code thay đổi
CMD ["air", "-c", ".air.toml"]

# --- Stage 3: Builder (Compile for Production) ---
FROM base AS builder
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/api/main.go

# --- Stage 4: Production (Smallest Image) ---
FROM alpine:3.20 AS production
WORKDIR /app
RUN apk add --no-cache ca-certificates tzdata
COPY --from=builder /app/main .

# Quan trọng: Không copy .env. Biến môi trường sẽ được truyền qua Compose hoặc Cloud Provider.
EXPOSE 8080
CMD ["./main"]
