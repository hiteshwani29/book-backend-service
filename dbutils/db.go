package dbutils

import (
	"context"
	"errors"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Book struct {
	ID     string `json:"id" bson:"_id,omitempty"`
	Title  string `json:"title" bson:"title"`
	Author string `json:"author" bson:"author"`
}

type DbService interface {
	CreateBook(book Book) (string, error)
	GetBooks(book Book) ([]Book, error)
}

type dbCol struct {
	collection *mongo.Collection
}

func GetDbService(collection *mongo.Collection) DbService {
	return &dbCol{collection: collection}
}

func (db dbCol) CreateBook(book Book) (string, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := db.collection.InsertOne(ctx, book)
	if err != nil {
		return "", err
	}

	book.ID = result.InsertedID.(primitive.ObjectID).Hex()
	return book.ID, err
}

func (db dbCol) GetBooks(book Book) ([]Book, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var books []Book
	cursor, err := db.collection.Find(ctx, bson.D{{}})
	if err != nil {
		log.Panicf("Error:: %v", err.Error())
		return []Book{}, errors.New("failed to get books")
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &books); err != nil {
		return []Book{}, errors.New("failed to decode book")
	}
	log.Printf("books::%v", books)
	if len(books) == 0 {
		return []Book{}, nil
	}
	
	return books, err
}
