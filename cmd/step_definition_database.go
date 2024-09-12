package main

import (
	"encoding/json"
	"fmt"
	"github.com/cucumber/godog"
	"strings"
)

func (s *StepsContext) RegisterDatabaseSteps(sc *godog.ScenarioContext) {
	sc.Step(`^SQL command`, s.executeSQL)
	sc.Step(`^SQL query "([^"]*)" result is equal to`, s.checkSQLqueryWithoutIgnore)
	sc.Step(`^SQL query "([^"]*)" result without the fields "([^"]*)" is equal to`, s.checkSQLqueryWithIgnoredFields)
}

func (s *StepsContext) executeSQL(sqlCommand string) error {
	_, err := s.database.Exec(sqlCommand)
	if err != nil {
		return err
	}
	return nil
}

func (s *StepsContext) checkSQLqueryWithoutIgnore(query, jsonString string) error {
	return s.checkSQLqueryWithIgnoredFields(query, jsonString, "")
}

func (s *StepsContext) checkSQLqueryWithIgnoredFields(query, ignoredFields, jsonString string) error {
	// Parse ignored fields into a map for quick lookup
	ignoredFieldsSet := make(map[string]struct{})
	if ignoredFields != "" {
		for _, field := range strings.Split(ignoredFields, ",") {
			ignoredFieldsSet[strings.TrimSpace(field)] = struct{}{}
		}
	}
	// Execute the SQL query
	rows, err := s.database.Query(query)
	if err != nil {
		return fmt.Errorf("error executing query: %w", err)
	}
	defer rows.Close()

	// Fetch column names
	columns, err := rows.Columns()
	if err != nil {
		return fmt.Errorf("error fetching columns: %w", err)
	}

	// Prepare a slice of maps to hold query results
	var resultRows []map[string]interface{}

	for rows.Next() {
		// Create a slice to hold column values
		columnValues := make([]interface{}, len(columns))
		columnPointers := make([]interface{}, len(columns))

		for i := range columnValues {
			columnPointers[i] = &columnValues[i]
		}

		// Scan the result into column pointers
		if err := rows.Scan(columnPointers...); err != nil {
			return fmt.Errorf("error scanning row: %w", err)
		}

		// Create a map to represent a row
		rowMap := make(map[string]interface{})
		for i, colName := range columns {
			val := columnValues[i]

			// Convert bytes to string for easier comparison
			if b, ok := val.([]byte); ok {
				rowMap[colName] = string(b)
			} else {
				rowMap[colName] = val
			}
		}

		resultRows = append(resultRows, rowMap)
	}
	// Remove ignored fields from the actual data
	for i := range resultRows {
		for ignoredField := range ignoredFieldsSet {
			delete(resultRows[i], ignoredField)
		}
	}
	// Remove ignored fields from the expected data
	var expectedData []map[string]interface{}
	if err := json.Unmarshal([]byte(jsonString), &expectedData); err != nil {
		return fmt.Errorf("error unmarshalling provided JSON string: %w", err)
	}
	for i := range expectedData {
		for ignoredField := range ignoredFieldsSet {
			delete(expectedData[i], ignoredField)
		}
	}
	// Convert query result to JSON
	queryResultJSON, err := json.Marshal(resultRows)
	if err != nil {
		return fmt.Errorf("error marshalling query result to JSON: %w", err)
	}
	//
	if match, err := compareJSONArrays(jsonString, string(queryResultJSON)); err != nil {
		return fmt.Errorf("error comparing JSON: %w", err)
	} else if !match {

		return fmt.Errorf("actual Expected: %s, actual: \n %s", jsonString, string(queryResultJSON))
	}
	return nil
}
