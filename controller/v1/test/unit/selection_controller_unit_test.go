package unit

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	v1 "pdf_service_api/controller/v1"
	"pdf_service_api/controller/v1/test/unit/mock"
	"pdf_service_api/models"
	"strings"
	"testing"
)

func TestSelectionUnit(t *testing.T) {
	t.Run("SelectionBoundsParsing", SelectionBoundsParsing)
	t.Run("NoSelectionBounds", NoSelectionBounds)
}

type addSelectionResponse struct {
	SelectionUUID uuid.UUID `json:"selectionUUID"`
}

func SelectionBoundsParsing(t *testing.T) {
	documentTestUUID := "b66fd223-515f-4503-80cc-2bdaa50ef474"

	mm := make(map[int][]models.SelectionBounds)
	mm[0] = make([]models.SelectionBounds, 2)
	mm[1] = make([]models.SelectionBounds, 2)

	toCreate := models.Selection{
		DocumentUUID:    func() *uuid.UUID { v := uuid.MustParse(documentTestUUID); return &v }(),
		IsComplete:      false,
		Settings:        nil,
		SelectionBounds: &mm,
	}

	mm[0][0] = models.SelectionBounds{
		SelectionMethod: nil,
		X1:              22,
		X2:              65,
		Y1:              24,
		Y2:              87,
	}

	mm[0][1] = models.SelectionBounds{
		SelectionMethod: nil,
		X1:              73,
		X2:              47,
		Y1:              28,
		Y2:              65,
	}

	mm[1][0] = models.SelectionBounds{
		SelectionMethod: nil,
		X1:              93,
		X2:              34,
		Y1:              16,
		Y2:              64,
	}

	mm[1][1] = models.SelectionBounds{
		SelectionMethod: nil,
		X1:              83,
		X2:              27,
		Y1:              36,
		Y2:              86,
	}

	jsonByteData, err := json.Marshal(toCreate)
	require.NoError(t, err)

	documentRepo := &mock.EmptyDocumentRepository{}
	selectionRepo := &mock.MapSelectionRepository{Repo: make(map[uuid.UUID]models.Selection)}
	selectionCtrl := &v1.SelectionController{SelectionRepository: selectionRepo}
	documentController := &v1.DocumentController{DocumentRepository: documentRepo}
	router := v1.SetupRouter(documentController, selectionCtrl, nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, httptest.NewRequest(
		"POST",
		"/api/v1/selections/",
		strings.NewReader(string(jsonByteData)),
	))

	assert.Equal(t, http.StatusOK, w.Result().StatusCode)

	response := addSelectionResponse{}
	err = json.NewDecoder(w.Body).Decode(&response)
	require.NoError(t, err)

	storedData := selectionRepo.Repo[response.SelectionUUID]
	assert.NotEqual(t, uuid.Nil, storedData.Uuid)
	assert.NotNil(t, storedData.SelectionBounds)
}

func NoSelectionBounds(t *testing.T) {
	documentTestUUID := "b66fd223-515f-4503-80cc-2bdaa50ef474"

	toCreate := models.Selection{
		Uuid:            uuid.New(),
		DocumentUUID:    func() *uuid.UUID { v := uuid.MustParse(documentTestUUID); return &v }(),
		IsComplete:      false,
		Settings:        nil,
		SelectionBounds: nil,
	}

	jsonByteData, err := json.Marshal(toCreate)
	require.NoError(t, err)

	documentRepo := &mock.EmptyDocumentRepository{}
	selectionRepo := &mock.MapSelectionRepository{Repo: make(map[uuid.UUID]models.Selection)}
	selectionCtrl := &v1.SelectionController{SelectionRepository: selectionRepo}
	documentController := &v1.DocumentController{DocumentRepository: documentRepo}
	router := v1.SetupRouter(documentController, selectionCtrl, nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, httptest.NewRequest(
		"POST",
		"/api/v1/selections/",
		strings.NewReader(string(jsonByteData)),
	))

	assert.Equal(t, http.StatusOK, w.Result().StatusCode)

	response := addSelectionResponse{}
	err = json.NewDecoder(w.Body).Decode(&response)
	require.NoError(t, err)

	storedData := selectionRepo.Repo[response.SelectionUUID]
	assert.NotEqual(t, uuid.Nil, storedData.Uuid)
	assert.Nil(t, storedData.SelectionBounds)
}
