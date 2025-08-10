package v1

import (
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"pdf_service_api/models"
)

// DocumentController injects the dependencies required for the controller implementations to operate.
type DocumentController struct {
	DocumentRepository models.DocumentRepository
}

// GetDocumentHandler godoc
// @Summary Get a document by UUID
// @Description get document details by its UUID
// @Tags documents
// @Accept json
// @Produce json
// @Param documentUUID query string true "Document UUID"
// @Success 200 {object} models.Document // Assuming you have a Document struct defined for your response
// @Failure 400 "Bad Request"
// @Failure 404 "Not Found"
// @Failure 500 "Internal Server Error"
// @Router /documents [get]
func (t DocumentController) GetDocumentHandler(c *gin.Context) {
	if id, isPresent := c.GetQuery("documentUUID"); isPresent {
		uid, err := uuid.Parse(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}

		document, err := t.DocumentRepository.GetDocumentByDocumentUUID(uid)
		if err != nil {
			switch {
			case errors.Is(err, sql.ErrNoRows):
				c.JSON(http.StatusNotFound, gin.H{"error": "Document with UUID " + uid.String() + " was found."})
				return
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"error": err})
				return
			}
		}

		c.JSON(200, gin.H{"document": document})
		return
	}

	c.JSON(http.StatusBadRequest, gin.H{"Error": "No param specified."})
}

// UploadDocumentHandler handles the HTTP POST request to upload a new document.
// It expects a JSON request body conforming to the CreateRequest struct,
// which should contain the document's base64 encoded string.
//
// Upon successful upload, it returns a 200 OK status with the UUID of the
// newly created document. If there's an error during request binding or
// document upload, it returns a 400 Bad Request status with an error message.
//
// @Summary Upload a new document
// @Description Uploads a document by receiving its base64 encoded string in the request body.
// @Tags documents
// @Accept  json
// @Produce  json
// @Param   request body v1.CreateRequest true "Document upload request"
// @Success 200 {object} map[string]string "Successful upload, returns the document UUID"
// @Failure 400 "Bad request, typically due to invalid input or upload failure"
// @Router /documents [post]
func (t DocumentController) UploadDocumentHandler(c *gin.Context) {
	body := &CreateRequest{}

	err := c.ShouldBindJSON(body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	newModel := models.Document{
		Uuid:          uuid.New(),
		PdfBase64:     &body.DocumentBase64String,
		SelectionData: nil,
	}

	err = t.DocumentRepository.UploadDocument(newModel)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	c.JSON(200, gin.H{"documentUUID": newModel.Uuid})
}

// DeleteDocumentHandler handles the HTTP DELETE request to delete a document by its UUID.
// It expects the document's UUID as a query parameter named "documentUUID".
//
// If the UUID is provided and valid, it attempts to delete the document from the repository.
// Upon successful deletion, it returns a 200 OK status with a success message.
// If the UUID is missing, invalid, or if an error occurs during deletion, it returns
// a 400 Bad Request status with an appropriate error message.
//
// @Summary Delete a document
// @Description Deletes a document based on the provided document UUID.
// @Tags documents
// @Accept  json
// @Produce  json
// @Param   documentUUID query string true "The UUID of the document to delete"
// @Success 200 {object} map[string]bool "Successful deletion"
// @Failure 400 "Bad request, typically due to missing/invalid UUID or deletion failure"
// @Router /documents [delete]
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
