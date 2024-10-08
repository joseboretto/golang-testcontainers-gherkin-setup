# Integration Testing in Go with Cucumber, Testcontainers, and HTTPMock

Integration testing is a crucial part of the software development lifecycle that ensures different modules of an
application work seamlessly together. In the Go (Golang) ecosystem, we can use various tools and libraries to facilitate
integration testing. This post will explore how to use Cucumber, Testcontainers for our PostgreSQL database, and
HTTPMock to write robust integration tests in Go.

# Why Integration Testing Matters

Integration testing verifies the interactions between different parts of your application, such as external services,
databases, and APIs. Unlike unit tests, which test individual components in isolation, integration tests ensure the
entire system functions as intended when combined.

# Tools Overview

1. **Cucumber**: A tool that supports Behavior-Driven Development (BDD). It allows you to write human-readable test
   scenarios in the Gherkin language, which bridges the gap between technical and non-technical stakeholders.
   1. https://github.com/cucumber/godog
2. **Testcontainers**: A library that provides lightweight, disposable containers for testing. It allows you to spin up
   real instances of databases, message brokers, or any other service your application depends on, making your tests
   more reliable and closer to real-world scenarios. We will use it to start our PostgresSQL database.
   1. https://github.com/testcontainers/testcontainers-go
3. **HTTPMock**: A simple HTTP mocking library that allows you to simulate HTTP interactions in your tests. It is useful
   when you want to mock responses from external HTTP services without actually making network requests.
   1. https://github.com/jarcoal/httpmock
    
# Use case

We will have 1 application that it’s in charge of executing the following logic:

1. Implement an API to create a book.
2. Check if the ISBN of the book exist in a third party API.
3. Store the book in the database
4. Send and email using a third-party API to the address  [helloworld@gmail.com](mailto:helloworld@gmail.com)

# Integration test architecture

We will treat our application as a  **black box**. This means we have to start the application and use it as any other
user, without mocking any part of the code.

![](https://miro.medium.com/v2/resize:fit:581/1*7IKtrF9tOXLYvuYFrncPMw.png)

The folder structure of project will be

```
- cmd/  
   - features/createBook.feature  
   - testAssets/testData/dev-db.sql
   - integration_test.go  
   - main.go  
   - step_definition_api.go  
   - step_definition_common.go  
   - step_definition_database.go  
   - step_definition_mock_server.go  
   - test_utils.go  
   - testcontainers_config.go  
- go.mod
```

# Step-by-Step Guide to Integration Testing in Go

Let’s walk through the process of setting up integration tests in Go using Cucumber, Testcontainers, and HTTPMock.

## 1. Installing Required Libraries

Next, install the necessary libraries

```bash
go get -u github.com/cucumber/godog  
go get -u github.com/testcontainers/testcontainers-go  
go get -u github.com/jarcoal/httpmock  
go get -u github.com/lib/pq  
go get -u gorm.io/driver/postgres
```

## 2. Update your main.go to control the execution fo the app

Now, we have to change a bit our main.go. Separate your main in 3 different functions.

1. `func mainHttpServerSetup(addr string, httpClient *http.Client) (*http.Server, func())`
    - This function is used in testcontainers_config.go allowing us control the execution of the app in
      integration_test.go
2. `main()`
    - The entry point of the actual application.

```go
package main

import (
	"fmt"
	servicebook "github.com/joseboretto/golang-testcontainers-gherkin-setup/internal/application/services/books"
	controller "github.com/joseboretto/golang-testcontainers-gherkin-setup/internal/infrastructure/controllers"
	controllerbook "github.com/joseboretto/golang-testcontainers-gherkin-setup/internal/infrastructure/controllers/books"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"os"
	"time"

	"log"
	"net/http"

	clientsbook "github.com/joseboretto/golang-testcontainers-gherkin-setup/internal/infrastructure/clients/book"
	persistancebook "github.com/joseboretto/golang-testcontainers-gherkin-setup/internal/infrastructure/persistance/book"
)

func main() {
	addr := ":8000"
	httpClient := &http.Client{
		Timeout: 5 * time.Second,
	}
	server, deferFn := mainHttpServerSetup(addr, httpClient)
	log.Println("Listing for requests at http://localhost" + addr)
	err := server.ListenAndServe()
	if err != nil {
		panic("Error staring the server: " + err.Error())
	}
	defer deferFn()
}

func mainHttpServerSetup(addr string, httpClient *http.Client) (*http.Server, func()) {
	db := getDatabaseConnection()
	// Migrate the schema
	err := db.AutoMigrate(&persistancebook.BookEntity{})
	if err != nil {
		panic("Error executing db.AutoMigrate" + err.Error())
	}
	// Clients
	checkIsbnClientHost := os.Getenv("CHECK_ISBN_CLIENT_HOST")
	checkIsbnClient := clientsbook.NewCheckIsbnClient(checkIsbnClientHost, httpClient)
	emailClientHost := os.Getenv("EMAIL_CLIENT_HOST")
	sendEmailClient := clientsbook.NewSendEmailClient(emailClientHost, httpClient)
	// repositories
	newCreateBookRepository := persistancebook.NewCreateBookRepository(db)
	// services
	createBookService := servicebook.NewCreateBookService(newCreateBookRepository, checkIsbnClient, sendEmailClient)
	// controllers
	bookController := controllerbook.NewBookController(createBookService)
	// routes
	controller.SetupRoutes(bookController)
	// Server
	server := &http.Server{Addr: addr, Handler: nil}
	// Defer function
	// Add all defer
	deferFn := func() {
		fmt.Println("closing database connection")
		sqlDB, err := db.DB()
		err = sqlDB.Close()
		if err != nil {
			panic("Error closing database connection: " + err.Error())
		}

	}
	return server, deferFn
}

func getDatabaseConnection() *gorm.DB {
	// Read database configuration from environment variables
	databaseUser := os.Getenv("DATABASE_USER")
	databasePassword := os.Getenv("DATABASE_PASSWORD")
	databaseName := os.Getenv("DATABASE_NAME")
	databaseHost := os.Getenv("DATABASE_HOST")
	databasePort := os.Getenv("DATABASE_PORT")
	databaseConnectionString := `host=` + databaseHost + ` user=` + databaseUser + ` password=` + databasePassword + ` dbname=` + databaseName + ` port=` + databasePort + ` sslmode=disable`
	fmt.Println("databaseConnectionString: ", databaseConnectionString)
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  databaseConnectionString,
		PreferSimpleProtocol: true, // disables implicit prepared statement usage
	}), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "myschema.", // schema name
			SingularTable: false,
		}})
	// Check if connection is successful
	if err != nil {
		panic("failed to connect database with error: " + err.Error() + "\n" + "Please check your database configuration: " + databaseConnectionString)
	}
	return db
}
```

## 3. Writing the Feature File (Cucumber)

Create a feature file (`createBook.feature`) inside the `features` directory with scenarios that describe the expected
behavior of your application.

```gherkin

Feature: Create book

   Background: Clean database
      Given SQL command
    """
    DELETE FROM myschema.books;
    """
      And reset mock server

   Scenario: Create a new book successfully
      Given a mock server request with method: "GET" and url: "https://api.isbncheck.com/isbn/0-061-96436-1"
      And a mock server response with status 200 and body
    """json
    {
       "id": "0-061-96436-1"
    }
    """
      And a mock server request with method: "POST" and url: "https://api.gmail.com/send-email" and body
    """json
    {
      "email" : "helloworld@gmail.com",
      "book" : {
        "isbn" : "0-061-96436-1",
        "title" : "The Art of Computer Programming"
      }
    }
    """
      And a mock server response with status 200 and body
    """json
    {
       "status": "OK"
    }
    """
      When API "POST" request is sent to "/api/v1/createBook" with payload
    """json
    {
      "isbn": "0-061-96436-1",
      "title": "The Art of Computer Programming"
    }
    """
      Then API response status code is 200 and payload is
    """json
    {
        "isbn": "0-061-96436-1",
        "title": "The Art of Computer Programming"
    }
    """
      And SQL query "SELECT * FROM myschema.books WHERE isbn = '0-061-96436-1'" result without the fields "created_at,deleted_at,updated_at" is equal to
    """json
    [
       {
          "id":1,
          "isbn":"0-061-96436-1",
          "title":"The Art of Computer Programming"
       }
    ]
    """
```

## 4. Implementing Step Definitions in Go

Now, create the Go files to implement the step definitions for the scenarios. We have multiple step definitions based on
the goals.

> The most important concept to remember is writing generic step definitions to re-use them on every use case.

#### step_definition_common.go

Share data between the app setup and the step definition.

```go
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

```

#### step_definition_api.go

Perform api request and check the response.

```go
package main

import (
	"bytes"
	"context"
	"fmt"
	"github.com/cucumber/godog"
	"io"
	"log"
	"net/http"
	"time"
)

func (s *StepsContext) RegisterApiSteps(sc *godog.ScenarioContext) {
	sc.Step(`^API "([^"]*)" request is sent to "([^"]*)" without payload$`, s.apiRequestIsSendWithoutPayload)
	sc.Step(`^API "([^"]*)" request is sent to "([^"]*)" with payload$`, s.apiRequestIsSendWithPayload)
	sc.Step(`^API response status code is (\d+) and payload is$`, s.apiResponseIs)
}

func (s *StepsContext) apiRequestIsSendWithoutPayload(method, url string) error {
	return s.apiRequestIsSendWithPayload(method, url, "")
}

func (s *StepsContext) apiRequestIsSendWithPayload(method, url, payloadJson string) error {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	req, err := http.NewRequestWithContext(context.Background(), method, s.mainHttpServerUrl+url, bytes.NewBufferString(payloadJson))
	if err != nil {
		return fmt.Errorf("failed to create a new HTTP request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	response, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute HTTP request: %w", err)
	}

	s.stepResponse = response
	return nil
}

func (s *StepsContext) apiResponseIs(expected int, expectedResponse string) error {
	defer s.stepResponse.Body.Close()

	if s.stepResponse.StatusCode != expected {
		return fmt.Errorf("expected status code %d but got %d", expected, s.stepResponse.StatusCode)
	}

	actualJson, err := getBody(s.stepResponse)
	if err != nil {
		return err
	}

	if match, err := compareJSON(expectedResponse, actualJson); err != nil {
		return fmt.Errorf("error comparing JSON: %w", err)
	} else if !match {
		log.Printf("Actual response body: %s", actualJson)
		return fmt.Errorf("response body does not match. Expected: %s, actual: %s", expectedResponse, actualJson)
	}

	return nil
}

func getBody(response *http.Response) (string, error) {
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}
	return string(body), nil
}

```

#### step_definition_database.go

store data in the database. Additionally, you have to implement the step definitions to check the data stored in the
database.

```go
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
```

#### test_utils.go

Mainly, this file contains the function to compare JSON strings.
```go

package main

import (
   "encoding/json"
   "fmt"
   "reflect"
)

// compareJSON compares two JSON strings and returns true if they are equal.
func compareJSON(jsonStr1, jsonStr2 string) (bool, error) {
   var obj1, obj2 map[string]interface{}

   // Unmarshal the first JSON string
   if err := json.Unmarshal([]byte(jsonStr1), &obj1); err != nil {
      return false, fmt.Errorf("error unmarshalling jsonStr1: %v", err)
   }

   // Unmarshal the second JSON string
   if err := json.Unmarshal([]byte(jsonStr2), &obj2); err != nil {
      return false, fmt.Errorf("error unmarshalling jsonStr2: %v", err)
   }
   // Compare the two maps
   return reflect.DeepEqual(obj1, obj2), nil
}

func compareJSONArrays(jsonStr1, jsonStr2 string) (bool, error) {
   var obj1, obj2 []map[string]interface{}

   // Unmarshal the first JSON string
   if err := json.Unmarshal([]byte(jsonStr1), &obj1); err != nil {
      return false, fmt.Errorf("error unmarshalling jsonStr1: %v", err)
   }

   // Unmarshal the second JSON string
   if err := json.Unmarshal([]byte(jsonStr2), &obj2); err != nil {
      return false, fmt.Errorf("error unmarshalling jsonStr2: %v", err)
   }
   // Compare the two maps
   return reflect.DeepEqual(obj1, obj2), nil
}

```

## 5. Configure test containers

#### testcontainers_config.go

> Configure  **TestContainersParams**  to match the envs and the config of you actual application. Additionally, **httpmock.ActivateNonDefault(mockClient)**  is key ot make  **httpmock**  works because it will only mock our third
> party http.client

```go
package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jarcoal/httpmock"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/docker/go-connections/nat"

	_ "github.com/lib/pq" // Import the postgres driver
	"github.com/testcontainers/testcontainers-go/wait"
)

type TestContainersParams struct {
	MainHttpServerAddress  string
	PostgresImage          string
	DatabaseName           string
	DatabaseHostEnvVar     string
	DatabasePortEnvVar     string
	DatabaseUserEnvVar     string
	DatabasePasswordEnvVar string
	DatabaseNameEnvVar     string
	DatabaseInitScript     string
	EnvironmentVariables   map[string]string
}

type TestContainersContext struct {
	MainHttpServer *http.Server
	Database       *sql.DB
	Params         *TestContainersParams
}

func NewTestContainersParams() *TestContainersParams {
	return &TestContainersParams{
		MainHttpServerAddress:  ":8000",
		PostgresImage:          "docker.io/postgres:16-alpine",
		DatabaseName:           "db",
		DatabaseHostEnvVar:     "DATABASE_HOST",
		DatabasePortEnvVar:     "DATABASE_PORT",
		DatabaseUserEnvVar:     "DATABASE_USER",
		DatabasePasswordEnvVar: "DATABASE_PASSWORD",
		DatabaseNameEnvVar:     "DATABASE_NAME",
		DatabaseInitScript:     filepath.Join(".", "testAssets", "testdata", "dev-db.sql"),
		EnvironmentVariables: map[string]string{
			"CHECK_ISBN_CLIENT_HOST": "https://api.isbncheck.com",
			"EMAIL_CLIENT_HOST":      "https://api.gmail.com",
		},
	}
}

func NewMainWithTestContainers(ctx context.Context) *TestContainersContext {
	params := NewTestContainersParams()
	// Set the environment variables
	setEnvVars(params.EnvironmentVariables)
	// Start the postgres container
	initPostgresContainer(ctx, params)
	// Create database connection
	db := getDatabaseConnectionTestContainers(params)
	// Mock the third-party API client
	mockClient := &http.Client{}
	httpmock.ActivateNonDefault(mockClient)
	// Build the app
	server, _ := mainHttpServerSetup(params.MainHttpServerAddress, mockClient)
	return &TestContainersContext{
		MainHttpServer: server,
		Database:       db,
		Params:         params,
	}
}

func getDatabaseConnectionTestContainers(params *TestContainersParams) *sql.DB {
	db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv(params.DatabaseHostEnvVar),
		os.Getenv(params.DatabasePortEnvVar),
		os.Getenv(params.DatabaseUserEnvVar),
		os.Getenv(params.DatabasePasswordEnvVar),
		os.Getenv(params.DatabaseNameEnvVar),
	))
	if err != nil {
		log.Fatalf("failed to connect to database: %s", err)
	}
	return db
}

// initPostgresContainer starts a postgres container and sets the required environment variables.
// Source: https://golang.testcontainers.org/modules/postgres/
func initPostgresContainer(ctx context.Context, params *TestContainersParams) *postgres.PostgresContainer {
	logTimeout := 10
	port := "5432/tcp"
	dbURL := func(host string, port nat.Port) string {
		return fmt.Sprintf("postgres://postgres:postgres@%s:%s/%s?sslmode=disable",
			host, port.Port(), params.DatabaseName)
	}

	postgresContainer, err := postgres.Run(ctx, params.PostgresImage,
		postgres.WithInitScripts(params.DatabaseInitScript),
		postgres.WithDatabase(params.DatabaseName),
		testcontainers.WithWaitStrategy(wait.ForSQL(nat.Port(port), "postgres", dbURL).
			WithStartupTimeout(time.Duration(logTimeout)*time.Second)),
	)
	if err != nil {
		log.Fatalf("failed to start postgresContainer: %s", err)
	}
	postgresHost, _ := postgresContainer.Host(ctx) //nolint:errcheck // non-critical
	// Set database environment variables
	setEnvVars(map[string]string{
		params.DatabaseHostEnvVar:     postgresHost,
		params.DatabasePortEnvVar:     getPostgresPort(ctx, postgresContainer),
		params.DatabaseUserEnvVar:     "postgres",
		params.DatabasePasswordEnvVar: "postgres",
		params.DatabaseNameEnvVar:     params.DatabaseName,
	})

	s, err := postgresContainer.ConnectionString(ctx)
	log.Printf("Postgres container started at: %s", s)
	if err != nil {
		log.Fatalf("error from initPostgresContainer - ConnectionString: %e", err)
	}

	return postgresContainer

}

func getPostgresPort(ctx context.Context, postgresContainer *postgres.PostgresContainer) string {
	port, e := postgresContainer.MappedPort(ctx, "5432/tcp")
	if e != nil {
		log.Fatalf("Failed to get postgresContainer.Ports: %s", e)
	}

	return port.Port()
}

// setEnvVars sets environment variables from a map.
func setEnvVars(envs map[string]string) {
	for key, value := range envs {
		if err := os.Setenv(key, value); err != nil {
			log.Fatalf("Failed to set environment variable %s: %v", key, err)
		}
	}
}

```

## 7. Configure cucumber to run the integration test

#### integration_test.go

```go
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

```

## 8. Execute the test

```bash
go test -v
=== RUN   TestFeatures
2024/09/12 16:09:52 github.com/testcontainers/testcontainers-go - Connected to docker: 
  Server Version: 24.0.7
  API Version: 1.43
  Operating System: Ubuntu 23.10
  Total Memory: 5910 MB
  Testcontainers for Go Version: v0.33.0
  Resolved Docker Host: unix:////Users/jose.boretto/.colima/docker.sock
  Resolved Docker Socket Path: /var/run/docker.sock
  Test SessionID: 1c9d3b9a409b2a340378c42d67678f67249a8f748d346a0dd3758b2ce56c3bda
  Test ProcessID: f1928595-5aaa-4f92-9c93-e2f2bad03364
2024/09/12 16:09:52 🐳 Creating container for image testcontainers/ryuk:0.8.1
2024/09/12 16:09:52 ✅ Container created: c0b1cf852231
2024/09/12 16:09:52 🐳 Starting container: c0b1cf852231
2024/09/12 16:09:52 ✅ Container started: c0b1cf852231
2024/09/12 16:09:52 ⏳ Waiting for container id c0b1cf852231 image: testcontainers/ryuk:0.8.1. Waiting for: &{Port:8080/tcp timeout:<nil> PollInterval:100ms skipInternalCheck:false}
2024/09/12 16:09:53 🔔 Container is ready: c0b1cf852231
2024/09/12 16:09:53 🐳 Creating container for image docker.io/postgres:16-alpine
2024/09/12 16:09:53 ✅ Container created: 3980cd9f254e
2024/09/12 16:09:53 🐳 Starting container: 3980cd9f254e
2024/09/12 16:09:53 ✅ Container started: 3980cd9f254e
2024/09/12 16:09:53 ⏳ Waiting for container id 3980cd9f254e image: docker.io/postgres:16-alpine. Waiting for: &{timeout:<nil> deadline:0x14000011148 Strategies:[0x140000aa780]}
2024/09/12 16:09:56 🔔 Container is ready: 3980cd9f254e
2024/09/12 16:09:56 Postgres container started at: postgres://postgres:postgres@localhost:32809/db?
databaseConnectionString:  host=localhost user=postgres password=postgres dbname=db port=32809 sslmode=disable
2024/09/12 16:09:56 Listening for requests at http://localhost:8000
Feature: Create book
=== RUN   TestFeatures/Create_a_new_book_successfully

  Background: Clean database
    Given SQL command                                                                                                                                  # <autogenerated>:1 -> *StepsContext
      """
      DELETE FROM myschema.books;
      """
    And reset mock server                                                                                                                              # <autogenerated>:1 -> *StepsContext

  Scenario: Create a new book successfully                                                                                                             # features/createBook.feature:10
    Given a mock server request with method: "GET" and url: "https://api.isbncheck.com/isbn/0-061-96436-1"                                             # <autogenerated>:1 -> *StepsContext
    And a mock server response with status 200 and body                                                                                                # <autogenerated>:1 -> *StepsContext
      """ json
      {
         "id": "0-061-96436-1"
      }
      """
    And a mock server request with method: "POST" and url: "https://api.gmail.com/send-email" and body                                                 # <autogenerated>:1 -> *StepsContext
      """ json
      {
        "email" : "helloworld@gmail.com",
        "book" : {
          "isbn" : "0-061-96436-1",
          "title" : "The Art of Computer Programming"
        }
      }
      """
    And a mock server response with status 200 and body                                                                                                # <autogenerated>:1 -> *StepsContext
      """ json
      {
         "status": "OK"
      }
      """

2024/09/12 16:09:57 /Users/jose.boretto/Documents/jose/golang-testcontainers-gherkin-setup/internal/infrastructure/persistance/book/create_book_repository.go:40 record not found
[3.278ms] [rows:0] SELECT * FROM "myschema"."books" WHERE isbn = '0-061-96436-1' AND "books"."deleted_at" IS NULL ORDER BY "books"."id" LIMIT 1
    When API "POST" request is sent to "/api/v1/createBook" with payload                                                                               # <autogenerated>:1 -> *StepsContext
      """ json
      {
        "isbn": "0-061-96436-1",
        "title": "The Art of Computer Programming"
      }
      """
    Then API response status code is 200 and payload is                                                                                                # <autogenerated>:1 -> *StepsContext
      """ json
      {
          "isbn": "0-061-96436-1",
          "title": "The Art of Computer Programming"
      }
      """
    And SQL query "SELECT * FROM myschema.books WHERE isbn = '0-061-96436-1'" result without the fields "created_at,deleted_at,updated_at" is equal to # <autogenerated>:1 -> *StepsContext
      """ json
      [
         {
            "id":1,
            "isbn":"0-061-96436-1",
            "title":"The Art of Computer Programming"
         }
      ]
      """

1 scenarios (1 passed)
9 steps (9 passed)
30.273417ms
--- PASS: TestFeatures (4.66s)
    --- PASS: TestFeatures/Create_a_new_book_successfully (0.03s)
PASS
ok      github.com/joseboretto/golang-testcontainers-gherkin-setup/cmd  5.087s
```

# Conclusion

Combining Cucumber, Testcontainers, and HTTPMock gives you a powerful suite of tools for writing integration tests in
Go. With this setup, you can ensure your application is well-tested, reliable, and ready for production deployment.

By leveraging the strengths of each tool, you can achieve more comprehensive test coverage and ensure that all parts of
your system work together seamlessly.

Feel free to experiment with these tools and extend the example to suit your own application’s needs!

Here you can find the repository with the
code:  [https://github.com/joseboretto/golang-testcontainers-gherkin-setup](https://github.com/joseboretto/golang-testcontainers-gherkin-setup)