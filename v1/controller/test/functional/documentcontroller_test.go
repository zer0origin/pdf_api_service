package functional

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"net/http"
	"net/http/httptest"
	"os"
	"pdf_service_api/database"
	v1 "pdf_service_api/v1"
	"pdf_service_api/v1/controller"
	"pdf_service_api/v1/models"
	"pdf_service_api/v1/models/requests"
	"pdf_service_api/v1/repositories"
	"strings"
	"testing"
	"time"
)

var dbUser = "user"
var dbPassword = "password"

func TestDatabaseConnection(t *testing.T) {
	ctx := context.Background()
	ctr, err := createTestContainerPostgres(ctx, "TestDatabaseConnection")
	assert.NoError(t, err)
	t.Cleanup(cleanUp(ctx, *ctr))

	dbConfig, err := createDbConfig(ctx, *ctr)
	assert.NoError(t, err)

	var databasePresent bool
	err = dbConfig.WithConnection(func(db *sql.DB) error { //This checks that the tables from the init script were created.
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

func TestGetDocumentHandler(t *testing.T) {
	t.Parallel()
	TestUUID := "b66fd223-515f-4503-80cc-2bdaa50ef474"

	ctx := context.Background()
	ctr, err := createTestContainerPostgres(ctx, "TestGetDatabase")
	assert.NoError(t, err)
	t.Cleanup(cleanUp(ctx, *ctr))

	dbConfig, err := createDbConfig(ctx, *ctr)
	assert.NoError(t, err)
	fmt.Println(dbConfig.ConUrl)

	repo := repositories.NewDocumentRepository(dbConfig)
	documentController := controller.NewDocumentController(repo)
	router := v1.SetupRouter(documentController)

	request := &requests.GetDocumentRequest{DocumentUuid: uuid.MustParse(TestUUID)}
	requestJSON, _ := json.Marshal(request)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, httptest.NewRequest(
		"GET",
		"/api/v1/documents/"+request.DocumentUuid.String(),
		strings.NewReader(string(requestJSON)),
	))

	responseDocument := &models.Document{}
	json.NewDecoder(w.Body).Decode(responseDocument)

	assert.Equal(t, http.StatusOK, w.Code, "Response should be 200")
	assert.Equal(t, TestUUID, responseDocument.Uuid.String(), "Response uuid does not match")
}

func createDbConfig(ctx context.Context, ctr postgres.PostgresContainer) (database.ConfigForDatabase, error) {
	connectionString, err := ctr.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		return database.ConfigForDatabase{}, err
	}

	dbConfig := database.ConfigForDatabase{
		ConUrl: connectionString,
	}

	return dbConfig, nil
}

func cleanUp(ctx context.Context, ctr postgres.PostgresContainer) func() {
	return func() {
		err := ctr.Terminate(ctx)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func createTestContainerPostgres(ctx context.Context, filename string) (ctr *postgres.PostgresContainer, err error) {
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
