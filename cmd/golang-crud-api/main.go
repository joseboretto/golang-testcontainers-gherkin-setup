package main

import (
	"fmt"
	servicebook "github.com/joseboretto/golang-crud-api/internal/application/services/books"
	controller "github.com/joseboretto/golang-crud-api/internal/infrastructure/controllers"
	controllerbook "github.com/joseboretto/golang-crud-api/internal/infrastructure/controllers/books"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"os"
	"time"

	"log"
	"net/http"

	clientsbook "github.com/joseboretto/golang-crud-api/internal/infrastructure/clients/book"
	persistancebook "github.com/joseboretto/golang-crud-api/internal/infrastructure/persistance/book"
)

func main() {
	addr := ":8000"
	httpClient := &http.Client{
		Timeout: 5 * time.Second,
	}
	server, deferFn := mainHttpServerSetup(addr, httpClient)
	log.Println("Listing for requests at http://localhost" + addr)
	err := server.ListenAndServe()
	if err != nil {
		panic("Error staring the server: " + err.Error())
	}
	defer deferFn()
}

func mainHttpServerSetup(addr string, httpClient *http.Client) (*http.Server, func()) {
	db := getDatabaseConnection()
	// Migrate the schema
	err := db.AutoMigrate(&persistancebook.BookEntity{})
	if err != nil {
		panic("Error executing db.AutoMigrate" + err.Error())
	}
	// Clients
	checkIsbnClientHost := os.Getenv("CHECK_ISBN_CLIENT_HOST")
	checkIsbnClient := clientsbook.NewCheckIsbnClient(checkIsbnClientHost, httpClient)
	// repositories
	newCreateBookRepository := persistancebook.NewCreateBookRepository(db)
	// services
	createBookService := servicebook.NewCreateBookService(newCreateBookRepository, checkIsbnClient)
	// controllers
	bookController := controllerbook.NewBookController(createBookService)
	// routes
	controller.SetupRoutes(bookController)
	// Server
	server := &http.Server{Addr: addr, Handler: nil}
	// Defer function
	// Add all defer
	deferFn := func() {
		fmt.Println("deferFn executed")
		// TODO: FIX IT
		/*
			 fmt.Println("closing database")
				sqlDB, err := db.DB()
				err = sqlDB.Close()
				if err != nil {
					panic("Error closing database connection: " + err.Error())
				}
		*/

	}
	return server, deferFn
}

func getDatabaseConnection() *gorm.DB {
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
	}), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "myschema.", // schema name
			SingularTable: false,
		}})
	// Check if connection is successful
	if err != nil {
		panic("failed to connect database with error: " + err.Error() + "\n" + "Please check your database configuration: " + databaseConnectionString)
	}
	return db
}
