package books

import "github.com/joseboretto/golang-testcontainers-gherkin-setup/internal/domain/models"

// CreateBookServiceInterface Inbound port
type CreateBookServiceInterface interface {
	CreateBook(book *models.Book) (*models.Book, error)
}

// CreateBookRepositoryInterface Outbound port
type CreateBookRepositoryInterface interface {
	SelectBookByIsbn(isbn string) (*models.Book, error)
	InsertBook(book *models.Book) (*models.Book, error)
}

// CheckIsbnClientInterface Outbound port
type CheckIsbnClientInterface interface {
	CheckIsbn(isbn string) (bool, error)
}

// SendEmailClient Outbound port
type SendEmailClientInterface interface {
	SendEmail(email string, book *models.Book) error
}
