package user

import (
	"errors"
	"go-fiber-api/pkg/core"
	"go-fiber-api/pkg/dto"
	"go-fiber-api/pkg/entities"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

// Service interface defines user business logic operations
type Service interface {
	FindAll() ([]entities.User, error)
	Register(registerDTO *dto.RegisterDTO) (*dto.AuthResponseDTO, error)
	Login(loginDTO *dto.LoginDTO) (*dto.AuthResponseDTO, error)
	GetUserByID(id string) (*entities.User, error)
	UpdateUser(user *dto.UpdateUserDTO) (*entities.User, error)
	RemoveUser(id string) error
}

type service struct {
	repo Repository
}

// NewService creates a new user service
func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

// FindAll finds all users in the database
func (s *service) FindAll() ([]entities.User, error) {
	return s.repo.FindAll()
}

// Register registers a new user
func (s *service) Register(registerDTO *dto.RegisterDTO) (*dto.AuthResponseDTO, error) {
	// Check if user already exists
	existingUser, _ := s.repo.FindByEmail(registerDTO.Email)
	if existingUser != nil {
		return nil, errors.New("user with this email already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registerDTO.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Create user
	user := &entities.User{
		Email:    registerDTO.Email,
		Password: string(hashedPassword),
		Name:     registerDTO.Name,
	}

	createdUser, err := s.repo.CreateUser(user)
	if err != nil {
		return nil, err
	}

	// Generate JWT token
	token, err := core.GenerateToken(createdUser.ID.Hex(), createdUser.Email)
	if err != nil {
		return nil, err
	}

	// Return auth response
	return &dto.AuthResponseDTO{
		Token: token,
		User: dto.UserInfoDTO{
			ID:    createdUser.ID.Hex(),
			Email: createdUser.Email,
			Name:  createdUser.Name,
		},
	}, nil
}

// Login authenticates a user and returns a JWT token
func (s *service) Login(loginDTO *dto.LoginDTO) (*dto.AuthResponseDTO, error) {
	// Find user by email
	user, err := s.repo.FindByEmail(loginDTO.Email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	// Compare password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginDTO.Password))
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	// Generate JWT token
	token, err := core.GenerateToken(user.ID.Hex(), user.Email)
	if err != nil {
		return nil, err
	}

	// Return auth response
	return &dto.AuthResponseDTO{
		Token: token,
		User: dto.UserInfoDTO{
			ID:    user.ID.Hex(),
			Email: user.Email,
			Name:  user.Name,
		},
	}, nil
}

// GetUserByID retrieves a user by ID
func (s *service) GetUserByID(id string) (*entities.User, error) {
	return s.repo.FindByID(id)
}

func (s *service) UpdateUser(dto *dto.UpdateUserDTO) (*entities.User, error) {

	id, err := primitive.ObjectIDFromHex(dto.ID)
	if err != nil {
		return nil, err
	}

	existing, err := s.repo.FindByID(id.Hex())
	if err != nil {
		return nil, err
	}

	if existing == nil {
		return nil, errors.New("user not found")
	}

	// update user
	user := &entities.User{
		ID:        id,
		Email:     dto.Email,
		Name:      dto.Name,
		Password:  existing.Password,
		CreatedAt: existing.CreatedAt,
		UpdatedAt: time.Now(),
	}

	return s.repo.UpdateUser(user)
}

func (s *service) RemoveUser(id string) error {
	return s.repo.DeleteUser(id)
}
