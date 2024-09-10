package book

import (
	"github.com/joseboretto/golang-crud-api/internal/domain/models"
	"github.com/joseboretto/golang-crud-api/internal/infrastructure/persistance"
)

type CreateBookRepository struct {
	database persistance.InMemoryKeyValueStorage
}

func NewCreateBookRepository(database persistance.InMemoryKeyValueStorage) *CreateBookRepository {
	return &CreateBookRepository{
		database: database,
	}
}

func (c *CreateBookRepository) InsertBook(book *models.Book) (*models.Book, error) {
	insert, err := c.database.Insert(book)
	if err != nil {
		return nil, err
	}
	return insert, nil
}

func (c *CreateBookRepository) SelectBookByIsbn(isbn string) (*models.Book, error) {
	books, _ := c.database.SelectBookByIsbn(isbn)
	return books, nil
}
