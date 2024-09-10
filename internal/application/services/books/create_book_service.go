package books

import (
	"errors"

	"github.com/joseboretto/golang-crud-api/internal/domain/models"
)

type CreateBookService struct {
	repository CreateBookRepositoryInterface
}

func NewCreateBookService(repository CreateBookRepositoryInterface) *CreateBookService {
	return &CreateBookService{
		repository: repository,
	}
}

func (s *CreateBookService) CreateBook(book *models.Book) (*models.Book, error) {
	// TODO: Add business logic here
	exist, err := s.repository.SelectBookByIsbn(book.Isbn)
	if err != nil {
		return nil, err
	}
	if exist != nil {
		return nil, errors.New("book already exist")
	}
	storedBook, err := s.repository.InsertBook(book)
	if err != nil {
		return nil, err
	}
	return storedBook, nil
}
