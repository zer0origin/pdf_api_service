package testutil

import (
	"context"
	"fmt"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"os"
	"pdf_service_api/database"
	"time"
)

func CreateDbConfig(ctx context.Context, ctr postgres.PostgresContainer) (database.ConfigForDatabase, error) {
	connectionString, err := ctr.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		return database.ConfigForDatabase{}, err
	}

	dbConfig := database.ConfigForDatabase{
		ConUrl: connectionString,
	}

	return dbConfig, nil
}

func CleanUp(ctx context.Context, ctr postgres.PostgresContainer) func() {
	return func() {
		err := ctr.Terminate(ctx)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func CreateTestContainerPostgres(ctx context.Context, filename string, dbUser string, dbPassword string) (ctr *postgres.PostgresContainer, err error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	sqlScript := wd + "/test-container/sql/" + filename + ".sql"

	ctr, err = postgres.Run(
		ctx,
		"postgres:16-alpine",
		postgres.WithInitScripts(sqlScript),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
		testcontainers.WithWaitStrategy(wait.ForLog("database system is ready to accept connections").
			WithOccurrence(2).
			WithStartupTimeout(30*time.Second)),
	)

	if err != nil {
		return nil, err
	}

	return ctr, nil
}
