package dto

type GetBookResponse struct {
	Title      string `json:"title"`
	TotalPages int    `json:"total_pages"`
	Isbn       string `json:"isbn"`
	Views      int    `json:"views"`
}
