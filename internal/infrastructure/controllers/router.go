package controllers

import (
	"net/http"

	"github.com/joseboretto/golang-testcontainers-gherkin-setup/internal/infrastructure/controllers/books"
)

func SetupRoutes(bookController *books.Controller) {
	http.HandleFunc("/api/v1/createBook", bookController.CreateBook)
}
