package dto

import "github.com/joseboretto/golang-crud-api/internal/domain/models"

func MapToBookModel(c *CreateBookRequest) *models.Book {
	return &models.Book{
		Title: c.Title,
		Isbn:  c.Isbn,
	}
}

func MapToCreateBookResponse(m *models.Book) *CreateBookResponse {
	return &CreateBookResponse{
		Title: m.Title,
		Isbn:  m.Isbn,
	}
}
