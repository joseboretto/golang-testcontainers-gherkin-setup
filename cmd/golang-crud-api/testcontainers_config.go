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

var appAddr = ":8000"

type TestcontainersConfig struct {
	App               *http.Server
	PostgresContainer *postgres.PostgresContainer
	DeferFn           func()
}

func GetTestcontainersConfig(ctx context.Context) *TestcontainersConfig {
	// todo: check if we need a singleton when multiple tests are running
	return mainWithTestContainers(ctx)
}

func (t TestcontainersConfig) TerminateAllContainers(ctx context.Context) {
	terminateAllContainers(t.PostgresContainer, ctx)
}

func mainWithTestContainers(ctx context.Context) *TestcontainersConfig {
	setAppEnvs()
	postgresContainer := initPostgresContainer(ctx)
	// Start the app
	// Set the mock client
	mockClient := &http.Client{}
	httpmock.ActivateNonDefault(mockClient)

	server, deferFn := httpServerSetup(appAddr, mockClient)
	defer deferFn()

	go func() {
		// Start and continue the execution.
		log.Println("Listing for requests at http://localhost" + appAddr)
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal("Error staring the server on", err)
			return
		}
	}()
	// Wait for the app to start
	time.Sleep(2 * time.Second) //nolint:gomnd // non-critical
	// Handle panic
	defer func() {
		if err := recover(); err != nil {
			log.Println("panic occurred:", err)
			terminateAllContainers(postgresContainer, ctx)
		}
	}()

	return &TestcontainersConfig{
		App:               server,
		PostgresContainer: postgresContainer,
		DeferFn:           deferFn,
	}
}

// setAppEnvs sets the environment variables required by the app.
func setAppEnvs() {
	setenv("CHECK_ISBN_CLIENT_HOST", "https://api.mybiz.com")
}

// initPostgresContainer starts a postgres container and sets the required environment variables.
// Source: https://golang.testcontainers.org/modules/postgres/
func initPostgresContainer(ctx context.Context) *postgres.PostgresContainer {
	dbName := "db"
	logTimeout := 10
	port := "5432/tcp"
	dbURL := func(host string, port nat.Port) string {
		return fmt.Sprintf("postgres://postgres:postgres@%s:%s/%s?sslmode=disable", //nolint:nosprintfhostport // non-critical
			host, port.Port(), dbName)
	}

	postgresContainer, err := postgres.Run(ctx, "docker.io/postgres:15.2-alpine",
		postgres.WithInitScripts(filepath.Join(".", "testAssets", "testdata", "dev-db.sql")),
		postgres.WithDatabase(dbName),
		testcontainers.WithWaitStrategy(wait.ForSQL(nat.Port(port), "postgres", dbURL).
			WithStartupTimeout(time.Duration(logTimeout)*time.Second)),
	)
	if err != nil {
		log.Fatalf("failed to start postgresContainer: %s", err)
	}
	postgresHost, _ := postgresContainer.Host(ctx) //nolint:errcheck // non-critical
	// Set envs. Source: persistence.Config + default postgres.RunContainer values
	setenv("DATABASE_HOST", postgresHost)
	setenv("DATABASE_PORT", getPostgresPort(ctx, postgresContainer))
	setenv("DATABASE_USER", "postgres")
	setenv("DATABASE_PASSWORD", "postgres")
	setenv("DATABASE_NAME", dbName)
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

func terminateAllContainers(postgresContainer *postgres.PostgresContainer, ctx context.Context) {
	if err := postgresContainer.Terminate(ctx); err != nil {
		log.Fatalf("failed to terminate postgresContainer: %s", err)
	}
}

func setenv(key, value string) {
	err := os.Setenv(key, value)

	if err != nil {
		log.Fatalf("error from setenv with %s:%s. error: %e", key, value, err)
	}
}
