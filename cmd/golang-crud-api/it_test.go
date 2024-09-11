package main

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/cucumber/godog"
)

func TestFeatures(t *testing.T) {
	// Run the godog test suite
	suite := godog.TestSuite{
		ScenarioInitializer: InitializeScenario,
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

func InitializeScenario(sc *godog.ScenarioContext) {
	// This address should match the address of the app in the testcontainers_config.go file
	stepsContext := NewStepsContext("http://localhost" + appAddr)
	// Register steps
	stepsContext.RegisterMockServerSteps(sc)
	stepsContext.RegisterDatabaseSteps(sc)
	stepsContext.RegisterApiSteps(sc)
}

// TestMain is the setup function for the tests in this package.
// It sets up the test environment, runs the tests and then tears down the environment.
// It uses testcontainers to manage the test environment.
// Docs: https://pkg.go.dev/testing#hdr-Main
func TestMain(m *testing.M) {
	// Setup code goes here
	ctx := context.Background()
	// Get testcontainersConfig
	testcontainersConfig := GetTestcontainersConfig(ctx)
	// Execute test
	code := m.Run()
	// Teardown code goes here
	err := testcontainersConfig.App.Shutdown(ctx)
	if err != nil {
		log.Fatalf("failed to App.Shutdown: %e", err)

		return
	}
	testcontainersConfig.TerminateAllContainers(ctx)
	testcontainersConfig.DeferFn()
	// Exit

	os.Exit(code)
}
