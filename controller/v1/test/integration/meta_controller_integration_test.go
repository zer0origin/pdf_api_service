package integration

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	v1 "pdf_service_api/controller/v1"
	"pdf_service_api/models"
	"pdf_service_api/postgres"
	"pdf_service_api/testutil"
	"strings"
	"testing"
)

func TestMetaIntegration(t *testing.T) {
	t.Run("get meta using a present uuid", getMetaPresentUUID)
	t.Run("update meta using a present uuid", updateMetaPresentUUID)
}

func getMetaPresentUUID(t *testing.T) {
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
		"/api/v1/meta/?metaUUID="+testUUID,
		nil,
	))

	assert.Equal(t, http.StatusOK, w.Result().StatusCode)
	assert.Equal(t, w.Body.String(), expected)
}

func updateMetaPresentUUID(t *testing.T) {
	testUUID := "b66fd223-515f-4503-80cc-2bdaa50ef474"
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

	newData := models.Meta{
		UUID:          uuid.MustParse(testUUID),
		NumberOfPages: func() *uint32 { v := new(uint32); *v = 43; return v }(),
		Height:        func() *float32 { v := new(float32); *v = 36.2; return v }(),
		Width:         func() *float32 { v := new(float32); *v = 29.5; return v }(),
		Images:        nil,
	}

	requestJSON, err := json.Marshal(newData)
	if err != nil {
		assert.FailNow(t, err.Error())
		return
	}

	w := httptest.NewRecorder()
	router.ServeHTTP(w, httptest.NewRequest(
		"PUT",
		"/api/v1/meta/?metaUUID="+testUUID,
		strings.NewReader(string(requestJSON)),
	))

	err = dbHandle.WithConnection(func(db *sql.DB) error {
		sqlStatement := `SELECT "Document_UUID", "Number_Of_Pages", "Height", "Width" FROM documentmeta_table WHERE "Document_UUID" = $1`
		row := db.QueryRow(sqlStatement, testUUID)

		var (
			uid     string
			noPages int32
			height  float32
			width   float32
		)

		err := row.Scan(&uid, &noPages, &height, &width)
		if err != nil {
			return err
		}

		assert.Equal(t, newData.UUID.String(), uid)
		assert.EqualValues(t, *newData.NumberOfPages, noPages)
		assert.EqualValues(t, *newData.Height, height)
		assert.EqualValues(t, *newData.Width, width)

		return nil
	})
	if err != nil {
		return
	}

	assert.Equal(t, http.StatusOK, w.Result().StatusCode)
}
