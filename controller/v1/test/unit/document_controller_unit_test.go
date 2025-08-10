package unit

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	v1 "pdf_service_api/controller/v1"
	"pdf_service_api/controller/v1/test/unit/mock"
	"pdf_service_api/models"
	"strings"
	"testing"
)

type UploadResponse struct {
	DocumentUUID uuid.UUID `json:"documentUUID"`
}

func TestDocumentControllerUnit(t *testing.T) {
	t.Run("Ping router", pingRouter)
	t.Run("Upload a document", uploadDocument)
	t.Run("Get document from present uuid", getDocumentFromPresentUUID)
}

func pingRouter(t *testing.T) {
	repo := &mock.MapDocumentRepository{Repo: make(map[uuid.UUID]models.Document)}
	documentController := &v1.DocumentController{DocumentRepository: repo}
	router := v1.SetupRouter(documentController, nil, nil)

	w := httptest.NewRecorder() //creates a recorder that records its mutations for later inspection in tests.
	req, _ := http.NewRequest("GET", "/ping", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "{\"message\":\"pong\"}", w.Body.String())
}

func uploadDocument(t *testing.T) {
	repo := &mock.MapDocumentRepository{Repo: make(map[uuid.UUID]models.Document)}
	documentController := &v1.DocumentController{DocumentRepository: repo}
	router := v1.SetupRouter(documentController, nil, nil)

	data := v1.CreateRequest{
		DocumentBase64String: "TEMP DOCUMENT",
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

func getDocumentFromPresentUUID(t *testing.T) {
	repo := &mock.MapDocumentRepository{Repo: make(map[uuid.UUID]models.Document)}
	documentController := &v1.DocumentController{DocumentRepository: repo}
	router := v1.SetupRouter(documentController, nil, nil)

	ExampleUUID := uuid.New()
	ExampleDocument := models.Document{
		Uuid:          ExampleUUID,
		PdfBase64:     func() *string { v := "TEMP DOCUMENT"; return &v }(),
		SelectionData: nil,
	}
	repo.Repo[ExampleUUID] = ExampleDocument
	expectedResponse := fmt.Sprintf(`{"document":{"documentUUID":"%s","pdfBase64":"%s"}}`, ExampleDocument.Uuid, *ExampleDocument.PdfBase64)

	request := &v1.GetDocumentRequest{DocumentUUID: &ExampleUUID}
	requestJSON, _ := json.Marshal(request)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, httptest.NewRequest(
		"GET",
		"/api/v1/documents/?documentUUID="+ExampleUUID.String(),
		strings.NewReader(string(requestJSON)),
	))

	assert.Equal(t, http.StatusOK, w.Code, "Response should be 200")
	assert.Equal(t, expectedResponse, w.Body.String())
}
