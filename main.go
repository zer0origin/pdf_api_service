package main

import (
	v2 "pdf_service_api/controller/v1"
	pg "pdf_service_api/postgres"
)

func main() {
	documentController := createDocumentController()
	router := v2.SetupRouter(documentController)
	_ = router.Run() // listen and serve on 0.0.0.0:8080
}

func createDocumentController() *v2.DocumentController {
	handler := pg.DatabaseHandler{DbConfig: pg.ConfigForDatabase{}}
	repository := pg.NewSelectionRepository(handler)
	selController := &v2.SelectionController{SelectionRepository: repository}

	documentController := &v2.DocumentController{
		DocumentRepository:  pg.NewDocumentRepository(handler),
		SelectionController: selController,
	}

	return documentController
}
