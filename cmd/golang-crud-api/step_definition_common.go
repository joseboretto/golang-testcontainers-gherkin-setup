package main

import (
	"gorm.io/gorm"
	"net/http"
)

type StepsContext struct {
	apiBaseUrl   string // http://localhost:8000
	database     *gorm.DB
	ResponseData *http.Response
}

func NewStepsContext(apiBaseUrl string) *StepsContext {
	db := getDatabaseConnection()
	return &StepsContext{
		apiBaseUrl: apiBaseUrl,
		database:   db,
	}
}
