package unit

import (
	"github.com/go-playground/assert/v2"
	"github.com/google/uuid"
	"net/http"
	"net/http/httptest"
	"pdf_service_api/models"
	"pdf_service_api/repositories/mock"
	"pdf_service_api/v1"
	"pdf_service_api/v1/controller"
	"testing"
)

/*
*
This test serves to test the ping functionality of the router.
*/
func TestUploadFileController(t *testing.T) {
	repo := &mock.MapRepository{Repo: make(map[uuid.UUID]models.Document)}
	documentController := &controller.DocumentController{DocumentRepository: repo}
	router := v1.SetupRouter(documentController)

	w := httptest.NewRecorder() //creates a recorder that records its mutations for later inspection in tests.
	req, _ := http.NewRequest("GET", "/ping", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "{\"message\":\"pong\"}", w.Body.String())
}
