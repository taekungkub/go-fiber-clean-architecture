package routes

import (
	"go-fiber-api/internal/api/handlers"
	"go-fiber-api/internal/api/middleware"
	"go-fiber-api/pkg/domain/book"

	"github.com/gofiber/fiber/v2"
)

// BookRouter is the Router for GoFiber App - all routes are protected with JWT
func BookRouter(app fiber.Router, service book.Service) {
	books := app.Group("/books", middleware.JWTMiddleware())

	books.Get("/", handlers.GetBooksWithPagination(service))
	books.Post("/", handlers.AddBook(service))
	books.Put("/", handlers.UpdateBook(service))
	books.Delete("/:id", handlers.RemoveBook(service))
	books.Post("/pagination", handlers.GetBooksWithPaginationBody(service))
}
