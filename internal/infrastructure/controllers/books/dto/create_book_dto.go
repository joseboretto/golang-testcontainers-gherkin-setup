package dto

type CreateBookRequest struct {
	Title      string `json:"title"`
	TotalPages int    `json:"total_pages"`
	Isbn       string `json:"isbn"`
}

type CreateBookResponse struct {
	Title      string `json:"title"`
	TotalPages int    `json:"total_pages"`
	Isbn       string `json:"isbn"`
}
