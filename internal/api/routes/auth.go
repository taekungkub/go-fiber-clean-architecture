package routes

import (
	"go-fiber-api/internal/api/handlers"
	"go-fiber-api/pkg/domain/user"

	"github.com/gofiber/fiber/v2"
)

// AuthRouter sets up authentication routes
func AuthRouter(app fiber.Router, service user.Service) {
	auth := app.Group("/auth")

	// Public routes
	auth.Post("/register", handlers.Register(service))
	auth.Post("/login", handlers.Login(service))
}

// ProtectedAuthRouter sets up protected authentication routes
func ProtectedAuthRouter(app fiber.Router, service user.Service) {
	auth := app.Group("/auth")

	// Protected routes
	auth.Get("/profile", handlers.GetProfile(service))

}
