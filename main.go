package main

import (
	"pdf_service_api/database"
	v1 "pdf_service_api/v1"
	"pdf_service_api/v1/controller"
	"pdf_service_api/v1/repositories"
)

func main() {
	dbConfig := database.ConfigForDatabase{}
	documentController := controller.NewDocumentController(repositories.NewDocumentRepository(dbConfig))
	router := v1.SetupRouter(documentController)
	_ = router.Run() // listen and serve on 0.0.0.0:8080
}
