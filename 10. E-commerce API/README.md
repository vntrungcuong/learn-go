# High-Performance E-commerce API (Go + Fiber)

![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)
![Framework](https://img.shields.io/badge/Fiber-v2-black?style=flat&logo=gofiber)
![Database](https://img.shields.io/badge/PostgreSQL-16-336791?style=flat&logo=postgresql)
![Platform](https://img.shields.io/badge/Platform-Windows%2011%20Native-0078D6?style=flat&logo=windows)

> **Project:** A high-throughput Backend API for a large-scale E-commerce platform.
> **Target:** Handling millions of concurrent users with sub-second latency (<100ms P99).
> **Infrastructure:** Native Windows 11 Server execution (No Docker/Container orchestration overhead).

## 📋 Overview

This project serves as the **Core API Engine**. It is architected to maximize raw hardware performance on Windows servers by bypassing the virtualization layer of Docker. It utilizes **Fiber** (built on `fasthttp`) for zero memory allocation routing and **pgx/v5** for high-performance database connection pooling.

### Key Features
* **High Concurrency:** Optimized for handling massive request throughput.
* **User Side:** Product browsing, Search, Cart management, Social Login (Google, FB, Apple).
* **Admin Side:** CMS for Products/Categories, Revenue Analytics.
* **Security:** JWT Authentication, Bcrypt Password Hashing, SQL Injection protection via parameterized queries.

---

## 🛠 Tech Stack & Architecture

The project follows **Clean Architecture** principles to ensure scalability and maintainability.

* **Language:** [Go (Golang)](https://go.dev/) (v1.21+)
* **Web Framework:** [Fiber v2](https://gofiber.io/) - Chosen for its superior performance over `net/http`.
* **Database:** PostgreSQL 16.
* **Driver:** [pgx/v5](https://github.com/jackc/pgx) - Native Go driver supporting binary format and efficient connection pooling.
* **OS:** Windows 11 (Bare metal).

### Project Structure
```text
/ecommerce-api
├── cmd
│   └── server
│       └── main.go        # Application Entry Point
├── internal
│   ├── config             # DB Connection & Env Config
│   ├── handlers           # HTTP Controllers (Fiber Handlers)
│   ├── models             # Data Structures & DTOs
│   ├── repository         # Data Access Layer (Raw SQL execution)
│   └── services           # Business Logic
├── .env                   # Environment Variables (Gitignored)
├── go.mod                 # Module definition
└── README.md              # Documentation
```

---
## 🚀 Installation & Setup

### Step 1: Clone and Initialize
Open your terminal (PowerShell or CMD) in the project root:

```powershell
# 1. Initialize module
go mod init ecommerce-api

# 2. Install High-Performance dependencies
go get github.com/gofiber/fiber/v2       # Web Framework
go get github.com/jackc/pgx/v5           # DB Driver
go get github.com/golang-jwt/jwt/v5      # JWT Auth
go get golang.org/x/crypto/bcrypt        # Password Hashing
go get github.com/joho/godotenv          # Env Config

# 3. Clean up dependencies
go mod tidy

---
## Database Setup
```powershell
$env:PGCLIENTENCODING='utf-8'
psql -h localhost -U postgres -d ecommerce_db -f .\sql\schema.sql
```

# 1. Tables
* users
* user_providers
* categories
* products
* orders
* order_items

# 2. Seed Data

```powershell
$env:PGCLIENTENCODING='utf-8'
psql -h localhost -U postgres -d ecommerce_db -f .\sql\seed.sql
```