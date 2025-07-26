package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"pdf_service_api/domain"
)

// DocumentController injects the dependencies required for the controller implementations to operate.
type DocumentController struct {
	DocumentRepository domain.DocumentRepository
}

// GetDocumentHandler gin handler function.
func (t DocumentController) GetDocumentHandler(c *gin.Context) {
	if id, isPresent := c.GetQuery("documentUUID"); isPresent {
		uid, err := uuid.Parse(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}

		document, err := t.DocumentRepository.GetDocumentByDocumentUUID(uid)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}

		c.JSON(200, document)
		return
	}

	c.JSON(http.StatusBadRequest, gin.H{"Error": "No param specified."})
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
	if id, isPresent := c.GetQuery("documentUUID"); isPresent {
		uid, err := uuid.Parse(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}

		err = t.DocumentRepository.DeleteDocumentById(uid)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}

		c.JSON(200, gin.H{"success": true})
		return
	}

	c.JSON(http.StatusBadRequest, gin.H{"Error": "No param specified."})
}

func (t DocumentController) SetupRouter(c *gin.RouterGroup) {
	c.POST("/", t.UploadDocumentHandler)
	c.PUT("/", t.UploadDocumentHandler)
	c.GET("/", t.GetDocumentHandler)
	c.DELETE("/", t.DeleteDocumentHandler)
}
