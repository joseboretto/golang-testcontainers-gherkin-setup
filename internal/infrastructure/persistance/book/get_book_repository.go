package book

import (
	"github.com/joseboretto/golang-crud-api/internal/domain/models"
	"github.com/joseboretto/golang-crud-api/internal/infrastructure/persistance"
)

type GetBookRepository struct {
	database persistance.InMemoryKeyValueStorage
}

func NewGetBookRepository(database persistance.InMemoryKeyValueStorage) *GetBookRepository {
	return &GetBookRepository{
		database: database,
	}
}

func (c *GetBookRepository) SelectBookByIsbn(isbn string) (*models.Book, error) {
	books, _ := c.database.SelectBookByIsbn(isbn)
	return books, nil
}

func (c *GetBookRepository) IncreaseBookViewsByIsbn(isbn string) {
	c.database.IncreaseBookViewsByIsbn(isbn)
}
