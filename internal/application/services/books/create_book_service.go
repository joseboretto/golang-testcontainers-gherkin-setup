package books

import (
	"errors"

	"github.com/joseboretto/golang-testcontainers-gherkin-setup/internal/domain/models"
)

type CreateBookService struct {
	repository               CreateBookRepositoryInterface
	checkIsbnClientInterface CheckIsbnClientInterface
	sendEmailClientInterface SendEmailClientInterface
}

func NewCreateBookService(repository CreateBookRepositoryInterface,
	checkIsbnClientInterface CheckIsbnClientInterface,
	sendEmailClientInterface SendEmailClientInterface) *CreateBookService {
	return &CreateBookService{
		repository:               repository,
		checkIsbnClientInterface: checkIsbnClientInterface,
		sendEmailClientInterface: sendEmailClientInterface,
	}
}

func (s *CreateBookService) CreateBook(book *models.Book) (*models.Book, error) {
	// Check if ISBN is valid
	isValid, err := s.checkIsbnClientInterface.CheckIsbn(book.Isbn)
	if err != nil {
		return nil, err
	}
	if !isValid {
		return nil, errors.New("isbn is not valid based on external service")
	}
	// Check if book already exist
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
	// Send email
	err = s.sendEmailClientInterface.SendEmail("helloworld@gmail.com", storedBook)
	if err != nil {
		return nil, err
	}
	return storedBook, nil
}
