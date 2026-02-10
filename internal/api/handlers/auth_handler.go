package handlers

import (
	"go-fiber-api/pkg/domain/user"
	"go-fiber-api/pkg/dto"

	"github.com/gofiber/fiber/v2"
)

// Register handles user registration
func Register(service user.Service) fiber.Handler {
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

		// Register user
		authResponse, err := service.Register(&registerDTO)
		if err != nil {
			return dto.SendError(c, fiber.StatusBadRequest, err.Error())
		}

		return c.Status(fiber.StatusCreated).JSON(dto.SuccessResponse{
			Error:   false,
			Message: "User registered successfully",
			Data:    authResponse,
		})
	}
}

// Login handles user login
func Login(service user.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var loginDTO dto.LoginDTO

		// Parse request body
		if err := c.BodyParser(&loginDTO); err != nil {
			return dto.SendError(c, fiber.StatusBadRequest, "Invalid request body")
		}

		// Validate DTO
		if errors := dto.ValidateStruct(&loginDTO); errors != nil {
			return dto.SendValidationError(c, errors)
		}

		// Login user
		authResponse, err := service.Login(&loginDTO)
		if err != nil {
			return dto.SendError(c, fiber.StatusUnauthorized, err.Error())
		}

		return dto.SendSuccess(c, "Login successful", authResponse)
	}
}

// GetProfile handles getting the current user's profile
func GetProfile(service user.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get user ID from context (set by JWT middleware)
		userID := c.Locals("userID").(string)

		// Get user
		user, err := service.GetUserByID(userID)
		if err != nil {
			return dto.SendError(c, fiber.StatusNotFound, "User not found")
		}

		return dto.SendSuccess(c, "Profile retrieved successfully", dto.UserInfoDTO{
			ID:    user.ID.Hex(),
			Email: user.Email,
			Name:  user.Name,
		})
	}
}
