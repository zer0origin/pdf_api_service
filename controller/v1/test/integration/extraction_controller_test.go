package integration

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	v1 "pdf_service_api/controller/v1"
	"pdf_service_api/service/dataapi"
	postgres2 "pdf_service_api/service/postgres"
	"pdf_service_api/testutil"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
)

func TestExtractionIntegration(t *testing.T) {
	t.Parallel()
	t.Run("Basic Endpoint functionality test", getTextFromSelectionUuidBase64NotIncluded)
}

func getTextFromSelectionUuidBase64NotIncluded(t *testing.T) {
	t.Parallel()
	selectionUids := []uuid.UUID{uuid.MustParse("a5fdea38-0a86-4c19-ae4f-c87a01bc860d"), uuid.MustParse("335a6b95-6707-4e2b-9c37-c76d017f6f97")}

	dataToSend := v1.ExtractUUIDsRequest{
		DocumentUid: uuid.MustParse("b66fd223-515f-4503-80cc-2bdaa50ef474"),
		OwnerUid:    uuid.MustParse("ea167a48-c1b3-46c4-911b-090e807132fc"),
		Uids:        selectionUids,
	}

	bytes, err := json.Marshal(dataToSend)
	require.NoError(t, err)

	ctx := context.Background()
	ctr, err := testutil.CreateTestContainerPostgresWithInitFileName(ctx, dbUser, dbPassword, "BreyerExample")
	require.NoError(t, err)
	defer testcontainers.TerminateContainer(ctr)

	//p, ctrTwo, err := testutil.CreateDataApiTestContainer()
	//require.NoError(t, err)
	//defer testcontainers.TerminateContainer(ctrTwo)

	connectionString, err := ctr.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	dbHandle := postgres2.DatabaseHandler{DbConfig: postgres2.ConfigForDatabase{ConUrl: connectionString}}
	dataApi := dataapi.DataService{BaseUrl: fmt.Sprintf("http://localhost:%d", 8080)}

	extractCtrl := &v1.ExtractionController{
		SelectionRepository: postgres2.NewSelectionRepository(dbHandle),
		DocumentRepository:  postgres2.NewDocumentRepository(dbHandle),
		DataService:         dataApi,
		Options:             v1.ExtractionOptions{GetBase64IfNotIncluded: true}}

	router := v1.SetupRouter(nil, nil, nil, extractCtrl)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, httptest.NewRequest(
		"POST",
		fmt.Sprintf("/api/v1/extract/basic"),
		strings.NewReader(string(bytes)),
	))

	fmt.Println(w.Body.String())
	assert.Equal(t, http.StatusOK, w.Result().StatusCode)
}
