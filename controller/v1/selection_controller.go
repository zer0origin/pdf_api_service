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
	if id, present := c.GetQuery("documentUUID"); present {
		uid, err := uuid.Parse(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Error": err})
			return
		}

		results, err := t.SelectionRepository.GetSelectionsByDocumentId(uid)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": err})
			return
		}

		c.JSON(200, results)
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
		return
	}

	if id, present := c.GetQuery("selectionUUID"); present {
		handleDeletion(id, t.SelectionRepository.DeleteSelectionBySelectionUUID)
	}

	if id, present := c.GetQuery("documentUUID"); present {
		handleDeletion(id, t.SelectionRepository.DeleteSelectionByDocumentUUID)
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
		DocumentID:      reqBody.DocumentID,
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
	c.DELETE("/:id", t.DeleteSelectionWhereSelectionUUID)
	c.POST("/", t.AddSelection)
	c.GET("/", t.GetSelection)
}
