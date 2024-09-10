package book

import (
	"fmt"
	"net/http"
	"time"
)

type CheckIsbnClient struct {
	httpClient    http.Client
	host          string
	checkIsbnPath string
}

func NewCheckIsbnClient(host string) *CheckIsbnClient {
	// Create a new http client
	httpClient := http.Client{
		Timeout: time.Second * 10,
	}
	return &CheckIsbnClient{
		httpClient:    httpClient,
		host:          host,
		checkIsbnPath: "/isbn/%s",
	}
}

func (c *CheckIsbnClient) CheckIsbn(isbn string) (bool, error) {
	// Create a new request
	url := c.host + fmt.Sprintf(c.checkIsbnPath, isbn)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return false, err
	}
	// Send request
	res, err := c.httpClient.Do(req)
	if err != nil {
		return false, err
	}
	// Check response status
	if res.StatusCode == http.StatusOK {
		return true, nil
	}
	return false, nil
}
