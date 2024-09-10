package books

import "github.com/joseboretto/golang-crud-api/internal/domain/models"

type CreateBookServiceInterface interface {
	CreateBook(book *models.Book) (*models.Book, error)
}

type CreateBookRepositoryInterface interface {
	SelectBookByIsbn(isbn string) (*models.Book, error)
	InsertBook(book *models.Book) (*models.Book, error)
}
