package document

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"pdf_service_api/database"
)

type Base64DocumentString struct { //TODO: Move!
	Data string `json:"pdfBase64"`
}

func uploadDocument(c *gin.Context) {
	body := &Base64DocumentString{}
	if err := c.ShouldBindJSON(&body); err == nil {
		fmt.Printf("Body Parsed!\n")
		fmt.Println(body.Data)

		u := uuid.Nil
		uploadDocumentSQL := createUploadDocumentFunction(body, &u) //create callback
		err := database.WithConnection(uploadDocumentSQL)
		if err != nil {
			c.JSON(200, gin.H{"Error": err})
		}

		if u != uuid.Nil {
			c.JSON(200, gin.H{"UUID": u.String()})
		} else {
			c.JSON(200, gin.H{"Error": "An error occurred while retrieving the inserted object's UUID"})
		}
	} else {
		c.JSON(http.StatusInternalServerError, err)
	}
}

func createUploadDocumentFunction(documentString *Base64DocumentString, insertedUUID *uuid.UUID) func(db *sql.DB) {
	return func(db *sql.DB) {
		sqlStatement := `insert into document_table values ($1, $2) returning "Document_UUID"`

		u := uuid.New()
		_, err := db.Exec(sqlStatement, u, documentString.Data)

		if err != nil {
			panic(err)
		}

		*insertedUUID = u
	}
}
