package main

import (
	v1 "pdf_service_api/v1"
)

func main() {
	documentController := controller.NewDocumentController(repositories.NewDocumentRepository())
	router := v1.SetupRouter(*documentController)
	_ = router.Run() // listen and serve on 0.0.0.0:8080
}
