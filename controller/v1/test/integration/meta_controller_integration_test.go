package integration

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	t.Parallel()
	t.Run("get meta using a present uuid", getMetaPresentUUID)
	t.Run("update meta using a present uuid", updateMetaPresentUUID)
	t.Run("update meta using a present uuid with new images", updateImageMetaPresentUUID)
}

func getMetaPresentUUID(t *testing.T) {
	t.Parallel()
	expectedObj := models.Meta{
		UUID:          uuid.MustParse("b66fd223-515f-4503-80cc-2bdaa50ef474"),
		NumberOfPages: func() *uint32 { v := uint32(31); return &v }(),
		Height:        func() *float32 { v := float32(1920); return &v }(),
		Width:         func() *float32 { v := float32(1080); return &v }(),
	}
	bytes, err := json.Marshal(expectedObj)
	require.NoError(t, err)

	ctx := context.Background()
	ctr, err := testutil.CreateTestContainerPostgresWithInitFileName(ctx, dbUser, dbPassword, "OneDocumentTableEntryTwoSelectionsAndMetaData")
	require.NoError(t, err)

	connectionString, err := ctr.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	dbHandle := postgres.DatabaseHandler{DbConfig: postgres.ConfigForDatabase{ConUrl: connectionString}}
	metaCtrl := &v1.MetaController{MetaRepository: postgres.NewMetaRepository(dbHandle)}
	router := v1.SetupRouter(nil, nil, metaCtrl)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, httptest.NewRequest(
		"GET",
		"/api/v1/meta/?metaUUID="+expectedObj.UUID.String(),
		nil,
	))

	assert.Equal(t, http.StatusOK, w.Result().StatusCode)
	assert.Equal(t, w.Body.String(), string(bytes))
}

func updateMetaPresentUUID(t *testing.T) {
	t.Parallel()
	testUUID := "b66fd223-515f-4503-80cc-2bdaa50ef474"

	ctx := context.Background()
	ctr, err := testutil.CreateTestContainerPostgresWithInitFileName(ctx, dbUser, dbPassword, "OneDocumentTableEntryTwoSelectionsAndMetaData")
	require.NoError(t, err)

	connectionString, err := ctr.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

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
	require.NoError(t, err)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, httptest.NewRequest(
		"PUT",
		"/api/v1/meta/?metaUUID="+testUUID,
		strings.NewReader(string(requestJSON)),
	))

	assert.Equal(t, http.StatusOK, w.Result().StatusCode)
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
}

func updateImageMetaPresentUUID(t *testing.T) {
	t.Parallel()
	testUUID := "b66fd223-515f-4503-80cc-2bdaa50ef474"
	expectedStr := `{"0":"Image0","1":"Image1","2":"Image2"}`

	ctx := context.Background()
	ctr, err := testutil.CreateTestContainerPostgresWithInitFileName(ctx, dbUser, dbPassword, "OneDocumentTableEntryTwoSelectionsAndMetaData")
	require.NoError(t, err)

	connectionString, err := ctr.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	dbHandle := postgres.DatabaseHandler{DbConfig: postgres.ConfigForDatabase{ConUrl: connectionString}}
	metaCtrl := &v1.MetaController{MetaRepository: postgres.NewMetaRepository(dbHandle)}
	router := v1.SetupRouter(nil, nil, metaCtrl)

	strArr := make(map[uint32]string, 0)
	strArr[0] = "Image0"
	strArr[1] = "Image1"
	strArr[2] = "Image2"

	newData := models.Meta{
		UUID:          uuid.MustParse(testUUID),
		NumberOfPages: nil,
		Height:        nil,
		Width:         nil,
		Images:        &strArr,
	}

	requestJSON, err := json.Marshal(newData)
	require.NoError(t, err)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, httptest.NewRequest(
		"PUT",
		"/api/v1/meta/?metaUUID="+testUUID,
		strings.NewReader(string(requestJSON)),
	))

	assert.Equal(t, http.StatusOK, w.Result().StatusCode)
	err = dbHandle.WithConnection(func(db *sql.DB) error {
		sqlStatement := `SELECT "Document_UUID", "Number_Of_Pages", "Height", "Width", "Images" FROM documentmeta_table WHERE "Document_UUID" = $1`
		row := db.QueryRow(sqlStatement, testUUID)

		var (
			uid     string
			noPages int32
			height  float32
			width   float32
			images  string
		)

		err := row.Scan(&uid, &noPages, &height, &width, &images)
		if err != nil {
			return err
		}

		assert.Equal(t, newData.UUID.String(), uid)
		assert.EqualValues(t, noPages, 31)
		assert.NotNil(t, noPages)
		assert.NotNil(t, height)
		assert.EqualValues(t, height, 1920)
		assert.NotNil(t, width)
		assert.EqualValues(t, width, 1080)
		assert.Equal(t, images, expectedStr)

		return nil
	})
	if err != nil {
		return
	}

}
