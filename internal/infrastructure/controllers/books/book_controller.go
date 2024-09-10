package books

import (
	"net/http"
	"strings"

	servicebook "github.com/joseboretto/golang-crud-api/internal/application/services/books"
	"github.com/joseboretto/golang-crud-api/internal/infrastructure/controllers/books/dto"
	"github.com/joseboretto/golang-crud-api/internal/infrastructure/controllers/utils"
)

type Controller struct {
	createBookServiceInterface  servicebook.CreateBookServiceInterface
	getAllBooksServiceInterface servicebook.GetAllBooksServiceInterface
	getBookServiceInterface     servicebook.GetBookServiceInterface
}

func NewBookController(createBookServiceInterface servicebook.CreateBookServiceInterface, getAllBooksServiceInterface servicebook.GetAllBooksServiceInterface, getBookServiceInterface *servicebook.GetBookService) *Controller {
	return &Controller{
		createBookServiceInterface:  createBookServiceInterface,
		getAllBooksServiceInterface: getAllBooksServiceInterface,
		getBookServiceInterface:     getBookServiceInterface,
	}
}

func (c *Controller) CreateBook(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if req.Method == "POST" {
		createBookRequest := new(dto.CreateBookRequest)
		err := utils.Decode(req, &createBookRequest)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// mapper
		bookDomain := dto.MapToBookModel(createBookRequest)
		// service
		createBook, err := c.createBookServiceInterface.CreateBook(bookDomain)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		// response
		createBookResponse := dto.MapToCreateBookResponse(createBook)

		if err = utils.Response(w, createBookResponse, http.StatusOK); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad Request"))
	}

}

func (c *Controller) GetBooks(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if req.Method == "GET" {
		// service
		books, err := c.getAllBooksServiceInterface.GetAllBooks()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		// response
		booksResponse := make([]*dto.GetBookResponse, 0, len(books))
		for _, value := range books {
			bookResponse := dto.MapToGetAllBookResponse(value)
			booksResponse = append(booksResponse, bookResponse)
		}

		if err = utils.Response(w, booksResponse, http.StatusOK); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad Request"))
	}

}

func (c *Controller) GetBookByIsbn(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if req.Method == "GET" {
		//
		isbn := strings.TrimPrefix(req.URL.Path, "/api/v1/getBookByIsbn/")
		if isbn == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("isbn required"))
			return
		}
		// service
		book, err := c.getBookServiceInterface.GetBook(isbn)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		// response
		bookResponse := dto.MapToGetAllBookResponse(book)

		if err = utils.Response(w, bookResponse, http.StatusOK); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad Request"))
	}

}
