package persistance

import (
	"fmt"
	"sync"

	"github.com/joseboretto/golang-crud-api/internal/domain/models"
)

type InMemoryKeyValueStorage struct {
	// This uses pointers over values because it can be nil
	// This is not thread-safe, and it's not meant to be used in production.
	// booksMap is a map of books, using the isbn as the key.
	// In this case, the database is using the model, but it should be an Entity.
	// map[string]*models.Book
	booksMap *sync.Map
}

func NewInMemoryKeyValueStorage() *InMemoryKeyValueStorage {
	return &InMemoryKeyValueStorage{
		booksMap: &sync.Map{},
	}
}

func (i *InMemoryKeyValueStorage) Insert(book *models.Book) (*models.Book, error) {
	i.booksMap.Store(book.Isbn, book)
	return book, nil
}

func (i *InMemoryKeyValueStorage) GetAll() []*models.Book {
	books := make([]*models.Book, 0)
	i.booksMap.Range(func(key, value interface{}) bool {
		bookValue, ok := value.(*models.Book)
		if !ok {
			fmt.Println("Conversion failed")
			return true
		}
		books = append(books, bookValue)
		return true
	})

	return books
}

func (i *InMemoryKeyValueStorage) SelectBookByIsbn(isbn string) (*models.Book, bool) {
	value, ok := i.booksMap.Load(isbn)
	bookValue, ok := value.(*models.Book)
	return bookValue, ok
}

func (i *InMemoryKeyValueStorage) IncreaseBookViewsByIsbn(isbn string) {
	value, _ := i.booksMap.Load(isbn)
	bookValue, _ := value.(*models.Book)
	bookValue.Views++
	i.booksMap.Store(bookValue.Isbn, bookValue)
}
