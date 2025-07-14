package integration

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"pdf_service_api/repositories"
	"pdf_service_api/testutil"
	v1 "pdf_service_api/v1"
	"pdf_service_api/v1/controller"
	"testing"
)

func TestSelectionsIntegration(t *testing.T) {
	t.Run("getSelections", getSelections)
	t.Run("deleteSelections", deleteSelections)
	t.Run("deleteSelectionsUuidDoesNotExist", deleteSelectionsUuidDoesNotExist)
}

// /api/v1/documents/:id/selections/
func getSelections(t *testing.T) {
	t.Parallel()

	expectedJsonResponse := `[{"selectionUUID":"a5fdea38-0a86-4c19-ae4f-c87a01bc860d","documentID":"b66fd223-515f-4503-80cc-2bdaa50ef474","selection_bounds":"{}"},{"selectionUUID":"335a6b95-6707-4e2b-9c37-c76d017f6f97","documentID":"b66fd223-515f-4503-80cc-2bdaa50ef474","selection_bounds":"{}"}]`
	ctx := context.Background()
	ctr, err := testutil.CreateTestContainerPostgres(ctx, "BasicSetupWithOneDocumentTableEntryAndTwoSelections", dbUser, dbPassword)
	if err != nil {
		assert.FailNow(t, err.Error())
	}
	t.Cleanup(testutil.CleanUp(ctx, *ctr))

	dbConfig, err := testutil.CreateDbConfig(ctx, *ctr)
	if err != nil {
		assert.FailNow(t, err.Error())
	}

	selectionController := &controller.SelectionController{SelectionRepository: repositories.NewSelectionRepository(dbConfig)}
	repo := repositories.NewDocumentRepository(dbConfig)
	documentController := &controller.DocumentController{DocumentRepository: repo, SelectionController: selectionController}
	router := v1.SetupRouter(documentController)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, httptest.NewRequest(
		"GET",
		fmt.Sprintf("/api/v1/documents/%s/selections/", "b66fd223-515f-4503-80cc-2bdaa50ef474"),
		nil,
	))

	fmt.Println(w.Body.String())
	assert.Equal(t, http.StatusOK, w.Result().StatusCode)
	assert.NotNil(t, w.Body.String())
	assert.NotContains(t, w.Body.String(), "Error")
	assert.Equal(t, expectedJsonResponse, w.Body.String())
}

// /api/v1/documents/:id/selections/
func deleteSelections(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	ctr, err := testutil.CreateTestContainerPostgres(ctx, "BasicSetupWithOneDocumentTableEntryAndTwoSelections", dbUser, dbPassword)
	if err != nil {
		assert.FailNow(t, err.Error())
	}
	t.Cleanup(testutil.CleanUp(ctx, *ctr))

	dbConfig, err := testutil.CreateDbConfig(ctx, *ctr)
	if err != nil {
		assert.FailNow(t, err.Error())
	}

	selectionController := &controller.SelectionController{SelectionRepository: repositories.NewSelectionRepository(dbConfig)}
	repo := repositories.NewDocumentRepository(dbConfig)
	documentController := &controller.DocumentController{DocumentRepository: repo, SelectionController: selectionController}
	router := v1.SetupRouter(documentController)

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
	}
	t.Cleanup(testutil.CleanUp(ctx, *ctr))

	dbConfig, err := testutil.CreateDbConfig(ctx, *ctr)
	if err != nil {
		assert.FailNow(t, err.Error())
	}

	selectionController := &controller.SelectionController{SelectionRepository: repositories.NewSelectionRepository(dbConfig)}
	repo := repositories.NewDocumentRepository(dbConfig)
	documentController := &controller.DocumentController{DocumentRepository: repo, SelectionController: selectionController}
	router := v1.SetupRouter(documentController)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, httptest.NewRequest(
		"DELETE",
		fmt.Sprintf("/api/v1/selections/%s", uuid.New().String()),
		nil,
	))

	fmt.Println(w.Body.String())
	assert.Equal(t, http.StatusOK, w.Result().StatusCode)
}
