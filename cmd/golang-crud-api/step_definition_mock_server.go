package main

import (
	"github.com/cucumber/godog"
	"github.com/jarcoal/httpmock"
	"io/ioutil"
	"net/http"
)

// RegisterMockServerSteps register all the step definition functions related to the mock server
func (s *StepsContext) RegisterMockServerSteps(ctx *godog.ScenarioContext) {
	ctx.Step(`^I setup a mock server for "([^"]*)" "([^"]*)" with response (\d+) and body$`, s.setupRegisterResponder)
}

// Function to setup a basic responder with httpmock using parameters
func (s *StepsContext) setupRegisterResponder(method, url string, statusCode int, responseBody string) error {
	// Register a basic responder
	httpmock.RegisterResponder(method, url,
		httpmock.NewStringResponder(statusCode, responseBody))

	return nil
}

// Function to setup a custom matcher responder with httpmock using parameters
func setupRegisterMatcherResponder(method, url, requestBody, responseBody string, statusCode, failureStatusCode int) error {
	// Register a custom matcher responder
	httpmock.RegisterResponder(method, url,
		func(req *http.Request) (*http.Response, error) {
			body, _ := ioutil.ReadAll(req.Body)
			if string(body) == requestBody {
				return httpmock.NewStringResponse(statusCode, responseBody), nil
			}
			return httpmock.NewStringResponse(failureStatusCode, `{"error": "invalid input"}`), nil
		})

	return nil
}
