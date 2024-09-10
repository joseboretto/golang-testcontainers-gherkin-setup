package book

import (
	"gorm.io/gorm"
)

type BookEntity struct {
	gorm.Model
	// Isbn: International Standard Book Number
	Isbn       string
	Title      string
	TotalPages int
	Views      int
}
