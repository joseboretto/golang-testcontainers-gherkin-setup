package main

import (
	"github.com/cucumber/godog"
)

func (s *StepsContext) RegisterDatabaseSteps(sc *godog.ScenarioContext) {
	sc.Step(`^SQL command`, s.executeSQL)
}

func (s *StepsContext) executeSQL(sqlCommand string) error {
	tx := s.database.Exec(sqlCommand)
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}
