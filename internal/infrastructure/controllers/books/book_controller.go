package books

import (
	servicebook "github.com/joseboretto/golang-crud-api/internal/application/services/books"
	"github.com/joseboretto/golang-crud-api/internal/infrastructure/controllers/books/dto"
	"github.com/joseboretto/golang-crud-api/internal/infrastructure/controllers/utils"
	"net/http"
)

type Controller struct {
	createBookServiceInterface servicebook.CreateBookServiceInterface
}

func NewBookController(createBookServiceInterface servicebook.CreateBookServiceInterface) *Controller {
	return &Controller{
		createBookServiceInterface: createBookServiceInterface,
	}
}

func (c *Controller) CreateBook(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if req.Method == "POST" {
		createBookRequest := new(dto.CreateBookRequest)
		err := utils.Decode(req, &createBookRequest)
		if err != nil {
			errorResponse := ErrorResponse{
				Message: "Bad Request. Invalid payload",
				Error:   err.Error(),
			}
			utils.Response(w, errorResponse, http.StatusBadRequest)
			return

		}

		// mapper
		bookDomain := dto.MapToBookModel(createBookRequest)
		// service
		createBook, err := c.createBookServiceInterface.CreateBook(bookDomain)
		if err != nil {
			errorResponse := ErrorResponse{
				Message: "Bad Request. Error creating book",
				Error:   err.Error(),
			}
			utils.Response(w, errorResponse, http.StatusBadRequest)
			return
		}
		// response
		createBookResponse := dto.MapToCreateBookResponse(createBook)

		if err = utils.Response(w, createBookResponse, http.StatusOK); err != nil {
			errorResponse := ErrorResponse{
				Message: "Internal Server Error",
				Error:   err.Error(),
			}
			utils.Response(w, errorResponse, http.StatusBadRequest)

			return
		}
	} else {
		errorResponse := ErrorResponse{
			Message: "Bad Request. Method not allowed. POST required",
			Error:   "",
		}
		utils.Response(w, errorResponse, http.StatusBadRequest)
	}

}
