package dto

import "github.com/joseboretto/golang-crud-api/internal/domain/models"

func MapToGetAllBookResponse(m *models.Book) *GetBookResponse {
	return &GetBookResponse{
		Title:      m.Title,
		TotalPages: m.TotalPages,
		Isbn:       m.Isbn,
		Views:      m.Views,
	}
}
