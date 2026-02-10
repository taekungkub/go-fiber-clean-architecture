package book

import (
	"errors"
	"go-fiber-api/pkg/entities"
	"math"
)

// Service is an interface from which our api module can access our repository of all our models
type Service interface {
	InsertBook(book *entities.Book) (*entities.Book, error)
	GetBooks() (*[]entities.Book, error)
	UpdateBook(book *entities.Book) (*entities.Book, error)
	RemoveBook(ID string) error
	GetBooksWithPagination(page int64, limit int64) (*[]entities.Book, int64, int64, error)
}

type service struct {
	repository Repository
}

// NewService is used to create a single instance of the service
func NewService(r Repository) Service {
	return &service{
		repository: r,
	}
}

// InsertBook is a service layer that helps insert book in BookShop
func (s *service) InsertBook(book *entities.Book) (*entities.Book, error) {

	// 1. Check author
	author, err := s.repository.GetBookByAuthor(book.Author)

	if err != nil {
		return nil, err
	}

	if author != nil {
		return nil, errors.New("author already exists")
	}

	return s.repository.CreateBook(book)
}

// FetchBooks is a service layer that helps fetch all books in BookShop
func (s *service) GetBooks() (*[]entities.Book, error) {
	return s.repository.GetBooks()
}

// UpdateBook is a service layer that helps update books in BookShop
// func (s *service) UpdateBook(book *entities.Book) (*entities.Book, error) {
// 	return s.repository.UpdateBook(book)
// }

func (s *service) UpdateBook(book *entities.Book) (*entities.Book, error) {

	existing, err := s.repository.GetBookByID(book.ID.Hex())
	if err != nil {
		return nil, err
	}

	if existing == nil {
		return nil, errors.New("book not found")
	}

	return s.repository.UpdateBook(book)
}

// RemoveBook is a service layer that helps remove books from BookShop
func (s *service) RemoveBook(ID string) error {
	return s.repository.DeleteBook(ID)
}

func (s *service) GetBooksWithPagination(page int64, limit int64) (*[]entities.Book, int64, int64, error) {

	books, total, err := s.repository.GetBooksWithPagination(page, limit)
	if err != nil {
		return nil, 0, 0, err
	}

	totalPages := int64(math.Ceil(float64(total) / float64(limit)))

	return books,
		total,
		totalPages,
		nil
}
