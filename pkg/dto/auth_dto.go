package dto

// RegisterDTO represents the data transfer object for user registration
type RegisterDTO struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6,max=100"`
	Name     string `json:"name" validate:"required,min=2,max=100"`
}

// LoginDTO represents the data transfer object for user login
type LoginDTO struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// AuthResponseDTO represents the response after successful authentication
type AuthResponseDTO struct {
	Token string      `json:"token"`
	User  UserInfoDTO `json:"user"`
}

// UserInfoDTO represents user information in responses
type UserInfoDTO struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}


