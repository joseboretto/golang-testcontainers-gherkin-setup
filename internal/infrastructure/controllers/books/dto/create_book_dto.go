package dto

type CreateBookRequest struct {
	Title string `json:"title"`
	Isbn  string `json:"isbn"`
}

type CreateBookResponse struct {
	Title string `json:"title"`
	Isbn  string `json:"isbn"`
}
