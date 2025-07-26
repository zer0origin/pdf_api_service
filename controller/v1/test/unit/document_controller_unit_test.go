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

type UploadResponse struct {
	DocumentUUID uuid.UUID `json:"documentUUID"`
}

func TestDocumentControllerUnit(t *testing.T) {
	t.Run("pingRouter", pingRouter)
	t.Run("uploadDocument", uploadDocument)
	t.Run("getDocument", getDocument)
}

func pingRouter(t *testing.T) {
	repo := &mock.MapDocumentRepository{Repo: make(map[uuid.UUID]domain.Document)}
	documentController := &v1.DocumentController{DocumentRepository: repo}
	router := v1.SetupRouter(documentController, nil, nil)

	w := httptest.NewRecorder() //creates a recorder that records its mutations for later inspection in tests.
	req, _ := http.NewRequest("GET", "/ping", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "{\"message\":\"pong\"}", w.Body.String())
}

func uploadDocument(t *testing.T) {
	repo := &mock.MapDocumentRepository{Repo: make(map[uuid.UUID]domain.Document)}
	documentController := &v1.DocumentController{DocumentRepository: repo}
	router := v1.SetupRouter(documentController, nil, nil)

	data := v1.UploadRequest{
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
	repo := &mock.MapDocumentRepository{Repo: make(map[uuid.UUID]domain.Document)}
	documentController := &v1.DocumentController{DocumentRepository: repo}
	router := v1.SetupRouter(documentController, nil, nil)

	ExampleUUID := uuid.New()
	ExampleDocument := domain.Document{
		Uuid:          ExampleUUID,
		PdfBase64:     func() *string { v := "TEMP DOCUMENT"; return &v }(),
		SelectionData: nil,
	}
	repo.Repo[ExampleUUID] = ExampleDocument

	request := &v1.GetDocumentRequest{DocumentUuid: ExampleUUID}
	requestJSON, _ := json.Marshal(request)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, httptest.NewRequest(
		"GET",
		"/api/v1/documents/?documentUUID="+ExampleUUID.String(),
		strings.NewReader(string(requestJSON)),
	))

	responseDocument := &domain.Document{}
	err := json.NewDecoder(w.Body).Decode(responseDocument)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, w.Code, "Response should be 200")
	assert.Equal(t, responseDocument.Uuid, ExampleDocument.Uuid, "Response uuid does not match")
	assert.Equal(t, responseDocument.PdfBase64, ExampleDocument.PdfBase64, "Response PdfBase64 does not match")
	assert.Equal(t, responseDocument.SelectionData, ExampleDocument.SelectionData, "Response SelectionData does not match")
}
