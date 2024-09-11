package main

import (
	"github.com/cucumber/godog"
	"gorm.io/gorm"
	"net/http"
)

type StepsContext struct {
	// Main setup
	mainHttpServerUrl string // http://localhost:8000
	database          *gorm.DB
	// Mock server setup
	stepMockServerRequestMethod *string
	stepMockServerRequestUrl    *string
	stepMockServerRequestBody   *string
	stepResponse                *http.Response
}

func NewStepsContext(mainHttpServerUrl string, sc *godog.ScenarioContext) *StepsContext {
	db := getDatabaseConnection()
	s := &StepsContext{
		mainHttpServerUrl: mainHttpServerUrl,
		database:          db,
	}
	// Register all the step definition function
	s.RegisterMockServerSteps(sc)
	s.RegisterDatabaseSteps(sc)
	s.RegisterApiSteps(sc)
	return s
}
