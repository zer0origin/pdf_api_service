package main

import (
	"fmt"
	"log"
	"os"
	v1 "pdf_service_api/controller/v1"
	pg "pdf_service_api/postgres"
)

var (
	dbUser     = os.Getenv("DATABASE_USER")
	dbPassword = os.Getenv("DATABASE_PASSWORD")
	dbPort     = os.Getenv("DATABASE_PORT")
	dbHost     = os.Getenv("DATABASE_HOST")
	dbDatabase = os.Getenv("DATABASE_DB")
)

func main() {
	errHandleFunction := func(str string) {
		panic("Database login credentials must be present.")
	}
	mustNotBeEmpty(errHandleFunction, dbUser, dbPassword, dbPort, dbHost, dbDatabase)

	dbHandler := pg.DatabaseHandler{DbConfig: pg.ConfigForDatabase{
		Host:     dbHost,
		Port:     dbPort,
		Username: dbUser,
		Password: dbPassword,
		Database: dbDatabase,
	}}

	err := dbHandler.RunInitScript()
	if err != nil {
		err = fmt.Errorf("failed to run init script: %s", err)
		panic(err)
	}

	documentCtrl := &v1.DocumentController{DocumentRepository: pg.NewDocumentRepository(dbHandler)}
	selectionCtrl := &v1.SelectionController{SelectionRepository: pg.NewSelectionRepository(dbHandler)}
	metaCtrl := &v1.MetaController{MetaRepository: pg.NewMetaRepository(dbHandler)}

	router := v1.SetupRouter(documentCtrl, selectionCtrl, metaCtrl)
	log.Fatal(router.Run(":8080"))
}

func mustNotBeEmpty(errorHandle func(string), a ...string) {
	for _, s := range a {
		if len(s) == 0 {
			errorHandle(s)
		}
	}
}
