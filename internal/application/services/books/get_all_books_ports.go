package books

import "github.com/joseboretto/golang-crud-api/internal/domain/models"

type GetAllBooksServiceInterface interface {
	GetAllBooks() ([]*models.Book, error)
}

type GetAllBooksRepositoryInterface interface {
	SelectBookByIsbn() ([]*models.Book, error)
}
