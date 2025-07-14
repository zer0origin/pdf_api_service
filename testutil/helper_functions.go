package testutil

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"os"
	v2 "pdf_service_api/controller/v1"
	pg "pdf_service_api/postgres"
	"testing"
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

	if ctr == nil {
		return nil, errors.New("failed to create postgres container")
	}

	p, _ := ctr.MappedPort(ctx, "5432")
	fmt.Printf("Postgres container listening to: %s\n", p)

	return ctr, nil
}

func CreateV1RouterAndPostgresContainer(t *testing.T, fileName string, dbUser string, dbPassword string) *gin.Engine {
	ctx := context.Background()
	ctr, err := CreateTestContainerPostgres(ctx, fileName, dbUser, dbPassword)
	if err != nil {
		assert.FailNow(t, err.Error())
	}
	t.Cleanup(CleanUp(ctx, *ctr))

	dbConfig, err := CreateDatabaseHandlerFromPostgresInfo(ctx, *ctr)
	if err != nil {
		assert.FailNow(t, err.Error())
	}

	selectionController := &v2.SelectionController{SelectionRepository: pg.NewSelectionRepository(dbConfig)}
	repo := pg.NewDocumentRepository(dbConfig)
	documentController := &v2.DocumentController{DocumentRepository: repo, SelectionController: selectionController}
	router := v2.SetupRouter(documentController)
	return router
}
