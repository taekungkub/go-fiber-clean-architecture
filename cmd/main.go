package main

import (
	"fmt"
	"go-fiber-api/internal/api/middleware"
	"go-fiber-api/internal/api/routes"
	"go-fiber-api/internal/config"
	"go-fiber-api/internal/database"
	"go-fiber-api/pkg/core"
	"go-fiber-api/pkg/domain/book"
	"go-fiber-api/pkg/domain/user"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	cfg := config.LoadConfig()

	fmt.Println("APP_ENV", cfg.APP_ENV)

	db, cancel, err := database.DBConnect(cfg)

	if err != nil {
		log.Fatal("Database Connection Error $s", err)
	}
	fmt.Println("Database connection success!")

	// Initialize JWT configuration
	core.JWTSecret = []byte(cfg.JWT_SECRET)
	core.JWTExpiry = time.Hour * time.Duration(cfg.JWT_EXPIRY)

	// Initialize book service
	bookCollection := db.Collection("books")
	bookRepo := book.NewRepo(bookCollection)
	bookService := book.NewService(bookRepo)

	// Initialize user service
	userCollection := db.Collection("users")
	userRepo := user.NewRepo(userCollection)
	userService := user.NewService(userRepo)

	app := fiber.New()
	app.Use(cors.New())
	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.Send([]byte("Welcome to the clean-architecture mongo book shop!"))
	})

	api := app.Group("/api")

	// Public auth routes (register, login)
	routes.AuthRouter(api, userService)

	// Protected book routes (require JWT)
	routes.BookRouter(api, bookService)
	routes.UserRouter(api, userService)

	// Protected auth routes (profile)
	protectedApi := api.Group("", middleware.JWTMiddleware())
	routes.ProtectedAuthRouter(protectedApi, userService)

	defer cancel()
	log.Fatal(app.Listen(":8080"))

}
