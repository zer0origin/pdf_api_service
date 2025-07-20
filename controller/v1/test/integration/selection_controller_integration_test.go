package integration

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	v1 "pdf_service_api/controller/v1"
	"pdf_service_api/postgres"
	"pdf_service_api/testutil"
	"strings"
	"testing"
)

func TestSelectionsIntegration(t *testing.T) {
	t.Run("getSelections", getSelections)
	t.Run("deleteSelections", deleteSelections)
	t.Run("deleteSelectionsUuidDoesNotExist", deleteSelectionsUuidDoesNotExist)
	t.Run("createNewSelection", createNewSelection)
}

func getSelections(t *testing.T) {
	testDocumentUuidString := "b66fd223-515f-4503-80cc-2bdaa50ef474"
	expectedJsonResponse := `[{"selectionUUID":"a5fdea38-0a86-4c19-ae4f-c87a01bc860d","documentID":"b66fd223-515f-4503-80cc-2bdaa50ef474"},{"selectionUUID":"335a6b95-6707-4e2b-9c37-c76d017f6f97","documentID":"b66fd223-515f-4503-80cc-2bdaa50ef474"}]`
	t.Parallel()

	ctx := context.Background()
	ctr, err := testutil.CreateTestContainerPostgres(ctx, "BasicSetupWithOneDocumentTableEntryAndTwoSelections", dbUser, dbPassword)
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

// /api/v1/documents/:id/selections/
func deleteSelections(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	ctr, err := testutil.CreateTestContainerPostgres(ctx, "BasicSetupWithOneDocumentTableEntryAndTwoSelections", dbUser, dbPassword)
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
	selectionCtrl := &v1.SelectionController{SelectionRepository: postgres.NewSelectionRepository(dbHandle)}
	router := v1.SetupRouter(nil, selectionCtrl, nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, httptest.NewRequest(
		"DELETE",
		fmt.Sprintf("/api/v1/selections/%s", "a5fdea38-0a86-4c19-ae4f-c87a01bc860d"),
		nil,
	))

	fmt.Println(w.Body.String())
	assert.Equal(t, http.StatusOK, w.Result().StatusCode)
}

// /api/v1/documents/:id/selections/
func deleteSelectionsUuidDoesNotExist(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	ctr, err := testutil.CreateTestContainerPostgres(ctx, "BasicSetupWithOneDocumentTableEntryAndTwoSelections", dbUser, dbPassword)
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
	selectionCtrl := &v1.SelectionController{SelectionRepository: postgres.NewSelectionRepository(dbHandle)}
	router := v1.SetupRouter(nil, selectionCtrl, nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, httptest.NewRequest(
		"DELETE",
		fmt.Sprintf("/api/v1/selections/%s", uuid.New().String()),
		nil,
	))

	fmt.Println(w.Body.String())
	assert.Equal(t, http.StatusOK, w.Result().StatusCode)
}

func createNewSelection(t *testing.T) {
	documentTestUUID := "b66fd223-515f-4503-80cc-2bdaa50ef474"
	t.Parallel()

	ctx := context.Background()
	ctr, err := testutil.CreateTestContainerPostgres(ctx, "BasicSetupWithOneDocumentTableEntryAndTwoSelections", dbUser, dbPassword)
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
	selectionCtrl := &v1.SelectionController{SelectionRepository: postgres.NewSelectionRepository(dbHandle)}
	router := v1.SetupRouter(nil, selectionCtrl, nil)

	request := &v1.AddNewSelectionRequest{
		DocumentID:      func() *uuid.UUID { v := uuid.MustParse(documentTestUUID); return &v }(),
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
