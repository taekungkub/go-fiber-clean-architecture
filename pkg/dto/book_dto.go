package dto

// CreateBookDTO represents the data transfer object for creating a book
type CreateBookDTO struct {
	Title  string `json:"title" validate:"required,min=1,max=200"`
	Author string `json:"author" validate:"max=100"`
}

// UpdateBookDTO represents the data transfer object for updating a book
type UpdateBookDTO struct {
	ID     string `json:"id" validate:"required"`
	Title  string `json:"title" validate:"required,min=1,max=200"`
	Author string `json:"author" validate:"max=100"`
}

// PaginationDTO represents the data transfer object for pagination requests
type PaginationDTO struct {
	Page  int `json:"page" validate:"required,gte=1"`
	Limit int `json:"limit" validate:"required,gte=1,lte=100"`
}
