package main

import (
	"github.com/gin-gonic/gin"
	v1 "pdf_service_api/v1"
)

func main() {
	router := gin.Default()
	apiV1Group := router.Group("/api/v1/")
	v1.SetupRoutes(apiV1Group)

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	_ = router.Run() // listen and serve on 0.0.0.0:8080
}
