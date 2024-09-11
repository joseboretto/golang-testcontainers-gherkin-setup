package book

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/joseboretto/golang-crud-api/internal/domain/models"
	"log"
	"net/http"
)

type SendEmailClient struct {
	httpClient    *http.Client
	host          string
	checkIsbnPath string
}

func NewSendEmailClient(host string, httpClient *http.Client) *SendEmailClient {
	return &SendEmailClient{
		httpClient:    httpClient,
		host:          host,
		checkIsbnPath: "/send-email",
	}
}

type SendEmailBookRequestBody struct {
	// Isbn: International Standard Book Number
	Isbn       string `json:"isbn"`
	Title      string `json:"title"`
	TotalPages int    `json:"total_pages"`
	Views      int    `json:"views"`
}

type SendEmailRequestBody struct {
	Email string                   `json:"email"`
	Book  SendEmailBookRequestBody `json:"book"`
}

func (c *SendEmailClient) SendEmail(email string, book *models.Book) error {
	// Create a new request
	sendEmailRequestBody := SendEmailRequestBody{
		Email: email,
		Book: SendEmailBookRequestBody{
			Isbn:       book.Isbn,
			Title:      book.Title,
			TotalPages: book.TotalPages,
			Views:      book.Views,
		},
	}
	sendEmailRequestBodyJson, err := json.Marshal(sendEmailRequestBody)
	if err != nil {
		log.Fatalf("Error marshalling sendEmailRequestBody on SendEmailClient.SendEmail: %s", err)
	}
	req, err := http.NewRequest(http.MethodPost, c.host+c.checkIsbnPath, bytes.NewReader(sendEmailRequestBodyJson))
	if err != nil {
		return err
	}
	// Send request
	res, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	// Check response status
	if res.StatusCode == http.StatusOK {
		return nil
	}
	return errors.New("Error sending email")
}
