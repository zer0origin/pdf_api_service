package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"pdf_service_api/repositories"
)

type SelectionController struct {
	SelectionRepository repositories.SelectionRepository
}

func (t SelectionController) getSelectionFromId(c *gin.Context) {
	uid, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err})
		return
	}

	results, err := t.SelectionRepository.GetSelectionBySelectionId(uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": err})
		return
	}

	c.JSON(200, results)
}

func (t SelectionController) deleteSelectionWhereSelectionUUID(c *gin.Context) {
	uid, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err})
		return
	}

	err = t.SelectionRepository.DeleteSelectionBySelectionUUIDFunction(uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": err})
		return
	}
}

func (t SelectionController) SetupRouterAppendToDocument(c *gin.RouterGroup) {
	c.GET("/", t.getSelectionFromId)
}

func (t SelectionController) SetupRouter(c *gin.RouterGroup) {
	c.DELETE("/:id", t.deleteSelectionWhereSelectionUUID)
}
