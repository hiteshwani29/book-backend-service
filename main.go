package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/BookApp/book-backend-service/dbutils"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DbService dbutils.DbService

func main() {
	// Connect to MongoDB
	host := os.Getenv("MONGO_HOST")
	username := os.Getenv("MONGO_USERNAME")
	password := os.Getenv("MONGO_PASSWORD")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	mongourl := fmt.Sprintf("mongodb://%v:%v@%v:27017", username, password, host)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongourl))
	if err != nil {
		log.Fatal(err)
	}

	// Select the database and collection
	db := client.Database("bookApp")
	DbService = dbutils.GetDbService(db.Collection("books"))

	// Create a new Gin router
	router := gin.Default()

	// CORS middleware configuration
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"https://app.example.com"} // Replace with your frontend domain
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}

	// Use the CORS middleware
	router.Use(cors.New(config))

	// Define API endpoints
	router.GET("/api/books", getBooks)
	router.GET("/api/books/:id", getBook)
	router.POST("/api/books", createBook)
	router.GET("/health", health)
	router.GET("/", mm)

	// Start the server
	router.Run(":8000")
}

func mm(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Service runnning"})
}

func health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Service healthy"})
}

func getBooks(c *gin.Context) {
	log.Printf("\nIn Get multi Books \n")
	books, err := DbService.GetBooks(dbutils.Book{})
	if err != nil {
		log.Panicf("Error:: %v", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Books not found"})
		return
	}
	log.Printf("books::%v", books)
	log.Printf("len(books):: %v", len(books))

	c.JSON(http.StatusOK, books)
}

func getBook(c *gin.Context) {
	id := c.Param("id")
	log.Printf("\nIn Get Book :: %v\n", id)
	books, err := DbService.GetBooks(dbutils.Book{ID: id})
	if err != nil {
		log.Panicf("Error:: %v", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Book not found"})
	}

	c.JSON(http.StatusOK, books)
}

func createBook(c *gin.Context) {
	log.Printf("\nIn Created Book \n")
	var book dbutils.Book
	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}
	id, err := DbService.CreateBook(book)
	if err != nil {
		log.Panicf("Error:: %v", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create book"})
		return
	}

	book.ID = id
	c.JSON(http.StatusCreated, book)
}
