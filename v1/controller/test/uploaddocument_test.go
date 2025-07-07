package test

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	v1 "pdf_service_api/v1"
	"pdf_service_api/v1/controller"
	"pdf_service_api/v1/controller/test/mock"
	"pdf_service_api/v1/models"
	"strings"
	"testing"
)

func TestUploadDocument(t *testing.T) {
	repo := &mock.MapRepository{Repo: make(map[uuid.UUID]models.Document)}
	documentController := controller.NewDocumentController(repo)
	router := v1.SetupRouter(documentController)

	data := models.Document{Uuid: uuid.New(),
		PdfBase64: func() *string { v := "TEMP DOCUMENT"; return &v }(),
	}
	documentJSON, _ := json.Marshal(data)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, httptest.NewRequest(
		"POST",
		"/api/v1/documents/",
		strings.NewReader(string(documentJSON)),
	))

	assert.Equal(t, http.StatusOK, w.Code, "Response should be 200")
	assert.Equal(t, data.Uuid, repo.Repo[data.Uuid].Uuid, "Data was not saved to repository")
}
