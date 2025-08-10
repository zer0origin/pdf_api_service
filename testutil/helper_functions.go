package testutil

import (
	"context"
	"errors"
	"fmt"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"os"
	pg "pdf_service_api/postgres"
	"time"
)

func CreateDatabaseHandlerFromPostgresInfo(ctx context.Context, ctr postgres.PostgresContainer) (pg.DatabaseHandler, error) {
	connectionString, err := ctr.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		return pg.DatabaseHandler{}, err
	}

	dbConfig := pg.DatabaseHandler{
		DbConfig: pg.ConfigForDatabase{
			ConUrl: connectionString,
		}}

	err = dbConfig.RunInitScript()
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
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

func CreateTestContainerPostgres(ctx context.Context, dbUser string, dbPassword string, filename string) (ctr *postgres.PostgresContainer, err error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	basicScript := wd + "/test-container/sql/BasicSetup.sql"
	sqlScript := wd + "/test-container/sql/" + filename + ".sql"

	ctr, err = postgres.Run(
		ctx,
		"postgres:16-alpine",
		postgres.WithOrderedInitScripts(basicScript, sqlScript),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
		testcontainers.WithWaitStrategy(wait.ForLog("database system is ready to accept connections").
			WithOccurrence(2).
			WithStartupTimeout(30*time.Second)),
	)

	if err != nil {
		return nil, err
	}

	if ctr == nil {
		return nil, errors.New("failed to create postgres container")
	}

	p, _ := ctr.MappedPort(ctx, "5432")
	fmt.Printf("Postgres container listening to: %s\n", p)

	return ctr, nil
}
