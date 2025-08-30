package integration

import (
	"context"
	"database/sql"
	_ "embed"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"net/http"
	"net/http/httptest"
	v1 "pdf_service_api/controller/v1"
	"pdf_service_api/models"
	"pdf_service_api/service/dataapi"
	pg "pdf_service_api/service/postgres"
	"pdf_service_api/testutil"
	"strings"
	"testing"
)

func TestMetaIntegration(t *testing.T) {
	t.Parallel()
	t.Run("get meta using a present uuid", getMetaPresentUUID)
	t.Run("get meta using a uuid not present in table", getMetaUUIDDoesNotExistInTable)
	t.Run("update meta using a present uuid", updateMetaPresentUUID)
	t.Run("update meta using a present uuid with new images", updateImageMetaPresentUUID)
	t.Run("Add meta using the DataApi to generate the meta data", addMetaBase64Included)
	t.Run("Add meta using the DataApi to generate the meta data", addMetaBase64Excluded)
}

func getMetaPresentUUID(t *testing.T) {
	t.Parallel()
	mm := make(map[uint32]string)
	expectedObj := models.Meta{
		DocumentUUID:  uuid.MustParse("b66fd223-515f-4503-80cc-2bdaa50ef474"),
		NumberOfPages: func() *uint32 { v := uint32(31); return &v }(),
		Height:        func() *float32 { v := float32(1920); return &v }(),
		Width:         func() *float32 { v := float32(1080); return &v }(),
		Images:        &mm,
	}
	bytes, err := json.Marshal(expectedObj)
	require.NoError(t, err)

	ctx := context.Background()
	ctr, err := testutil.CreateTestContainerPostgresWithInitFileName(ctx, dbUser, dbPassword, "OneDocumentTableEntryTwoSelectionsAndMetaData")
	require.NoError(t, err)
	defer testcontainers.TerminateContainer(ctr)

	connectionString, err := ctr.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	dbHandle := pg.DatabaseHandler{DbConfig: pg.ConfigForDatabase{ConUrl: connectionString}}
	metaCtrl := &v1.MetaController{MetaRepository: pg.NewMetaRepository(dbHandle)}
	router := v1.SetupRouter(nil, nil, metaCtrl)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, httptest.NewRequest(
		"GET",
		"/api/v1/meta/?documentUUID="+expectedObj.DocumentUUID.String()+"&ownerUUID=f701aa7e-10e9-48b9-83f1-6b035a5b7564",
		nil,
	))

	assert.Equal(t, http.StatusOK, w.Result().StatusCode)
	assert.Equal(t, string(bytes), w.Body.String())
}

func getMetaUUIDDoesNotExistInTable(t *testing.T) {
	t.Parallel()
	expectedObj := "{\"error\":\"data not found\"}"

	ctx := context.Background()
	ctr, err := testutil.CreateTestContainerPostgresWithInitFileName(ctx, dbUser, dbPassword, "OneDocumentTableEntryTwoSelectionsAndMetaData")
	require.NoError(t, err)
	defer testcontainers.TerminateContainer(ctr)

	connectionString, err := ctr.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	dbHandle := pg.DatabaseHandler{DbConfig: pg.ConfigForDatabase{ConUrl: connectionString}}
	metaCtrl := &v1.MetaController{MetaRepository: pg.NewMetaRepository(dbHandle)}
	router := v1.SetupRouter(nil, nil, metaCtrl)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, httptest.NewRequest(
		"GET",
		"/api/v1/meta/?documentUUID="+"b66fa223-515f-4503-80cc-2bdaa50ef474"+"&ownerUUID=f701aa7e-10e9-48b9-83f1-6b035a5b7564",
		nil,
	))

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, expectedObj, w.Body.String())
}

func updateMetaPresentUUID(t *testing.T) {
	t.Parallel()
	testUUID := "b66fd223-515f-4503-80cc-2bdaa50ef474"

	ctx := context.Background()
	ctr, err := testutil.CreateTestContainerPostgresWithInitFileName(ctx, dbUser, dbPassword, "OneDocumentTableEntryTwoSelectionsAndMetaData")
	require.NoError(t, err)
	defer testcontainers.TerminateContainer(ctr)

	connectionString, err := ctr.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	dbHandle := pg.DatabaseHandler{DbConfig: pg.ConfigForDatabase{ConUrl: connectionString}}
	metaCtrl := &v1.MetaController{MetaRepository: pg.NewMetaRepository(dbHandle)}
	router := v1.SetupRouter(nil, nil, metaCtrl)

	newData := models.Meta{
		DocumentUUID:  uuid.MustParse(testUUID),
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
		"/api/v1/meta/?documentUUID="+testUUID,
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

		assert.Equal(t, newData.DocumentUUID.String(), uid)
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
	defer testcontainers.TerminateContainer(ctr)

	connectionString, err := ctr.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	dbHandle := pg.DatabaseHandler{DbConfig: pg.ConfigForDatabase{ConUrl: connectionString}}
	metaCtrl := &v1.MetaController{MetaRepository: pg.NewMetaRepository(dbHandle)}
	router := v1.SetupRouter(nil, nil, metaCtrl)

	strArr := make(map[uint32]string, 0)
	strArr[0] = "Image0"
	strArr[1] = "Image1"
	strArr[2] = "Image2"

	newData := models.Meta{
		DocumentUUID:  uuid.MustParse(testUUID),
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
		"/api/v1/meta/?documentUUID="+testUUID,
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

		assert.Equal(t, newData.DocumentUUID.String(), uid)
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

func addMetaBase64Included(t *testing.T) {
	if *testutil.SkipDataApiIntegrationTest {
		t.Skip("Skipping test due to flags")
	}

	documentUUID := "b66fd223-515f-4503-80cc-2bdaa50ef474"
	ownerUUID := "ea167a48-c1b3-46c4-911b-090e807132fc"
	request := v1.AddMetaRequest{
		DocumentUUID:         uuid.MustParse(documentUUID),
		OwnerUUID:            uuid.MustParse(ownerUUID),
		DocumentBase64String: func() *string { return &testutil.HundredPagesPdfInBase64 }(),
		OwnerType:            1,
	}
	requestBytes, err := json.Marshal(request)
	require.NoError(t, err)

	p, ctr, err := testutil.CreateDataApiTestContainer()
	require.NoError(t, err)
	defer testcontainers.TerminateContainer(ctr)

	ctx := context.Background()
	pgCtr, err := testutil.CreateTestContainerPostgresWithInitFileName(ctx, dbUser, dbPassword, "OneDocumentTableEntryWithRealDocument")
	require.NoError(t, err)
	defer testcontainers.TerminateContainer(pgCtr)

	dbHandle, err := testutil.CreateDatabaseHandlerFromPostgresInfo(ctx, *pgCtr)
	require.NoError(t, err)

	srv := dataapi.DataService{BaseUrl: fmt.Sprintf("http://localhost:%d", p.Int())}
	metaCtrl := &v1.MetaController{DataService: srv, MetaRepository: pg.NewMetaRepository(dbHandle)}
	router := v1.SetupRouter(nil, nil, metaCtrl)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, httptest.NewRequest(
		"POST",
		"/api/v1/meta/",
		strings.NewReader(string(requestBytes)),
	))

	assert.Equal(t, http.StatusOK, w.Code)
	if w.Code != http.StatusOK {
		fmt.Println(w.Body.String())
	}

	dbHandle.WithConnection(func(db *sql.DB) error {
		row := db.QueryRow(`SELECT "Number_Of_Pages" FROM postgres.public.documentmeta_table WHERE "Document_UUID" = $1`, request.DocumentUUID)

		var number int8
		err := row.Scan(&number)
		require.NoError(t, err)
		assert.EqualValues(t, number, 101)
		return nil
	})
}

func addMetaBase64Excluded(t *testing.T) {
	if *testutil.SkipDataApiIntegrationTest {
		t.Skip("Skipping test due to flags")
	}

	documentUUID := "b66fd223-515f-4503-80cc-2bdaa50ef474"
	ownerUUID := "ea167a48-c1b3-46c4-911b-090e807132fc"
	request := v1.AddMetaRequest{
		DocumentUUID: uuid.MustParse(documentUUID),
		OwnerUUID:    uuid.MustParse(ownerUUID),
		OwnerType:    1,
	}
	requestBytes, err := json.Marshal(request)
	require.NoError(t, err)

	p, ctr, err := testutil.CreateDataApiTestContainer()
	require.NoError(t, err)
	defer testcontainers.TerminateContainer(ctr)

	ctx := context.Background()
	pgCtr, err := testutil.CreateTestContainerPostgresWithInitFileName(ctx, dbUser, dbPassword, "OneDocumentTableEntryWithRealDocument")
	require.NoError(t, err)
	defer testcontainers.TerminateContainer(pgCtr)

	dbHandle, err := testutil.CreateDatabaseHandlerFromPostgresInfo(ctx, *pgCtr)
	require.NoError(t, err)

	srv := dataapi.DataService{BaseUrl: fmt.Sprintf("http://localhost:%d", p.Int())}
	metaCtrl := &v1.MetaController{DataService: srv, MetaRepository: pg.NewMetaRepository(dbHandle), DocumentRepository: pg.NewDocumentRepository(dbHandle)}
	router := v1.SetupRouter(nil, nil, metaCtrl)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, httptest.NewRequest(
		"POST",
		"/api/v1/meta/",
		strings.NewReader(string(requestBytes)),
	))

	assert.Equal(t, http.StatusOK, w.Code)
	if w.Code != http.StatusOK {
		fmt.Println(w.Body.String())
	}

	dbHandle.WithConnection(func(db *sql.DB) error {
		row := db.QueryRow(`SELECT "Number_Of_Pages" FROM postgres.public.documentmeta_table WHERE "Document_UUID" = $1`, request.DocumentUUID)

		var number int8
		err := row.Scan(&number)
		require.NoError(t, err)
		assert.EqualValues(t, number, 101)
		return nil
	})
}
