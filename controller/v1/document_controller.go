package v1

import (
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"pdf_service_api/models"
	"slices"
	"strconv"
)

// DocumentController injects the dependencies required for the controller implementations to operate.
type DocumentController struct {
	DocumentRepository models.DocumentRepository
}

// GetDocumentHandler
// @Summary Get documents
// @Description Retrieves document details. Documents can be fetched either by their unique Document UUID or by an Owner UUID.
// @Description Optional exclusion parameters can be used to omit specific fields from the response.
// @Tags documents
// @Accept json
// @Produce json
// @Param documentUUID query string false "The unique identifier of the document to retrieve. If provided"
// @Param ownerUUID query string true "The unique identifier of the owner whose documents are to be retrieved."
// @Param exclude query []string false "Fields to exclude from the response. Allowed values: `documentTitle`, `timeCreated`, `ownerUUID`, `ownerType`, `base64`." collectionFormat(multi)
// @Success 200 {object} object{documents=[]models.Document} "Successfully retrieved document(s)."
// @Failure 400 {object} object{error=string} "Bad Request: Invalid UUID format or no valid parameters specified."
// @Failure 404 {object} object{error=string} "Not Found: No document(s) found for the given UUID."
// @Failure 500 {object} object{error=string} "Internal Server Error: An unexpected error occurred on the server."
// @Router /documents [get]
func (t DocumentController) GetDocumentHandler(c *gin.Context) {
	exclude := make(map[string]bool)
	if values, present := c.GetQueryArray("exclude"); present {
		if slices.Contains(values, "documentTitle") {
			exclude["documentTitle"] = true
		}

		if slices.Contains(values, "timeCreated") {
			exclude["timeCreated"] = true
		}

		if slices.Contains(values, "ownerUUID") {
			exclude["ownerUUID"] = true
		}

		if slices.Contains(values, "ownerType") {
			exclude["ownerType"] = true
		}

		if slices.Contains(values, "base64") {
			exclude["pdfBase64"] = true
		}
	}

	var limit int8 = 100
	if values, present := c.GetQuery("limit"); present {
		number, err := strconv.ParseInt(values, 10, 8)
		if err != nil {
			return
		}

		limit = int8(number)
	}

	var offset int8 = 0
	if values, present := c.GetQuery("offset"); present {
		number, err := strconv.ParseInt(values, 10, 8)
		if err != nil {
			return
		}

		offset = int8(number)
	}

	documentUidStr, isDocumentUuidPresent := c.GetQuery("documentUUID")
	ownerUidStr, isOwnerUuidPresent := c.GetQuery("ownerUUID")

	if !isOwnerUuidPresent {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Required OwnerUuid is missing"})
		return
	}

	ownerUid, err := uuid.Parse(ownerUidStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if isDocumentUuidPresent {
		documentUid, err := uuid.Parse(documentUidStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		document, err := t.DocumentRepository.GetDocumentByDocumentUUID(documentUid, ownerUid, exclude)
		if err != nil {
			switch {
			case errors.Is(err, sql.ErrNoRows):
				c.JSON(http.StatusNotFound, gin.H{"error": "Document with documentUUID " + documentUid.String() + " was not found."})
				return
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}

		c.JSON(200, gin.H{"documents": []models.Document{document}})
		return
	}

	documents, err := t.DocumentRepository.GetDocumentByOwnerUUID(ownerUid, limit, offset, exclude)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			c.JSON(http.StatusNotFound, gin.H{"error": "Document with ownerUUID " + ownerUid.String() + " was not found."})
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(200, gin.H{"documents": documents})
	return
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newModel := models.Document{
		Uuid:          uuid.New(),
		Base64:        &body.Base64,
		DocumentTitle: body.DocumentTitle,
		OwnerUUID:     body.OwnerUUID,
		OwnerType:     body.OwnerType,
		SelectionData: nil,
	}

	err = t.DocumentRepository.UploadDocument(newModel)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
// @Param   ownerUUID query string true "The UUID of the owner of the document that is getting deleted"
// @Success 200 {object} map[string]bool "Successful deletion"
// @Failure 400 "Bad request, typically due to missing/invalid UUID or deletion failure"
// @Router /documents [delete]
func (t DocumentController) DeleteDocumentHandler(c *gin.Context) {
	ownerUuidStr, isPresent := c.GetQuery("ownerUUID")
	if !isPresent {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ownerUUID was not present"})
		return
	}

	documentUuidStr, isPresent := c.GetQuery("documentUUID")
	if !isPresent {
		c.JSON(http.StatusBadRequest, gin.H{"error": "documentUUID was not present"})
		return
	}

	ownerUuid, err := uuid.Parse(ownerUuidStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	documentUuid, err := uuid.Parse(documentUuidStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = t.DocumentRepository.DeleteDocumentById(documentUuid, ownerUuid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"success": true})
	return
}

func (t DocumentController) SetupRouter(c *gin.RouterGroup) {
	c.POST("/", t.UploadDocumentHandler)
	c.PUT("/", t.UploadDocumentHandler)
	c.GET("/", t.GetDocumentHandler)
	c.DELETE("/", t.DeleteDocumentHandler)
}
