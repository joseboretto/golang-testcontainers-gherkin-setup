package book

import (
	"github.com/joseboretto/golang-crud-api/internal/domain/models"
	"github.com/joseboretto/golang-crud-api/internal/infrastructure/persistance"
)

type GetAllBooksRepository struct {
	database persistance.InMemoryKeyValueStorage
}

func NewGetAllBooksRepository(database persistance.InMemoryKeyValueStorage) *GetAllBooksRepository {
	return &GetAllBooksRepository{
		database: database,
	}
}

func (c *GetAllBooksRepository) SelectBookByIsbn() ([]*models.Book, error) {
	return c.database.GetAll(), nil
}
