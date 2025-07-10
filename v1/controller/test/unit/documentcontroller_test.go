package unit

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"pdf_service_api/models"
	"pdf_service_api/models/requests"
	"pdf_service_api/repositories/mock"
	v1 "pdf_service_api/v1"
	"pdf_service_api/v1/controller"
	"strings"
	"testing"
)

type UploadResponse struct {
	DocumentUUID uuid.UUID `json:"documentUUID"`
}

func TestUploadDocument(t *testing.T) {
	t.Parallel()
	repo := &mock.MapRepository{Repo: make(map[uuid.UUID]models.Document)}
	documentController := controller.NewDocumentController(repo)
	router := v1.SetupRouter(documentController)

	data := requests.UploadRequest{
		DocumentBase64String: func() *string { v := "TEMP DOCUMENT"; return &v }(),
	}
	documentJSON, _ := json.Marshal(data)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, httptest.NewRequest(
		"POST",
		"/api/v1/documents/",
		strings.NewReader(string(documentJSON)),
	))

	responseUUID := UploadResponse{}
	err := json.NewDecoder(w.Body).Decode(&responseUUID)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, w.Code, "Response should be 200")
	assert.NotEqual(t, uuid.Nil, responseUUID.DocumentUUID)
}

func TestGetDocument(t *testing.T) {
	t.Parallel()
	repo := &mock.MapRepository{Repo: make(map[uuid.UUID]models.Document)}
	documentController := controller.NewDocumentController(repo)
	router := v1.SetupRouter(documentController)

	ExampleUUID := uuid.New()
	ExampleDocument := models.Document{
		Uuid:          ExampleUUID,
		PdfBase64:     func() *string { v := "TEMP DOCUMENT"; return &v }(),
		SelectionData: nil,
	}
	repo.Repo[ExampleUUID] = ExampleDocument

	request := &requests.GetDocumentRequest{DocumentUuid: ExampleUUID}
	requestJSON, _ := json.Marshal(request)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, httptest.NewRequest(
		"GET",
		"/api/v1/documents/"+ExampleUUID.String(),
		strings.NewReader(string(requestJSON)),
	))

	responseDocument := &models.Document{}
	json.NewDecoder(w.Body).Decode(responseDocument)

	assert.Equal(t, http.StatusOK, w.Code, "Response should be 200")
	assert.Equal(t, responseDocument.Uuid, ExampleDocument.Uuid, "Response uuid does not match")
	assert.Equal(t, responseDocument.PdfBase64, ExampleDocument.PdfBase64, "Response PdfBase64 does not match")
	assert.Equal(t, responseDocument.SelectionData, ExampleDocument.SelectionData, "Response SelectionData does not match")
}
