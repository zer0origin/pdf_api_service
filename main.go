package main

import (
	"log"
	"os"
	v1 "pdf_service_api/controller/v1"
	pg "pdf_service_api/postgres"
)

var (
	dbUser     string = os.Getenv("DATABASE_USER")
	dbPassword string = os.Getenv("DATABASE_PASSWORD")
	dbPort     string = os.Getenv("DATABASE_PORT")
	dbHost     string = os.Getenv("DATABASE_HOST")
	dbDatabase string = os.Getenv("DATABASE_DB")
)

func main() {
	dbHandler := pg.DatabaseHandler{DbConfig: pg.ConfigForDatabase{
		Host:     dbHost,
		Port:     dbPort,
		Username: dbUser,
		Password: dbPassword,
		Database: dbDatabase,
	}}

	documentCtrl := &v1.DocumentController{DocumentRepository: pg.NewDocumentRepository(dbHandler)}
	selectionCtrl := &v1.SelectionController{SelectionRepository: pg.NewSelectionRepository(dbHandler)}
	metaCtrl := &v1.MetaController{MetaRepository: pg.NewMetaRepository(dbHandler)}

	router := v1.SetupRouter(documentCtrl, selectionCtrl, metaCtrl)
	log.Fatal(router.Run(":8080"))
}
