package integration

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	v1 "pdf_service_api/controller/v1"
	postgres2 "pdf_service_api/service/postgres"
	"pdf_service_api/testutil"
	"strings"
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
	uuids := []string{"a5fdea38-0a86-4c19-ae4f-c87a01bc860d", "335a6b95-6707-4e2b-9c37-c76d017f6f97"}
	bytes, err := json.Marshal(uuids)
	require.NoError(t, err)

	ctx := context.Background()
	ctr, err := testutil.CreateTestContainerPostgresWithInitFileName(ctx, dbUser, dbPassword, "OneDocumentTableEntryAndTwoSelections")
	require.NoError(t, err)
	defer testcontainers.TerminateContainer(ctr)

	connectionString, err := ctr.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	dbHandle := postgres2.DatabaseHandler{DbConfig: postgres2.ConfigForDatabase{ConUrl: connectionString}}
	selectionCtrl := &v1.SelectionController{SelectionRepository: postgres2.NewSelectionRepository(dbHandle)}
	extractCtrl := &v1.ExtractionController{SelectionRepository: postgres2.NewSelectionRepository(dbHandle)}
	router := v1.SetupRouter(nil, selectionCtrl, nil, extractCtrl)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, httptest.NewRequest(
		"POST",
		fmt.Sprintf("/api/v1/extract/basic"),
		strings.NewReader(string(bytes)),
	))

	fmt.Println(w.Body.String())
	assert.Equal(t, http.StatusOK, w.Result().StatusCode)
}
