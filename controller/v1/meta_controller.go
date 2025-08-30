package v1

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"pdf_service_api/models"
	"pdf_service_api/service/dataapi"
)

type MetaController struct {
	DocumentRepository models.DocumentRepository
	MetaRepository     models.MetaRepository
	DataService        dataapi.DataService
}

// AddMeta handles the HTTP POST request to add new metadata.
// It expects a JSON request body conforming to the AddMetaRequest struct,
// which should contain the NumberOfPages, Height, Width, and Images for the new metadata.
//
// A new UUID will be generated for the metadata.
// Upon successful creation, it returns a 200 OK status with the UUID of the
// newly created metadata. If there's an error during request binding or
// metadata creation, it returns a 400 Bad Request or 500 Internal Server Error
// status with an error message.
//
// @Summary Add new metadata
// @Description Creates new metadata with a generated UUID.
// @Tags meta
// @Accept  json
// @Produce  json
// @Param   request body v1.AddMetaRequest true "Metadata creation request"
// @Success 200 {object} map[string]uuid.UUID "Successful creation, returns the metadata UUID"
// @Failure 400 "Bad request, typically due to invalid input"
// @Failure 500 "Internal server error, typically due to database issues"
// @Router /meta [post]
func (t MetaController) AddMeta(c *gin.Context) {
	body := &AddMetaRequest{}
	if err := c.ShouldBindJSON(body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if body.DocumentBase64String == nil {
		exclude := make(models.Exclude)
		exclude.TimeCreated(true).OwnerUUID(true).OwnerType(true).DocumentTitle(true)
		document, err := t.DocumentRepository.GetDocumentByDocumentUUID(body.DocumentUUID, body.OwnerUUID, exclude)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		body.DocumentBase64String = document.PdfBase64
	}

	request, err := t.DataService.SendMetaRequest(*body.DocumentBase64String)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Errorf("error sending SendMetaRequest: %v", err.Error())})
		return
	}

	request.DocumentUUID = body.DocumentUUID
	request.OwnerUUID = &body.OwnerUUID
	request.OwnerType = &body.OwnerType

	err = t.MetaRepository.AddMeta(request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

// UpdateMeta handles the HTTP PUT request to update existing metadata.
// It expects a JSON request body conforming to the UpdateMetaRequest struct,
// which should contain the UUID of the metadata to be updated, and the fields
// to be modified (NumberOfPages, Height, Width, Images). Note that these fields
// are pointers in the `models.Meta` struct, allowing for partial updates.
//
// Upon successful update, it returns a 200 OK status with an empty JSON object.
// If there's an error during request binding or metadata update, it returns
// a 400 Bad Request or 500 Internal Server Error status with an error message.
//
// @Summary Update existing metadata
// @Description Updates specific fields of an existing metadata entry.
// @Tags meta
// @Accept  json
// @Produce  json
// @Param   request body UpdateMetaRequest true "Metadata update request"
// @Success 200 "Successful update"
// @Failure 400 "Bad request, typically due to invalid input"
// @Failure 500 "Internal server error, typically due to database issues"
// @Router /meta [put]
func (t MetaController) UpdateMeta(c *gin.Context) {
	if id, isPresent := c.GetQuery("documentUUID"); isPresent {
		uid, err := uuid.Parse(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		body := &UpdateMetaRequest{}
		if err := c.ShouldBindJSON(body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		model := models.Meta{
			DocumentUUID:  body.UUID,
			NumberOfPages: body.NumberOfPages,
			Height:        body.Height,
			Width:         body.Width,
			Images:        body.Images,
		}

		if err := t.MetaRepository.UpdateMeta(uid, model); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{})

		return
	}

	c.JSON(http.StatusBadRequest, gin.H{"Error": "No param specified."})
}

// DeleteMeta handles the HTTP DELETE request to remove metadata.
// It expects a JSON request body conforming to the DeleteMetaRequest struct,
// which should contain the UUID of the metadata to be deleted.
//
// Upon successful deletion, it returns a 200 OK status with an empty JSON object.
// If there's an error during request binding or metadata deletion, it returns
// a 400 Bad Request or 500 Internal Server Error status with an error message.
//
// @Summary Delete metadata by UUID
// @Description Deletes metadata based on the provided UUID in the request body.
// @Tags meta
// @Accept  json
// @Produce  json
// @Param   request body v1.DeleteMetaRequest true "Metadata deletion request"
// @Success 200 {object} map[string]bool "Successful deletion"
// @Failure 400 "Bad request, typically due to invalid input"
// @Failure 500 "Internal server error, typically due to database issues"
// @Router /meta [delete]
func (t MetaController) DeleteMeta(c *gin.Context) {
	body := &DeleteMetaRequest{}
	if err := c.ShouldBindJSON(body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	model := models.Meta{
		DocumentUUID: body.UUID,
	}

	if err := t.MetaRepository.DeleteMeta(model); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

// GetMeta handles the HTTP GET request to retrieve metadata by its UUID.
// It expects the metadata's UUID as a query parameter named "id".
//
// Upon successful retrieval, it returns a 200 OK status with the metadata object.
// If the UUID is missing, invalid, or if an error occurs during retrieval, it returns
// a 400 Bad Request or 500 Internal Server Error status with an appropriate error message.
//
// @Summary Get metadata by UUID
// @Description Retrieves metadata associated with a given UUID.
// @Tags meta
// @Accept  json
// @Produce  json
// @Param   documentUUID query string true "The UUID of the metadata to retrieve"
// @Param   ownerUUID query string true "The UUID of the owner of the metadata to retrieve"
// @Success 200 {object} models.Meta "Successful retrieval of metadata"
// @Success 404 data not found
// @Failure 400 "Bad request, typically due to missing/invalid UUID"
// @Failure 500 "Internal server error, typically due to database issues"
// @Router /meta [get]
func (t MetaController) GetMeta(c *gin.Context) {
	documentUid, isPresent := c.GetQuery("documentUUID")
	if !isPresent {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Required param documentUUID missing!"})
		return
	}

	ownerUid, isPresent := c.GetQuery("ownerUUID")
	if !isPresent {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Required param ownerUUID missing!"})
		return
	}

	documentUUID, err := uuid.Parse(documentUid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ownerUUID, err := uuid.Parse(ownerUid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	data, err := t.MetaRepository.GetMeta(documentUUID, ownerUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "data not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, data)
	return
}

func (t MetaController) SetupRouter(c *gin.RouterGroup) {
	c.GET("/", t.GetMeta)
	c.POST("/", t.AddMeta)
	c.PUT("/", t.UpdateMeta)
	c.DELETE("/", t.DeleteMeta)
}
