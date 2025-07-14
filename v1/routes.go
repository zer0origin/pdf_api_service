package v1

import (
	"github.com/gin-gonic/gin"
	"pdf_service_api/v1/controller"
)

func SetupRouter(documentController *controller.DocumentController) *gin.Engine {
	router := gin.Default()
	router.GET("/ping", controller.OnPing)

	apiV1Group := router.Group("/api/v1/")
	documentGroup := apiV1Group.Group("/documents")
	documentController.SetupRouter(documentGroup)

	if documentController.SelectionController != nil {
		selectionGroup := apiV1Group.Group("/selections")
		documentController.SelectionController.SetupRouter(selectionGroup)
	}

	return router
}
