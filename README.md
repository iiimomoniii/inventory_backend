# Inventory_Backend

Go + Fiber + PostgreSQL

## Project Structure

inventory_backend/
├── config/
│   └── config.go
├── handler/
│   └── product_handler.go
├── model/
│   └── product.go
├── repository/
│   ├── product_repository.go
│   └── inmemory_repository.go
├── routes/
│   ├── bootstrap.go
│   └── route_apiserver.go
├── service/
│   ├── product_service.go
│   └── product_service_test.go
├── go.mod
└── main.go

## Install Go Packages

```bash
# 1. HTTP Framework
go get github.com/gofiber/fiber/v2

# 2. ตรวจสอบว่า go.mod และ go.sum อัพเดตแล้ว
go mod tidy
```

## Run

```bash
go run main.go
```

## Environment Variables

| ตัวแปร | Default | คำอธิบาย |
|---|---|---|
| APP_PORT | :8080 | port ที่ server รัน |
| APP_NAME | inventory-api | ชื่อ app |

## API Endpoints

| Method | Path | Description |
|---|---|---|
| GET | /live | Health check |
| POST | /products/search | Search products |
| GET | /products/:id | Get product by ID |
| POST | /products | Create product |
| PUT | /products/:id | Update product |
| DELETE | /products/:id | Delete product |