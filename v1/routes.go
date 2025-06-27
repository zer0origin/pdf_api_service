package v1

import "github.com/gin-gonic/gin"

func SetupRoutes(router *gin.RouterGroup) {
	// All routes defined here are relative to the group passed in.
	router.GET("/upload", uploadDocument)
}
