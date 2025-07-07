package repositories

import (
	"database/sql"
	"github.com/google/uuid"
	"pdf_service_api/database"
	"pdf_service_api/v1/models"
)

type DocumentRepository interface {
	UploadDocument(document models.Document) (uuid.UUID, error)
}

type documentRepository struct {
}

func NewDocumentRepository() DocumentRepository {
	return documentRepository{}
}

func (d documentRepository) UploadDocument(document models.Document) (uuid.UUID, error) {
	u := uuid.New()
	uploadDocumentSQL := createUploadDocumentFunction(&document, u) //create callback
	err := database.WithConnection(uploadDocumentSQL)
	if err != nil {
		return uuid.Nil, err
	}

	return u, nil
}

// SQL Query to database
func createUploadDocumentFunction(document *models.Document, insertedUUID uuid.UUID) func(db *sql.DB) error {
	return func(db *sql.DB) error {
		sqlStatement := `insert into document_table values ($1, $2) returning "Document_UUID"`
		_, err := db.Exec(sqlStatement, insertedUUID, document.PdfBase64, insertedUUID)

		if err != nil {
			return err
		}

		return nil
	}
}
