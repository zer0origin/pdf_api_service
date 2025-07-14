package unit

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	v1 "pdf_service_api/controller/v1"
	"pdf_service_api/controller/v1/test/unit/mock"
	"pdf_service_api/domain"
	"strings"
	"testing"
)

func TestSelectionUnit(t *testing.T) {
	t.Run("TestSelectionBoundsParsing", TestSelectionBoundsParsing)
}

func TestSelectionBoundsParsing(t *testing.T) {
	documentTestUUID := "b66fd223-515f-4503-80cc-2bdaa50ef474"

	mm := make(map[int][]domain.SelectionBounds)
	mm[0] = make([]domain.SelectionBounds, 2)
	mm[1] = make([]domain.SelectionBounds, 2)

	toCreate := domain.Selection{
		Uuid:            uuid.New(),
		DocumentID:      func() *uuid.UUID { v := uuid.MustParse(documentTestUUID); return &v }(),
		IsComplete:      false,
		Settings:        nil,
		SelectionBounds: &mm,
	}

	mm[0][0] = domain.SelectionBounds{
		SelectionMethod: nil,
		X1:              22,
		X2:              65,
		Y1:              24,
		Y2:              87,
	}

	mm[0][1] = domain.SelectionBounds{
		SelectionMethod: nil,
		X1:              73,
		X2:              47,
		Y1:              28,
		Y2:              65,
	}

	mm[1][0] = domain.SelectionBounds{
		SelectionMethod: nil,
		X1:              93,
		X2:              34,
		Y1:              16,
		Y2:              64,
	}

	mm[1][1] = domain.SelectionBounds{
		SelectionMethod: nil,
		X1:              83,
		X2:              27,
		Y1:              36,
		Y2:              86,
	}

	jsonByteData, err := json.Marshal(toCreate)
	if err != nil {
		assert.FailNow(t, err.Error())
	}

	documentRepo := &mock.MapDocumentRepository{Repo: make(map[uuid.UUID]domain.Document)}
	selectionRepo := &mock.MapSelectionRepository{Repo: make(map[uuid.UUID]domain.Selection)}
	selectionCtrl := &v1.SelectionController{SelectionRepository: selectionRepo}
	documentController := &v1.DocumentController{DocumentRepository: documentRepo, SelectionController: selectionCtrl}
	router := v1.SetupRouter(documentController)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, httptest.NewRequest(
		"POST",
		"/api/v1/selections/",
		strings.NewReader(string(jsonByteData)),
	))

	assert.Equal(t, http.StatusOK, w.Result().StatusCode)
	assert.NotNil(t, selectionRepo.Repo[toCreate.Uuid])
}
