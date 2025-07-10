package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"pdf_service_api/models"
	"pdf_service_api/models/requests"
	"pdf_service_api/repositories"
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
	body := &requests.UploadRequest{}

	err := c.ShouldBindJSON(body)
	if err != nil {
		return
	}

	newModel := models.Document{
		Uuid:          uuid.New(),
		PdfBase64:     body.DocumentBase64String,
		SelectionData: nil,
	}

	err = t.DocumentRepository.UploadDocument(newModel)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(200, gin.H{"documentUUID": newModel.Uuid})
}

func (t DocumentController) SetupRouter(c *gin.RouterGroup) {
	c.POST("/", t.UploadDocumentHandler)
	c.PUT("/", t.UploadDocumentHandler)
	c.GET("/:id", t.GetDocumentHandler)
}
