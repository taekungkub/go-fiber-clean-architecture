package handlers

import (
	"go-fiber-api/pkg/domain/book"
	"go-fiber-api/pkg/dto"
	"go-fiber-api/pkg/entities"

	"github.com/gofiber/fiber/v2"
)

// AddBook is handler/controller which creates Books in the BookShop
func AddBook(service book.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var createBookDTO dto.CreateBookDTO

		// Parse request body
		if err := c.BodyParser(&createBookDTO); err != nil {
			return dto.SendError(c, fiber.StatusBadRequest, "Invalid request body")
		}

		// Validate DTO
		if errors := dto.ValidateStruct(&createBookDTO); errors != nil {
			return dto.SendValidationError(c, errors)
		}

		// Convert DTO to entity
		bookEntity := &entities.Book{
			Title:  createBookDTO.Title,
			Author: createBookDTO.Author,
		}

		result, err := service.InsertBook(bookEntity)
		if err != nil {
			return dto.SendError(c, fiber.StatusBadRequest, err.Error())
		}

		return dto.SendSuccess(c, "Created book successfully", result)
	}
}

// GetBooks is handler/controller which fetches all books from the BookShop
func GetBooks(service book.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		result, err := service.GetBooks()
		if err != nil {
			return dto.SendError(c, fiber.StatusInternalServerError, "Failed to get books")
		}

		return dto.SendSuccess(c, "Retrieved books successfully", result)
	}
}

// UpdateBook is handler/controller which updates data of Books in the BookShop
func UpdateBook(service book.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var updateBookDTO dto.UpdateBookDTO

		// Parse request body
		if err := c.BodyParser(&updateBookDTO); err != nil {
			return dto.SendError(c, fiber.StatusBadRequest, "Invalid request body")
		}

		// Validate DTO
		if errors := dto.ValidateStruct(&updateBookDTO); errors != nil {
			return dto.SendValidationError(c, errors)
		}

		// Convert DTO to entity
		bookEntity := &entities.Book{
			Title:  updateBookDTO.Title,
			Author: updateBookDTO.Author,
		}

		result, err := service.UpdateBook(bookEntity)
		if err != nil {
			return dto.SendError(c, fiber.StatusInternalServerError, "Failed to update book")
		}

		return dto.SendSuccess(c, "Updated book successfully", result)
	}
}

// RemoveBook is handler/controller which removes Books from the BookShop
func RemoveBook(service book.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		idParam := c.Params("id")

		if idParam == "" {
			return dto.SendError(c, fiber.StatusBadRequest, "ID is required")
		}

		if err := service.RemoveBook(idParam); err != nil {
			return dto.SendError(c, fiber.StatusInternalServerError, "Failed to delete book")
		}

		return dto.SendSuccess(c, "Deleted book successfully", nil)
	}
}

// GetBooksWithPagination handles pagination via query parameters
func GetBooksWithPagination(service book.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		page := c.QueryInt("page", 1)
		limit := c.QueryInt("limit", 10)

		// Validate pagination parameters
		paginationDTO := dto.PaginationDTO{
			Page:  page,
			Limit: limit,
		}

		if errors := dto.ValidateStruct(&paginationDTO); errors != nil {
			return dto.SendValidationError(c, errors)
		}

		result, total, totalPages, err := service.GetBooksWithPagination(int64(page), int64(limit))
		if err != nil {
			return dto.SendError(c, fiber.StatusInternalServerError, "Failed to get books")
		}

		return c.JSON(fiber.Map{
			"error":      false,
			"message":    "Retrieved books successfully",
			"data":       result,
			"total":      total,
			"totalPages": totalPages,
		})
	}
}

// GetBooksWithPaginationBody handles pagination via request body
func GetBooksWithPaginationBody(service book.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var paginationDTO dto.PaginationDTO

		// Parse request body
		if err := c.BodyParser(&paginationDTO); err != nil {
			return dto.SendError(c, fiber.StatusBadRequest, "Invalid request body")
		}

		// Validate DTO
		if errors := dto.ValidateStruct(&paginationDTO); errors != nil {
			return dto.SendValidationError(c, errors)
		}

		result, total, totalPages, err := service.GetBooksWithPagination(int64(paginationDTO.Page), int64(paginationDTO.Limit))
		if err != nil {
			return dto.SendError(c, fiber.StatusInternalServerError, "Failed to get books")
		}

		return c.JSON(fiber.Map{
			"error":      false,
			"message":    "Retrieved books successfully",
			"data":       result,
			"total":      total,
			"totalPages": totalPages,
		})
	}
}
