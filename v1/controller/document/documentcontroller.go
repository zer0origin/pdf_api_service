package document

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"pdf_service_api/v1/models/requests"
)

func getDocument(c *gin.Context) {
	body := &requests.GetDocumentRequest{}

	if err := c.ShouldBindJSON(body); err == nil {
		fmt.Printf("Body Parsed!\n")
		c.JSON(200, body)
	}
}

func SetupRouter(c *gin.RouterGroup) {
	c.POST("/", nil)
	c.PUT("/", uploadDocument)
	c.GET("/:id", getDocument)
}
