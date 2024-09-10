package main

import (
	"fmt"
	servicebook "github.com/joseboretto/golang-crud-api/internal/application/services/books"
	controller "github.com/joseboretto/golang-crud-api/internal/infrastructure/controllers"
	controllerbook "github.com/joseboretto/golang-crud-api/internal/infrastructure/controllers/books"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"

	"log"
	"net/http"

	clientsbook "github.com/joseboretto/golang-crud-api/internal/infrastructure/clients/book"
	persistancebook "github.com/joseboretto/golang-crud-api/internal/infrastructure/persistance/book"
)

func main() {
	databaseConnectionString, db, err := getDatabaseConnection()
	if err != nil {
		panic("failed to connect database with error: " + err.Error() + "\n" + "Please check your database configuration: " + databaseConnectionString)
	}
	// Migrate the schema
	err = db.AutoMigrate(&persistancebook.BookEntity{})
	if err != nil {
		panic("Error executing db.AutoMigrate" + err.Error())
	}
	// Clients
	checkIsbnClientHost := os.Getenv("CHECK_ISBN_CLIENT_HOST")
	checkIsbnClient := clientsbook.NewCheckIsbnClient(checkIsbnClientHost)
	// repositories
	newCreateBookRepository := persistancebook.NewCreateBookRepository(db)
	// services
	createBookService := servicebook.NewCreateBookService(newCreateBookRepository, checkIsbnClient)
	// controllers
	bookController := controllerbook.NewBookController(createBookService)
	// routes
	controller.SetupRoutes(bookController)

	log.Println("Listing for requests at http://localhost:8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func getDatabaseConnection() (string, *gorm.DB, error) {
	// Read database configuration from environment variables
	databaseUser := os.Getenv("DATABASE_USER")
	databasePassword := os.Getenv("DATABASE_PASSWORD")
	databaseName := os.Getenv("DATABASE_NAME")
	databaseHost := os.Getenv("DATABASE_HOST")
	databasePort := os.Getenv("DATABASE_PORT")
	databaseConnectionString := `host=` + databaseHost + ` user=` + databaseUser + ` password=` + databasePassword + ` dbname=` + databaseName + ` port=` + databasePort + ` sslmode=disable`
	fmt.Println("databaseConnectionString: ", databaseConnectionString)
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  databaseConnectionString,
		PreferSimpleProtocol: true, // disables implicit prepared statement usage
	}), &gorm.Config{})
	return databaseConnectionString, db, err
}
