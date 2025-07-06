package v1

import (
	"github.com/gin-gonic/gin"
	"pdf_service_api/v1/controller"
	"pdf_service_api/v1/controller/document"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()
	router.GET("/ping", controller.OnPing)

	apiV1Group := router.Group("/api/v1/")
	documentGroup := apiV1Group.Group("/documents")
	document.SetupRouter(documentGroup)

	return router
}
