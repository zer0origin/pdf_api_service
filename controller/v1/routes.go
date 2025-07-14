package v1

import (
	"github.com/gin-gonic/gin"
)

func SetupRouter(documentController *DocumentController) *gin.Engine {
	router := gin.Default()
	router.GET("/ping", OnPing)

	apiV1Group := router.Group("/api/v1/")
	documentGroup := apiV1Group.Group("/documents")
	documentController.SetupRouter(documentGroup)

	if documentController.SelectionController != nil {
		selectionGroup := apiV1Group.Group("/selections")
		documentController.SelectionController.SetupRouter(selectionGroup)
	}

	return router
}
