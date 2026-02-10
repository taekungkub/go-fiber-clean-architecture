# Go Fiber API - Clean Architecture

## 📚 Clean Architecture Overview

Clean Architecture เป็นรูปแบบการออกแบบซอฟต์แวร์ที่แยก business logic ออกจาก framework และ infrastructure ทำให้โค้ดมีความยืดหยุ่น ทดสอบง่าย และบำรุงรักษาได้ดี

### 🎯 หลักการสำคัญ (Dependency Rule)

```
┌─────────────────────────────────────┐
│   Presentation Layer (Handler)      │  ← HTTP, gRPC, CLI
├─────────────────────────────────────┤
│   Business Logic Layer (Service)    │  ← Use Cases, Rules
├─────────────────────────────────────┤
│   Data Access Layer (Repository)    │  ← Database, Cache
└─────────────────────────────────────┘
```

**Dependency Rule:** ชั้นใน (inner layer) ไม่ควรรู้จักชั้นนอก (outer layer)
- Repository ไม่รู้จัก Service
- Service ไม่รู้จัก Handler
- Business Logic ไม่ผูกกับ Framework

---

## 🗂️ Project Structure

```
go-fiber-api/
├── cmd/
│   └── main.go              → Entry point, setup dependencies
├── internal/
│   ├── api/
│   │   ├── handlers/        → HTTP handlers (Presentation Layer)
│   │   ├── middleware/      → Authentication, logging, etc.
│   │   └── routes/          → Route definitions
│   ├── config/              → Configuration
├── pkg/
│   └── core/                → Shared utilities
│   └── dto/                 → Data Transfer Objects
│   ├── domain/
│   │   ├── book/          
│   │       └── repository   → Data Access Layer
│   │       └── service      → Business Logic Layer
│   └── entities/            → Business entities/models

└── go.mod
```

### 📁 Directory Explanation

| Directory | Purpose | Visibility |
|-----------|---------|------------|
| `cmd/` | Entry point (main.go) | Public |
| `internal/` | ใช้ได้เฉพาะในโปรเจกต์นี้เท่านั้น | Private |
| `pkg/` | Library / utilities ที่ reuse ได้ | Public |

---

## 🏗️ Clean Architecture Layers

### 1️⃣ Repository Layer (Data Access)

**Location:** `internal/repository/`

**หน้าที่:** 
- คุยกับ Database (CRUD operations)
- Query และ Transaction
- ไม่รู้เรื่อง HTTP, JSON, Fiber

**ตัวอย่าง:**
```go
// internal/repository/user_repository.go
type UserRepository interface {
    Create(user *entities.User) error
    FindByID(id uint) (*entities.User, error)
    FindByEmail(email string) (*entities.User, error)
    Update(user *entities.User) error
    Delete(id uint) error
}
```

**ทำไมต้องแยก Repository?**
- ✅ ใช้ซ้ำได้ (API, worker, cron job)
- ✅ Test ง่าย (mock database)
- ✅ ไม่ผูกกับ framework
- ✅ เปลี่ยน database ได้ง่าย

---

### 2️⃣ Service Layer (Business Logic)

**Location:** `internal/service/`

**หน้าที่:**
- Business rules และ validation
- Combine หลาย repository
- Use case logic
- ไม่รู้ว่า request มาจาก HTTP, gRPC หรือ CLI

**ตัวอย่าง:**
```go
// internal/service/user_service.go
type UserService interface {
    Register(dto *dto.RegisterRequest) error
    Login(dto *dto.LoginRequest) (string, error)
    GetProfile(userID uint) (*dto.UserResponse, error)
}
```

**ทำไมต้องแยก Service?**
- ✅ เป็น core ของระบบ
- ✅ ไม่ผูกกับ transport layer (HTTP/gRPC)
- ✅ เรียกใช้จาก handler ไหนก็ได้
- ✅ Test business logic แยกจาก HTTP

---

### 3️⃣ Handler Layer (Presentation)

**Location:** `internal/api/handlers/`

**หน้าที่:**
- รับ HTTP request
- Validate input (DTO)
- เรียก Service
- ส่ง HTTP response

**ตัวอย่าง:**
```go
// internal/api/handlers/user_handler.go
func (h *UserHandler) Register(c *fiber.Ctx) error {
    var req dto.RegisterRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
    }
    
    if err := h.userService.Register(&req); err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }
    
    return c.Status(201).JSON(fiber.Map{"message": "User registered"})
}
```

**ทำไมต้องแยก Handler?**
- ✅ แยก HTTP logic จาก business logic
- ✅ เปลี่ยน framework ได้ง่าย
- ✅ รองรับหลาย transport (HTTP, gRPC, WebSocket)

---

## 🚀 Quick Start

### 1️⃣ สร้างโปรเจกต์
```bash
mkdir go-fiber-api
cd go-fiber-api
go mod init go-fiber-api
```

### 2️⃣ ติดตั้ง Dependencies
```bash
go get github.com/gofiber/fiber/v2
go get gorm.io/gorm
go get gorm.io/driver/postgres
go get github.com/golang-jwt/jwt/v5
```

### 3️⃣ สร้างไฟล์ main.go
```go
// cmd/main.go
package main

import (
    "log"
    "github.com/gofiber/fiber/v2"
)

func main() {
    app := fiber.New()
    
    app.Get("/", func(c *fiber.Ctx) error {
        return c.SendString("Hello, Clean Architecture! 🏗️")
    })
    
    log.Fatal(app.Listen(":3000"))
}
```

### 4️⃣ รันโปรแกรม
```bash
go run cmd/main.go
```

### 5️⃣ ทดสอบ
เปิด Browser แล้วเข้า URL: http://localhost:3000  
คุณจะเห็นข้อความ: **Hello, Clean Architecture! 🏗️**

---

## 📋 Best Practices

### ✅ DO
- ใช้ Interface สำหรับ Repository และ Service
- แยก DTO (Data Transfer Object) จาก Entity
- Dependency Injection ผ่าน constructor
- Error handling ที่ชัดเจน
- Unit test แต่ละ layer แยกกัน

### ❌ DON'T
- ห้าม Repository เรียก Service
- ห้าม Repository รู้จัก HTTP
- ห้าม hardcode configuration
- ห้าม business logic ใน Handler
- ห้าม SQL query ใน Service

---

## 🔄 Data Flow Example

```
HTTP Request
    ↓
Handler (validate input)
    ↓
Service (business logic)
    ↓
Repository (database)
    ↓
Database
    ↓
Repository (return entity)
    ↓
Service (transform to DTO)
    ↓
Handler (return JSON)
    ↓
HTTP Response
```

---

## 📚 Additional Resources

- [Clean Architecture by Uncle Bob](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Go Fiber Documentation](https://docs.gofiber.io/)
- [GORM Documentation](https://gorm.io/docs/)

---

## 🐹 Go Commands Cheat Sheet

### 🔧 ตรวจสอบเวอร์ชัน
```bash
go version
```

### 📦 จัดการ Module
```bash
go mod init <module-name> # เริ่มต้นโปรเจกต์ใหม่
go mod tidy               # จัดการ dependencies ให้ตรงกับโค้ด
```

### การรันและสร้างไฟล์
```bash
go run .                  # รันโปรเจกต์ในโฟลเดอร์ปัจจุบัน
go build -o main .        # คอมไพล์โค้ดเป็นไฟล์ executable
```

### การทดสอบ (Testing)
```bash
go test ./...             # รัน test ทั้งหมดในโปรเจกต์
go test -v ./...          # รัน test พร้อมแสดงรายละเอียด
go test -cover ./...      # ตรวจสอบ test coverage
```

###  เครื่องมือจัดการโค้ด
```bash
go fmt ./...              # จัดฟอร์แมตโค้ดให้เป็นมาตรฐาน
go get <package-url>      # ดาวน์โหลดและติดตั้ง package
go vet ./...              # ตรวจสอบความผิดปกติของโค้ด
```




