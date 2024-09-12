package main

import (
	"context"
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/cucumber/godog"
)

func TestFeatures(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize test containers configuration
	testcontainersConfig := NewMainWithTestContainers(ctx)

	// Channel to notify when the server is ready
	serverReady := make(chan struct{})

	// Start the HTTP server in a separate goroutine
	go func() {
		log.Println("Listening for requests at http://localhost" + testcontainersConfig.Params.MainHttpServerAddress)
		// Notify that the server is ready
		close(serverReady)

		if err := testcontainersConfig.MainHttpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error starting the server: %v", err)
		}
	}()

	// Wait for the server to start
	<-serverReady

	// Allow a brief moment for the server to initialize
	time.Sleep(500 * time.Millisecond)

	// Run the godog test suite
	suite := godog.TestSuite{
		ScenarioInitializer: func(sc *godog.ScenarioContext) {
			// This address should match the address of the app in the testcontainers_config.go file
			mainHttpServerUrl := "http://localhost" + testcontainersConfig.Params.MainHttpServerAddress
			NewStepsContext(mainHttpServerUrl, testcontainersConfig.Database, sc)
		},
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"features"}, // Edit this path locally to execute only the feature files you want to test.
			TestingT: t,                    // Testing instance that will run subtests.
		},
	}

	if suite.Run() != 0 {
		t.Fatal("Non-zero status returned, failed to run feature tests")
	}

	// Gracefully shutdown the server after tests are done
	if err := testcontainersConfig.MainHttpServer.Shutdown(ctx); err != nil {
		log.Fatalf("Error shutting down the server: %v", err)
	}
}
