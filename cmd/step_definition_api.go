package main

import (
	"bytes"
	"context"
	"fmt"
	"github.com/cucumber/godog"
	"io"
	"log"
	"net/http"
	"time"
)

func (s *StepsContext) RegisterApiSteps(sc *godog.ScenarioContext) {
	sc.Step(`^API "([^"]*)" request is sent to "([^"]*)" without payload$`, s.apiRequestIsSendWithoutPayload)
	sc.Step(`^API "([^"]*)" request is sent to "([^"]*)" with payload$`, s.apiRequestIsSendWithPayload)
	sc.Step(`^response status code is (\d+) and payload is$`, s.apiResponseIs)
}

func (s *StepsContext) apiRequestIsSendWithoutPayload(method, url string) error {
	return s.apiRequestIsSendWithPayload(method, url, "")
}

func (s *StepsContext) apiRequestIsSendWithPayload(method, url, payloadJson string) error {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	req, err := http.NewRequestWithContext(context.Background(), method, s.mainHttpServerUrl+url, bytes.NewBufferString(payloadJson))
	if err != nil {
		return fmt.Errorf("failed to create a new HTTP request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	response, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute HTTP request: %w", err)
	}

	s.stepResponse = response
	return nil
}

func (s *StepsContext) apiResponseIs(expected int, expectedResponse string) error {
	defer s.stepResponse.Body.Close()

	if s.stepResponse.StatusCode != expected {
		return fmt.Errorf("expected status code %d but got %d", expected, s.stepResponse.StatusCode)
	}

	actualJson, err := getBody(s.stepResponse)
	if err != nil {
		return err
	}

	if match, err := compareJSON(expectedResponse, actualJson); err != nil {
		return fmt.Errorf("error comparing JSON: %w", err)
	} else if !match {
		log.Printf("Actual response body: %s", actualJson)
		return fmt.Errorf("response body does not match. Expected: %s, actual: %s", expectedResponse, actualJson)
	}

	return nil
}

func getBody(response *http.Response) (string, error) {
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}
	return string(body), nil
}
