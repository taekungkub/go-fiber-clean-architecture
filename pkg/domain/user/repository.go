package user

import (
	"context"
	"errors"
	"go-fiber-api/pkg/entities"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository interface {
	FindAll() ([]entities.User, error)
	CreateUser(user *entities.User) (*entities.User, error)
	FindByEmail(email string) (*entities.User, error)
	FindByID(id string) (*entities.User, error)
	UpdateUser(user *entities.User) (*entities.User, error)
	DeleteUser(id string) error
}

type repository struct {
	Collection *mongo.Collection
}

func NewRepo(collection *mongo.Collection) Repository {
	return &repository{
		Collection: collection,
	}
}

func (r *repository) FindAll() ([]entities.User, error) {
	var users []entities.User
	cursor, err := r.Collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	if err = cursor.All(context.Background(), &users); err != nil {
		return nil, err
	}
	return users, nil
}

func (r *repository) CreateUser(user *entities.User) (*entities.User, error) {
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	result, err := r.Collection.InsertOne(context.Background(), user)
	if err != nil {
		return nil, err
	}

	user.ID = result.InsertedID.(primitive.ObjectID)
	return user, nil
}

func (r *repository) FindByEmail(email string) (*entities.User, error) {
	var user entities.User
	err := r.Collection.FindOne(context.Background(), bson.M{"email": email}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *repository) FindByID(id string) (*entities.User, error) {
	objectID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return nil, err
	}

	var user entities.User
	err = r.Collection.FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&user)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (r *repository) UpdateUser(user *entities.User) (*entities.User, error) {

	query := bson.M{"_id": user.ID}
	update := bson.M{"$set": user}

	_, err := r.Collection.UpdateOne(context.Background(), query, update)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *repository) DeleteUser(id string) error {

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid id format")
	}

	var user entities.User

	err = r.Collection.
		FindOneAndDelete(
			context.Background(),
			bson.M{"_id": objID},
		).
		Decode(&user)

	if err == mongo.ErrNoDocuments {
		return errors.New("user not found")
	}

	if err != nil {
		return err
	}

	return nil

}
