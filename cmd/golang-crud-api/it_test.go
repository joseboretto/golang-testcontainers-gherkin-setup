package main

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/cucumber/godog"
)

func TestFeatures(t *testing.T) {
	ctx := context.Background()
	// Get testcontainersConfig
	testcontainersConfig := NewMainWithTestContainers(ctx)
	// Start the app within a goroutine to avoid blocking the test execution
	go func() {
		// Start and continue the execution.
		log.Println("Listing for requests at http://localhost" + testcontainersConfig.Params.MainHttpServerAddress)
		err := testcontainersConfig.MainHttpServer.ListenAndServe()
		if err != nil {
			log.Fatal("Error staring the server on", err)
			return
		}
	}()
	// Wait for the app to start
	time.Sleep(2 * time.Second)
	// Run the godog test suite
	suite := godog.TestSuite{
		ScenarioInitializer: func(sc *godog.ScenarioContext) {
			// This address should match the address of the app in the testcontainers_config.go file
			mainHttpServerUrl := "http://localhost" + testcontainersConfig.Params.MainHttpServerAddress
			NewStepsContext(mainHttpServerUrl, sc)
		},
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"features"}, // Edit this path locally to execute only the feature files you want to test. (features/seller-catalog-stream/getProductsV2.feature)
			TestingT: t,                    // Testing instance that will run subtests.
		},
	}
	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}
