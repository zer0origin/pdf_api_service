package unit

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	v1 "pdf_service_api/controller/v1"
	"pdf_service_api/controller/v1/test/unit/mock"
	"pdf_service_api/models"
	"testing"
)

func TestPingControllerUnit(t *testing.T) {
	t.Run("Ping router", pingRouter)
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
