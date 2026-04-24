package v1

import (
	"net/http"
	"pdf_service_api/models"
	"pdf_service_api/service/dataapi"

	"github.com/gin-gonic/gin"
)

type ExtractionController struct {
	SelectionRepository models.SelectionRepository
	DocumentRepository  models.DocumentRepository
	DataService         dataapi.DataService
	Options             ExtractionOptions
}

// ExtractionOptions provides optional parameters that change the behaviour of the controller
type ExtractionOptions struct {
	GetBase64IfNotIncluded bool
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
	body := &ExtractUUIDsRequest{} //TODO: Change this to reflect the new json schema, also allow tests to be disabled.
	if err := c.ShouldBindBodyWithJSON(body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if body.Base64EncodedDocument == "" {
		if t.Options.GetBase64IfNotIncluded {
			excludes := models.DocumentExcludes{}.DocumentTitle(true).OwnerType(true).OwnerUUID(true).TimeCreated(true)

			document, err := t.DocumentRepository.GetDocumentByDocumentUUID(body.DocumentUid, body.OwnerUid, excludes)
			if err != nil {
				c.Status(http.StatusBadRequest)
				return
			}

			body.Base64EncodedDocument = *document.PdfBase64
		} else {
			c.Status(http.StatusBadRequest)
			return
		}
	}

	res, err := t.SelectionRepository.GetMapOfSelectionsBySelectionUUID(body.Uids)
	if err != nil {
		return
	}

	req := dataapi.ExtractionRequest{
		DocumentUid:           body.DocumentUid,
		Base64EncodedDocument: body.Base64EncodedDocument,
		Selections:            res,
	}

	err = t.DataService.SendBasicExtractionRequest(req)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	c.Status(200)
}

func (t ExtractionController) SetupRouter(c *gin.RouterGroup) {
	c.POST("/basic", t.extractAsText)
}
