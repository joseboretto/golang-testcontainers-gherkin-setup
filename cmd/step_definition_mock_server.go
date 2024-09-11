package main

import (
	"fmt"
	"github.com/cucumber/godog"
	"github.com/jarcoal/httpmock"
	"io"
	"log"
	"net/http"
)

// RegisterMockServerSteps registers all the step definition functions related to the mock server
func (s *StepsContext) RegisterMockServerSteps(ctx *godog.ScenarioContext) {
	ctx.Step(`^a mock server request with method: "([^"]*)" and url: "([^"]*)"$`, s.storeMockServerMethodAndUrlInStepContext)
	ctx.Step(`^a mock server request with method: "([^"]*)" and url: "([^"]*)" and body$`, s.storeMockServerMethodAndUrlAndRequestBodyInStepContext)
	ctx.Step(`^a mock server response with status (\d+) and body$`, s.setupRegisterResponder)
	ctx.Step(`^reset mock server$`, s.resetMockServer)
}

func (s *StepsContext) storeMockServerMethodAndUrlInStepContext(method, url string) error {
	s.stepMockServerRequestMethod = &method
	s.stepMockServerRequestUrl = &url
	s.stepMockServerRequestBody = nil
	return nil
}

func (s *StepsContext) storeMockServerMethodAndUrlAndRequestBodyInStepContext(method, url, body string) error {
	s.stepMockServerRequestMethod = &method
	s.stepMockServerRequestUrl = &url
	s.stepMockServerRequestBody = &body
	return nil
}

func (s *StepsContext) setupRegisterResponder(statusCode int, responseBody string) error {
	if s.stepMockServerRequestMethod == nil || s.stepMockServerRequestUrl == nil {
		log.Fatal("stepMockServerRequestMethod or stepMockServerRequestUrl is nil. You have to setup the storeMockServerMethodAndUrlInStepContext step first")
	}
	if s.stepMockServerRequestBody != nil {
		s.registerResponderWithBodyCheck(statusCode, responseBody)
	} else {
		httpmock.RegisterResponder(*s.stepMockServerRequestMethod, *s.stepMockServerRequestUrl, httpmock.NewStringResponder(statusCode, responseBody))
	}
	return nil
}

func (s *StepsContext) registerResponderWithBodyCheck(statusCode int, responseBody string) {
	httpmock.RegisterResponder(*s.stepMockServerRequestMethod, *s.stepMockServerRequestUrl,
		func(req *http.Request) (*http.Response, error) {
			body, err := io.ReadAll(req.Body)
			if err != nil {
				return nil, fmt.Errorf("failed to read request body: %w", err)
			}

			if match, err := compareJSON(*s.stepMockServerRequestBody, string(body)); err != nil {
				return nil, fmt.Errorf("error comparing JSON: %w", err)
			} else if !match {
				log.Printf("Actual request body without escapes: %s", body)
				return nil, fmt.Errorf("request body does not match for method: %s and url: %s. Expected: %s, actual: %s",
					*s.stepMockServerRequestMethod, *s.stepMockServerRequestUrl, *s.stepMockServerRequestBody, string(body))
			}

			return httpmock.NewStringResponse(statusCode, responseBody), nil
		})
}

func (s *StepsContext) resetMockServer() error {
	httpmock.Reset()
	// Reset the step context
	s.stepMockServerRequestMethod = nil
	s.stepMockServerRequestUrl = nil
	s.stepMockServerRequestBody = nil
	return nil
}
