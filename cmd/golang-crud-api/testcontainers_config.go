package main

import (
	"context"
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
	EnvironmentVariables   map[string]string
}

type TestContainersContext struct {
	MainHttpServer *http.Server
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
		EnvironmentVariables: map[string]string{
			"CHECK_ISBN_CLIENT_HOST": "https://api.mybiz.com",
		},
	}
}

func NewMainWithTestContainers(ctx context.Context) *TestContainersContext {
	params := NewTestContainersParams()
	// Set the environment variables
	setAppEnvs(params)
	// Start the postgres container
	initPostgresContainer(ctx, params)
	// Mock the third-party API client
	mockClient := &http.Client{}
	httpmock.ActivateNonDefault(mockClient)
	// Build the app
	server, _ := mainHttpServerSetup(params.MainHttpServerAddress, mockClient)
	return &TestContainersContext{
		MainHttpServer: server,
		Params:         params,
	}
}

// setAppEnvs sets the environment variables required by the app.
func setAppEnvs(params *TestContainersParams) {
	for key, value := range params.EnvironmentVariables {
		setenv(key, value)
	}
}

// initPostgresContainer starts a postgres container and sets the required environment variables.
// Source: https://golang.testcontainers.org/modules/postgres/
func initPostgresContainer(ctx context.Context, params *TestContainersParams) *postgres.PostgresContainer {
	logTimeout := 10
	port := "5432/tcp"
	dbURL := func(host string, port nat.Port) string {
		return fmt.Sprintf("postgres://postgres:postgres@%s:%s/%s?sslmode=disable", //nolint:nosprintfhostport // non-critical
			host, port.Port(), params.DatabaseName)
	}

	postgresContainer, err := postgres.Run(ctx, params.PostgresImage,
		postgres.WithInitScripts(filepath.Join(".", "testAssets", "testdata", "dev-db.sql")),
		postgres.WithDatabase(params.DatabaseName),
		testcontainers.WithWaitStrategy(wait.ForSQL(nat.Port(port), "postgres", dbURL).
			WithStartupTimeout(time.Duration(logTimeout)*time.Second)),
	)
	if err != nil {
		log.Fatalf("failed to start postgresContainer: %s", err)
	}
	postgresHost, _ := postgresContainer.Host(ctx) //nolint:errcheck // non-critical
	// Set envs
	setenv(params.DatabaseHostEnvVar, postgresHost)
	setenv(params.DatabasePortEnvVar, getPostgresPort(ctx, postgresContainer))
	setenv(params.DatabaseUserEnvVar, "postgres")
	setenv(params.DatabasePasswordEnvVar, "postgres")
	setenv(params.DatabaseNameEnvVar, params.DatabaseName)
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

func setenv(key, value string) {
	err := os.Setenv(key, value)

	if err != nil {
		log.Fatalf("error from setenv with %s:%s. error: %e", key, value, err)
	}
}
