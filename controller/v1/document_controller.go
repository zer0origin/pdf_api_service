package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"pdf_service_api/domain"
)

// DocumentController injects the dependencies required for the controller implementations to operate.
type DocumentController struct {
	DocumentRepository  domain.DocumentRepository
	SelectionController *SelectionController
}

// GetDocumentHandler gin handler function.
func (t DocumentController) GetDocumentHandler(c *gin.Context) {
	getUUID := uuid.MustParse(c.Param("id"))

	document, err := t.DocumentRepository.GetDocumentById(getUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(200, document)
}

// UploadDocumentHandler gin handler function
func (t DocumentController) UploadDocumentHandler(c *gin.Context) {
	body := &UploadRequest{}

	err := c.ShouldBindJSON(body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	newModel := domain.Document{
		Uuid:          uuid.New(),
		PdfBase64:     body.DocumentBase64String,
		SelectionData: nil,
	}

	err = t.DocumentRepository.UploadDocument(newModel)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	c.JSON(200, gin.H{"documentUUID": newModel.Uuid})
}

func (t DocumentController) DeleteDocumentHandler(c *gin.Context) {
	deleteUUID := uuid.Nil

	id := c.Param("id")
	deleteUUID = uuid.MustParse(id)

	err := t.DocumentRepository.DeleteDocumentById(deleteUUID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	c.JSON(200, gin.H{"success": true})
}

func (t DocumentController) SetupRouter(c *gin.RouterGroup) {
	c.POST("/", t.UploadDocumentHandler)
	c.PUT("/", t.UploadDocumentHandler)
	c.GET("/:id", t.GetDocumentHandler)
	c.DELETE("/:id", t.DeleteDocumentHandler)

	if t.SelectionController == nil {
		return
	}

	selectionGroup := c.Group("/:id/selections")
	selectionGroup.GET("/", t.SelectionController.getSelectionFromId)
}
