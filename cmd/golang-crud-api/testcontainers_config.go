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
	DatabaseInitScript     string
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
