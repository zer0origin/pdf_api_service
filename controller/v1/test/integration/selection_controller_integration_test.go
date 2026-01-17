package integration

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	v1 "pdf_service_api/controller/v1"
	"pdf_service_api/models"
	postgres2 "pdf_service_api/service/postgres"
	"pdf_service_api/testutil"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
)

func TestSelectionsIntegration(t *testing.T) {
	t.Parallel()
	t.Run("Get selection from a present document uuid", getSelectionsFromPresentDocumentUUID)
	t.Run("Get selection from a nonexistent document uuid", getSelectionsFromInvalidDocumentUUID)
	t.Run("Get selection from a present selection uuid", getSelectionFromPresentSelectionUUID)
	t.Run("Get selection from a nonexistent selection uuid", getSelectionsFromNonExistentDocumentUUID)
	t.Run("Delete selections by selection uuid", deleteSelectionsBySelectionUUID)
	t.Run("Delete selections by present document uuid", deleteSelectionsByDocumentUUID)
	t.Run("Delete selections by nonexistent selection uuid", deleteDelectionByNonexistentSelectionUUID)
	t.Run("Create new selection", createNewSelection)
	t.Run("Create lots of new selection", CreateNewSelectionWithCoordinatesBulk)
	t.Run("Checks that the route fails correctly, providing the correct information", CreateNewSelectionWithCoordinatesBulkFailure)
	t.Run("Create new selection that includes a page key", CreateNewSelectionWithPageKey)
	t.Run("Create new selection that includes some coordinates", CreateNewSelectionWithCoordinates)
}

func getSelectionFromPresentSelectionUUID(t *testing.T) {
	t.Parallel()
	testDocumentUuidString := "a5fdea38-0a86-4c19-ae4f-c87a01bc860d"
	expectedJsonResponse := `{"selections":[{"selectionUUID":"a5fdea38-0a86-4c19-ae4f-c87a01bc860d","documentUUID":"b66fd223-515f-4503-80cc-2bdaa50ef474"}]}`

	ctx := context.Background()
	ctr, err := testutil.CreateTestContainerPostgresWithInitFileName(ctx, dbUser, dbPassword, "OneDocumentTableEntryAndTwoSelections")
	require.NoError(t, err)
	defer testcontainers.TerminateContainer(ctr)

	connectionString, err := ctr.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	dbHandle := postgres2.DatabaseHandler{DbConfig: postgres2.ConfigForDatabase{ConUrl: connectionString}}
	selectionCtrl := &v1.SelectionController{SelectionRepository: postgres2.NewSelectionRepository(dbHandle)}
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
	t.Parallel()
	testDocumentUuidString := "b66fd223-515f-4503-80cc-2bdaa50ef474"
	expectedJsonResponse := `{"selections":[{"selectionUUID":"a5fdea38-0a86-4c19-ae4f-c87a01bc860d","documentUUID":"b66fd223-515f-4503-80cc-2bdaa50ef474"},{"selectionUUID":"335a6b95-6707-4e2b-9c37-c76d017f6f97","documentUUID":"b66fd223-515f-4503-80cc-2bdaa50ef474"}]}`

	ctx := context.Background()
	ctr, err := testutil.CreateTestContainerPostgresWithInitFileName(ctx, dbUser, dbPassword, "OneDocumentTableEntryAndTwoSelections")
	require.NoError(t, err)
	defer testcontainers.TerminateContainer(ctr)

	connectionString, err := ctr.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	dbHandle := postgres2.DatabaseHandler{DbConfig: postgres2.ConfigForDatabase{ConUrl: connectionString}}
	selectionCtrl := &v1.SelectionController{SelectionRepository: postgres2.NewSelectionRepository(dbHandle)}
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
	t.Parallel()
	testDocumentUuidString := uuid.Nil.String()
	expectedJsonResponse := `{"selections":[]}`

	ctx := context.Background()
	ctr, err := testutil.CreateTestContainerPostgresWithInitFileName(ctx, dbUser, dbPassword, "OneDocumentTableEntryAndTwoSelections")
	require.NoError(t, err)
	defer testcontainers.TerminateContainer(ctr)

	connectionString, err := ctr.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	dbHandle := postgres2.DatabaseHandler{DbConfig: postgres2.ConfigForDatabase{ConUrl: connectionString}}
	selectionCtrl := &v1.SelectionController{SelectionRepository: postgres2.NewSelectionRepository(dbHandle)}
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
	t.Parallel()
	testDocumentUuidString := uuid.New().String()
	expectedJsonResponse := `{"selections":[]}`

	ctx := context.Background()
	ctr, err := testutil.CreateTestContainerPostgresWithInitFileName(ctx, dbUser, dbPassword, "OneDocumentTableEntryAndTwoSelections")
	require.NoError(t, err)
	defer testcontainers.TerminateContainer(ctr)

	connectionString, err := ctr.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	dbHandle := postgres2.DatabaseHandler{DbConfig: postgres2.ConfigForDatabase{ConUrl: connectionString}}
	selectionCtrl := &v1.SelectionController{SelectionRepository: postgres2.NewSelectionRepository(dbHandle)}
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
	defer testcontainers.TerminateContainer(ctr)

	connectionString, err := ctr.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	dbHandle := postgres2.DatabaseHandler{DbConfig: postgres2.ConfigForDatabase{ConUrl: connectionString}}
	selectionCtrl := &v1.SelectionController{SelectionRepository: postgres2.NewSelectionRepository(dbHandle)}
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
	t.Parallel()
	documentTestUUID := "b66fd223-515f-4503-80cc-2bdaa50ef474"

	ctx := context.Background()
	ctr, err := testutil.CreateTestContainerPostgresWithInitFileName(ctx, dbUser, dbPassword, "OneDocumentTableEntryAndTwoSelections")
	require.NoError(t, err)
	defer testcontainers.TerminateContainer(ctr)

	connectionString, err := ctr.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	dbHandle := postgres2.DatabaseHandler{DbConfig: postgres2.ConfigForDatabase{ConUrl: connectionString}}
	selectionCtrl := &v1.SelectionController{SelectionRepository: postgres2.NewSelectionRepository(dbHandle)}
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
	defer testcontainers.TerminateContainer(ctr)

	connectionString, err := ctr.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	dbHandle := postgres2.DatabaseHandler{DbConfig: postgres2.ConfigForDatabase{ConUrl: connectionString}}
	selectionCtrl := &v1.SelectionController{SelectionRepository: postgres2.NewSelectionRepository(dbHandle)}
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
	t.Parallel()
	documentTestUUID := "b66fd223-515f-4503-80cc-2bdaa50ef474"

	ctx := context.Background()
	ctr, err := testutil.CreateTestContainerPostgresWithInitFileName(ctx, dbUser, dbPassword, "OneDocumentTableEntryAndTwoSelections")
	require.NoError(t, err)
	defer testcontainers.TerminateContainer(ctr)

	connectionString, err := ctr.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	dbHandle := postgres2.DatabaseHandler{DbConfig: postgres2.ConfigForDatabase{ConUrl: connectionString}}
	selectionCtrl := &v1.SelectionController{SelectionRepository: postgres2.NewSelectionRepository(dbHandle)}
	router := v1.SetupRouter(nil, selectionCtrl, nil)

	request := &v1.AddNewSelectionRequest{
		DocumentUUID: func() *uuid.UUID { v := uuid.MustParse(documentTestUUID); return &v }(),
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

func CreateNewSelectionWithPageKey(t *testing.T) {
	t.Parallel()
	documentTestUUID := "b66fd223-515f-4503-80cc-2bdaa50ef474"

	ctx := context.Background()
	ctr, err := testutil.CreateTestContainerPostgresWithInitFileName(ctx, dbUser, dbPassword, "OneDocumentTableEntryAndTwoSelections")
	require.NoError(t, err)
	defer testcontainers.TerminateContainer(ctr)

	connectionString, err := ctr.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	dbHandle := postgres2.DatabaseHandler{DbConfig: postgres2.ConfigForDatabase{ConUrl: connectionString}}
	selectionCtrl := &v1.SelectionController{SelectionRepository: postgres2.NewSelectionRepository(dbHandle)}
	router := v1.SetupRouter(nil, selectionCtrl, nil)

	request := &v1.AddNewSelectionRequest{
		DocumentUUID: func() *uuid.UUID { v := uuid.MustParse(documentTestUUID); return &v }(),
		PageKey:      "TestPage",
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

func CreateNewSelectionWithCoordinates(t *testing.T) {
	t.Parallel()
	documentTestUUID := "b66fd223-515f-4503-80cc-2bdaa50ef474"

	ctx := context.Background()
	ctr, err := testutil.CreateTestContainerPostgresWithInitFileName(ctx, dbUser, dbPassword, "OneDocumentTableEntry")
	require.NoError(t, err)
	defer testcontainers.TerminateContainer(ctr)

	connectionString, err := ctr.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	dbHandle := postgres2.DatabaseHandler{DbConfig: postgres2.ConfigForDatabase{ConUrl: connectionString}}
	selectionCtrl := &v1.SelectionController{SelectionRepository: postgres2.NewSelectionRepository(dbHandle)}
	router := v1.SetupRouter(nil, selectionCtrl, nil)

	request := &v1.AddNewSelectionRequest{
		DocumentUUID: func() *uuid.UUID { v := uuid.MustParse(documentTestUUID); return &v }(),
		Coordinates: &models.Coordinates{
			X1: 43.122,
			Y1: 52.125,
			X2: 13,
			Y2: 27.853,
		},
	}

	requestJSON, _ := json.Marshal(request)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, httptest.NewRequest(
		"POST",
		"/api/v1/selections/",
		strings.NewReader(string(requestJSON)),
	))

	storedCoordinates := &models.Coordinates{}
	err = dbHandle.WithConnection(func(db *sql.DB) error {
		data := ""
		rows := db.QueryRow(`SELECT "Coordinates" FROM selection_table WHERE "Document_UUID"=$1 AND "Coordinates" IS NOT NULL `, documentTestUUID)
		err := rows.Scan(&data)
		require.NoError(t, err)

		err = json.Unmarshal([]byte(data), storedCoordinates)
		require.NoError(t, err)

		return nil
	})

	assert.Equal(t, request.Coordinates.X1, storedCoordinates.X1)
	assert.Equal(t, request.Coordinates.X2, storedCoordinates.X2)
	assert.Equal(t, request.Coordinates.Y1, storedCoordinates.Y1)
	assert.Equal(t, request.Coordinates.Y2, storedCoordinates.Y2)
}

func CreateNewSelectionWithCoordinatesBulk(t *testing.T) {
	t.Parallel()
	documentTestUUID := "b66fd223-515f-4503-80cc-2bdaa50ef474"

	ctx := context.Background()
	ctr, err := testutil.CreateTestContainerPostgresWithInitFileName(ctx, dbUser, dbPassword, "OneDocumentTableEntry")
	require.NoError(t, err)
	defer testcontainers.TerminateContainer(ctr)

	connectionString, err := ctr.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	dbHandle := postgres2.DatabaseHandler{DbConfig: postgres2.ConfigForDatabase{ConUrl: connectionString}}
	selectionCtrl := &v1.SelectionController{SelectionRepository: postgres2.NewSelectionRepository(dbHandle)}
	router := v1.SetupRouter(nil, selectionCtrl, nil)

	request := make([]v1.AddNewSelectionRequest, 2)
	request[0] = v1.AddNewSelectionRequest{
		DocumentUUID: func() *uuid.UUID { v := uuid.MustParse(documentTestUUID); return &v }(),
		Coordinates: &models.Coordinates{
			X1: 43.122,
			Y1: 52.125,
			X2: 13,
			Y2: 27.853,
		},
	}

	request[1] = v1.AddNewSelectionRequest{
		DocumentUUID: func() *uuid.UUID { v := uuid.MustParse(documentTestUUID); return &v }(),
		Coordinates: &models.Coordinates{
			X1: 73,
			Y1: 76,
			X2: 65,
			Y2: 34,
		},
	}

	requestJSON, _ := json.Marshal(request)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, httptest.NewRequest(
		"POST",
		"/api/v1/selections/bulk",
		strings.NewReader(string(requestJSON)),
	))
	assert.Equal(t, http.StatusOK, w.Result().StatusCode)

	_ = dbHandle.WithConnection(func(db *sql.DB) error {
		data := 0
		rows := db.QueryRow(`SELECT count("Coordinates") FROM selection_table WHERE "Document_UUID"=$1 AND "Coordinates" IS NOT NULL `, documentTestUUID)
		err := rows.Scan(&data)
		require.NoError(t, err)
		fmt.Println(data)
		assert.Equal(t, 2, data)
		return nil
	})
}

func CreateNewSelectionWithCoordinatesBulkFailure(t *testing.T) {
	t.Parallel()
	documentTestUUID := "1c705cd7-146e-4569-bd02-bde5f3a015c5"

	ctx := context.Background()
	ctr, err := testutil.CreateTestContainerPostgresWithInitFileName(ctx, dbUser, dbPassword, "OneDocumentTableEntry")
	require.NoError(t, err)
	defer testcontainers.TerminateContainer(ctr)

	connectionString, err := ctr.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	dbHandle := postgres2.DatabaseHandler{DbConfig: postgres2.ConfigForDatabase{ConUrl: connectionString}}
	selectionCtrl := &v1.SelectionController{SelectionRepository: postgres2.NewSelectionRepository(dbHandle)}
	router := v1.SetupRouter(nil, selectionCtrl, nil)

	request := make([]v1.AddNewSelectionRequest, 2)
	request[0] = v1.AddNewSelectionRequest{
		DocumentUUID: func() *uuid.UUID { v := uuid.MustParse(documentTestUUID); return &v }(),
		Coordinates: &models.Coordinates{
			X1: 43.122,
			Y1: 52.125,
			X2: 13,
			Y2: 27.853,
		},
	}

	request[1] = v1.AddNewSelectionRequest{
		DocumentUUID: func() *uuid.UUID { v := uuid.MustParse(documentTestUUID); return &v }(),
		Coordinates: &models.Coordinates{
			X1: 73,
			Y1: 76,
			X2: 65,
			Y2: 34,
		},
	}

	requestJSON, _ := json.Marshal(request)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, httptest.NewRequest(
		"POST",
		"/api/v1/selections/bulk",
		strings.NewReader(string(requestJSON)),
	))

	assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
	str := w.Body.String()
	assert.Contains(t, str, "error")
	fmt.Println(str)

	_ = dbHandle.WithConnection(func(db *sql.DB) error {
		data := 0
		rows := db.QueryRow(`SELECT count("Coordinates") FROM selection_table WHERE "Document_UUID"=$1 AND "Coordinates" IS NOT NULL `, documentTestUUID)
		err := rows.Scan(&data)
		require.NoError(t, err)
		fmt.Println(data)
		assert.Equal(t, 0, data)
		return nil
	})
}
