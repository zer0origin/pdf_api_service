package v1

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	database "pdf_service_api/database"
)

type Base64DocumentString struct {
	DocumentString string `json:"pdfBase64"`
}

func uploadDocument(c *gin.Context) {
	body := &Base64DocumentString{}
	if err := c.ShouldBindJSON(&body); err != nil {
		fmt.Printf("Document received!\n")

		err := database.WithConnection(uploadDocumentSQL)
		if err != nil {
			c.JSON(200, gin.H{"Error": err})
		}
	} else {
		c.JSON(http.StatusInternalServerError, err)
	}
}

func uploadDocumentSQL(db *sql.DB) {
}
