package books

import "github.com/joseboretto/golang-crud-api/internal/domain/models"

type GetBookService struct {
	repository GetBookRepositoryInterface
}

func NewGetBookService(repository GetBookRepositoryInterface) *GetBookService {
	return &GetBookService{
		repository: repository,
	}
}

func (s *GetBookService) GetBook(isbn string) (*models.Book, error) {
	books, err := s.repository.SelectBookByIsbn(isbn)
	if err != nil {
		return nil, err
	}
	s.repository.IncreaseBookViewsByIsbn(isbn)
	return books, nil
}
