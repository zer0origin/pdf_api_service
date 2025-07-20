package v1

import (
	"github.com/gin-gonic/gin"
)

func SetupRouter(documentController *DocumentController, selectionController *SelectionController, metaController *MetaController) *gin.Engine {
	router := gin.Default()
	router.GET("/ping", OnPing)
	apiV1Group := router.Group("/api/v1/")

	if documentController != nil {
		documentGroup := apiV1Group.Group("/documents")
		documentController.SetupRouter(documentGroup)
	}

	if selectionController != nil {
		selectionGroup := apiV1Group.Group("/selections")
		selectionController.SetupRouter(selectionGroup)
	}

	if metaController != nil {
		metaGroup := apiV1Group.Group("/meta")
		metaController.SetupRouter(metaGroup)
	}

	return router
}
