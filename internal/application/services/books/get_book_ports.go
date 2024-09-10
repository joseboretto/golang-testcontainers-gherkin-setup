package books

import "github.com/joseboretto/golang-crud-api/internal/domain/models"

type GetBookServiceInterface interface {
	GetBook(isbn string) (*models.Book, error)
}

type GetBookRepositoryInterface interface {
	SelectBookByIsbn(isbn string) (*models.Book, error)
	IncreaseBookViewsByIsbn(isbn string)
}
