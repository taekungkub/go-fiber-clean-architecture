package handlers

import (
	"go-fiber-api/pkg/domain/user"
	"go-fiber-api/pkg/dto"

	"github.com/gofiber/fiber/v2"
)

func GetUsers(service user.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		result, err := service.FindAll()

		if err != nil {
			return dto.SendError(c, fiber.StatusInternalServerError, "Failed to get users")
		}

		return dto.SendSuccess(c, "Get users successfully", result)
	}
}

func CreateUser(service user.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var registerDTO dto.RegisterDTO

		// Parse request body
		if err := c.BodyParser(&registerDTO); err != nil {
			return dto.SendError(c, fiber.StatusBadRequest, "Invalid request body")
		}

		// Validate DTO
		if errors := dto.ValidateStruct(&registerDTO); errors != nil {
			return dto.SendValidationError(c, errors)
		}

		result, err := service.Register(&registerDTO)
		if err != nil {
			return dto.SendError(c, fiber.StatusBadRequest, err.Error())
		}

		return dto.SendSuccess(c, "Registered successfully", result)
	}
}

func UpdateUser(service user.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var updateUserDTO dto.UpdateUserDTO

		// Parse request body
		if err := c.BodyParser(&updateUserDTO); err != nil {
			return dto.SendError(c, fiber.StatusBadRequest, "Invalid request body")
		}

		// Validate DTO
		if errors := dto.ValidateStruct(&updateUserDTO); errors != nil {
			return dto.SendValidationError(c, errors)
		}

		result, err := service.UpdateUser(&updateUserDTO)

		if err != nil {
			return dto.SendError(c, fiber.StatusInternalServerError, err.Error())
		}

		return dto.SendSuccess(c, "Updated user successfully", result)
	}

}

func RemoveUser(service user.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		idParam := c.Params("id")

		if err := service.RemoveUser(idParam); err != nil {
			return dto.SendError(c, fiber.StatusInternalServerError, err.Error())
		}

		return dto.SendSuccess(c, "Deleted user successfully", nil)
	}
}
