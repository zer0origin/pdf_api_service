package unit

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"pdf_service_api/models"
	"pdf_service_api/repositories/mock"
	v1 "pdf_service_api/v1"
	"pdf_service_api/v1/controller"
	"strings"
	"testing"
)

type UploadResponse struct {
	DocumentUUID uuid.UUID `json:"documentUUID"`
}

func TestDocumentControllerUnit(t *testing.T) {
	t.Run("pingRouter", pingRouter)
	t.Run("uploadDocument", uploadDocument)
	t.Run("getDocument", getDocument)
}

func pingRouter(t *testing.T) {
	repo := &mock.MapRepository{Repo: make(map[uuid.UUID]models.Document)}
	documentController := &controller.DocumentController{DocumentRepository: repo}
	router := v1.SetupRouter(documentController)

	w := httptest.NewRecorder() //creates a recorder that records its mutations for later inspection in tests.
	req, _ := http.NewRequest("GET", "/ping", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "{\"message\":\"pong\"}", w.Body.String())
}

func uploadDocument(t *testing.T) {
	repo := &mock.MapRepository{Repo: make(map[uuid.UUID]models.Document)}
	documentController := &controller.DocumentController{DocumentRepository: repo}
	router := v1.SetupRouter(documentController)

	data := models.UploadRequest{
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

func getDocument(t *testing.T) {
	repo := &mock.MapRepository{Repo: make(map[uuid.UUID]models.Document)}
	documentController := &controller.DocumentController{DocumentRepository: repo}
	router := v1.SetupRouter(documentController)

	ExampleUUID := uuid.New()
	ExampleDocument := models.Document{
		Uuid:          ExampleUUID,
		PdfBase64:     func() *string { v := "TEMP DOCUMENT"; return &v }(),
		SelectionData: nil,
	}
	repo.Repo[ExampleUUID] = ExampleDocument

	request := &models.GetDocumentRequest{DocumentUuid: ExampleUUID}
	requestJSON, _ := json.Marshal(request)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, httptest.NewRequest(
		"GET",
		"/api/v1/documents/"+ExampleUUID.String(),
		strings.NewReader(string(requestJSON)),
	))

	responseDocument := &models.Document{}
	err := json.NewDecoder(w.Body).Decode(responseDocument)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, w.Code, "Response should be 200")
	assert.Equal(t, responseDocument.Uuid, ExampleDocument.Uuid, "Response uuid does not match")
	assert.Equal(t, responseDocument.PdfBase64, ExampleDocument.PdfBase64, "Response PdfBase64 does not match")
	assert.Equal(t, responseDocument.SelectionData, ExampleDocument.SelectionData, "Response SelectionData does not match")
}
