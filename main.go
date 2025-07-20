package main

import (
	v1 "pdf_service_api/controller/v1"
	pg "pdf_service_api/postgres"
)

func main() {
	dbHandler := pg.DatabaseHandler{DbConfig: pg.ConfigForDatabase{}}

	documentCtrl := &v1.DocumentController{DocumentRepository: pg.NewDocumentRepository(dbHandler)}
	selectionCtrl := &v1.SelectionController{SelectionRepository: pg.NewSelectionRepository(dbHandler)}
	metaCtrl := &v1.MetaController{MetaRepository: pg.NewMetaRepository(dbHandler)}

	router := v1.SetupRouter(documentCtrl, selectionCtrl, metaCtrl)
	_ = router.Run() // listen and serve on 0.0.0.0:8080
}
