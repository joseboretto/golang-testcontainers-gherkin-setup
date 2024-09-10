package book

import (
	"github.com/joseboretto/golang-crud-api/internal/domain/models"
	"gorm.io/gorm"
)

type GetBookRepository struct {
	database *gorm.DB
}

func NewGetBookRepository(database *gorm.DB) *GetBookRepository {
	return &GetBookRepository{
		database: database,
	}
}

func (c *GetBookRepository) SelectBookByIsbn(isbn string) (*models.Book, error) {
	//
	var bookEntity BookEntity
	tx := c.database.Find(&bookEntity).Where("isbn = ?", isbn)
	if tx.Error != nil {
		return nil, tx.Error
	}

	book := &models.Book{
		Isbn:       bookEntity.Isbn,
		Title:      bookEntity.Title,
		TotalPages: bookEntity.TotalPages,
		Views:      bookEntity.Views,
	}
	return book, nil
}

func (c *GetBookRepository) IncreaseBookViewsByIsbn(isbn string) {
	// TODO: Implement this method
}
