package main

import (
	"github.com/cucumber/godog"
	"gorm.io/gorm"
	"net/http"
)

type StepsContext struct {
	mainHttpServerUrl string // http://localhost:8000
	database          *gorm.DB
	responseData      *http.Response
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
