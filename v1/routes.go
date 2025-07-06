package v1

import "github.com/gin-gonic/gin"

func SetupRouter() *gin.Engine {
	router := gin.Default()
	router.GET("/ping", onPing)
	apiV1Group := router.Group("/api/v1/")
	apiV1Group.GET("/upload", uploadDocument)

	return router
}
