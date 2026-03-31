package v1

import (
	"fmt"
	"net/http"
	"pdf_service_api/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ExtractionController struct {
	SelectionRepository models.SelectionRepository
}

// extraction handles the HTTP POST request to take selections and retrieve the text.
//
// @Summary Get text inside selection UUIDs
// @Description  Provide a bunch of selection UUIDs and receive the text within them.
// @Tags extract
// @Accept json
// @Produce json
// @Success 200 {object} map[string][]string "Successful retrieval of selections"
// @Success 200
// @Failure 400 "Bad request, typically due to invalid input"
// @Failure 500 "Internal server error, typically due to database issues"
// @Router /extract/basic [post]
func (t ExtractionController) extractAsText(c *gin.Context) {
	body := &ExtractUUIDsRequest{}
	if err := c.ShouldBindBodyWithJSON(body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for _, element := range *body {
		stringUUID, err := uuid.Parse(element)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		res, err := t.SelectionRepository.GetSelectionBySelectionUUID(stringUUID)
		if err != nil {
			return
		}

		fmt.Println(res)
	}
}

func (t ExtractionController) SetupRouter(c *gin.RouterGroup) {
	c.POST("/basic", t.extractAsText)
}
