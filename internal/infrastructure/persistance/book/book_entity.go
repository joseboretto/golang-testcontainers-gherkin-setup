package book

import (
	"gorm.io/gorm"
)

type BookEntity struct {
	gorm.Model
	// Isbn: International Standard Book Number
	Isbn       string `gorm:"column:isbn;unique"`
	Title      string `gorm:"column:title"`
	TotalPages int    `gorm:"column:total_pages"`
	Views      int    `gorm:"column:views"`
}

func (BookEntity) TableName() string {
	return "myschema.books"
}
