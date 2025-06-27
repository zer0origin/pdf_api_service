package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Base64DocumentString struct {
	DocumentString string `json:"pdfBase64"`
}

func uploadDocument(c *gin.Context) {
	body := &Base64DocumentString{}
	if err := c.ShouldBindJSON(&body); err != nil {
		fmt.Printf("Document received!\n")
	} else {
		c.JSON(http.StatusInternalServerError, err)
	}
}
