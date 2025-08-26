package integration

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	v1 "pdf_service_api/controller/v1"
	postgres2 "pdf_service_api/service/postgres"
	"pdf_service_api/testutil"
	"strings"
	"testing"
)

var dbUser = "user"
var dbPassword = "password"

func TestDocumentIntegration(t *testing.T) {
	t.Parallel()
	t.Run("Test database connection", databaseConnection)
	t.Run("Get document with present document uuid", getDocumentWithDocumentUUID)
	t.Run("Get document with present document uuid and no owner uuid", getDocumentWithNoOwnerUuid)
	t.Run("Get document with present document uuid with excludes param", getDocumentWithDocumentUUIDExcludeBase64)
	t.Run("Get document with present owner uuid", getDocumentWithOwnerUUID)
	t.Run("Get document with present owner uuid with limit 1 and offset 0 set", getDocumentWithOwnerUUIDWithLimit1AndOffset0)
	t.Run("Get document with present owner uuid with limit 1 and offset 1 set", getDocumentWithOwnerUUIDWithLimit1AndOffset1)
	t.Run("Get document with present owner uuid with limit 1 and offset 2 set", getDocumentWithOwnerUUIDWithLimit1AndOffset2)
	t.Run("Get document with present owner uuid with limit 1 and offset 10 set", getDocumentWithOwnerUUIDWithLimit1AndOffset10)
	t.Run("Get document with present owner uuid with limit 2 and offset 0 set", getDocumentWithOwnerUUIDWithLimit2AndOffset0)
	t.Run("Get document with present owner uuid with limit 2 and offset 1 set", getDocumentWithOwnerUUIDWithLimit2AndOffset1)
	t.Run("Get document with present owner uuid with excludes params", getDocumentWithOwnerUUIDWithExcludes)
	t.Run("Get document with nonexistent document uuid", getDocumentWithNonexistentDocumentUUID)
	t.Run("Upload a new document", uploadDocument)
	t.Run("Upload a new document with document title", uploadDocumentWithTitle)
	t.Run("Delete existing document", deleteDocument)
}

func databaseConnection(t *testing.T) {
	t.Parallel()
	wd, err := os.Getwd()
	rd, err := os.Executable()
	dir := filepath.Dir(rd)
	fmt.Printf("%s", wd)
	fmt.Printf("%s", rd)
	fmt.Printf("%s", dir)

	ctx := context.Background()
	ctr, err := testutil.CreateTestContainerPostgres(ctx, dbUser, dbPassword)
	defer testcontainers.TerminateContainer(ctr)
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

func getDocumentWithDocumentUUID(t *testing.T) {
	t.Parallel()
	documentTestUUID := uuid.MustParse("b66fd223-515f-4503-80cc-2bdaa50ef474")
	ownerTestUUID := uuid.MustParse("ea167a48-c1b3-46c4-911b-090e807132fc")
	expectedResponse := fmt.Sprintf(`{"documents":[{"documentUUID":"%s","documentTitle":"Fake Title","timeCreated":"2022-10-10T11:30:30Z","ownerUUID":"ea167a48-c1b3-46c4-911b-090e807132fc","ownerType":1,"base64":"Fake document for testing"}]}`, documentTestUUID)

	ctx := context.Background()
	ctr, err := testutil.CreateTestContainerPostgresWithInitFileName(ctx, dbUser, dbPassword, "OneDocumentTableEntry")
	require.NoError(t, err)
	defer testcontainers.TerminateContainer(ctr)

	connectionString, err := ctr.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	dbHandle := postgres2.DatabaseHandler{DbConfig: postgres2.ConfigForDatabase{ConUrl: connectionString}}
	documentCtrl := &v1.DocumentController{DocumentRepository: postgres2.NewDocumentRepository(dbHandle)}
	router := v1.SetupRouter(documentCtrl, nil, nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, httptest.NewRequest(
		"GET",
		fmt.Sprintf("/api/v1/documents/?documentUUID=%s&ownerUUID=%s", documentTestUUID.String(), ownerTestUUID.String()),
		nil,
	))

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, expectedResponse, w.Body.String())
}

func getDocumentWithDocumentUUIDExcludeBase64(t *testing.T) {
	t.Parallel()
	documentTestUUID := uuid.MustParse("b66fd223-515f-4503-80cc-2bdaa50ef474")
	ownerTestUUID := uuid.MustParse("ea167a48-c1b3-46c4-911b-090e807132fc")
	expectedResponse := fmt.Sprintf(`{"documents":[{"documentUUID":"%s","documentTitle":"Fake Title"}]}`, documentTestUUID)

	ctx := context.Background()
	ctr, err := testutil.CreateTestContainerPostgresWithInitFileName(ctx, dbUser, dbPassword, "OneDocumentTableEntry")
	require.NoError(t, err)
	defer testcontainers.TerminateContainer(ctr)

	connectionString, err := ctr.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	dbHandle := postgres2.DatabaseHandler{DbConfig: postgres2.ConfigForDatabase{ConUrl: connectionString}}
	documentCtrl := &v1.DocumentController{DocumentRepository: postgres2.NewDocumentRepository(dbHandle)}
	router := v1.SetupRouter(documentCtrl, nil, nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, httptest.NewRequest(
		"GET",
		fmt.Sprintf("/api/v1/documents/?exclude=timeCreated&exclude=ownerUUID&exclude=ownerType&exclude=base64&documentUUID=%s&ownerUUID=%s", documentTestUUID.String(), ownerTestUUID.String()),
		nil,
	))

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, expectedResponse, w.Body.String())
}

func getDocumentWithNoOwnerUuid(t *testing.T) {
	t.Parallel()
	documentTestUUID := uuid.MustParse(uuid.Nil.String())
	expectedResponse := "{\"error\":\"Required OwnerUuid is missing\"}"

	ctx := context.Background()
	ctr, err := testutil.CreateTestContainerPostgresWithInitFileName(ctx, dbUser, dbPassword, "OneDocumentTableEntry")
	require.NoError(t, err)
	defer testcontainers.TerminateContainer(ctr)

	connectionString, err := ctr.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	dbHandle := postgres2.DatabaseHandler{DbConfig: postgres2.ConfigForDatabase{ConUrl: connectionString}}
	documentCtrl := &v1.DocumentController{DocumentRepository: postgres2.NewDocumentRepository(dbHandle)}
	router := v1.SetupRouter(documentCtrl, nil, nil)

	request := &v1.GetDocumentRequest{DocumentUUID: &documentTestUUID}
	requestJSON, _ := json.Marshal(request)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, httptest.NewRequest(
		"GET",
		"/api/v1/documents/?documentUUID="+request.DocumentUUID.String(),
		strings.NewReader(string(requestJSON)),
	))

	fmt.Println(w.Body.String())
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, expectedResponse, w.Body.String())
}

func getDocumentWithNonexistentDocumentUUID(t *testing.T) {
	t.Parallel()
	documentTestUUID := uuid.MustParse(uuid.Nil.String())
	ownerTestUUID := uuid.MustParse("ea167a48-c1b3-46c4-911b-090e807132fc")
	expectedResponse := fmt.Sprintf(`{"error":"Document with documentUUID %s was not found."}`, documentTestUUID)

	ctx := context.Background()
	ctr, err := testutil.CreateTestContainerPostgresWithInitFileName(ctx, dbUser, dbPassword, "OneDocumentTableEntry")
	require.NoError(t, err)
	defer testcontainers.TerminateContainer(ctr)

	connectionString, err := ctr.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	dbHandle := postgres2.DatabaseHandler{DbConfig: postgres2.ConfigForDatabase{ConUrl: connectionString}}
	documentCtrl := &v1.DocumentController{DocumentRepository: postgres2.NewDocumentRepository(dbHandle)}
	router := v1.SetupRouter(documentCtrl, nil, nil)

	request := &v1.GetDocumentRequest{DocumentUUID: &documentTestUUID}
	requestJSON, _ := json.Marshal(request)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, httptest.NewRequest(
		"GET",
		"/api/v1/documents/?documentUUID="+request.DocumentUUID.String()+"&ownerUUID="+ownerTestUUID.String(),
		strings.NewReader(string(requestJSON)),
	))

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, expectedResponse, w.Body.String())
}

func getDocumentWithOwnerUUID(t *testing.T) {
	t.Parallel()
	ownerTestUUID := uuid.MustParse("4ce6af41-6cb5-4b02-a671-9fce16ea688d")
	expectedResponse := "{\"documents\":[{\"documentUUID\":\"b66fd223-515f-4503-80cc-2bdaa50ef474\",\"documentTitle\":\"Fake Title\",\"timeCreated\":\"2022-10-10T11:30:31Z\",\"ownerUUID\":\"4ce6af41-6cb5-4b02-a671-9fce16ea688d\",\"ownerType\":1,\"base64\":\"1\"},{\"documentUUID\":\"b5b7f18e-aed3-4eb7-aca8-79bcedf03d1b\",\"timeCreated\":\"2022-10-10T11:30:30Z\",\"ownerUUID\":\"4ce6af41-6cb5-4b02-a671-9fce16ea688d\",\"ownerType\":1,\"base64\":\"2\"},{\"documentUUID\":\"489fc81f-a087-457e-b8b4-ef9ad571d954\",\"timeCreated\":\"2022-10-10T11:30:29Z\",\"ownerUUID\":\"4ce6af41-6cb5-4b02-a671-9fce16ea688d\",\"ownerType\":1,\"base64\":\"3\"}]}"

	ctx := context.Background()
	ctr, err := testutil.CreateTestContainerPostgresWithInitFileName(ctx, dbUser, dbPassword, "UserTable")
	require.NoError(t, err)
	defer testcontainers.TerminateContainer(ctr)

	connectionString, err := ctr.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	dbHandle := postgres2.DatabaseHandler{DbConfig: postgres2.ConfigForDatabase{ConUrl: connectionString}}
	documentCtrl := &v1.DocumentController{DocumentRepository: postgres2.NewDocumentRepository(dbHandle)}
	router := v1.SetupRouter(documentCtrl, nil, nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, httptest.NewRequest(
		"GET",
		"/api/v1/documents/?ownerUUID="+ownerTestUUID.String(),
		nil,
	))

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, expectedResponse, w.Body.String())
}

func getDocumentWithOwnerUUIDWithLimit1AndOffset0(t *testing.T) {
	t.Parallel()
	ownerTestUUID := uuid.MustParse("4ce6af41-6cb5-4b02-a671-9fce16ea688d")
	expectedResponse := "{\"documents\":[{\"documentUUID\":\"b66fd223-515f-4503-80cc-2bdaa50ef474\",\"documentTitle\":\"Fake Title\",\"timeCreated\":\"2022-10-10T11:30:31Z\",\"ownerUUID\":\"4ce6af41-6cb5-4b02-a671-9fce16ea688d\",\"ownerType\":1,\"base64\":\"1\"}]}"

	ctx := context.Background()
	ctr, err := testutil.CreateTestContainerPostgresWithInitFileName(ctx, dbUser, dbPassword, "UserTable")
	require.NoError(t, err)
	defer testcontainers.TerminateContainer(ctr)

	connectionString, err := ctr.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	dbHandle := postgres2.DatabaseHandler{DbConfig: postgres2.ConfigForDatabase{ConUrl: connectionString}}
	documentCtrl := &v1.DocumentController{DocumentRepository: postgres2.NewDocumentRepository(dbHandle)}
	router := v1.SetupRouter(documentCtrl, nil, nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, httptest.NewRequest(
		"GET",
		"/api/v1/documents/?limit=1&offset=0&ownerUUID="+ownerTestUUID.String(),
		nil,
	))

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, expectedResponse, w.Body.String())
}

func getDocumentWithOwnerUUIDWithLimit1AndOffset1(t *testing.T) {
	t.Parallel()
	ownerTestUUID := uuid.MustParse("4ce6af41-6cb5-4b02-a671-9fce16ea688d")
	expectedResponse := "{\"documents\":[{\"documentUUID\":\"b5b7f18e-aed3-4eb7-aca8-79bcedf03d1b\",\"timeCreated\":\"2022-10-10T11:30:30Z\",\"ownerUUID\":\"4ce6af41-6cb5-4b02-a671-9fce16ea688d\",\"ownerType\":1,\"base64\":\"2\"}]}"

	ctx := context.Background()
	ctr, err := testutil.CreateTestContainerPostgresWithInitFileName(ctx, dbUser, dbPassword, "UserTable")
	require.NoError(t, err)
	defer testcontainers.TerminateContainer(ctr)

	connectionString, err := ctr.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	dbHandle := postgres2.DatabaseHandler{DbConfig: postgres2.ConfigForDatabase{ConUrl: connectionString}}
	documentCtrl := &v1.DocumentController{DocumentRepository: postgres2.NewDocumentRepository(dbHandle)}
	router := v1.SetupRouter(documentCtrl, nil, nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, httptest.NewRequest(
		"GET",
		"/api/v1/documents/?limit=1&offset=1&ownerUUID="+ownerTestUUID.String(),
		nil,
	))

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, expectedResponse, w.Body.String())
}

func getDocumentWithOwnerUUIDWithLimit1AndOffset2(t *testing.T) {
	t.Parallel()
	ownerTestUUID := uuid.MustParse("4ce6af41-6cb5-4b02-a671-9fce16ea688d")
	expectedResponse := "{\"documents\":[{\"documentUUID\":\"489fc81f-a087-457e-b8b4-ef9ad571d954\",\"timeCreated\":\"2022-10-10T11:30:29Z\",\"ownerUUID\":\"4ce6af41-6cb5-4b02-a671-9fce16ea688d\",\"ownerType\":1,\"base64\":\"3\"}]}"

	ctx := context.Background()
	ctr, err := testutil.CreateTestContainerPostgresWithInitFileName(ctx, dbUser, dbPassword, "UserTable")
	require.NoError(t, err)
	defer testcontainers.TerminateContainer(ctr)

	connectionString, err := ctr.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	dbHandle := postgres2.DatabaseHandler{DbConfig: postgres2.ConfigForDatabase{ConUrl: connectionString}}
	documentCtrl := &v1.DocumentController{DocumentRepository: postgres2.NewDocumentRepository(dbHandle)}
	router := v1.SetupRouter(documentCtrl, nil, nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, httptest.NewRequest(
		"GET",
		"/api/v1/documents/?limit=1&offset=2&ownerUUID="+ownerTestUUID.String(),
		nil,
	))

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, expectedResponse, w.Body.String())
}

func getDocumentWithOwnerUUIDWithLimit1AndOffset10(t *testing.T) {
	t.Parallel()
	ownerTestUUID := uuid.MustParse("4ce6af41-6cb5-4b02-a671-9fce16ea688d")
	expectedResponse := "{\"documents\":[]}"

	ctx := context.Background()
	ctr, err := testutil.CreateTestContainerPostgresWithInitFileName(ctx, dbUser, dbPassword, "UserTable")
	require.NoError(t, err)
	defer testcontainers.TerminateContainer(ctr)

	connectionString, err := ctr.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	dbHandle := postgres2.DatabaseHandler{DbConfig: postgres2.ConfigForDatabase{ConUrl: connectionString}}
	documentCtrl := &v1.DocumentController{DocumentRepository: postgres2.NewDocumentRepository(dbHandle)}
	router := v1.SetupRouter(documentCtrl, nil, nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, httptest.NewRequest(
		"GET",
		"/api/v1/documents/?limit=1&offset=10&ownerUUID="+ownerTestUUID.String(),
		nil,
	))

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, expectedResponse, w.Body.String())
}

func getDocumentWithOwnerUUIDWithLimit2AndOffset0(t *testing.T) {
	t.Parallel()
	ownerTestUUID := uuid.MustParse("4ce6af41-6cb5-4b02-a671-9fce16ea688d")
	expectedResponse := "{\"documents\":[{\"documentUUID\":\"b66fd223-515f-4503-80cc-2bdaa50ef474\",\"documentTitle\":\"Fake Title\",\"timeCreated\":\"2022-10-10T11:30:31Z\",\"ownerUUID\":\"4ce6af41-6cb5-4b02-a671-9fce16ea688d\",\"ownerType\":1,\"base64\":\"1\"},{\"documentUUID\":\"b5b7f18e-aed3-4eb7-aca8-79bcedf03d1b\",\"timeCreated\":\"2022-10-10T11:30:30Z\",\"ownerUUID\":\"4ce6af41-6cb5-4b02-a671-9fce16ea688d\",\"ownerType\":1,\"base64\":\"2\"}]}"

	ctx := context.Background()
	ctr, err := testutil.CreateTestContainerPostgresWithInitFileName(ctx, dbUser, dbPassword, "UserTable")
	require.NoError(t, err)
	defer testcontainers.TerminateContainer(ctr)

	connectionString, err := ctr.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	dbHandle := postgres2.DatabaseHandler{DbConfig: postgres2.ConfigForDatabase{ConUrl: connectionString}}
	documentCtrl := &v1.DocumentController{DocumentRepository: postgres2.NewDocumentRepository(dbHandle)}
	router := v1.SetupRouter(documentCtrl, nil, nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, httptest.NewRequest(
		"GET",
		"/api/v1/documents/?limit=2&offset=0&ownerUUID="+ownerTestUUID.String(),
		nil,
	))

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, expectedResponse, w.Body.String())
}

func getDocumentWithOwnerUUIDWithLimit2AndOffset1(t *testing.T) {
	t.Parallel()
	ownerTestUUID := uuid.MustParse("4ce6af41-6cb5-4b02-a671-9fce16ea688d")
	expectedResponse := "{\"documents\":[{\"documentUUID\":\"b5b7f18e-aed3-4eb7-aca8-79bcedf03d1b\",\"timeCreated\":\"2022-10-10T11:30:30Z\",\"ownerUUID\":\"4ce6af41-6cb5-4b02-a671-9fce16ea688d\",\"ownerType\":1,\"base64\":\"2\"},{\"documentUUID\":\"489fc81f-a087-457e-b8b4-ef9ad571d954\",\"timeCreated\":\"2022-10-10T11:30:29Z\",\"ownerUUID\":\"4ce6af41-6cb5-4b02-a671-9fce16ea688d\",\"ownerType\":1,\"base64\":\"3\"}]}"

	ctx := context.Background()
	ctr, err := testutil.CreateTestContainerPostgresWithInitFileName(ctx, dbUser, dbPassword, "UserTable")
	require.NoError(t, err)
	defer testcontainers.TerminateContainer(ctr)

	connectionString, err := ctr.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	dbHandle := postgres2.DatabaseHandler{DbConfig: postgres2.ConfigForDatabase{ConUrl: connectionString}}
	documentCtrl := &v1.DocumentController{DocumentRepository: postgres2.NewDocumentRepository(dbHandle)}
	router := v1.SetupRouter(documentCtrl, nil, nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, httptest.NewRequest(
		"GET",
		"/api/v1/documents/?limit=2&offset=1&ownerUUID="+ownerTestUUID.String(),
		nil,
	))

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, expectedResponse, w.Body.String())
}

func getDocumentWithOwnerUUIDWithExcludes(t *testing.T) {
	t.Parallel()
	ownerTestUUID := uuid.MustParse("4ce6af41-6cb5-4b02-a671-9fce16ea688d")
	expectedResponse := "{\"documents\":[{\"documentUUID\":\"b66fd223-515f-4503-80cc-2bdaa50ef474\",\"documentTitle\":\"Fake Title\",\"ownerUUID\":\"4ce6af41-6cb5-4b02-a671-9fce16ea688d\"},{\"documentUUID\":\"b5b7f18e-aed3-4eb7-aca8-79bcedf03d1b\",\"ownerUUID\":\"4ce6af41-6cb5-4b02-a671-9fce16ea688d\"},{\"documentUUID\":\"489fc81f-a087-457e-b8b4-ef9ad571d954\",\"ownerUUID\":\"4ce6af41-6cb5-4b02-a671-9fce16ea688d\"}]}"

	ctx := context.Background()
	ctr, err := testutil.CreateTestContainerPostgresWithInitFileName(ctx, dbUser, dbPassword, "UserTable")
	require.NoError(t, err)
	defer testcontainers.TerminateContainer(ctr)

	connectionString, err := ctr.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	dbHandle := postgres2.DatabaseHandler{DbConfig: postgres2.ConfigForDatabase{ConUrl: connectionString}}
	documentCtrl := &v1.DocumentController{DocumentRepository: postgres2.NewDocumentRepository(dbHandle)}
	router := v1.SetupRouter(documentCtrl, nil, nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, httptest.NewRequest(
		"GET",
		"/api/v1/documents/?exclude=base64&exclude=timeCreated&exclude=ownerType&ownerUUID="+ownerTestUUID.String(),
		nil,
	))

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, expectedResponse, w.Body.String())
}

type UploadResponse struct {
	DocumentUUID uuid.UUID `json:"documentUUID"`
}

func uploadDocument(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	ctr, err := testutil.CreateTestContainerPostgres(ctx, dbUser, dbPassword)
	require.NoError(t, err)
	defer testcontainers.TerminateContainer(ctr)

	connectionString, err := ctr.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	dbHandle := postgres2.DatabaseHandler{DbConfig: postgres2.ConfigForDatabase{ConUrl: connectionString}}
	documentCtrl := &v1.DocumentController{DocumentRepository: postgres2.NewDocumentRepository(dbHandle)}
	router := v1.SetupRouter(documentCtrl, nil, nil)
	request := &v1.CreateRequest{Base64: "THIS IS A TEST DOCUMENT"}
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

func uploadDocumentWithTitle(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	ctr, err := testutil.CreateTestContainerPostgres(ctx, dbUser, dbPassword)
	require.NoError(t, err)
	defer testcontainers.TerminateContainer(ctr)

	connectionString, err := ctr.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	dbHandle := postgres2.DatabaseHandler{DbConfig: postgres2.ConfigForDatabase{ConUrl: connectionString}}
	documentCtrl := &v1.DocumentController{DocumentRepository: postgres2.NewDocumentRepository(dbHandle)}
	router := v1.SetupRouter(documentCtrl, nil, nil)
	request := &v1.CreateRequest{DocumentTitle: func() *string { v := "Document Title"; return &v }(), Base64: "THIS IS A TEST DOCUMENT"}
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

	err = dbHandle.WithConnection(func(db *sql.DB) error {
		row := db.QueryRow(`SELECT 1 FROM document_table WHERE "Document_Title" = $1`, "Document Title")

		var number int8
		err := row.Scan(&number)
		return err
	})
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
	ctr, err := testutil.CreateTestContainerPostgresWithInitFileName(ctx, dbUser, dbPassword, "OneDocumentTableEntry")
	require.NoError(t, err)
	defer testcontainers.TerminateContainer(ctr)

	dbHandle, err := testutil.CreateDatabaseHandlerFromPostgresInfo(ctx, *ctr)
	require.NoError(t, err)

	documentCtrl := &v1.DocumentController{DocumentRepository: postgres2.NewDocumentRepository(dbHandle)}
	router := v1.SetupRouter(documentCtrl, nil, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, httptest.NewRequest(
		"DELETE",
		fmt.Sprintf("/api/v1/documents/?documentUUID=%s&ownerUUID=%s", "b66fd223-515f-4503-80cc-2bdaa50ef474", "ea167a48-c1b3-46c4-911b-090e807132fc"),
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
