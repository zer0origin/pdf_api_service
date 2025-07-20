package integration

import (
	"context"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	v1 "pdf_service_api/controller/v1"
	"pdf_service_api/postgres"
	"pdf_service_api/testutil"
	"testing"
)

func TestMetaIntegration(t *testing.T) {
	t.Run("getMetaDataFromDatabase", getMetaDataFromDatabase)
}

func getMetaDataFromDatabase(t *testing.T) {
	testUUID := "b66fd223-515f-4503-80cc-2bdaa50ef474"
	expected := "{\"UUID\":\"b66fd223-515f-4503-80cc-2bdaa50ef474\",\"NumberOfPages\":31,\"Height\":1920,\"Width\":1080,\"Images\":null}"
	t.Parallel()

	ctx := context.Background()
	ctr, err := testutil.CreateTestContainerPostgres(ctx, "BasicSetupWithOneDocumentTableEntryTwoSelectionsAndMetaData", dbUser, dbPassword)
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
	metaCtrl := &v1.MetaController{MetaRepository: postgres.NewMetaRepository(dbHandle)}
	router := v1.SetupRouter(nil, nil, metaCtrl)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, httptest.NewRequest(
		"GET",
		"/api/v1/meta/?id="+testUUID,
		nil,
	))

	assert.Equal(t, http.StatusOK, w.Result().StatusCode)
	assert.Equal(t, w.Body.String(), expected)
}
