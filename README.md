# Inventory_Backend

Go + Fiber + PostgreSQL

## Project Structure

```
inventory_backend/
├── config/
│   └── config.go
├── handler/
│   └── product_handler.go
├── locales/
│   ├── global.json
│   └── product.json
├── middleware/
│   ├── auth.go
│   ├── cors.go
│   └── i18n.go
├── model/
│   ├── product.go
│   └── response.go
├── repository/
│   ├── product_repository.go
│   └── inmemory_repository.go
├── route/
│   ├── bootstrap.go
│   └── route_apiserver.go
├── service/
│   ├── product_service.go
│   └── product_service_test.go
├── utils/
│   └── utils.go
├── go.mod
└── main.go
```

## Install Go Packages

```bash
# HTTP Framework
go get github.com/gofiber/fiber/v2

# JWT
go get github.com/golang-jwt/jwt/v4

#
go get gopkg.in/yaml.v3

go get github.com/spf13/viper

go get github.com/subosito/gotenv

go get github.com/golang-migrate/migrate/v4
go get github.com/golang-migrate/migrate/v4/database/postgres
go get github.com/golang-migrate/migrate/v4/source/file

go get github.com/lib/pq

go get golang.org/x/crypto/bcrypt

go get github.com/gofiber/fiber/v2/middleware/limiter@v2.52.13

go install github.com/swaggo/swag/cmd/swag@latest
go get github.com/swaggo/fiber-swagger
go get github.com/swaggo/files
# ตรวจสอบ go.mod และ go.sum
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
| POST | /products | Create product (single) |
| POST | /products/create/items | Create products (bulk) |
| PUT | /products/:id | Update product |
| DELETE | /products/:id | Delete product |

## Error Code

### Global (GLB)

| Code | Name | EN | TH |
|---|---|---|---|
| GLB001 | NOT_FOUND | Record not found | ไม่พบข้อมูล |
| GLB002 | INTERNAL_ERROR | Internal server error | เกิดข้อผิดพลาดภายในระบบ |
| GLB003 | INVALID_BODY | Invalid request body | รูปแบบข้อมูลไม่ถูกต้อง |
| GLB004 | UNAUTHORIZED | Unauthorized | ไม่มีสิทธิ์เข้าถึง |

### Product (PRD)

| Code | Name | EN | TH |
|---|---|---|---|
| PRD001 | NAME_REQUIRED | Name is required | กรุณากรอกชื่อสินค้า |
| PRD002 | CATEGORY_REQUIRED | Category is required | กรุณาเลือกหมวดหมู่ |
| PRD003 | PRICE_INVALID | Price must be greater than 0 | ราคาต้องมากกว่า 0 |
| PRD004 | STOCK_INVALID | Stock must be >= 0 | จำนวนต้องไม่ติดลบ |
| PRD005 | INVALID_ID | Invalid ID | รหัสไม่ถูกต้อง |
| PRD006 | INVALID_PAGE | Invalid page number | หมายเลขหน้าไม่ถูกต้อง |
| PRD007 | INVALID_PRICE_RANGE | minPrice cannot be greater than maxPrice | ราคาต่ำสุดต้องไม่มากกว่าราคาสูงสุด |
| PRD008 | NO_FIELDS_TO_UPDATE | At least one field must be provided | กรุณาระบุข้อมูลที่ต้องการอัปเดต |
| PRD009 | PARTIAL_FAILED | Some items failed to process | บางรายการดำเนินการไม่สำเร็จ |

## Response Format

### Single

```json
{
  "status": 400,
  "message": "Name is required",
  "data": {
    "timestamp": "2026-05-18T11:00:00+07:00",
    "status": 400,
    "errorCode": "PRD001",
    "error": "NAME_REQUIRED",
    "messageEN": "Name is required",
    "messageTH": "กรุณากรอกชื่อสินค้า",
    "path": "/products"
  }
}
```

### Bulk (POST /products/create/items)

```json
{
  "status": 207,
  "message": "Some items failed to process",
  "data": {
    "timestamp": "2026-05-18T11:00:00+07:00",
    "status": 207,
    "totalSuccess": 2,
    "totalError": 1,
    "results": [
      { "index": 0, "status": "success", "data": { "id": 1, "name": "Apple" } },
      { "index": 1, "status": "error", "errorCode": "PRD001", "error": "NAME_REQUIRED", "messageEN": "Name is required", "messageTH": "กรุณากรอกชื่อสินค้า" }
    ],
    "path": "/products/create/items"
  }
}
```

### Bulk Status

| Status | HTTP | เมื่อไหร่ |
|---|---|---|
| success | 200 | ทุกอันสำเร็จ |
| Partial success | 207 | บางอันสำเร็จ บางอันผิด |
| All items failed | 400 | ทุกอันผิด |

## เพิ่ม Module ใหม่

```
1. สร้าง locales/user.json   ← USR001, USR002...
2. เพิ่ม LoadLocaleFile ใน middleware/i18n.go
3. สร้าง model, repository, service, handler
4. เพิ่ม route ใน route_apiserver.go
```
