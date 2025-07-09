package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"pdf_service_api/v1/models"
	"pdf_service_api/v1/models/requests"
	"pdf_service_api/v1/repositories"
)

// DocumentController injects the dependencies required for the controller implementations to operate.
type DocumentController struct {
	DocumentRepository repositories.DocumentRepository
}

// NewDocumentController creates a new instance of the repository using the injected repositories.DocumentRepository
func NewDocumentController(repository repositories.DocumentRepository) *DocumentController {
	return &DocumentController{DocumentRepository: repository}
}

// GetDocumentHandler gin handler function.
func (t DocumentController) GetDocumentHandler(c *gin.Context) {
	body := &requests.GetDocumentRequest{}

	if err := c.ShouldBindJSON(body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	document, err := t.DocumentRepository.GetDocumentById(body.DocumentUuid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(200, document)
}

// UploadDocumentHandler gin handler function
func (t DocumentController) UploadDocumentHandler(c *gin.Context) {
	body := &models.Document{}

	if err := c.ShouldBindJSON(body); err != nil { //Assume all requests bodies are json.
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	if body.Uuid == uuid.Nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "uuid is required"})
		return
	}

	if *body.PdfBase64 == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "pdfBase64 is required"})
		return
	}

	documentUUID, err := t.DocumentRepository.UploadDocument(*body)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(200, gin.H{"documentUUID": documentUUID})
}

func (t DocumentController) SetupRouter(c *gin.RouterGroup) {
	c.POST("/", t.UploadDocumentHandler)
	c.PUT("/", t.UploadDocumentHandler)
	c.GET("/:id", t.GetDocumentHandler)
}
