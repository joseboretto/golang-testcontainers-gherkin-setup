package controllers

import (
	"net/http"

	"github.com/joseboretto/golang-crud-api/internal/infrastructure/controllers/books"
)

func SetupRoutes(bookController *books.Controller) {
	http.HandleFunc("/api/v1/createBook", bookController.CreateBook)
	http.HandleFunc("/api/v1/getBooks", bookController.GetBooks)
	http.HandleFunc("/api/v1/getBookByIsbn/", bookController.GetBookByIsbn)
}
