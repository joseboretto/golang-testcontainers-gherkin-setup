package main

import (
	"database/sql"
	"github.com/cucumber/godog"
	"net/http"
)

type StepsContext struct {
	// Main setup
	mainHttpServerUrl string // http://localhost:8000
	database          *sql.DB
	// Mock server setup
	stepMockServerRequestMethod *string
	stepMockServerRequestUrl    *string
	stepMockServerRequestBody   *string
	stepResponse                *http.Response
}

func NewStepsContext(mainHttpServerUrl string, database *sql.DB, sc *godog.ScenarioContext) *StepsContext {
	s := &StepsContext{
		mainHttpServerUrl: mainHttpServerUrl,
		database:          database,
	}
	// Register all the step definition function
	s.RegisterMockServerSteps(sc)
	s.RegisterDatabaseSteps(sc)
	s.RegisterApiSteps(sc)
	return s
}
