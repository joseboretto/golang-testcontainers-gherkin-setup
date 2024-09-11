package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/cucumber/godog"
	"io"
	"log"
	"net/http"
	"time"
)

func (s *StepsContext) RegisterApiSteps(sc *godog.ScenarioContext) {
	sc.Step(`^API "([^"]*)" request is sent to "([^"]*)" without payload`, s.apiRequestIsSendWithoutPayload)
	sc.Step(`^API "([^"]*)" request is sent to "([^"]*)" with payload`, s.apiRequestIsSendWithPayload)
	sc.Step(`^response status code is (\d+) and payload is$`, s.apiResponseIs)
}

func (s *StepsContext) apiRequestIsSendWithoutPayload(method, url string) error {
	return s.apiRequestIsSendWithPayload(method, url, "")
}

func (s *StepsContext) apiRequestIsSendWithPayload(method, url, payloadJson string) error {
	client := &http.Client{
		Timeout: time.Duration(5) * time.Second,
	}

	// Create a new HTTP request
	req, err := http.NewRequestWithContext(context.Background(),
		method,
		s.mainHttpServerUrl+url,
		bytes.NewBufferString(payloadJson))
	if err != nil {
		return fmt.Errorf("failed to create a new HTTP request: %w", err)
	}

	// Set content type to JSON. Modify this if your payload is not JSON
	req.Header.Set("Content-Type", "application/json")

	// Execute the HTTP request
	response, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute HTTP request: %w", err)
	}
	// store the response in the context
	s.stepResponse = response

	return nil
}

func (s *StepsContext) apiResponseIs(expected int, expectedResponse string) error {
	if s.stepResponse.StatusCode != expected {
		fmt.Println("Response:" + getBody(s.stepResponse))
		return fmt.Errorf("apiResponseIs: Expected %d but got %d", expected, s.stepResponse.StatusCode)
	}

	actualJson := getBody(s.stepResponse)

	if compare, err := compareJSON(expectedResponse, actualJson); err != nil {
		return err
	} else if !compare {
		log.Printf("Actual respponse body without escapes: %s", actualJson)
		return errors.New(fmt.Sprintf("Response body doesnt match. Expected body: %s, \n actual body: %v", expectedResponse, actualJson))
	}
	return nil
}

func getBody(response *http.Response) string {
	body, err1 := io.ReadAll(response.Body)
	if err1 != nil {
		panic(err1)
	}
	// close response body
	err2 := response.Body.Close()
	if err2 != nil {
		panic(err2)
	}

	return string(body)
}
