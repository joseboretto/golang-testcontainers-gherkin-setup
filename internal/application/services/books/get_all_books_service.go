package books

import "github.com/joseboretto/golang-crud-api/internal/domain/models"

type GetAllBooksService struct {
	repository GetAllBooksRepositoryInterface
}

func NewGetAllBooksService(repository GetAllBooksRepositoryInterface) *GetAllBooksService {
	return &GetAllBooksService{
		repository: repository,
	}
}

func (s *GetAllBooksService) GetAllBooks() ([]*models.Book, error) {
	books, err := s.repository.SelectBookByIsbn()
	if err != nil {
		return nil, err
	}
	return books, nil
}
