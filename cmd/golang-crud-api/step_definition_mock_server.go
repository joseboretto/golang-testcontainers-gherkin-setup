package main

import (
	"errors"
	"fmt"
	"github.com/cucumber/godog"
	"github.com/jarcoal/httpmock"
	"io"
	"log"
	"net/http"
)

// RegisterMockServerSteps register all the step definition functions related to the mock server
func (s *StepsContext) RegisterMockServerSteps(ctx *godog.ScenarioContext) {

	ctx.Step(`^a mock server request with method: "([^"]*)" and url: "([^"]*)"$`, s.storeMockServerMethodAndUrlInStepContext)
	ctx.Step(`^a mock server request with method: "([^"]*)" and url: "([^"]*)" and body$`, s.storeMockServerMethodAndUrlAndRequestBodyInStepContext)
	ctx.Step(`^a mock server response with status (\d+) and body$`, s.setupRegisterResponder)
	ctx.Step(`^reset mock server$`, s.resetMockServer)

}

// Function to setup a basic responder with httpmock using parameters
func (s *StepsContext) storeMockServerMethodAndUrlInStepContext(method, url string) error {
	s.stepMockServerRequestMethod = method
	s.stepMockServerRequestUrl = url
	s.stepMockServerRequestBody = nil
	return nil
}

func (s *StepsContext) storeMockServerMethodAndUrlAndRequestBodyInStepContext(method, url, body string) error {
	s.stepMockServerRequestMethod = method
	s.stepMockServerRequestUrl = url
	s.stepMockServerRequestBody = &body
	return nil
}

// Function to setup a basic responder with httpmock using parameters
func (s *StepsContext) setupRegisterResponder(statusCode int, responseBody string) error {
	if s.stepMockServerRequestBody != nil {
		httpmock.RegisterResponder(s.stepMockServerRequestMethod, s.stepMockServerRequestUrl,
			func(req *http.Request) (*http.Response, error) {
				body, _ := io.ReadAll(req.Body)
				// Compare json body
				if compare, err := compareJSON(*s.stepMockServerRequestBody, string(body)); err != nil {
					return nil, err
				} else if !compare {
					log.Printf("Actual request body without escapes: %s", body)
					return nil, errors.New(fmt.Sprintf("Request body doesnt match for method: %s and url: %s. Expected body: %s, \n actual body: %v",
						s.stepMockServerRequestMethod, s.stepMockServerRequestUrl, *s.stepMockServerRequestBody, string(body)))
				} else {
					return httpmock.NewStringResponse(statusCode, responseBody), nil
				}
			})
	} else {
		httpmock.RegisterResponder(s.stepMockServerRequestMethod, s.stepMockServerRequestUrl,
			httpmock.NewStringResponder(statusCode, responseBody))
	}
	return nil
}

func (s *StepsContext) resetMockServer() error {
	httpmock.Reset()
	return nil
}
