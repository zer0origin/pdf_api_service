package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"pdf_service_api/domain"
)

type SelectionController struct {
	SelectionRepository domain.SelectionRepository
}

func (t SelectionController) GetSelection(c *gin.Context) {
	getSelection := func(id string, passedServiceGetFunction func(uid uuid.UUID) ([]domain.Selection, error)) {
		uid, err := uuid.Parse(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Error": err})
			return
		}

		results, err := passedServiceGetFunction(uid)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": err})
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

func (t SelectionController) DeleteSelection(c *gin.Context) {
	handleDeletion := func(id string, serviceFunction func(uid uuid.UUID) error) {
		uid, err := uuid.Parse(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Error": err})
			return
		}

		err = serviceFunction(uid)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": err})
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

func (t SelectionController) AddSelection(c *gin.Context) {
	reqBody := &AddNewSelectionRequest{}

	if err := c.ShouldBindJSON(reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err})
		return
	}

	toCreate := domain.Selection{
		Uuid:            uuid.New(),
		DocumentUUID:    reqBody.DocumentID,
		IsComplete:      reqBody.IsComplete,
		Settings:        reqBody.Settings,
		SelectionBounds: reqBody.SelectionBounds,
	}

	err := t.SelectionRepository.AddNewSelection(toCreate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": err})
		return
	}

	c.JSON(200, gin.H{"selectionUUID": toCreate.Uuid.String()})
}

func (t SelectionController) SetupRouter(c *gin.RouterGroup) {
	c.DELETE("/", t.DeleteSelection)
	c.POST("/", t.AddSelection)
	c.GET("/", t.GetSelection)
}
