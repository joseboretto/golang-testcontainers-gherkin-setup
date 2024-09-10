package book

import (
	"github.com/joseboretto/golang-crud-api/internal/domain/models"
	"gorm.io/gorm"
)

type GetAllBooksRepository struct {
	database *gorm.DB
}

func NewGetAllBooksRepository(database *gorm.DB) *GetAllBooksRepository {
	return &GetAllBooksRepository{
		database: database,
	}
}

func (c *GetAllBooksRepository) SelectAllBook() (books []*models.Book, err error) {
	//
	var booksEntity []BookEntity
	tx := c.database.Find(&booksEntity)
	if tx.Error != nil {
		return nil, tx.Error
	}
	for _, bookEntity := range booksEntity {
		book := &models.Book{
			Isbn:       bookEntity.Isbn,
			Title:      bookEntity.Title,
			TotalPages: bookEntity.TotalPages,
			Views:      bookEntity.Views,
		}
		books = append(books, book)
	}
	return books, nil
}
