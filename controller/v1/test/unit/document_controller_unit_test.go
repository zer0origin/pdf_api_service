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
	t.Run("Upload a document", uploadDocument)
	t.Run("Get document from present document uuid", getDocumentFromPresentDocumentUUID)
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

func getDocumentFromPresentDocumentUUID(t *testing.T) {
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
	expectedResponse := fmt.Sprintf(`{"documents":[{"documentUUID":"%s","pdfBase64":"%s"}]}`, ExampleDocument.Uuid, *ExampleDocument.PdfBase64)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, httptest.NewRequest(
		"GET",
		"/api/v1/documents/?documentUUID="+ExampleUUID.String(),
		nil,
	))

	assert.Equal(t, http.StatusOK, w.Code, "Response should be 200")
	assert.Equal(t, expectedResponse, w.Body.String())
}
