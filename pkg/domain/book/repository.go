package book

import (
	"context"
	"go-fiber-api/pkg/entities"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Repository interface allows us to access the CRUD Operations in mongo here.
type Repository interface {
	CreateBook(book *entities.Book) (*entities.Book, error)
	GetBooks() (*[]entities.Book, error)
	UpdateBook(book *entities.Book) (*entities.Book, error)
	DeleteBook(ID string) error
	GetBookByID(id string) (*entities.Book, error)
	GetBooksWithPagination(page int64, limit int64) (*[]entities.Book, int64, error)
	GetBookByAuthor(author string) (*entities.Book, error)
}
type repository struct {
	Collection *mongo.Collection
}

// NewRepo is the single instance repo that is being created.
func NewRepo(collection *mongo.Collection) Repository {
	return &repository{
		Collection: collection,
	}
}

// CreateBook is a mongo repository that helps to create books
func (r *repository) CreateBook(book *entities.Book) (*entities.Book, error) {
	book.ID = primitive.NewObjectID()
	book.CreatedAt = time.Now()
	book.UpdatedAt = time.Now()
	_, err := r.Collection.InsertOne(context.Background(), book)
	if err != nil {
		return nil, err
	}
	return book, nil
}

// GetBooks is a mongo repository that helps to fetch books
func (r *repository) GetBooks() (*[]entities.Book, error) {
	var books []entities.Book = make([]entities.Book, 0)

	cursor, err := r.Collection.Find(context.Background(), bson.D{})
	if err != nil {
		return nil, err
	}

	if err := cursor.All(context.Background(), &books); err != nil {
		return nil, err
	}
	return &books, nil
}

// UpdateBook is a mongo repository that helps to update books
func (r *repository) UpdateBook(book *entities.Book) (*entities.Book, error) {
	book.UpdatedAt = time.Now()
	_, err := r.Collection.UpdateOne(context.Background(), bson.M{"_id": book.ID}, bson.M{"$set": book})
	if err != nil {
		return nil, err
	}
	return book, nil
}

// DeleteBook is a mongo repository that helps to delete books
func (r *repository) DeleteBook(ID string) error {
	bookID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		return err
	}
	_, err = r.Collection.DeleteOne(context.Background(), bson.M{"_id": bookID})
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) GetBookByID(id string) (*entities.Book, error) {
	var book entities.Book
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	result := r.Collection.FindOne(context.Background(), bson.M{"_id": objID})
	if err := result.Decode(&book); err != nil {
		return nil, err
	}
	return &book, nil
}

func (r *repository) GetBooksWithPagination(page int64, limit int64) (*[]entities.Book, int64, error) {
	var books []entities.Book = make([]entities.Book, 0)

	opts := options.Find().
		SetLimit(limit).
		SetSkip(int64((page - 1) * limit)).
		SetSort(bson.D{{Key: "_id", Value: -1}}) // newest first

	cursor, err := r.Collection.Find(context.Background(), bson.D{}, opts)
	if err != nil {
		return nil, 0, err
	}

	if err := cursor.All(context.Background(), &books); err != nil {
		return nil, 0, err
	}

	// นับทั้งหมด (เอาไปคำนวณ totalPage)
	total, err := r.Collection.CountDocuments(context.Background(), bson.M{})
	if err != nil {
		return nil, 0, err
	}

	return &books, total, nil
}

func (r *repository) GetBookByAuthor(author string) (*entities.Book, error) {
	var book entities.Book

	err := r.Collection.
		FindOne(context.Background(), bson.M{"author": author}).
		Decode(&book)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // 👈 no doc = OK
		}
		return nil, err // real DB error
	}

	return &book, nil
}
