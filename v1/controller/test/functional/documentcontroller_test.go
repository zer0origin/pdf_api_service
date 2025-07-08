package functional

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"os"
	"pdf_service_api/database"
	"testing"
	"time"
)

var dbUser = "user"
var dbPassword = "password"

func TestDatabaseConnection(t *testing.T) {
	ctx := context.Background()
	ctr, err := createTestContainerPostgres(ctx)
	fmt.Println(ctr.ConnectionString(ctx, "sslmode=disable"))
	t.Cleanup(func() {
		err := ctr.Terminate(ctx)
		if err != nil {
			fmt.Println(err)
		}
	})

	port, err := ctr.Container.MappedPort(ctx, "5432/tcp")
	assert.NoError(t, err, "could not find a port for postgres")
	portStr := port.Port()

	dbConfig := database.ConfigForDatabase{
		Port: portStr,
	}

	assert.NoError(t, err, "Error creating postgres container")

	var databasePresent bool

	err = dbConfig.WithConnection(func(db *sql.DB) error {
		sqlStatement := "SELECT EXISTS (SELECT FROM information_schema.tables WHERE  table_schema = $1 AND table_name   = $2);"
		row := db.QueryRow(sqlStatement, "public", "document_table")
		err := row.Scan(&databasePresent)
		if err != nil {
			return err
		}
		return nil
	})
	assert.NoError(t, err, "Error connecting to database")
	assert.True(t, databasePresent, "Database should exists")
}

func createTestContainerPostgres(ctx context.Context) (ctr *postgres.PostgresContainer, err error) {

	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	sqlScript := wd + "/test-container/sql/init.sql"

	ctr, err = postgres.Run(
		ctx,
		"postgres:16-alpine",
		postgres.WithInitScripts(sqlScript),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
		testcontainers.WithExposedPorts("5432/tcp"),
		testcontainers.WithWaitStrategy(wait.ForLog("database system is ready to accept connections").
			WithOccurrence(2).
			WithStartupTimeout(30*time.Second)),
	)

	if err != nil {
		return nil, err
	}

	return ctr, nil
}
