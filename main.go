package main

import (
	"pdf_service_api/database"
	"pdf_service_api/repositories"
	v1 "pdf_service_api/v1"
	"pdf_service_api/v1/controller"
)

func main() {
	documentController := createDocumentController()
	router := v1.SetupRouter(documentController)
	_ = router.Run() // listen and serve on 0.0.0.0:8080
}

func createDocumentController() *controller.DocumentController {
	dbConfig := database.ConfigForDatabase{}
	repository := repositories.NewSelectionRepository(dbConfig)
	selController := &controller.SelectionController{SelectionRepository: repository}

	documentController := &controller.DocumentController{
		DocumentRepository:  repositories.NewDocumentRepository(dbConfig),
		SelectionController: selController,
	}

	return documentController
}
