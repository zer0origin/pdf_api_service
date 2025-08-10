package integration

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	v1 "pdf_service_api/controller/v1"
	"pdf_service_api/postgres"
	"pdf_service_api/testutil"
	"strings"
	"testing"
)

var dbUser = "user"
var dbPassword = "password"

func TestDocumentIntegration(t *testing.T) {
	t.Run("Test database connection", databaseConnection)
	t.Run("Get document with present uuid", getDocumentWithPresentUUID)
	t.Run("Get document with nonexistent uuid", documentWithNonexistentUUID)
	t.Run("Upload new document", uploadDocument)
	t.Run("Delete existing document", deleteDocument)
}

func databaseConnection(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	ctr, err := testutil.CreateTestContainerPostgres(ctx, dbUser, dbPassword, "BasicSetup")
	require.NoError(t, err)
	t.Cleanup(testutil.CleanUp(ctx, *ctr))

	dbConfig, err := testutil.CreateDatabaseHandlerFromPostgresInfo(ctx, *ctr)
	require.NoError(t, err)

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
	assert.NoError(t, err, "Error connecting to postgres")
	assert.True(t, databasePresent, "Database should exists")
}

func getDocumentWithPresentUUID(t *testing.T) {
	documentTestUUID := uuid.MustParse("b66fd223-515f-4503-80cc-2bdaa50ef474")
	expectedResponse := fmt.Sprintf(`{"document":{"documentUUID":"%s","pdfBase64":"Fake document for testing"}}`, documentTestUUID)
	t.Parallel()

	ctx := context.Background()
	ctr, err := testutil.CreateTestContainerPostgres(ctx, dbUser, dbPassword, "OneDocumentTableEntry")
	require.NoError(t, err)

	connectionString, err := ctr.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	dbHandle := postgres.DatabaseHandler{DbConfig: postgres.ConfigForDatabase{ConUrl: connectionString}}
	documentCtrl := &v1.DocumentController{DocumentRepository: postgres.NewDocumentRepository(dbHandle)}
	router := v1.SetupRouter(documentCtrl, nil, nil)

	request := &v1.GetDocumentRequest{DocumentUUID: &documentTestUUID}
	requestJSON, _ := json.Marshal(request)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, httptest.NewRequest(
		"GET",
		"/api/v1/documents/?documentUUID="+request.DocumentUUID.String(),
		strings.NewReader(string(requestJSON)),
	))

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, expectedResponse, w.Body.String())
}

func documentWithNonexistentUUID(t *testing.T) {
	documentTestUUID := uuid.MustParse(uuid.Nil.String())
	expectedResponse := fmt.Sprintf(`{"error":"Document with UUID %s was found."}`, documentTestUUID)
	t.Parallel()

	ctx := context.Background()
	ctr, err := testutil.CreateTestContainerPostgres(ctx, dbUser, dbPassword, "OneDocumentTableEntry")
	require.NoError(t, err)

	connectionString, err := ctr.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	dbHandle := postgres.DatabaseHandler{DbConfig: postgres.ConfigForDatabase{ConUrl: connectionString}}
	documentCtrl := &v1.DocumentController{DocumentRepository: postgres.NewDocumentRepository(dbHandle)}
	router := v1.SetupRouter(documentCtrl, nil, nil)

	request := &v1.GetDocumentRequest{DocumentUUID: &documentTestUUID}
	requestJSON, _ := json.Marshal(request)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, httptest.NewRequest(
		"GET",
		"/api/v1/documents/?documentUUID="+request.DocumentUUID.String(),
		strings.NewReader(string(requestJSON)),
	))

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, expectedResponse, w.Body.String())
}

type UploadResponse struct {
	DocumentUUID uuid.UUID `json:"documentUUID"`
}

func uploadDocument(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	ctr, err := testutil.CreateTestContainerPostgres(ctx, dbUser, dbPassword, "BasicSetup")
	require.NoError(t, err)

	connectionString, err := ctr.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	dbHandle := postgres.DatabaseHandler{DbConfig: postgres.ConfigForDatabase{ConUrl: connectionString}}
	documentCtrl := &v1.DocumentController{DocumentRepository: postgres.NewDocumentRepository(dbHandle)}
	router := v1.SetupRouter(documentCtrl, nil, nil)
	request := &v1.CreateRequest{DocumentBase64String: "THIS IS A TEST DOCUMENT"}
	requestJSON, _ := json.Marshal(request)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, httptest.NewRequest(
		"POST",
		"/api/v1/documents/",
		strings.NewReader(string(requestJSON)),
	))

	response := UploadResponse{}
	err = json.NewDecoder(w.Body).Decode(&response)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, w.Code, "Response should be 200")
	assert.NotEqual(t, uuid.Nil, response.DocumentUUID)
}

type DeleteResponse struct {
	Success bool `json:"success"`
}

func deleteDocument(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	ctr, err := testutil.CreateTestContainerPostgres(ctx, dbUser, dbPassword, "OneDocumentTableEntry")
	require.NoError(t, err)

	connectionString, err := ctr.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	dbHandle := postgres.DatabaseHandler{DbConfig: postgres.ConfigForDatabase{ConUrl: connectionString}}
	documentCtrl := &v1.DocumentController{DocumentRepository: postgres.NewDocumentRepository(dbHandle)}
	router := v1.SetupRouter(documentCtrl, nil, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, httptest.NewRequest(
		"DELETE",
		fmt.Sprintf("/api/v1/documents/?documentUUID=%s", "b66fd223-515f-4503-80cc-2bdaa50ef474"),
		nil,
	))

	fmt.Println(w.Body.String())
	response := DeleteResponse{}
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		assert.FailNow(t, err.Error())
		return
	}

	assert.Equal(t, http.StatusOK, w.Code, "Response should be 200")
	assert.True(t, response.Success)
}
