package integration

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	v1 "pdf_service_api/controller/v1"
	postgres2 "pdf_service_api/service/postgres"
	"pdf_service_api/testutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
)

func TestExtractionIntegration(t *testing.T) {
	t.Parallel()
	t.Run("Basic Endpoint functionality test", getTextFromSelectionUUID)
}

func getTextFromSelectionUUID(t *testing.T) {
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
	router := v1.SetupRouter(nil, selectionCtrl, nil, nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, httptest.NewRequest(
		"GET",
		fmt.Sprintf("/api/v1/selections/?selectionUUID=%s", testDocumentUuidString),
		nil,
	))

	fmt.Println(w.Body.String())
	assert.Equal(t, http.StatusOK, w.Result().StatusCode)
	assert.NotNil(t, w.Body.String())
	assert.NotContains(t, w.Body.String(), "error")
	assert.Equal(t, expectedJsonResponse, w.Body.String(), "Body does not match expected output.")
}
