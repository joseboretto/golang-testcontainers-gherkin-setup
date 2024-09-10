package main

import (
	servicebook "github.com/joseboretto/golang-crud-api/internal/application/services/books"
	controller "github.com/joseboretto/golang-crud-api/internal/infrastructure/controllers"
	controllerbook "github.com/joseboretto/golang-crud-api/internal/infrastructure/controllers/books"
	"github.com/joseboretto/golang-crud-api/internal/infrastructure/persistance"

	"log"
	"net/http"

	persistancebook "github.com/joseboretto/golang-crud-api/internal/infrastructure/persistance/book"
)

func main() {
	// Dependency Injection
	database := persistance.NewInMemoryKeyValueStorage()
	// repositories
	newCreateBookRepository := persistancebook.NewCreateBookRepository(*database)
	newGetAllBooksRepository := persistancebook.NewGetAllBooksRepository(*database)
	newGetBookRepository := persistancebook.NewGetBookRepository(*database)
	// services
	createBookService := servicebook.NewCreateBookService(newCreateBookRepository)
	getAllBooksService := servicebook.NewGetAllBooksService(newGetAllBooksRepository)
	newGetBookService := servicebook.NewGetBookService(newGetBookRepository)
	// controllers
	bookController := controllerbook.NewBookController(createBookService, getAllBooksService, newGetBookService)
	// routes
	controller.SetupRoutes(bookController)

	log.Println("Listing for requests at http://localhost:8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
