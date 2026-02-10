package routes

import (
	"go-fiber-api/internal/api/handlers"
	"go-fiber-api/pkg/domain/user"

	"github.com/gofiber/fiber/v2"
)

func UserRouter(app fiber.Router, service user.Service) {
	users := app.Group("/users")

	users.Get("/", handlers.GetUsers(service))
	users.Post("/", handlers.CreateUser(service))
	users.Put("/", handlers.UpdateUser(service))
	users.Delete("/:id", handlers.RemoveUser(service))
}
