package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cucumber/godog"
	"io"
	"net/http"
	"reflect"
	"time"
)

func (s *StepsContext) RegisterApiSteps(sc *godog.ScenarioContext) {
	sc.Step(`^API "([^"]*)" request is sent to "([^"]*)" without payload`, s.apiRequestIsSendWithoutPayload)
	sc.Step(`^API "([^"]*)" request is sent to "([^"]*)" with payload`, s.apiRequestIsSendWithPayload)
	sc.Step(`^response status code is (\d+)$`, s.httpStatusIsEqualTo)
	sc.Step(`^response body is$`, s.responseBodyIsEqualTo)
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
		s.apiBaseUrl+url,
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
	s.ResponseData = response

	return nil
}

func (s *StepsContext) httpStatusIsEqualTo(expected int) error {
	if s.ResponseData.StatusCode != expected {
		fmt.Println("Response:" + getBody(s.ResponseData))
		return fmt.Errorf("httpStatusIsEqualTo: Expected %d but got %d", expected, s.ResponseData.StatusCode)
	}

	return nil
}

func (s *StepsContext) responseBodyIsEqualTo(expectedJson string) error {
	actualJson := getBody(s.ResponseData)
	var obj1, obj2 map[string]interface{}
	err1 := json.Unmarshal([]byte(expectedJson), &obj1)
	err2 := json.Unmarshal([]byte(actualJson), &obj2)

	if err1 != nil || err2 != nil {
		return errors.New("error unmarshalling JSON")
	}
	if !reflect.DeepEqual(obj1, obj2) {
		actualJsonPretty, err := json.MarshalIndent(obj2, "", "    ")
		if err != nil {
			return err
		}

		return fmt.Errorf("responseBody is NOT equal to expected. Actual JSON response: \n %s", actualJsonPretty)

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
