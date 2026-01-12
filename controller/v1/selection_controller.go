package v1

import (
	"net/http"
	"pdf_service_api/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SelectionController struct {
	SelectionRepository models.SelectionRepository
}

// GetSelection handles the HTTP GET request to retrieve selections based on either
// a document UUID or a selection UUID.
//
// It expects either "documentUUID" or "selectionUUID" as a query parameter.
// If "documentUUID" is provided, it fetches all selections associated with that document.
// If "selectionUUID" is provided, it fetches selections matching that specific selection UUID.
//
// Upon successful retrieval, it returns a 200 OK status with a JSON array of selections.
// If no parameter is specified, the UUID is invalid, or an error occurs during retrieval,
// it returns a 400 Bad Request or 500 Internal Server Error status with an error message.
//
// @Summary Get selections by document or selection UUID
// @Description Retrieves selections based on either a document's UUID or a specific selection's UUID.
// @Tags selections
// @Accept  json
// @Produce  json
// @Param   documentUUID query string false "The UUID of the document to retrieve selections for"
// @Param   selectionUUID query string false "The UUID of the specific selection to retrieve"
// @Success 200 {object} map[string][]models.Selection "Successful retrieval of selections"
// @Failure 400 "Bad request, typically due to missing/invalid UUID parameter"
// @Failure 500 "Internal server error, typically due to database issues"
// @Router /selections [get]
func (t SelectionController) GetSelection(c *gin.Context) {
	getSelection := func(id string, passedServiceGetFunction func(uid uuid.UUID) ([]models.Selection, error)) {
		uid, err := uuid.Parse(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		results, err := passedServiceGetFunction(uid)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{"selections": results})
	}

	if id, isPresent := c.GetQuery("documentUUID"); isPresent {
		getSelection(id, t.SelectionRepository.GetSelectionsByDocumentUUID)
		return
	}

	if id, isPresent := c.GetQuery("selectionUUID"); isPresent {
		getSelection(id, t.SelectionRepository.GetSelectionsBySelectionUUID)
		return
	}

	c.JSON(http.StatusBadRequest, gin.H{"Error": "No param specified."})
}

// DeleteSelection handles the HTTP DELETE request to remove selections.
// It allows deletion by either a specific selection UUID or by a document UUID,
// which will delete all selections associated with that document.
//
// It expects either "selectionUUID" or "documentUUID" as a query parameter.
// If "selectionUUID" is provided, it deletes the specific selection.
// If "documentUUID" is provided, it deletes all selections belonging to that document.
//
// Upon successful deletion, it returns a 200 OK status with a success message.
// If no parameter is specified, the UUID is invalid, or an error occurs during deletion,
// it returns a 400 Bad Request or 500 Internal Server Error status with an error message.
//
// @Summary Delete selections by selection or document UUID
// @Description Deletes selections based on a specific selection UUID or all selections associated with a document UUID.
// @Tags selections
// @Accept  json
// @Produce  json
// @Param   selectionUUID query string false "The UUID of the specific selection to delete"
// @Param   documentUUID query string false "The UUID of the document whose selections are to be deleted"
// @Success 200 {object} map[string]bool "Successful deletion"
// @Failure 400 "Bad request, typically due to missing/invalid UUID parameter"
// @Failure 500 "Internal server error, typically due to database issues"
// @Router /selections [delete]
func (t SelectionController) DeleteSelection(c *gin.Context) {
	handleDeletion := func(id string, serviceFunction func(uid uuid.UUID) error) {
		uid, err := uuid.Parse(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err = serviceFunction(uid)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{"success": true})
	}

	if id, isPresent := c.GetQuery("selectionUUID"); isPresent {
		handleDeletion(id, t.SelectionRepository.DeleteSelectionBySelectionUUID)
		return
	}

	if id, isPresent := c.GetQuery("documentUUID"); isPresent {
		handleDeletion(id, t.SelectionRepository.DeleteSelectionByDocumentUUID)
		return
	}

	c.JSON(http.StatusBadRequest, gin.H{"Error": "No param specified."})
}

// AddSelection handles the HTTP POST request to add a new selection.
// It expects a JSON request body conforming to the AddNewSelectionRequest struct,
// which should include the DocumentUUID, IsComplete status, Settings, and SelectionBounds
// for the new selection.
//
// A new UUID will be generated for the selection.
// Upon successful creation, it returns a 200 OK status with the UUID of the
// newly created selection. If there's an error during request binding or
// selection creation, it returns a 400 Bad Request or 500 Internal Server Error
// status with an error message.
//
// @Summary Add a new selection
// @Description Creates a new selection associated with a document.
// @Tags selections
// @Accept  json
// @Produce  json
// @Param   request body v1.AddNewSelectionRequest true "Selection creation request"
// @Success 200 {object} map[string]uuid.UUID "Successful creation, returns the selection UUID"
// @Failure 400 "Bad request, typically due to invalid input"
// @Failure 500 "Internal server error, typically due to database issues"
// @Router /selections [post]
func (t SelectionController) AddSelection(c *gin.Context) {
	reqBody := &AddNewSelectionRequest{}

	if err := c.ShouldBindJSON(reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	toCreate := models.Selection{
		Uuid:            uuid.New(),
		DocumentUUID:    reqBody.DocumentUUID,
		SelectionBounds: reqBody.SelectionBounds,
		PageKey:         &reqBody.PageKey,
	}

	err := t.SelectionRepository.AddNewSelection(toCreate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"selectionUUID": toCreate.Uuid.String()})
}

func (t SelectionController) SetupRouter(c *gin.RouterGroup) {
	c.DELETE("/", t.DeleteSelection)
	c.POST("/", t.AddSelection)
	c.GET("/", t.GetSelection)
}
