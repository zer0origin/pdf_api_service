package integration

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	v1 "pdf_service_api/controller/v1"
	"pdf_service_api/domain"
	"pdf_service_api/postgres"
	"pdf_service_api/testutil"
	"strings"
	"testing"
)

var dbUser = "user"
var dbPassword = "password"

func TestDocumentIntegration(t *testing.T) {
	t.Run("databaseConnection", databaseConnection)
	t.Run("getDocumentHandler", getDocumentHandler)
	t.Run("uploadDocument", uploadDocument)
	t.Run("deleteDocument", deleteDocument)
}

func databaseConnection(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	ctr, err := testutil.CreateTestContainerPostgres(ctx, "BasicSetup", dbUser, dbPassword)
	if err != nil {
		assert.FailNow(t, err.Error())
	}
	t.Cleanup(testutil.CleanUp(ctx, *ctr))

	dbConfig, err := testutil.CreateDatabaseHandlerFromPostgresInfo(ctx, *ctr)
	if err != nil {
		assert.FailNow(t, err.Error())
	}

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

func getDocumentHandler(t *testing.T) {
	documentTestUUID := "b66fd223-515f-4503-80cc-2bdaa50ef474"
	t.Parallel()

	ctx := context.Background()
	ctr, err := testutil.CreateTestContainerPostgres(ctx, "BasicSetupWithOneDocumentTableEntry", dbUser, dbPassword)
	if err != nil {
		assert.FailNow(t, err.Error())
		return
	}

	connectionString, err := ctr.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		assert.FailNow(t, err.Error())
		return
	}

	dbHandle := postgres.DatabaseHandler{DbConfig: postgres.ConfigForDatabase{ConUrl: connectionString}}
	documentCtrl := &v1.DocumentController{DocumentRepository: postgres.NewDocumentRepository(dbHandle)}
	router := v1.SetupRouter(documentCtrl, nil, nil)

	request := &v1.GetDocumentRequest{DocumentUuid: uuid.MustParse(documentTestUUID)}
	requestJSON, _ := json.Marshal(request)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, httptest.NewRequest(
		"GET",
		"/api/v1/documents/?documentUUID="+request.DocumentUuid.String(),
		strings.NewReader(string(requestJSON)),
	))

	responseDocument := &domain.Document{}
	if err := json.NewDecoder(w.Body).Decode(responseDocument); err != nil {
		assert.FailNow(t, err.Error())
	}

	assert.Equal(t, http.StatusOK, w.Code, "Response should be 200")
	assert.Equal(t, documentTestUUID, responseDocument.Uuid.String(), "Response uuid does not match")
}

type UploadResponse struct {
	DocumentUUID uuid.UUID `json:"documentUUID"`
}

func uploadDocument(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	ctr, err := testutil.CreateTestContainerPostgres(ctx, "BasicSetup", dbUser, dbPassword)
	if err != nil {
		assert.FailNow(t, err.Error())
		return
	}

	connectionString, err := ctr.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		assert.FailNow(t, err.Error())
		return
	}

	dbHandle := postgres.DatabaseHandler{DbConfig: postgres.ConfigForDatabase{ConUrl: connectionString}}
	documentCtrl := &v1.DocumentController{DocumentRepository: postgres.NewDocumentRepository(dbHandle)}
	router := v1.SetupRouter(documentCtrl, nil, nil)
	request := &v1.UploadRequest{DocumentBase64String: func() *string { v := "THIS IS A TEST DOCUMENT"; return &v }()}
	requestJSON, _ := json.Marshal(request)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, httptest.NewRequest(
		"POST",
		"/api/v1/documents/",
		strings.NewReader(string(requestJSON)),
	))

	response := UploadResponse{}
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		assert.FailNow(t, err.Error())
		return
	}

	assert.Equal(t, http.StatusOK, w.Code, "Response should be 200")
	assert.NotEqual(t, uuid.Nil, response.DocumentUUID)
}

type DeleteResponse struct {
	Success bool `json:"success"`
}

func deleteDocument(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	ctr, err := testutil.CreateTestContainerPostgres(ctx, "BasicSetupWithOneDocumentTableEntry", dbUser, dbPassword)
	if err != nil {
		assert.FailNow(t, err.Error())
		return
	}

	connectionString, err := ctr.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		assert.FailNow(t, err.Error())
		return
	}

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
