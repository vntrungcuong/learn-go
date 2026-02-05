package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

// DB là biến toàn cục (Singleton) để các package khác gọi tới
var DB *pgxpool.Pool

// ConnectDB khởi tạo kết nối PostgreSQL với Pool tối ưu
func ConnectDB() {
	// 1. Load file .env
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found")
	}

	// 2. Lấy thông tin kết nối từ env
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	sslMode := os.Getenv("SSL_MODE")

	// 3. Tạo connection string
	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		dbUser, dbPassword, dbHost, dbPort, dbName, sslMode)

	// 4. Parse config from DSN
	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		log.Fatal("Unable to parse connection string:", err)
	}

	// 5. TUNNING CONNECTION POOL (Quan trọng cho High Concurrency)

	// MaxConns: Số lượng kết nối tối đa
	// Lưu ý: Phải nhỏ hơn 'max_connections' trong file postgresql.conf
	maxConns, _ := strconv.Atoi(os.Getenv("DB_MAX_CONNS"))

	if maxConns == 0 {
		maxConns = 50 // Default
	}

	config.MaxConns = int32(maxConns)

	// MinConns: Số lượng kết nối duy trì tối thiểu (Warm start)
	minConns, _ := strconv.Atoi(os.Getenv("DB_MIN_CONNS"))

	if minConns == 0 {
		minConns = 5 // Default
	}

	config.MinConns = int32(minConns)

	// Thời gian sống tối đa của một kết nối (tránh issue memory leak bên DB)
	config.MaxConnLifetime = time.Hour

	// Thời gian tối đa một kết nối rảnh rỗi trước khi bị đóng
	config.MaxConnIdleTime = time.Minute * 30

	// 6. Create pool connection
	// Dùng context.Background() vì kết nối DB tồn tại suốt vòng đời ứng dụng
	DB, err = pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Fatal("Unable to create connection pool:", err)
	}

	// 7. Test connection
	if err := DB.Ping(context.Background()); err != nil {
		log.Fatal("Unable to ping database:", err)
	}

	log.Println("Database connection pool created successfully")
}

func ClostDB() {
	if DB != nil {
		DB.Close()

		fmt.Println("Database connection pool closed successfully")
	}
}
