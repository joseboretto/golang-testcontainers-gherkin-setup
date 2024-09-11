package book

import (
	"errors"
	"github.com/joseboretto/golang-crud-api/internal/domain/models"
	"gorm.io/gorm"
)

type CreateBookRepository struct {
	database *gorm.DB
}

func NewCreateBookRepository(database *gorm.DB) *CreateBookRepository {
	return &CreateBookRepository{
		database: database,
	}
}

func (c *CreateBookRepository) InsertBook(book *models.Book) (*models.Book, error) {
	// Map model to entity
	bookEntity := BookEntity{
		Isbn:  book.Isbn,
		Title: book.Title,
	}
	// Insert entity
	result := c.database.Create(&bookEntity)
	if result.Error != nil {
		return nil, result.Error
	}
	// Map entity to model
	insert := &models.Book{
		Isbn:  bookEntity.Isbn,
		Title: bookEntity.Title,
	}
	return insert, nil
}

func (c *CreateBookRepository) SelectBookByIsbn(isbn string) (*models.Book, error) {
	bookEntity := BookEntity{}
	result := c.database.Where("isbn = ?", isbn).First(&bookEntity)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	// Map entity to model
	book := &models.Book{
		Isbn:  bookEntity.Isbn,
		Title: bookEntity.Title,
	}
	return book, nil
}
