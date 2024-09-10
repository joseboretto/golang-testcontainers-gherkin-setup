package main

import (
	servicebook "github.com/joseboretto/golang-crud-api/internal/application/services/books"
	controller "github.com/joseboretto/golang-crud-api/internal/infrastructure/controllers"
	controllerbook "github.com/joseboretto/golang-crud-api/internal/infrastructure/controllers/books"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"

	"log"
	"net/http"

	persistancebook "github.com/joseboretto/golang-crud-api/internal/infrastructure/persistance/book"
)

func main() {
	// Read database configuration from environment variables
	databaseUser := os.Getenv("DATABASE_USER")
	databasePassword := os.Getenv("DATABASE_PASSWORD")
	databaseName := os.Getenv("DATABASE_NAME")
	databaseHost := os.Getenv("DATABASE_HOST")
	databasePort := os.Getenv("DATABASE_PORT")
	databaseConnectionString := `host=` + databaseHost + ` user=` + databaseUser + ` password=` + databasePassword + ` dbname=` + databaseName + ` port=` + databasePort + ` sslmode=disable`
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  databaseConnectionString,
		PreferSimpleProtocol: true, // disables implicit prepared statement usage
	}), &gorm.Config{})
	if err != nil {
		panic("failed to connect database with error: " + err.Error() + "\n" + "Please check your database configuration: " + databaseConnectionString)
	}
	// Migrate the schema
	err = db.AutoMigrate(&persistancebook.BookEntity{})
	if err != nil {
		panic("Error executing db.AutoMigrate" + err.Error())
	}

	// repositories
	newCreateBookRepository := persistancebook.NewCreateBookRepository(db)
	newGetAllBooksRepository := persistancebook.NewGetAllBooksRepository(db)
	newGetBookRepository := persistancebook.NewGetBookRepository(db)
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
