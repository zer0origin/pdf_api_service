package integration

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
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

func TestSelectionsIntegration(t *testing.T) {
	t.Run("Get selection from a present document uuid", getSelectionsFromPresentDocumentUUID)
	t.Run("Get selection from a nonexistent document uuid", getSelectionsFromInvalidDocumentUUID)
	t.Run("Get selection from a present selection uuid", getSelectionFromPresentSelectionUUID)
	t.Run("Get selection from a nonexistent selection uuid", getSelectionsFromNonExistentDocumentUUID)
	t.Run("Delete selections by selection uuid", deleteSelectionsBySelectionUUID)
	t.Run("Delete selections by present document uuid", deleteSelectionsByDocumentUUID)
	t.Run("Delete selections by nonexistent selection uuid", deleteDelectionByNonexistentSelectionUUID)
	t.Run("Create new selection", createNewSelection)
}

func getSelectionFromPresentSelectionUUID(t *testing.T) {
	testDocumentUuidString := "a5fdea38-0a86-4c19-ae4f-c87a01bc860d"
	expectedJsonResponse := `{"selections":[{"selectionUUID":"a5fdea38-0a86-4c19-ae4f-c87a01bc860d","documentUUID":"b66fd223-515f-4503-80cc-2bdaa50ef474"}]}`
	t.Parallel()

	ctx := context.Background()
	ctr, err := testutil.CreateTestContainerPostgresWithInitFileName(ctx, dbUser, dbPassword, "OneDocumentTableEntryAndTwoSelections")
	require.NoError(t, err)

	connectionString, err := ctr.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	dbHandle := postgres.DatabaseHandler{DbConfig: postgres.ConfigForDatabase{ConUrl: connectionString}}
	selectionCtrl := &v1.SelectionController{SelectionRepository: postgres.NewSelectionRepository(dbHandle)}
	router := v1.SetupRouter(nil, selectionCtrl, nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, httptest.NewRequest(
		"GET",
		fmt.Sprintf("/api/v1/selections/?selectionUUID=%s", testDocumentUuidString),
		nil,
	))

	fmt.Println(w.Body.String())
	assert.Equal(t, http.StatusOK, w.Result().StatusCode)
	assert.NotNil(t, w.Body.String())
	assert.NotContains(t, w.Body.String(), "Error")
	assert.Equal(t, expectedJsonResponse, w.Body.String(), "Body does not match expected output.")
}

func getSelectionsFromPresentDocumentUUID(t *testing.T) {
	testDocumentUuidString := "b66fd223-515f-4503-80cc-2bdaa50ef474"
	expectedJsonResponse := `{"selections":[{"selectionUUID":"a5fdea38-0a86-4c19-ae4f-c87a01bc860d","documentUUID":"b66fd223-515f-4503-80cc-2bdaa50ef474"},{"selectionUUID":"335a6b95-6707-4e2b-9c37-c76d017f6f97","documentUUID":"b66fd223-515f-4503-80cc-2bdaa50ef474"}]}`
	t.Parallel()

	ctx := context.Background()
	ctr, err := testutil.CreateTestContainerPostgresWithInitFileName(ctx, dbUser, dbPassword, "OneDocumentTableEntryAndTwoSelections")
	require.NoError(t, err)

	connectionString, err := ctr.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	dbHandle := postgres.DatabaseHandler{DbConfig: postgres.ConfigForDatabase{ConUrl: connectionString}}
	selectionCtrl := &v1.SelectionController{SelectionRepository: postgres.NewSelectionRepository(dbHandle)}
	router := v1.SetupRouter(nil, selectionCtrl, nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, httptest.NewRequest(
		"GET",
		fmt.Sprintf("/api/v1/selections/?documentUUID=%s", testDocumentUuidString),
		nil,
	))

	assert.Equal(t, http.StatusOK, w.Result().StatusCode)
	assert.NotNil(t, w.Body.String())
	assert.NotContains(t, w.Body.String(), "Error")
	assert.Equal(t, expectedJsonResponse, w.Body.String(), "Body does not match expected output.")
}

func getSelectionsFromNonExistentDocumentUUID(t *testing.T) {
	testDocumentUuidString := uuid.Nil.String()
	expectedJsonResponse := `{"selections":[]}`
	t.Parallel()

	ctx := context.Background()
	ctr, err := testutil.CreateTestContainerPostgresWithInitFileName(ctx, dbUser, dbPassword, "OneDocumentTableEntryAndTwoSelections")
	require.NoError(t, err)

	connectionString, err := ctr.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	dbHandle := postgres.DatabaseHandler{DbConfig: postgres.ConfigForDatabase{ConUrl: connectionString}}
	selectionCtrl := &v1.SelectionController{SelectionRepository: postgres.NewSelectionRepository(dbHandle)}
	router := v1.SetupRouter(nil, selectionCtrl, nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, httptest.NewRequest(
		"GET",
		fmt.Sprintf("/api/v1/selections/?documentUUID=%s", testDocumentUuidString),
		nil,
	))

	assert.Equal(t, http.StatusOK, w.Result().StatusCode)
	assert.NotNil(t, w.Body.String())
	assert.NotContains(t, w.Body.String(), "Error")
	assert.Equal(t, expectedJsonResponse, w.Body.String(), "Body does not match expected output.")
}

func getSelectionsFromInvalidDocumentUUID(t *testing.T) {
	testDocumentUuidString := uuid.New().String()
	expectedJsonResponse := `{"selections":[]}`
	t.Parallel()

	ctx := context.Background()
	ctr, err := testutil.CreateTestContainerPostgresWithInitFileName(ctx, dbUser, dbPassword, "OneDocumentTableEntryAndTwoSelections")
	require.NoError(t, err)

	connectionString, err := ctr.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	dbHandle := postgres.DatabaseHandler{DbConfig: postgres.ConfigForDatabase{ConUrl: connectionString}}
	selectionCtrl := &v1.SelectionController{SelectionRepository: postgres.NewSelectionRepository(dbHandle)}
	router := v1.SetupRouter(nil, selectionCtrl, nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, httptest.NewRequest(
		"GET",
		fmt.Sprintf("/api/v1/selections/?documentUUID=%s", testDocumentUuidString),
		nil,
	))

	fmt.Println(w.Body.String())
	assert.Equal(t, http.StatusOK, w.Result().StatusCode)
	assert.NotNil(t, w.Body.String())
	assert.NotContains(t, w.Body.String(), "Error")
	assert.Equal(t, expectedJsonResponse, w.Body.String(), "Body does not match expected output.")
}

func deleteSelectionsBySelectionUUID(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	ctr, err := testutil.CreateTestContainerPostgresWithInitFileName(ctx, dbUser, dbPassword, "OneDocumentTableEntryAndTwoSelections")
	require.NoError(t, err)

	connectionString, err := ctr.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	dbHandle := postgres.DatabaseHandler{DbConfig: postgres.ConfigForDatabase{ConUrl: connectionString}}
	selectionCtrl := &v1.SelectionController{SelectionRepository: postgres.NewSelectionRepository(dbHandle)}
	router := v1.SetupRouter(nil, selectionCtrl, nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, httptest.NewRequest(
		"DELETE",
		fmt.Sprintf("/api/v1/selections/?selectionUUID=%s", "a5fdea38-0a86-4c19-ae4f-c87a01bc860d"),
		nil,
	))

	fmt.Println(w.Body.String())
	assert.Equal(t, http.StatusOK, w.Result().StatusCode)
}

func deleteSelectionsByDocumentUUID(t *testing.T) {
	documentTestUUID := "b66fd223-515f-4503-80cc-2bdaa50ef474"
	t.Parallel()

	ctx := context.Background()
	ctr, err := testutil.CreateTestContainerPostgresWithInitFileName(ctx, dbUser, dbPassword, "OneDocumentTableEntryAndTwoSelections")
	require.NoError(t, err)

	connectionString, err := ctr.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	dbHandle := postgres.DatabaseHandler{DbConfig: postgres.ConfigForDatabase{ConUrl: connectionString}}
	selectionCtrl := &v1.SelectionController{SelectionRepository: postgres.NewSelectionRepository(dbHandle)}
	router := v1.SetupRouter(nil, selectionCtrl, nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, httptest.NewRequest(
		"DELETE",
		fmt.Sprintf("/api/v1/selections/?documentUUID=%s", documentTestUUID),
		nil,
	))

	assert.Equal(t, http.StatusOK, w.Result().StatusCode)

	err = dbHandle.WithConnection(func(db *sql.DB) error {
		rows := db.QueryRow(`SELECT * FROM selection_table WHERE "Document_UUID"=$1`, documentTestUUID)
		err := rows.Scan()
		if err != nil {
			return err
		}

		return nil
	})

	if !errors.Is(err, sql.ErrNoRows) {
		assert.FailNow(t, err.Error())
	}
}

func deleteDelectionByNonexistentSelectionUUID(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	ctr, err := testutil.CreateTestContainerPostgresWithInitFileName(ctx, dbUser, dbPassword, "OneDocumentTableEntryAndTwoSelections")
	require.NoError(t, err)

	connectionString, err := ctr.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	dbHandle := postgres.DatabaseHandler{DbConfig: postgres.ConfigForDatabase{ConUrl: connectionString}}
	selectionCtrl := &v1.SelectionController{SelectionRepository: postgres.NewSelectionRepository(dbHandle)}
	router := v1.SetupRouter(nil, selectionCtrl, nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, httptest.NewRequest(
		"DELETE",
		fmt.Sprintf("/api/v1/selections/?selectionUUID=%s", uuid.New().String()),
		nil,
	))

	fmt.Println(w.Body.String())
	assert.Equal(t, http.StatusOK, w.Result().StatusCode)
}

func createNewSelection(t *testing.T) {
	documentTestUUID := "b66fd223-515f-4503-80cc-2bdaa50ef474"
	t.Parallel()

	ctx := context.Background()
	ctr, err := testutil.CreateTestContainerPostgresWithInitFileName(ctx, dbUser, dbPassword, "OneDocumentTableEntryAndTwoSelections")
	require.NoError(t, err)

	connectionString, err := ctr.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	dbHandle := postgres.DatabaseHandler{DbConfig: postgres.ConfigForDatabase{ConUrl: connectionString}}
	selectionCtrl := &v1.SelectionController{SelectionRepository: postgres.NewSelectionRepository(dbHandle)}
	router := v1.SetupRouter(nil, selectionCtrl, nil)

	request := &v1.AddNewSelectionRequest{
		DocumentUUID:    func() *uuid.UUID { v := uuid.MustParse(documentTestUUID); return &v }(),
		IsComplete:      false,
		Settings:        nil,
		SelectionBounds: nil,
	}

	requestJSON, _ := json.Marshal(request)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, httptest.NewRequest(
		"POST",
		"/api/v1/selections/",
		strings.NewReader(string(requestJSON)),
	))

	assert.Equal(t, http.StatusOK, w.Result().StatusCode)
	assert.NotContains(t, w.Body.String(), "Error")
}
