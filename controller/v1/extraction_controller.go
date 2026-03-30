package v1

import (
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
// @Summary
// @Description
// @Tags extraction
// @Accept
// @Produce
// @Param
// @Success 200
// @Failure 400 "Bad request, typically due to invalid input"
// @Failure 500 "Internal server error, typically due to database issues"
// @Router
func (t SelectionController) extractAsText(c *gin.Context) {
	body := &ExtractUUIDsRequest{}
	if err := c.ShouldBindJSON(body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for _, element := range body.SelectionUUIDs {
		stringUUID, err := uuid.Parse(element)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		_, err = t.SelectionRepository.GetSelectionBySelectionUUID(stringUUID)
		if err != nil {
			return
		}
	}

}
