package v1

import (
	"github.com/gin-gonic/gin"
	"pdf_service_api/service"
)

type DataServiceIntegrationController struct {
	dataService service.DataService
}

// GenerateMeta handles the HTTP PORT request to return the generated metadata based on the data service's json response.
// It expects a pdf document as a string.
//
// Upon successful retrieval, it returns a 200 OK status with the metadata object.
// a 400 Bad Request or 500 Internal Server Error status with an appropriate error message.
//
// @Summary Get metadata by UUID
// @Description Retrieves metadata associated with a given UUID.
// @Tags meta
// @Accept  json
// @Produce  json
// @Param   base64 query string true "The base64 of the document"
// @Success 200 {object} models.Meta "Successful creation of metadata"
// @Failure 400 "Bad request, typically due to missing param"
// @Failure 500 "Internal server error"
// @Router /generatemeta [post]
func (t DataServiceIntegrationController) GenerateMeta(c *gin.Context) {

}
